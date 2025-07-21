package settings

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	apiKeyInput   textinput.Model
	baseURLInput  textinput.Model
	focusedInput  int
	width, height int
}

func New() *Model {
	apiKey := textinput.New()
	apiKey.Placeholder = "Sua chave de API da OpenAI"
	apiKey.Focus()

	baseURL := textinput.New()
	baseURL.Placeholder = "URL base da API (opcional)"

	return &Model{
		apiKeyInput:  apiKey,
		baseURLInput: baseURL,
		focusedInput: 0,
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" {
				// Save settings or move to next view
				return m, tea.Quit // for now, just quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusedInput--
			} else {
				m.focusedInput++
			}

			if m.focusedInput > 1 {
				m.focusedInput = 0
			}
			if m.focusedInput < 0 {
				m.focusedInput = 1
			}

			cmds := make([]tea.Cmd, 2)
			if m.focusedInput == 0 {
				cmds[0] = m.apiKeyInput.Focus()
				m.baseURLInput.Blur()
			} else {
				cmds[0] = m.baseURLInput.Focus()
				m.apiKeyInput.Blur()
			}
			return m, tea.Batch(cmds...)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.focusedInput == 0 {
		m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
	} else {
		m.baseURLInput, cmd = m.baseURLInput.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Configurações da API da OpenAI",
		m.apiKeyInput.View(),
		m.baseURLInput.View(),
		"(tab para mudar de campo, esc para sair)",
	)
}

func (m *Model) ApiKey() string {
	return m.apiKeyInput.Value()
}

func (m *Model) BaseURL() string {
	return m.baseURLInput.Value()
}

func (m *Model) CanGoBack() bool {
	return true
}
