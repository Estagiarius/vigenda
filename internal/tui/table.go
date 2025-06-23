package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type TableModel struct {
	table table.Model
}

func NewTableModel(columns []table.Column, rows []table.Row) TableModel {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return TableModel{table: t}
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

// Helper function to create and run a simple table
func ShowTable(columns []table.Column, rows []table.Row, output io.Writer) {
	p := tea.NewProgram(NewTableModel(columns, rows), tea.WithOutput(output))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(output, "Error running table program: %v\n", err)
	}
}

// Example Usage (can be moved to a main or test file)
/*
func main() {
    columns := []table.Column{
        {Title: "ID", Width: 4},
        {Title: "Name", Width: 10},
        {Title: "Age", Width: 5},
    }

    rows := []table.Row{
        {"1", "Alice", "30"},
        {"2", "Bob", "24"},
        {"3", "Charlie", "35"},
    }

    ShowTable(columns, rows, os.Stdout)
}
*/
