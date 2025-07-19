package app

import (
	"context"
	"testing"
	"time"
	"vigenda/internal/models"
	"vigenda/internal/service"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock Service Implementations ---

type mockTaskService struct{ service.TaskService }
func (m *mockTaskService) ListAllTasks(ctx context.Context) ([]models.Task, error) { return []models.Task{}, nil }
func (m *mockTaskService) GetUpcomingActiveTasks(ctx context.Context, userID int64, fromDate time.Time, limit int) ([]models.Task, error) { return []models.Task{}, nil }
func (m *mockTaskService) ListAllActiveTasks(ctx context.Context) ([]models.Task, error) { return []models.Task{}, nil }

type mockClassService struct{ service.ClassService }
func (m *mockClassService) ListAllClasses(ctx context.Context) ([]models.Class, error) { return []models.Class{}, nil }
func (m *mockClassService) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) { return []models.Student{}, nil }


type mockAssessmentService struct{ service.AssessmentService }
func (m *mockAssessmentService) ListAllAssessments(ctx context.Context) ([]models.Assessment, error) { return []models.Assessment{}, nil }


type mockQuestionService struct{ service.QuestionService }
type mockProofService struct{ service.ProofService }
type mockLessonService struct{ service.LessonService }
func (m *mockLessonService) GetLessonsForDate(ctx context.Context, userID int64, date time.Time) ([]models.Lesson, error) { return []models.Lesson{}, nil }

type mockSubjectService struct{ service.SubjectService }

// Helper to create a new model with all mock services
func newTestAppModel() *Model {
	return New(
		&mockTaskService{},
		&mockClassService{},
		&mockAssessmentService{},
		&mockQuestionService{},
		&mockProofService{},
		&mockLessonService{},
		&mockSubjectService{},
	)
}

func TestNewModel_InitialState(t *testing.T) {
	m := newTestAppModel()
	assert.Equal(t, DashboardView, m.currentView, "Initial view should be DashboardView")
	require.Greater(t, len(m.list.Items()), 0, "List should have items")
	assert.Contains(t, m.list.Items()[0].(menuItem).Title(), "Painel de Controle")
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
	m.list.Select(1) // Select "Gerenciar Tarefas"
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	nextModelTea, _ := m.Update(enterMsg)
	m = nextModelTea.(*Model)

	assert.Equal(t, TaskManagementView, m.currentView, "View should change to TaskManagementView on Enter")

	// Simulate pressing 'esc' to go back
	if m.tasksModel != nil {
		// This part of the test is tricky without knowing the sub-model's internal state management.
		// We assume the sub-model has a mechanism to signal it can go back.
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

	// Navigate to a sub-view (TaskManagementView)
	m.list.Select(1) // "Gerenciar Tarefas"
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

	var cmd tea.Cmd
	updatedModelTea, cmd = m.Update(enterMsg)
	m = updatedModelTea.(*Model)

	if cmd != nil {
		msg := cmd()
		updatedModelTea, _ = m.Update(msg)
		m = updatedModelTea.(*Model)
	}

	viewOutput = m.View()
	assert.Contains(t, viewOutput, "Pressione 'esc' para voltar ao menu principal.")
}


func TestModel_Update_WindowSize(t *testing.T) {
	m := newTestAppModel()
	newWidth, newHeight := 80, 24
	sizeMsg := tea.WindowSizeMsg{Width: newWidth, Height: newHeight}
	nextModelTea, _ := m.Update(sizeMsg)
	m = nextModelTea.(*Model)

	assert.Equal(t, newWidth, m.width, "Model width should be updated")
	assert.Equal(t, newHeight, m.height, "Model height should be updated")
}

func TestMenuItem_Interface(t *testing.T) {
	item := menuItem{title: "Test Title", view: DashboardView}
	assert.Equal(t, "Test Title", item.Title())
	assert.Equal(t, "Test Title", item.FilterValue())
	assert.Equal(t, "", item.Description())
}

func TestView_String(t *testing.T) {
	assert.Equal(t, "Menu Principal", DashboardView.String())
	assert.Equal(t, "Gerenciar Tarefas", TaskManagementView.String())
}
