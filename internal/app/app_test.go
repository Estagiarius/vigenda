package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"vigenda/internal/service"
	"vigenda/internal/app/dashboard"
	"vigenda/internal/repository/stubs"
)

func setupTestApp(t *testing.T) (*Model, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	taskRepo := stubs.NewMockTaskRepository(ctrl)
	classRepo := stubs.NewMockClassRepository(ctrl)
	assessmentRepo := stubs.NewMockAssessmentRepository(ctrl)
	questionRepo := stubs.NewMockQuestionRepository(ctrl)
	subjectRepo := stubs.NewMockSubjectRepository(ctrl)
	lessonRepo := stubs.NewMockLessonRepository(ctrl)

	taskService := service.NewTaskService(taskRepo)
	classService := service.NewClassService(classRepo, subjectRepo)
	assessmentService := service.NewAssessmentService(assessmentRepo, classRepo)
	questionService := service.NewQuestionService(questionRepo, subjectRepo)
	proofService := service.NewProofService(questionRepo)
	lessonService := service.NewLessonService(lessonRepo, classRepo)

	// Criar e inicializar o dashboard explicitamente
	dashboardModel := dashboard.New(taskService, classService, assessmentService, lessonService)
	dashboardModel.Init() // Para carregar os dados iniciais

	app := New(taskService, classService, assessmentService, questionService, proofService, lessonService)
	app.dashboardModel = dashboardModel
	app.currentView = ConcreteDashboardView // Inicia no dashboard concreto
	return app, ctrl
}

func TestApp_InitialViewIsConcreteDashboard(t *testing.T) {
	app, ctrl := setupTestApp(t)
	defer ctrl.Finish()
	assert.Equal(t, ConcreteDashboardView, app.currentView)
}

func TestApp_NavigateToMenuAndSelectTaskView(t *testing.T) {
	app, ctrl := setupTestApp(t)
	defer ctrl.Finish()

	// 1. Pressiona 'm' para abrir o menu
	msgMenu := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}}
	updatedModel, _ := app.Update(msgMenu)
	app = updatedModel.(*Model)
	assert.Equal(t, DashboardView, app.currentView)

	// 2. Seleciona "Gerenciar Tarefas"
	for i, item := range app.list.Items() {
		if item.(menuItem).view == TaskManagementView {
			app.list.Select(i)
			break
		}
	}

	// 3. Pressiona Enter para confirmar a seleção
	msgEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = app.Update(msgEnter)
	app = updatedModel.(*Model)

	// A view deve mudar para TaskManagementView
	assert.Equal(t, TaskManagementView, app.currentView)
}

func TestApp_ReturnToDashboardFromTaskView(t *testing.T) {
	app, ctrl := setupTestApp(t)
	defer ctrl.Finish()

	// Define o estado inicial como TaskManagementView
	app.currentView = TaskManagementView

	// Simula pressionar Esc
	msgEsc := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := app.Update(msgEsc)
	app = updatedModel.(*Model)

	// Deve voltar para o ConcreteDashboardView
	assert.Equal(t, ConcreteDashboardView, app.currentView)
}

func TestApp_QuitApplicationFromMenu(t *testing.T) {
	app, ctrl := setupTestApp(t)
	defer ctrl.Finish()

	// 1. Abre o menu
	app.currentView = DashboardView

	// 2. Simula pressionar 'q' para sair
	msgQuit := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd := app.Update(msgQuit)
	app = updatedModel.(*Model)

	assert.True(t, app.quitting)
	assert.NotNil(t, cmd) // Deve retornar um comando de Quit
	if quitCmd, ok := cmd().(tea.QuitMsg); !ok {
		assert.Fail(t, "O comando retornado não foi tea.QuitMsg", "Recebido: %T", quitCmd)
	}
}

func TestApp_GlobalQuitWithCtrlC(t *testing.T) {
	app, ctrl := setupTestApp(t)
	defer ctrl.Finish()
	app.currentView = TaskManagementView // Em uma view qualquer

	// Simula pressionar Ctrl+C
	msgCtrlC := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd := app.Update(msgCtrlC)
	app = updatedModel.(*Model)

	assert.True(t, app.quitting)
	assert.NotNil(t, cmd)
	if quitCmd, ok := cmd().(tea.QuitMsg); !ok {
		assert.Fail(t, "O comando retornado não foi tea.QuitMsg", "Recebido: %T", quitCmd)
	}
}
