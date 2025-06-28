# Visão Geral do Projeto

Bem-vindo ao nosso projeto! Este arquivo README fornece uma visão geral do projeto, instruções sobre como começar e informações sobre como contribuir.

## Descrição

**Vigenda** é uma aplicação de linha de comando (CLI) desenvolvida em Go, projetada para ajudar estudantes a gerenciar suas atividades acadêmicas. Ela oferece uma Interface de Texto do Usuário (TUI) para interações e parece incluir funcionalidades para gerenciamento de tarefas, aulas, disciplinas, avaliações e acompanhamento de progresso. O nome "Vigenda" sugere uma combinação de "vida" e "agenda", focando na organização da vida acadêmica.

## Instalação e Execução

Siga estas instruções para configurar o ambiente de desenvolvimento e executar o projeto. Para instruções mais detalhadas, consulte o arquivo `INSTALLATION.MD`.

### Pré-requisitos

Antes de começar, certifique-se de ter o seguinte instalado:

- **Go:** Versão 1.23.0 ou superior (conforme `go.mod`, toolchain go1.24.3).
- **GCC:** Necessário para algumas dependências Go (como `go-sqlite3`) ou para o processo de build.
- **Git:** Para clonar o repositório.

Consulte `INSTALLATION.MD` para links e instruções de instalação detalhadas para cada pré-requisito.

### Passos de Instalação

1.  **Clone o repositório:**
    ```bash
    git clone https://github.com/usuario/vigenda.git # Substitua pela URL correta do repositório
    cd vigenda
    ```
    (Assumindo que o diretório do projeto é `vigenda` com base no `cmd/vigenda` e `go.mod`)

2.  **Instale as dependências Go:**
    Dentro do diretório do projeto (`vigenda`), execute:
    ```bash
    go mod tidy
    ```
    Isso irá baixar as dependências listadas no arquivo `go.mod` e remover as não utilizadas.

### Executando o Projeto

Para executar a aplicação a partir do código fonte (a partir da raiz do projeto `vigenda`):
```bash
go run ./cmd/vigenda/main.go [comando_ou_argumentos]
```
Exemplo: `go run ./cmd/vigenda/main.go tarefa listar`

Para construir o binário e depois executá-lo:
Consulte o script `build.sh` para as opções de build para diferentes plataformas. Um exemplo básico é:
```bash
go build -o vigenda_cli ./cmd/vigenda/main.go
./vigenda_cli [comando_ou_argumentos]
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

[Especifique a licença sob a qual o projeto é distribuído, por exemplo: MIT, Apache 2.0, etc.]

## Contato

[Opcional: Adicione informações de contato ou links para canais de comunicação do projeto.]
