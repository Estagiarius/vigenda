# Visão Geral do Projeto

Bem-vindo ao nosso projeto! Este arquivo README fornece uma visão geral do projeto, instruções sobre como começar e informações sobre como contribuir.

## Descrição

**Vigenda** é uma aplicação de linha de comando (CLI) desenvolvida em Go, projetada para ajudar professores e estudantes a gerenciar atividades acadêmicas. A principal forma de interação é através de uma Interface de Texto do Usuário (TUI) robusta, acessada executando `vigenda` sem subcomandos. A aplicação também oferece alguns subcomandos CLI para funcionalidades específicas.

Funcionalidades incluem:
- Gerenciamento de Tarefas
- Gestão de Turmas e Alunos
- Planejamento de Aulas
- Criação e Gestão de Avaliações
- Lançamento de Notas
- Banco de Questões e Geração de Provas
- Ferramentas de Produtividade (como sessões de foco)

O nome "Vigenda" sugere uma combinação de "vida" e "agenda", focando na organização da vida acadêmica.

## Interação Principal: TUI (Interface de Texto do Usuário)

A forma primária de utilizar o Vigenda é através de sua interface interativa no terminal:

```bash
go run ./cmd/vigenda/main.go
```
Ou, se você construiu o binário (ex: `vigenda_cli`):
```bash
./vigenda_cli
```
Isso iniciará a TUI, que oferece um menu principal para acessar todas as funcionalidades de forma interativa, incluindo a criação e gerenciamento de disciplinas, turmas, alunos, aulas, avaliações, etc.

## Comandos CLI Adicionais

Além da TUI principal, alguns subcomandos estão disponíveis para acesso direto a funcionalidades específicas. Para uma lista completa e atualizada, utilize:
```bash
go run ./cmd/vigenda/main.go --help
# ou
./vigenda_cli --help
```
Exemplos de subcomandos incluem (mas não se limitam a):
- `tarefa listar`: Para listar tarefas.
- `tarefa add`: Para adicionar tarefas rapidamente.
- `avaliacao criar`: Para criar avaliações.
- `bancoq add`: Para adicionar questões ao banco.

Consulte `vigenda [comando] --help` para detalhes sobre cada subcomando.

## Instalação

Siga estas instruções para configurar o ambiente de desenvolvimento e executar o projeto. Para instruções mais detalhadas, consulte o arquivo `INSTALLATION.MD`.

### Pré-requisitos

Antes de começar, certifique-se de ter o seguinte instalado:

- **Go:** Versão 1.23.0 ou superior (conforme `go.mod`, toolchain go1.24.3).
- **GCC:** Necessário para algumas dependências Go (como `go-sqlite3`) ou para o processo de build.
- **Git:** Para clonar o repositório.

Consulte `INSTALLATION.MD` para links e instruções de instalação detalhadas para cada pré-requisito.

### Passos de Instalação

1.  **Clone o repositório:**
    (Se você tiver acesso ao repositório Git, substitua `[URL_DO_REPOSITORIO_GIT]` pela URL correta.)
    ```bash
    git clone [URL_DO_REPOSITORIO_GIT] vigenda
    cd vigenda
    ```
    Se você recebeu os arquivos do projeto de outra forma, apenas navegue até o diretório raiz do projeto.

2.  **Instale as dependências Go:**
    Dentro do diretório raiz do projeto (`vigenda`), execute:
    ```bash
    go mod tidy
    ```
    Isso irá baixar as dependências listadas no arquivo `go.mod` e remover as não utilizadas.

### Executando o Projeto

Conforme mencionado, a principal forma de interação é através da TUI:
```bash
go run ./cmd/vigenda/main.go
```
Para usar subcomandos CLI específicos:
```bash
go run ./cmd/vigenda/main.go [subcomando] [flags]
```
Exemplo: `go run ./cmd/vigenda/main.go tarefa listar --all`

Para construir o binário e depois executá-lo:
Consulte o script `build.sh` para as opções de build para diferentes plataformas. Um exemplo básico para sua plataforma atual é:
```bash
go build -o vigenda_cli ./cmd/vigenda/main.go
./vigenda_cli # Para iniciar a TUI
./vigenda_cli [subcomando] [flags] # Para subcomandos CLI
```
O script `build.sh` gera os executáveis no diretório `dist/`.

## Como Contribuir

Estamos abertos a contribuições da comunidade! Se você deseja contribuir, por favor, siga estas diretrizes básicas:

1. **Leia nosso Guia de Contribuição:** Antes de começar, familiarize-se com nosso arquivo `CONTRIBUTING.md`. Ele contém informações detalhadas sobre nosso processo de desenvolvimento, padrões de codificação (incluindo o uso de `gofmt`/`goimports`), como escrever mensagens de commit e o fluxo de Pull Requests.
2. **Crie uma Issue:** Se você encontrou um bug ou tem uma ideia para uma nova funcionalidade, crie uma issue em nosso rastreador de issues.
3. **Faça um Fork e Crie um Branch:** Faça um fork do repositório e crie um novo branch para suas alterações.
4. **Desenvolva e Teste:** Implemente suas alterações e certifique-se de que todos os testes passam.
5. **Envie um Pull Request:** Envie um pull request com uma descrição clara das suas alterações.

Agradecemos por seu interesse em contribuir!

## Licença

Este projeto é distribuído sob a licença [INSIRA A LICENÇA AQUI, ex: MIT, Apache 2.0, etc.]. Consulte o arquivo `LICENSE` para mais detalhes. (Nota: O arquivo `LICENSE` não foi encontrado na listagem inicial, mas é uma prática padrão. Se não houver um, esta seção pode precisar ser ajustada ou o arquivo `LICENSE` criado).

## Contato

Para dúvidas, sugestões ou reporte de bugs, por favor, crie uma "Issue" no rastreador de issues do projeto (se disponível publicamente) ou entre em contato com o mantenedor do projeto através de [INSIRA O MÉTODO DE CONTATO AQUI, ex: email@example.com].
