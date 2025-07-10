package styles

import "github.com/charmbracelet/lipgloss"

// Shared TUI Colors
var Colors = struct {
	Red    lipgloss.TerminalColor
	Green  lipgloss.TerminalColor
	Subtle lipgloss.TerminalColor
	// Add other common colors if needed
	ErrorColor   lipgloss.Color // Example for specific semantic color
	SuccessColor lipgloss.Color
	HelpColor    lipgloss.Color
}{
	Red:    lipgloss.Color("9"),
	Green:  lipgloss.Color("10"),
	Subtle: lipgloss.Color("241"),

	ErrorColor:   lipgloss.Color("9"),   // Consistent Red
	SuccessColor: lipgloss.Color("10"),  // Consistent Green
	HelpColor:    lipgloss.Color("240"), // Slightly different gray for help
}
