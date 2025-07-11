// Package tui fornece componentes e utilitários para a Interface de Texto do Usuário.
// Este arquivo (prompt.go) define um modelo BubbleTea para solicitar entrada de texto do usuário.
package tui

import (
	"bufio" // Para ler de stdin não-TTY.
	"fmt"
	"io"
	"os" // Necessário para isatty.IsTerminal e os.Stdout/os.Stdin.
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty" // Para detecção de TTY.
)

// PromptModel é o modelo BubbleTea para um prompt de entrada de texto.
// Ele gerencia o estado do campo de texto, mensagens de erro e o valor submetido.
type PromptModel struct {
	prompt      string          // prompt é a mensagem exibida ao usuário.
	textInput   textinput.Model // textInput é o componente de entrada de texto do bubbles.
	err         error           // err armazena qualquer erro ocorrido durante a execução do prompt.
	quitting    bool            // quitting é true quando o modelo está prestes a sair.
	submitted   bool            // submitted é true se o usuário submeteu um valor (pressionou Enter).
	SubmittedCh chan string     // SubmittedCh é um canal para enviar o valor submetido de volta para o chamador.
}

// Estilos para o prompt, usando lipgloss.
var (
	promptStyle    = lipgloss.NewStyle().Padding(0, 1)                                // Estilo para o texto do prompt.
	focusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))            // Estilo para elementos focados (ex: cursor, botão).
	blurredStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))            // Estilo para elementos não focados.
	noStyle        = lipgloss.NewStyle()                                              // Estilo vazio.
	helpStyle      = blurredStyle.Copy()                                              // Estilo para texto de ajuda.
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))              // Estilo para mensagens de erro (vermelho).
	cursorModeHelp = helpStyle.Render("modo cursor está habilitado")                  // Mensagem de ajuda para modo cursor (não usado ativamente na View).
	focusedButton  = focusedStyle.Copy().Render("[ Submeter ]")                       // Aparência do botão de submissão quando focado.
	blurredButton  = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submeter")) // Aparência do botão de submissão quando não focado.
)

// NewPromptModel cria uma nova instância de PromptModel com o texto de prompt fornecido.
// Configura o componente textinput com placeholder, foco inicial, limite de caracteres e largura.
func NewPromptModel(promptText string) PromptModel {
	ti := textinput.New()
	ti.Placeholder = "Digite sua resposta aqui..."
	ti.Focus() // Define o foco inicial para o campo de texto.
	ti.CharLimit = 256
	ti.Width = 50
	ti.Prompt = "┃ " // Caractere de prompt antes do campo de texto.

	return PromptModel{
		prompt:      promptText,
		textInput:   ti,
		err:         nil,
		SubmittedCh: make(chan string, 1), // Canal bufferizado para evitar bloqueio.
	}
}

// Init é o comando inicial para o modelo PromptModel.
// Retorna textinput.Blink para iniciar a animação do cursor no campo de texto.
func (m PromptModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update lida com mensagens (eventos) para o PromptModel.
// Processa entradas de teclado (Enter, Esc, Ctrl+C) e atualiza o estado do textInput.
// Retorna o modelo atualizado e um comando (geralmente tea.Quit ou um comando do textInput).
func (m PromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Só submete se houver algum valor.
			// TODO: Adicionar validação opcional para permitir submissão de valor vazio se necessário.
			if m.textInput.Value() != "" {
				m.submitted = true
				m.quitting = true
				m.SubmittedCh <- m.textInput.Value() // Envia o valor submetido pelo canal.
				return m, tea.Quit                   // Encerra o programa BubbleTea para este prompt.
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			m.SubmittedCh <- "" // Envia string vazia para indicar que o usuário desistiu.
			close(m.SubmittedCh) // Fecha o canal para sinalizar que não haverá mais valores.
			return m, tea.Quit
		}

	case error: // Lida com mensagens de erro.
		m.err = msg
		return m, nil
	}

	// Atualiza o modelo do campo de texto e obtém qualquer comando resultante.
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renderiza o estado atual do PromptModel como uma string para exibição no terminal.
// Mostra o texto do prompt, o campo de entrada de texto e mensagens de ajuda ou erro.
func (m PromptModel) View() string {
	if m.quitting {
		// Se estiver saindo após submissão, pode-se optar por mostrar brevemente o valor submetido.
		// Atualmente, retorna string vazia para limpar a tela rapidamente.
		// if m.submitted {
		// 	return fmt.Sprintf("%s\n%s%s\n", m.prompt, m.textInput.Prompt, m.textInput.Value())
		// }
		return "" // Limpa a visualização ao sair.
	}

	var viewBuilder strings.Builder

	viewBuilder.WriteString(promptStyle.Render(m.prompt) + "\n")
	viewBuilder.WriteString(m.textInput.View() + "\n\n")
	// O botão é mais visual do que funcional em um prompt simples, mas mantido como exemplo.
	if m.textInput.Focused() {
		viewBuilder.WriteString(focusedButton)
	} else {
		viewBuilder.WriteString(blurredButton)
	}

	if m.err != nil {
		viewBuilder.WriteString("\n" + errorStyle.Render(fmt.Sprintf("Erro: %v", m.err)))
	}

	help := helpStyle.Render("Pressione Enter para submeter, Esc ou Ctrl+C para sair.")
	viewBuilder.WriteString("\n\n" + help)

	return viewBuilder.String()
}

// GetInput é uma função de alto nível que executa um prompt interativo para obter uma linha de entrada do usuário.
// Parâmetros:
//   - promptText: A mensagem a ser exibida para o usuário.
//   - output: O io.Writer para a saída da TUI (geralmente os.Stdout).
//   - inputReader: O io.Reader para a entrada da TUI (geralmente os.Stdin).
//
// Comportamento:
//   - Se inputReader (geralmente os.Stdin) não for um terminal TTY (ex: entrada redirecionada/piped),
//     ele tentará ler uma única linha diretamente do inputReader.
//   - Se for um TTY, ele iniciará um programa BubbleTea com o PromptModel para uma entrada interativa.
//
// Retorna:
//   - A string de entrada do usuário (sem espaços em branco no início/fim se lida de não-TTY).
//   - Um erro se o prompt for abortado pelo usuário (Esc/Ctrl+C) ou se ocorrer outro erro.
func GetInput(promptText string, output io.Writer, inputReader io.Reader) (string, error) {
	// Verifica se inputReader (normalmente os.Stdin) é um terminal.
	if f, ok := inputReader.(*os.File); ok && !isatty.IsTerminal(f.Fd()) && !isatty.IsCygwinTerminal(f.Fd()) {
		// Não é um TTY (ex: entrada via pipe: echo "texto" | ./app).
		// Imprime o prompt para stderr para não interferir com a saída de dados se o programa
		// também estiver produzindo saída para stdout que está sendo redirecionada.
		fmt.Fprintln(os.Stderr, promptText)

		scanner := bufio.NewScanner(inputReader)
		if scanner.Scan() {
			return strings.TrimSpace(scanner.Text()), nil // Retorna a linha lida.
		}
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("GetInput: erro ao ler de entrada não-TTY: %w", err)
		}
		// Se nada foi lido (ex: pipe vazio), retorna string vazia e nenhum erro.
		// O chamador pode precisar tratar isso como um erro se uma entrada for obrigatória.
		return "", nil
	}

	// É um TTY, ou inputReader não é os.Stdin; executa o prompt interativo BubbleTea.
	model := NewPromptModel(promptText)
	programOpts := []tea.ProgramOption{tea.WithOutput(output)}

	if inputReader != nil {
		programOpts = append(programOpts, tea.WithInput(inputReader))
	}

	p := tea.NewProgram(model, programOpts...)

	// Executa o programa BubbleTea e espera que ele termine.
	// O valor submetido será enviado através do canal `model.SubmittedCh`.
	// Uma alternativa seria capturar o modelo final retornado por p.Run().
	// No entanto, a comunicação via canal é explícita aqui.
	go func() {
		// É importante tratar o erro de Run, embora em um prompt simples
		// seja menos provável, a menos que haja problemas de configuração de terminal.
		if _, errRun := p.Run(); errRun != nil {
			// Logar ou tratar o erro de p.Run() é importante.
			// Para este prompt, se p.Run() falhar, o canal pode não receber valor.
			// Considerar fechar o canal ou enviar um erro específico por ele.
			// Por simplicidade, aqui apenas logamos.
			fmt.Fprintf(os.Stderr, "Error running prompt program: %v\n", errRun)
			// Tentar fechar o canal se não foi fechado pelo Update.
			// Isso é um pouco arriscado se Update ainda puder escrever.
			// A lógica de fechamento no Update(KeyEsc/CtrlC) é mais segura.
		}
	}()

	// Aguarda o valor submetido ou uma string vazia (se o usuário desistiu).
	submittedValue, ok := <-model.SubmittedCh
	if !ok {
		// Canal foi fechado sem valor (ex: Ctrl+C e canal fechado no Update).
		return "", fmt.Errorf("prompt abortado pelo usuário")
	}

	if model.err != nil { // Verifica se houve algum erro no modelo.
		return "", fmt.Errorf("GetInput: erro no modelo do prompt: %w", model.err)
	}

	if !model.submitted && submittedValue == "" { // Usuário desistiu (Esc/Ctrl+C).
		return "", fmt.Errorf("prompt abortado pelo usuário")
	}

	return submittedValue, nil
}

// main_example é uma função de exemplo para demonstrar o uso de GetInput.
// Ela é renomeada para evitar conflitos se este pacote for usado como biblioteca.
func main_example() {
	fmt.Println("Iniciando prompt de exemplo...")
	userInput, err := GetInput("Qual é o seu nome?", os.Stdout, os.Stdin)
	if err != nil {
		fmt.Printf("\nErro no prompt: %v\n", err)
		return
	}
	if userInput != "" {
		fmt.Printf("\nOlá, %s!\n", userInput)
	} else {
		fmt.Println("\nNenhuma entrada recebida ou prompt abortado.")
	}
}
