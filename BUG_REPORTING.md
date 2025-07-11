# Reportando Bugs no Vigenda

Agradecemos por ajudar a melhorar o Vigenda! Se você encontrar um bug, esta página descreve a melhor forma de reportá-lo.

## Antes de Reportar

1.  **Verifique a Documentação:** Certifique-se de que o comportamento que você está observando não é o esperado, consultando o [Manual do Usuário](./docs/user_manual/README.md) e o [FAQ](./docs/faq/README.md).
2.  **Procure por Issues Existentes:** Verifique o rastreador de issues do projeto (se disponível publicamente) para ver se o bug já foi reportado. Se sim, você pode adicionar comentários ou reações à issue existente.

## Como Reportar um Bug

Se o bug ainda não foi reportado, por favor, crie uma nova "Issue" no rastreador de issues do projeto. Se o projeto não tiver um rastreador público ou se você preferir, entre em contato conforme as instruções no `README.md` principal.

Ao reportar um bug, por favor, inclua o máximo de detalhes possível:

1.  **Título Claro e Descritivo:** Um bom título ajuda a entender rapidamente a natureza do problema.
    *   Exemplo Ruim: "Programa quebrou"
    *   Exemplo Bom: "Comando `tarefa listar --classid X` falha se a turma não tiver tarefas"

2.  **Versão do Vigenda:** Se possível, especifique a versão do Vigenda que você está usando (ex: commit hash se compilado da fonte, ou número da versão se for um release).

3.  **Ambiente:**
    *   Sistema Operacional e versão (ex: Ubuntu 22.04, Windows 10 Pro 22H2, macOS Sonoma 14.1).
    *   Versão do Go (se compilando da fonte).
    *   Tipo de terminal usado (se relevante para bugs na TUI).

4.  **Passos para Reproduzir o Bug:**
    Detalhe, passo a passo, as ações que levam ao bug. Quanto mais preciso, melhor.
    *   Exemplo:
        1.  Executei `vigenda turma importar-alunos 1 alunos.csv`.
        2.  O arquivo `alunos.csv` continha a seguinte linha mal formatada: `,,,`
        3.  A aplicação encerrou com um panic em vez de mostrar uma mensagem de erro amigável.

5.  **Comportamento Esperado:**
    O que você esperava que acontecesse?
    *   Exemplo: "Esperava que a aplicação mostrasse uma mensagem de erro indicando que a linha X no CSV está mal formatada e continuasse o processamento das outras linhas ou parasse graciosamente."

6.  **Comportamento Atual (Observado):**
    O que realmente aconteceu?
    *   Exemplo: "A aplicação travou e exibiu a seguinte mensagem de erro/stack trace: [copie e cole o erro aqui]."
    *   Se for um bug visual na TUI, descreva o que você vê. Screenshots (se possível anexar na issue) podem ser muito úteis.

7.  **Logs (Se Aplicável):**
    Se houver logs relevantes no arquivo `vigenda.log` (localizado geralmente em `~/.config/vigenda/` ou no diretório atual), por favor, anexe a seção pertinente.

## O que Acontece Depois?

Após reportar um bug:
*   Os mantenedores do projeto revisarão a issue.
*   Eles podem fazer perguntas adicionais para entender melhor o problema.
*   A issue será priorizada e, quando possível, uma correção será implementada.
*   Você será notificado sobre o progresso através da issue.

Obrigado por sua contribuição para tornar o Vigenda uma ferramenta melhor!
