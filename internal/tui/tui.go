// Package tui implements the Text User Interface (TUI) for the Vigenda application.
// It uses the Bubble Tea library and its components to create an interactive CLI experience.
package tui

import (
	"context"
	"fmt"
	"log"
	"strconv" // Added for parsing ClassID
	"strings"
	"time"
	"vigenda/internal/app"
	"vigenda/internal/models"
	"vigenda/internal/service"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
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
	focusedInputStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredInputStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	labelStyle         = lipgloss.NewStyle().Bold(true).MarginRight(1)
)

// KeyMap
type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Select       key.Binding // Also used for Submit in forms
	Back         key.Binding // Also used for Cancel in forms
	Quit         key.Binding
	Help         key.Binding
	CompleteTask key.Binding
	AddTask      key.Binding
	NextField    key.Binding
	PrevField    key.Binding
}

var DefaultKeyMap = KeyMap{
	Up:           key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "cima")),
	Down:         key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "baixo")),
	Select:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "slcnar/submeter")),
	Back:         key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "voltar/cancelar")), // Unified Back/Cancel
	Quit:         key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q/ctrl+c", "sair")),
	Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "ajuda")),
	CompleteTask: key.NewBinding(key.WithKeys("x", "c"), key.WithHelp("x/c", "concluir")),
	AddTask:      key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add tarefa")),
	NextField:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "próx. campo")),
	PrevField:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "campo ant.")),
}

// Constants for task form fields
const (
	taskFormFieldTitle       = "title"
	taskFormFieldDescription = "description"
	taskFormFieldDueDate     = "dueDate"
	taskFormFieldClassID     = "classID"
)

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
	taskFormInputs       map[string]textinput.Model
	taskFormFocusedField string
	taskFormOrder        []string
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

	formOrder := []string{
		taskFormFieldTitle,
		taskFormFieldDescription,
		taskFormFieldDueDate,
		taskFormFieldClassID,
	}

	m := Model{
		classService:      cs,
		taskService:       ts,
		assessmentService: as,
		spinner:           s,
		keys:              DefaultKeyMap,
		currentView:       app.MainMenuView,
		isLoading:         false,
		list:              mainList,
		taskFormInputs:    make(map[string]textinput.Model),
		taskFormOrder:     formOrder,
	}
	log.Println("TUI: NewTUIModel - currentView:", m.currentView.String())
	return m
}

func (m *Model) initializeTaskForm() {
	m.taskFormInputs = make(map[string]textinput.Model)
	var ti textinput.Model

	ti = textinput.New()
	ti.Placeholder = "Título da Tarefa (obrigatório)"
	ti.CharLimit = 120; ti.Width = 50
	ti.PromptStyle = focusedInputStyle
	m.taskFormInputs[taskFormFieldTitle] = ti

	ti = textinput.New()
	ti.Placeholder = "Descrição (opcional)"
	ti.CharLimit = 256; ti.Width = 50
	ti.PromptStyle = blurredInputStyle
	m.taskFormInputs[taskFormFieldDescription] = ti

	ti = textinput.New()
	ti.Placeholder = "AAAA-MM-DD (opcional)"
	ti.CharLimit = 10; ti.Width = 20
	ti.PromptStyle = blurredInputStyle
	m.taskFormInputs[taskFormFieldDueDate] = ti

	ti = textinput.New()
	ti.Placeholder = "ID da Turma (opcional, numérico)"
	ti.CharLimit = 10; ti.Width = 20
	ti.PromptStyle = blurredInputStyle
	m.taskFormInputs[taskFormFieldClassID] = ti

	if len(m.taskFormOrder) > 0 {
		m.taskFormFocusedField = m.taskFormOrder[0]
		focusedInput := m.taskFormInputs[m.taskFormFocusedField]
		focusedInput.Focus()
		m.taskFormInputs[m.taskFormFocusedField] = focusedInput
	}
    m.list.SetItems(nil)
    m.statusMessage = ""
    m.err = nil
}


const statusMessageDuration = 3 * time.Second
type clearStatusMessageMsg struct{}

func setStatusMessageCmd() tea.Cmd { // Removed duration, using const
	return tea.Tick(statusMessageDuration, func(t time.Time) tea.Msg {
		return clearStatusMessageMsg{}
	})
}

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

func (m *Model) loadAndDisplayTasks() tea.Cmd {
	m.isLoading = true
	return func() tea.Msg {
		tasks, err := m.taskService.ListAllActiveTasks(context.Background())
		if err != nil {
			return errMsg{err: err, context: "loading tasks for view"}
		}
		return tasksForViewLoadedMsg(tasks)
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

// Helper to change focus in form
func (m *Model) changeFormFocus(forward bool) {
    currentIndex := -1
    for i, fieldName := range m.taskFormOrder {
        if fieldName == m.taskFormFocusedField {
            currentIndex = i
            break
        }
    }

    // Blur current
    input := m.taskFormInputs[m.taskFormFocusedField]
    input.Blur()
	input.PromptStyle = blurredInputStyle
    m.taskFormInputs[m.taskFormFocusedField] = input

    if forward {
        currentIndex = (currentIndex + 1) % len(m.taskFormOrder)
    } else {
        currentIndex = (currentIndex - 1 + len(m.taskFormOrder)) % len(m.taskFormOrder)
    }

    // Focus next
    m.taskFormFocusedField = m.taskFormOrder[currentIndex]
    input = m.taskFormInputs[m.taskFormFocusedField]
    input.Focus()
	input.PromptStyle = focusedInputStyle
    m.taskFormInputs[m.taskFormFocusedField] = input
}


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Helper to navigate back to TaskManagementView from form
	goBackToTaskManagementView := func() tea.Model {
		m.currentView = app.TaskManagementView
		m.list.Title = app.TaskManagementView.String()
		m.taskFormInputs = make(map[string]textinput.Model) // Clear form
		m.taskFormFocusedField = ""
		m.isLoading = true // To show loading for task list
		cmds = append(cmds, m.loadAndDisplayTasks())
		return m
	}


	goBackToMainMenu := func() tea.Model {
		m.currentView = app.MainMenuView
		m.list.Title = app.MainMenuView.String()
		m.statusMessage = ""
		m.dashboardItems = nil; m.upcomingTasks = nil; m.classes = nil; m.students = nil; m.selectedClass = nil
		m.taskFormInputs = make(map[string]textinput.Model)
		m.taskFormFocusedField = ""
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
				cmds = append(cmds, m.loadAndDisplayTasks())
			}
			cmds = append(cmds, setStatusMessageCmd())
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
					case app.DashboardView:
						m.currentView = selected.targetView
						m.list.Title = selected.targetView.String()
						cmds = append(cmds, m.loadAndDisplayTasks())
					case app.TaskManagementView:
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

		case app.TaskManagementView:
			if key.Matches(msg, m.keys.AddTask) {
				m.initializeTaskForm()
				m.currentView = app.CreateTaskFormView
				m.list.Title = app.CreateTaskFormView.String()
				m.isLoading = false
			} else if key.Matches(msg, m.keys.CompleteTask) { handleTaskCompletion()
			} else if key.Matches(msg, m.keys.Back) { return goBackToMainMenu(), tea.Batch(cmds...)
			} else {
				m.list, cmd = m.list.Update(msg); cmds = append(cmds, cmd)
			}

		case app.DashboardView:
			if key.Matches(msg, m.keys.CompleteTask) { handleTaskCompletion() }
			else if key.Matches(msg, m.keys.Back) { return goBackToMainMenu(), tea.Batch(cmds...) }
			m.list, cmd = m.list.Update(msg); cmds = append(cmds, cmd)

		case app.CreateTaskFormView:
			if key.Matches(msg, m.keys.Back) { // Esc to cancel form
				return goBackToTaskManagementView(), tea.Batch(cmds...)
			} else if key.Matches(msg, m.keys.NextField) {
				m.changeFormFocus(true)
			} else if key.Matches(msg, m.keys.PrevField) {
				m.changeFormFocus(false)
			} else if key.Matches(msg, m.keys.Select) { // Enter to submit
				// Validation
				title := m.taskFormInputs[taskFormFieldTitle].Value()
				if title == "" {
					m.statusMessage = "Título da tarefa é obrigatório."
					cmds = append(cmds, setStatusMessageCmd())
					return m, tea.Batch(cmds...)
				}
				description := m.taskFormInputs[taskFormFieldDescription].Value()
				dueDateStr := m.taskFormInputs[taskFormFieldDueDate].Value()
				classIDStr := m.taskFormInputs[taskFormFieldClassID].Value()

				var dueDate *time.Time
				if dueDateStr != "" {
					parsedDate, err := time.Parse("2006-01-02", dueDateStr)
					if err != nil {
						m.statusMessage = fmt.Sprintf("Data de prazo inválida (use AAAA-MM-DD): %v", err)
						cmds = append(cmds, setStatusMessageCmd())
						return m, tea.Batch(cmds...)
					}
					dueDate = &parsedDate
				}

				var classID *int64
				if classIDStr != "" {
					cid, err := strconv.ParseInt(classIDStr, 10, 64)
					if err != nil {
						m.statusMessage = fmt.Sprintf("ID da Turma inválido (deve ser numérico): %v", err)
						cmds = append(cmds, setStatusMessageCmd())
						return m, tea.Batch(cmds...)
					}
					classID = &cid
				}

				// Call service
				_, err := m.taskService.CreateTask(context.Background(), title, description, classID, dueDate)
				if err != nil {
					m.statusMessage = fmt.Sprintf("Erro ao criar tarefa: %v", err)
					cmds = append(cmds, setStatusMessageCmd())
				} else {
					m.statusMessage = fmt.Sprintf("Tarefa '%s' criada com sucesso!", title)
					cmds = append(cmds, setStatusMessageCmd())
					// Return to task list view after successful creation
					return goBackToTaskManagementView(), tea.Batch(cmds...)
				}
			} else { // Default: pass key to focused input
				if ti, ok := m.taskFormInputs[m.taskFormFocusedField]; ok {
					var inputCmd tea.Cmd
					newTi, inputCmd := ti.Update(msg)
					m.taskFormInputs[m.taskFormFocusedField] = newTi
					cmds = append(cmds, inputCmd)
				}
			}


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

		case app.AssessmentManagementView, app.QuestionBankView, app.ProofGenerationView:
			if key.Matches(msg, m.keys.Back) { return goBackToMainMenu(), tea.Batch(cmds...) }
			m.list, cmd = m.list.Update(msg); cmds = append(cmds, cmd)
		}

	case spinner.TickMsg:
		if m.isLoading { m.spinner, cmd = m.spinner.Update(msg); cmds = append(cmds, cmd) }

	case clearStatusMessageMsg: m.statusMessage = ""

	case mainMenuLoadedMsg:
		m.isLoading = false; m.mainMenuItems = []list.Item(msg)
		m.list.SetItems(m.mainMenuItems); m.list.Title = app.MainMenuView.String()

	case tasksForViewLoadedMsg:
		m.isLoading = false; m.upcomingTasks = []models.Task(msg)
		m.dashboardItems = []list.Item{}
		if len(m.upcomingTasks) == 0 {
			m.dashboardItems = append(m.dashboardItems, placeholderItem{title: "Nenhuma tarefa ativa encontrada."})
		} else {
			for _, task := range m.upcomingTasks {
				m.dashboardItems = append(m.dashboardItems, dashboardTaskItem{task})
			}
		}
		m.list.SetItems(m.dashboardItems)

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

		inputWidth := msg.Width - h - 4 // Default width for inputs
		if m.currentView == app.CreateTaskFormView {
			for k := range m.taskFormInputs {
				ti := m.taskFormInputs[k]
				// Customize width per field if needed, e.g., description longer
				if k == taskFormFieldDescription {
					ti.Width = inputWidth // Or a larger fixed value
				} else if k == taskFormFieldDueDate || k == taskFormFieldClassID {
					ti.Width = inputWidth / 2 // Shorter fields
				} else {
					ti.Width = inputWidth
				}
				m.taskFormInputs[k] = ti
			}
		} else { // For list-based views
			listHeight := msg.Height - v - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView()) - statusMessageHeight
			m.list.SetSize(msg.Width-h, listHeight)
		}
	}
	return m, tea.Batch(cmds...)
}

// View
func (m Model) View() string {
	if m.err != nil { return docStyle.Render(fmt.Sprintf("Erro: %v\n\nSair (q/ctrl+c).", m.err)) }

	var sb strings.Builder
	sb.WriteString(m.headerView() + "\n")
	if m.statusMessage != "" { sb.WriteString(statusMessageStyle.Render(m.statusMessage) + "\n") }

	if m.isLoading && m.currentView != app.CreateTaskFormView { // Don't show general loading for form view
		title := m.list.Title; if title == "" { title = m.currentView.String() }
		sb.WriteString(fmt.Sprintf("\n%s Carregando %s...", m.spinner.View(), title))
	} else {
		if m.currentView == app.DashboardView {
			sb.WriteString(titleStyle.Copy().Bold(false).Underline(true).Render("Tarefas Pendentes") + "\n")
		}

		if m.currentView == app.CreateTaskFormView {
			formContent := []string{}
			for _, fieldName := range m.taskFormOrder {
				inputModel := m.taskFormInputs[fieldName]
				label := ""
				switch fieldName {
				case taskFormFieldTitle: label = "Título:"
				case taskFormFieldDescription: label = "Descrição:"
				case taskFormFieldDueDate: label = "Prazo (AAAA-MM-DD):"
				case taskFormFieldClassID: label = "ID da Turma:"
				}
				// Apply focus/blur styles to prompt, not placeholder
				if inputModel.Focused() {
					inputModel.PromptStyle = focusedInputStyle
				} else {
					inputModel.PromptStyle = blurredInputStyle
				}
				formContent = append(formContent, labelStyle.Render(label))
				formContent = append(formContent, inputModel.View()+"\n")
			}
			sb.WriteString(lipgloss.JoinVertical(lipgloss.Left, formContent...))
		} else {
		    sb.WriteString(m.list.View())
		}
	}
	sb.WriteString("\n" + m.footerView())
	return docStyle.Render(sb.String())
}

func (m Model) headerView() string {
	titleStr := ""
	if m.currentView == app.CreateTaskFormView {
		titleStr = app.CreateTaskFormView.String()
	} else {
		titleStr = m.list.Title
		if titleStr == "" { titleStr = m.currentView.String() }
		if m.currentView == app.StudentListView && m.selectedClass != nil {
			 titleStr = fmt.Sprintf("%s - %s", app.StudentListView.String(), m.selectedClass.Name)
		}
	}
	return titleStyle.Render(titleStr)
}

func (m Model) footerView() string {
	parts := []string{}
	if m.currentView == app.CreateTaskFormView {
		parts = append(parts, m.keys.NextField.Help().Key+": "+m.keys.NextField.Help().Desc)
		parts = append(parts, m.keys.PrevField.Help().Key+": "+m.keys.PrevField.Help().Desc)
		parts = append(parts, m.keys.Select.Help().Key+": submeter")
	} else { // For list based views
		parts = append(parts, m.keys.Up.Help().Key+"/"+m.keys.Down.Help().Key+": "+m.keys.Up.Help().Desc)
		parts = append(parts, m.keys.Select.Help().Key+": "+m.keys.Select.Help().Desc)
	}

	if m.currentView == app.DashboardView || m.currentView == app.TaskManagementView {
		parts = append(parts, m.keys.CompleteTask.Help().Key+": "+m.keys.CompleteTask.Help().Desc)
	}
	if m.currentView == app.TaskManagementView && m.currentView != app.CreateTaskFormView {
		parts = append(parts, m.keys.AddTask.Help().Key+": "+m.keys.AddTask.Help().Desc)
	}

	if m.currentView != app.MainMenuView {
		parts = append(parts, m.keys.Back.Help().Key+": "+m.keys.Back.Help().Desc)
	}
	parts = append(parts, m.keys.Quit.Help().Key+": "+m.keys.Quit.Help().Desc)
	return helpStyle.Render(strings.Join(parts, " | "))
}

// List Item Types
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
type tasksForViewLoadedMsg []models.Task
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
