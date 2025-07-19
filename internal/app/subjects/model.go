package subjects

import (
	"vigenda/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	subjectService service.SubjectService
	// outros campos para a TUI de disciplinas
}

func New(ss service.SubjectService) *Model {
	return &Model{
		subjectService: ss,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	return "Gerenciamento de Disciplinas (TODO)"
}

func (m *Model) CanGoBack() bool {
	return true
}
