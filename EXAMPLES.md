# Exemplos de Uso do Vigenda

Este arquivo fornece exemplos práticos de como usar o Vigenda, focando na interação através da Interface de Texto do Usuário (TUI) e nos comandos CLI disponíveis.

**Nota:**
- Assumimos que você já [instalou o Vigenda e configurou seu ambiente](INSTALLATION.MD).
- Para iniciar a TUI: `go run ./cmd/vigenda/main.go` ou `./vigenda_cli` (se compilado).
- Para comandos CLI: `vigenda <comando> [flags]` (substitua `vigenda` por `go run ./cmd/vigenda/main.go` ou `./vigenda_cli` conforme sua execução).
- IDs (`<disciplina_id>`, `<turma_id>`, etc.) são ilustrativos. Os IDs reais serão gerados pela aplicação.

## 1. Interação Principal via TUI

A maioria das operações de criação, listagem detalhada, edição e remoção de entidades como **Disciplinas, Turmas, Alunos (exceto importação/status), e Aulas** é realizada através da **TUI principal**.

**Exemplo de fluxo na TUI:**

1.  **Iniciar o Vigenda (TUI):**
    ```bash
    vigenda
    ```
2.  **No Menu Principal da TUI:**
    *   Navegue até "Disciplinas" (ou nome similar) e selecione "Criar Nova Disciplina". Preencha os detalhes solicitados.
    *   Com uma disciplina criada (ex: "História", ID: 1), navegue até "Turmas", selecione a disciplina "História" e depois "Criar Nova Turma". Preencha os detalhes (ex: "História 9A"). Anote o ID da turma criada (ex: ID: 1).
    *   Para adicionar alunos manualmente ou gerenciar aulas, explore as opções dentro da respectiva turma/disciplina na TUI.

## 2. Gerenciamento de Tarefas (CLI e TUI)

As tarefas podem ser gerenciadas tanto pela TUI quanto por comandos CLI.

### Criar uma nova tarefa (CLI)
Para a Turma ID `1` (criada via TUI):
```bash
vigenda tarefa add "Revisar Capítulos 1-3 de História Antiga" --classid 1 --duedate 2024-09-15
```
Criar uma tarefa pessoal (sem ID de turma):
```bash
vigenda tarefa add "Comprar canetas novas" --description "Canetas azuis e pretas ponta fina"
```

### Listar tarefas (CLI)
Listar tarefas ativas para a Turma ID `1`:
```bash
vigenda tarefa listar --classid 1
```
Listar todas as tarefas (de todas as turmas e do sistema/bugs):
```bash
vigenda tarefa listar --all
```

### Marcar uma tarefa como concluída (CLI)
Supondo que a tarefa "Revisar Capítulos 1-3" tem ID `5` (obtido ao criar ou listar):
```bash
vigenda tarefa complete 5
```

## 3. Gerenciamento de Alunos (CLI para operações em lote)

As turmas devem ser criadas primeiro via TUI.

### Importar alunos para uma turma (CLI)
Supondo que a Turma "História 9A" (ID `1`) existe:
```bash
vigenda turma importar-alunos 1 lista_alunos_hist9a.csv
```
Consulte o `docs/user_manual/README.md` para o formato do arquivo CSV.

### Atualizar status de um aluno (CLI)
Supondo que o Aluno ID `101` (de uma turma existente) precisa ser atualizado:
```bash
vigenda turma atualizar-status 101 transferido
```

## 4. Gerenciamento de Avaliações (CLI e TUI)

Disciplinas e Turmas devem ser criadas primeiro via TUI.

### Criar uma nova avaliação (CLI)
Para a Turma ID `1` ("História 9A"):
```bash
vigenda avaliacao criar "Prova Mensal - Unidade 1" --classid 1 --term "1º Bimestre" --weight 4.0 --date 2024-08-20
```

### Lançar notas para uma avaliação (CLI, interativo)
Para a Avaliação ID `3` (criada acima). Este comando abrirá um prompt interativo:
```bash
vigenda avaliacao lancar-notas 3
```

### Calcular média da turma (CLI)
Para a Turma ID `1`:
```bash
vigenda avaliacao media-turma 1
```

## 5. Banco de Questões e Provas (CLI)

Disciplinas devem ser criadas primeiro via TUI.

### Adicionar questões ao banco (CLI)
Supondo que o arquivo `historia_questoes.json` existe e a disciplina "História" (ID 1) está cadastrada:
```bash
vigenda bancoq add historia_questoes.json
```
No arquivo JSON, as questões devem referenciar o nome da disciplina existente (ex: "História"). Consulte o `docs/user_manual/README.md` para o formato JSON.

### Gerar uma prova (CLI)
Para a Disciplina "História" (ID `1`):
```bash
vigenda prova gerar --subjectid 1 --easy 5 --medium 3 --hard 2 --output prova_historia_u1.txt
```

## 6. Workflow Combinado (TUI e CLI)

1.  **Na TUI:**
    *   Inicie `vigenda`.
    *   Crie a disciplina "Física Moderna". Anote seu ID (ex: `SubjectID = 2`).
    *   Dentro de "Física Moderna", crie a turma "Física Moderna - T01". Anote seu ID (ex: `ClassID = 3`).

2.  **No CLI:**
    *   Adicione uma tarefa para a turma:
        ```bash
        vigenda tarefa add "Resolver lista de exercícios sobre Relatividade" --classid 3 --duedate 2024-10-01
        ```
        (Suponha que esta tarefa receba o ID `TaskID = 10`)
    *   Importe alunos para a turma:
        ```bash
        vigenda turma importar-alunos 3 alunos_fisica_t01.csv
        ```
    *   Crie uma avaliação:
        ```bash
        vigenda avaliacao criar "P1 - Relatividade Especial" --classid 3 --term "1" --weight 5.0 --date 2024-10-15
        ```
        (Suponha que esta avaliação receba o ID `AssessmentID = 4`)
    *   Lance as notas (será interativo):
        ```bash
        vigenda avaliacao lancar-notas 4
        ```
    *   Adicione questões de Física ao banco (o JSON deve referenciar a disciplina "Física Moderna"):
        ```bash
        vigenda bancoq add questoes_fisica_relatividade.json
        ```
    *   Gere uma prova:
        ```bash
        vigenda prova gerar --subjectid 2 --easy 3 --medium 2 --output prova_relatividade.txt
        ```

Estes exemplos ilustram como combinar a TUI para gerenciamento de entidades base e a CLI para operações rápidas ou em lote. Consulte sempre `vigenda --help` e `vigenda <comando> --help` para a lista mais atualizada de comandos e opções.
