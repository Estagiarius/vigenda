# Exemplos de Uso do Vigenda

Este arquivo fornece exemplos práticos de como usar o Vigenda, cobrindo tanto a interação via Interface de Texto do Usuário (TUI) quanto os comandos diretos da Interface de Linha de Comando (CLI).

**Nota:**
- Assumimos que você já [instalou o Vigenda e configurou seu ambiente](INSTALLATION.MD).
- Para iniciar a TUI: `go run ./cmd/vigenda/main.go` ou `./vigenda_cli` (se você compilou o binário).
- Para comandos CLI: `vigenda <comando> [flags]` (substitua `vigenda` pela forma como você executa a aplicação).
- IDs (`<disciplina_id>`, `<turma_id>`, etc.) são ilustrativos. Os IDs reais serão gerados e exibidos pela aplicação.

## 1. Interação Principal via TUI (Interface de Texto do Usuário)

A TUI é a forma mais completa e interativa de usar o Vigenda, especialmente para criar e gerenciar entidades como Disciplinas, Turmas, Alunos, Aulas e Avaliações.

**Fluxo de Exemplo: Configurando uma Nova Disciplina e Turma**

1.  **Iniciar o Vigenda (TUI):**
    ```bash
    vigenda
    ```
    Isso abrirá o Menu Principal.

2.  **Criar uma Nova Disciplina:**
    *   No Menu Principal, use as setas para navegar até a opção "Gerenciar Disciplinas" (ou similar) e pressione Enter.
    *   Escolha a opção "Criar Nova Disciplina".
    *   Digite o nome da disciplina quando solicitado (ex: "Biologia Celular") e pressione Enter.
    *   A TUI confirmará a criação e poderá exibir o ID da nova disciplina. Anote-o se precisar para comandos CLI posteriores.

3.  **Criar uma Nova Turma para a Disciplina:**
    *   Volte ao Menu Principal (geralmente com a tecla `Esc`).
    *   Navegue até "Gerenciar Turmas e Alunos".
    *   O sistema poderá listar as disciplinas existentes. Selecione "Biologia Celular".
    *   Escolha "Criar Nova Turma".
    *   Digite o nome da turma (ex: "BIO-101 Manhã 2024S2") e pressione Enter.
    *   A TUI confirmará e exibirá o ID da nova turma. Anote-o.

4.  **Adicionar Alunos à Turma (via TUI):**
    *   Dentro da visualização da turma "BIO-101 Manhã 2024S2", procure uma opção como "Adicionar Aluno" ou "Gerenciar Alunos".
    *   Siga os prompts para inserir o nome completo, número de matrícula (opcional) e status do aluno. Repita para cada aluno.
    *   (Como alternativa para muitos alunos, veja a importação via CSV na seção CLI abaixo).

5.  **Criar um Plano de Aula (via TUI):**
    *   Ainda na visualização da turma "BIO-101 Manhã 2024S2", encontre uma opção como "Gerenciar Aulas" ou "Planos de Aula".
    *   Escolha "Criar Nova Aula".
    *   Preencha o título da aula (ex: "Introdução à Mitocôndria"), o conteúdo do plano (Markdown é suportado) e a data/hora agendada.

**Outras Operações na TUI:**
*   **Dashboard:** Acesse o "Painel de Controle" no Menu Principal para uma visão geral de tarefas e aulas futuras.
*   **Edição/Remoção:** A maioria das entidades criadas pode ser editada ou removida através das opções correspondentes nos menus da TUI.

## 2. Exemplos de Comandos CLI

Os comandos CLI são úteis para operações rápidas, scripts ou quando a interatividade da TUI não é necessária. Muitas operações CLI dependem de IDs de entidades (disciplinas, turmas) que você pode obter através da TUI ou como resultado de outros comandos.

### 2.1. Gerenciamento de Tarefas

#### Criar uma tarefa para uma turma
(Suponha que a Turma ID `1` foi criada via TUI)
```bash
vigenda tarefa add "Estudar Ciclo de Krebs" --classid 1 --duedate 2024-09-20 --description "Ler páginas 45-55 do livro base."
```

#### Criar uma tarefa pessoal (sem associação a turma)
```bash
vigenda tarefa add "Comprar novo caderno para laboratório" --priority alta
```
O sistema pode pedir interativamente a descrição se não fornecida.

#### Listar tarefas
Listar tarefas ativas para a Turma ID `1`:
```bash
vigenda tarefa listar --classid 1
```
Listar todas as tarefas ativas (de todas as turmas e do sistema/bugs):
```bash
vigenda tarefa listar --all
```
A TUI também oferece uma visualização de tarefas, geralmente mais rica.

#### Marcar uma tarefa como concluída
(Suponha que a tarefa "Estudar Ciclo de Krebs" tem ID `5`)
```bash
vigenda tarefa complete 5
```

#### Atualizar uma tarefa
(Suponha que a tarefa ID `5` precisa de uma nova data)
```bash
# Para atualizar via CLI, você precisaria de um comando 'tarefa atualizar'
# que permita modificar campos específicos. Se não existir, a TUI é a alternativa.
# Exemplo conceitual (verifique a disponibilidade do comando com 'vigenda tarefa --help'):
# vigenda tarefa atualizar --id 5 --duedate 2024-09-22
# Atualmente, a atualização detalhada de tarefas é melhor realizada via TUI.
```
A forma mais garantida de atualizar todos os campos de uma tarefa é através da TUI.

#### Deletar uma tarefa
(Suponha que a tarefa ID `7` foi criada por engano)
```bash
# Exemplo conceitual (verifique a disponibilidade do comando com 'vigenda tarefa --help'):
# vigenda tarefa deletar --id 7
# Atualmente, a deleção de tarefas é melhor realizada via TUI.
```

### 2.2. Gerenciamento de Alunos (Operações em Lote)

(Turmas devem ser criadas primeiro via TUI)

#### Importar alunos para uma turma via CSV
Supondo que a Turma "BIO-101 Manhã 2024S2" (ID `1`) existe:
```bash
vigenda turma importar-alunos 1 alunos_bio101.csv
```
O arquivo `alunos_bio101.csv` deve ter colunas como `numero_chamada` (opcional), `nome_completo` (obrigatório), `situacao` (opcional). Consulte `docs/user_manual/README.md` para o formato exato.

#### Atualizar status de um aluno
Supondo que o Aluno ID `101` (previamente importado ou adicionado via TUI na Turma ID `1`) mudou de status:
```bash
vigenda turma atualizar-status 101 transferido
```

### 2.3. Gerenciamento de Avaliações e Notas

(Disciplinas e Turmas devem ser criadas primeiro via TUI)

#### Criar uma nova avaliação
Para a Turma ID `1` ("BIO-101 Manhã 2024S2"):
```bash
vigenda avaliacao criar "Prova Parcial 1 - Estrutura Celular" --classid 1 --term "1" --weight 3.5 --date 2024-09-30
```

#### Lançar notas para uma avaliação (interativo)
Para a Avaliação ID `3` (criada acima). Este comando iniciará um prompt interativo no terminal para inserir as notas de cada aluno da turma:
```bash
vigenda avaliacao lancar-notas 3
```

#### Calcular média da turma
Para a Turma ID `1`. Exibe a média ponderada dos alunos com base nas avaliações e notas lançadas.
```bash
vigenda avaliacao media-turma 1
```

### 2.4. Banco de Questões e Geração de Provas

(Disciplinas devem ser criadas primeiro via TUI)

#### Adicionar questões ao banco a partir de um arquivo JSON
Supondo que o arquivo `questoes_biologia_celular.json` existe e a disciplina "Biologia Celular" (ID `SubjectID=1`, por exemplo) está cadastrada:
```bash
vigenda bancoq add questoes_biologia_celular.json
```
No arquivo JSON, as questões devem referenciar o nome da disciplina existente (ex: `"disciplina": "Biologia Celular"`). Consulte `docs/user_manual/README.md` para o formato JSON detalhado.

#### Gerar uma prova
Para a Disciplina "Biologia Celular" (ID `1`):
```bash
vigenda prova gerar --subjectid 1 --easy 5 --medium 3 --hard 2 --topic "Membrana Plasmática" --output prova_biocel_membrana.txt
```
Isso tentará gerar uma prova com 5 questões fáceis, 3 médias e 2 difíceis sobre "Membrana Plasmática" da disciplina "Biologia Celular" e salvará em `prova_biocel_membrana.txt`.

## 3. Workflow Combinado (TUI e CLI)

Este exemplo mostra como você pode usar a TUI para configurações iniciais e a CLI para tarefas subsequentes.

1.  **Na TUI (`vigenda`):**
    *   Crie a disciplina "Química Orgânica". Anote seu ID (ex: `SubjectID = 3`).
    *   Dentro de "Química Orgânica", crie a turma "QO-202 Tarde". Anote seu ID (ex: `ClassID = 4`).
    *   Adicione alguns alunos à turma "QO-202 Tarde" manualmente ou prepare um CSV para importação.

2.  **No CLI:**
    *   Se preparou um CSV (`alunos_qo202.csv`), importe os alunos:
        ```bash
        vigenda turma importar-alunos 4 alunos_qo202.csv
        ```
    *   Adicione uma tarefa para a turma:
        ```bash
        vigenda tarefa add "Lista de exercícios: Nomenclatura de Alcanos" --classid 4 --duedate 2024-10-10
        ```
    *   Crie uma avaliação:
        ```bash
        vigenda avaliacao criar "Teste 1 - Hidrocarbonetos" --classid 4 --term "1" --weight 3.0 --date 2024-10-20
        ```
        (Suponha que esta avaliação receba o ID `AssessmentID = 5`)
    *   Após a aplicação do teste, lance as notas (será interativo):
        ```bash
        vigenda avaliacao lancar-notas 5
        ```
    *   Adicione questões de Química Orgânica ao banco (o JSON deve referenciar a disciplina "Química Orgânica"):
        ```bash
        vigenda bancoq add questoes_quimica_organica.json
        ```
    *   Gere uma prova para estudo ou futura avaliação:
        ```bash
        vigenda prova gerar --subjectid 3 --medium 7 --hard 3 --output prova_qo_hidrocarbonetos.txt
        ```
    *   Consulte a média da turma:
        ```bash
        vigenda avaliacao media-turma 4
        ```

Estes exemplos visam ilustrar como o Vigenda pode ser utilizado. Para uma lista completa e atualizada de comandos e suas opções, use `vigenda --help` e `vigenda <comando> --help`. A TUI também oferece ajuda contextual em muitas de suas telas.
