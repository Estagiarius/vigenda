package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss" // Added missing import
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModel_InitialState(t *testing.T) {
	m := New()
	assert.Equal(t, DashboardView, m.currentView, "Initial view should be DashboardView")
	require.Greater(t, m.list.Index(), -1, "List should have items") // list.Items() is not public
	assert.Contains(t, m.list.Items()[0].(menuItem).Title(), DashboardView.String(), "First item should be Dashboard")
}

func TestModel_Update_Quit(t *testing.T) {
	m := New()
	// Test 'q' from DashboardView
	qMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	nextModel, cmd := m.Update(qMsg)
	assert.True(t, nextModel.(Model).quitting, "Model should be quitting on 'q'")
	assert.NotNil(t, cmd, "A command (tea.Quit) should be returned on 'q'") // Check cmd is not nil

	// Test 'ctrl+c'
	m = New() // Reset model
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	nextModel, cmd = m.Update(ctrlCMsg)
	assert.True(t, nextModel.(Model).quitting, "Model should be quitting on 'ctrl+c'")
	assert.NotNil(t, cmd, "A command (tea.Quit) should be returned on 'ctrl+c'") // Check cmd is not nil
}

func TestModel_Update_NavigateToSubViewAndBack(t *testing.T) {
	m := New()
	initialView := m.currentView

	// Simulate selecting the second item (TaskManagementView)
	m.list.Select(1) // Select "Gerenciar Tarefas"
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	nextModel, _ := m.Update(enterMsg)
	m = nextModel.(Model)

	assert.Equal(t, TaskManagementView, m.currentView, "View should change to TaskManagementView on Enter")

	// Simulate pressing 'esc' to go back
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	nextModel, _ = m.Update(escMsg)
	m = nextModel.(Model)

	assert.Equal(t, initialView, m.currentView, "View should change back to DashboardView on Esc")

	// Simulate pressing 'q' from subview to go back
	m.currentView = TaskManagementView // Go to subview again
	qMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	nextModel, _ = m.Update(qMsg)
	m = nextModel.(Model)
	assert.Equal(t, initialView, m.currentView, "View should change back to DashboardView on 'q' from subview")

}

func TestModel_View_Content(t *testing.T) {
	m := New()
	// Set initial size for consistent rendering of the list
	updatedModel, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updatedModel.(Model)

	// Initial view (Dashboard/Menu)
	viewOutput := m.View()
	assert.Contains(t, viewOutput, m.list.Title, "View should contain list title in DashboardView")
	assert.Contains(t, viewOutput, "Navegue com ↑/↓", "View should contain help text for menu")


	// Navigate to a sub-view
	m.list.Select(1) // "Gerenciar Tarefas"
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	nextModel, _ := m.Update(enterMsg)
	m = nextModel.(Model)

	viewOutput = m.View()
	expectedSubstring := "Você está na visão: " + TaskManagementView.String()
	assert.Contains(t, viewOutput, expectedSubstring, "View should show placeholder for TaskManagementView")
	assert.Contains(t, viewOutput, "Pressione 'esc' ou 'q' para voltar", "View should contain help text for subview")

}

func TestModel_Update_WindowSize(t *testing.T) {
	m := New()
	newWidth, newHeight := 80, 24
	sizeMsg := tea.WindowSizeMsg{Width: newWidth, Height: newHeight}
	nextModel, _ := m.Update(sizeMsg)
	m = nextModel.(Model)

	assert.Equal(t, newWidth, m.width, "Model width should be updated")
	assert.Equal(t, newHeight, m.height, "Model height should be updated")
	// Check if list size was also updated (approximate check)
	// This depends on internal calculations of list.SetSize
	// We know appStyle padding is 1,2. Vertical padding 2. Title height 1. Help text height 1. Total 4.
	// list.SetSize(msg.Width-appStyle.GetHorizontalPadding(), msg.Height-appStyle.GetVerticalPadding()-lipgloss.Height(m.list.Title)-2)
	// Horizontal padding is 2*2 = 4. So list width = 80 - 4 = 76
	// Vertical padding is 1*2 = 2. Title height = 1 (title + margin). Help text = 1. Total elements to subtract = 2 + 1 + 1 = 4
	// So list height = 24 - (appStyle vertical padding: 2) - (title height with margin: 2) - (help text: 1) = 24 - 2 - 2 -1 = 19
	// Actually, it's msg.Height-appStyle.GetVerticalPadding()-lipgloss.Height(m.list.Title)-2
	// msg.Height - 2 - 2 - 2 = 18. Let's recheck lipgloss.Height(m.list.Title)
	// m.list.Title is "Vigenda - Menu Principal". l.Styles.Title has MarginBottom(1). So height is 2.
	// So, msg.Height (24) - appStyle.GetVerticalPadding() (2) - lipgloss.Height(m.list.Title) (2) - 2 (for help text and general spacing) = 18
	assert.Equal(t, newWidth-appStyle.GetHorizontalPadding(), m.list.Width(), "List width should be updated based on window size")
	assert.Equal(t, newHeight-appStyle.GetVerticalPadding()-lipgloss.Height(m.list.Title)-2, m.list.Height(), "List height should be updated based on window size")
}

func TestMenuItem_Interface(t *testing.T) {
	item := menuItem{title: "Test Title", view: DashboardView}
	assert.Equal(t, "Test Title", item.Title())
	assert.Equal(t, "Test Title", item.FilterValue()) // FilterValue currently returns title
	assert.Equal(t, "", item.Description()) // Description is empty
}

func TestView_String(t *testing.T) {
	assert.Equal(t, "Dashboard", DashboardView.String())
	assert.Equal(t, "Gerenciar Tarefas", TaskManagementView.String())
	// Add more if new views are added and string representations are critical
}

// Helper to simulate a key press
func simulateKeyPress(m tea.Model, key rune) (tea.Model, tea.Cmd) {
	return m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}})
}

// Helper to simulate an Enter key press
func simulateEnterPress(m tea.Model) (tea.Model, tea.Cmd) {
	return m.Update(tea.KeyMsg{Type: tea.KeyEnter})
}

// Helper to simulate an Esc key press
func simulateEscPress(m tea.Model) (tea.Model, tea.Cmd) {
	return m.Update(tea.KeyMsg{Type: tea.KeyEsc})
}

// Helper for Ctrl+C
func simulateCtrlCPress(m tea.Model) (tea.Model, tea.Cmd) {
	return m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
}

// Example of using helpers (can be expanded into actual tests)
func TestModel_Update_WithHelpers(t *testing.T) {
	m := New()

	// Select second item ("Gerenciar Tarefas")
	m.list.Select(1)
	nextModel, _ := simulateEnterPress(m)
	m = nextModel.(Model)
	assert.Equal(t, TaskManagementView, m.currentView)

	nextModel, _ = simulateEscPress(m)
	m = nextModel.(Model)
	assert.Equal(t, DashboardView, m.currentView)

	nextModel, cmd := simulateCtrlCPress(m)
	assert.True(t, nextModel.(Model).quitting)
	assert.NotNil(t, cmd) // Check that tea.Quit (a non-nil cmd) is returned
}
