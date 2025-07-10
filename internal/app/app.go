package app

import (
	"fmt"
	"log" // Adicionado para logging
	// "os" // Removido import não utilizado

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"vigenda/internal/app/assessments" // Import for the assessments module
	"vigenda/internal/app/classes"     // Import for the classes module
	"vigenda/internal/app/dashboard"   // Import for the new dashboard module
	"vigenda/internal/app/proofs"      // Import for the proofs module
	"vigenda/internal/app/questions"   // Import for the questions module
	"vigenda/internal/app/tasks"       // Import for the tasks module
	"vigenda/internal/service"         // Import for service interfaces
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	// Define other styles as needed
)

// Model represents the main application model.
type Model struct {
	list             list.Model
	currentView      View
	tasksModel       *tasks.Model
	classesModel     *classes.Model
	assessmentsModel *assessments.Model
	questionsModel   *questions.Model
	proofsModel      *proofs.Model
	dashboardModel   *dashboard.Model // Added dashboard model field
	// Add other sub-models here as they are developed
	width    int
	height   int
	quitting bool
	err      error // To store any critical errors for display

	// Services - these will be injected
	taskService       service.TaskService
	classService      service.ClassService
	assessmentService service.AssessmentService
	questionService   service.QuestionService
	proofService      service.ProofService
	lessonService     service.LessonService // Adicionado
	// ... other services
}

// Init initializes the application model.
// Changed to pointer receiver
func (m *Model) Init() tea.Cmd {
	return nil // No initial command, sub-models handle their own Init
}

// New creates a new instance of the application model.
// It requires services to be injected for its sub-models.
// Changed to return *Model
func New(ts service.TaskService, cs service.ClassService, as service.AssessmentService, qs service.QuestionService, ps service.ProofService, ls service.LessonService /* add other services as params */) *Model {
	// Define menu items using the View enum for safer mapping
	menuItems := []list.Item{
		// DashboardView (as main menu) does not need an item for itself, it IS the list.
		menuItem{title: ConcreteDashboardView.String(), view: ConcreteDashboardView}, // New menu item for actual dashboard
		menuItem{title: TaskManagementView.String(), view: TaskManagementView},
		menuItem{title: ClassManagementView.String(), view: ClassManagementView},
		menuItem{title: AssessmentManagementView.String(), view: AssessmentManagementView},
		menuItem{title: QuestionBankView.String(), view: QuestionBankView},
		menuItem{title: ProofGenerationView.String(), view: ProofGenerationView},
	}

	l := list.New(menuItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Vigenda - Menu Principal" // Title for the main menu (DashboardView)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "sair")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "selecionar")),
		}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys // Keep it simple for now

	// Initialize sub-models
	tm := tasks.New(ts)
	cm := classes.New(cs)
	am := assessments.New(as)
	qm := questions.New(qs)
	pm := proofs.New(ps)
	dshModel := dashboard.New(ts, cs, as, ls) // Initialize dashboard model, passing necessary services

	return &Model{ // Return pointer
		list:              l,
		currentView:       DashboardView, // Start with the main menu (DashboardView acts as the container)
		tasksModel:        tm,
		taskService:       ts,
		classesModel:      cm,
		classService:      cs,
		assessmentsModel:  am,
		assessmentService: as,
		questionsModel:    qm,
		questionService:   qs,
		proofsModel:       pm,
		proofService:      ps,
		lessonService:     ls,         // Store lessonService
		dashboardModel:    dshModel, // Assign initialized dashboard model
	}
}

// menuItem holds the title and the corresponding view for a menu entry.
type menuItem struct {
	title string
	view  View
}

func (i menuItem) FilterValue() string { return i.title }
func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return "" } // Could add descriptions later

// Update handles messages and updates the model.
// Changed to pointer receiver
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("AppModel: Update GLOBAL - Recebida msg tipo %T Valor: %v", msg, msg)
	var cmds []tea.Cmd // Use a slice to collect commands

	// First switch for messages that AppModel handles directly or globally
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Adjust list size for main menu
		listHeight := msg.Height - appStyle.GetVerticalPadding() - lipgloss.Height(m.list.Title) - lipgloss.Height(m.list.Help.View(m.list)) - 2 // Approximate height for help
		m.list.SetSize(msg.Width-appStyle.GetHorizontalPadding(), listHeight)

		// Propagate WindowSizeMsg to all submodels so they can resize
		var subCmd tea.Cmd
		var tempModel tea.Model

		// Dashboard model
		tempModel, subCmd = m.dashboardModel.Update(msg)
		m.dashboardModel = tempModel.(*dashboard.Model)
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.tasksModel.Update(msg)
		m.tasksModel = tempModel.(*tasks.Model)
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.classesModel.Update(msg)
		m.classesModel = tempModel.(*classes.Model)
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.assessmentsModel.Update(msg)
		m.assessmentsModel = tempModel.(*assessments.Model)
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.questionsModel.Update(msg)
		m.questionsModel = tempModel.(*questions.Model)
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.proofsModel.Update(msg)
		m.proofsModel = tempModel.(*proofs.Model)
		cmds = append(cmds, subCmd)
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		// Global quit
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))) {
			m.quitting = true
			cmds = append(cmds, tea.Quit)
			return m, tea.Batch(cmds...)
		}

		// Handle interactions based on the current view
		// If in main menu (DashboardView), handle list navigation and selection
		if m.currentView == DashboardView {
			var listCmd tea.Cmd
			m.list, listCmd = m.list.Update(msg)
			cmds = append(cmds, listCmd)

			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				selectedItem, ok := m.list.SelectedItem().(menuItem)
				if ok { // No need to check selectedItem.view != DashboardView, as DashboardView is the container
					m.currentView = selectedItem.view
					log.Printf("AppModel: Mudando para view %s (%d)", m.currentView.String(), m.currentView)
					// Dispatch Init command for the selected submodel
					switch m.currentView {
					case ConcreteDashboardView: // When "Painel de Controle" is selected
						cmds = append(cmds, m.dashboardModel.Init())
					case TaskManagementView:
						cmds = append(cmds, m.tasksModel.Init())
					case ClassManagementView:
						cmds = append(cmds, m.classesModel.Init())
					case AssessmentManagementView:
						cmds = append(cmds, m.assessmentsModel.Init())
					case QuestionBankView:
						cmds = append(cmds, m.questionsModel.Init())
					case ProofGenerationView:
						cmds = append(cmds, m.proofsModel.Init())
					}
				}
			} else if key.Matches(msg, key.NewBinding(key.WithKeys("q"))) {
				m.quitting = true
				cmds = append(cmds, tea.Quit)
			}
			return m, tea.Batch(cmds...)
		}
		// If not in DashboardView (main menu), KeyMsg will be passed to the active submodel's Update method below.

	case error: // Catch global errors (e.g., from Init functions of submodels if not handled there)
		m.err = msg
		log.Printf("AppModel: Erro global recebido: %v", msg)
		// Optionally, switch to an error view or prepare to quit
		// For now, just store the error. The View method can display it.
		return m, tea.Batch(cmds...)
	}

	// Second stage: Delegate message to the active submodel if not handled above
	var submodelCmd tea.Cmd
	var updatedModel tea.Model

	switch m.currentView {
	case ConcreteDashboardView:
		updatedModel, submodelCmd = m.dashboardModel.Update(msg)
		m.dashboardModel = updatedModel.(*dashboard.Model)
		cmds = append(cmds, submodelCmd)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.dashboardModel.IsFocused() {
				m.currentView = DashboardView // Go back to main menu
				log.Println("AppModel: Voltando para o Menu Principal a partir do Painel de Controle.")
			}
		}
	case TaskManagementView:
		updatedModel, submodelCmd = m.tasksModel.Update(msg)
		m.tasksModel = updatedModel.(*tasks.Model)
		cmds = append(cmds, submodelCmd)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.tasksModel.IsFocused() {
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para o Menu Principal a partir de Gerenciar Tarefas.")
			}
		}
	case ClassManagementView:
		updatedModel, submodelCmd = m.classesModel.Update(msg)
		m.classesModel = updatedModel.(*classes.Model)
		cmds = append(cmds, submodelCmd)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.classesModel.IsFocused() {
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para o Menu Principal a partir de Gerenciar Turmas.")
			}
		}
	case AssessmentManagementView:
		updatedModel, submodelCmd = m.assessmentsModel.Update(msg)
		m.assessmentsModel = updatedModel.(*assessments.Model)
		cmds = append(cmds, submodelCmd)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.assessmentsModel.IsFocused() {
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para o Menu Principal a partir de Gerenciar Avaliações.")
			}
		}
	case QuestionBankView:
		updatedModel, submodelCmd = m.questionsModel.Update(msg)
		m.questionsModel = updatedModel.(*questions.Model)
		cmds = append(cmds, submodelCmd)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.questionsModel.IsFocused() {
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para o Menu Principal a partir do Banco de Questões.")
			}
		}
	case ProofGenerationView:
		updatedModel, submodelCmd = m.proofsModel.Update(msg)
		m.proofsModel = updatedModel.(*proofs.Model)
		cmds = append(cmds, submodelCmd)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.proofsModel.IsFocused() {
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para o Menu Principal a partir de Gerar Provas.")
			}
		}
	// DashboardView (main menu) list updates are handled in the tea.KeyMsg section of the first switch.
	}

	return m, tea.Batch(cmds...)
}

// View renders the application's UI.
// Changed to pointer receiver
func (m *Model) View() string {
	if m.quitting {
		return appStyle.Render("Saindo do Vigenda...\n")
	}
	if m.err != nil {
		// More robust error display
		errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
		return appStyle.Render(fmt.Sprintf("Ocorreu um erro crítico: %v\nPressione Ctrl+C para sair.", errorStyle.Render(m.err.Error())))
	}

	var viewContent string
	var help string // Help text can be context-dependent

	switch m.currentView {
	case DashboardView: // This is the Main Menu
		viewContent = m.list.View()
		// Help for the main menu is usually part of the list component itself or can be added
		help = m.list.Help.View(m.list) // Use built-in help view of the list
	case ConcreteDashboardView: // The actual dashboard display
		viewContent = m.dashboardModel.View()
		// Help for ConcreteDashboardView might be part of its own View or defined here
		// For now, assuming its View includes necessary help, or it's simple 'esc' to go back.
		help = "\nPressione 'esc' para voltar ao menu principal."
	case TaskManagementView:
		viewContent = m.tasksModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal." // Simplified help
	case ClassManagementView:
		viewContent = m.classesModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal."
	case AssessmentManagementView:
		viewContent = m.assessmentsModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal."
	case QuestionBankView:
		viewContent = m.questionsModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal."
	case ProofGenerationView:
		viewContent = m.proofsModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal."
	default: // Should not happen if all views are handled
		viewContent = fmt.Sprintf("Visão desconhecida: %s (%d)", m.currentView.String(), m.currentView)
		help = "\nPressione 'esc' ou 'q' para tentar voltar ao menu principal."
	}

	// Combine content and help.
	// Ensure help is only added if viewContent is not already filling the screen or managing its own help.
	finalRender := lipgloss.JoinVertical(lipgloss.Left,
		viewContent,
		lipgloss.NewStyle().MarginTop(1).Render(help), // Add margin to separate help
	)
	return appStyle.Render(finalRender)
}

// StartApp is a helper to run the BubbleTea program.
// It requires services to be passed for initializing the main model.
func StartApp(ts service.TaskService, cs service.ClassService, as service.AssessmentService, qs service.QuestionService, ps service.ProofService, ls service.LessonService /*, other services */) {
	// New now returns *Model, so model is already a pointer.
	model := New(ts, cs, as, qs, ps, ls /*, other services */)
	// tea.NewProgram expects tea.Model, and *Model now implements tea.Model.
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		// Use log for consistency, as file logging is set up.
		log.Fatalf("Error running BubbleTea program: %v", err)
		// os.Exit(1) // log.Fatalf will exit
	}
}
