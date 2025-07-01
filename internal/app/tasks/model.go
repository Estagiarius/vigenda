package tasks

import (
	"context" // Required for service calls
	"fmt"     // For formatting data into table rows
	"strconv" // For parsing class ID
	"strings" // For form view
	"time"    // For parsing due date

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"vigenda/internal/models"
	"vigenda/internal/service"
	// "vigenda/internal/tui" // If using shared prompt, but better to embed form logic here
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// FormState represents the current state of the task form (creating or editing).
type FormState int

const (
	NoForm FormState = iota
	CreatingTask
	EditingTask // Future use
	ViewingDetail
)

type Model struct {
	taskService service.TaskService
	table       table.Model
	isLoading   bool
	err         error
	formState   FormState // Used for Create/Edit forms and now ViewDetail state

	// Form fields (for creating/editing)
	// titleInput       textinput.Model // Will be inputs[0]
	// descriptionInput textinput.Model // Will be inputs[1]
	// dueDateInput     textinput.Model // Will be inputs[2]
	// classIDInput     textinput.Model // Will be inputs[3]
	inputs     []textinput.Model // Holds all form inputs
	focusIndex int

	selectedTaskForDetail *models.Task // Store the task whose details are being viewed
	editingTaskID         int64        // ID of the task being edited
	taskIDToDelete        int64        // ID of the task pending delete confirmation
	confirmingDelete      bool         // True if waiting for delete confirmation

	width  int // For layout
	height int // For layout
}

// Define messages for async operations
type tasksLoadedMsg struct {
	tasks []models.Task
	err   error
}

// fetchedTaskDetailMsg is sent when a single task's details are fetched (for viewing or editing).
type fetchedTaskDetailMsg struct {
	task    *models.Task
	err     error
	forEdit bool // Indicates if the fetch was for editing
}

// taskUpdatedMsg is sent when a task is successfully updated.
type taskUpdatedMsg struct{}

// taskUpdateFailedMsg is sent when task update fails.
type taskUpdateFailedMsg struct{ err error }

// taskDeletedMsg is sent when a task is successfully deleted.
type taskDeletedMsg struct{} // Can be empty, success is implied

// taskDeleteFailedMsg is sent when task deletion fails.
type taskDeleteFailedMsg struct{ err error }

// taskMarkedCompletedMsg is sent when a task is successfully marked as completed.
type taskMarkedCompletedMsg struct{} // Can include taskID if needed for specific UI updates

// taskMarkCompleteFailedMsg is sent when marking a task as completed fails.
type taskMarkCompleteFailedMsg struct{ err error }

// loadTasksCmd is a command that fetches tasks from the service.
func (m *Model) loadTasksCmd() tea.Msg {
	tasks, err := m.taskService.ListAllActiveTasks(context.Background())
	return tasksLoadedMsg{tasks: tasks, err: err}
}

// fetchTaskForDetailCmd fetches a single task by ID for detail view or editing.
func (m *Model) fetchTaskForDetailCmd(taskID int64, forEditing bool) tea.Cmd {
	return func() tea.Msg {
		task, err := m.taskService.GetTaskByID(context.Background(), taskID)
		return fetchedTaskDetailMsg{task: task, err: err, forEdit: forEditing}
	}
}

// updateTaskCmd updates an existing task.
func (m *Model) updateTaskCmd(taskToUpdate *models.Task) tea.Cmd {
	return func() tea.Msg {
		err := m.taskService.UpdateTask(context.Background(), taskToUpdate)
		if err != nil {
			return taskUpdateFailedMsg{err}
		}
		return taskUpdatedMsg{}
	}
}

// deleteTaskCmd deletes a task by its ID.
func (m *Model) deleteTaskCmd(taskID int64) tea.Cmd {
	return func() tea.Msg {
		err := m.taskService.DeleteTask(context.Background(), taskID)
		if err != nil {
			return taskDeleteFailedMsg{err}
		}
		return taskDeletedMsg{}
	}
}

// markTaskCompleteCmd marks a task as completed.
func (m *Model) markTaskCompleteCmd(taskID int64) tea.Cmd {
	return func() tea.Msg {
		err := m.taskService.MarkTaskAsCompleted(context.Background(), taskID)
		if err != nil {
			return taskMarkCompleteFailedMsg{err}
		}
		return taskMarkedCompletedMsg{}
	}
}

// New creates a new task management model.
func New(taskService service.TaskService) *Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Título", Width: 30},
		{Title: "Prazo", Width: 10},
		{Title: "ID Turma", Width: 8},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

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

	// Form inputs
	ti := textinput.New()
	ti.Placeholder = "Título da Tarefa"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Prompt = "Título: "

	di := textinput.New()
	di.Placeholder = "Descrição (opcional)"
	di.CharLimit = 256
	di.Width = 50
	di.Prompt = "Descrição: "

	ddi := textinput.New()
	ddi.Placeholder = "DD/MM/YYYY (opcional)" // Changed placeholder
	ddi.CharLimit = 10
	ddi.Width = 20
	ddi.Prompt = "Prazo: "

	ci := textinput.New()
	ci.Placeholder = "ID da Turma (opcional, numérico)"
	ci.CharLimit = 10
	ci.Width = 20
	ci.Prompt = "ID Turma: "

	inputs := make([]textinput.Model, 4)
	inputs[0] = ti
	inputs[1] = di
	inputs[2] = ddi
	inputs[3] = ci

	return &Model{
		taskService:           taskService,
		table:                 t,
		isLoading:             true,
		formState:             NoForm,
		inputs:                inputs,
		focusIndex:            0,
		selectedTaskForDetail: nil,
		editingTaskID:         0,
		taskIDToDelete:        0,
		confirmingDelete:      false,
	}
}

// Init is called when the model becomes active.
func (m *Model) Init() tea.Cmd {
	m.isLoading = true
	m.err = nil
	m.formState = NoForm // Reset form, detail view, and edit state
	m.selectedTaskForDetail = nil
	m.editingTaskID = 0
	m.taskIDToDelete = 0
	m.confirmingDelete = false
	// m.resetFormInputs() // resetFormInputs will be called when entering form state
	return m.loadTasksCmd
}

// taskCreatedMsg is sent when a task is successfully created.
type taskCreatedMsg struct{ task models.Task }

// taskCreationFailedMsg is sent when task creation fails.
type taskCreationFailedMsg struct{ err error }

// createTaskCmd creates a new task.
func (m *Model) createTaskCmd(title, description string, classID *int64, dueDate *time.Time) tea.Cmd {
	return func() tea.Msg {
		task, err := m.taskService.CreateTask(context.Background(), title, description, classID, dueDate)
		if err != nil {
			return taskCreationFailedMsg{err}
		}
		return taskCreatedMsg{task}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tasksLoadedMsg:
		m.isLoading = false
		m.err = msg.err
		if msg.err == nil {
			rows := make([]table.Row, len(msg.tasks))
			for i, task := range msg.tasks {
				dueDate := "N/A"
				if task.DueDate != nil {
					dueDate = task.DueDate.Format("02/01/2006")
				}
				classIDStr := "N/A"
				if task.ClassID != nil && *task.ClassID != 0 {
					classIDStr = fmt.Sprintf("%d", *task.ClassID)
				}
				rows[i] = table.Row{fmt.Sprintf("%d", task.ID), task.Title, dueDate, classIDStr}
			}
			m.table.SetRows(rows)
		}
		return m, nil

	case taskCreatedMsg:
		m.formState = NoForm
		m.resetFormInputs()
		m.err = nil // Clear any previous form error
		return m, m.loadTasksCmd // Refresh list

	case taskCreationFailedMsg:
		m.err = msg.err // Display error, keep form open for correction
		// Don't reset form or change state, user might want to correct input
		return m, nil

	case taskUpdatedMsg:
		m.formState = NoForm
		m.resetFormInputs()
		m.editingTaskID = 0
		m.err = nil
		return m, m.loadTasksCmd // Refresh list

	case taskUpdateFailedMsg:
		m.err = msg.err // Display error, keep form open for correction
		return m, nil

	case taskDeletedMsg:
		m.isLoading = false
		m.confirmingDelete = false
		m.taskIDToDelete = 0
		m.err = nil
		return m, m.loadTasksCmd // Refresh list

	case taskDeleteFailedMsg:
		m.isLoading = false
		m.confirmingDelete = false // Keep confirmation context or clear based on UX choice
		m.err = msg.err
		return m, nil

	case taskMarkedCompletedMsg:
		m.isLoading = false
		m.err = nil
		// Potentially add a temporary success message here if desired
		// e.g., m.statusMessage = "Tarefa marcada como concluída!"
		// then clear statusMessage after a short delay or on next key press.
		return m, m.loadTasksCmd // Refresh list to remove it from active tasks

	case taskMarkCompleteFailedMsg:
		m.isLoading = false
		m.err = msg.err // Display error to the user
		return m, nil

	case tea.KeyMsg:
		if m.confirmingDelete {
			switch msg.String() {
			case "s", "S":
				m.isLoading = true
				m.err = nil
				cmds = append(cmds, m.deleteTaskCmd(m.taskIDToDelete))
				// confirmingDelete will be reset by taskDeletedMsg or taskDeleteFailedMsg
				return m, tea.Batch(cmds...)
			case "n", "N", "esc":
				m.confirmingDelete = false
				m.taskIDToDelete = 0
				m.err = nil
				return m, nil
			}
		} else if m.formState == CreatingTask || m.formState == EditingTask {
			switch msg.String() {
			case "ctrl+c", "esc":
				m.formState = NoForm
				m.resetFormInputs()
				m.editingTaskID = 0
				m.err = nil
				return m, nil
			case "enter":
				if m.focusIndex == len(m.inputs)-1 { // Last input, submit
					title := m.inputs[0].Value() // Title from inputs[0]
					if title == "" {
						m.err = fmt.Errorf("título não pode ser vazio")
						return m, nil
					}
					description := m.inputs[1].Value() // Description from inputs[1]
					var classID *int64
					if m.inputs[3].Value() != "" { // ClassID from inputs[3]
						cid, err := strconv.ParseInt(m.inputs[3].Value(), 10, 64)
						if err != nil {
							m.err = fmt.Errorf("ID da turma inválido: %v", err)
							return m, nil
						}
						classID = &cid
					}
					var dueDate *time.Time
					if m.inputs[2].Value() != "" { // DueDate from inputs[2]
						parsedDate, err := time.Parse("02/01/2006", m.inputs[2].Value()) // Changed format string
						if err != nil {
							m.err = fmt.Errorf("formato de data inválido (use DD/MM/YYYY): %v", err) // Changed error message
							return m, nil
						}
						dueDate = &parsedDate
					}

					m.isLoading = true
					m.err = nil

					var submitCmd tea.Cmd
					if m.formState == CreatingTask {
						submitCmd = m.createTaskCmd(title, description, classID, dueDate)
					} else if m.formState == EditingTask {
						updatedTask := &models.Task{
							ID:          m.editingTaskID,
							UserID:      m.selectedTaskForDetail.UserID,
							Title:       title,
							Description: description,
							ClassID:     classID,
							DueDate:     dueDate,
							IsCompleted: m.selectedTaskForDetail.IsCompleted,
						}
						submitCmd = m.updateTaskCmd(updatedTask)
					}
					return m, submitCmd // Return immediately with the submission command
				} else { // Enter on a field that is not the last one: navigate to next field
					m.nextInput()
					// The subsequent input update loop will handle the Enter key for the newly focused input,
					// which is usually a no-op for value change but might affect cursor/blink.
				}
			case "tab", "shift+tab", "up", "down":
				if msg.String() == "up" || msg.String() == "shift+tab" {
					m.prevInput()
				} else {
					m.nextInput()
				}
			}
			for i := range m.inputs {
				if i == m.focusIndex {m.inputs[i].Focus()} else {m.inputs[i].Blur()}
			}
			newInputs := make([]textinput.Model, len(m.inputs))
			for i := range m.inputs {
				newInputs[i], cmd = m.inputs[i].Update(msg)
				cmds = append(cmds, cmd)
			}
			m.inputs = newInputs
			return m, tea.Batch(cmds...)

		} else if m.formState == ViewingDetail {
			switch msg.String() {
			case "ctrl+c", "esc", "q":
				m.formState = NoForm
				m.selectedTaskForDetail = nil
				m.err = nil
				return m, nil
			}
		} else { // NoForm state (table view)
			switch msg.String() {
			case "a":
				m.formState = CreatingTask
				m.resetFormInputs()
				m.focusIndex = 0
				if len(m.inputs) > 0 {m.inputs[m.focusIndex].Focus()}
				m.err = nil
				return m, textinput.Blink
			case "e": // Edit selected task
				if len(m.table.Rows()) > 0 && m.table.Cursor() >= 0 && m.table.Cursor() < len(m.table.Rows()) {
					selectedRow := m.table.SelectedRow()
					taskIDStr := selectedRow[0]
					taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
					if err != nil {
						m.err = fmt.Errorf("erro ao parsear ID da tarefa para editar: %v", err)
						return m, nil
					}
					cmds = append(cmds, m.fetchTaskForDetailCmd(taskID, true)) // Fetch for editing
					m.isLoading = true
				}
			case "c": // Mark task as completed
				if len(m.table.Rows()) > 0 && m.table.Cursor() >= 0 && m.table.Cursor() < len(m.table.Rows()) {
					selectedRow := m.table.SelectedRow()
					taskIDStr := selectedRow[0] // Assuming ID is the first column
					taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
					if err != nil {
						m.err = fmt.Errorf("erro ao parsear ID da tarefa para completar: %v", err)
					} else {
						m.isLoading = true
						m.err = nil
						cmds = append(cmds, m.markTaskCompleteCmd(taskID))
					}
				}
			case "d": // Delete selected task
				if len(m.table.Rows()) > 0 && m.table.Cursor() >= 0 && m.table.Cursor() < len(m.table.Rows()) {
					selectedRow := m.table.SelectedRow()
					taskIDStr := selectedRow[0] // Assuming ID is the first column
					taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
					if err != nil {
						m.err = fmt.Errorf("erro ao parsear ID da tarefa para excluir: %v", err)
					} else {
						m.taskIDToDelete = taskID
						m.confirmingDelete = true
						m.err = nil // Clear previous errors before showing confirmation
					}
				}
			case "v", "enter":
				if len(m.table.Rows()) > 0 && m.table.Cursor() >= 0 && m.table.Cursor() < len(m.table.Rows()) {
					selectedRow := m.table.SelectedRow()
					taskIDStr := selectedRow[0]
					taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
					if err != nil {
						m.err = fmt.Errorf("erro ao parsear ID da tarefa selecionada: %v", err)
						return m, nil
					}
					cmds = append(cmds, m.fetchTaskForDetailCmd(taskID, false)) // Fetch for viewing
					m.isLoading = true
				}
			}
		}

	case fetchedTaskDetailMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
			m.formState = NoForm
		} else {
			m.selectedTaskForDetail = msg.task // Store for reference (e.g., UserID, IsCompleted)
			if msg.forEdit {
				m.editingTaskID = msg.task.ID
				m.inputs[0].SetValue(msg.task.Title)             // Title
				m.inputs[1].SetValue(msg.task.Description)       // Description
				if msg.task.DueDate != nil {
					m.inputs[2].SetValue(msg.task.DueDate.Format("02/01/2006")) // DueDate - Changed format
				} else {
					m.inputs[2].SetValue("")
				}
				if msg.task.ClassID != nil {
					m.inputs[3].SetValue(strconv.FormatInt(*msg.task.ClassID, 10)) // ClassID
				} else {
					m.inputs[3].SetValue("")
				}
				m.formState = EditingTask
				m.focusIndex = 0 // Start focus on the first field
				if len(m.inputs) > 0 { m.inputs[0].Focus() }
				m.err = nil
				cmds = append(cmds, textinput.Blink)
			} else { // For viewing
				m.formState = ViewingDetail
				m.err = nil
			}
		}
		return m, tea.Batch(cmds...)


	case error:
		m.isLoading = false
		m.err = msg
		return m, nil
	}

	if m.formState == NoForm {
		var updatedTable table.Model
		updatedTable, cmd = m.table.Update(msg) // This handles table navigation (up/down keys)
		m.table = updatedTable
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

// View renders the task management UI.
func (m *Model) View() string {
	if m.confirmingDelete {
		if m.isLoading { // Loading during delete operation itself
			return fmt.Sprintf("Excluindo tarefa ID %d...", m.taskIDToDelete)
		}
		if m.err != nil { // Error occurred during delete attempt
			return fmt.Sprintf("Erro ao excluir tarefa ID %d: %v\n\nPressione Esc para voltar à lista.", m.taskIDToDelete, m.err)
		}
		return fmt.Sprintf("Tem certeza que deseja excluir a tarefa ID %d? (s/n)", m.taskIDToDelete)
	}

	if m.err != nil && m.formState != CreatingTask && m.formState != ViewingDetail && m.formState != EditingTask {
		return fmt.Sprintf("Erro: %v\n\nPressione 'a' para adicionar, 'e' para editar, 'd' para excluir, 'v' para ver detalhes, 'esc' para sair desta tela.", m.err)
	}

	switch m.formState {
	case CreatingTask, EditingTask:
		return m.viewForm()
	case ViewingDetail:
		if m.isLoading {
			return "Carregando detalhes da tarefa..."
		}
		if m.err != nil {
			return fmt.Sprintf("Erro ao ver detalhes: %v\n\nPressione 'esc' para voltar.", m.err)
		}
		return m.viewTaskDetail()
	default: // NoForm (Table view)
		if m.isLoading {
			return "Carregando tarefas..."
		}
		var help strings.Builder
		help.WriteString("\n\n")
		help.WriteString("  'a': Adicionar Nova Tarefa\n")
		help.WriteString("  'e': Editar Tarefa Selecionada\n")
		help.WriteString("  'd': Excluir Tarefa Selecionada\n")
		help.WriteString("  'c': Concluir Tarefa Selecionada\n")
		help.WriteString("  'v' ou Enter: Ver Detalhes da Tarefa")
		// 'esc' to go back is handled by the main app model.
		return baseStyle.Render(m.table.View()) + help.String()
	}
}

// viewForm renders the task creation/editing form.
func (m *Model) viewForm() string {
	var b strings.Builder
	formTitle := "Nova Tarefa"
	if m.formState == EditingTask {
		formTitle = fmt.Sprintf("Editando Tarefa (ID: %d)", m.editingTaskID)
	}
	b.WriteString(fmt.Sprintf("%s (Pressione Enter para avançar, Esc para cancelar)\n\n", formTitle))

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if m.focusIndex == i {
			b.WriteString(" <") // Indicator for focused field
		}
		b.WriteString("\n")
	}
	b.WriteString("\nPressione Enter no último campo para Salvar.")
	if m.err != nil { // Display form-specific errors
		b.WriteString(fmt.Sprintf("\n\nErro: %v", m.err))
	}
	return baseStyle.Render(b.String())
}

// viewTaskDetail renders the detailed view of a selected task.
func (m *Model) viewTaskDetail() string {
	if m.selectedTaskForDetail == nil {
		return "Nenhuma tarefa selecionada para ver detalhes.\nPressione Esc para voltar."
	}

	var b strings.Builder
	task := m.selectedTaskForDetail

	b.WriteString(fmt.Sprintf("Detalhes da Tarefa (ID: %d)\n", task.ID))
	b.WriteString(strings.Repeat("-", 30) + "\n") // Separator

	b.WriteString(fmt.Sprintf("Título: %s\n", task.Title))
	b.WriteString(fmt.Sprintf("Descrição: %s\n", task.Description)) // Display full description

	dueDateStr := "N/A"
	if task.DueDate != nil {
		dueDateStr = task.DueDate.Format("02/01/2006")
	}
	b.WriteString(fmt.Sprintf("Prazo: %s\n", dueDateStr))

	classIDStr := "N/A (Tarefa Geral)"
	if task.ClassID != nil && *task.ClassID != 0 {
		classIDStr = fmt.Sprintf("%d", *task.ClassID)
	}
	b.WriteString(fmt.Sprintf("ID Turma: %s\n", classIDStr))
	b.WriteString(fmt.Sprintf("Concluída: %t\n", task.IsCompleted)) // Show completion status

	b.WriteString(strings.Repeat("-", 30) + "\n")
	b.WriteString("\nPressione Esc para voltar à lista.")

	// Apply some styling. The content can be wrapped if too long.
	// For simplicity, direct rendering. For complex layouts, consider lipgloss.JoinVertical/Horizontal.
	detailStyle := lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("63"))
	contentWidth := m.width - detailStyle.GetHorizontalFrameSize() - 4 // Extra padding for content

	// Word wrap for description if it's too long
	// This is a basic implementation. For more robust wrapping, a library might be better.
	descLines := strings.Split(task.Description, "\n")
	wrappedDesc := ""
	for _, line := range descLines {
		if contentWidth > 0 && len(line) > contentWidth {
			// Simple wrap logic
			for len(line) > contentWidth {
				wrappedDesc += line[:contentWidth] + "\n"
				line = line[contentWidth:]
			}
		}
		wrappedDesc += line + "\n"
	}


	finalView := fmt.Sprintf("Detalhes da Tarefa (ID: %d)\n%s\n", task.ID, strings.Repeat("-", 30))
	finalView += fmt.Sprintf("Título: %s\n", task.Title)
	finalView += fmt.Sprintf("Descrição:\n%s\n", strings.TrimSpace(wrappedDesc))
	finalView += fmt.Sprintf("Prazo: %s\n", dueDateStr)
	finalView += fmt.Sprintf("ID Turma: %s\n", classIDStr)
	finalView += fmt.Sprintf("Concluída: %t\n", task.IsCompleted)
	finalView += fmt.Sprintf("%s\n\nPressione Esc para voltar à lista.", strings.Repeat("-", 30))


	return detailStyle.Render(finalView)
}


// nextInput moves focus to the next text input field
func (m *Model) nextInput() {
	m.inputs[m.focusIndex].Blur()
	m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
	m.inputs[m.focusIndex].Focus()
}

// prevInput moves focus to the previous text input field
func (m *Model) prevInput() {
	m.inputs[m.focusIndex].Blur()
	m.focusIndex--
	if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs) - 1
	}
	m.inputs[m.focusIndex].Focus()
}

// resetFormInputs clears all input fields and resets focus.
func (m *Model) resetFormInputs() {
	if len(m.inputs) < 4 {
		// This case should ideally not happen if New() initializes inputs correctly.
		// Log an error or handle appropriately if it does.
		// For now, let's ensure it doesn't panic.
		// This might indicate a need to re-initialize inputs if they are ever nil or too short.
		// However, New() should prevent this.
		return
	}
	m.inputs[0].Reset() // Title
	m.inputs[1].Reset() // Description
	m.inputs[2].Reset() // DueDate
	m.inputs[3].Reset() // ClassID
	// No need to reassign m.inputs if they are just being reset.
	// m.inputs = []textinput.Model{m.titleInput, m.descriptionInput, m.dueDateInput, m.classIDInput}
	m.focusIndex = 0
	if len(m.inputs) > 0 { // Ensure inputs slice is not empty before focusing
		m.inputs[0].Focus()
	}
	m.err = nil
}


// SetSize allows the main app model to adjust the size of this component.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	// Adjust table size
	// The baseStyle adds 2 for horizontal and 2 for vertical borders.
	// The trailing newline in View() for table takes 1 line.
	// The help text ("Pressione 'a'...") takes about 2 lines.
	// Effective table height = height - borderV - helpTextLines - viewNewline
	tableWidth := width - baseStyle.GetHorizontalFrameSize()
	tableHeight := height - baseStyle.GetVerticalFrameSize() - 3 // Approximation for help text + newline
	if tableHeight < 1 { tableHeight = 1 } // Minimum height for table
	m.table.SetWidth(tableWidth)
	m.table.SetHeight(tableHeight)


	// Adjust form input sizes (can be fixed or relative to width)
	inputWidth := width - baseStyle.GetHorizontalFrameSize() - 10 // Some padding
	if inputWidth < 20 { inputWidth = 20 }
	for i := range m.inputs {
		m.inputs[i].Width = inputWidth
	}
}

// IsFocused returns true if the model is currently in a form input state,
// which might mean it wants to trap 'esc' locally.
func (m *Model) IsFocused() bool {
	return m.formState != NoForm
}

// IsLoading returns true if the model is currently loading data (e.g. tasks list)
func (m *Model) IsLoading() bool {
    return m.isLoading
}
