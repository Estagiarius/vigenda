// Package tui fornece componentes e utilitários para a Interface de Texto do Usuário.
// Este arquivo (table.go) define um modelo BubbleTea para exibir dados em formato de tabela.
package tui

import (
	"fmt"
	"io" // Para io.Writer na função ShowTable.
	// "os" // Comentado pois os.Stdout é usado apenas no exemplo.

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// baseStyle é um estilo lipgloss base para a borda da tabela.
var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).       // Define uma borda normal.
	BorderForeground(lipgloss.Color("240")) // Define a cor da borda como cinza claro.

// TableModel é o modelo BubbleTea para exibir uma tabela.
// Ele encapsula o componente table.Model do pacote bubbles.
type TableModel struct {
	table table.Model // table é o componente de tabela subjacente.
}

// NewTableModel cria uma nova instância de TableModel com as colunas e linhas fornecidas.
// Configura a tabela com foco inicial, altura padrão e estilos personalizados para cabeçalho e linha selecionada.
func NewTableModel(columns []table.Column, rows []table.Row) TableModel {
	t := table.New(
		table.WithColumns(columns), // Define as colunas da tabela.
		table.WithRows(rows),       // Define as linhas de dados da tabela.
		table.WithFocused(true),    // Define a tabela para ter foco inicialmente (permite navegação).
		table.WithHeight(7),        // Define uma altura padrão para a tabela (pode ser ajustada).
	)

	// Personaliza os estilos da tabela.
	s := table.DefaultStyles()
	s.Header = s.Header. // Estilo para o cabeçalho da tabela.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")). // Cor da borda do cabeçalho.
				BorderBottom(true).                      // Adiciona uma borda inferior ao cabeçalho.
				Bold(false)                              // Remove negrito padrão do cabeçalho.
	s.Selected = s.Selected. // Estilo para a linha atualmente selecionada.
				Foreground(lipgloss.Color("229")). // Cor do texto da linha selecionada (branco brilhante).
				Background(lipgloss.Color("57")).  // Cor de fundo da linha selecionada (roxo).
				Bold(false)                         // Remove negrito da linha selecionada.
	t.SetStyles(s) // Aplica os estilos personalizados à tabela.

	return TableModel{table: t}
}

// Init é o comando inicial para o TableModel.
// Para este modelo simples de tabela, não há comando inicial necessário, então retorna nil.
func (m TableModel) Init() tea.Cmd {
	return nil
}

// Update lida com mensagens (eventos) para o TableModel.
// Processa entradas de teclado para navegação na tabela, foco/blur e sair.
// Retorna o modelo atualizado e um comando (geralmente do componente table ou tea.Quit).
func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc": // Tecla Escape.
			if m.table.Focused() {
				m.table.Blur() // Remove o foco da tabela.
			} else {
				m.table.Focus() // Coloca o foco na tabela.
			}
		case "q", "ctrl+c": // Teclas 'q' ou Ctrl+C.
			return m, tea.Quit // Encerra o programa BubbleTea.
			// TODO: Adicionar navegação (cima, baixo, enter para selecionar linha) se necessário
			// e passar essas teclas para m.table.Update(msg) se a tabela estiver focada.
		}
	}
	// Atualiza o modelo da tabela interna e obtém qualquer comando resultante.
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renderiza o estado atual do TableModel como uma string para exibição no terminal.
// Envolve a visualização da tabela interna com o baseStyle.
func (m TableModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

// ShowTable é uma função auxiliar para criar e executar rapidamente um programa BubbleTea
// que exibe uma tabela simples com as colunas e linhas fornecidas.
// Útil para exibir dados tabulares de forma isolada ou para depuração.
// Parâmetros:
//   - columns: As definições de coluna para a tabela.
//   - rows: Os dados da linha para a tabela.
//   - output: O io.Writer para a saída da TUI (geralmente os.Stdout).
func ShowTable(columns []table.Column, rows []table.Row, output io.Writer) {
	model := NewTableModel(columns, rows)
	p := tea.NewProgram(model, tea.WithOutput(output))
	if _, err := p.Run(); err != nil {
		// Em caso de erro ao executar o programa, imprime o erro no output.
		// Isso é importante para depuração, especialmente se a TUI não iniciar.
		fmt.Fprintf(output, "Erro ao executar o programa da tabela: %v\n", err)
	}
}

// Exemplo de Uso (pode ser movido para um arquivo main_example.go ou _test.go)
/*
package main // Ou um pacote de teste

import (
	"os"
	"vigenda/internal/tui" // Ajuste o caminho de importação conforme necessário

	"github.com/charmbracelet/bubbles/table"
)

func main() {
    columns := []table.Column{
        {Title: "ID", Width: 4},
        {Title: "Nome", Width: 10},
        {Title: "Idade", Width: 5},
    }

    rows := []table.Row{
        {"1", "Alice", "30"},
        {"2", "Bob", "24"},
        {"3", "Charlie", "35"},
    }

    // Usa tui.ShowTable para exibir a tabela.
    // Note que ShowTable é bloqueante até que o programa BubbleTea (tabela) seja encerrado.
    tui.ShowTable(columns, rows, os.Stdout)
}
*/
