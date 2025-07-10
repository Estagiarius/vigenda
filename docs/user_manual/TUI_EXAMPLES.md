# Exemplos de Uso da TUI (Interface de Texto do Usuário) do Vigenda

Este documento fornece exemplos passo a passo de como realizar tarefas comuns utilizando a Interface de Texto do Usuário (TUI) principal do Vigenda. Para iniciar a TUI, execute `vigenda` (ou `go run ./cmd/vigenda/main.go`) no seu terminal.

Use as **teclas de seta** para navegar, **Enter** para selecionar, e **Esc** para voltar ou cancelar.

## 1. Criando sua Primeira Disciplina e Turma

Este fluxo é fundamental para começar a organizar seu conteúdo.

1.  **Inicie o Vigenda:**
    ```bash
    vigenda
    ```
    Você verá o "Menu Principal".

2.  **Criar uma Nova Disciplina:**
    *   No Menu Principal, navegue com as setas até a opção que gerencia disciplinas (ex: "Gerenciar Disciplinas", "Disciplinas") e pressione **Enter**.
    *   Procure por uma opção como "Criar Nova Disciplina" ou "Adicionar Disciplina". Selecione-a.
    *   Será solicitado o nome da disciplina. Digite, por exemplo, `História Antiga` e pressione **Enter**.
    *   A TUI deve confirmar a criação. Anote o ID da disciplina se exibido, pode ser útil. Pressione **Esc** para voltar ao menu de disciplinas ou ao menu principal.

3.  **Criar uma Nova Turma para a Disciplina:**
    *   No Menu Principal, navegue até a opção de gerenciamento de turmas (ex: "Gerenciar Turmas e Alunos", "Turmas").
    *   O sistema pode listar suas disciplinas. Selecione "História Antiga" (a disciplina que você acabou de criar).
    *   Dentro da visualização da disciplina "História Antiga", procure e selecione a opção "Criar Nova Turma".
    *   Digite o nome da turma, por exemplo, `HIS101 - Matutino 2024S1` e pressione **Enter**.
    *   A TUI confirmará a criação. Anote o ID da turma.

## 2. Adicionando Alunos a uma Turma

Depois de criar uma turma, você pode adicionar alunos.

1.  **Navegue até a Turma:**
    *   No Menu Principal da TUI, vá para "Gerenciar Turmas e Alunos".
    *   Selecione a disciplina à qual a turma pertence (ex: "História Antiga").
    *   Selecione a turma desejada (ex: "HIS101 - Matutino 2024S1").

2.  **Adicionar um Aluno:**
    *   Dentro da visualização da turma, procure uma opção como "Gerenciar Alunos" ou "Adicionar Aluno". Selecione-a.
    *   Se houver uma sub-opção "Adicionar Novo Aluno", selecione-a.
    *   Siga os prompts para inserir:
        *   Nome Completo do Aluno (ex: `João da Silva`)
        *   Número de Matrícula/Chamada (opcional, ex: `2024001`)
        *   Status (geralmente 'ativo' por padrão)
    *   Confirme a adição.
    *   Repita para cada aluno que deseja adicionar manualmente.

    **Alternativa:** Para adicionar muitos alunos, use o comando CLI `vigenda turma importar-alunos <ID_DA_TURMA> arquivo.csv`. Veja `EXAMPLES.md` e o Manual do Usuário principal.

## 3. Criando e Gerenciando Tarefas (via TUI)

1.  **Navegue para Gerenciamento de Tarefas:**
    *   No Menu Principal, selecione "Gerenciar Tarefas".

2.  **Criar uma Nova Tarefa:**
    *   Escolha "Criar Nova Tarefa".
    *   Preencha os campos solicitados:
        *   **Título:** (ex: `Ler Capítulo 5 do Livro Texto`)
        *   **Descrição:** (opcional, ex: `Foco nas seções sobre o Império Romano.`)
        *   **Associar à Turma?** (Sim/Não): Se Sim, selecione a disciplina e depois a turma.
        *   **Data de Conclusão:** (opcional, formato AAAA-MM-DD)
    *   Confirme para salvar a tarefa.

3.  **Listar e Visualizar Tarefas:**
    *   Na tela de "Gerenciar Tarefas", você geralmente verá uma lista de tarefas.
    *   Pode haver opções para filtrar por turma, status (pendente, concluída) ou data.
    *   Selecione uma tarefa para ver seus detalhes.

4.  **Marcar Tarefa como Concluída:**
    *   Na lista de tarefas, selecione a tarefa que deseja concluir.
    *   Procure por uma opção como "Marcar como Concluída" ou um atalho de teclado indicado.

5.  **Editar ou Excluir Tarefa:**
    *   Selecione a tarefa desejada na lista.
    *   Procure por opções como "Editar Tarefa" ou "Excluir Tarefa".

## 4. Gerenciando Avaliações (via TUI)

1.  **Navegue para Gerenciamento de Avaliações:**
    *   No Menu Principal, selecione "Gerenciar Avaliações e Notas".

2.  **Criar uma Nova Avaliação:**
    *   Escolha "Criar Nova Avaliação".
    *   Selecione a Disciplina e a Turma para a qual a avaliação se aplica.
    *   Preencha os detalhes da avaliação:
        *   **Nome:** (ex: `Prova Mensal - Unidade 1`)
        *   **Período/Bimestre:** (ex: `1` ou `1º Bimestre`)
        *   **Peso:** (ex: `4.0`)
        *   **Data da Avaliação:** (opcional, AAAA-MM-DD)
    *   Confirme para salvar.

3.  **Lançar Notas para uma Avaliação:**
    *   Na tela de "Gerenciar Avaliações", selecione a avaliação para a qual deseja lançar notas.
    *   Escolha uma opção como "Lançar/Editar Notas".
    *   A TUI listará os alunos da turma associada. Para cada aluno, insira a nota e pressione Enter.
    *   Siga as instruções na tela para salvar ou cancelar.

4.  **Consultar Médias da Turma:**
    *   Ainda em "Gerenciar Avaliações" ou em uma seção específica de "Relatórios", pode haver uma opção para "Calcular Média da Turma".
    *   Selecione a turma para ver as médias calculadas.

## 5. Usando o Banco de Questões e Gerando Provas (via TUI)

1.  **Navegue para o Banco de Questões:**
    *   No Menu Principal, selecione "Banco de Questões".

2.  **Adicionar Questões (se a TUI suportar criação individual):**
    *   Procure uma opção "Adicionar Nova Questão".
    *   Preencha os campos: Disciplina associada, Tópico (opcional), Tipo (Múltipla Escolha, Dissertativa), Dificuldade, Enunciado, Opções (se múltipla escolha) e Resposta Correta.
    *   **Nota:** A importação em lote de questões é mais eficientemente feita pelo comando CLI `vigenda bancoq add arquivo.json`.

3.  **Listar/Filtrar Questões:**
    *   A TUI pode oferecer opções para listar questões, possivelmente filtrando por disciplina, tópico ou dificuldade.

4.  **Navegue para Geração de Provas:**
    *   No Menu Principal, selecione "Gerar Provas".

5.  **Gerar uma Prova:**
    *   Selecione a Disciplina para a qual deseja gerar a prova.
    *   Especifique os critérios:
        *   Tópico (opcional).
        *   Número de questões fáceis, médias e difíceis.
    *   Indique se deseja salvar a prova em um arquivo (e forneça o nome do arquivo) ou apenas visualizá-la.
    *   Confirme para gerar a prova. A prova gerada será exibida na tela ou salva no arquivo especificado.

Estes são apenas alguns exemplos dos fluxos de trabalho comuns na TUI do Vigenda. A melhor forma de aprender é explorando os menus e opções disponíveis diretamente na aplicação. Consulte também o [Manual do Usuário Completo](../user_manual/README.md) para descrições detalhadas de cada funcionalidade.
