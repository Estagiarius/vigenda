package app

// View represents the current active view/module in the TUI.
type View int

const (
	MainMenuView View = iota // Main menu is now the first view
	DashboardView
	TaskManagementView
	ClassManagementView
	AssessmentManagementView
	QuestionBankView
	ProofGenerationView
	// Add other views as needed
	// StudentListView is a temporary view, might be refactored later
	StudentListView // Represents app.View(99) for now
)

func (v View) String() string {
	// Ensure this array matches the order and count of the View constants
	return [...]string{
		"Menu Principal",
		"Dashboard",
		"Gerenciar Tarefas",
		"Gerenciar Turmas e Alunos",
		"Gerenciar Avaliações e Notas",
		"Banco de Questões",
		"Gerar Provas",
		"Alunos da Turma", // String for StudentListView
	}[v]
}
