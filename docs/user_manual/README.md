# Manual do Usuário do Vigenda

Bem-vindo ao Manual do Usuário do Vigenda! Este documento detalha todas as funcionalidades do Vigenda, ajudando você a organizar suas atividades pedagógicas de forma eficiente.

## Sumário

1.  [Introdução](#introducao)
2.  [Dashboard Interativo](#dashboard-interativo)
3.  [Gestão de Tarefas](#gestao-de-tarefas)
    *   [Adicionar Tarefa (`vigenda tarefa add`)](#adicionar-tarefa-vigenda-tarefa-add)
    *   [Listar Tarefas (`vigenda tarefa listar`)](#listar-tarefas-vigenda-tarefa-listar)
    *   [Completar Tarefa (`vigenda tarefa complete`)](#completar-tarefa-vigenda-tarefa-complete)
4.  [Gestão de Turmas e Alunos](#gestao-de-turmas-e-alunos)
    *   [Criar Turma (`vigenda turma criar`)](#criar-turma-vigenda-turma-criar)
    *   [Importar Alunos (`vigenda turma importar-alunos`)](#importar-alunos-vigenda-turma-importar-alunos)
    *   [Atualizar Status do Aluno (`vigenda turma atualizar-status`)](#atualizar-status-do-aluno-vigenda-turma-atualizar-status)
5.  [Gestão de Avaliações e Notas](#gestao-de-avaliacoes-e-notas)
    *   [Criar Avaliação (`vigenda avaliacao criar`)](#criar-avaliacao-vigenda-avaliacao-criar)
    *   [Lançar Notas (`vigenda avaliacao lancar-notas`)](#lancar-notas-vigenda-avaliacao-lancar-notas)
    *   [Calcular Média da Turma (`vigenda avaliacao media-turma`)](#calcular-media-da-turma-vigenda-avaliacao-media-turma)
6.  [Banco de Questões e Geração de Provas](#banco-de-questoes-e-geracao-de-provas)
    *   [Adicionar Questões ao Banco (`vigenda bancoq add`)](#adicionar-questoes-ao-banco-vigenda-bancoq-add)
    *   [Gerar Prova (`vigenda prova gerar`)](#gerar-prova-vigenda-prova-gerar)
7.  [Formatos de Ficheiros de Importação](#formatos-de-ficheiros-de-importacao)
    *   [Importação de Alunos (CSV)](#importacao-de-alunos-csv)
    *   [Importação de Questões (JSON)](#importacao-de-questoes-json)
8.  [Configuração da Base de Dados](#configuracao-da-base-de-dados)
    *   [Tipos de Base de Dados Suportados](#tipos-de-base-de-dados-suportados)
    *   [Variáveis de Ambiente para Configuração](#variaveis-de-ambiente-para-configuracao)
    *   [Configuração Específica para SQLite](#configuracao-especifica-para-sqlite)
    *   [Configuração Específica para PostgreSQL](#configuracao-especifica-para-postgresql)
    *   [Exemplos de Configuração](#exemplos-de-configuracao)
    *   [Migrações de Esquema](#migracoes-de-esquema-schema-migrations)

## 1. Introdução

O Vigenda é uma aplicação de linha de comando (CLI) projetada para ajudar professores a gerenciar tarefas, aulas, avaliações e outras atividades relacionadas ao ensino. Ele visa oferecer uma maneira focada e eficiente de organização, especialmente útil para aqueles que podem se beneficiar de uma interface direta e sem distrações.

## 2. Dashboard Interativo

Ao executar o Vigenda sem subcomandos, você acessa o Dashboard Interativo:

```bash
./vigenda
```

Este dashboard fornece uma visão geral da sua agenda do dia, tarefas urgentes e outras notificações importantes, permitindo que você comece o dia bem informado.

## 3. Gestão de Tarefas

Gerencie suas tarefas pedagógicas de forma simples e eficaz.

### Adicionar Tarefa (`vigenda tarefa add`)

Crie novas tarefas com descrições, prazos e associação a turmas.

**Uso:**
```bash
./vigenda tarefa add "Descrição da Tarefa" --classid <ID_DA_TURMA> --duedate <AAAA-MM-DD> [outras opções]
```

**Argumentos e Opções:**
*   `"Descrição da Tarefa"`: O texto que descreve a tarefa (obrigatório).
*   `--classid <ID_DA_TURMA>`: O ID numérico da turma à qual esta tarefa está associada.
*   `--duedate <AAAA-MM-DD>`: A data de vencimento da tarefa no formato ano-mês-dia.
*   `--priority <prioridade>`: (Opcional) Nível de prioridade (ex: alta, média, baixa).
*   `--notes "Notas adicionais"`: (Opcional) Observações extras sobre a tarefa.

**Exemplo:**
```bash
./vigenda tarefa add "Preparar apresentação sobre a Revolução Industrial" --classid 1 --duedate 2024-08-15
```
Isso adiciona uma tarefa para a turma com ID 1, com o prazo final em 15 de agosto de 2024.

### Listar Tarefas (`vigenda tarefa listar`)

Visualize tarefas pendentes, com a opção de filtrar por turma.

**Uso:**
```bash
./vigenda tarefa listar [--classid <ID_DA_TURMA>] [--status <status>]
```

**Opções:**
*   `--classid <ID_DA_TURMA>`: (Opcional) Filtra tarefas pela ID da turma especificada.
*   `--status <status>`: (Opcional) Filtra tarefas pelo status (ex: `pendente`, `concluida`, `atrasada`). Por padrão, lista as pendentes.
*   `--all`: (Opcional) Lista todas as tarefas, independentemente do status.

**Exemplo:**
```bash
./vigenda tarefa listar --classid 1
```
Lista todas as tarefas pendentes para a turma com ID 1.

```bash
./vigenda tarefa listar --status concluida
```
Lista todas as tarefas concluídas de todas as turmas.

### Completar Tarefa (`vigenda tarefa complete`)

Marque uma ou mais tarefas como concluídas.

**Uso:**
```bash
./vigenda tarefa complete <ID_DA_TAREFA_1> [ID_DA_TAREFA_2 ...]
```

**Argumentos:**
*   `<ID_DA_TAREFA>`: O ID numérico da tarefa que você deseja marcar como concluída. Você pode fornecer múltiplos IDs separados por espaços.

**Exemplo:**
```bash
./vigenda tarefa complete 42
```
Marca a tarefa com ID 42 como concluída.

```bash
./vigenda tarefa complete 10 11 15
```
Marca as tarefas com IDs 10, 11 e 15 como concluídas.

## 4. Gestão de Turmas e Alunos

Organize suas turmas e os alunos pertencentes a cada uma.

### Criar Turma (`vigenda turma criar`)

Crie novas turmas, associando-as a disciplinas e períodos letivos.

**Uso:**
```bash
./vigenda turma criar "Nome da Turma" --subjectid <ID_DA_DISCIPLINA> [--year <ANO>] [--period <PERIODO>]
```

**Argumentos e Opções:**
*   `"Nome da Turma"`: O nome da turma (obrigatório, ex: "Matemática 9A", "História - Tarde").
*   `--subjectid <ID_DA_DISCIPLINA>`: O ID numérico da disciplina à qual esta turma pertence (obrigatório). (Nota: O sistema deve ter um mecanismo para gerenciar disciplinas e seus IDs).
*   `--year <ANO>`: (Opcional) O ano letivo da turma (ex: 2024).
*   `--period <PERIODO>`: (Opcional) O período/semestre da turma (ex: "1º Semestre", "Anual").

**Exemplo:**
```bash
./vigenda turma criar "Física 2B" --subjectid 3 --year 2024 --period "2º Semestre"
```

### Importar Alunos (`vigenda turma importar-alunos`)

Importe uma lista de alunos para uma turma específica a partir de um ficheiro CSV.

**Uso:**
```bash
./vigenda turma importar-alunos <ID_DA_TURMA> <CAMINHO_DO_ARQUIVO_CSV>
```

**Argumentos:**
*   `<ID_DA_TURMA>`: O ID numérico da turma para a qual os alunos serão importados.
*   `<CAMINHO_DO_ARQUIVO_CSV>`: O caminho para o ficheiro CSV contendo os dados dos alunos.

**Formato do CSV:** Consulte a seção [Importação de Alunos (CSV)](#importacao-de-alunos-csv) para detalhes sobre a estrutura do arquivo.

**Exemplo:**
```bash
./vigenda turma importar-alunos 2 alunos_turma_2B.csv
```

### Atualizar Status do Aluno (`vigenda turma atualizar-status`)

Modifique o status de um aluno dentro de uma turma (ex: ativo, inativo, transferido).

**Uso:**
```bash
./vigenda turma atualizar-status --alunoid <ID_DO_ALUNO> --novostatus <NOVO_STATUS>
```
ou
```bash
./vigenda turma atualizar-status --nomealuno "Nome Completo do Aluno" --turmaid <ID_DA_TURMA> --novostatus <NOVO_STATUS>
```

**Opções:**
*   `--alunoid <ID_DO_ALUNO>`: O ID numérico do aluno cujo status será atualizado.
*   `--nomealuno "Nome Completo do Aluno"`: O nome completo do aluno (usado em conjunto com `--turmaid` se o ID do aluno não for conhecido).
*   `--turmaid <ID_DA_TURMA>`: O ID da turma do aluno (necessário se estiver usando `--nomealuno`).
*   `--novostatus <NOVO_STATUS>`: O novo status para o aluno. Valores permitidos: `ativo`, `inativo`, `transferido`.

**Exemplo:**
```bash
./vigenda turma atualizar-status --alunoid 101 --novostatus transferido
```
Atualiza o status do aluno com ID 101 para "transferido".

```bash
./vigenda turma atualizar-status --nomealuno "Maria Silva" --turmaid 2 --novostatus inativo
```
Atualiza o status de "Maria Silva" na turma com ID 2 para "inativo".

## 5. Gestão de Avaliações e Notas

Defina avaliações, lance notas e calcule médias.

### Criar Avaliação (`vigenda avaliacao criar`)

Defina novas avaliações para as turmas, especificando nome, período e peso.

**Uso:**
```bash
./vigenda avaliacao criar "Nome da Avaliação" --classid <ID_DA_TURMA> --term <PERIODO_AVALIATIVO> --weight <PESO> [--date <AAAA-MM-DD>]
```

**Argumentos e Opções:**
*   `"Nome da Avaliação"`: O nome da avaliação (obrigatório, ex: "Prova Mensal - Unidade 1", "Trabalho em Grupo").
*   `--classid <ID_DA_TURMA>`: O ID numérico da turma para a qual esta avaliação se aplica (obrigatório).
*   `--term <PERIODO_AVALIATIVO>`: O período avaliativo ao qual esta avaliação pertence (ex: "1º Bimestre", "Trimestre Final").
*   `--weight <PESO>`: O peso desta avaliação no cálculo da média final (ex: 1.0, 2.5, 3.0).
*   `--date <AAAA-MM-DD>`: (Opcional) A data em que a avaliação será aplicada ou entregue.

**Exemplo:**
```bash
./vigenda avaliacao criar "Seminário sobre Ecossistemas" --classid 3 --term "2º Trimestre" --weight 2.0 --date 2024-09-10
```

### Lançar Notas (`vigenda avaliacao lancar-notas`)

Lance as notas dos alunos de forma interativa para uma avaliação específica. O sistema listará os alunos da turma associada à avaliação e solicitará a nota para cada um.

**Uso:**
```bash
./vigenda avaliacao lancar-notas <ID_DA_AVALIACAO>
```

**Argumentos:**
*   `<ID_DA_AVALIACAO>`: O ID numérico da avaliação para a qual as notas serão lançadas.

**Exemplo:**
```bash
./vigenda avaliacao lancar-notas 5
```
Ao executar, o sistema iniciará um prompt interativo para inserir as notas dos alunos para a avaliação de ID 5.

### Calcular Média da Turma (`vigenda avaliacao media-turma`)

Calcule a média geral de uma turma com base nas avaliações e seus pesos definidos para um determinado período avaliativo ou para o resultado final.

**Uso:**
```bash
./vigenda avaliacao media-turma --classid <ID_DA_TURMA> [--term <PERIODO_AVALIATIVO>]
```

**Opções:**
*   `--classid <ID_DA_TURMA>`: O ID numérico da turma para a qual a média será calculada (obrigatório).
*   `--term <PERIODO_AVALIATIVO>`: (Opcional) Especifica o período avaliativo para o qual a média deve ser calculada. Se omitido, pode calcular uma média geral ou final, dependendo da implementação.

**Exemplo:**
```bash
./vigenda avaliacao media-turma --classid 3 --term "2º Trimestre"
```
Calcula e exibe as médias dos alunos da turma ID 3 para o 2º Trimestre.

## 6. Banco de Questões e Geração de Provas

Mantenha um banco de questões organizado e gere provas personalizadas.

### Adicionar Questões ao Banco (`vigenda bancoq add`)

Importe questões para o seu banco de dados a partir de um ficheiro JSON.

**Uso:**
```bash
./vigenda bancoq add <CAMINHO_DO_ARQUIVO_JSON>
```

**Argumentos:**
*   `<CAMINHO_DO_ARQUIVO_JSON>`: O caminho para o ficheiro JSON contendo as questões.

**Formato do JSON:** Consulte a seção [Importação de Questões (JSON)](#importacao-de-questoes-json) para detalhes sobre a estrutura do arquivo.

**Exemplo:**
```bash
./vigenda bancoq add /data/questoes/historia_moderna.json
```

### Gerar Prova (`vigenda prova gerar`)

Crie provas personalizadas selecionando questões do banco por disciplina, tópico e nível de dificuldade.

**Uso:**
```bash
./vigenda prova gerar --subjectid <ID_DA_DISCIPLINA> [--topic "Tópico"] [--easy <NUM>] [--medium <NUM>] [--hard <NUM>] [--output <ARQUIVO_SAIDA>]
```

**Opções:**
*   `--subjectid <ID_DA_DISCIPLINA>`: O ID da disciplina para a qual a prova será gerada (obrigatório).
*   `--topic "Tópico"`: (Opcional) Filtra questões por um tópico específico dentro da disciplina.
*   `--easy <NUM>`: (Opcional) Número de questões de dificuldade "fácil" a incluir.
*   `--medium <NUM>`: (Opcional) Número de questões de dificuldade "média" a incluir.
*   `--hard <NUM>`: (Opcional) Número de questões de dificuldade "difícil" a incluir.
*   `--total <NUM>`: (Opcional) Número total de questões para a prova (o sistema tentará balancear as dificuldades se os números específicos não forem fornecidos).
*   `--output <ARQUIVO_SAIDA>`: (Opcional) Nome do arquivo onde a prova gerada será salva (ex: `prova_historia_01.txt`). Se omitido, a prova pode ser exibida no console.

**Exemplo:**
```bash
./vigenda prova gerar --subjectid 1 --easy 5 --medium 3 --hard 2 --output prova_final_hist_9A.txt
```
Gera uma prova para a disciplina com ID 1, contendo 5 questões fáceis, 3 médias e 2 difíceis, salvando o resultado em `prova_final_hist_9A.txt`.

## 7. Formatos de Ficheiros de Importação

Esta seção detalha os formatos esperados para importação de dados.

### Importação de Alunos (CSV)

O comando `vigenda turma importar-alunos` espera um ficheiro CSV com as seguintes colunas:

*   `numero_chamada` (opcional): Número de chamada do aluno.
*   `nome_completo`: Nome completo do aluno (obrigatório).
*   `situacao` (opcional): Status do aluno. Valores permitidos: `ativo`, `inativo`, `transferido`. Se omitido, o padrão é `ativo`.

**Exemplo (`alunos.csv`):**
```csv
numero_chamada,nome_completo,situacao
1,"Ana Beatriz Costa","ativo"
2,"Bruno Dias","ativo"
,"Carlos Eduardo Lima",
4,"Daniel Mendes","transferido"
```
No exemplo acima, "Carlos Eduardo Lima" terá `numero_chamada` nulo e `situacao` definida como `ativo` por padrão.

### Importação de Questões (JSON)

O comando `vigenda bancoq add` espera um ficheiro JSON contendo uma lista (array) de objetos, onde cada objeto representa uma questão.

**Estrutura de cada objeto de questão:**

*   `disciplina` (string, obrigatório): Nome da disciplina à qual a questão pertence (Ex: "História").
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

## 8. Configuração da Base de Dados

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

*   `VIGENDA_DB_HOST`: Endereço do servidor PostgreSQL. (Padrão: `localhost`)
*   `VIGENDA_DB_PORT`: Porta do servidor PostgreSQL. (Padrão: `5432`)
*   `VIGENDA_DB_USER`: Nome de utilizador para a conexão. (Obrigatório)
*   `VIGENDA_DB_PASSWORD`: Senha para o utilizador.
*   `VIGENDA_DB_NAME`: Nome da base de dados PostgreSQL. (Obrigatório)
*   `VIGENDA_DB_SSLMODE`: Modo de SSL para a conexão PostgreSQL. (Padrão: `disable`)

### Exemplos de Configuração

#### SQLite (Caminho Personalizado)
```bash
export VIGENDA_DB_TYPE="sqlite"
export VIGENDA_DB_PATH="/var/data/vigenda_production.db"
./vigenda
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
export VIGENDA_DB_DSN="postgresql://vigenda_user:super_secret_password@my.postgres.server.com:5433/vigenda_prod_db?sslmode=require"
./vigenda
```
(Nota: Ao usar DSN, `VIGENDA_DB_TYPE` pode ser omitido se o driver puder inferir o tipo.)

### Migrações de Esquema (Schema Migrations)

*   **SQLite**: O Vigenda tentará aplicar o esquema inicial (`internal/database/migrations/001_initial_schema.sql`) automaticamente se a base de dados parecer vazia.
*   **PostgreSQL**: As migrações de esquema devem ser geridas externamente. O Vigenda não tentará criar tabelas ou modificar o esquema numa base de dados PostgreSQL existente. Certifique-se de que o esquema apropriado já foi aplicado.

---

Para mais informações sobre como instalar e começar a usar o Vigenda rapidamente, consulte o [Guia de Introdução](../getting_started/README.md).
Se tiver perguntas comuns, visite nosso [FAQ](../faq/README.md).
Para exemplos práticos, explore nossos [Tutoriais](../tutorials/README.md).
