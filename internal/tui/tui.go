// Package tui implements the Text User Interface (TUI) for the Vigenda application.
// It uses the Bubble Tea library and its components to create an interactive CLI experience.
package tui

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"vigenda/internal/app"
	"vigenda/internal/models"
	"vigenda/internal/service"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	docStyle           = lipgloss.NewStyle().Margin(1, 2)
	titleStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))  // Purple
	selectedStyle      = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("62")).PaddingLeft(1)
	deselectedStyle    = lipgloss.NewStyle().PaddingLeft(1)
	helpStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray
	statusMessageStyle = lipgloss.NewStyle().MarginTop(1).MarginBottom(1).Foreground(lipgloss.Color("202"))
)

// KeyMap
type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Select       key.Binding
	Back         key.Binding
	Quit         key.Binding
	Help         key.Binding
	CompleteTask key.Binding
}

var DefaultKeyMap = KeyMap{
	Up:           key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "cima")),
	Down:         key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "baixo")),
	Select:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "slcnar")),
	Back:         key.NewBinding(key.WithKeys("esc", "backspace"), key.WithHelp("esc", "voltar")),
	Quit:         key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q/ctrl+c", "sair")),
	Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "ajuda")),
	CompleteTask: key.NewBinding(key.WithKeys("x", "c"), key.WithHelp("x/c", "concluir")),
}

// Model
type Model struct {
	classService      service.ClassService
	taskService       service.TaskService
	assessmentService service.AssessmentService
	list              list.Model
	spinner           spinner.Model
	keys              KeyMap
	currentView       app.View
	isLoading         bool
	err               error
	statusMessage     string
	mainMenuItems     []list.Item
	classes           []models.Class
	students          []models.Student
	dashboardItems    []list.Item
	upcomingTasks     []models.Task
	selectedClass     *models.Class
}

// mainMenuItem
type mainMenuItem struct {
	title       string
	targetView  app.View
	description string
}

func (mmi mainMenuItem) Title() string       { return mmi.title }
func (mmi mainMenuItem) Description() string { return mmi.description }
func (mmi mainMenuItem) FilterValue() string { return mmi.title }

func NewTUIModel(cs service.ClassService, ts service.TaskService, as service.AssessmentService) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	delegate := list.NewDefaultDelegate()
	mainList := list.New(nil, delegate, 0, 0)
	mainList.SetShowHelp(false)
	m := Model{
		classService:      cs,
		taskService:       ts,
		assessmentService: as,
		spinner:           s,
		keys:              DefaultKeyMap,
		currentView:       app.MainMenuView,
		isLoading:         false,
		list:              mainList,
	}
	log.Println("TUI: NewTUIModel - currentView:", m.currentView.String())
	return m
}

const statusMessageDuration = 3 * time.Second
type clearStatusMessageMsg struct{}

func setStatusMessageCmd(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return clearStatusMessageMsg{}
	})
}

// Command/Data Loading functions
func (m *Model) loadMainMenuItemsCmd() tea.Cmd {
	m.isLoading = true
	return func() tea.Msg {
		items := []list.Item{
			mainMenuItem{title: app.DashboardView.String(), targetView: app.DashboardView, description: "Visão geral e tarefas pendentes."},
			mainMenuItem{title: app.TaskManagementView.String(), targetView: app.TaskManagementView, description: "Criar, listar e gerenciar tarefas."},
			mainMenuItem{title: app.ClassManagementView.String(), targetView: app.ClassManagementView, description: "Gerenciar turmas e alunos."},
			mainMenuItem{title: app.AssessmentManagementView.String(), targetView: app.AssessmentManagementView, description: "Gerenciar avaliações e notas."},
			mainMenuItem{title: app.QuestionBankView.String(), targetView: app.QuestionBankView, description: "Acessar banco de questões."},
			mainMenuItem{title: app.ProofGenerationView.String(), targetView: app.ProofGenerationView, description: "Gerar provas."},
		}
		return mainMenuLoadedMsg(items)
	}
}

func (m *Model) loadAndDisplayTasks() tea.Cmd { // Renamed from loadDashboardData for clarity
	m.isLoading = true
	return func() tea.Msg {
		tasks, err := m.taskService.ListAllActiveTasks(context.Background())
		if err != nil {
			return errMsg{err: err, context: "loading tasks for view"}
		}
		return tasksForViewLoadedMsg(tasks) // New message type
	}
}

func (m *Model) loadClasses() tea.Cmd {
	m.isLoading = true
	return func() tea.Msg {
		classes, err := m.classService.ListAllClasses(context.Background())
		if err != nil {
			return errMsg{err: err, context: "loading classes"}
		}
		return classesLoadedMsg(classes)
	}
}

func (m *Model) loadStudentsForClass(classID int64) tea.Cmd {
	m.isLoading = true
	return func() tea.Msg {
		students, err := m.classService.GetStudentsByClassID(context.Background(), classID)
		if err != nil {
			return errMsg{err: err, context: fmt.Sprintf("loading students for class %d", classID)}
		}
		return studentsLoadedMsg(students)
	}
}

// Init
func (m Model) Init() tea.Cmd {
	var initialCmd tea.Cmd
	if m.currentView == app.MainMenuView {
		initialCmd = m.loadMainMenuItemsCmd()
	} else {
		log.Printf("TUI: Init - Visão inicial inesperada: %s. Carregando menu principal.", m.currentView.String())
		initialCmd = m.loadMainMenuItemsCmd()
	}
	return tea.Batch(m.spinner.Tick, initialCmd)
}

// Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	goBackToMainMenu := func() tea.Model {
		m.currentView = app.MainMenuView
		m.list.Title = app.MainMenuView.String()
		m.statusMessage = ""
		m.dashboardItems = nil; m.upcomingTasks = nil; m.classes = nil; m.students = nil; m.selectedClass = nil
		cmds = append(cmds, m.loadMainMenuItemsCmd())
		return m
	}

	setPlaceholderView := func(targetView app.View) {
		m.currentView = targetView
		m.list.Title = targetView.String()
		m.statusMessage = ""
		m.list.SetItems([]list.Item{placeholderItem{title: fmt.Sprintf("Visão %s em desenvolvimento.", targetView.String())}})
		m.isLoading = false
	}

	handleTaskCompletion := func() {
		selectedItem := m.list.SelectedItem()
		if taskItem, ok := selectedItem.(dashboardTaskItem); ok {
			err := m.taskService.MarkTaskAsCompleted(context.Background(), taskItem.Task.ID)
			if err != nil {
				m.statusMessage = fmt.Sprintf("Erro: %v", err)
			} else {
				m.statusMessage = fmt.Sprintf("Tarefa '%s' concluída!", taskItem.Title())
				m.isLoading = true
				// Reload tasks for the current view (Dashboard or TaskManagement)
				cmds = append(cmds, m.loadAndDisplayTasks())
			}
			cmds = append(cmds, setStatusMessageCmd(statusMessageDuration))
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) { return m, tea.Quit }

		switch m.currentView {
		case app.MainMenuView:
			if key.Matches(msg, m.keys.Select) {
				selected, ok := m.list.SelectedItem().(mainMenuItem)
				if ok {
					m.list.SetItems(nil); m.isLoading = true; m.statusMessage = ""
					switch selected.targetView {
					case app.DashboardView, app.TaskManagementView: // Both will list tasks for now
						m.currentView = selected.targetView
						m.list.Title = selected.targetView.String()
						cmds = append(cmds, m.loadAndDisplayTasks())
					case app.ClassManagementView:
						m.currentView = app.ClassManagementView
						m.list.Title = app.ClassManagementView.String()
						cmds = append(cmds, m.loadClasses())
					case app.AssessmentManagementView, app.QuestionBankView, app.ProofGenerationView:
						setPlaceholderView(selected.targetView)
					default: setPlaceholderView(selected.targetView)
					}
				}
			} else if key.Matches(msg, m.keys.Back) { return m, tea.Quit }
			m.list, cmd = m.list.Update(msg); cmds = append(cmds, cmd)

		case app.DashboardView, app.TaskManagementView: // Combined logic for task listing views
			if key.Matches(msg, m.keys.CompleteTask) { handleTaskCompletion() }
			if key.Matches(msg, m.keys.Back) { return goBackToMainMenu(), tea.Batch(cmds...) }
			m.list, cmd = m.list.Update(msg); cmds = append(cmds, cmd)

		case app.ClassManagementView:
			if key.Matches(msg, m.keys.Select) {
				selected, ok := m.list.SelectedItem().(listItemClass)
				if ok {
					m.selectedClass = &selected.Class; m.currentView = app.StudentListView
					m.list.Title = fmt.Sprintf("%s - %s", app.StudentListView.String(), m.selectedClass.Name)
					m.list.SetItems(nil); m.students = nil; m.isLoading = true
					cmds = append(cmds, m.loadStudentsForClass(selected.ID()))
				}
			} else if key.Matches(msg, m.keys.Back) { return goBackToMainMenu(), tea.Batch(cmds...) }
			m.list, cmd = m.list.Update(msg); cmds = append(cmds, cmd)

		case app.StudentListView:
			if key.Matches(msg, m.keys.Back) {
				m.currentView = app.ClassManagementView; m.list.Title = app.ClassManagementView.String()
				m.selectedClass = nil; m.list.SetItems(nil); m.students = nil; m.isLoading = true
				cmds = append(cmds, m.loadClasses())
			}
			m.list, cmd = m.list.Update(msg); cmds = append(cmds, cmd)

		case app.AssessmentManagementView, app.QuestionBankView, app.ProofGenerationView: // Placeholder views
			if key.Matches(msg, m.keys.Back) { return goBackToMainMenu(), tea.Batch(cmds...) }
			m.list, cmd = m.list.Update(msg); cmds = append(cmds, cmd)
		}

	case spinner.TickMsg:
		if m.isLoading { m.spinner, cmd = m.spinner.Update(msg); cmds = append(cmds, cmd) }

	case clearStatusMessageMsg: m.statusMessage = ""

	case mainMenuLoadedMsg:
		m.isLoading = false; m.mainMenuItems = []list.Item(msg)
		m.list.SetItems(m.mainMenuItems); m.list.Title = app.MainMenuView.String()

	case tasksForViewLoadedMsg: // Handles tasks for Dashboard and TaskManagementView
		m.isLoading = false; m.upcomingTasks = []models.Task(msg)
		m.dashboardItems = []list.Item{} // Use dashboardItems for consistency in display
		if len(m.upcomingTasks) == 0 {
			m.dashboardItems = append(m.dashboardItems, placeholderItem{title: "Nenhuma tarefa ativa encontrada."})
		} else {
			for _, task := range m.upcomingTasks {
				m.dashboardItems = append(m.dashboardItems, dashboardTaskItem{task})
			}
		}
		m.list.SetItems(m.dashboardItems)
		// Title is already set by navigation logic

	case classesLoadedMsg:
		m.isLoading = false; m.classes = []models.Class(msg)
		items := make([]list.Item, len(m.classes))
		if len(m.classes) == 0 { items = append(items, placeholderItem{title: "Nenhuma turma encontrada."}) } else {
			for i, c := range m.classes { items[i] = listItemClass{c} }
		}
		m.list.SetItems(items)

	case studentsLoadedMsg:
		m.isLoading = false; m.students = []models.Student(msg)
		items := make([]list.Item, len(m.students))
		if len(m.students) == 0 { items = append(items, placeholderItem{title: "Nenhum aluno encontrado."}) } else {
			for i, s := range m.students { items[i] = listItemStudent{s} }
		}
		m.list.SetItems(items)

	case errMsg:
		m.isLoading = false; m.err = msg.err

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize(); statusMessageHeight := 0
		if m.statusMessage != "" { statusMessageHeight = lipgloss.Height(statusMessageStyle.Render(m.statusMessage)) +1 }
		listHeight := msg.Height - v - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView()) - statusMessageHeight
		m.list.SetSize(msg.Width-h, listHeight)
	}
	return m, tea.Batch(cmds...)
}

// View
func (m Model) View() string {
	if m.err != nil { return docStyle.Render(fmt.Sprintf("Erro: %v\n\nSair (q/ctrl+c).", m.err)) }

	var sb strings.Builder
	sb.WriteString(m.headerView() + "\n")
	if m.statusMessage != "" { sb.WriteString(statusMessageStyle.Render(m.statusMessage) + "\n") }

	if m.isLoading {
		title := m.list.Title; if title == "" { title = m.currentView.String() }
		sb.WriteString(fmt.Sprintf("\n%s Carregando %s...", m.spinner.View(), title))
	} else {
		if m.currentView == app.DashboardView { // Explicit section title for dashboard
			sb.WriteString(titleStyle.Copy().Bold(false).Underline(true).Render("Tarefas Pendentes") + "\n")
		}
		sb.WriteString(m.list.View())
	}
	sb.WriteString("\n" + m.footerView())
	return docStyle.Render(sb.String())
}

func (m Model) headerView() string {
	titleStr := m.list.Title; if titleStr == "" { titleStr = m.currentView.String() }
	if m.currentView == app.StudentListView && m.selectedClass != nil {
		 titleStr = fmt.Sprintf("%s - %s", app.StudentListView.String(), m.selectedClass.Name)
	}
	return titleStyle.Render(titleStr)
}

func (m Model) footerView() string {
	parts := []string{
		m.keys.Up.Help().Key+"/"+m.keys.Down.Help().Key+": "+m.keys.Up.Help().Desc,
		m.keys.Select.Help().Key+": "+m.keys.Select.Help().Desc,
	}
	if m.currentView == app.DashboardView || m.currentView == app.TaskManagementView {
		parts = append(parts, m.keys.CompleteTask.Help().Key+": "+m.keys.CompleteTask.Help().Desc)
	}
	if m.currentView != app.MainMenuView {
		parts = append(parts, m.keys.Back.Help().Key+": "+m.keys.Back.Help().Desc)
	}
	parts = append(parts, m.keys.Quit.Help().Key+": "+m.keys.Quit.Help().Desc)
	return helpStyle.Render(strings.Join(parts, " | "))
}

// List Item Types (placeholder, dashboardTaskItem, etc.)
type placeholderItem struct{ title, description string }
func (p placeholderItem) Title() string       { return p.title }
func (p placeholderItem) Description() string { return p.description }
func (p placeholderItem) FilterValue() string { return p.title }

type dashboardTaskItem struct{ models.Task }
func (dti dashboardTaskItem) Title() string { return dti.Task.Title }
func (dti dashboardTaskItem) Description() string {
	desc := ""; if dti.Task.DueDate != nil { desc = fmt.Sprintf("Prazo: %s", dti.Task.DueDate.Format("02/01/2006")) } else { desc = "Sem prazo" }
	if dti.Task.ClassID != nil { desc += fmt.Sprintf(" (Turma ID: %d)", *dti.Task.ClassID) }
	return desc
}
func (dti dashboardTaskItem) FilterValue() string { return dti.Task.Title }

type listItemClass struct{ models.Class }
func (lic listItemClass) Title() string { return lic.Name }
func (lic listItemClass) Description() string { return fmt.Sprintf("ID: %d, Disciplina ID: %d", lic.Class.ID, lic.SubjectID) }
func (lic listItemClass) FilterValue() string { return lic.Name }
func (lic listItemClass) ID() int64 { return lic.Class.ID }

type listItemStudent struct{ models.Student }
func (lis listItemStudent) Title() string { return lis.FullName }
func (lis listItemStudent) Description() string { return fmt.Sprintf("Matrícula: %s, Status: %s", lis.EnrollmentID, lis.Status) }
func (lis listItemStudent) FilterValue() string { return lis.FullName }
func (lis listItemStudent) ID() int64 { return lis.Student.ID }


// Messages
type errMsg struct{ err error; context string }
func (e errMsg) Error() string { return fmt.Sprintf("ctx: %s, err: %v", e.context, e.err) }

type mainMenuLoadedMsg []list.Item
type tasksForViewLoadedMsg []models.Task // Used by Dashboard and TaskManagementView
type classesLoadedMsg []models.Class
type studentsLoadedMsg []models.Student

// Start
func Start(classService service.ClassService, taskService service.TaskService, assessmentService service.AssessmentService) error {
	if classService == nil || taskService == nil || assessmentService == nil {
		return fmt.Errorf("serviços essenciais ausentes para iniciar a TUI")
	}
	m := NewTUIModel(classService, taskService, assessmentService)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Printf("Erro ao executar programa Bubble Tea: %v", err)
		return err
	}
	return nil
}
