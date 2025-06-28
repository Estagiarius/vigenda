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
	ClassAverageView
	ListAssessmentsView // For showing assessments in a table
)

// Model represents the assessments management model.
type Model struct {
	assessmentService service.AssessmentService
	// classService service.ClassService // May need for student listing
	state ViewState

	list       list.Model  // For actions or selecting assessments/classes
	table      table.Model // For displaying assessments or students
	textInputs []textinput.Model
	focusIndex int

	// Data
	assessments     []models.Assessment
	currentClassID  *int64 // For context when creating or listing assessments for a class
	currentAssessmentID *int64 // For context when entering grades or calculating average
	studentsForGrading []models.Student // For EnterGradesView
	gradesInput        map[int64]textinput.Model // studentID -> textinput for grade

	isLoading bool
	err       error
	message   string

	width  int
	height int
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
type classAverageCalculatedMsg struct {
	average float64
	className string // Or ClassID
	err     error
}


// --- Cmds ---
func (m *Model) loadAssessmentsCmd(classID int64) tea.Cmd {
	return func() tea.Msg {
		// Assuming a service method like ListAssessmentsByClass exists
		// If not, this needs adjustment. For now, this is a placeholder.
		// assessments, err := m.assessmentService.ListAssessmentsByClass(context.Background(), classID)
		// For a generic list, maybe ListAllAssessments if that's more appropriate first
		assessments, err := m.assessmentService.ListAllAssessments(context.Background()) // Placeholder
		return assessmentsLoadedMsg{assessments: assessments, err: err}
	}
}


func New(assessmentService service.AssessmentService /*, classService service.ClassService */) Model {
	actionItems := []list.Item{
		actionItem{title: "Listar Avaliações", description: "Visualizar todas as avaliações (pode pedir turma)."},
		actionItem{title: "Criar Nova Avaliação", description: "Adicionar uma nova avaliação para uma turma."},
		actionItem{title: "Lançar Notas", description: "Lançar/editar notas de alunos para uma avaliação."},
		actionItem{title: "Calcular Média da Turma", description: "Calcular a média geral de uma turma em avaliações."},
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

	return Model{
		assessmentService: assessmentService,
		// classService: classService,
		state:      ListView,
		list:       l,
		table:      tbl,
		textInputs: inputs,
		isLoading:  false,
		gradesInput: make(map[int64]textinput.Model),
	}
}

func (m Model) Init() tea.Cmd {
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


func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				selected, ok := m.list.SelectedItem().(actionItem)
				if ok {
					m.err = nil
					m.message = ""
					m.resetForms()
					switch selected.title {
					case "Listar Avaliações":
						m.isLoading = true
						m.state = ListAssessmentsView // Dedicated view for table
						// We might need to ask for Class ID first or show all.
						// For now, let's assume we list all.
						// A better UX would be to list classes, select one, then list its assessments.
						// This is a simplified path for now.
						cmds = append(cmds, m.loadAssessmentsCmd(0)) // 0 for all, or adapt service
					case "Criar Nova Avaliação":
						m.state = CreateAssessmentView
						m.setupCreateAssessmentForm()
					case "Lançar Notas":
						// This needs to first ask for Assessment ID
						m.state = EnterGradesView // Intermediate state to ask for ID
						m.setupEnterAssessmentIDForm("Lançar Notas para Avaliação ID:")
					case "Calcular Média da Turma":
						// This needs to first ask for Class ID
						m.state = ClassAverageView // Intermediate state to ask for ID
						m.setupEnterClassIDForm("Calcular Média para Turma ID:")
					}
				}
			}

		case ListAssessmentsView: // Navigating the table of assessments
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
			// Potentially Enter here could select an assessment for details/actions

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

		case EnterGradesView: // Could be asking for Assessment ID, or showing student list
			if len(m.studentsForGrading) == 0 { // Still asking for Assessment ID
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
							// Now load students for this assessment
							cmds = append(cmds, m.loadStudentsForGradingCmd(assessmentID))
						}
					} else {
						m.textInputs[0], cmd = m.textInputs[0].Update(msg)
						cmds = append(cmds, cmd)
					}
				} else {
					m.textInputs[0], cmd = m.textInputs[0].Update(msg)
					cmds = append(cmds, cmd)
				}
			} else { // Displaying students and grade inputs
				// This part needs a more complex navigation (table for students, inputs for grades)
				// For now, just acknowledge it's complex.
				// Handle navigation between grade inputs, and a final submit.
				// This is a placeholder for a more robust implementation.
				if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
					// Potentially submit all grades
					m.isLoading = true
					cmds = append(cmds, m.submitGradesCmd())
				} else {
					// Handle focus and input for grade fields (m.gradesInput)
					// This is simplified. A real implementation would manage focus across many inputs.
					// For now, let's assume any key press is for the "current" grade input if one were active.
				}
			}

		case ClassAverageView: // Asking for Class ID
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				if m.focusIndex == 1 { // Submit button for ID input
					m.isLoading = true
					classIDStr := m.textInputs[0].Value()
					// classID, err := strconv.ParseInt(classIDStr, 10, 64) // Comentado para evitar erro de não utilizado
					_, err := strconv.ParseInt(classIDStr, 10, 64)
					if err != nil {
						m.err = fmt.Errorf("ID da Turma inválido: %w", err)
						m.isLoading = false
					} else {
						// cmds = append(cmds, m.calculateClassAverageCmd(classID))
						// Placeholder for actual call
						m.err = fmt.Errorf("Cálculo de média da turma TUI não totalmente implementado.")
						m.isLoading = false
					}
				} else {
					m.textInputs[0], cmd = m.textInputs[0].Update(msg)
					cmds = append(cmds, cmd)
				}
			} else {
				m.textInputs[0], cmd = m.textInputs[0].Update(msg)
				cmds = append(cmds, cmd)
			}
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
			for _, s := range msg.students {
				ti := textinput.New()
				ti.Placeholder = "Nota (ex: 7.5)"
				ti.CharLimit = 5
				ti.Width = 10
				// ti.Validate = isFloatOrEmpty // TODO: Implement validator
				m.gradesInput[s.ID] = ti
			}
			// Focus the first grade input or a general "submit" area
			// This part of UI interaction is complex and simplified here.
		}

	case gradesEnteredMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.message = "Notas lançadas com sucesso!"
			m.state = ListView
			m.list.Select(-1)
			m.studentsForGrading = nil // Clear student list
			m.gradesInput = make(map[int64]textinput.Model) // Clear grade inputs
		}

	case error:
		m.err = msg
		m.isLoading = false

	case tea.WindowSizeMsg:
		m.width = msg.Width - baseStyle.GetHorizontalFrameSize()
		m.height = msg.Height - baseStyle.GetVerticalFrameSize() -1

		listHeight := m.height - lipgloss.Height(m.list.Title) - 2
		m.list.SetSize(m.width, listHeight)

		m.table.SetWidth(m.width)
		tableHeight := m.height - 6
		if tableHeight < 5 { tableHeight = 5 }
		m.table.SetHeight(tableHeight)

		inputWidth := m.width - 20 // More padding for forms
		if inputWidth < 20 { inputWidth = 20 }
		for i := range m.textInputs {
			if i < len(m.textInputs) {
				m.textInputs[i].Width = inputWidth
			}
		}
		// Also resize gradesInput if active
		for _, ti := range m.gradesInput {
			ti.Width = 10 // Fixed small width for grade inputs
		}
	}

	// Update table if it's the component in focus (e.g. for scrolling)
	if m.state == ListAssessmentsView {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}


	return m, tea.Batch(cmds...)
}


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
			b.WriteString(fmt.Sprintf("Lançando notas para Avaliação ID %d\n", *m.currentAssessmentID))
			// This needs a proper table or formatted list for students and their grade inputs
			b.WriteString("Aluno ID | Nome do Aluno        | Nota\n")
			b.WriteString(strings.Repeat("-", 40) + "\n")
			for _, s := range m.studentsForGrading {
				gradeInputView := ""
				if gi, ok := m.gradesInput[s.ID]; ok {
					gradeInputView = gi.View()
				}
				b.WriteString(fmt.Sprintf("%-8d | %-20s | %s\n", s.ID, s.FullName, gradeInputView))
			}
			b.WriteString("\n[ Submeter Todas as Notas (Enter) ] [ Cancelar (Esc) ]\n")
			b.WriteString("Navegação entre campos de nota e submissão final ainda é simplificada.\n")
		}

	case ClassAverageView: // Asking for Class ID
		b.WriteString("Calcular Média da Turma\n\n")
		b.WriteString(m.textInputs[0].View() + "\n") // Class ID input
		submitButton := "[ Calcular Média ]"
		if m.focusIndex == 1 { // Assuming one input + submit
			submitButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(submitButton)
		}
		b.WriteString("\n" + submitButton + "\n\n")

	default:
		b.WriteString("Visualização de Avaliações Desconhecida")
	}

	return baseStyle.Render(b.String())
}

// --- Form Setup and Submission Logic ---
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
	} else if (m.state == EnterGradesView && len(m.studentsForGrading) == 0) || m.state == ClassAverageView {
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

func (m *Model) updateInputFocusStyle() tea.Cmd {
	numInputs := 0
	// Determine active inputs based on state
	if m.state == CreateAssessmentView {
		numInputs = 4
	} else if (m.state == EnterGradesView && len(m.studentsForGrading) == 0) || m.state == ClassAverageView {
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

func (m *Model) updateFormInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	// Update focused text input
	activeInputs := 0
	if m.state == CreateAssessmentView { activeInputs = 4 }
	if (m.state == EnterGradesView && len(m.studentsForGrading) == 0) || m.state == ClassAverageView { activeInputs = 1}


	if m.focusIndex < activeInputs {
		var cmd tea.Cmd
		m.textInputs[m.focusIndex], cmd = m.textInputs[m.focusIndex].Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func (m *Model) submitCreateAssessmentFormCmd() tea.Cmd {
	name := m.textInputs[0].Value()
	classIDStr := m.textInputs[1].Value()
	termStr := m.textInputs[2].Value()
	weightStr := m.textInputs[3].Value()

	if name == "" || classIDStr == "" || termStr == "" || weightStr == "" {
		m.err = fmt.Errorf("Todos os campos são obrigatórios.")
		m.isLoading = false
		return nil
	}
	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil { m.err = fmt.Errorf("ID da Turma inválido"); m.isLoading = false; return nil }
	term, err := strconv.Atoi(termStr)
	if err != nil { m.err = fmt.Errorf("Período inválido"); m.isLoading = false; return nil }
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil { m.err = fmt.Errorf("Peso inválido"); m.isLoading = false; return nil }

	return func() tea.Msg {
		asm, err := m.assessmentService.CreateAssessment(context.Background(), name, classID, term, weight)
		return assessmentCreatedMsg{assessment: asm, err: err}
	}
}

func (m *Model) loadStudentsForGradingCmd(assessmentID int64) tea.Cmd {
	return func() tea.Msg {
		// This requires assessmentService to have a method like GetStudentsForAssessmentGrading
		// Which in turn might need to fetch the assessment, then its class, then students of that class.
		// For now, this is a placeholder for that complex logic.
		// students, assessmentName, err := m.assessmentService.GetStudentsAndAssessmentNameForGrading(context.Background(), assessmentID)
		// return studentsForGradingLoadedMsg{students: students, assessmentName: assessmentName, err: err}

		// Simulating a failure as the service method is complex and likely not there yet.
		return studentsForGradingLoadedMsg{err: fmt.Errorf("serviço para carregar alunos para avaliação não implementado")}
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


func (m *Model) SetSize(width, height int) {
	m.width = width - baseStyle.GetHorizontalFrameSize()
	m.height = height - baseStyle.GetVerticalFrameSize() -1

	listHeight := m.height - lipgloss.Height(m.list.Title) - 2
	m.list.SetSize(m.width, listHeight)

	m.table.SetWidth(m.width)
	tableHeight := m.height - 6
	if tableHeight < 5 { tableHeight = 5 }
	m.table.SetHeight(tableHeight)

	inputWidth := m.width / 2
	if inputWidth < 20 { inputWidth = 20}
	for i := range m.textInputs {
		if i < len(m.textInputs) {
			m.textInputs[i].Width = inputWidth
		}
	}
	for _, ti := range m.gradesInput {
		ti.Width = 10 // Keep grade inputs small
	}
}

func (m Model) IsFocused() bool {
	// Focused if in any form input state
	return m.state == CreateAssessmentView ||
	       (m.state == EnterGradesView && len(m.studentsForGrading) == 0) || // inputting assessment ID
	       (m.state == EnterGradesView && len(m.studentsForGrading) > 0) || // inputting grades (complex focus)
	       m.state == ClassAverageView // inputting class ID
}
