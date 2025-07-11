# Manual do Usuário do Vigenda

Bem-vindo ao Manual do Usuário do Vigenda! Este documento detalha todas as funcionalidades do Vigenda, uma ferramenta de linha de comando (CLI) criada para simplificar e otimizar a organização das atividades pedagógicas de professores. Seja gerenciando tarefas, turmas, avaliações ou seu banco de questões, o Vigenda oferece uma maneira focada e eficiente de manter tudo sob controle.

## Sumário

1.  [Introdução](#introducao)
    *   [Obtendo Ajuda](#obtendo-ajuda)
2.  [Conceitos Fundamentais](#conceitos-fundamentais)
    *   [Navegação e Atalhos Globais](#navegacao-e-atalhos-globais)
    *   [Entendendo os IDs](#entendendo-os-ids)
    *   [Feedback dos Comandos](#feedback-dos-comandos)
3.  [Menu Principal e Painel de Controle (Dashboard)](#menu-principal-e-painel-de-controle-dashboard)
4.  [Gestão de Tarefas](#gestao-de-tarefas)
    *   [Adicionar Tarefa (`vigenda tarefa add`)](#adicionar-tarefa-vigenda-tarefa-add)
    *   [Listar Tarefas (`vigenda tarefa listar`)](#listar-tarefas-vigenda-tarefa-listar)
    *   [Editar Tarefa (`vigenda tarefa editar`)](#editar-tarefa-vigenda-tarefa-editar)
    *   [Completar Tarefa (`vigenda tarefa complete`)](#completar-tarefa-vigenda-tarefa-complete)
    *   [Remover Tarefa (`vigenda tarefa remover`)](#remover-tarefa-vigenda-tarefa-remover)
5.  [Gestão de Turmas e Alunos](#gestao-de-turmas-e-alunos)
    *   [Criar Turma (`vigenda turma criar`)](#criar-turma-vigenda-turma-criar)
    *   [Listar Turmas (`vigenda turma listar`)](#listar-turmas-vigenda-turma-listar)
    *   [Editar Turma (`vigenda turma editar`)](#editar-turma-vigenda-turma-editar)
    *   [Remover Turma (`vigenda turma remover`)](#remover-turma-vigenda-turma-remover)
    *   [Importar Alunos (`vigenda turma importar-alunos`)](#importar-alunos-vigenda-turma-importar-alunos)
    *   [Listar Alunos (`vigenda aluno listar`)](#listar-alunos-vigenda-aluno-listar)
    *   [Atualizar Status do Aluno (`vigenda turma atualizar-status`)](#atualizar-status-do-aluno-vigenda-turma-atualizar-status)
    *   [Editar Aluno (`vigenda aluno editar`)](#editar-aluno-vigenda-aluno-editar)
    *   [Remover Aluno (`vigenda aluno remover`)](#remover-aluno-vigenda-aluno-remover)
6.  [Gestão de Avaliações e Notas](#gestao-de-avaliacoes-e-notas)
    *   [Criar Avaliação (`vigenda avaliacao criar`)](#criar-avaliacao-vigenda-avaliacao-criar)
    *   [Listar Avaliações (`vigenda avaliacao listar`)](#listar-avaliacoes-vigenda-avaliacao-listar)
    *   [Editar Avaliação (`vigenda avaliacao editar`)](#editar-avaliacao-vigenda-avaliacao-editar)
    *   [Remover Avaliação (`vigenda avaliacao remover`)](#remover-avaliacao-vigenda-avaliacao-remover)
    *   [Lançar Notas (`vigenda avaliacao lancar-notas`)](#lancar-notas-vigenda-avaliacao-lancar-notas)
    *   [Calcular Média da Turma (`vigenda avaliacao media-turma`)](#calcular-media-da-turma-vigenda-avaliacao-media-turma)
7.  [Banco de Questões e Geração de Provas](#banco-de-questoes-e-geracao-de-provas)
    *   [Adicionar Questões ao Banco (`vigenda bancoq add`)](#adicionar-questoes-ao-banco-vigenda-bancoq-add)
    *   [Listar Questões do Banco (`vigenda bancoq listar`)](#listar-questoes-do-banco-vigenda-bancoq-listar)
    *   [Editar Questão no Banco (`vigenda bancoq editar`)](#editar-questao-no-banco-vigenda-bancoq-editar)
    *   [Remover Questão do Banco (`vigenda bancoq remover`)](#remover-questao-do-banco-vigenda-bancoq-remover)
    *   [Gerar Prova (`vigenda prova gerar`)](#gerar-prova-vigenda-prova-gerar)
    *   [Gerar Gabarito (`vigenda prova gabarito`)](#gerar-gabarito-vigenda-prova-gabarito)
8.  [Formatos de Ficheiros de Importação](#formatos-de-ficheiros-de-importacao)
    *   [Importação de Alunos (CSV)](#importacao-de-alunos-csv)
    *   [Importação de Questões (JSON)](#importacao-de-questoes-json)
9.  [Configuração da Base de Dados](#configuracao-da-base-de-dados)
    *   [Tipos de Base de Dados Suportados](#tipos-de-base-de-dados-suportados)
    *   [Variáveis de Ambiente para Configuração](#variaveis-de-ambiente-para-configuracao)
    *   [Configuração Específica para SQLite](#configuracao-especifica-para-sqlite)
    *   [Configuração Específica para PostgreSQL](#configuracao-especifica-para-postgresql)
    *   [Exemplos de Configuração](#exemplos-de-configuracao)
    *   [Migrações de Esquema](#migracoes-de-esquema-schema-migrations)
10. [Dicas de Uso e Boas Práticas](#dicas-de-uso-e-boas-praticas)
11. [Solução de Problemas Comuns (FAQ)](#solucao-de-problemas-comuns-faq)

## 1. Introdução

O Vigenda é uma aplicação de linha de comando (CLI) projetada para ajudar professores a gerenciar tarefas, aulas, avaliações e outras atividades relacionadas ao ensino. Ele visa oferecer uma maneira focada e eficiente de organização, especialmente útil para aqueles que podem se beneficiar de uma interface direta e sem distrações. Com o Vigenda, você pode simplificar seu fluxo de trabalho pedagógico e ter mais tempo para o que realmente importa: ensinar.

### Obtendo Ajuda

A qualquer momento, você pode obter ajuda sobre os comandos do Vigenda diretamente na linha de comando:

*   Para uma visão geral de todos os comandos disponíveis e opções globais:
    ```bash
    ./vigenda --help
    ```
*   Para ajuda específica sobre um comando (por exemplo, `tarefa add`):
    ```bash
    ./vigenda tarefa add --help
    ```
    Ou para um submódulo (por exemplo, `tarefa`):
    ```bash
    ./vigenda tarefa --help
    ```

## 2. Conceitos Fundamentais

Antes de mergulhar nos comandos específicos, alguns conceitos são importantes para entender como o Vigenda funciona.

### Navegação e Atalhos Globais

*   **Menu Principal:** Ao executar `./vigenda` sem argumentos, você acessa o menu principal interativo. Use as setas para cima/baixo e `Enter` para selecionar.
*   **Voltar ao Menu:** De qualquer módulo ou painel, pressione `Esc` para retornar ao Menu Principal.
*   **Recarregar Painel de Controle:** Dentro do Painel de Controle, pressione `r` para atualizar os dados exibidos.

### Entendendo os IDs

Muitas operações no Vigenda (como editar, remover, ou associar itens) dependem de **IDs** numéricos únicos. Por exemplo, cada turma, aluno, tarefa ou avaliação terá seu próprio ID.

*   **Como obter IDs?**
    *   Geralmente, quando você cria um novo item (ex: uma turma com `vigenda turma criar`), o Vigenda exibirá uma mensagem de sucesso contendo o ID do item recém-criado. **É uma boa prática anotar esses IDs.**
    *   Comandos de listagem (ex: `vigenda turma listar`, `vigenda tarefa listar`) também exibirão os IDs dos itens.
*   **Por que IDs?** Usar IDs garante que você está referenciando o item exato que deseja modificar ou usar, evitando ambiguidades, especialmente se você tiver itens com nomes parecidos.

### Feedback dos Comandos

*   **Mensagens de Sucesso:** Após a execução bem-sucedida de um comando que cria ou modifica dados, o Vigenda geralmente exibirá uma mensagem de confirmação, muitas vezes incluindo o ID do item afetado. Por exemplo: `Tarefa "Preparar aula" (ID: 15) adicionada com sucesso.`
*   **Mensagens de Erro:** Se um comando não puder ser executado, o Vigenda fornecerá uma mensagem de erro. Leia atentamente para entender o problema. Erros comuns incluem:
    *   Falta de argumentos obrigatórios (ex: não fornecer a descrição de uma tarefa).
    *   ID não encontrado (ex: tentar editar uma tarefa com um ID que não existe).
    *   Formato de dados inválido (ex: uma data em formato incorreto).
    *   Problemas de permissão ou conexão com o banco de dados.
    Se a mensagem de erro não for clara, tente o comando `--help` para revisar o uso correto.

## 3. Menu Principal e Painel de Controle (Dashboard)

Ao executar o Vigenda sem subcomandos, você acessa o **Menu Principal** interativo:

```bash
./vigenda
```

O Menu Principal lista todas as principais funcionalidades do Vigenda, permitindo que você navegue para diferentes módulos da aplicação.

Uma das opções neste menu é o **"Painel de Controle"**. Ao selecioná-lo, você acessa o dashboard da aplicação, que fornece uma visão geral da sua agenda do dia, tarefas urgentes e outras notificações importantes.

O Painel de Controle é uma tela dinâmica que exibe informações relevantes para o seu dia a dia. Atualmente, ele foca em:
*   **Tarefas Próximas:** Uma lista de tarefas com prazos futuros, ajudando você a priorizar.
*   **Outras Seções (Potenciais):** Funcionalidades como "Aulas de Hoje", "Próximas Avaliações", "Notificações" e "Resumo de Progresso" estão planejadas para futuras versões.

A aparência e o conteúdo exato do Painel de Controle podem evoluir. Consulte sempre a versão mais recente da documentação ou use o comando `--help` para novidades.

(As dicas de uso foram movidas para a seção [Navegação e Atalhos Globais](#navegacao-e-atalhos-globais)).

## 4. Gestão de Tarefas

Gerencie suas tarefas pedagógicas de forma simples e eficaz.

> **Nota sobre IDs:** Para editar, completar ou remover tarefas, você precisará do ID da tarefa. Este ID é fornecido quando a tarefa é criada com `vigenda tarefa add` ou pode ser obtido com `vigenda tarefa listar`.

### Adicionar Tarefa (`vigenda tarefa add`)

Crie novas tarefas com descrições, prazos e associação a turmas.

**Uso:**
```bash
./vigenda tarefa add "Descrição da Tarefa" --classid <ID_DA_TURMA> --duedate <AAAA-MM-DD> [outras opções]
```

**Argumentos e Opções:**
*   `"Descrição da Tarefa"`: (Obrigatório) O texto que descreve a tarefa.
*   `--classid <ID_DA_TURMA>`: (Obrigatório) O ID numérico da turma à qual esta tarefa está associada. Use `vigenda turma listar` para encontrar o ID da turma.
*   `--duedate <AAAA-MM-DD>`: (Obrigatório) A data de vencimento da tarefa no formato ano-mês-dia.
*   `--priority <prioridade>`: (Opcional) Nível de prioridade. Valores permitidos: `baixa`, `media`, `alta`. Padrão: `media`.
*   `--notes "Notas adicionais"`: (Opcional) Observações extras sobre a tarefa.

**Exemplo:**
```bash
./vigenda tarefa add "Preparar apresentação sobre a Revolução Industrial" --classid 1 --duedate 2024-08-15 --priority alta
```
**Feedback Esperado:**
```
Tarefa "Preparar apresentação sobre a Revolução Industrial" (ID: 23) adicionada com sucesso para a turma ID 1.
```

#### Exemplos Avançados e Dicas (`tarefa add`):
*   **Adicionar uma tarefa com notas detalhadas:**
    ```bash
    ./vigenda tarefa add "Finalizar correção das provas" --classid 2 --duedate 2024-07-25 --priority alta --notes "Lembrar de verificar questões dissertativas com critério X. Publicar notas até dia 26."
    ```
*   **Adicionar uma tarefa recorrente (Dica):**
    O Vigenda não possui suporte nativo a tarefas recorrentes. Para tarefas que se repetem, como "Preparar aula de segunda-feira", adicione a primeira ocorrência. Ao completá-la, adicione a próxima manualmente ou use um script de shell para automatizar se necessário.
    ```bash
    ./vigenda tarefa add "Preparar aula de Segunda (Semana 1)" --classid 1 --duedate 2024-08-05
    ```

### Listar Tarefas (`vigenda tarefa listar`)

Visualize tarefas, com a opção de filtrar por turma ou status. Este comando é essencial para obter os IDs das tarefas.

**Uso:**
```bash
./vigenda tarefa listar [--classid <ID_DA_TURMA>] [--status <status>]
```

**Opções:**
*   `--classid <ID_DA_TURMA>`: (Opcional) Filtra tarefas pelo ID da turma especificada.
*   `--status <status>`: (Opcional) Filtra tarefas pelo status. Valores permitidos: `pendente`, `concluida`, `atrasada`. Se omitido, lista as tarefas `pendente`.
*   `--all`: (Opcional) Lista todas as tarefas, independentemente do status (combina `pendente`, `concluida`, `atrasada`).
*   `--priority <prioridade>`: (Opcional) Filtra tarefas pelo nível de prioridade. Valores permitidos: `baixa`, `media`, `alta`.
*   `--days <numero>`: (Opcional) Filtra tarefas cujo prazo (`duedate`) está nos próximos `<numero>` dias (ex: `--days 7` para tarefas da próxima semana). *(Nota: Funcionalidade avançada, pode não estar em todas as versões)*.

**Exemplo de Saída:**
```
ID | Descrição                               | Turma ID | Prazo      | Prioridade | Status
---|-----------------------------------------|----------|------------|------------|----------
23 | Preparar apresentação Rev. Industrial   | 1        | 2024-08-15 | alta       | pendente
24 | Corrigir provas de História             | 1        | 2024-08-10 | media      | pendente
28 | Lançar notas da P1 de Matemática        | 2        | 2024-08-12 | alta       | pendente
```

**Exemplos de Uso:**
*   **Listar todas as tarefas pendentes para a turma com ID 1:**
    ```bash
    ./vigenda tarefa listar --classid 1 --status pendente
    ```
*   **Listar todas as tarefas concluídas de todas as turmas:**
    ```bash
    ./vigenda tarefa listar --status concluida
    ```
*   **Listar todas as tarefas da turma 2, independentemente do status:**
    ```bash
    ./vigenda tarefa listar --classid 2 --all
    ```
*   **Listar todas as tarefas de alta prioridade:**
    ```bash
    ./vigenda tarefa listar --priority alta
    ```

### Editar Tarefa (`vigenda tarefa editar`)

Modifique os detalhes de uma tarefa existente.

**Uso:**
```bash
./vigenda tarefa editar <ID_DA_TAREFA> [--desc "Nova Descrição"] [--classid <NOVO_ID_TURMA>] [--duedate <NOVA_DATA>] [--priority <NOVA_PRIORIDADE>] [--notes "Novas Notas"] [--status <NOVO_STATUS>]
```

**Argumentos e Opções:**
*   `<ID_DA_TAREFA>`: (Obrigatório) O ID da tarefa a ser editada.
*   `--desc "Nova Descrição"`: (Opcional) Novo texto de descrição da tarefa.
*   `--classid <NOVO_ID_TURMA>`: (Opcional) Novo ID da turma associada.
*   `--duedate <NOVA_DATA>`: (Opcional) Nova data de vencimento (AAAA-MM-DD).
*   `--priority <NOVA_PRIORIDADE>`: (Opcional) Novo nível de prioridade (`baixa`, `media`, `alta`).
*   `--notes "Novas Notas"`: (Opcional) Novas notas adicionais.
*   `--status <NOVO_STATUS>`: (Opcional) Novo status da tarefa (`pendente`, `concluida`, `atrasada`). Usar `vigenda tarefa complete` é geralmente preferível para marcar como concluída.

**Exemplo:**
```bash
./vigenda tarefa editar 23 --desc "Preparar e revisar apresentação sobre a Revolução Industrial" --priority media
```
**Feedback Esperado:**
```
Tarefa ID 23 atualizada com sucesso.
```
*(Nota: Este comando é uma adição sugerida para maior flexibilidade. Verifique `./vigenda tarefa editar --help` para confirmar sua disponibilidade e opções exatas na sua versão.)*

### Completar Tarefa (`vigenda tarefa complete`)

Marque uma ou mais tarefas como concluídas. Isso geralmente muda o status da tarefa para `concluida`.

**Uso:**
```bash
./vigenda tarefa complete <ID_DA_TAREFA_1> [ID_DA_TAREFA_2 ...]
```

**Argumentos:**
*   `<ID_DA_TAREFA>`: (Obrigatório) O ID numérico da tarefa a marcar como concluída. Pode fornecer múltiplos IDs.

**Exemplo:**
```bash
./vigenda tarefa complete 23
```
**Feedback Esperado:**
```
Tarefa ID 23 marcada como concluída.
```

```bash
./vigenda tarefa complete 10 11 15
```
**Feedback Esperado:**
```
Tarefa ID 10 marcada como concluída.
Tarefa ID 11 marcada como concluída.
Tarefa ID 15 marcada como concluída.
```

#### Dicas (`tarefa complete`):
*   **Completar e adicionar tarefa de acompanhamento (Shell):**
    ```bash
    ./vigenda tarefa complete 25 && ./vigenda tarefa add "Revisar feedback da tarefa 25" --classid 1 --duedate 2024-08-01
    ```
*   **Reabrir uma tarefa:** Para marcar uma tarefa como `pendente` novamente após ter sido completada, use `vigenda tarefa editar <ID_DA_TAREFA> --status pendente`.

### Remover Tarefa (`vigenda tarefa remover`)

Exclui permanentemente uma tarefa do sistema. Use com cuidado.

**Uso:**
```bash
./vigenda tarefa remover <ID_DA_TAREFA_1> [ID_DA_TAREFA_2 ...] [--force]
```

**Argumentos e Opções:**
*   `<ID_DA_TAREFA>`: (Obrigatório) O ID da tarefa a ser removida. Pode fornecer múltiplos IDs.
*   `--force`: (Opcional) Remove a tarefa sem pedir confirmação. Por padrão, o sistema pode pedir para confirmar a exclusão.

**Exemplo:**
```bash
./vigenda tarefa remover 45
```
**Feedback Esperado (sem --force):**
```
Você tem certeza que quer remover a tarefa ID 45 ("Descrição da Tarefa")? [s/N]: s
Tarefa ID 45 removida com sucesso.
```

```bash
./vigenda tarefa remover 46 --force
```
**Feedback Esperado (com --force):**
```
Tarefa ID 46 removida com sucesso.
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda tarefa remover --help` para confirmar sua disponibilidade.)*


## 5. Gestão de Turmas e Alunos

Organize suas turmas e os alunos pertencentes a cada uma.

> **Nota sobre IDs:** Os IDs de Turma são essenciais para associar alunos, tarefas e avaliações. Os IDs de Aluno são usados para operações específicas de alunos como atualizar status ou editar dados. Obtenha-os na criação ou via comandos de listagem.

### Criar Turma (`vigenda turma criar`)

Crie novas turmas, que servem como contêineres para alunos, tarefas e avaliações.

**Uso:**
```bash
./vigenda turma criar "Nome da Turma" --subjectid <ID_DA_DISCIPLINA> [--year <ANO>] [--period <PERIODO>]
```

**Argumentos e Opções:**
*   `"Nome da Turma"`: (Obrigatório) O nome da turma (ex: "Matemática 9A", "História - Tarde").
*   `--subjectid <ID_DA_DISCIPLINA>`: (Obrigatório) O ID numérico da disciplina à qual esta turma pertence. *(Nota: O Vigenda pode futuramente incluir um módulo de gestão de disciplinas. Por enquanto, assume-se que IDs de disciplina são gerenciados externamente ou implicitamente).*
*   `--year <ANO>`: (Opcional) O ano letivo da turma (ex: 2024).
*   `--period <PERIODO>`: (Opcional) O período/semestre da turma (ex: "1º Semestre", "Anual").
*   `--notes <NOTAS>`: (Opcional) Observações ou notas adicionais sobre a turma.

**Exemplo:**
```bash
./vigenda turma criar "Física 2B" --subjectid 3 --year 2024 --period "2º Semestre" --notes "Turma avançada"
```
**Feedback Esperado:**
```
Turma "Física 2B" (ID: 7) criada com sucesso.
```

#### Dicas (`turma criar`):
*   **Criar múltiplas turmas (Shell Scripting):**
    Para criar várias turmas com um padrão, você pode usar um loop de shell:
    ```bash
    # Exemplo para shell bash/zsh
    for LETRA in A B C; do \
      ./vigenda turma criar "Matemática 1${LETRA}" --subjectid 2 --year 2024; \
    done
    ```
*   **Verificar existência:** Antes de criar, use `vigenda turma listar` para ver as turmas já existentes e evitar duplicatas.

### Listar Turmas (`vigenda turma listar`)

Visualize todas as turmas cadastradas, exibindo seus IDs e detalhes.

**Uso:**
```bash
./vigenda turma listar [--year <ANO>] [--period <PERIODO>] [--subjectid <ID_DISCIPLINA>]
```

**Opções:**
*   `--year <ANO>`: (Opcional) Filtra turmas pelo ano letivo.
*   `--period <PERIODO>`: (Opcional) Filtra turmas pelo período/semestre.
*   `--subjectid <ID_DISCIPLINA>`: (Opcional) Filtra turmas pelo ID da disciplina.

**Exemplo de Saída:**
```
ID | Nome da Turma | Disciplina ID | Ano  | Período      | Notas
---|---------------|---------------|------|--------------|----------------
1  | Matemática 9A | 2             | 2024 | Anual        |
7  | Física 2B     | 3             | 2024 | 2º Semestre  | Turma avançada
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda turma listar --help` para confirmar sua disponibilidade e formato de saída.)*

### Editar Turma (`vigenda turma editar`)

Modifique os detalhes de uma turma existente.

**Uso:**
```bash
./vigenda turma editar <ID_DA_TURMA> [--name "Novo Nome"] [--subjectid <NOVO_ID_DISCIPLINA>] [--year <NOVO_ANO>] [--period <NOVO_PERIODO>] [--notes <NOVAS_NOTAS>]
```

**Argumentos e Opções:**
*   `<ID_DA_TURMA>`: (Obrigatório) O ID da turma a ser editada.
*   `--name "Novo Nome"`: (Opcional) Novo nome para a turma.
*   `--subjectid <NOVO_ID_DISCIPLINA>`: (Opcional) Novo ID de disciplina.
*   `--year <NOVO_ANO>`: (Opcional) Novo ano letivo.
*   `--period <NOVO_PERIODO>`: (Opcional) Novo período/semestre.
*   `--notes <NOVAS_NOTAS>`: (Opcional) Novas observações para a turma. Para remover notas existentes, passe aspas vazias: `--notes ""`.

**Exemplo:**
```bash
./vigenda turma editar 7 --notes "Turma avançada com foco em laboratório"
```
**Feedback Esperado:**
```
Turma ID 7 atualizada com sucesso.
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda turma editar --help` para confirmar sua disponibilidade.)*

### Remover Turma (`vigenda turma remover`)

Exclui permanentemente uma turma e **todos os seus dados associados** (alunos, tarefas, avaliações). Esta ação é irreversível. Use com extrema cautela.

**Uso:**
```bash
./vigenda turma remover <ID_DA_TURMA_1> [ID_DA_TURMA_2 ...] --force
```

**Argumentos e Opções:**
*   `<ID_DA_TURMA>`: (Obrigatório) O ID da turma a ser removida. Pode fornecer múltiplos IDs.
*   `--force`: (Obrigatório para este comando destrutivo) Confirma que você entende as consequências da remoção. Sem `--force`, o comando não será executado.

**Exemplo:**
```bash
./vigenda turma remover 7 --force
```
**Feedback Esperado:**
```
Turma ID 7 e todos os seus dados associados foram removidos com sucesso.
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda turma remover --help` para confirmar sua disponibilidade e o requisito do `--force`.)*


### Importar Alunos (`vigenda turma importar-alunos`)

Importe uma lista de alunos para uma turma específica a partir de um ficheiro CSV. Se um aluno com o mesmo `nome_completo` (ou `numero_chamada`, se usado como chave) já existir na turma, seus dados podem ser atualizados ou a importação daquele aluno pode ser ignorada, dependendo da implementação.

**Uso:**
```bash
./vigenda turma importar-alunos <ID_DA_TURMA> <CAMINHO_DO_ARQUIVO_CSV> [--update-existing]
```

**Argumentos e Opções:**
*   `<ID_DA_TURMA>`: (Obrigatório) O ID numérico da turma para a qual os alunos serão importados.
*   `<CAMINHO_DO_ARQUIVO_CSV>`: (Obrigatório) O caminho para o ficheiro CSV contendo os dados dos alunos.
*   `--update-existing`: (Opcional) Se presente, atualiza os dados de alunos já existentes na turma que correspondem aos do CSV. Caso contrário, alunos duplicados podem ser ignorados ou causar um erro. (Comportamento exato a ser verificado com `--help`).

**Formato do CSV:** Consulte a seção [Importação de Alunos (CSV)](#importacao-de-alunos-csv) para detalhes.

**Exemplo:**
```bash
./vigenda turma importar-alunos 2 alunos_turma_2B.csv --update-existing
```
**Feedback Esperado:**
```
Importação de alunos para a turma ID 2 concluída.
Alunos novos adicionados: 15
Alunos existentes atualizados: 3
Linhas com erro: 1 (verifique o arquivo alunos_turma_2B.csv.errors para detalhes)
```
(O feedback pode variar, incluindo a geração de um arquivo de log para erros).

#### Dicas (`turma importar-alunos`):
*   **Lidar com Erros:** Verifique a saída do comando para mensagens de erro. Se um arquivo de log de erros for gerado, examine-o para corrigir problemas no CSV.
*   **Reimportação:** Para atualizar uma lista grande de alunos, o uso de `--update-existing` (se disponível) é recomendado.

### Listar Alunos (`vigenda aluno listar`)

Visualize os alunos de uma turma específica, exibindo seus IDs e detalhes.

**Uso:**
```bash
./vigenda aluno listar --classid <ID_DA_TURMA>
```

**Opções:**
*   `--classid <ID_DA_TURMA>`: (Obrigatório) O ID da turma cujos alunos você deseja listar.

**Exemplo de Saída:**
```
ID Aluno | Nome Completo        | Nº Chamada | Situação
---------|----------------------|------------|----------
101      | Alice Wonderland     | 1          | ativo
102      | Bob The Builder      | 2          | ativo
103      | Charles Xavier       | 3          | inativo
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda aluno listar --help` para confirmar sua disponibilidade e formato de saída.)*

### Atualizar Status do Aluno (`vigenda turma atualizar-status`)

Modifique o status de um aluno dentro de uma turma (ex: de `ativo` para `inativo` ou `transferido`).

**Uso:**
```bash
./vigenda turma atualizar-status --alunoid <ID_DO_ALUNO> --novostatus <NOVO_STATUS>
```
ou
```bash
./vigenda turma atualizar-status --nomealuno "Nome Completo do Aluno" --turmaid <ID_DA_TURMA> --novostatus <NOVO_STATUS>
```

**Opções:**
*   `--alunoid <ID_DO_ALUNO>`: (Obrigatório, a menos que `--nomealuno` e `--turmaid` sejam usados) O ID numérico do aluno cujo status será atualizado. Obtenha com `vigenda aluno listar --classid <ID>`.
*   `--nomealuno "Nome Completo do Aluno"`: (Opcional) O nome completo do aluno. Usado em conjunto com `--turmaid` se o ID do aluno não for conhecido. Se houver alunos com nomes idênticos na mesma turma, o comando pode falhar ou afetar o primeiro encontrado; usar `--alunoid` é mais seguro.
*   `--turmaid <ID_DA_TURMA>`: (Obrigatório se estiver usando `--nomealuno`) O ID da turma do aluno.
*   `--novostatus <NOVO_STATUS>`: (Obrigatório) O novo status para o aluno. Valores permitidos: `ativo`, `inativo`, `transferido`, `graduado`.

**Exemplo:**
```bash
./vigenda turma atualizar-status --alunoid 101 --novostatus transferido
```
**Feedback Esperado:**
```
Status do aluno ID 101 atualizado para "transferido".
```

```bash
./vigenda turma atualizar-status --nomealuno "Maria Silva" --turmaid 2 --novostatus inativo
```
**Feedback Esperado:**
```
Status do aluno "Maria Silva" (ID: 108) na turma ID 2 atualizado para "inativo".
```
(O sistema idealmente retornaria o ID do aluno afetado se encontrado pelo nome).

#### Dicas (`turma atualizar-status`):
*   **Precisão:** Para evitar ambiguidades, prefira usar `--alunoid`.
*   **Múltiplos Alunos:** Para atualizar status de vários alunos, você pode precisar executar o comando repetidamente ou usar um script de shell.

### Editar Aluno (`vigenda aluno editar`)

Modifique os detalhes de um aluno existente, como nome ou número de chamada.

**Uso:**
```bash
./vigenda aluno editar <ID_DO_ALUNO> [--nome "Novo Nome Completo"] [--numchamada <NOVO_NUM_CHAMADA>]
```

**Argumentos e Opções:**
*   `<ID_DO_ALUNO>`: (Obrigatório) O ID do aluno a ser editado.
*   `--nome "Novo Nome Completo"`: (Opcional) Novo nome completo do aluno.
*   `--numchamada <NOVO_NUM_CHAMADA>`: (Opcional) Novo número de chamada do aluno.

**Exemplo:**
```bash
./vigenda aluno editar 101 --nome "Alice Wonderland Jr."
```
**Feedback Esperado:**
```
Dados do aluno ID 101 atualizados com sucesso.
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda aluno editar --help` para confirmar sua disponibilidade.)*

### Remover Aluno (`vigenda aluno remover`)

Exclui permanentemente um aluno de uma turma. Isso pode também remover notas associadas a este aluno. Use com cuidado.

**Uso:**
```bash
./vigenda aluno remover <ID_DO_ALUNO> [--turmaid <ID_DA_TURMA>] [--force]
```

**Argumentos e Opções:**
*   `<ID_DO_ALUNO>`: (Obrigatório) O ID do aluno a ser removido.
*   `--turmaid <ID_DA_TURMA>`: (Opcional, mas recomendado) Especifica a turma da qual o aluno será removido, especialmente se o aluno pudesse estar em múltiplas turmas (embora o modelo atual pareça ser um aluno por turma via importação). Se não fornecido, pode remover o aluno de todas as turmas ou requerer que o aluno esteja em apenas uma.
*   `--force`: (Opcional) Remove o aluno sem pedir confirmação.

**Exemplo:**
```bash
./vigenda aluno remover 101 --turmaid 1 --force
```
**Feedback Esperado:**
```
Aluno ID 101 removido da turma ID 1 com sucesso.
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda aluno remover --help` para confirmar sua disponibilidade e o impacto nos dados associados.)*


## 6. Gestão de Avaliações e Notas

Defina avaliações, lance notas e calcule médias.

> **Nota sobre IDs:** IDs de Avaliação são cruciais para lançar notas ou editar/remover avaliações. Obtenha-os na criação ou com `vigenda avaliacao listar`.

### Criar Avaliação (`vigenda avaliacao criar`)

Defina novas avaliações para as turmas, especificando nome, período, peso e data.

**Uso:**
```bash
./vigenda avaliacao criar "Nome da Avaliação" --classid <ID_DA_TURMA> --term <PERIODO_AVALIATIVO> --weight <PESO> [--date <AAAA-MM-DD>]
```

**Argumentos e Opções:**
*   `"Nome da Avaliação"`: (Obrigatório) O nome da avaliação (ex: "Prova Mensal - Unidade 1", "Trabalho em Grupo").
*   `--classid <ID_DA_TURMA>`: (Obrigatório) O ID numérico da turma para a qual esta avaliação se aplica.
*   `--term <PERIODO_AVALIATIVO>`: (Obrigatório) O período avaliativo (ex: "1º Bimestre", "Trimestre Final"). Use termos consistentes.
*   `--weight <PESO>`: (Obrigatório) O peso desta avaliação no cálculo da média final (ex: 1.0, 2.5, 0.0 para avaliações formativas). Deve ser um número não negativo.
*   `--date <AAAA-MM-DD>`: (Opcional) A data em que a avaliação será aplicada ou entregue. Formato: ano-mês-dia.

**Exemplo:**
```bash
./vigenda avaliacao criar "Seminário sobre Ecossistemas" --classid 3 --term "2º Trimestre" --weight 2.0 --date 2024-09-10
```
**Feedback Esperado:**
```
Avaliação "Seminário sobre Ecossistemas" (ID: 12) criada com sucesso para a turma ID 3.
```

#### Dicas (`avaliacao criar`):
*   **Nomes Padronizados:** Use nomes como "P1 - [Nome da Disciplina]", "T1 - [Nome da Disciplina]" para fácil identificação.
*   **Planejamento Antecipado:** Adicione todas as avaliações do período no início para melhor organização.
*   **Avaliações Formativas:** Use `--weight 0` para avaliações que não contam para a média final, mas cujas notas você deseja registrar.

### Listar Avaliações (`vigenda avaliacao listar`)

Visualize as avaliações de uma turma, opcionalmente filtrando por período.

**Uso:**
```bash
./vigenda avaliacao listar --classid <ID_DA_TURMA> [--term <PERIODO_AVALIATIVO>]
```

**Opções:**
*   `--classid <ID_DA_TURMA>`: (Obrigatório) O ID da turma cujas avaliações você deseja listar.
*   `--term <PERIODO_AVALIATIVO>`: (Opcional) Filtra avaliações pelo período especificado.

**Exemplo de Saída:**
```
ID | Nome da Avaliação             | Turma ID | Período      | Peso | Data
---|-------------------------------|----------|--------------|------|-----------
10 | P1 - História 9A              | 1        | 1º Bimestre  | 5.0  | 2024-03-15
11 | T1 - Trabalho: Egito Antigo   | 1        | 1º Bimestre  | 3.0  | 2024-03-29
12 | Participacao - 1º Bimestre    | 1        | 1º Bimestre  | 2.0  |
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda avaliacao listar --help` para confirmar sua disponibilidade.)*

### Editar Avaliação (`vigenda avaliacao editar`)

Modifique os detalhes de uma avaliação existente.

**Uso:**
```bash
./vigenda avaliacao editar <ID_DA_AVALIACAO> [--name "Novo Nome"] [--term <NOVO_PERIODO>] [--weight <NOVO_PESO>] [--date <NOVA_DATA>]
```

**Argumentos e Opções:**
*   `<ID_DA_AVALIACAO>`: (Obrigatório) O ID da avaliação a ser editada.
*   `--name "Novo Nome"`: (Opcional) Novo nome para a avaliação.
*   `--term <NOVO_PERIODO>`: (Opcional) Novo período avaliativo.
*   `--weight <NOVO_PESO>`: (Opcional) Novo peso para a avaliação.
*   `--date <NOVA_DATA>`: (Opcional) Nova data para a avaliação (AAAA-MM-DD).

**Exemplo:**
```bash
./vigenda avaliacao editar 12 --weight 2.5 --date 2024-04-05
```
**Feedback Esperado:**
```
Avaliação ID 12 atualizada com sucesso.
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda avaliacao editar --help` para confirmar sua disponibilidade.)*

### Remover Avaliação (`vigenda avaliacao remover`)

Exclui permanentemente uma avaliação e todas as notas associadas a ela. Use com cuidado.

**Uso:**
```bash
./vigenda avaliacao remover <ID_DA_AVALIACAO_1> [ID_DA_AVALIACAO_2 ...] --force
```

**Argumentos e Opções:**
*   `<ID_DA_AVALIACAO>`: (Obrigatório) O ID da avaliação a ser removida. Pode fornecer múltiplos IDs.
*   `--force`: (Opcional, mas recomendado para evitar remoção acidental) Remove a avaliação sem pedir confirmação.

**Exemplo:**
```bash
./vigenda avaliacao remover 10 --force
```
**Feedback Esperado:**
```
Avaliação ID 10 e suas notas associadas foram removidas com sucesso.
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda avaliacao remover --help` para confirmar sua disponibilidade.)*


### Lançar Notas (`vigenda avaliacao lancar-notas`)

Lance as notas dos alunos de forma interativa para uma avaliação específica. O sistema listará os alunos da turma associada à avaliação e solicitará a nota para cada um.

**Uso:**
```bash
./vigenda avaliacao lancar-notas <ID_DA_AVALIACAO>
```

**Argumentos:**
*   `<ID_DA_AVALIACAO>`: (Obrigatório) O ID numérico da avaliação para a qual as notas serão lançadas.

**Exemplo:**
```bash
./vigenda avaliacao lancar-notas 12
```
Ao executar, o sistema iniciará um prompt interativo:
```
Lançando notas para a Avaliação: "Participacao - 1º Bimestre" (ID: 12)
Turma: História 9A (ID: 1)
--------------------------------------------------
Aluno: Alice Wonderland (ID: 101)
Nota (deixe em branco para pular, 's' para sair e salvar, 'c' para cancelar): 8.5
--------------------------------------------------
Aluno: Bob The Builder (ID: 102)
Nota (deixe em branco para pular, 's' para sair e salvar, 'c' para cancelar): 7.0
...
Notas lançadas com sucesso para a avaliação ID 12.
```

#### Dicas (`avaliacao lancar-notas`):
*   **Navegação Interativa:** O prompt deve indicar claramente para qual aluno a nota está sendo inserida. Siga as instruções na tela para inserir notas, pular alunos (deixando em branco e pressionando Enter), salvar e sair (`s`), ou cancelar a operação (`c`).
*   **Valores de Nota:** O sistema validará os valores de nota inseridos (ex: geralmente números, possivelmente dentro de um intervalo como 0-10 ou 0-100). Se um valor inválido for inserido, uma mensagem de erro será mostrada, e você poderá tentar novamente para o mesmo aluno.
*   **Sessões Longas:** Para turmas grandes, se a sessão for interrompida, você geralmente pode retomar o lançamento executando o comando novamente. O sistema pode pular alunos que já têm notas ou permitir a sobrescrita.
*   **Correção de Notas:**
    *   **Durante o lançamento:** Se você errar uma nota e perceber antes de ir para o próximo aluno, algumas interfaces podem permitir voltar (verifique as instruções do prompt).
    *   **Após o lançamento:** Para corrigir uma nota já lançada, execute `vigenda avaliacao lancar-notas <ID_DA_AVALIACAO>` novamente. O sistema geralmente permite sobrescrever a nota existente para alunos específicos quando você chegar neles no prompt. Alternativamente, um comando `vigenda avaliacao editar-nota --alunoid <ID> --avaliacaoid <ID> --nota <NOVA_NOTA>` seria ideal para edições pontuais (verifique se existe).

### Calcular Média da Turma (`vigenda avaliacao media-turma`)

Calcule a média geral de uma turma com base nas avaliações e seus pesos definidos para um determinado período avaliativo, ou a média final.

**Uso:**
```bash
./vigenda avaliacao media-turma --classid <ID_DA_TURMA> [--term <PERIODO_AVALIATIVO>]
```

**Opções:**
*   `--classid <ID_DA_TURMA>`: (Obrigatório) O ID numérico da turma para a qual a média será calculada.
*   `--term <PERIODO_AVALIATIVO>`: (Opcional) Especifica o período avaliativo para o cálculo. Se omitido, o Vigenda pode calcular uma média geral de todos os períodos ou requerer um termo específico (ex: "Final"). Verifique com `--help`.
*   `--alunoid <ID_ALUNO>`: (Opcional) Calcula a média para um aluno específico dentro da turma e período.

**Exemplo:**
```bash
./vigenda avaliacao media-turma --classid 1 --term "1º Bimestre"
```
**Exemplo de Saída:**
```
Médias para Turma ID 1 ("História 9A"), Período: 1º Bimestre
-----------------------------------------------------------
Aluno ID | Nome Aluno        | Média Bimestral
-----------------------------------------------------------
101      | Alice Wonderland  | 8.75
102      | Bob The Builder   | 7.20
...
-----------------------------------------------------------
* Cálculo baseado nas avaliações com notas lançadas e pesos definidos para o período.
```

#### Detalhes do Cálculo (`avaliacao media-turma`):
A média é calculada usando a fórmula de média ponderada:
*Média = (Nota1 * Peso1 + Nota2 * Peso2 + ... + NotaN * PesoN) / (Peso1 + Peso2 + ... + PesoN)*

**Considerações Importantes:**
*   **Período (`--term`):** Apenas avaliações do período especificado (e da turma) são incluídas.
*   **Notas Ausentes:** (Comportamento a ser confirmado/implementado) Idealmente, o Vigenda deve:
    1.  **Opção Padrão (Recomendada):** Não calcular a média para um aluno se ele tiver notas pendentes em avaliações com peso maior que zero no período. Exibir "Pendente" ou similar.
    2.  **Opção Configurável/Alternativa:** Tratar notas ausentes como zero (se explicitamente configurado ou como um fallback claro).
    *Atualmente, o comportamento exato precisa ser verificado. Teste ou consulte `./vigenda avaliacao media-turma --help`.*
*   **Média Final:** Para calcular uma média final de todos os períodos, pode ser necessário um comando específico ou um valor especial para `--term` (ex: `--term "FINAL"`).

#### Dicas (`avaliacao media-turma`):
*   **Exportar Médias:** Para salvar as médias, redirecione a saída para um arquivo:
    ```bash
    ./vigenda avaliacao media-turma --classid 1 --term "1º Bimestre" > medias_turma1_bim1.txt
    ```
*   **Impacto de Novas Notas:** Execute `media-turma` após lançar/editar notas para ver o impacto imediato.

## 7. Banco de Questões e Geração de Provas

Mantenha um banco de questões organizado e gere provas personalizadas.

> **Nota sobre IDs:** IDs de Disciplina são usados para filtrar questões. IDs de Questão (se expostos pelo sistema) seriam usados para editar/remover questões individuais.

### Adicionar Questões ao Banco (`vigenda bancoq add`)

Importe questões para o seu banco de dados a partir de um ficheiro JSON.

**Uso:**
```bash
./vigenda bancoq add <CAMINHO_DO_ARQUIVO_JSON> [--subjectid <ID_DISCIPLINA_PADRAO>]
```

**Argumentos e Opções:**
*   `<CAMINHO_DO_ARQUIVO_JSON>`: (Obrigatório) O caminho para o ficheiro JSON contendo as questões.
*   `--subjectid <ID_DISCIPLINA_PADRAO>`: (Opcional) Se as questões no JSON não especificarem uma `disciplina` (ou `subject_id`), este ID será usado como padrão. Se o JSON contiver a informação de disciplina, ela terá precedência.

**Formato do JSON:** Consulte a seção [Importação de Questões (JSON)](#importacao-de-questoes-json) para detalhes.

**Exemplo:**
```bash
./vigenda bancoq add /data/questoes/historia_moderna.json
```
**Feedback Esperado:**
```
Importação de questões de /data/questoes/historia_moderna.json concluída.
Questões novas adicionadas: 25
Questões atualizadas (com base em ID interno, se houver): 3
Questões com erro: 1 (verifique o arquivo historia_moderna.json.errors)
```

#### Dicas (`bancoq add`):
*   **Múltiplos Arquivos:** Importe de vários arquivos executando o comando para cada um.
*   **Organização:** Mantenha arquivos JSON por disciplina ou tópico para facilitar a gestão.
*   **Duplicatas/Atualizações:** Verifique se o sistema previne duplicatas exatas ou se permite atualizar questões com base em algum identificador único na questão (além do ID gerado pelo banco). Idealmente, cada questão no JSON poderia ter um `id_externo` opcional para facilitar atualizações.

### Listar Questões do Banco (`vigenda bancoq listar`)

Visualize questões do banco, com filtros por disciplina, tópico, tipo ou dificuldade.

**Uso:**
```bash
./vigenda bancoq listar [--subjectid <ID_DISCIPLINA>] [--topic "Tópico"] [--type <TIPO>] [--difficulty <DIFICULDADE>] [--limit <NUM>]
```

**Opções:**
*   `--subjectid <ID_DISCIPLINA>`: (Opcional) Filtra por ID da disciplina.
*   `--topic "Tópico"`: (Opcional) Filtra por tópico (busca parcial ou exata, dependendo da implementação).
*   `--type <TIPO>`: (Opcional) Filtra por tipo de questão. Valores: `multipla_escolha`, `dissertativa`.
*   `--difficulty <DIFICULDADE>`: (Opcional) Filtra por dificuldade. Valores: `facil`, `media`, `dificil`.
*   `--limit <NUM>`: (Opcional) Limita o número de questões exibidas.

**Exemplo de Saída (Resumida):**
```
ID Questão | Disciplina ID | Tópico              | Tipo              | Dificuldade | Enunciado (início)
-----------|---------------|---------------------|-------------------|-------------|---------------------
501        | 1             | Revolução Francesa  | multipla_escolha  | media       | Qual destes eventos...
502        | 1             | Revolução Francesa  | dissertativa      | dificil     | Discorra sobre o impacto...
...
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda bancoq listar --help` para confirmar sua disponibilidade e formato de saída. O ID da Questão seria gerado pelo banco.)*

### Editar Questão no Banco (`vigenda bancoq editar`)

Modifique os detalhes de uma questão existente no banco. Requer o ID da questão.

**Uso:**
```bash
./vigenda bancoq editar <ID_DA_QUESTAO> [--subjectid <ID>] [--topic "Novo Tópico"] [--type <NOVO_TIPO>] [--difficulty <NOVA_DIFICULDADE>] [--enunciado "Novo Enunciado"] [--opcoes '["Nova Op1", "Nova Op2"]'] [--resposta "Nova Resposta"]
```

**Argumentos e Opções:**
*   `<ID_DA_QUESTAO>`: (Obrigatório) O ID da questão a ser editada (obtido de `bancoq listar`).
*   Demais opções permitem alterar os respectivos campos da questão. Para `--opcoes`, o formato JSON em string é esperado.

**Exemplo:**
```bash
./vigenda bancoq editar 501 --difficulty dificil --enunciado "Considerando o contexto da Revolução Francesa, qual destes eventos é o principal estopim?"
```
**Feedback Esperado:**
```
Questão ID 501 atualizada com sucesso.
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda bancoq editar --help` para confirmar sua disponibilidade.)*

### Remover Questão do Banco (`vigenda bancoq remover`)

Exclui permanentemente uma questão do banco.

**Uso:**
```bash
./vigenda bancoq remover <ID_DA_QUESTAO_1> [ID_DA_QUESTAO_2 ...] --force
```

**Argumentos e Opções:**
*   `<ID_DA_QUESTAO>`: (Obrigatório) O ID da questão a ser removida.
*   `--force`: (Opcional) Remove sem confirmação.

**Exemplo:**
```bash
./vigenda bancoq remover 501 --force
```
**Feedback Esperado:**
```
Questão ID 501 removida do banco com sucesso.
```
*(Nota: Este comando é uma adição sugerida. Verifique `./vigenda bancoq remover --help` para confirmar sua disponibilidade.)*

### Gerar Prova (`vigenda prova gerar`)

Crie provas personalizadas selecionando questões do banco.

**Uso:**
```bash
./vigenda prova gerar --subjectid <ID_DA_DISCIPLINA> [--topic "Tópico"] [--easy <NUM>] [--medium <NUM>] [--hard <NUM>] [--total <NUM>] [--type <TIPO_QUESTAO>] [--output <ARQUIVO_SAIDA.txt>] [--title "Título da Prova"]
```

**Opções:**
*   `--subjectid <ID_DA_DISCIPLINA>`: (Obrigatório) ID da disciplina.
*   `--topic "Tópico"`: (Opcional) Filtra questões por tópico.
*   `--easy <NUM>`, `--medium <NUM>`, `--hard <NUM>`: (Opcional) Número exato de questões por dificuldade.
*   `--total <NUM>`: (Opcional) Número total de questões. Se usado sem especificar dificuldades, o sistema tentará um balanceamento (ex: ~1/3 de cada, se disponíveis). Se usado com especificações de dificuldade, o `--total` deve ser a soma delas.
*   `--type <TIPO_QUESTAO>`: (Opcional) Filtra por tipo de questão (`multipla_escolha`, `dissertativa`, `todas`). Padrão: `todas`.
*   `--output <ARQUIVO_SAIDA.txt>`: (Opcional) Salva a prova em um arquivo de texto. Se omitido, exibe no console. **Este arquivo não inclui respostas.**
*   `--title "Título da Prova"`: (Opcional) Define um título para a prova, que aparecerá no cabeçalho do arquivo gerado.

**Exemplo:**
```bash
./vigenda prova gerar --subjectid 1 --topic "Revolução Francesa" --easy 2 --medium 2 --hard 1 --total 5 --output prova_rev_fr_01.txt --title "Prova 1: Revolução Francesa"
```
**Feedback Esperado:**
```
Prova "Prova 1: Revolução Francesa" gerada com 5 questões e salva em prova_rev_fr_01.txt.
Aviso: Não foram encontradas questões suficientes para 'hard' (solicitado 1, encontrado 0). A prova foi gerada com as questões disponíveis.
```
(O feedback incluirá avisos se o número solicitado de questões não puder ser atendido).

#### Dicas (`prova gerar`):
*   **Balanceamento com `--total`:** Se apenas `--total` for fornecido, o Vigenda tentará uma distribuição equilibrada de dificuldades disponíveis. Se você especificar, por exemplo, `--easy 5 --total 3`, o sistema priorizará as especificações de dificuldade e informará se o total não puder ser atingido.
*   **Insuficiência de Questões:** O sistema informará se não houver questões suficientes no banco para atender aos critérios. A prova será gerada com o máximo de questões que puder encontrar.
*   **Formato de Saída:** O arquivo `.txt` é formatado para fácil leitura e impressão, incluindo cabeçalho (com título, disciplina), questões numeradas, enunciados e opções (para múltipla escolha) ou espaço para resposta. **Não inclui respostas.**

### Gerar Gabarito (`vigenda prova gabarito`)

Gera um arquivo de gabarito (respostas) para uma prova previamente definida ou com base nos mesmos critérios de `prova gerar`.

**Uso:**
```bash
./vigenda prova gabarito --subjectid <ID_DA_DISCIPLINA> [--topic "Tópico"] [--easy <NUM>] [--medium <NUM>] [--hard <NUM>] [--total <NUM>] [--type <TIPO_QUESTAO>] --output <ARQUIVO_GABARITO.txt> [--title "Gabarito da Prova X"]
```
As opções de filtragem de questões são as mesmas de `vigenda prova gerar`. O importante é que os critérios de seleção de questões sejam **idênticos** aos usados para gerar a prova para que o gabarito corresponda.

**Opções Adicionais:**
*   `--output <ARQUIVO_GABARITO.txt>`: (Obrigatório) Nome do arquivo onde o gabarito será salvo.
*   `--for-prova <ARQUIVO_PROVA.txt>`: (Opcional, funcionalidade avançada) Se o Vigenda salvar metadados da prova gerada, esta opção poderia usar esses metadados para gerar o gabarito exato daquela prova. (Verificar com `--help`).

**Exemplo:**
```bash
./vigenda prova gabarito --subjectid 1 --topic "Revolução Francesa" --easy 2 --medium 2 --hard 1 --total 5 --output gabarito_rev_fr_01.txt --title "Gabarito - Prova 1: Revolução Francesa"
```
**Feedback Esperado:**
```
Gabarito "Gabarito - Prova 1: Revolução Francesa" gerado com 5 questões e salvo em gabarito_rev_fr_01.txt.
```
*(Nota: Este comando e a sua capacidade de corresponder exatamente a uma prova gerada são adições sugeridas. A forma mais simples é garantir que os parâmetros de seleção de questões sejam idênticos entre `prova gerar` e `prova gabarito`.)*


## 8. Formatos de Ficheiros de Importação

Esta seção detalha os formatos esperados para importação de dados.

### Importação de Alunos (CSV)

O comando `vigenda turma importar-alunos` espera um ficheiro CSV com as seguintes colunas (o cabeçalho na primeira linha é recomendado):

*   `nome_completo` (obrigatório): Nome completo do aluno.
*   `numero_chamada` (opcional): Número de chamada do aluno. Se fornecido, pode ser usado como um identificador secundário ou para ordenação. Verifique se o sistema o trata como único dentro da turma.
*   `email` (opcional): Endereço de e-mail do aluno.
*   `telefone` (opcional): Número de telefone do aluno.
*   `situacao` (opcional): Status inicial do aluno.
    *   Valores permitidos: `ativo`, `inativo`, `transferido`, `graduado`.
    *   Se omitido ou deixado em branco, o padrão é `ativo`.
*   Outros campos personalizados (opcional): O Vigenda pode permitir colunas adicionais que são armazenadas como metadados do aluno (verificar com `--help`).

**Exemplo (`alunos.csv`):**
```csv
nome_completo,numero_chamada,email,situacao
"Ana Beatriz Costa",1,"ana.costa@email.com","ativo"
"Bruno Dias",2,,"ativo"
"Carlos Eduardo Lima",,,"inativo"
"Daniel Mendes",4,"dani.mendes@email.com","transferido"
```
No exemplo acima:
*   "Ana Beatriz Costa" tem todos os campos preenchidos.
*   "Bruno Dias" não tem e-mail, mas sua situação é `ativo`.
*   "Carlos Eduardo Lima" não tem número de chamada nem e-mail, e está `inativo`.
*   "Daniel Mendes" está `transferido`.

**Observações Importantes:**
*   **Codificação:** Use codificação UTF-8 para o arquivo CSV para garantir a correta importação de caracteres especiais.
*   **Delimitador:** O delimitador padrão é vírgula (`,`). Se o seu CSV usa outro delimitador (ex: ponto e vírgula `;`), verifique se o comando `importar-alunos` possui uma opção para especificar o delimitador (ex: `--delimiter ";" `).
*   **Aspas:** Nomes ou outros campos contendo vírgulas devem estar entre aspas duplas (ex: `"Silva, João"`).

### Importação de Questões (JSON)

O comando `vigenda bancoq add` espera um ficheiro JSON contendo uma lista (array) de objetos, onde cada objeto representa uma questão.

**Estrutura de cada objeto de questão:**

*   `id_externo` (string, opcional): Um identificador único que você atribui à questão. Útil para futuras atualizações da questão se o sistema suportar a edição baseada neste ID.
*   `disciplina` (string, obrigatório): Nome da disciplina à qual a questão pertence (Ex: "História"). Este nome deve corresponder a uma disciplina existente no sistema ou ser usado em conjunto com a opção `--subjectid` no comando `bancoq add` se a disciplina não estiver no JSON.
*   `subject_id` (integer, opcional): ID numérico da disciplina. Se fornecido, tem precedência sobre o campo `disciplina` em nome.
*   `topico` (string, opcional): Tópico específico da questão dentro da disciplina (Ex: "Revolução Francesa", "Geometria Espacial").
*   `tipo` (string, obrigatório): Tipo da questão.
    *   Valores permitidos: `multipla_escolha`, `dissertativa`, `verdadeiro_falso`.
*   `dificuldade` (string, obrigatório): Nível de dificuldade.
    *   Valores permitidos: `facil`, `media`, `dificil`. (Pode haver mais níveis dependendo da configuração).
*   `enunciado` (string, obrigatório): O texto da questão. Pode conter formatação Markdown básica se o sistema de exibição de provas suportar.
*   `opcoes` (array de strings, obrigatório para `multipla_escolha` e `verdadeiro_falso`):
    *   Para `multipla_escolha`: Uma lista das opções de resposta (ex: `["Paris", "Londres", "Berlim"]`).
    *   Para `verdadeiro_falso`: Deve conter duas opções, geralmente `["Verdadeiro", "Falso"]` ou `["Certo", "Errado"]`.
*   `resposta_correta` (string ou array de strings, obrigatório):
    *   Para `multipla_escolha` (resposta única): O texto exato de uma das `opcoes` (ex: `"Paris"`).
    *   Para `multipla_escolha` (múltiplas respostas corretas, se suportado): Um array com os textos das opções corretas (ex: `["Paris", "Berlim"]`). (Verificar se o sistema suporta).
    *   Para `dissertativa`: O gabarito ou uma descrição da resposta esperada.
    *   Para `verdadeiro_falso`: O texto exato de uma das `opcoes` (ex: `"Verdadeiro"`).
*   `tags` (array de strings, opcional): Palavras-chave ou tags para classificar melhor a questão (Ex: `["século XVIII", "europa", "politica"]`).
*   `tempo_estimado_min` (integer, opcional): Tempo estimado em minutos para responder à questão.

**Exemplo (`questoes.json`):**
```json
[
  {
    "id_externo": "HIST-RF-001",
    "disciplina": "História",
    "topico": "Revolução Francesa",
    "tipo": "multipla_escolha",
    "dificuldade": "media",
    "enunciado": "Qual destes eventos é considerado o estopim da Revolução Francesa?",
    "opcoes": [
      "A Queda da Bastilha",
      "A convocação dos Estados Gerais",
      "O Juramento da Quadra de Tênis",
      "A Fuga de Varennes"
    ],
    "resposta_correta": "A Queda da Bastilha",
    "tags": ["revolucao francesa", "marco inicial"],
    "tempo_estimado_min": 2
  },
  {
    "disciplina": "Matemática",
    "subject_id": 2,
    "topico": "Álgebra",
    "tipo": "dissertativa",
    "dificuldade": "facil",
    "enunciado": "Explique o que é uma equação de primeiro grau e dê um exemplo.",
    "resposta_correta": "Uma equação de primeiro grau é uma igualdade matemática que envolve uma ou mais incógnitas com expoente 1 e que pode ser escrita na forma ax + b = 0, onde a e b são constantes e a ≠ 0. Exemplo: 2x + 3 = 7, que pode ser reescrita como 2x - 4 = 0.",
    "tempo_estimado_min": 5
  },
  {
    "disciplina": "Ciências",
    "topico": "Fotossíntese",
    "tipo": "verdadeiro_falso",
    "dificuldade": "facil",
    "enunciado": "A fotossíntese ocorre apenas em plantas.",
    "opcoes": ["Verdadeiro", "Falso"],
    "resposta_correta": "Falso",
    "tags": ["biologia celular", "seres vivos"]
  }
]
```

## 9. Configuração da Base de Dados

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
    *   **Padrão**: Um ficheiro `vigenda.db` no diretório de configuração do utilizador (ex: `~/.config/vigenda/vigenda.db` no Linux, ou `%APPDATA%\Vigenda\vigenda.db` no Windows) ou no diretório atual se o diretório de configuração não for acessível/gravável. Certifique-se de que o Vigenda tem permissões de escrita para o local escolhido.
    *   **Exemplo**: `export VIGENDA_DB_PATH="/var/data/meu_vigenda.db"`
    *   **Backup e Restauração (SQLite):**
        *   **Backup:** A forma mais simples e segura de fazer backup de um banco de dados SQLite é copiar o arquivo `.db` para um local seguro. É altamente recomendável que a aplicação Vigenda **não esteja em execução** ao copiar o arquivo para garantir a consistência dos dados.
            ```bash
            # Certifique-se que o Vigenda não está rodando
            cp ~/.config/vigenda/vigenda.db /caminho/do/backup/vigenda_backup_$(date +%Y%m%d_%H%M%S).db
            ```
        *   **Restauração:** Para restaurar, com o Vigenda parado, substitua o arquivo `.db` ativo pelo arquivo de backup desejado.
            ```bash
            # Certifique-se que o Vigenda não está rodando
            cp /caminho/do/backup/vigenda_backup_YYYYMMDD_HHMMSS.db ~/.config/vigenda/vigenda.db
            ```
        *   **Backup Online (Avançado):** O `sqlite3` CLI oferece um comando de backup online que pode ser usado mesmo com a aplicação rodando, pois lida com locks.
            ```bash
            sqlite3 ~/.config/vigenda/vigenda.db ".backup /caminho/do/backup/vigenda_online_backup.db"
            ```
            No entanto, para máxima segurança, um backup offline (com a aplicação parada) é preferível.

#### Configuração Específica para PostgreSQL

Se `VIGENDA_DB_TYPE` for `postgres` e `VIGENDA_DB_DSN` não for fornecida, as seguintes variáveis são usadas para construir a DSN:

*   `VIGENDA_DB_HOST`: Endereço do servidor PostgreSQL. (Padrão: `localhost`)
*   `VIGENDA_DB_PORT`: Porta do servidor PostgreSQL. (Padrão: `5432`)
*   `VIGENDA_DB_USER`: (Obrigatório) Nome de utilizador para a conexão.
*   `VIGENDA_DB_PASSWORD`: (Obrigatório se o servidor exigir) Senha para o utilizador.
*   `VIGENDA_DB_NAME`: (Obrigatório) Nome da base de dados PostgreSQL.
*   `VIGENDA_DB_SSLMODE`: Modo de SSL para a conexão PostgreSQL. (Padrão: `disable`. Valores comuns: `require`, `verify-ca`, `verify-full`).

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

*   **SQLite**: O Vigenda tentará aplicar o esquema inicial (localizado em `internal/database/migrations/001_initial_schema.sql`) automaticamente na primeira vez que se conectar a um arquivo de banco de dados SQLite vazio ou não existente.
*   **PostgreSQL**: Para PostgreSQL, as migrações de esquema devem ser gerenciadas e aplicadas externamente **antes** de usar o Vigenda com o banco de dados. O Vigenda não tentará criar tabelas ou modificar o esquema em um banco de dados PostgreSQL existente. Você pode usar ferramentas de migração padrão de PostgreSQL (como `psql` para executar os scripts SQL de `internal/database/migrations/` ou ferramentas como Flyway/Liquibase) para aplicar o esquema necessário. Consulte a [Documentação do Desenvolvedor (`docs/developer/README.md`)](docs/developer/README.md#migracoes-de-base-de-dados) para mais detalhes sobre a estrutura das migrações.

### Configuração Avançada com Docker

Usar o Vigenda com Docker pode simplificar o gerenciamento de dependências, a configuração do ambiente e a implantação.

#### 1. Construindo uma Imagem Docker para o Vigenda

Você precisará de um `Dockerfile` na raiz do seu projeto Vigenda. Aqui está um exemplo básico (assumindo que seu projeto é em Go e usa Go Modules):

```dockerfile
# Estágio 1: Build
# Use uma imagem Go compatível com sua versão do Go.
FROM golang:1.23-alpine AS builder

# Instalar dependências do sistema necessárias para compilação (ex: gcc para CGO, git)
RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

# Copiar go.mod e go.sum para otimizar o cache de dependências do Docker
COPY go.mod go.sum ./
# Baixar dependências
RUN go mod download
RUN go mod verify

# Copiar o restante do código fonte da aplicação
COPY . .

# Compilar a aplicação Vigenda.
# CGO_ENABLED=1 é geralmente necessário para go-sqlite3.
# ldflags para reduzir o tamanho do binário são opcionais.
# O output é /vigenda dentro do builder.
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /vigenda ./cmd/vigenda/

# Estágio 2: Imagem final leve
# Use uma imagem base pequena como alpine.
FROM alpine:latest

# (Opcional) Adicionar certificados CA, se sua aplicação precisar fazer chamadas HTTPS
# RUN apk add --no-cache ca-certificates

# Copiar o binário compilado do estágio de build para a imagem final
COPY --from=builder /vigenda /usr/local/bin/vigenda

# (Opcional, mas recomendado) Copiar os arquivos de migração SQL para a imagem,
# especialmente se você quiser que o container possa inicializar um banco SQLite do zero.
# COPY internal/database/migrations /app/internal/database/migrations
# Se fizer isso, certifique-se que VIGENDA_DB_PATH aponte para um local onde o esquema
# possa ser aplicado se o DB não existir.

# Definir o diretório de trabalho padrão para o container.
# É uma boa prática usar um diretório dedicado para dados persistentes se houver.
WORKDIR /data

# Definir o entrypoint padrão para o container.
# Isso permite que você execute o container como se fosse o próprio binário.
ENTRYPOINT ["vigenda"]

# (Opcional) Definir um CMD padrão se você quiser que o container execute algo
# específico ao iniciar sem argumentos (ex: exibir ajuda).
CMD ["--help"]
```

Para construir a imagem Docker, navegue até o diretório raiz do seu projeto (onde o `Dockerfile` está) e execute:
```bash
docker build -t vigenda-app .
```
(Substitua `vigenda-app` pelo nome que desejar para sua imagem).

#### 2. Executando o Vigenda com SQLite em Docker

Para persistir os dados do SQLite entre execuções do container, você **deve** montar um volume do host para o container onde o arquivo do banco de dados SQLite será armazenado.

*   **Exemplo de execução com volume para SQLite:**
    Supondo que seu `Dockerfile` tenha `WORKDIR /data` e você configure `VIGENDA_DB_PATH` para salvar o banco dentro deste diretório.
    ```bash
    # Criar um diretório no host para persistir os dados do SQLite, se ainda não existir
    # mkdir -p $(pwd)/my_vigenda_sqlite_data

    docker run -it --rm \
      -v $(pwd)/my_vigenda_sqlite_data:/data \
      -e VIGENDA_DB_TYPE="sqlite" \
      -e VIGENDA_DB_PATH="/data/vigenda_production.db" \
      vigenda-app tarefa listar
    ```
    Neste exemplo:
    *   `-v $(pwd)/my_vigenda_sqlite_data:/data`: Mapeia o subdiretório `my_vigenda_sqlite_data` do seu diretório atual no host para o diretório `/data` dentro do container. O arquivo `vigenda_production.db` será criado/lido neste local do host.
    *   `-e VIGENDA_DB_PATH="/data/vigenda_production.db"`: Informa ao Vigenda dentro do container para usar este caminho para o arquivo do banco de dados.
    *   `--rm`: Remove o container após a execução. Se você quiser que o container persista (por exemplo, se fosse um serviço), você omitiria `--rm` e poderia dar um nome com `--name`.
    *   `-it`: Para interação com o CLI.

#### 3. Executando o Vigenda com PostgreSQL em Docker (Conectando a um Servidor PostgreSQL)

Se você tem um servidor PostgreSQL rodando (seja no host, em outro container Docker, ou em um serviço de nuvem), você configura o Vigenda dentro do container para se conectar a ele usando as variáveis de ambiente `VIGENDA_DB_*` apropriadas.

**Exemplo: Vigenda em um container conectando a um PostgreSQL em outro container:**

1.  **Crie uma rede Docker customizada (ponte):**
    Isso permite que os containers se comuniquem usando seus nomes como hostnames.
    ```bash
    docker network create vigenda-net
    ```

2.  **Inicie o container do PostgreSQL (Exemplo):**
    ```bash
    docker run --name my-postgres-db --network vigenda-net \
      -e POSTGRES_USER=vigenda_user \
      -e POSTGRES_PASSWORD=your_strong_password \
      -e POSTGRES_DB=vigenda_database \
      -p 5432:5432 \
      -v $(pwd)/pgdata:/var/lib/postgresql/data \
      -d postgres:15-alpine
    ```
    *   `--name my-postgres-db`: Nome do container PostgreSQL.
    *   `--network vigenda-net`: Conecta à rede criada.
    *   `-e ...`: Define usuário, senha e nome do banco de dados.
    *   `-p 5432:5432`: Mapeia a porta do PostgreSQL para o host (opcional se apenas o container Vigenda precisar acessá-lo).
    *   `-v $(pwd)/pgdata:/var/lib/postgresql/data`: Persiste os dados do PostgreSQL no host.
    *   `-d postgres:15-alpine`: Executa em modo detached usando a imagem oficial do PostgreSQL.
    *   **Lembre-se de aplicar as migrações de esquema ao `my-postgres-db` externamente.**

3.  **Execute o container do Vigenda conectado ao PostgreSQL:**
    ```bash
    docker run -it --rm --network vigenda-net \
      -e VIGENDA_DB_TYPE="postgres" \
      -e VIGENDA_DB_HOST="my-postgres-db" \
      -e VIGENDA_DB_PORT="5432" \
      -e VIGENDA_DB_USER="vigenda_user" \
      -e VIGENDA_DB_PASSWORD="your_strong_password" \
      -e VIGENDA_DB_NAME="vigenda_database" \
      -e VIGENDA_DB_SSLMODE="disable" \
      vigenda-app tarefa listar
    ```
    *   `--network vigenda-net`: Conecta o container Vigenda à mesma rede.
    *   `VIGENDA_DB_HOST="my-postgres-db"`: O Vigenda usa o nome do container do PostgreSQL como hostname. Docker DNS resolverá isso dentro da rede `vigenda-net`.

Esta seção fornece uma visão geral. Configurações de Docker podem se tornar mais complexas dependendo dos requisitos específicos de segurança, volumes, e orquestração (ex: Docker Compose).

## 10. Dicas de Uso e Boas Práticas

*   **Consistência é Chave:** Use nomes consistentes para turmas, períodos avaliativos e tópicos de questões. Isso facilitará a filtragem e a organização.
*   **Anote os IDs:** Ao criar turmas, avaliações, etc., o Vigenda geralmente retorna um ID. Anote esses IDs, pois são necessários para muitas operações de edição ou referência. Comandos de listagem (`listar`) são seus amigos para encontrar IDs posteriormente.
*   **Backup Regular (SQLite):** Se estiver usando SQLite, faça backups regulares do seu arquivo `vigenda.db`, especialmente antes de atualizações do Vigenda ou grandes importações de dados. Consulte a seção [Configuração Específica para SQLite](#configuracao-especifica-para-sqlite) para mais detalhes.
*   **Arquivos de Importação:** Mantenha seus arquivos CSV de alunos e JSON de questões organizados e com backup. Eles são a fonte dos seus dados e podem ser úteis para reimportações ou correções em massa.
*   **Explore com `--help`:** Cada comando e subcomando do Vigenda oferece uma opção `--help` que detalha seu uso, argumentos e opções. Use-a frequentemente para descobrir todas as capacidades.
*   **Scripts:** Para tarefas repetitivas (ex: criar múltiplas turmas com um padrão, atualizar status de vários alunos), considere usar scripts de shell (Bash, PowerShell, etc.) que invocam os comandos do Vigenda.

## 11. Solução de Problemas Comuns (FAQ)

*   **P: Esqueci o ID de uma turma/avaliação/tarefa. Como posso encontrá-lo?**
    *   **R:** Use os comandos de listagem apropriados:
        *   `./vigenda turma listar` para turmas.
        *   `./vigenda avaliacao listar --classid <ID_DA_TURMA>` para avaliações de uma turma.
        *   `./vigenda tarefa listar [--classid <ID_DA_TURMA>]` para tarefas.
        *   `./vigenda aluno listar --classid <ID_DA_TURMA>` para alunos de uma turma.
        *   `./vigenda bancoq listar [--subjectid <ID_DISCIPLINA>]` para questões.

*   **P: O comando `X` não funciona ou dá um erro de "argumento faltando".**
    *   **R:** Execute `./vigenda <comando> --help` (ex: `./vigenda tarefa add --help`) para ver todos os argumentos e opções obrigatórios e opcionais, juntamente com seus formatos esperados. Certifique-se de que está fornecendo todas as informações necessárias.

*   **P: Ao importar um arquivo CSV/JSON, recebo um erro de "arquivo não encontrado".**
    *   **R:** Verifique se o caminho para o arquivo está correto. Se estiver usando um caminho relativo, certifique-se de que está executando o Vigenda a partir do diretório correto. Caminhos absolutos são geralmente mais seguros. Verifique também as permissões de leitura do arquivo.

*   **P: As médias calculadas por `avaliacao media-turma` não parecem corretas.**
    *   **R:** Verifique os seguintes pontos:
        1.  Todas as notas para o período foram lançadas?
        2.  Os pesos (`--weight`) das avaliações foram definidos corretamente ao criá-las?
        3.  O período (`--term`) especificado no comando `media-turma` é exatamente o mesmo usado ao criar as avaliações?
        4.  Consulte a seção [Detalhes do Cálculo (`avaliacao media-turma`)](#detalhes-do-calculo-avaliacao-media-turma) para entender como as notas ausentes são tratadas.

*   **P: Posso editar uma questão diretamente no banco de dados se eu souber SQL?**
    *   **R:** Embora tecnicamente possível (especialmente com SQLite), **não é recomendado** editar o banco de dados diretamente. Isso pode levar a inconsistências de dados, corromper o banco ou causar comportamento inesperado no Vigenda, especialmente se houver lógica na aplicação para manter a integridade dos dados. Sempre prefira usar os comandos fornecidos pelo Vigenda para interagir com os dados (ex: `bancoq editar`, se disponível).

---

Para mais informações sobre como instalar e começar a usar o Vigenda rapidamente, consulte o [Guia de Introdução (`docs/getting_started/README.md`)](../getting_started/README.md).
Para exemplos práticos, explore nossos [Tutoriais (`docs/tutorials/`)](../tutorials/).
Se você é um desenvolvedor, a [Documentação do Desenvolvedor (`docs/developer/README.md`)](developer/README.md) pode ser útil.
