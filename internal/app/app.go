// Package app contém a lógica principal da Interface de Texto do Usuário (TUI)
// da aplicação Vigenda, utilizando o framework BubbleTea.
// Este arquivo (app.go) define o modelo principal da aplicação TUI,
// que gerencia as diferentes visualizações (telas/módulos) e suas interações.
package app

import (
	"fmt"
	"log" // Para logging interno do ciclo de vida da TUI.

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// Importações dos submódulos de visualização da TUI.
	"vigenda/internal/app/assessments"
	"vigenda/internal/app/classes"
	"vigenda/internal/app/dashboard"
	"vigenda/internal/app/proofs"
	"vigenda/internal/app/questions"
	"vigenda/internal/app/tasks"
	"vigenda/internal/service" // Importa as interfaces de serviço.
)

var (
	// appStyle define um estilo base para o contêiner principal da aplicação TUI.
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

// Model (ou AppModel) é o modelo raiz da aplicação TUI.
// Ele gerencia o estado global da TUI, como a visualização atual,
// e contém instâncias dos sub-modelos para cada funcionalidade principal.
// Também armazena as dependências de serviço injetadas.
type Model struct {
	list list.Model // list é usado para o menu principal quando currentView é DashboardView.

	currentView View // currentView rastreia qual módulo/tela está ativo.

	// Sub-modelos para cada funcionalidade principal da TUI.
	// Cada sub-modelo é um programa BubbleTea independente em sua essência.
	tasksModel       *tasks.Model
	classesModel     *classes.Model
	assessmentsModel *assessments.Model
	questionsModel   *questions.Model
	proofsModel      *proofs.Model
	dashboardModel   *dashboard.Model // Modelo para o painel de controle.

	width    int  // width da janela do terminal.
	height   int  // height da janela do terminal.
	quitting bool // quitting é true se a aplicação está em processo de encerramento.
	err      error // err armazena erros críticos que podem precisar ser exibidos.

	// Instâncias de serviço injetadas, usadas pelos sub-modelos.
	taskService       service.TaskService
	classService      service.ClassService
	assessmentService service.AssessmentService
	questionService   service.QuestionService
	proofService      service.ProofService
	lessonService     service.LessonService
}

// Init é o método de inicialização para o Model principal da aplicação.
// Conforme a filosofia BubbleTea, pode retornar um tea.Cmd para executar
// tarefas iniciais (ex: carregar dados). Neste caso, como os sub-modelos
// têm seus próprios Inits, o Init do AppModel principal pode não precisar
// fazer muito inicialmente, exceto se houver um estado global a ser carregado.
// Atualmente, retorna nil, indicando nenhuma ação inicial imediata neste nível.
func (m *Model) Init() tea.Cmd {
	// O estado inicial (DashboardView com a lista de menu) é configurado em New.
	// Os Inits dos sub-modelos são chamados quando a visualização muda para eles.
	return nil
}

// New é a função construtora para o Model principal da aplicação TUI.
// Recebe todas as dependências de serviço necessárias para a operação dos
// seus sub-modelos. Configura o menu principal (lista de itens) e inicializa
// todos os sub-modelos. Retorna um ponteiro para o Model configurado.
func New(
	ts service.TaskService, cs service.ClassService,
	as service.AssessmentService, qs service.QuestionService,
	ps service.ProofService, ls service.LessonService,
) *Model {
	// Define os itens do menu principal. Cada item tem um título e uma View associada.
	menuItems := []list.Item{
		menuItem{title: ConcreteDashboardView.String(), view: ConcreteDashboardView},
		menuItem{title: TaskManagementView.String(), view: TaskManagementView},
		menuItem{title: ClassManagementView.String(), view: ClassManagementView},
		menuItem{title: AssessmentManagementView.String(), view: AssessmentManagementView},
		menuItem{title: QuestionBankView.String(), view: QuestionBankView},
		menuItem{title: ProofGenerationView.String(), view: ProofGenerationView},
	}

	// Cria o componente de lista para o menu principal.
	l := list.New(menuItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Vigenda - Menu Principal"
	l.SetShowStatusBar(false)       // Não usa a barra de status padrão da lista.
	l.SetFilteringEnabled(false)    // Desabilita filtragem para o menu principal.
	l.Styles.Title = lipgloss.NewStyle().Bold(true).MarginBottom(1) // Estiliza o título do menu.
	// Define teclas de ajuda adicionais para a lista do menu.
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "sair")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "selecionar")),
		}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys // Mantém simples por enquanto.

	// Inicializa todos os sub-modelos, injetando suas respectivas dependências de serviço.
	tm := tasks.New(ts)
	cm := classes.New(cs)
	am := assessments.New(as, cs) // Passa ClassService aqui
	qm := questions.New(qs)
	pm := proofs.New(ps)
	dshModel := dashboard.New(ts, cs, as, ls)

	// Retorna a instância do Model principal.
	return &Model{
		list:              l,
		currentView:       DashboardView, // A visualização inicial é o menu principal (DashboardView).
		tasksModel:        tm,
		taskService:       ts,
		classesModel:      cm,
		classService:      cs,
		assessmentsModel:  am,
		assessmentService: as,
		questionsModel:    qm,
		questionService:   qs,
		proofsModel:       pm,
		proofService:      ps,
		lessonService:     ls,
		dashboardModel:    dshModel,
	}
}

// menuItem é uma struct helper que implementa a interface list.Item.
// Usada para popular a lista do menu principal.
type menuItem struct {
	title string // title é o texto exibido para o item de menu.
	view  View   // view é a constante View associada a este item de menu.
}

// FilterValue retorna o valor usado pela lista para filtrar itens (o título neste caso).
func (i menuItem) FilterValue() string { return i.title }

// Title retorna o título do item de menu.
func (i menuItem) Title() string { return i.title }

// Description retorna a descrição do item de menu (vazio neste caso).
func (i menuItem) Description() string { return "" }

// Update é o coração do ciclo de vida do BubbleTea para o Model principal.
// Ele processa mensagens (eventos) como entradas de teclado, redimensionamento de janela,
// ou mensagens customizadas de comandos.
// Retorna o modelo atualizado e um tea.Cmd para quaisquer operações subsequentes.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("AppModel: Update GLOBAL - Recebida msg tipo %T Valor: %v", msg, msg)
	var cmds []tea.Cmd // Slice para acumular comandos a serem executados.

	// Primeiro switch: lida com mensagens globais ou específicas do AppModel.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // Janela do terminal foi redimensionada.
		m.width = msg.Width
		m.height = msg.Height
		// Ajusta o tamanho da lista do menu principal.
		listHeight := msg.Height - appStyle.GetVerticalPadding() - lipgloss.Height(m.list.Title) - lipgloss.Height(m.list.Help.View(m.list)) - 2
		m.list.SetSize(msg.Width-appStyle.GetHorizontalPadding(), listHeight)

		// Propaga a mensagem de redimensionamento para todos os sub-modelos ativos.
		var subCmd tea.Cmd
		var tempModel tea.Model

		tempModel, subCmd = m.dashboardModel.Update(msg)
		m.dashboardModel = tempModel.(*dashboard.Model)
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.tasksModel.Update(msg)
		m.tasksModel = tempModel.(*tasks.Model)
		cmds = append(cmds, subCmd)
		// ... (repetir para todos os outros sub-modelos) ...
		tempModel, subCmd = m.classesModel.Update(msg)
		m.classesModel = tempModel.(*classes.Model)
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.assessmentsModel.Update(msg)
		m.assessmentsModel = tempModel.(*assessments.Model)
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.questionsModel.Update(msg)
		m.questionsModel = tempModel.(*questions.Model)
		cmds = append(cmds, subCmd)

		tempModel, subCmd = m.proofsModel.Update(msg)
		m.proofsModel = tempModel.(*proofs.Model)
		cmds = append(cmds, subCmd)

		return m, tea.Batch(cmds...)

	case tea.KeyMsg: // Mensagem de tecla pressionada.
		// Atalho global para sair (Ctrl+C).
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))) {
			m.quitting = true
			cmds = append(cmds, tea.Quit)
			return m, tea.Batch(cmds...)
		}

		// Se a visualização atual é o menu principal (DashboardView).
		if m.currentView == DashboardView {
			var listCmd tea.Cmd
			m.list, listCmd = m.list.Update(msg) // Deixa o componente de lista lidar com a navegação.
			cmds = append(cmds, listCmd)

			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) { // Seleção de item de menu.
				selectedItem, ok := m.list.SelectedItem().(menuItem)
				if ok {
					m.currentView = selectedItem.view // Muda para a visualização selecionada.
					log.Printf("AppModel: Mudando para view %s (%d)", m.currentView.String(), m.currentView)
					// Dispara o comando Init do sub-modelo correspondente.
					switch m.currentView {
					case ConcreteDashboardView:
						cmds = append(cmds, m.dashboardModel.Init())
					case TaskManagementView:
						cmds = append(cmds, m.tasksModel.Init())
					case ClassManagementView:
						cmds = append(cmds, m.classesModel.Init())
					case AssessmentManagementView:
						cmds = append(cmds, m.assessmentsModel.Init())
					case QuestionBankView:
						cmds = append(cmds, m.questionsModel.Init())
					case ProofGenerationView:
						cmds = append(cmds, m.proofsModel.Init())
					}
				}
			} else if key.Matches(msg, key.NewBinding(key.WithKeys("q"))) { // Sair do menu principal.
				m.quitting = true
				cmds = append(cmds, tea.Quit)
			}
			return m, tea.Batch(cmds...)
		}
		// Se não estiver no DashboardView, a tecla será passada para o sub-modelo ativo abaixo.

	case error: // Captura erros globais (ex: de Inits de sub-modelos).
		m.err = msg
		log.Printf("AppModel: Erro global recebido: %v", msg)
		return m, tea.Batch(cmds...) // Armazena o erro para exibição.
	}

	// Segundo estágio: Delega a mensagem para o sub-modelo ativo se não foi tratada globalmente.
	var submodelCmd tea.Cmd
	var updatedSubModel tea.Model // Usar tea.Model para o tipo retornado por Update.

	switch m.currentView {
	case ConcreteDashboardView:
		updatedSubModel, submodelCmd = m.dashboardModel.Update(msg)
		m.dashboardModel = updatedSubModel.(*dashboard.Model)
		// A lógica de 'esc' para o dashboard é um caso especial, pois ele pode ter foco interno.
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if !m.dashboardModel.IsFocused() { // Se o dashboard não tem mais foco interno, volta ao menu.
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para o Menu Principal a partir do Painel de Controle.")
			}
		}
	case TaskManagementView:
		updatedSubModel, submodelCmd = m.tasksModel.Update(msg)
		m.tasksModel = updatedSubModel.(*tasks.Model)
		// Se o sub-modelo de tarefas sinalizar que deve voltar (ex: após 'esc' no nível raiz),
		// ele deve resetar seu próprio estado e o AppModel o trará de volta ao menu.
		// Vamos assumir que o sub-modelo gerencia seu estado e nós apenas trocamos a view.
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if m.tasksModel.CanGoBack() { // Supondo que o modelo de tarefas tenha um método CanGoBack.
				m.currentView = DashboardView
			}
		}
	case ClassManagementView:
		updatedSubModel, submodelCmd = m.classesModel.Update(msg)
		m.classesModel = updatedSubModel.(*classes.Model)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if m.classesModel.CanGoBack() { // Supondo que o modelo de turmas tenha um método CanGoBack.
				m.currentView = DashboardView
			}
		}
	case AssessmentManagementView:
		updatedSubModel, submodelCmd = m.assessmentsModel.Update(msg)
		m.assessmentsModel = updatedSubModel.(*assessments.Model)
		// O modelo de avaliações tem sua própria navegação interna com 'esc'.
		// Só voltamos ao menu se o 'esc' for pressionado na visualização principal (ListView) do sub-modelo.
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if m.assessmentsModel.IsAtRoot() { // Supondo que o modelo de avaliações tenha um método IsAtRoot.
				m.currentView = DashboardView
				log.Println("AppModel: Voltando para o Menu Principal a partir de Gerenciar Avaliações.")
			}
		}
	case QuestionBankView:
		updatedSubModel, submodelCmd = m.questionsModel.Update(msg)
		m.questionsModel = updatedSubModel.(*questions.Model)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if m.questionsModel.CanGoBack() { // Supondo que o modelo de questões tenha um método CanGoBack.
				m.currentView = DashboardView
			}
		}
	case ProofGenerationView:
		updatedSubModel, submodelCmd = m.proofsModel.Update(msg)
		m.proofsModel = updatedSubModel.(*proofs.Model)
		if km, ok := msg.(tea.KeyMsg); ok && key.Matches(km, key.NewBinding(key.WithKeys("esc"))) {
			if m.proofsModel.CanGoBack() { // Supondo que o modelo de provas tenha um método CanGoBack.
				m.currentView = DashboardView
			}
		}
	}
	cmds = append(cmds, submodelCmd) // Adiciona comando do sub-modelo.

	return m, tea.Batch(cmds...)
}

// View renderiza a interface do usuário com base no estado atual do Model principal.
// Se estiver saindo ou houver um erro crítico, exibe mensagens apropriadas.
// Caso contrário, delega a renderização para o View() do sub-modelo ativo
// ou exibe o menu principal.
func (m *Model) View() string {
	if m.quitting {
		return appStyle.Render("Saindo do Vigenda...\n")
	}
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")) // Vermelho para erro.
		return appStyle.Render(fmt.Sprintf("Ocorreu um erro crítico: %v\nPressione Ctrl+C para sair.", errorStyle.Render(m.err.Error())))
	}

	var viewContent string
	var help string // Texto de ajuda, pode ser contextual.

	switch m.currentView {
	case DashboardView: // Se a view atual é o menu principal.
		viewContent = m.list.View()
		help = m.list.Help.View(m.list) // Usa a ajuda embutida do componente de lista.
	case ConcreteDashboardView:
		viewContent = m.dashboardModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal."
	case TaskManagementView:
		viewContent = m.tasksModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal."
	case ClassManagementView:
		viewContent = m.classesModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal."
	case AssessmentManagementView:
		viewContent = m.assessmentsModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal."
	case QuestionBankView:
		viewContent = m.questionsModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal."
	case ProofGenerationView:
		viewContent = m.proofsModel.View()
		help = "\nPressione 'esc' para voltar ao menu principal."
	default: // Caso uma view desconhecida seja definida.
		viewContent = fmt.Sprintf("Visão desconhecida: %s (%d)", m.currentView.String(), m.currentView)
		help = "\nPressione 'esc' ou 'q' para tentar voltar ao menu principal."
	}

	// Junta o conteúdo da view principal com o texto de ajuda.
	finalRender := lipgloss.JoinVertical(lipgloss.Left,
		viewContent,
		lipgloss.NewStyle().MarginTop(1).Render(help), // Adiciona margem para separar a ajuda.
	)
	return appStyle.Render(finalRender) // Aplica o estilo global da aplicação.
}

// StartApp é a função ponto de entrada para iniciar a aplicação TUI Vigenda.
// Ela cria uma nova instância do Model principal, injetando todas as dependências de serviço,
// e então inicia o programa BubbleTea.
// Esta função é tipicamente chamada pelo comando raiz da CLI quando nenhuma subcomando é fornecido.
func StartApp(
	ts service.TaskService, cs service.ClassService,
	as service.AssessmentService, qs service.QuestionService,
	ps service.ProofService, ls service.LessonService,
) {
	model := New(ts, cs, as, qs, ps, ls)
	// tea.WithAltScreen() usa o buffer alternativo do terminal, preservando o histórico do shell.
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		// Usa log.Fatalf para registrar o erro e sair, garantindo que o erro seja logado no arquivo.
		log.Fatalf("Erro ao executar o programa BubbleTea: %v", err)
	}
}
