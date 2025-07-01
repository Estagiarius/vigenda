package tasks

import (
	"context"
	// "errors" // Not using direct errors from here in new tests, testify/mock handles errors
	"fmt"
	"testing"
	"time"
	"vigenda/internal/models"
	"vigenda/internal/service"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	// "github.com/stretchr/testify/require" // Not using require directly in these tests
)

// MockTaskService is a mock implementation of service.TaskService using testify/mock
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error) {
	args := m.Called(ctx, title, description, classID, dueDate)
	// Handle case where Task might be zero value if error is not nil
	if task, ok := args.Get(0).(models.Task); ok {
		return task, args.Error(1)
	}
	return models.Task{}, args.Error(1)
}

func (m *MockTaskService) ListAllActiveTasks(ctx context.Context) ([]models.Task, error) {
	args := m.Called(ctx)
	if tasks, ok := args.Get(0).([]models.Task); ok {
		return tasks, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) MarkTaskAsCompleted(ctx context.Context, taskID int64) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func (m *MockTaskService) GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error) {
	args := m.Called(ctx, taskID)
	if taskPtr, ok := args.Get(0).(*models.Task); ok {
		return taskPtr, args.Error(1)
	}
	// Handle if Get(0) is nil (e.g. when error is returned and task is nil)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// This case should ideally not be hit if mock is set up correctly
	return nil, fmt.Errorf("mock GetTaskByID returned non-pointer type for task")
}

func (m *MockTaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, taskID int64) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

// Ensure MockTaskService implements the interface
var _ service.TaskService = (*MockTaskService)(nil)


func TestTasksModel_Init(t *testing.T) {
	mockService := new(MockTaskService)
	model := New(mockService)

	mockService.On("ListAllActiveTasks", mock.Anything).Return([]models.Task{}, nil)

	cmd := model.Init()
	assert.NotNil(t, cmd, "Init should return a command")

	msg := cmd() // Execute the command
	_, ok := msg.(tasksLoadedMsg)
	assert.True(t, ok, "Command from Init should produce a tasksLoadedMsg")

	mockService.AssertExpectations(t)
}

func TestTasksModel_CreateTask_SubmitForm(t *testing.T) {
	mockService := new(MockTaskService)
	model := New(mockService)

	// Setup model to be in creating state with some input values
	model.formState = CreatingTask
	model.titleInput.SetValue("New Task Title")
	model.descriptionInput.SetValue("New Task Description")
	// Optional fields like dueDate and classID can be empty or set
	model.focusIndex = len(model.inputs) -1 // Simulate focus on the last input for submission

	expectedTask := models.Task{ID: 1, Title: "New Task Title", Description: "New Task Description", UserID: 1} // Assuming UserID 1 for now
	mockService.On("CreateTask", mock.Anything, "New Task Title", "New Task Description", (*int64)(nil), (*time.Time)(nil)).Return(expectedTask, nil)

	// Simulate pressing Enter to submit the form
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updatedModel.(*Model)

	assert.True(t, m.isLoading, "Model should be loading after submitting create form")
	assert.NotNil(t, cmd, "Update should return a command for task creation")

	// Execute the command returned by the form submission (which is createTaskCmd)
	createMsg := cmd()
	assert.IsType(t, taskCreatedMsg{}, createMsg, "Message should be taskCreatedMsg")

	// Simulate the model receiving the taskCreatedMsg
	// Need to mock ListAllActiveTasks for the refresh call
	mockService.On("ListAllActiveTasks", mock.Anything).Return([]models.Task{expectedTask}, nil)
	updatedModelAfterCreate, refreshCmd := m.Update(createMsg)
	m = updatedModelAfterCreate.(*Model)

	assert.Equal(t, NoForm, m.formState, "Form state should reset to NoForm")
	assert.Nil(t, m.err, "Error should be nil after successful creation")
	assert.NotNil(t, refreshCmd, "Should return a command to refresh tasks")

	refreshMsg := refreshCmd()
	_, ok := refreshMsg.(tasksLoadedMsg)
	assert.True(t, ok, "Refresh command should produce tasksLoadedMsg")

	mockService.AssertExpectations(t)
}

func TestTasksModel_UpdateTask_SubmitForm(t *testing.T) {
	mockService := new(MockTaskService)
	originalTask := &models.Task{ID: 1, Title: "Original Title", Description: "Original Desc", UserID: 1, IsCompleted: false}

	model := New(mockService)
	model.formState = EditingTask
	model.editingTaskID = originalTask.ID
	model.selectedTaskForDetail = originalTask // Store original task to preserve UserID, IsCompleted

	model.titleInput.SetValue("Updated Title")
	model.descriptionInput.SetValue("Updated Desc")
	model.focusIndex = len(model.inputs) - 1

	// Expected task to be sent to UpdateTask service method
	// Note: UserID and IsCompleted should be preserved from originalTask
	taskWithUpdates := &models.Task{
		ID:          originalTask.ID,
		UserID:      originalTask.UserID,
		Title:       "Updated Title",
		Description: "Updated Desc",
		ClassID:     nil, // Assuming not set in form for this test
		DueDate:     nil, // Assuming not set in form for this test
		IsCompleted: originalTask.IsCompleted,
	}

	mockService.On("UpdateTask", mock.Anything, mock.MatchedBy(func(task *models.Task) bool {
		return task.ID == taskWithUpdates.ID &&
			   task.Title == taskWithUpdates.Title &&
			   task.Description == taskWithUpdates.Description &&
			   task.UserID == taskWithUpdates.UserID && // Ensure UserID is preserved
			   task.IsCompleted == taskWithUpdates.IsCompleted // Ensure IsCompleted is preserved
	})).Return(nil) // Successful update

	// Simulate Enter on last field
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updatedModel.(*Model)
	assert.True(t, m.isLoading)
	assert.NotNil(t, cmd)

	updateMsg := cmd() // Execute updateTaskCmd
	assert.IsType(t, taskUpdatedMsg{}, updateMsg)

	// Simulate model receiving taskUpdatedMsg
	mockService.On("ListAllActiveTasks", mock.Anything).Return([]models.Task{*taskWithUpdates}, nil) // For refresh
	updatedModelAfterUpdate, refreshCmd := m.Update(updateMsg)
	m = updatedModelAfterUpdate.(*Model)

	assert.Equal(t, NoForm, m.formState)
	assert.Equal(t, int64(0), m.editingTaskID)
	assert.NotNil(t, refreshCmd)

	mockService.AssertExpectations(t)
}


func TestTasksModel_KeyBindings_InTableView(t *testing.T) {
	mockService := new(MockTaskService)
	initialTasks := []models.Task{{ID: 1, Title: "Task 1", UserID: 1, Description: "Desc 1"}}

	// Initial load for Init()
	mockService.On("ListAllActiveTasks", mock.Anything).Return(initialTasks, nil).Once()
	model := New(mockService)
	setupCmd := model.Init()
	setupMsg := setupCmd()
	model.Update(setupMsg) // Populate table for selection

	model.table.Select(0) // Select the first row

	// Test 'a' -> CreatingTask
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	assert.Equal(t, CreatingTask, m.(*Model).formState)
	model.formState = NoForm // Reset for next test

	// Test 'e' -> EditingTask (via fetching task)
	mockService.On("GetTaskByID", mock.Anything, int64(1)).Return(&initialTasks[0], nil).Once()
	m, cmdEdit := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	assert.True(t, m.(*Model).isLoading)
	fetchEditMsg := cmdEdit()                                  // Execute fetchTaskForDetailCmd
	m, _ = m.(*Model).Update(fetchEditMsg) // Process fetchedTaskDetailMsg
	assert.Equal(t, EditingTask, m.(*Model).formState)
	assert.Equal(t, int64(1), m.(*Model).editingTaskID)
	model.formState = NoForm; model.editingTaskID = 0 // Reset

	// Test 'd' -> ConfirmingDelete
	m, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	assert.True(t, m.(*Model).confirmingDelete)
	assert.Equal(t, int64(1), m.(*Model).taskIDToDelete)
	model.confirmingDelete = false; model.taskIDToDelete = 0 // Reset

	// Test 'c' -> Mark Complete
	mockService.On("MarkTaskAsCompleted", mock.Anything, int64(1)).Return(nil).Once()
	mockService.On("ListAllActiveTasks", mock.Anything).Return([]models.Task{}, nil).Once() // After completion, list is empty
	m, cmdComplete := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	assert.True(t, m.(*Model).isLoading)
	completeMsg := cmdComplete()                                // Execute markTaskCompleteCmd
	m, cmdRefreshComplete := m.(*Model).Update(completeMsg) // Process taskMarkedCompletedMsg
	assert.False(t, m.(*Model).isLoading)
	assert.NotNil(t, cmdRefreshComplete)
	model.formState = NoForm // Reset

	// Test 'v' -> ViewingDetail (via fetching task)
	mockService.On("GetTaskByID", mock.Anything, int64(1)).Return(&initialTasks[0], nil).Once()
	m, cmdView := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	assert.True(t, m.(*Model).isLoading)
	fetchViewMsg := cmdView()                                // Execute fetchTaskForDetailCmd
	m, _ = m.(*Model).Update(fetchViewMsg) // Process fetchedTaskDetailMsg
	assert.Equal(t, ViewingDetail, m.(*Model).formState)
	assert.NotNil(t, m.(*Model).selectedTaskForDetail)

	mockService.AssertExpectations(t)
}


func TestTasksModel_DeleteTask_ConfirmYes(t *testing.T) {
	mockService := new(MockTaskService)
	model := New(mockService)
	model.confirmingDelete = true
	model.taskIDToDelete = 1

	mockService.On("DeleteTask", mock.Anything, int64(1)).Return(nil).Once()
	mockService.On("ListAllActiveTasks", mock.Anything).Return([]models.Task{}, nil).Once() // For refresh

	m, cmdDelete := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	assert.True(t, m.(*Model).isLoading)

	deleteMsg := cmdDelete() // Execute deleteTaskCmd
	m, cmdRefresh := m.(*Model).Update(deleteMsg) // Process taskDeletedMsg

	assert.False(t, m.(*Model).confirmingDelete)
	assert.Equal(t, int64(0), m.(*Model).taskIDToDelete)
	assert.False(t, m.(*Model).isLoading)
	assert.NotNil(t, cmdRefresh)

	mockService.AssertExpectations(t)
}

func TestTasksModel_DeleteTask_ConfirmNo(t *testing.T) {
	mockService := new(MockTaskService)
	model := New(mockService)
	model.confirmingDelete = true
	model.taskIDToDelete = 1

	m, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	assert.False(t, m.(*Model).confirmingDelete)
	assert.Equal(t, int64(0), m.(*Model).taskIDToDelete)
	assert.Nil(t, cmd) // No command for 'n'
	mockService.AssertNotCalled(t, "DeleteTask", mock.Anything, mock.Anything)
}

// Further tests could cover:
// - Error handling for each command (e.g., CreateTask returns error)
// - Form input validation messages (e.g., empty title, invalid date)
// - Navigation within the form (nextInput, prevInput) - harder to test in isolation
// - IsFocused behavior under different states
// - View rendering for different states (though this is complex and snapshot-like)
// - SetSize behavior (how it affects table and form inputs)
// - Handling of tea.WindowSizeMsg
// - More detailed checks of what data is passed to service methods in UpdateTask.
// - Edge cases for table selections (e.g., empty table).
