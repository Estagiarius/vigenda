package assessments

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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

// ViewState defines the current state of the assessments view
type ViewState int

const (
	ListView ViewState = iota // Main action list
	CreateAssessmentView
	EnterGradesView // This will be complex: needs to list students, then input grades
	FinalGradesView
	ListAssessmentsView // For showing assessments in a table
)

// Model represents the assessments management model.
type Model struct {
	assessmentService service.AssessmentService
	classService      service.ClassService // May need for student listing
	state             ViewState

	list list.Model // For actions or selecting assessments/classes
	table      table.Model // For displaying assessments or students
	textInputs []textinput.Model
	focusIndex int

	// Data
	assessments     []models.Assessment
	currentClassID  *int64 // For context when creating or listing assessments for a class
	currentAssessmentID *int64 // For context when entering grades or calculating average
	studentsForGrading []models.Student // For EnterGradesView
	gradesInput        map[int64]textinput.Model // studentID -> textinput for grade
	gradeFocusIndex int // New: To track focus on grade inputs

	// Popup state
	isPopupVisible bool
	popup          popupModel

	isLoading bool
	err       error
	message   string

	width  int
	height int
}

// popupModel holds the state for the calculation settings popup.
type popupModel struct {
	inputs     []textinput.Model
	focusIndex int
}

// --- Messages ---
type assessmentsLoadedMsg struct {
	assessments []models.Assessment
	err         error
}
type assessmentCreatedMsg struct {
	assessment models.Assessment
	err        error
}
type studentsForGradingLoadedMsg struct {
	students []models.Student
	assessmentName string
	err      error
}
type gradesEnteredMsg struct {
	err error
}

type assessmentDeletedMsg struct {
	assessmentID int64
	err          error
}

type studentsForFinalGradesLoadedMsg struct {
	students []models.Student
	err      error
}

type finalGradesEnteredMsg struct {
	err error
}

type classAverageCalculatedMsg struct {
	averages map[int64]float64
	err      error
}


// --- Cmds ---
func (m *Model) calculateAverageCmd(classID int64, termsStr string) tea.Cmd {
	return func() tea.Msg {
		var terms []int
		if termsStr != "" {
			parts := strings.Split(termsStr, ",")
			for _, p := range parts {
				term, err := strconv.Atoi(strings.TrimSpace(p))
				if err != nil {
					return classAverageCalculatedMsg{err: fmt.Errorf("período inválido: %s", p)}
				}
				terms = append(terms, term)
			}
		}

		averages, err := m.assessmentService.CalculateClassAverage(context.Background(), classID, terms)
		if err != nil {
			return classAverageCalculatedMsg{err: err}
		}
		return classAverageCalculatedMsg{averages: averages}
	}
}

func (m *Model) deleteAssessmentCmd(assessmentID int64) tea.Cmd {
	return func() tea.Msg {
		err := m.assessmentService.DeleteAssessment(context.Background(), assessmentID)
		return assessmentDeletedMsg{assessmentID: assessmentID, err: err}
	}
}

func (m *Model) loadAssessmentsCmd() tea.Cmd {
	m.isLoading = true
	return func() tea.Msg {
		assessments, err := m.assessmentService.ListAllAssessments(context.Background())
		return assessmentsLoadedMsg{assessments: assessments, err: err}
	}
}


func New(assessmentService service.AssessmentService, classService service.ClassService) *Model {
	actionItems := []list.Item{
		actionItem{title: "Listar Avaliações", description: "Visualizar todas as avaliações (pode pedir turma)."},
		actionItem{title: "Criar Nova Avaliação", description: "Adicionar uma nova avaliação para uma turma."},
		actionItem{title: "Lançar Notas", description: "Lançar/editar notas de alunos para uma avaliação."},
		actionItem{title: "Lançar/Calcular Notas Finais", description: "Lançar manualmente ou calcular a média final de uma turma."},
	}
	l := list.New(actionItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Gerenciar Avaliações e Notas"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	l.AdditionalShortHelpKeys = func() []key.Binding{
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "selecionar")),
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "voltar")),
		}
	}

	cols := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Nome", Width: 30},
		{Title: "ID Turma", Width: 10},
		{Title: "Período", Width: 10},
		{Title: "Peso", Width: 8},
	}
	tbl := table.New(table.WithColumns(cols), table.WithFocused(true), table.WithHeight(10))
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	tbl.SetStyles(s)

	inputs := make([]textinput.Model, 4) // Max 4 inputs for CreateAssessment
	for i := range inputs {
		ti := textinput.New()
		ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		ti.CharLimit = 64
		inputs[i] = ti
	}

	return &Model{ // Ensure this returns a pointer
		assessmentService: assessmentService,
		classService:      classService,
		state:             ListView,
		list:              l,
		table:             tbl,
		textInputs:        inputs,
		isLoading:         false,
		gradesInput:       make(map[int64]textinput.Model),
	}
}

func (m *Model) loadStudentsForFinalGradesCmd(classID int64) tea.Cmd {
	return func() tea.Msg {
		students, err := m.classService.GetStudentsByClassID(context.Background(), classID)
		if err != nil {
			return studentsForFinalGradesLoadedMsg{err: err}
		}
		return studentsForFinalGradesLoadedMsg{students: students, err: nil}
	}
}

func (m *Model) submitFinalGradesCmd() tea.Cmd {
	if m.currentClassID == nil {
		return func() tea.Msg { return finalGradesEnteredMsg{err: fmt.Errorf("ID da turma não definido")} }
	}

	grades := make(map[int64]float64)
	for studentID, ti := range m.gradesInput {
		gradeStr := ti.Value()
		if gradeStr == "" {
			continue
		}
		grade, err := strconv.ParseFloat(gradeStr, 64)
		if err != nil {
			return func() tea.Msg { return finalGradesEnteredMsg{err: fmt.Errorf("nota inválida para aluno ID %d: '%s'", studentID, gradeStr)} }
		}
		grades[studentID] = grade
	}

	if len(grades) == 0 {
		return func() tea.Msg { return finalGradesEnteredMsg{err: fmt.Errorf("nenhuma nota foi inserida")} }
	}

	return func() tea.Msg {
		err := m.assessmentService.EnterFinalGrades(context.Background(), *m.currentClassID, grades)
		return finalGradesEnteredMsg{err: err}
	}
}

// Changed to pointer receiver
func (m *Model) Init() tea.Cmd {
	// It's good practice to ensure fields are in a known state at Init.
	// Many of these are already set by New or resetForms, but explicit here is fine.
	m.state = ListView
	m.err = nil
	m.message = ""
	m.currentClassID = nil
	m.currentAssessmentID = nil
	m.list.Select(-1) // Deselect any previous action
	return nil // No initial data loading from main menu
}

type actionItem struct { // Re-defined locally, or use a shared TUI components package
	title, description string
}

func (i actionItem) Title() string       { return i.title }
func (i actionItem) Description() string { return i.description }
func (i actionItem) FilterValue() string { return i.title }

// Changed to pointer receiver
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) {
			if m.state == ListView {
				return m, nil // Let parent model handle 'esc' from main action list
			}
			// Go back to previous state or main action list
			if m.state == EnterGradesView && len(m.studentsForGrading) > 0 { // If in grade entry, Esc might go to assessment selection or main list
				m.state = ListView // Simplified: back to main action list
				m.list.Title = "Gerenciar Avaliações e Notas"
			} else {
				m.state = ListView
			}
			m.err = nil
			m.message = ""
			m.resetForms()
			m.list.Select(-1) // Deselect
			return m, nil
		}

		switch m.state {
		case ListView:
			// m.list, cmd = m.list.Update(msg)
			var updatedList list.Model
			updatedList, cmd = m.list.Update(msg)
			m.list = updatedList
			cmds = append(cmds, cmd)

			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				selected, ok := m.list.SelectedItem().(actionItem)
				if ok {
					m.err = nil
					m.message = ""
					m.resetForms()
					switch selected.title {
					case "Listar Avaliações":
						m.state = ListAssessmentsView
						cmds = append(cmds, m.loadAssessmentsCmd())
					case "Criar Nova Avaliação":
						m.state = CreateAssessmentView
						m.setupCreateAssessmentForm()
					case "Lançar Notas":
						m.state = EnterGradesView
						m.setupEnterAssessmentIDForm("Lançar Notas para Avaliação ID:")
					case "Calcular Média da Turma":
						m.state = FinalGradesView
						m.setupEnterClassIDForm("Lançar/Calcular Notas Finais para Turma ID:")
					}
				}
			}

		case ListAssessmentsView: // Navigating the table of assessments
			if key.Matches(msg, key.NewBinding(key.WithKeys("d"))) {
				selectedRow := m.table.SelectedRow()
				if len(selectedRow) > 0 {
					assessmentID, _ := strconv.ParseInt(selectedRow[0], 10, 64)
					m.isLoading = true
					m.message = fmt.Sprintf("Deletando avaliação ID %d...", assessmentID)
					cmds = append(cmds, m.deleteAssessmentCmd(assessmentID))
				}
			} else if key.Matches(msg, key.NewBinding(key.WithKeys("e"))) {
				// Placeholder for edit functionality
			} else {
				var updatedTable table.Model
				updatedTable, cmd = m.table.Update(msg)
				m.table = updatedTable
				cmds = append(cmds, cmd)
			}

		case CreateAssessmentView:
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				if m.focusIndex == len(m.textInputs) { // "Submit"
					m.isLoading = true
					cmds = append(cmds, m.submitCreateAssessmentFormCmd())
				} else {
					cmds = append(cmds, m.updateFocus())
				}
			} else {
				cmds = append(cmds, m.updateFormInputs(msg))
			}

		case EnterGradesView:
			if len(m.studentsForGrading) == 0 { // Still asking for Assessment ID
				m.err = nil // Clear previous errors on new keypress
				if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
					if m.focusIndex == 1 { // Submit button for ID input
						m.isLoading = true
						assessmentIDStr := m.textInputs[0].Value()
						assessmentID, err := strconv.ParseInt(assessmentIDStr, 10, 64)
						if err != nil {
							m.err = fmt.Errorf("ID da Avaliação inválido: %w", err)
							m.isLoading = false
						} else {
							m.currentAssessmentID = &assessmentID
							cmds = append(cmds, m.loadStudentsForGradingCmd(assessmentID))
						}
					} else { // Focus is on the input field itself
						m.textInputs[0], cmd = m.textInputs[0].Update(msg)
						cmds = append(cmds, cmd)
					}
				} else { // Other keys for the input field
					m.textInputs[0], cmd = m.textInputs[0].Update(msg)
					cmds = append(cmds, cmd)
				}
			} else { // Displaying students and grade inputs
				cmds = append(cmds, m.updateGradeInputs(msg))
			}

		case FinalGradesView:
			if len(m.studentsForGrading) == 0 { // Asking for Class ID
				if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
					if m.focusIndex == 1 { // Submit button for ID input
						m.isLoading = true
						classIDStr := m.textInputs[0].Value()
						classID, err := strconv.ParseInt(classIDStr, 10, 64)
						if err != nil {
							m.err = fmt.Errorf("ID da Turma inválido: %w", err)
							m.isLoading = false
						} else {
							m.currentClassID = &classID
							cmds = append(cmds, m.loadStudentsForFinalGradesCmd(classID))
						}
					} else { // Focus on input
						m.textInputs[0], cmd = m.textInputs[0].Update(msg)
						cmds = append(cmds, cmd)
					}
				} else { // Other keys for input
					m.textInputs[0], cmd = m.textInputs[0].Update(msg)
					cmds = append(cmds, cmd)
				}
			} else {
				cmds = append(cmds, m.updateGradeInputs(msg))
			}
		}

	if m.isPopupVisible {
		cmd = m.updatePopup(msg)
		return m, cmd
	}

	// Handle async results
	case assessmentsLoadedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.assessments = msg.assessments
			rows := make([]table.Row, len(m.assessments))
			for i, asm := range m.assessments {
				rows[i] = table.Row{
					fmt.Sprintf("%d", asm.ID),
					asm.Name,
					fmt.Sprintf("%d", asm.ClassID),
					fmt.Sprintf("%d", asm.Term),
					fmt.Sprintf("%.1f", asm.Weight),
				}
			}
			m.table.SetRows(rows)
		}

	case assessmentCreatedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.message = fmt.Sprintf("Avaliação '%s' criada com sucesso!", msg.assessment.Name)
			m.state = ListView // Go back to action list
			m.list.Select(-1)
		}

	case studentsForGradingLoadedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.studentsForGrading = msg.students
			m.message = fmt.Sprintf("Alunos carregados para avaliação: %s. Insira as notas.", msg.assessmentName)
			m.gradesInput = make(map[int64]textinput.Model)
			for i, s := range msg.students { // Use index for focus logic if needed
				ti := textinput.New()
				ti.Placeholder = "Nota (ex: 7.5)"
				ti.CharLimit = 5
				ti.Width = 10
				m.gradesInput[s.ID] = ti
				if i == 0 { // Focus the first grade input
					// This needs a proper focus management system for multiple inputs
				}
			}
		}

	case gradesEnteredMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.message = "Notas lançadas com sucesso!"
			m.state = ListView
			m.list.Select(-1)
			m.studentsForGrading = nil
			m.gradesInput = make(map[int64]textinput.Model)
		}

	case assessmentDeletedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = fmt.Errorf("Erro ao deletar avaliação ID %d: %w", msg.assessmentID, msg.err)
		} else {
			m.message = fmt.Sprintf("Avaliação ID %d deletada com sucesso.", msg.assessmentID)
			// Refresh the list
			cmds = append(cmds, m.loadAssessmentsCmd())
		}

	case finalGradesEnteredMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.message = "Notas finais salvas com sucesso!"
			m.state = ListView
			m.list.Select(-1)
			m.studentsForGrading = nil
			m.gradesInput = make(map[int64]textinput.Model)
		}

	case classAverageCalculatedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.message = "Médias calculadas com sucesso!"
			for studentID, avg := range msg.averages {
				if ti, ok := m.gradesInput[studentID]; ok {
					ti.SetValue(fmt.Sprintf("%.2f", avg))
					ti.Blur() // Disable editing after calculation
					m.gradesInput[studentID] = ti
				}
			}
		}

	case studentsForFinalGradesLoadedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.studentsForGrading = msg.students
			m.message = "Insira as notas finais."
			m.gradesInput = make(map[int64]textinput.Model)
			for _, s := range msg.students {
				ti := textinput.New()
				ti.Placeholder = "Nota Final"
				ti.CharLimit = 5
				ti.Width = 10
				m.gradesInput[s.ID] = ti
			}
		}

	case error:
		m.err = msg
		m.isLoading = false

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height) // Use the SetSize method
	}

	// Update table if it's the component in focus (e.g. for scrolling)
	if m.state == ListAssessmentsView && m.table.Focused() { // Check if table is focused
		var updatedTable table.Model
		updatedTable, cmd = m.table.Update(msg)
		m.table = updatedTable
		cmds = append(cmds, cmd)
	}


	return m, tea.Batch(cmds...)
}

// Changed to pointer receiver
func (m *Model) View() string {
	var b strings.Builder

	if m.isLoading {
		b.WriteString("Carregando...")
		return baseStyle.Render(b.String())
	}
	if m.err != nil {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Erro: %v\n\n", m.err)))
	}
	if m.message != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(fmt.Sprintf("%s\n\n", m.message)))
	}

	var mainContent string
	switch m.state {
	case ListView:
		b.WriteString(m.list.View())

	case ListAssessmentsView:
		b.WriteString("Avaliações Cadastradas:\n")
		b.WriteString(m.table.View())
		b.WriteString("\n(Navegue com ↑/↓, 'esc' para voltar às ações)")

	case CreateAssessmentView:
		b.WriteString("Criar Nova Avaliação\n\n")
		for i := range m.textInputs {
			b.WriteString(m.textInputs[i].View() + "\n")
		}
		submitButton := "[ Criar Avaliação ]"
		if m.focusIndex == len(m.textInputs) {
			submitButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(submitButton)
		}
		b.WriteString("\n" + submitButton + "\n\n")
		b.WriteString("(Use Tab/Shift+Tab ou ↑/↓ para navegar, Enter para submeter, Esc para cancelar)")

	case EnterGradesView:
		if len(m.studentsForGrading) == 0 { // Asking for Assessment ID
			b.WriteString("Lançar Notas para Avaliação\n\n")
			b.WriteString(m.textInputs[0].View() + "\n") // Assessment ID input
			submitButton := "[ Carregar Alunos ]"
			if m.focusIndex == 1 { // Assuming one input + submit
				submitButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(submitButton)
			}
			b.WriteString("\n" + submitButton + "\n\n")
		} else { // Displaying students for grade entry
			b.WriteString(fmt.Sprintf("Lançando notas para Avaliação ID %d: %s\n\n", *m.currentAssessmentID, m.message))

			// The view logic will need to be more sophisticated, rendering inputs separately or
			// building a string representation that includes the input's view.
			// For now, let's build a simpler string-based view.
			b.WriteString(lipgloss.NewStyle().Bold(true).Render("Aluno                         Nota\n"))
			b.WriteString(strings.Repeat("-", 45) + "\n")
			for _, s := range m.studentsForGrading {
				gradeInputView := ""
				if ti, ok := m.gradesInput[s.ID]; ok {
					gradeInputView = ti.View()
				}
				// Simple layout, can be improved with lipgloss.JoinHorizontal, etc.
				b.WriteString(fmt.Sprintf("%-30s %s\n", s.FullName, gradeInputView))
			}

			b.WriteString("\n[ Salvar Notas (Ctrl+S) ] [ Cancelar (Esc) ]\n")
			b.WriteString("Use ↑/↓ para navegar, Enter/Tab para editar, Esc para sair da edição.\n")
		}

	case FinalGradesView:
		if len(m.studentsForGrading) == 0 { // Still asking for Class ID
			b.WriteString("Lançar/Calcular Notas Finais da Turma\n\n")
			b.WriteString(m.textInputs[0].View() + "\n") // Class ID input
			submitButton := "[ Carregar Alunos ]"
			if m.focusIndex == 1 { // Assuming one input + submit
				submitButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(submitButton)
			}
			b.WriteString("\n" + submitButton + "\n\n")
		} else { // Displaying students for final grade entry
			b.WriteString(fmt.Sprintf("Lançando notas finais para Turma ID %d\n\n", *m.currentClassID))
			b.WriteString(lipgloss.NewStyle().Bold(true).Render("Aluno                         Nota Final\n"))
			b.WriteString(strings.Repeat("-", 45) + "\n")
			for _, s := range m.studentsForGrading {
				gradeInputView := ""
				if ti, ok := m.gradesInput[s.ID]; ok {
					gradeInputView = ti.View()
				}
				b.WriteString(fmt.Sprintf("%-30s %s\n", s.FullName, gradeInputView))
			}
			b.WriteString("\n[ Salvar (Ctrl+S) ] [ Calcular Média (Ctrl+M) ] [ Cancelar (Esc) ]\n")
		}

	default:
		b.WriteString("Visualização de Avaliações Desconhecida")
	}

	mainContent = b.String()

	if m.isPopupVisible {
		popupContent := lipgloss.NewStyle().
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("69")).
			Padding(1).
			Render(m.popup.inputs[0].View() + "\n\n[ Calcular ] [ Cancelar (Esc) ]")

		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			popupContent,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("240")),
		)
	}

	return baseStyle.Render(mainContent)
}

// --- Form Setup and Submission Logic ---
// Changed to pointer receiver
func (m *Model) resetForms() {
	for i := range m.textInputs {
		m.textInputs[i].Reset()
		m.textInputs[i].Blur()
	}
	m.gradesInput = make(map[int64]textinput.Model)
	m.studentsForGrading = nil
	m.focusIndex = 0
	m.err = nil
	m.message = ""
}

func (m *Model) setupCreateAssessmentForm() {
	m.focusIndex = 0
	m.textInputs = make([]textinput.Model, 4) // Name, ClassID, Term, Weight

	placeholders := []string{"Nome da Avaliação", "ID da Turma", "Período (ex: 1)", "Peso (ex: 3.0)"}
	validators := []func(string)error{nil, isNumber, isNumber, isFloatOrEmpty}

	for i, p := range placeholders {
		m.textInputs[i] = textinput.New()
		m.textInputs[i].Placeholder = p
		m.textInputs[i].CharLimit = 50
		m.textInputs[i].Width = m.width / 2
		if validators[i] != nil {
			m.textInputs[i].Validate = validators[i]
		}
	}
	m.textInputs[0].Focus()
	m.updateInputFocusStyle()
}

// Changed to pointer receiver
func (m *Model) setupEnterAssessmentIDForm(prompt string) {
	m.focusIndex = 0
	m.textInputs = make([]textinput.Model, 1)
	m.textInputs[0] = textinput.New()
	m.textInputs[0].Placeholder = prompt // "ID da Avaliação"
	m.textInputs[0].Focus()
	m.textInputs[0].CharLimit = 10
	m.textInputs[0].Width = m.width / 2
	m.textInputs[0].Validate = isNumber
	m.updateInputFocusStyle()
}

// Changed to pointer receiver
func (m *Model) setupEnterClassIDForm(prompt string) {
	m.focusIndex = 0
	m.textInputs = make([]textinput.Model, 1)
	m.textInputs[0] = textinput.New()
	m.textInputs[0].Placeholder = prompt // "ID da Turma"
	m.textInputs[0].Focus()
	m.textInputs[0].CharLimit = 10
	m.textInputs[0].Width = m.width / 2
	m.textInputs[0].Validate = isNumber
	m.updateInputFocusStyle()
}

func (m *Model) updateFocus() tea.Cmd {
	// Determine number of active inputs for current form
	numInputs := 0
	if m.state == CreateAssessmentView {
		numInputs = 4
	} else if (m.state == EnterGradesView && len(m.studentsForGrading) == 0) || m.state == FinalGradesView {
		numInputs = 1 // Single ID input
	} else if m.state == EnterGradesView && len(m.studentsForGrading) > 0 {
		// Complex: focus between grade inputs and a submit button
		// This needs a more specific focus management system.
		// For now, this generic updateFocus won't apply well here.
		return nil
	}


	m.focusIndex = (m.focusIndex + 1) % (numInputs + 1) // +1 for submit "button"
	return m.updateInputFocusStyle()
}

// Changed to pointer receiver
func (m *Model) updateInputFocusStyle() tea.Cmd {
	numInputs := 0
	// Determine active inputs based on state
	if m.state == CreateAssessmentView {
		numInputs = 4
	} else if (m.state == EnterGradesView && len(m.studentsForGrading) == 0) || m.state == FinalGradesView {
		numInputs = 1
	} // Other states might not use these textInputs directly or have their own focus logic


	cmds := make([]tea.Cmd, numInputs)
	for i := 0; i < numInputs; i++ {
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

// Changed to pointer receiver
func (m *Model) updateFormInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	// Update focused text input
	activeInputs := 0
	if m.state == CreateAssessmentView {
		activeInputs = 4
	}
	if (m.state == EnterGradesView && len(m.studentsForGrading) == 0) || m.state == FinalGradesView {
		activeInputs = 1
	}


	if m.focusIndex < activeInputs {
		var cmd tea.Cmd
		m.textInputs[m.focusIndex], cmd = m.textInputs[m.focusIndex].Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func (m *Model) submitCreateAssessmentFormCmd() tea.Cmd {
	return func() tea.Msg {
		name := m.textInputs[0].Value()
		classIDStr := m.textInputs[1].Value()
		termStr := m.textInputs[2].Value()
		weightStr := m.textInputs[3].Value()

		if name == "" || classIDStr == "" || termStr == "" || weightStr == "" {
			return assessmentCreatedMsg{err: fmt.Errorf("todos os campos são obrigatórios")}
		}

		classID, err := strconv.ParseInt(classIDStr, 10, 64)
		if err != nil {
			return assessmentCreatedMsg{err: fmt.Errorf("ID da turma inválido: '%s'", classIDStr)}
		}

		term, err := strconv.Atoi(termStr)
		if err != nil {
			return assessmentCreatedMsg{err: fmt.Errorf("período inválido: '%s'", termStr)}
		}

		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			return assessmentCreatedMsg{err: fmt.Errorf("peso inválido: '%s'", weightStr)}
		}

		asm, err := m.assessmentService.CreateAssessment(context.Background(), name, classID, term, weight)
		return assessmentCreatedMsg{assessment: asm, err: err}
	}
}

func (m *Model) loadStudentsForGradingCmd(assessmentID int64) tea.Cmd {
	return func() tea.Msg {
		students, assessment, err := m.assessmentService.GetStudentsForGrading(context.Background(), assessmentID)
		if err != nil {
			return studentsForGradingLoadedMsg{err: err}
		}
		return studentsForGradingLoadedMsg{
			students:       students,
			assessmentName: assessment.Name,
			err:            nil,
		}
	}
}

func (m *Model) submitGradesCmd() tea.Cmd {
	if m.currentAssessmentID == nil {
		return func() tea.Msg { return gradesEnteredMsg{err: fmt.Errorf("ID da avaliação não definido")} }
	}

	grades := make(map[int64]float64)
	for studentID, ti := range m.gradesInput {
		gradeStr := ti.Value()
		if gradeStr == "" { continue } // Skip empty grades or handle as 0?
		grade, err := strconv.ParseFloat(gradeStr, 64)
		if err != nil {
			return func() tea.Msg { return gradesEnteredMsg{err: fmt.Errorf("nota inválida para aluno ID %d: '%s'", studentID, gradeStr)} }
		}
		grades[studentID] = grade
	}

	if len(grades) == 0 {
		return func() tea.Msg { return gradesEnteredMsg{err: fmt.Errorf("nenhuma nota foi inserida")} }
	}

	return func() tea.Msg {
		err := m.assessmentService.EnterGrades(context.Background(), *m.currentAssessmentID, grades)
		return gradesEnteredMsg{err: err}
	}
}


// --- Helpers & Validators ---
func isNumber(s string) error {
	if s == "" { return nil }
	if _, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("deve ser um número inteiro")
	}
	return nil
}
func isFloatOrEmpty(s string) error {
	if s == "" { return nil }
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return fmt.Errorf("deve ser um número (ex: 7.5)")
	}
	return nil
}

// SetSize method was already using a pointer receiver, which is correct.
func (m *Model) SetSize(width, height int) {
	m.width = width - baseStyle.GetHorizontalFrameSize()
	m.height = height - baseStyle.GetVerticalFrameSize() -1 // Adjusted for potential message line

	// Adjust list (main action list)
	listTitleHeight := lipgloss.Height(m.list.Title)
	// Assuming some help text height for the list view if applicable
	listHelpHeight := 2
	availableHeightForList := m.height - listTitleHeight - listHelpHeight
	if availableHeightForList < 0 { availableHeightForList = 0 }
	m.list.SetSize(m.width, availableHeightForList)

	// Adjust table (for listing assessments)
	// Assuming table has a title/header of its own if m.state == ListAssessmentsView
	tableHeaderHeight := 1 // if table has its own title line rendered by the model's View
	tableHelpHeight := 1   // if table view has help text
	availableHeightForTable := m.height - tableHeaderHeight - tableHelpHeight
	if availableHeightForTable < 5 { availableHeightForTable = 5} // Min height for table
	m.table.SetWidth(m.width)
	m.table.SetHeight(availableHeightForTable)


	// Adjust textInputs based on current state or a general approach
	// This width calculation can be dynamic based on form structure in View()
	inputRegionWidth := m.width - 4 // General padding for forms
	if inputRegionWidth < 20 { inputRegionWidth = 20}

	for i := range m.textInputs {
		// Check if textInput is actually part of the current form to avoid nil pointer if textInputs is resized
		if i < len(m.textInputs) && m.textInputs[i].Placeholder != "" { // Basic check if it's an active input
			m.textInputs[i].Width = inputRegionWidth
		}
	}

	// Adjust gradesInput (these are typically smaller)
	for studentID := range m.gradesInput {
		ti := m.gradesInput[studentID]
		ti.Width = 10 // Keep grade inputs small and fixed width
		m.gradesInput[studentID] = ti
	}
}

func (m *Model) setupPopup() {
	m.popup.focusIndex = 0
	m.popup.inputs = make([]textinput.Model, 1) // Just one input for now: terms
	ti := textinput.New()
	ti.Placeholder = "Períodos a incluir (ex: 1,2,3)"
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 30
	m.popup.inputs[0] = ti
}

func (m *Model) updatePopup(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, key.NewBinding(key.WithKeys("esc"))):
			m.isPopupVisible = false
			return nil
		case key.Matches(keyMsg, key.NewBinding(key.WithKeys("enter"))):
			m.isPopupVisible = false
			m.isLoading = true
			termsStr := m.popup.inputs[0].Value()
			return m.calculateAverageCmd(*m.currentClassID, termsStr)
		}
	}

	// Handle text input
	var cmd tea.Cmd
	m.popup.inputs[0], cmd = m.popup.inputs[0].Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) updateGradeInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	if key.Matches(keyMsg, key.NewBinding(key.WithKeys("ctrl+m"))) {
		m.isPopupVisible = true
		m.setupPopup()
		return nil
	}

	// When not in text input mode, navigate between students
	if m.gradeFocusIndex < len(m.studentsForGrading) {
		focusedStudent := m.studentsForGrading[m.gradeFocusIndex]
		if ti, exists := m.gradesInput[focusedStudent.ID]; !exists || !ti.Focused() {
			switch {
			case key.Matches(keyMsg, key.NewBinding(key.WithKeys("up"))):
				if m.gradeFocusIndex > 0 {
					m.gradeFocusIndex--
				}
			case key.Matches(keyMsg, key.NewBinding(key.WithKeys("down"))):
				if m.gradeFocusIndex < len(m.studentsForGrading)-1 {
					m.gradeFocusIndex++
				}
			case key.Matches(keyMsg, key.NewBinding(key.WithKeys("enter"), key.WithKeys("tab"))):
				// Focus the text input for the current student
				if ti, exists := m.gradesInput[focusedStudent.ID]; exists {
					cmds = append(cmds, ti.Focus())
					m.gradesInput[focusedStudent.ID] = ti
				}
			case key.Matches(keyMsg, key.NewBinding(key.WithKeys("ctrl+s"))):
				m.isLoading = true
				if m.state == EnterGradesView {
					cmds = append(cmds, m.submitGradesCmd())
				} else if m.state == FinalGradesView {
					cmds = append(cmds, m.submitFinalGradesCmd())
				}
			}
		}
	}

	// Update the focused text input
	for i, student := range m.studentsForGrading {
		if ti, exists := m.gradesInput[student.ID]; exists && ti.Focused() {
			// Handle Enter and Esc within the text input
			if key.Matches(keyMsg, key.NewBinding(key.WithKeys("enter"))) {
				ti.Blur()
				// Move to next student
				if i < len(m.studentsForGrading)-1 {
					m.gradeFocusIndex = i + 1
				}
			} else if key.Matches(keyMsg, key.NewBinding(key.WithKeys("esc"))) {
				ti.Blur()
			} else {
				var cmd tea.Cmd
				m.gradesInput[student.ID], cmd = ti.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	// Refocus based on gradeFocusIndex
	for i, student := range m.studentsForGrading {
		ti := m.gradesInput[student.ID]
		if i == m.gradeFocusIndex && !ti.Focused() {
			// Visual cue for focus, without cursor
			ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		} else {
			ti.PromptStyle = lipgloss.NewStyle()
		}
		m.gradesInput[student.ID] = ti
	}

	return tea.Batch(cmds...)
}
