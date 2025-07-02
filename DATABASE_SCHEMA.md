# Documentação do Esquema do Banco de Dados (Vigenda)

Este documento descreve o esquema do banco de dados SQLite utilizado pela aplicação Vigenda. O esquema é definido no arquivo `internal/database/migrations/001_initial_schema.sql`.

## Visão Geral

O banco de dados é projetado para armazenar informações relacionadas a usuários, disciplinas (matérias), turmas, estudantes, aulas, avaliações, notas, tarefas e um banco de questões.

## Tabelas

### 1. `users`

Armazena as informações dos usuários da aplicação.

-   **Propósito:** Gerenciar contas de usuário para acesso à aplicação.
-   **Colunas:**
    -   `id` (INTEGER, PRIMARY KEY AUTOINCREMENT): Identificador único do usuário.
    -   `username` (TEXT, NOT NULL UNIQUE): Nome de usuário para login. Deve ser único.
    -   `password_hash` (TEXT, NOT NULL): Hash da senha do usuário.

### 2. `subjects`

Armazena as disciplinas (matérias) cadastradas pelos usuários.

-   **Propósito:** Organizar o conteúdo e as atividades por disciplina.
-   **Colunas:**
    -   `id` (INTEGER, PRIMARY KEY AUTOINCREMENT): Identificador único da disciplina.
    -   `user_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `users(id)`. Indica a qual usuário a disciplina pertence.
        -   `ON DELETE CASCADE`: Se um usuário for deletado, suas disciplinas também serão.
    -   `name` (TEXT, NOT NULL): Nome da disciplina (ex: "Matemática", "História").

### 3. `classes`

Armazena as turmas associadas a uma disciplina de um usuário.

-   **Propósito:** Agrupar estudantes e atividades dentro de uma disciplina específica.
-   **Colunas:**
    -   `id` (INTEGER, PRIMARY KEY AUTOINCREMENT): Identificador único da turma.
    -   `user_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `users(id)`.
        -   `ON DELETE CASCADE`: Se um usuário for deletado, suas turmas também serão.
    -   `subject_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `subjects(id)`. Indica a qual disciplina esta turma pertence.
        -   `ON DELETE CASCADE`: Se uma disciplina for deletada, suas turmas também serão.
    -   `name` (TEXT, NOT NULL): Nome da turma (ex: "Turma 9A - 2025", "Cálculo I - Engenharia Civil").
    -   `created_at` (TIMESTAMP, NOT NULL, DEFAULT CURRENT_TIMESTAMP): Data e hora de criação do registro.
    -   `updated_at` (TIMESTAMP, NOT NULL, DEFAULT CURRENT_TIMESTAMP): Data e hora da última atualização do registro.

### 4. `students`

Armazena informações sobre os estudantes de uma turma.

-   **Propósito:** Gerenciar a lista de estudantes por turma e registrar seu status.
-   **Colunas:**
    -   `id` (INTEGER, PRIMARY KEY AUTOINCREMENT): Identificador único do estudante.
    -   `class_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `classes(id)`. Indica a qual turma o estudante pertence.
        -   `ON DELETE CASCADE`: Se uma turma for deletada, seus estudantes também serão.
    -   `full_name` (TEXT, NOT NULL): Nome completo do estudante.
    -   `enrollment_id` (TEXT): Número de matrícula ou de chamada do estudante (opcional).
    -   `status` (TEXT, NOT NULL, DEFAULT 'ativo'): Situação do estudante na turma. Valores permitidos incluem 'ativo', 'inativo', 'transferido'.
    -   `created_at` (TIMESTAMP, NOT NULL, DEFAULT CURRENT_TIMESTAMP): Data e hora de criação do registro.
    -   `updated_at` (TIMESTAMP, NOT NULL, DEFAULT CURRENT_TIMESTAMP): Data e hora da última atualização do registro.

### 5. `lessons`

Armazena os planos de aula para cada turma.

-   **Propósito:** Organizar o conteúdo programático e o cronograma das aulas.
-   **Colunas:**
    -   `id` (INTEGER, PRIMARY KEY AUTOINCREMENT): Identificador único da aula.
    -   `class_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `classes(id)`. Indica a qual turma esta aula pertence.
        -   `ON DELETE CASCADE`: Se uma turma for deletada, suas aulas também serão.
    -   `title` (TEXT, NOT NULL): Título da aula.
    -   `plan_content` (TEXT): Conteúdo do plano de aula, preferencialmente em formato Markdown.
    -   `scheduled_at` (TIMESTAMP, NOT NULL): Data e hora agendada para a aula.

### 6. `assessments`

Armazena informações sobre as avaliações (provas, trabalhos) de uma turma.

-   **Propósito:** Gerenciar as avaliações aplicadas, seus pesos e datas.
-   **Colunas:**
    -   `id` (INTEGER, PRIMARY KEY AUTOINCREMENT): Identificador único da avaliação.
    -   `class_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `classes(id)`. Indica a qual turma esta avaliação pertence.
        -   `ON DELETE CASCADE`: Se uma turma for deletada, suas avaliações também serão.
    -   `name` (TEXT, NOT NULL): Nome da avaliação (ex: "Prova Bimestral 1", "Trabalho de História Moderna").
    -   `term` (INTEGER, NOT NULL): Período da avaliação (ex: 1, 2, 3, 4 para bimestres/trimestres).
    -   `weight` (REAL, NOT NULL): Peso da avaliação na composição da nota final (ex: 4.0).
    -   `assessment_date` (DATE): Data da aplicação da avaliação.

### 7. `grades`

Armazena as notas dos estudantes em cada avaliação.

-   **Propósito:** Registrar o desempenho individual dos estudantes nas avaliações.
-   **Colunas:**
    -   `id` (INTEGER, PRIMARY KEY AUTOINCREMENT): Identificador único do registro de nota.
    -   `assessment_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `assessments(id)`. Indica a qual avaliação esta nota se refere.
        -   `ON DELETE CASCADE`: Se uma avaliação for deletada, suas notas também serão.
    -   `student_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `students(id)`. Indica a qual estudante esta nota pertence.
        -   `ON DELETE CASCADE`: Se um estudante for deletado, suas notas também serão.
    -   `grade` (REAL, NOT NULL): Nota obtida pelo estudante na avaliação.

### 8. `tasks`

Armazena tarefas gerais do usuário, que podem ou não estar associadas a uma turma específica.

-   **Propósito:** Gerenciar a lista de afazeres do usuário, incluindo tarefas acadêmicas ou pessoais.
-   **Colunas:**
    -   `id` (INTEGER, PRIMARY KEY AUTOINCREMENT): Identificador único da tarefa.
    -   `user_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `users(id)`. Indica a qual usuário a tarefa pertence.
        -   `ON DELETE CASCADE`: Se um usuário for deletado, suas tarefas também serão.
    -   `class_id` (INTEGER, NULLABLE): Chave estrangeira opcional referenciando `classes(id)`. Permite associar uma tarefa a uma turma específica.
        -   `ON DELETE CASCADE`: Se uma turma for deletada, suas tarefas associadas também serão.
    -   `title` (TEXT, NOT NULL): Título da tarefa.
    -   `description` (TEXT): Descrição detalhada da tarefa.
    -   `due_date` (TIMESTAMP): Data e hora de vencimento da tarefa.
    -   `is_completed` (BOOLEAN, NOT NULL, DEFAULT 0): Indica se a tarefa foi concluída (0 para não, 1 para sim).

### 9. `questions`

Armazena um banco de questões que podem ser usadas em avaliações ou estudos.

-   **Propósito:** Criar um repositório de questões reutilizáveis para elaboração de provas e atividades.
-   **Colunas:**
    -   `id` (INTEGER, PRIMARY KEY AUTOINCREMENT): Identificador único da questão.
    -   `user_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `users(id)`. Indica a qual usuário a questão pertence.
        -   `ON DELETE CASCADE`: Se um usuário for deletado, suas questões também serão.
    -   `subject_id` (INTEGER, NOT NULL): Chave estrangeira referenciando `subjects(id)`. Indica a qual disciplina esta questão está relacionada.
        -   `ON DELETE CASCADE`: Se uma disciplina for deletada, suas questões também serão.
    -   `topic` (TEXT): Tópico específico da disciplina ao qual a questão se refere (opcional).
    -   `type` (TEXT, NOT NULL): Tipo da questão. Valores comuns: 'multipla_escolha', 'dissertativa'.
    -   `difficulty` (TEXT, NOT NULL): Nível de dificuldade da questão. Valores comuns: 'facil', 'media', 'dificil'.
    -   `statement` (TEXT, NOT NULL): Enunciado da questão.
    -   `options` (TEXT): Para questões de múltipla escolha, um array JSON (como string) contendo as opções. Ex: `["Opção A", "Opção B", "Opção C"]`. Para outros tipos, pode ser NULO.
    -   `correct_answer` (TEXT, NOT NULL): Resposta correta da questão. Para múltipla escolha, pode ser o texto da opção correta ou um índice. Para dissertativas, um gabarito ou palavras-chave.

## Relacionamentos Principais (Resumo)

-   Um `user` pode ter várias `subjects`.
-   Uma `subject` (de um `user`) pode ter várias `classes`.
-   Uma `class` pode ter vários `students`.
-   Uma `class` pode ter várias `lessons`.
-   Uma `class` pode ter várias `assessments`.
-   Uma `assessment` pode ter várias `grades` (uma por `student`).
-   Um `user` pode ter várias `tasks`. Uma `task` pode opcionalmente pertencer a uma `class`.
-   Um `user` pode ter várias `questions`. Uma `question` pertence a uma `subject`.

Este esquema forma a base para o gerenciamento de informações acadêmicas no Vigenda. Modificações ou adições futuras ao esquema seriam tratadas através de novos arquivos de migração.
