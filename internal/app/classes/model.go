package classes

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"vigenda/internal/models"
	"vigenda/internal/service"
	// "vigenda/internal/tui" // For potential GetInput like utility
	"context"
	"strconv"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// ViewState defines the current state of the classes view
type ViewState int

const (
	ListView ViewState = iota
	CreateClassView
	ImportStudentsView
	UpdateStudentStatusView
	// Potentially add DetailView for a specific class
)

// Model represents the classes management model.
type Model struct {
	classService service.ClassService
	state        ViewState
	list         list.Model // For listing classes or actions
	table        table.Model // For displaying classes or students

	// Form fields for creating/editing
	textInputs []textinput.Model
	focusIndex int

	isLoading bool
	err       error
	message   string // For success/info messages

	// Data
	classes  []models.Class
	students []models.Student // For when viewing students of a class or importing

	width  int
	height int
}

// Messages for async operations or view switching
type classesLoadedMsg struct {
	classes []models.Class
	err     error
}
type studentsLoadedMsg struct {
	students []models.Student
	err      error
}
type classCreatedMsg struct {
	class models.Class
	err   error
}
type studentsImportedMsg struct {
	count int
	err   error
}
type studentStatusUpdatedMsg struct {
	err error
}

func (m *Model) loadClassesCmd() tea.Msg {
	classes, err := m.classService.ListAllClasses(context.Background()) // Assuming ListAllClasses exists
	return classesLoadedMsg{classes: classes, err: err}
}

// New creates a new classes management model.
func New(classService service.ClassService) Model {
	// For ListView - main actions
	actionItems := []list.Item{
		actionItem{title: "Listar Todas as Turmas", description: "Visualizar todas as turmas cadastradas."},
		actionItem{title: "Criar Nova Turma", description: "Adicionar uma nova turma ao sistema."},
		actionItem{title: "Importar Alunos para Turma", description: "Adicionar alunos a uma turma a partir de um arquivo CSV."},
		actionItem{title: "Atualizar Status de Aluno", description: "Modificar o status de um aluno (ativo, inativo, etc.)."},
	}
	l := list.New(actionItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Gerenciar Turmas e Alunos"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "selecionar")),
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "voltar ao menu")),
		}
	}

	// For table view (listing classes)
	cols := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Nome da Turma", Width: 30},
		{Title: "ID Disciplina", Width: 15},
		{Title: "Nome Disciplina", Width: 25}, // Assuming we can get this
	}
	tbl := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	tbl.SetStyles(s)


	// Text inputs for forms
	inputs := make([]textinput.Model, 3) // Max 3 inputs for now (e.g. Create Class: Name, SubjectID; Import: ClassID, Path)
	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		t.CharLimit = 64 // General purpose limit
		inputs[i] = t
	}


	return Model{
		classService: classService,
		state:        ListView, // Start with the list of actions
		list:         l,
		table:        tbl,
		textInputs:   inputs,
		isLoading:    false, // Not loading initially, actions will trigger loading
	}
}

func (m Model) Init() tea.Cmd {
	m.state = ListView // Reset to main action list view
	m.err = nil
	m.message = ""
	// No initial data loading, user chooses an action first.
	// If we wanted to show a list of classes by default, we'd call loadClassesCmd here.
	return nil
}

// actionItem implements list.Item for the main action list.
type actionItem struct {
	title, description string
}

func (i actionItem) Title() string       { return i.title }
func (i actionItem) Description() string { return i.description }
func (i actionItem) FilterValue() string { return i.title }


func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global keys for this model, like 'esc' to go back a state or to main menu
		if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) {
			if m.state == ListView { // If in main action list, esc does nothing here (handled by parent)
				return m, nil // Or signal parent to go to main menu
			}
			// If in a sub-view (form, table), go back to ListView
			m.state = ListView
			m.err = nil
			m.message = ""
			m.resetInputs()
			// Reset list selection to avoid re-triggering same action
			m.list.Select(-1) // Deselect
			return m, nil
		}

		// State-specific key handling
		switch m.state {
		case ListView:
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				selected, ok := m.list.SelectedItem().(actionItem)
				if ok {
					m.err = nil
					m.message = ""
					switch selected.title {
					case "Listar Todas as Turmas":
						m.isLoading = true
						m.state = ListView // Stay in ListView but show table
						m.table.SetColumns([]table.Column{ // Columns for classes
							{Title: "ID", Width: 5},
							{Title: "Nome da Turma", Width: 30},
							{Title: "ID Disciplina", Width: 15},
							// {Title: "Nome Disciplina", Width: 25}, // TODO: Fetch subject name
						})
						cmds = append(cmds, m.loadClassesCmd)
					case "Criar Nova Turma":
						m.state = CreateClassView
						m.setupCreateClassForm()
					case "Importar Alunos para Turma":
						m.state = ImportStudentsView
						m.setupImportStudentsForm()
					case "Atualizar Status de Aluno":
						m.state = UpdateStudentStatusView
						m.setupUpdateStudentStatusForm()
					}
				}
			}
		case CreateClassView, ImportStudentsView, UpdateStudentStatusView:
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				if m.focusIndex == len(m.textInputs) { // "Submit"
					m.isLoading = true
					m.err = nil
					cmds = append(cmds, m.submitFormCmd())
				} else { // Focus next input
					cmds = append(cmds, m.updateFocus())
				}
			} else if key.Matches(msg, key.NewBinding(key.WithKeys("up"))) {
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = len(m.textInputs)
				}
				m.updateInputFocusStyle()
			} else if key.Matches(msg, key.NewBinding(key.WithKeys("down"))) {
				m.focusIndex++
				if m.focusIndex > len(m.textInputs) {
					m.focusIndex = 0
				}
				m.updateInputFocusStyle()
			} else { // Pass to focused text input
				if m.focusIndex < len(m.textInputs) {
					m.textInputs[m.focusIndex], cmd = m.textInputs[m.focusIndex].Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		}

	// Handle results of async operations
	case classesLoadedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.classes = msg.classes
		rows := make([]table.Row, len(m.classes))
		for i, cls := range m.classes {
			rows[i] = table.Row{
				fmt.Sprintf("%d", cls.ID),
				cls.Name,
				fmt.Sprintf("%d", cls.SubjectID),
				// cls.SubjectName, // TODO
			}
		}
		m.table.SetRows(rows)
		// We are already in ListView state, but now showing the table
		// No need to change m.state here unless we had a dedicated TableView state

	case classCreatedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.message = fmt.Sprintf("Turma '%s' (ID: %d) criada com sucesso!", msg.class.Name, msg.class.ID)
			m.state = ListView // Go back to action list
			m.resetInputs()
		}

	case studentsImportedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.message = fmt.Sprintf("%d alunos importados com sucesso!", msg.count)
			m.state = ListView
			m.resetInputs()
		}

	case studentStatusUpdatedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.message = "Status do aluno atualizado com sucesso!"
			m.state = ListView
			m.resetInputs()
		}

	case error: // Generic error message
		m.err = msg
		m.isLoading = false


	// Propagate WindowSizeMsg to components
	case tea.WindowSizeMsg:
		m.width = msg.Width - baseStyle.GetHorizontalFrameSize() // Available width for content
		m.height = msg.Height - baseStyle.GetVerticalFrameSize()   // Available height for content

		// Adjust list size
		listHeight := m.height - lipgloss.Height(m.list.Title) - 2 // Some padding/margin
		m.list.SetSize(m.width, listHeight)

		// Adjust table size
		// Table height: available height minus title (if any), error/message line, help line.
		// This needs careful calculation based on what's rendered in View() for this state.
		m.table.SetWidth(m.width)
		m.table.SetHeight(m.height - 5) // Rough estimate, adjust as needed

		// Adjust text input widths if they are visible
		inputWidth := m.width - 4 // Some padding for inputs
		if inputWidth < 20 { inputWidth = 20 } // Minimum width
		for i := range m.textInputs {
			m.textInputs[i].Width = inputWidth
		}
	}

	// Update table if it's the component in focus (e.g. for scrolling)
	// This is a bit tricky if the table isn't always "active"
	// For now, we assume if classes are loaded, table might be active.
	if m.state == ListView && len(m.classes) > 0 {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}


	return m, tea.Batch(cmds...)
}

func (m *Model) setupCreateClassForm() {
	m.focusIndex = 0
	m.textInputs = make([]textinput.Model, 2) // Name, SubjectID

	m.textInputs[0] = textinput.New()
	m.textInputs[0].Placeholder = "Nome da Turma"
	m.textInputs[0].Focus()
	m.textInputs[0].CharLimit = 50
	m.textInputs[0].Width = m.width / 2 // Example width

	m.textInputs[1] = textinput.New()
	m.textInputs[1].Placeholder = "ID da Disciplina (numérico)"
	m.textInputs[1].CharLimit = 10
	m.textInputs[1].Width = m.width / 2
	m.textInputs[1].Validate = isNumber // Simple numeric validation

	m.updateInputFocusStyle()
}

func (m *Model) setupImportStudentsForm() {
	m.focusIndex = 0
	m.textInputs = make([]textinput.Model, 2) // ClassID, CSV File Path

	m.textInputs[0] = textinput.New()
	m.textInputs[0].Placeholder = "ID da Turma (numérico)"
	m.textInputs[0].Focus()
	m.textInputs[0].CharLimit = 10
	m.textInputs[0].Width = m.width / 2
	m.textInputs[0].Validate = isNumber

	m.textInputs[1] = textinput.New()
	m.textInputs[1].Placeholder = "Caminho para o arquivo CSV"
	m.textInputs[1].CharLimit = 255
	m.textInputs[1].Width = m.width / 2

	m.updateInputFocusStyle()
}

func (m *Model) setupUpdateStudentStatusForm() {
	m.focusIndex = 0
	m.textInputs = make([]textinput.Model, 2) // StudentID, NewStatus

	m.textInputs[0] = textinput.New()
	m.textInputs[0].Placeholder = "ID do Aluno (numérico)"
	m.textInputs[0].Focus()
	m.textInputs[0].CharLimit = 10
	m.textInputs[0].Width = m.width / 2
	m.textInputs[0].Validate = isNumber

	m.textInputs[1] = textinput.New()
	m.textInputs[1].Placeholder = "Novo Status (ativo, inativo, transferido)"
	m.textInputs[1].CharLimit = 20
	m.textInputs[1].Width = m.width / 2
	// TODO: Add validation for status values

	m.updateInputFocusStyle()
}

func (m *Model) resetInputs() {
	for i := range m.textInputs {
		m.textInputs[i].Reset()
		m.textInputs[i].Blur()
	}
	m.focusIndex = 0
}


func (m *Model) updateFocus() tea.Cmd {
	m.focusIndex = (m.focusIndex + 1) % (len(m.textInputs) + 1) // +1 for submit "button"
	return m.updateInputFocusStyle()
}

func (m *Model) updateInputFocusStyle() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textInputs))
	for i := 0; i < len(m.textInputs); i++ {
		if i == m.focusIndex {
			cmds[i] = m.textInputs[i].Focus()
			m.textInputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		} else {
			m.textInputs[i].Blur()
			m.textInputs[i].PromptStyle = lipgloss.NewStyle()
		}
	}
	return tea.Batch(cmds...)
}


func isNumber(s string) error {
	if s == "" { return nil } // Allow empty for now, validate on submit
	if _, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("deve ser um número")
	}
	return nil
}

func (m *Model) submitFormCmd() tea.Cmd {
	switch m.state {
	case CreateClassView:
		className := m.textInputs[0].Value()
		subjectIDStr := m.textInputs[1].Value()
		if className == "" || subjectIDStr == "" {
			m.err = fmt.Errorf("Nome da turma e ID da disciplina são obrigatórios.")
			m.isLoading = false
			return nil
		}
		subjectID, err := strconv.ParseInt(subjectIDStr, 10, 64)
		if err != nil {
			m.err = fmt.Errorf("ID da disciplina inválido: %w", err)
			m.isLoading = false
			return nil
		}
		return func() tea.Msg {
			cls, err := m.classService.CreateClass(context.Background(), className, subjectID)
			return classCreatedMsg{class: cls, err: err}
		}
	case ImportStudentsView:
		classIDStr := m.textInputs[0].Value()
		csvPath := m.textInputs[1].Value()
		if classIDStr == "" || csvPath == "" {
			m.err = fmt.Errorf("ID da Turma e Caminho do CSV são obrigatórios.")
			m.isLoading = false
			return nil
		}
		// classID, err := strconv.ParseInt(classIDStr, 10, 64) // Comentado para evitar erro de não utilizado
		_, err := strconv.ParseInt(classIDStr, 10, 64)
		if err != nil {
			m.err = fmt.Errorf("ID da Turma inválido: %w", err)
			m.isLoading = false
			return nil
		}
		// Reading file and calling service. This is complex for a simple TUI.
		// Ideally, the service would take the path and handle file reading.
		// For now, we'll assume the service needs byte data. This is a simplification.
		// In a real app, you'd use a file picker or more robust path handling.
		// For this example, we'll skip the actual file reading in the TUI model itself.
		// The CLI version does os.ReadFile. The TUI might need a different approach or
		// this functionality might be simplified/deferred.
		// Let's assume the service can take the path for now for the sake of the TUI structure.
		// This requires service.ImportStudentsFromCSVByPath(ctx, classID, path)
		// If not, this part needs to be more complex or rethought.
		// For now, let's simulate the call if such a method existed:
		// return func() tea.Msg {
		//  count, err := m.classService.ImportStudentsFromCSVByPath(context.Background(), classID, csvPath)
		// 	return studentsImportedMsg{count: count, err: err}
		// }
		// Since it likely doesn't, we show an error.
		m.err = fmt.Errorf("Importação de CSV via TUI não implementada completamente (leitura de arquivo local).")
		m.isLoading = false
		return nil

	case UpdateStudentStatusView:
		studentIDStr := m.textInputs[0].Value()
		newStatus := m.textInputs[1].Value()
		if studentIDStr == "" || newStatus == "" {
			m.err = fmt.Errorf("ID do Aluno e Novo Status são obrigatórios.")
			m.isLoading = false
			return nil
		}
		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			m.err = fmt.Errorf("ID do Aluno inválido: %w", err)
			m.isLoading = false
			return nil
		}
		// Validate status (basic check)
		validStatuses := []string{"ativo", "inativo", "transferido"}
		isValidStatus := false
		for _, vs := range validStatuses {
			if newStatus == vs {
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			m.err = fmt.Errorf("Status inválido. Use: ativo, inativo, transferido.")
			m.isLoading = false
			return nil
		}

		return func() tea.Msg {
			err := m.classService.UpdateStudentStatus(context.Background(), studentID, newStatus)
			return studentStatusUpdatedMsg{err: err}
		}
	}
	return nil
}


func (m Model) View() string {
	var s strings.Builder

	if m.isLoading {
		s.WriteString("Carregando...")
		return baseStyle.Render(s.String())
	}
	if m.err != nil {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Erro: %v\n\n", m.err)))
	}
	if m.message != "" {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(fmt.Sprintf("%s\n\n", m.message)))
	}

	switch m.state {
	case ListView:
		// If classes are loaded, show table, otherwise show action list
		if len(m.classes) > 0 && m.err == nil { // Check err == nil to not show table on load error
			s.WriteString("Turmas Cadastradas:\n")
			s.WriteString(m.table.View())
			s.WriteString("\n(Navegue com ↑/↓, 'esc' para voltar às ações)")
		} else if m.err != nil && len(m.classes) == 0 { // Error occurred and no classes loaded
			s.WriteString("(Pressione 'esc' para voltar às ações)")
		} else { // Show action list
			s.WriteString(m.list.View())
			// Help is part of list.AdditionalShortHelpKeys
		}

	case CreateClassView:
		s.WriteString("Criar Nova Turma\n\n")
		for i := range m.textInputs {
			s.WriteString(m.textInputs[i].View() + "\n")
		}
		submitButton := "[ Criar Turma ]"
		if m.focusIndex == len(m.textInputs) {
			submitButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(submitButton)
		}
		s.WriteString("\n" + submitButton + "\n\n")
		s.WriteString("(Use Tab/Shift+Tab ou ↑/↓ para navegar, Enter para submeter, Esc para cancelar)")

	case ImportStudentsView:
		s.WriteString("Importar Alunos para Turma\n\n")
		for i := range m.textInputs {
			s.WriteString(m.textInputs[i].View() + "\n")
		}
		submitButton := "[ Importar Alunos ]"
		if m.focusIndex == len(m.textInputs) {
			submitButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(submitButton)
		}
		s.WriteString("\n" + submitButton + "\n\n")
		s.WriteString("(Use Tab/Shift+Tab ou ↑/↓ para navegar, Enter para submeter, Esc para cancelar)")
		s.WriteString("\nNota: A leitura direta de arquivos do sistema pela TUI é complexa.\nEsta funcionalidade pode ser limitada ou exigir que o arquivo esteja acessível de uma forma específica.")


	case UpdateStudentStatusView:
		s.WriteString("Atualizar Status de Aluno\n\n")
		for i := range m.textInputs {
			s.WriteString(m.textInputs[i].View() + "\n")
		}
		submitButton := "[ Atualizar Status ]"
		if m.focusIndex == len(m.textInputs) {
			submitButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(submitButton)
		}
		s.WriteString("\n" + submitButton + "\n\n")
		s.WriteString("(Use Tab/Shift+Tab ou ↑/↓ para navegar, Enter para submeter, Esc para cancelar)")

	default:
		s.WriteString("Visualização de Turmas Desconhecida")
	}

	return baseStyle.Render(s.String())
}

// SetSize allows the main app model to adjust the size of this component.
func (m *Model) SetSize(width, height int) {
	m.width = width - baseStyle.GetHorizontalFrameSize()
	m.height = height - baseStyle.GetVerticalFrameSize() -1 // -1 for the final newline in main app view

	// Update components that depend on size
	listHeight := m.height - lipgloss.Height(m.list.Title) - 2
	m.list.SetSize(m.width, listHeight)

	m.table.SetWidth(m.width)
	// Dynamic height for table, considering title, error/message, help
	// For simplicity, let's make it a bit smaller than full height.
	tableHeight := m.height - 6
	if tableHeight < 5 { tableHeight = 5 }
	m.table.SetHeight(tableHeight)

	inputWidth := m.width - 4
	if inputWidth < 20 { inputWidth = 20 }
	for i := range m.textInputs {
		if i < len(m.textInputs) { // Check slice bounds
			m.textInputs[i].Width = inputWidth
		}
	}
}

// IsFocused can be used by the parent model to determine if this sub-model
// should capture global keys like 'esc' or 'q'.
// For now, returning false means the parent model will handle 'esc' to navigate away.
// If this model had its own popups or modes that 'esc' should close first,
// this would return true in those cases.
func (m Model) IsFocused() bool {
	// If a form is active (i.e., textInputs are being used), consider it focused.
	return m.state == CreateClassView || m.state == ImportStudentsView || m.state == UpdateStudentStatusView
}
