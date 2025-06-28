package app

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"vigenda/internal/app/tasks" // Import for the tasks module
	"vigenda/internal/service"   // Import for service interfaces
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	// Define other styles as needed
)

// Model represents the main application model.
type Model struct {
	list        list.Model
	currentView View
	tasksModel  tasks.Model // Add tasks model
	// Add other sub-models here as they are developed e.g. classesModel classes.Model
	width    int
	height   int
	quitting bool
	err      error // To store any critical errors for display

	// Services - these will be injected
	taskService service.TaskService
	// classService service.ClassService
	// ... other services
}

// Init initializes the application model.
func (m Model) Init() tea.Cmd {
	return nil // No initial command, sub-models handle their own Init
}

// New creates a new instance of the application model.
// It requires services to be injected for its sub-models.
func New(ts service.TaskService /* add other services as params */) Model {
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
	// cm := classes.New(cs) // Example for future

	return Model{
		list:         l,
		currentView:  DashboardView, // Start with the main menu
		tasksModel:   tm,
		taskService:  ts,
		// classesModel: cm,
		// classService: cs,
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
		// if m.classesModel != nil { m.classesModel.SetSize(subViewWidth, subViewHeight) }

		return m, nil

	case tea.KeyMsg:
		// Global keybindings
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))) {
			m.quitting = true
			return m, tea.Quit
		}

		// View-specific keybindings
		switch m.currentView {
		case DashboardView: // Main Menu
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("q"))):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				selectedItem, ok := m.list.SelectedItem().(menuItem)
				if ok {
					// DashboardView itself is the menu, so selecting it does nothing new
					// or could refresh a dashboard if it were more complex.
					// For other items, switch the view.
					if selectedItem.view != DashboardView {
						m.currentView = selectedItem.view
						// If switching to a view that needs initialization (like loading data)
						if m.currentView == TaskManagementView {
							cmd = m.tasksModel.Init() // Trigger task loading
						}
						// Add similar blocks for other views if they need explicit Init on switch
					}
				}
				return m, cmd // Return cmd which might come from sub-model Init
			}
			m.list, cmd = m.list.Update(msg) // Update the list component

		case TaskManagementView:
			var updatedTasksModel tasks.Model
			updatedTasksModel, cmd = m.tasksModel.Update(msg)
			m.tasksModel = updatedTasksModel
			// Check for 'esc' or 'q' to navigate back, if not handled by sub-model
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q"))) {
				if !m.tasksModel.IsFocused() { // Example: if sub-model has focus concept
					m.currentView = DashboardView
				}
			}
		// Add cases for other views like ClassManagementView etc.
		// default: // Other views
		//  if key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q"))) {
		//    m.currentView = DashboardView
		//    return m, nil
		//  }
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
		// Help text for task view might be part of tasksModel.View() or defined here
		help = "\nNavegue na tabela com ↑/↓. Pressione 'esc' ou 'q' para voltar ao menu."
	// Add cases for other views
	default: // Placeholder for other views not yet fully integrated
		viewContent = fmt.Sprintf("Você está na visão: %s\n\nPressione 'esc' ou 'q' para voltar ao menu principal.", m.currentView.String())
		help = "\nPressione 'esc' ou 'q' para voltar ao menu."
	}

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		viewContent,
		lipgloss.NewStyle().MarginTop(1).Render(help),
	))
}

// StartApp is a helper to run the BubbleTea program.
// It requires services to be passed for initializing the main model.
func StartApp(ts service.TaskService /*, other services */) {
	model := New(ts /*, other services */)
	p := tea.NewProgram(model, tea.WithAltScreen()) // Using AltScreen is common for TUIs
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar o Vigenda TUI: %v\n", err)
		os.Exit(1)
	}
}
