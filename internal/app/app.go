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
	dashboardModel   *dashboard.Model // Added dashboard model
	tasksModel       *tasks.Model
	classesModel     *classes.Model
	assessmentsModel *assessments.Model
	questionsModel   *questions.Model
	proofsModel      *proofs.Model
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
	return nil // No initial command, sub-models handle their own Init
}

// New creates a new instance of the application model.
// It requires services to be injected for its sub-models.
// Changed to return *Model
func New(ts service.TaskService, cs service.ClassService, as service.AssessmentService, qs service.QuestionService, ps service.ProofService /* add other services as params */) *Model {
	// Define menu items using the View enum for safer mapping
	// The DashboardView itself is now a dedicated view, not the menu.
	// The menu will be the initial state, and selecting "Dashboard" from it (or by default)
	// will switch to the dashboard.Model's view.
	// For now, let's assume starting directly on the dashboard view.
	// If a menu is desired first, currentView would be something like MenuView,
	// and DashboardView would be an item in that menu.

	// Initialize sub-models, including the new dashboard model
	dashModel := dashboard.New(ts, cs) // Pass necessary services
	tm := tasks.New(ts)
	cm := classes.New(cs)
	am := assessments.New(as)
	qm := questions.New(qs)
	pm := proofs.New(ps)

	// For the main application list (menu), if we want one:
	// Let's adjust so that app.Model's list is for the main menu,
	// and DashboardView is one of the views managed by app.Model.
	// The original request implies dashboard is a specific view, not the menu itself.
	menuItems := []list.Item{
		// menuItem{title: "Ver Dashboard", view: DashboardView}, // Item to explicitly go to dashboard
		// Or, if Dashboard is the default initial view, it might not need a menu item if there's no other "main menu" list.
		// Given the existing structure, DashboardView was already 0 (iota).
		// Let's keep the menu items as they were, but the behavior of DashboardView will change.
		// It will now render dashboardModel.View() instead of list.View().
		menuItem{title: "Dashboard", view: DashboardView}, // This will now show the actual dashboard content
		menuItem{title: TaskManagementView.String(), view: TaskManagementView},
		menuItem{title: ClassManagementView.String(), view: ClassManagementView},
		menuItem{title: AssessmentManagementView.String(), view: AssessmentManagementView},
		menuItem{title: QuestionBankView.String(), view: QuestionBankView},
		menuItem{title: ProofGenerationView.String(), view: ProofGenerationView},
	}

	l := list.New(menuItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Vigenda - Menu Principal" // This title might be confusing if DashboardView shows dashboard content
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "sair")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "selecionar")),
		}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys

	return &Model{ // Return pointer
		list:              l,             // This list is for the main navigation menu
		currentView:       DashboardView, // Start with the Dashboard View active
		dashboardModel:    dashModel,
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
		var tempModel tea.Model // tempModel can be reused

		// Dashboard Model
		if m.dashboardModel != nil { // Check if initialized
			tempModel, subCmd = m.dashboardModel.Update(msg)
			if tempModel != nil { // Ensure model is not nil before assertion
				m.dashboardModel = tempModel.(*dashboard.Model)
				cmds = append(cmds, subCmd)
			}
		}

		// Tasks Model
		if m.tasksModel != nil {
			tempModel, subCmd = m.tasksModel.Update(msg)
			if tempModel != nil {
				m.tasksModel = tempModel.(*tasks.Model) // Corrected
				cmds = append(cmds, subCmd)
			}
		}

		// Classes Model
		if m.classesModel != nil {
			tempModel, subCmd = m.classesModel.Update(msg)
			if tempModel != nil {
				m.classesModel = tempModel.(*classes.Model) // Corrected
				cmds = append(cmds, subCmd)
			}
		}

		// Assessments Model
		if m.assessmentsModel != nil {
			tempModel, subCmd = m.assessmentsModel.Update(msg)
			if tempModel != nil {
				m.assessmentsModel = tempModel.(*assessments.Model) // Corrected
				cmds = append(cmds, subCmd)
			}
		}

		// Questions Model
		if m.questionsModel != nil {
			tempModel, subCmd = m.questionsModel.Update(msg)
			if tempModel != nil {
				m.questionsModel = tempModel.(*questions.Model) // Corrected
				cmds = append(cmds, subCmd)
			}
		}

		// Proofs Model
		if m.proofsModel != nil {
			tempModel, subCmd = m.proofsModel.Update(msg)
			if tempModel != nil {
				m.proofsModel = tempModel.(*proofs.Model) // Corrected
				cmds = append(cmds, subCmd)
			}
		}
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		// Global quit
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))) {
			m.quitting = true
			cmds = append(cmds, tea.Quit)
			return m, tea.Batch(cmds...)
		}

		// If currentView is DashboardView (content is shown):
		if m.currentView == DashboardView {
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) {
				// 'esc' from Dashboard content. Placeholder: quit.
				// TODO: Implement logic to show the main menu (m.list)
				// This would involve changing currentView to a new "MenuState"
				// and that state would handle m.list.
				log.Println("AppModel: 'esc' from DashboardView (content). Placeholder: Quitting.")
				m.quitting = true
				cmds = append(cmds, tea.Quit)
			} else if key.Matches(msg, key.NewBinding(key.WithKeys("m"))) { // Example: 'm' to toggle main menu
				// Placeholder for showing main menu list
				log.Println("AppModel: 'm' pressed in DashboardView. TODO: Implement switch to main menu list view.")
				// Example: m.currentView = MainUserMenuView (a new enum)
				// cmds = append(cmds, m.list.Focus()) // If m.list is shown
			}
			// Keys for dashboardModel itself are handled by its own Update via the second switch.
			return m, tea.Batch(cmds...)
		}

		// If currentView is NOT DashboardView (content), it's a sub-module.
		// 'esc' from sub-modules is handled in their specific cases in the second switch
		// (e.g., returning to DashboardView content).
		// Other keys are passed to the sub-module via the second switch.

		// The logic for handling m.list (main menu) needs a dedicated state.
		// For example, if currentView == MainUserMenuView:
		//   var listCmd tea.Cmd
		//   m.list, listCmd = m.list.Update(msg)
		//   cmds = append(cmds, listCmd)
		//   if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
		//     selectedItem, ok := m.list.SelectedItem().(menuItem)
		//     if ok {
		//       m.currentView = selectedItem.view // Switch to Dashboard (content), Tasks, Classes etc.
		//       log.Printf("AppModel: Switched to view %v from menu", m.currentView)
		//       // Call Init for the selected model
		//       switch m.currentView {
		//       case DashboardView:
		//         if m.dashboardModel != nil { cmds = append(cmds, m.dashboardModel.Init()) }
		//       case TaskManagementView:
		//         if m.tasksModel != nil { cmds = append(cmds, m.tasksModel.Init()) }
		//       // ... other cases
		//       }
		//     }
		//   } else if key.Matches(msg, key.NewBinding(key.WithKeys("q"))) { // Quit from menu
		//     m.quitting = true
		//     cmds = append(cmds, tea.Quit)
		//   }
		//   return m, tea.Batch(cmds...)
		// This block is commented out as MainUserMenuView is not yet defined.

	case error:
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
	switch m.currentView {
	case DashboardView: // DashboardView (iota 0) now shows dashboard.Model content
		var updatedModel tea.Model
		var subCmd tea.Cmd
		if m.dashboardModel != nil {
			updatedModel, subCmd = m.dashboardModel.Update(msg)
			if updatedModel != nil {
				m.dashboardModel = updatedModel.(*dashboard.Model)
				cmds = append(cmds, subCmd)
			}
		}
		// 'esc' pressed in DashboardView (content) should navigate away (e.g., show main menu or quit).
		// This specific 'esc' handling for DashboardView (content) could be here or in the KeyMsg block.
		// For now, assuming KeyMsg block handles 'esc' from DashboardView (content).
	case TaskManagementView:
		var updatedModel tea.Model
		var submodelCmd tea.Cmd // Declarar localmente
		updatedModel, submodelCmd = m.tasksModel.Update(msg) // msg is the original tea.Msg
		m.tasksModel = updatedModel.(*tasks.Model) // Corrected type assertion to pointer
		cmds = append(cmds, submodelCmd)
		// Handle 'esc' to go back to dashboard (which now means dashboard content)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.tasksModel.IsFocused() { // Check if submodel itself wants to handle Esc
				m.currentView = DashboardView // Return to DashboardView (content)
				log.Println("AppModel: Voltando para DashboardView (content) a partir de TaskManagementView.")
				if m.dashboardModel != nil {
					cmds = append(cmds, m.dashboardModel.Init()) // Re-initialize dashboard content
				}
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
	case DashboardView: // DashboardView (iota 0) now shows dashboard.Model content
		if m.dashboardModel != nil {
			viewContent = m.dashboardModel.View()
		} else {
			viewContent = "Dashboard model não inicializado."
		}
		help = "\nDashboard. Pressione 'esc' para (placeholder: sair), 'm' para (placeholder: menu)."
	case TaskManagementView:
		viewContent = m.tasksModel.View()
		help = "\nNavegue na tabela com ↑/↓. Pressione 'esc' para voltar ao dashboard."
	case ClassManagementView:
		viewContent = m.classesModel.View()
		help = "\nNavegue com ↑/↓, Enter para selecionar. 'esc' para voltar ao dashboard."
	case AssessmentManagementView:
		viewContent = m.assessmentsModel.View()
		help = "\nNavegue com ↑/↓, Enter para selecionar. 'esc' para voltar ao dashboard."
	case QuestionBankView:
		viewContent = m.questionsModel.View()
		help = "\nUse Enter para selecionar/submeter. 'esc' para voltar ao dashboard."
	case ProofGenerationView:
		viewContent = m.proofsModel.View()
		help = "\nUse Tab/Setas para navegar no formulário. 'esc' para voltar ao dashboard."
	default:
		// This case should ideally not be hit if all views are handled.
		// If it's a new view state for the main menu list, it would be handled here.
		viewContent = fmt.Sprintf("Visão desconhecida: %s. Pressione 'q' para sair.", m.currentView.String())
		help = "\nPressione 'q' para sair."
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
