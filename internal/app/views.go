// Package app contém a lógica principal da Interface de Texto do Usuário (TUI)
// da aplicação Vigenda, utilizando o framework BubbleTea.
// Este arquivo (views.go) define os diferentes estados de visualização (telas/módulos) da TUI.
package app

// View é um tipo enumerado que representa os diferentes módulos ou telas
// que podem ser exibidos na Interface de Texto do Usuário (TUI) principal.
// Cada valor de View corresponde a uma funcionalidade específica da aplicação.
type View int

const (
	// DashboardView representa o estado inicial ou o menu principal da TUI.
	// Atua como um contêiner para navegação para outras visualizações.
	DashboardView View = iota

	// TaskManagementView representa a tela para gerenciamento de tarefas.
	TaskManagementView

	// ClassManagementView representa a tela para gerenciamento de turmas e alunos.
	ClassManagementView

	// AssessmentManagementView representa a tela para gerenciamento de avaliações e notas.
	AssessmentManagementView

	// QuestionBankView representa a tela para gerenciamento do banco de questões.
	QuestionBankView

	// ProofGenerationView representa a tela para geração de provas.
	ProofGenerationView

	// ConcreteDashboardView representa a tela de conteúdo real do painel de controle/dashboard.
	// Distingue-se de DashboardView (menu principal) para permitir uma navegação clara.
	ConcreteDashboardView

	// StudentView é um exemplo de uma sub-visualização, possivelmente para listar ou editar alunos.
	// O seu uso e contexto exato podem depender de como o ClassManagementView é implementado.
	// NOTA: Este valor (99) está fora da sequência iota e foi usado em tui.go;
	// considerar reavaliar sua necessidade ou integrá-lo melhor ao enum.
	StudentView View = 99 // Usado em tui.go, verificar se é necessário aqui ou se é específico daquele contexto.
	// TODO: Adicionar outras visualizações conforme são desenvolvidas, mantendo o iota ou usando valores explícitos.
)

// String retorna a representação textual amigável de uma View.
// Usado para exibir títulos de menu ou identificar a visualização atual.
// É crucial que a ordem e o número de strings neste array correspondam
// exatamente às constantes View definidas acima (exceto para valores explícitos como StudentView).
func (v View) String() string {
	switch v {
	case DashboardView:
		return "Menu Principal" // DashboardView atua como o menu principal.
	case TaskManagementView:
		return "Gerenciar Tarefas"
	case ClassManagementView:
		return "Gerenciar Turmas e Alunos"
	case AssessmentManagementView:
		return "Gerenciar Avaliações e Notas"
	case QuestionBankView:
		return "Banco de Questões"
	case ProofGenerationView:
		return "Gerar Provas"
	case ConcreteDashboardView:
		return "Painel de Controle"
	case StudentView: // Caso para o valor explícito
		return "Visualizar Alunos" // Ou um nome mais apropriado
	default:
		return "Visualização Desconhecida"
	}
}
