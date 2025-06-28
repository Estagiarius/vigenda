# Documentação para Agentes de IA

Olá, Agente! Este documento fornece diretrizes e informações para ajudá-lo a entender e trabalhar eficientemente com este codebase.

## 1. Visão Geral do Projeto

-   **Propósito Principal:** [Descreva resumidamente o que o projeto faz. Ex: "Este é um serviço de API REST para gerenciar tarefas de usuários."]
-   **Tecnologias Chave:**
    -   Linguagem Principal: Go (instalado via `apt-get install golang-go`)
    -   Compilador C: GCC (instalado via `apt-get install gcc`)
    -   Framework Principal: [Ainda não definido, a ser preenchido conforme o projeto evolui]
    -   Banco de Dados: [Ainda não definido, a ser preenchido conforme o projeto evolui]
    -   Outras tecnologias importantes: [Ainda não definido, a ser preenchido conforme o projeto evolui]
-   **Estrutura do Repositório (Documentação Inicial):**
    -   `README.md`: Contém a visão geral do projeto, instruções de instalação e como contribuir. **Leia-o primeiro.**
    -   `API_DOCUMENTATION.md`: Modelo para documentação de API. Preencha se o projeto expuser uma API.
    -   `TECHNICAL_SPECIFICATION.MD`: Modelo para especificações técnicas. Detalhe as escolhas de arquitetura e design aqui.
    -   `INSTALLATION.MD`: Instruções detalhadas para configurar o ambiente de desenvolvimento.
    -   `TESTING.MD`: Como executar e escrever testes. Adapte para as ferramentas de teste escolhidas.
    -   `CONTRIBUTING.MD`: Diretrizes para contribuição humana e de IA, incluindo padrões de codificação e processo de PR. **Siga estas diretrizes rigorosamente.**
    -   `CHANGELOG.MD`: Registro de mudanças nas versões. Mantenha-o atualizado.
    -   `AGENTS.md`: Este arquivo. Consulte-o para diretrizes específicas para IA.
    -   `cmd/`: Diretório comum para aplicações principais (entrypoints) em projetos Go.
    -   `internal/` ou `pkg/`: Diretórios comuns para código fonte principal em projetos Go.
        -   `internal/`: Código que não deve ser importado por outros projetos.
        -   `pkg/`: Código que pode ser importado por outros projetos (se houver bibliotecas reutilizáveis).
    -   Outros diretórios (`api/`, `handlers/`, `models/`, `services/`, `store/`, `configs/`, `scripts/`, `tests/`, `web/`, `ui/`, `frontend/`, `docs/`) devem ser criados conforme necessário e sua estrutura explicada aqui e na `TECHNICAL_SPECIFICATION.MD`.
-   **Documentação Adicional:**
    -   Sempre consulte o `README.md`, `INSTALLATION.MD`, `CONTRIBUTING.MD` e `TESTING.MD` antes de iniciar qualquer tarefa.
    -   A `TECHNICAL_SPECIFICATION.MD` deve ser atualizada com decisões de design à medida que são tomadas.
    -   Se houver uma API, a `API_DOCUMENTATION.MD` é crucial.

## 2. Configuração do Ambiente

-   Siga **rigorosamente** as instruções em `INSTALLATION.MD` para configurar seu ambiente.
-   As dependências principais já instaladas são Go e GCC.
-   Preste atenção especial à configuração de variáveis de ambiente e arquivos de configuração (ex: `.env`, `config.json`) quando eles forem introduzidos no projeto.

## 3. Tarefas Comuns e Como Abordá-las

### 3.1. Entendendo o Código Existente (Quando Houver)
-   Comece pela função ou módulo principal relacionado à sua tarefa.
-   Use as ferramentas de busca (`grep`) para encontrar definições de funções, tipos e usos.
-   Leia os comentários no código. Se não houver, adicione comentários claros conforme as diretrizes em `CONTRIBUTING.MD`.
-   Analise os testes relacionados ao código que você está investigando. Eles demonstram como o código deve ser usado e qual é o comportamento esperado.

### 3.2. Implementando Novas Funcionalidades
1.  **Planeje:**
    -   Certifique-se de que a funcionalidade está bem definida.
    -   Identifique os arquivos e módulos que precisarão ser modificados ou criados.
    -   Considere como a nova funcionalidade se encaixa na arquitetura existente (consulte e atualize `TECHNICAL_SPECIFICATION.MD`).
2.  **Escreva Testes Primeiro (TDD/BDD quando possível):**
    -   Consulte `TESTING.MD` para tipos de testes e ferramentas.
    -   Escreva testes unitários para a nova lógica.
    -   Se a funcionalidade envolver um endpoint de API, escreva testes de integração para esse endpoint.
3.  **Implemente o Código:**
    -   Siga as convenções de estilo de código e padrões de design descritos em `CONTRIBUTING.MD` e `TECHNICAL_SPECIFICATION.MD`.
    -   **Comente seu código:** Explique a lógica complexa, as decisões de design e o propósito das funções e estruturas públicas. Use o formato de comentário especificado (ex: Godoc para Go: `// MinhaFuncao faz X e Y.`).
    -   **Tratamento de Erros:** Implemente tratamento de erros robusto. Retorne erros apropriados e forneça contexto.
    -   **Logging:** Adicione logs úteis para depuração e monitoramento. Use o logger configurado no projeto (a ser definido).
4.  **Execute os Testes:**
    -   Execute todos os testes (unitários, integração) para garantir que suas alterações não quebraram nada e que seus novos testes passam.
    -   Verifique a cobertura de teste, se configurada.

### 3.3. Corrigindo Bugs
1.  **Reproduza o Bug:** Entenda claramente como o bug ocorre.
2.  **Escreva um Teste que Falhe:** Crie um teste que reproduza o bug. Este teste deve falhar com o código atual.
3.  **Corrija o Código:** Implemente a correção para o bug.
4.  **Execute os Testes:** Verifique se o novo teste (e todos os outros) agora passa.

### 3.4. Adicionando Comentários no Código
-   **Objetivo:** Tornar o código mais fácil de entender para outros desenvolvedores (humanos e IA).
-   **Onde Comentar:**
    -   **Funções/Métodos Públicos (Exportados):** Descreva o que a função faz, seus parâmetros, o que ela retorna e quaisquer efeitos colaterais importantes. (Ex: formato Godoc para Go: `// MinhaFuncao faz X e Y.`)
    -   **Estruturas de Dados Públicas (Exportadas):** Descreva o propósito da estrutura e seus campos.
    -   **Blocos de Lógica Complexa:** Adicione comentários inline para explicar partes do código que não são imediatamente óbvias. Explique o "porquê" de uma determinada abordagem, se não for claro.
    -   **Decisões de Design Importantes:** Se uma escolha específica de implementação foi feita por uma razão particular (performance, evitar um problema conhecido), documente isso.
-   **O que NÃO Comentar (Geralmente):**
    -   Código óbvio (ex: `// incrementa i`).
    -   Comentários que apenas repetem o que o código faz (ex: `// atribui x a y`).
-   **Estilo:** Siga o estilo de comentário da linguagem e do projeto (ver `CONTRIBUTING.MD`). Mantenha os comentários atualizados com as mudanças no código.

### 3.5. Refatorando Código
-   Certifique-se de que há testes adequados cobrindo o código a ser refatorado.
-   Faça pequenas alterações incrementais e execute os testes frequentemente.
-   O objetivo da refatoração é melhorar a clareza, manutenibilidade ou desempenho sem alterar o comportamento externo.

## 4. Padrões e Convenções Específicas do Projeto

-   **Estilo de Código:** Siga rigorosamente as diretrizes em `CONTRIBUTING.MD`. Use as ferramentas de formatação e linting especificadas (`gofmt`/`goimports`, `golangci-lint` para Go, a serem configuradas).
    -   Execute os linters/formatadores antes de submeter o código.
-   **Mensagens de Commit:** Use o formato de [Conventional Commits](https://www.conventionalcommits.org/) conforme detalhado em `CONTRIBUTING.MD`.
    -   Ex: `feat: adiciona endpoint de criação de usuário`
    -   Ex: `fix: corrige cálculo de total em pedidos`
    -   Ex: `docs: atualiza documentação da API para endpoint X`
-   **Branching:** Crie branches descritivas para suas tarefas (ex: `feature/TASK-ID-descricao` ou `fix/TASK-ID-descricao`).
-   **Pull Requests (PRs):**
    -   Forneça descrições claras e concisas das alterações.
    -   Referencie as issues relacionadas.
    -   Certifique-se de que todos os testes e verificações de CI passam (CI a ser configurado).

## 5. Ferramentas e Comandos Úteis

-   **Listar arquivos:** `ls -la [diretório]` (ou `ls` se o diretório for o atual)
-   **Ler arquivos:** `cat [arquivo]` ou use a ferramenta `read_files(["caminho/do/arquivo"])`
-   **Buscar texto em arquivos (grep):** `grep -R "termo_de_busca" [diretório_ou_arquivo]`
    -   Ex: `grep -R "CreateUser" internal/`
-   **Executar testes:** Veja `TESTING.MD`. Comandos comuns a serem definidos:
    -   Go: `go test ./...`, `go test -run NomeDoTesteEspecífico`, `go test -coverprofile=c.out ./... && go tool cover -html=c.out`
-   **Executar linters:** Veja `TESTING.MD`.
    -   Go: `golangci-lint run ./...` (após configuração)
-   **Compilar/Executar a aplicação:** Veja `INSTALLATION.MD` ou `README.md`.
    -   Go: `go run cmd/server/main.go` (ou o entrypoint relevante, a ser criado)
    -   Go (build): `go build -o minhaapp ./cmd/server/` (a ser criado)

## 6. O que Evitar

-   **Alterar arquivos fora do escopo da tarefa sem uma boa razão e sem documentar.**
-   **Introduzir dependências desnecessárias sem discuti-las ou documentá-las na `TECHNICAL_SPECIFICATION.MD`.**
-   **Comentar código funcional em vez de removê-lo (se não for mais necessário).**
-   **Ignorar falhas de teste ou de linter.**
-   **Submeter código sem testá-lo adequadamente conforme `TESTING.MD`.**
-   **Escrever mensagens de commit vagas como "correções" ou "atualizações". Siga `CONTRIBUTING.MD`.**

## 7. Se Você Ficar Preso

1.  **Releia a Tarefa e os Documentos:** Certifique-se de que entendeu completamente o requisito e consultou `README.md`, `INSTALLATION.MD`, `CONTRIBUTING.MD`, `TESTING.MD` e este `AGENTS.md`.
2.  **Pesquise Erros:** Copie e cole mensagens de erro em um mecanismo de busca ou na sua base de conhecimento.
3.  **Simplifique o Problema:** Tente isolar a parte do código que está causando o problema. Crie um caso de teste mínimo, se aplicável.
4.  **Consulte a Documentação da Linguagem/Framework:** Verifique a documentação oficial do Go (golang.org) ou de quaisquer bibliotecas/frameworks que venham a ser utilizados.
5.  **Peça Ajuda (se aplicável e configurado):** Use a ferramenta `request_user_input` descrevendo o problema, o que você já tentou (referenciando os documentos consultados) e onde está preso.

---
Lembre-se, seu objetivo é produzir código de alta qualidade e documentação que seja fácil de manter e entender. Siga estas diretrizes e use os documentos fornecidos para guiá-lo. Atualize esta documentação e as outras conforme o projeto evolui. Boa codificação!
