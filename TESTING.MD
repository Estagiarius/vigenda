# Documentação de Testes

Este documento descreve a estratégia de testes para o projeto Vigenda, os tipos de testes implementados e como executá-los. Testes são cruciais para garantir a qualidade, estabilidade e manutenibilidade do código.

## 1. Filosofia de Testes

Nosso objetivo é ter uma cobertura de testes abrangente que nos dê confiança para refatorar o código e adicionar novas funcionalidades sem introduzir regressões. Priorizamos:

-   **Testes Unitários:** Para verificar pequenas unidades de código (funções, métodos) de forma isolada.
-   **Testes de Integração:** Para verificar a interação entre diferentes componentes ou módulos do sistema, especialmente a CLI com as camadas de serviço, repositório e o banco de dados.

Incentivamos a escrita de testes *antes* ou *durante* o desenvolvimento do código (Test-Driven Development - TDD, ou Behavior-Driven Development - BDD), sempre que prático.

## 2. Tipos de Testes e Ferramentas

### 2.1. Testes Unitários

-   **Propósito:** Validar a lógica de funções e métodos individuais dentro de cada pacote (ex: `internal/service`, `internal/repository`, `internal/tui`, `internal/config`). Devem ser rápidos e independentes de dependências externas; mocks/stubs são usados para isolar o código sob teste (ex: mockar a camada de repositório ao testar a camada de serviço).
-   **Localização:** Em arquivos com sufixo `_test.go` no mesmo pacote do código que está sendo testado (ex: `internal/service/task_service_test.go`).
-   **Ferramentas Utilizadas:**
    -   Pacote `testing` nativo do Go.
    -   `github.com/stretchr/testify/assert`: Para asserções fluentes (ex: `assert.Equal(t, expected, actual)`).
    -   `github.com/stretchr/testify/require`: Similar ao `assert`, mas interrompe a execução do teste na primeira falha, útil para verificações críticas.
    -   `github.com/stretchr/testify/mock`: Para criar dublês de teste (mocks) de interfaces (ex: para mockar `repository.TaskRepositoryInterface` ao testar `service.TaskService`).
-   **Como Escrever:**
    -   Foco em um único aspecto ou comportamento por função de teste.
    -   Use nomes descritivos para as funções de teste, começando com `Test` (ex: `TestTaskService_CreateTask_Success`).
    -   Siga o padrão Arrange-Act-Assert (AAA):
        1.  **Arrange:** Configure as condições iniciais e as entradas para o teste (incluindo mocks).
        2.  **Act:** Execute a unidade de código que está sendo testada.
        3.  **Assert:** Verifique se os resultados (valores retornados, estado dos mocks) são os esperados.

### 2.2. Testes de Integração

-   **Propósito:** Verificar se diferentes partes do sistema funcionam corretamente juntas. No Vigenda, isso primariamente envolve testar a interação da CLI (construída ou via `go run`) com as camadas de serviço, repositório e o banco de dados SQLite real.
-   **Localização:** No diretório `tests/integration/`, como `cli_integration_test.go`.
-   **Ferramentas Utilizadas:**
    -   Pacote `testing` nativo do Go.
    -   `os/exec` para executar o binário da CLI compilado ou `go run ./cmd/vigenda/main.go ...`.
    -   `bytes.Buffer` para capturar a saída padrão (stdout) e erro padrão (stderr) da CLI.
    -   Pacote `github.com/mattn/go-sqlite3` para configurar, popular e inspecionar um banco de dados de teste SQLite dedicado para cada execução de teste ou suíte.
    -   Comparação com "Golden Files" (arquivos de texto em `tests/integration/golden_files/`) para verificar saídas de texto complexas ou esperadas.
-   **Como Escrever:**
    -   Cada teste de integração ou suíte de testes de integração geralmente configura um banco de dados de teste limpo e temporário (arquivos `.db` em `tests/integration/test_dbs/` são exemplos, mas o ideal é que os testes criem seus próprios bancos de dados temporários).
    -   Executa um ou mais comandos da CLI, simulando a entrada do usuário.
    -   Verifica a saída da CLI (stdout, stderr), o código de saída, e/ou o estado do banco de dados após a execução dos comandos para garantir que a integração funcionou como esperado.

### 2.3. Testes End-to-End (E2E) para CLI/TUI

-   Para uma aplicação CLI com TUI como Vigenda, os testes de integração (`tests/integration/cli_integration_test.go`) que executam a aplicação e verificam sua saída e estado do banco de dados já cobrem uma parte significativa do que seria um teste E2E.
-   Testar a TUI interativa de forma totalmente automatizada é desafiador. Abordagens podem incluir:
    -   Testar os componentes da TUI (`internal/tui`) de forma unitária, verificando seus modelos e lógica de atualização.
    -   Para fluxos críticos da TUI, podem ser considerados testes que enviam sequências de teclas simuladas e capturam o estado da tela (snapshot testing), embora isso possa ser frágil. Atualmente, esta não é uma prática estabelecida no projeto.

### 2.4. Testes de Performance

-   **Propósito:** Avaliar a responsividade e o uso de recursos da aplicação sob condições específicas (ex: grande volume de dados no banco).
-   **Ferramentas Utilizadas:**
    -   **Benchmarks do Go:** Usando o tipo `testing.B` em arquivos `_test.go` para medir o desempenho de funções críticas.
        ```go
        // Exemplo de benchmark em algum_test.go
        // func BenchmarkNomeDaFuncaoCritica(b *testing.B) {
        //     // Setup necessário
        //     for i := 0; i < b.N; i++ {
        //         NomeDaFuncaoCritica(parametro) // Função a ser benchmarkada
        //     }
        // }
        ```
    -   **Ferramentas de Profiling do Go (`pprof`):** Para identificar gargalos de CPU e memória.
        -   Comandos: `go test -bench=. -cpuprofile cpu.prof -memprofile mem.prof ./...` e depois `go tool pprof cpu.prof`.
-   Atualmente, não há benchmarks formais ou testes de performance automatizados no repositório, mas podem ser adicionados conforme necessário para otimizar gargalos.

### 2.5. Linters e Análise Estática

-   **Propósito:** Detectar problemas de estilo de código, possíveis bugs, "code smells" e garantir a consistência do código sem executá-lo.
-   **Ferramentas Utilizadas:**
    -   **`gofmt` / `goimports`:** Ferramentas padrão do Go para formatação de código. `goimports` também organiza os imports. É esperado que todo o código Go seja formatado com uma dessas ferramentas.
        -   Comando para verificar formatação (sem alterar arquivos): `goimports -l .` (lista arquivos que precisam de formatação) ou `gofmt -l .`
        -   Comando para formatar arquivos no local: `goimports -w .` ou `gofmt -w .`
    -   **`golangci-lint`:** Um agregador rápido de linters para Go. Altamente recomendado para manter a qualidade do código.
        -   **Instalação:** Consulte `INSTALLATION.MD` ou o site oficial do `golangci-lint`.
        -   **Configuração:** Idealmente, configurado com um arquivo `.golangci.yml` ou `.golangci.yaml` na raiz do projeto. Se ausente, usa configurações padrão.
        -   **Comando para executar:** `golangci-lint run ./...`
        -   Para tentar corrigir automaticamente alguns problemas (quando suportado pelos linters): `golangci-lint run --fix ./...`

## 3. Como Executar os Testes

### 3.1. Executando Todos os Testes (Unitários e de Integração)

Para executar todos os testes definidos no projeto (arquivos `_test.go` em todos os subdiretórios, incluindo `tests/integration`):
```bash
go test ./...
```
Para uma saída mais detalhada, incluindo os nomes dos testes executados e o status (PASS/FAIL):
```bash
go test -v ./...
```

### 3.2. Executando Testes de um Pacote Específico

Para executar testes apenas de um pacote específico (ex: o pacote `internal/service`):
```bash
go test ./internal/service/...
```
Ou para executar testes apenas no diretório de integração:
```bash
go test ./tests/integration/...
```
Adicione `-v` para saída verbosa, se desejado (ex: `go test -v ./internal/service/...`).

### 3.3. Executando um Teste Específico por Nome

Use o flag `-run` com uma expressão regular que corresponda ao nome da função de teste (ou um conjunto de testes). O nome deve incluir o prefixo `Test` (para testes) ou `Benchmark` (para benchmarks).

-   **Para executar uma função de teste específica (ex: `TestTaskService_CreateTask_Success` no pacote `internal/service`):**
    ```bash
    go test -v -run ^TestTaskService_CreateTask_Success$ ./internal/service/...
    ```
-   **Para executar todos os testes cujo nome comece com `TestTask` dentro do pacote `internal/service`:**
    ```bash
    go test -v -run ^TestTask ./internal/service/...
    ```
-   **Para executar todos os testes de integração em `tests/integration` que contenham `DashboardOutput` no nome:**
    ```bash
    go test -v -run DashboardOutput ./tests/integration/...
    ```
Lembre-se que `-run` aceita uma expressão regular. Usar `^NomeDoTeste$` garante que apenas o teste com esse nome exato seja executado.

### 3.4. Executando Benchmarks de Performance

Para executar todos os benchmarks no projeto:
```bash
go test -bench=. ./...
```
Para executar um benchmark específico (ex: `BenchmarkMinhaFuncao`):
```bash
go test -bench=^BenchmarkMinhaFuncao$ ./caminho/do/pacote_com_benchmark/...
```
Para obter dados de profiling durante a execução dos benchmarks:
```bash
go test -bench=. -cpuprofile cpu.prof -memprofile mem.prof ./...
# Em seguida, analise com: go tool pprof cpu.prof ou go tool pprof mem.prof
```

### 3.5. Verificando Cobertura de Teste

É uma boa prática verificar a cobertura de código dos testes para identificar áreas não testadas.

-   **Gerar perfil de cobertura e exibir no terminal:**
    ```bash
    go test -cover ./...
    ```
    Isso mostrará a porcentagem de cobertura para cada pacote testado.

-   **Gerar relatório de cobertura em HTML:**
    1.  Primeiro, gere o perfil de cobertura:
        ```bash
        go test -coverprofile=coverage.out ./...
        ```
    2.  Depois, visualize o relatório em HTML:
        ```bash
        go tool cover -html=coverage.out
        ```
        Este comando abrirá o relatório de cobertura detalhado no seu navegador padrão, destacando as linhas de código cobertas e não cobertas.

Nossa meta de cobertura é **manter acima de 80%** para pacotes críticos como `internal/service` e `internal/repository`, e buscar uma boa cobertura geral para os demais pacotes. A qualidade e relevância dos testes são mais importantes do que a porcentagem de cobertura por si só.

## 4. Ambiente de Teste

-   **Banco de Dados de Teste:** Os testes de integração (`tests/integration/`) utilizam arquivos de banco de dados SQLite separados, frequentemente criados ou limpos no setup de cada teste ou suíte para garantir isolamento. Veja os arquivos em `tests/integration/test_dbs/` como exemplos de estado de banco de dados para testes específicos, embora o ideal seja que os testes criem seus próprios bancos de dados temporários.
-   **Mocks e Stubs:** Para testes unitários, especialmente na camada de serviço, as dependências de repositório são mockadas usando `testify/mock`. Isso isola a unidade de código sob teste e torna os testes mais rápidos e determinísticos.
-   **Variáveis de Ambiente:** Geralmente, não são necessárias variáveis de ambiente específicas para executar os testes, a menos que um teste específico as configure para simular diferentes cenários.

## 5. Adicionando Novos Testes

-   Ao corrigir um bug, escreva primeiro um teste que reproduza o bug. Este teste deve falhar com o código existente e passar após a correção.
-   Ao adicionar uma nova funcionalidade, escreva testes (unitários e, se aplicável, de integração) que cubram os principais casos de uso, casos de borda e cenários de erro da funcionalidade.
-   Certifique-se de que os novos testes sigam os padrões e ferramentas descritos neste documento.
-   Consulte o `CONTRIBUTING.MD` para mais diretrizes sobre como escrever e submeter código, incluindo testes.

---

Manter uma suíte de testes robusta é responsabilidade de todos os contribuidores. Certifique-se de que seus Pull Requests incluem testes apropriados para as alterações feitas e que todos os testes (incluindo os existentes) continuam passando.
