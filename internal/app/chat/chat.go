package chat

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	openai "github.com/sashabaranov/go-openai"
	"vigenda/internal/app/settings"
)

var (
	senderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	botStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
)

type Model struct {
	viewport    viewport.Model
	textarea    textarea.Model
	messages    []openai.ChatCompletionMessage
	client      *openai.Client
	settings    *settings.Model
	err         error
	width       int
	height      int
	isStreaming bool
}

func New(settings *settings.Model) *Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = -1
	ta.SetWidth(50)
	ta.SetHeight(1)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(50, 10)

	return &Model{
		textarea: ta,
		viewport: vp,
		messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Você é um agente de viagem. Seja descritivo e gentil.",
			},
		},
		settings: settings,
	}
}

func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textarea.Value() != "" && !m.isStreaming {
				m.messages = append(m.messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: m.textarea.Value(),
				})
				m.renderMessages()
				m.isStreaming = true
				return m, m.streamResponse()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(m.View())
		m.textarea.SetWidth(msg.Width)
		m.renderMessages()

	case streamMsg:
		m.messages[len(m.messages)-1].Content += msg.content
		m.renderMessages()
		return m, m.streamResponse()

	case streamEndMsg:
		m.isStreaming = false
		if msg.err != nil {
			m.err = msg.err
		}
		return m, nil

	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(taCmd, vpCmd)
}

func (m *Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		m.textarea.View(),
	)
}

type streamMsg struct{ content string }
type streamEndMsg struct{ err error }

func (m *Model) streamResponse() tea.Cmd {
	return func() tea.Msg {
		config := openai.DefaultConfig(m.settings.ApiKey())
		if m.settings.BaseURL() != "" {
			config.BaseURL = m.settings.BaseURL()
		}
		m.client = openai.NewClientWithConfig(config)

		ctx := context.Background()
		req := openai.ChatCompletionRequest{
			Model:    "sabia-3",
			Messages: m.messages,
			Stream:   true,
		}
		stream, err := m.client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			return streamEndMsg{err: err}
		}

		// Add an empty bot message to append to
		if len(m.messages) == 0 || m.messages[len(m.messages)-1].Role != openai.ChatMessageRoleAssistant {
			m.messages = append(m.messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "",
			})
		}

		resp, err := stream.Recv()
		if err != nil {
			stream.Close()
			return streamEndMsg{err: err}
		}

		if len(resp.Choices) > 0 {
			return streamMsg{content: resp.Choices[0].Delta.Content}
		}

		stream.Close()
		return streamEndMsg{err: nil}
	}
}

func (m *Model) renderMessages() {
	var sb strings.Builder
	for _, msg := range m.messages {
		if msg.Role == openai.ChatMessageRoleSystem {
			continue
		}
		var style lipgloss.Style
		var role string
		if msg.Role == openai.ChatMessageRoleUser {
			style = senderStyle
			role = "You"
		} else {
			style = botStyle
			role = "Bot"
		}
		sb.WriteString(style.Render(role) + ":\n")
		sb.WriteString(msg.Content + "\n\n")
	}
	m.viewport.SetContent(sb.String())
	m.viewport.GotoBottom()
}

func (m *Model) CanGoBack() bool {
	return true
}
