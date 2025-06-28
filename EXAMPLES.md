# Exemplos de Uso da CLI Vigenda

Este arquivo fornece exemplos práticos de como usar a Interface de Linha de Comando (CLI) do Vigenda para realizar tarefas comuns e cenários mais avançados.

**Nota:**
- Assumimos que você já [instalou o Vigenda e configurou seu ambiente](INSTALLATION.MD).
- Os comandos podem ser executados usando `go run ./cmd/vigenda/main.go <comando> [flags]` ou, se você construiu o binário, `./vigenda_cli <comando> [flags]` (ou o caminho para o binário em `dist/`). Para simplificar, usaremos `vigenda <comando> [flags]` nos exemplos abaixo.
- Os IDs (`<disciplina_id>`, `<turma_id>`, `<tarefa_id>`, etc.) são exemplos e devem ser substituídos pelos IDs reais gerados pela aplicação ao criar entidades.
- Alguns comandos podem abrir uma Interface de Texto do Usuário (TUI) para entrada de dados interativa se não forem fornecidos todos os argumentos via flags.

## 1. Gerenciamento de Disciplinas (Subjects)

### Criar uma nova disciplina
```bash
vigenda disciplina criar --nome "Cálculo I"
vigenda disciplina criar --nome "História Antiga"
```
(Isso provavelmente abrirá uma TUI para confirmar ou adicionar mais detalhes, ou se houver um comando direto não interativo)

Se o comando for totalmente CLI:
```bash
# Assumindo que o comando `disciplina criar` aceita o nome diretamente
# ou que ele abre uma TUI para entrada de nome.
# O comando exato pode variar dependendo da implementação do Cobra.
# Exemplo conceitual:
# vigenda disciplina criar --nome "Física Quântica"
```
*(Nota: A implementação atual dos comandos `criar` parece focar na TUI. Os exemplos CLI diretos são conceituais se flags diretas não estiverem implementadas para todos os campos.)*

### Listar todas as disciplinas
```bash
vigenda disciplina listar
```

## 2. Gerenciamento de Turmas (Classes)

### Criar uma nova turma para uma disciplina
Primeiro, obtenha o ID da disciplina (ex: após listar as disciplinas). Suponha que "Cálculo I" tem ID `1`.
```bash
# Comando para criar turma (pode abrir TUI)
vigenda turma criar --disciplinaID 1 --nome "Cálculo I - Turma A 2024"
```

### Listar todas as turmas (geral ou por disciplina)
```bash
vigenda turma listar
```
Para listar turmas de uma disciplina específica (ID `1`):
```bash
vigenda turma listar --disciplinaID 1
```

## 3. Gerenciamento de Tarefas (Tasks)

### Criar uma nova tarefa para uma turma
Suponha que a turma "Cálculo I - Turma A 2024" tem ID `1`.
```bash
# Comando para criar tarefa (pode abrir TUI)
vigenda tarefa criar --turmaID 1 --titulo "Lista de Exercícios 1 - Limites" --descricao "Resolver os exercícios da seção 2.3 do livro." --dataEntrega "2024-08-15T23:59:00"
```
Se a data de entrega for opcional ou puder ser definida depois:
```bash
vigenda tarefa criar --turmaID 1 --titulo "Preparar apresentação sobre Derivadas"
```

### Criar uma tarefa pessoal (não associada a uma turma)
```bash
# Comando para criar tarefa (pode abrir TUI)
vigenda tarefa criar --titulo "Comprar livro de Álgebra Linear" --dataEntrega "2024-08-01"
```

### Listar todas as tarefas
```bash
vigenda tarefa listar
```

### Listar tarefas de uma turma específica (ID `1`)
```bash
vigenda tarefa listar --turmaID 1
```

### Listar tarefas concluídas ou pendentes
```bash
vigenda tarefa listar --status pendente
vigenda tarefa listar --status concluida
```

### Marcar uma tarefa como concluída (ID da tarefa `5`)
```bash
vigenda tarefa atualizar --id 5 --concluida
```
Ou, se houver um comando dedicado:
```bash
vigenda tarefa concluir --id 5
```

### Atualizar detalhes de uma tarefa (ID da tarefa `5`)
```bash
# Comando para atualizar tarefa (pode abrir TUI para edição)
vigenda tarefa atualizar --id 5 --titulo "Lista de Exercícios 1 - Limites e Continuidade" --descricao "Resolver os exercícios das seções 2.3 e 2.4."
```

## 4. Sessões de Foco (Focus)

### Iniciar uma sessão de foco para uma tarefa específica
Suponha que a tarefa "Lista de Exercícios 1" tem ID `5`.
```bash
vigenda foco iniciar --tarefaID 5 --duracao 25m
```
Iniciar uma sessão de foco sem associar a uma tarefa:
```bash
vigenda foco iniciar --duracao 45m --atividade "Leitura de Artigo Científico"
```
(A flag `--atividade` é conceitual e dependeria da implementação do comando `foco iniciar`)

### Visualizar histórico de sessões de foco
```bash
vigenda foco historico
```

## 5. Gerenciamento de Avaliações e Notas (Assessments & Grades) - Conceitual

*(Os comandos exatos para avaliações e notas dependerão da implementação detalhada dos comandos Cobra e da TUI. Abaixo estão exemplos conceituais.)*

### Criar uma nova avaliação para uma turma (ID da turma `1`)
```bash
# Pode abrir TUI
vigenda avaliacao criar --turmaID 1 --nome "Prova Bimestral 1" --bimestre 1 --peso 4.0 --data "2024-08-20"
```

### Listar avaliações de uma turma
```bash
vigenda avaliacao listar --turmaID 1
```

### Lançar notas para uma avaliação (ID da avaliação `3`)
Este comando provavelmente abrirá uma TUI para entrada interativa das notas dos alunos da turma associada à avaliação.
```bash
vigenda nota lancar --avaliacaoID 3
```

### Consultar notas de um aluno (ID do aluno `10`) ou de uma turma (ID da turma `1`)
```bash
vigenda nota consultar --alunoID 10
vigenda nota consultar --turmaID 1
```

## 6. Banco de Questões (Questions) - Conceitual

### Adicionar uma nova questão a uma disciplina (ID da disciplina `1`)
```bash
# Pode abrir TUI para detalhes da questão
vigenda questao criar --disciplinaID 1 --tipo "multipla_escolha" --dificuldade "media" --topico "Limites Laterais" --enunciado "Qual o limite de f(x) = 1/x quando x tende a 0 pela direita?" --opcoes '["infinito", "-infinito", "0", "1"]' --respostaCorreta "infinito"
```

### Listar questões de uma disciplina
```bash
vigenda questao listar --disciplinaID 1
```
Filtrar por dificuldade ou tópico:
```bash
vigenda questao listar --disciplinaID 1 --dificuldade "dificil"
vigenda questao listar --disciplinaID 1 --topico "Limites Laterais"
```

## Combinando Comandos (Workflow Example)

1.  **Criar uma disciplina:**
    ```bash
    vigenda disciplina criar --nome "Programação Orientada a Objetos"
    ```
    (Suponha que o ID gerado para "Programação Orientada a Objetos" seja `2`)

2.  **Criar uma turma para esta disciplina:**
    ```bash
    vigenda turma criar --disciplinaID 2 --nome "POO - T01 Manhã 2024.2"
    ```
    (Suponha que o ID gerado para a turma seja `3`)

3.  **Adicionar uma tarefa para esta turma:**
    ```bash
    vigenda tarefa criar --turmaID 3 --titulo "Trabalho Prático 1: Classes e Objetos" --descricao "Implementar sistema de biblioteca em Java." --dataEntrega "2024-09-10T23:59:00"
    ```
    (Suponha que o ID gerado para a tarefa seja `7`)

4.  **Iniciar uma sessão de foco para esta tarefa:**
    ```bash
    vigenda foco iniciar --tarefaID 7 --duracao 1h30m
    ```

5.  **Listar tarefas pendentes da turma:**
    ```bash
    vigenda tarefa listar --turmaID 3 --status pendente
    ```

6.  **Marcar a tarefa como concluída:**
    ```bash
    vigenda tarefa concluir --id 7
    ```

Estes exemplos visam ilustrar o potencial de uso da CLI Vigenda. Consulte a saída de `vigenda --help` e `vigenda <comando> --help` para obter a lista completa de comandos, subcomandos e flags disponíveis.
