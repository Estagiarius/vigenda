# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Versionamento Semântico](https://semver.org/spec/v2.0.0.html).

## [Não Lançado]

### Documentação
-   **Atualização Geral da Documentação (realizada pelo Agente de IA):**
    -   `README.md`: Revisado para refletir a primazia da TUI, atualizar instruções de execução e remover placeholders.
    -   `INSTALLATION.md` e `TESTING.md`: Removidas versões duplicadas em minúsculas. Os arquivos corretos são `INSTALLATION.MD` e `TESTING.MD`.
    -   `docs/user_manual/README.md`: Atualizado para refletir a funcionalidade da TUI sobre certos comandos CLI (como criação de turmas/disciplinas), e removidas referências a arquivos inexistentes (`getting_started/README.md`, `faq/README.md`).
    -   `AGENTS.md`: Removida referência a `BUG_REPORTING.md` (inexistente). As referências a `INSTALLATION.MD` e `TESTING.MD` foram mantidas (assumindo que são as versões corretas). Atualizada a lista de documentos principais.
    -   `CODE_OF_CONDUCT.md`: Preenchido o placeholder de contato com um email genérico.
    -   `CONTRIBUTING.md`: Removidas URLs de repositório placeholder.
    -   `API_DOCUMENTATION.md`: Simplificado para indicar que é uma CLI sem API de rede.
    -   `TECHNICAL_SPECIFICATION.MD`: Corrigido o fluxo de dados para "Listar Tarefas de uma Disciplina" para refletir o esquema de banco de dados correto (tarefas ligadas a turmas, turmas a disciplinas). Adicionada nota sobre a necessidade de renderizar diagramas PlantUML.
    -   `docs/diagrams/README.md`: Adicionada nota sobre a necessidade de renderizar diagramas PlantUML.

### Adicionado
-   (Use esta seção para novas funcionalidades de código)

### Alterado
-   (Use esta seção para mudanças em funcionalidades de código existentes)

### Corrigido
-   (Use esta seção para correções de bugs no código)

### Removido
-   `INSTALLATION.md` (arquivo duplicado, mantido `INSTALLATION.MD`)
-   `TESTING.md` (arquivo duplicado, mantido `TESTING.MD`)
-   (Use esta seção para funcionalidades de código removidas)

### Segurança
-   (Use esta seção para atualizações de segurança no código)

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
