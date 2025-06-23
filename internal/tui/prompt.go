package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type PromptModel struct {
	prompt      string
	textInput   textinput.Model
	err         error
	quitting    bool
	submitted   bool
	SubmittedCh chan string // Channel to send the submitted value
}

var (
	promptStyle     = lipgloss.NewStyle().Padding(0, 1)
	focusedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle         = lipgloss.NewStyle()
	helpStyle       = blurredStyle.Copy()
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red for errors
	cursorModeHelp  = helpStyle.Render("cursor mode is enabled")
	focusedButton   = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton   = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func NewPromptModel(promptText string) PromptModel {
	ti := textinput.New()
	ti.Placeholder = "Type your answer here..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50
	ti.Prompt = "┃ "

	return PromptModel{
		prompt:      promptText,
		textInput:   ti,
		err:         nil,
		SubmittedCh: make(chan string, 1),
	}
}

func (m PromptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.textInput.Value() != "" {
				m.submitted = true
				m.quitting = true
				m.SubmittedCh <- m.textInput.Value() // Send value before quitting
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			m.SubmittedCh <- "" // Send empty string if user quits
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m PromptModel) View() string {
	if m.quitting {
		if m.submitted {
			// Keep the prompt and submitted value visible briefly
			return fmt.Sprintf("%s\n%s%s\n", m.prompt, m.textInput.Prompt, m.textInput.Value())
		}
		return "" // Return empty when quitting without submission or after showing submitted value
	}

	var viewBuilder strings.Builder

	viewBuilder.WriteString(promptStyle.Render(m.prompt) + "\n")
	viewBuilder.WriteString(m.textInput.View() + "\n\n")
	viewBuilder.WriteString(blurredButton) // Placeholder for submit button

	if m.err != nil {
		viewBuilder.WriteString("\n" + errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	help := helpStyle.Render("Press Enter to submit, Esc or Ctrl+C to quit.")
	viewBuilder.WriteString("\n\n" + help)


	return viewBuilder.String()
}

// GetInput runs a prompt and returns the user's input or an error.
// It blocks until the user submits or quits.
func GetInput(promptText string, output io.Writer, input io.Reader) (string, error) {
	model := NewPromptModel(promptText)

	opts := []tea.ProgramOption{tea.WithOutput(output)}
	if input != nil {
		opts = append(opts, tea.WithInput(input))
	}

	p := tea.NewProgram(model, opts...)

	// Goroutine to run the program and wait for the result
	go func() {
		if _, err := p.Run(); err != nil {
			// Handle error, perhaps by sending it through a channel if needed
			// For now, we assume the main function will handle it or it's logged.
			fmt.Fprintf(output, "Error running prompt program: %v\n", err)
		}
	}()

	// Block until a value is received from the channel
	submittedValue := <-model.SubmittedCh
	close(model.SubmittedCh)


	if submittedValue == "" && !model.submitted { // User quit without submitting
		return "", fmt.Errorf("prompt aborted by user")
	}
	return submittedValue, nil
}

// Example Usage (can be moved to a main or test file)
/*
func main() {
    // Example: GetInput
    fmt.Println("Starting prompt...")
    userInput, err := GetInput("What is your name?", os.Stdout, os.Stdin)
    if err != nil {
        fmt.Printf("Prompt error: %v\n", err)
        return
    }
    if userInput != "" {
        fmt.Printf("\nHello, %s!\n", userInput)
    } else {
        fmt.Println("\nNo input received.")
    }
}
*/
