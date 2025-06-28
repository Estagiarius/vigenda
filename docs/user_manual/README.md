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
Ele pode exibir informações como:
*   **Tarefas com prazo para hoje ou atrasadas:** Um lembrete visual das suas prioridades imediatas.
*   **Próximos eventos ou aulas:** (Se a funcionalidade de agenda estiver implementada) Um vislumbre do que está por vir.
*   **Notificações do sistema:** Alertas sobre importações de alunos concluídas, erros na última operação, etc.
*   **Resumo de progresso:** (Potencialmente) Percentual de tarefas concluídas na semana, média de notas recentes, etc.

A aparência e o conteúdo exato podem variar dependendo da versão e configuração do Vigenda.

### Dicas de Uso para o Dashboard
*   **Comece seu dia aqui:** Execute `vigenda` como seu primeiro comando para ter um panorama rápido.
*   **Verifique regularmente:** Se você usa o Vigenda ao longo do dia, revisitar o dashboard pode ajudar a manter o foco.

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

#### Exemplos Avançados e Dicas (`tarefa add`):
*   **Adicionar uma tarefa de alta prioridade com notas detalhadas:**
    ```bash
    ./vigenda tarefa add "Finalizar correção das provas - URGENTE" --classid 2 --duedate 2024-07-25 --priority alta --notes "Lembrar de verificar questões dissertativas com critério X. Publicar notas até dia 26."
    ```
*   **Adicionar uma tarefa recorrente (conceitualmente, Vigenda pode não suportar diretamente, mas você pode adaptá-la):**
    Para tarefas que se repetem, como "Preparar aula de segunda-feira", você pode adicionar a primeira ocorrência e, ao completá-la, adicionar a próxima manualmente ou usar um script externo.
    ```bash
    ./vigenda tarefa add "Preparar aula de Segunda (Semana 1)" --classid 1 --duedate 2024-08-05
    ```

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

#### Exemplos Avançados e Dicas (`tarefa listar`):
*   **Listar todas as tarefas de uma turma específica, incluindo concluídas e pendentes:**
    ```bash
    ./vigenda tarefa listar --classid 1 --all
    ```
*   **Listar todas as tarefas com prioridade "alta" (assumindo que o campo prioridade é filtrável, o que pode não ser o caso - verificar com `--help`):**
    ```bash
    ./vigenda tarefa listar --priority alta
    ```
    Se não houver filtro direto por prioridade, você pode precisar listar todas e filtrar visualmente ou com ferramentas de linha de comando como `grep`.
*   **Listar tarefas com prazo para os próximos 7 dias (requer funcionalidade mais avançada de filtro por data):**
    O Vigenda pode não ter um filtro de data tão granular. Se não, você listaria as tarefas pendentes e verificaria as datas manualmente.
    ```bash
    ./vigenda tarefa listar --status pendente
    ```

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

#### Exemplos Avançados e Dicas (`tarefa complete`):
*   **Completar uma tarefa e imediatamente adicionar uma tarefa de acompanhamento:**
    Não há um comando direto para isso, mas você pode encadear comandos no seu shell (se suportado):
    ```bash
    ./vigenda tarefa complete 25 && ./vigenda tarefa add "Revisar feedback da tarefa 25" --classid 1 --duedate 2024-08-01
    ```
    (O `&&` significa que o segundo comando só roda se o primeiro for bem-sucedido).
*   **Marcar uma tarefa como "não concluída" ou "reabrir" (se suportado):**
    O Vigenda pode ter um comando como `tarefa reabrir <ID>` ou `tarefa atualizar --id <ID> --status pendente`. Se não, você pode precisar remover a tarefa e adicioná-la novamente ou editar o banco de dados diretamente (não recomendado para usuários comuns). Verifique com `./vigenda tarefa --help`.

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

#### Exemplos Avançados e Dicas (`turma criar`):
*   **Criar múltiplas turmas para a mesma disciplina rapidamente:**
    Você pode usar um loop de shell se estiver criando várias turmas com um padrão (ex: Turma A, B, C):
    ```bash
    # Exemplo para shell bash/zsh
    for LETRA in A B C; do \
      ./vigenda turma criar "Matemática 1${LETRA}" --subjectid 2 --year 2024; \
    done
    ```
*   **Verificar se uma turma já existe antes de criar:**
    O Vigenda pode retornar um erro se a turma já existir. Não há um comando `turma existe` explícito geralmente; você tentaria criar e observaria a resposta, ou usaria `turma listar` para verificar.

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

#### Exemplos Avançados e Dicas (`turma importar-alunos`):
*   **Importar para uma turma recém-criada em um único fluxo (shell):**
    Supondo que `turma criar` retorne o ID da turma de alguma forma que possa ser capturada (isso é avançado e depende do shell e da saída do programa).
    ```bash
    # Exemplo conceitual, a captura de ID pode ser complexa
    # TURMA_ID=$(./vigenda turma criar "Nova Turma X" --subjectid 4 --year 2024 | grep "ID da Turma:" | awk '{print $NF}')
    # if [ -n "$TURMA_ID" ]; then
    #   ./vigenda turma importar-alunos $TURMA_ID alunos_novos.csv
    # else
    #   echo "Falha ao criar turma ou obter ID."
    # fi
    ```
*   **Lidar com erros durante a importação:**
    O Vigenda deve reportar erros se uma linha no CSV estiver mal formatada ou se um aluno já existir com dados conflitantes. Verifique a saída do comando para tais mensagens. Alguns sistemas podem gerar um arquivo de log de erros.
*   **Reimportar uma lista de alunos (atualizar dados):**
    Verifique se o comando `importar-alunos` atualiza os alunos existentes ou apenas adiciona novos. Se ele não atualizar, você pode precisar de um comando `turma atualizar-aluno` ou gerenciar isso manualmente. O `AGENTS.md` sugere `atualizar-status`, mas não uma atualização completa de dados via CSV.

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

#### Exemplos Avançados e Dicas (`turma atualizar-status`):
*   **Atualizar status de múltiplos alunos:**
    O comando parece ser para um aluno por vez. Para múltiplos alunos, você precisaria executar o comando repetidamente, possivelmente com um script.
    ```bash
    # Exemplo de script para atualizar status de alunos listados em um arquivo
    # while IFS=, read -r alunoid novostatus; do
    #   ./vigenda turma atualizar-status --alunoid "$alunoid" --novostatus "$novostatus"
    # done < lista_atualizacao_status.txt
    # (lista_atualizacao_status.txt conteria linhas como: 101,inativo)
    ```
*   **Listar alunos para encontrar o ID ou nome correto:**
    Pode ser necessário um comando como `vigenda turma listar-alunos --classid <ID_TURMA>` para obter os IDs ou confirmar nomes antes de atualizar o status. Se não existir, a importação original ou o dashboard podem ser as fontes dessa informação.

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

#### Exemplos Avançados e Dicas (`avaliacao criar`):
*   **Criar avaliações com nomes padronizados para facilitar a organização:**
    Ex: "P1 - [Nome da Disciplina]", "T1 - [Nome da Disciplina]".
    ```bash
    ./vigenda avaliacao criar "P1 - História 9A" --classid 5 --term "1º Bimestre" --weight 10.0 --date 2024-03-15
    ```
*   **Planejar todas as avaliações do período de uma vez:**
    Se você já tem o plano de avaliação do semestre/ano, pode adicioná-las todas no início para melhor organização.
*   **Avaliações sem peso (formadoras/diagnósticas):**
    Você pode usar `--weight 0` se a avaliação não conta para a média final, mas você ainda quer registrar as notas ou participação.
    ```bash
    ./vigenda avaliacao criar "Diagnóstica - Revolução Francesa" --classid 5 --term "1º Bimestre" --weight 0 --date 2024-02-20
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

#### Exemplos Avançados e Dicas (`avaliacao lancar-notas`):
*   **Navegação no modo interativo:**
    O prompt deve indicar claramente para qual aluno a nota está sendo inserida. Pode haver opções para pular um aluno, voltar ao anterior, ou sair e salvar o progresso. (Verifique a interface do prompt).
*   **Valores de nota aceitos:**
    Confirme quais são os valores de nota válidos (ex: 0 a 10, 0 a 100, conceitos como A, B, C). O sistema deve validar a entrada.
*   **Lançar notas para um grande número de alunos:**
    O modo interativo é ideal para precisão. Para grandes turmas, certifique-se de que o processo é estável e que você pode salvar o progresso intermitentemente se for uma sessão longa (o Vigenda pode salvar automaticamente após cada entrada ou ao sair).
*   **Corrigir uma nota lançada errada:**
    Se você errar uma nota durante o lançamento interativo, verifique se pode corrigi-la antes de finalizar. Se precisar corrigir depois, pode ser necessário usar este mesmo comando `lancar-notas` novamente para a mesma avaliação (ele pode permitir sobrescrever) ou um comando específico como `avaliacao atualizar-nota`.

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

#### Detalhes do Cálculo (`avaliacao media-turma`):
O cálculo da média geralmente segue a fórmula de média ponderada:

*Média = (Nota1 * Peso1 + Nota2 * Peso2 + ... + NotaN * PesoN) / (Peso1 + Peso2 + ... + PesoN)*

Onde:
*   `NotaX` é a nota do aluno na avaliação X.
*   `PesoX` é o peso da avaliação X, definido durante a criação da avaliação com `vigenda avaliacao criar --weight <PESO>`.

**Considerações:**
*   **Período (`--term`):** Se um período é especificado (ex: "1º Bimestre"), apenas as avaliações pertencentes a esse período e à turma informada são incluídas no cálculo.
*   **Notas Ausentes:** É importante entender como o sistema trata alunos que não têm nota lançada para uma ou mais avaliações consideradas no cálculo.
    *   **Opção 1 (Comum):** O aluno pode não ter sua média calculada para o período até que todas as notas sejam lançadas.
    *   **Opção 2:** A avaliação sem nota pode ser ignorada (não recomendável pois distorce a média).
    *   **Opção 3:** Uma nota padrão (ex: 0) pode ser atribuída para avaliações sem nota lançada (o comportamento deve ser documentado ou claro para o usuário).
    Consulte a documentação específica do `AGENTS.md` ou teste para confirmar o comportamento exato do Vigenda.
*   **Média Final:** Se a opção `--term` for omitida, o sistema pode calcular uma média final considerando todas as avaliações da turma ao longo de todos os períodos, ou pode requerer um termo específico como "Final" ou "Anual".

#### Exemplos Avançados e Dicas (`avaliacao media-turma`):
*   **Calcular média para um único aluno (se suportado):**
    Verifique se há uma flag como `--alunoid <ID_ALUNO>` para focar em um aluno específico. Caso contrário, você terá que procurar o aluno na lista de saída.
*   **Exportar as médias:**
    O comando exibirá as médias no console. Para salvar em um arquivo, você pode usar o redirecionamento de saída do shell:
    ```bash
    ./vigenda avaliacao media-turma --classid 3 --term "2º Trimestre" > medias_turma3_trim2.txt
    ```
*   **Verificar o impacto de uma avaliação recém-lançada:**
    Após lançar notas para uma nova avaliação, execute `media-turma` para ver como as médias dos alunos foram afetadas.

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

#### Exemplos Avançados e Dicas (`bancoq add`):
*   **Importar múltiplos arquivos de questões:**
    Execute o comando para cada arquivo JSON que você possui.
    ```bash
    ./vigenda bancoq add questoes_hist_antiga.json
    ./vigenda bancoq add questoes_hist_media.json
    ```
*   **Organizar questões em arquivos JSON por disciplina ou tópico:**
    Isso facilita a manutenção e a importação seletiva.
*   **Verificar se questões foram duplicadas:**
    O sistema pode ou não verificar duplicatas (com base no enunciado ou em um ID interno, se houver). Se você importar o mesmo arquivo duas vezes, verifique se as questões são duplicadas ou se o sistema as ignora inteligentemente.
*   **Atualizar questões existentes:**
    O comando `bancoq add` é primariamente para adicionar. Para atualizar uma questão existente, você pode precisar de um comando `bancoq update` (que não está listado no `AGENTS.md`) ou remover a questão antiga e adicionar a nova. Verifique se o `AGENTS.md` ou `--help` oferece mais opções para gerenciamento de questões individuais no banco.

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

#### Exemplos Avançados e Dicas (`prova gerar`):
*   **Gerar uma prova com um número total específico de questões, balanceando dificuldades automaticamente:**
    ```bash
    ./vigenda prova gerar --subjectid 2 --total 10 --output prova_mat_10q.txt
    ```
    O sistema tentará selecionar uma mistura de dificuldades se os números específicos (`--easy`, `--medium`, `--hard`) não forem fornecidos. O critério de balanceamento dependerá da implementação.

*   **Gerar uma prova focada em um tópico específico:**
    ```bash
    ./vigenda prova gerar --subjectid 1 --topic "Revolução Francesa" --total 5 --output prova_rev_francesa.txt
    ```

*   **Gerar uma prova apenas com questões dissertativas (se o tipo de questão for um filtro):**
    Verifique com `./vigenda prova gerar --help` se existe uma opção como `--type dissertativa`. Se não, o tipo é definido no JSON da questão e a seleção por tipo pode não ser um filtro direto na geração.
    ```bash
    # Exemplo hipotético se o filtro existir
    # ./vigenda prova gerar --subjectid 1 --type dissertativa --total 3
    ```

*   **O que acontece se não houver questões suficientes?**
    O sistema deve informar se não puder atender ao pedido. Por exemplo, se você pedir 10 questões difíceis e só houver 3 no banco para aquela disciplina/tópico.
    A saída pode ser uma prova com menos questões do que o solicitado, junto com um aviso, ou um erro informando a insuficiência.

*   **Formato do arquivo de saída (`--output`):**
    O arquivo de saída será provavelmente um arquivo de texto simples, formatado para fácil leitura e impressão. Pode incluir:
    *   Cabeçalho (Nome da Disciplina, Prova, etc.)
    *   Numeração das questões.
    *   Enunciados.
    *   Opções (para múltipla escolha).
    *   Espaço para respostas (para dissertativas).
    *   **Importante:** Verifique se o arquivo de saída inclui ou não as respostas corretas. Geralmente, a prova para o aluno não as inclui, mas pode haver uma opção para gerar um gabarito.

*   **Gerar diferentes versões da mesma prova (se suportado):**
    Se o sistema tiver uma funcionalidade de embaralhar questões ou opções, gerar a prova múltiplas vezes com os mesmos parâmetros poderia (idealmente) produzir versões ligeiramente diferentes. Isso não é uma funcionalidade padrão em CLIs simples, a menos que explicitamente mencionado.

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
    *   **Exemplo SQLite DSN**: `file:/caminho/absoluto/para/meu_vigenda.db?cache=shared&mode=rwc&_journal_mode=WAL`
        *   `cache=shared`: Permite que múltiplas conexões (dentro do mesmo processo ou diferentes processos, dependendo do sistema operacional e como o SQLite foi compilado) acessem o mesmo banco de dados em memória.
        *   `mode=rwc`: Abre o banco de dados para leitura, escrita e o cria se não existir.
        *   `_journal_mode=WAL`: (Write-Ahead Logging) É um modo de jornalismo que pode oferecer melhor concorrência e performance em muitos casos, comparado ao `DELETE` (padrão) ou `TRUNCATE`.
    *   **Exemplo PostgreSQL DSN**: `postgres://utilizador:senha@localhost:5432/nome_da_base?sslmode=disable&connect_timeout=10`
        *   `sslmode=disable`: Para desenvolvimento local, desabilita SSL. Para produção, use `require`, `verify-ca`, ou `verify-full`.
        *   `connect_timeout=10`: Tempo em segundos para aguardar por uma conexão bem-sucedida.

#### Configuração Específica para SQLite

Se `VIGENDA_DB_TYPE` for `sqlite` (ou não estiver definida) e `VIGENDA_DB_DSN` não for fornecida, a seguinte variável é usada:

*   `VIGENDA_DB_PATH`: Caminho para o ficheiro da base de dados SQLite.
    *   **Padrão**: Um ficheiro `vigenda.db` no diretório de configuração do utilizador (ex: `~/.config/vigenda/vigenda.db` no Linux) ou no diretório atual se o diretório de configuração não for acessível.
    *   **Exemplo**: `export VIGENDA_DB_PATH="/caminho/para/sua/vigenda.db"`
    *   **Backup e Restauração (SQLite):**
        *   **Backup:** Para fazer backup de um banco de dados SQLite, basta copiar o arquivo `.db` para um local seguro enquanto a aplicação não estiver ativamente escrevendo nele. Para garantir consistência, é melhor parar a aplicação Vigenda antes de copiar o arquivo, se possível.
            ```bash
            # Exemplo de backup (com Vigenda parado ou em momento de baixa atividade)
            cp ~/.config/vigenda/vigenda.db /caminho/do/backup/vigenda_backup_$(date +%Y%m%d).db
            ```
        *   **Restauração:** Para restaurar, substitua o arquivo `.db` existente pelo arquivo de backup (com a aplicação Vigenda parada).
            ```bash
            # Exemplo de restauração (com Vigenda parado)
            cp /caminho/do/backup/vigenda_backup_YYYYMMDD.db ~/.config/vigenda/vigenda.db
            ```
        *   **Ferramentas:** O próprio `sqlite3` CLI pode ser usado para backups online:
            ```bash
            sqlite3 ~/.config/vigenda/vigenda.db ".backup /caminho/do/backup/vigenda_online_backup.db"
            ```
            Este método é geralmente seguro para ser executado mesmo com a aplicação rodando, pois lida com os locks apropriadamente.

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

**Potenciais Problemas Comuns com PostgreSQL:**
*   **Falha na Conexão:**
    *   Verifique se o servidor PostgreSQL está rodando e acessível a partir da máquina onde o Vigenda está sendo executado.
    *   Confirme se `VIGENDA_DB_HOST` e `VIGENDA_DB_PORT` estão corretos.
    *   Firewalls podem estar bloqueando a conexão.
*   **Autenticação Falhou:**
    *   Certifique-se de que `VIGENDA_DB_USER` e `VIGENDA_DB_PASSWORD` estão corretos.
    *   Verifique o arquivo `pg_hba.conf` no servidor PostgreSQL para garantir que o método de autenticação para o usuário e host de conexão está configurado corretamente (ex: `md5`, `scram-sha-256`).
*   **Base de Dados Não Existe:**
    *   A base de dados especificada em `VIGENDA_DB_NAME` deve existir no servidor PostgreSQL. Crie-a manualmente se necessário (`CREATE DATABASE nome_da_base;`).
*   **SSL/TLS:**
    *   Se o servidor PostgreSQL exigir SSL (`sslmode=require` ou superior), certifique-se de que os certificados necessários estão configurados corretamente no cliente, se aplicável (para `verify-ca` ou `verify-full`). Para conexões simples, `sslmode=prefer` pode ser uma opção, ou `disable` se o servidor não usar SSL.

### Migrações de Esquema (Schema Migrations)

*   **SQLite**: O Vigenda tentará aplicar o esquema inicial (`internal/database/migrations/001_initial_schema.sql`) automaticamente se a base de dados parecer vazia.
*   **PostgreSQL**: As migrações de esquema devem ser geridas externamente. O Vigenda não tentará criar tabelas ou modificar o esquema numa base de dados PostgreSQL existente. Certifique-se de que o esquema apropriado já foi aplicado.

### Configuração Avançada com Docker

Usar o Vigenda com Docker pode simplificar o gerenciamento de dependências e a implantação.

#### 1. Construindo uma Imagem Docker para o Vigenda

Você precisará de um `Dockerfile` na raiz do seu projeto Vigenda. Aqui está um exemplo básico:

```dockerfile
# Estágio 1: Build
FROM golang:1.23-alpine AS builder

# Instalar dependências do sistema (GCC para CGO, git para baixar módulos)
RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

# Copiar go.mod e go.sum e baixar dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o restante do código fonte
COPY . .

# Compilar a aplicação Vigenda
# O -ldflags="-s -w" reduz o tamanho do binário
# CGO_ENABLED=1 é necessário para go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /vigenda ./cmd/vigenda/

# Estágio 2: Imagem final leve
FROM alpine:latest

# Copiar o binário compilado do estágio de build
COPY --from=builder /vigenda /usr/local/bin/vigenda

# (Opcional) Copiar migrações se você quiser que o container SQLite as execute
# COPY internal/database/migrations /app/internal/database/migrations

# Definir o diretório de trabalho (opcional, mas bom para consistência)
WORKDIR /data

# Definir o entrypoint padrão
ENTRYPOINT ["vigenda"]

# (Opcional) Definir um CMD padrão se você quiser que ele execute algo ao iniciar sem argumentos
# CMD ["--help"]
```

Para construir a imagem:
```bash
docker build -t vigenda-app .
```

#### 2. Executando o Vigenda com SQLite em Docker

Para persistir os dados do SQLite, você deve montar um volume do host para o container.

*   **Se o Vigenda salva `vigenda.db` em `/data` (conforme `WORKDIR /data` no Dockerfile):**
    ```bash
    docker run -it --rm \
      -v $(pwd)/vigenda_data:/data \
      -e VIGENDA_DB_PATH="/data/vigenda.db" \
      vigenda-app tarefa listar
    ```
    Neste exemplo:
    *   `-v $(pwd)/vigenda_data:/data`: Mapeia um diretório `vigenda_data` no seu host (diretório atual) para `/data` dentro do container. O arquivo `vigenda.db` será salvo aqui.
    *   `-e VIGENDA_DB_PATH="/data/vigenda.db"`: Garante que o Vigenda dentro do container salve o banco de dados no diretório montado.

*   **Se o Vigenda salva em `~/.config/vigenda/` por padrão e você quer manter esse comportamento dentro do container (requer um usuário no container):**
    Se o Dockerfile não criar um usuário específico, ele rodará como root. Para mapear `~/.config` do container, você pode precisar ajustar o Dockerfile ou os caminhos. Uma abordagem mais simples é forçar o `VIGENDA_DB_PATH` como acima.

#### 3. Executando o Vigenda com PostgreSQL em Docker (Conectando a um Servidor PostgreSQL)

Se você tem um servidor PostgreSQL rodando (seja no host, em outro container, ou na rede), você configura o Vigenda via variáveis de ambiente para se conectar a ele.

**Exemplo com um PostgreSQL rodando em outro container chamado `my-postgres-db`:**

Primeiro, certifique-se que sua rede Docker permite a comunicação entre os containers. Geralmente, eles precisam estar na mesma rede customizada.

```bash
# Criar uma rede Docker (se ainda não tiver uma)
docker network create vigenda-net

# Rodar o container do PostgreSQL (exemplo)
docker run --name my-postgres-db --network vigenda-net \
  -e POSTGRES_USER=vigenda_user \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=vigenda_db \
  -p 5432:5432 \
  -d postgres:15

# Rodar o container do Vigenda
docker run -it --rm --network vigenda-net \
  -e VIGENDA_DB_TYPE="postgres" \
  -e VIGENDA_DB_HOST="my-postgres-db" \
  -e VIGENDA_DB_PORT="5432" \
  -e VIGENDA_DB_USER="vigenda_user" \
  -e VIGENDA_DB_PASSWORD="secret" \
  -e VIGENDA_DB_NAME="vigenda_db" \
  -e VIGENDA_DB_SSLMODE="disable" \
  vigenda-app tarefa listar
```
Notas:
*   `--network vigenda-net`: Conecta ambos os containers à mesma rede.
*   `VIGENDA_DB_HOST="my-postgres-db"`: O Vigenda usa o nome do container do PostgreSQL como host. Docker DNS resolverá isso.
*   As migrações para PostgreSQL devem ser aplicadas ao `my-postgres-db` externamente.

Esta seção fornece uma visão geral. Configurações de Docker podem se tornar complexas dependendo dos requisitos específicos de segurança, volumes e rede.

---

Para mais informações sobre como instalar e começar a usar o Vigenda rapidamente, consulte o [Guia de Introdução](../getting_started/README.md).
Se tiver perguntas comuns, visite nosso [FAQ](../faq/README.md).
Para exemplos práticos, explore nossos [Tutoriais](../tutorials/README.md).
