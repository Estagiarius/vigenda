# Tutorial 02: Gerenciando Avaliações e Lançando Notas Detalhadamente

Este tutorial detalha como gerenciar o ciclo de vida das avaliações no Vigenda, desde a criação de múltiplas avaliações para uma turma, atribuindo pesos, lançando notas de forma interativa, até o cálculo da média da turma para verificar o impacto.

**Pré-requisitos:**
*   Vigenda instalado e funcionando. Consulte [**INSTALLATION.MD**](../../INSTALLATION.MD).
*   Uma **disciplina** e uma **turma** já devem ter sido criadas através da **Interface de Texto do Usuário (TUI)** do Vigenda.
    *   Para iniciar a TUI: execute `vigenda` (ou `go run ./cmd/vigenda/main.go`).
    *   No menu da TUI, crie uma disciplina (ex: "História") e, dentro dela, uma turma (ex: "História 9A"). Anote o ID da turma (vamos assumir `ID = 1` para este tutorial).
*   Alunos devem ter sido importados para a turma `ID = 1` (ex: usando o comando `vigenda turma importar-alunos 1 alunos.csv`).
*   Conhecimento básico dos comandos CLI `vigenda avaliacao criar`, `vigenda avaliacao lancar-notas`, e `vigenda avaliacao media-turma`. Consulte o [Manual do Usuário](../../docs/user_manual/README.md) para detalhes.

## Cenário

Vamos simular a gestão de avaliações para a turma `ID = 1` (ex: "História 9A") durante o "1º Bimestre". Planejamos três avaliações:
1.  **P1**: Prova escrita, peso 5.0.
2.  **T1**: Trabalho em grupo, peso 3.0.
3.  **Participacao**: Nota de participação, peso 2.0.

## Passo 1: Criar as Avaliações

Primeiro, vamos criar cada uma das avaliações planejadas para a Turma ID 1, no período "1º Bimestre".

1.  **Criar a P1 (Prova 1):**
    ```bash
    ./vigenda avaliacao criar "P1 - História 9A" --classid 1 --term "1º Bimestre" --weight 5.0 --date 2024-03-15
    ```
    O sistema deve retornar o ID da avaliação criada. Anote-o. Ex: `Avaliação "P1 - História 9A" criada com sucesso. ID: 10`

2.  **Criar o T1 (Trabalho 1):**
    ```bash
    ./vigenda avaliacao criar "T1 - Trabalho: Egito Antigo" --classid 1 --term "1º Bimestre" --weight 3.0 --date 2024-03-29
    ```
    Anote o ID. Ex: `Avaliação "T1 - Trabalho: Egito Antigo" criada com sucesso. ID: 11`

3.  **Criar a nota de Participação:**
    ```bash
    ./vigenda avaliacao criar "Participacao - 1º Bimestre" --classid 1 --term "1º Bimestre" --weight 2.0
    ```
    (A data é opcional, para participação pode não ser uma data única).
    Anote o ID. Ex: `Avaliação "Participacao - 1º Bimestre" criada com sucesso. ID: 12`

**Verificação (Opcional):**
Se houvesse um comando `vigenda avaliacao listar --classid 1 --term "1º Bimestre"`, você poderia usá-lo para ver as avaliações criadas.

## Passo 2: Lançar Notas para a Primeira Avaliação (P1)

Vamos supor que a P1 (ID 10) já ocorreu e você tem as notas.

Execute o comando para lançar notas:
```bash
./vigenda avaliacao lancar-notas 10
```

O Vigenda iniciará um prompt interativo. Para cada aluno da Turma ID 1, ele pedirá a nota:
```
Lançando notas para a Avaliação: P1 - História 9A (ID: 10)
Aluno: Alice Wonderland (ID: 101) - Nota: 7.5
Aluno: Bob The Builder (ID: 102) - Nota: 8.0
Aluno: Charles Xavier (ID: 103) - Nota: 6.0
... (continue para todos os alunos) ...
Notas lançadas com sucesso para a avaliação ID 10.
```
Digite a nota para cada aluno e pressione Enter. Siga as instruções do prompt para navegar, corrigir ou finalizar.

## Passo 3: Calcular a Média Parcial (Após P1)

Após lançar as notas da P1, você pode querer ver como estão as médias (ainda parciais, pois faltam T1 e Participação).

```bash
./vigenda avaliacao media-turma 1
```
*(Nota: A implementação CLI atual do comando `media-turma` não suporta a flag `--term`. A média será calculada com base em todas as avaliações com notas lançadas para a turma ID 1).*

A saída mostrará a média de cada aluno considerando apenas a P1 (se apenas ela tiver notas no momento).
Exemplo de saída:
```
Médias para Turma ID 1, Período: 1º Bimestre
-------------------------------------------
Aluno ID | Nome Aluno        | Média Parcial
-------------------------------------------
101      | Alice Wonderland  | 7.50
102      | Bob The Builder   | 8.00
103      | Charles Xavier    | 6.00
...
-------------------------------------------
* Cálculo baseado nas avaliações com notas lançadas no período.
```

## Passo 4: Lançar Notas para as Demais Avaliações (T1 e Participação)

Repita o processo do Passo 2 para as outras avaliações, usando seus respectivos IDs:

1.  **Lançar notas para T1 (ID 11):**
    ```bash
    ./vigenda avaliacao lancar-notas 11
    ```
    Insira as notas para o trabalho.

2.  **Lançar notas para Participacao (ID 12):**
    ```bash
    ./vigenda avaliacao lancar-notas 12
    ```
    Insira as notas de participação.

## Passo 5: Calcular a Média Final do Bimestre

Com todas as notas do "1º Bimestre" lançadas, calcule novamente a média da turma.

```bash
./vigenda avaliacao media-turma 1
```
*(Lembre-se: a flag `--term` não é usada pelo comando CLI atual. A média incluirá todas as avaliações com notas para a turma ID 1).*

Agora, a saída deverá refletir a média ponderada de todas as três avaliações (P1, T1, Participacao) que tiveram notas lançadas:
```
Médias para Turma ID 1, Período: 1º Bimestre
-------------------------------------------
Aluno ID | Nome Aluno        | Média Final Bimestre
-------------------------------------------
101      | Alice Wonderland  | X.XX  (Ex: (7.5*5 + T1_nota*3 + Part_nota*2) / (5+3+2) )
102      | Bob The Builder   | Y.YY
103      | Charles Xavier    | Z.ZZ
...
-------------------------------------------
* Cálculo baseado nas avaliações com notas lançadas no período.
```

**Entendendo o Cálculo:**
A média de cada aluno será calculada usando a fórmula da média ponderada:
`(Nota_P1 * Peso_P1 + Nota_T1 * Peso_T1 + Nota_Participacao * Peso_Participacao) / (Peso_P1 + Peso_T1 + Peso_Participacao)`
No nosso exemplo, o divisor da soma dos pesos é `5 + 3 + 2 = 10`.

## Dicas Adicionais

*   **Corrigindo Notas:** Se você precisar corrigir uma nota após o lançamento inicial, você pode geralmente executar o comando `vigenda avaliacao lancar-notas <ID_AVALIACAO>` novamente. O sistema pode permitir que você sobrescreva a nota anterior para alunos específicos. Verifique se há um comando `avaliacao atualizar-nota` para edições mais pontuais.
*   **Alunos sem Nota:** Observe como o sistema lida com alunos que possam não ter uma nota em alguma avaliação (ex: aluno faltou à P1). A média pode não ser calculada para esse aluno, ou ele pode receber uma nota zero, dependendo da política da escola e da implementação do Vigenda. Isso foi discutido na seção [Detalhes do Cálculo (`avaliacao media-turma`)](../../user_manual/README.md#detalhes-do-calculo-avaliacao-media-turma) do Manual do Usuário.
*   **Consistência nos Nomes/Termos:** Use nomes consistentes para os períodos (`--term`) para garantir que os cálculos de média agrupem as avaliações corretas.

## Conclusão

Este tutorial demonstrou como gerenciar múltiplas avaliações dentro de um período letivo, desde a criação até o lançamento de notas e o cálculo de médias ponderadas. Essa organização permite um acompanhamento claro do desempenho dos alunos.

Pratique com suas próprias turmas e avaliações para se familiarizar com o fluxo!
