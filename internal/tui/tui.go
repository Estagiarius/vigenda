// Package tui implements the Text User Interface (TUI) for the Vigenda application.
// It uses the Bubble Tea library and its components to create an interactive CLI experience.
package tui

import (
	"context"
	"fmt"
	"log"
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
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")) // Purple
	selectedStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("62")).PaddingLeft(1)
	deselectedStyle = lipgloss.NewStyle().PaddingLeft(1)
	// helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray // Commented out to resolve re-declaration
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
	classService service.ClassService

	// TUI components
	list    list.Model
	spinner spinner.Model
	keys    KeyMap

	// State
	currentView app.View
	isLoading   bool
	err         error

	// Data for views
	classes  []models.Class
	students []models.Student
	// selectedClass stores the currently selected class when navigating to students view
	selectedClass *models.Class
}

// NewTUIModel creates a new TUI model.
func NewTUIModel(cs service.ClassService) Model {
	log.Printf("TUI: NewTUIModel - Chamado. ClassService is nil: %t", cs == nil)
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := Model{
		classService: cs,
		spinner:      s,
		keys:         DefaultKeyMap,
		currentView:  app.DashboardView, // Start with Dashboard or main menu
		isLoading:    false,
	}
	log.Println("TUI: NewTUIModel - Modelo TUI parcialmente inicializado, chamando loadInitialData.")
	m.loadInitialData() // Load initial list of items (e.g., main menu options or classes)
	log.Println("TUI: NewTUIModel - loadInitialData chamado, retornando modelo.")
	return m
}

func (m *Model) loadInitialData() {
	m.isLoading = true
	// Initially, we might load main menu items or directly classes if that's the primary view
	// For "Gerenciar Turmas e Alunos", we start by listing classes.
	m.currentView = app.ClassManagementView // Set to class management view for now
	m.loadClasses()
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
		log.Printf("TUI: loadStudentsForClass (async) - Tentando carregar alunos para a turma ID %d do serviço.", classID)
		students, err := m.classService.GetStudentsByClassID(context.Background(), classID)
		if err != nil {
			log.Printf("TUI: loadStudentsForClass (async) - Erro ao carregar alunos para a turma ID %d: %v", classID, err)
			return errMsg{err: err, context: fmt.Sprintf("loading students for class %d", classID)}
		}
		log.Printf("TUI: loadStudentsForClass (async) - Alunos carregados com sucesso para a turma ID %d: %d alunos.", classID, len(students))
		return studentsLoadedMsg(students)
	}
}

// Init initializes the TUI model.
func (m Model) Init() tea.Cmd {
	// return m.spinner.Tick // Start spinner if initially loading
	return tea.Batch(m.spinner.Tick, m.loadClasses())
}

// Update handles messages and updates the TUI model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			m.list.CursorUp()
		case key.Matches(msg, m.keys.Down):
			m.list.CursorDown()
		case key.Matches(msg, m.keys.Select):
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				switch m.currentView {
				case app.ClassManagementView:
					if class, ok := selectedItem.(listItemClass); ok {
						m.selectedClass = &class.Class
						m.currentView = app.View(99) // Temporary view for students
						cmds = append(cmds, m.loadStudentsForClass(class.ID()))
					}
				case app.View(99): // Student view (temporary)
					// Handle student selection if needed, or just go back
				}
			}
		case key.Matches(msg, m.keys.Back):
			if m.currentView == app.View(99) { // Student view
				m.currentView = app.ClassManagementView
				m.selectedClass = nil
				// Reload classes to reset the list
				cmds = append(cmds, m.loadClasses())
			} else if m.currentView == app.ClassManagementView {
				// TODO: Go back to a previous menu or exit if this is the top level
				// For now, quit if at class view and back is pressed.
				// return m, tea.Quit
			}
		}

	case spinner.TickMsg:
		if m.isLoading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case classesLoadedMsg:
		log.Println("TUI: Update - Recebida classesLoadedMsg.")
		m.isLoading = false
		m.classes = []models.Class(msg)
		items := make([]list.Item, len(m.classes))
		for i, c := range m.classes {
			items[i] = listItemClass{c}
		}
		m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
		m.list.Title = "Turmas"
		m.list.SetShowHelp(false) // We'll use our own help display
		log.Printf("TUI: Update - Lista de turmas atualizada com %d itens.", len(items))

	case studentsLoadedMsg:
		log.Println("TUI: Update - Recebida studentsLoadedMsg.")
		m.isLoading = false
		m.students = []models.Student(msg)
		items := make([]list.Item, len(m.students))
		for i, s := range m.students {
			items[i] = listItemStudent{s}
		}
		m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
		if m.selectedClass != nil {
			m.list.Title = fmt.Sprintf("Alunos da Turma: %s", m.selectedClass.Name)
		} else {
			m.list.Title = "Alunos"
		}
		m.list.SetShowHelp(false)
		log.Printf("TUI: Update - Lista de alunos atualizada com %d itens.", len(items))

	case errMsg:
		log.Printf("TUI: Update - Recebida errMsg. Contexto: '%s', Erro: %v", msg.context, msg.err)
		m.isLoading = false
		m.err = msg.err // Store the actual error
		// Log the error
		// Potentially return m, tea.Quit or display error to user
		return m, tea.Quit // For now, quit on error

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-lipgloss.Height(m.headerView())-lipgloss.Height(m.footerView()))
	}

	// Handle list updates if not loading and no specific key handled above
	if !m.isLoading {
		var listCmd tea.Cmd
		m.list, listCmd = m.list.Update(msg)
		cmds = append(cmds, listCmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the TUI.
func (m Model) View() string {
	if m.err != nil {
		// Simplified error view
		return docStyle.Render(fmt.Sprintf("Ocorreu um erro: %v\n\nPressione qualquer tecla para sair.", m.err))
	}

	if m.isLoading {
		loadingText := fmt.Sprintf("%s Carregando...", m.spinner.View())
		if m.currentView == app.ClassManagementView && len(m.classes) == 0 {
			loadingText += "\n\nNenhuma turma encontrada ainda."
		} else if m.currentView == app.View(99) && len(m.students) == 0 && m.selectedClass != nil {
			loadingText += fmt.Sprintf("\n\nNenhum aluno encontrado para a turma %s.", m.selectedClass.Name)
		}
		return docStyle.Render(loadingText)
	}

	// If not loading and no error, render the current view
	return docStyle.Render(m.headerView() + "\n" + m.list.View() + "\n" + m.footerView())
}

func (m Model) headerView() string {
	var title string
	switch m.currentView {
	case app.ClassManagementView:
		title = titleStyle.Render("Gerenciar Turmas")
	case app.View(99): // Student view
		if m.selectedClass != nil {
			title = titleStyle.Render(fmt.Sprintf("Alunos - %s", m.selectedClass.Name))
		} else {
			title = titleStyle.Render("Alunos")
		}
	default:
		title = titleStyle.Render(m.currentView.String())
	}
	return title
}

func (m Model) footerView() string {
	// Basic help. Could be more context-aware.
	return helpStyle.Render(fmt.Sprintf("↑/↓: Navegar | Enter: Selecionar | Esc: Voltar | q: Sair"))
}

// Custom list items
type listItemClass struct {
	models.Class
}

func (lic listItemClass) Title() string       { return lic.Name }
func (lic listItemClass) Description() string { return fmt.Sprintf("ID: %d, Disciplina ID: %d", lic.Class.ID, lic.SubjectID) } // Use lic.Class.ID
func (lic listItemClass) FilterValue() string { return lic.Name }
func (lic listItemClass) ID() int64           { return lic.Class.ID }

type listItemStudent struct {
	models.Student
}

func (lis listItemStudent) Title() string { return lis.FullName }
func (lis listItemStudent) Description() string {
	return fmt.Sprintf("Matrícula: %s, Status: %s", lis.EnrollmentID, lis.Status)
}
func (lis listItemStudent) FilterValue() string { return lis.FullName }
func (lis listItemStudent) ID() int64           { return lis.Student.ID }

// Messages for async operations
// errMsg now includes a context string to identify the source of the error.
type errMsg struct {
	err     error
	context string // e.g., "loading classes", "loading students for class X"
}

func (e errMsg) Error() string {
	return fmt.Sprintf("context: %s, error: %v", e.context, e.err)
}

type classesLoadedMsg []models.Class
type studentsLoadedMsg []models.Student

// Start runs the TUI.
func Start(classService service.ClassService) error {
	log.Printf("TUI: Start - Função Start chamada. ClassService is nil: %t", classService == nil)
	if classService == nil {
		// Usar log.Fatalf fará com que a aplicação encerre, o que é apropriado aqui.
		// O log irá para o arquivo de log configurado antes do encerramento.
		log.Fatalf("TUI: Start - ClassService não pode ser nulo para iniciar a TUI.")
	}
	m := NewTUIModel(classService)
	p := tea.NewProgram(m, tea.WithAltScreen()) // Use AltScreen for better TUI experience
	log.Println("TUI: Start - Iniciando programa Bubble Tea (p.Run()).")
	_, err := p.Run()
	if err != nil {
		log.Printf("TUI: Start - Erro ao executar o programa Bubble Tea: %v", err)
	} else {
		log.Println("TUI: Start - Programa Bubble Tea finalizado sem erros.")
	}
	return err
}
