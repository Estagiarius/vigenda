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
    -   `TECHNICAL_SPECIFICATION.MD`: Detalhes da arquitetura em camadas, fluxo de dados (com diagramas Mermaid), descrição detalhada de módulos internos, padrões de design. Consulte para entender a organização.
    -   `DATABASE_SCHEMA.md`: Documentação detalhada do esquema do banco de dados SQLite, incluindo tabelas, colunas e relacionamentos. **Consulte ao trabalhar com repositórios ou modelos.**
    -   `INSTALLATION.MD`: Instruções detalhadas para configurar o ambiente de desenvolvimento (Go, GCC, ferramentas) e seção expandida de solução de problemas.
    -   `TESTING.MD`: Como executar testes unitários, de integração, benchmarks e linters. Inclui comandos específicos.
    -   `EXAMPLES.md`: Exemplos práticos de uso da CLI Vigenda para diversos cenários.
    -   `CONTRIBUTING.MD`: Diretrizes para contribuição (humanos e IA), padrões de codificação Go (incluindo Godoc), formato de commit (Conventional Commits), processo de PR. **SIGA RIGOROSAMENTE.**
    -   `CODE_OF_CONDUCT.md`: Código de Conduta para Contribuidores.
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
        -   `internal/app/`: Contém a lógica principal da Interface de Texto do Usuário (TUI) baseada em BubbleTea.
            -   `internal/app/app.go`: Modelo principal da aplicação TUI, gerenciando as diferentes visualizações (views).
            -   `internal/app/views.go`: Define as diferentes visualizações/módulos da TUI.
            -   `internal/app/[nome_do_modulo]/[nome_do_modulo].go`: Padrão para componentes de TUI modulares (ex: `internal/app/dashboard/dashboard.go`, `internal/app/tasks/tasks.go`). Cada módulo geralmente implementa seu próprio `Model`, `Init`, `Update`, `View`.
        -   `internal/tui/`: Pode conter componentes TUI mais genéricos ou uma estrutura TUI mais antiga/alternativa (ex: `prompt.go`, `table.go`). O desenvolvimento TUI principal atual está focado em `internal/app/`.
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
    -   Identifique os arquivos e módulos que precisarão ser modificados ou criados (ex: novo comando Cobra, novo método de serviço, alterações no repositório, novo componente TUI como `internal/app/meu_novo_modulo/meu_novo_modulo.go`).
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


6.  **Criar uma avaliação para a Turma ID 1:**
    ```bash
    ./vigenda avaliacao criar "Prova Mensal - Unidade 1" --classid 1 --term 1 --weight 3.0
    ```

7.  **Lançar notas para a Avaliação ID 1 (será interativo):**
    ```bash
    ./vigenda avaliacao lancar-notas 1
    ```

8.  **Adicionar questões de um arquivo `historia_questoes.json` ao banco:**
    (Consulte "Formatos de Ficheiros de Importação" abaixo para a estrutura do JSON)
    ```bash
    ./vigenda bancoq add historia_questoes.json
    ```

9.  **Gerar uma prova para a Disciplina ID 1 com 5 questões fáceis e 3 médias:**
    ```bash
    ./vigenda prova gerar --subjectid 1 --easy 5 --medium 3
    ```

Para obter ajuda sobre qualquer comando específico e suas opções, use a flag `--help`:
```bash
./vigenda tarefa add --help
./vigenda turma importar-alunos --help
```

## Formatos de Ficheiros de Importação

### 1. Importação de Alunos (CSV)

O comando `vigenda turma importar-alunos` espera um ficheiro CSV com as seguintes colunas:

*   `numero_chamada` (opcional): Número de chamada do aluno.
*   `nome_completo`: Nome completo do aluno (obrigatório).
*   `situacao` (opcional): Status do aluno. Valores permitidos: `ativo`, `inativo`, `transferido`. Se omitido, o padrão é `ativo`.

**Exemplo (`alunos.csv`):**
```csv
numero_chamada,nome_completo,situacao
1,"Ana Beatriz Costa","ativo"
2,"Bruno Dias","ativo"
,"Carlos Eduardo Lima", # numero_chamada omitido, situacao será 'ativo' por padrão
4,"Daniel Mendes","transferido"
```

### 2. Importação de Questões (JSON)

O comando `vigenda bancoq add` espera um ficheiro JSON contendo uma lista (array) de objetos, onde cada objeto representa uma questão.

**Estrutura de cada objeto de questão:**

*   `disciplina` (string, obrigatório): Nome da disciplina à qual a questão pertence (Ex: "História"). *Nota: O sistema tentará encontrar um ID de disciplina correspondente. Idealmente, o sistema deveria permitir referenciar por ID de disciplina diretamente ou ter um mecanismo para criar/mapear disciplinas.*
*   `topico` (string, opcional): Tópico específico da questão dentro da disciplina (Ex: "Revolução Francesa").
*   `tipo` (string, obrigatório): Tipo da questão. Valores permitidos: `multipla_escolha`, `dissertativa`.
*   `dificuldade` (string, obrigatório): Nível de dificuldade. Valores permitidos: `facil`, `media`, `dificil`.
*   `enunciado` (string, obrigatório): O texto da questão.
*   `opcoes` (array de strings, obrigatório para `multipla_escolha`): Uma lista das opções de resposta.
*   `resposta_correta` (string, obrigatório): O texto da resposta correta. Para `multipla_escolha`, deve corresponder exatamente a uma das `opcoes`.

**Exemplo (`questoes.json`):**
```json
[
  {
    "disciplina": "História",
    "topico": "Revolução Francesa",
    "tipo": "multipla_escolha",
    "dificuldade": "media",
    "enunciado": "Qual destes eventos é considerado o estopim da Revolução Francesa?",
    "opcoes": [
      "A Queda da Bastilha",
      "A convocação dos Estados Gerais",
      "O Juramento da Quadra de Tênis"
    ],
    "resposta_correta": "A Queda da Bastilha"
  },
  {
    "disciplina": "Matemática",
    "topico": "Álgebra",
    "tipo": "dissertativa",
    "dificuldade": "facil",
    "enunciado": "Explique o que é uma equação de primeiro grau e dê um exemplo.",
    "resposta_correta": "Uma equação de primeiro grau é uma igualdade que envolve uma ou mais incógnitas com expoente 1. Exemplo: 2x + 3 = 7."
  }
]
```

## Configuração da Base de Dados

O Vigenda suporta diferentes tipos de bases de dados, configuráveis através de variáveis de ambiente.

### Tipos de Base de Dados Suportados

*   **SQLite** (padrão): Leve, baseada em ficheiro, ideal para uso individual.
*   **PostgreSQL**: Robusta, para cenários com múltiplos utilizadores ou maior volume de dados.

### Variáveis de Ambiente para Configuração

As seguintes variáveis de ambiente podem ser usadas para configurar a conexão com a base de dados:

*   `VIGENDA_DB_TYPE`: Especifica o tipo de base de dados.
    *   Valores: `sqlite` (padrão), `postgres`.
*   `VIGENDA_DB_DSN`: Uma string de conexão (Data Source Name) completa. Se esta variável for definida, ela tem precedência sobre as variáveis individuais abaixo.
    *   **Exemplo SQLite DSN**: `file:/caminho/absoluto/para/meu_vigenda.db?cache=shared&mode=rwc`
    *   **Exemplo PostgreSQL DSN**: `postgres://utilizador:senha@localhost:5432/nome_da_base?sslmode=disable`

#### Configuração Específica para SQLite

Se `VIGENDA_DB_TYPE` for `sqlite` (ou não estiver definida) e `VIGENDA_DB_DSN` não for fornecida, a seguinte variável é usada:

*   `VIGENDA_DB_PATH`: Caminho para o ficheiro da base de dados SQLite.
    *   **Padrão**: Um ficheiro `vigenda.db` no diretório de configuração do utilizador (ex: `~/.config/vigenda/vigenda.db` no Linux) ou no diretório atual se o diretório de configuração não for acessível.
    *   **Exemplo**: `export VIGENDA_DB_PATH="/caminho/para/sua/vigenda.db"`

#### Configuração Específica para PostgreSQL

Se `VIGENDA_DB_TYPE` for `postgres` e `VIGENDA_DB_DSN` não for fornecida, as seguintes variáveis são usadas para construir a DSN:

*   `VIGENDA_DB_HOST`: Endereço do servidor PostgreSQL.
    *   Padrão: `localhost`
*   `VIGENDA_DB_PORT`: Porta do servidor PostgreSQL.
    *   Padrão: `5432`
*   `VIGENDA_DB_USER`: Nome de utilizador para a conexão. (Obrigatório)
*   `VIGENDA_DB_PASSWORD`: Senha para o utilizador. (Pode ser vazia se o método de autenticação permitir)
*   `VIGENDA_DB_NAME`: Nome da base de dados PostgreSQL. (Obrigatório)
*   `VIGENDA_DB_SSLMODE`: Modo de SSL para a conexão PostgreSQL.
    *   Padrão: `disable`
    *   Outros valores comuns: `require`, `verify-ca`, `verify-full`.

### Exemplos de Configuração

#### SQLite (Caminho Personalizado)

Se você quiser usar SQLite mas num local específico:
```bash
export VIGENDA_DB_TYPE="sqlite"
export VIGENDA_DB_PATH="/var/data/vigenda_production.db"
./vigenda
```
Ou, de forma mais concisa, se `VIGENDA_DB_TYPE` for omitido (assume SQLite):
```bash
export VIGENDA_DB_PATH="/var/data/vigenda_production.db"
./vigenda ...
```

#### PostgreSQL (Usando Variáveis Individuais)

```bash
export VIGENDA_DB_TYPE="postgres"
export VIGENDA_DB_HOST="my.postgres.server.com"
export VIGENDA_DB_PORT="5433"
export VIGENDA_DB_USER="vigenda_user"
export VIGENDA_DB_PASSWORD="super_secret_password"
export VIGENDA_DB_NAME="vigenda_prod_db"
export VIGENDA_DB_SSLMODE="require"
./vigenda
```

#### PostgreSQL (Usando DSN Completa)

```bash
export VIGENDA_DB_TYPE="postgres" # Ou pode ser inferido pela DSN se o driver souber
export VIGENDA_DB_DSN="postgresql://vigenda_user:super_secret_password@my.postgres.server.com:5433/vigenda_prod_db?sslmode=require"
./vigenda
```

**Nota sobre Migrações de Esquema (Schema Migrations):**
*   Para **SQLite**, o Vigenda tentará aplicar o esquema inicial (`001_initial_schema.sql`) automaticamente se a base de dados parecer vazia (ex: a tabela `users` não existir).
*   Para **PostgreSQL**, as migrações de esquema devem ser geridas externamente (ex: usando ferramentas como `goose`, `migrate`, `Flyway`, ou scripts SQL manuais). O Vigenda não tentará criar tabelas ou modificar o esquema numa base de dados PostgreSQL existente. Certifique-se de que o esquema definido em `internal/database/migrations/001_initial_schema.sql` (ou uma versão compatível) já foi aplicado à sua base de dados PostgreSQL antes de executar a aplicação.

Por padrão, o Vigenda cria e utiliza um ficheiro de base de dados SQLite chamado `vigenda.db` no diretório de configuração do utilizador ou no diretório atual.

Você pode especificar um caminho diferente para o ficheiro da base de dados SQLite (se estiver a usar SQLite e não uma DSN completa) definindo a variável de ambiente `VIGENDA_DB_PATH`:

```bash
export VIGENDA_DB_PATH="/caminho/para/sua/vigenda.db"
./vigenda ...
```

## Contribuições

Este projeto é atualmente mantido para um propósito específico. No entanto, sugestões e discussões sobre melhorias são bem-vindas (se um canal de comunicação for estabelecido, como issues em um repositório Git).

## Documentação do Usuário

Para um guia completo sobre como usar o Vigenda, incluindo detalhes sobre todos os comandos, configuração e exemplos práticos, consulte nossa documentação do usuário:

*   **[Manual do Usuário](./docs/user_manual/README.md)**: Um guia detalhado sobre todas as funcionalidades.
*   **[Guia de Introdução](./docs/getting_started/README.md)**: Para uma instalação rápida e os primeiros passos.
*   **[FAQ (Perguntas Frequentes)](./docs/faq/README.md)**: Respostas para as dúvidas mais comuns.
*   **[Tutoriais](./docs/tutorials/README.md)**: Exemplos práticos passo a passo.
*   **[Guia do Desenvolvedor](./docs/developer/README.md)**: Informações sobre a arquitetura do projeto, como configurar o ambiente de desenvolvimento e diretrizes para contribuição (para desenvolvedores).

## Licença

Este projeto não possui uma licença de código aberto definida no momento. Todos os direitos são reservados.

## Reporte de Bugs

Para informações sobre como reportar bugs, como eles são analisados e gerenciados, por favor consulte o arquivo [BUG_REPORTING.md](BUG_REPORTING.md).

Lembre-se, seu objetivo é produzir código de alta qualidade e documentação que seja fácil de manter e entender. Siga estas diretrizes e use os documentos fornecidos para guiá-lo. Atualize esta documentação e as outras conforme o projeto evolui. Boa codificação!

