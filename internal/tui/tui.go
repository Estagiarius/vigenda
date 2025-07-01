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
	width       int // Largura da janela do terminal
	height      int // Altura da janela do terminal

	// Estado específico da DashboardView
	dashboardMainMenuItems  []string // Ex: {"Dashboard", "Tarefas", "Turmas", ...}
	dashboardMainMenuCursor int      // Item focado no menu lateral

	// Data for views
	classes       []models.Class
	students      []models.Student
	selectedClass *models.Class // selectedClass stores the currently selected class when navigating to students view
}

// NewTUIModel creates a new TUI model.
func NewTUIModel(cs service.ClassService) Model {
	log.Println("TUI: NewTUIModel - INÍCIO")
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	mainMenuItems := []string{
		"Dashboard",
		"Gerenciar Tarefas",
		"Gerenciar Turmas e Alunos",
		"Gerenciar Avaliações e Notas",
		"Banco de Questões",
		"Gerar Provas",
	}

	m := Model{
		classService:            cs,
		spinner:                 s,
		keys:                    DefaultKeyMap,
		currentView:             app.DashboardView, // Start with DashboardView
		isLoading:               false,
		dashboardMainMenuItems:  mainMenuItems,
		dashboardMainMenuCursor: 0,
	}
	log.Printf("TUI: NewTUIModel - Modelo TUI antes de loadInitialData. currentView: %s", m.currentView.String())
	m.loadInitialData() // Load initial data (primarily for Dashboard now)
	log.Printf("TUI: NewTUIModel - Modelo TUI após loadInitialData. currentView: %s", m.currentView.String())
	log.Println("TUI: NewTUIModel - FIM")
	return m
}

func (m *Model) loadInitialData() {
	log.Printf("TUI: loadInitialData - INÍCIO. currentView: %s", m.currentView.String())
	// This function will be responsible for loading initial data for the Dashboard.
	// For now, it can be a no-op or load some mock data for dashboard panels if we add them.
	// If DashboardView needs to load something asynchronously, this is where it would start.
	// m.isLoading = true
	// cmds = append(cmds, m.loadDashboardDataCmd())
	log.Printf("TUI: loadInitialData - FIM. currentView: %s. Nenhuma ação de carregamento específica para Dashboard.", m.currentView.String())
}

// loadDashboardDataCmd would be a tea.Cmd to load dashboard specific data.
// func (m *Model) loadDashboardDataCmd() tea.Cmd {
// 	 return func() tea.Msg {
// 	 	 // Simulate loading data
// 	 	 time.Sleep(time.Millisecond * 500)
// 	 	 return dashboardDataLoadedMsg{ /* ... some data ... */ }
// 	 }
// }

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
	log.Printf("TUI: Init - INÍCIO. currentView: %s", m.currentView.String())
	// return m.spinner.Tick // Start spinner if initially loading
	// No initial data loading command for dashboard yet, can be added later.
	// If dashboard needs async data: return tea.Batch(m.spinner.Tick, m.loadDashboardDataCmd())
	log.Printf("TUI: Init - FIM. currentView: %s", m.currentView.String())
	return m.spinner.Tick
}

// Helper function to check if the current view uses the m.list component
func (m *Model) currentViewUsesList() bool {
	switch m.currentView {
	case app.ClassManagementView, app.View(99), // app.View(99) is student list
		app.TaskManagementView, app.AssessmentManagementView,
		app.QuestionBankView, app.ProofGenerationView:
		return true
	default:
		return false
	}
}

// Update handles messages and updates the TUI model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	log.Printf("TUI: Update - INÍCIO. currentView: %s. Mensagem: %T", m.currentView.String(), msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // Always handle window size changes
		log.Printf("TUI: Update - Recebido tea.WindowSizeMsg. Width: %d, Height: %d", msg.Width, msg.Height)
		m.width = msg.Width
		m.height = msg.Height
		// If we are in a list-based view, update the list size
		if m.currentViewUsesList() {
			h, v := docStyle.GetFrameSize()
			headerHeight := lipgloss.Height(m.headerView())
			footerHeight := lipgloss.Height(m.footerView())
			m.list.SetSize(m.width-h, m.height-v-headerHeight-footerHeight)
		}

	case tea.KeyMsg:
		if m.currentView == app.DashboardView {
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Up):
				if m.dashboardMainMenuCursor > 0 {
					m.dashboardMainMenuCursor--
				}
			case key.Matches(msg, m.keys.Down):
				if m.dashboardMainMenuCursor < len(m.dashboardMainMenuItems)-1 {
					m.dashboardMainMenuCursor++
				}
			case key.Matches(msg, m.keys.Select): // Enter
				selectedMenuItem := m.dashboardMainMenuItems[m.dashboardMainMenuCursor]
				log.Printf("Dashboard: Item selecionado no menu: %s", selectedMenuItem)
				switch selectedMenuItem {
				case "Dashboard":
					// Already on dashboard, maybe refresh data in the future
					// cmds = append(cmds, m.loadDashboardDataCmd())
					log.Println("Dashboard: 'Dashboard' selecionado, nenhuma mudança de view.")
				case "Gerenciar Tarefas":
					m.currentView = app.TaskManagementView
					m.list.Title = "Gerenciar Tarefas" // Setup list for the new view
					// TODO: cmds = append(cmds, m.loadTaskViewDataCmd())
					log.Println("Dashboard: Navegando para Gerenciar Tarefas.")
				case "Gerenciar Turmas e Alunos":
					m.currentView = app.ClassManagementView
					m.list.Title = "Turmas" // Setup list for the new view
					cmds = append(cmds, m.loadClasses()) // This will set isLoading = true
					log.Println("Dashboard: Navegando para Gerenciar Turmas e Alunos.")
				case "Gerenciar Avaliações e Notas":
					m.currentView = app.AssessmentManagementView
					m.list.Title = "Gerenciar Avaliações e Notas" // Setup list
					// TODO: Load data for this view
					log.Println("Dashboard: Navegando para Gerenciar Avaliações e Notas.")
				case "Banco de Questões":
					m.currentView = app.QuestionBankView
					m.list.Title = "Banco de Questões" // Setup list
					// TODO: Load data for this view
					log.Println("Dashboard: Navegando para Banco de Questões.")
				case "Gerar Provas":
					m.currentView = app.ProofGenerationView
					m.list.Title = "Gerar Provas" // Setup list
					// TODO: Load data for this view
					log.Println("Dashboard: Navegando para Gerar Provas.")
				}
			// TODO: Add keybindings for Tab to navigate between dashboard panels if implemented
			}
		} else { // Logic for non-Dashboard views (those that might use m.list)
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Back):
				if m.currentView == app.View(99) { // Student list view
					m.currentView = app.ClassManagementView
					m.selectedClass = nil
					m.list.Title = "Turmas" // Reset title
					cmds = append(cmds, m.loadClasses())
					log.Printf("TUI: Voltando para ClassManagementView a partir da lista de alunos.")
				} else if m.currentViewUsesList() { // For other top-level list views
					m.currentView = app.DashboardView
					log.Printf("TUI: Voltando para DashboardView a partir de %s.", m.currentView.String())
					// No need to explicitly load dashboard data unless it's dynamic and needs refresh on back.
				}
				// If no specific back logic (e.g., already on Dashboard or a view with no back path), do nothing.
			default:
				// If the current view uses the list, pass the message to the list component.
				// This handles list's internal keybindings (like up/down navigation, selection).
				if m.currentViewUsesList() && !m.isLoading {
					var listCmd tea.Cmd
					m.list, listCmd = m.list.Update(msg) // This includes Select for list items
					cmds = append(cmds, listCmd)

					// Handle item selection from the list (if Enter was pressed and list handled it)
					// Need to check if the msg was 'enter' and if an item was actually selected by the list.
					// This is a bit tricky as list.Update() consumes the Enter key.
					// We might need to check for a specific message type from the list if it had one,
					// or check the selected item *after* list.Update().
					// For ClassManagementView -> StudentView navigation, it's already specific.
					if key.Matches(msg, m.keys.Select) { // Check if the original key was 'select'
						selectedListItem := m.list.SelectedItem()
						if selectedListItem != nil {
							if class, ok := selectedListItem.(listItemClass); ok && m.currentView == app.ClassManagementView {
								m.selectedClass = &class.Class
								m.currentView = app.View(99) // Temporary view for students
								m.list.Title = fmt.Sprintf("Alunos da Turma: %s", m.selectedClass.Name)
								cmds = append(cmds, m.loadStudentsForClass(class.ID()))
								log.Printf("TUI: Navegando para lista de alunos da turma %s.", m.selectedClass.Name)
							}
							// TODO: Add similar logic for item selection in other list-based views if needed
						}
					}
				}
			}
		}

	case spinner.TickMsg:
		if m.isLoading { // Only tick spinner if loading
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

	// case tea.WindowSizeMsg: // MOVED to the top of the switch
	// 	h, v := docStyle.GetFrameSize()
	// 	m.list.SetSize(msg.Width-h, msg.Height-v-lipgloss.Height(m.headerView())-lipgloss.Height(m.footerView()))
	}

	// Handle list updates if not loading and no specific key handled above,
	// ONLY if the current view actually uses the list.
	if m.currentViewUsesList() && !m.isLoading {
		var listCmd tea.Cmd
		// Ensure m.list is not nil before calling Update.
		// Though, if currentViewUsesList is true, m.list should have been initialized
		// when transitioning to that view (e.g., in classesLoadedMsg or studentsLoadedMsg).
		// However, a defensive check or ensuring list is always initialized might be good.
		// For now, assuming it's initialized if currentViewUsesList() is true.
		m.list, listCmd = m.list.Update(msg)
		cmds = append(cmds, listCmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the TUI.
func (m Model) View() string {
	log.Printf("TUI: View - INÍCIO. currentView: %s, isLoading: %t, err: %v", m.currentView.String(), m.isLoading, m.err)

	if m.err != nil {
		log.Println("TUI: View - Renderizando erro.")
		// Simplified error view
		return docStyle.Render(fmt.Sprintf("Ocorreu um erro: %v\n\nPressione qualquer tecla para sair.", m.err))
	}

	// Temporariamente, vamos remover o isLoading da condição de forçar a dashboard para ver se é isso.
	// if m.isLoading {
	// 	log.Printf("TUI: View - Renderizando tela de carregamento para %s.", m.currentView.String())
	// 	loadingText := fmt.Sprintf("%s Carregando...", m.spinner.View())
	// 	if m.currentView == app.ClassManagementView && len(m.classes) == 0 {
	// 		loadingText += "\n\nNenhuma turma encontrada ainda."
	// 	} else if m.currentView == app.View(99) && len(m.students) == 0 && m.selectedClass != nil {
	// 		loadingText += fmt.Sprintf("\n\nNenhum aluno encontrado para a turma %s.", m.selectedClass.Name)
	// 	} else if m.currentView == app.DashboardView && m.isLoading { // Specific loading for dashboard if any
	// 		loadingText = fmt.Sprintf("%s Carregando Dashboard...", m.spinner.View())
	// 	}
	// 	return docStyle.Render(loadingText)
	// }

	if m.currentView == app.DashboardView {
		log.Println("TUI: View - Chamando renderDashboardView().")
		return m.renderDashboardView()
	}

	// Se não for DashboardView, verificar isLoading para outras views
	if m.isLoading {
		log.Printf("TUI: View - Renderizando tela de carregamento para %s.", m.currentView.String())
		loadingText := fmt.Sprintf("%s Carregando %s...", m.spinner.View(), m.currentView.String())
		// Adicionar condições específicas de texto de carregamento se necessário
		return docStyle.Render(loadingText)
	}

	log.Printf("TUI: View - Renderizando view de lista para %s.", m.currentView.String())
	return docStyle.Render(m.headerView() + "\n" + m.list.View() + "\n" + m.footerView())
}

func (m Model) headerView() string {
	var titleStr string
	switch m.currentView {
	case app.DashboardView:
		titleStr = "Dashboard Principal Vigenda"
	case app.ClassManagementView:
		titleStr = "Gerenciar Turmas"
	case app.View(99): // Student view
		if m.selectedClass != nil {
			titleStr = fmt.Sprintf("Alunos - %s", m.selectedClass.Name)
		} else {
			titleStr = "Alunos"
		}
	default: // For other views like TaskManagement, AssessmentManagement, etc.
		titleStr = m.currentView.String()
	}
	return titleStyle.Render(titleStr)
}

func (m Model) footerView() string {
	// Basic help. Could be more context-aware.
	// return helpStyle.Render(fmt.Sprintf("↑/↓: Navegar | Enter: Selecionar | Esc: Voltar | q: Sair"))
	// Updated to be context-aware
	var help string
	switch m.currentView {
	case app.DashboardView:
		help = "↑/↓: Menu | Enter: Selecionar | q: Sair" // TODO: Add Tab for panels later
	default: // For list-based views
		help = "↑/↓: Navegar | Enter: Selecionar | Esc: Voltar | q: Sair"
	}
	return helpStyle.Render(help)
}

func (m Model) renderDashboardView() string {
	// 1. Renderizar o Menu Lateral Esquerdo
	menuWidth := 25 // Largura do menu lateral (aumentada para nomes completos)
	menuItemsContent := []string{}

	// Adiciona um título ao menu
	menuTitle := titleStyle.Copy().Width(menuWidth).Align(lipgloss.Center).Render("Navegação")
	menuItemsContent = append(menuItemsContent, menuTitle, "") // Adiciona espaço após o título

	for i, item := range m.dashboardMainMenuItems {
		if i == m.dashboardMainMenuCursor {
			menuItemsContent = append(menuItemsContent, selectedStyle.Width(menuWidth-2).Render("▶ "+item)) // -2 for padding/border of selectedStyle
		} else {
			menuItemsContent = append(menuItemsContent, deselectedStyle.Width(menuWidth-2).Render("  "+item)) // -2 for padding of deselectedStyle
		}
	}
	menuPanelStyle := lipgloss.NewStyle().
		Width(menuWidth).
		// Height(m.height - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView()) - 2). // Adjust height dynamically if needed, -2 for docStyle margin
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")). // Purple border
		Padding(1, 0)

	menuPanel := menuPanelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, menuItemsContent...))

	// 2. Renderizar os Painéis de Conteúdo da Dashboard (Placeholder por enquanto)
	contentWidth := m.width - menuWidth - docStyle.GetHorizontalFrameSize() - 2 // -2 for potential spacing between panels
	if contentWidth < 10 {
		contentWidth = 10 // Minimum content width
	}

	placeholderContent := fmt.Sprintf("Conteúdo da Dashboard Principal.\n\nLargura disponível para conteúdo: %d\n\nFuturos painéis:\n- Foco Atual\n- Agenda do Dia\n- Pendências Urgentes", contentWidth)
	contentPanelStyle := lipgloss.NewStyle().
		Width(contentWidth).
		// Height(m.height - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView()) - 2). // Match menu height approx.
		Border(lipgloss.NormalBorder()).
		Padding(1)

	contentPanel := contentPanelStyle.Render(placeholderContent)

	// 3. Combinar o menu e os painéis de conteúdo usando lipgloss.JoinHorizontal
	dashboardBody := lipgloss.JoinHorizontal(
		lipgloss.Top,
		menuPanel,
		contentPanel,
	)

	// Header e Footer específicos da Dashboard
	// O headerView() e footerView() já são contextuais, então podemos usá-los.
	// Apenas garantindo que eles retornem o texto correto para DashboardView.
	finalHeader := m.headerView() // headerView() já é contextual e renderiza o título da Dashboard
	finalFooter := m.footerView() // footerView() já é contextual e renderiza as dicas da Dashboard

	// Envolve o corpo da dashboard com docStyle para margens consistentes com outras views
	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Left, finalHeader, dashboardBody, finalFooter))
}

// Custom list items
type listItemClass struct {
	models.Class
}

func (lic listItemClass) Title() string       { return lic.Name }
func (lic listItemClass) Description() string { return fmt.Sprintf("ID: %d, Disciplina ID: %d", lic.ID, lic.SubjectID) }
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
