package app

// View represents the current active view/module in the TUI.
type View int

const (
	DashboardView View = iota // Represents the main menu container
	TaskManagementView
	ClassManagementView
	AssessmentManagementView
	QuestionBankView
	ProofGenerationView
	ConcreteDashboardView // Represents the actual dashboard content view
	LessonManagementView  // Added Lesson Management View
	// Add other views as needed
)

func (v View) String() string {
	// Ensure this array matches the order and number of constants in View
	return [...]string{
		"Menu Principal", // DashboardView is the main menu
		"Gerenciar Tarefas",
		"Gerenciar Turmas e Alunos",
		"Gerenciar Avaliações e Notas",
		"Banco de Questões",
		"Gerar Provas",
		"Painel de Controle",     // String for ConcreteDashboardView
		"Gerenciar Aulas/Lições", // String for LessonManagementView
	}[v]
}
