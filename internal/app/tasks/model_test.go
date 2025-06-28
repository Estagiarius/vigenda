package tasks

import (
	"context"
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vigenda/internal/models"
	// "vigenda/internal/service" // Not directly used, mock implements the interface
)

// MockTaskService for testing
type mockTaskService struct {
	tasksToReturn []models.Task
	errorToReturn error
}

func (m *mockTaskService) CreateTask(ctx context.Context, title string, description string, classID *int64, dueDate *time.Time) (models.Task, error) {
	// Return a zero-value Task and an error for the mock, as it's not the focus of these tests.
	return models.Task{}, errors.New("not implemented in mock for CreateTask")
}
func (m *mockTaskService) ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error) {
	return nil, errors.New("not implemented in mock")
}
func (m *mockTaskService) ListAllActiveTasks(ctx context.Context) ([]models.Task, error) {
	if m.errorToReturn != nil {
		return nil, m.errorToReturn
	}
	return m.tasksToReturn, nil
}
func (m *mockTaskService) MarkTaskAsCompleted(ctx context.Context, taskID int64) error {
	return errors.New("not implemented in mock")
}
func (m *mockTaskService) GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error) {
	return nil, errors.New("not implemented in mock")
}


func TestTasksModel_New(t *testing.T) {
	mockService := &mockTaskService{}
	m := New(mockService)

	assert.NotNil(t, m.table, "Table should be initialized")
	assert.Equal(t, mockService, m.taskService, "TaskService should be stored")
	assert.True(t, m.isLoading, "isLoading should be true on New")

	expectedColumns := []string{"ID", "Título", "Prazo", "Concluída", "ID Turma"}
	actualColumns := m.table.Columns()
	require.Len(t, actualColumns, len(expectedColumns), "Number of columns should match")
	for i, col := range actualColumns {
		assert.Equal(t, expectedColumns[i], col.Title, "Column title mismatch")
	}
}

func TestTasksModel_Init(t *testing.T) {
	m := New(&mockTaskService{})
	// Init sets isLoading to true and returns the command
	cmd := m.Init()
	require.NotNil(t, cmd, "Init should return a command")
	assert.True(t, m.isLoading, "isLoading should be true after Init is called")

	// Check if the command is loadTasksCmd by evaluating it
	msg := cmd()
	_, ok := msg.(tasksLoadedMsg)
	assert.True(t, ok, "Init command should produce a tasksLoadedMsg")
}

func TestTasksModel_Update_TasksLoadedSuccessfully(t *testing.T) {
	now := time.Now()
	classID := int64(1)
	sampleTasks := []models.Task{
		{ID: 1, Title: "Task 1", DueDate: &now, IsCompleted: false, ClassID: &classID},
		{ID: 2, Title: "Task 2", IsCompleted: true},
	}
	mockService := &mockTaskService{tasksToReturn: sampleTasks}
	m := New(mockService)

	// Simulate WindowSizeMsg first for table to have dimensions
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Simulate tasksLoadedMsg
	msg := tasksLoadedMsg{tasks: sampleTasks, err: nil}
	updatedModel, cmd := m.Update(msg)
	require.Nil(t, cmd, "No command should be returned after loading tasks")

	m = updatedModel // No type assertion needed, it's already tasks.Model
	assert.Nil(t, m.err, "Error should be nil on successful load")
	assert.False(t, m.isLoading, "isLoading should be false after tasks are loaded")
	assert.Len(t, m.table.Rows(), len(sampleTasks), "Table should have correct number of rows")

	// Check first row content
	expectedRow1 := []string{
		"1",
		"Task 1",
		now.Format("02/01/2006"),
		"false",
		"1",
	}
	assert.Equal(t, expectedRow1, []string(m.table.Rows()[0]), "First row content mismatch")
}

func TestTasksModel_Update_TasksLoadedWithError(t *testing.T) {
	loadErr := errors.New("failed to load tasks")
	mockService := &mockTaskService{errorToReturn: loadErr}
	m := New(mockService)

	msg := tasksLoadedMsg{tasks: nil, err: loadErr}
	updatedModel, cmd := m.Update(msg)
	require.Nil(t, cmd)

	m = updatedModel // No type assertion needed
	assert.Equal(t, loadErr, m.err, "Error should be set")
	assert.False(t, m.isLoading, "isLoading should be false after loading error")
	assert.Empty(t, m.table.Rows(), "Table should be empty on error")
}

func TestTasksModel_View(t *testing.T) {
	t.Run("View when loading", func(t *testing.T) {
		mockService := &mockTaskService{}
		m := New(mockService) // isLoading is true by default
		// m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24}) // Size doesn't affect loading message

		view := m.View()
		assert.Contains(t, view, "Carregando tarefas...", "View should show loading message")
	})

	t.Run("View with loaded tasks", func(t *testing.T) {
		sampleTasks := []models.Task{{ID: 1, Title: "View Task"}}
		mockService := &mockTaskService{tasksToReturn: sampleTasks}
		m := New(mockService)
		m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24}) // Set size
		// Correctly call loadTasksCmd which is a tea.Cmd (a function)
		// and then pass its result (a tasksLoadedMsg) to Update.
		loadedMsg := m.loadTasksCmd()
		m, _ = m.Update(loadedMsg)


		view := m.View()
		assert.Contains(t, view, "View Task", "View should contain task title")
		assert.NotContains(t, view, "Erro ao carregar tarefas", "View should not show error message")
	})

	t.Run("View with error", func(t *testing.T) {
		loadErr := errors.New("view load error")
		mockService := &mockTaskService{errorToReturn: loadErr}
		m := New(mockService)
		// Correctly call loadTasksCmd and pass its result to Update
		loadedMsgWithError := m.loadTasksCmd()
		m, _ = m.Update(loadedMsgWithError)


		view := m.View()
		assert.Contains(t, view, "Erro ao carregar tarefas: view load error", "View should show error message")
	})
}

func TestTasksModel_SetSize(t *testing.T) {
	m := New(&mockTaskService{})
	newWidth, newHeight := 100, 30
	m.SetSize(newWidth, newHeight)

	// baseStyle has NormalBorder, horizontal frame size is 2.
	expectedTableWidth := newWidth - 2 // baseStyle.GetHorizontalFrameSize()
	// tasks.Model.View() is baseStyle.Render(m.table.View()) + "\n".
	// Vertical: 2 for baseStyle border, 1 for \n, 1 for table header. Total 4.
	expectedTableViewportHeight := newHeight - 4

	assert.Equal(t, expectedTableWidth, m.table.Width(), "Table width should be updated")
	assert.Equal(t, expectedTableViewportHeight, m.table.Height(), "Table viewport height should be updated")
}

func TestTasksModel_IsFocused(t *testing.T) {
	m := New(&mockTaskService{})
	assert.False(t, m.IsFocused(), "IsFocused should return false for now")
}

func TestTasksModel_Update_TableNavigation(t *testing.T) {
	sampleTasks := []models.Task{
		{ID: 1, Title: "Task A"},
		{ID: 2, Title: "Task B"},
		{ID: 3, Title: "Task C"},
	}
	mockService := &mockTaskService{tasksToReturn: sampleTasks}
	m := New(mockService)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	loadedMsg := m.loadTasksCmd() // Get the message from the command
	m, _ = m.Update(loadedMsg)    // Send the message to Update

	require.Len(t, m.table.Rows(), 3)
	assert.Equal(t, 0, m.table.Cursor(), "Cursor should initially be at the first row")

	// Simulate ArrowDown
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown}) // Corrected to tea.KeyDown
	assert.Equal(t, 1, m.table.Cursor(), "Cursor should move down")

	// Simulate ArrowUp
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp}) // Corrected to tea.KeyUp
	assert.Equal(t, 0, m.table.Cursor(), "Cursor should move up")
}
