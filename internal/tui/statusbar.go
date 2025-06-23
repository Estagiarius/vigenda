package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Height(1).
			Padding(0, 1).
			Background(lipgloss.Color("236")). // Dark gray background
			Foreground(lipgloss.Color("250"))  // Light gray text

	statusTextLeftStyle = lipgloss.NewStyle().Align(lipgloss.Left)
	statusTextRightStyle = lipgloss.NewStyle().Align(lipgloss.Right)

	separator = " | "
)

type StatusBarModel struct {
	width         int
	status        string // Main status message
	ephemeralMsg  string // Temporary message (e.g., "Copied!")
	ephemeralTime time.Time
	ephemeralTTL  time.Duration // Time to live for ephemeral messages
	rightText     string        // Text to display on the right (e.g., time, version)
}

func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{
		status:       "Ready",
		ephemeralTTL: 2 * time.Second, // Default TTL for ephemeral messages
		rightText:    time.Now().Format("15:04:05"),
	}
}

func (m StatusBarModel) Init() tea.Cmd {
	// Update time every second
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return UpdateTimeMsg(t)
	})
}

// UpdateTimeMsg is sent when the time should be updated in the status bar.
type UpdateTimeMsg time.Time

// SetStatusMsg is a message to update the main status.
type SetStatusMsg string

// SetEphemeralStatusMsg is a message to show a temporary status.
type SetEphemeralStatusMsg struct {
	Text string
	TTL  time.Duration // Optional: if zero, uses default TTL
}


func (m StatusBarModel) Update(msg tea.Msg) (StatusBarModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case UpdateTimeMsg:
		m.rightText = time.Time(msg).Format("15:04:05")
		// Re-tick for next update
		cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return UpdateTimeMsg(t)
		}))
	case SetStatusMsg:
		m.status = string(msg)
		m.clearEphemeral() // Clear any temporary message
	case SetEphemeralStatusMsg:
		m.ephemeralMsg = msg.Text
		m.ephemeralTime = time.Now()
		if msg.TTL > 0 {
			m.ephemeralTTL = msg.TTL
		}
		// Schedule a clear for the ephemeral message
		cmds = append(cmds, tea.Tick(m.ephemeralTTL, func(t time.Time) tea.Msg {
			return ClearEphemeralMsg{}
		}))

	case ClearEphemeralMsg:
		if time.Since(m.ephemeralTime) >= m.ephemeralTTL { // Ensure it's the right clear
			m.clearEphemeral()
		}
	}
	return m, tea.Batch(cmds...)
}

// ClearEphemeralMsg is a message to clear the ephemeral status.
type ClearEphemeralMsg struct{}

func (m *StatusBarModel) clearEphemeral() {
	m.ephemeralMsg = ""
}

func (m StatusBarModel) View() string {
	if m.width == 0 {
		return "" // Don't render if width is not set
	}

	displayStatus := m.status
	if m.ephemeralMsg != "" {
		if time.Since(m.ephemeralTime) < m.ephemeralTTL {
			displayStatus = m.ephemeralMsg
		} else {
			// This case should ideally be handled by ClearEphemeralMsg
			// but as a fallback, we ensure it's cleared if Update missed it.
			// No, we don't want to mutate in View. The Update should handle clearing.
			// If we reach here and it's not cleared, it means Update didn't get ClearEphemeralMsg yet.
		}
	}

	left := statusTextLeftStyle.Render(displayStatus)
	right := statusTextRightStyle.Render(m.rightText)

	// Calculate available width for the left status, considering the right text and separator
	// This is a simplified approach. For more complex layouts, consider advanced string manipulation.
	availableWidth := m.width - lipgloss.Width(right) - lipgloss.Width(separator)
	if availableWidth < 0 {
		availableWidth = 0
	}

	// Truncate left text if too long
	if lipgloss.Width(left) > availableWidth {
		// This is a naive truncation, doesn't handle multi-byte characters well.
		// For robust truncation, consider using a library or more sophisticated logic.
		// For now, we'll just take a substring.
		// A better approach for rune-correct truncation:
		runes := []rune(displayStatus)
		if len(runes) > availableWidth {
			left = statusTextLeftStyle.Render(string(runes[:availableWidth-3]) + "...")
		} else {
             left = statusTextLeftStyle.Render(string(runes))
        }
	}


	// Ensure the total content does not exceed the bar width
    // This is tricky with lipgloss alignment. We'll render the full bar and let lipgloss handle it.
    // The statusBarStyle will clip if content is too wide.

    // We need to fill the space between left and right.
    // Calculate the width of the gap
    gapWidth := m.width - lipgloss.Width(left) - lipgloss.Width(right)
    if gapWidth < 0 {
        gapWidth = 0 // Should not happen if truncation is correct
    }
    gap := strings.Repeat(" ", gapWidth)


	return statusBarStyle.Width(m.width).Render(left + gap + right)
}

// Example of integrating StatusBarModel into a larger Bubble Tea application:
/*
type MainModel struct {
    statusBar StatusBarModel
    // other parts of your model
}

func NewMainModel() MainModel {
    return MainModel{
        statusBar: NewStatusBarModel(),
        // initialize other parts
    }
}

func (m MainModel) Init() tea.Cmd {
    return m.statusBar.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd

    // Process messages for the main model components
    // ...

    // Update status bar
    newStatusBar, statusBarCmd := m.statusBar.Update(msg)
    m.statusBar = newStatusBar
    if statusBarCmd != nil {
        cmds = append(cmds, statusBarCmd)
    }

    // Example: Sending a status update to the status bar
    // if someCondition {
    //     cmds = append(cmds, func() tea.Msg { return SetStatusMsg("New status achieved!") })
    // }
    // if someOtherCondition {
    //     cmds = append(cmds, func() tea.Msg { return SetEphemeralStatusMsg{Text: "Action Done!", TTL: 3 * time.Second} })
    // }


    return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
    // Render other parts of your UI
    mainContentView := "Your main content here..."

    // Combine with status bar view
    return lipgloss.JoinVertical(lipgloss.Left,
        mainContentView,
        m.statusBar.View(),
    )
}

func main() {
    model := NewMainModel()
    p := tea.NewProgram(model, tea.WithAltScreen()) // Example with AltScreen
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error running program: %v", err)
    }
}
*/
