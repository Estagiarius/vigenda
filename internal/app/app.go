package app

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	// Define other styles as needed
)

// Model represents the main application model.
type Model struct {
	list         list.Model
	currentView  View
	width        int
	height       int
	quitting     bool
	err          error // To store any critical errors for display
	// Potentially, sub-models for different views will be added here
	// e.g., tasksModel tasks.Model
}

// Init initializes the application model.
func (m Model) Init() tea.Cmd {
	return nil // No initial command
}

// New creates a new instance of the application model.
func New() Model {
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

	return Model{
		list:        l,
		currentView: DashboardView, // Start with the main menu
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
		m.list.SetSize(msg.Width-appStyle.GetHorizontalPadding(), msg.Height-appStyle.GetVerticalPadding()-lipgloss.Height(m.list.Title)-2)
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
						// Here you might load data or initialize the sub-model for the new view
						// For now, just switching the view enum is enough for placeholder views.
					}
				}
				return m, nil
			}
			m.list, cmd = m.list.Update(msg) // Update the list component
		default: // Other views (TaskManagementView, ClassManagementView, etc.)
			if key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q"))) {
				// Go back to the main menu from any other view
				m.currentView = DashboardView
				// You might want to reset the state of the sub-view here
				return m, nil
			}
			// Later, pass messages to sub-models:
			// switch m.currentView {
			// case TaskManagementView:
			//   m.tasksModel, cmd = m.tasksModel.Update(msg)
			// }
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

	switch m.currentView {
	case DashboardView: // Main Menu
		viewContent = m.list.View()
	default: // Placeholder for other views
		viewContent = fmt.Sprintf("Você está na visão: %s\n\nPressione 'esc' ou 'q' para voltar ao menu principal.", m.currentView.String())
	}

	// Basic help footer
	help := "\nNavegue com ↑/↓, selecione com Enter."
	if m.currentView != DashboardView {
		help = "\nPressione 'esc' ou 'q' para voltar ao menu."
	}


	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		viewContent,
		lipgloss.NewStyle().MarginTop(1).Render(help),
	))
}

// StartApp is a helper to run the BubbleTea program.
func StartApp() {
	p := tea.NewProgram(New(), tea.WithAltScreen()) // Using AltScreen is common for TUIs
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar o Vigenda TUI: %v\n", err)
		os.Exit(1)
	}
}
