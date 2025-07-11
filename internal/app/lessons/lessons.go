package lessons

import (
	"context" // Added import
	"fmt"     // Added import
	"strconv" // Added import (was missing from previous diff, but likely needed for ParseInt)
	"strings" // Added import (was missing from previous diff, but likely needed for strings.Builder)
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss" // Adicionado para estilos
	"vigenda/internal/common/styles"    // Changed import
	"vigenda/internal/models"
	"vigenda/internal/service"
)

// ViewState define os diferentes estados da UI de lições.
type ViewState int

const (
	ListView ViewState = iota
	CreatingView
	EditingView
	ConfirmDeleteView
	ErrorView // Para exibir erros de forma mais proeminente
)

// formFocusableInput é uma interface para unificar textinput e textarea para navegação
type formFocusableInput interface {
	Focus() tea.Cmd
	Blur()
	Focused() bool
	SetValue(string)
	Value() string
}

// Model representa o estado do módulo de gerenciamento de lições.
type Model struct {
	lessonService service.LessonService
	classService  service.ClassService // Para validar ClassID e buscar nomes de turmas
	userID        int64                // Para operações de serviço que o exigem

	currentView ViewState
	table       table.Model

	// Formulário para criar/editar
	classIDInput     textinput.Model // ID da Turma
	titleInput       textinput.Model // Título da Lição
	planContentInput textarea.Model  // Conteúdo do Plano de Aula
	scheduledAtInput textinput.Model // Data Agendada (ex: YYYY-MM-DD HH:MM)

	// formFocusOrder gerencia os inputs focáveis no formulário
	// 0: classIDInput, 1: titleInput, 2: planContentInput, 3: scheduledAtInput
	formFocusOrder []formFocusableInput
	focusedIndex   int // Índice do input focado em formFocusOrder

	allLessons     []models.Lesson  // Cache das lições carregadas
	selectedLesson *models.Lesson   // Lição selecionada para edição/exclusão
	classCache     map[int64]string // Cache de nomes de turmas para exibição na tabela

	isLoading      bool
	errorMessage   string // Para armazenar mensagens de erro
	successMessage string // Para armazenar mensagens de sucesso
	width          int
	height         int
	keyMap         KeyMap // Renomeado para KeyMap para seguir convenção
	help           help.Model

	// Estilos (podem vir de um pacote tui comum)
	listHeaderStyle lipgloss.Style
	errorStyle      lipgloss.Style
	successStyle    lipgloss.Style
	helpStyle       lipgloss.Style
}

func New(ls service.LessonService, cs service.ClassService, initialUserID int64, width, height int) Model {
	// --- Table Initialization ---
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Título", Width: 30},
		{Title: "Turma", Width: 20},         // Exibirá Nome da Turma (ou ID se nome não disponível)
		{Title: "Data Agendada", Width: 16}, // Formato YYYY-MM-DD HH:MM
	}
	// Heights an widths need to be adjusted dynamically or after first WindowSizeMsg
	tbl := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10), // Placeholder height
	)
	// Estilo da tabela (pode ser mais elaborado)
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
	tbl.SetStyles(s)

	// --- Form Inputs Initialization ---
	classIDInput := textinput.New()
	classIDInput.Placeholder = "ID da Turma (Numérico)"
	classIDInput.CharLimit = 10
	classIDInput.Width = 20
	// classIDInput.Focus() // O primeiro campo foca por padrão

	titleInput := textinput.New()
	titleInput.Placeholder = "Ex: Introdução às Funções Trigonométricas"
	titleInput.CharLimit = 100
	titleInput.Width = 50

	planContentInput := textarea.New()
	planContentInput.Placeholder = "Detalhes do plano de aula em Markdown..."
	planContentInput.SetWidth(50) // Será ajustado dinamicamente
	planContentInput.SetHeight(5) // Será ajustado dinamicamente

	scheduledAtInput := textinput.New()
	scheduledAtInput.Placeholder = "YYYY-MM-DD HH:MM (ex: 2024-07-21 14:30)"
	scheduledAtInput.CharLimit = 16 // Length of "YYYY-MM-DD HH:MM"
	scheduledAtInput.Width = 20

	// Wrapper para textarea.Model para adaptá-la a formFocusableInput
	// No entanto, textarea.Model já possui Focus(), Blur(), Focused(), SetValue(), Value()
	// então a conversão direta para a interface deve funcionar se os métodos tiverem a mesma assinatura.
	// Se não, um wrapper seria:
	// type textAreaWrapper struct { textarea.Model }
	// func (t *textAreaWrapper) Focus() tea.Cmd { return t.Model.Focus() } ... etc.
	// Mas vamos tentar direto primeiro.

	// --- Table and Help Initialization (moved before m so they can be used) ---
	tbl := table.New(
		table.WithColumns(columns), // columns defined earlier
		table.WithFocused(true),
		table.WithHeight(10), // Placeholder height
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
	tbl.SetStyles(s)

	helpModel := help.New()
	helpModel.Width = width // Initial width

	// Estilos básicos (idealmente de um pacote de tema ou tui)
	listHeaderStyle := lipgloss.NewStyle().Bold(true).Padding(0, 1)
	errorStyle := lipgloss.NewStyle().Foreground(styles.Colors.ErrorColor)     // Use new package
	successStyle := lipgloss.NewStyle().Foreground(styles.Colors.SuccessColor) // Use new package
	helpStyle := lipgloss.NewStyle().Foreground(styles.Colors.HelpColor)       // Use new package

	// --- Construct the Model instance ---
	m := Model{
		lessonService:    ls,
		classService:     cs,
		userID:           initialUserID,
		currentView:      ListView,
		table:            tbl,
		classIDInput:     classIDInput,
		titleInput:       titleInput,
		planContentInput: planContentInput,
		scheduledAtInput: scheduledAtInput,
		// formFocusOrder will be set next
		focusedIndex:     0,
		isLoading:        true, // Começa carregando
		help:             helpModel,
		keyMap:           DefaultKeyMap(),
		width:            width,
		height:           height,
		classCache:       make(map[int64]string),
		listHeaderStyle:  listHeaderStyle,
		errorStyle:       errorStyle,
		successStyle:     successStyle,
		helpStyle:        helpStyle,
	}

	// --- Initialize formFocusOrder using pointers to m's fields ---
	m.formFocusOrder = []formFocusableInput{
		&m.classIDInput,
		&m.titleInput,
		&m.planContentInput,
		&m.scheduledAtInput,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	// Retorna um comando que envia uma mensagem lessonsLoadedMsg com uma lista vazia.
	// Isso fará com que isLoading seja definido como false e a UI seja renderizada.
	// A carga real das lições pode ser acionada por uma ação do usuário (ex: refresh)
	// ou ao selecionar uma turma (se implementado).
	return func() tea.Msg {
		return lessonsLoadedMsg{lessons: []models.Lesson{}} // Envia lista vazia
	}
}

// --- Comandos ---

func fetchLessonsCmd(service service.LessonService, userID int64, classID *int64) tea.Cmd {
	return func() tea.Msg {
		// Por enquanto, vamos buscar todas as lições do usuário se classID for nil.
		// TODO: Adicionar filtro por turma (classID) se fornecido.
		// A interface LessonService.GetLessonsByClassID espera um classID não-ponteiro.
		// E LessonService.GetLessonsForDate espera um UserID.
		// Precisamos de um método no serviço como GetLessons(userID, ?classID) ou ajustar.
		// Temporariamente, se não houver classID, usamos GetLessonsForDate para um período amplo (improvável de ser útil)
		// ou criamos um método no serviço para "todas as lições do usuário".
		// Para este exemplo, vamos assumir que queremos todas as lições de uma turma específica,
		// ou todas as lições do usuário se nenhuma turma for especificada (requer novo método de serviço).
		// Por simplicidade, vamos apenas simular uma busca por classID se fornecido,
		// ou um erro/lista vazia se não.
		// Ideal: ctx := context.Background()
		// Se quisermos todas as lições do usuário, precisaríamos de algo como:
		// lessons, err := service.GetAllLessonsByUserID(ctx, userID)
		// Por enquanto, vamos apenas retornar uma lista vazia para simular e focar na UI.
		// Se classID for fornecido, podemos tentar buscar por ele.
		// if classID != nil {
		//  lessons, err := service.GetLessonsByClassID(context.Background(), *classID)
		//  if err != nil {
		//   return errMsg{err}
		//  }
		//  return lessonsLoadedMsg{lessons}
		// }
		// return lessonsLoadedMsg{lessons: []models.Lesson{}} // Placeholder
		//
		// Vamos simular uma busca que pode falhar ou ter sucesso.
		// Em uma implementação real, chamaríamos o serviço aqui.
		// Exemplo: lessons, err := ls.lessonService.GetLessonsByClassID(context.TODO(), someClassID)
		// O serviço LessonService tem GetLessonsByClassID(ctx, classID) e GetLessonsForDate(ctx, userID, date)
		// Para buscar "todas as lições do usuário logado", precisaríamos de um novo método no serviço
		// ou iterar sobre as turmas do usuário e depois buscar lições por turma.
		// Por simplicidade, se classID for nil, vamos usar GetLessonsForDate para um período muito amplo
		// como uma aproximação de "todas as lições visíveis para o usuário".
		// O userID já está disponível no model e é passado para este comando.
		ctx := context.Background()
		var lessons []models.Lesson
		var err error

		if classID != nil {
			lessons, err = service.GetLessonsByClassID(ctx, *classID)
			if err != nil {
				return errMsg{fmt.Errorf("falha ao buscar lições por turma: %w", err)}
			}
		} else {
			// Para buscar todas as lições do usuário, precisamos de um método de serviço apropriado.
			// Ex: service.GetAllLessonsByUserID(ctx, userID)
			// Simulação: Usar GetLessonsForDate com um range amplo.
			// Esta não é uma solução ideal, mas demonstra o uso do serviço.
			// Um método dedicado no serviço seria melhor.
	// startDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) // endDate was unused. startDate also unused with current logic.
	// endDate := time.Date(2070, 1, 1, 0, 0, 0, 0, time.UTC) // Unused
	// lessons, err = service.GetLessonsForDate(ctx, userID, today) // This line was problematic and is removed/handled by the else block
			// Para um range, o repo tem GetLessonsByDateRange.
			// O serviço precisa expor isso.
			// Assumindo que LessonService é atualizado para ter GetLessonsByDateRange(ctx, userID, startDate, endDate)
			// ou que GetLessonsForDate é alterado para um range, ou que temos outro método.
			// Por agora, vamos assumir que GetLessonsForDate com a data de início do range funciona para este exemplo
			// ou que o plano é buscar todas as lições de um usuário (o que GetLessonsForDate não faz diretamente para todas).
			// Vamos usar um método hipotético GetAllUserLessons (que chamaria o repo GetLessonsByDateRange com userID)
			// if service has GetAllUserLessons(ctx context.Context, userID int64) ([]models.Lesson, error)
			// lessons, err = service.GetAllUserLessons(ctx, userID)
			// Se não, a melhor aproximação com os métodos atuais é iterar turmas ou usar um range amplo no repo.
			// Dado o serviço atual, a busca sem classID é mais complexa.
			// Para este exemplo, vamos simplificar: se classID não for fornecido, retornaremos uma lista vazia
			// e o usuário DEVE filtrar por turma ou uma funcionalidade de "ver todas as minhas lições"
			// precisaria de um backend mais robusto.
			// Para o propósito deste TUI, vamos focar em carregar por turma se o ID for dado,
			// ou carregar as lições para uma data específica usando `GetLessonsForDate`.
			// O plano original sugere "buscar todas as lições para o usuário ou as lições da primeira turma disponível".
			// O `fetchLessonsCmd` atual é chamado com `classID = nil` em `fetchLessons()`.
			// Vamos simular que o serviço tem um `GetAllLessons(ctx, userID)`
			// Para o mock, vamos usar o GetLessonsForDate com a data de hoje e um range grande.
			// Esta parte precisa de um design de serviço mais claro para "todas as lições do usuário".
			// Para manter o progresso, se classID for nil, retornaremos as lições de hoje para o usuário.
			// startDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) // endDate was unused
			today := time.Now()
			lessons, err = service.GetLessonsForDate(ctx, userID, today)
			if err != nil {
				return errMsg{fmt.Errorf("falha ao buscar lições para hoje: %w", err)}
			}
		}
		return lessonsLoadedMsg{lessons: lessons}
	}
}

func createLessonCmd(service service.LessonService, lesson models.Lesson, userID int64) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		// O serviço CreateLesson já lida com a validação de propriedade da turma com base no userID (placeholder no serviço).
		// Precisamos garantir que o userID correto seja passado ou obtido no serviço.
		// O modelo Lesson não tem UserID, mas ClassID sim, e Class tem UserID.
		// O serviço CreateLesson já faz essa checagem (com userID placeholder).
		createdLesson, err := service.CreateLesson(ctx, lesson.ClassID, lesson.Title, lesson.PlanContent, lesson.ScheduledAt)
		if err != nil {
			return errMsg{fmt.Errorf("falha ao criar lição: %w", err)}
		}
		return lessonCreatedMsg{lesson: createdLesson}
	}
}

func updateLessonCmd(service service.LessonService, lesson models.Lesson, userID int64) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		// O serviço UpdateLesson também deve validar a propriedade.
		updatedLesson, err := service.UpdateLesson(ctx, lesson.ID, lesson.Title, lesson.PlanContent, lesson.ScheduledAt)
		if err != nil {
			return errMsg{fmt.Errorf("falha ao atualizar lição: %w", err)}
		}
		return lessonUpdatedMsg{lesson: updatedLesson}
	}
}

func deleteLessonCmd(service service.LessonService, lessonID int64, userID int64) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		// O serviço DeleteLesson também deve validar a propriedade.
		err := service.DeleteLesson(ctx, lessonID)
		if err != nil {
			return errMsg{fmt.Errorf("falha ao excluir lição: %w", err)}
		}
		return lessonDeletedMsg{}
	}
}

// fetchClassNameCmd busca o nome de uma turma.
// Usado para popular o cache de nomes de turmas.
func fetchClassNameCmd(classService service.ClassService, classID int64, userID int64) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		// O ClassService.GetClassByID não necessariamente valida o UserID.
		// A validação de propriedade da turma deve ocorrer no LessonService ao operar sobre lições.
		// Para exibir o nome, apenas buscamos. Se a turma não existir, o serviço retornará erro.
		class, err := classService.GetClassByID(ctx, classID)
		if err != nil {
			// Não tratar como erro fatal para a lista de lições, apenas não mostrará o nome.
			// Retornar uma mensagem de erro específica ou simplesmente o ID como nome.
			return classNameLoadedMsg{classID: classID, name: fmt.Sprintf("ID: %d (erro ao buscar nome)", classID), err: err}
		}
		return classNameLoadedMsg{classID: classID, name: class.Name}
	}
}

// --- Mensagens Adicionais ---
type classNameLoadedMsg struct {
	classID int64
	name    string
	err     error // Added field to carry error from fetchClassNameCmd
}

// focusFormMsg é usada para sinalizar que o foco deve ir para o formulário
type focusFormMsg struct{}

// GoBackMsg é usada para sinalizar ao app.Model para voltar ao menu principal.
type GoBackMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd // Initialize cmds to an empty slice

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		// Ajustar altura da tabela e largura/altura do textarea
		headerHeight := lipgloss.Height(m.ViewHeader())
		footerHeight := lipgloss.Height(m.ViewFooter()) // Help + possible status line
		formHeight := 0
		if m.currentView == CreatingView || m.currentView == EditingView {
			formHeight = lipgloss.Height(m.ViewForm())
		}
		m.table.SetHeight(m.height - headerHeight - footerHeight - formHeight - 1) // -1 for margin/border
		m.planContentInput.SetWidth(m.width - 6)                                   // Ajustar conforme padding do form
		// Altura do textarea pode ser uma porcentagem ou fixa
		m.planContentInput.SetHeight(m.height / 4)

	case lessonsLoadedMsg:
		m.isLoading = false
		m.allLessons = msg.lessons
		var rows []table.Row
		var classCmds []tea.Cmd
		for _, lesson := range m.allLessons {
			className, ok := m.classCache[lesson.ClassID]
			if !ok {
				// Se não estiver no cache, buscar e adicionar ao cache
				// Poderia otimizar para buscar apenas uma vez por ClassID único na lista
				if _, exists := m.classCache[lesson.ClassID]; !exists { // Evitar múltiplas buscas para o mesmo ID
					m.classCache[lesson.ClassID] = "Carregando..." // Placeholder
					classCmds = append(classCmds, fetchClassNameCmd(m.classService, lesson.ClassID, m.userID))
				}
				className = "Carregando..."
			}
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", lesson.ID),
				lesson.Title,
				className,
				lesson.ScheduledAt.Format("2006-01-02 15:04"),
			})
		}
		m.table.SetRows(rows)
		if len(classCmds) > 0 {
			cmds = append(cmds, tea.Batch(classCmds...))
		}
		m.errorMessage = "" // Limpar erro anterior ao carregar com sucesso

	case classNameLoadedMsg:
		m.classCache[msg.classID] = msg.name
		// Re-renderizar a tabela se estiver na ListView para mostrar o nome da turma atualizado
		if m.currentView == ListView {
			// Reconstruir linhas da tabela com nomes de turma atualizados
			var rows []table.Row
			for _, lesson := range m.allLessons {
				className := m.classCache[lesson.ClassID] // Deve existir agora ou ser o nome carregado
				rows = append(rows, table.Row{
					fmt.Sprintf("%d", lesson.ID),
					lesson.Title,
					className,
					lesson.ScheduledAt.Format("2006-01-02 15:04"),
				})
			}
			m.table.SetRows(rows)
		}

	case lessonCreatedMsg:
		m.isLoading = false
		m.successMessage = fmt.Sprintf("Lição '%s' criada com sucesso!", msg.lesson.Title)
		m.currentView = ListView
		// Recarregar todas as lições para incluir a nova
		cmds = append(cmds, m.fetchLessons())

	case lessonUpdatedMsg:
		m.isLoading = false
		m.successMessage = fmt.Sprintf("Lição '%s' atualizada com sucesso!", msg.lesson.Title)
		m.currentView = ListView
		cmds = append(cmds, m.fetchLessons())

	case lessonDeletedMsg:
		m.isLoading = false
		m.successMessage = "Lição excluída com sucesso!"
		m.currentView = ListView
		cmds = append(cmds, m.fetchLessons())

	case errMsg:
		m.isLoading = false
		m.errorMessage = msg.Error()
		// Não mudar para ErrorView automaticamente, deixar a view atual mostrar o erro
		// m.currentView = ErrorView // Ou mostrar na view atual
		// Limpar mensagem de sucesso se houver erro
		m.successMessage = ""
		return m, nil // Não processar mais comandos se houver erro

	case tea.KeyMsg:
		// Limpar mensagens de erro/sucesso em qualquer tecla, exceto quando visualizando erro
		if m.currentView != ErrorView {
			m.errorMessage = ""
			m.successMessage = ""
		}

		switch m.currentView {
		case ListView:
			cmd = m.handleListViewKeyPress(msg)
		case CreatingView, EditingView:
			cmd = m.handleFormViewKeyPress(msg)
		case ConfirmDeleteView:
			cmd = m.handleConfirmDeleteViewKeyPress(msg)
		case ErrorView: // Se ErrorView for um estado que captura todas as teclas
			if key.Matches(msg, m.keyMap.Back) {
				m.currentView = ListView // Voltar para a lista
				m.errorMessage = ""      // Limpar erro
			}
		}
		cmds = append(cmds, cmd)

	// --- Tratamento de foco e submissão do formulário ---
	case focusFormMsg: // Mensagem para focar no primeiro campo do formulário
		if m.currentView == CreatingView || m.currentView == EditingView {
			if len(m.formFocusOrder) > 0 {
				m.focusedIndex = 0
				cmds = append(cmds, m.formFocusOrder[m.focusedIndex].Focus())
			}
		}

	}

	// Processar inputs do formulário se estiver em modo de formulário e focado
	if m.currentView == CreatingView || m.currentView == EditingView {
		var formCmd tea.Cmd
		// Atualizar o input focado.
		// A msg é passada para o input que está atualmente focado.
		// Os métodos Update dos inputs (textinput, textarea) retornam o modelo atualizado e um comando.
		switch m.focusedIndex {
		case 0: // classIDInput
			m.classIDInput, formCmd = m.classIDInput.Update(msg)
		case 1: // titleInput
			m.titleInput, formCmd = m.titleInput.Update(msg)
		case 2: // planContentInput
			m.planContentInput, formCmd = m.planContentInput.Update(msg)
		case 3: // scheduledAtInput
			m.scheduledAtInput, formCmd = m.scheduledAtInput.Update(msg)
		}
		cmds = append(cmds, formCmd)
	}

	// Atualizar tabela se estiver visível
	if m.currentView == ListView {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Atualizar a view da ajuda
	// A ajuda contextual depende do m.currentView, então chamamos UpdateCurrentViewForHelp
	UpdateCurrentViewForHelp(m.currentView) // Atualiza a variável global para keymap
	m.help.Update(msg)                      // Isso pode não ser necessário se a ajuda é estática por view

	return m, tea.Batch(cmds...)
}

// --- Funções Auxiliares para Update ---

func (m *Model) fetchLessons() tea.Cmd {
	m.isLoading = true
	m.errorMessage = ""
	m.successMessage = ""
	// Por enquanto, classID é nil, o que significa que fetchLessonsCmd precisa lidar com isso
	// ou precisamos de um estado para a turma selecionada.
	// Para o exemplo, vamos assumir que sempre buscamos "todas" as lições do usuário.
	// O serviço GetLessonsByDateRange precisa de UserID, startDate, endDate.
	// Para "todas", podemos usar um range muito amplo ou um novo método de serviço.
	// Vamos usar o fetchLessonsCmd que já tem uma lógica mock.
	return fetchLessonsCmd(m.lessonService, m.userID, nil)
}

func (m *Model) handleListViewKeyPress(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keyMap.Back): // Esc
		// Sinalizar para o app.Model para voltar ao menu principal
		return func() tea.Msg { return GoBackMsg{} }

	case key.Matches(msg, m.keyMap.Quit):
		return tea.Quit

	case key.Matches(msg, m.keyMap.New): // 'n'
		m.currentView = CreatingView
		m.resetFormFields()
		m.focusedIndex = 0 // Focar no primeiro campo
		return func() tea.Msg { return focusFormMsg{} }

	case key.Matches(msg, m.keyMap.Edit): // 'e'
		if len(m.table.Rows()) > 0 && m.table.Cursor() < len(m.allLessons) {
			// Garantir que o cursor da tabela esteja dentro dos limites de allLessons
			idx := m.table.Cursor()
			if idx >= 0 && idx < len(m.allLessons) {
				m.selectedLesson = &m.allLessons[idx]
				m.populateFormFields(*m.selectedLesson)
				m.currentView = EditingView
				m.focusedIndex = 0 // Focar no primeiro campo
				return func() tea.Msg { return focusFormMsg{} }
			}
		}

	case key.Matches(msg, m.keyMap.Delete): // 'd'
		if len(m.table.Rows()) > 0 && m.table.Cursor() < len(m.allLessons) {
			idx := m.table.Cursor()
			if idx >= 0 && idx < len(m.allLessons) {
				m.selectedLesson = &m.allLessons[idx]
				m.currentView = ConfirmDeleteView
			}
		}

	case key.Matches(msg, m.keyMap.Refresh): // 'r'
		return m.fetchLessons()

	// Adicionar navegação na tabela (cima, baixo, pgup, pgdown, home, end)
	// Estes são geralmente tratados pela própria table.Model.Update,
	// mas precisamos garantir que a tabela receba a mensagem.
	default:
		// Se não for uma tecla específica do ListView, passar para a tabela
		// var cmd tea.Cmd
		// m.table, cmd = m.table.Update(msg)
		// return cmd
		// O update da tabela já está sendo feito no final da função Update principal
	}
	return nil
}

func (m *Model) handleFormViewKeyPress(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keyMap.Cancel): // Esc
		m.currentView = ListView
		m.resetFormFields()
		m.errorMessage = ""
		m.successMessage = ""
		// Desfocar o campo atual
		if m.focusedIndex >= 0 && m.focusedIndex < len(m.formFocusOrder) {
			m.formFocusOrder[m.focusedIndex].Blur()
		}
		return nil

	case key.Matches(msg, m.keyMap.Submit): // Enter (ou Ctrl+S se configurado)
		// Se o textarea estiver focado, Enter deve criar nova linha, não submeter.
		// A submissão pode ser por Tab no último campo e Enter, ou uma tecla dedicada (Ctrl+S).
		// Por ora, se o Enter for pressionado e o último campo (scheduledAtInput) estiver focado, submete.
		// Ou se o textarea não estiver focado.
		if !m.planContentInput.Focused() || (m.focusedIndex == len(m.formFocusOrder)-1 && m.formFocusOrder[m.focusedIndex].Focused()) {
			return m.submitForm()
		}
		// Se o textarea estiver focado e Enter for pressionado, o textarea.Update() tratará (nova linha).

	case key.Matches(msg, m.keyMap.NextField): // Tab
		return m.focusNextInput()

	case key.Matches(msg, m.keyMap.PrevField): // Shift+Tab
		return m.focusPrevInput()

		// Teclas de navegação dentro do textarea (cima/baixo) devem ser tratadas pelo textarea
		// Se o textarea estiver focado, ele deve consumir essas teclas.
		// Se não, podemos ter navegação cima/baixo entre campos aqui.
		// Por simplicidade, Tab/Shift+Tab são os principais para navegação entre campos.
		// O textinput/textarea.Update já é chamado no final da função Update.
	}
	return nil
}

func (m *Model) handleConfirmDeleteViewKeyPress(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keyMap.Submit): // 's', 'y', ou Enter para confirmar
		if m.selectedLesson != nil {
			m.isLoading = true
			cmd := deleteLessonCmd(m.lessonService, m.selectedLesson.ID, m.userID)
			m.selectedLesson = nil // Limpar após iniciar a exclusão
			return cmd
		}
	case key.Matches(msg, m.keyMap.Cancel): // 'n' ou Esc para cancelar
		m.currentView = ListView
		m.selectedLesson = nil
		m.errorMessage = ""
	}
	return nil
}

func (m *Model) resetFormFields() {
	m.classIDInput.SetValue("")
	m.classIDInput.Blur() // Garantir que não esteja focado
	m.titleInput.SetValue("")
	m.titleInput.Blur()
	m.planContentInput.SetValue("")
	m.planContentInput.Blur()
	m.scheduledAtInput.SetValue("")
	m.scheduledAtInput.Blur()
	m.selectedLesson = nil
	m.errorMessage = "" // Limpar erros ao resetar formulário
	m.focusedIndex = -1 // Nenhum campo focado inicialmente ao resetar
}

func (m *Model) populateFormFields(lesson models.Lesson) {
	m.classIDInput.SetValue(fmt.Sprintf("%d", lesson.ClassID))
	m.titleInput.SetValue(lesson.Title)
	m.planContentInput.SetValue(lesson.PlanContent)
	m.scheduledAtInput.SetValue(lesson.ScheduledAt.Format("2006-01-02 15:04"))
}

func (m *Model) submitForm() tea.Cmd {
	// Validação básica
	if m.classIDInput.Value() == "" {
		m.errorMessage = "ID da Turma é obrigatório."
		m.currentView = EditingView                     // Manter na view do formulário para corrigir
		return func() tea.Msg { return focusFormMsg{} } // Focar no primeiro campo
	}
	classID, err := strconv.ParseInt(m.classIDInput.Value(), 10, 64)
	if err != nil {
		m.errorMessage = "ID da Turma deve ser um número."
		m.currentView = EditingView
		m.classIDInput.Focus()
		return nil
	}
	if m.titleInput.Value() == "" {
		m.errorMessage = "Título da lição é obrigatório."
		m.currentView = EditingView
		m.titleInput.Focus()
		return nil
	}
	// Validação de data/hora
	scheduledAtStr := m.scheduledAtInput.Value()
	if scheduledAtStr == "" { // Permitir data agendada vazia? O modelo tem time.Time, não *time.Time
		m.errorMessage = "Data agendada é obrigatória." // Assumindo que é obrigatória
		m.currentView = EditingView
		m.scheduledAtInput.Focus()
		return nil
	}
	scheduledAt, err := time.Parse("2006-01-02 15:04", scheduledAtStr)
	if err != nil {
		m.errorMessage = "Data agendada inválida. Use o formato YYYY-MM-DD HH:MM."
		m.currentView = EditingView
		m.scheduledAtInput.Focus()
		return nil
	}

	// Validar se a turma (classID) existe e pertence ao usuário (usando classService)
	// O userID da lição é inferido pelo userID da turma.
	// O LessonService já faz uma validação de propriedade da turma ao criar/atualizar lições,
	// mas uma checagem antecipada aqui na TUI pode dar feedback mais rápido.
	// No entanto, LessonService.validateUserOwnsClass usa um userID placeholder (1).
	// Para uma validação real aqui, precisaríamos do userID correto.
	// Por agora, vamos apenas validar a existência da turma.
	// A validação de propriedade efetiva ocorrerá no LessonService.
	_, err = m.classService.GetClassByID(context.Background(), classID)
	if err != nil {
		m.errorMessage = fmt.Sprintf("Turma com ID %d não encontrada: %v", classID, err)
		if m.currentView == CreatingView { // Manter na view correta
			m.currentView = CreatingView
		} else {
			m.currentView = EditingView
		}
		m.classIDInput.Focus()
		return nil
	}

	lessonData := models.Lesson{
		ClassID:     classID,
		Title:       m.titleInput.Value(),
		PlanContent: m.planContentInput.Value(),
		ScheduledAt: scheduledAt,
	}
	m.isLoading = true
	// Desfocar todos os campos do formulário antes de submeter
	for _, input := range m.formFocusOrder {
		input.Blur()
	}

	if m.currentView == EditingView && m.selectedLesson != nil {
		lessonData.ID = m.selectedLesson.ID
		return updateLessonCmd(m.lessonService, lessonData, m.userID)
	}
	return createLessonCmd(m.lessonService, lessonData, m.userID)
}

func (m *Model) focusNextInput() tea.Cmd {
	if len(m.formFocusOrder) == 0 {
		return nil
	}
	m.formFocusOrder[m.focusedIndex].Blur()
	m.focusedIndex = (m.focusedIndex + 1) % len(m.formFocusOrder)
	return m.formFocusOrder[m.focusedIndex].Focus()
}

func (m *Model) focusPrevInput() tea.Cmd {
	if len(m.formFocusOrder) == 0 {
		return nil
	}
	m.formFocusOrder[m.focusedIndex].Blur()
	m.focusedIndex--
	if m.focusedIndex < 0 {
		m.focusedIndex = len(m.formFocusOrder) - 1
	}
	return m.formFocusOrder[m.focusedIndex].Focus()
}

// ViewHeader, ViewFooter, ViewForm (para ajudar no cálculo de altura da tabela)
func (m Model) ViewHeader() string {
	// Placeholder para um possível título ou cabeçalho do módulo
	return "" // Ex: m.listHeaderStyle.Render("Gerenciamento de Lições")
}
func (m Model) ViewFooter() string {
	return m.help.View(m.keyMap)
}
func (m Model) ViewForm() string {
	// Placeholder para a altura ocupada pelo formulário
	// Isso é complexo de calcular precisamente sem renderizar.
	// Uma aproximação: número de linhas dos inputs + padding.
	if m.currentView == CreatingView || m.currentView == EditingView {
		return "form placeholder for height calculation" // Altura aproximada de 5-10 linhas
	}
	return ""
}

func (m Model) View() string {
	var s strings.Builder

	// Cabeçalho Global (Opcional, pode ser parte do app.Model)
	// s.WriteString(m.listHeaderStyle.Render("Gerenciamento de Lições") + "\n\n")

	if m.isLoading {
		s.WriteString("Carregando lições...")
		return s.String()
	}

	// Exibir mensagem de erro global, se houver e não estiver em ErrorView dedicada
	// Se ErrorView for usada, essa lógica pode ser diferente.
	if m.errorMessage != "" && m.currentView != ErrorView {
		s.WriteString(m.errorStyle.Render("Erro: "+m.errorMessage) + "\n")
	}
	if m.successMessage != "" {
		s.WriteString(m.successStyle.Render(m.successMessage) + "\n")
	}

	switch m.currentView {
	case ListView:
		s.WriteString(m.viewListView())
	case CreatingView:
		s.WriteString(m.viewFormView("Nova Lição"))
	case EditingView:
		if m.selectedLesson != nil {
			s.WriteString(m.viewFormView(fmt.Sprintf("Editando Lição ID: %d", m.selectedLesson.ID)))
		} else {
			// This case should ideally not be reached if logic is correct,
			// but as a safeguard:
			s.WriteString(m.errorStyle.Render("Erro: Nenhuma lição selecionada para edição."))
			// Consider also setting m.currentView = ListView and m.errorMessage
		}
	case ConfirmDeleteView:
		s.WriteString(m.viewConfirmDeleteView())
	case ErrorView: // Se ErrorView for um estado principal
		s.WriteString(m.errorStyle.Render(fmt.Sprintf("Ocorreu um erro: %s", m.errorMessage)))
		s.WriteString("\nPressione 'esc' para voltar.")
	default:
		s.WriteString("Visualização desconhecida.")
	}

	// Rodapé com ajuda
	s.WriteString("\n" + m.helpStyle.Render(m.help.View(m.keyMap)))

	return s.String()
}

func (m Model) viewListView() string {
	return m.table.View()
}

func (m Model) viewFormView(title string) string {
	var form strings.Builder
	form.WriteString(lipgloss.NewStyle().Bold(true).Render(title) + "\n\n")

	inputs := []string{
		fmt.Sprintf("ID da Turma:\n%s", m.classIDInput.View()),
		fmt.Sprintf("Título:\n%s", m.titleInput.View()),
		fmt.Sprintf("Plano de Aula (Markdown):\n%s", m.planContentInput.View()),
		fmt.Sprintf("Data Agendada (YYYY-MM-DD HH:MM):\n%s", m.scheduledAtInput.View()),
	}

	for i, inputView := range inputs {
		style := lipgloss.NewStyle().PaddingBottom(1)
		if i == m.focusedIndex { // Destacar o campo focado (opcional)
			// style = style.Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("205"))
		}
		form.WriteString(style.Render(inputView))
		// Adicionar mensagens de erro específicas do campo aqui, se necessário
	}

	// Adicionar dica de submissão
	form.WriteString("\n" + m.helpStyle.Render("Pressione Tab/Shift+Tab para navegar, Enter para submeter (ou nova linha no plano). Esc para cancelar."))

	// Envolver o formulário em uma caixa ou aplicar margem
	// return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1,2).Render(form.String())
	return form.String()
}

func (m Model) viewConfirmDeleteView() string {
	if m.selectedLesson == nil {
		return m.errorStyle.Render("Erro: Nenhuma lição selecionada para exclusão.")
	}
	title := fmt.Sprintf("Confirmar Exclusão da Lição: '%s' (ID: %d)?", m.selectedLesson.Title, m.selectedLesson.ID)
	return lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Render(title),
		"\nEsta ação não pode ser desfeita.",
		m.helpStyle.Render("[s/enter] para confirmar, [n/esc] para cancelar"),
	)
}

// KeyMap define os atalhos de teclado para o módulo de lições.
// Renomeado de keyMap para KeyMap para seguir convenção de exportação se necessário,
// embora aqui seja usado internamente.
type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	GotoTop     key.Binding
	GotoBottom  key.Binding
	Filter      key.Binding // Para filtrar a lista
	ClearFilter key.Binding // Para limpar o filtro

	Back    key.Binding // Esc: sair da view atual (form, delete confirm) ou do módulo
	New     key.Binding // n: abrir formulário de nova lição
	Edit    key.Binding // e: abrir formulário de edição da lição selecionada
	Delete  key.Binding // d: abrir confirmação de exclusão
	Select  key.Binding // Enter: selecionar item na lista (pode ser para editar ou ver detalhes)
	Refresh key.Binding // r: recarregar lista de lições

	// Form specific
	NextField key.Binding // Tab
	PrevField key.Binding // Shift+Tab
	Submit    key.Binding // Enter no form (ou Ctrl+s)
	Cancel    key.Binding // Esc no form (sinônimo de Back)

	// Universal
	Help key.Binding // ?
	Quit key.Binding // q ou Ctrl+c
}

// ShortHelp retorna a ajuda curta para o estado atual.
func (k KeyMap) ShortHelp() []key.Binding {
	switch currentView { // Precisa de acesso ao m.currentView ou passar como arg
	case ListView:
		return []key.Binding{k.New, k.Edit, k.Delete, k.Refresh, k.Help, k.Back}
	case CreatingView, EditingView:
		return []key.Binding{k.Submit, k.Cancel, k.NextField, k.PrevField, k.Help}
	case ConfirmDeleteView:
		return []key.Binding{k.Submit, k.Cancel, k.Help}
	default:
		return []key.Binding{k.Help, k.Quit}
	}
}

// FullHelp retorna a ajuda completa para o estado atual.
func (k KeyMap) FullHelp() [][]key.Binding {
	switch currentView { // Precisa de acesso ao m.currentView ou passar como arg
	case ListView:
		return [][]key.Binding{
			{k.Up, k.Down, k.PageUp, k.PageDown},
			{k.New, k.Edit, k.Delete, k.Select},
			{k.Filter, k.ClearFilter, k.Refresh},
			{k.Help, k.Back},
		}
	case CreatingView, EditingView:
		return [][]key.Binding{
			{k.Submit, k.Cancel},
			{k.NextField, k.PrevField},
			{k.Help},
		}
	case ConfirmDeleteView:
		return [][]key.Binding{
			{k.Submit, k.Cancel},
			{k.Help},
		}
	default:
		return [][]key.Binding{{k.Help, k.Quit}}
	}
}

// DefaultKeyMap retorna os atalhos de teclado padrão.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:          key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "mover para cima")),
		Down:        key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "mover para baixo")),
		PageUp:      key.NewBinding(key.WithKeys("pgup"), key.WithHelp("pgup", "página para cima")),
		PageDown:    key.NewBinding(key.WithKeys("pgdown"), key.WithHelp("pgdn", "página para baixo")),
		GotoTop:     key.NewBinding(key.WithKeys("home"), key.WithHelp("home", "ir para o topo")),
		GotoBottom:  key.NewBinding(key.WithKeys("end"), key.WithHelp("end", "ir para o fim")),
		Filter:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filtrar")),
		ClearFilter: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "limpar filtro")), // Sobrescreve 'esc' em modo de filtro

		Back:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "voltar/cancelar")),
		New:     key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "nova")),
		Edit:    key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "editar")),
		Delete:  key.NewBinding(key.WithKeys("d", "delete"), key.WithHelp("d/del", "excluir")),
		Select:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "selecionar/confirmar")),
		Refresh: key.NewBinding(key.WithKeys("r", "ctrl+r"), key.WithHelp("r", "recarregar")),

		NextField: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "próximo campo")),
		PrevField: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "campo anterior")),
		Submit:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submeter formulário")), // No form, enter é submit
		Cancel:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancelar formulário")),     // No form, esc é cancelar

		Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "ajuda")),
		Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "sair")),
	}
}

// Variáveis globais para ShortHelp e FullHelp, pois precisam de m.currentView
// Isso é um pouco problemático, idealmente KeyMap.ShortHelp/FullHelp receberia o ViewState.
// Por enquanto, vamos deixar assim e refatorar se necessário.
var currentView ViewState

// UpdateCurrentViewForHelp atualiza a view atual para que a ajuda seja exibida corretamente.
// Esta é uma solução temporária para o problema de acoplamento.
func UpdateCurrentViewForHelp(vs ViewState) {
	currentView = vs
}

// TODO: Definir mensagens específicas para lições (lessonsLoadedMsg, lessonCreatedMsg, etc.)
type lessonsLoadedMsg struct {
	lessons []models.Lesson
}

type lessonCreatedMsg struct {
	lesson models.Lesson
}

type lessonUpdatedMsg struct {
	lesson models.Lesson
}

type lessonDeletedMsg struct{}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }
