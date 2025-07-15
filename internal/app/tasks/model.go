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
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var strikethroughStyle = lipgloss.NewStyle().Strikethrough(true)

// ViewState represents the current main view/state within the tasks model.
type ViewState int

const (
	TableView ViewState = iota
	FormView
	DetailView
	ConfirmDeleteView
)

// FormState represents the current state of the task form (creating or editing).
// This is used when currentView is FormView.
type FormState int

const (
	CreatingTask FormState = iota
	EditingTask
)

// FocusedTable indicates which table (pending or completed) has focus.
type FocusedTable int

const (
	PendingTableFocus FocusedTable = iota
	CompletedTableFocus
)

type Model struct {
	taskService         service.TaskService
	pendingTasksTable   table.Model
	completedTasksTable table.Model
	isLoading           bool
	err                 error
	currentView         ViewState
	formSubState        FormState    // If currentView is FormView, this specifies if creating or editing.
	focusedTable        FocusedTable

	inputs     []textinput.Model // Holds all form inputs
	focusIndex int

	selectedTaskForDetail *models.Task
	editingTaskID         int64
	taskIDToDelete        int64
	// confirmingDelete      bool         // This state is now handled by currentView = ConfirmDeleteView

	width  int
	height int
}

// Define messages for async operations
type tasksLoadedMsg struct {
	tasks []models.Task
	err   error
}
type fetchedTaskDetailMsg struct {
	task    *models.Task
	err     error
	forEdit bool
}
// taskCreatedMsg is sent when a task is successfully created.
type taskCreatedMsg struct{ task models.Task }

// taskCreationFailedMsg is sent when task creation fails.
type taskCreationFailedMsg struct{ err error }

type taskUpdatedMsg struct{}
type taskUpdateFailedMsg struct{ err error }
type taskDeletedMsg struct{}
type taskDeleteFailedMsg struct{ err error }
type taskMarkedCompletedMsg struct{}
type taskMarkCompleteFailedMsg struct{ err error }

func (m *Model) loadTasksCmd() tea.Msg {
	tasks, err := m.taskService.ListAllTasks(context.Background())
	return tasksLoadedMsg{tasks: tasks, err: err}
}

func (m *Model) fetchTaskForDetailCmd(taskID int64, forEditing bool) tea.Cmd {
	return func() tea.Msg {
		task, err := m.taskService.GetTaskByID(context.Background(), taskID)
		return fetchedTaskDetailMsg{task: task, err: err, forEdit: forEditing}
	}
}

// createTaskCmd creates a new task.
func (m *Model) createTaskCmd(title, description string, classID *int64, dueDate *time.Time) tea.Cmd {
	return func() tea.Msg {
		task, err := m.taskService.CreateTask(context.Background(), title, description, classID, dueDate)
		if err != nil {
			return taskCreationFailedMsg{err: err}
		}
		return taskCreatedMsg{task: task}
	}
}

func (m *Model) updateTaskCmd(taskToUpdate *models.Task) tea.Cmd {
	return func() tea.Msg {
		err := m.taskService.UpdateTask(context.Background(), taskToUpdate)
		if err != nil {
			return taskUpdateFailedMsg{err}
		}
		return taskUpdatedMsg{}
	}
}

func (m *Model) deleteTaskCmd(taskID int64) tea.Cmd {
	return func() tea.Msg {
		err := m.taskService.DeleteTask(context.Background(), taskID)
		if err != nil {
			return taskDeleteFailedMsg{err}
		}
		return taskDeletedMsg{}
	}
}

func (m *Model) markTaskCompleteCmd(taskID int64) tea.Cmd {
	return func() tea.Msg {
		err := m.taskService.MarkTaskAsCompleted(context.Background(), taskID)
		if err != nil {
			return taskMarkCompleteFailedMsg{err}
		}
		return taskMarkedCompletedMsg{}
	}
}

func New(taskService service.TaskService) *Model {
	pendingColumns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Título", Width: 30},
		{Title: "Prazo", Width: 10},
		{Title: "ID Turma", Width: 8},
	}
	pendingTable := table.New(
		table.WithColumns(pendingColumns),
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
	pendingTable.SetStyles(s)

	completedColumns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Título", Width: 30},
		{Title: "Prazo", Width: 10},
		{Title: "ID Turma", Width: 8},
	}
	completedTable := table.New(
		table.WithColumns(completedColumns),
		table.WithFocused(false),
		table.WithHeight(5),
	)
	completedTable.SetStyles(s)

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
	ddi.Placeholder = "DD/MM/YYYY (opcional)"
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
		pendingTasksTable:     pendingTable,
		completedTasksTable:   completedTable,
		isLoading:             true,
		currentView:           TableView,
		formSubState:          CreatingTask,
		focusedTable:          PendingTableFocus,
		inputs:                inputs,
		focusIndex:            0,
		selectedTaskForDetail: nil,
		editingTaskID:         0,
		taskIDToDelete:        0,
	}
}

func (m *Model) Init() tea.Cmd {
	m.isLoading = true
	m.err = nil
	m.currentView = TableView
	m.focusedTable = PendingTableFocus
	m.selectedTaskForDetail = nil
	m.editingTaskID = 0
	m.taskIDToDelete = 0
	return m.loadTasksCmd
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tasksLoadedMsg:
		m.isLoading = false
		m.err = msg.err
		if msg.err == nil {
			pendingRows := []table.Row{}
			completedRows := []table.Row{}
			for _, task := range msg.tasks {
				dueDate := "N/A"
				if task.DueDate != nil {
					dueDate = task.DueDate.Format("02/01/2006")
				}
				classIDStr := "N/A"
				if task.ClassID != nil && *task.ClassID != 0 {
					classIDStr = fmt.Sprintf("%d", *task.ClassID)
				}

				titleCell := task.Title
				if task.IsCompleted {
					titleCell = strikethroughStyle.Render(task.Title)
				}
				row := table.Row{fmt.Sprintf("%d", task.ID), titleCell, dueDate, classIDStr}

				if task.IsCompleted {
					completedRows = append(completedRows, row)
				} else {
					pendingRows = append(pendingRows, row)
				}
			}
			m.pendingTasksTable.SetRows(pendingRows)
			m.completedTasksTable.SetRows(completedRows)
		} else {
			m.pendingTasksTable.SetRows([]table.Row{})
			m.completedTasksTable.SetRows([]table.Row{})
		}
		return m, nil

	case taskCreatedMsg:
		m.currentView = TableView
		m.resetFormInputs()
		m.err = nil
		return m, m.loadTasksCmd

	case taskCreationFailedMsg:
		m.err = msg.err
		return m, nil

	case taskUpdatedMsg:
		m.currentView = TableView
		m.resetFormInputs()
		m.editingTaskID = 0
		m.err = nil
		return m, m.loadTasksCmd

	case taskUpdateFailedMsg:
		m.err = msg.err
		return m, nil

	case taskDeletedMsg:
		m.isLoading = false
		m.currentView = TableView
		m.taskIDToDelete = 0
		m.err = nil
		return m, m.loadTasksCmd

	case taskDeleteFailedMsg:
		m.isLoading = false
		m.err = msg.err
		return m, nil

	case taskMarkedCompletedMsg:
		m.isLoading = false
		m.err = nil
		return m, m.loadTasksCmd

	case taskMarkCompleteFailedMsg:
		m.isLoading = false
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch m.currentView {
		case ConfirmDeleteView:
			switch msg.String() {
			case "s", "S":
				m.isLoading = true
				m.err = nil
				cmds = append(cmds, m.deleteTaskCmd(m.taskIDToDelete))
				return m, tea.Batch(cmds...)
			case "n", "N", "esc":
				m.currentView = TableView
				m.taskIDToDelete = 0
				m.err = nil
				return m, nil
			}
		case FormView:
			switch msg.String() {
			case "ctrl+c", "esc":
				m.currentView = TableView
				m.resetFormInputs()
				m.editingTaskID = 0
				m.err = nil
				return m, nil
			case "enter":
				if m.focusIndex == len(m.inputs)-1 {
					title := m.inputs[0].Value()
					if title == "" {
						m.err = fmt.Errorf("título não pode ser vazio")
						return m, nil
					}
					description := m.inputs[1].Value()
					var classID *int64
					if m.inputs[3].Value() != "" {
						cid, errConv := strconv.ParseInt(m.inputs[3].Value(), 10, 64)
						if errConv != nil {
							m.err = fmt.Errorf("ID da turma inválido: %v", errConv)
							return m, nil
						}
						classID = &cid
					}
					var dueDate *time.Time
					if m.inputs[2].Value() != "" {
						parsedDate, errConv := time.Parse("02/01/2006", m.inputs[2].Value())
						if errConv != nil {
							m.err = fmt.Errorf("formato de data inválido (use DD/MM/YYYY): %v", errConv)
							return m, nil
						}
						dueDate = &parsedDate
					}

					m.isLoading = true
					m.err = nil

					var submitCmd tea.Cmd
					if m.formSubState == CreatingTask {
						submitCmd = m.createTaskCmd(title, description, classID, dueDate)
					} else if m.formSubState == EditingTask {
						if m.selectedTaskForDetail == nil {
							m.err = fmt.Errorf("erro interno: dados da tarefa original não encontrados para edição")
							m.isLoading = false
							return m, nil
						}
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
					return m, submitCmd
				} else {
					m.nextInput()
				}
			case "tab", "shift+tab", "up", "down":
				if msg.String() == "up" || msg.String() == "shift+tab" {
					m.prevInput()
				} else {
					m.nextInput()
				}
			}
			tmpCmds := []tea.Cmd{}
			for i := range m.inputs {
				if i == m.focusIndex {m.inputs[i].Focus()} else {m.inputs[i].Blur()}
				var updatedInput textinput.Model
				updatedInput, cmd = m.inputs[i].Update(msg)
				m.inputs[i] = updatedInput
				tmpCmds = append(tmpCmds, cmd)
			}
			cmds = append(cmds, tmpCmds...)
			return m, tea.Batch(cmds...)

		case DetailView:
			switch msg.String() {
			case "ctrl+c", "esc", "q":
				m.currentView = TableView
				m.selectedTaskForDetail = nil
				m.err = nil
				return m, nil
			}
		case TableView:
			activeTable := &m.pendingTasksTable
			if m.focusedTable == CompletedTableFocus {
				activeTable = &m.completedTasksTable
			}

			switch msg.String() {
			case "a":
				m.currentView = FormView
				m.formSubState = CreatingTask
				m.resetFormInputs()
				m.err = nil
				return m, textinput.Blink
			case "e":
				if m.focusedTable == PendingTableFocus && len(activeTable.Rows()) > 0 && activeTable.Cursor() >= 0 && activeTable.Cursor() < len(activeTable.Rows()) {
					selectedRow := activeTable.SelectedRow()
					taskIDStr := selectedRow[0]
					taskID, errConv := strconv.ParseInt(taskIDStr, 10, 64)
					if errConv != nil {
						m.err = fmt.Errorf("erro ao parsear ID da tarefa para editar: %v", errConv)
						return m, nil
					}
					cmds = append(cmds, m.fetchTaskForDetailCmd(taskID, true))
					m.isLoading = true
				}
			case "c":
				if m.focusedTable == PendingTableFocus && len(m.pendingTasksTable.Rows()) > 0 && m.pendingTasksTable.Cursor() < len(m.pendingTasksTable.Rows()) {
					selectedRow := m.pendingTasksTable.SelectedRow()
					taskIDStr := selectedRow[0]
					taskID, errConv := strconv.ParseInt(taskIDStr, 10, 64)
					if errConv != nil {
						m.err = fmt.Errorf("erro ao parsear ID da tarefa para completar: %v", errConv)
					} else {
						m.isLoading = true
						m.err = nil
						cmds = append(cmds, m.markTaskCompleteCmd(taskID))
					}
				}
			case "d":
				if len(activeTable.Rows()) > 0 && activeTable.Cursor() >= 0 && activeTable.Cursor() < len(activeTable.Rows()) {
					selectedRow := activeTable.SelectedRow()
					taskIDStr := selectedRow[0]
					taskID, errConv := strconv.ParseInt(taskIDStr, 10, 64)
					if errConv != nil {
						m.err = fmt.Errorf("erro ao parsear ID da tarefa para excluir: %v", errConv)
					} else {
						m.taskIDToDelete = taskID
						m.currentView = ConfirmDeleteView
						m.err = nil
					}
				}
			case "v", "enter":
				if len(activeTable.Rows()) > 0 && activeTable.Cursor() >= 0 && activeTable.Cursor() < len(activeTable.Rows()) {
					selectedRow := activeTable.SelectedRow()
					taskIDStr := selectedRow[0]
					taskID, errConv := strconv.ParseInt(taskIDStr, 10, 64)
					if errConv != nil {
						m.err = fmt.Errorf("erro ao parsear ID da tarefa selecionada: %v", errConv)
						return m, nil
					}
					cmds = append(cmds, m.fetchTaskForDetailCmd(taskID, false))
					m.isLoading = true
				}
			case "tab":
				if m.focusedTable == PendingTableFocus {
					m.focusedTable = CompletedTableFocus
					m.pendingTasksTable.Blur()
					m.completedTasksTable.Focus()
				} else {
					m.focusedTable = PendingTableFocus
					m.completedTasksTable.Blur()
					m.pendingTasksTable.Focus()
				}
			case "up", "k", "down", "j":
				var updatedTableModel table.Model // Corrected type
				if m.focusedTable == PendingTableFocus {
					updatedTableModel, cmd = m.pendingTasksTable.Update(msg)
					m.pendingTasksTable = updatedTableModel // Corrected assignment
				} else {
					updatedTableModel, cmd = m.completedTasksTable.Update(msg)
					m.completedTasksTable = updatedTableModel // Corrected assignment
				}
				cmds = append(cmds, cmd)
			}
		}

	case fetchedTaskDetailMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
			m.currentView = TableView
		} else {
			if msg.task == nil {
				m.err = fmt.Errorf("detalhes da tarefa não encontrados, mas nenhum erro reportado")
				m.currentView = TableView
				return m, nil
			}
			m.selectedTaskForDetail = msg.task
			if msg.forEdit {
				m.editingTaskID = msg.task.ID
				m.inputs[0].SetValue(msg.task.Title)
				m.inputs[1].SetValue(msg.task.Description)
				if msg.task.DueDate != nil {
					m.inputs[2].SetValue(msg.task.DueDate.Format("02/01/2006"))
				} else {
					m.inputs[2].SetValue("")
				}
				if msg.task.ClassID != nil {
					m.inputs[3].SetValue(strconv.FormatInt(*msg.task.ClassID, 10))
				} else {
					m.inputs[3].SetValue("")
				}
				m.currentView = FormView
				m.formSubState = EditingTask
				m.focusIndex = 0
				if len(m.inputs) > 0 { m.inputs[0].Focus() }
				m.err = nil
				cmds = append(cmds, textinput.Blink)
			} else {
				m.currentView = DetailView
				m.err = nil
			}
		}
		return m, tea.Batch(cmds...)

	case error:
		m.isLoading = false
		m.err = msg
		return m, nil
	}

	if m.currentView == TableView {
		if keyMsg, ok := msg.(tea.KeyMsg); !ok || (ok && keyMsg.String() == "") { // Process non-string KeyMsgs or other Msgs for tables
			var pendCmd, compCmd tea.Cmd
			// Pass the original msg to both tables; they'll decide if it's relevant.
			// This handles WindowSizeMsg and other potential messages.
			m.pendingTasksTable, pendCmd = m.pendingTasksTable.Update(msg)
			cmds = append(cmds, pendCmd)

			m.completedTasksTable, compCmd = m.completedTasksTable.Update(msg)
			cmds = append(cmds, compCmd)
		}
	}
	return m, tea.Batch(cmds...)
}

// View renders the task management UI.
func (m *Model) View() string {
	if m.currentView == ConfirmDeleteView {
		if m.isLoading {
			return fmt.Sprintf("Excluindo tarefa ID %d...", m.taskIDToDelete)
		}
		if m.err != nil {
			return fmt.Sprintf("Erro ao excluir tarefa ID %d: %v\n\nPressione Esc para voltar à lista.", m.taskIDToDelete, m.err)
		}
		return fmt.Sprintf("Tem certeza que deseja excluir a tarefa ID %d? (s/n)", m.taskIDToDelete)
	}

	if m.err != nil && m.currentView != FormView && m.currentView != DetailView {
		return fmt.Sprintf("Erro: %v\n\nPressione 'a' para adicionar, 'e' para editar, 'd' para excluir, 'v' para ver detalhes, 'esc' para sair desta tela.", m.err)
	}

	switch m.currentView {
	case FormView:
		return m.viewForm()
	case DetailView:
		if m.isLoading {
			return "Carregando detalhes da tarefa..."
		}
		if m.err != nil {
			return fmt.Sprintf("Erro ao ver detalhes: %v\n\nPressione 'esc' para voltar.", m.err)
		}
		return m.viewTaskDetail()
	default: // TableView
		if m.isLoading {
			return "Carregando tarefas..."
		}
		pendingHeader := "Tarefas Pendentes"
		if m.focusedTable == PendingTableFocus {
			pendingHeader = lipgloss.NewStyle().Bold(true).SetString(pendingHeader).String()
		}
		completedHeader := "Tarefas Concluídas"
		if m.focusedTable == CompletedTableFocus {
			completedHeader = lipgloss.NewStyle().Bold(true).SetString(completedHeader).String()
		}

		tablesView := lipgloss.JoinVertical(lipgloss.Left,
			pendingHeader,
			baseStyle.Render(m.pendingTasksTable.View()),
			"\n"+completedHeader,
			baseStyle.Render(m.completedTasksTable.View()),
		)

		var help strings.Builder
		help.WriteString("\n\n")
		help.WriteString("  'a': Adicionar | 'e': Editar (pendentes) | 'd': Excluir | 'c': Concluir (pendentes)\n")
		help.WriteString("  'v'|Enter: Detalhes | Tab: Mudar Tabela Focada")
		return tablesView + help.String()
	}
}

// viewForm renders the task creation/editing form.
func (m *Model) viewForm() string {
	var b strings.Builder
	formTitle := "Nova Tarefa"
	if m.formSubState == EditingTask {
		formTitle = fmt.Sprintf("Editando Tarefa (ID: %d)", m.editingTaskID)
	}
	b.WriteString(fmt.Sprintf("%s (Pressione Enter para avançar, Esc para cancelar)\n\n", formTitle))

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if m.focusIndex == i {
			b.WriteString(" <")
		}
		b.WriteString("\n")
	}
	b.WriteString("\nPressione Enter no último campo para Salvar.")
	if m.err != nil {
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

	title := task.Title
	desc := task.Description
	dueDateStr := "N/A"
	if task.DueDate != nil {
		dueDateStr = task.DueDate.Format("02/01/2006")
	}
	classIDStr := "N/A (Tarefa Geral)"
	if task.ClassID != nil && *task.ClassID != 0 {
		classIDStr = fmt.Sprintf("%d", *task.ClassID)
	}
	statusStr := "Pendente"
	if task.IsCompleted {
		statusStr = "Concluída"
		title = strikethroughStyle.Render(title)
	}

	b.WriteString(fmt.Sprintf("Detalhes da Tarefa (ID: %d)\n", task.ID))
	b.WriteString(strings.Repeat("-", 30) + "\n")

	b.WriteString(fmt.Sprintf("Título: %s\n", title))
	b.WriteString(fmt.Sprintf("Descrição: %s\n", desc))

	b.WriteString(fmt.Sprintf("Prazo: %s\n", dueDateStr))
	b.WriteString(fmt.Sprintf("ID Turma: %s\n", classIDStr))
	b.WriteString(fmt.Sprintf("Status: %s\n", statusStr))

	b.WriteString(strings.Repeat("-", 30) + "\n")
	b.WriteString("\nPressione Esc para voltar à lista.")

	detailStyle := lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("63"))
	contentWidth := m.width - detailStyle.GetHorizontalFrameSize() - 4

	descLines := strings.Split(task.Description, "\n")
	wrappedDesc := ""
	for _, line := range descLines {
		if contentWidth > 0 && len(line) > contentWidth {
			for len(line) > contentWidth {
				wrappedDesc += line[:contentWidth] + "\n"
				line = line[contentWidth:]
			}
		}
		wrappedDesc += line + "\n"
	}

	finalView := fmt.Sprintf("Detalhes da Tarefa (ID: %d)\n%s\n", task.ID, strings.Repeat("-", 30))
	finalView += fmt.Sprintf("Título: %s\n", title)
	finalView += fmt.Sprintf("Descrição:\n%s\n", strings.TrimSpace(wrappedDesc))
	finalView += fmt.Sprintf("Prazo: %s\n", dueDateStr)
	finalView += fmt.Sprintf("ID Turma: %s\n", classIDStr)
	finalView += fmt.Sprintf("Status: %s\n", statusStr)
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
		return
	}
	m.inputs[0].Reset() // Title
	m.inputs[1].Reset() // Description
	m.inputs[2].Reset() // DueDate
	m.inputs[3].Reset() // ClassID
	m.focusIndex = 0
	if len(m.inputs) > 0 {
		m.inputs[0].Focus()
	}
	m.err = nil
}

// SetSize allows the main app model to adjust the size of this component.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	availableHeight := height - baseStyle.GetVerticalFrameSize() - 5
	if availableHeight < 2 { availableHeight = 2} // Corrected variable name typo

	pendingTableHeight := availableHeight / 2
	if len(m.completedTasksTable.Rows()) == 0 && availableHeight > 1 { // Changed TotalRows to len(Rows)
	    pendingTableHeight = availableHeight -1
		if pendingTableHeight < 1 { pendingTableHeight = 1 }
	}
	completedTableHeight := availableHeight - pendingTableHeight

	if pendingTableHeight < 1 { pendingTableHeight = 1 }
	if completedTableHeight < 1 { completedTableHeight = 1 }

	tableWidth := width - baseStyle.GetHorizontalFrameSize()
	m.pendingTasksTable.SetWidth(tableWidth)
	m.pendingTasksTable.SetHeight(pendingTableHeight)
	m.completedTasksTable.SetWidth(tableWidth)
	m.completedTasksTable.SetHeight(completedTableHeight)

	inputWidth := width - baseStyle.GetHorizontalFrameSize() - 10
	if inputWidth < 20 { inputWidth = 20 }
	for i := range m.inputs {
		m.inputs[i].Width = inputWidth
	}
}

// CanGoBack returns true if the model is in a state where 'esc' should return to the main menu.
func (m *Model) CanGoBack() bool {
	return m.currentView == TableView
}

// IsLoading returns true if the model is currently loading data.
func (m *Model) IsLoading() bool {
    return m.isLoading
}
