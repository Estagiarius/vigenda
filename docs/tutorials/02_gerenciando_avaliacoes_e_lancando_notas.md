# Tutorial 02: Gerenciando Avaliações e Lançando Notas Detalhadamente

Este tutorial detalha como gerenciar o ciclo de vida das avaliações no Vigenda, desde a criação de múltiplas avaliações para uma turma, atribuindo pesos, lançando notas de forma interativa, até o cálculo da média da turma para verificar o impacto.

**Pré-requisitos:**
*   Vigenda instalado e configurado (consulte o [Guia de Introdução](../../docs/getting_started/README.md)).
*   Uma turma criada no Vigenda. Se você ainda não tem uma, pode consultar o [Manual do Usuário](../../user_manual/README.md#criar-turma-vigenda-turma-criar) ou o tutorial sobre criação de turmas (se disponível). Para este tutorial, vamos assumir que existe uma turma com `ID = 1` (ex: "História 9A") e que ela já tem alunos importados.
*   Familiaridade com os seguintes comandos do Vigenda (consulte o Manual do Usuário para detalhes):
    *   `vigenda avaliacao criar`
    *   `vigenda avaliacao listar` (Importante para obter IDs de avaliação)
    *   `vigenda avaliacao lancar-notas`
    *   `vigenda avaliacao media-turma`

## Cenário Proposto

Vamos gerenciar as avaliações da turma "História 9A" (`ID = 1`) para o período letivo "1º Bimestre". Planejamos três instrumentos avaliativos com pesos diferentes:
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
    Anote o ID. Ex: `Avaliação "Participacao - 1º Bimestre" (ID: 12) criada com sucesso para a turma ID 1.`

**Verificação:**
Após criar as avaliações, é uma boa prática listá-las para confirmar e obter os IDs corretos para os próximos passos.
```bash
./vigenda avaliacao listar --classid 1 --term "1º Bimestre"
```
A saída deve ser algo como:
```
ID | Nome da Avaliação             | Turma ID | Período      | Peso | Data
---|-------------------------------|----------|--------------|------|-----------
10 | P1 - História 9A              | 1        | 1º Bimestre  | 5.0  | 2024-03-15
11 | T1 - Trabalho: Egito Antigo   | 1        | 1º Bimestre  | 3.0  | 2024-03-29
12 | Participacao - 1º Bimestre    | 1        | 1º Bimestre  | 2.0  |
```
Guarde bem esses IDs (10, 11, 12 no nosso exemplo).

## Passo 2: Lançar Notas para a Primeira Avaliação (P1)

Vamos supor que a P1 (ID 10 em nosso exemplo) já ocorreu e você tem as notas dos alunos.

Execute o comando para lançar as notas:
```bash
./vigenda avaliacao lancar-notas 10
```

O Vigenda iniciará um prompt interativo. Para cada aluno da Turma ID 1, ele solicitará a nota:
```
Lançando notas para a Avaliação: "P1 - História 9A" (ID: 10)
Turma: História 9A (ID: 1)
--------------------------------------------------
Aluno: Alice Wonderland (ID: 101)
Nota (deixe em branco para pular, 's' para sair e salvar, 'c' para cancelar): 7.5
--------------------------------------------------
Aluno: Bob The Builder (ID: 102)
Nota (deixe em branco para pular, 's' para sair e salvar, 'c' para cancelar): 8.0
--------------------------------------------------
Aluno: Charles Xavier (ID: 103)
Nota (deixe em branco para pular, 's' para sair e salvar, 'c' para cancelar): 6.0
... (continue para todos os alunos) ...
Notas lançadas com sucesso para a avaliação ID 10.
```
Digite a nota para cada aluno e pressione Enter.
*   **Pular Aluno:** Deixe o campo de nota em branco e pressione Enter.
*   **Sair e Salvar:** Digite `s` e pressione Enter para salvar o progresso e sair.
*   **Cancelar:** Digite `c` e pressione Enter para sair sem salvar as alterações desta sessão de lançamento.

## Passo 3: Calcular a Média Parcial (Após Lançar Notas da P1)

Após lançar as notas da P1, você pode querer ver como estão as médias (ainda parciais, pois faltam T1 e Participação).

```bash
./vigenda avaliacao media-turma --classid 1 --term "1º Bimestre"
```

A saída mostrará a média de cada aluno considerando apenas a P1, pois é a única com notas lançadas no período.
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
./vigenda avaliacao media-turma --classid 1 --term "1º Bimestre"
```

Agora, a saída deverá refletir a média ponderada de todas as três avaliações (P1, T1, Participacao):
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
No nosso exemplo, o divisor da soma dos pesos é `5.0 + 3.0 + 2.0 = 10.0`.

## Dicas Adicionais e Boas Práticas

*   **Corrigindo Notas Lançadas:**
    *   Se você errar uma nota durante o lançamento interativo e perceber antes de finalizar para aquele aluno, algumas interfaces podem permitir apagar e redigitar.
    *   Se precisar corrigir uma nota após a sessão de lançamento, a forma mais comum é executar o comando `vigenda avaliacao lancar-notas <ID_DA_AVALIACAO>` novamente. Ao chegar no aluno desejado, insira a nova nota. O sistema geralmente sobrescreve a anterior.
    *   Para edições muito pontuais, verifique no Manual do Usuário se existe um comando como `vigenda avaliacao editar-nota --alunoid <ID_ALUNO> --avaliacaoid <ID_AVALIACAO> --nota <NOVA_NOTA>`.

*   **Alunos sem Nota (Ausentes ou Pendentes):**
    *   O comportamento padrão do Vigenda ao calcular médias com notas ausentes é crucial. Idealmente, alunos com notas pendentes em avaliações obrigatórias (peso > 0) não teriam sua média calculada para o período, ou a média seria exibida como "Pendente".
    *   **Teste este cenário:** Crie uma avaliação, lance nota para alguns alunos mas não para todos, e então calcule a média da turma para observar o comportamento. Consulte a seção [Detalhes do Cálculo (`avaliacao media-turma`)](../../user_manual/README.md#detalhes-do-calculo-avaliacao-media-turma) no Manual do Usuário para a política oficial do Vigenda sobre isso.

*   **Consistência nos Nomes dos Períodos (`--term`):**
    Use exatamente os mesmos nomes para os períodos avaliativos (ex: "1º Bimestre", "2º Trimestre") ao criar avaliações e ao calcular médias. Qualquer variação (ex: "Bimestre 1" vs "1º Bimestre") fará com que as avaliações não sejam agrupadas corretamente.

*   **Backup dos Dados:** Lembre-se da importância de fazer backups regulares do seu banco de dados Vigenda, especialmente se estiver usando SQLite.

## Conclusão

Este tutorial demonstrou o fluxo completo de gerenciamento de múltiplas avaliações para uma turma dentro de um período letivo: criação com pesos, listagem para verificação, lançamento interativo de notas e cálculo de médias ponderadas. Essa organização sistemática é fundamental para um acompanhamento preciso e eficiente do desempenho dos alunos ao longo do tempo.

Experimente com suas próprias turmas e diferentes cenários de avaliação para se tornar proficiente no uso dessas funcionalidades do Vigenda!
