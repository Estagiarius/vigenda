package assessments

import (
	"fmt"
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

var (
	dbOperationTimeout = 5 * time.Second
)

type ViewState int

const (
	ActionListView ViewState = iota
	ListAssessmentsView
	CreateAssessmentView
	EditAssessmentView
	DeleteConfirmView
	EnterGradesView
)

type Model struct {
	assessmentService service.AssessmentService
	state             ViewState
	keys              KeyMap
	list              table.Model // Main table for listing assessments
	formInputs        struct {
		inputs     []textinput.Model
		focusIndex int
	}
	allAssessments     []models.Assessment
	selectedAssessment *models.Assessment

	// For grade entry
	gradingSheet *service.GradingSheet
	gradesTable  table.Model
	gradeInputs  map[int64]textinput.Model // studentID -> textinput
	gradeFocusIndex int

	isLoading bool
	width     int
	height    int
	err       error
}

func New(assessmentService service.AssessmentService) *Model {
	m := &Model{
		assessmentService: assessmentService,
		state:             ActionListView,
		keys:              DefaultKeyMap,
		isLoading:         true,
	}
	m.setupAssessmentListTable()
	m.setupGradeEntryTable()
	m.setupFormInputs(4) // For assessment form
	return m
}

func (m *Model) Init() tea.Cmd {
	m.state = ListAssessmentsView
	m.isLoading = true
	return m.fetchAssessmentsCmd
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		cmd = m.handleKeyPress(msg)
	case fetchedAssessmentsMsg:
		cmd = m.handleFetchedAssessments(msg)
	case assessmentCreatedMsg:
		cmd = m.handleAssessmentCreated(msg)
	case assessmentUpdatedMsg:
		cmd = m.handleAssessmentUpdated(msg)
	case assessmentDeletedMsg:
		cmd = m.handleAssessmentDeleted(msg)
	case fetchedGradingSheetMsg:
		cmd = m.handleFetchedGradingSheet(msg)
	case gradesEnteredMsg:
		cmd = m.handleGradesEntered(msg)
	case errMsg:
		m.err = msg.err
		m.isLoading = false
	}
	cmds = append(cmds, cmd)

	if m.state == CreateAssessmentView || m.state == EditAssessmentView {
		inputCmds := m.updateFormInputs(msg)
		cmds = append(cmds, inputCmds...)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if m.isLoading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "Carregando...")
	}
	var b strings.Builder
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).PaddingBottom(1)
		b.WriteString(errorStyle.Render(fmt.Sprintf("Erro: %v", m.err)))
	}

	switch m.state {
	case ListAssessmentsView:
		b.WriteString(lipgloss.NewStyle().Bold(true).MarginBottom(1).Render("Lista de Avaliações"))
		b.WriteString(m.list.View())
		help := "↑/↓: Navegar | Enter: Lançar Notas | n: Nova | e: Editar | d: Deletar | q/Esc: Voltar"
		b.WriteString("\n" + lipgloss.NewStyle().Faint(true).Render(help))

	case CreateAssessmentView, EditAssessmentView:
		title := "Nova Avaliação"
		if m.state == EditAssessmentView && m.selectedAssessment != nil {
			title = fmt.Sprintf("Editando Avaliação: %s", m.selectedAssessment.Name)
		}
		b.WriteString(lipgloss.NewStyle().Bold(true).MarginBottom(1).Render(title) + "\n")
		for _, input := range m.formInputs.inputs {
			b.WriteString(input.View() + "\n")
		}
		b.WriteString("\n" + lipgloss.NewStyle().Faint(true).Render("Tab: Próximo | Shift+Tab: Anterior | Enter: Salvar | Esc: Cancelar"))

	case DeleteConfirmView:
		if m.selectedAssessment != nil {
			b.WriteString(lipgloss.NewStyle().Bold(true).MarginBottom(1).Render(fmt.Sprintf("Confirmar Exclusão: %s?", m.selectedAssessment.Name)))
			b.WriteString("Esta ação não pode ser desfeita.\n\n")
			b.WriteString(lipgloss.NewStyle().Faint(true).Render("Pressione 's' para confirmar, 'n' ou 'Esc' para cancelar."))
		}

	case EnterGradesView:
		// To be implemented in the next step
		b.WriteString("Tela de Lançamento de Notas (em construção)")

	}
	return b.String()
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	if key.Matches(msg, m.keys.Quit) {
		// Let the parent model handle quitting
		return nil
	}
	switch m.state {
	case ListAssessmentsView:
		return m.handleListAssessmentsKeys(msg)
	case CreateAssessmentView, EditAssessmentView:
		return m.handleAssessmentFormKeys(msg)
	case DeleteConfirmView:
		return m.handleDeleteConfirmKeys(msg)
	case EnterGradesView:
		// To be implemented
	}
	return nil
}

// Key Handlers
func (m *Model) handleListAssessmentsKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
		m.prepareAssessmentForm(nil)
		return textinput.Blink
	case key.Matches(msg, key.NewBinding(key.WithKeys("e"))):
		if idx := m.list.Cursor(); idx >= 0 && idx < len(m.allAssessments) {
			m.selectedAssessment = &m.allAssessments[idx]
			m.prepareAssessmentForm(m.selectedAssessment)
			return textinput.Blink
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("d"))):
		if idx := m.list.Cursor(); idx >= 0 && idx < len(m.allAssessments) {
			m.selectedAssessment = &m.allAssessments[idx]
			m.state = DeleteConfirmView
			m.err = nil
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
		// Placeholder for grade entry
		if idx := m.list.Cursor(); idx >= 0 && idx < len(m.allAssessments) {
			m.selectedAssessment = &m.allAssessments[idx]
			m.state = EnterGradesView
			m.isLoading = true
			m.err = nil
			return m.fetchGradingSheetCmd(m.selectedAssessment.ID)
		}
	default:
		var tableCmd tea.Cmd
		m.list, tableCmd = m.list.Update(msg)
		cmds = append(cmds, tableCmd)
	}
	return tea.Batch(cmds...)
}

func (m *Model) handleAssessmentFormKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
		m.state = ListAssessmentsView
		m.err = nil
		m.selectedAssessment = nil
		m.list.Focus()
	case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
		if m.formInputs.focusIndex == len(m.formInputs.inputs) { // On submit
			return m.submitAssessmentForm()
		}
		m.nextFormInput()
	case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
		m.nextFormInput()
	case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
		m.prevFormInput()
	}
	return nil
}

func (m *Model) handleDeleteConfirmKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "s", "S":
		if m.selectedAssessment != nil {
			m.isLoading = true
			return m.deleteAssessmentCmd(m.selectedAssessment.ID)
		}
	case "n", "N", "esc":
		m.state = ListAssessmentsView
		m.selectedAssessment = nil
		m.err = nil
		m.list.Focus()
	}
	return nil
}

// Message Handlers
func (m *Model) handleFetchedAssessments(msg fetchedAssessmentsMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = msg.err
		return nil
	}
	m.allAssessments = msg.assessments
	rows := make([]table.Row, len(m.allAssessments))
	for i, a := range m.allAssessments {
		rows[i] = table.Row{
			fmt.Sprintf("%d", a.ID),
			a.Name,
			fmt.Sprintf("%d", a.ClassID),
			fmt.Sprintf("%d", a.Term),
			fmt.Sprintf("%.2f", a.Weight),
		}
	}
	m.list.SetRows(rows)
	return nil
}

func (m *Model) handleAssessmentCreated(msg assessmentCreatedMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = msg.err
		return nil
	}
	m.state = ListAssessmentsView
	m.err = nil
	m.list.Focus()
	return m.fetchAssessmentsCmd
}

func (m *Model) handleAssessmentUpdated(msg assessmentUpdatedMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = msg.err
		return nil
	}
	m.state = ListAssessmentsView
	m.err = nil
	m.selectedAssessment = nil
	m.list.Focus()
	return m.fetchAssessmentsCmd
}

func (m *Model) handleAssessmentDeleted(msg assessmentDeletedMsg) tea.Cmd {
	m.isLoading = false
	if msg.err != nil {
		m.err = msg.err
		return nil
	}
	m.state = ListAssessmentsView
	m.err = nil
	m.selectedAssessment = nil
	m.list.Focus()
	return m.fetchAssessmentsCmd
}

// To be implemented
func (m *Model) handleFetchedGradingSheet(msg fetchedGradingSheetMsg) tea.Cmd {
	m.isLoading = false
	m.err = msg.err
	return nil
}
func (m *Model) handleGradesEntered(msg gradesEnteredMsg) tea.Cmd {
	m.isLoading = false
	m.err = msg.err
	return nil
}
