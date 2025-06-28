package classes

import (
	"context" // Adicionado para chamadas de serviço
	"fmt"
	"strconv" // Adicionado para conversão de subjectID
	"strings" // Adicionado para strings.Builder

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"vigenda/internal/models"
	"vigenda/internal/service"
)

type ViewState int

const (
	ListView ViewState = iota
	CreatingView
	// DetailsView // Para o futuro
)

var (
	columnTitleID        = "ID"
	columnTitleName      = "Nome da Turma"
	columnTitleSubjectID = "ID Disciplina"
)

// Model representa o modelo para a gestão de turmas.
type Model struct {
	classService service.ClassService
	state        ViewState
	table        table.Model
	createForm   struct {
		nameInput      textinput.Model
		subjectIDInput textinput.Model
		focusIndex     int
	}
	allClasses []models.Class // Para armazenar as turmas carregadas
	isLoading    bool
	width        int
	height       int
	err          error
}

// New cria um novo modelo para a gestão de turmas.
func New(cs service.ClassService) Model {
	// Tabela para listar turmas
	t := table.New(
		table.WithColumns([]table.Column{
			{Title: columnTitleID, Width: 5},
			{Title: columnTitleName, Width: 30},
			{Title: columnTitleSubjectID, Width: 15},
		}),
		table.WithRows([]table.Row{}), // Inicialmente vazia
		table.WithFocused(true),
		table.WithHeight(10), // Altura será ajustada
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

	// Formulário de criação
	nameInput := textinput.New()
	nameInput.Placeholder = "Nome da Nova Turma"
	nameInput.Focus() // Foco inicial pode ser removido se o estado inicial não for CreatingView
	nameInput.CharLimit = 100
	nameInput.Width = 30

	subjectIDInput := textinput.New()
	subjectIDInput.Placeholder = "ID da Disciplina (ex: 1)"
	subjectIDInput.CharLimit = 10
	subjectIDInput.Width = 20

	return Model{
		classService: cs,
		state:        ListView, // Estado inicial é a lista
		table:        t,
		createForm: struct {
			nameInput      textinput.Model
			subjectIDInput textinput.Model
			focusIndex     int
		}{
			nameInput:      nameInput,
			subjectIDInput: subjectIDInput,
			focusIndex:     0, // Foco no primeiro campo quando o formulário abrir
		},
		isLoading: true, // Começa carregando
	}
}

// Init carrega os dados iniciais para a gestão de turmas.
func (m Model) Init() tea.Cmd {
	m.isLoading = true // Garante que isLoading seja true ao iniciar
	return m.fetchClassesCmd
}

// Update lida com mensagens e atualiza o modelo.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height) // Chama SetSize para recalcular layout

	case tea.KeyMsg:
		switch m.state {
		case ListView:
			switch {
			// case key.Matches(msg, key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "voltar"))):
			// 'q' ou 'esc' para voltar ao menu é tratado em app.go
			case key.Matches(msg, key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "nova turma"))):
				m.state = CreatingView
				m.createForm.focusIndex = 0
				m.createForm.nameInput.Focus()
				m.createForm.subjectIDInput.Blur()
				m.createForm.nameInput.SetValue("")
				m.createForm.subjectIDInput.SetValue("")
				m.err = nil // Limpa erros anteriores
				return m, textinput.Blink
			default:
				m.table, cmd = m.table.Update(msg)
				cmds = append(cmds, cmd)
			}
		case CreatingView:
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancelar"))):
				m.state = ListView
				m.err = nil
				m.createForm.nameInput.Blur()
				m.createForm.subjectIDInput.Blur()
				return m, nil // Não precisa de comando, apenas muda o estado
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "salvar"))):
				// Salvar apenas se o foco estiver no último campo ou se houver apenas um campo
				// (Neste caso, com dois campos, quando o foco está no segundo)
				if m.createForm.focusIndex == 1 {
					name := strings.TrimSpace(m.createForm.nameInput.Value())
					subjectIDStr := strings.TrimSpace(m.createForm.subjectIDInput.Value())

					if name == "" || subjectIDStr == "" {
						m.err = fmt.Errorf("nome da turma e ID da disciplina são obrigatórios")
						return m, nil
					}
					m.isLoading = true
					cmds = append(cmds, m.createClassCmd(name, subjectIDStr))
				} else {
					m.nextFormInput() // Avança para o próximo campo
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "próximo"))),
				key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "anterior"))):
				if msg.String() == "shift+tab" {
					m.prevFormInput()
				} else {
					m.nextFormInput()
				}
			}
			// Atualiza o campo focado
			var focusedInputCmd tea.Cmd
			if m.createForm.focusIndex == 0 {
				m.createForm.nameInput, focusedInputCmd = m.createForm.nameInput.Update(msg)
			} else {
				m.createForm.subjectIDInput, focusedInputCmd = m.createForm.subjectIDInput.Update(msg)
			}
			cmds = append(cmds, focusedInputCmd)
		}

	case fetchedClassesMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
			m.allClasses = nil
			m.table.SetRows([]table.Row{}) // Limpa a tabela em caso de erro
		} else {
			m.err = nil
			m.allClasses = msg.classes
			var rows []table.Row
			for _, cls := range m.allClasses {
				rows = append(rows, table.Row{
					fmt.Sprintf("%d", cls.ID),
					cls.Name,
					fmt.Sprintf("%d", cls.SubjectID),
				})
			}
			m.table.SetRows(rows)
		}

	case classCreatedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = fmt.Errorf("erro ao criar turma: %w", msg.err)
			// Mantém o estado CreatingView para o usuário corrigir
		} else {
			m.state = ListView
			m.err = nil
			m.isLoading = true // Ativa isLoading para o fetchClassesCmd
			cmds = append(cmds, m.fetchClassesCmd)
		}

	case errMsg:
		m.err = msg.err
		m.isLoading = false
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) nextFormInput() {
	m.createForm.focusIndex = (m.createForm.focusIndex + 1) % 2 // 2 campos
	if m.createForm.focusIndex == 0 {
		m.createForm.nameInput.Focus()
		m.createForm.subjectIDInput.Blur()
	} else {
		m.createForm.nameInput.Blur()
		m.createForm.subjectIDInput.Focus()
	}
}

func (m *Model) prevFormInput() {
	m.createForm.focusIndex = (m.createForm.focusIndex - 1 + 2) % 2 // 2 campos
	if m.createForm.focusIndex == 0 {
		m.createForm.nameInput.Focus()
		m.createForm.subjectIDInput.Blur()
	} else {
		m.createForm.nameInput.Blur()
		m.createForm.subjectIDInput.Focus()
	}
}

// View renderiza a UI para a gestão de turmas.
func (m Model) View() string {
	var b strings.Builder

	if m.isLoading {
		// Centralizar a mensagem de carregamento (aproximadamente)
		loadingView := lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			"Carregando...",
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}))
		return loadingView
	}

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).PaddingBottom(1)
		b.WriteString(errorStyle.Render(fmt.Sprintf("Erro: %v", m.err)))
	}

	switch m.state {
	case ListView:
		titleStyle := lipgloss.NewStyle().Bold(true).MarginBottom(1)
		b.WriteString(titleStyle.Render("Lista de Turmas"))
		b.WriteString(m.table.View())
		helpStyle := lipgloss.NewStyle().Faint(true).MarginTop(1)
		b.WriteString(helpStyle.Render("Pressione 'n' para Nova Turma, ↑/↓ para navegar, 'q' ou 'esc' para voltar ao menu."))
	case CreatingView:
		titleStyle := lipgloss.NewStyle().Bold(true).MarginBottom(1)
		b.WriteString(titleStyle.Render("Nova Turma"))

		b.WriteString(m.createForm.nameInput.View())
		b.WriteString("\n")
		b.WriteString(m.createForm.subjectIDInput.View())
		b.WriteString("\n\n")

		helpStyle := lipgloss.NewStyle().Faint(true)
		b.WriteString(helpStyle.Render("Pressione Tab para navegar, Enter para salvar (no último campo), Esc para cancelar."))
	}
	return b.String()
}

// SetSize ajusta o tamanho do modelo e seus componentes.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Ajustes gerais de layout baseados no estado
	titleHeight := 1 // Para o título principal da view (ex: "Lista de Turmas")
	errorHeight := 0
	if m.err != nil {
		errorHeight = strings.Count(fmt.Sprintf("Erro: %v", m.err), "\n") + 1 + 1 // +1 for padding
	}
	helpHeight := 1 // Para a linha de ajuda

	remainingHeight := height - titleHeight - errorHeight - helpHeight

	if m.state == ListView {
		// Para a tabela, precisamos subtrair a altura do cabeçalho da tabela
		// A altura da tabela é o corpo + cabeçalho. table.Height() é só o corpo.
		// table.View() renderiza tudo.
		// Vamos dar um espaço fixo para o cabeçalho e bordas, e o resto para as linhas.
		tableHeaderAndBorderHeight := 3 // Estimativa
		tableBodyHeight := remainingHeight - tableHeaderAndBorderHeight
		if tableBodyHeight < 1 {
			tableBodyHeight = 1
		}
		m.table.SetHeight(tableBodyHeight)
		m.table.SetWidth(width - 4) // Margens laterais
	}

	if m.state == CreatingView {
		// Para o formulário, distribuir espaço para inputs
		// Altura dos inputs é geralmente 1 linha (ou 3 com bordas)
		// Largura dos inputs
		inputWidth := width - 4 // Margens
		if inputWidth < 20 {
			inputWidth = 20
		}
		m.createForm.nameInput.Width = inputWidth
		m.createForm.subjectIDInput.Width = inputWidth
	}
}

// IsFocused indica se o modelo de turmas tem algum input focado.
func (m Model) IsFocused() bool {
	return m.state == CreatingView // Se estiver no formulário, considera focado para 'esc' local
}

// Comandos e Mensagens

type fetchedClassesMsg struct {
	classes []models.Class
	err     error
}

type classCreatedMsg struct {
	createdClass models.Class // Pode ser útil para feedback, embora não usado diretamente agora
	err          error
}

type errMsg struct{ err error }

// Error torna errMsg em um tipo de erro válido.
func (e errMsg) Error() string { return e.err.Error() }

func (m Model) fetchClassesCmd() tea.Msg {
	// context.Background() é geralmente ok para operações TUI que não são canceláveis pelo usuário
	// de forma granular, mas se houver operações longas, um contexto com timeout/cancel pode ser melhor.
	classes, err := m.classService.ListAllClasses(context.Background())
	if err != nil {
		// Retorna um erro que pode ser tratado no Update
		return errMsg{fmt.Errorf("falha ao buscar turmas: %w", err)}
	}
	return fetchedClassesMsg{classes: classes, err: nil}
}

func (m Model) createClassCmd(name string, subjectIDStr string) tea.Cmd {
	return func() tea.Msg {
		subjectID, convErr := strconv.ParseInt(subjectIDStr, 10, 64)
		if convErr != nil {
			return errMsg{fmt.Errorf("ID da disciplina inválido ('%s'): %w", subjectIDStr, convErr)}
		}

		// Validação adicional, embora já feita antes de chamar o comando
		if name == "" {
			return errMsg{fmt.Errorf("nome da turma não pode ser vazio")}
		}
		// UserID é tratado pelo serviço/repositório
		createdClass, err := m.classService.CreateClass(context.Background(), name, subjectID)
		if err != nil {
			// Retorna um erro específico para criação, ou pode ser errMsg também
			return classCreatedMsg{err: fmt.Errorf("serviço falhou ao criar turma: %w", err)}
		}
		return classCreatedMsg{createdClass: createdClass, err: nil}
	}
}
