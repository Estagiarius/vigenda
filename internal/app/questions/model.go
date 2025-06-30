package questions

import (
	"context"
	"fmt"
	"os" // For os.ReadFile (if service needs bytes)
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// "vigenda/internal/models" // Might need if listing questions
	"vigenda/internal/service"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// ViewState defines the current state of the questions view
type ViewState int

const (
	ActionListView ViewState = iota // Main action list (e.g., "Add Questions from JSON")
	AddQuestionsFormView
	// ListQuestionsView // Optional: For listing questions, potentially with filters
)

// Model represents the question bank management model.
type Model struct {
	questionService service.QuestionService
	state           ViewState
	list            list.Model
	textInputs      []textinput.Model
	focusIndex      int
	isLoading       bool
	err             error
	message         string
	width           int
	height          int
}

// --- Messages ---
type questionsAddedMsg struct {
	count int
	err   error
}

// --- Cmds ---
// (No initial data loading command like loadQuestionsCmd unless we implement listing)

func New(questionService service.QuestionService) *Model { // Return *Model
	actionItems := []list.Item{
		actionItem{title: "Adicionar Questões de JSON", description: "Importar questões de um arquivo JSON."},
		// actionItem{title: "Listar Questões", description: "Visualizar e filtrar questões do banco."},
	}
	l := list.New(actionItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Banco de Questões"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	// Se você quiser configurar key bindings personalizados para a lista:
	// keyBindings := []key.Binding{
	//     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "selecionar")),
	//     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "voltar")),
	// }
	// Você precisaria aplicar esses bindings conforme a API da biblioteca

	// Text input for JSON file path
	inputs := make([]textinput.Model, 1)
	ti := textinput.New()
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ti.CharLimit = 255 // Path length
	ti.Placeholder = "Caminho para o arquivo JSON de questões"
	inputs[0] = ti

	return &Model{ // Corrected to return a pointer
		questionService: questionService,
		state:           ActionListView,
		list:            l,
		textInputs:      inputs,
		isLoading:       false,
	}
}

// Changed to pointer receiver
func (m *Model) Init() tea.Cmd {
	m.state = ActionListView
	m.err = nil
	m.message = ""
	m.resetForms()
	m.list.Select(-1) // Deselect any previous action
	return nil
}

type actionItem struct { // Local definition
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
			if m.state == ActionListView {
				return m, nil // Let parent model handle 'esc' from main action list
			}
			// Go back to action list view from form
			m.state = ActionListView
			m.err = nil
			m.message = ""
			m.resetForms()
			m.list.Select(-1)
			return m, nil
		}

		switch m.state {
		case ActionListView:
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
					case "Adicionar Questões de JSON":
						m.state = AddQuestionsFormView
						m.setupAddQuestionsForm()
						cmds = append(cmds, m.textInputs[0].Focus())
					// case "Listar Questões":
					// m.state = ListQuestionsView
					// m.isLoading = true
					// cmds = append(cmds, m.loadQuestionsCmd()) // If implemented
					}
				}
			}

		case AddQuestionsFormView:
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				m.isLoading = true
				cmds = append(cmds, m.submitAddQuestionsFormCmd())
			} else {
				// m.textInputs[0], cmd = m.textInputs[0].Update(msg)
				var updatedInput textinput.Model
				updatedInput, cmd = m.textInputs[0].Update(msg)
				m.textInputs[0] = updatedInput
				cmds = append(cmds, cmd)
			}
		}

	// Handle async results
	case questionsAddedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.message = fmt.Sprintf("%d questões adicionadas com sucesso!", msg.count)
			m.state = ActionListView // Go back to action list
			m.list.Select(-1)
			m.resetForms()
		}

	case error:
		m.err = msg
		m.isLoading = false

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height) // Use the SetSize method
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

	switch m.state {
	case ActionListView:
		b.WriteString(m.list.View())

	case AddQuestionsFormView:
		b.WriteString("Adicionar Questões de Arquivo JSON\n\n")
		b.WriteString(m.textInputs[0].View() + "\n\n")
		b.WriteString("(Pressione Enter para submeter, Esc para cancelar)\n")
		b.WriteString("Nota: Certifique-se que o caminho do arquivo é acessível pela aplicação.")


	// case ListQuestionsView:
	//  b.WriteString("Listagem de Questões (Não Implementado)\n")
	//  b.WriteString(m.table.View()) // If using a table
	//  b.WriteString("\n(Navegue com ↑/↓, 'esc' para voltar às ações)")

	default:
		b.WriteString("Visualização do Banco de Questões Desconhecida")
	}

	return baseStyle.Render(b.String())
}

// Changed to pointer receiver
func (m *Model) resetForms() {
	if len(m.textInputs) > 0 {
		m.textInputs[0].Reset()
		m.textInputs[0].Blur()
	}
	m.focusIndex = 0 // Although not used for navigation in this simple form
	m.err = nil
	m.message = ""
}

// Changed to pointer receiver
func (m *Model) setupAddQuestionsForm() {
	m.focusIndex = 0 // Only one input, so focus is implicitly on it or a submit action
	if len(m.textInputs) > 0 {
		m.textInputs[0].Reset()
		m.textInputs[0].Placeholder = "Caminho para o arquivo JSON de questões"
		// m.textInputs[0].Focus() // This will be called in Update when switching state
	}
}

// Changed to pointer receiver
func (m *Model) submitAddQuestionsFormCmd() tea.Cmd {
	jsonPath := m.textInputs[0].Value()
	if jsonPath == "" {
		m.err = fmt.Errorf("O caminho do arquivo JSON é obrigatório.")
		m.isLoading = false
		return nil
	}

	return func() tea.Msg {
		// The service QuestionService.AddQuestionsFromJSON expects jsonData []byte.
		// We need to read the file here.
		jsonData, err := os.ReadFile(jsonPath)
		if err != nil {
			return questionsAddedMsg{err: fmt.Errorf("falha ao ler arquivo JSON '%s': %w", jsonPath, err)}
		}
		count, err := m.questionService.AddQuestionsFromJSON(context.Background(), jsonData)
		return questionsAddedMsg{count: count, err: err}
	}
}

func (m *Model) SetSize(width, height int) {
	m.width = width - baseStyle.GetHorizontalFrameSize()
	m.height = height - baseStyle.GetVerticalFrameSize() - 1

	listHeight := m.height - lipgloss.Height(m.list.Title) - 2
	m.list.SetSize(m.width, listHeight)

	// Table size (if ListQuestionsView is implemented)
	// m.table.SetWidth(m.width)
	// tableHeight := m.height - 6
	// if tableHeight < 5 { tableHeight = 5 }
	// m.table.SetHeight(tableHeight)

	inputWidth := m.width - 4 // Ensure this aligns with how View renders padding
	if inputWidth < 20 {
		inputWidth = 20
	}
	if len(m.textInputs) > 0 {
		m.textInputs[0].Width = inputWidth
	}
}

// Changed to pointer receiver for consistency
func (m *Model) IsFocused() bool {
	// Focused if in the form input state
	return m.state == AddQuestionsFormView
}
