// Package tui contém componentes e lógica relacionados à Interface de Texto do Usuário (TUI)
// da aplicação Vigenda. Este pacote específico (tui.go) define um modelo BubbleTea
// que parece ser um protótipo ou uma parte de uma interface para gerenciamento de turmas e alunos.
// A TUI principal e mais completa da aplicação é gerenciada pelo pacote `internal/app`.
// Os componentes definidos aqui podem ser genéricos ou legados.
package tui

import (
	"context"
	"fmt"
	"log"
	"vigenda/internal/app"
	"vigenda/internal/models"
	"vigenda/internal/service"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	docStyle        = lipgloss.NewStyle().Margin(1, 2)
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")) // Purple
	selectedStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("62")).PaddingLeft(1)
	deselectedStyle = lipgloss.NewStyle().PaddingLeft(1)
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray
)

// KeyMap define um conjunto de mapeamentos de teclas comuns usados na TUI.
// Cada campo representa uma ação e está associado a uma ou mais teclas.
type KeyMap struct {
	Up     key.Binding // Up move a seleção para cima.
	Down   key.Binding // Down move a seleção para baixo.
	Select key.Binding // Select confirma a seleção atual ou executa uma ação.
	Back   key.Binding // Back retorna à visualização anterior ou cancela uma ação.
	Quit   key.Binding // Quit encerra a aplicação TUI.
	Help   key.Binding // Help exibe/oculta informações de ajuda.
}

// DefaultKeyMap fornece um conjunto padrão de mapeamentos de teclas para a TUI.
var DefaultKeyMap = KeyMap{
	Up:     key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "mover para cima")),
	Down:   key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "mover para baixo")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "selecionar")),
	Back:   key.NewBinding(key.WithKeys("esc", "backspace"), key.WithHelp("esc/bksp", "voltar")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q/ctrl+c", "sair")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "ajuda")),
}

// Model representa o estado da interface TUI definida neste arquivo.
// Este modelo parece ser focado no gerenciamento de turmas (classes) e alunos (students),
// servindo como um exemplo ou uma parte de uma TUI maior.
// A TUI principal da aplicação é gerenciada por `internal/app.AppModel`.
type Model struct {
	classService service.ClassService // classService é a dependência para interagir com a lógica de negócios de turmas.

	list    list.Model    // list é o componente de lista usado para exibir turmas ou alunos.
	spinner spinner.Model // spinner é usado para indicar atividades de carregamento.
	keys    KeyMap        // keys contém os mapeamentos de teclas para este modelo.

	currentView app.View // currentView indica qual visualização (turmas ou alunos) está ativa.
	isLoading   bool       // isLoading é true quando dados estão sendo carregados assincronamente.
	err         error      // err armazena qualquer erro ocorrido durante operações.

	classes       []models.Class  // classes armazena a lista de turmas carregadas.
	students      []models.Student // students armazena a lista de alunos carregados para uma turma selecionada.
	selectedClass *models.Class // selectedClass armazena a turma atualmente selecionada ao visualizar seus alunos.
}

// NewTUIModel cria e inicializa uma nova instância do Model da TUI.
// Requer um `service.ClassService` para interagir com a lógica de negócios.
// Este construtor configura o spinner, define as teclas padrão e inicia o carregamento
// dos dados iniciais (lista de turmas).
func NewTUIModel(cs service.ClassService) Model {
	log.Printf("TUI(tui.go): NewTUIModel - Chamado. ClassService is nil: %t", cs == nil)
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := Model{
		classService: cs,
		spinner:      s,
		keys:         DefaultKeyMap,
		currentView:  app.DashboardView, // Start with Dashboard or main menu
		isLoading:    false,
	}
	log.Println("TUI: NewTUIModel - Modelo TUI parcialmente inicializado, chamando loadInitialData.")
	m.loadInitialData() // Carrega os dados iniciais (ex: lista de turmas).
	log.Println("TUI(tui.go): NewTUIModel - loadInitialData chamado, retornando modelo.")
	return m
}

// loadInitialData é chamado durante a inicialização do modelo para carregar
// o estado inicial da visualização, como a lista de turmas.
// Define isLoading como true e dispara o comando para carregar as turmas.
func (m *Model) loadInitialData() {
	m.isLoading = true
	// Este modelo TUI específico parece focado no gerenciamento de turmas.
	m.currentView = app.ClassManagementView // Define a visualização inicial.
	m.loadClasses()                         // Dispara o carregamento das turmas.
}

// loadClasses retorna um tea.Cmd que, quando executado, carrega a lista de todas as turmas
// do ClassService. Em caso de sucesso, envia uma classesLoadedMsg. Em caso de erro,
// envia uma errMsg.
func (m *Model) loadClasses() tea.Cmd {
	m.isLoading = true
	log.Println("TUI(tui.go): loadClasses - Iniciando carregamento de turmas.")
	return func() tea.Msg {
		log.Println("TUI(tui.go): loadClasses (cmd) - Tentando carregar turmas do serviço.")
		classes, err := m.classService.ListAllClasses(context.Background())
		if err != nil {
			log.Printf("TUI(tui.go): loadClasses (cmd) - Erro ao carregar turmas: %v", err)
			return errMsg{err: err, context: "carregando turmas"}
		}
		log.Printf("TUI(tui.go): loadClasses (cmd) - Turmas carregadas com sucesso: %d turmas.", len(classes))
		return classesLoadedMsg(classes)
	}
}

// loadStudentsForClass retorna um tea.Cmd para carregar os alunos de uma turma específica.
// classID é o ID da turma cujos alunos serão carregados.
// Envia studentsLoadedMsg em caso de sucesso ou errMsg em caso de erro.
func (m *Model) loadStudentsForClass(classID int64) tea.Cmd {
	m.isLoading = true
	log.Printf("TUI(tui.go): loadStudentsForClass - Iniciando carregamento de alunos para a turma ID %d.", classID)
	return func() tea.Msg {
		log.Printf("TUI(tui.go): loadStudentsForClass (cmd) - Tentando carregar alunos para a turma ID %d.", classID)
		students, err := m.classService.GetStudentsByClassID(context.Background(), classID)
		if err != nil {
			log.Printf("TUI(tui.go): loadStudentsForClass (cmd) - Erro ao carregar alunos para a turma ID %d: %v", classID, err)
			return errMsg{err: err, context: fmt.Sprintf("carregando alunos para turma %d", classID)}
		}
		log.Printf("TUI(tui.go): loadStudentsForClass (cmd) - Alunos carregados para turma ID %d: %d alunos.", classID, len(students))
		return studentsLoadedMsg(students)
	}
}

// Init é o comando inicial executado quando o programa BubbleTea inicia.
// Para este modelo, ele inicia o spinner (se estiver carregando) e dispara o
// carregamento inicial da lista de turmas.
func (m Model) Init() tea.Cmd {
	// Inicia o spinner e o comando para carregar as turmas.
	return tea.Batch(m.spinner.Tick, m.loadClasses())
}

// Update é a função principal que lida com todas as mensagens (eventos) recebidas pela TUI,
// como entradas de teclado, ticks de spinner, ou mensagens de conclusão de comandos assíncronos.
// Retorna o modelo atualizado e um comando a ser executado (pode ser nil).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd // Coleção de comandos a serem executados.

	switch msg := msg.(type) {
	case tea.KeyMsg: // Mensagem de tecla pressionada.
		switch {
		case key.Matches(msg, m.keys.Quit): // Se 'q' ou 'ctrl+c' for pressionado.
			return m, tea.Quit // Encerra o programa.
		case key.Matches(msg, m.keys.Up): // Tecla para cima.
			m.list.CursorUp() // Move o cursor da lista para cima.
		case key.Matches(msg, m.keys.Down): // Tecla para baixo.
			m.list.CursorDown() // Move o cursor da lista para baixo.
		case key.Matches(msg, m.keys.Select): // Tecla Enter.
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				switch m.currentView {
				case app.ClassManagementView: // Se na visualização de gerenciamento de turmas.
					if class, ok := selectedItem.(listItemClass); ok {
						m.selectedClass = &class.Class // Armazena a turma selecionada.
						m.currentView = app.StudentView // Muda para a visualização de alunos (valor de exemplo).
						cmds = append(cmds, m.loadStudentsForClass(class.ID())) // Carrega alunos desta turma.
					}
				case app.StudentView: // Se na visualização de alunos.
					// TODO: Implementar seleção de aluno ou ação.
				}
			}
		case key.Matches(msg, m.keys.Back): // Tecla Esc ou Backspace.
			if m.currentView == app.StudentView { // Se na visualização de alunos.
				m.currentView = app.ClassManagementView // Volta para a visualização de turmas.
				m.selectedClass = nil
				cmds = append(cmds, m.loadClasses()) // Recarrega a lista de turmas.
			} else if m.currentView == app.ProofGenerationView { // Exemplo de outra view.
				m.currentView = app.DashboardView // Volta para o dashboard.
			} else if m.currentView == app.ClassManagementView {
				// TODO: Implementar volta para menu anterior ou sair se for o nível mais alto.
			}
		}

	case spinner.TickMsg: // Mensagem de tick do spinner (animação).
		if m.isLoading { // Só atualiza o spinner se estiver carregando.
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case classesLoadedMsg: // Mensagem indicando que as turmas foram carregadas.
		log.Println("TUI(tui.go): Update - Recebida classesLoadedMsg.")
		m.isLoading = false
		m.classes = []models.Class(msg)
		items := make([]list.Item, len(m.classes))
		for i, c := range m.classes {
			items[i] = listItemClass{c} // Converte para o tipo de item da lista.
		}
		m.list = list.New(items, list.NewDefaultDelegate(), 0, 0) // Cria/atualiza a lista.
		m.list.Title = "Turmas"
		m.list.SetShowHelp(false) // Desabilita ajuda padrão da lista.
		log.Printf("TUI(tui.go): Update - Lista de turmas atualizada com %d itens.", len(items))

	case studentsLoadedMsg: // Mensagem indicando que os alunos foram carregados.
		log.Println("TUI(tui.go): Update - Recebida studentsLoadedMsg.")
		m.isLoading = false
		m.students = []models.Student(msg)
		items := make([]list.Item, len(m.students))
		for i, s := range m.students {
			items[i] = listItemStudent{s}
		}
		m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
		if m.selectedClass != nil {
			m.list.Title = fmt.Sprintf("Alunos da Turma: %s", m.selectedClass.Name)
		} else {
			m.list.Title = "Alunos"
		}
		m.list.SetShowHelp(false)
		log.Printf("TUI(tui.go): Update - Lista de alunos atualizada com %d itens.", len(items))

	case errMsg: // Mensagem de erro.
		log.Printf("TUI(tui.go): Update - Recebida errMsg. Contexto: '%s', Erro: %v", msg.context, msg.err)
		m.isLoading = false
		m.err = msg.err // Armazena o erro para exibição.
		// TODO: Implementar melhor exibição de erro na TUI em vez de quitar.
		return m, tea.Quit // Por enquanto, encerra em caso de erro.

	case tea.WindowSizeMsg: // Mensagem de redimensionamento da janela.
		h, v := docStyle.GetFrameSize()
		// Ajusta o tamanho da lista com base no tamanho da janela e outros elementos.
		m.list.SetSize(msg.Width-h, msg.Height-v-lipgloss.Height(m.headerView())-lipgloss.Height(m.footerView()))
	}

	// Se não estiver carregando, permite que o componente de lista processe a mensagem.
	if !m.isLoading {
		var listCmd tea.Cmd
		m.list, listCmd = m.list.Update(msg)
		cmds = append(cmds, listCmd)
	}

	return m, tea.Batch(cmds...) // Retorna o modelo atualizado e todos os comandos acumulados.
}

// View renderiza a interface do usuário com base no estado atual do modelo.
// Retorna uma string que representa a TUI a ser exibida no terminal.
func (m Model) View() string {
	if m.err != nil { // Se houver um erro, exibe a mensagem de erro.
		return docStyle.Render(fmt.Sprintf("Ocorreu um erro: %v\n\nPressione qualquer tecla para sair.", m.err))
	}

	if m.isLoading { // Se estiver carregando, exibe o spinner e uma mensagem.
		loadingText := fmt.Sprintf("%s Carregando...", m.spinner.View())
		// Adiciona mensagens contextuais de carregamento.
		if m.currentView == app.ClassManagementView && len(m.classes) == 0 {
			loadingText += "\n\nNenhuma turma encontrada ainda."
		} else if m.currentView == app.StudentView && len(m.students) == 0 && m.selectedClass != nil {
			loadingText += fmt.Sprintf("\n\nNenhum aluno encontrado para a turma %s.", m.selectedClass.Name)
		}
		return docStyle.Render(loadingText)
	}

	// Se não houver erro e não estiver carregando, renderiza a visualização principal.
	return docStyle.Render(m.headerView() + "\n" + m.list.View() + "\n" + m.footerView())
}

// headerView renderiza o cabeçalho da visualização atual.
func (m Model) headerView() string {
	var title string
	switch m.currentView {
	case app.ClassManagementView:
		title = titleStyle.Render("Gerenciar Turmas")
	case app.StudentView: // Visualização de alunos (valor de exemplo)
		if m.selectedClass != nil {
			title = titleStyle.Render(fmt.Sprintf("Alunos - %s", m.selectedClass.Name))
		} else {
			title = titleStyle.Render("Alunos")
		}
	default:
		title = titleStyle.Render(m.currentView.String()) // Usa o nome da view padrão.
	}
	return title
}

// footerView renderiza o rodapé com informações de ajuda.
func (m Model) footerView() string {
	return helpStyle.Render(fmt.Sprintf("↑/↓: Navegar | Enter: Selecionar | Esc: Voltar | q: Sair"))
}

// listItemClass é um adaptador para usar models.Class com o componente list.Model.
// Implementa a interface list.Item.
type listItemClass struct {
	models.Class // Incorpora models.Class para acesso direto aos campos.
}

// Title retorna o título do item da lista (nome da turma).
func (lic listItemClass) Title() string { return lic.Name }

// Description retorna a descrição do item da lista (ID da turma e ID da disciplina).
func (lic listItemClass) Description() string {
	return fmt.Sprintf("ID: %d, Disciplina ID: %d", lic.Class.ID, lic.SubjectID)
}

// FilterValue retorna o valor usado para filtrar a lista (nome da turma).
func (lic listItemClass) FilterValue() string { return lic.Name }

// ID retorna o ID da turma (usado internamente, não faz parte de list.Item mas útil).
func (lic listItemClass) ID() int64 { return lic.Class.ID }

// listItemStudent é um adaptador para usar models.Student com o componente list.Model.
// Implementa a interface list.Item.
type listItemStudent struct {
	models.Student // Incorpora models.Student.
}

// Title retorna o título do item da lista (nome completo do aluno).
func (lis listItemStudent) Title() string { return lis.FullName }

// Description retorna a descrição do item da lista (matrícula e status do aluno).
func (lis listItemStudent) Description() string {
	return fmt.Sprintf("Matrícula: %s, Status: %s", lis.EnrollmentID, lis.Status)
}

// FilterValue retorna o valor usado para filtrar a lista (nome completo do aluno).
func (lis listItemStudent) FilterValue() string { return lis.FullName }

// ID retorna o ID do aluno (usado internamente).
func (lis listItemStudent) ID() int64 { return lis.Student.ID }

// errMsg é uma mensagem para encapsular erros ocorridos durante operações assíncronas (tea.Cmd).
// Inclui um contexto para identificar a origem do erro.
type errMsg struct {
	err     error  // O erro original.
	context string // Contexto onde o erro ocorreu (ex: "carregando turmas").
}

// Error implementa a interface error.
func (e errMsg) Error() string {
	return fmt.Sprintf("contexto: %s, erro: %v", e.context, e.err)
}

// classesLoadedMsg é uma mensagem enviada quando a lista de turmas é carregada com sucesso.
// Contém a slice de turmas.
type classesLoadedMsg []models.Class

// studentsLoadedMsg é uma mensagem enviada quando a lista de alunos de uma turma é carregada.
// Contém a slice de alunos.
type studentsLoadedMsg []models.Student

// Start inicia e executa o programa TUI definido neste arquivo.
// Recebe um ClassService para operações de dados.
// Esta função é provavelmente um ponto de entrada para uma seção específica da TUI
// ou um exemplo, já que a TUI principal é iniciada por `app.StartApp`.
// Retorna um erro se o programa BubbleTea falhar ao executar.
func Start(classService service.ClassService) error {
	log.Printf("TUI(tui.go): Start - Função Start chamada. ClassService is nil: %t", classService == nil)
	if classService == nil {
		log.Fatalf("TUI(tui.go): Start - ClassService não pode ser nulo para iniciar este modelo TUI.")
	}
	m := NewTUIModel(classService)
	p := tea.NewProgram(m, tea.WithAltScreen()) // Usa AltScreen para uma melhor experiência TUI.
	log.Println("TUI(tui.go): Start - Iniciando programa Bubble Tea (p.Run()).")
	_, err := p.Run()
	if err != nil {
		log.Printf("TUI(tui.go): Start - Erro ao executar o programa Bubble Tea: %v", err)
	} else {
		log.Println("TUI(tui.go): Start - Programa Bubble Tea finalizado sem erros.")
	}
	return err
}
