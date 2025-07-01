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
	"vigenda/internal/app/proofs"      // Import for the proofs module
	"vigenda/internal/app/questions"   // Import for the questions module
	"vigenda/internal/app/tasks"       // Import for the tasks module
	"vigenda/internal/service"         // Import for service interfaces
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	// Define other styles as needed
)

	"vigenda/internal/service"         // Import for service interfaces
	"vigenda/internal/tui"             // Import for the TUI model we developed
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	// Define other styles as needed
)

// Model represents the main application model.
type Model struct {
	// list             list.Model // REMOVED - Dashboard TUI will handle its own menu
	currentView View
	dashboardTUI *tui.Model // ADDED - Our TUI model for the dashboard

	tasksModel       *tasks.Model       // Changed to pointer
	classesModel     *classes.Model     // Changed to pointer
	assessmentsModel *assessments.Model // Changed to pointer
	questionsModel   *questions.Model   // Changed to pointer
	proofsModel      *proofs.Model      // Changed to pointer
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
	// ... other services
}

// Init initializes the application model.
// Changed to pointer receiver
func (m *Model) Init() tea.Cmd {
	// If the dashboard TUI has an Init, we might want to call it here,
	// or let its first Update handle initialization.
	// For now, tui.Model.Init() returns m.spinner.Tick, which is fine.
	if m.currentView == DashboardView && m.dashboardTUI != nil {
		return m.dashboardTUI.Init()
	}
	return nil // No initial command for other sub-models here, they handle their own Init upon view switch
}

// New creates a new instance of the application model.
// It requires services to be injected for its sub-models.
// Changed to return *Model
func New(ts service.TaskService, cs service.ClassService, as service.AssessmentService, qs service.QuestionService, ps service.ProofService /* add other services as params */) *Model {
	log.Println("app.New: Criando novo app.Model")

	// Initialize our tui.Model for the Dashboard
	// tui.NewTUIModel currently expects classService. This might need to be expanded
	// if the dashboard requires more services.
	dashboardTUIModel := tui.NewTUIModel(cs) // Pass classService

	// Initialize other sub-models
	tm := tasks.New(ts)
	cm := classes.New(cs)
	am := assessments.New(as)
	qm := questions.New(qs)
	cm := classes.New(cs)
	am := assessments.New(as)
	qm := questions.New(qs)
	pm := proofs.New(ps) // Initialize proofs model

	return &Model{ // Return pointer
		currentView:       DashboardView, // Start with the main menu (which is now our tui.Model)
		dashboardTUI:      &dashboardTUIModel, // Store the pointer to tui.Model
		tasksModel:        tm, // tasks.New now returns *tasks.Model
		taskService:       ts,
		classesModel:      cm, // classes.New now returns *classes.Model
		classService:      cs,
		assessmentsModel:  am, // assessments.New now returns *assessments.Model
		assessmentService: as,
		questionsModel:    qm, // questions.New now returns *questions.Model
		questionService:   qs,
		proofsModel:       pm, // proofs.New now returns *proofs.Model
		proofService:      ps,
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
	// var cmd tea.Cmd // Removida declaração no escopo da função
	var cmds []tea.Cmd // Use a slice to collect commands

	// First switch for messages that AppModel handles directly or globally
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		listHeight := msg.Height - appStyle.GetVerticalPadding() - lipgloss.Height(m.list.Title) - 2
		m.list.SetSize(msg.Width-appStyle.GetHorizontalPadding(), listHeight)

		// Propagate WindowSizeMsg to all submodels so they can resize
		var subCmd tea.Cmd
		var tempModel tea.Model

		if m.dashboardTUI != nil {
			// dashboardTUI.Update returns tui.Model, not *tui.Model
			// but tea.Model is an interface, so direct assignment works if it fits.
			// However, tui.Model.Update returns (tea.Model, tea.Cmd)
			// and we need to update m.dashboardTUI with the new state.
			// Since tui.Model is not a pointer in its own Update, we need to handle this carefully.
			// For now, let's assume tui.Model's Update method correctly modifies itself if it's a value type,
			// or we get the updated model if it's a pointer type.
			// The tui.Model.Update method we wrote returns (tea.Model, tea.Cmd)
			// where tea.Model is the updated tui.Model.
			updatedDashboardModel, dashboardCmd := m.dashboardTUI.Update(msg)
			// If tui.Model.Update returns the updated model as tea.Model, we need to cast it back.
			// This cast might be problematic if tui.Model is not a pointer and Update returns a modified copy.
			// Let's assume tui.NewTUIModel returns a value and tui.Update returns a new value.
			// So m.dashboardTUI should store the returned tea.Model.
			// This needs careful handling of pointer vs value types.
			// For simplicity, let's assume tui.Model's Update handles internal state correctly
			// and we mainly care about the command. If tui.Model needs to be replaced,
			// its Update method should return the new instance.
			// Our tui.Model.Update returns (tea.Model, tea.Cmd)
			// Let's make dashboardTUI a value type tui.Model for now if NewTUIModel returns value.
			// Current tui.NewTUIModel returns tui.Model (value).
			// So, m.dashboardTUI.Update(msg) would be on a copy if dashboardTUI is not a pointer.
			// Let's make m.dashboardTUI a pointer, and ensure NewTUIModel is consistent.
			// NewTUIModel returns tui.Model (value). So we store *tui.Model.
			// Then dashboardTUI.Update will be called on the pointer.
			// tui.Model.Update should be func (m Model) ... to return the new state.
			// This is getting complex. Let's simplify:
			// 1. tui.Model.Update will modify its own state if it's a pointer receiver.
			// 2. Or, tui.Model.Update returns the new tui.Model.
			// Our tui.Model.Update is `func (m Model) Update(...) (tea.Model, tea.Cmd)`
			// This means it returns a new `Model` instance.
			if m.currentView == DashboardView && m.dashboardTUI != nil {
				var dashboardCmd tea.Cmd
				*m.dashboardTUI, dashboardCmd = m.dashboardTUI.Update(msg) // Update m.dashboardTUI with the new model
				cmds = append(cmds, dashboardCmd)
			}
		}

		// Propagate to other sub-models (unchanged)
		tempModel, subCmd = m.tasksModel.Update(msg)
		if temp, ok := tempModel.(*tasks.Model); ok { m.tasksModel = temp }
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.classesModel.Update(msg)
		if temp, ok := tempModel.(*classes.Model); ok { m.classesModel = temp }
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.assessmentsModel.Update(msg)
		if temp, ok := tempModel.(*assessments.Model); ok { m.assessmentsModel = temp }
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.questionsModel.Update(msg)
		if temp, ok := tempModel.(*questions.Model); ok { m.questionsModel = temp }
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.proofsModel.Update(msg)
		if temp, ok := tempModel.(*proofs.Model); ok { m.proofsModel = temp }
		cmds = append(cmds, subCmd)
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		// Global quit
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))) {
			m.quitting = true
			cmds = append(cmds, tea.Quit)
			return m, tea.Batch(cmds...)
		}

		if m.currentView == DashboardView && m.dashboardTUI != nil {
			// Delegate KeyMsg to the dashboard TUI
			// The dashboardTUI's Update method will handle its internal navigation (menu cursor)
			// and will also decide if a view change is needed (e.g., user selected "Gerenciar Tarefas")
			updatedDashboardModel, dashboardCmd := m.dashboardTUI.Update(msg)
			*m.dashboardTUI = updatedDashboardModel.(tui.Model) // Update with the new state
			cmds = append(cmds, dashboardCmd)

			// Now, check if the dashboardTUI signaled a desire to change the app-level view
			// This requires a new mechanism, e.g., a message or a state in tui.Model
			// For now, we assume tui.Model's own `currentView` changes, and we need to reflect that
			// in app.Model. This is a temporary, somewhat coupled approach.
			// A better way: tui.Model.Update returns a specific tea.Msg for view change.

			// Let's assume tui.Model now has a method `GetRequestedAppView() (app.View, bool)`
			// Or, simpler for now: if tui.Model's internal currentView is NOT DashboardView,
			// it means it wants to switch.
			if m.dashboardTUI.CurrentView() != DashboardView {
				requestedView := m.dashboardTUI.CurrentView()
				log.Printf("AppModel: DashboardTUI solicitou mudança para view %s", requestedView.String())
				// Reset dashboardTUI's internal view back to Dashboard, so it's ready next time
				m.dashboardTUI.ResetToDashboardView() // Needs to be implemented in tui.Model

				m.currentView = requestedView // Change app.Model's currentView
				// Dispatch Init command for the selected submodel
				switch m.currentView {
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
			// Handle 'q' to quit directly from dashboard TUI if it doesn't handle it internally
			// (Our tui.Model handles 'q' for its own context, this might be redundant or a fallback)
			if key.Matches(msg, key.NewBinding(key.WithKeys("q"))) && m.currentView == DashboardView {
                 // Check if dashboardTUI already handled quit
                 // This is tricky. For now, let app.Model handle global quit from dashboard.
				m.quitting = true
				cmds = append(cmds, tea.Quit)
			}
			return m, tea.Batch(cmds...)
		}
		// If not in DashboardView, KeyMsg will be passed to the submodel delegation below.

	case error: // This case should be after specific message types if they can also be errors.
		m.err = msg
		log.Printf("AppModel: Erro global recebido: %v", msg)
		// Optionally, you might want to switch to an error view or quit
		// For now, just store the error. The View method can display it.
		// cmds = append(cmds, tea.Quit) // Uncomment to quit on any unhandled error
		return m, tea.Batch(cmds...) // Return accumulated commands

		// Default case for the first switch: if the msg type wasn't WindowSizeMsg, KeyMsg (for dashboard), or error,
		// it will fall through to the submodel delegation logic.
	}

	// Second stage: Delegate message to the active submodel if not handled above
	// This is where messages like `classes.fetchedClassesMsg` should be handled.
	// var submodelCmd tea.Cmd // Declarar localmente dentro de cada case se necessário
	switch m.currentView {
	case TaskManagementView:
		var updatedModel tea.Model
		var submodelCmd tea.Cmd // Declarar localmente
		updatedModel, submodelCmd = m.tasksModel.Update(msg) // msg is the original tea.Msg
		m.tasksModel = updatedModel.(*tasks.Model) // Corrected type assertion to pointer
		cmds = append(cmds, submodelCmd)
		// Handle 'esc' to go back to dashboard
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.tasksModel.IsFocused() { // Check if submodel itself wants to handle Esc
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para DashboardView a partir de TaskManagementView.")
			}
		}
	case ClassManagementView:
		log.Printf("AppModel: Update (delegação) - CurrentView=ClassManagementView, encaminhando msg tipo %T para ClassesModel.Update", msg)
		var updatedModel tea.Model
		var submodelCmd tea.Cmd // Declarar localmente
		updatedModel, submodelCmd = m.classesModel.Update(msg) // msg is the original tea.Msg
		m.classesModel = updatedModel.(*classes.Model) // Corrected type assertion to pointer
		cmds = append(cmds, submodelCmd)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.classesModel.IsFocused() {
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para DashboardView a partir de ClassManagementView.")
			}
		}
	case AssessmentManagementView:
		var updatedModel tea.Model
		var submodelCmd tea.Cmd // Declarar localmente
		updatedModel, submodelCmd = m.assessmentsModel.Update(msg)
		m.assessmentsModel = updatedModel.(*assessments.Model) // Corrected type assertion to pointer
		cmds = append(cmds, submodelCmd)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.assessmentsModel.IsFocused() {
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para DashboardView a partir de AssessmentManagementView.")
			}
		}
	case QuestionBankView:
		var updatedModel tea.Model
		var submodelCmd tea.Cmd // Declarar localmente
		updatedModel, submodelCmd = m.questionsModel.Update(msg)
		m.questionsModel = updatedModel.(*questions.Model) // Corrected type assertion to pointer
		cmds = append(cmds, submodelCmd)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.questionsModel.IsFocused() {
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para DashboardView a partir de QuestionBankView.")
			}
		}
	case ProofGenerationView:
		var updatedModel tea.Model
		var submodelCmd tea.Cmd // Declarar localmente
		updatedModel, submodelCmd = m.proofsModel.Update(msg)
		m.proofsModel = updatedModel.(*proofs.Model) // Corrected type assertion to pointer
		cmds = append(cmds, submodelCmd)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.proofsModel.IsFocused() {
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para DashboardView a partir de ProofGenerationView.")
			}
		}
		// Note: DashboardView itself doesn't have a submodel Update to call here,
		// its interactions (list navigation) are handled in the tea.KeyMsg case of the first switch.
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
		return appStyle.Render(fmt.Sprintf("Ocorreu um erro: %v\nPressione 'q' para sair ou 'esc' para voltar ao menu.", m.err))
	}

	var viewContent string
	var help string

	switch m.currentView {
	case DashboardView: // Main Menu is now our tui.Model's DashboardView
		if m.dashboardTUI != nil {
			viewContent = m.dashboardTUI.View()
			// Help text for dashboardTUI is handled by its own footerView,
			// so we might not need to append a global help string here,
			// or we can make it very minimal.
			// For now, let dashboardTUI handle its own help.
			help = "" // Or a very generic app-level help if needed
		} else {
			viewContent = "Erro: Dashboard TUI não inicializado."
			help = "\nPressione Ctrl+C para sair."
		}
	case TaskManagementView:
		viewContent = m.tasksModel.View()
		help = "\nNavegue na tabela com ↑/↓. Pressione 'esc' para voltar ao menu."
	case ClassManagementView:
		viewContent = m.classesModel.View()
		help = "\nNavegue com ↑/↓, Enter para selecionar. 'esc' para voltar/cancelar."
	case AssessmentManagementView:
		viewContent = m.assessmentsModel.View()
		help = "\nNavegue com ↑/↓, Enter para selecionar. 'esc' para voltar/cancelar."
	case QuestionBankView:
		viewContent = m.questionsModel.View()
		help = "\nUse Enter para selecionar/submeter. 'esc' para voltar/cancelar."
	case ProofGenerationView:
		viewContent = m.proofsModel.View()
		help = "\nUse Tab/Setas para navegar no formulário. 'esc' para voltar."
	default:
		viewContent = fmt.Sprintf("Você está na visão: %s\n\nPressione 'esc' ou 'q' para voltar ao menu principal.", m.currentView.String())
		help = "\nPressione 'esc' ou 'q' para voltar ao menu."
	}

	finalRender := appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		viewContent,
		lipgloss.NewStyle().MarginTop(1).Render(help),
	))
	// fmt.Printf("[LOG app.Model] View(): Final render string length: %d\n", len(finalRender)) // Potencialmente muito verboso
	return finalRender
}

// StartApp is a helper to run the BubbleTea program.
// It requires services to be passed for initializing the main model.
func StartApp(ts service.TaskService, cs service.ClassService, as service.AssessmentService, qs service.QuestionService, ps service.ProofService /*, other services */) {
	// New now returns *Model, so model is already a pointer.
	model := New(ts, cs, as, qs, ps /*, other services */)
	// tea.NewProgram expects tea.Model, and *Model now implements tea.Model.
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		// Use log for consistency, as file logging is set up.
		log.Fatalf("Error running BubbleTea program: %v", err)
		// os.Exit(1) // log.Fatalf will exit
	}
}
