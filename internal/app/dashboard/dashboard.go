package dashboard

import (
	"context"
	"fmt"
	"strings"
	"time"
	"vigenda/internal/models"
	"vigenda/internal/service"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// General Styles
	docStyle = lipgloss.NewStyle().Padding(1, 2) // Overall padding for the dashboard view
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).MarginBottom(1).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).PaddingBottom(1)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true) // Red for errors
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(1) // Gray for help text
	loadingStyle = lipgloss.NewStyle().MarginTop(1).MarginBottom(1)

	// Section Styles
	sectionTitleStyle = lipgloss.NewStyle().Bold(true).MarginBottom(1).Foreground(lipgloss.Color("220")) // Yellowish for section titles
	sectionBoxStyle   = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238")). // Lighter gray border
				Padding(1).
				MarginBottom(1)

	// Item Styles
	taskStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("104")) // Light purple for tasks
	taskDoneStyle    = lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("240")) // Gray and strikethrough for done tasks
	lessonStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))  // Green for lessons
	notificationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange for notifications

	// Layout styles
	columnStyle = lipgloss.NewStyle().Padding(0, 1) // Padding between columns
)

// Model represents the state of the dashboard TUI component.
// It holds data like upcoming tasks and today's lessons, and manages
// their fetching and display.
type Model struct {
	taskService    service.TaskService // Service for fetching task-related data.
	classService   service.ClassService  // Service for fetching class and lesson-related data.
	// Add other services as needed

	upcomingTasks []models.Task
	todaysLessons []models.Lesson // Assuming models.Lesson exists or will be created
	// notifications []string // Placeholder for notifications

	isLoading bool
	spinner   spinner.Model
	width     int
	height    int
	err       error                 // Stores any error encountered during data fetching or processing.
}

// New creates a new dashboard model, initializing it with necessary services
// and setting its initial state to loading.
func New(ts service.TaskService, cs service.ClassService) *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot // Sets the visual style of the spinner.
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &Model{
		taskService:  ts,
		classService: cs,
		isLoading:    true,
		spinner:      s,
	}
}

// Init loads the initial data for the dashboard.
func (m *Model) Init() tea.Cmd {
	m.isLoading = true
	cmds := []tea.Cmd{m.spinner.Tick}
	// Assuming a default user ID or a way to get current user ID
	// For now, let's assume userID 1 for fetching user-specific data like tasks.
	// Lessons might be tied to a class, which is tied to a user/subject.
	// This part might need adjustment based on how user context is handled.
	cmds = append(cmds, m.fetchDashboardData(context.Background(), 1)) // Example userID = 1
	return tea.Batch(cmds...)
}

// dashboardDataLoadedMsg is a message sent when dashboard data (tasks, lessons)
// has been successfully fetched.
type dashboardDataLoadedMsg struct {
	tasks   []models.Task
	lessons []models.Lesson
	// notifications []string // Placeholder for future notification data
}

// dashboardErrorMsg is a message sent when an error occurs during
// dashboard data fetching.
type dashboardErrorMsg struct{ err error }

// Error implements the error interface for dashboardErrorMsg.
func (e dashboardErrorMsg) Error() string { return e.err.Error() }

// fetchDashboardData creates a tea.Cmd that, when run, fetches all necessary data
// for the dashboard, such as upcoming tasks and today's lessons.
// It returns either a dashboardDataLoadedMsg on success or a dashboardErrorMsg on failure.
// Note: The current implementation of upcoming tasks fetching is temporary and filters
// all active tasks manually. This should be replaced by a more specific service call.
func (m *Model) fetchDashboardData(ctx context.Context, userID int64) tea.Cmd {
	return func() tea.Msg {
		// Fetch upcoming tasks
		// TaskService needs a method like GetUpcomingTasksByUserID or GetUpcomingTasks
		// For now, let's assume a method that fetches tasks for a user that are not completed and ordered by due date.
		// Let's placeholder with a conceptual method name.
		// tasks, err := m.taskService.GetUpcomingTasks(ctx, userID, 5)

		// Using ListAllActiveTasks as a temporary measure, then filtering.
		// This is NOT ideal for performance and should be replaced with a dedicated service method.
		allTasks, err := m.taskService.ListAllActiveTasks(ctx) // Assuming this gets tasks for the "current" user or all if no user context in service
		if err != nil {
			return dashboardErrorMsg{err: fmt.Errorf("failed to fetch tasks: %w", err)}
		}
		// Filter and limit tasks manually (TEMPORARY)
		var upcomingTasks []models.Task
		count := 0
		for _, task := range allTasks {
			if !task.IsCompleted && task.DueDate != nil && time.Now().Before(*task.DueDate) { // Only future tasks
				// This example doesn't filter by UserID directly as ListAllActiveTasks might not support it.
				// Proper implementation would handle UserID in the service/repository layer.
				if task.UserID == userID { // Manual filter by UserID
					upcomingTasks = append(upcomingTasks, task)
					count++
					if count >= 5 {
						break
					}
				}
			}
		}


		// Fetch today's lessons (assuming models.Lesson and service method exist)
		// lessons, err := m.classService.GetTodaysLessons(ctx, userID)
		// if err != nil {
		// 	return dashboardErrorMsg{fmt.Errorf("failed to fetch lessons: %w", err)}
		// }
		// Placeholder for lessons as models.Lesson might not exist or be fully defined yet.
		// var todaysLessons []models.Lesson = []models.Lesson{} // Empty for now

		// Actually fetch today's lessons
		todaysLessons, err := m.classService.GetTodaysLessons(ctx, userID)
		if err != nil {
			// Return specific error for lesson fetching failure
			return dashboardErrorMsg{err: fmt.Errorf("failed to fetch lessons: %w", err)}
		}

		// Fetch notifications (placeholder)
		// var notifications []string = []string{"Bem-vindo ao Vigenda!"}

		return dashboardDataLoadedMsg{
			tasks:   upcomingTasks,
			lessons: todaysLessons,
			// notifications: notifications,
		}
	}
}

// Update handles messages for the dashboard model. It processes incoming messages
// (like data loaded, errors, or window size changes) and updates the model state accordingly.
// It returns the updated model and any command to be executed.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Adjust layout based on size if needed

	case spinner.TickMsg:
		if m.isLoading {
			var spinCmd tea.Cmd
			m.spinner, spinCmd = m.spinner.Update(msg)
			cmds = append(cmds, spinCmd)
		}

	case dashboardDataLoadedMsg:
		m.isLoading = false
		m.upcomingTasks = msg.tasks
		m.todaysLessons = msg.lessons
		// m.notifications = msg.notifications

	case dashboardErrorMsg:
		m.isLoading = false
		m.err = msg.err
		// Log error or display it prominently

	// No specific key messages for dashboard navigation yet, assuming it's display-only for now.
	// If interactive elements are added, handle keys here.
	}

	return m, tea.Batch(cmds...)
}

// View renders the dashboard UI.
func (m *Model) View() string {
	if m.err != nil {
		return docStyle.Render(errorStyle.Render(fmt.Sprintf("Erro ao carregar o dashboard: %v", m.err)))
	}

	if m.isLoading {
		return docStyle.Render(loadingStyle.Render(fmt.Sprintf("%s Carregando dados do dashboard...", m.spinner.View())))
	}

	// Calculate available width for content, considering padding of docStyle
	contentWidth := m.width - docStyle.GetHorizontalPadding()

	// --- Upcoming Tasks Section ---
	var tasksContent strings.Builder
	tasksContent.WriteString(sectionTitleStyle.Render("Tarefas Próximas"))
	if len(m.upcomingTasks) == 0 {
		tasksContent.WriteString("\nNenhuma tarefa próxima encontrada.")
	} else {
		for _, task := range m.upcomingTasks {
			dueDateStr := "Sem prazo"
			if task.DueDate != nil {
				// Check if the task is overdue
				if time.Now().After(*task.DueDate) {
					dueDateStr = task.DueDate.Format("02/01 (Mon)") + " (Atrasada!)"
				} else {
					dueDateStr = task.DueDate.Format("02/01 (Mon)")
				}
			}
			taskDisplay := fmt.Sprintf("%s %s (%s)", taskStyle.Render("□"), task.Title, dueDateStr)
			tasksContent.WriteString("\n" + taskDisplay)
		}
	}
	tasksSectionRendered := sectionBoxStyle.Width(contentWidth).Render(tasksContent.String())

	// --- Today's Lessons Section ---
	var lessonsContent strings.Builder
	lessonsContent.WriteString(sectionTitleStyle.Render("Aulas de Hoje"))
	if len(m.todaysLessons) == 0 {
		lessonsContent.WriteString("\nNenhuma aula para hoje.")
	} else {
		for _, lesson := range m.todaysLessons {
			scheduledTimeStr := lesson.ScheduledAt.Format("15:04")
			lessonDisplay := fmt.Sprintf("%s %s (Turma %d) às %s", lessonStyle.Render("●"), lesson.Title, lesson.ClassID, scheduledTimeStr)
			lessonsContent.WriteString("\n" + lessonDisplay)
		}
	}
	lessonsSectionRendered := sectionBoxStyle.Width(contentWidth).Render(lessonsContent.String())

	// --- Notifications Section (Placeholder) ---
	// var notificationsContent strings.Builder
	// notificationsContent.WriteString(sectionTitleStyle.Render("Notificações"))
	// notificationsContent.WriteString("\n" + notificationStyle.Render("Bem-vindo ao Vigenda!"))
	// notificationsContent.WriteString("\n" + notificationStyle.Render("Nova versão disponível em breve."))
	// notificationsSectionRendered := sectionBoxStyle.Width(contentWidth).Render(notificationsContent.String())

	// --- Main Content Area ---
	// For now, stack sections vertically. Horizontal layout can be complex.
	mainContent := lipgloss.JoinVertical(lipgloss.Left,
		tasksSectionRendered,
		lessonsSectionRendered,
		// notificationsSectionRendered, // Uncomment when notifications are implemented
	)

	// --- Final Assembly ---
	// Using lipgloss.JoinVertical for overall structure including title and help
	fullView := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Width(contentWidth).Render("Dashboard Vigenda"),
		mainContent,
		helpStyle.Width(contentWidth).Render("Pressione 'esc' para (placeholder: sair), 'm' para (placeholder: menu)."),
	)

	return docStyle.Render(fullView)
}

// IsFocused can be used by the parent model to determine if this submodel
// is currently handling input (e.g., if it has a text input active).
// For a display-only dashboard, this might always return false.
func (m *Model) IsFocused() bool {
	return false // No input fields in this basic dashboard, so it never has primary focus for input.
}

// Placeholder for models.Lesson if it doesn't exist.
// Remove this if models.Lesson is properly defined in models/models.go
// type Lesson struct {
// 	Title       string
// 	ScheduledAt time.Time
// }
