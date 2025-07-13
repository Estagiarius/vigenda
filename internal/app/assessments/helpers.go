package assessments

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"vigenda/internal/models"
)

func (m *Model) setupAssessmentListTable() {
	cols := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Nome", Width: 30},
		{Title: "ID Turma", Width: 10},
		{Title: "Período", Width: 10},
		{Title: "Peso", Width: 10},
	}
	m.list = table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	m.list.SetStyles(s)
}

func (m *Model) setupGradeEntryTable() {
	cols := []table.Column{
		{Title: "ID Aluno", Width: 10},
		{Title: "Nome", Width: 30},
		{Title: "Nota", Width: 15},
	}
	m.gradesTable = table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	m.gradesTable.SetStyles(s)
}

func (m *Model) setupFormInputs(count int) {
	m.formInputs.inputs = make([]textinput.Model, count)
	for i := 0; i < count; i++ {
		ti := textinput.New()
		ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		ti.CharLimit = 100
		m.formInputs.inputs[i] = ti
	}
}

func (m *Model) prepareAssessmentForm(assessment *models.Assessment) {
	m.setupFormInputs(4)
	placeholders := []string{"Nome da Avaliação", "ID da Turma", "Período (ex: 1)", "Peso (ex: 3.0)"}
	for i, p := range placeholders {
		m.formInputs.inputs[i].Placeholder = p
		if assessment != nil {
			switch i {
			case 0:
				m.formInputs.inputs[i].SetValue(assessment.Name)
			case 1:
				m.formInputs.inputs[i].SetValue(strconv.FormatInt(assessment.ClassID, 10))
			case 2:
				m.formInputs.inputs[i].SetValue(strconv.Itoa(assessment.Term))
			case 3:
				m.formInputs.inputs[i].SetValue(fmt.Sprintf("%.2f", assessment.Weight))
			}
		}
	}

	if assessment == nil {
		m.state = CreateAssessmentView
	} else {
		m.state = EditAssessmentView
	}
	m.formInputs.focusIndex = 0
	m.formInputs.inputs[0].Focus()
	m.err = nil
}

func (m *Model) updateFormInputs(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.formInputs.inputs))
	// Only update the focused input
	if m.formInputs.focusIndex < len(m.formInputs.inputs) {
		m.formInputs.inputs[m.formInputs.focusIndex], cmds[m.formInputs.focusIndex] = m.formInputs.inputs[m.formInputs.focusIndex].Update(msg)
	}
	return cmds
}


func (m *Model) nextFormInput() {
    if m.formInputs.focusIndex < len(m.formInputs.inputs) {
	    m.formInputs.inputs[m.formInputs.focusIndex].Blur()
    }
	m.formInputs.focusIndex = (m.formInputs.focusIndex + 1) % (len(m.formInputs.inputs) + 1) // +1 for submit
	if m.formInputs.focusIndex < len(m.formInputs.inputs) {
		m.formInputs.inputs[m.formInputs.focusIndex].Focus()
	}
}

func (m *Model) prevFormInput() {
	if m.formInputs.focusIndex < len(m.formInputs.inputs) {
		m.formInputs.inputs[m.formInputs.focusIndex].Blur()
	}
	m.formInputs.focusIndex--
	if m.formInputs.focusIndex < 0 {
		m.formInputs.focusIndex = len(m.formInputs.inputs)
	}
	if m.formInputs.focusIndex < len(m.formInputs.inputs) {
		m.formInputs.inputs[m.formInputs.focusIndex].Focus()
	}
}

func (m *Model) submitAssessmentForm() tea.Cmd {
	name := m.formInputs.inputs[0].Value()
	classIDStr := m.formInputs.inputs[1].Value()
	termStr := m.formInputs.inputs[2].Value()
	weightStr := m.formInputs.inputs[3].Value()

	if name == "" || classIDStr == "" || termStr == "" || weightStr == "" {
		m.err = fmt.Errorf("todos os campos são obrigatórios")
		return nil
	}

	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil {
		m.err = fmt.Errorf("ID da turma inválido")
		return nil
	}

	term, err := strconv.Atoi(termStr)
	if err != nil {
		m.err = fmt.Errorf("período inválido")
		return nil
	}

	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		m.err = fmt.Errorf("peso inválido")
		return nil
	}

	m.isLoading = true
	if m.state == CreateAssessmentView {
		return m.createAssessmentCmd(name, classID, term, weight)
	}
	if m.state == EditAssessmentView && m.selectedAssessment != nil {
		return m.updateAssessmentCmd(m.selectedAssessment.ID, name, classID, term, weight)
	}
	return nil
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	contentHeight := height - 5
	if contentHeight < 1 {
		contentHeight = 1
	}
	tableWidth := width - 4
	if tableWidth < 20 {
		tableWidth = 20
	}

	m.list.SetHeight(contentHeight)
	m.list.SetWidth(tableWidth)
	m.gradesTable.SetHeight(contentHeight)
	m.gradesTable.SetWidth(tableWidth)

	inputWidth := width - 8
	if inputWidth < 20 {
		inputWidth = 20
	}
	for i := range m.formInputs.inputs {
		m.formInputs.inputs[i].Width = inputWidth
	}
}

func (m *Model) IsFocused() bool {
	return m.state == CreateAssessmentView || m.state == EditAssessmentView
}
