# Análise do Código do Projeto Vigenda

## `cmd/vigenda/main.go`

**Propósito e Funcionalidade:**

Este arquivo é o ponto de entrada principal da aplicação Vigenda CLI. Ele é responsável por:
*   Definir a estrutura de comandos da CLI usando a biblioteca Cobra.
*   Inicializar a conexão com o banco de dados.
*   Configurar e inicializar os serviços da aplicação.
*   Lançar a interface TUI (Text User Interface) principal baseada em BubbleTea.
*   Gerenciar o logging da aplicação para um arquivo.

**Estruturas de Dados e Funções Chave:**

*   `db *sql.DB`: Pool de conexão global com o banco de dados.
*   `logFile *os.File`: Descritor de arquivo para o log.
*   Variáveis globais para cada serviço: `taskService`, `classService`, `assessmentService`, `questionService`, `proofService`.
*   `rootCmd`: O comando raiz do Cobra, que é executado quando `vigenda` é chamado sem subcomandos. Sua função `Run` inicia a aplicação TUI.
    *   `PersistentPreRunE`: Uma função executada antes de qualquer comando. Ela garante que o logging e a conexão com o banco de dados sejam configurados. Se o banco de dados ainda não estiver inicializado, ele estabelece a conexão (suportando SQLite e PostgreSQL, com SQLite como padrão) e, em seguida, chama `initializeServices`.
*   `initializeServices(db *sql.DB)`: Inicializa todas as instâncias de serviço, injetando seus respectivos repositórios (que, por sua vez, recebem a conexão `db`).
*   Comandos Cobra (`taskCmd`, `classCmd`, `assessmentCmd`, `questionBankCmd`, `proofCmd`):
    *   Cada um define um grupo de funcionalidades (ex: `vigenda tarefa ...`).
    *   Possuem subcomandos para ações específicas (ex: `taskAddCmd`, `taskListCmd`, `taskCompleteCmd`).
    *   Os subcomandos `Run` geralmente:
        *   Parseiam argumentos e flags.
        *   Interagem com o usuário para obter entradas adicionais (usando `tui.GetInput` para prompts simples).
        *   Chamam os métodos de serviço apropriados.
        *   Imprimem resultados ou mensagens de erro no console.
*   `setupLogging()`: Configura o logging para gravar em um arquivo `vigenda.log` no diretório de configuração do usuário ou no diretório atual como fallback. Inclui timestamps e informações de arquivo/linha no log.
*   `init()`: Função padrão do Go que é executada na inicialização do pacote. Aqui, ela é usada para:
    *   Definir flags para os comandos Cobra.
    *   Adicionar subcomandos aos seus comandos pais.
    *   Adicionar os comandos principais ao `rootCmd`.
*   `main()`: Ponto de entrada da aplicação. Executa o `rootCmd` e garante que o arquivo de log seja fechado ao final da execução.

**Decisões de Arquitetura e Funcionalidade:**

1.  **CLI com Cobra:** A escolha do Cobra facilita a criação de uma CLI robusta com subcomandos, flags e documentação de ajuda gerada automaticamente. Isso torna a aplicação extensível e fácil de usar a partir do terminal.
2.  **Interface TUI com BubbleTea:** Para interações mais complexas e uma experiência de usuário mais rica do que uma CLI pura, o `rootCmd` (sem subcomandos) lança uma aplicação BubbleTea. Isso permite uma navegação baseada em menus e formulários interativos.
3.  **Camada de Serviço e Repositório:** A aplicação segue um design em camadas, separando a lógica de apresentação (CLI/TUI) da lógica de negócios (serviços) e do acesso a dados (repositórios). Isso promove a modularidade e testabilidade.
    *   `main.go` interage principalmente com a camada de serviço.
4.  **Gerenciamento de Dependências (Serviços):** Os serviços são inicializados centralmente em `initializeServices` e passados para as partes da aplicação que precisam deles (como a TUI ou os manipuladores de comando Cobra).
5.  **Configuração de Banco de Dados Flexível:** O `PersistentPreRunE` permite configurar o tipo de banco de dados (SQLite, PostgreSQL) e o DSN através de variáveis de ambiente, oferecendo flexibilidade para diferentes ambientes de implantação. SQLite é o padrão para facilitar o uso local.
6.  **Logging em Arquivo:** A decisão de logar em um arquivo ajuda na depuração e no rastreamento de problemas, especialmente para uma aplicação CLI/TUI onde o stdout/stderr é usado para interação com o usuário.
7.  **Prompts Interativos para Comandos CLI:** Para alguns comandos CLI que requerem entrada adicional (ex: descrição de uma tarefa), `tui.GetInput` é usado para solicitar informações ao usuário de forma interativa, melhorando a usabilidade.
8.  **Tratamento de Erros:** Os erros dos serviços são geralmente capturados e impressos no console para o usuário. O logging também captura erros mais detalhados.

**Como se Encaixa no Projeto:**

`main.go` é o orquestrador da aplicação. Ele define a interface do usuário (seja CLI via subcomandos Cobra ou a TUI principal) e conecta essa interface à lógica de negócios subjacente, garantindo que todas as dependências (como banco de dados e serviços) estejam prontas antes que qualquer comando seja executado. Ele serve como a "cola" que une as diferentes partes do sistema Vigenda.
A transição para uma TUI principal (quando `vigenda` é executado sem subcomandos) indica uma preferência por uma experiência interativa para as funcionalidades centrais, enquanto os subcomandos Cobra oferecem uma maneira mais direta e scriptável de acessar funcionalidades específicas.

## `internal/app/app.go`

**Propósito e Funcionalidade:**

Este arquivo define o modelo principal e a lógica da aplicação TUI (Text User Interface) usando a biblioteca BubbleTea. Ele gerencia a navegação entre diferentes "visões" ou módulos da aplicação (como gerenciamento de tarefas, turmas, etc.) e delega a lógica específica de cada visão para seus respectivos submodelos.

**Estruturas de Dados e Funções Chave:**

*   `Model struct`: A estrutura principal do modelo BubbleTea para a aplicação. Contém:
    *   `list list.Model`: Um componente de lista para o menu principal de navegação.
    *   `currentView View`: Enum que indica a visão atualmente ativa.
    *   Ponteiros para os submodelos de cada módulo: `tasksModel`, `classesModel`, `assessmentsModel`, `questionsModel`, `proofsModel`.
    *   Dimensões da janela (`width`, `height`).
    *   Estado de saída (`quitting`) e erros (`err`).
    *   Instâncias dos serviços (`taskService`, `classService`, etc.) injetadas para serem usadas pelos submodelos.
*   `New(...) *Model`: Construtor para o `Model`. Inicializa a lista do menu principal com itens correspondentes às diferentes `View`s da aplicação. Também inicializa todos os submodelos, injetando os serviços necessários.
*   `menuItem struct`: Define um item para a lista do menu principal, associando um título a uma `View`.
*   `(m *Model) Init() tea.Cmd`: Método de inicialização do modelo principal. Geralmente retorna `nil` ou comanda a inicialização de um submodelo ativo.
*   `(m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)`: O coração da lógica de atualização do BubbleTea. Processa mensagens (eventos de teclado, redimensionamento de janela, mensagens de submodelos, erros).
    *   **Tratamento Global:** Lida com `tea.WindowSizeMsg` para redimensionar a si mesmo e propagar para os submodelos. Trata `tea.KeyMsg` para sair globalmente (`ctrl+c`).
    *   **Navegação no Menu Principal:** Se `currentView == DashboardView`, gerencia a navegação na lista do menu. Ao pressionar "Enter", muda `currentView` para a visão selecionada e envia o comando `Init()` do submodelo correspondente.
    *   **Delegação para Submodelos:** Se não estiver no `DashboardView`, encaminha a mensagem `msg` para o `Update` do submodelo atualmente ativo (`tasksModel`, `classesModel`, etc.).
    *   **Retorno ao Menu:** Se a tecla "Esc" for pressionada em uma subvisão e o submodelo não estiver "focado" (ou seja, não estiver em um estado interno que precise tratar "Esc", como um formulário), `currentView` é alterada de volta para `DashboardView`.
*   `(m *Model) View() string`: Renderiza a UI com base no estado atual.
    *   Se estiver saindo (`quitting`) ou houver um erro (`err`), exibe mensagens apropriadas.
    *   Caso contrário, renderiza a `View()` do submodelo ativo ou a lista do menu principal se `currentView == DashboardView`.
    *   Adiciona texto de ajuda contextual na parte inferior da tela.
*   `StartApp(...)`: Função auxiliar que cria uma nova instância do `Model` (com os serviços injetados) e inicia o programa BubbleTea.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Arquitetura Baseada em Submodelos (Model-View-Update - MVU):** A aplicação TUI é estruturada em torno de um modelo principal (`app.Model`) que gerencia e delega para submodelos (`tasks.Model`, `classes.Model`, etc.). Cada submodelo é responsável por sua própria lógica de estado, atualização e renderização, seguindo o padrão MVU do BubbleTea. Isso promove a modularidade e facilita o gerenciamento de funcionalidades complexas.
2.  **Navegação Centralizada:** O `app.Model` atua como um roteador, controlando qual subvisão está ativa e como o usuário navega entre elas (menu principal e tecla "Esc" para voltar).
3.  **Injeção de Dependência (Serviços):** Os serviços da camada de negócios são injetados no `app.Model` e, subsequentemente, passados para os submodelos que precisam deles. Isso desacopla a TUI da lógica de negócios e facilita os testes.
4.  **Gerenciamento de Estado Global vs. Local:** O `app.Model` lida com o estado global da TUI (visão atual, erro global), enquanto os submodelos gerenciam seu estado local específico.
5.  **Comunicação por Mensagens:** A interação entre componentes e a atualização de estado ocorrem através do sistema de mensagens do BubbleTea (`tea.Msg`). Comandos (`tea.Cmd`) são usados para realizar operações assíncronas (como carregar dados de um serviço) que resultarão em novas mensagens.
6.  **Reutilização de Componentes BubbleTea:** Utiliza componentes padrão como `list.Model` para menus, e cada submodelo provavelmente usa outros componentes (como `table.Model`, `textinput.Model`).

**Como se Encaixa no Projeto:**

`app.go` é o núcleo da interface interativa do Vigenda quando o comando `vigenda` é executado sem subcomandos. Ele fornece a estrutura principal da TUI, permitindo que o usuário navegue e interaja com as diferentes funcionalidades da aplicação (tarefas, turmas, avaliações, etc.) de forma visual e interativa no terminal. Ele se integra com a camada de serviço para buscar e manipular dados.

## `internal/app/app_test.go`

**Propósito e Funcionalidade:**

Este arquivo contém testes unitários para o `internal/app/app.go`, focando no comportamento do modelo principal da TUI (`app.Model`).

**Estruturas de Dados e Funções Chave:**

*   **Mock Services:** Define estruturas de mock para cada interface de serviço (`mockTaskService`, `mockClassService`, etc.). Essas mocks implementam as interfaces de serviço, mas geralmente não fazem nada ou retornam dados/erros predefinidos, permitindo isolar o `app.Model` de dependências reais durante os testes.
*   `newTestAppModel() *Model`: Função helper que cria uma instância de `app.Model` com todos os serviços mockados.
*   **Test Functions (usando `testify/assert` e `testify/require`):**
    *   `TestNewModel_InitialState`: Verifica se o modelo é inicializado corretamente (visão inicial, itens da lista do menu).
    *   `TestModel_Update_Quit`: Testa se o modelo transita para o estado `quitting` e retorna o comando `tea.Quit` ao receber as teclas 'q' (no dashboard) ou 'ctrl+c'.
    *   `TestModel_Update_NavigateToSubViewAndBack`: Simula a navegação do menu principal para uma subvisão (ex: `TaskManagementView`) e o retorno usando "Esc". Verifica se `currentView` é atualizado corretamente.
    *   `TestModel_View_Content`: Verifica se o conteúdo renderizado (`View()`) contém os elementos esperados, como o título da lista no `DashboardView` e o texto de ajuda apropriado para diferentes visões. Também simula a navegação e o processamento de comandos de submodelos para garantir que a visão do submodelo seja renderizada.
    *   `TestModel_Update_WindowSize`: Testa se o modelo atualiza suas dimensões (`width`, `height`) e redimensiona a lista do menu corretamente ao receber uma `tea.WindowSizeMsg`.
    *   `TestMenuItem_Interface`: Testa os métodos da struct `menuItem` (usada na lista do menu).
    *   `TestView_String`: Testa o método `String()` do enum `View` para garantir que os nomes das visões sejam retornados corretamente.
    *   Funções helper de simulação de teclas: `simulateKeyPress`, `simulateEnterPress`, `simulateEscPress`, `simulateCtrlCPress` para simplificar a escrita dos testes de atualização.
    *   `TestModel_Update_WithHelpers`: Reimplementa alguns testes de atualização usando as helpers de simulação de teclas.

**Decisões de Arquitetura e Funcionalidade (nos Testes):**

1.  **Mocks para Dependências:** O uso de mocks para os serviços é crucial para testes unitários, pois permite testar a lógica do `app.Model` isoladamente, sem depender do comportamento real dos serviços ou do banco de dados.
2.  **Testes de Estado e Transição:** Os testes focam em verificar se o estado do modelo (`currentView`, `quitting`, `err`) é atualizado corretamente em resposta a diferentes mensagens (especialmente `tea.KeyMsg`).
3.  **Testes de Renderização (Básicos):** `TestModel_View_Content` faz verificações básicas no conteúdo renderizado. Testes de UI mais complexos exigiriam "snapshot testing" ou ferramentas mais especializadas, mas para TUI, verificar a presença de strings chave é uma abordagem comum.
4.  **Simulação de Eventos:** Os testes simulam eventos do BubbleTea (como `tea.KeyMsg`, `tea.WindowSizeMsg`) para acionar a lógica de `Update()`.
5.  **Helpers para Simulação de Teclas:** A introdução de funções como `simulateEnterPress` torna os testes mais legíveis e menos repetitivos.

**Como se Encaixa no Projeto:**

`app_test.go` garante a corretude e a robustez do componente central da TUI. Ao testar a navegação, o gerenciamento de estado e a resposta a eventos, esses testes ajudam a prevenir regressões e a validar o comportamento esperado da interface principal da aplicação.

## `internal/app/views.go`

**Propósito e Funcionalidade:**

Este arquivo define o enum `View`, que é usado em toda a aplicação TUI (`internal/app/*`) para representar as diferentes telas ou módulos disponíveis para o usuário. Ele também fornece um método `String()` para obter uma representação textual amigável de cada valor do enum, usada principalmente para títulos de menu e logs.

**Estruturas de Dados e Funções Chave:**

*   `View int`: Define `View` como um tipo inteiro, a base para o enum.
*   `const (...)`: Bloco de constantes que define os valores do enum para cada visão:
    *   `DashboardView`
    *   `TaskManagementView`
    *   `ClassManagementView`
    *   `AssessmentManagementView`
    *   `QuestionBankView`
    *   `ProofGenerationView`
*   `func (v View) String() string`: Método associado ao tipo `View`. Retorna uma string descritiva para cada valor do enum. Por exemplo, `DashboardView.String()` retorna "Dashboard".

**Decisões de Arquitetura e Funcionalidade:**

1.  **Enum para Visões:** Usar um enum para representar as visões é uma prática comum que melhora a legibilidade e a segurança de tipo em comparação com o uso de strings ou inteiros mágicos. Torna o código mais fácil de entender e manter, pois as visões são claramente definidas e referenciadas.
2.  **Método `String()` para Representação:** Fornecer um método `String()` permite que o enum seja facilmente convertido para uma forma legível por humanos, útil para UIs (títulos de menu, cabeçalhos de tela) e logging.

**Como se Encaixa no Projeto:**

`views.go` fornece um tipo centralizado e seguro para gerenciar os diferentes estados de tela da aplicação TUI. O `app.Model` usa este enum para rastrear a visão atual e para construir o menu de navegação. Os submodelos também podem se referir a essas constantes se precisarem interagir ou sinalizar transições de visão.
## `internal/app/assessments/model.go`

**Propósito e Funcionalidade:**

Este arquivo define o modelo BubbleTea e a lógica para o módulo de gerenciamento de avaliações e notas dentro da TUI. Ele permite ao usuário listar avaliações, criar novas avaliações, lançar notas para alunos e calcular médias de turma.

**Estruturas de Dados e Funções Chave:**

*   `ViewState int`: Enum que define os diferentes estados internos da visão de avaliações:
    *   `ListView`: Lista de ações principais (criar, lançar notas, etc.).
    *   `CreateAssessmentView`: Formulário para criar uma nova avaliação.
    *   `EnterGradesView`: Lógica para selecionar uma avaliação e, em seguida, inserir notas para os alunos dessa avaliação.
    *   `ClassAverageView`: Lógica para selecionar uma turma e calcular sua média.
    *   `ListAssessmentsView`: Visão para exibir uma tabela de avaliações existentes.
*   `Model struct`: O modelo BubbleTea para este módulo. Contém:
    *   `assessmentService service.AssessmentService`: Serviço para interagir com a lógica de negócios de avaliações.
    *   `state ViewState`: O estado atual da visão.
    *   Componentes BubbleTea: `list.Model` (para menus de ação/seleção), `table.Model` (para exibir avaliações/alunos), `textInputs []textinput.Model` (para formulários).
    *   `focusIndex int`: Para gerenciar o foco em formulários com múltiplos campos.
    *   Dados: `assessments []models.Assessment`, `currentClassID *int64`, `currentAssessmentID *int64`, `studentsForGrading []models.Student`, `gradesInput map[int64]textinput.Model` (map de ID do aluno para campo de input de nota).
    *   `isLoading bool`, `err error`, `message string`: Para feedback ao usuário.
    *   `width`, `height`: Dimensões da visão.
*   **Mensagens (tea.Msg):**
    *   `assessmentsLoadedMsg`, `assessmentCreatedMsg`, `studentsForGradingLoadedMsg`, `gradesEnteredMsg`, `classAverageCalculatedMsg`: Mensagens para comunicar resultados de operações assíncronas (chamadas de serviço) de volta para o `Update`.
*   **Comandos (tea.Cmd):**
    *   `loadAssessmentsCmd`, `submitCreateAssessmentFormCmd`, `loadStudentsForGradingCmd`, `submitGradesCmd`: Funções que retornam `tea.Cmd` para executar operações assíncronas (geralmente chamadas de serviço).
*   `New(...) *Model`: Construtor para o `Model`. Inicializa a lista de ações, a tabela e os inputs de texto.
*   `(m *Model) Init() tea.Cmd`: Reseta o estado do modelo para a visão inicial (lista de ações).
*   `(m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)`: Lida com mensagens.
    *   **Teclas:** Trata "Esc" para voltar ao estado anterior ou à lista de ações. Em cada `state`, trata "Enter" para submeter formulários ou selecionar itens, e outras teclas para navegação em formulários ou tabelas.
    *   **Mensagens Assíncronas:** Processa os resultados das operações de serviço (`assessmentsLoadedMsg`, etc.), atualizando o estado, dados e UI.
    *   **Redimensionamento:** `tea.WindowSizeMsg` chama `SetSize`.
*   `(m *Model) View() string`: Renderiza a UI com base no `state` atual.
    *   Exibe mensagens de carregamento, erro ou sucesso.
    *   Renderiza a lista de ações, formulários de criação/entrada de ID, tabela de avaliações, ou a interface de lançamento de notas para alunos.
*   **Funções de Formulário:**
    *   `resetForms()`, `setupCreateAssessmentForm()`, `setupEnterAssessmentIDForm()`, `setupEnterClassIDForm()`: Preparam os campos de texto para diferentes formulários.
    *   `updateFocus()`, `updateInputFocusStyle()`, `updateFormInputs()`: Gerenciam a navegação e atualização dos campos de texto em formulários.
*   `SetSize(width, height int)`: Ajusta o tamanho dos componentes internos (lista, tabela, inputs) com base no tamanho da janela.
*   `IsFocused() bool`: Indica se o modelo está em um estado que deve capturar eventos de teclado (ex: dentro de um formulário).
*   Validadores (`isNumber`, `isFloatOrEmpty`): Funções simples para validar entradas de texto.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Submodelo Dedicado:** Segue o padrão de ter um submodelo específico para a funcionalidade de avaliações, mantendo o código organizado e modular dentro da TUI principal.
2.  **Gerenciamento de Estado Interno:** Usa um enum `ViewState` para controlar as diferentes telas e interações dentro do módulo de avaliações.
3.  **Formulários Interativos:** Emprega `textinput.Model` para criar formulários para entrada de dados (criar avaliação, IDs).
4.  **Tabela para Listagem:** Usa `table.Model` para exibir uma lista de avaliações de forma estruturada.
5.  **Interação Assíncrona com Serviços:** As chamadas para `assessmentService` são feitas de forma assíncrona através de `tea.Cmd`, com os resultados sendo processados como `tea.Msg`.
6.  **Feedback ao Usuário:** Fornece feedback visual para estados de carregamento, mensagens de erro e sucesso.
7.  **Interface de Lançamento de Notas (Simplificada):** A funcionalidade de lançar notas envolve carregar alunos e, em seguida, apresentar campos de entrada para cada um. A navegação e submissão de múltiplas notas em uma única tela TUI é um desafio e a implementação atual parece ser uma versão simplificada, com potencial para melhorias na navegação entre os campos de nota.

**Como se Encaixa no Projeto:**

Este módulo é uma das principais funcionalidades do Vigenda, permitindo ao professor gerenciar o ciclo de vida das avaliações. Ele se integra com o `assessmentService` para realizar operações de backend e com o `app.Model` principal para navegação. A clareza na apresentação das informações e a facilidade na entrada de dados são cruciais para a usabilidade desta seção.

## `internal/app/classes/model.go`

**Propósito e Funcionalidade:**

Este arquivo define o modelo BubbleTea e a lógica para o módulo de gerenciamento de turmas e alunos na TUI. Ele permite ao usuário listar, criar, editar e deletar turmas, bem como visualizar, adicionar, editar e deletar alunos dentro de uma turma selecionada.

**Estruturas de Dados e Funções Chave:**

*   `ViewState int`: Enum para os estados da visão de turmas: `ListView` (lista de turmas), `CreatingView`, `EditingClassView`, `DeletingClassConfirmView`, `DetailsView` (detalhes da turma com lista de alunos), `AddingStudentView`, `EditingStudentView`, `DeletingStudentConfirmView`.
*   `FocusTarget int`: Enum para controlar o foco dentro da `DetailsView` (entre informações da turma e a tabela de alunos).
*   `Model struct`: O modelo BubbleTea para este módulo. Contém:
    *   `classService service.ClassService`: Serviço para interagir com a lógica de negócios de turmas/alunos.
    *   `state ViewState`: O estado atual da visão.
    *   `table table.Model`: Tabela para listar turmas.
    *   `studentsTable table.Model`: Tabela para listar alunos na `DetailsView`.
    *   `formInputs`: Struct aninhada para gerenciar os campos de texto (`[]textinput.Model`) e o índice de foco (`focusIndex`) para os formulários de turma e aluno.
    *   Dados: `allClasses []models.Class`, `selectedClass *models.Class`, `selectedStudent *models.Student`, `classStudents []models.Student`.
    *   `detailsViewFocusTarget FocusTarget`: Controla qual componente tem foco na tela de detalhes.
    *   `isLoading bool`, `width int`, `height int`, `err error`.
*   **Mensagens (tea.Msg):**
    *   `fetchedClassesMsg`, `classCreatedMsg`, `classUpdatedMsg`, `classDeletedMsg`: Para operações CRUD de turmas.
    *   `fetchedClassStudentsMsg`, `studentAddedMsg`, `studentUpdatedMsg`, `studentDeletedMsg`: Para operações CRUD de alunos.
    *   `errMsg`: Para erros gerais.
*   **Comandos (tea.Cmd):**
    *   Funções que retornam `tea.Cmd` para executar operações de serviço de forma assíncrona (ex: `fetchClassesCmd`, `createClassCmd`, `fetchClassStudentsCmd`).
*   `New(...) *Model`: Construtor. Inicializa as tabelas de turmas e alunos, e define o estado inicial.
*   `(m *Model) Init() tea.Cmd`: Carrega a lista inicial de turmas.
*   `(m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)`: Processa mensagens.
    *   **Teclas:** Chama `handleKeyPress` que, por sua vez, delega para manipuladores específicos do estado (`handleListViewKeys`, `handleClassFormKeys`, etc.). Esses manipuladores gerenciam a navegação, seleção, e transições de estado (ex: 'n' para nova turma, 'Enter' para detalhes/salvar).
    *   **Mensagens Assíncronas:** Processa os resultados das operações de serviço, atualizando os dados (`allClasses`, `classStudents`) e as tabelas.
    *   **Redimensionamento:** `tea.WindowSizeMsg` chama `SetSize`.
*   `(m *Model) View() string`: Renderiza a UI com base no `state`.
    *   Exibe mensagens de carregamento e erro.
    *   Renderiza a tabela de turmas, formulários de criação/edição de turma/aluno, diálogos de confirmação de exclusão, ou a tela de detalhes da turma (incluindo a tabela de alunos).
    *   Fornece texto de ajuda contextual.
*   **Manipuladores de Teclas Específicos do Estado:** (ex: `handleListViewKeys`, `handleDetailsViewKeys`)
    *   `handleListViewKeys`: Navegação na tabela de turmas, 'n' (novo), 'e' (editar), 'd' (deletar), 'Enter' (detalhes).
    *   `handleClassFormKeys`, `handleStudentFormKeys`: Navegação nos campos do formulário (Tab, Shift+Tab), 'Enter' (salvar), 'Esc' (cancelar).
    *   `handleDeleteClassConfirmKeys`, `handleDeleteStudentConfirmKeys`: Confirmação 's'/'n'.
    *   `handleDetailsViewKeys`: 'Esc' (voltar), 'a' (add aluno), 'Tab' (focar tabela de alunos), e, dentro da tabela de alunos, navegação e 'e' (editar aluno), 'd' (deletar aluno).
*   **Gerenciamento de Formulários:**
    *   `resetFormInputs()`, `prepareClassForm(...)`, `prepareStudentForm(...)`: Configuram os `textinput.Model` para os formulários.
    *   `nextFormInput()`, `prevFormInput()`: Gerenciam a navegação de foco entre os campos.
*   `SetSize(width, height int)`: Ajusta o tamanho das tabelas e campos de formulário.
*   `IsFocused() bool`: Indica se o modelo está em um estado de formulário.

**Decisões de Arquitetura e Funcionalidade:**

1.  **CRUD Completo para Turmas e Alunos:** O módulo oferece uma interface TUI abrangente para todas as operações básicas de gerenciamento de turmas e seus respectivos alunos.
2.  **Múltiplas Visões e Estados:** A complexidade é gerenciada através de múltiplos `ViewState`s, cada um com sua própria lógica de UI e interação.
3.  **Foco Contextual em Detalhes:** A `DetailsView` implementa um sistema de foco (`detailsViewFocusTarget`) para permitir que o usuário interaja alternadamente com as informações da turma e a lista de alunos.
4.  **Reutilização de Formulários:** A estrutura `formInputs` é usada tanto para formulários de turma quanto de aluno, com funções `prepare...Form` adaptando os campos e placeholders.
5.  **Confirmação para Exclusão:** Ações destrutivas como deletar turmas ou alunos requerem uma etapa de confirmação para evitar exclusões acidentais.
6.  **Carregamento Assíncrono de Dados:** Dados como a lista de turmas e a lista de alunos de uma turma específica são carregados assincronamente.

**Como se Encaixa no Projeto:**

Este é um módulo central para a organização do professor, permitindo a estruturação de suas turmas e o acompanhamento dos alunos. Ele depende fortemente do `classService` para persistir e recuperar dados. A interface TUI visa tornar essas operações administrativas eficientes e acessíveis diretamente do terminal.

## `internal/app/proofs/model.go`

**Propósito e Funcionalidade:**

Este arquivo define o modelo BubbleTea e a lógica para o módulo de geração de provas na TUI. Ele permite ao usuário especificar critérios (ID da disciplina, tópico, número de questões por dificuldade) e, em seguida, visualiza a prova gerada.

**Estruturas de Dados e Funções Chave:**

*   `ViewState int`: Enum para os estados da visão de provas:
    *   `FormView`: Formulário para inserir os critérios de geração da prova.
    *   `ProofView`: Visão para exibir as questões da prova gerada.
*   `Model struct`: O modelo BubbleTea para este módulo. Contém:
    *   `proofService service.ProofService`: Serviço para interagir com a lógica de negócios de geração de provas.
    *   `state ViewState`: O estado atual da visão.
    *   `textInputs []textinput.Model`: Cinco campos de texto para: ID da Disciplina, Tópico, Qtd. Fáceis, Qtd. Médias, Qtd. Difíceis.
    *   `focusIndex int`: Para gerenciar o foco entre os campos do formulário e o botão de submissão.
    *   `generatedProof []models.Question`: Armazena as questões da prova após a geração.
    *   `isLoading bool`, `err error`, `message string`: Para feedback ao usuário.
    *   `width`, `height`: Dimensões da visão.
*   **Mensagens (tea.Msg):**
    *   `proofGeneratedMsg`: Mensagem para comunicar o resultado (prova ou erro) da operação de geração de prova.
*   **Comandos (tea.Cmd):**
    *   `generateProofCmd()`: Função que retorna `tea.Cmd` para executar a geração da prova de forma assíncrona. Coleta os valores dos `textInputs`, constrói `service.ProofCriteria` e chama `proofService.GenerateProof`.
*   `New(...) *Model`: Construtor. Inicializa os campos de texto do formulário.
*   `(m *Model) Init() tea.Cmd`: Reseta o estado do modelo para `FormView`, limpa erros/mensagens e prepara o formulário.
*   `(m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)`: Lida com mensagens.
    *   **Teclas:**
        *   "Esc": Se em `ProofView`, volta para `FormView`. Se em `FormView`, é tratado pelo `app.Model` principal para voltar ao menu.
        *   Em `FormView`:
            *   "Enter": Se o foco estiver no "botão" de submissão (simulado pelo `focusIndex == len(textInputs)`), inicia a geração da prova. Caso contrário, move o foco para o próximo campo.
            *   "Up/Shift+Tab", "Down/Tab": Navega entre os campos do formulário e o botão de submissão.
            *   Outras teclas: Passadas para o `textinput.Model` focado.
    *   **Mensagens Assíncronas:**
        *   `proofGeneratedMsg`: Se sucesso e a prova tiver questões, muda para `ProofView` e armazena a prova. Se erro ou nenhuma questão, exibe o erro e permanece em `FormView`.
    *   **Redimensionamento:** `tea.WindowSizeMsg` chama `SetSize`.
*   `(m *Model) View() string`: Renderiza a UI com base no `state`.
    *   Exibe mensagens de carregamento, erro ou informativas.
    *   Se `FormView`: Renderiza os campos de texto e um botão de submissão simulado. Inclui ajuda para navegação.
    *   Se `ProofView`: Renderiza o título "Prova Gerada" e, em seguida, cada questão formatada (enunciado, tipo, dificuldade, opções, resposta). Inclui ajuda para voltar.
*   `resetForm()`: Limpa os campos de texto, reseta o foco para o primeiro campo.
*   `updateInputFocusStyle()`: Atualiza o estilo visual dos campos de texto para indicar qual está focado.
*   `isNumberOrEmpty()`: Validador simples para campos numéricos que podem ser vazios.
*   `SetSize(width, height int)`: Ajusta a largura dos campos de texto.
*   `IsFocused() bool`: Retorna `true` se estiver em `FormView`, indicando que o módulo deve capturar entradas de teclado.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Interface de Formulário:** Utiliza múltiplos `textinput.Model` para coletar os critérios de geração da prova de forma interativa.
2.  **Dois Estados Principais:** Alterna entre um formulário de entrada (`FormView`) e uma visualização dos resultados (`ProofView`).
3.  **Geração Assíncrona:** A geração da prova (que pode envolver consultas ao banco de dados) é feita assincronamente.
4.  **Exibição Formatada da Prova:** As questões da prova são renderizadas de forma legível, incluindo detalhes como tipo, dificuldade, opções (se múltipla escolha) e a resposta correta.
5.  **Feedback Detalhado:** Fornece mensagens de erro específicas se a geração falhar ou se não houver questões suficientes.

**Como se Encaixa no Projeto:**

Este módulo oferece uma funcionalidade chave para professores, permitindo-lhes criar avaliações customizadas a partir do banco de questões. Ele depende do `proofService` (que por sua vez usa o `questionRepository`) para a lógica de seleção de questões. A TUI visa tornar esse processo de criação de provas simples e direto.

## `internal/app/questions/model.go`

**Propósito e Funcionalidade:**

Este arquivo define o modelo BubbleTea e a lógica para o módulo de gerenciamento do banco de questões na TUI. Atualmente, sua principal funcionalidade implementada é permitir ao usuário adicionar novas questões a partir de um arquivo JSON.

**Estruturas de Dados e Funções Chave:**

*   `ViewState int`: Enum que define os estados internos da visão de questões:
    *   `ActionListView`: Lista de ações principais (ex: "Adicionar Questões de JSON").
    *   `AddQuestionsFormView`: Formulário para inserir o caminho do arquivo JSON.
    *   `ListQuestionsView` (Comentado/Opcional): Poderia ser usado para listar questões existentes.
*   `Model struct`: O modelo BubbleTea para este módulo. Contém:
    *   `questionService service.QuestionService`: Serviço para interagir com a lógica de negócios do banco de questões.
    *   `state ViewState`: O estado atual da visão.
    *   `list list.Model`: Componente de lista para o menu de ações.
    *   `textInputs []textinput.Model`: Um campo de texto para o caminho do arquivo JSON.
    *   `focusIndex int`: Gerencia o foco (relevante se houvesse múltiplos campos ou botões).
    *   `isLoading bool`, `err error`, `message string`: Para feedback ao usuário.
    *   `width`, `height`: Dimensões da visão.
*   **Mensagens (tea.Msg):**
    *   `questionsAddedMsg`: Mensagem para comunicar o resultado (número de questões adicionadas ou erro) da operação de adição de questões.
*   **Comandos (tea.Cmd):**
    *   `submitAddQuestionsFormCmd()`: Função que retorna `tea.Cmd` para executar a adição de questões de forma assíncrona. Lê o caminho do arquivo do `textInputs[0]`, lê o conteúdo do arquivo, e chama `questionService.AddQuestionsFromJSON`.
*   `New(...) *Model`: Construtor. Inicializa a lista de ações e o campo de texto para o caminho do arquivo.
*   `(m *Model) Init() tea.Cmd`: Reseta o estado do modelo para `ActionListView`.
*   `actionItem struct`: Item para a `list.Model`.
*   `(m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)`: Lida com mensagens.
    *   **Teclas:**
        *   "Esc": Se em `AddQuestionsFormView`, volta para `ActionListView`. Se em `ActionListView`, é tratado pelo `app.Model` principal.
        *   Em `ActionListView`: "Enter" para selecionar uma ação (atualmente, apenas "Adicionar Questões de JSON", que muda o estado para `AddQuestionsFormView`).
        *   Em `AddQuestionsFormView`: "Enter" para submeter o formulário (chama `submitAddQuestionsFormCmd`). Outras teclas são passadas para o campo de texto.
    *   **Mensagens Assíncronas:**
        *   `questionsAddedMsg`: Se sucesso, exibe mensagem de sucesso e volta para `ActionListView`. Se erro, exibe o erro.
    *   **Redimensionamento:** `tea.WindowSizeMsg` chama `SetSize`.
*   `(m *Model) View() string`: Renderiza a UI com base no `state`.
    *   Exibe mensagens de carregamento, erro ou sucesso.
    *   Se `ActionListView`: Renderiza a lista de ações.
    *   Se `AddQuestionsFormView`: Renderiza o campo de texto para o caminho do arquivo JSON e instruções.
*   `resetForms()`: Limpa o campo de texto.
*   `setupAddQuestionsForm()`: Prepara o campo de texto para o caminho do arquivo.
*   `SetSize(width, height int)`: Ajusta o tamanho da lista e do campo de texto.
*   `IsFocused() bool`: Retorna `true` se estiver em `AddQuestionsFormView`.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Foco na Importação JSON:** A funcionalidade TUI atual para o banco de questões está centrada na importação de questões via arquivos JSON, alinhando-se com o comando CLI `vigenda bancoq add`.
2.  **Interface Simples:** A interface para adicionar questões é direta, solicitando apenas o caminho do arquivo.
3.  **Feedback de Operação:** O usuário recebe mensagens sobre o sucesso ou falha da importação.
4.  **Potencial para Expansão:** A estrutura com `ActionListView` e o comentado `ListQuestionsView` sugere que o módulo foi projetado com a intenção de adicionar mais funcionalidades TUI para o banco de questões no futuro (como listar, editar ou criar questões individualmente pela TUI).

**Como se Encaixa no Projeto:**

Este módulo fornece a interface TUI para popular o banco de dados de questões, que é um componente essencial para a funcionalidade de geração de provas. Ele depende do `questionService` para processar o arquivo JSON e persistir as questões.

## `internal/app/tasks/model.go`

**Propósito e Funcionalidade:**

Este arquivo define o modelo BubbleTea e a lógica para o módulo de gerenciamento de tarefas na TUI. Ele permite ao usuário visualizar tarefas pendentes e concluídas, adicionar novas tarefas, editar tarefas existentes, marcar tarefas como concluídas e excluir tarefas.

**Estruturas de Dados e Funções Chave:**

*   `ViewState int`: Enum para os estados principais da visão de tarefas: `TableView`, `FormView`, `DetailView`, `ConfirmDeleteView`.
*   `FormState int`: Enum para o sub-estado do formulário: `CreatingTask`, `EditingTask`.
*   `FocusedTable int`: Enum para indicar qual tabela (pendentes ou concluídas) tem o foco: `PendingTableFocus`, `CompletedTableFocus`.
*   `Model struct`: O modelo BubbleTea para este módulo. Contém:
    *   `taskService service.TaskService`: Serviço para interagir com a lógica de negócios de tarefas.
    *   `pendingTasksTable table.Model`: Tabela para exibir tarefas pendentes.
    *   `completedTasksTable table.Model`: Tabela para exibir tarefas concluídas.
    *   `isLoading bool`, `err error`: Para feedback ao usuário.
    *   `currentView ViewState`, `formSubState FormState`, `focusedTable FocusedTable`: Gerenciamento de estado.
    *   `inputs []textinput.Model`: Quatro campos de texto para o formulário de tarefa (Título, Descrição, Prazo, ID da Turma).
    *   `focusIndex int`: Para gerenciar o foco nos campos do formulário.
    *   `selectedTaskForDetail *models.Task`: Tarefa selecionada para visualização de detalhes ou edição.
    *   `editingTaskID int64`, `taskIDToDelete int64`: IDs para operações de edição/exclusão.
    *   `width`, `height`: Dimensões da visão.
*   **Mensagens (tea.Msg):**
    *   `tasksLoadedMsg`, `fetchedTaskDetailMsg`, `taskCreatedMsg`, `taskCreationFailedMsg`, `taskUpdatedMsg`, `taskUpdateFailedMsg`, `taskDeletedMsg`, `taskDeleteFailedMsg`, `taskMarkedCompletedMsg`, `taskMarkCompleteFailedMsg`: Mensagens para comunicar resultados de operações assíncronas.
*   **Comandos (tea.Cmd):**
    *   Funções que retornam `tea.Cmd` para executar operações de serviço de forma assíncrona (ex: `loadTasksCmd`, `createTaskCmd`, `fetchTaskForDetailCmd`, `updateTaskCmd`, `deleteTaskCmd`, `markTaskCompleteCmd`).
*   `New(...) *Model`: Construtor. Inicializa as tabelas, campos de formulário e estado inicial.
*   `(m *Model) Init() tea.Cmd`: Carrega a lista inicial de todas as tarefas.
*   `(m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)`: Lida com mensagens.
    *   **Mensagens Assíncronas:** Processa os resultados das operações de serviço, atualizando os dados das tabelas, mudando de estado (ex: volta para `TableView` após criar/editar) ou exibindo erros.
    *   **Teclas:**
        *   Em `ConfirmDeleteView`: 's' para confirmar exclusão, 'n'/'Esc' para cancelar.
        *   Em `FormView`: "Esc" para cancelar e voltar para `TableView`. "Enter" para submeter o formulário (criar/editar tarefa) ou avançar para o próximo campo. "Tab"/"Shift+Tab" ou "Up"/"Down" para navegar entre os campos.
        *   Em `DetailView`: "Esc"/"q" para voltar para `TableView`.
        *   Em `TableView`:
            *   'a': Mudar para `FormView` (Criar Tarefa).
            *   'e': Mudar para `FormView` (Editar Tarefa selecionada, apenas pendentes). Carrega detalhes da tarefa primeiro.
            *   'c': Marcar tarefa pendente selecionada como concluída.
            *   'd': Mudar para `ConfirmDeleteView` para a tarefa selecionada.
            *   'v'/'Enter': Mudar para `DetailView` para a tarefa selecionada. Carrega detalhes da tarefa.
            *   'Tab': Alternar foco entre tabela de pendentes e concluídas.
            *   Teclas de navegação (Up, Down, j, k): Passadas para a tabela focada.
    *   **Redimensionamento:** `tea.WindowSizeMsg` chama `SetSize`.
*   `(m *Model) View() string`: Renderiza a UI com base no `currentView`.
    *   Se `ConfirmDeleteView`: Exibe diálogo de confirmação.
    *   Se `FormView`: Chama `viewForm()` para renderizar o formulário.
    *   Se `DetailView`: Chama `viewTaskDetail()` para renderizar os detalhes da tarefa.
    *   Se `TableView` (padrão): Renderiza os cabeçalhos e as tabelas de tarefas pendentes e concluídas, e a ajuda.
*   `viewForm()`: Renderiza o formulário de criação/edição de tarefa.
*   `viewTaskDetail()`: Renderiza a visualização detalhada de uma tarefa.
*   `nextInput()`, `prevInput()`: Gerenciam a navegação de foco nos campos do formulário.
*   `resetFormInputs()`: Limpa os campos do formulário.
*   `SetSize(width, height int)`: Ajusta o tamanho das tabelas e campos de formulário.
*   `IsFocused() bool`: Retorna `true` se a visão não for `TableView`, indicando que o módulo deve capturar "Esc".
*   `IsLoading() bool`: Retorna o estado de carregamento.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Visualização Dupla de Tarefas:** Apresenta tarefas pendentes e concluídas em tabelas separadas, com a capacidade de alternar o foco entre elas. Tarefas concluídas são exibidas com estilo tachado.
2.  **Operações CRUD Completas:** Fornece uma interface TUI para criar, visualizar (listar e detalhar), editar, marcar como concluída e excluir tarefas.
3.  **Formulário Reutilizável:** Um único conjunto de campos de formulário (`inputs`) é usado tanto para criar novas tarefas quanto para editar tarefas existentes, com o estado `formSubState` controlando a lógica de submissão.
4.  **Feedback Visual e de Estado:** Usa `isLoading` para feedback de carregamento e `err` para exibir mensagens de erro.
5.  **Navegação Intuitiva:** Utiliza teclas comuns (Tab, Enter, Esc, setas) para navegação em tabelas e formulários.

**Como se Encaixa no Projeto:**

O módulo de tarefas é fundamental para a organização pessoal e pedagógica do professor. Ele permite o rastreamento de pendências, sejam elas gerais ou associadas a turmas específicas. A integração com `taskService` garante que as operações TUI sejam refletidas no backend. A clareza na distinção entre tarefas pendentes e concluídas, e a facilidade de manipulação, são aspectos chave de sua usabilidade.
## `internal/config/config.go`

**Propósito e Funcionalidade:**

Este arquivo parece ser um placeholder ou um início para uma funcionalidade de gerenciamento de configuração (por exemplo, a partir de um arquivo `config.toml`). No entanto, em seu estado atual, ele está praticamente vazio, contendo apenas a declaração do pacote.

**Estruturas de Dados e Funções Chave:**

*   Nenhuma estrutura de dados ou função significativa definida no código fornecido.

**Decisões de Arquitetura e Funcionalidade:**

*   A existência do arquivo sugere que havia uma intenção de ter um sistema de configuração mais elaborado, possivelmente para carregar configurações de um arquivo TOML ou similar, o que é uma prática comum para aplicações configuráveis.
*   Atualmente, a configuração principal (como DSN do banco de dados) é gerenciada através de variáveis de ambiente em `cmd/vigenda/main.go`.

**Como se Encaixa no Projeto:**

No momento, este arquivo não desempenha um papel ativo significativo. Se a aplicação evoluir para necessitar de configurações mais complexas do que as gerenciadas por variáveis de ambiente (por exemplo, configurações de TUI, preferências do usuário, etc.), este pacote seria o local ideal para implementar essa lógica.

## `internal/database/connection.go`

**Propósito e Funcionalidade:**

Este arquivo é responsável por estabelecer e gerenciar a conexão com o banco de dados. Ele abstrai os detalhes da conexão, permitindo que diferentes tipos de banco de dados (SQLite, PostgreSQL) sejam usados com uma configuração comum.

**Estruturas de Dados e Funções Chave:**

*   `DBConfig struct`: Define os parâmetros necessários para a conexão:
    *   `DBType string`: Tipo do banco de dados ("sqlite" ou "postgres").
    *   `DSN string`: Data Source Name para a conexão.
*   `GetDBConnection(config DBConfig) (*sql.DB, error)`:
    *   Recebe uma `DBConfig`.
    *   Seleciona o driver do banco de dados (`sqlite3` ou `postgres`) com base em `config.DBType`.
    *   Se `config.DSN` estiver vazio para SQLite, usa `DefaultSQLitePath()` para obter um caminho padrão.
    *   Abre a conexão com o banco de dados usando `sql.Open`.
    *   Verifica a conexão usando `db.Ping()`.
    *   Para SQLite, chama `applySQLiteSchema()` para garantir que o schema inicial seja aplicado se o banco de dados for novo (verificando a existência da tabela 'users').
*   `applySQLiteSchema(db *sql.DB) error`:
    *   Lê o arquivo `migrations/001_initial_schema.sql` embutido no binário (usando `embed.FS` de `database.go`).
    *   Executa o SQL do schema no banco de dados SQLite.
*   `DefaultSQLitePath() string`:
    *   Retorna um caminho padrão para o arquivo do banco de dados SQLite. Tenta criar um diretório `vigenda` dentro do diretório de configuração do usuário (`os.UserConfigDir()`). Se falhar, usa `vigenda.db` no diretório de trabalho atual.
*   `DefaultDbPath() string`: Função legada ou alternativa que chama `DefaultSQLitePath()`.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Abstração da Conexão:** Centraliza a lógica de conexão, tornando mais fácil para o resto da aplicação obter uma instância `*sql.DB` sem se preocupar com os detalhes específicos do driver ou DSN.
2.  **Suporte a Múltiplos Bancos de Dados:** Projetado para suportar SQLite e PostgreSQL, permitindo flexibilidade na implantação. SQLite é usado como padrão para simplicidade.
3.  **Schema Embutido para SQLite:** O schema inicial do SQLite é embutido e aplicado automaticamente, facilitando a configuração inicial e garantindo que o banco de dados tenha a estrutura correta sem a necessidade de arquivos SQL externos em tempo de execução.
4.  **Caminho Padrão Inteligente para SQLite:** Tenta usar o diretório de configuração do usuário, que é uma prática recomendada para armazenar dados de aplicação.
5.  **Verificação de Conexão (Ping):** Garante que a conexão com o banco de dados seja válida antes de retorná-la.

**Como se Encaixa no Projeto:**

`connection.go` é um componente crucial da camada de acesso a dados. Ele é usado por `cmd/vigenda/main.go` para inicializar a conexão principal com o banco de dados que será usada por todos os repositórios. A aplicação automática do schema para SQLite simplifica muito a configuração para novos usuários.

## `internal/database/database.go`

**Propósito e Funcionalidade:**

Este arquivo complementa `connection.go` e parece ter sido parte de uma estrutura anterior ou alternativa para gerenciamento de banco de dados e migrações. Ele define uma `embed.FS` para os arquivos de migração SQL e inclui funções para aplicar essas migrações, especificamente para SQLite.

**Estruturas de Dados e Funções Chave:**

*   `migrationsFS embed.FS`: Usando a diretiva `go:embed`, este campo embute todos os arquivos `.sql` do diretório `migrations` diretamente no binário da aplicação. Isso garante que os scripts de schema estejam sempre disponíveis.
*   `DBConfig_database struct`: Uma struct de configuração de banco de dados, similar (e potencialmente redundante) à `DBConfig` em `connection.go`. O sufixo `_database` sugere uma tentativa de evitar conflito de nomes.
*   `GetDBConnection_database(config DBConfig_database) (*sql.DB, error)`: Similar à `GetDBConnection` em `connection.go`. Abre uma conexão, faz ping, e para SQLite, chama `applyMigrations_database` e `SeedData`.
*   `applyMigrations_database(db *sql.DB) error`:
    *   Lê os arquivos `.sql` do `migrationsFS` embutido.
    *   Executa cada script de migração no banco de dados SQLite.
    *   Inclui lógica para tentar executar statements individualmente se a execução do bloco falhar (útil para scripts SQL com múltiplos statements separados por ';').
*   `DefaultSQLitePath_database() string`: Similar à `DefaultSQLitePath` em `connection.go`, para determinar o caminho padrão do arquivo SQLite.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Migrações Embutidas:** Embutir os arquivos de migração SQL no binário é uma excelente prática para garantir que a aplicação possa configurar seu próprio schema sem depender de arquivos externos, tornando a implantação mais simples e robusta.
2.  **Aplicação de Migrações:** A função `applyMigrations_database` fornece um mecanismo para executar os scripts de schema. A tentativa de executar statements individualmente em caso de falha do bloco aumenta a robustez contra pequenas variações na sintaxe SQL permitida por `db.Exec()`.
3.  **Seed de Dados (Opcional):** A chamada para `SeedData` (de `seed.go`) após as migrações para SQLite permite popular o banco de dados com dados de exemplo, o que é útil para desenvolvimento e demonstração.
4.  **Potencial Redundância:** A existência de `DBConfig_database` e `GetDBConnection_database` que são muito semelhantes às suas contrapartes em `connection.go` sugere uma possível refatoração ou consolidação pendente. `connection.go` parece ser a versão mais "atual" ou a que está sendo usada primariamente pelo `main.go`.

**Como se Encaixa no Projeto:**

`database.go` (juntamente com `connection.go` e `seed.go`) forma a espinha dorsal da configuração e inicialização do banco de dados. A capacidade de embutir e aplicar migrações automaticamente é uma grande vantagem para a portabilidade e facilidade de configuração do Vigenda, especialmente ao usar SQLite. A lógica de `applyMigrations_database` é usada por `connection.go` para o schema inicial do SQLite.

## `internal/database/seed.go`

**Propósito e Funcionalidade:**

Este arquivo é responsável por popular o banco de dados com dados de exemplo (seeding) caso ele esteja vazio. Isso é útil para desenvolvimento, testes e para fornecer uma experiência inicial rica para novos usuários.

**Estruturas de Dados e Funções Chave:**

*   `SeedData(db *sql.DB) error`:
    *   Verifica se o banco de dados já foi populado (contando registros na tabela `users`). Se já houver dados, não faz nada.
    *   Inicia uma transação no banco de dados.
    *   **Insere Dados de Exemplo:**
        1.  Um usuário de exemplo (`demo_user`).
        2.  Disciplinas de exemplo ("Matemática", "História") associadas ao `demo_user`.
        3.  Turmas de exemplo ("Turma A - Matemática", "Turma B - Matemática", "Turma Única - História") associadas às disciplinas criadas.
        4.  Alunos de exemplo distribuídos nas turmas criadas, com diferentes status ("ativo", "inativo", "transferido").
    *   Faz commit da transação se todas as inserções forem bem-sucedidas. Caso contrário, faz rollback.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Preenchimento Condicional:** Os dados só são inseridos se o banco de dados estiver presumivelmente vazio (verificando a tabela `users`), evitando a duplicação de dados em execuções subsequentes.
2.  **Uso de Transações:** Todas as operações de seeding são envolvidas em uma transação, garantindo que o banco de dados seja populado de forma atômica (ou tudo ou nada).
3.  **Dados de Exemplo Abrangentes:** Os dados de exemplo cobrem as principais entidades do sistema (usuários, disciplinas, turmas, alunos), fornecendo um conjunto útil para testar e demonstrar as funcionalidades da aplicação.
4.  **Placeholder de Senha:** Conscientemente usa uma senha placeholder (`hashed_password_placeholder`) com um comentário sobre a necessidade de hashing real em produção, o que é uma boa prática de segurança.

**Como se Encaixa no Projeto:**

`seed.go` desempenha um papel importante na experiência de desenvolvimento e na configuração inicial da aplicação, especialmente quando se usa SQLite. Ao popular o banco de dados com dados de exemplo, ele permite que os desenvolvedores e testadores comecem a usar e verificar as funcionalidades da aplicação imediatamente, sem a necessidade de inserir manualmente uma grande quantidade de dados preliminares. É chamado por `GetDBConnection_database` em `database.go` após a aplicação das migrações para SQLite.

## `internal/models/models.go`

**Propósito e Funcionalidade:**

Este arquivo define as estruturas Go (structs) que representam as entidades de dados da aplicação Vigenda. Essas estruturas são usadas para transferir dados entre as camadas da aplicação (TUI, serviços, repositórios) e para mapear os dados de e para o banco de dados.

**Estruturas de Dados Chave:**

Cada struct representa uma tabela ou conceito no banco de dados:

*   `User`: Representa um usuário da aplicação.
    *   `ID`, `Username`, `PasswordHash` (marcado com `json:"-"` para não ser exposto em APIs JSON).
*   `Subject`: Representa uma disciplina (ex: Matemática).
    *   `ID`, `UserID` (proprietário da disciplina), `Name`.
*   `Class`: Representa uma turma dentro de uma disciplina.
    *   `ID`, `UserID` (proprietário, indiretamente via Subject), `SubjectID`, `Name`, `CreatedAt`, `UpdatedAt`.
*   `Student`: Representa um aluno em uma turma.
    *   `ID`, `ClassID`, `FullName`, `EnrollmentID` (matrícula/chamada), `Status` ('ativo', 'inativo', etc.), `CreatedAt`, `UpdatedAt`.
*   `Lesson`: Representa uma aula planejada.
    *   `ID`, `ClassID`, `Title`, `PlanContent` (Markdown), `ScheduledAt`.
*   `Assessment`: Representa uma avaliação (prova, trabalho).
    *   `ID`, `ClassID`, `Name`, `Term` (bimestre/período), `Weight` (peso na média), `AssessmentDate` (ponteiro para permitir nulo).
*   `Grade`: Representa a nota de um aluno em uma avaliação.
    *   `ID`, `AssessmentID`, `StudentID`, `Grade` (valor numérico).
*   `Task`: Representa uma tarefa ou item "to-do".
    *   `ID`, `UserID` (proprietário), `ClassID` (opcional, ponteiro para nulo), `Title`, `Description` (opcional), `DueDate` (opcional, ponteiro para nulo), `IsCompleted`.
*   `Question`: Representa uma questão no banco de questões.
    *   `ID`, `UserID` (proprietário), `SubjectID`, `Topic` (opcional), `Type` ('multipla_escolha', 'dissertativa'), `Difficulty` ('facil', 'media', 'dificil'), `Statement` (enunciado), `Options` (JSON string para múltipla escolha, ponteiro para nulo), `CorrectAnswer`.
*   `ModelError`: Tipo customizado para erros específicos da camada de modelo (ex: `ErrClassNotFound`).

**Decisões de Arquitetura e Funcionalidade:**

1.  **Estruturas de Dados Claras:** As structs são bem definidas e refletem de perto o schema do banco de dados (conforme Artefacto 3.2).
2.  **Tags JSON:** Muitas structs incluem tags `json:"..."` para controlar como seriam serializadas/desserializadas para JSON, o que é útil se a aplicação tiver uma API HTTP ou importar/exportar dados JSON. O uso de `omitempty` para campos opcionais e `json:"-"` para campos sensíveis como `PasswordHash` são boas práticas.
3.  **Uso de Ponteiros para Campos Opcionais/Nulos:** Campos como `ClassID` e `DueDate` em `Task`, `AssessmentDate` em `Assessment`, e `Options` em `Question` são definidos como ponteiros (`*int64`, `*time.Time`, `*string`). Isso permite que eles representem valores NULOS no banco de dados, o que é crucial para campos opcionais.
4.  **Timestamps:** `CreatedAt` e `UpdatedAt` em `Class` e `Student` são úteis para rastrear quando os registros foram criados e modificados pela última vez. (Nota: O schema da DB no Artefacto 3.2 não mostra consistentemente `created_at`/`updated_at` para todas as tabelas que poderiam se beneficiar deles, como `tasks` ou `questions`. As structs aqui estão um pouco à frente do schema nesse aspecto, ou o schema precisa ser atualizado.)
5.  **Tipo de Erro Personalizado:** A introdução de `ModelError` permite erros mais específicos e semânticos originados da camada de modelo ou repositório.

**Como se Encaixa no Projeto:**

`models.go` é fundamental para todo o projeto. Essas estruturas de dados são a linguagem comum usada por todas as camadas da aplicação para representar e manipular os dados do Vigenda. Qualquer alteração no schema do banco de dados provavelmente exigirá uma alteração correspondente neste arquivo. Eles são a base para a lógica de negócios nos serviços e para as operações de persistência nos repositórios.
## `internal/repository/assessment_repository.go`

**Propósito e Funcionalidade:**

Este arquivo implementa a interface `AssessmentRepository`, fornecendo métodos concretos para interagir com a tabela `assessments` e tabelas relacionadas (como `grades`, e indiretamente `students` para obter informações para cálculo de médias) no banco de dados. Ele lida com todas as operações de Create, Read, Update, Delete (CRUD) para avaliações e notas.

**Estruturas de Dados e Funções Chave:**

*   `assessmentRepository struct`: Contém um campo `db *sql.DB` para a conexão com o banco de dados.
*   `NewAssessmentRepository(db *sql.DB) AssessmentRepository`: Construtor que retorna uma nova instância de `assessmentRepository`.
*   `CreateAssessment(ctx context.Context, assessment *models.Assessment) (int64, error)`: Insere uma nova avaliação no banco de dados. Retorna o ID da avaliação criada.
*   `GetAssessmentByID(ctx context.Context, assessmentID int64) (*models.Assessment, error)`: Busca uma avaliação pelo seu ID.
*   `GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error)`: Busca todos os alunos ativos de uma determinada turma. Usado, por exemplo, para listar alunos ao lançar notas.
*   `EnterGrade(ctx context.Context, grade *models.Grade) error`: Insere ou atualiza uma nota para um aluno em uma avaliação. Usa `ON CONFLICT DO UPDATE` para lidar com casos onde a nota já existe.
*   `GetGradesByClassID(ctx context.Context, classID int64) ([]models.Grade, []models.Assessment, []models.Student, error)`: Uma função complexa que busca todas as notas, todas as avaliações e todos os alunos de uma determinada turma. Esses dados são usados pelo serviço para calcular a média da turma.
*   `ListAllAssessments(ctx context.Context) ([]models.Assessment, error)`: Retorna uma lista de todas as avaliações cadastradas no sistema.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Implementação Direta de SQL:** Os métodos usam SQL raw para interagir com o banco de dados. Isso dá controle total sobre as queries, mas requer cuidado com a segurança (ex: evitar SQL injection, embora o uso de placeholders `?` com `ExecContext`/`QueryRowContext` mitigue isso).
2.  **Context Propagation:** Todos os métodos aceitam `context.Context` como primeiro argumento, o que é uma boa prática para cancelamento e timeouts.
3.  **Tratamento de Erros Específico:** Erros como `sql.ErrNoRows` são verificados para fornecer mensagens de erro mais semânticas (ex: "no assessment found").
4.  **Gerenciamento de Nulos:** Usa `sql.NullString`, `sql.NullInt64`, `sql.NullTime` quando apropriado para lidar com colunas que podem ser NULAS no banco de dados e mapeá-las corretamente para os ponteiros ou tipos básicos nos `models`.
5.  **Eficiência em `GetGradesByClassID`:** Este método busca múltiplos conjuntos de dados (notas, avaliações, alunos) que são necessários para calcular médias. Embora possa parecer que faz várias queries, elas são direcionadas e os resultados são combinados na camada de serviço.

**Como se Encaixa no Projeto:**

Este repositório é a ponte entre a lógica de negócios relacionada a avaliações e notas (na camada de serviço) e o banco de dados. Ele abstrai os detalhes da execução de SQL, permitindo que os serviços operem em termos de objetos de modelo (`models.Assessment`, `models.Grade`).

## `internal/repository/class_repository.go`

**Propósito e Funcionalidade:**

Este arquivo implementa a interface `ClassRepository`. Ele fornece métodos para realizar operações CRUD (Create, Read, Update, Delete) para turmas (`classes`) e alunos (`students`) no banco de dados.

**Estruturas de Dados e Funções Chave:**

*   `classRepository struct`: Contém um campo `db *sql.DB`.
*   `NewClassRepository(db *sql.DB) ClassRepository`: Construtor.
*   **Operações de Turma (Class):**
    *   `CreateClass(ctx context.Context, class *models.Class) (int64, error)`: Insere uma nova turma, definindo `created_at` e `updated_at`.
    *   `GetClassByID(ctx context.Context, id int64) (*models.Class, error)`: Busca uma turma pelo ID.
    *   `UpdateClass(ctx context.Context, class *models.Class) error`: Atualiza os dados de uma turma (nome, subject_id) e `updated_at`. Verifica se a turma pertence ao usuário (com base no `user_id` no `models.Class` passado).
    *   `DeleteClass(ctx context.Context, classID int64, userID int64) error`: Deleta uma turma. Requer `userID` para verificar a propriedade.
    *   `ListAllClasses(ctx context.Context) ([]models.Class, error)`: Lista todas as turmas, ordenadas pelo nome. Inclui logging detalhado das etapas da query.
*   **Operações de Aluno (Student):**
    *   `AddStudent(ctx context.Context, student *models.Student) (int64, error)`: Adiciona um novo aluno a uma turma, definindo `created_at` e `updated_at`. Lida com `enrollment_id` opcional.
    *   `GetStudentByID(ctx context.Context, studentID int64) (*models.Student, error)`: Busca um aluno pelo ID.
    *   `UpdateStudent(ctx context.Context, student *models.Student) error`: Atualiza os dados de um aluno (nome, enrollment_id, status) e `updated_at`. A query também verifica `class_id` para garantir que o aluno não seja movido para outra turma por este método.
    *   `UpdateStudentStatus(ctx context.Context, studentID int64, status string) error`: Atualiza especificamente o status de um aluno.
    *   `DeleteStudent(ctx context.Context, studentID int64, classID int64) error`: Deleta um aluno. Requer `classID` para uma verificação de consistência/propriedade na query.
    *   `GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error)`: Lista todos os alunos de uma turma específica, ordenados pelo nome completo.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Responsabilidade Separada:** Mantém a lógica de acesso a dados para turmas e alunos dentro de um único repositório, pois estão intimamente relacionados.
2.  **Verificação de Propriedade (UserID):** Alguns métodos (`UpdateClass`, `DeleteClass`) incluem `userID` em suas queries ou verificações para garantir que as operações sejam realizadas pelo proprietário do registro, uma forma de controle de acesso no nível do repositório.
3.  **Timestamps Automáticos:** `CreatedAt` e `UpdatedAt` são gerenciados pelos métodos de criação e atualização.
4.  **Logging Detalhado em `ListAllClasses`:** A função `ListAllClasses` possui `log.Printf` statements que detalham o processo de execução da query, o que é útil para depuração.
5.  **Tratamento de Nulos para `enrollment_id`:** Usa `sql.NullString` para o campo `enrollment_id` do aluno, que é opcional.

**Como se Encaixa no Projeto:**

O `ClassRepository` é essencial para as funcionalidades de gerenciamento de turmas e alunos. Ele é consumido pelo `ClassService` para executar a lógica de negócios e, em última instância, pela TUI para apresentar e manipular esses dados.

## `internal/repository/question_repository.go`

**Propósito e Funcionalidade:**

Este arquivo implementa a interface `QuestionRepository`, fornecendo métodos para interagir com a tabela `questions` no banco de dados. Ele lida com a adição e busca de questões com base em vários critérios.

**Estruturas de Dados e Funções Chave:**

*   `questionRepository struct`: Contém um campo `db *sql.DB`.
*   `NewQuestionRepository(db *sql.DB) QuestionRepository`: Construtor.
*   `AddQuestion(ctx context.Context, question *models.Question) (int64, error)`: Insere uma nova questão no banco. Lida com o campo `options` (que é um JSON string) usando `sql.NullString`.
*   `GetQuestionsByCriteria(ctx context.Context, criteria QuestionQueryCriteria) ([]models.Question, error)`: Busca questões com base em `SubjectID`, `Topic` (opcional) e `Difficulty`. Permite limitar o número de resultados e ordena aleatoriamente (`ORDER BY RANDOM()`).
*   `GetQuestionsByCriteriaProofGeneration(ctx context.Context, criteria ProofCriteria) ([]models.Question, error)`: Uma função mais especializada para a geração de provas. Ela busca um número específico de questões para cada nível de dificuldade (`EasyCount`, `MediumCount`, `HardCount`) dentro de uma transação. Também filtra por `SubjectID` e `Topic` (opcional).
*   Funções helper (privadas, não exportadas, mas conceituais no contexto do código): `marshalOptions` e `unmarshalOptions`. Embora não estejam explicitamente no arquivo, a lógica de lidar com o campo `options` (que é um `*string` no modelo, representando um JSON) implica tal marshalling/unmarshalling ou tratamento cuidadoso de strings JSON. O código atual trata `options` como `sql.NullString` e atribui `options.String` diretamente a `q.Options` (que é `*string`).

**Decisões de Arquitetura e Funcionalidade:**

1.  **Busca Flexível de Questões:** `GetQuestionsByCriteria` oferece uma maneira flexível de consultar questões, útil para diferentes cenários de filtragem. A ordenação aleatória é útil para variar a seleção de questões.
2.  **Busca Otimizada para Geração de Provas:** `GetQuestionsByCriteriaProofGeneration` é projetada para atender eficientemente aos requisitos da geração de provas, buscando o número exato de questões por dificuldade dentro de uma transação para garantir consistência.
3.  **Tratamento de JSON em `options`:** O campo `options` é armazenado como uma string JSON no banco de dados (tipo TEXT). O repositório lida com isso usando `sql.NullString` e atribui a string diretamente ao campo `*string` no modelo `models.Question`. A serialização/desserialização para um slice de strings (`[]string`) ocorreria na camada de serviço ou na TUI, se necessário.
4.  **Uso de Transações:** `GetQuestionsByCriteriaProofGeneration` usa uma transação (`BeginTx`) para garantir que todas as buscas por diferentes níveis de dificuldade sejam tratadas atomicamente.

**Como se Encaixa no Projeto:**

O `QuestionRepository` é a base para as funcionalidades de banco de questões e geração de provas. Ele é usado pelo `QuestionService` e `ProofService` para buscar e adicionar questões, permitindo que a aplicação construa avaliações dinâmicas.

## `internal/repository/repository.go`

**Propósito e Funcionalidade:**

Este arquivo define as **interfaces** para todos os repositórios na camada de acesso a dados. Essas interfaces estabelecem os contratos que as implementações concretas dos repositórios (como `taskRepository.go`, `classRepository.go`, etc.) devem seguir. Usar interfaces aqui é uma prática fundamental de design que promove o desacoplamento e a testabilidade.

**Estruturas de Dados e Funções Chave (Interfaces):**

*   `QuestionQueryCriteria struct`: Define os critérios para buscar questões (usado por `QuestionRepository`).
    *   `SubjectID`, `Topic *string`, `Difficulty`, `Limit`.
*   `ProofCriteria struct`: Define os critérios para buscar questões para geração de provas (usado por `QuestionRepository`).
    *   `SubjectID`, `Topic *string`, `EasyCount`, `MediumCount`, `HardCount`.
*   `QuestionRepository interface`:
    *   `GetQuestionsByCriteria(...)`
    *   `AddQuestion(...)`
    *   `GetQuestionsByCriteriaProofGeneration(...)`
*   `SubjectRepository interface`:
    *   `GetOrCreateByNameAndUser(...)`
*   `TaskRepository interface`:
    *   `CreateTask(...)`, `GetTaskByID(...)`, `GetTasksByClassID(...)`, `GetAllTasks(...)`, `MarkTaskCompleted(...)`, `UpdateTask(...)`, `DeleteTask(...)`.
*   `ClassRepository interface`: (A diretiva `go:generate mockgen` indica que mocks são gerados para esta interface)
    *   Métodos para CRUD de turmas e alunos: `CreateClass`, `GetClassByID`, `AddStudent`, `UpdateStudentStatus`, `ListAllClasses`, `GetStudentsByClassID`, `UpdateClass`, `DeleteClass`, `GetStudentByID`, `UpdateStudent`, `DeleteStudent`.
*   `AssessmentRepository interface`:
    *   `CreateAssessment(...)`, `GetAssessmentByID(...)`, `EnterGrade(...)`, `GetGradesByClassID(...)`, `ListAllAssessments(...)`.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Princípio de Inversão de Dependência:** As camadas de serviço dependem dessas interfaces de repositório, não das implementações concretas. Isso permite que as implementações dos repositórios sejam trocadas (ex: de um repositório SQL real para um mock em testes) sem alterar a camada de serviço.
2.  **Contratos Claros:** As interfaces definem claramente quais operações de dados estão disponíveis para cada entidade, servindo como documentação e garantindo consistência.
3.  **Testabilidade:** Facilita o teste unitário dos serviços, pois os repositórios podem ser facilmente mockados (como indicado pelo `go:generate mockgen` para `ClassRepository`).
4.  **Modularidade:** Separa a definição da funcionalidade de acesso a dados de sua implementação.

**Como se Encaixa no Projeto:**

`repository.go` é um arquivo central para a arquitetura da aplicação. Ele define o "o quê" da camada de acesso a dados, enquanto os outros arquivos no mesmo diretório (como `task_repository.go`) fornecem o "como". Os serviços interagem exclusivamente com estas interfaces.

## `internal/repository/subject_repository.go`

**Propósito e Funcionalidade:**

Este arquivo implementa a interface `SubjectRepository`. Atualmente, ele fornece um método principal: `GetOrCreateByNameAndUser`, que busca uma disciplina pelo nome e ID do usuário, criando-a se não existir.

**Estruturas de Dados e Funções Chave:**

*   `subjectRepository struct`: Contém um campo `db *sql.DB`.
*   `NewSubjectRepository(db *sql.DB) SubjectRepository`: Construtor.
*   `GetOrCreateByNameAndUser(ctx context.Context, name string, userID int64) (models.Subject, error)`:
    *   Primeiro, tenta buscar uma disciplina com o nome e `userID` fornecidos.
    *   Se encontrada, retorna a disciplina.
    *   Se ocorrer um erro `sql.ErrNoRows` (disciplina não encontrada), ela cria uma nova disciplina com o nome e `userID` fornecidos.
    *   Se ocorrer qualquer outro erro durante a busca, retorna o erro.
    *   Retorna a disciplina encontrada ou recém-criada.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Lógica "Get or Create":** O método `GetOrCreateByNameAndUser` é uma operação comum e útil, especialmente ao importar dados ou criar entidades que dependem de outras que podem ou não já existir (como questões que dependem de disciplinas). Isso evita a necessidade de verificar a existência e depois criar em duas etapas separadas na camada de serviço.
2.  **Consistência de Dados:** Ajuda a evitar a criação de disciplinas duplicadas para o mesmo usuário com o mesmo nome.
3.  **Simplicidade Atual:** O repositório é atualmente simples, focado nesta única funcionalidade. Poderia ser expandido com métodos CRUD completos para disciplinas se necessário.

**Como se Encaixa no Projeto:**

O `SubjectRepository` é usado principalmente pelo `QuestionService` quando se adicionam questões a partir de um JSON. Se o JSON especifica nomes de disciplinas em vez de IDs, este repositório pode ser usado para encontrar o `SubjectID` correspondente ou criar uma nova disciplina se necessário, garantindo que as questões sejam associadas corretamente.

## `internal/repository/task_repository.go`

**Propósito e Funcionalidade:**

Este arquivo implementa a interface `TaskRepository`, fornecendo métodos para interagir com a tabela `tasks` no banco de dados. Ele lida com todas as operações CRUD para tarefas.

**Estruturas de Dados e Funções Chave:**

*   `taskRepository struct`: Contém um campo `db *sql.DB`.
*   `NewTaskRepository(db *sql.DB) TaskRepository`: Construtor.
*   `CreateTask(ctx context.Context, task *models.Task) (int64, error)`: Insere uma nova tarefa. Lida com `ClassID` e `DueDate` opcionais (nulos). Os campos `created_at` e `updated_at` não estão presentes no schema da tabela `tasks` (Artefacto 3.2) e, portanto, não são definidos aqui.
*   `GetTaskByID(ctx context.Context, id int64) (*models.Task, error)`: Busca uma tarefa pelo ID.
*   `GetTasksByClassID(ctx context.Context, classID int64) ([]models.Task, error)`: Busca todas as tarefas associadas a um `ClassID` específico.
*   `GetAllTasks(ctx context.Context) ([]models.Task, error)`: Busca todas as tarefas do banco de dados.
*   `MarkTaskCompleted(ctx context.Context, taskID int64) error`: Define `is_completed = true` para uma tarefa. Não define `updated_at` pois não existe na tabela.
*   `DeleteTask(ctx context.Context, taskID int64) error`: Deleta uma tarefa pelo ID.
*   `UpdateTask(ctx context.Context, task *models.Task) error`: Atualiza todos os campos de uma tarefa existente. Não define `updated_at`.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Schema da DB (Tasks):** A implementação reflete o schema da tabela `tasks` do Artefacto 3.2, que não inclui `created_at` ou `updated_at`. Se esses campos fossem adicionados ao schema, o repositório precisaria ser atualizado para gerenciá-los.
2.  **Tratamento de Campos Nulos:** Usa `sql.NullInt64`, `sql.NullString`, `sql.NullTime` para mapear corretamente os campos opcionais (`ClassID`, `Description`, `DueDate`) de/para o banco de dados.
3.  **Operações CRUD Completas:** Fornece um conjunto completo de métodos para gerenciar o ciclo de vida das tarefas.
4.  **Feedback de `RowsAffected`:** Muitos métodos de atualização/exclusão verificam `RowsAffected` para garantir que a operação teve o efeito esperado (ex: que um registro foi realmente encontrado e modificado/deletado).

**Como se Encaixa no Projeto:**

O `TaskRepository` é a camada de persistência para a funcionalidade de gerenciamento de tarefas. Ele é usado pelo `TaskService` para executar a lógica de negócios e interagir com o banco de dados. É fundamental para que os usuários possam criar, listar, atualizar e concluir suas tarefas.
## `internal/service/assessment_service.go`

**Propósito e Funcionalidade:**

Este arquivo implementa a interface `AssessmentService`, definindo a lógica de negócios para o gerenciamento de avaliações e notas. Ele atua como um intermediário entre a camada de apresentação (CLI/TUI) e a camada de repositório (`AssessmentRepository`, `ClassRepository`).

**Estruturas de Dados e Funções Chave:**

*   `assessmentServiceImpl struct`: Contém instâncias de `AssessmentRepository` e `ClassRepository`.
*   `NewAssessmentService(assessmentRepo repository.AssessmentRepository, classRepo repository.ClassRepository) AssessmentService`: Construtor que injeta os repositórios necessários.
*   `CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error)`:
    *   Valida os parâmetros de entrada (nome não vazio, IDs e valores positivos).
    *   Cria uma struct `models.Assessment`.
    *   Chama `assessmentRepo.CreateAssessment` para persistir a avaliação.
    *   Retorna a avaliação criada com seu ID.
*   `EnterGrades(ctx context.Context, assessmentID int64, studentGrades map[int64]float64) error`:
    *   Valida `assessmentID` e se `studentGrades` não está vazio.
    *   Opcionalmente (e recomendado), valida se `assessmentID` existe e se os `studentID`s são válidos para a turma da avaliação (usando `assessmentRepo.GetAssessmentByID` e `classRepo.GetStudentsByClassID`). A implementação atual faz a validação do `assessmentID`.
    *   Itera sobre `studentGrades`, criando uma struct `models.Grade` para cada entrada.
    *   Chama `assessmentRepo.EnterGrade` para persistir cada nota.
*   `CalculateClassAverage(ctx context.Context, classID int64) (float64, error)`:
    *   Valida `classID`.
    *   Chama `assessmentRepo.GetGradesByClassID` para obter todas as notas, todas as avaliações e todos os alunos da turma.
    *   Realiza a lógica de cálculo da média ponderada:
        *   Para cada aluno ativo, calcula sua média individual considerando o peso de cada avaliação em que ele tem nota.
        *   Calcula a média geral da turma somando as médias individuais dos alunos ativos e dividindo pelo número de alunos ativos.
    *   Lida com casos onde não há alunos ou avaliações.
*   `ListAllAssessments(ctx context.Context) ([]models.Assessment, error)`:
    *   Chama `assessmentRepo.ListAllAssessments` para buscar todas as avaliações.
    *   Retorna a lista de avaliações.
*   `findStudent(students []models.Student, studentID int64) (models.Student, bool)`: Função helper (privada) para buscar um aluno em uma slice pelo ID.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Lógica de Negócios Clara:** Separa a lógica de validação, orquestração de chamadas ao repositório e cálculos (como a média da turma) da simples persistência de dados.
2.  **Injeção de Dependência:** Recebe instâncias de repositório, o que facilita testes e desacoplamento.
3.  **Validação de Entradas:** Realiza validações básicas nos parâmetros de entrada dos métodos.
4.  **Cálculo de Média Complexo:** A função `CalculateClassAverage` implementa a lógica de cálculo de média ponderada, que pode ser não trivial, demonstrando o papel do serviço em realizar essas operações.
5.  **Tratamento de Alunos Ativos:** No cálculo da média, considera apenas alunos com status "ativo".

**Como se Encaixa no Projeto:**

O `AssessmentService` é crucial para as funcionalidades de avaliação do Vigenda. Ele é consumido pela CLI (`cmd/vigenda/main.go`) e pela TUI (`internal/app/assessments/model.go`) para executar ações como criar avaliações, lançar notas e visualizar médias. Ele garante que as regras de negócio sejam aplicadas antes de interagir com o banco de dados.

## `internal/service/assessment_service_test.go`

**Propósito e Funcionalidade:**

Este arquivo destina-se a conter testes unitários para `AssessmentService`. No entanto, o conteúdo fornecido é um esqueleto com `// TODO: Implement test` para as principais funções.

**Estruturas de Dados e Funções Chave:**

*   `MockAssessmentRepository struct`: Define uma estrutura de mock para `AssessmentRepository` (embora a convenção seja usar mocks gerados ou de bibliotecas como `testify/mock`).
    *   Contém campos de função (ex: `CreateAssessmentFunc`) que podem ser sobrescritos nos testes para simular o comportamento do repositório.
*   Funções de teste (esqueletos): `TestCreateAssessment`, `TestEnterGrades`, `TestCalculateClassAverage`.

**Decisões de Arquitetura e Funcionalidade (nos Testes):**

1.  **Testes Unitários Planejados:** A existência do arquivo e dos esqueletos de teste indica a intenção de testar unitariamente a camada de serviço.
2.  **Mocking de Repositório:** A `MockAssessmentRepository` demonstra a abordagem de usar mocks para isolar o serviço de dependências de banco de dados reais durante os testes.

**Como se Encaixa no Projeto:**

Uma vez implementados, esses testes seriam vitais para garantir a corretude da lógica de negócios em `AssessmentService`. Eles verificariam se as validações funcionam, se os repositórios são chamados corretamente e se os cálculos (como médias) são precisos.

## `internal/service/class_service.go`

**Propósito e Funcionalidade:**

Este arquivo implementa a interface `ClassService`, definindo a lógica de negócios para o gerenciamento de turmas e alunos. Ele coordena as interações com `ClassRepository` e, potencialmente, `SubjectRepository`.

**Estruturas de Dados e Funções Chave:**

*   `classServiceImpl struct`: Contém instâncias de `ClassRepository` e `SubjectRepository`.
*   `NewClassService(classRepo repository.ClassRepository, subjectRepo repository.SubjectRepository) ClassService`: Construtor.
*   **Operações de Turma:**
    *   `CreateClass(ctx context.Context, name string, subjectID int64) (models.Class, error)`: Valida entradas, assume um `UserID` placeholder (1), cria a struct `models.Class`, chama `classRepo.CreateClass`, e então `classRepo.GetClassByID` para retornar a turma completa com timestamps.
    *   `UpdateClass(ctx context.Context, classID int64, name string, subjectID int64) (models.Class, error)`: Valida entradas, busca a turma original para verificar propriedade (comparando `UserID`), atualiza os campos, chama `classRepo.UpdateClass`, e busca novamente para retornar a turma atualizada.
    *   `DeleteClass(ctx context.Context, classID int64) error`: Valida `classID`, busca a turma (para verificação de existência/propriedade), e chama `classRepo.DeleteClass`.
    *   `GetClassByID(ctx context.Context, classID int64) (models.Class, error)`: Valida `classID`, chama `classRepo.GetClassByID`. (A verificação de `UserID` comentada sugere um controle de acesso futuro ou que ele é feito no repositório).
    *   `ListAllClasses(ctx context.Context) ([]models.Class, error)`: Chama `classRepo.ListAllClasses`. Inclui logging detalhado.
*   **Operações de Aluno:**
    *   `AddStudent(ctx context.Context, classID int64, fullName string, enrollmentID string, status string) (models.Student, error)`: Valida entradas (default 'ativo' para status), cria `models.Student`, chama `classRepo.AddStudent`, e busca o aluno recém-criado para retornar o objeto completo.
    *   `GetStudentByID(ctx context.Context, studentID int64) (models.Student, error)`: Valida `studentID`, chama `classRepo.GetStudentByID`.
    *   `UpdateStudent(ctx context.Context, studentID int64, fullName string, enrollmentID string, status string) (models.Student, error)`: Valida entradas, busca o aluno original, atualiza seus campos, chama `classRepo.UpdateStudent`, e busca novamente para retornar o aluno atualizado.
    *   `DeleteStudent(ctx context.Context, studentID int64) error`: Valida `studentID`, busca o aluno (para obter seu `ClassID`, que é usado na chamada `classRepo.DeleteStudent`), e chama `classRepo.DeleteStudent`.
    *   `UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error`: Valida `studentID` e `newStatus`, chama `classRepo.UpdateStudentStatus`.
    *   `GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error)`: Valida `classID`, chama `classRepo.GetStudentsByClassID`.
*   `ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error)`:
    *   Valida `classID`.
    *   Parseia os dados CSV (esperando colunas: `enrollment_id`, `full_name`, `status` opcional).
    *   Para cada registro válido no CSV, chama o método `AddStudent` do próprio serviço para adicionar o aluno.
    *   Loga e continua em caso de erro ao adicionar um aluno individual do CSV.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Validação de Parâmetros:** A maioria dos métodos começa com validações dos argumentos recebidos.
2.  **UserID Placeholder:** Consistentemente usa `UserID = 1` como placeholder, indicando que a integração com um sistema de autenticação real para obter o `UserID` do contexto é um trabalho futuro.
3.  **Operações CRUD e Lógica de Negócios:** Fornece uma camada de abstração sobre o repositório, adicionando validações e, em alguns casos, orquestrando múltiplas chamadas ao repositório (ex: `CreateClass` que cria e depois busca).
4.  **Importação CSV Robusta (Parcialmente):** A importação de CSV tenta continuar mesmo se registros individuais falharem, logando os erros.
5.  **Retorno de Objetos Completos:** Métodos de criação e atualização geralmente buscam o objeto recém-criado/atualizado do repositório para retornar um `models.Class` ou `models.Student` completo, incluindo IDs gerados e timestamps.

**Como se Encaixa no Projeto:**

`ClassService` é uma camada de serviço vital que gerencia a lógica de negócios para turmas e alunos. É usado pela TUI (`internal/app/classes/model.go`) e pela CLI (`cmd/vigenda/main.go`) para todas as operações relacionadas a essas entidades.

## `internal/service/class_service_test.go`

**Propósito e Funcionalidade:**

Este arquivo contém testes unitários para `ClassService`, utilizando mocks gerados (via `go.uber.org/mock/gomock`) para o `ClassRepository`. O objetivo é testar a lógica dentro do `ClassService` isoladamente.

**Estruturas de Dados e Funções Chave:**

*   Usa `stubs.NewMockClassRepository(ctrl)` para criar um mock do `ClassRepository`.
*   **Test Functions:** Cada função testa um método público do `ClassService`:
    *   `TestClassServiceImpl_CreateClass`: Verifica se `CreateClass` chama `classRepo.CreateClass` e `classRepo.GetClassByID` com os argumentos corretos e retorna a classe esperada.
    *   `TestClassServiceImpl_UpdateClass`: Testa a lógica de atualização, incluindo as chamadas a `GetClassByID` (para buscar o original e o atualizado) e `UpdateClass`.
    *   `TestClassServiceImpl_DeleteClass`: Verifica se `GetClassByID` (para verificação de existência) e `DeleteClass` são chamados.
    *   `TestClassServiceImpl_AddStudent`: Testa a adição de aluno, incluindo chamadas a `AddStudent` e `GetStudentByID`.
    *   `TestClassServiceImpl_UpdateStudent`: Semelhante a `UpdateClass`, mas para alunos.
    *   `TestClassServiceImpl_DeleteStudent`: Testa a exclusão de aluno, verificando a chamada a `GetStudentByID` (para obter `ClassID`) e `DeleteStudent`.
    *   `TestClassServiceImpl_ListAllClasses`: Verifica se `ListAllClasses` do repositório é chamado e os resultados são repassados.
    *   `TestClassServiceImpl_GetStudentsByClassID`: Semelhante, para `GetStudentsByClassID`.
    *   `TestUpdateStudentStatus`: Verifica a chamada a `UpdateStudentStatus` no repositório.
*   **Uso de `gomock`:**
    *   `ctrl := gomock.NewController(t)`: Cria um controlador gomock.
    *   `mockClassRepo.EXPECT().MethodName(...).Return(...).Times(1)`: Define as expectativas para as chamadas aos métodos do repositório mockado, incluindo os argumentos esperados e os valores de retorno.
    *   `gomock.Any()`: Usado para argumentos cujo valor exato não é crítico para o teste (como `context.Context` ou ponteiros para structs onde apenas alguns campos são relevantes).
    *   `.DoAndReturn(...)`: Permite fornecer uma função customizada para simular o comportamento do método mockado e verificar os argumentos passados.

**Decisões de Arquitetura e Funcionalidade (nos Testes):**

1.  **Testes Baseados em Mock:** Os testes dependem fortemente de mocks para `ClassRepository`, permitindo focar na lógica do serviço (validações, orquestração de chamadas) sem depender do banco de dados.
2.  **Verificação de Interações:** Os testes não apenas verificam o resultado final, mas também se os métodos corretos do repositório foram chamados com os argumentos esperados (`AssertExpectations`).
3.  **Cobertura de Métodos:** O objetivo é ter uma função de teste para cada método público do `ClassService`.
4.  **UserID Fixo em Testes:** Os testes usam um `UserID` fixo (geralmente 1) ao definir expectativas, alinhando-se com o placeholder usado no `ClassService`.

**Como se Encaixa no Projeto:**

Estes testes são cruciais para garantir a confiabilidade da lógica de negócios de gerenciamento de turmas e alunos. Eles ajudam a detectar bugs e regressões no `ClassService` à medida que o código evolui.

## `internal/service/proof_service.go`

**Propósito e Funcionalidade:**

Este arquivo implementa a interface `ProofService`, que é responsável pela lógica de negócios de geração de provas. Ele usa o `QuestionRepository` para buscar questões com base em critérios específicos.

**Estruturas de Dados e Funções Chave:**

*   `proofServiceImpl struct`: Contém uma instância de `QuestionRepository`.
*   `NewProofService(qr repository.QuestionRepository) ProofService`: Construtor.
*   `GenerateProof(ctx context.Context, criteria ProofCriteria) ([]models.Question, error)`:
    *   Valida os critérios de entrada (`SubjectID` não zero, contagens de dificuldade não negativas, pelo menos uma contagem > 0).
    *   Converte `service.ProofCriteria` para `repository.ProofCriteria`.
    *   Chama `questionRepo.GetQuestionsByCriteriaProofGeneration` para buscar todas as questões necessárias (para todas as dificuldades) em uma única chamada ao repositório.
    *   Verifica se o número de questões retornadas para cada nível de dificuldade corresponde ao solicitado. Se não, retorna um erro indicando qual dificuldade não teve questões suficientes.
    *   Retorna a lista de questões que compõem a prova. A ordem é determinada pelo repositório (que usa `ORDER BY RANDOM()`).

**Decisões de Arquitetura e Funcionalidade:**

1.  **Validação de Critérios no Serviço:** O serviço realiza validações nos critérios de geração da prova antes de prosseguir.
2.  **Delegação ao Repositório Especializado:** Utiliza o método `GetQuestionsByCriteriaProofGeneration` do repositório, que é otimizado para buscar o número exato de questões por dificuldade, possivelmente dentro de uma transação.
3.  **Verificação de Suficiência de Questões:** O serviço verifica se o repositório retornou o número esperado de questões para cada nível de dificuldade e retorna um erro se não for o caso. Isso garante que a prova gerada atenda aos requisitos.
4.  **Ordem Aleatória:** A aleatoriedade da seleção de questões é delegada ao repositório.

**Como se Encaixa no Projeto:**

`ProofService` é uma parte essencial da funcionalidade de geração de provas. Ele é usado pela CLI (`cmd/vigenda/main.go`) e pela TUI (`internal/app/proofs/model.go`) para criar provas com base nos critérios fornecidos pelo usuário.

## `internal/service/proof_service_test.go`

**Propósito e Funcionalidade:**

Este arquivo contém testes unitários para `ProofService`, usando um mock para `QuestionRepository` para testar a lógica de geração de provas isoladamente.

**Estruturas de Dados e Funções Chave:**

*   `MockQuestionRepository struct`: Define um mock para `QuestionRepository` usando `testify/mock`. (Nota: Este mock é definido localmente neste arquivo de teste. Se `QuestionService` também o usasse, seria melhor ter um mock compartilhado.)
    *   Implementa os métodos da interface `QuestionRepository`, permitindo definir expectativas para as chamadas.
*   **Test Functions:**
    *   `TestProofService_GenerateProof`: Contém subtestes para diferentes cenários:
        *   `"success"`: Testa o caso de sucesso onde questões suficientes são encontradas para todas as dificuldades. Verifica se `questionRepo.GetQuestionsByCriteriaProofGeneration` é chamado corretamente e se as questões retornadas estão corretas.
        *   `"error_no_difficulty_count"`: Testa a validação de que pelo menos uma contagem de dificuldade deve ser > 0.
        *   `"error_fetching_easy_questions"` (e similares para medium/hard): Testa o cenário onde o repositório retorna um erro ao tentar buscar questões.
        *   `"error_not_enough_easy_questions"` (e similares para medium/hard): Testa o cenário onde o repositório retorna menos questões do que o solicitado para uma dificuldade específica, e o serviço identifica e reporta esse erro.
        *   `"success_only_medium_questions"`: Testa um caso onde apenas questões de uma dificuldade são solicitadas.
*   **Uso de `testify/mock` e `testify/assert`:**
    *   `mockRepo.On("MethodName", ...).Return(...).Once()`: Configura as expectativas para as chamadas ao repositório mockado.
    *   `assert.NoError(t, err)`, `assert.Error(t, err)`, `assert.EqualError(t, err, "...")`, `assert.Len(t, questions, ...)`, `assert.Contains(t, questions, ...)`: Usados para verificar os resultados e erros.

**Decisões de Arquitetura e Funcionalidade (nos Testes):**

1.  **Testes de Cenário:** Cobre vários cenários, incluindo sucesso, erros de validação de entrada, erros do repositório e casos onde não há questões suficientes.
2.  **Mock Detalhado para `GetQuestionsByCriteriaProofGeneration`:** Os testes configuram o mock para retornar diferentes conjuntos de questões ou erros com base nos critérios passados, permitindo testar a lógica de tratamento de respostas do serviço.
3.  **Foco na Lógica do Serviço:** Os testes validam que o `ProofService` interpreta corretamente os resultados do repositório e aplica sua própria lógica de validação (como verificar se o número de questões é suficiente).

**Como se Encaixa no Projeto:**

Esses testes garantem que a funcionalidade de geração de provas seja confiável, que lide corretamente com diferentes respostas do banco de dados (via mock do repositório) e que os critérios sejam validados adequadamente.

## `internal/service/question_service.go`

**Propósito e Funcionalidade:**

Este arquivo implementa a interface `QuestionService`, fornecendo a lógica de negócios para o gerenciamento do banco de questões, incluindo a adição de questões a partir de JSON e a geração de "testes" (que parece ser uma funcionalidade similar ou precursora da geração de provas).

**Estruturas de Dados e Funções Chave:**

*   `questionServiceImpl struct`: Contém instâncias de `QuestionRepository` e `SubjectRepository`.
*   `NewQuestionService(qr repository.QuestionRepository, sr repository.SubjectRepository) QuestionService`: Construtor.
*   `AddQuestionsFromJSON(ctx context.Context, jsonData []byte) (int, error)`:
    *   Desserializa o `jsonData` para um slice de uma struct anônima que representa o formato esperado das questões no JSON (incluindo `SubjectName` em vez de `SubjectID`).
    *   Valida campos obrigatórios em cada questão do JSON.
    *   Para cada questão:
        *   Usa `subjectRepo.GetOrCreateByNameAndUser` para obter/criar o `SubjectID` com base no `SubjectName` e `UserID` do JSON. (Se `subjectRepo` for nil, retorna um erro).
        *   Se o tipo for "multipla_escolha", garante que `Options` não seja nulo/vazio e serializa as opções para uma string JSON.
        *   Cria uma `models.Question` e a preenche com os dados processados (incluindo o `UserID` do JSON).
        *   Chama `questionRepo.AddQuestion` para cada questão.
    *   Retorna o número de questões adicionadas com sucesso.
*   `GenerateTest(ctx context.Context, criteria TestCriteria) ([]models.Question, error)`:
    *   Valida que pelo menos uma contagem de dificuldade seja > 0.
    *   Para cada nível de dificuldade com contagem > 0 (Easy, Medium, Hard):
        *   Chama `questionRepo.GetQuestionsByCriteria` para buscar o número solicitado de questões daquela dificuldade, filtrando por `SubjectID` e `Topic` (opcional).
        *   Se não houver questões suficientes para um nível de dificuldade, retorna um erro.
    *   Agrega as questões de todas as dificuldades e as retorna.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Processamento de JSON Flexível:** `AddQuestionsFromJSON` lida com um formato JSON que usa nomes de disciplinas (`SubjectName`) e os converte para `SubjectID` usando o `SubjectRepository`. Também lida com o campo `options` que pode ser `any` no JSON e o serializa para string.
2.  **Criação de Disciplina Implícita:** A dependência do `SubjectRepository` com `GetOrCreateByNameAndUser` permite que novas disciplinas sejam criadas automaticamente se mencionadas no JSON de importação de questões e ainda não existirem.
3.  **Validação Detalhada na Importação JSON:** Realiza várias validações nos dados do JSON para garantir a integridade das questões.
4.  **Geração de Teste por Partes:** `GenerateTest` busca questões para cada nível de dificuldade separadamente e as combina. Isso difere de `ProofService.GenerateProof` que usa um método de repositório mais especializado (`GetQuestionsByCriteriaProofGeneration`). Isso pode indicar que `GenerateTest` é uma funcionalidade mais antiga ou com requisitos ligeiramente diferentes.
5.  **UserID da Questão do JSON:** `AddQuestionsFromJSON` agora usa corretamente o `UserID` fornecido em cada objeto de questão no JSON para popular `models.Question.UserID`.

**Como se Encaixa no Projeto:**

`QuestionService` é fundamental para popular e utilizar o banco de questões. `AddQuestionsFromJSON` fornece um mecanismo crucial para importar questões em lote. `GenerateTest` (semelhante a `ProofService.GenerateProof`) permite a criação de conjuntos de questões para avaliações.

## `internal/service/question_service_test.go`

**Propósito e Funcionalidade:**

Este arquivo contém testes unitários para `QuestionService`, utilizando mocks para `QuestionRepository` e `SubjectRepository` para isolar a lógica do serviço.

**Estruturas de Dados e Funções Chave:**

*   `MockQuestionRepository` e `MockSubjectRepository`: Mocks definidos localmente (ou referenciados se estiverem no mesmo pacote de teste) usando `testify/mock`.
    *   `MockSubjectRepository` inclui `GetOrCreateByNameAndUser`.
*   **Test Functions:**
    *   `TestQuestionService_AddQuestionsFromJSON`:
        *   `"success"`: Testa a importação bem-sucedida de múltiplas questões (dissertativa e múltipla escolha). Verifica se `subjectRepo.GetOrCreateByNameAndUser` e `questionRepo.AddQuestion` são chamados corretamente.
        *   `"error_invalid_json"`: Testa o tratamento de JSON malformado.
        *   `"error_empty_json_array"`: Testa o JSON com um array vazio.
        *   `"error_missing_required_field_statement"` (e similares para outros campos): Testa a validação de campos obrigatórios.
        *   `"error_missing_options_for_multiple_choice"`: Testa a validação de opções para questões de múltipla escolha.
        *   `"error_empty_options_for_multiple_choice"`: Testa a validação de opções vazias.
        *   `"error_on_add_question_repository"`: Testa o cenário onde `questionRepo.AddQuestion` retorna um erro.
    *   `TestQuestionService_GenerateTest`: Similar aos testes de `ProofService.GenerateProof`.
        *   `"success"`: Testa a geração bem-sucedida.
        *   `"error_no_difficulty_count"`: Validação de contagens de dificuldade.
        *   `"error_fetching_easy_questions"`: Erro do repositório ao buscar questões.
        *   `"error_not_enough_medium_questions"`: Repositório retorna menos questões do que o solicitado.
*   **Uso de `testify/mock` e `testify/assert`:** Semelhante aos outros arquivos de teste de serviço.

**Decisões de Arquitetura e Funcionalidade (nos Testes):**

1.  **Foco em `AddQuestionsFromJSON`:** Dada a complexidade do processamento de JSON e interação com múltiplos repositórios (potencialmente), esta função recebe atenção detalhada nos testes.
2.  **Mock de `SubjectRepository`:** Os testes para `AddQuestionsFromJSON` agora mockam `subjectRepo.GetOrCreateByNameAndUser` para simular a conversão de nome de disciplina para ID.
3.  **Testes Abrangentes de Validação JSON:** Cobre vários cenários de erro relacionados ao formato e conteúdo do JSON de entrada.
4.  **Teste de `GenerateTest`:** Embora semelhante a `ProofService`, é importante testar esta funcionalidade separadamente, pois ela usa uma abordagem diferente para buscar questões do repositório (`GetQuestionsByCriteria` chamado múltiplas vezes).

**Como se Encaixa no Projeto:**

Esses testes asseguram que o `QuestionService` possa importar questões de JSON de forma confiável, validando os dados e interagindo corretamente com os repositórios. Também verificam a lógica de geração de testes.

## `internal/service/service.go`

**Propósito e Funcionalidade:**

Este arquivo define as **interfaces** para todos os serviços na camada de lógica de negócios. Essas interfaces estabelecem os contratos que as implementações concretas dos serviços (como `task_service.go`, `class_service.go`, etc.) devem seguir.

**Estruturas de Dados e Funções Chave (Interfaces):**

*   `TaskService interface`:
    *   Métodos para CRUD de tarefas: `CreateTask`, `ListActiveTasksByClass`, `ListAllActiveTasks`, `ListAllTasks`, `MarkTaskAsCompleted`, `GetTaskByID`, `UpdateTask`, `DeleteTask`.
*   `ClassService interface`:
    *   Métodos para CRUD de turmas e alunos: `CreateClass`, `ImportStudentsFromCSV`, `UpdateStudentStatus`, `GetClassByID`, `ListAllClasses`, `GetStudentsByClassID`, `UpdateClass`, `DeleteClass`, `AddStudent`, `GetStudentByID`, `UpdateStudent`, `DeleteStudent`.
*   `AssessmentService interface`:
    *   `CreateAssessment`, `EnterGrades`, `CalculateClassAverage`, `ListAllAssessments`.
*   `QuestionService interface`:
    *   `AddQuestionsFromJSON`, `GenerateTest`.
*   `TestCriteria struct`: Usado por `QuestionService.GenerateTest`.
    *   `SubjectID`, `Topic *string`, `EasyCount`, `MediumCount`, `HardCount`.
*   `ProofService interface`:
    *   `GenerateProof`.
*   `ProofCriteria struct`: Usado por `ProofService.GenerateProof`.
    *   `SubjectID`, `Topic *string`, `EasyCount`, `MediumCount`, `HardCount`.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Princípio de Inversão de Dependência:** A camada de apresentação (CLI/TUI) depende dessas interfaces de serviço, não das implementações concretas. Isso permite que as implementações dos serviços sejam alteradas ou mockadas sem afetar a camada de apresentação.
2.  **Contratos Claros para Lógica de Negócios:** As interfaces definem claramente quais operações de negócios estão disponíveis para cada domínio funcional.
3.  **Testabilidade:** Facilita o teste unitário da camada de apresentação, pois os serviços podem ser facilmente mockados.
4.  **Modularidade:** Separa a definição da lógica de negócios de sua implementação.

**Como se Encaixa no Projeto:**

`service.go` é um arquivo central para a arquitetura da aplicação, definindo os limites e contratos da camada de lógica de negócios. A CLI e a TUI interagem com a aplicação através dessas interfaces de serviço.

## `internal/service/stubs.go`

**Propósito e Funcionalidade:**

Este arquivo fornece implementações "stub" (simplificadas, muitas vezes com comportamento fixo ou mínimo) das interfaces de serviço. Stubs são úteis para:
*   Permitir que a aplicação compile e execute mesmo que as implementações completas dos serviços ou repositórios ainda não estejam prontas.
*   Testes básicos da CLI ou TUI onde o comportamento exato da lógica de negócios não é o foco principal, mas a conectividade entre as camadas sim.
*   Desenvolvimento inicial e prototipagem.

**Estruturas de Dados e Funções Chave:**

Para cada interface de serviço em `service.go`, há uma implementação stub correspondente:

*   `stubTaskService struct`:
    *   Contém um `repository.TaskRepository` (que também pode ser um stub).
    *   Implementa os métodos de `TaskService`, geralmente imprimindo uma mensagem de log e chamando o método correspondente no `taskRepo` stub (que pode fazer uma operação de DB simples ou nada).
*   `stubClassService struct`:
    *   Contém um `repository.ClassRepository`.
    *   Implementação similar para métodos de `ClassService`. `ImportStudentsFromCSV` tem uma lógica básica de parsing CSV.
*   `stubAssessmentService struct`:
    *   Contém um `repository.AssessmentRepository`.
    *   Métodos como `CalculateClassAverage` retornam um valor fixo.
*   `stubQuestionService struct`:
    *   Contém `repository.QuestionRepository` e `repository.SubjectRepository`.
    *   Métodos retornam valores fixos ou listas vazias.
*   `stubProofService struct`:
    *   Contém `repository.QuestionRepository`.
    *   `GenerateProof` chama `questionRepo.GetQuestionsByCriteria` (que, se o repo também for um stub, retornará dados limitados).

**Decisões de Arquitetura e Funcionalidade:**

1.  **Facilitar Desenvolvimento Incremental:** Permite que diferentes partes da aplicação sejam desenvolvidas e testadas independentemente. Por exemplo, a CLI pode ser desenvolvida usando serviços stub enquanto os serviços reais e repositórios estão em andamento.
2.  **Comportamento Mínimo Viável:** Os stubs geralmente não contêm lógica de negócios complexa, mas sim o suficiente para que os métodos possam ser chamados e retornem tipos válidos, evitando que a aplicação quebre.
3.  **Uso de Repositórios Stub (Implícito):** Os serviços stub frequentemente dependem de repositórios stub (como `repository.NewStubTaskRepository`), continuando a cadeia de stubs até a camada de dados.
4.  **Logging para Rastreamento:** Muitos métodos stub imprimem mensagens (`fmt.Printf`) indicando que foram chamados, o que ajuda a rastrear o fluxo durante testes com stubs.

**Como se Encaixa no Projeto:**

Os stubs em `stubs.go` foram provavelmente usados durante as fases iniciais de desenvolvimento para montar a estrutura da aplicação e testar a integração entre a CLI/TUI e a camada de serviço antes que as implementações completas dos serviços e repositórios estivessem prontas. Eles são substituídos pelas implementações reais (`taskServiceImpl`, etc.) em `cmd/vigenda/main.go` para a execução normal da aplicação.

## `internal/service/task_service.go`

**Propósito e Funcionalidade:**

Este arquivo implementa a interface `TaskService`, definindo a lógica de negócios para o gerenciamento de tarefas. Ele interage com o `TaskRepository` para persistência. Uma característica notável é a tentativa de criar "tarefas de bug" automaticamente quando ocorrem erros inesperados.

**Estruturas de Dados e Funções Chave:**

*   `taskServiceImpl struct`: Contém uma instância de `repository.TaskRepository`.
*   `NewTaskService(repo repository.TaskRepository) TaskService`: Construtor.
*   `logError`: Função helper interna para logar erros (atualmente para `os.Stderr`).
*   `handleErrorAndCreateBugTask`: Função helper que loga um erro original e tenta criar uma nova tarefa de bug no sistema (com `UserID = 0` e `ClassID = nil`) usando `createTaskInternal`.
*   `createTaskInternal`: Versão interna de `CreateTask` que chama diretamente o repositório, usada para evitar recursão na criação de tarefas de bug.
*   `CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error)`:
    *   Valida se o título não está vazio.
    *   Assume um `UserID` placeholder (1).
    *   Chama `createTaskInternal`. Se ocorrer um erro, chama `handleErrorAndCreateBugTask`.
*   `UpdateTask(ctx context.Context, task *models.Task) error`:
    *   Valida se o título não está vazio.
    *   Chama `repo.UpdateTask`. Se ocorrer um erro (e não for "não encontrado" ou "sem alterações"), chama `handleErrorAndCreateBugTask`.
*   `DeleteTask(ctx context.Context, taskID int64) error`:
    *   Chama `repo.DeleteTask`. Se ocorrer um erro (e não for "não encontrado"), chama `handleErrorAndCreateBugTask`.
*   `ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error)`:
    *   Chama `repo.GetTasksByClassID`.
    *   Filtra as tarefas para retornar apenas aquelas onde `IsCompleted == false`.
    *   Se erro do repositório, chama `handleErrorAndCreateBugTask`.
*   `ListAllTasks(ctx context.Context) ([]models.Task, error)`:
    *   Chama `repo.GetAllTasks`. Se erro, chama `handleErrorAndCreateBugTask`.
*   `ListAllActiveTasks(ctx context.Context) ([]models.Task, error)`:
    *   Chama `repo.GetAllTasks`.
    *   Filtra para retornar apenas tarefas não concluídas.
    *   Se erro do repositório, chama `handleErrorAndCreateBugTask`.
*   `MarkTaskAsCompleted(ctx context.Context, taskID int64) error`:
    *   Chama `repo.MarkTaskCompleted`. Se erro, chama `handleErrorAndCreateBugTask`.
*   `GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error)`:
    *   Chama `repo.GetTaskByID`.
    *   Se o erro for "não encontrado" (especificamente `sql.ErrNoRows` ou contendo "no task found"), loga e retorna um erro específico.
    *   Para outros erros, chama `handleErrorAndCreateBugTask`.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Criação Automática de Tarefas de Bug:** A tentativa de criar tarefas de bug para erros inesperados do sistema é uma abordagem interessante para auto-rastreamento de problemas. `UserID = 0` e `ClassID = nil` são usados para identificar essas tarefas de sistema.
2.  **Validação de Entrada:** Realiza validações básicas (ex: título não vazio).
3.  **UserID Fixo:** Assim como outros serviços, usa um `UserID` placeholder (1) para operações normais, indicando a necessidade de integração com autenticação.
4.  **Filtragem na Camada de Serviço:** `ListActiveTasksByClass` e `ListAllActiveTasks` realizam a filtragem de tarefas (para ativas) na camada de serviço após buscar um conjunto maior do repositório. Idealmente, se o desempenho fosse crítico, o repositório ofereceria métodos para buscar diretamente apenas tarefas ativas.
5.  **Tratamento Diferenciado de Erros "Não Encontrado":** `GetTaskByID` (e parcialmente `UpdateTask`/`DeleteTask`) trata erros de "não encontrado" de forma diferente, não criando tarefas de bug para eles, pois podem ser resultado de entrada inválida do usuário.

**Como se Encaixa no Projeto:**

`TaskService` é o principal componente para toda a lógica de negócios relacionada a tarefas. Ele é usado pela CLI e TUI. A funcionalidade de criação de tarefas de bug é uma característica distintiva que visa melhorar a manutenção e depuração do sistema.

## `internal/service/task_service_test.go`

**Propósito e Funcionalidade:**

Este arquivo contém testes unitários para `TaskService`, utilizando um mock para `TaskRepository` para isolar a lógica do serviço. Os testes verificam o comportamento do serviço, incluindo a criação de tarefas de bug em cenários de erro.

**Estruturas de Dados e Funções Chave:**

*   `MockTaskRepository struct`: Mock para `TaskRepository`.
    *   Contém campos de função para cada método da interface.
    *   `CreatedBugTasks []models.Task`: Um slice para armazenar tarefas de bug que foram "criadas" pelo mock, permitindo que os testes verifiquem se elas foram geradas corretamente.
*   **Test Functions:**
    *   `TestTaskService_CreateTask`:
        *   `"successful task creation"`: Verifica a criação bem-sucedida.
        *   `"repository error on create"`: Simula um erro do repositório na criação normal e verifica se uma tarefa de bug é criada.
        *   `"validation error (empty title)"`: Verifica se erros de validação não criam tarefas de bug.
    *   `TestTaskService_ListActiveTasksByClass`:
        *   `"successful listing"`: Verifica a listagem e filtragem corretas.
        *   `"repository error on listing"`: Simula erro do repositório e verifica a criação de tarefa de bug.
    *   `TestTaskService_ListAllTasks` e `TestTaskService_ListAllActiveTasks`: Testes semelhantes para listagem global, incluindo tratamento de erro e criação de bug.
    *   `TestTaskService_MarkTaskAsCompleted`: Testa o sucesso e o erro (com criação de bug) ao marcar tarefa como concluída.
    *   `TestTaskService_GetTaskByID`:
        *   `"successful get"`: Busca bem-sucedida.
        *   `"task not found"`: Simula erro "não encontrado" do repositório e verifica se NENHUMA tarefa de bug é criada.
        *   `"repository error on get"`: Simula outro erro do repositório e verifica a criação de tarefa de bug.
    *   `TestTaskService_UpdateTask` e `TestTaskService_DeleteTask`: Testam sucesso, erro de validação, erro "não encontrado" (sem bug task) e outros erros do repositório (com bug task).
*   `setupTestDB`: Função helper (não usada ativamente pelos testes com mock, mas presente) para configurar um banco de dados SQLite em memória para testes que poderiam precisar dele.
*   `TestMain`: Função `TestMain` padrão, atualmente sem setup/teardown global significativo para os testes baseados em mock.

**Decisões de Arquitetura e Funcionalidade (nos Testes):**

1.  **Foco na Lógica de Tratamento de Erros e Bug Tasks:** Uma parte significativa dos testes é dedicada a verificar se as tarefas de bug são criadas (ou não) corretamente em diferentes cenários de erro.
2.  **Mocking Detalhado:** O `MockTaskRepository` é configurado em cada subteste para simular comportamentos específicos do repositório (retornar dados, retornar erros específicos).
3.  **Verificação de Efeitos Colaterais (Bug Tasks):** Os testes inspecionam `mockRepo.CreatedBugTasks` para confirmar a criação e o conteúdo das tarefas de bug.
4.  **Distinção de Tipos de Erro:** Os testes validam que erros de "não encontrado" são tratados de forma diferente de erros inesperados do sistema no que diz respeito à criação de tarefas de bug.

**Como se Encaixa no Projeto:**

Esses testes são essenciais para garantir que o `TaskService` funcione conforme o esperado, especialmente sua lógica de tratamento de erros e a funcionalidade de criação automática de tarefas de bug, que é uma característica importante do sistema.
## `internal/tui/prompt.go`

**Propósito e Funcionalidade:**

Este arquivo define um componente TUI reutilizável para solicitar entrada de texto do usuário de forma interativa no terminal. Ele usa um `textinput.Model` do BubbleTea para capturar a entrada. A função principal `GetInput` pode lidar tanto com entrada interativa (quando stdin é um TTY) quanto com entrada não interativa (quando stdin é redirecionado, por exemplo, de um pipe).

**Estruturas de Dados e Funções Chave:**

*   `PromptModel struct`: O modelo BubbleTea para o prompt.
    *   `prompt string`: O texto da pergunta a ser exibida.
    *   `textInput textinput.Model`: O campo de entrada de texto.
    *   `err error`, `quitting bool`, `submitted bool`: Estado do prompt.
    *   `SubmittedCh chan string`: Um canal para enviar o valor submetido de volta para o chamador de `GetInput`.
*   `NewPromptModel(promptText string) PromptModel`: Construtor para `PromptModel`.
*   `(m PromptModel) Init() tea.Cmd`: Retorna `textinput.Blink` para iniciar o cursor.
*   `(m PromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)`:
    *   Trata `tea.KeyEnter`: Se o input não estiver vazio, marca como submetido, envia o valor para `SubmittedCh` e encerra.
    *   Trata `tea.KeyCtrlC`, `tea.KeyEsc`: Marca como encerrando, envia string vazia para `SubmittedCh` e encerra.
    *   Passa outras mensagens para `textInput.Update()`.
*   `(m PromptModel) View() string`: Renderiza o prompt, o campo de entrada e a ajuda.
*   `GetInput(promptText string, output io.Writer, inputReader io.Reader) (string, error)`:
    *   Verifica se `inputReader` (geralmente `os.Stdin`) é um terminal TTY usando `isatty.IsTerminal`.
    *   **Se não for TTY (ex: entrada via pipe):** Lê uma única linha diretamente de `inputReader` usando `bufio.Scanner`.
    *   **Se for TTY:** Cria uma `PromptModel`, inicia um novo programa BubbleTea (`tea.NewProgram`), e executa-o (`p.Run()`). O valor submetido (ou string vazia em caso de desistência) é esperado do `SubmittedCh` da `PromptModel` após `p.Run()` completar, mas a implementação atual de `GetInput` não usa o canal `SubmittedCh` diretamente; em vez disso, ela espera que `p.Run()` retorne o modelo final e então acessa `finalPromptModel.textInput.Value()` se `finalPromptModel.submitted` for verdadeiro. *Nota: Há uma pequena inconsistência aqui, pois `PromptModel.Update` envia para `SubmittedCh`, mas `GetInput` não o lê. No entanto, como `p.Run()` é bloqueante e retorna o modelo final, o valor pode ser recuperado do modelo.* A lógica atual de `GetInput` para o caso TTY é mais simples: executa `p.Run()` e, se não houver erro e o modelo final indicar submissão, retorna o valor do `textInput`. Se o usuário desistir, um erro é retornado.
*   `main_example()`: Uma função de exemplo (comentada) mostrando como usar `GetInput`.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Componente de Prompt Reutilizável:** Fornece uma maneira padronizada de obter entrada de texto do usuário na TUI.
2.  **Suporte a Entradas Interativas e Não Interativas:** A capacidade de `GetInput` de lidar com TTYs e pipes torna-o versátil para diferentes casos de uso da CLI (interação direta vs. scripting).
3.  **Uso de BubbleTea para Interatividade:** Para entradas TTY, aproveita os recursos do BubbleTea para uma boa experiência de edição de texto.
4.  **Comunicação de Resultado:** A `PromptModel` usa um canal (`SubmittedCh`) para comunicar o resultado, embora a função `GetInput` atualmente recupere o valor do modelo final retornado por `p.Run()`.

**Como se Encaixa no Projeto:**

`prompt.go` é um utilitário TUI usado em `cmd/vigenda/main.go` para obter entradas adicionais para certos comandos CLI (como a descrição de uma tarefa ao usar `vigenda tarefa add` sem a flag `--description`). Ele melhora a interatividade dos comandos CLI que podem se beneficiar de entrada de texto mais elaborada.

## `internal/tui/prompt_test.go`

**Propósito e Funcionalidade:**

Este arquivo contém testes unitários para o componente `PromptModel` de `prompt.go`.

**Estruturas de Dados e Funções Chave:**

*   **Test Functions:**
    *   `TestNewPromptModel`: Verifica a inicialização correta de `PromptModel`.
    *   `TestPromptModel_Update_Enter`: Testa o comportamento ao pressionar Enter com entrada não vazia (submissão, encerramento, valor enviado ao canal).
    *   `TestPromptModel_Update_Enter_Empty`: Testa que Enter com entrada vazia não submete nem encerra.
    *   `TestPromptModel_Update_Quit`: Testa as teclas de desistência (Ctrl+C, Esc), verificando o estado de encerramento e o envio de string vazia ao canal.
    *   `TestPromptModel_Update_TextInput`: Verifica se a entrada de texto normal atualiza o valor do `textInput`.
    *   `TestPromptModel_View`: Verifica se a renderização contém os elementos esperados (texto do prompt, valor do input, botão de submissão) e como ela se comporta ao encerrar.
    *   `TestGetInput_Simulated`: Tenta testar a função `GetInput`.
        *   Subteste `"User submits input"`: Tenta simular a entrada do usuário para `GetInput` usando `io.Pipe` para fornecer dados ao `tea.Program`. Esta é uma forma de teste de integração para `GetInput`.

**Decisões de Arquitetura e Funcionalidade (nos Testes):**

1.  **Teste de Modelo e Lógica de Atualização:** Foca em testar o comportamento do `PromptModel` em resposta a diferentes eventos de tecla.
2.  **Verificação de Canal:** Os testes para Enter e Quit verificam se o valor correto (ou string vazia) é enviado para `SubmittedCh`.
3.  **Teste de `GetInput` com Pipe:** `TestGetInput_Simulated` tenta uma abordagem mais integrada para testar `GetInput`, fornecendo um `io.Reader` controlado. Isso é importante para verificar o comportamento de fallback para não-TTY.
4.  **Limitações de Teste TUI:** Reconhece a dificuldade de testar completamente aplicações BubbleTea interativas em testes unitários e foca nos aspectos controláveis.

**Como se Encaixa no Projeto:**

Os testes em `prompt_test.go` ajudam a garantir que o componente de prompt funcione corretamente, tanto em sua lógica interna de modelo quanto na função `GetInput` que o encapsula, especialmente em como lida com diferentes tipos de fontes de entrada.

## `internal/tui/statusbar.go`

**Propósito e Funcionalidade:**

Este arquivo define um componente TUI de barra de status reutilizável usando BubbleTea. A barra de status pode exibir uma mensagem de status principal, mensagens temporárias (efêmeras) e um texto à direita (como a hora atual).

**Estruturas de Dados e Funções Chave:**

*   `StatusBarModel struct`: O modelo BubbleTea para a barra de status.
    *   `width int`: Largura da barra.
    *   `status string`: Mensagem de status principal.
    *   `ephemeralMsg string`, `ephemeralTime time.Time`, `ephemeralTTL time.Duration`: Para mensagens temporárias.
    *   `rightText string`: Texto à direita (ex: hora).
*   `NewStatusBarModel() StatusBarModel`: Construtor.
*   `(m StatusBarModel) Init() tea.Cmd`: Inicia um `tea.Tick` para atualizar a hora a cada segundo.
*   **Mensagens (tea.Msg):**
    *   `UpdateTimeMsg time.Time`: Para atualizar a hora.
    *   `SetStatusMsg string`: Para definir a mensagem de status principal.
    *   `SetEphemeralStatusMsg struct { Text string; TTL time.Duration }`: Para definir uma mensagem efêmera.
    *   `ClearEphemeralMsg struct{}`: Para sinalizar a limpeza de uma mensagem efêmera (usado internamente após o TTL).
*   `(m StatusBarModel) Update(msg tea.Msg) (StatusBarModel, tea.Cmd)`:
    *   Trata `tea.WindowSizeMsg` para atualizar `width`.
    *   Trata `UpdateTimeMsg` para atualizar `rightText` (hora) e agenda o próximo `Tick`.
    *   Trata `SetStatusMsg` para definir `status` e limpar qualquer mensagem efêmera.
    *   Trata `SetEphemeralStatusMsg` para definir `ephemeralMsg`, `ephemeralTime`, `ephemeralTTL`, e agenda um `Tick` para `ClearEphemeralMsg`.
    *   Trata `ClearEphemeralMsg`: Se o TTL da mensagem efêmera realmente expirou, limpa `ephemeralMsg`.
*   `clearEphemeral()`: Helper para limpar a mensagem efêmera.
*   `(m StatusBarModel) View() string`: Renderiza a barra de status.
    *   Decide se exibe a mensagem de status principal ou a efêmera (se ativa e não expirada).
    *   Formata e alinha o texto da esquerda (status) e da direita (hora).
    *   Calcula o espaço (`gap`) entre o texto da esquerda e da direita para preencher a largura.
    *   Trunca o texto da esquerda se for muito longo para caber.
    *   Aplica estilos (`statusBarStyle`, etc.) usando `lipgloss`.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Componente de UI Reutilizável:** A barra de status é projetada para ser um componente genérico que pode ser integrado em aplicações BubbleTea maiores.
2.  **Mensagens Efêmeras com TTL:** Suporta mensagens temporárias que desaparecem automaticamente após um tempo, útil para notificações breves (ex: "Salvo!").
3.  **Atualização Dinâmica da Hora:** A hora é atualizada automaticamente a cada segundo.
4.  **Layout Flexível com Lipgloss:** Usa `lipgloss` para estilização e tenta gerenciar o layout do texto à esquerda e à direita, incluindo truncamento e preenchimento de espaço.
5.  **Comunicação Baseada em Mensagens:** A barra de status é atualizada através do envio de mensagens específicas (`SetStatusMsg`, `SetEphemeralStatusMsg`).

**Como se Encaixa no Projeto:**

A `StatusBarModel` pode ser usada por qualquer modelo BubbleTea principal (como o `app.Model` ou os submodelos de tarefas, turmas, etc.) para fornecer feedback contextual e informações persistentes (como a hora) ao usuário na parte inferior da tela. Ela melhora a experiência do usuário fornecendo um local consistente para mensagens de status.

## `internal/tui/statusbar_test.go`

**Propósito e Funcionalidade:**

Este arquivo contém testes unitários para o componente `StatusBarModel` de `statusbar.go`.

**Estruturas de Dados e Funções Chave:**

*   **Test Functions:**
    *   `TestNewStatusBarModel`: Verifica a inicialização correta.
    *   `TestStatusBarModel_Init`: Verifica se `Init` retorna um comando (para o `Tick` da hora).
    *   `TestStatusBarModel_Update_WindowSize`: Testa a atualização da largura.
    *   `TestStatusBarModel_Update_Time`: Testa a atualização da hora e o re-agendamento do `Tick`.
    *   `TestStatusBarModel_Update_SetStatus`: Testa a definição da mensagem de status principal e a limpeza de mensagens efêmeras.
    *   `TestStatusBarModel_Update_SetEphemeralStatus`: Testa a definição de mensagens efêmeras e o agendamento de sua limpeza.
    *   `TestStatusBarModel_Update_ClearEphemeral`: Testa a lógica de limpeza de mensagens efêmeras (se o TTL passou ou não).
    *   `TestStatusBarModel_View`: Testa a renderização em diferentes cenários: status padrão, com mensagem efêmera, com mensagem efêmera expirada, sem largura definida, e com truncamento de status longo.
    *   `TestStatusBar_GapCalculation`: Um teste conceitual para verificar se os componentes esquerdo, direito e o espaço entre eles estão presentes na renderização.

**Decisões de Arquitetura e Funcionalidade (nos Testes):**

1.  **Cobertura de Estados e Mensagens:** Os testes cobrem como o modelo responde a todas as mensagens que ele pode processar e como seu estado interno (`status`, `ephemeralMsg`, `rightText`) muda.
2.  **Teste de Lógica de Tempo (TTL):** `TestStatusBarModel_Update_ClearEphemeral` e partes de `TestStatusBarModel_View` verificam se as mensagens efêmeras são exibidas e ocultadas corretamente com base em seu TTL.
3.  **Teste de Renderização (View):** `TestStatusBarModel_View` verifica se a string renderizada contém os textos esperados e se o truncamento funciona. Testar a renderização exata de componentes estilizados com `lipgloss` pode ser complexo, então os testes focam na presença de conteúdo chave.
4.  **Simulação de Ticks:** Os testes que envolvem tempo (como `UpdateTimeMsg` ou `ClearEphemeralMsg` que são acionados por `Tick`) simulam o recebimento dessas mensagens diretamente.

**Como se Encaixa no Projeto:**

Os testes para `StatusBarModel` garantem que este componente de UI reutilizável funcione de forma confiável, exibindo as informações corretas e lidando com mensagens efêmeras e atualizações de tempo como esperado.

## `internal/tui/table.go`

**Propósito e Funcionalidade:**

Este arquivo define um wrapper simples em torno do componente `table.Model` do BubbleTea, fornecendo um `TableModel` básico que pode ser usado para exibir dados tabulares. Ele também inclui uma função helper `ShowTable` para exibir rapidamente uma tabela.

**Estruturas de Dados e Funções Chave:**

*   `TableModel struct`: Contém uma `table.Model` interna.
*   `NewTableModel(columns []table.Column, rows []table.Row) TableModel`: Construtor que cria e configura uma `table.Model` com as colunas, linhas e estilos padrão fornecidos.
*   `(m TableModel) Init() tea.Cmd`: Retorna `nil`.
*   `(m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)`:
    *   Trata `tea.KeyEsc` para alternar o foco da tabela (`Blur()`/`Focus()`).
    *   Trata `tea.KeyCtrlC`, `'q'` para encerrar.
    *   Passa outras mensagens para `m.table.Update(msg)`.
*   `(m TableModel) View() string`: Renderiza a tabela envolvida por um `baseStyle` (que adiciona uma borda).
*   `ShowTable(columns []table.Column, rows []table.Row, output io.Writer)`: Função helper que cria um `TableModel` com os dados fornecidos e executa um novo programa BubbleTea para exibi-lo. Útil para visualizações rápidas de tabelas.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Wrapper Simples:** Fornece uma maneira um pouco mais estruturada de usar `table.Model` dentro de uma aplicação BubbleTea maior, encapsulando a configuração de estilo e o comportamento básico de foco.
2.  **Estilização Padrão:** Aplica um conjunto de estilos padrão (cabeçalho, linha selecionada) à tabela.
3.  **Helper `ShowTable`:** Conveniente para depuração ou para exibir dados tabulares de forma isolada sem integrá-los a um modelo TUI maior.

**Como se Encaixa no Projeto:**

O `TableModel` ou diretamente o `table.Model` do BubbleTea é usado extensivamente por outros módulos da TUI (como `tasks/model.go`, `classes/model.go`) para exibir listas de tarefas, turmas, alunos, etc. Este arquivo `table.go` fornece um modelo básico e um helper que podem ter sido usados para prototipagem ou para componentes de tabela mais simples. Os módulos mais complexos parecem implementar sua própria lógica de tabela ou usar `table.Model` diretamente.

## `internal/tui/table_test.go`

**Propósito e Funcionalidade:**

Este arquivo contém testes unitários para o `TableModel` de `table.go`.

**Estruturas de Dados e Funções Chave:**

*   **Test Functions:**
    *   `TestNewTableModel`: Verifica a inicialização correta do `TableModel` (foco, número de colunas/linhas).
    *   `TestTableModel_Update`: Testa o tratamento de teclas:
        *   'q' para encerrar.
        *   'Esc' para alternar foco.
        *   Outras teclas (como setas) são passadas para a tabela interna (o teste verifica que não há crash).
    *   `TestTableModel_View`: Verifica se a renderização da tabela contém os cabeçalhos e dados das células, e se as bordas do `baseStyle` estão presentes.
    *   `TestShowTable`: Testa conceitualmente a função `ShowTable` verificando a saída `View()` do modelo que ela criaria. Não executa `p.Run()` diretamente devido à sua natureza bloqueante e interativa.

**Decisões de Arquitetura e Funcionalidade (nos Testes):**

1.  **Teste do Wrapper:** Foca em testar a lógica adicionada pelo `TableModel` (manipulação de foco com Esc) e a configuração da tabela interna.
2.  **Teste de Renderização (Básico):** Verifica a presença de conteúdo esperado na string renderizada.
3.  **Limitações no Teste de `ShowTable`:** Reconhece que testar uma função que executa `tea.NewProgram().Run()` é complexo em um teste unitário e, em vez disso, testa o modelo que `ShowTable` instancia.

**Como se Encaixa no Projeto:**

Os testes garantem que o `TableModel` básico funcione como esperado, especialmente sua lógica de alternância de foco e a renderização básica. Isso é importante se este componente for usado como base para outras exibições de tabela na aplicação.

## `internal/tui/tui.go`

**Propósito e Funcionalidade:**

Este arquivo parece ser uma tentativa inicial ou um esqueleto para a interface TUI principal, focada especificamente no gerenciamento de turmas e alunos. Ele se sobrepõe significativamente em intenção e funcionalidade com `internal/app/app.go` e `internal/app/classes/model.go`. Dada a estrutura mais completa e modular de `internal/app/`, este arquivo `internal/tui/tui.go` pode ser uma versão mais antiga, um protótipo, ou um componente específico que talvez não esteja totalmente integrado ou tenha sido substituído pela abordagem em `internal/app/`.

**Estruturas de Dados e Funções Chave:**

*   `KeyMap`, `DefaultKeyMap`: Definições de keybindings (semelhantes às de `internal/app/app.go`, mas não usadas diretamente no `Model` aqui).
*   `Model struct`: Modelo BubbleTea.
    *   `classService service.ClassService`.
    *   `list list.Model`: Para exibir turmas ou alunos.
    *   `spinner spinner.Model`.
    *   `currentView app.View`: Usa o mesmo enum de `internal/app/views.go`.
    *   `isLoading bool`, `err error`.
    *   `classes []models.Class`, `students []models.Student`.
    *   `selectedClass *models.Class`.
*   `NewTUIModel(cs service.ClassService) Model`: Construtor.
*   `loadInitialData()`, `loadClasses() tea.Cmd`, `loadStudentsForClass(classID int64) tea.Cmd`: Funções para carregar dados.
*   `(m Model) Init() tea.Cmd`: Inicia o carregamento de turmas.
*   `(m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)`: Lida com mensagens.
    *   Lógica de navegação entre visualização de turmas e visualização de alunos de uma turma selecionada.
    *   Processa `classesLoadedMsg` e `studentsLoadedMsg` para popular a `list.Model`.
    *   Lida com `errMsg` e `tea.WindowSizeMsg`.
*   `(m Model) View() string`: Renderiza a UI.
    *   `headerView()` e `footerView()`: Helpers para renderizar cabeçalho e rodapé.
*   `listItemClass`, `listItemStudent`: Structs para itens da lista (semelhantes às de `internal/app/classes/model.go`).
*   `errMsg`, `classesLoadedMsg`, `studentsLoadedMsg`: Tipos de mensagem.
*   `Start(classService service.ClassService) error`: Inicia o programa BubbleTea com este modelo.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Foco em Turmas/Alunos:** Este modelo TUI está especificamente focado na funcionalidade de gerenciamento de turmas e alunos.
2.  **Estrutura MVU:** Segue o padrão Model-View-Update do BubbleTea.
3.  **Carregamento Assíncrono:** Carrega dados de turmas e alunos de forma assíncrona.
4.  **Navegação Simples:** Permite visualizar uma lista de turmas e, ao selecionar uma, visualizar seus alunos. O retorno à lista de turmas também é implementado.
5.  **Potencial Redundância/Conflito:** A funcionalidade aqui é um subconjunto do que é gerenciado por `internal/app/app.go` (como roteador principal) e `internal/app/classes/model.go` (como o submodelo específico para turmas). O `app.Model` em `internal/app/app.go` é projetado para ser o modelo TUI principal que delega para submódulos como `classes.Model`. Se `internal/tui/tui.go` for usado diretamente, ele não se encaixa nessa arquitetura de submódulos. `cmd/vigenda/main.go` chama `app.StartApp`, que usa `internal/app/app.Model`.

**Como se Encaixa no Projeto:**

Dado que `cmd/vigenda/main.go` usa `app.StartApp` (que inicializa `internal/app/app.Model`), o arquivo `internal/tui/tui.go` provavelmente **não é o ponto de entrada TUI principal atualmente ativo**. Pode ser:
*   Um protótipo anterior.
*   Um componente que foi refatorado e movido para dentro de `internal/app/classes/`.
*   Uma tentativa de uma TUI alternativa ou um módulo que não foi totalmente integrado da mesma forma que os submódulos em `internal/app/`.

Seu `Start` e `NewTUIModel` são semelhantes ao `StartApp` e `New` de `internal/app/app.go`, mas focados apenas em `ClassService`. A TUI definida aqui é menos abrangente do que a arquitetura modular apresentada em `internal/app/`.

## `test_prompt/main.go`

**Propósito e Funcionalidade:**

Este é um pequeno programa Go localizado fora da estrutura principal da aplicação (`cmd/` ou `internal/`), provavelmente usado para testar isoladamente o componente `tui.GetInput` de `internal/tui/prompt.go`. Ele demonstra como `GetInput` se comporta tanto com entrada de terminal interativa quanto com entrada redirecionada (pipe).

**Estruturas de Dados e Funções Chave:**

*   `main()`:
    *   Imprime mensagens de log para `stderr` indicando o início do teste e o status do TTY de `os.Stdin`.
    *   Chama `tui.GetInput("Enter data:", os.Stdout, os.Stdin)` para solicitar entrada.
    *   Imprime a entrada recebida (prefixada com "Received_from_prompt:") ou um erro para `stdout`.

**Decisões de Arquitetura e Funcionalidade:**

1.  **Teste Isolado de Componente:** Permite testar `tui.GetInput` sem a complexidade do resto da aplicação Vigenda.
2.  **Demonstração de Dupla Funcionalidade:** Projetado para ser executado interativamente (`go run test_prompt/main.go`) e com entrada via pipe (`echo "MyPipedInput" | go run test_prompt/main.go`) para verificar ambos os caminhos de código em `GetInput`.
3.  **Uso de Stderr para Logs de Teste:** Mensagens de diagnóstico são enviadas para `stderr` para não interferir com a saída do `GetInput` que vai para `stdout` (no caso de pipe).

**Como se Encaixa no Projeto:**

Serve como um utilitário de desenvolvimento e teste para o componente `tui.GetInput`. Não faz parte do binário final da aplicação Vigenda, mas é útil para verificar o comportamento do prompt em diferentes cenários de execução.

## `tests/integration/cli_integration_test.go`

**Propósito e Funcionalidade:**

Este arquivo contém testes de integração para a aplicação Vigenda CLI. O objetivo é construir o binário da CLI e executá-lo com diferentes argumentos, comparando a saída (stdout) com "golden files" (arquivos de saída esperada). Ele também lida com a configuração de bancos de dados de teste para cada cenário.

**Estruturas de Dados e Funções Chave:**

*   Variáveis globais `binName` e `binPath` para o nome e caminho do binário compilado.
*   `TestMain(m *testing.M)`:
    *   Determina o nome do binário com base no SO.
    *   Cria um diretório temporário para o binário.
    *   Compila o binário da CLI (`cmd/vigenda/main.go`) usando `go build -a -o ...`. O `-a` força a reconstrução.
    *   Executa os testes (`m.Run()`).
    *   Limpa o diretório temporário do binário.
*   `setupTestDB(t *testing.T, testName string) string`:
    *   Cria um diretório `test_dbs` se não existir.
    *   Cria um arquivo de banco de dados SQLite específico para o teste (ex: `vigenda_test_TestDashboardOutput.db`), removendo qualquer um existente para garantir um estado limpo.
    *   Define as variáveis de ambiente `VIGENDA_DB_TYPE="sqlite"` e `VIGENDA_DB_PATH` para que a CLI use este banco de dados.
    *   Aplica o schema inicial ao banco de dados de teste lendo `internal/database/migrations/001_initial_schema.sql`.
    *   Retorna o caminho para o arquivo do banco de dados.
*   `seedDB(t *testing.T, dbPath string, statements []string)`:
    *   Abre o banco de dados de teste especificado.
    *   Executa uma lista de statements SQL para popular o banco com dados de teste.
*   `runCLI(t *testing.T, args ...string) (string, string, error)`:
    *   Executa o binário compilado com os argumentos fornecidos, usando `exec.CommandContext` com um timeout de 30 segundos.
    *   Captura stdout e stderr.
    *   Retorna stdout, stderr e o erro da execução.
*   `assertGoldenFile(t *testing.T, actualOutput string, goldenFilePath string)`:
    *   Lê o conteúdo do `goldenFilePath`.
    *   Normaliza quebras de linha e remove espaços em branco extras do início/fim tanto da saída real quanto da esperada.
    *   Compara a saída real normalizada com a esperada. Se diferentes, falha o teste e imprime ambas.
*   **Test Functions (exemplos):**
    *   `TestDashboardOutput`:
        *   Configura um banco de dados de teste (embora a saída atual do dashboard seja estática).
        *   Executa `vigenda` (sem argumentos).
        *   Compara stdout com `golden_files/dashboard_output.txt`.
    *   `TestNotasLancarOutput` (Placeholder): Comentado, observa a complexidade de testar TUI interativo.
    *   `TestRelatorioProgressoTurmaOutput` (Placeholder): Comentado, observa a necessidade de setup de DB.
    *   `TestTarefaListarTurmaOutput`:
        *   Configura um banco de dados de teste.
        *   Popula o banco com dados específicos (usuário, disciplina, turma, tarefas) usando `seedDB`.
        *   Executa `vigenda tarefa listar --classid 1`.
        *   Compara stdout com `golden_files/tarefa_listar_turma_output.txt`.
    *   `TestFocoIniciarOutput` (Placeholder): Comentado, observa a dificuldade com saídas sensíveis ao tempo.

**Decisões de Arquitetura e Funcionalidade (nos Testes):**

1.  **Testes End-to-End (Black-Box):** Os testes tratam a CLI como uma caixa preta, executando o binário compilado e verificando sua saída, sem inspecionar o código interno diretamente.
2.  **Compilação Dinâmica da CLI:** A CLI é compilada no início dos testes (`TestMain`) para garantir que a versão mais recente do código seja testada.
3.  **Bancos de Dados de Teste Isolados:** Cada função de teste que interage com o banco de dados usa `setupTestDB` para criar um banco de dados SQLite limpo e específico para o teste, garantindo o isolamento entre os testes. As variáveis de ambiente são usadas para direcionar a CLI para esses bancos de dados.
4.  **Seeding de Dados Específico do Teste:** `seedDB` permite que cada teste popule o banco de dados com o estado exato necessário para seu cenário.
5.  **Golden File Testing:** A comparação com "golden files" é uma técnica comum para testar a saída de CLIs, onde a saída esperada é armazenada em arquivos. A normalização da saída antes da comparação ajuda a evitar falhas espúrias devido a diferenças de quebra de linha ou espaçamento.
6.  **Timeout para Comandos CLI:** O uso de `exec.CommandContext` com timeout previne que os testes fiquem bloqueados indefinidamente se a CLI travar.
7.  **Reconhecimento de Desafios:** Os placeholders e comentários para testes de TUI interativo e saídas sensíveis ao tempo reconhecem as limitações e desafios inerentes a esses tipos de teste em um framework de integração simples.

**Como se Encaixa no Projeto:**

Os testes de integração em `cli_integration_test.go` são cruciais para verificar o comportamento da aplicação Vigenda como um todo, do ponto de vista do usuário da CLI. Eles garantem que os diferentes comandos funcionem corretamente, interajam adequadamente com o banco de dados (configurado para o teste) e produzam a saída esperada para cenários específicos. Eles complementam os testes unitários, validando a integração entre as diferentes camadas da aplicação.
