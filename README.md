# Vigenda

## Visão Geral

Bem-vindo ao Vigenda!

**Vigenda** é uma aplicação de linha de comando desenvolvida em Go, com uma **Interface de Texto do Usuário (TUI)** robusta, projetada para ajudar professores e educadores a gerenciar suas atividades acadêmicas de forma eficiente.

O nome "Vigenda" sugere uma combinação de "vida" e "agenda", focando na organização da vida acadêmica.

## Funcionalidades

- **Gestão de Tarefas:** Crie e gerencie tarefas, associando-as a turmas específicas.
- **Gestão de Turmas e Disciplinas:** Organize seu conteúdo por disciplinas e turmas.
- **Gestão de Alunos:** Mantenha um registro dos alunos em cada turma.
- **Gestão de Avaliações e Notas:** Crie avaliações, lance notas e calcule médias ponderadas.
- **Banco de Questões:** Crie um banco de questões para gerar provas personalizadas.
- **Geração de Provas:** Gere provas em formato de texto a partir do seu banco de questões.

## Interação Principal: TUI (Interface de Texto do Usuário)

A forma primária de utilizar o Vigenda é através de sua interface interativa no terminal. Para iniciá-la, execute o seguinte comando no diretório raiz do projeto:

```bash
go run ./cmd/vigenda/main.go
```

Ou, se você construiu o binário:

```bash
./vigenda
```

Isso iniciará a TUI, que oferece um menu principal para acessar todas as funcionalidades de forma interativa.

## Instalação

Para instruções detalhadas sobre como instalar e configurar o ambiente de desenvolvimento, consulte o arquivo [**INSTALLATION.MD**](./INSTALLATION.MD).

### Pré-requisitos

- **Go:** Versão 1.23 ou superior.
- **GCC:** Compilador C para CGO (usado pela dependência `go-sqlite3`).
- **Git:** Para clonar o repositório.

### Passos de Instalação

1.  **Clone o repositório** (se aplicável).
2.  **Instale as dependências Go:**
    ```bash
    go mod tidy
    ```

## Como Contribuir

Estamos abertos a contribuições! Se você deseja contribuir, por favor, leia nosso [**Guia de Contribuição**](./CONTRIBUTING.md). Ele contém informações detalhadas sobre nosso processo de desenvolvimento, padrões de codificação e fluxo de Pull Requests.

## Chat com IA

O Vigenda inclui um recurso de chat que utiliza a API da OpenAI. Para usar este recurso, você precisa configurar sua chave de API e, opcionalmente, uma URL base da API.

1.  **Inicie o Vigenda:**
    ```bash
    go run ./cmd/vigenda/main.go
    ```
2.  **Acesse as Configurações:** No menu principal, selecione a opção "Configurações".
3.  **Insira suas credenciais:**
    *   No campo "Sua chave de API da OpenAI", insira sua chave de API.
    *   Se você estiver usando um proxy ou um endpoint personalizado, insira a URL no campo "URL base da API (opcional)".
4.  **Salve as configurações:** Pressione `enter` para salvar as configurações.
5.  **Use o Chat:** Agora você pode acessar a tela "Chat com IA" e interagir com o modelo de linguagem.

## Documentação Adicional

- **[Documentação do Usuário](./docs/USER_DOCUMENTATION.md):** Um guia completo para todas as funcionalidades.
- **[Documentação do Desenvolvedor](./docs/developer/README.md):** Informações sobre a arquitetura do projeto e como estendê-lo.
- **[Exemplos de Uso da TUI](./docs/user_manual/TUI_EXAMPLES.md):** Um guia rápido para realizar as tarefas mais comuns.
- **[Relato de Bugs](./BUG_REPORTING.md):** Como relatar bugs e problemas.
- **[Esquema do Banco de Dados](./DATABASE_SCHEMA.md):** Detalhes sobre a estrutura do banco de dados.
- **[Especificação Técnica](./TECHNICAL_SPECIFICATION.MD):** Informações técnicas sobre o projeto.
- **[Testes](./TESTING.MD):** Como executar os testes do projeto.
