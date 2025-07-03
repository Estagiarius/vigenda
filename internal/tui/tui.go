// Package tui implements the Text User Interface (TUI) for the Vigenda application.
// It uses the Bubble Tea library and its components to create an interactive CLI experience.
package tui

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time" // Added for status message timeout
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
	statusMessageStyle = lipgloss.NewStyle().MarginTop(1).MarginBottom(1).Foreground(lipgloss.Color("202")) // Orange, with margin
)

// KeyMap defines the keybindings for the TUI.
type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Select       key.Binding
	Back         key.Binding
	Quit         key.Binding
	Help         key.Binding
	CompleteTask key.Binding
}

// DefaultKeyMap provides the default keybindings.
var DefaultKeyMap = KeyMap{
	Up:           key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "navegar para cima")),
	Down:         key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "navegar para baixo")),
	Select:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "selecionar")),
	Back:         key.NewBinding(key.WithKeys("esc", "backspace"), key.WithHelp("esc/bksp", "voltar")),
	Quit:         key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q/ctrl+c", "sair")),
	Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "ajuda")),
	CompleteTask: key.NewBinding(key.WithKeys("x", "c"), key.WithHelp("x/c", "concluir tarefa")),
}

// Model represents the state of the TUI.
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

type mainMenuItem struct {
	title       string
	targetView  app.View
	description string
}

func (mmi mainMenuItem) Title() string       { return mmi.title }
func (mmi mainMenuItem) Description() string { return mmi.description }
func (mmi mainMenuItem) FilterValue() string { return mmi.title }

func NewTUIModel(cs service.ClassService, ts service.TaskService, as service.AssessmentService) Model {
	log.Printf("TUI: NewTUIModel - Chamado. Services Initialized - Class: %t, Task: %t, Assessment: %t", cs != nil, ts != nil, as != nil)
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
	log.Println("TUI: NewTUIModel - Modelo TUI inicializado, currentView:", m.currentView.String())
	return m
}

const statusMessageDuration = 3 * time.Second

type clearStatusMessageMsg struct{}

func setStatusMessageCmd(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return clearStatusMessageMsg{}
	})
}

func (m *Model) loadMainMenuItemsCmd() tea.Cmd {
	m.isLoading = true
	return func() tea.Msg {
		items := []list.Item{
			mainMenuItem{title: app.DashboardView.String(), targetView: app.DashboardView, description: "Visão geral e tarefas."},
			mainMenuItem{title: app.ClassManagementView.String(), targetView: app.ClassManagementView, description: "Gerenciar turmas e alunos."},
			mainMenuItem{title: app.AssessmentManagementView.String(), targetView: app.AssessmentManagementView, description: "Gerenciar avaliações e notas."},
			mainMenuItem{title: app.QuestionBankView.String(), targetView: app.QuestionBankView, description: "Acessar banco de questões."},
			mainMenuItem{title: app.ProofGenerationView.String(), targetView: app.ProofGenerationView, description: "Gerar provas."},
		}
		log.Printf("TUI: loadMainMenuItemsCmd - Itens do menu principal preparados: %d", len(items))
		return mainMenuLoadedMsg(items)
	}
}

func (m *Model) loadDashboardData() tea.Cmd {
	m.isLoading = true
	log.Println("TUI: loadDashboardData - Iniciando carregamento de dados do dashboard.")
	return func() tea.Msg {
		log.Println("TUI: loadDashboardData (async) - Tentando carregar tarefas ativas.")
		tasks, err := m.taskService.ListAllActiveTasks(context.Background())
		if err != nil {
			log.Printf("TUI: loadDashboardData (async) - Erro: %v", err)
			return errMsg{err: err, context: "loading dashboard tasks"}
		}
		log.Printf("TUI: loadDashboardData (async) - Tarefas carregadas: %d.", len(tasks))
		return dashboardTasksLoadedMsg(tasks)
	}
}

func (m *Model) loadClasses() tea.Cmd {
	m.isLoading = true
	log.Println("TUI: loadClasses - Iniciando carregamento de turmas.")
	return func() tea.Msg {
		log.Println("TUI: loadClasses (async) - Tentando carregar turmas.")
		classes, err := m.classService.ListAllClasses(context.Background())
		if err != nil {
			log.Printf("TUI: loadClasses (async) - Erro: %v", err)
			return errMsg{err: err, context: "loading classes"}
		}
		log.Printf("TUI: loadClasses (async) - Turmas carregadas: %d.", len(classes))
		return classesLoadedMsg(classes)
	}
}

func (m *Model) loadStudentsForClass(classID int64) tea.Cmd {
	m.isLoading = true
	log.Printf("TUI: loadStudentsForClass - Carregando alunos para turma ID %d.", classID)
	return func() tea.Msg {
		log.Printf("TUI: loadStudentsForClass (async) - Tentando carregar alunos para turma %d.", classID)
		students, err := m.classService.GetStudentsByClassID(context.Background(), classID)
		if err != nil {
			log.Printf("TUI: loadStudentsForClass (async) - Erro: %v", err)
			return errMsg{err: err, context: fmt.Sprintf("loading students for class %d", classID)}
		}
		log.Printf("TUI: loadStudentsForClass (async) - Alunos carregados: %d.", len(students))
		return studentsLoadedMsg(students)
	}
}

func (m Model) Init() tea.Cmd {
	log.Printf("TUI: Init - Current View: %s", m.currentView.String())
	var initialCmd tea.Cmd
	switch m.currentView {
	case app.MainMenuView:
		initialCmd = m.loadMainMenuItemsCmd()
	default:
		log.Printf("TUI: Init - Visão inicial não é MainMenu (%s). Carregando menu principal.", m.currentView.String())
		initialCmd = m.loadMainMenuItemsCmd()
	}
	return tea.Batch(m.spinner.Tick, initialCmd)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	goBackToMainMenu := func() tea.Model {
		m.currentView = app.MainMenuView
		m.list.Title = app.MainMenuView.String()
		m.statusMessage = ""
		m.dashboardItems = nil
		m.upcomingTasks = nil
		m.classes = nil
		m.students = nil
		m.selectedClass = nil
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		switch m.currentView {
		case app.MainMenuView:
			switch {
			case key.Matches(msg, m.keys.Up):
				m.list.CursorUp()
			case key.Matches(msg, m.keys.Down):
				m.list.CursorDown()
			case key.Matches(msg, m.keys.Select):
				selected, ok := m.list.SelectedItem().(mainMenuItem)
				if ok {
					log.Printf("TUI: MainMenu - Selecionado: %s, Navegando para: %s", selected.title, selected.targetView.String())
					m.list.SetItems(nil)
					m.isLoading = true
					m.statusMessage = ""

					switch selected.targetView {
					case app.DashboardView:
						m.currentView = app.DashboardView
						m.list.Title = app.DashboardView.String()
						cmds = append(cmds, m.loadDashboardData())
					case app.ClassManagementView:
						m.currentView = app.ClassManagementView
						m.list.Title = app.ClassManagementView.String()
						cmds = append(cmds, m.loadClasses())
					case app.AssessmentManagementView, app.QuestionBankView, app.ProofGenerationView:
						setPlaceholderView(selected.targetView)
					default:
						log.Printf("TUI: MainMenu - Comportamento de navegação para %s não definido.", selected.targetView.String())
						setPlaceholderView(selected.targetView)
					}
				}
			case key.Matches(msg, m.keys.Back):
				return m, tea.Quit
			default:
				if !m.isLoading {
					m.list, cmd = m.list.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		case app.DashboardView:
			switch {
			case key.Matches(msg, m.keys.CompleteTask):
				selectedItem := m.list.SelectedItem()
				if taskItem, ok := selectedItem.(dashboardTaskItem); ok {
					log.Printf("TUI: Dashboard - Tentando concluir tarefa ID: %d", taskItem.Task.ID)
					err := m.taskService.MarkTaskAsCompleted(context.Background(), taskItem.Task.ID)
					if err != nil {
						m.statusMessage = fmt.Sprintf("Erro ao concluir '%s': %v", taskItem.Title(), err)
						log.Printf("TUI: Dashboard - Erro ao concluir tarefa: %v", err)
					} else {
						m.statusMessage = fmt.Sprintf("Tarefa '%s' marcada como concluída!", taskItem.Title())
						log.Printf("TUI: Dashboard - Tarefa '%s' concluída.", taskItem.Title())
						m.isLoading = true
						cmds = append(cmds, m.loadDashboardData())
					}
					cmds = append(cmds, setStatusMessageCmd(statusMessageDuration))
				}
			case key.Matches(msg, m.keys.Back):
				return goBackToMainMenu(), tea.Batch(cmds...)
			default:
				if !m.isLoading {
					m.list, cmd = m.list.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		case app.ClassManagementView:
			switch {
			case key.Matches(msg, m.keys.Select):
				selected, ok := m.list.SelectedItem().(listItemClass)
				if ok {
					m.selectedClass = &selected.Class
					m.currentView = app.StudentListView
					m.list.Title = fmt.Sprintf("%s - %s", app.StudentListView.String(), m.selectedClass.Name)
					m.list.SetItems(nil)
					m.students = nil
					m.isLoading = true
					cmds = append(cmds, m.loadStudentsForClass(selected.ID()))
				}
			case key.Matches(msg, m.keys.Back):
				return goBackToMainMenu(), tea.Batch(cmds...)
			default:
				if !m.isLoading {
					m.list, cmd = m.list.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		case app.StudentListView:
			switch {
			case key.Matches(msg, m.keys.Back):
				m.currentView = app.ClassManagementView
				m.list.Title = app.ClassManagementView.String()
				m.selectedClass = nil
				m.list.SetItems(nil)
				m.students = nil
				m.isLoading = true
				cmds = append(cmds, m.loadClasses())
			default:
				if !m.isLoading {
					m.list, cmd = m.list.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		case app.AssessmentManagementView, app.QuestionBankView, app.ProofGenerationView:
			switch {
			case key.Matches(msg, m.keys.Back):
				return goBackToMainMenu(), tea.Batch(cmds...)
			default:
				if !m.isLoading {
					m.list, cmd = m.list.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		default:
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

	case clearStatusMessageMsg:
		m.statusMessage = ""
		log.Println("TUI: Update - Mensagem de status limpa.")

	case mainMenuLoadedMsg:
		log.Println("TUI: Update - Recebida mainMenuLoadedMsg.")
		m.isLoading = false
		m.mainMenuItems = []list.Item(msg)
		m.list.SetItems(m.mainMenuItems)
		m.list.Title = app.MainMenuView.String()
		log.Printf("TUI: Update - Menu Principal atualizado com %d itens.", len(m.mainMenuItems))

	case dashboardTasksLoadedMsg:
		log.Println("TUI: Update - Recebida dashboardTasksLoadedMsg.")
		m.isLoading = false
		m.upcomingTasks = []models.Task(msg)
		m.dashboardItems = []list.Item{}
		if len(m.upcomingTasks) == 0 {
			m.dashboardItems = append(m.dashboardItems, placeholderItem{title: "Nenhuma tarefa ativa encontrada."})
		} else {
			for _, task := range m.upcomingTasks {
				m.dashboardItems = append(m.dashboardItems, dashboardTaskItem{task})
			}
		}
		m.list.SetItems(m.dashboardItems)
		log.Printf("TUI: Update - Dashboard atualizado com %d tarefas.", len(m.upcomingTasks))

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
		log.Printf("TUI: Update - Lista de alunos atualizada com %d itens.", len(items))

	case errMsg:
		log.Printf("TUI: Update - Recebida errMsg. Contexto: '%s', Erro: %v", msg.context, msg.err)
		m.isLoading = false
		m.err = msg.err

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		// Adjust for header, footer, and potential status message line
		statusMessageHeight := 0
		if m.statusMessage != "" {
			statusMessageHeight = lipgloss.Height(statusMessageStyle.Render(m.statusMessage)) +1 // +1 for potential newline/spacing
		}
		listHeight := msg.Height - v - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView()) - statusMessageHeight
		m.list.SetSize(msg.Width-h, listHeight)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.err != nil {
		errorText := fmt.Sprintf("Ocorreu um erro:\n%v\n\nPressione 'q' ou 'ctrl+c' para sair.", m.err)
		return docStyle.Render(errorText)
	}

	var viewContentBuilder strings.Builder
	viewContentBuilder.WriteString(m.headerView() + "\n")

	if m.statusMessage != "" {
		viewContentBuilder.WriteString(statusMessageStyle.Render(m.statusMessage) + "\n")
	}

	if m.isLoading {
		titleForLoading := m.list.Title
		if titleForLoading == "" {
			titleForLoading = m.currentView.String()
		}
		loadingText := fmt.Sprintf("\n%s Carregando %s...", m.spinner.View(), titleForLoading)
		viewContentBuilder.WriteString(loadingText)
	} else {
		if m.currentView == app.DashboardView && len(m.dashboardItems) > 0 { // Check if dashboardItems is not empty
			// Only add section title if there are items to show, or always show it
			// For now, let's always show it if it's the dashboard view and not loading.
			// The list itself will show "Nenhuma tarefa..." if empty.
			viewContentBuilder.WriteString(titleStyle.Copy().Bold(false).Underline(true).Render("Tarefas Pendentes") + "\n")
		}
		viewContentBuilder.WriteString(m.list.View())
	}

	viewContentBuilder.WriteString("\n" + m.footerView())
	return docStyle.Render(viewContentBuilder.String())
}

func (m Model) headerView() string {
	titleStr := m.list.Title
	if titleStr == "" {
		titleStr = m.currentView.String()
	}
	if m.currentView == app.StudentListView && m.selectedClass != nil {
		 titleStr = fmt.Sprintf("%s - %s", app.StudentListView.String(), m.selectedClass.Name)
	}
	return titleStyle.Render(titleStr)
}

func (m Model) footerView() string {
	var helpParts []string
	helpParts = append(helpParts, m.keys.Up.Help().Key+"/"+m.keys.Down.Help().Key+": "+m.keys.Up.Help().Desc)
	helpParts = append(helpParts, m.keys.Select.Help().Key+": "+m.keys.Select.Help().Desc)

	if m.currentView == app.DashboardView {
		helpParts = append(helpParts, m.keys.CompleteTask.Help().Key+": "+m.keys.CompleteTask.Help().Desc)
	}

	if m.currentView != app.MainMenuView {
		helpParts = append(helpParts, m.keys.Back.Help().Key+": "+m.keys.Back.Help().Desc)
	}
	helpParts = append(helpParts, m.keys.Quit.Help().Key+": "+m.keys.Quit.Help().Desc)

	return helpStyle.Render(strings.Join(helpParts, " | "))
}

type placeholderItem struct {
	title       string
	description string
}

func (p placeholderItem) Title() string       { return p.title }
func (p placeholderItem) Description() string { return p.description }
func (p placeholderItem) FilterValue() string { return p.title }

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
		desc += fmt.Sprintf(" (Turma ID: %d)", *dti.Task.ClassID)
	}
	return desc
}
func (dti dashboardTaskItem) FilterValue() string { return dti.Task.Title }

type dashboardAssessmentItem struct {
	models.Assessment
}

func (dai dashboardAssessmentItem) Title() string { return dai.Assessment.Name }
func (dai dashboardAssessmentItem) Description() string {
	return fmt.Sprintf("Turma ID: %d, Período: %d, Peso: %.1f",
		dai.Assessment.ClassID, dai.Assessment.Term, dai.Assessment.Weight)
}
func (dai dashboardAssessmentItem) FilterValue() string { return dai.Assessment.Name }

type listItemClass struct {
	models.Class
}

func (lic listItemClass) Title() string       { return lic.Name }
func (lic listItemClass) Description() string { return fmt.Sprintf("ID: %d, Disciplina ID: %d", lic.Class.ID, lic.SubjectID) }
func (lic listItemClass) FilterValue() string { return lic.Name }
func (lic listItemClass) ID() int64           { return lic.Class.ID }

type listItemStudent struct {
	models.Student
}

func (lis listItemStudent) Title() string { return lis.FullName }
func (lis listItemStudent) Description() string {
	return fmt.Sprintf("Matrícula: %s, Status: %s", lis.EnrollmentID, lis.Status)
}
func (lis listItemStudent) ID() int64           { return lis.Student.ID }
func (lis listItemStudent) FilterValue() string { return lis.FullName }

type errMsg struct {
	err     error
	context string
}

func (e errMsg) Error() string { return fmt.Sprintf("context: %s, error: %v", e.context, e.err) }

type mainMenuLoadedMsg []list.Item
type classesLoadedMsg []models.Class
type studentsLoadedMsg []models.Student
type dashboardTasksLoadedMsg []models.Task

func Start(classService service.ClassService, taskService service.TaskService, assessmentService service.AssessmentService) error {
	log.Printf("TUI: Start - Função Start chamada. Services - Class: %t, Task: %t, Assessment: %t",
		classService != nil, taskService != nil, assessmentService != nil)

	if classService == nil || taskService == nil || assessmentService == nil {
		var missing []string
		if classService == nil { missing = append(missing, "ClassService") }
		if taskService == nil { missing = append(missing, "TaskService") }
		if assessmentService == nil { missing = append(missing, "AssessmentService") }
		fatalMsg := fmt.Sprintf("TUI: Start - Serviços essenciais ausentes: %v. A TUI não pode iniciar.", strings.Join(missing, ", "))
		log.Println(fatalMsg)
		return fmt.Errorf(fatalMsg)
	}

	m := NewTUIModel(classService, taskService, assessmentService)
	p := tea.NewProgram(m, tea.WithAltScreen())

	log.Println("TUI: Start - Iniciando programa Bubble Tea (p.Run()).")
	_, err := p.Run()
	if err != nil {
		log.Printf("TUI: Start - Erro ao executar o programa Bubble Tea: %v", err)
		return err
	}

	log.Println("TUI: Start - Programa Bubble Tea finalizado.")
	return nil
}
