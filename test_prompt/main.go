package main

import (
	"fmt"
	"os"
	"vigenda/internal/tui" // Assuming vigenda is the module name
	"github.com/mattn/go-isatty"
)

func main() {
	fmt.Fprintln(os.Stderr, "Testing prompt.GetInput with potentially piped data...")
	isTerminal := isatty.IsTerminal(os.Stdin.Fd())
	isCygwinTerminal := isatty.IsCygwinTerminal(os.Stdin.Fd())
	fmt.Fprintf(os.Stderr, "os.Stdin is TTY: %v, CygwinTTY: %v\n", isTerminal, isCygwinTerminal)

	// To test piped: echo "MyPipedInput" | go run test_prompt/main.go
	// To test interactive: go run test_prompt/main.go (and type then enter)
	input, err := tui.GetInput("Enter data:", os.Stdout, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from GetInput: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Received_from_prompt:'%s'\n", input) // Use a unique marker
}
