# Guia do Desenvolvedor e Contribuidor do Vigenda

Bem-vindo à seção de desenvolvedores do Vigenda! Este documento fornece uma visão geral da arquitetura do projeto, como configurar o ambiente de desenvolvimento e diretrizes para contribuição.

## Sumário
1.  [Visão Geral da Arquitetura](#visao-geral-da-arquitetura)
    *   [Estrutura de Diretórios Principal](#estrutura-de-diretorios-principal)
    *   [Fluxo de Interação](#fluxo-de-interacao)
2.  [Configuração do Ambiente de Desenvolvimento](#configuracao-do-ambiente-de-desenvolvimento)
    *   [Pré-requisitos](#pre-requisitos)
    *   [Obtendo o Código](#obtendo-o-codigo)
    *   [Compilando](#compilando)
3.  [Executando Testes](#executando-testes)
4.  [Diretrizes para Contribuição](#diretrizes-para-contribuicao)
    *   [Estilo de Código](#estilo-de-codigo)
    *   [Propondo Mudanças](#propondo-mudancas)
    *   [Mensagens de Commit](#mensagens-de-commit)
5.  [Estrutura Detalhada (Exemplos)](#estrutura-detalhada-exemplos)
    *   [`internal/config`](#internalconfig)
    *   [`internal/database`](#internaldatabase)
    *   [`internal/models`](#internalmodels)
    *   [`internal/repository`](#internalrepository)
    *   [`internal/service`](#internalservice)
    *   [`internal/tui`](#internaltui)
    *   [`cmd/vigenda`](#cmdvigenda)

## 1. Visão Geral da Arquitetura

O Vigenda é uma aplicação CLI escrita em Go. Ele segue uma estrutura de projeto que tenta separar as preocupações para facilitar a manutenção e o desenvolvimento.

### Estrutura de Diretórios Principal

A estrutura de diretórios mais relevante para desenvolvedores inclui:

*   **`cmd/vigenda/`**: Contém o ponto de entrada principal da aplicação (`main.go`). É responsável por inicializar configurações, serviços, e o framework CLI (provavelmente Cobra ou similar) que parseia os comandos e argumentos.
*   **`internal/`**: Este diretório contém a maior parte da lógica de negócios e código da aplicação que não se destina a ser importado por outros projetos.
    *   **`config/`**: Gerenciamento de configurações da aplicação (ex: variáveis de ambiente para banco de dados).
    *   **`database/`**: Lógica de conexão com o banco de dados, migrações de esquema SQL.
    *   **`models/`**: Definições das estruturas de dados (structs Go) que representam as entidades do domínio (ex: Tarefa, Turma, Aluno, Avaliação, Questão).
    *   **`repository/`**: Camada de acesso a dados. Contém a lógica para interagir diretamente com o banco de dados (operações CRUD - Create, Read, Update, Delete) para os modelos.
    *   **`service/`**: Camada de serviço que contém a lógica de negócios da aplicação. Coordena as operações, utilizando os repositórios para buscar ou persistir dados e aplicando regras de negócio.
    *   **`tui/`**: (Text User Interface) Componentes e lógica para a interface de usuário no terminal (ex: prompts interativos, tabelas de visualização, dashboard).
*   **`pkg/`**: (Se existir) Código que é seguro para ser importado por projetos externos. (Atualmente, o `ls` não mostrou este diretório, então pode não ser usado).
*   **`tests/`**: Testes de integração e, possivelmente, E2E (End-to-End).
    *   **`integration/`**: Testes que verificam a interação entre diferentes partes do sistema, incluindo o banco de dados.
*   **`scripts/`** ou **`tools/`**: (Se existir) Scripts auxiliares para build, testes, etc. (O `build.sh` está na raiz).
*   **`docs/`**: Documentação do projeto.
*   **`go.mod`, `go.sum`**: Gerenciamento de dependências do Go.
*   **`AGENTS.md`**: Instruções específicas para o agente de desenvolvimento.

### Fluxo de Interação (Exemplo Simplificado)

1.  Usuário executa um comando (ex: `./vigenda tarefa add "Nova tarefa" ...`).
2.  `cmd/vigenda/main.go` e o framework CLI (ex: Cobra) parseiam o comando e os argumentos.
3.  O handler do comando no CLI chama o método apropriado no serviço correspondente (ex: `TaskService.CreateTask(...)`).
4.  O `TaskService` valida os dados de entrada e aplica regras de negócio.
5.  O `TaskService` chama o `TaskRepository` para persistir a nova tarefa no banco de dados.
6.  O `TaskRepository` executa a query SQL (ou usa um ORM) para inserir os dados.
7.  O resultado da operação retorna pela mesma cadeia (Repository -> Service -> CLI Handler).
8.  O CLI Handler (ou um componente TUI chamado pelo handler) exibe a resposta para o usuário.

## 2. Configuração do Ambiente de Desenvolvimento

### Pré-requisitos
*   **Go**: Versão 1.23 ou superior (conforme `AGENTS.md`).
*   **GCC**: Compilador C para CGO (usado pela dependência `go-sqlite3`).
*   **Git**: Para clonar o repositório (se aplicável) e gerenciar versões.
*   (Opcional) **Docker**: Para testes de integração com diferentes bancos de dados ou para criar ambientes isolados.
*   (Opcional) **Make**: Se um `Makefile` for usado para automatizar tarefas comuns.

### Obtendo o Código
Se o projeto estiver em um repositório Git:
```bash
git clone <URL_DO_REPOSITORIO_VIGENDA>
cd vigenda
```
Caso contrário, certifique-se de ter a estrutura de arquivos do projeto.

### Compilando
Para compilar para desenvolvimento local:
```bash
go build -o vigenda_dev ./cmd/vigenda/
```
Isso cria um executável `vigenda_dev` no diretório raiz. Use um nome diferente do de produção (`vigenda`) para evitar conflitos.

O script `build.sh` pode ser usado para compilações de release ou cross-compilação.

## 3. Executando Testes

O Vigenda possui testes unitários (geralmente junto aos arquivos de código, ex: `*_test.go`) e testes de integração (em `tests/integration/`).

Para executar todos os testes unitários e de integração (a partir do diretório raiz do projeto):
```bash
go test ./...
```
*   `./...` instrui o Go a executar testes em todos os pacotes do projeto.

Para executar testes de um pacote específico:
```bash
go test ./internal/service/
```

Para testes de integração, pode ser necessário configurar um banco de dados de teste. Verifique se há instruções específicas em `tests/integration/README.md` ou no código de teste sobre como o banco de dados de teste é gerenciado (ex: criado e destruído automaticamente, ou se requer configuração manual). O `AGENTS.md` menciona que o Vigenda usa SQLite por padrão e pode criar o arquivo `vigenda.db` automaticamente, o que simplifica os testes.

Os testes de integração em `tests/integration/cli_integration_test.go` parecem usar "golden files" para comparar saídas de comandos CLI, o que é uma boa prática.

## 4. Diretrizes para Contribuição

### Estilo de Código
*   Siga as convenções padrão do Go (ex: `gofmt` para formatação).
*   Use `golangci-lint` se o projeto o utilizar para garantir a qualidade do código (verifique se há um arquivo de configuração como `.golangci.yml`).
*   Mantenha as linhas com um comprimento razoável (geralmente abaixo de 100-120 caracteres).
*   Escreva comentários claros e concisos onde necessário. Documente funções públicas.

### Propondo Mudanças
1.  **Crie uma Issue:** Se estiver propondo uma nova funcionalidade ou uma correção de bug significativa, crie uma issue no sistema de rastreamento do projeto (se houver) para discussão.
2.  **Crie um Branch:** Crie um novo branch para suas alterações a partir do branch principal (ex: `main` ou `develop`).
    ```bash
    git checkout -b feature/minha-nova-funcionalidade
    # ou
    git checkout -b fix/meu-bug-fix
    ```
3.  **Desenvolva e Teste:** Implemente suas alterações. Escreva novos testes (unitários e/ou de integração) para cobrir seu código. Certifique-se de que todos os testes existentes continuam passando.
4.  **Faça Commit:** Faça commits atômicos e com mensagens claras (veja abaixo).
5.  **Push e Pull Request:** Envie seu branch para o repositório remoto e abra um Pull Request (PR) para o branch principal. Descreva suas mudanças no PR e referencie a issue original, se houver.

### Mensagens de Commit
Siga um estilo de mensagem de commit consistente. Um formato comum é o [Conventional Commits](https://www.conventionalcommits.org/):
```
<tipo>[escopo opcional]: <descrição>

[corpo opcional]

[rodapé opcional]
```
Exemplos:
*   `feat(tarefa): Adicionar suporte para prioridade de tarefas`
*   `fix(db): Corrigir query de listagem de alunos para postgres`
*   `docs(readme): Atualizar instruções de instalação`
*   `test(service): Adicionar testes unitários para AssessmentService`

## 5. Estrutura Detalhada (Exemplos)

Esta seção dá uma ideia do que esperar em cada diretório principal dentro de `internal/`.

### `internal/config`
(`config.go`)
*   Provavelmente define uma struct `Config`.
*   Funções para carregar configurações de variáveis de ambiente (ex: `VIGENDA_DB_TYPE`, `VIGENDA_DB_DSN`).
*   Pode usar bibliotecas como `godotenv` para carregar arquivos `.env` em desenvolvimento.

### `internal/database`
(`connection.go`, `database.go`, `migrations/001_initial_schema.sql`)
*   `connection.go`: Lógica para estabelecer a conexão com o banco de dados (SQLite, PostgreSQL) usando os DSNs da configuração.
*   `database.go`: Pode conter a interface do banco de dados ou funções de ajuda.
*   `migrations/`: Arquivos SQL para criar/atualizar o esquema do banco. O `001_initial_schema.sql` define a estrutura inicial das tabelas.

### `internal/models`
(`models.go`)
*   Contém as structs Go que representam as entidades do banco de dados.
    ```go
    // Exemplo (pode ser diferente no código real)
    type Task struct {
        ID          int64     `json:"id"`
        Description string    `json:"description"`
        DueDate     time.Time `json:"due_date"`
        ClassID     int64     `json:"class_id"`
        IsCompleted bool      `json:"is_completed"`
        // ... outros campos
    }
    ```

### `internal/repository`
(ex: `task_repository.go`, `class_repository.go`)
*   Define interfaces para as operações de banco de dados (ex: `TaskRepositoryInterface`).
*   Implementações dessas interfaces que usam o objeto de conexão do banco (`sql.DB` ou similar) para executar queries SQL.
    ```go
    // Exemplo de método em uma implementação de repositório
    func (r *taskRepoImpl) GetTaskByID(ctx context.Context, id int64) (*models.Task, error) {
        // Lógica SQL para buscar tarefa por ID
    }
    ```

### `internal/service`
(ex: `task_service.go`, `assessment_service.go`, `*_test.go`)
*   Define interfaces para os serviços (ex: `TaskServiceInterface`).
*   Implementações que orquestram a lógica de negócios, chamando métodos dos repositórios.
*   Contém validações, transformações de dados, etc.
*   Os arquivos `*_test.go` aqui devem focar em testes unitários para a lógica de negócios, usando mocks/stubs para a camada de repositório.

### `internal/tui`
(`prompt.go`, `table.go`, `statusbar.go`, `tui.go`)
*   Componentes para construir a interface de linha de comando.
*   `prompt.go`: Funções para solicitar entrada do usuário de forma interativa.
*   `table.go`: Funções para exibir dados em formato de tabela no console.
*   `statusbar.go`: (Potencialmente) para exibir informações de status.
*   `tui.go`: Pode ser o ponto de entrada ou coordenador dos componentes TUI.

### `cmd/vigenda`
(`main.go`)
*   Inicialização de tudo: configuração, conexão com DB, repositórios, serviços.
*   Configuração da CLI usando uma biblioteca como Cobra: definição de comandos, subcomandos, flags.
*   Cada comando da CLI terá uma função `Run` (ou similar) que chama o método de serviço apropriado.

---
Este guia deve fornecer um bom ponto de partida para entender e contribuir com o Vigenda. Consulte sempre o `AGENTS.md` para quaisquer instruções específicas e explore o código para obter detalhes mais profundos.
