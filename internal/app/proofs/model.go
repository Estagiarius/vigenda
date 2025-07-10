package proofs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"vigenda/internal/models"
	"vigenda/internal/service"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// ViewState defines the current state of the proofs view
type ViewState int

const (
	FormView  ViewState = iota // View for inputting proof generation criteria
	ProofView                  // View for displaying the generated proof
)

// Model represents the proof generation model.
type Model struct {
	proofService service.ProofService
	state        ViewState

	textInputs []textinput.Model // For SubjectID, Topic, Easy, Medium, Hard counts
	focusIndex int

	generatedProof []models.Question // Stores the generated questions
	isLoading      bool
	err            error
	message        string // For informational messages or context

	width  int
	height int
}

// --- Messages ---
type proofGeneratedMsg struct {
	proof []models.Question
	err   error
}

// --- Cmds ---
func (m *Model) generateProofCmd() tea.Cmd {
	subjectIDStr := m.textInputs[0].Value()
	topic := m.textInputs[1].Value() // Optional
	easyCountStr := m.textInputs[2].Value()
	mediumCountStr := m.textInputs[3].Value()
	hardCountStr := m.textInputs[4].Value()

	if subjectIDStr == "" {
		return func() tea.Msg { return proofGeneratedMsg{err: fmt.Errorf("ID da Disciplina é obrigatório")} }
	}
	subjectID, err := strconv.ParseInt(subjectIDStr, 10, 64)
	if err != nil {
		return func() tea.Msg { return proofGeneratedMsg{err: fmt.Errorf("ID da Disciplina inválido: %w", err)} }
	}

	easyCount, _ := strconv.Atoi(easyCountStr)     // Default to 0 if empty or invalid
	mediumCount, _ := strconv.Atoi(mediumCountStr) // Default to 0
	hardCount, _ := strconv.Atoi(hardCountStr)     // Default to 0

	if easyCount == 0 && mediumCount == 0 && hardCount == 0 {
		return func() tea.Msg {
			return proofGeneratedMsg{err: fmt.Errorf("pelo menos uma contagem de dificuldade deve ser maior que zero")}
		}
	}

	criteria := service.ProofCriteria{
		SubjectID:   subjectID,
		EasyCount:   easyCount,
		MediumCount: mediumCount,
		HardCount:   hardCount,
	}
	if topic != "" {
		criteria.Topic = &topic
	}

	return func() tea.Msg {
		proof, err := m.proofService.GenerateProof(context.Background(), criteria)
		return proofGeneratedMsg{proof: proof, err: err}
	}
}

func New(proofService service.ProofService) *Model { // Return *Model
	inputs := make([]textinput.Model, 5) // SubjectID, Topic, Easy, Medium, Hard
	placeholders := []string{
		"ID da Disciplina (obrigatório)",
		"Tópico (opcional)",
		"Qtd. Fáceis (ex: 5)",
		"Qtd. Médias (ex: 3)",
		"Qtd. Difíceis (ex: 2)",
	}
	charLimits := []int{10, 50, 3, 3, 3}
	validators := []func(string) error{
		isNumberOrEmpty, // Subject ID is mandatory, will be checked on submit
		nil,             // Topic is optional text
		isNumberOrEmpty, // Counts are numeric, can be empty (implies 0)
		isNumberOrEmpty,
		isNumberOrEmpty,
	}

	for i := range inputs {
		ti := textinput.New()
		ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		ti.Placeholder = placeholders[i]
		ti.CharLimit = charLimits[i]
		if validators[i] != nil {
			ti.Validate = validators[i]
		}
		inputs[i] = ti
	}

	return &Model{ // Corrected to return a pointer
		proofService: proofService,
		state:        FormView,
		textInputs:   inputs,
		isLoading:    false,
	}
}

// Changed to pointer receiver
func (m *Model) Init() tea.Cmd {
	m.state = FormView
	m.err = nil
	m.message = "Insira os critérios para gerar a prova."
	m.generatedProof = nil
	m.resetForm() // resetForm now handles focus
	// if len(m.textInputs) > 0 { // Focus is handled by resetForm
	//     return m.textInputs[0].Focus()
	// }
	return nil
}

// Changed to pointer receiver
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) {
			if m.state == ProofView { // If viewing proof, Esc goes back to form
				m.state = FormView
				m.err = nil
				m.message = "Insira os critérios para gerar a prova."
				m.generatedProof = nil // Clear previous proof
				m.resetForm()          // Reset form fields and focus
				// No need to explicitly return focus cmd here as resetForm handles it.
				return m, nil // Return m directly
			}
			// If in FormView, Esc is handled by parent model to go to main menu
			return m, nil // Return m directly
		}

		switch m.state {
		case FormView:
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				if m.focusIndex == len(m.textInputs) { // "Submit" action
					m.isLoading = true
					m.err = nil
					m.message = ""
					cmds = append(cmds, m.generateProofCmd())
				} else { // Move focus to next input
					m.focusIndex = (m.focusIndex + 1) % (len(m.textInputs) + 1) // +1 for submit "button"
					cmds = append(cmds, m.updateInputFocusStyle())
				}
			} else if key.Matches(msg, key.NewBinding(key.WithKeys("up", "shift+tab"))) {
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = len(m.textInputs) // Wrap around to submit
				}
				cmds = append(cmds, m.updateInputFocusStyle())
			} else if key.Matches(msg, key.NewBinding(key.WithKeys("down", "tab"))) {
				m.focusIndex++
				if m.focusIndex > len(m.textInputs) {
					m.focusIndex = 0 // Wrap around to first input
				}
				cmds = append(cmds, m.updateInputFocusStyle())
			} else { // Pass to focused text input
				if m.focusIndex < len(m.textInputs) {
					// m.textInputs[m.focusIndex], cmd = m.textInputs[m.focusIndex].Update(msg)
					var updatedInput textinput.Model
					updatedInput, cmd = m.textInputs[m.focusIndex].Update(msg)
					m.textInputs[m.focusIndex] = updatedInput
					cmds = append(cmds, cmd)
				}
			}

		case ProofView:
			// Currently, no interaction in proof view other than Esc
			// pass_through := true // TODO: remove this variable
			// _ = pass_through      // TODO: remove this variable
		}

	// Handle async results
	case proofGeneratedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
			m.state = FormView // Stay in form view to show error
		} else {
			if len(msg.proof) == 0 {
				m.err = fmt.Errorf("Nenhuma questão encontrada para os critérios fornecidos.")
				m.state = FormView
			} else {
				m.generatedProof = msg.proof
				m.state = ProofView
				m.message = fmt.Sprintf("Prova gerada com %d questões.", len(m.generatedProof))
			}
		}

	case error: // Generic error message
		m.err = msg
		m.isLoading = false
		m.state = FormView // Revert to form on other errors

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height) // Use SetSize method
	}

	return m, tea.Batch(cmds...)
}

// Changed to pointer receiver
func (m *Model) View() string {
	var b strings.Builder

	if m.isLoading {
		b.WriteString("Gerando prova...")
		return baseStyle.Render(b.String())
	}
	if m.err != nil {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Erro: %v\n\n", m.err)))
	}
	if m.message != "" && m.state == FormView { // Show general messages only in FormView to avoid clutter in ProofView
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(fmt.Sprintf("%s\n\n", m.message)))
	}

	switch m.state {
	case FormView:
		b.WriteString("Gerar Nova Prova\n\n")
		for i := range m.textInputs {
			b.WriteString(m.textInputs[i].View() + "\n")
		}
		submitButton := "[ Gerar Prova ]"
		if m.focusIndex == len(m.textInputs) { // If submit "button" is focused
			submitButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(submitButton)
		}
		b.WriteString("\n" + submitButton + "\n\n")
		b.WriteString("(Use Tab/Shift+Tab ou ↑/↓ para navegar, Enter para submeter, Esc para voltar ao menu principal)")

	case ProofView:
		titleStyle := lipgloss.NewStyle().Bold(true).MarginBottom(1)
		b.WriteString(titleStyle.Render(fmt.Sprintf("Prova Gerada (%d questões)", len(m.generatedProof))))
		b.WriteString("\n\n")

		questionStyle := lipgloss.NewStyle().MarginBottom(1)
		optionStyle := lipgloss.NewStyle().MarginLeft(2)
		answerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("28")).MarginLeft(2) // Green

		for i, q := range m.generatedProof {
			b.WriteString(questionStyle.Render(fmt.Sprintf("Q%d (%s, %s): %s", i+1, q.Difficulty, q.Type, q.Statement)))
			b.WriteString("\n")
			if q.Options != nil && *q.Options != "" && *q.Options != "null" {
				var opts []string
				// It's safer to check Unmarshal error
				if json.Unmarshal([]byte(*q.Options), &opts) == nil {
					for j, opt := range opts {
						b.WriteString(optionStyle.Render(fmt.Sprintf("  %c) %s", 'a'+j, opt)))
						b.WriteString("\n")
					}
				} else {
					b.WriteString(optionStyle.Render(fmt.Sprintf("  Opções: %s (formato inválido)", *q.Options)))
					b.WriteString("\n")
				}
			}
			b.WriteString(answerStyle.Render(fmt.Sprintf("   R: %s", q.CorrectAnswer)))
			b.WriteString("\n\n")
		}
		b.WriteString("\n(Pressione Esc para voltar ao formulário de geração)")

	default:
		b.WriteString("Visualização de Geração de Prova Desconhecida")
	}

	return baseStyle.Render(b.String())
}

// Changed to pointer receiver (already was, just confirming)
func (m *Model) resetForm() {
	for i := range m.textInputs {
		m.textInputs[i].Reset()
		m.textInputs[i].Blur()
	}
	m.focusIndex = 0
	if len(m.textInputs) > 0 {
		m.textInputs[0].Focus()
	}
	m.updateInputFocusStyle()
}

func (m *Model) updateInputFocusStyle() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textInputs))
	for i := 0; i < len(m.textInputs); i++ {
		if i == m.focusIndex {
			cmds[i] = m.textInputs[i].Focus()
			m.textInputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // Highlight focused
		} else {
			m.textInputs[i].Blur()
			m.textInputs[i].PromptStyle = lipgloss.NewStyle() // Default style
		}
	}
	return tea.Batch(cmds...)
}

// isNumberOrEmpty validator remains a package-level function, no receiver needed.
func isNumberOrEmpty(s string) error {
	if s == "" {
		return nil // Allow empty, implies 0 for counts
	}
	if _, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("deve ser um número")
	}
	return nil
}

// SetSize method already uses a pointer receiver.
func (m *Model) SetSize(width, height int) {
	m.width = width - baseStyle.GetHorizontalFrameSize()
	m.height = height - baseStyle.GetVerticalFrameSize() - 1 // Adjusted for potential message line

	inputLayoutWidth := m.width - 4 // General padding for input area
	if inputLayoutWidth < 30 {
		inputLayoutWidth = 30
	} // Min width for placeholders

	for i := range m.textInputs {
		m.textInputs[i].Width = inputLayoutWidth
	}
}

// Changed to pointer receiver for consistency
func (m *Model) IsFocused() bool {
	// The form is always the primary interaction if this model is active,
	// until a proof is generated. When proof is shown, it's more of a display state.
	return m.state == FormView
}
