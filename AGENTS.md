# Vigenda

**Vigenda** é uma aplicação de linha de comando (CLI) projetada para ajudar professores, especialmente aqueles com TDAH, a organizar tarefas, aulas, avaliações e outras atividades pedagógicas de forma eficiente e focada.

## Funcionalidades Principais

O Vigenda oferece um conjunto de ferramentas para simplificar o dia a dia do professor:

*   **Dashboard Interativo**: Ao iniciar o `vigenda`, visualize rapidamente sua agenda do dia, tarefas mais urgentes e notificações importantes.
*   **Gestão de Tarefas**:
    *   `vigenda tarefa add`: Crie novas tarefas com descrições, prazos e associação a turmas.
    *   `vigenda tarefa listar`: Visualize tarefas pendentes, filtrando por turma.
    *   `vigenda tarefa complete`: Marque tarefas como concluídas.
*   **Gestão de Turmas e Alunos**:
    *   `vigenda turma criar`: Crie novas turmas, associando-as a disciplinas.
    *   `vigenda turma importar-alunos`: Importe listas de alunos para uma turma a partir de um ficheiro CSV.
    *   `vigenda turma atualizar-status`: Modifique o status de um aluno (ativo, inativo, transferido).
*   **Gestão de Avaliações e Notas**:
    *   `vigenda avaliacao criar`: Defina novas avaliações para as turmas, especificando nome, período e peso.
    *   `vigenda avaliacao lancar-notas`: Lance as notas dos alunos de forma interativa para uma avaliação específica.
    *   `vigenda avaliacao media-turma`: Calcule a média geral de uma turma com base nas avaliações e seus pesos.
*   **Banco de Questões e Geração de Provas**:
    *   `vigenda bancoq add`: Importe questões para o seu banco de dados a partir de um ficheiro JSON.
    *   `vigenda prova gerar`: Crie provas personalizadas selecionando questões do banco por disciplina, tópico e nível de dificuldade.
*   **(Futuro) Gestão de Agenda**: `vigenda agenda` para gerenciar eventos e aulas.
*   **(Futuro) Modo Foco**: `vigenda foco iniciar` para sessões de trabalho cronometradas.

## Instalação

### Pré-requisitos

*   **Go**: Versão 1.23 ou superior. Você pode verificar sua versão com `go version`.
*   **GCC**: Um compilador C como o GCC é necessário para a dependência `go-sqlite3`.
    *   Em sistemas Debian/Ubuntu: `sudo apt-get install gcc`
    *   Em macOS: Xcode Command Line Tools (geralmente já instalado ou solicitado ao tentar compilar).
    *   Em Windows: MinGW/TDM-GCC.

### Compilando a Partir do Código Fonte

1.  **Clone o repositório (ou obtenha os arquivos do projeto):**
    ```bash
    # Exemplo se fosse um repositório git
    # git clone https://example.com/vigenda.git
    # cd vigenda
    ```

2.  **Compile o projeto:**
    Navegue até o diretório raiz do projeto onde o `go.mod` está localizado e execute:
    ```bash
    go build -o vigenda ./cmd/vigenda/
    ```
    Isso irá gerar um executável chamado `vigenda` (ou `vigenda.exe` no Windows) no diretório atual.

3.  **(Opcional) Adicione ao PATH:**
    Para usar o `vigenda` de qualquer lugar no seu terminal, mova o executável para um diretório que esteja no seu PATH do sistema (ex: `/usr/local/bin` ou `~/bin` em Linux/macOS) ou adicione o diretório atual ao seu PATH.

### Compilação Cruzada (Cross-Compilation)

O projeto inclui um script `build.sh` para facilitar a compilação cruzada para diferentes sistemas operacionais e arquiteturas.

**Pré-requisitos para Cross-Compilation:**

*   **Para Linux (dentro de um ambiente Linux):** `gcc` (geralmente já instalado).
*   **Para Windows (compilando de Linux):** `mingw-w64`. Instale com `sudo apt-get install mingw-w64`.
*   **Para macOS (compilando de Linux):** A compilação cruzada para macOS a partir do Linux para projetos que usam CGo (como este, devido ao `go-sqlite3`) é complexa e requer um SDK do macOS e um compilador Clang configurado para cross-compilation (ex: via `osxcross`). O script `build.sh` atual não suporta totalmente a compilação para macOS a partir do Linux devido a essas dependências. Recomenda-se compilar para macOS diretamente em uma máquina macOS.

**Usando o script de build:**

1.  **Torne o script executável (se ainda não o fez):**
    ```bash
    chmod +x build.sh
    ```
2.  **Execute o script:**
    ```bash
    ./build.sh
    ```
    Os binários compilados serão colocados no diretório `dist/`, nomeados de acordo com o sistema operacional e arquitetura (ex: `dist/vigenda-linux-amd64`, `dist/vigenda-windows-amd64.exe`).

## Guia de Início Rápido

Aqui estão alguns exemplos de como usar os comandos mais comuns do Vigenda:

1.  **Ver o Dashboard:**
    ```bash
    ./vigenda
    ```

2.  **Adicionar uma nova tarefa para a Turma ID 1:**
    ```bash
    ./vigenda tarefa add "Preparar slides para aula de Segunda Guerra" --classid 1 --duedate 2024-07-20
    ```

3.  **Listar tarefas da Turma ID 1:**
    ```bash
    ./vigenda tarefa listar --classid 1
    ```

4.  **Criar uma nova turma chamada "História 9A" para a Disciplina ID 1:**
    ```bash
    ./vigenda turma criar "História 9A" --subjectid 1
    ```

5.  **Importar alunos para a Turma ID 1 a partir de um arquivo `alunos.csv`:**
    (Consulte "Formatos de Ficheiros de Importação" abaixo para a estrutura do CSV)
    ```bash
    ./vigenda turma importar-alunos 1 alunos.csv
    ```

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

## Licença

Este projeto não possui uma licença de código aberto definida no momento. Todos os direitos são reservados.

## Reporte de Bugs

Para informações sobre como reportar bugs, como eles são analisados e gerenciados, por favor consulte o arquivo [BUG_REPORTING.md](BUG_REPORTING.md).
