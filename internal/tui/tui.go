// Package tui implements the Text User Interface (TUI) for the Vigenda application.
// It uses the Bubble Tea library and its components to create an interactive CLI experience.
package tui

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time" // Added for dashboard item descriptions
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
	docStyle        = lipgloss.NewStyle().Margin(1, 2)
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))  // Purple
	selectedStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("62")).PaddingLeft(1)
	deselectedStyle = lipgloss.NewStyle().PaddingLeft(1)
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray
)

// KeyMap defines the keybindings for the TUI.
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	Quit   key.Binding
	Help   key.Binding
}

// DefaultKeyMap provides the default keybindings.
var DefaultKeyMap = KeyMap{
	Up:     key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "move up")),
	Down:   key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "move down")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Back:   key.NewBinding(key.WithKeys("esc", "backspace"), key.WithHelp("esc/bksp", "back")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q/ctrl+c", "quit")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
}

// Model represents the state of the TUI.
type Model struct {
	// Services
	classService      service.ClassService
	taskService       service.TaskService
	assessmentService service.AssessmentService

	// TUI components
	list    list.Model
	spinner spinner.Model
	keys    KeyMap

	// State
	currentView app.View
	isLoading   bool
	err         error

	// Data for views
	classes           []models.Class
	students          []models.Student
	dashboardItems    []list.Item // Holds items for the dashboard list
	upcomingTasks     []models.Task
	recentAssessments []models.Assessment // Example, not yet loaded

	// selectedClass stores the currently selected class when navigating to students view
	selectedClass *models.Class
}

// NewTUIModel creates a new TUI model.
func NewTUIModel(cs service.ClassService, ts service.TaskService, as service.AssessmentService) Model {
	log.Printf("TUI: NewTUIModel - Chamado. ClassService: %t, TaskService: %t, AssessmentService: %t", cs != nil, ts != nil, as != nil)
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	delegate := list.NewDefaultDelegate()
	mainList := list.New(nil, delegate, 0, 0) // Items will be set in Update
	mainList.SetShowHelp(false)               // Using custom footer for help

	m := Model{
		classService:      cs,
		taskService:       ts,
		assessmentService: as,
		spinner:           s,
		keys:              DefaultKeyMap,
		currentView:       app.DashboardView, // Start with DashboardView
		isLoading:         true,              // Start in loading state for initial data
		list:              mainList,
	}
	log.Println("TUI: NewTUIModel - Modelo TUI inicializado.")
	return m
}

// loadDashboardData fetches data for the dashboard.
func (m *Model) loadDashboardData() tea.Cmd {
	m.isLoading = true
	log.Println("TUI: loadDashboardData - Iniciando carregamento de dados do dashboard.")
	return func() tea.Msg {
		log.Println("TUI: loadDashboardData (async) - Tentando carregar tarefas ativas para o dashboard.")
		tasks, err := m.taskService.ListAllActiveTasks(context.Background())
		if err != nil {
			log.Printf("TUI: loadDashboardData (async) - Erro ao carregar tarefas: %v", err)
			return errMsg{err: err, context: "loading dashboard tasks"}
		}
		log.Printf("TUI: loadDashboardData (async) - Tarefas do dashboard carregadas: %d tarefas.", len(tasks))
		return dashboardTasksLoadedMsg(tasks)
	}
}

func (m *Model) loadClasses() tea.Cmd {
	m.isLoading = true
	log.Println("TUI: loadClasses - Iniciando carregamento de turmas.")
	return func() tea.Msg {
		log.Println("TUI: loadClasses (async) - Tentando carregar turmas do serviço.")
		classes, err := m.classService.ListAllClasses(context.Background())
		if err != nil {
			log.Printf("TUI: loadClasses (async) - Erro ao carregar turmas: %v", err)
			return errMsg{err: err, context: "loading classes"}
		}
		log.Printf("TUI: loadClasses (async) - Turmas carregadas com sucesso: %d turmas.", len(classes))
		return classesLoadedMsg(classes)
	}
}

func (m *Model) loadStudentsForClass(classID int64) tea.Cmd {
	m.isLoading = true
	log.Printf("TUI: loadStudentsForClass - Iniciando carregamento de alunos para a turma ID %d.", classID)
	return func() tea.Msg {
		log.Printf("TUI: loadStudentsForClass (async) - Tentando carregar alunos para a turma ID %d.", classID)
		students, err := m.classService.GetStudentsByClassID(context.Background(), classID)
		if err != nil {
			log.Printf("TUI: loadStudentsForClass (async) - Erro ao carregar alunos: %v", err)
			return errMsg{err: err, context: fmt.Sprintf("loading students for class %d", classID)}
		}
		log.Printf("TUI: loadStudentsForClass (async) - Alunos carregados: %d.", len(students))
		return studentsLoadedMsg(students)
	}
}

// Init initializes the TUI model.
func (m Model) Init() tea.Cmd {
	var initialCmd tea.Cmd
	switch m.currentView {
	case app.DashboardView:
		initialCmd = m.loadDashboardData()
	case app.ClassManagementView:
		initialCmd = m.loadClasses()
	default:
		log.Printf("TUI: Init - Unknown initial view: %s. Defaulting to loadDashboardData.", m.currentView.String())
		initialCmd = m.loadDashboardData()
	}
	return tea.Batch(m.spinner.Tick, initialCmd)
}

// Update handles messages and updates the TUI model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		switch m.currentView {
		case app.DashboardView:
			switch {
			case key.Matches(msg, m.keys.Up):
				m.list.CursorUp()
			case key.Matches(msg, m.keys.Down):
				m.list.CursorDown()
			case key.Matches(msg, m.keys.Select):
				selectedListItem := m.list.SelectedItem()
				if selectedListItem != nil {
					log.Printf("Dashboard item selected: %s", selectedListItem.FilterValue())
					// Future: Navigate based on item type.
					// For example, if item is dashboardTaskItem, switch to TaskManagementView.
					// if taskItem, ok := selectedListItem.(dashboardTaskItem); ok {
					// m.currentView = app.TaskManagementView
					// Call a method to load data for taskItem.Task.ID
					// cmds = append(cmds, m.loadTaskDetailsCmd(taskItem.Task.ID))
					// }
				}
			case key.Matches(msg, m.keys.Back):
				return m, tea.Quit // On Dashboard, 'Back' quits.
			default:
				if !m.isLoading { // Pass unhandled keys to list
					m.list, cmd = m.list.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		case app.ClassManagementView:
			switch {
			case key.Matches(msg, m.keys.Up):
				m.list.CursorUp()
			case key.Matches(msg, m.keys.Down):
				m.list.CursorDown()
			case key.Matches(msg, m.keys.Select):
				selectedItem := m.list.SelectedItem()
				if class, ok := selectedItem.(listItemClass); ok {
					m.selectedClass = &class.Class
					m.currentView = app.View(99) // Temporary view for students
					m.list.Title = fmt.Sprintf("Alunos da Turma: %s", m.selectedClass.Name)
					m.list.SetItems(nil) // Clear current items
					m.students = nil     // Clear student data
					cmds = append(cmds, m.loadStudentsForClass(class.ID()))
				}
			case key.Matches(msg, m.keys.Back):
				m.currentView = app.DashboardView
				m.list.Title = "Dashboard"
				m.list.SetItems(nil)     // Clear current items
				m.dashboardItems = nil   // Clear dashboard data
				cmds = append(cmds, m.loadDashboardData())
			default:
				if !m.isLoading {
					m.list, cmd = m.list.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		case app.View(99): // Student view (temporary)
			switch {
			case key.Matches(msg, m.keys.Up):
				m.list.CursorUp()
			case key.Matches(msg, m.keys.Down):
				m.list.CursorDown()
			case key.Matches(msg, m.keys.Select):
				// Student selection logic here
			case key.Matches(msg, m.keys.Back):
				m.currentView = app.ClassManagementView
				m.list.Title = "Turmas"
				m.selectedClass = nil
				m.list.SetItems(nil) // Clear student items
				m.classes = nil      // Clear class data to force reload
				cmds = append(cmds, m.loadClasses())
			default:
				if !m.isLoading {
					m.list, cmd = m.list.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		default: // Default for other views primarily using the list
			if !m.isLoading {
				m.list, cmd = m.list.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case spinner.TickMsg:
		if m.isLoading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case dashboardTasksLoadedMsg:
		log.Println("TUI: Update - Recebida dashboardTasksLoadedMsg.")
		m.isLoading = false
		m.upcomingTasks = []models.Task(msg)

		m.dashboardItems = []list.Item{} // Reset before populating
		if len(m.upcomingTasks) == 0 {
			m.dashboardItems = append(m.dashboardItems, placeholderItem{title: "Nenhuma tarefa ativa encontrada."})
		} else {
			for _, task := range m.upcomingTasks {
				m.dashboardItems = append(m.dashboardItems, dashboardTaskItem{task})
			}
		}
		// If other dashboard data (e.g., assessments) is loaded, append it here.
		// Example:
		// for _, assessment := range m.recentAssessments {
		// m.dashboardItems = append(m.dashboardItems, dashboardAssessmentItem{assessment})
		// }
		m.list.SetItems(m.dashboardItems)
		m.list.Title = "Dashboard"
		log.Printf("TUI: Update - Dashboard atualizado com %d itens.", len(m.dashboardItems))

	case classesLoadedMsg:
		log.Println("TUI: Update - Recebida classesLoadedMsg.")
		m.isLoading = false
		m.classes = []models.Class(msg)
		items := make([]list.Item, len(m.classes))
		if len(m.classes) == 0 {
			items = append(items, placeholderItem{title: "Nenhuma turma encontrada."})
		} else {
			for i, c := range m.classes {
				items[i] = listItemClass{c}
			}
		}
		m.list.SetItems(items)
		m.list.Title = "Turmas"
		log.Printf("TUI: Update - Lista de turmas atualizada com %d itens.", len(items))

	case studentsLoadedMsg:
		log.Println("TUI: Update - Recebida studentsLoadedMsg.")
		m.isLoading = false
		m.students = []models.Student(msg)
		items := make([]list.Item, len(m.students))
		if len(m.students) == 0 {
			items = append(items, placeholderItem{title: "Nenhum aluno encontrado nesta turma."})
		} else {
			for i, s := range m.students {
				items[i] = listItemStudent{s}
			}
		}
		m.list.SetItems(items)
		// Title for student list is set when navigating to it.
		log.Printf("TUI: Update - Lista de alunos atualizada com %d itens.", len(items))

	case errMsg:
		log.Printf("TUI: Update - Recebida errMsg. Contexto: '%s', Erro: %v", msg.context, msg.err)
		m.isLoading = false
		m.err = msg.err // Store error to be displayed by View()
		// No automatic quit; View() will show the error.

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		listHeight := msg.Height - v - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView())
		m.list.SetSize(msg.Width-h, listHeight)
	}

	return m, tea.Batch(cmds...)
}

// View renders the TUI.
func (m Model) View() string {
	if m.err != nil {
		errorText := fmt.Sprintf("Ocorreu um erro:\n%v\n\nPressione 'q' ou 'ctrl+c' para sair.", m.err)
		return docStyle.Render(errorText)
	}

	if m.isLoading {
		loadingText := fmt.Sprintf("%s Carregando %s...", m.spinner.View(), m.currentView.String())
		return docStyle.Render(loadingText)
	}

	mainContent := m.list.View()
	return docStyle.Render(m.headerView() + "\n" + mainContent + "\n" + m.footerView())
}

func (m Model) headerView() string {
	titleStr := m.currentView.String()
	if m.currentView == app.View(99) { // Student view (temporary naming)
		if m.selectedClass != nil {
			titleStr = fmt.Sprintf("Alunos - %s", m.selectedClass.Name)
		} else {
			titleStr = "Alunos"
		}
	}
	return titleStyle.Render(titleStr)
}

func (m Model) footerView() string {
	var helpParts []string
	helpParts = append(helpParts, m.keys.Up.Help().Key+"/"+m.keys.Down.Help().Key+": "+m.keys.Up.Help().Desc)
	helpParts = append(helpParts, m.keys.Select.Help().Key+": "+m.keys.Select.Help().Desc)

	// Show 'Back' option only if not on the Dashboard
	if m.currentView != app.DashboardView {
		helpParts = append(helpParts, m.keys.Back.Help().Key+": "+m.keys.Back.Help().Desc)
	}
	helpParts = append(helpParts, m.keys.Quit.Help().Key+": "+m.keys.Quit.Help().Desc)
	// helpParts = append(helpParts, m.keys.Help.Help().Key+": "+m.keys.Help.Help().Desc) // If general help toggle is active

	return helpStyle.Render(strings.Join(helpParts, " | "))
}

// --- Custom list items ---

// placeholderItem is used when a list is empty.
type placeholderItem struct {
	title       string
	description string
}

func (p placeholderItem) Title() string       { return p.title }
func (p placeholderItem) Description() string { return p.description }
func (p placeholderItem) FilterValue() string { return p.title }

// dashboardTaskItem for displaying tasks on the dashboard.
type dashboardTaskItem struct {
	models.Task
}

func (dti dashboardTaskItem) Title() string { return dti.Task.Title }
func (dti dashboardTaskItem) Description() string {
	desc := ""
	if dti.Task.DueDate != nil {
		desc = fmt.Sprintf("Prazo: %s", dti.Task.DueDate.Format("02/01/2006"))
	} else {
		desc = "Sem prazo definido"
	}
	if dti.Task.ClassID != nil {
		// For a better UX, ClassName could be fetched and displayed.
		desc += fmt.Sprintf(" (Turma ID: %d)", *dti.Task.ClassID)
	}
	return desc
}
func (dti dashboardTaskItem) FilterValue() string { return dti.Task.Title }

// dashboardAssessmentItem for displaying assessments (example, not fully used yet).
type dashboardAssessmentItem struct {
	models.Assessment
}

func (dai dashboardAssessmentItem) Title() string { return dai.Assessment.Name }
func (dai dashboardAssessmentItem) Description() string {
	return fmt.Sprintf("Turma ID: %d, Período: %d, Peso: %.1f",
		dai.Assessment.ClassID, dai.Assessment.Term, dai.Assessment.Weight)
}
func (dai dashboardAssessmentItem) FilterValue() string { return dai.Assessment.Name }

// listItemClass for displaying classes.
type listItemClass struct {
	models.Class
}

func (lic listItemClass) Title() string       { return lic.Name }
func (lic listItemClass) Description() string { return fmt.Sprintf("ID: %d, Disciplina ID: %d", lic.Class.ID, lic.SubjectID) }
func (lic listItemClass) FilterValue() string { return lic.Name }
func (lic listItemClass) ID() int64           { return lic.Class.ID }

// listItemStudent for displaying students.
type listItemStudent struct {
	models.Student
}

func (lis listItemStudent) Title() string { return lis.FullName }
func (lis listItemStudent) Description() string {
	return fmt.Sprintf("Matrícula: %s, Status: %s", lis.EnrollmentID, lis.Status)
}
func (lis listItemStudent) FilterValue() string { return lis.FullName }
func (lis listItemStudent) ID() int64           { return lis.Student.ID }

// --- Messages for async operations ---

type errMsg struct {
	err     error
	context string
}

func (e errMsg) Error() string {
	return fmt.Sprintf("context: %s, error: %v", e.context, e.err)
}

type classesLoadedMsg []models.Class
type studentsLoadedMsg []models.Student
type dashboardTasksLoadedMsg []models.Task
// Example for future: type dashboardAssessmentsLoadedMsg []models.Assessment

// Start runs the TUI.
func Start(classService service.ClassService, taskService service.TaskService, assessmentService service.AssessmentService) error {
	log.Printf("TUI: Start - Função Start chamada. Services - Class: %t, Task: %t, Assessment: %t",
		classService != nil, taskService != nil, assessmentService != nil)

	if classService == nil || taskService == nil || assessmentService == nil {
		var missing []string
		if classService == nil { missing = append(missing, "ClassService") }
		if taskService == nil { missing = append(missing, "TaskService") }
		if assessmentService == nil { missing = append(missing, "AssessmentService") }

		fatalMsg := fmt.Sprintf("TUI: Start - Serviços essenciais ausentes: %v. A TUI não pode iniciar.", strings.Join(missing, ", "))
		log.Println(fatalMsg) // Log as info/error, don't os.Exit here
		return fmt.Errorf(fatalMsg) // Return an error to the caller
	}

	m := NewTUIModel(classService, taskService, assessmentService)
	// Consider tea.WithMouseCellMotion() for mouse support if needed later.
	p := tea.NewProgram(m, tea.WithAltScreen())

	log.Println("TUI: Start - Iniciando programa Bubble Tea (p.Run()).")
	finalModel, err := p.Run() // p.Run() is blocking
	if err != nil {
		log.Printf("TUI: Start - Erro ao executar o programa Bubble Tea: %v", err)
		return err
	}

	log.Println("TUI: Start - Programa Bubble Tea finalizado.")
	if finalModel != nil {
		// Access the final model state if needed, e.g., to check m.err
		// finalState, ok := finalModel.(Model)
		// if ok && finalState.err != nil {
		// log.Printf("TUI: Start - Final model state contained an error: %v", finalState.err)
		// return finalState.err
		// }
	}
	return nil // Success
}
