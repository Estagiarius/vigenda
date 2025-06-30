package tasks

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"context" // Required for service calls
	"fmt"     // For formatting data into table rows
	// "time"    // For formatting dates // No longer directly used for formatting here

	"github.com/charmbracelet/lipgloss"
	"vigenda/internal/models" // Required for models.Task
	"vigenda/internal/service"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	taskService service.TaskService
	table       table.Model
	isLoading   bool // To show a loading indicator
	err         error
	// Future: add states for focused, creating, editing, etc.
}

// Define messages for async operations
type tasksLoadedMsg struct {
	tasks []models.Task
	err   error
}

// loadTasksCmd is a command that fetches tasks from the service.
func (m Model) loadTasksCmd() tea.Msg {
	tasks, err := m.taskService.ListAllActiveTasks(context.Background()) // Or another appropriate listing method
	return tasksLoadedMsg{tasks: tasks, err: err}
}

// New creates a new task management model.
// It requires a TaskService to interact with the backend.
func New(taskService service.TaskService) Model {
	// _ = time.Now() // Diagnostic to ensure 'time' package is seen as used.
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Título", Width: 30},
		// {Title: "Descrição", Width: 40}, // Description can be long, maybe show in a detail view
		{Title: "Prazo", Width: 10},
		// {Title: "Concluída", Width: 10}, // Tasks listed are active, so always false.
		{Title: "ID Turma", Width: 8},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10), // Initial height, will be adjusted by WindowSizeMsg
	)

	// Setup styles for the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// isLoading is true by default as Init will be called to load tasks.
	return Model{
		taskService: taskService,
		table:       t,
		isLoading:   true,
	}
}

// Init is called when the model becomes active. It starts the process of loading tasks.
func (m Model) Init() tea.Cmd {
	m.isLoading = true // Explicitly set loading to true
	m.err = nil        // Clear any previous errors
	return m.loadTasksCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { // Changed return type to tea.Model
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tasksLoadedMsg:
		m.isLoading = false // Finished loading
		m.err = msg.err     // Store error if any
		if msg.err != nil {
			return m, nil
		}
		rows := make([]table.Row, len(msg.tasks))
		for i, task := range msg.tasks {
			dueDate := "N/A"
			if task.DueDate != nil {
				dueDate = task.DueDate.Format("02/01/2006")
			}
			classID := "N/A" // Default for tasks not associated with a class (e.g., system bugs)
			if task.ClassID != nil && *task.ClassID != 0 { // Check if ClassID is non-nil and not zero
				classID = fmt.Sprintf("%d", *task.ClassID)
			}
			rows[i] = table.Row{
				fmt.Sprintf("%d", task.ID),
				task.Title,
				// task.Description, // If description column is re-added
				dueDate,
				// fmt.Sprintf("%t", task.IsCompleted), // Removed as we only show active tasks
				classID,
			}
		}
		m.table.SetRows(rows)
		return m, nil // No further command after loading

	case tea.KeyMsg:
		// For now, only handle table navigation.
		// Later, add keys for add, edit, delete, complete.
		// 'q' or 'esc' to go back will be handled by the main app model.
	case error: // This handles errors passed as messages
		m.err = msg
		return m, nil
	}

	m.table, cmd = m.table.Update(msg) // This updates the table's internal state (like cursor)
	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		// Basic error display for now
		return "Erro ao carregar tarefas: " + m.err.Error()
	}
	if m.isLoading {
		return "Carregando tarefas..." // Simple loading message
	}
	return baseStyle.Render(m.table.View()) + "\n"
}

// SetSize allows the main app model to adjust the size of this component.
func (m *Model) SetSize(width, height int) {
	// height is the total available height for the tasks.Model's View.
	// tasks.Model.View() is baseStyle.Render(m.table.View()) + "\n".
	// baseStyle border takes 2 vertical lines. Trailing \n takes 1 line.
	// So, height available for m.table.View() within the border is height - 2 (border) - 1 (\n) = height - 3.
	// This available height must accommodate the table's header (1 line) + its viewport.
	// So, table viewport height = (height - 3) - 1 = height - 4.
	// However, testing reveals m.table.Height() is consistently 2 less than value passed to m.table.SetHeight().
	// To achieve a viewport of (height - 4), we pass (height - 4) + 2 = height - 2.
	m.table.SetWidth(width - baseStyle.GetHorizontalFrameSize())
	m.table.SetHeight(height - 2)
}

// IsFocused returns false for now, indicating that the task view itself isn't managing focus
// in a way that would prevent the main app from handling 'esc'/'q' for navigation.
// This can be expanded if the tasks view gets its own input fields or modals.
func (m Model) IsFocused() bool {
	return false
}
