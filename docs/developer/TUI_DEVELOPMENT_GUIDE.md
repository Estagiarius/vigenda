# Guia de Desenvolvimento da TUI (Interface de Texto do Usuário)

Este guia destina-se a desenvolvedores que desejam entender, modificar ou estender a Interface de Texto do Usuário (TUI) da aplicação Vigenda. A TUI é construída usando o framework [Bubble Tea](https://github.com/charmbracelet/bubbletea) e seus componentes associados da Charm.

## 1. Visão Geral da Arquitetura da TUI

A TUI do Vigenda segue o padrão **The Elm Architecture (TEA)**, que é a base do Bubble Tea. Os principais conceitos são:

-   **Model:** Representa o estado completo da sua aplicação TUI ou de um componente específico. Em Vigenda, temos um `app.Model` principal (`internal/app/app.go`) que gerencia o estado geral e contém sub-modelos para cada funcionalidade/tela principal (ex: `tasks.Model`, `classes.Model`).
-   **View:** Uma função (`View() string`) que recebe o estado atual do `Model` e retorna uma string representando a interface do usuário a ser renderizada no terminal. A renderização é feita usando `lipgloss` para estilização.
-   **Update:** Uma função (`Update(msg tea.Msg) (tea.Model, tea.Cmd)`) que processa mensagens (eventos) e atualiza o `Model`. Mensagens podem ser entradas do usuário (teclado, mouse), respostas de operações assíncronas, ticks de temporizador, etc. O `Update` retorna o modelo atualizado e um `tea.Cmd` opcional (um comando a ser executado, como uma chamada de I/O).
-   **Messages (`tea.Msg`):** Representam eventos que podem ocorrer na aplicação (ex: `tea.KeyMsg` para teclas, `tea.WindowSizeMsg` para redimensionamento, ou mensagens customizadas para indicar a conclusão de uma tarefa assíncrona).
-   **Commands (`tea.Cmd`):** Representam efeitos colaterais que sua aplicação precisa executar (ex: fazer uma chamada HTTP, ler um arquivo, aguardar um tempo). Quando um comando termina, ele geralmente envia uma mensagem de volta para a função `Update`.

### Estrutura Principal em Vigenda:

-   **`internal/app/app.go` (`app.Model`):** É o coração da TUI.
    -   Gerencia a `currentView` (qual tela/módulo está ativo).
    -   Contém o menu principal (uma `list.Model`).
    -   Mantém instâncias dos sub-modelos para cada funcionalidade (ex: `tasksModel *tasks.Model`).
    -   Delega mensagens e renderização para o sub-modelo ativo.
    -   Lida com navegação global (voltar ao menu, sair).
-   **`internal/app/views.go`:** Define o enum `View` que identifica cada tela/módulo principal.
-   **Submódulos em `internal/app/` (ex: `internal/app/tasks/`, `internal/app/dashboard/`):**
    -   Cada submódulo define seu próprio `Model` BubbleTea, com seus próprios `Init`, `Update`, e `View`.
    -   Interagem com os `services` para buscar/manipular dados.
    -   Usam componentes de `bubbles` (como `list.Model`, `textinput.Model`, `spinner.Model`, `table.Model`) e `lipgloss` para estilização.
-   **`internal/tui/`:** Contém componentes TUI mais genéricos ou legados (ex: `prompt.go`, `table.go`, `statusbar.go`) que podem ser usados por diferentes partes da aplicação.

## 2. Adicionando uma Nova Visualização/Módulo Principal

Siga estes passos para adicionar uma nova funcionalidade principal à TUI (ex: um novo item de menu que leva a uma nova tela).

### Passo 2.1: Definir a Nova View

No arquivo `internal/app/views.go`:
1.  Adicione uma nova constante ao enum `View`. Siga o padrão `iota` se for sequencial ou atribua um valor explícito.
    ```go
    const (
        // ... outras views ...
        NovaFuncionalidadeView View = iota // Ou um valor explícito
    )
    ```
2.  Atualize o método `String()` para incluir o nome amigável da sua nova view:
    ```go
    func (v View) String() string {
        switch v {
        // ... outros casos ...
        case NovaFuncionalidadeView:
            return "Minha Nova Funcionalidade"
        default:
            return "Visualização Desconhecida"
        }
    }
    ```

### Passo 2.2: Criar o Pacote do Novo Módulo

1.  Crie um novo diretório em `internal/app/`. Por exemplo: `internal/app/novafunc/`.
2.  Dentro deste diretório, crie um arquivo para o seu modelo, ex: `novafunc.go`.

### Passo 2.3: Definir o Modelo do Novo Módulo

No arquivo `internal/app/novafunc/novafunc.go`:

```go
package novafunc

import (
	tea "github.com/charmbracelet/bubbletea"
	// Importe outros pacotes necessários (models, services, bubbles, lipgloss)
)

// Model representa o estado do módulo "Nova Funcionalidade".
type Model struct {
	// Campos para o estado do seu módulo:
	// Ex: dados carregados, estado de formulários, spinner, erros, etc.
	// Ex: data []models.AlgumTipo
	// Ex: isLoading bool
	// Ex: errMsg error
	// Ex: textInput textinput.Model
	// Ex: parent AppModel // Se precisar referenciar o modelo pai para navegação ou estado global
}

// New cria uma nova instância do Model para "Nova Funcionalidade".
// Injete quaisquer dependências de serviço necessárias.
func New(/* ... serviços ... */) *Model {
	return &Model{
		// Inicialize os campos
	}
}

// Init é chamado quando este módulo se torna ativo.
// Pode retornar um tea.Cmd para carregar dados iniciais.
func (m *Model) Init() tea.Cmd {
	// Ex: return m.loadDataCmd()
	return nil
}

// Update lida com mensagens para o módulo.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	// var cmds []tea.Cmd // Se precisar de múltiplos comandos

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Lide com entradas de teclado específicas para este módulo.
		// Ex: navegação em lista, submissão de formulário.
		// Se 'esc' for pressionado e o módulo não estiver focado em um subcomponente,
		// pode-se enviar uma mensagem para o AppModel para voltar ao menu (ver abaixo).
		// Para isso, o modelo pode ter um método IsFocused().

	// Lide com outras mensagens (ex: dados carregados, erros).
	// case dadosCarregadosMsg:
	//    m.data = msg.data
	//    m.isLoading = false
	// case errMsg:
	//    m.errMsg = msg.err
	}

	// Atualize componentes internos (ex: textInput, list)
	// m.textInput, cmd = m.textInput.Update(msg)
	// cmds = append(cmds, cmd)

	// return m, tea.Batch(cmds...)
	return m, cmd
}

// View renderiza a UI para o módulo "Nova Funcionalidade".
func (m *Model) View() string {
	if m.isLoading {
		return "Carregando Nova Funcionalidade..."
	}
	if m.errMsg != nil {
		return fmt.Sprintf("Erro: %v", m.errMsg)
	}
	// Construa e retorne a string da sua UI.
	return "Conteúdo da Minha Nova Funcionalidade"
}

// IsFocused (opcional): pode ser usado pelo AppModel para determinar
// se a tecla 'esc' deve voltar ao menu principal ou ser tratada internamente.
// func (m *Model) IsFocused() bool {
//    return m.textInput.Focused() // Exemplo
// }
```

### Passo 2.4: Integrar o Novo Módulo no `app.Model`

No arquivo `internal/app/app.go`:

1.  **Importe o novo pacote:**
    ```go
    import (
        // ... outros imports ...
        "vigenda/internal/app/novafunc" // Seu novo pacote
    )
    ```
2.  **Adicione um campo para o novo modelo na struct `Model`:**
    ```go
    type Model struct {
        // ... outros campos e sub-modelos ...
        novaFuncModel *novafunc.Model
    }
    ```
3.  **Instancie o novo modelo na função `New()`:**
    ```go
    func New(/* ... serviços ... */) *Model {
        // ... inicialização de outros modelos ...
        nfm := novafunc.New(/* injete serviços se necessário */)

        return &Model{
            // ... outras atribuições ...
            novaFuncModel: nfm,
        }
    }
    ```
4.  **Adicione um item de menu para a nova view na lista do `AppModel` (em `New()`):**
    ```go
    menuItems := []list.Item{
        // ... outros itens de menu ...
        menuItem{title: NovaFuncionalidadeView.String(), view: NovaFuncionalidadeView},
    }
    ```
5.  **No método `Update()` do `AppModel`:**
    *   Na seção que lida com a seleção de item de menu (`key.Matches(msg, key.NewBinding(key.WithKeys("enter"))`):
        ```go
        switch m.currentView {
        // ... outros casos ...
        case NovaFuncionalidadeView:
            cmds = append(cmds, m.novaFuncModel.Init()) // Chama o Init do seu novo modelo
        }
        ```
    *   Na seção que delega mensagens para o sub-modelo ativo:
        ```go
        switch m.currentView {
        // ... outros casos ...
        case NovaFuncionalidadeView:
            // Lógica para atualizar o sub-modelo e potencialmente voltar ao menu com 'esc'
            if km, ok := msg.(tea.KeyMsg); ok { // Se for uma mensagem de tecla
                var shouldReturn bool
                updatedSubModel, submodelCmd, shouldReturn = processSubmodelUpdate(m.novaFuncModel, km) // Use a função helper existente ou similar
                m.novaFuncModel = updatedSubModel.(*novafunc.Model)
                if shouldReturn {
                    m.currentView = DashboardView // Volta para o menu principal
                    log.Println("AppModel: Voltando para o Menu Principal a partir de Nova Funcionalidade.")
                }
            } else { // Para outras mensagens (não-tecla)
                updatedSubModel, submodelCmd = m.novaFuncModel.Update(msg)
                m.novaFuncModel = updatedSubModel.(*novafunc.Model)
            }
            cmds = append(cmds, submodelCmd)
        }
        ```
        *(Adapte a lógica de `processSubmodelUpdate` ou crie uma similar se o seu sub-modelo tiver um comportamento de foco específico para a tecla `Esc`)*

6.  **No método `View()` do `AppModel`:**
    ```go
    switch m.currentView {
    // ... outros casos ...
    case NovaFuncionalidadeView:
        viewContent = m.novaFuncModel.View()
        help = "\nPressione 'esc' para voltar ao menu principal." // Ou ajuda específica do módulo
    }
    ```

## 3. Usando Componentes `bubbles` e `lipgloss`

-   **`bubbles`:** Fornece componentes TUI prontos:
    -   `list.Model`: Para listas selecionáveis (usado no menu principal e em muitas sub-visualizações).
    -   `textinput.Model`: Para campos de entrada de texto.
    -   `textarea.Model`: Para entrada de texto multi-linha.
    -   `spinner.Model`: Para indicadores de carregamento.
    -   `table.Model`: Para exibir dados tabulares.
    -   `paginator.Model`: Para paginar conteúdo.
    -   `progress.Model`: Para barras de progresso.
    -   `viewport.Model`: Para conteúdo rolável.
    -   `help.Model`: Para exibir mensagens de ajuda contextuais baseadas em `key.Binding`.
    Cada componente tem seus próprios métodos `Init`, `Update`, `View` e geralmente é incorporado como um campo no seu `Model`.
-   **`lipgloss`:** Usado para estilizar strings. Permite definir cores de primeiro plano/fundo, negrito, itálico, sublinhado, bordas, padding, margin, alinhamento, etc.
    -   Crie estilos com `lipgloss.NewStyle().Foreground(lipgloss.Color("...")).Bold(true)...`
    -   Aplique estilos com `meuEstilo.Render("meu texto")`.
    -   Use `lipgloss.JoinVertical(...)` e `lipgloss.JoinHorizontal(...)` para compor layouts.

## 4. Gerenciamento de Estado e Mensagens

-   **Estado Centralizado vs. Distribuído:**
    -   O `AppModel` mantém o estado global (como `currentView`).
    -   Cada sub-modelo (`tasks.Model`, `classes.Model`, etc.) gerencia seu próprio estado interno.
-   **Mensagens (`tea.Msg`):**
    -   **Padrão:** `tea.KeyMsg`, `tea.MouseMsg`, `tea.WindowSizeMsg`.
    -   **Customizadas:** Defina suas próprias structs de mensagem para comunicar resultados de operações assíncronas ou outros eventos específicos da aplicação.
        ```go
        type dadosCarregadosMsg struct { dados []meuTipo }
        type erroOperacaoMsg struct { err error }
        ```
-   **Comandos (`tea.Cmd`):**
    -   Usados para executar operações que podem levar tempo ou ter efeitos colaterais (ex: chamadas de serviço).
    -   Uma função que retorna `tea.Msg` é um `tea.Cmd`.
        ```go
        func carregarDadosCmd(serv service.MeuServico) tea.Cmd {
            return func() tea.Msg {
                dados, err := serv.BuscarDados(context.Background())
                if err != nil {
                    return erroOperacaoMsg{err: err}
                }
                return dadosCarregadosMsg{dados: dados}
            }
        }
        ```
    -   Use `tea.Batch(...)` para agrupar múltiplos comandos.
-   **Comunicação entre Modelos:**
    -   **Pai para Filho:** O `AppModel` pode passar dados para sub-modelos via métodos ou atualizando campos diretamente (com cautela).
    -   **Filho para Pai:** Um sub-modelo pode enviar uma mensagem customizada que o `AppModel` escuta e processa. Alternativamente, o `Update` do sub-modelo pode retornar um estado/flag que o `AppModel` verifica para tomar decisões (ex: mudar de view).

## 5. Dicas e Melhores Práticas

-   **Mantenha os Modelos Pequenos:** Divida funcionalidades complexas em sub-modelos menores e mais gerenciáveis.
-   **Desacoplamento:** Use interfaces de serviço para interagir com a lógica de negócios, mantendo os modelos da TUI focados na apresentação.
-   **Tratamento de Erros:** Exiba erros de forma clara para o usuário. Use mensagens de erro customizadas.
-   **Feedback ao Usuário:** Use spinners ou mensagens para indicar operações em andamento.
-   **Navegação Clara:** Garanta que o usuário sempre saiba como voltar ou sair. A tecla `Esc` é um padrão comum para "voltar".
-   **Logging:** Adicione logs (como visto em `app.go` e `task_service.go`) para ajudar na depuração do fluxo da TUI, especialmente para o processamento de mensagens e comandos.

Este guia fornece uma base para o desenvolvimento da TUI no Vigenda. Consulte a documentação oficial do Bubble Tea e dos componentes Charm para informações mais detalhadas.
