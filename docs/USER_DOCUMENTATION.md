# Documentação do Usuário do Vigenda

## Introdução

Bem-vindo ao Vigenda!

O Vigenda é uma aplicação de linha de comando com uma **Interface de Texto do Usuário (TUI)**, projetada para ajudar professores e educadores a gerenciar suas atividades acadêmicas de forma eficiente. Com ele, você pode organizar turmas, planejar aulas, gerenciar tarefas, criar avaliações, manter um banco de questões e gerar provas personalizadas.

Esta documentação serve como um guia completo para todas as funcionalidades que o Vigenda oferece.

## Como Iniciar

Para começar a usar o Vigenda, abra seu terminal e execute o seguinte comando no diretório de instalação:

```bash
./vigenda
```

Isso iniciará a interface principal do Vigenda, onde você poderá navegar por todas as funcionalidades usando as **teclas de seta**, selecionar opções com **Enter** e voltar aos menus anteriores com **Esc**.

## Funcionalidades Principais

O menu principal do Vigenda oferece acesso às seguintes seções:

### 1. Painel de Controle (Dashboard)

A tela inicial que fornece uma visão geral rápida das suas informações mais importantes, como tarefas pendentes e próximas avaliações.

### 2. Gestão de Tarefas

Crie, visualize, edite e marque tarefas como concluídas. As tarefas podem ser gerais ou associadas a turmas específicas, ajudando a manter o controle das atividades dos alunos.

### 3. Gestão de Turmas e Disciplinas

O coração da organização acadêmica.
- **Crie Disciplinas:** Agrupe seu conteúdo por matéria (ex: "Matemática", "História").
- **Crie Turmas:** Organize seus alunos em turmas dentro de cada disciplina (ex: "Turma 101 - Manhã").
- **Adicione Alunos:** Gerencie a lista de alunos para cada turma.

### 4. Gestão de Avaliações e Notas

Controle todo o ciclo de vida das avaliações:
- **Crie Avaliações:** Defina avaliações com nome, período, peso e data.
- **Lance Notas:** Insira as notas dos alunos para cada avaliação de forma interativa.
- **Calcule Médias:** Visualize as médias ponderadas da turma automaticamente.

### 5. Banco de Questões

Construa um repositório centralizado de questões para suas provas e atividades.
- **Adicione Questões:** Crie questões de múltipla escolha ou dissertativas, classificando-as por disciplina, tópico e nível de dificuldade.
- **Importe em Lote:** Use arquivos JSON para importar dezenas de questões de uma só vez (via linha de comando).

### 6. Geração de Provas

Gere provas personalizadas em segundos a partir do seu banco de questões.
- **Filtre Questões:** Selecione questões por disciplina, tópico e dificuldade.
- **Gere Arquivos:** Crie arquivos de texto com a prova formatada, pronta para impressão.

## Tutoriais e Exemplos

Para um guia passo a passo sobre como usar essas funcionalidades, consulte os seguintes documentos:
- **[Exemplos de Uso da TUI](./user_manual/TUI_EXAMPLES.md):** Um guia rápido para realizar as tarefas mais comuns diretamente na interface.
- **[Tutorial: Gerenciando Avaliações](./tutorials/02_gerenciando_avaliacoes_e_lancando_notas.md):** Um passo a passo detalhado sobre o fluxo de avaliações.
- **[Tutorial: Dominando o Banco de Questões](./tutorials/03_dominando_banco_questoes_e_geracao_provas.md):** Um guia completo para criar e usar o banco de questões e gerar provas.

## Instalação

Para instruções sobre como instalar o Vigenda em seu sistema, consulte o arquivo [**INSTALLATION.MD**](../INSTALLATION.MD).

## Suporte

Se encontrar problemas ou tiver dúvidas, consulte a seção de [**Relato de Bugs**](../BUG_REPORTING.MD).
