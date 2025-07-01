package app

import (
	"fmt"
	"log" // Adicionado para logging
	"os"

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

// Model represents the main application model.
type Model struct {
	list             list.Model
	currentView      View
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
	return nil // No initial command, sub-models handle their own Init
}

// New creates a new instance of the application model.
// It requires services to be injected for its sub-models.
// Changed to return *Model
func New(ts service.TaskService, cs service.ClassService, as service.AssessmentService, qs service.QuestionService, ps service.ProofService /* add other services as params */) *Model {
	// Define menu items using the View enum for safer mapping
	menuItems := []list.Item{
		menuItem{title: DashboardView.String(), view: DashboardView}, // Dashboard is the menu itself
		menuItem{title: TaskManagementView.String(), view: TaskManagementView},
		menuItem{title: ClassManagementView.String(), view: ClassManagementView},
		menuItem{title: AssessmentManagementView.String(), view: AssessmentManagementView},
		menuItem{title: QuestionBankView.String(), view: QuestionBankView},
		menuItem{title: ProofGenerationView.String(), view: ProofGenerationView},
	}

	l := list.New(menuItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Vigenda - Menu Principal"
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
	pm := proofs.New(ps) // Initialize proofs model

	return &Model{ // Return pointer
		list:              l,
		currentView:       DashboardView, // Start with the main menu
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
	var cmd tea.Cmd
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
		// We don't strictly need to capture the updated models here if WindowSizeMsg only affects size
		// and doesn't return a new model instance, but it's safer if it might.
		var tempModel tea.Model
		tempModel, subCmd = m.tasksModel.Update(msg)
		m.tasksModel = tempModel.(*tasks.Model) // Corrected
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.classesModel.Update(msg)
		m.classesModel = tempModel.(*classes.Model) // Corrected
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.assessmentsModel.Update(msg)
		m.assessmentsModel = tempModel.(*assessments.Model) // Corrected
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.questionsModel.Update(msg)
		m.questionsModel = tempModel.(*questions.Model) // Corrected
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.proofsModel.Update(msg)
		m.proofsModel = tempModel.(*proofs.Model) // Corrected
		cmds = append(cmds, subCmd)
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		// Global quit
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))) {
			m.quitting = true
			cmds = append(cmds, tea.Quit)
			return m, tea.Batch(cmds...)
		}

		if m.currentView == DashboardView {
			// Handle navigation for the main menu
			var listCmd tea.Cmd
			m.list, listCmd = m.list.Update(msg)
			cmds = append(cmds, listCmd)

			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				selectedItem, ok := m.list.SelectedItem().(menuItem)
				if ok && selectedItem.view != DashboardView {
					m.currentView = selectedItem.view
					log.Printf("AppModel: Mudando para view %v", m.currentView)
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
			} else if key.Matches(msg, key.NewBinding(key.WithKeys("q"))) { // 'q' to quit from dashboard
				m.quitting = true
				cmds = append(cmds, tea.Quit)
			}
			return m, tea.Batch(cmds...)
		}
		// If not in DashboardView, KeyMsg will be passed to the submodel delegation below.

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
	var submodelCmd tea.Cmd
	switch m.currentView {
	case TaskManagementView:
		var updatedModel tea.Model
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
	case DashboardView: // Main Menu
		viewContent = m.list.View()
		help = "\nNavegue com ↑/↓, selecione com Enter. Pressione 'q' para sair."
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
