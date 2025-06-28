# Documentação para Agentes de IA

Olá, Agente! Este documento fornece diretrizes e informações para ajudá-lo a entender e trabalhar eficientemente com este codebase do projeto **Vigenda**.

## 1. Visão Geral do Projeto "Vigenda"

-   **Propósito Principal:** Vigenda é uma aplicação de linha de comando (CLI) com uma Interface de Texto do Usuário (TUI), desenvolvida em Go. Seu objetivo é ajudar estudantes a gerenciar suas atividades acadêmicas, como tarefas, aulas, disciplinas e avaliações.
-   **Tecnologias Chave:**
    -   Linguagem Principal: **Go** (versão 1.23.0, toolchain go1.24.3 - conforme `go.mod`).
    -   Compilador C: **GCC** (necessário para dependências CGO como `mattn/go-sqlite3`).
    -   Framework CLI: **`spf13/cobra`**.
    -   Framework TUI: **`charmbracelet/bubbletea`**, com componentes de **`charmbracelet/bubbles`** e estilização via **`charmbracelet/lipgloss`**.
    -   Banco de Dados: **SQLite** (usando o driver `mattn/go-sqlite3`).
    -   Testes: Pacote `testing` nativo do Go, `stretchr/testify` (assert, require, mock).
-   **Estrutura do Repositório Principal (Pacotes Go):**
    -   `README.md`: Visão geral, instalação básica, como contribuir. **LEIA PRIMEIRO.**
    -   `TECHNICAL_SPECIFICATION.MD`: Detalhes da arquitetura em camadas, fluxo de dados, padrões de design. Consulte para entender a organização.
    -   `INSTALLATION.MD`: Instruções detalhadas para configurar o ambiente de desenvolvimento, incluindo Go, GCC e ferramentas.
    -   `TESTING.MD`: Como executar testes unitários, de integração, benchmarks e linters. Inclui comandos específicos.
    -   `CONTRIBUTING.MD`: Diretrizes para contribuição (humanos e IA), padrões de codificação Go, formato de commit (Conventional Commits), processo de PR. **SIGA RIGOROSAMENTE.**
    -   `CHANGELOG.MD`: Registro de mudanças.
    -   `AGENTS.md`: Este arquivo.
    -   `go.mod`, `go.sum`: Gerenciamento de dependências Go.
    -   `build.sh`: Script para construir binários para múltiplas plataformas.
    -   `cmd/vigenda/main.go`: Ponto de entrada da aplicação CLI.
    -   `internal/`: Contém o código principal da aplicação, não destinado a ser importado por outros projetos.
        -   `internal/config/config.go`: Carregamento e gerenciamento de configurações.
        -   `internal/database/database.go`: Conexão com o banco de dados SQLite.
        -   `internal/database/migrations/`: Migrações de esquema do banco de dados (ex: `001_initial_schema.sql`).
        -   `internal/models/models.go`: Definições das estruturas de dados (structs Go) do domínio.
        -   `internal/repository/`: Camada de acesso a dados (operações CRUD), abstraindo o SQLite. Arquivos como `task_repository.go`.
        -   `internal/service/`: Camada de lógica de negócios, orquestrando operações. Arquivos como `task_service.go`.
        -   `internal/tui/`: Componentes e lógica da Interface de Texto do Usuário (Bubbletea). Arquivos como `prompt.go`, `table.go`, `statusbar.go`.
    -   `tests/`:
        -   `tests/integration/cli_integration_test.go`: Testes de integração da CLI.
        -   `tests/integration/golden_files/`: Arquivos de saída esperada para testes de integração.
        -   `tests/integration/test_dbs/`: Bancos de dados de teste pré-populados ou exemplos.
    -   `dist/`: Diretório onde os binários compilados pelo `build.sh` são colocados (ignorado pelo Git).
-   **Documentação Adicional:**
    -   Para qualquer tarefa, comece consultando `README.md`, `INSTALLATION.MD`, `CONTRIBUTING.MD`, `TESTING.MD` e `TECHNICAL_SPECIFICATION.MD`.
    -   A `API_DOCUMENTATION.md` não é relevante no momento, pois Vigenda é uma CLI e não expõe uma API de rede.

## 2. Configuração do Ambiente

-   Siga **rigorosamente** as instruções em `INSTALLATION.MD` para configurar seu ambiente.
    -   Certifique-se de ter a versão correta do **Go** e **GCC** instalados.
    -   Instale `goimports` e `golangci-lint` conforme as instruções em `INSTALLATION.MD`.
-   O banco de dados SQLite é baseado em arquivo e gerenciado pela aplicação; nenhuma configuração manual de servidor de banco de dados é necessária.

## 3. Tarefas Comuns e Como Abordá-las

### 3.1. Entendendo o Código Existente
-   Comece pela função ou módulo principal relacionado à sua tarefa (ex: um comando em `cmd/vigenda/main.go`, um serviço em `internal/service/`, ou um componente TUI em `internal/tui/`).
-   Use as ferramentas de busca (`grep`) para encontrar definições de funções, tipos e seus usos.
-   Leia os comentários no código (Godoc para funções públicas, comentários inline para lógica complexa). Se não houver ou estiverem desatualizados, considere adicionar/atualizar conforme as diretrizes em `CONTRIBUTING.MD`.
-   Analise os testes relacionados (arquivos `_test.go`). Eles demonstram como o código deve ser usado e qual é o comportamento esperado.

### 3.2. Implementando Novas Funcionalidades
1.  **Planeje:**
    -   Certifique-se de que a funcionalidade está bem definida.
    -   Identifique os arquivos e módulos que precisarão ser modificados ou criados (ex: novo comando Cobra, novo método de serviço, alterações no repositório, novo componente TUI).
    -   Considere como a nova funcionalidade se encaixa na arquitetura existente (consulte `TECHNICAL_SPECIFICATION.MD`).
2.  **Escreva Testes Primeiro (TDD/BDD quando possível):**
    -   Consulte `TESTING.MD` para tipos de testes e ferramentas.
    -   Escreva testes unitários para a nova lógica de serviço e repositório.
    -   Se a funcionalidade envolve um novo comando CLI ou altera um existente, adicione ou atualize os testes de integração em `tests/integration/`.
3.  **Implemente o Código:**
    -   Siga as convenções de estilo de código e padrões de design descritos em `CONTRIBUTING.MD` e `TECHNICAL_SPECIFICATION.MD`.
    -   **Comente seu código:** Use Godoc para todas as funções, tipos e variáveis exportadas. Adicione comentários inline para lógica complexa.
    -   **Tratamento de Erros:** Implemente tratamento de erros robusto. Retorne erros apropriados e use `fmt.Errorf` com `%w` para envolver erros quando apropriado.
    -   **Logging:** Adicione logs úteis para depuração e monitoramento (atualmente, o logging é simples; melhorias podem ser uma tarefa futura).
4.  **Execute os Testes e Linters:**
    -   Execute `go test ./...` para garantir que todos os testes passam.
    -   Execute `goimports -w .` para formatar o código.
    -   Execute `golangci-lint run ./...` para verificar problemas de lint. Corrija-os.
    -   Verifique a cobertura de teste (`go test -coverprofile=c.out ./... && go tool cover -html=c.out`).

### 3.3. Corrigindo Bugs
1.  **Reproduza o Bug:** Entenda claramente como o bug ocorre. Se possível, use os testes de integração existentes ou crie um novo para reproduzir o bug.
2.  **Escreva um Teste que Falhe:** Crie ou modifique um teste que demonstre o bug. Este teste deve falhar com o código atual.
3.  **Corrija o Código:** Implemente a correção para o bug.
4.  **Execute os Testes:** Verifique se o novo teste (e todos os outros) agora passa. Certifique-se de que o linter também passa.

### 3.4. Adicionando Comentários no Código
-   Siga as diretrizes de `CONTRIBUTING.MD` e as convenções Godoc.
-   Comente todas as funções, tipos, constantes e variáveis exportadas.
-   Explique o "porquê" da lógica complexa, não apenas o "o quê".

### 3.5. Refatorando Código
-   Certifique-se de que há testes adequados cobrindo o código a ser refatorado. A cobertura de testes é sua rede de segurança.
-   Faça pequenas alterações incrementais e execute os testes (`go test ./...`) frequentemente.
-   O objetivo da refatoração é melhorar a clareza, manutenibilidade ou desempenho sem alterar o comportamento externo observável.

## 4. Padrões e Convenções Específicas do Projeto Vigenda

-   **Estilo de Código:** Siga rigorosamente as diretrizes em `CONTRIBUTING.MD`.
    -   **Formatação:** Use `goimports -w .` ou `gofmt -w .` para formatar seu código antes de commitar.
    -   **Linting:** Execute `golangci-lint run ./...` e corrija quaisquer problemas reportados. Consulte `INSTALLATION.MD` para instalar `golangci-lint`.
-   **Mensagens de Commit:** Use o formato de [Conventional Commits](https://www.conventionalcommits.org/) conforme detalhado em `CONTRIBUTING.MD`.
    -   Ex: `feat(service): adiciona método para calcular progresso da disciplina`
    -   Ex: `fix(tui): corrige exibição de datas em tarefas`
    -   Ex: `docs(readme): atualiza instruções de build`
    -   Ex: `test(repository): adiciona testes para SubjectRepository`
-   **Branching:** Crie branches descritivas para suas tarefas (ex: `feature/TASK-ID-descricao-curta` ou `fix/TASK-ID-bug-especifico`).
-   **Pull Requests (PRs):**
    -   Forneça descrições claras e concisas das alterações no PR.
    -   Referencie as issues relacionadas (ex: `Closes #123`).
    -   Certifique-se de que todos os testes (`go test ./...`) e linters (`golangci-lint run ./...`) passam localmente antes de enviar o PR. (CI/CD ainda não configurado).

## 5. Ferramentas e Comandos Úteis Específicos do Vigenda

-   **Listar arquivos:** `ls -R` (para listagem recursiva e ter uma ideia da estrutura).
-   **Ler arquivos:** Use a ferramenta `read_files(["caminho/do/arquivo.go", "outro/arquivo.md"])`.
-   **Buscar texto em arquivos (grep):**
    -   Para buscar uma definição de função: `grep -R "func CreateTask(" internal/`
    -   Para buscar usos de um modelo: `grep -R "models.Task" internal/`
-   **Executar todos os testes:**
    ```bash
    go test ./...
    ```
    Para saída verbosa:
    ```bash
    go test -v ./...
    ```
-   **Executar testes de um pacote específico (ex: `internal/service`):**
    ```bash
    go test ./internal/service/...
    ```
-   **Executar um teste específico por nome (ex: `TestTaskService_CreateTask_Success` no pacote `service`):**
    ```bash
    go test -v -run ^TestTaskService_CreateTask_Success$ ./internal/service/...
    ```
-   **Gerar e visualizar cobertura de testes:**
    ```bash
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out
    ```
-   **Executar linters (`golangci-lint`):**
    ```bash
    golangci-lint run ./...
    ```
    Para tentar correções automáticas:
    ```bash
    golangci-lint run --fix ./...
    ```
-   **Formatar código (`goimports`):**
    ```bash
    goimports -w .
    ```
-   **Executar a aplicação (a partir da raiz do projeto):**
    ```bash
    go run ./cmd/vigenda/main.go [comando] [flags]
    ```
    Ex: `go run ./cmd/vigenda/main.go tarefa listar`
-   **Construir a aplicação (script):**
    ```bash
    ./build.sh
    ```
    (Os binários estarão em `dist/`)
-   **Construir a aplicação (manual para a plataforma atual):**
    ```bash
    go build -o vigenda_cli ./cmd/vigenda/main.go
    ```

## 6. O que Evitar

-   **Alterar arquivos fora do escopo da tarefa sem uma boa razão e sem documentar claramente no PR.**
-   **Introduzir dependências (`go.mod`) desnecessárias sem discuti-las ou documentá-las na `TECHNICAL_SPECIFICATION.MD`.**
-   **Comentar código funcional em vez de removê-lo (se não for mais necessário). Use o controle de versão para histórico.**
-   **Ignorar falhas de teste ou de linter. Eles existem para ajudar a manter a qualidade.**
-   **Submeter código sem testá-lo adequadamente conforme `TESTING.MD`.**
-   **Escrever mensagens de commit vagas como "correções" ou "atualizações". Siga o padrão Conventional Commits detalhado em `CONTRIBUTING.MD`.**

## 7. Se Você Ficar Preso

1.  **Releia a Tarefa e os Documentos:** Certifique-se de que entendeu completamente o requisito e consultou `README.md`, `INSTALLATION.MD`, `CONTRIBUTING.MD`, `TESTING.MD`, `TECHNICAL_SPECIFICATION.MD` e este `AGENTS.md`.
2.  **Pesquise Erros:** Copie e cole mensagens de erro em um mecanismo de busca ou na sua base de conhecimento.
3.  **Simplifique o Problema:** Tente isolar a parte do código que está causando o problema. Crie um caso de teste mínimo, se aplicável.
4.  **Consulte a Documentação da Linguagem/Framework:** Verifique a documentação oficial do Go (golang.org), Cobra, Bubbletea, etc.
5.  **Peça Ajuda (se aplicável e configurado):** Use a ferramenta `request_user_input` descrevendo o problema, o que você já tentou (referenciando os documentos consultados) e onde está preso.

---
Lembre-se, seu objetivo é produzir código de alta qualidade e documentação que seja fácil de manter e entender. Siga estas diretrizes e use os documentos fornecidos para guiá-lo. Atualize esta documentação e as outras conforme o projeto evolui. Boa codificação!
