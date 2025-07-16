# Exemplos de Uso do Vigenda

Este arquivo fornece exemplos práticos de como usar a **Interface de Texto do Usuário (TUI)** do Vigenda, que é a forma principal e recomendada de interagir com a aplicação.

**Nota:**
- Assumimos que você já [instalou o Vigenda e configurou seu ambiente](INSTALLATION.MD).
- Para iniciar a TUI, execute `vigenda` no seu terminal.
- Use as **teclas de seta** para navegar, **Enter** para selecionar e **Esc** para voltar.

---

## 1. Fluxo de Exemplo: Configurando uma Nova Disciplina e Turma

1.  **Iniciar o Vigenda (TUI):**
    ```bash
    vigenda
    ```
    Isso abrirá o Menu Principal.

2.  **Criar uma Nova Disciplina:**
    *   No Menu Principal, navegue até **"Gerenciar Turmas e Disciplinas"** e pressione **Enter**.
    *   Escolha a opção **"Criar Nova Disciplina"**.
    *   Digite o nome da disciplina (ex: "Biologia Celular") e pressione **Enter**.

3.  **Criar uma Nova Turma:**
    *   Na lista de disciplinas, selecione "Biologia Celular".
    *   Escolha **"Criar Nova Turma"**.
    *   Digite o nome da turma (ex: "BIO-101 Manhã 2024") e pressione **Enter**.

4.  **Adicionar Alunos à Turma:**
    *   Na tela da turma, selecione **"Adicionar Novo Aluno"**.
    *   Siga os prompts para inserir o nome completo e a matrícula do aluno.

5.  **Criar uma Avaliação:**
    *   Ainda na tela da turma, selecione **"Criar Nova Avaliação"**.
    *   Preencha o nome (ex: "Prova Parcial 1"), o período, o peso e a data.

6.  **Lançar Notas:**
    *   Selecione a avaliação criada e escolha **"Lançar/Editar Notas"**.
    *   Insira a nota para cada aluno da lista.

---

## 2. Operações Adicionais na TUI

### Gerenciamento de Tarefas

*   No Menu Principal, selecione **"Gerenciar Tarefas"**.
*   Use as opções para criar novas tarefas, associá-las a turmas, definir prazos e marcá-las como concluídas.

### Banco de Questões e Geração de Provas

*   **Importar Questões (via CLI):** A maneira mais eficiente de popular seu banco de questões é usando um arquivo JSON e o comando de linha de comando. (Veja a seção CLI abaixo).
*   **Gerar Provas (na TUI):**
    *   No Menu Principal, selecione **"Gerar Provas"**.
    *   Selecione a disciplina.
    *   Defina os critérios (tópico, número de questões por dificuldade).
    *   Digite o nome do arquivo de saída (ex: `prova_biologia.txt`).
    *   A prova será gerada no arquivo de texto especificado.

---

## 3. Exemplos de Comandos CLI (Uso Avançado/Scripts)

Embora a TUI seja a interface principal, alguns comandos CLI são úteis para automação e operações em lote.

### Adicionar questões ao banco a partir de um arquivo JSON

Este é o principal uso recomendado para a CLI.
```bash
vigenda bancoq add minhas_questoes.json
```
O arquivo JSON deve conter uma lista de questões, cada uma com disciplina, tópico, tipo, dificuldade, enunciado e resposta.

### Importar alunos para uma turma via CSV

Supondo que a Turma com ID `1` foi criada via TUI:
```bash
vigenda turma importar-alunos 1 alunos.csv
```
O arquivo `alunos.csv` deve ter as colunas `nome_completo` e `matricula`.

### Listar tarefas rapidamente

```bash
vigenda tarefa listar --classid 1
```

Para uma lista completa de comandos e suas opções, use `vigenda --help` e `vigenda <comando> --help`.
