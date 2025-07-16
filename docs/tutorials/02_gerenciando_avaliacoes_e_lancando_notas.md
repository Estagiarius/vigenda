# Tutorial 02: Gerenciando Avaliações e Lançando Notas com a TUI

Este tutorial detalha como gerenciar o ciclo de vida das avaliações no Vigenda usando a **Interface de Texto do Usuário (TUI)**, desde a criação de múltiplas avaliações para uma turma, atribuindo pesos, lançando notas de forma interativa, até o cálculo da média da turma.

**Pré-requisitos:**
*   Vigenda instalado e funcionando. Consulte [**INSTALLATION.MD**](../../INSTALLATION.MD).
*   Uma **disciplina** e uma **turma** já devem ter sido criadas. Se não, siga o guia em [**Exemplos de Uso da TUI**](../../docs/user_manual/TUI_EXAMPLES.md).
*   Alunos já devem ter sido adicionados à turma.

## Cenário

Vamos simular a gestão de avaliações para a turma "História 9A" durante o "1º Bimestre". Planejamos três avaliações:
1.  **P1**: Prova escrita, peso 5.0.
2.  **T1**: Trabalho em grupo, peso 3.0.
3.  **Participacao**: Nota de participação, peso 2.0.

## Passo 1: Iniciar o Vigenda e Navegar até Avaliações

1.  Abra seu terminal e execute o Vigenda:
    ```bash
    ./vigenda
    ```
2.  No "Menu Principal", use as setas para navegar até **"Gerenciar Avaliações e Notas"** e pressione **Enter**.

## Passo 2: Criar as Avaliações

Você será apresentado a uma tela para selecionar a disciplina e a turma.
1.  Selecione a disciplina (ex: "História").
2.  Selecione a turma (ex: "História 9A").

Agora, dentro da tela de avaliações da turma, vamos criar cada uma das avaliações planejadas.

1.  **Criar a P1 (Prova 1):**
    *   Navegue e selecione a opção **"Criar Nova Avaliação"**.
    *   Preencha os campos solicitados:
        *   **Nome:** `P1 - Prova Bimestral`
        *   **Período/Bimestre:** `1º Bimestre`
        *   **Peso:** `5.0`
        *   **Data:** `2024-03-15` (formato AAAA-MM-DD)
    *   Confirme para salvar. A "P1" aparecerá na lista de avaliações.

2.  **Criar o T1 (Trabalho 1):**
    *   Selecione novamente **"Criar Nova Avaliação"**.
    *   Preencha os campos:
        *   **Nome:** `T1 - Trabalho: Egito Antigo`
        *   **Período/Bimestre:** `1º Bimestre`
        *   **Peso:** `3.0`
        *   **Data:** `2024-03-29`
    *   Confirme para salvar.

3.  **Criar a nota de Participação:**
    *   Selecione **"Criar Nova Avaliação"** mais uma vez.
    *   Preencha os campos:
        *   **Nome:** `Participacao - 1º Bimestre`
        *   **Período/Bimestre:** `1º Bimestre`
        *   **Peso:** `2.0`
        *   **Data:** (pode deixar em branco)
    *   Confirme para salvar.

Agora você deve ver as três avaliações listadas para a turma "História 9A".

## Passo 3: Lançar Notas para a Primeira Avaliação (P1)

1.  Na lista de avaliações, navegue até **"P1 - Prova Bimestral"** e pressione **Enter**.
2.  No menu da avaliação, selecione a opção **"Lançar/Editar Notas"**.
3.  A TUI exibirá a lista de alunos matriculados na turma. Para cada aluno, digite a nota e pressione **Enter** para ir para o próximo.
    ```
    Lançando notas para a Avaliação: P1 - Prova Bimestral
    - Aluno: Alice Wonderland | Nota: 7.5
    - Aluno: Bob The Builder  | Nota: 8.0
    - Aluno: Charles Xavier   | Nota: 6.0
    ... (continue para todos os alunos) ...
    ```
4.  Após inserir a última nota, siga as instruções na tela para **salvar** as notas.

## Passo 4: Calcular a Média Parcial (Após P1)

Após lançar as notas da P1, você pode verificar o impacto nas médias.

1.  Pressione **Esc** para voltar à tela da turma (onde as avaliações são listadas).
2.  Procure e selecione a opção **"Calcular Média da Turma"**.
3.  A TUI exibirá uma tabela com a média de cada aluno, considerando apenas as avaliações que já têm notas lançadas. Neste ponto, a média será baseada apenas na P1.

## Passo 5: Lançar Notas para as Demais Avaliações (T1 e Participação)

Repita o processo do Passo 3 para as outras avaliações:

1.  Na lista de avaliações, selecione **"T1 - Trabalho: Egito Antigo"** e pressione **Enter**.
2.  Escolha **"Lançar/Editar Notas"** e insira as notas do trabalho para cada aluno. Salve ao final.
3.  Volte, selecione **"Participacao - 1º Bimestre"**, pressione **Enter**.
4.  Escolha **"Lançar/Editar Notas"** e insira as notas de participação. Salve ao final.

## Passo 6: Calcular a Média Final do Bimestre

Com todas as notas do "1º Bimestre" lançadas, vamos verificar a média final ponderada.

1.  Volte para a tela da turma (onde as avaliações são listadas).
2.  Selecione **"Calcular Média da Turma"** novamente.

Agora, a TUI exibirá a média final de cada aluno, calculada com base nas notas e pesos das três avaliações.

**Entendendo o Cálculo:**
A média de cada aluno é calculada usando a fórmula da média ponderada:
`(Nota_P1 * Peso_P1 + Nota_T1 * Peso_T1 + Nota_Participacao * Peso_Participacao) / (Soma dos Pesos)`
No nosso exemplo, o divisor (soma dos pesos) é `5 + 3 + 2 = 10`.

## Dicas Adicionais

*   **Corrigindo Notas:** Para corrigir uma nota, basta navegar até a avaliação, selecionar "Lançar/Editar Notas" novamente e inserir a nota correta para o aluno desejado.
*   **Alunos sem Nota:** Observe como o sistema lida com alunos sem nota em uma avaliação. Dependendo da configuração, a média pode não ser calculada para esse aluno ou ele pode receber uma nota zero.

## Conclusão

Este tutorial demonstrou como usar a TUI do Vigenda para gerenciar múltiplas avaliações dentro de um período letivo de forma integrada. Você aprendeu a criar avaliações, lançar notas interativamente e calcular as médias ponderadas, tudo dentro da interface do Vigenda.
