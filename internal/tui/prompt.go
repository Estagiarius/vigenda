package tui

import (
	"fmt"
	"io"
	"os" // Required for isatty.IsTerminal
	"bufio" // For reading from non-TTY stdin

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty" // For TTY detection
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
	// helpStyle is now defined in tui.go and accessible within the package.
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red for errors
	cursorModeHelp  = helpStyle.Render("cursor mode is enabled") // This will use helpStyle from tui.go
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
func GetInput(promptText string, output io.Writer, inputReader io.Reader) (string, error) {
	// Check if stdin is a TTY. inputReader is typically os.Stdin.
	// This check is most relevant when inputReader is os.Stdin.
	// If a different reader is passed, we might assume it's non-TTY or requires interactive handling.
	// For simplicity, we check specifically for os.Stdin.
	if f, ok := inputReader.(*os.File); ok && !isatty.IsTerminal(f.Fd()) && !isatty.IsCygwinTerminal(f.Fd()) {
		// Not a TTY, attempt to read directly.
		// This is a simplified way to handle piped input for a single-line prompt.
		// It assumes the entire piped content is the desired input for the prompt.
		// Print the prompt message to stderr so it doesn't interfere with piped output.
		fmt.Fprintln(os.Stderr, promptText)

		scanner := bufio.NewScanner(inputReader)
		if scanner.Scan() {
			return strings.TrimSpace(scanner.Text()), nil
		}
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading from non-TTY input: %w", err)
		}
		// If nothing was scanned (e.g., empty pipe), return empty string.
		// This behavior might need adjustment based on requirements (e.g., error on empty pipe).
		return "", nil
	}

	// It's a TTY, or not os.Stdin, run the interactive BubbleTea prompt.
	model := NewPromptModel(promptText)

	opts := []tea.ProgramOption{tea.WithOutput(output)}
	// Ensure the input for BubbleTea is correctly set.
	// If inputReader was os.Stdin and it's a TTY, BubbleTea will use it.
	// If inputReader was something else, BubbleTea will use that.
	if inputReader != nil {
		opts = append(opts, tea.WithInput(inputReader))
	}

	p := tea.NewProgram(model, opts...)

	// It's generally better to run p.Run() in the same goroutine
	// and handle its model directly, especially for simple blocking prompts.
	// The channel communication with a goroutine can be tricky for termination.

	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running prompt program: %w", err)
	}

	// Type assert to get the final model state
	finalPromptModel, ok := finalModel.(PromptModel)
	if !ok {
		return "", fmt.Errorf("internal error: prompt model has unexpected type")
	}

	if finalPromptModel.submitted {
		return finalPromptModel.textInput.Value(), nil
	}
	// If not submitted, it means the user quit (Esc, Ctrl+C)
	return "", fmt.Errorf("prompt aborted by user")
}

// Example Usage (can be moved to a main or test file)

func main_example() { // Renamed to avoid conflict if this file is part of a library build
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

/*
// To run this example:
// 1. Save this file (e.g., as prompt_example.go or ensure it's in a main package context)
// 2. Add `import "os"` if not already present at the top level for os.Stdout, os.Stdin
// 3. Call main_example() from a real main function.
// Example main.go for testing:
package main
import (
	"fmt"
	"os"
	"vigenda/internal/tui" // Adjust import path as necessary
)
func main() {
	fmt.Println("Testing prompt.GetInput with piped data:")
	// Simulate piped input scenario:
	// echo "TestName" | go run your_main_test_program.go

	// For direct execution test without pipe, os.Stdin will be the terminal.
	// To test piped: echo "MyPipedInput" | go run main_test.go
	input, err := tui.GetInput("Enter data:", os.Stdout, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from GetInput: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Received: '%s'\n", input)
}
*/
