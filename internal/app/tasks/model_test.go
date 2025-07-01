package tasks

import (
	"context"
	"fmt"
	"testing"
	"time"
	"vigenda/internal/models"
	"vigenda/internal/service"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskService is a mock implementation of service.TaskService using testify/mock
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error) {
	args := m.Called(ctx, title, description, classID, dueDate)
	if task, ok := args.Get(0).(models.Task); ok {
		return task, args.Error(1)
	}
	return models.Task{}, args.Error(1)
}

func (m *MockTaskService) ListAllTasks(ctx context.Context) ([]models.Task, error) {
	args := m.Called(ctx)
	if tasks, ok := args.Get(0).([]models.Task); ok {
		return tasks, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error) {
	args := m.Called(ctx, classID)
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if taskPtr, ok := args.Get(0).(*models.Task); ok {
		return taskPtr, args.Error(1)
	}
	return nil, fmt.Errorf("mock GetTaskByID returned non-pointer type or unexpected type for task")
}

func (m *MockTaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, taskID int64) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

var _ service.TaskService = (*MockTaskService)(nil)

func TestTasksModel_Init(t *testing.T) {
	mockService := new(MockTaskService)
	model := New(mockService)
	mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{}, nil)
	cmd := model.Init()
	assert.NotNil(t, cmd)
	msg := cmd()
	_, ok := msg.(tasksLoadedMsg)
	assert.True(t, ok)
	mockService.AssertExpectations(t)
}

func TestTasksModel_PopulateTables_PendingAndCompleted(t *testing.T) {
	mockService := new(MockTaskService)
	model := New(mockService)
	model.SetSize(80,24)

	task1 := models.Task{ID: 1, Title: "Pending Task", IsCompleted: false}
	task2 := models.Task{ID: 2, Title: "Completed Task", IsCompleted: true}
	allTasks := []models.Task{task1, task2}

	updatedModel, _ := model.Update(tasksLoadedMsg{tasks: allTasks, err: nil})
	m := updatedModel.(*Model)

	assert.Len(t, m.pendingTasksTable.Rows(), 1, "Pending table should have 1 row")
	assert.Equal(t, "Pending Task", m.pendingTasksTable.Rows()[0][1], "Pending task title mismatch")

	assert.Len(t, m.completedTasksTable.Rows(), 1, "Completed table should have 1 row")
	expectedCompletedTitle := lipgloss.NewStyle().Strikethrough(true).Render("Completed Task")
	assert.Equal(t, expectedCompletedTitle, m.completedTasksTable.Rows()[0][1], "Completed task title mismatch or not strikethrough")
}


func TestTasksModel_KeyBindings_InTableView_TabFocusSwitch(t *testing.T) {
	mockService := new(MockTaskService)
	model := New(mockService)
	model.SetSize(80,24)
	mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{}, nil).Once()
	model.Update(model.Init()())

	assert.Equal(t, PendingTableFocus, model.focusedTable)
	assert.True(t, model.pendingTasksTable.Focused())
	assert.False(t, model.completedTasksTable.Focused())

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.Equal(t, CompletedTableFocus, model.focusedTable)
	assert.False(t, model.pendingTasksTable.Focused())
	assert.True(t, model.completedTasksTable.Focused())

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.Equal(t, PendingTableFocus, model.focusedTable)
	assert.True(t, model.pendingTasksTable.Focused())
	assert.False(t, model.completedTasksTable.Focused())
	mockService.AssertExpectations(t)
}


func TestTasksModel_MarkTaskCompleted_MovesToCompletedTable(t *testing.T) {
	mockService := new(MockTaskService)
	pendingTask := models.Task{ID: 1, Title: "Task to complete", IsCompleted: false, UserID: 1}

	mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{pendingTask}, nil).Once()
	model := New(mockService)
	model.SetSize(80,24)
	model.Update(model.Init()())

	assert.Len(t, model.pendingTasksTable.Rows(), 1)
	assert.Len(t, model.completedTasksTable.Rows(), 0)
	model.pendingTasksTable.SetCursor(0)

	mockService.On("MarkTaskAsCompleted", mock.Anything, pendingTask.ID).Return(nil).Once()

	completedTask := pendingTask
	completedTask.IsCompleted = true
	mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{completedTask}, nil).Once()

	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	assert.True(t, model.isLoading)
	assert.NotNil(t, cmd)

	msg := cmd()
	assert.IsType(t, taskMarkedCompletedMsg{}, msg)
	model, cmd = model.Update(msg)

	assert.NotNil(t, cmd)
	msg = cmd()
	assert.IsType(t, tasksLoadedMsg{}, msg)
	model, _ = model.Update(msg)

	assert.False(t, model.isLoading)
	assert.Len(t, model.pendingTasksTable.Rows(), 0, "Pending table should be empty")
	assert.Len(t, model.completedTasksTable.Rows(), 1, "Completed table should have 1 row")
	assert.Equal(t, strikethroughStyle.Render(pendingTask.Title), model.completedTasksTable.Rows()[0][1])

	mockService.AssertExpectations(t)
}


func TestTasksModel_CreateTask_SubmitForm(t *testing.T) {
	mockService := new(MockTaskService)
	model := New(mockService)

	model.currentView = FormView // Set initial state for form
	model.formSubState = CreatingTask
	model.inputs[0].SetValue("New Task Title")
	model.inputs[1].SetValue("New Task Description")
	model.focusIndex = len(model.inputs) -1

	expectedTask := models.Task{ID: 1, Title: "New Task Title", Description: "New Task Description", UserID: 1}
	mockService.On("CreateTask", mock.Anything, "New Task Title", "New Task Description", (*int64)(nil), (*time.Time)(nil)).Return(expectedTask, nil)

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updatedModel.(*Model)

	assert.True(t, m.isLoading)
	assert.NotNil(t, cmd)

	createMsg := cmd()
	assert.IsType(t, taskCreatedMsg{}, createMsg)

	mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{expectedTask}, nil)
	updatedModelAfterCreate, refreshCmd := m.Update(createMsg)
	m = updatedModelAfterCreate.(*Model)

	assert.Equal(t, TableView, m.currentView) // Check currentView resets
	assert.Nil(t, m.err)
	assert.NotNil(t, refreshCmd)

	refreshMsg := refreshCmd()
	_, ok := refreshMsg.(tasksLoadedMsg)
	assert.True(t, ok)

	mockService.AssertExpectations(t)
}

func TestTasksModel_UpdateTask_SubmitForm(t *testing.T) {
	mockService := new(MockTaskService)
	originalTask := &models.Task{ID: 1, Title: "Original Title", Description: "Original Desc", UserID: 1, IsCompleted: false}

	model := New(mockService)
	model.currentView = FormView // Set initial state for form
	model.formSubState = EditingTask
	model.editingTaskID = originalTask.ID
	model.selectedTaskForDetail = originalTask

	model.inputs[0].SetValue("Updated Title")
	model.inputs[1].SetValue("Updated Desc")
	model.focusIndex = len(model.inputs) - 1

	taskWithUpdates := &models.Task{
		ID:          originalTask.ID, UserID:      originalTask.UserID,
		Title:       "Updated Title", Description: "Updated Desc",
		ClassID:     nil, DueDate:     nil,
		IsCompleted: originalTask.IsCompleted,
	}

	mockService.On("UpdateTask", mock.Anything, mock.MatchedBy(func(task *models.Task) bool {
		return task.ID == taskWithUpdates.ID && task.Title == taskWithUpdates.Title &&
			   task.Description == taskWithUpdates.Description && task.UserID == taskWithUpdates.UserID &&
			   task.IsCompleted == taskWithUpdates.IsCompleted
	})).Return(nil)

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updatedModel.(*Model)
	assert.True(t, m.isLoading)
	assert.NotNil(t, cmd)

	updateMsg := cmd()
	assert.IsType(t, taskUpdatedMsg{}, updateMsg)

	mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{*taskWithUpdates}, nil)
	updatedModelAfterUpdate, refreshCmd := m.Update(updateMsg)
	m = updatedModelAfterUpdate.(*Model)

	assert.Equal(t, TableView, m.currentView) // Check currentView resets
	assert.Equal(t, int64(0), m.editingTaskID)
	assert.NotNil(t, refreshCmd)

	mockService.AssertExpectations(t)
}


func TestTasksModel_KeyBindings_CRUD_OnFocusedTable(t *testing.T) {
	mockService := new(MockTaskService)
	task1Pending := models.Task{ID: 1, Title: "Pending Task 1", UserID: 1, Description: "Desc P1"}
	task2Completed := models.Task{ID: 2, Title: "Completed Task 1", UserID: 1, Description: "Desc C1", IsCompleted: true}

	mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{task1Pending, task2Completed}, nil).Once()
	model := New(mockService)
	model.SetSize(80, 30)
	model.Update(model.Init()())

	assert.Equal(t, PendingTableFocus, model.focusedTable)
	model.pendingTasksTable.SetCursor(0)

	// Edit 'e' - should work on pending
	mockService.On("GetTaskByID", mock.Anything, task1Pending.ID).Return(&task1Pending, nil).Once()
	modelAfterEditKey, cmdEdit := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	assert.True(t, modelAfterEditKey.(*Model).isLoading)
	assert.NotNil(t, cmdEdit)
	modelAfterEditKey.(*Model).Update(cmdEdit())
	assert.Equal(t, FormView, modelAfterEditKey.(*Model).currentView)
	assert.Equal(t, EditingTask, modelAfterEditKey.(*Model).formSubState)
	assert.Equal(t, task1Pending.ID, modelAfterEditKey.(*Model).editingTaskID)

	model = New(mockService)
	mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{task1Pending, task2Completed}, nil).Once()
	model.SetSize(80,30)
	model.Update(model.Init()())
	model.pendingTasksTable.SetCursor(0)

	modelAfterDeleteKey, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	assert.Equal(t, ConfirmDeleteView, modelAfterDeleteKey.(*Model).currentView)
	assert.Equal(t, task1Pending.ID, modelAfterDeleteKey.(*Model).taskIDToDelete)

	model = New(mockService)
	mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{task1Pending, task2Completed}, nil).Once()
	model.SetSize(80,30)
	model.Update(model.Init()())

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.Equal(t, CompletedTableFocus, model.focusedTable)
	model.completedTasksTable.SetCursor(0)

	modelBeforeEditAttempt := *model
	modelAfterEditKeyOnCompleted, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	assert.Equal(t, modelBeforeEditAttempt.currentView, modelAfterEditKeyOnCompleted.(*Model).currentView)
	assert.Equal(t, modelBeforeEditAttempt.formSubState, modelAfterEditKeyOnCompleted.(*Model).formSubState)

	modelAfterDeleteKeyOnCompleted, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	assert.Equal(t, ConfirmDeleteView, modelAfterDeleteKeyOnCompleted.(*Model).currentView)
	assert.Equal(t, task2Completed.ID, modelAfterDeleteKeyOnCompleted.(*Model).taskIDToDelete)

	mockService.AssertExpectations(t)
}


func TestTasksModel_ViewDetails_OnFocusedTable(t *testing.T) {
    mockService := new(MockTaskService)
    taskPending := models.Task{ID: 1, Title: "Pending Detail", UserID: 1}
    taskCompleted := models.Task{ID: 2, Title: "Completed Detail", UserID: 1, IsCompleted: true}

    mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{taskPending, taskCompleted}, nil).Once()
    model := New(mockService)
    model.SetSize(80,30)
    model.Update(model.Init()())
    model.pendingTasksTable.SetCursor(0)
    assert.Equal(t, PendingTableFocus, model.focusedTable)

    mockService.On("GetTaskByID", mock.Anything, taskPending.ID).Return(&taskPending, nil).Once()
    model, cmdViewPending := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
    assert.True(t, model.isLoading)
    model.Update(cmdViewPending())
    assert.Equal(t, DetailView, model.currentView)
    assert.Equal(t, taskPending.ID, model.selectedTaskForDetail.ID)

    model = New(mockService)
    mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{taskPending, taskCompleted}, nil).Once()
    model.SetSize(80,30)
    model.Update(model.Init()())
    model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
    model.completedTasksTable.SetCursor(0)
    assert.Equal(t, CompletedTableFocus, model.focusedTable)

    mockService.On("GetTaskByID", mock.Anything, taskCompleted.ID).Return(&taskCompleted, nil).Once()
    model, cmdViewCompleted := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
    assert.True(t, model.isLoading)
    model.Update(cmdViewCompleted())
    assert.Equal(t, DetailView, model.currentView)
    assert.Equal(t, taskCompleted.ID, model.selectedTaskForDetail.ID)

    mockService.AssertExpectations(t)
}

func TestTasksModel_DeleteTask_ConfirmYes(t *testing.T) {
	mockService := new(MockTaskService)
	model := New(mockService)
	model.currentView = ConfirmDeleteView // Set state for delete confirmation
	model.taskIDToDelete = 1

	mockService.On("DeleteTask", mock.Anything, int64(1)).Return(nil).Once()
	mockService.On("ListAllTasks", mock.Anything).Return([]models.Task{}, nil).Once()

	m, cmdDelete := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	assert.True(t, m.(*Model).isLoading)

	deleteMsg := cmdDelete()
	m, cmdRefresh := m.(*Model).Update(deleteMsg)

	assert.Equal(t, TableView, m.(*Model).currentView) // Check currentView resets
	assert.Equal(t, int64(0), m.(*Model).taskIDToDelete)
	assert.False(t, m.(*Model).isLoading)
	assert.NotNil(t, cmdRefresh)

	mockService.AssertExpectations(t)
}

func TestTasksModel_DeleteTask_ConfirmNo(t *testing.T) {
	mockService := new(MockTaskService)
	model := New(mockService)
	model.currentView = ConfirmDeleteView // Set state for delete confirmation
	model.taskIDToDelete = 1

	m, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	assert.Equal(t, TableView, m.(*Model).currentView) // Check currentView resets
	assert.Equal(t, int64(0), m.(*Model).taskIDToDelete)
	assert.Nil(t, cmd)
	mockService.AssertNotCalled(t, "DeleteTask", mock.Anything, mock.Anything)
}
