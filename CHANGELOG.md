# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Versionamento Semântico](https://semver.org/spec/v2.0.0.html).

## [Não Lançado]

### Adicionado
-   Estrutura inicial da documentação do desenvolvedor (TASK-DOC-02):
    -   `README.md`: Visão geral do projeto, instalação e contribuição.
    -   `API_DOCUMENTATION.md`: Modelo para documentação de API (marcado como não aplicável para CLI).
    -   `TECHNICAL_SPECIFICATION.MD`: Modelo para especificações técnicas.
    -   `CONTRIBUTING.md`: Guia detalhado para contribuições.
    -   `CHANGELOG.md`: Este arquivo, para rastrear mudanças.
    -   `INSTALLATION.MD`: Modelo para documentação de instalação.
    -   `TESTING.MD`: Modelo para documentação de testes.
    -   `AGENTS.md`: Modelo para documentação específica para agentes de IA.
-   Configuração inicial do ambiente de desenvolvimento (instalação de Go e GCC).

### Alterado
-   **Aprofundamento da Documentação do Desenvolvedor:**
    -   `README.md`: Preenchido com detalhes específicos do projeto Vigenda, incluindo descrição, pré-requisitos (Go 1.23.0, GCC), instruções de instalação (`go mod tidy`), e execução (`go run`, `build.sh`).
    -   `TECHNICAL_SPECIFICATION.MD`: Detalhada a arquitetura em camadas do Vigenda, componentes principais (pacotes Go), fluxo de dados com diagrama textual (Mermaid), escolhas tecnológicas (Go, Cobra, Bubbletea, SQLite), padrões de design, considerações de segurança, escalabilidade, desempenho e estratégia de testes.
    -   `INSTALLATION.MD`: Especificados os sistemas operacionais suportados, versões de Go e GCC, instruções para MinGW (cross-compilação Windows), `go mod tidy`, instalação de `goimports` e `golangci-lint`. Detalhados os comandos de build (`./build.sh`, `go build`) e execução.
    -   `TESTING.MD`: Detalhados os tipos de testes (Unitários, Integração, E2E, Performance), ferramentas (pacote `testing`, Testify, `os/exec`, Golden Files, `golangci-lint`, `goimports`), e comandos específicos para executar todos os testes, testes de pacotes, testes por nome, benchmarks e verificar cobertura de teste.
    -   `AGENTS.md`: Adaptado significativamente para o projeto Vigenda, detalhando a visão geral, tecnologias, estrutura de diretórios Go, configuração de ambiente, e comandos úteis específicos (`go test`, `go run`, `go build`, `golangci-lint run`, `goimports -w`).

### Corrigido
-   (Nenhum bug corrigido ainda)

### Removido
-   (Nada removido ainda)

### Segurança
-   (Nenhuma atualização de segurança ainda)

---

## [0.1.0] - YYYY-MM-DD (Exemplo de Versão Inicial)

### Adicionado
-   Funcionalidade inicial X.
-   Funcionalidade inicial Y.

### Alterado
-   Melhoria na performance do módulo Z.

### Corrigido
-   Correção de bug que causava problema A ao fazer B.

---

**Como usar este Changelog:**

-   **Adicionado (Added):** para novas funcionalidades.
-   **Alterado (Changed):** para mudanças em funcionalidades existentes.
-   **Corrigido (Fixed):** para quaisquer correções de bugs.
-   **Removido (Removed):** para funcionalidades removidas.
-   **Segurança (Security):** em caso de vulnerabilidades.
-   **Obsoleto (Deprecated):** para funcionalidades que serão removidas em breve.

Mantenha uma seção `[Não Lançado]` no topo para acumular mudanças para o próximo lançamento.
Quando uma nova versão for lançada, renomeie a seção `[Não Lançado]` para a nova versão (ex: `[1.0.0] - YYYY-MM-DD`) e crie uma nova seção `[Não Lançado]` vazia acima dela.
