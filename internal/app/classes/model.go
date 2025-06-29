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
	DetailsView
)

type FocusTarget int

const (
	FocusTargetNone FocusTarget = iota
	FocusTargetStudentsTable
	FocusTargetStatusSelector // Será usado na Parte 2
)

var (
	columnTitleID        = "ID"
	columnTitleName      = "Nome da Turma"
	columnTitleSubjectID = "ID Disciplina"

	// Colunas para a tabela de alunos
	studentColumnTitleID         = "ID Aluno"
	studentColumnTitleEnrollment = "Nº Chamada"
	studentColumnTitleFullName   = "Nome Completo"
	studentColumnTitleStatus     = "Status"
)

// Model representa o modelo para a gestão de turmas.
type Model struct {
	classService  service.ClassService
	state         ViewState
	table         table.Model // Tabela de turmas
	studentsTable table.Model // Tabela de alunos para DetailsView
	createForm    struct {
		nameInput      textinput.Model
		subjectIDInput textinput.Model
		focusIndex     int
	}
	allClasses    []models.Class // Para armazenar as turmas carregadas
	selectedClass *models.Class  // Turma selecionada para DetailsView
	classStudents []models.Student // Alunos da turma selecionada

	detailsViewFocusTarget FocusTarget // Novo campo

	isLoading bool
	width     int
	height    int
	err       error
}

// New cria um novo modelo para a gestão de turmas.
func New(cs service.ClassService) Model {
	// Tabela para listar turmas
	classTable := table.New(
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
	classTable.SetStyles(s)

	// Tabela para listar alunos (em DetailsView)
	studentsTable := table.New(
		table.WithColumns([]table.Column{
			{Title: studentColumnTitleID, Width: 8},
			{Title: studentColumnTitleEnrollment, Width: 10},
			{Title: studentColumnTitleFullName, Width: 30},
			{Title: studentColumnTitleStatus, Width: 10},
		}),
		table.WithRows([]table.Row{}),
		table.WithFocused(false), // Foco inicial na tabela de turmas ou no form
		table.WithHeight(10),
	)
	studentsTable.SetStyles(s) // Reutiliza o mesmo estilo básico

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
		classService:  cs,
		state:         ListView, // Estado inicial é a lista
		table:         classTable,
		studentsTable: studentsTable,
		createForm: struct {
			nameInput      textinput.Model
			subjectIDInput textinput.Model
			focusIndex     int
		}{
			nameInput:      nameInput,
			subjectIDInput: subjectIDInput,
			focusIndex:     0, // Foco no primeiro campo quando o formulário abrir
		},
		detailsViewFocusTarget: FocusTargetNone, // Explícito para clareza
		isLoading:              true,            // Começa carregando
	}
}

// Init carrega os dados iniciais para a gestão de turmas.
func (m Model) Init() tea.Cmd {
	// m.isLoading = true // isLoading já é true por padrão em New e reafirmado aqui
	return m.fetchClassesCmd
}


// Update lida com mensagens e atualiza o modelo.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		// Se estivermos em DetailsView e studentsTable estiver focada, ela deve lidar com as teclas primeiro.
		// (Isso será relevante quando adicionarmos foco à studentsTable)
		// if m.state == DetailsView && m.studentsTable.Focused() {
		// 	var studentsTableCmd tea.Cmd
		// 	m.studentsTable, studentsTableCmd = m.studentsTable.Update(msg)
		// 	cmds = append(cmds, studentsTableCmd)
		// 	return m, tea.Batch(cmds...)
		// }


		switch m.state {
		case ListView:
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "detalhes"))):
				if len(m.allClasses) > 0 && m.table.Cursor() < len(m.allClasses) {
					selected := m.allClasses[m.table.Cursor()]
					m.selectedClass = &selected // Armazena um ponteiro para a turma selecionada
					m.state = DetailsView
					m.isLoading = true // Mostrar carregamento para os alunos
					m.err = nil        // Limpar erros anteriores
					// Limpar alunos e tabela de alunos anteriores
					m.classStudents = nil
					m.studentsTable.SetRows([]table.Row{})
					// Disparar comando para buscar alunos da turma selecionada
					if m.selectedClass != nil { // Garantir que selectedClass não é nil
						cmds = append(cmds, m.fetchClassStudentsCmd(m.selectedClass.ID))
					} else {
						// Isso não deveria acontecer se a lógica de seleção estiver correta
						cmds = append(cmds, func() tea.Msg { return errMsg{fmt.Errorf("turma selecionada é nil antes de buscar alunos")} })
					}
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "nova turma"))):
				m.state = CreatingView
				m.createForm.focusIndex = 0
				m.createForm.nameInput.Focus()
				m.createForm.subjectIDInput.Blur()
				m.createForm.nameInput.SetValue("")
				m.createForm.subjectIDInput.SetValue("")
				m.err = nil
				m.selectedClass = nil // Limpa qualquer seleção anterior
				return m, textinput.Blink
			default:
				m.table, cmd = m.table.Update(msg)
				cmds = append(cmds, cmd)
			}
		// ... outros cases para CreatingView, DetailsView (para 'esc')...
		case CreatingView:
			// ... (lógica existente) ...
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancelar"))):
				m.state = ListView
				m.err = nil
				m.createForm.nameInput.Blur()
				m.createForm.subjectIDInput.Blur()
				m.table.Focus() // Devolve o foco para a tabela de turmas
				return m, nil
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "salvar"))):
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
					m.nextFormInput()
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "próximo"))),
				key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "anterior"))):
				if msg.String() == "shift+tab" {
					m.prevFormInput()
				} else {
					m.nextFormInput()
				}
			}
			var focusedInputCmd tea.Cmd
			if m.createForm.focusIndex == 0 {
				m.createForm.nameInput, focusedInputCmd = m.createForm.nameInput.Update(msg)
			} else {
				m.createForm.subjectIDInput, focusedInputCmd = m.createForm.subjectIDInput.Update(msg)
			}
			cmds = append(cmds, focusedInputCmd)

		case DetailsView:
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "voltar"))):
				if m.detailsViewFocusTarget == FocusTargetStudentsTable {
					m.detailsViewFocusTarget = FocusTargetNone
					m.studentsTable.Blur()
				} else { // FocusTargetNone ou outros futuros focos não tratados aqui
					m.state = ListView
					m.selectedClass = nil
					m.classStudents = nil
					m.studentsTable.SetRows([]table.Row{})
					m.err = nil
					m.table.Focus()
				}

			case key.Matches(msg, key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "focar/navegar"))):
				if m.detailsViewFocusTarget == FocusTargetNone {
					m.detailsViewFocusTarget = FocusTargetStudentsTable
					m.studentsTable.Focus() // table.Focus() não retorna cmd, mas é necessário para o estado interno da tabela
				} else if m.detailsViewFocusTarget == FocusTargetStudentsTable {
					// Se já estiver na tabela, Tab pode remover o foco (ou ir para próximo elemento focável)
					m.detailsViewFocusTarget = FocusTargetNone
					m.studentsTable.Blur()
				}
				// Se houvesse outros elementos focáveis, Tab poderia ciclar entre eles.

			case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "remover foco"))):
				if m.detailsViewFocusTarget == FocusTargetStudentsTable {
					m.detailsViewFocusTarget = FocusTargetNone
					m.studentsTable.Blur()
				}
				// Se houvesse outros elementos focáveis, Shift+Tab poderia ciclar na ordem inversa.

			default:
				if m.detailsViewFocusTarget == FocusTargetStudentsTable {
					var studentsTableCmd tea.Cmd
					m.studentsTable, studentsTableCmd = m.studentsTable.Update(msg)
					cmds = append(cmds, studentsTableCmd)
				}
				// Lógica para outras teclas se nenhum elemento específico estiver focado
				// ou se o foco estiver em outros elementos (como o futuro seletor de status)
			}
		}


	// ... cases para fetchedClassesMsg, classCreatedMsg, errMsg ...
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
		} else {
			m.state = ListView
			m.err = nil
			m.isLoading = true
			m.table.Focus()
			cmds = append(cmds, m.fetchClassesCmd)
		}

	case errMsg:
		m.err = msg.err
		m.isLoading = false

	// Novo case para processar os alunos buscados
	case fetchedClassStudentsMsg:
		m.isLoading = false // Finaliza o estado de carregamento (de alunos)
		if msg.err != nil {
			m.err = msg.err // Exibe o erro se a busca de alunos falhar
			m.classStudents = nil
			m.studentsTable.SetRows([]table.Row{})
		} else {
			m.err = nil // Limpa erros anteriores se a busca for bem-sucedida
			m.classStudents = msg.students
			var rows []table.Row
			if len(m.classStudents) == 0 {
				// Adiciona uma linha indicando que não há alunos, se desejar
				// Ou simplesmente deixa a tabela vazia.
				// rows = append(rows, table.Row{"---", "Nenhum aluno encontrado", "---", "---"})
			} else {
				for _, student := range m.classStudents {
					rows = append(rows, table.Row{
						fmt.Sprintf("%d", student.ID),
						student.EnrollmentID, // Já é string ou pode ser ""
						student.FullName,
						student.Status,
					})
				}
			}
			m.studentsTable.SetRows(rows)
			// Opcionalmente, focar a tabela de alunos aqui se for a próxima interação principal
			// m.studentsTable.Focus()
		}
	} // Fim do switch msg.(type)

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

	case DetailsView: // Novo case
		if m.selectedClass == nil {
			// Isso não deveria acontecer se a lógica de estado estiver correta
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Erro: Nenhuma turma selecionada."))
		} else {
			titleStyle := lipgloss.NewStyle().Bold(true).PaddingBottom(1)
			detailHeaderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Cinza claro

			b.WriteString(titleStyle.Render(fmt.Sprintf("Detalhes da Turma: %s", m.selectedClass.Name)))
			b.WriteString(fmt.Sprintf("%s %d\n", detailHeaderStyle.Render("ID da Turma:"), m.selectedClass.ID))
			b.WriteString(fmt.Sprintf("%s %d\n\n", detailHeaderStyle.Render("ID da Disciplina:"), m.selectedClass.SubjectID))

			b.WriteString(lipgloss.NewStyle().Bold(true).Render("Alunos:"))
			if m.isLoading { // Se estiver carregando alunos especificamente para esta view
				b.WriteString("\nCarregando alunos...")
			} else if len(m.classStudents) == 0 && m.err == nil { // m.err == nil para não sobrescrever erro de busca
				b.WriteString("\nNenhum aluno encontrado para esta turma.")
			} else if m.err != nil && len(m.classStudents) == 0 {
				// Erro já impresso globalmente
			} else {
				studentsTableRender := m.studentsTable.View()
				if m.detailsViewFocusTarget == FocusTargetStudentsTable {
					// Adiciona uma borda sutil para indicar foco na tabela de alunos
					focusBorderStyle := lipgloss.NewStyle().
						Border(lipgloss.RoundedBorder()).
						BorderForeground(lipgloss.Color("63")). // Magenta para foco
						Padding(0,0) // Padding (0,1) pode ser muito se a tabela já tem suas margens

					studentsTableRender = focusBorderStyle.Render(studentsTableRender)
				}
				b.WriteString("\n" + studentsTableRender)
			}
		}
		helpStyle := lipgloss.NewStyle().Faint(true).MarginTop(1)
		currentHelp := "Pressione 'Esc' para voltar."
		if m.detailsViewFocusTarget == FocusTargetNone {
			currentHelp += " Use Tab para focar a tabela de alunos."
		} else if m.detailsViewFocusTarget == FocusTargetStudentsTable {
			currentHelp += " Use Setas para navegar, Shift+Tab ou Esc para sair do foco da tabela."
			// Futuramente: "... 's' para mudar status."
		}
		b.WriteString("\n\n" + helpStyle.Render(currentHelp))

	} // Fim do switch m.state
	return b.String()
}

// SetSize ajusta o tamanho do modelo e seus componentes.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Alturas comuns
	titleHeight := 1
	errorHeight := 0
	if m.err != nil {
		errorHeight = strings.Count(fmt.Sprintf("Erro: %v", m.err), "\n") + 1 + 1 // +1 for padding
	}
	helpHeight := 1

	// Altura disponível para o conteúdo principal da view
	contentHeight := height - titleHeight - errorHeight - helpHeight
	if contentHeight < 0 { contentHeight = 0 }


	switch m.state {
	case ListView:
		tableHeaderAndBorderHeight := 3
		tableBodyHeight := contentHeight - tableHeaderAndBorderHeight
		if tableBodyHeight < 1 { tableBodyHeight = 1 }
		m.table.SetHeight(tableBodyHeight)
		m.table.SetWidth(width - 4)

	case CreatingView:
		inputWidth := width - 4
		if inputWidth < 20 { inputWidth = 20 }
		m.createForm.nameInput.Width = inputWidth
		m.createForm.subjectIDInput.Width = inputWidth

	case DetailsView:
		// Altura para os detalhes da turma (ID, Nome, ID Disciplina) - estimativa
		classDetailsRenderedHeight := 4 // Título da view + ID Turma + ID Disciplina + linha em branco

		// Altura para o cabeçalho "Alunos:"
		studentsSectionHeaderHeight := 1

		// Altura disponível para a tabela de alunos
		studentsTableAvailableHeight := contentHeight - classDetailsRenderedHeight - studentsSectionHeaderHeight

		studentsTableHeaderAndBorderHeight := 3 // Estimativa para cabeçalho e bordas da studentsTable
		studentsTableBodyHeight := studentsTableAvailableHeight - studentsTableHeaderAndBorderHeight
		if studentsTableBodyHeight < 1 { studentsTableBodyHeight = 1}

		m.studentsTable.SetHeight(studentsTableBodyHeight)
		m.studentsTable.SetWidth(width - 4) // Margens laterais
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

type fetchedClassStudentsMsg struct {
	students []models.Student
	err      error
}

func (m Model) fetchClassStudentsCmd(classID int64) tea.Cmd {
	return func() tea.Msg {
		if m.classService == nil { // Checagem defensiva
			return errMsg{fmt.Errorf("classService não inicializado")}
		}
		// selectedClass é verificado antes de chamar este comando, mas uma checagem de classID aqui é boa.
		if classID == 0 {
			return errMsg{fmt.Errorf("ID da turma inválido (0) para buscar alunos")}
		}
		students, err := m.classService.GetStudentsByClassID(context.Background(), classID)
		if err != nil {
			return errMsg{fmt.Errorf("falha ao buscar alunos para a turma ID %d: %w", classID, err)}
		}
		return fetchedClassStudentsMsg{students: students, err: nil}
	}
}

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
