// Package tui fornece componentes e utilitários para a Interface de Texto do Usuário.
// Este arquivo (statusbar.go) define um modelo BubbleTea para uma barra de status personalizável.
package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Estilos para a barra de status e seus textos, usando lipgloss.
var (
	statusBarStyle = lipgloss.NewStyle().
			Height(1).                               // Altura fixa de 1 linha.
			Padding(0, 1).                          // Padding horizontal.
			Background(lipgloss.Color("236")).      // Cor de fundo (cinza escuro).
			Foreground(lipgloss.Color("250"))       // Cor do texto (cinza claro).
	statusTextLeftStyle  = lipgloss.NewStyle().Align(lipgloss.Left)  // Alinha o texto à esquerda.
	statusTextRightStyle = lipgloss.NewStyle().Align(lipgloss.Right) // Alinha o texto à direita.
	separator            = " | "                                      // Separador usado entre seções da barra (não usado ativamente na View atual).
)

// StatusBarModel é o modelo BubbleTea para a barra de status.
// Gerencia a largura, mensagens de status (permanentes e efêmeras) e texto à direita (ex: hora).
type StatusBarModel struct {
	width         int           // width armazena a largura atual do terminal/componente pai.
	status        string        // status é a mensagem principal exibida na barra.
	ephemeralMsg  string        // ephemeralMsg é uma mensagem temporária (ex: "Copiado!").
	ephemeralTime time.Time     // ephemeralTime registra quando a mensagem efêmera foi definida.
	ephemeralTTL  time.Duration // ephemeralTTL define por quanto tempo a mensagem efêmera é exibida.
	rightText     string        // rightText é o texto exibido no lado direito da barra (ex: hora atual).
}

// NewStatusBarModel cria uma nova instância de StatusBarModel com valores padrão.
// Define um status inicial "Ready", um TTL padrão para mensagens efêmeras e a hora atual.
func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{
		status:       "Pronto",
		ephemeralTTL: 2 * time.Second, // TTL padrão de 2 segundos para mensagens efêmeras.
		rightText:    time.Now().Format("15:04:05"),
	}
}

// Init é o comando inicial para o StatusBarModel.
// Inicia um ticker para atualizar a hora na barra de status a cada segundo.
func (m StatusBarModel) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return UpdateTimeMsg(t) // Envia uma UpdateTimeMsg a cada segundo.
	})
}

// UpdateTimeMsg é uma mensagem enviada para sinalizar que a hora na barra de status deve ser atualizada.
// Contém o novo timestamp.
type UpdateTimeMsg time.Time

// SetStatusMsg é uma mensagem para definir a mensagem de status principal (permanente).
// O payload é a string da nova mensagem de status.
type SetStatusMsg string

// SetEphemeralStatusMsg é uma mensagem para exibir uma mensagem de status temporária (efêmera).
// Contém o texto da mensagem e um TTL opcional.
type SetEphemeralStatusMsg struct {
	Text string        // Text é o conteúdo da mensagem efêmera.
	TTL  time.Duration // TTL (Time To Live) opcional; se zero, usa o ephemeralTTL padrão do modelo.
}

// ClearEphemeralMsg é uma mensagem para forçar a limpeza da mensagem efêmera.
type ClearEphemeralMsg struct{}

// Update lida com mensagens (eventos) para o StatusBarModel.
// Processa redimensionamento de janela, atualizações de hora, e definição/limpeza de mensagens de status.
// Retorna o modelo atualizado e quaisquer comandos a serem executados.
func (m StatusBarModel) Update(msg tea.Msg) (StatusBarModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // Mensagem de redimensionamento da janela.
		m.width = msg.Width // Atualiza a largura da barra.
	case UpdateTimeMsg: // Mensagem para atualizar a hora.
		m.rightText = time.Time(msg).Format("15:04:05") // Formata a nova hora.
		// Agenda o próximo tick para atualização da hora.
		cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return UpdateTimeMsg(t)
		}))
	case SetStatusMsg: // Define uma nova mensagem de status principal.
		m.status = string(msg)
		m.clearEphemeral() // Limpa qualquer mensagem efêmera existente.
	case SetEphemeralStatusMsg: // Define uma nova mensagem efêmera.
		m.ephemeralMsg = msg.Text
		m.ephemeralTime = time.Now()
		if msg.TTL > 0 { // Se um TTL específico foi fornecido, usa-o.
			m.ephemeralTTL = msg.TTL
		}
		// Agenda um comando para limpar a mensagem efêmera após seu TTL.
		cmds = append(cmds, tea.Tick(m.ephemeralTTL, func(t time.Time) tea.Msg {
			return ClearEphemeralMsg{}
		}))

	case ClearEphemeralMsg: // Limpa a mensagem efêmera.
		// Verifica se o tempo desde que a mensagem foi definida é maior ou igual ao TTL,
		// para garantir que estamos limpando a mensagem correta (caso múltiplos clears sejam agendados).
		if time.Since(m.ephemeralTime) >= m.ephemeralTTL {
			m.clearEphemeral()
		}
	}
	return m, tea.Batch(cmds...) // Retorna o modelo e os comandos acumulados.
}

// clearEphemeral redefine a mensagem efêmera para uma string vazia.
func (m *StatusBarModel) clearEphemeral() {
	m.ephemeralMsg = ""
}

// View renderiza o estado atual da StatusBarModel como uma string.
// Exibe a mensagem de status (efêmera, se ativa, ou a principal) à esquerda
// e o rightText (hora) à direita, preenchendo o espaço entre eles.
func (m StatusBarModel) View() string {
	if m.width == 0 {
		return "" // Não renderiza se a largura não estiver definida (evita pânico/layout quebrado).
	}

	displayStatus := m.status
	// Se houver uma mensagem efêmera e ela ainda não expirou, exibe-a.
	if m.ephemeralMsg != "" && time.Since(m.ephemeralTime) < m.ephemeralTTL {
		displayStatus = m.ephemeralMsg
	}
	// Nota: A limpeza explícita da ephemeralMsg é feita no Update via ClearEphemeralMsg.
	// Não é necessário lógica de expiração aqui na View para evitar mutação de estado.

	left := statusTextLeftStyle.Render(displayStatus)
	right := statusTextRightStyle.Render(m.rightText)

	// Calcula o espaço disponível para o texto esquerdo, considerando o texto direito.
	// Esta é uma abordagem simplificada para o preenchimento.
	availableWidth := m.width - lipgloss.Width(right)
	// Se o texto esquerdo for muito longo, trunca-o.
	// Esta truncagem é básica e pode não lidar bem com caracteres multibyte.
	// Para truncagem robusta, bibliotecas ou lógica mais sofisticada seriam necessárias.
	if lipgloss.Width(left) > availableWidth {
		runes := []rune(displayStatus)
		if availableWidth > 3 && len(runes) > (availableWidth-3) { // Garante espaço para "..."
			left = statusTextLeftStyle.Render(string(runes[:availableWidth-3]) + "...")
		} else if availableWidth > 0 && len(runes) > availableWidth { // Trunca sem "..." se muito pouco espaço
			left = statusTextLeftStyle.Render(string(runes[:availableWidth]))
		} else if availableWidth <= 0 { // Sem espaço
			left = ""
		}
		// Se availableWidth for suficiente, 'left' já está correto.
	}

	// Calcula o espaço de preenchimento entre o texto esquerdo e direito.
	gapWidth := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gapWidth < 0 {
		gapWidth = 0 // Evita preenchimento negativo se os cálculos não forem perfeitos.
	}
	gap := strings.Repeat(" ", gapWidth)

	// Renderiza a barra de status completa com a largura definida.
	return statusBarStyle.Width(m.width).Render(left + gap + right)
}
