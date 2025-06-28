package app

// View represents the current active view/module in the TUI.
type View int

const (
	DashboardView View = iota
	TaskManagementView
	ClassManagementView
	AssessmentManagementView
	QuestionBankView
	ProofGenerationView
	// Add other views as needed
)

func (v View) String() string {
	return [...]string{
		"Dashboard",
		"Gerenciar Tarefas",
		"Gerenciar Turmas e Alunos",
		"Gerenciar Avaliações e Notas",
		"Banco de Questões",
		"Gerar Provas",
	}[v]
}
