# Exemplos de Uso da TUI (Interface de Texto do Usuário)

Este documento fornece exemplos passo a passo de como realizar tarefas comuns utilizando a Interface de Texto do Usuário (TUI) do Vigenda.

**Como usar:**
- Para iniciar a TUI, execute `vigenda` (ou `go run ./cmd/vigenda/main.go`) no seu terminal.
- Use as **teclas de seta** para navegar.
- Pressione **Enter** para selecionar uma opção.
- Pressione **Esc** para voltar ao menu anterior ou cancelar uma ação.

---

## 1. Criando sua Primeira Disciplina e Turma

Este é o primeiro passo para organizar seu conteúdo.

1.  **Inicie o Vigenda** para ver o "Menu Principal".
2.  **Navegue até "Gerenciar Turmas e Disciplinas"** e pressione **Enter**.
3.  **Criar uma Nova Disciplina:**
    *   Na tela de disciplinas, selecione a opção **"Criar Nova Disciplina"**.
    *   Digite o nome da disciplina (ex: `História Antiga`) e pressione **Enter**.
    *   A nova disciplina aparecerá na lista. Pressione **Esc** para voltar se necessário.
4.  **Criar uma Nova Turma:**
    *   Na lista de disciplinas, selecione a disciplina que você acabou de criar (ex: "História Antiga") e pressione **Enter**.
    *   Dentro da tela da disciplina, selecione **"Criar Nova Turma"**.
    *   Digite o nome da turma (ex: `HIS101 - Matutino 2024`) e pressione **Enter**.
    *   A nova turma aparecerá na lista.

---

## 2. Adicionando Alunos a uma Turma

1.  **Navegue até a Turma:**
    *   No "Menu Principal", vá para **"Gerenciar Turmas e Disciplinas"**.
    *   Selecione a disciplina (ex: "História Antiga").
    *   Selecione a turma (ex: "HIS101 - Matutino 2024").
2.  **Adicionar um Aluno:**
    *   Dentro da tela da turma, selecione **"Adicionar Novo Aluno"**.
    *   Siga os prompts para inserir o **Nome Completo do Aluno** e, opcionalmente, um número de matrícula.
    *   Confirme a adição. O novo aluno aparecerá na lista de alunos da turma.
    *   Repita para cada aluno.

**Dica:** Para adicionar muitos alunos de uma vez, é mais rápido usar o comando de linha de comando: `vigenda turma importar-alunos <ID_DA_TURMA> arquivo.csv`.

---

## 3. Criando e Gerenciando Tarefas

1.  **Navegue para Tarefas:**
    *   No "Menu Principal", selecione **"Gerenciar Tarefas"**.
2.  **Criar uma Nova Tarefa:**
    *   Escolha **"Criar Nova Tarefa"**.
    *   Preencha os campos solicitados: Título, Descrição, Data de Conclusão e, se desejar, associe a tarefa a uma turma específica.
    *   Confirme para salvar.
3.  **Listar e Marcar como Concluída:**
    *   A tela de tarefas exibirá uma lista de todas as suas tarefas.
    *   Selecione uma tarefa na lista.
    *   Use a opção **"Marcar como Concluída"** (ou um atalho de teclado indicado) para atualizar seu status.

---

## 4. Gerenciando Avaliações e Notas

1.  **Navegue para Avaliações:**
    *   No "Menu Principal", selecione **"Gerenciar Avaliações e Notas"**.
2.  **Selecione a Turma:**
    *   Escolha a Disciplina e, em seguida, a Turma para a qual deseja gerenciar avaliações.
3.  **Criar uma Nova Avaliação:**
    *   Na tela da turma, escolha **"Criar Nova Avaliação"**.
    *   Preencha os detalhes: Nome (ex: `Prova Mensal`), Período (ex: `1º Bimestre`), Peso (ex: `4.0`), e Data.
    *   Confirme para salvar. A avaliação aparecerá na lista.
4.  **Lançar Notas:**
    *   Selecione a avaliação desejada na lista e pressione **Enter**.
    *   Escolha a opção **"Lançar/Editar Notas"**.
    *   A TUI listará os alunos da turma. Para cada um, insira a nota e pressione **Enter**.
    *   Salve as notas no final do processo.
5.  **Consultar Médias:**
    *   Na tela da turma (onde as avaliações são listadas), selecione **"Calcular Média da Turma"** para ver um relatório das médias ponderadas.

---

## 5. Usando o Banco de Questões e Gerando Provas

1.  **Navegue para o Banco de Questões:**
    *   No "Menu Principal", selecione **"Banco de Questões"**. Aqui você pode visualizar as questões que já foram importadas.
    *   **Nota:** A criação de questões é mais eficiente via importação de arquivos JSON pela linha de comando (`vigenda bancoq add arquivo.json`).
2.  **Navegue para Geração de Provas:**
    *   No "Menu Principal", selecione **"Gerar Provas"**.
3.  **Gerar uma Prova:**
    *   **Selecione a Disciplina**.
    *   **Defina os Critérios:** Informe o tópico (opcional), o número de questões por dificuldade (fácil, média, difícil) e o nome do arquivo de saída (ex: `minha_prova.txt`).
    *   Selecione **"Gerar Prova"** e confirme.
    *   A prova será criada no arquivo de texto especificado, pronta para ser impressa.

Estes exemplos cobrem os fluxos de trabalho mais comuns na TUI do Vigenda. Explore os menus para descobrir mais detalhes e opções.
