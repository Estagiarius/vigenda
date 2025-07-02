package classes

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"vigenda/internal/models"
	"vigenda/internal/service"
)

type ViewState int

const (
	ListView ViewState = iota
	CreatingView
	EditingClassView
	DeletingClassConfirmView
	DetailsView
	AddingStudentView
	EditingStudentView
	DeletingStudentConfirmView
)

type FocusTarget int

const (
	FocusTargetNone FocusTarget = iota
	FocusTargetStudentsTable
	FocusTargetClassForm
	FocusTargetStudentForm
)

var (
	columnTitleID        = "ID"
	columnTitleName      = "Nome da Turma"
	columnTitleSubjectID = "ID Disciplina"
	columnTitleCreatedAt = "Criada em"
	columnTitleUpdatedAt = "Atualizada em"

	dbOperationTimeout = 5 * time.Second

	studentColumnTitleID         = "ID Aluno"
	studentColumnTitleEnrollment = "Nº Chamada"
	studentColumnTitleFullName   = "Nome Completo"
	studentColumnTitleStatus     = "Status"
	studentColumnTitleCreatedAt  = "Criado em"
	studentColumnTitleUpdatedAt  = "Atualizado em"
)

type Model struct {
	classService service.ClassService
	state        ViewState
	table        table.Model
	formInputs   struct {
		inputs     []textinput.Model
		focusIndex int
	}
	allClasses           []models.Class
	selectedClass        *models.Class
	selectedStudent      *models.Student
	classStudents        []models.Student
	studentsTable        table.Model
	detailsViewFocusTarget FocusTarget
	isLoading            bool
	width                int
	height               int
	err                  error
}

func New(cs service.ClassService) *Model {
	log.Println("ClassesModel: New")
	classTable := table.New(
		table.WithColumns([]table.Column{
			{Title: columnTitleID, Width: 5},
			{Title: columnTitleName, Width: 25},
			{Title: columnTitleSubjectID, Width: 15},
			{Title: columnTitleCreatedAt, Width: 18},
			{Title: columnTitleUpdatedAt, Width: 18},
		}),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	classTable.SetStyles(s)

	studentsTable := table.New(
		table.WithColumns([]table.Column{
			{Title: studentColumnTitleID, Width: 8},
			{Title: studentColumnTitleEnrollment, Width: 10},
			{Title: studentColumnTitleFullName, Width: 25},
			{Title: studentColumnTitleStatus, Width: 10},
			{Title: studentColumnTitleCreatedAt, Width: 18},
			{Title: studentColumnTitleUpdatedAt, Width: 18},
		}),
		table.WithRows([]table.Row{}),
		table.WithFocused(false),
		table.WithHeight(10),
	)
	studentsTable.SetStyles(s)

	return &Model{
		classService:  cs,
		state:         ListView,
		table:         classTable,
		studentsTable: studentsTable,
		formInputs:    struct{ inputs []textinput.Model; focusIndex int }{inputs: []textinput.Model{}, focusIndex: 0},
		isLoading:     true,
	}
}

func (m *Model) Init() tea.Cmd {
	m.isLoading = true
	return m.fetchClassesCmd
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		currentCmd := m.handleKeyPress(msg)
		cmds = append(cmds, currentCmd)
	case fetchedClassesMsg:
		cmd = m.handleFetchedClasses(msg)
	case classCreatedMsg:
		cmd = m.handleClassCreated(msg)
	case classUpdatedMsg:
		cmd = m.handleClassUpdated(msg)
	case classDeletedMsg:
		cmd = m.handleClassDeleted(msg)
	case fetchedClassStudentsMsg:
		cmd = m.handleFetchedClassStudents(msg)
	case studentAddedMsg:
		cmd = m.handleStudentAdded(msg)
	case studentUpdatedMsg:
		cmd = m.handleStudentUpdated(msg)
	case studentDeletedMsg:
		cmd = m.handleStudentDeleted(msg)
	case errMsg:
		m.err = msg.err
		m.isLoading = false
	}
	cmds = append(cmds, cmd)

	// Update focused form input if any form is active
	if (m.state == CreatingView || m.state == EditingClassView || m.state == AddingStudentView || m.state == EditingStudentView) &&
		len(m.formInputs.inputs) > 0 && m.formInputs.focusIndex >= 0 && m.formInputs.focusIndex < len(m.formInputs.inputs) {
		// Only update if the message is a KeyMsg and it wasn't an action key already handled by form logic
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			isActionKey := key.Matches(keyMsg, key.NewBinding(key.WithKeys("enter", "esc", "tab", "shift+tab")))
			if !isActionKey {
				var inputCmd tea.Cmd
				m.formInputs.inputs[m.formInputs.focusIndex], inputCmd = m.formInputs.inputs[m.formInputs.focusIndex].Update(msg)
				cmds = append(cmds, inputCmd)
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	switch m.state {
	case ListView:
		return m.handleListViewKeys(msg)
	case CreatingView, EditingClassView:
		return m.handleClassFormKeys(msg)
	case DeletingClassConfirmView:
		return m.handleDeleteClassConfirmKeys(msg)
	case DetailsView:
		return m.handleDetailsViewKeys(msg)
	case AddingStudentView, EditingStudentView:
		return m.handleStudentFormKeys(msg)
	case DeletingStudentConfirmView:
		return m.handleDeleteStudentConfirmKeys(msg)
	}
	return nil
}

func (m *Model) View() string {
	var b strings.Builder
	if m.isLoading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "Carregando...")
	}
	if m.err != nil {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).PaddingBottom(1).Render(fmt.Sprintf("Erro: %v", m.err)))
	}

	switch m.state {
	case ListView:
		b.WriteString(lipgloss.NewStyle().Bold(true).MarginBottom(1).Render("Lista de Turmas"))
		b.WriteString(m.table.View())
		help := "↑/↓: Navegar | Enter: Detalhes | n: Nova | e: Editar | d: Deletar | q/Esc: Voltar"
		b.WriteString(lipgloss.NewStyle().Faint(true).MarginTop(1).Render(help))
	case CreatingView, EditingClassView:
		title := "Nova Turma"
		if m.state == EditingClassView && m.selectedClass != nil {
			title = fmt.Sprintf("Editando Turma: %s (ID: %d)", m.selectedClass.Name, m.selectedClass.ID)
		}
		b.WriteString(lipgloss.NewStyle().Bold(true).MarginBottom(1).Render(title))
		for _, input := range m.formInputs.inputs {
			b.WriteString(input.View() + "\n")
		}
		b.WriteString("\n" + lipgloss.NewStyle().Faint(true).Render("Tab: Próximo | Shift+Tab: Anterior | Enter: Salvar | Esc: Cancelar"))
	case DeletingClassConfirmView:
		if m.selectedClass != nil {
			b.WriteString(lipgloss.NewStyle().Bold(true).MarginBottom(1).Render(fmt.Sprintf("Confirmar Exclusão da Turma: %s?", m.selectedClass.Name)))
			b.WriteString("Esta ação não pode ser desfeita.\n\n")
			b.WriteString(lipgloss.NewStyle().Faint(true).Render("Pressione 's' para confirmar, 'n' ou 'Esc' para cancelar."))
		}
	case DetailsView:
		if m.selectedClass == nil {
			b.WriteString("Erro: Nenhuma turma selecionada.")
		} else {
			headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
			b.WriteString(lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Detalhes da Turma: %s\n", m.selectedClass.Name)))
			b.WriteString(fmt.Sprintf("%s %d\n", headerStyle.Render("ID:"), m.selectedClass.ID))
			b.WriteString(fmt.Sprintf("%s %d\n", headerStyle.Render("ID Disciplina:"), m.selectedClass.SubjectID))
			b.WriteString(fmt.Sprintf("%s %s\n", headerStyle.Render("Criada em:"), m.selectedClass.CreatedAt.Format("02/01/2006 15:04")))
			b.WriteString(fmt.Sprintf("%s %s\n\n", headerStyle.Render("Atualizada em:"), m.selectedClass.UpdatedAt.Format("02/01/2006 15:04")))
			b.WriteString(lipgloss.NewStyle().Bold(true).Render("Alunos:\n"))
			if len(m.classStudents) == 0 && m.err == nil { // Check m.err too
				b.WriteString("Nenhum aluno encontrado.")
			} else {
				tableRender := m.studentsTable.View()
				if m.detailsViewFocusTarget == FocusTargetStudentsTable {
					tableRender = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Render(tableRender)
				}
				b.WriteString(tableRender)
			}
			help := "Esc: Voltar | a: Add Aluno"
			if len(m.classStudents) > 0 { // Only show table focus keys if there are students
				if m.detailsViewFocusTarget == FocusTargetNone {
					help += " | Tab: Focar Tabela Alunos"
				} else {
					help += " | ↑/↓: Nav Alunos | e: Editar Aluno | d: Deletar Aluno | Shift+Tab/Esc: Sair Foco"
				}
			}
			b.WriteString("\n\n" + lipgloss.NewStyle().Faint(true).Render(help))
		}
	case AddingStudentView, EditingStudentView:
		title := "Adicionar Novo Aluno"
		if m.state == EditingStudentView && m.selectedStudent != nil {
			title = fmt.Sprintf("Editando Aluno: %s (ID: %d)", m.selectedStudent.FullName, m.selectedStudent.ID)
		}
		b.WriteString(lipgloss.NewStyle().Bold(true).MarginBottom(1).Render(title))
		for _, input := range m.formInputs.inputs {
			b.WriteString(input.View() + "\n")
		}
		b.WriteString("\n" + lipgloss.NewStyle().Faint(true).Render("Tab: Próximo | Shift+Tab: Anterior | Enter: Salvar | Esc: Cancelar"))
	case DeletingStudentConfirmView:
		if m.selectedStudent != nil {
			b.WriteString(lipgloss.NewStyle().Bold(true).MarginBottom(1).Render(fmt.Sprintf("Confirmar Exclusão do Aluno: %s?", m.selectedStudent.FullName)))
			b.WriteString("Esta ação não pode ser desfeita.\n\n")
			b.WriteString(lipgloss.NewStyle().Faint(true).Render("Pressione 's' para confirmar, 'n' ou 'Esc' para cancelar."))
		}
	}
	return b.String()
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	contentHeight := height - 5 // Approx for title, error, help
	if contentHeight < 1 { contentHeight = 1 }
	tableWidth := width - 4
	if tableWidth < 20 { tableWidth = 20 }

	m.table.SetHeight(contentHeight)
	m.table.SetWidth(tableWidth)
	m.studentsTable.SetHeight(contentHeight - 7) // Adjusted for class details
	m.studentsTable.SetWidth(tableWidth)

	inputWidth := width - 8
	if inputWidth < 20 { inputWidth = 20 }
	for i := range m.formInputs.inputs {
		m.formInputs.inputs[i].Width = inputWidth
	}
}

func (m *Model) IsFocused() bool {
	return m.state == CreatingView || m.state == EditingClassView || m.state == AddingStudentView || m.state == EditingStudentView
}

// Key Handlers
func (m *Model) handleListViewKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
		m.prepareClassForm(nil)
		return textinput.Blink
	case key.Matches(msg, key.NewBinding(key.WithKeys("e"))):
		if idx := m.table.Cursor(); idx < len(m.allClasses) {
			m.selectedClass = &m.allClasses[idx]
			m.prepareClassForm(m.selectedClass)
			return textinput.Blink
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("d"))):
		if idx := m.table.Cursor(); idx < len(m.allClasses) {
			m.selectedClass = &m.allClasses[idx]
			m.state = DeletingClassConfirmView
			m.err = nil
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
		if idx := m.table.Cursor(); idx < len(m.allClasses) {
			m.selectedClass = &m.allClasses[idx]
			m.state = DetailsView
			m.isLoading = true
			m.err = nil
			m.classStudents = nil
			m.studentsTable.SetRows([]table.Row{})
			cmds = append(cmds, m.fetchClassStudentsCmd(m.selectedClass.ID))
		}
	default:
		var tableCmd tea.Cmd
		m.table, tableCmd = m.table.Update(msg)
		cmds = append(cmds, tableCmd)
	}
	return tea.Batch(cmds...)
}

func (m *Model) handleClassFormKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
		m.state = ListView
		m.err = nil
		m.resetFormInputs()
		m.selectedClass = nil
		m.table.Focus()
	case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
		if m.formInputs.focusIndex == len(m.formInputs.inputs)-1 {
			name := strings.TrimSpace(m.formInputs.inputs[0].Value())
			subjectIDStr := strings.TrimSpace(m.formInputs.inputs[1].Value())
			if name == "" || subjectIDStr == "" {
				m.err = fmt.Errorf("nome e ID da disciplina obrigatórios")
				return nil
			}
			m.isLoading = true
			if m.state == CreatingView {
				return m.createClassCmd(name, subjectIDStr)
			} else if m.state == EditingClassView && m.selectedClass != nil {
				return m.updateClassCmd(m.selectedClass.ID, name, subjectIDStr)
			}
		} else {
			m.nextFormInput()
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
		m.nextFormInput()
	case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
		m.prevFormInput()
	}
	return nil // Input updates are handled in main Update loop
}

func (m *Model) handleDeleteClassConfirmKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "s", "S":
		if m.selectedClass != nil {
			m.isLoading = true
			return m.deleteClassCmd(m.selectedClass.ID)
		}
	case "n", "N", "esc":
		m.state = ListView
		m.selectedClass = nil
		m.err = nil
		m.table.Focus()
	}
	return nil
}

func (m *Model) handleDetailsViewKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
		if m.detailsViewFocusTarget == FocusTargetStudentsTable {
			m.detailsViewFocusTarget = FocusTargetNone
			m.studentsTable.Blur()
		} else {
			m.state = ListView
			m.selectedClass = nil
			m.err = nil
			m.table.Focus()
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("a"))):
		if m.selectedClass != nil {
			m.prepareStudentForm(nil)
			return textinput.Blink
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
		if m.detailsViewFocusTarget == FocusTargetNone && len(m.classStudents) > 0 {
			m.detailsViewFocusTarget = FocusTargetStudentsTable
			m.studentsTable.Focus()
		} else { // Blur or cycle
			m.detailsViewFocusTarget = FocusTargetNone
			m.studentsTable.Blur()
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
		if m.detailsViewFocusTarget == FocusTargetStudentsTable {
			m.detailsViewFocusTarget = FocusTargetNone
			m.studentsTable.Blur()
		}
	default:
		if m.detailsViewFocusTarget == FocusTargetStudentsTable {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("e"))):
				if idx := m.studentsTable.Cursor(); idx < len(m.classStudents) {
					m.selectedStudent = &m.classStudents[idx]
					m.prepareStudentForm(m.selectedStudent)
					return textinput.Blink
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("d"))):
				if idx := m.studentsTable.Cursor(); idx < len(m.classStudents) {
					m.selectedStudent = &m.classStudents[idx]
					m.state = DeletingStudentConfirmView
					m.err = nil
				}
			default:
				var tableCmd tea.Cmd
				m.studentsTable, tableCmd = m.studentsTable.Update(msg)
				cmds = append(cmds, tableCmd)
			}
		}
	}
	return tea.Batch(cmds...)
}

func (m *Model) handleStudentFormKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
		m.state = DetailsView
		m.err = nil
		m.resetFormInputs()
		m.selectedStudent = nil
		// Re-fetch students for details view to be up-to-date
		if m.selectedClass != nil {
			m.isLoading = true
			return m.fetchClassStudentsCmd(m.selectedClass.ID)
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
		if m.formInputs.focusIndex == len(m.formInputs.inputs)-1 { // Last field
			fullName := strings.TrimSpace(m.formInputs.inputs[0].Value())
			enrollmentID := strings.TrimSpace(m.formInputs.inputs[1].Value())
			status := strings.ToLower(strings.TrimSpace(m.formInputs.inputs[2].Value()))
			if fullName == "" || status == "" {
				m.err = fmt.Errorf("nome completo e status são obrigatórios")
				return nil
			}
			// Basic status validation
			if status != "ativo" && status != "inativo" && status != "transferido" {
				m.err = fmt.Errorf("status inválido: use ativo, inativo ou transferido")
				return nil
			}
			m.isLoading = true
			if m.state == AddingStudentView && m.selectedClass != nil {
				return m.addStudentCmd(m.selectedClass.ID, fullName, enrollmentID, status)
			} else if m.state == EditingStudentView && m.selectedStudent != nil {
				return m.updateStudentCmd(m.selectedStudent.ID, fullName, enrollmentID, status)
			}
		} else {
			m.nextFormInput()
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
		m.nextFormInput()
	case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
		m.prevFormInput()
	}
	return nil // Input updates are handled in main Update loop
}

func (m *Model) handleDeleteStudentConfirmKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "s", "S":
		if m.selectedStudent != nil {
			m.isLoading = true
			return m.deleteStudentCmd(m.selectedStudent.ID)
		}
	case "n", "N", "esc":
		m.state = DetailsView
		m.selectedStudent = nil
		m.err = nil
		// Re-fetch students to update view
		if m.selectedClass != nil {
			m.isLoading = true
			return m.fetchClassStudentsCmd(m.selectedClass.ID)
		}
	}
	return nil
}

// Form Management
func (m *Model) resetFormInputs() {
	m.formInputs.inputs = []textinput.Model{}
	m.formInputs.focusIndex = 0
}

func (m *Model) prepareClassForm(classToEdit *models.Class) {
	m.resetFormInputs()
	nameInput := textinput.New()
	nameInput.Placeholder = "Nome da Turma"
	nameInput.CharLimit = 100
	subjectIDInput := textinput.New()
	subjectIDInput.Placeholder = "ID da Disciplina (ex: 1)"
	subjectIDInput.CharLimit = 10

	if classToEdit != nil {
		nameInput.SetValue(classToEdit.Name)
		subjectIDInput.SetValue(strconv.FormatInt(classToEdit.SubjectID, 10))
		m.state = EditingClassView
	} else {
		m.state = CreatingView
	}
	m.formInputs.inputs = []textinput.Model{nameInput, subjectIDInput}
	m.formInputs.inputs[0].Focus()
	m.formInputs.focusIndex = 0
	m.err = nil
}

func (m *Model) prepareStudentForm(studentToEdit *models.Student) {
	m.resetFormInputs()
	fullNameInput := textinput.New()
	fullNameInput.Placeholder = "Nome Completo do Aluno"
	enrollmentIDInput := textinput.New()
	enrollmentIDInput.Placeholder = "Nº Chamada/Matrícula (opcional)"
	statusInput := textinput.New()
	statusInput.Placeholder = "Status (ativo, inativo, transferido)"

	if studentToEdit != nil {
		fullNameInput.SetValue(studentToEdit.FullName)
		enrollmentIDInput.SetValue(studentToEdit.EnrollmentID)
		statusInput.SetValue(studentToEdit.Status)
		m.state = EditingStudentView
	} else {
		statusInput.SetValue("ativo") // Default
		m.state = AddingStudentView
	}
	m.formInputs.inputs = []textinput.Model{fullNameInput, enrollmentIDInput, statusInput}
	m.formInputs.inputs[0].Focus()
	m.formInputs.focusIndex = 0
	m.err = nil
}

func (m *Model) nextFormInput() {
	if len(m.formInputs.inputs) == 0 { return }
	m.formInputs.inputs[m.formInputs.focusIndex].Blur()
	m.formInputs.focusIndex = (m.formInputs.focusIndex + 1) % len(m.formInputs.inputs)
	m.formInputs.inputs[m.formInputs.focusIndex].Focus()
}

func (m *Model) prevFormInput() {
	if len(m.formInputs.inputs) == 0 { return }
	m.formInputs.inputs[m.formInputs.focusIndex].Blur()
	m.formInputs.focusIndex = (m.formInputs.focusIndex - 1 + len(m.formInputs.inputs)) % len(m.formInputs.inputs)
	m.formInputs.inputs[m.formInputs.focusIndex].Focus()
}

// Message Handlers (CRUD results)
func (m *Model) handleFetchedClasses(msg fetchedClassesMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = msg.err
		m.allClasses = nil
	} else {
		m.err = nil
		m.allClasses = msg.classes
	}
	var rows []table.Row
	for _, cls := range m.allClasses {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", cls.ID),
			cls.Name,
			fmt.Sprintf("%d", cls.SubjectID),
			cls.CreatedAt.Format("02/01/06 15:04"),
			cls.UpdatedAt.Format("02/01/06 15:04"),
		})
	}
	m.table.SetRows(rows)
	return nil
}

func (m *Model) handleClassCreated(msg classCreatedMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = fmt.Errorf("criar turma: %w", msg.err)
		return nil
	}
	m.state = ListView
	m.err = nil
	m.resetFormInputs()
	m.table.Focus()
	return m.fetchClassesCmd
}

func (m *Model) handleClassUpdated(msg classUpdatedMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = fmt.Errorf("atualizar turma: %w", msg.err)
		return nil
	}
	m.state = ListView
	m.err = nil
	m.resetFormInputs()
	m.selectedClass = nil
	m.table.Focus()
	return m.fetchClassesCmd
}

func (m *Model) handleClassDeleted(msg classDeletedMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = fmt.Errorf("deletar turma: %w", msg.err)
		return nil
	}
	m.state = ListView
	m.err = nil
	m.selectedClass = nil
	m.table.Focus()
	return m.fetchClassesCmd
}

func (m *Model) handleFetchedClassStudents(msg fetchedClassStudentsMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = msg.err
		m.classStudents = nil
	} else {
		m.err = nil
		m.classStudents = msg.students
	}
	var rows []table.Row
	for _, std := range m.classStudents {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", std.ID),
			std.EnrollmentID,
			std.FullName,
			std.Status,
			std.CreatedAt.Format("02/01/06 15:04"),
			std.UpdatedAt.Format("02/01/06 15:04"),
		})
	}
	m.studentsTable.SetRows(rows)
	if len(rows) == 0 { // If no students, blur table
		m.detailsViewFocusTarget = FocusTargetNone
		m.studentsTable.Blur()
	} else if m.detailsViewFocusTarget == FocusTargetStudentsTable { // Re-apply focus if it was intended
		m.studentsTable.Focus()
	}
	return nil
}

func (m *Model) handleStudentAdded(msg studentAddedMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = fmt.Errorf("add aluno: %w", msg.err)
		return nil
	}
	m.state = DetailsView
	m.err = nil
	m.resetFormInputs()
	if m.selectedClass != nil {
		m.isLoading = true
		return m.fetchClassStudentsCmd(m.selectedClass.ID)
	}
	return nil
}

func (m *Model) handleStudentUpdated(msg studentUpdatedMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = fmt.Errorf("atualizar aluno: %w", msg.err)
		return nil
	}
	m.state = DetailsView
	m.err = nil
	m.resetFormInputs()
	m.selectedStudent = nil
	if m.selectedClass != nil {
		m.isLoading = true
		return m.fetchClassStudentsCmd(m.selectedClass.ID)
	}
	return nil
}

func (m *Model) handleStudentDeleted(msg studentDeletedMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = fmt.Errorf("deletar aluno: %w", msg.err)
		return nil
	}
	m.state = DetailsView
	m.err = nil
	m.selectedStudent = nil
	if m.selectedClass != nil {
		m.isLoading = true
		return m.fetchClassStudentsCmd(m.selectedClass.ID)
	}
	return nil
}

// Commands
type fetchedClassesMsg struct { classes []models.Class; err error }
type classCreatedMsg struct { createdClass models.Class; err error }
type classUpdatedMsg struct { updatedClass models.Class; err error }
type classDeletedMsg struct { err error }
type fetchedClassStudentsMsg struct { students []models.Student; err error }
type studentAddedMsg struct { addedStudent models.Student; err error }
type studentUpdatedMsg struct { updatedStudent models.Student; err error }
type studentDeletedMsg struct { err error }
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m *Model) fetchClassesCmd() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
	defer cancel()
	classes, err := m.classService.ListAllClasses(ctx)
	return fetchedClassesMsg{classes: classes, err: err}
}

func (m *Model) createClassCmd(name string, subjectIDStr string) tea.Cmd {
	return func() tea.Msg {
		subjectID, convErr := strconv.ParseInt(subjectIDStr, 10, 64)
		if convErr != nil { return errMsg{fmt.Errorf("ID disciplina inválido: %w", convErr)} }
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		created, err := m.classService.CreateClass(ctx, name, subjectID)
		return classCreatedMsg{createdClass: created, err: err}
	}
}

func (m *Model) updateClassCmd(id int64, name string, subjectIDStr string) tea.Cmd {
	return func() tea.Msg {
		subjectID, convErr := strconv.ParseInt(subjectIDStr, 10, 64)
		if convErr != nil { return errMsg{fmt.Errorf("ID disciplina inválido: %w", convErr)} }
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		updated, err := m.classService.UpdateClass(ctx, id, name, subjectID)
		return classUpdatedMsg{updatedClass: updated, err: err}
	}
}

func (m *Model) deleteClassCmd(id int64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		err := m.classService.DeleteClass(ctx, id)
		return classDeletedMsg{err: err}
	}
}

func (m *Model) fetchClassStudentsCmd(classID int64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		students, err := m.classService.GetStudentsByClassID(ctx, classID)
		return fetchedClassStudentsMsg{students: students, err: err}
	}
}

func (m *Model) addStudentCmd(classID int64, fullName, enrollmentID, status string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		added, err := m.classService.AddStudent(ctx, classID, fullName, enrollmentID, status)
		return studentAddedMsg{addedStudent: added, err: err}
	}
}

func (m *Model) updateStudentCmd(id int64, fullName, enrollmentID, status string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		updated, err := m.classService.UpdateStudent(ctx, id, fullName, enrollmentID, status)
		return studentUpdatedMsg{updatedStudent: updated, err: err}
	}
}

func (m *Model) deleteStudentCmd(id int64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		err := m.classService.DeleteStudent(ctx, id)
		return studentDeletedMsg{err: err}
	}
}
