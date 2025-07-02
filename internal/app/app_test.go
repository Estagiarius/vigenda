package app

import (
	"context"
	"testing"
	"time"
	"vigenda/internal/models"
	"vigenda/internal/service"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock Service Implementations ---

type mockTaskService struct{}

func (m *mockTaskService) CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error) {
	return models.Task{}, nil
}
func (m *mockTaskService) ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error) {
	return nil, nil
}
func (m *mockTaskService) ListAllActiveTasks(ctx context.Context) ([]models.Task, error) {
	return nil, nil
}
func (m *mockTaskService) MarkTaskAsCompleted(ctx context.Context, taskID int64) error { return nil }
func (m *mockTaskService) GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error) {
	return nil, nil
}
func (m *mockTaskService) UpdateTask(ctx context.Context, task *models.Task) error { return nil }
func (m *mockTaskService) DeleteTask(ctx context.Context, taskID int64) error      { return nil }

type mockClassService struct{}

func (m *mockClassService) CreateClass(ctx context.Context, name string, subjectID int64) (models.Class, error) {
	return models.Class{}, nil
}
func (m *mockClassService) ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error) {
	return 0, nil
}
func (m *mockClassService) UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error {
	return nil
}
func (m *mockClassService) GetClassByID(ctx context.Context, classID int64) (*models.Class, error) { // Interface expects *models.Class
	return &models.Class{}, nil
}
func (m *mockClassService) ListAllClasses(ctx context.Context) ([]models.Class, error) { return nil, nil }
func (m *mockClassService) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
	return nil, nil
}

type mockAssessmentService struct{}

func (m *mockAssessmentService) CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error) {
	return models.Assessment{}, nil
}
func (m *mockAssessmentService) EnterGrades(ctx context.Context, assessmentID int64, studentGrades map[int64]float64) error {
	return nil
}
func (m *mockAssessmentService) CalculateClassAverage(ctx context.Context, classID int64) (float64, error) {
	return 0, nil
}
func (m *mockAssessmentService) ListAllAssessments(ctx context.Context) ([]models.Assessment, error) {
	return nil, nil
}

type mockQuestionService struct{}

func (m *mockQuestionService) AddQuestionsFromJSON(ctx context.Context, jsonData []byte) (int, error) {
	return 0, nil
}
func (m *mockQuestionService) GenerateTest(ctx context.Context, criteria service.TestCriteria) ([]models.Question, error) {
	return nil, nil
}

type mockProofService struct{}

func (m *mockProofService) GenerateProof(ctx context.Context, criteria service.ProofCriteria) ([]models.Question, error) {
	return nil, nil
}

// Helper to create a new model with all mock services
func newTestAppModel() *Model {
	return New(
		&mockTaskService{},
		&mockClassService{},
		&mockAssessmentService{},
		&mockQuestionService{},
		&mockProofService{},
	)
}

func TestNewModel_InitialState(t *testing.T) {
	m := newTestAppModel()
	assert.Equal(t, DashboardView, m.currentView, "Initial view should be DashboardView")
	require.Greater(t, len(m.list.Items()), 0, "List should have items")
	assert.Contains(t, m.list.Items()[0].(menuItem).Title(), DashboardView.String(), "First item should be Dashboard")
}

func TestModel_Update_Quit(t *testing.T) {
	m := newTestAppModel()
	// Test 'q' from DashboardView
	qMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	nextModelTea, cmd := m.Update(qMsg)
	nextModel := nextModelTea.(*Model)
	assert.True(t, nextModel.quitting, "Model should be quitting on 'q'")
	assert.NotNil(t, cmd, "A command (tea.Quit) should be returned on 'q'")

	// Test 'ctrl+c'
	m = newTestAppModel() // Reset model
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	nextModelTea, cmd = m.Update(ctrlCMsg)
	nextModel = nextModelTea.(*Model)
	assert.True(t, nextModel.quitting, "Model should be quitting on 'ctrl+c'")
	assert.NotNil(t, cmd, "A command (tea.Quit) should be returned on 'ctrl+c'")
}

func TestModel_Update_NavigateToSubViewAndBack(t *testing.T) {
	m := newTestAppModel()
	initialView := m.currentView

	// Simulate selecting the second item (TaskManagementView)
	m.list.SetCursor(1) // Select "Gerenciar Tarefas"
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	nextModelTea, _ := m.Update(enterMsg)
	m = nextModelTea.(*Model)

	assert.Equal(t, TaskManagementView, m.currentView, "View should change to TaskManagementView on Enter")

	// Simulate pressing 'esc' to go back
	// For sub-views, 'esc' is handled by the sub-model if IsFocused, or by app model if not.
	// Assuming tasksModel.IsFocused() is false when just viewing the table.
	if m.tasksModel != nil { // Ensure tasksModel is initialized
		m.tasksModel.formState = NoForm // Ensure it's not in a focused state like form editing
	}
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	nextModelTea, _ = m.Update(escMsg)
	m = nextModelTea.(*Model)

	assert.Equal(t, initialView, m.currentView, "View should change back to DashboardView on Esc")
}


func TestModel_View_Content(t *testing.T) {
	m := newTestAppModel()
	// Set initial size for consistent rendering of the list
	updatedModelTea, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updatedModelTea.(*Model)

	// Initial view (Dashboard/Menu)
	viewOutput := m.View()
	assert.Contains(t, viewOutput, m.list.Title, "View should contain list title in DashboardView")
	assert.Contains(t, viewOutput, "Navegue com ↑/↓", "View should contain help text for menu")

	// Navigate to a sub-view (TaskManagementView)
	m.list.SetCursor(1) // "Gerenciar Tarefas"
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

	// If tasksModel.Init() dispatches a command, we need to handle its message
	var cmd tea.Cmd
	updatedModelTea, cmd = m.Update(enterMsg)
	m = updatedModelTea.(*Model)

	if cmd != nil {
		msg := cmd() // Execute command from tasksModel.Init()
		updatedModelTea, _ = m.Update(msg) // Process the message (e.g., tasksLoadedMsg)
		m = updatedModelTea.(*Model)
	}

	viewOutput = m.View()
	// tasksModel.View() will render its own content.
	// For this test, check that the main app help text changes.
	assert.Contains(t, viewOutput, "Navegue na tabela com ↑/↓. Pressione 'esc' para voltar ao menu.", "View should contain help text for subview")
}


func TestModel_Update_WindowSize(t *testing.T) {
	m := newTestAppModel()
	newWidth, newHeight := 80, 24
	sizeMsg := tea.WindowSizeMsg{Width: newWidth, Height: newHeight}
	nextModelTea, _ := m.Update(sizeMsg)
	m = nextModelTea.(*Model)

	assert.Equal(t, newWidth, m.width, "Model width should be updated")
	assert.Equal(t, newHeight, m.height, "Model height should be updated")

	expectedListWidth := newWidth - appStyle.GetHorizontalPadding()
	// Calculation based on app.go: listHeight := msg.Height - appStyle.GetVerticalPadding() - lipgloss.Height(m.list.Title) - 2
	// appStyle.GetVerticalPadding() is 2 (1 top, 1 bottom)
	// lipgloss.Height(m.list.Title) for "Vigenda - Menu Principal" with MarginBottom(1) is 2
	// So, listHeight = 24 - 2 - 2 - 2 = 18
	expectedListHeight := newHeight - appStyle.GetVerticalPadding() - lipgloss.Height(m.list.Title) - 2

	assert.Equal(t, expectedListWidth, m.list.Width(), "List width should be updated based on window size")
	assert.Equal(t, expectedListHeight, m.list.Height(), "List height should be updated based on window size")
}

func TestMenuItem_Interface(t *testing.T) {
	item := menuItem{title: "Test Title", view: DashboardView}
	assert.Equal(t, "Test Title", item.Title())
	assert.Equal(t, "Test Title", item.FilterValue())
	assert.Equal(t, "", item.Description())
}

func TestView_String(t *testing.T) {
	assert.Equal(t, "Dashboard", DashboardView.String())
	assert.Equal(t, "Gerenciar Tarefas", TaskManagementView.String())
}

func simulateKeyPress(m *Model, key rune) (*Model, tea.Cmd) {
	nextModelTea, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}})
	return nextModelTea.(*Model), cmd
}

func simulateEnterPress(m *Model) (*Model, tea.Cmd) {
	nextModelTea, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	return nextModelTea.(*Model), cmd
}

func simulateEscPress(m *Model) (*Model, tea.Cmd) {
	nextModelTea, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	return nextModelTea.(*Model), cmd
}

func simulateCtrlCPress(m *Model) (*Model, tea.Cmd) {
	nextModelTea, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	return nextModelTea.(*Model), cmd
}

func TestModel_Update_WithHelpers(t *testing.T) {
	m := newTestAppModel()

	m.list.SetCursor(1) // Select "Gerenciar Tarefas"
	var cmd tea.Cmd
	m, cmd = simulateEnterPress(m)
	// Handle potential command from sub-model Init
	if cmd != nil {
		msg := cmd()
		m, _ = m.Update(msg)
	}
	assert.Equal(t, TaskManagementView, m.currentView)

	// To test 'esc' from subview, ensure subview is not focused (e.g., tasksModel.IsFocused() is false)
	if m.tasksModel != nil {
		m.tasksModel.formState = NoForm // Not in a form
	}
	m, _ = simulateEscPress(m)
	assert.Equal(t, DashboardView, m.currentView)

	m, cmd = simulateCtrlCPress(m)
	assert.True(t, m.quitting)
	assert.NotNil(t, cmd)
}
