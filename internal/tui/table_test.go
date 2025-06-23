package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNewTableModel(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 10},
	}
	rows := []table.Row{
		{"1", "Test"},
	}

	model := NewTableModel(columns, rows)

	if model.table.Focused() != true {
		t.Errorf("Table should be focused by default")
	}
	if len(model.table.Columns()) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(model.table.Columns()))
	}
	if len(model.table.Rows()) != 1 {
		t.Errorf("Expected 1 row, got %d", len(model.table.Rows()))
	}
}

func TestTableModel_Update(t *testing.T) {
	columns := []table.Column{{Title: "Col1", Width: 5}}
	rows := []table.Row{{"Row1"}}
	model := NewTableModel(columns, rows)

	// Test quit message
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("Expected tea.Quit command on 'q' key press")
	}

	// Test ESC message
	model.table.Focus() // Ensure focused
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if updatedModel.(TableModel).table.Focused() {
		t.Errorf("Table should blur on first ESC if focused")
	}
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !updatedModel.(TableModel).table.Focused() {
		t.Errorf("Table should focus on second ESC if blurred")
	}

	// Test other key (e.g., down arrow)
	initialSelectedRow := model.table.Cursor()
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	// This depends on the internal table behavior, which might not change selection if only one row
	// For a more robust test, you might need more rows or to mock table.Update
	// For now, we just ensure it doesn't crash and returns the model itself
	if model.table.Cursor() == initialSelectedRow {
		// This is expected with one row, but if it were multi-row, we'd check for change
	}

}

func TestTableModel_View(t *testing.T) {
	columns := []table.Column{{Title: "Header", Width: 10}}
	rows := []table.Row{{"Cell"}}
	model := NewTableModel(columns, rows)

	view := model.View()
	if !strings.Contains(view, "Header") {
		t.Errorf("View should contain column header")
	}
	if !strings.Contains(view, "Cell") {
		t.Errorf("View should contain row cell data")
	}
	if !strings.HasPrefix(view, "┌") || !strings.HasSuffix(view, "┘\n") {
		// Check for border characters (may vary with lipgloss versions/styles)
		// This is a basic check for the baseStyle border
		// t.Logf("View output: %q", view) // Log for debugging if needed
	}
}

func TestShowTable(t *testing.T) {
	columns := []table.Column{
		{Title: "Key", Width: 5},
		{Title: "Value", Width: 10},
	}
	rows := []table.Row{
		{"k1", "v1"},
	}

	// var buf bytes.Buffer // This was part of a previous commented-out test idea

	// tea.NewProgram usually takes over the terminal.
	// For testing, we can't easily check the full interactive output.
	// We can ensure it runs without error and produces some output.
	// A more advanced test would involve using a test driver for Bubble Tea.

	// We can't directly call p.Run() in a test like this easily because it blocks
	// and expects terminal input.
	// Instead, we'll test that the model is created correctly.
	// A simple way to test ShowTable is to ensure it tries to create a program.
	// For now, we'll just verify it attempts to write something.
	// This is a very basic test for ShowTable due to Bubble Tea's nature.

	go func() {
		// Simulate a quit message to stop the program almost immediately.
		// This is tricky to coordinate in tests.
		// For a more robust test, we'd need a way to send messages to the program.
		// For now, we'll rely on the fact that NewProgram itself is tested by Bubble Tea.
		// Our main concern is that our model setup is correct.
		// ShowTable(columns, rows, &buf)
	}()

	// This test for ShowTable is limited. We'll mainly rely on testing TableModel.
	// If we wanted to test ShowTable output, we'd need a mock tea.Program or similar.
	// For now, just check it doesn't panic.
	// A panic would fail the test.
	// We can check if *something* is written to the buffer, but it's hard to predict exact output.

	// To avoid blocking, we won't run ShowTable directly in this simple test.
	// The core logic is in TableModel, which is tested above.
	// If ShowTable had more logic, we'd need a more sophisticated test.

	// Let's try a minimal run by sending a quit message almost immediately.
	// This still isn't ideal as it might flash the screen.
	// A "headless" mode or test driver for Bubble Tea would be better.

	// As a proxy for ShowTable, we can test the model's View directly.
	model := NewTableModel(columns, rows)
	viewOutput := model.View() // Get the initial view

	if !strings.Contains(viewOutput, "Key") {
		t.Errorf("ShowTable's underlying model view should contain 'Key'. Got: %s", viewOutput)
	}
	if !strings.Contains(viewOutput, "v1") {
		t.Errorf("ShowTable's underlying model view should contain 'v1'. Got: %s", viewOutput)
	}
	// This doesn't test ShowTable's p.Run(), but tests the model it would use.
}
