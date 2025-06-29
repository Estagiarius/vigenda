package app

import (
	"fmt"
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
	tasksModel       tasks.Model       // Tasks sub-model
	classesModel     classes.Model     // Classes sub-model
	assessmentsModel assessments.Model // Assessments sub-model
	questionsModel   questions.Model   // Questions sub-model
	proofsModel      proofs.Model      // Proofs sub-model
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
func (m Model) Init() tea.Cmd {
	fmt.Println("[LOG app.Model] Init() called.")
	return nil // No initial command, sub-models handle their own Init
}

// New creates a new instance of the application model.
// It requires services to be injected for its sub-models.
func New(ts service.TaskService, cs service.ClassService, as service.AssessmentService, qs service.QuestionService, ps service.ProofService /* add other services as params */) Model {
	fmt.Println("[LOG app.Model] New() called.")
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

	return Model{
		list:              l,
		currentView:       DashboardView, // Start with the main menu
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

// View type and String method (moved from views.go for simplicity in logging, can be kept separate)
type View int

const (
	DashboardView View = iota
	TaskManagementView
	ClassManagementView
	AssessmentManagementView
	QuestionBankView
	ProofGenerationView
	// Add other views here
)

func (v View) String() string {
	return [...]string{
		"Dashboard",
		"Gerenciamento de Tarefas",
		"Gerenciamento de Turmas",
		"Gerenciamento de Avaliações",
		"Banco de Questões",
		"Geração de Provas",
	}[v]
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
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Adjust list size, leaving space for title and help
		listHeight := msg.Height - appStyle.GetVerticalPadding() - lipgloss.Height(m.list.Title) - 2
		m.list.SetSize(msg.Width-appStyle.GetHorizontalPadding(), listHeight)

		// Propagate size to sub-models
		// Subtracting space used by main app's padding/title, etc.
		// This might need adjustment based on how much space the main app's chrome takes.
		subViewWidth := msg.Width - appStyle.GetHorizontalPadding()
		subViewHeight := msg.Height - appStyle.GetVerticalPadding() // Example: if subview takes full height within padding

		m.tasksModel.SetSize(subViewWidth, subViewHeight)
		m.classesModel.SetSize(subViewWidth, subViewHeight)
		m.assessmentsModel.SetSize(subViewWidth, subViewHeight)
		m.questionsModel.SetSize(subViewWidth, subViewHeight)
		m.proofsModel.SetSize(subViewWidth, subViewHeight) // Propagate to proofsModel too

		return m, nil

	case tea.KeyMsg:
		// Global keybindings
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))) {
			fmt.Println("[LOG app.Model] Update(): Ctrl+C pressed. Quitting.")
			m.quitting = true
			return m, tea.Quit
		}
		fmt.Printf("[LOG app.Model] Update(): Received KeyMsg: %s, CurrentView: %s\n", msg.String(), m.currentView.String())


		// View-specific keybindings
		switch m.currentView {
		case DashboardView: // Main Menu
			fmt.Printf("[LOG app.Model] Update(): Key '%s' in DashboardView\n", msg.String())
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("q"))):
				fmt.Println("[LOG app.Model] Update(): 'q' pressed in Dashboard. Quitting.")
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				fmt.Println("[LOG app.Model] Update(): 'enter' pressed in Dashboard.")
				selectedItem, ok := m.list.SelectedItem().(menuItem)
				if ok {
					fmt.Printf("[LOG app.Model] Update(): Selected item: '%s', Target View: %s\n", selectedItem.title, selectedItem.view.String())
					if selectedItem.view != DashboardView {
						fmt.Printf("[LOG app.Model] Update(): Changing currentView from %s to %s\n", m.currentView.String(), selectedItem.view.String())
						m.currentView = selectedItem.view
						fmt.Printf("[LOG app.Model] Update(): currentView is now %s\n", m.currentView.String())

						// If switching to a view that needs initialization (like loading data)
						if m.currentView == TaskManagementView {
							fmt.Println("[LOG app.Model] Update(): Initializing TaskManagementView.")
							cmd = m.tasksModel.Init()
						} else if m.currentView == ClassManagementView {
							fmt.Println("[LOG app.Model] Update(): Initializing ClassManagementView.")
							cmd = m.classesModel.Init()
						} else if m.currentView == AssessmentManagementView {
							fmt.Println("[LOG app.Model] Update(): Initializing AssessmentManagementView.")
							cmd = m.assessmentsModel.Init()
						} else if m.currentView == QuestionBankView {
							fmt.Println("[LOG app.Model] Update(): Initializing QuestionBankView.")
							cmd = m.questionsModel.Init()
						} else if m.currentView == ProofGenerationView {
							fmt.Println("[LOG app.Model] Update(): Initializing ProofGenerationView.")
							cmd = m.proofsModel.Init()
						}
						if cmd != nil {
							fmt.Printf("[LOG app.Model] Update(): Sub-model Init() returned a command.\n")
						} else {
							fmt.Printf("[LOG app.Model] Update(): Sub-model Init() returned nil command.\n")
						}
					} else {
						fmt.Println("[LOG app.Model] Update(): Dashboard selected, no view change or init.")
					}
				} else {
					fmt.Println("[LOG app.Model] Update(): Enter pressed but no item selected or item is not menuItem.")
				}
				return m, cmd
			}
			fmt.Println("[LOG app.Model] Update(): Updating list component in DashboardView.")
			m.list, cmd = m.list.Update(msg)

		case TaskManagementView:
			var updatedTasksModel tasks.Model
			updatedTasksModel, cmd = m.tasksModel.Update(msg)
			m.tasksModel = updatedTasksModel
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) { // 'q' might be for sub-model actions
				if !m.tasksModel.IsFocused() {
					m.currentView = DashboardView
				}
			}

		case ClassManagementView:
			var updatedClassesModel classes.Model
			updatedClassesModel, cmd = m.classesModel.Update(msg)
			m.classesModel = updatedClassesModel
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) {
				if !m.classesModel.IsFocused() {
					m.currentView = DashboardView
				}
			}

		case AssessmentManagementView:
			var updatedAssessmentsModel assessments.Model
			updatedAssessmentsModel, cmd = m.assessmentsModel.Update(msg)
			m.assessmentsModel = updatedAssessmentsModel
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) {
				if !m.assessmentsModel.IsFocused() {
					m.currentView = DashboardView
				}
			}

		case QuestionBankView:
			var updatedQuestionsModel questions.Model
			updatedQuestionsModel, cmd = m.questionsModel.Update(msg)
			m.questionsModel = updatedQuestionsModel
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) {
				if !m.questionsModel.IsFocused() {
					m.currentView = DashboardView
				}
			}

		case ProofGenerationView:
			var updatedProofsModel proofs.Model
			updatedProofsModel, cmd = m.proofsModel.Update(msg)
			m.proofsModel = updatedProofsModel
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) {
				if !m.proofsModel.IsFocused() {
					m.currentView = DashboardView
				}
			}

		default: // Other views (if any become active without specific handling yet)
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q"))) {
				m.currentView = DashboardView // Default back to dashboard
				return m, nil
			}
		}

	case error:
		m.err = msg
		return m, nil
	}
	return m, cmd
}

// View renders the application's UI.
func (m Model) View() string {
	if m.quitting {
		fmt.Println("[LOG app.Model] View(): Quitting.")
		return appStyle.Render("Saindo do Vigenda...\n")
	}
	if m.err != nil {
		fmt.Printf("[LOG app.Model] View(): Rendering error: %v\n", m.err)
		return appStyle.Render(fmt.Sprintf("Ocorreu um erro: %v\nPressione 'q' para sair ou 'esc' para voltar ao menu.", m.err))
	}

	fmt.Printf("[LOG app.Model] View(): Rendering currentView: %s\n", m.currentView.String())
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
	fmt.Println("[LOG app] StartApp(): called.")
	model := New(ts, cs, as, qs, ps /*, other services */)
	fmt.Println("[LOG app] StartApp(): Model created. Initializing BubbleTea program.")
	p := tea.NewProgram(model, tea.WithAltScreen())
	fmt.Println("[LOG app] StartApp(): Program created. Running...")
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[LOG app] StartApp(): Error running BubbleTea program: %v\n", err)
		os.Exit(1)
	}
}
