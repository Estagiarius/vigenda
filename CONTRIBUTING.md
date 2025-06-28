# Guia de Contribuição

Agradecemos seu interesse em contribuir para o nosso projeto! Este guia detalha como você pode fazer contribuições significativas. Seguir estas diretrizes ajuda a manter a qualidade do código e a eficiência do processo de desenvolvimento.

## Como Começar

- **Leia o `README.md`:** Certifique-se de ter lido o arquivo `README.md` para entender o propósito do projeto e como configurá-lo.
- **Verifique as `Issues`:** Procure por issues abertas, especialmente aquelas marcadas como `good first issue` ou `help wanted`. Se você encontrar uma que gostaria de resolver, comente nela para que possamos atribuí-la a você e evitar trabalho duplicado.
- **Proponha Novas Ideias:** Se você tem uma ideia para uma nova funcionalidade ou melhoria, crie uma nova issue para discuti-la com a equipe antes de começar a trabalhar nela.

## Processo de Desenvolvimento

1.  **Fork o Repositório:**
    Crie um fork pessoal do repositório principal na sua conta do GitHub (ou plataforma similar).

2.  **Clone o Fork Localmente:**
    ```bash
    git clone https://github.com/SEU_USUARIO/NOME_DO_REPOSITORIO.git
    cd NOME_DO_REPOSITORIO
    ```

3.  **Configure o Repositório Remoto (Upstream):**
    Adicione o repositório original como um remote chamado `upstream`.
    ```bash
    git remote add upstream https://github.com/ORGANIZACAO_OU_USUARIO_ORIGINAL/NOME_DO_REPOSITORIO.git
    ```
    Mantenha seu fork atualizado com o repositório principal:
    ```bash
    git fetch upstream
    git rebase upstream/main # ou a branch principal do projeto
    ```

4.  **Crie um Branch:**
    Crie um novo branch para suas alterações. Use um nome descritivo para o branch (por exemplo, `feature/nova-funcionalidade` ou `fix/bug-reportado`).
    ```bash
    git checkout -b nome-do-seu-branch
    ```

5.  **Escreva o Código:**
    -   **Siga o Estilo de Código:** Adira às convenções de estilo de código do projeto (detalhadas abaixo).
    -   **Comente seu Código:** Adicione comentários claros e concisos, especialmente em partes complexas da lógica. Explique o "porquê" do código, não apenas o "o quê".
        -   **Comentários de Função/Método:** Descreva o propósito da função, seus parâmetros e o que ela retorna.
        -   **Comentários Inline:** Use para explicar seções específicas ou decisões de design dentro do código.
    -   **Escreva Testes:** Adicione testes unitários e de integração para suas alterações. Certifique-se de que todos os testes existentes e os novos passam.
    -   **Mantenha os Commits Atômicos:** Faça commits pequenos e lógicos. Cada commit deve representar uma unidade de trabalho coesa.

6.  **Teste Suas Alterações:**
    Execute todos os testes para garantir que suas alterações não quebraram nada.
    ```bash
    # Exemplo de comando para rodar testes (ajuste conforme o projeto)
    go test ./...
    # ou
    # npm test
    ```
    Consulte `TESTING.md` para mais detalhes sobre como executar os testes.

7.  **Faça o Commit das Suas Alterações:**
    Use mensagens de commit claras e descritivas. Siga o formato convencional (veja a seção "Mensagens de Commit").
    ```bash
    git add .
    git commit -m "feat: Adiciona nova funcionalidade X"
    ```

8.  **Envie Suas Alterações (Push):**
    Envie suas alterações para o seu fork no GitHub.
    ```bash
    git push origin nome-do-seu-branch
    ```

9.  **Crie um Pull Request (PR):**
    -   Abra um Pull Request do seu branch no seu fork para o branch principal (geralmente `main` ou `master`) do repositório original.
    -   Forneça um título claro e uma descrição detalhada das suas alterações no PR.
    -   Referencie qualquer issue relacionada (por exemplo, "Closes #123").
    -   Se o seu PR for um trabalho em andamento, marque-o como "Draft" ou adicione "[WIP]" ao título.

10. **Revisão de Código:**
    -   Um ou mais mantenedores revisarão seu PR.
    -   Esteja preparado para responder a perguntas e fazer alterações com base no feedback.
    -   Após a aprovação e a passagem de quaisquer verificações de CI, seu PR será mesclado.

## Diretrizes de Estilo de Código

[Esta seção deve ser adaptada para as linguagens e ferramentas específicas do projeto.]

### Geral

-   **Consistência:** Mantenha a consistência com o estilo de código existente no projeto.
-   **Clareza:** Escreva código que seja fácil de ler e entender.
-   **Simplicidade:** Prefira soluções simples e diretas sempre que possível.

### [Linguagem Específica, ex: Go]

-   **Formatação:** Use `gofmt` ou `goimports` para formatar seu código Go.
-   **Nomenclatura:** Siga as convenções de nomenclatura do Go (por exemplo, `camelCase` para variáveis locais, `PascalCase` para identificadores exportados).
-   **Comentários:**
    -   Comente todos os identificadores exportados (funções, tipos, constantes, variáveis).
    -   Use `//` para comentários de linha e `/* */` para comentários de bloco quando apropriado.
-   **Tratamento de Erros:** Verifique e trate erros explicitamente. Evite ignorar erros.
-   **Linters:** [Se aplicável, mencione linters como `golangci-lint` e como executá-los.]

### [Linguagem Específica, ex: JavaScript/TypeScript]

-   **Formatação:** Use Prettier ou ESLint com uma configuração de formatação.
-   **Nomenclatura:** [Especifique convenções, por exemplo, `camelCase` para variáveis e funções, `PascalCase` para classes e componentes.]
-   **Linters:** Use ESLint com um conjunto de regras definido (por exemplo, Airbnb, StandardJS).
-   **Comentários:** Use JSDoc para documentar funções, classes e módulos.

## Mensagens de Commit

Siga o padrão de [Conventional Commits](https://www.conventionalcommits.org/):

```
<tipo>[escopo opcional]: <descrição concisa em letras minúsculas>

[corpo opcional do commit com mais detalhes]

[rodapé opcional com BREAKING CHANGE ou referências a issues]
```

**Tipos Comuns:**

-   `feat`: Uma nova funcionalidade.
-   `fix`: Uma correção de bug.
-   `docs`: Alterações apenas na documentação.
-   `style`: Alterações que não afetam o significado do código (espaçamento, formatação, ponto e vírgula ausente, etc.).
-   `refactor`: Uma alteração de código que não corrige um bug nem adiciona uma funcionalidade.
-   `perf`: Uma alteração de código que melhora o desempenho.
-   `test`: Adicionando testes ausentes ou corrigindo testes existentes.
-   `build`: Alterações que afetam o sistema de build ou dependências externas.
-   `ci`: Alterações nos nossos arquivos e scripts de configuração de CI.
-   `chore`: Outras alterações que não modificam arquivos `src` ou de teste (ex: atualização de dependências).

**Exemplo:**

```
feat(api): adiciona endpoint para buscar usuários por ID

Implementa a lógica para recuperar um usuário específico
do banco de dados usando seu identificador único.

Fixes #42
```

## Gestão de Issues

-   **Relatando Bugs:**
    -   Use o template de bug (se disponível) ao criar uma issue.
    -   Forneça o máximo de detalhes possível: passos para reproduzir, comportamento esperado, comportamento atual, ambiente (versão do SO, versão do software, etc.).
-   **Sugerindo Melhorias:**
    -   Use o template de feature request (se disponível).
    -   Explique claramente o problema que a melhoria resolveria e a solução proposta.

## Conduta

Esperamos que todos os contribuidores sigam nosso [Código de Conduta](CODE_OF_CONDUCT.md) (se existir, caso contrário, adicione uma breve seção sobre comportamento respeitoso).

Obrigado por contribuir! Sua ajuda é muito valiosa.
