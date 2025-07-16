# Guia de Contribuição

Agradecemos seu interesse em contribuir para o nosso projeto! Este guia detalha como você pode fazer contribuições significativas. Seguir estas diretrizes ajuda a manter a qualidade do código e a eficiência do processo de desenvolvimento.

## Como Começar

- **Leia o `README.md`:** Certifique-se de ter lido o arquivo `README.md` para entender o propósito do projeto e como configurá-lo.
- **Verifique as `Issues`:** Procure por issues abertas, especialmente aquelas marcadas como `good first issue` ou `help wanted`. Se você encontrar uma que gostaria de resolver, comente nela para que possamos atribuí-la a você e evitar trabalho duplicado.
- **Proponha Novas Ideias:** Se você tem uma ideia para uma nova funcionalidade ou melhoria, crie uma nova issue para discuti-la com a equipe antes de começar a trabalhar nela.

## Processo de Desenvolvimento

1.  **Fork o Repositório:**
    Crie um fork pessoal do repositório principal.

2.  **Clone o Fork Localmente:**
    ```bash
    git clone <URL_DO_SEU_FORK>
    cd <NOME_DO_REPOSITORIO>
    ```

3.  **Configure o Repositório Remoto (Upstream):**
    Adicione o repositório original como um remote chamado `upstream`.
    ```bash
    git remote add upstream <URL_DO_REPOSITORIO_ORIGINAL>
    ```
    Mantenha seu fork atualizado com o repositório principal:
    ```bash
    git fetch upstream
    git rebase upstream/main
    ```

4.  **Crie um Branch:**
    Crie um novo branch para suas alterações. Use um nome descritivo (ex: `feature/nova-funcionalidade` ou `fix/bug-reportado`).
    ```bash
    git checkout -b nome-do-seu-branch
    ```

5.  **Escreva o Código:**
    -   **Siga o Estilo de Código:** Adira às convenções de estilo de código do projeto.
    -   **Comente seu Código:** Adicione comentários claros e concisos.
    -   **Escreva Testes:** Adicione testes unitários e de integração para suas alterações.
    -   **Mantenha os Commits Atômicos:** Faça commits pequenos e lógicos.

6.  **Teste Suas Alterações:**
    Execute todos os testes para garantir que suas alterações não quebraram nada.
    ```bash
    go test ./...
    ```
    Consulte `TESTING.MD` para mais detalhes.

7.  **Faça o Commit das Suas Alterações:**
    Use mensagens de commit claras e descritivas, seguindo o padrão de [Conventional Commits](https://www.conventionalcommits.org/).
    ```bash
    git add .
    git commit -m "feat: adiciona nova funcionalidade X"
    ```

8.  **Envie Suas Alterações (Push):**
    ```bash
    git push origin nome-do-seu-branch
    ```

9.  **Crie um Pull Request (PR):**
    -   Abra um Pull Request do seu branch para o branch principal do repositório original.
    -   Forneça um título claro e uma descrição detalhada das suas alterações.
    -   Referencie qualquer issue relacionada (ex: "Closes #123").

10. **Revisão de Código:**
    -   Um ou mais mantenedores revisarão seu PR.
    -   Esteja preparado para responder a perguntas e fazer alterações com base no feedback.

## Diretrizes de Estilo de Código

-   **Formatação:** Use `gofmt` ou `goimports` para formatar seu código Go.
-   **Nomenclatura:** Siga as convenções de nomenclatura do Go.
-   **Comentários:** Comente todos os identificadores exportados.
-   **Tratamento de Erros:** Verifique e trate erros explicitamente.
-   **Linters:** Use `golangci-lint` para garantir a qualidade do código.

## Mensagens de Commit

Siga o padrão de [Conventional Commits](https://www.conventionalcommits.org/).

**Tipos Comuns:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`.

**Exemplo:**
```
feat(api): adiciona endpoint para buscar usuários por ID

Implementa a lógica para recuperar um usuário específico
do banco de dados usando seu identificador único.

Fixes #42
```

## Gestão de Issues

-   **Relatando Bugs:** Use o template de bug (se disponível) e forneça o máximo de detalhes possível.
-   **Sugerindo Melhorias:** Use o template de feature request (se disponível) e explique claramente o problema e a solução proposta.

## Conduta

Esperamos que todos os contribuidores sigam nosso [Código de Conduta](./CODE_OF_CONDUCT.md).

Obrigado por contribuir!
