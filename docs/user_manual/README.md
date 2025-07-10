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

O Vigenda é uma aplicação de linha de comando (CLI) desenvolvida em Go, projetada para ajudar professores e estudantes a gerenciar suas atividades acadêmicas. A principal forma de interação é através de uma **Interface de Texto do Usuário (TUI)** robusta, acessada executando `vigenda` (ou `go run ./cmd/vigenda/main.go`) sem subcomandos.

Esta TUI oferece um menu principal para acessar todas as funcionalidades de forma interativa, incluindo a criação e gerenciamento de disciplinas, turmas, alunos, aulas, avaliações, etc.

Adicionalmente, alguns subcomandos CLI estão disponíveis para acesso direto a funcionalidades específicas.

## 2. Usando a Interface Principal (TUI)

Para iniciar a TUI do Vigenda:
```bash
# Se estiver executando do código fonte
go run ./cmd/vigenda/main.go

# Ou, se você construiu o binário (ex: ./vigenda_cli)
./vigenda_cli
```
Isso abrirá o menu principal da aplicação. Use as teclas de seta para navegar e Enter para selecionar uma opção.

### 2.1. Painel de Controle (Dashboard)
Uma das opções centrais da TUI é o "Painel de Controle". Ele fornece uma visão geral da sua agenda, tarefas urgentes e outras notificações.

### 2.2. Outras Funcionalidades da TUI
A TUI permite gerenciar:
*   **Disciplinas:** Criar, listar, editar e remover disciplinas.
*   **Turmas:** Criar turmas dentro de disciplinas, listar, editar.
*   **Alunos:** Adicionar alunos a turmas (além da importação por CSV).
*   **Aulas:** Planejar e visualizar aulas.
*   **Avaliações:** Criar e gerenciar avaliações (além do comando CLI).
*   E mais. Explore os menus para descobrir todas as funcionalidades.

**Dicas de Navegação na TUI:**
*   **Setas (Cima/Baixo):** Navegar entre itens de menu ou campos.
*   **Enter:** Selecionar uma opção ou confirmar uma entrada.
*   **Esc (Escape):** Voltar ao menu anterior ou sair de um formulário/visualização.
*   **Teclas de Atalho Específicas:** Algumas telas podem ter teclas de atalho (ex: 'r' para recarregar no dashboard). Observe as dicas na tela.

## 3. Comandos CLI Detalhados

Embora a TUI seja a interface principal para muitas operações de criação e edição, os seguintes subcomandos CLI oferecem acesso direto a funcionalidades importantes. Para obter ajuda sobre qualquer comando, use `vigenda [comando] --help`.

### Gestão de Tarefas

#### Adicionar Tarefa (`vigenda tarefa add`)
Cria rapidamente uma nova tarefa.
**Uso:**
```bash
./vigenda tarefa add "Descrição da Tarefa" [--classid ID_DA_TURMA] [--duedate AAAA-MM-DD] [--description "Detalhes"]
```
*   `"Descrição da Tarefa"`: Título/descrição curta (obrigatório).
*   `--classid ID_DA_TURMA`: (Opcional) ID da turma para associar a tarefa.
*   `--duedate AAAA-MM-DD`: (Opcional) Data de conclusão.
*   `--description "Detalhes"`: (Opcional) Descrição mais longa. Se não fornecida e o sistema detectar um terminal interativo, pode solicitar.

**Exemplo:**
```bash
./vigenda tarefa add "Preparar slides Aula 5" --classid 1 --duedate 2024-08-15
```

#### Listar Tarefas (`vigenda tarefa listar`)
Visualiza tarefas.
**Uso:**
```bash
./vigenda tarefa listar [--classid ID_DA_TURMA] [--all]
```
*   `--classid ID_DA_TURMA`: (Opcional) Filtra tarefas pela ID da turma. Se esta flag for usada, `--all` é ignorada.
*   `--all`: (Opcional) Lista todas as tarefas de todas as turmas e também tarefas do sistema (bugs) que não têm `classid`.
*   Se nenhuma flag for fornecida, o comportamento padrão pode variar (consulte `cmd/vigenda/main.go` ou teste; idealmente, listaria tarefas do usuário atual ou pediria um filtro). A implementação atual em `main.go` exige `--classid` OU `--all`.

**Exemplos:**
```bash
./vigenda tarefa listar --classid 1 # Tarefas da turma 1
./vigenda tarefa listar --all       # Todas as tarefas (de todas as turmas e do sistema)
```

#### Completar Tarefa (`vigenda tarefa complete`)
Marca uma tarefa como concluída.
**Uso:**
```bash
./vigenda tarefa complete ID_DA_TAREFA
```
**Exemplo:**
```bash
./vigenda tarefa complete 42
```

### Gestão de Turmas e Alunos

A criação e edição detalhada de turmas e disciplinas é primariamente feita via TUI. Os comandos CLI abaixo são para operações específicas.

#### Importar Alunos (`vigenda turma importar-alunos`)
Importa alunos de um CSV para uma turma existente. A turma deve ser criada previamente via TUI.
**Uso:**
```bash
./vigenda turma importar-alunos ID_DA_TURMA CAMINHO_DO_ARQUIVO_CSV
```
**Formato do CSV:** Veja a seção [Importação de Alunos (CSV)](#importacao-de-alunos-csv).
**Exemplo:**
```bash
./vigenda turma importar-alunos 2 alunos_turma_2B.csv
```

#### Atualizar Status do Aluno (`vigenda turma atualizar-status`)
Modifica o status de um aluno. O aluno e a turma devem existir.
**Uso:**
```bash
./vigenda turma atualizar-status ID_DO_ALUNO NOVO_STATUS
```
*   `NOVO_STATUS`: Valores permitidos: `ativo`, `inativo`, `transferido`.
**Exemplo:**
```bash
./vigenda turma atualizar-status 101 transferido
```

### Gestão de Avaliações e Notas

A criação de disciplinas e turmas, que são pré-requisitos para avaliações, é feita via TUI.

#### Criar Avaliação (`vigenda avaliacao criar`)
Define uma nova avaliação para uma turma existente.
**Uso:**
```bash
./vigenda avaliacao criar "Nome da Avaliação" --classid ID_DA_TURMA --term PERIODO --weight PESO [--date AAAA-MM-DD]
```
*   `"Nome da Avaliação"`: Obrigatório.
*   `--classid ID_DA_TURMA`: Obrigatório. A turma deve existir.
*   `--term PERIODO`: Obrigatório (ex: "1", "2", "1º Bimestre").
*   `--weight PESO`: Obrigatório (ex: 3.0, 5.0).
*   `--date AAAA-MM-DD`: (Opcional) Data da avaliação.

**Exemplo:**
```bash
./vigenda avaliacao criar "P1 - Unidade I" --classid 1 --term "1" --weight 4.0 --date 2024-03-15
```

#### Lançar Notas (`vigenda avaliacao lancar-notas`)
Inicia um prompt interativo para lançar notas de uma avaliação existente.
**Uso:**
```bash
./vigenda avaliacao lancar-notas ID_DA_AVALIACAO
```
**Exemplo:**
```bash
./vigenda avaliacao lancar-notas 5
```

#### Calcular Média da Turma (`vigenda avaliacao media-turma`)
Calcula a média ponderada de uma turma existente.
**Uso:**
```bash
./vigenda avaliacao media-turma ID_DA_TURMA
```
**Exemplo:**
```bash
./vigenda avaliacao media-turma 1
```
*(Nota: A flag `--term` para este comando não está presente na implementação CLI atual em `main.go`, portanto a média calculada será geral para a turma, considerando todas as suas avaliações com notas lançadas.)*

### Banco de Questões e Geração de Provas

A criação de disciplinas, pré-requisito para associar questões, é feita via TUI.

#### Adicionar Questões ao Banco (`vigenda bancoq add`)
Importa questões de um arquivo JSON para disciplinas existentes.
**Uso:**
```bash
./vigenda bancoq add CAMINHO_DO_ARQUIVO_JSON
```
**Formato do JSON:** Veja a seção [Importação de Questões (JSON)](#importacao-de-questoes-json). As disciplinas mencionadas no JSON devem existir.
**Exemplo:**
```bash
./vigenda bancoq add questoes_geografia.json
```

#### Gerar Prova (`vigenda prova gerar`)
Cria uma prova selecionando questões do banco para uma disciplina existente.
**Uso:**
```bash
./vigenda prova gerar --subjectid ID_DA_DISCIPLINA [--topic "Tópico"] --easy NUM --medium NUM --hard NUM [--output ARQUIVO.txt]
```
*   `--subjectid ID_DA_DISCIPLINA`: Obrigatório. A disciplina deve existir.
*   `--easy/medium/hard NUM`: Número de questões por dificuldade. Pelo menos uma contagem deve ser > 0.
*   `--output ARQUIVO.txt`: (Opcional) Salva a prova em um arquivo.

**Exemplo:**
```bash
./vigenda prova gerar --subjectid 1 --easy 5 --medium 3 --hard 2 --output prova_hist.txt
```

## 4. Formatos de Ficheiros de Importação
