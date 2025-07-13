package dashboard

import (
	"context"
	"fmt"
	"log" // Adicionado para depuração
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"vigenda/internal/models"
	"vigenda/internal/service"
	// Outros imports necessários como models, services podem ser adicionados depois
)

// Model representa o estado do componente Dashboard.
type Model struct {
	// Dimensões da área de visualização
	width  int
	height int

	// Serviços para buscar dados
	// Estes serão injetados através do construtor New.
	taskService       service.TaskService
	classService      service.ClassService       // Para buscar aulas, por exemplo (pode ser removido se LessonService for suficiente)
	assessmentService service.AssessmentService  // Para buscar avaliações
	lessonService     service.LessonService      // Adicionado para buscar lições

	// Dados a serem exibidos no Dashboard
	// Estes campos serão populados pelas respostas dos serviços.
	upcomingTasks       []models.Task        // Tarefas com prazos futuros
	todaysLessons       []models.Lesson      // Lições agendadas para o dia atual
	upcomingAssessments []models.Assessment  // Avaliações agendadas futuramente
	// Poderíamos adicionar mais, como:
	// recentGrades      []models.GradeSummary // Resumo de notas recentes lançadas
	// systemMessages    []string              // Mensagens importantes do sistema ou lembretes

	// Estado da UI do Dashboard
	isLoading bool  // True enquanto os dados estiverem sendo carregados
	err       error // Armazena qualquer erro ocorrido durante o carregamento de dados
	// selectedSection int // Se o dashboard tivesse seções navegáveis internamente
}

// New cria uma nova instância do Dashboard Model.
// Parâmetros:
//   ts: Instância de TaskService para buscar dados de tarefas.
//   cs: Instância de ClassService (pode ser removido se não for mais usado diretamente pelo dashboard).
//   as: Instância de AssessmentService para buscar dados de avaliações.
//   ls: Instância de LessonService para buscar dados de lições.
func New(ts service.TaskService, cs service.ClassService, as service.AssessmentService, ls service.LessonService) *Model {
	return &Model{
		taskService:       ts,
		classService:      cs, // Manter por enquanto, pode ser removido se não usado
		assessmentService: as,
		lessonService:     ls,
		isLoading:         true, // Inicia em estado de carregamento por padrão
		// upcomingTasks, todaysLessons, upcomingAssessments são inicializados como slices vazios (nil)
	}
}

// --- Mensagens para comunicação interna do dashboard ---

// upcomingTasksLoadedMsg é enviada quando as tarefas futuras são carregadas.
type upcomingTasksLoadedMsg struct{ tasks []models.Task }

// todaysLessonsLoadedMsg é enviada quando as lições de hoje são carregadas.
type todaysLessonsLoadedMsg struct{ lessons []models.Lesson }

// upcomingAssessmentsLoadedMsg é enviada quando as próximas avaliações são carregadas.
type upcomingAssessmentsLoadedMsg struct{ assessments []models.Assessment }

// dashboardErrorMsg é enviada quando ocorre um erro ao buscar dados para o dashboard.
type dashboardErrorMsg struct{ err error }

// --- Funções de Comando para buscar dados ---

func (m *Model) fetchUpcomingTasks() tea.Cmd {
	return func() tea.Msg {
		// TODO: Implementar a lógica real de busca no TaskService.
		// Exemplo: tasks, err := m.taskService.GetUpcomingTasks(context.Background(), 5, time.Now())
		// Por enquanto, retorna dados mockados ou vazios.
		// Supondo que TaskService tenha um método como ListActiveTasks (ou similar)
		// e precisaremos filtrar por data.
		// Usar o novo método GetUpcomingActiveTasks do TaskService.
		// TODO: Obter UserID do contexto/configuração quando a autenticação estiver implementada.
		userID := int64(1) // Placeholder
		limit := 5         // Mostrar até 5 tarefas futuras

		// Usar o início do dia atual para fromDate para incluir todas as tarefas de hoje.
		now := time.Now()
		fromDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		tasks, err := m.taskService.GetUpcomingActiveTasks(context.Background(), userID, fromDate, limit)
		if err != nil {
			return dashboardErrorMsg{fmt.Errorf("buscar tarefas futuras: %w", err)}
		}
		return upcomingTasksLoadedMsg{tasks: tasks}
	}
}

func (m *Model) fetchTodaysLessons() tea.Cmd {
	return func() tea.Msg {
		// TODO: Obter UserID do contexto/configuração quando a autenticação estiver implementada.
		userID := int64(1) // Placeholder
		lessons, err := m.lessonService.GetLessonsForDate(context.Background(), userID, time.Now())
		if err != nil {
			return dashboardErrorMsg{fmt.Errorf("buscar lições de hoje: %w", err)}
		}
		return todaysLessonsLoadedMsg{lessons: lessons}
	}
}

func (m *Model) fetchUpcomingAssessments() tea.Cmd {
	return func() tea.Msg {
		// TODO: Implementar a lógica real de busca no AssessmentService.
		// Exemplo: assessments, err := m.assessmentService.GetUpcomingAssessments(context.Background(), time.Now(), 7*24*time.Hour) // Próximos 7 dias
		// Por enquanto, retorna dados mockados ou vazios.
		// Supondo que AssessmentService precise de um método específico.
		// Por agora, vamos simular.
		// assessments, err := m.assessmentService.ListAllAssessments(context.Background()) // Se tal método existir
		// if err != nil {
		// 	return dashboardErrorMsg{fmt.Errorf("buscar todas as avaliações: %w", err)}
		// }
		// Placeholder: Filtro manual
		// var upcoming []models.Assessment
		// now := time.Now()
		// for _, assessment := range assessments {
		// 	if assessment.Date.After(now) { // Supondo que assessment.Date exista e seja time.Time
		// 		upcoming = append(upcoming, assessment)
		// 	}
		// }
		// Retornando vazio por enquanto
		return upcomingAssessmentsLoadedMsg{assessments: []models.Assessment{}}
	}
}

// Init é chamado quando o modelo é iniciado.
// Retorna comandos para carregar os dados iniciais do dashboard.
func (m *Model) Init() tea.Cmd {
	m.isLoading = true
	m.err = nil // Limpar erros anteriores
	return tea.Batch(
		m.fetchUpcomingTasks(),
		m.fetchTodaysLessons(), // Corrigido
		m.fetchUpcomingAssessments(),
	)
}

// Update lida com mensagens e atualiza o modelo.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd // Usar para acumular múltiplos comandos

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case upcomingTasksLoadedMsg:
		m.upcomingTasks = msg.tasks
		// Log para depuração
		log.Printf("DashboardModel: upcomingTasksLoadedMsg - %d tarefas recebidas", len(msg.tasks))
		for i, task := range msg.tasks {
			log.Printf("  Tarefa %d: ID %d, Título '%s', DueDate: %v", i+1, task.ID, task.Title, task.DueDate)
		}
		// Não definir isLoading = false aqui, esperar todas as cargas
	case todaysLessonsLoadedMsg: // Corrigido
		m.todaysLessons = msg.lessons // Corrigido
	case upcomingAssessmentsLoadedMsg:
		m.upcomingAssessments = msg.assessments
		// Assumindo que esta é a última mensagem de dados esperada do batch em Init:
		m.isLoading = false
		m.err = nil // Limpar erro se a última carga foi bem-sucedida

	case dashboardErrorMsg:
		m.err = msg.err
		m.isLoading = false // Parar o carregamento em caso de erro

	case tea.KeyMsg:
		// 'esc' para voltar ao menu é tratado pelo app.Model.
		// Adicionar 'r' para recarregar dados do dashboard:
		if msg.String() == "r" {
			m.isLoading = true
			m.err = nil // Limpar erro antes de recarregar
			cmds = append(cmds, m.Init()) // Re-chama Init para buscar dados novamente
		}
	}

	return m, tea.Batch(cmds...)
}

// View renderiza a UI do dashboard.
func (m *Model) View() string {
	if m.isLoading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "Carregando dados do dashboard...")
	}

	if m.err != nil {
		// Estilo para mensagem de erro
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
		errorMsg := fmt.Sprintf("Erro ao carregar dashboard:\n%v\n\nPressione 'r' para tentar novamente ou 'esc' para voltar.", m.err)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, errorStyle.Render(errorMsg))
	}

	// Estilos
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).MarginBottom(1)
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("208")).MarginTop(1)
	listItemStyle := lipgloss.NewStyle().PaddingLeft(2)
	noDataStyle := lipgloss.NewStyle().Italic(true).PaddingLeft(2)
	helpStyle := lipgloss.NewStyle().Faint(true).MarginTop(1)

	// Construtor de String para a View
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("Painel de Controle Vigenda") + "\n")

	// Seção: Tarefas Próximas
	sb.WriteString(sectionTitleStyle.Render("Tarefas Próximas") + "\n")
	if len(m.upcomingTasks) == 0 {
		sb.WriteString(noDataStyle.Render("Nenhuma tarefa próxima encontrada.") + "\n")
	} else {
		for _, task := range m.upcomingTasks {
			dueDateStr := "N/A"
			if task.DueDate != nil {
				dueDateStr = task.DueDate.Format("02/01/2006")
			}
			sb.WriteString(listItemStyle.Render(fmt.Sprintf("• %s (Prazo: %s)", task.Title, dueDateStr)) + "\n")
		}
	}

	// Seção: Aulas de Hoje
	sb.WriteString(sectionTitleStyle.Render("Aulas de Hoje") + "\n")
	if len(m.todaysLessons) == 0 {
		sb.WriteString(noDataStyle.Render("Nenhuma aula programada para hoje.") + "\n")
	} else {
		for _, lesson := range m.todaysLessons {
			sb.WriteString(listItemStyle.Render(fmt.Sprintf("• %s (%s)", lesson.Title, lesson.ScheduledAt.Format("15:04"))))
			// TODO: Adicionar nome da turma se lesson.ClassID for usado para buscar detalhes da turma
			// Ex: className := m.getClassName(lesson.ClassID) // precisaria de acesso ao classService ou um mapa
			sb.WriteString("\n")
		}
	}

	// Seção: Próximas Avaliações
	sb.WriteString(sectionTitleStyle.Render("Próximas Avaliações") + "\n")
	if len(m.upcomingAssessments) == 0 {
		sb.WriteString(noDataStyle.Render("Nenhuma avaliação próxima encontrada.") + "\n")
	} else {
		for _, assessment := range m.upcomingAssessments {
			assessmentDateStr := "N/D"
			if assessment.AssessmentDate != nil {
				assessmentDateStr = assessment.AssessmentDate.Format("02/01/2006")
			}
			sb.WriteString(listItemStyle.Render(fmt.Sprintf("• %s (Turma ID: %d, Data: %s)", assessment.Name, assessment.ClassID, assessmentDateStr)) + "\n")
		}
	}

	sb.WriteString(helpStyle.Render("Pressione 'r' para recarregar. Pressione 'esc' para voltar ao menu."))

	// Usar Place para melhor controle do layout geral, especialmente se o conteúdo for menor que a tela.
	// Para conteúdo que pode exceder a altura, o Place pode não ser ideal sem scroll.
	// Por agora, vamos assumir que o conteúdo cabe.
	return lipgloss.NewStyle().Padding(1,2).Render(sb.String())
	// return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, sb.String(), lipgloss.WithMaxHeight(m.height), lipgloss.WithMaxWidth(m.width))
}

// IsFocused indica se o dashboard tem algum componente interno focado (como um input de texto).
// Isso é usado pelo app.Model para decidir se 'esc' deve voltar ao menu ou ser tratado pelo dashboard.
// Para um dashboard que apenas exibe dados, isso geralmente será false.
func (m *Model) IsFocused() bool {
	return false
}

func (m *Model) IsLoading() bool {
	return m.isLoading
}

// --- Mensagens para carregamento de dados (exemplos) ---
// Estas seriam definidas mais concretamente quando a lógica de fetch for implementada.

// type upcomingTasksLoadedMsg struct {
// 	tasks []models.Task
// }

// type todaysClassesLoadedMsg struct {
// 	classes []models.Class
// }

// type upcomingAssessmentsLoadedMsg struct {
// 	assessments []models.Assessment
// }

// type dashboardErrorMsg struct {
// 	err error
// }

// --- Funções de Fetch (exemplos, a serem implementados) ---
// Estes seriam chamados em Init e retornariam tea.Cmd

// func (m *Model) fetchUpcomingTasks() tea.Cmd {
// 	return func() tea.Msg {
// 		// Simulação:
// 		// tasks, err := m.taskService.GetUpcomingTasks(context.Background(), 5)
// 		// if err != nil {
// 		// 	return dashboardErrorMsg{err}
// 		// }
// 		// return upcomingTasksLoadedMsg{tasks}
// 		return upcomingTasksLoadedMsg{tasks: []models.Task{}} // Placeholder
// 	}
// }

// func (m *Model) fetchTodaysClasses() tea.Cmd {
// 	return func() tea.Msg {
// 		// Lógica para buscar aulas de hoje usando m.classService
// 		return todaysClassesLoadedMsg{classes: []models.Class{}} // Placeholder
// 	}
// }

// func (m *Model) fetchUpcomingAssessments() tea.Cmd {
// 	return func() tea.Msg {
// 		// Lógica para buscar próximas avaliações usando m.assessmentService
// 		return upcomingAssessmentsLoadedMsg{assessments: []models.Assessment{}} // Placeholder
// 	}
// }
