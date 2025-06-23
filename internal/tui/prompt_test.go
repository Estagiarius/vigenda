package tui

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewPromptModel(t *testing.T) {
	promptText := "Enter your name:"
	model := NewPromptModel(promptText)

	if model.prompt != promptText {
		t.Errorf("Expected prompt text '%s', got '%s'", promptText, model.prompt)
	}
	if model.textInput.Placeholder != "Type your answer here..." {
		t.Errorf("Unexpected placeholder: '%s'", model.textInput.Placeholder)
	}
	if !model.textInput.Focused() {
		t.Errorf("TextInput should be focused by default")
	}
	if model.SubmittedCh == nil {
		t.Errorf("SubmittedCh should be initialized")
	}
}

func TestPromptModel_Update_Enter(t *testing.T) {
	model := NewPromptModel("Test prompt")
	testValue := "hello world"
	model.textInput.SetValue(testValue)

	// Test Enter key
	updatedModelInterface, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updatedModel := updatedModelInterface.(PromptModel)

	if !updatedModel.submitted {
		t.Errorf("Model should be in submitted state after Enter")
	}
	if !updatedModel.quitting {
		t.Errorf("Model should be in quitting state after Enter")
	}
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("Expected tea.Quit command on Enter key press with non-empty input")
	}

	// Check if value was sent to channel
	select {
	case val := <-updatedModel.SubmittedCh:
		if val != testValue {
			t.Errorf("Expected value '%s' from SubmittedCh, got '%s'", testValue, val)
		}
	case <-time.After(100 * time.Millisecond): // Timeout to prevent test hanging
		t.Errorf("Timeout waiting for value from SubmittedCh")
	}
}

func TestPromptModel_Update_Enter_Empty(t *testing.T) {
	model := NewPromptModel("Test prompt")
	model.textInput.SetValue("") // Empty input

	// Test Enter key with empty input
	updatedModelInterface, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updatedModel := updatedModelInterface.(PromptModel)

	if updatedModel.submitted {
		t.Errorf("Model should not be in submitted state after Enter with empty input")
	}
	if updatedModel.quitting { // It should not quit if input is empty
		t.Errorf("Model should not be in quitting state after Enter with empty input")
	}
	if cmd != nil { // No tea.Quit command expected
		t.Errorf("Expected no command on Enter key press with empty input, got %v", cmd)
	}
}


func TestPromptModel_Update_Quit(t *testing.T) {
	testCases := []struct {
		name    string
		keyType tea.KeyType
		keyRunes []rune
	}{
		{"CtrlC", tea.KeyCtrlC, nil},
		{"Esc", tea.KeyEsc, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := NewPromptModel("Test prompt")
			updatedModelInterface, cmd := model.Update(tea.KeyMsg{Type: tc.keyType, Runes: tc.keyRunes})
			updatedModel := updatedModelInterface.(PromptModel)

			if !updatedModel.quitting {
				t.Errorf("Model should be in quitting state after %s", tc.name)
			}
			if cmd == nil || cmd() != tea.Quit() {
				t.Errorf("Expected tea.Quit command on %s key press", tc.name)
			}
			// Check if empty value was sent to channel
			select {
			case val := <-updatedModel.SubmittedCh:
				if val != "" {
					t.Errorf("Expected empty value from SubmittedCh on quit, got '%s'", val)
				}
			case <-time.After(100 * time.Millisecond):
				t.Errorf("Timeout waiting for value from SubmittedCh on quit")
			}
		})
	}
}

func TestPromptModel_Update_TextInput(t *testing.T) {
	model := NewPromptModel("Test prompt")
	// Simulate typing 'a'
	updatedModelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	updatedModel := updatedModelInterface.(PromptModel)

	if updatedModel.textInput.Value() != "a" {
		t.Errorf("Expected textInput value 'a', got '%s'", updatedModel.textInput.Value())
	}
}

func TestPromptModel_View(t *testing.T) {
	promptText := "Your favorite color?"
	model := NewPromptModel(promptText)
	model.textInput.SetValue("Blue")

	view := model.View()

	if !strings.Contains(view, promptText) {
		t.Errorf("View should contain prompt text '%s'. Got: %s", promptText, view)
	}
	if !strings.Contains(view, "Blue") { // Check if current input value is in view
		t.Errorf("View should contain current textInput value 'Blue'. Got: %s", view)
	}
	if !strings.Contains(view, "[ Submit ]") { // Check for submit button text
		t.Errorf("View should contain submit button text. Got: %s", view)
	}

	// Test view when quitting
	model.quitting = true
	model.submitted = false // Quit without submitting
	viewAfterQuit := model.View()
	if viewAfterQuit != "" {
		t.Errorf("View should be empty when quitting without submission. Got: %s", viewAfterQuit)
	}

	model.submitted = true // Quit after submitting
	model.textInput.SetValue("Final Answer")
	viewAfterSubmitQuit := model.View()
	expectedView := fmt.Sprintf("%s\n%s%s\n", model.prompt, model.textInput.Prompt, "Final Answer")
	if viewAfterSubmitQuit != expectedView {
		t.Errorf("View after submit and quit is not as expected.\nExpected:\n%s\nGot:\n%s", expectedView, viewAfterSubmitQuit)
	}
}


func TestGetInput_Simulated(t *testing.T) {
	promptText := "Enter test input:"

	// Simulate user input by preparing a reader
	// This is a simplified simulation. Real Bubble Tea program runs its own loop.
	// For GetInput, the tea.Program runs in a goroutine.
	// We need to send a value through the channel as if the program ran and user submitted.

	t.Run("User submits input", func(t *testing.T) {
		// For this test, we can't easily use a real input stream with tea.Program.
		// Instead, we'll test the logic by directly manipulating a model instance
		// and checking the channel behavior, which GetInput relies on.

		model := NewPromptModel(promptText)
		expectedInput := "simulated input"

		go func() {
			// Simulate the program running and user typing then pressing enter
			// This goroutine will act like the tea.Program
			time.Sleep(10 * time.Millisecond) // Give GetInput time to start listening
			model.textInput.SetValue(expectedInput)
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			if cmd != nil && cmd() == tea.Quit() {
				// The model's Update method sends to SubmittedCh upon successful submission
				// model.SubmittedCh <- model.textInput.Value() // This is done by Update
			} else {
				// Should not happen in this simulation if Update works correctly
				// model.SubmittedCh <- "error in simulation"
			}
		}()

		// GetInput will create its own model, so we need a different approach.
		// We can't directly inject keypresses into the tea.Program started by GetInput
		// without a more complex test setup (e.g., using bubbletea.WithInput(fakeInput)).

		// A more practical test for GetInput would involve creating a Program
		// with a custom input that sends key events, or testing its sub-components.
		// Given the structure of GetInput, testing its behavior with a real tea.Program
		// without actual terminal interaction is tricky.

		// Let's try a version of GetInput that we can control more easily for testing,
		// or accept that testing the full GetInput flow is an integration test.

		// For now, we'll test the happy path using a short timeout and assuming the goroutine
		// in GetInput works. This is more of an integration-style test for GetInput.
		var outputBuf bytes.Buffer

		// To make this testable without hanging, we need to ensure the tea.Program inside GetInput quits.
		// One way is to send a quit message or simulate input.
		// For this specific test, we will rely on the timeout of the channel read within GetInput
		// if the program doesn't quit as expected, or the channel logic within GetInput.

		// Test case: User provides input
		// This requires a way to feed "test\n" into the program's input.
		// Bubble Tea's default input is os.Stdin. We can use tea.WithInput for testing.

		inputReader, inputWriter := io.Pipe()

		go func() {
			defer inputWriter.Close()
			time.Sleep(50 * time.Millisecond) // give program time to start
			fmt.Fprintln(inputWriter, "test input") // Simulate typing "test input" and pressing enter
		}()

		userInput, err := GetInput(promptText, &outputBuf, inputReader)
		if err != nil {
			t.Fatalf("GetInput returned an error: %v. Output: %s", err, outputBuf.String())
		}
		if userInput != "test input" {
			t.Errorf("Expected input 'test input', got '%s'. Output: %s", userInput, outputBuf.String())
		}
	})


	t.Run("User aborts prompt", func(t *testing.T) {
		var outputBuf bytes.Buffer
		// The following inputReader, inputWriter, and the first go func were problematic and unused
		// in the actual logic path being tested by escInputReader.
		// inputReader, inputWriter := io.Pipe()
		// go func() {
		//	defer inputWriter.Close()
		//	time.Sleep(50 * time.Millisecond)
		// ... (rest of the original go func content commented out) ...
		// }()

		// Simulating an abort is tricky. If we just close the input pipe,
		// the textinput might not register it as an abort signal.
		// The GetInput function relies on PromptModel sending "" to SubmittedCh on abort.
		// If we want to test GetInput's error path for abort:
		// We'd need to ensure its internal PromptModel instance's SubmittedCh receives "".

		// For this test, let's assume the internal program quits immediately (e.g., error or immediate Esc).
		// This is still not a perfect simulation of user pressing Esc.
		// A more robust test would be to use bubbletea's Program.Send() if we had access to it.

		// Due to the complexities of simulating interactive terminal behavior,
		// this part of the test for GetInput might be less reliable or require
		// more advanced Bubble Tea testing techniques.
		// The core logic of submission and quitting is tested in PromptModel_Update.
		// GetInput is a wrapper that runs the program.

		// Let's simplify: if GetInput is called and the user (somehow) quits immediately
		// such that "" is sent to the channel *and* model.submitted is false,
		// it should return an error.
		// This is what happens if tea.KeyEsc is the first thing processed.

		// Re-simulate with Esc key by providing specific input
		escInputReader, escInputWriter := io.Pipe()
		go func() {
			defer escInputWriter.Close()
			time.Sleep(50 * time.Millisecond)
			// Sending escape character. Note: Bubbletea might need a proper TTY
			// for this to be interpreted as KeyEsc.
			// For non-TTY, KeyEsc might not be generated this way.
			// This part of the test is environment-dependent.
			escInputWriter.Write([]byte{0x1b}) // ANSI escape character
		}()

		// Test with a timeout because if Esc is not handled, it might hang.
		done := make(chan struct{})
		var err error
		var userInput string

		go func() {
			userInput, err = GetInput("Abort test", &outputBuf, escInputReader)
			close(done)
		}()

		select {
		case <-done:
			if err == nil {
				t.Errorf("Expected an error when prompt is aborted, got nil. UserInput: '%s'", userInput)
			} else if !strings.Contains(err.Error(), "prompt aborted by user") {
				t.Errorf("Expected error message 'prompt aborted by user', got '%v'", err)
			}
		case <-time.After(5 * time.Second): // Increased timeout further to 5 seconds
			t.Errorf("GetInput timed out during abort test. This might indicate Esc was not processed.")
		}
	})
}
