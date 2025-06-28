# Tutorial 01: Criando sua Primeira Turma e Importando Alunos

Este tutorial guiará você pelos passos para criar sua primeira turma no Vigenda e, em seguida, importar uma lista de alunos para essa turma usando um arquivo CSV.

**Pré-requisitos:**
*   Vigenda instalado e funcionando. Consulte o [Guia de Introdução](../../getting_started/README.md) se precisar de ajuda.
*   Acesso ao terminal ou prompt de comando.
*   (Opcional) Um editor de texto simples para criar o arquivo CSV (como Bloco de Notas, VS Code, gedit, etc.).

## Passo 1: Entendendo as Disciplinas (Suposição)

Antes de criar uma turma, o Vigenda geralmente associa turmas a "Disciplinas" (ex: Matemática, História, Ciências). Para este tutorial, vamos **supor** que já existe uma disciplina cadastrada no sistema com `ID = 1` e nome "História Geral".

*Se o seu sistema Vigenda requer que você crie disciplinas primeiro, você precisaria de um comando como `vigenda disciplina criar "História Geral"` ou similar. Verifique a documentação do seu Vigenda específico ou use um ID de disciplina que você sabe que existe.*

## Passo 2: Criar uma Nova Turma

Vamos criar uma turma chamada "História 9º Ano A" para o ano de 2024, associada à disciplina "História Geral" (que estamos supondo ter ID 1).

Abra seu terminal e execute o seguinte comando (ajuste `./vigenda` se necessário, dependendo de onde seu executável está e se ele está no PATH):

```bash
./vigenda turma criar "História 9º Ano A" --subjectid 1 --year 2024
```

**O que acontece:**
*   `./vigenda turma criar`: Chama o comando para criar uma turma.
*   `"História 9º Ano A"`: É o nome que você está dando para a sua turma.
*   `--subjectid 1`: Associa esta turma à disciplina com ID 1 (nossa "História Geral" suposta).
*   `--year 2024`: (Opcional, mas bom para organização) Define o ano letivo da turma.

**Saída Esperada:**
O Vigenda deve responder com uma mensagem indicando que a turma foi criada com sucesso e, crucialmente, **informará o ID da nova turma**. Algo como:
```
Turma "História 9º Ano A" criada com sucesso. ID da Turma: 5
```
**Anote este ID da Turma (por exemplo, `5`).** Você precisará dele no próximo passo. Vamos usar `5` como exemplo daqui para frente. Se o seu ID for diferente, use o seu.

## Passo 3: Preparar o Arquivo CSV de Alunos

Agora, vamos criar um arquivo CSV (Comma Separated Values - Valores Separados por Vírgula) com a lista de alunos para importar.

Crie um novo arquivo de texto e nomeie-o, por exemplo, `alunos_historia_9a.csv`.
O conteúdo do arquivo deve seguir este formato:

```csv
numero_chamada,nome_completo,situacao
1,"Alice Wonderland","ativo"
2,"Bob The Builder","ativo"
3,"Charles Xavier","inativo"
,"Diana Prince","ativo"
5,"Edward Scissorhands","transferido"
```

**Detalhes das colunas:**
*   `numero_chamada` (opcional): O número de chamada do aluno. Pode ser deixado em branco (como para "Diana Prince").
*   `nome_completo` (obrigatório): O nome completo do aluno. Se o nome contiver vírgulas, coloque-o entre aspas duplas (embora neste exemplo simples não seja necessário, é uma boa prática).
*   `situacao` (opcional): O status do aluno. Valores válidos são `ativo`, `inativo`, `transferido`. Se omitido ou deixado em branco, o padrão geralmente é `ativo`.

Salve este arquivo em um local que você possa acessar facilmente pelo terminal (por exemplo, no mesmo diretório onde está o executável `vigenda` ou em um subdiretório `dados/`).

## Passo 4: Importar os Alunos para a Turma

Com o arquivo CSV pronto e o ID da turma em mãos (lembre-se, estamos usando `5` como exemplo), execute o comando de importação:

```bash
./vigenda turma importar-alunos 5 alunos_historia_9a.csv
```

**O que acontece:**
*   `./vigenda turma importar-alunos`: Chama o comando para importar alunos.
*   `5`: É o ID da turma que você anotou no Passo 2 (substitua pelo seu ID real).
*   `alunos_historia_9a.csv`: É o nome do arquivo CSV que você criou no Passo 3. Se o arquivo estiver em outro diretório, forneça o caminho completo (ex: `dados/alunos_historia_9a.csv`).

**Saída Esperada:**
O Vigenda deve processar o arquivo e informar quantos alunos foram importados, ou se houve algum erro. Uma mensagem de sucesso pode ser:
```
Importação de alunos para a turma ID 5 concluída.
Alunos processados: 5
Alunos importados com sucesso: 5
Alunos com erro: 0
```

## Passo 5: Verificar (Opcional)

Atualmente, o `AGENTS.md` não especifica um comando direto para listar alunos de uma turma. No entanto, se tal comando existir (ex: `vigenda turma listar-alunos --classid 5`), você poderia usá-lo aqui para verificar se os alunos foram adicionados corretamente.

Outra forma de verificar, indiretamente, seria ao tentar lançar notas para uma avaliação nesta turma; os nomes dos alunos importados deveriam aparecer.

## Conclusão

Parabéns! Você criou sua primeira turma e importou alunos para ela usando o Vigenda.

**Resumo dos comandos utilizados:**
1.  `./vigenda turma criar "Nome da Turma" --subjectid <ID_DISCIPLINA> [--year <ANO>]`
2.  `./vigenda turma importar-alunos <ID_DA_TURMA_CRIADA> <NOME_DO_ARQUIVO.csv>`

Agora você está pronto para adicionar tarefas, criar avaliações e gerenciar suas atividades pedagógicas para esta turma!

---
Próximos passos sugeridos:
*   Explore como adicionar tarefas para esta turma no [Manual do Usuário](../../user_manual/README.md#gestao-de-tarefas).
*   Consulte outros [Tutoriais](./README.md).
