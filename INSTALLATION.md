# Documentação de Instalação

Este documento fornece instruções detalhadas sobre como configurar o ambiente de desenvolvimento para este projeto e instalar todas as dependências necessárias.

## 1. Pré-requisitos

Antes de prosseguir com a instalação, certifique-se de que seu sistema atende aos seguintes pré-requisitos:

### 1.1. Sistema Operacional

-   **Linux:** Ubuntu 20.04+, Fedora 34+ ou distribuições Linux equivalentes.
-   **macOS:** 11.0 (Big Sur) ou superior. (Cross-compilação para macOS a partir do Linux pode exigir OSXCross ou configuração similar devido ao CGO).
-   **Windows:** Windows 10+ (WSL2 é recomendado para uma experiência de desenvolvimento semelhante ao Linux, ou use os binários compilados para Windows).

### 1.2. Software Essencial

-   **Git:** Necessário para clonar o repositório.
    -   Verifique a instalação: `git --version`
    -   Instruções de instalação: [https://git-scm.com/book/en/v2/Getting-Started-Installing-Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

-   **Go (Golang):** Versão 1.23.0 ou superior (toolchain go1.24.3, conforme `go.mod`).
    -   Verifique a instalação: `go version`
    -   Instruções de instalação: [https://golang.org/doc/install](https://golang.org/doc/install)
    -   Certifique-se de que `$GOROOT` está configurado corretamente e que `$GOPATH/bin` (se ainda usar GOPATH) e o diretório `bin` da instalação do Go (ex: `/usr/local/go/bin` ou `~/go/bin`) estão no seu `PATH`.

-   **GCC (GNU Compiler Collection):** Necessário para compilar dependências CGo, como o driver `mattn/go-sqlite3`.
    -   Verifique a instalação: `gcc --version`
    -   Instruções de instalação:
        -   **Linux (Debian/Ubuntu):** `sudo apt update && sudo apt install build-essential`
        -   **Linux (Fedora):** `sudo dnf groupinstall "Development Tools"`
        -   **macOS:** Instale as Xcode Command Line Tools: `xcode-select --install`
        -   **Windows:** É necessário um compilador C compatível com CGO, como o MinGW-w64. O script `build.sh` tenta usar `x86_64-w64-mingw32-gcc` para builds Windows. Instale-o (ex: via MSYS2) e adicione ao PATH.

-   **(Opcional para Cross-Compilação Windows a partir do Linux/macOS) Mingw-w64:**
    -   Exemplo de instalação no Ubuntu: `sudo apt install mingw-w64`
    -   O script `build.sh` usa `CC=x86_64-w64-mingw32-gcc` para compilar para Windows de 64 bits.

### 1.3. Variáveis de Ambiente (Opcional)

-   Para cross-compilação CGO, a variável `CC` pode precisar ser definida para apontar para o compilador C correto para o alvo (ex: `CC=x86_64-w64-mingw32-gcc` ao compilar para Windows 64-bit a partir do Linux). O script `build.sh` tenta fazer isso.

## 2. Clonando o Repositório

1.  Abra seu terminal ou prompt de comando.
2.  Navegue até o diretório onde você deseja clonar o projeto.
3.  Execute o seguinte comando (substitua pela URL correta se necessário):
    ```bash
    git clone https://github.com/seu-usuario/vigenda.git
    cd vigenda
    ```

## 3. Instalando Dependências do Projeto

### 3.1. Dependências Go

O projeto usa Go Modules para gerenciamento de dependências. O arquivo `go.mod` lista todas as dependências diretas e indiretas.

1.  **Navegue até a raiz do projeto clonado (`vigenda`).**
2.  **Execute o seguinte comando para sincronizar as dependências:**
    ```bash
    go mod tidy
    ```
    Este comando garante que o arquivo `go.mod` e `go.sum` estejam consistentes, baixando as dependências necessárias e removendo as não utilizadas.

### 3.2. Ferramentas de Desenvolvimento Go (Opcional, mas Recomendado)

Para formatação de código e linting, que são boas práticas de desenvolvimento:

-   **`goimports`** (para formatação e organização automática de imports):
    ```bash
    go install golang.org/x/tools/cmd/goimports@latest
    ```
-   **`golangci-lint`** (um agregador rápido de linters para Go):
    -   Verifique o site oficial para o método de instalação mais recente: [https://golangci-lint.run/usage/install/](https://golangci-lint.run/usage/install/)
    -   Exemplo de instalação (pode variar dependendo do seu SO e preferências):
        ```bash
        # Linux/macOS (binário)
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2
        ```
        (Substitua `v1.57.2` pela versão mais recente ou desejada).
        Certifique-se de que `$(go env GOPATH)/bin` ou `~/go/bin` está no seu `PATH`.

### 3.3. Dependências do Sistema Adicionais

-   Nenhuma outra dependência de sistema é explicitamente necessária além de Go e GCC para a compilação e execução básica do Vigenda.

## 4. Configuração do Ambiente

### 4.1. Arquivos de Configuração

-   Atualmente, Vigenda não requer arquivos de configuração externos para sua operação básica. As configurações, como o caminho do banco de dados, são gerenciadas internamente (consulte `internal/config/config.go` e `internal/database/database.go`).
-   Se configurações personalizadas forem introduzidas no futuro, esta seção será atualizada.

### 4.2. Configuração do Banco de Dados

-   O Vigenda utiliza **SQLite** como banco de dados, que é baseado em arquivo.
-   O arquivo do banco de dados (por exemplo, `vigenda.db` ou similar, dependendo da lógica em `internal/database/database.go`) é criado e gerenciado automaticamente pela aplicação no diretório de dados apropriado (geralmente no diretório de configuração do usuário ou no diretório do projeto).
-   As migrações de esquema do banco de dados estão localizadas em `internal/database/migrations/` (ex: `001_initial_schema.sql`) e são aplicadas automaticamente pela aplicação na inicialização, se o banco de dados ainda não estiver no esquema esperado. Não há comando manual de migração para o usuário executar.

## 5. Construindo e Executando a Aplicação

### 5.1. Executando a Partir do Código Fonte (Desenvolvimento)

Para executar a aplicação Vigenda diretamente do código fonte (a partir da raiz do projeto `vigenda`):
```bash
go run ./cmd/vigenda/main.go [comando_da_cli_vigenda] [argumentos_do_comando]
```
Exemplos:
```bash
go run ./cmd/vigenda/main.go --help
go run ./cmd/vigenda/main.go tarefa listar
go run ./cmd/vigenda/main.go foco iniciar --duracao 25m --tarefa "Estudar Go"
```

### 5.2. Construindo o Binário

O projeto inclui um script `build.sh` para construir binários para diferentes plataformas (Linux, Windows, macOS).

-   **Para executar o script de build (necessita de `bash` e possivelmente compiladores C específicos para cross-compilação, como MinGW para Windows):**
    ```bash
    chmod +x build.sh # Certifique-se de que o script é executável
    ./build.sh
    ```
    Os binários serão gerados no diretório `dist/` (ex: `dist/vigenda-linux-amd64`, `dist/vigenda-windows-amd64.exe`).

-   **Para construir manualmente para sua plataforma atual (exemplo para Linux):**
    Use o comando `go build`. O flag `-o` especifica o nome do arquivo de saída.
    ```bash
    go build -o vigenda_cli ./cmd/vigenda/main.go
    ```
    Isso criará um executável chamado `vigenda_cli` no diretório raiz do projeto.

    Para builds menores, removendo informações de debug e símbolos (similar ao `build.sh`):
    ```bash
    go build -ldflags="-s -w" -o vigenda_cli ./cmd/vigenda/main.go
    ```

### 5.3. Executando o Binário Compilado

Após construir o binário (ex: `vigenda_cli` ou um dos binários em `dist/`):
```bash
./vigenda_cli [comando_da_cli_vigenda] [argumentos_do_comando]
# Exemplo usando um binário de dist:
# ./dist/vigenda-linux-amd64 tarefa listar
```

## 6. Verificando a Instalação

Após seguir os passos de instalação de dependências e, opcionalmente, de build:

1.  **Execute os Testes:**
    Consulte o arquivo `TESTING.MD` para instruções detalhadas sobre como executar os testes. O comando principal é:
    ```bash
    go test ./...
    ```
    Uma suíte de testes passando geralmente indica uma configuração correta do ambiente de desenvolvimento.

2.  **Execute a Aplicação:**
    Tente executar um comando básico da aplicação para verificar se ela inicia e responde:
    ```bash
    go run ./cmd/vigenda/main.go --help
    # ou, se você construiu o binário:
    # ./vigenda_cli --help
    ```
    Isso deve exibir a mensagem de ajuda da CLI, confirmando que a aplicação pode ser executada corretamente.

## 7. Solução de Problemas Comuns

-   **Erro `command not found: go`:** Certifique-se de que o Go está instalado e que o diretório `bin` do Go (ex: `/usr/local/go/bin` ou `~/go/bin`) está no `PATH` do seu sistema.
-   **Problemas de Permissão:** Ao instalar ferramentas globais ou pacotes de sistema, você pode precisar de privilégios de `sudo`.
-   **Dependências CGo Falhando (especialmente em cross-compilação):**
    -   Verifique se o GCC (ou o compilador C correto para o alvo, como MinGW para Windows) está instalado e configurado corretamente.
    -   Para cross-compilação, pode ser necessário definir a variável de ambiente `CC` para o compilador C do sistema alvo. O script `build.sh` tenta fazer isso para Windows.
    -   Consulte a documentação do `mattn/go-sqlite3` para requisitos específicos de CGO e cross-compilação.
-   **Erro de Conexão com o Banco de Dados:** Como o Vigenda usa SQLite, os erros geralmente estão relacionados a permissões de arquivo no diretório onde o arquivo `.db` é armazenado, ou um arquivo de banco de dados corrompido.

Se você encontrar outros problemas, verifique as `Issues` do projeto no GitHub (se disponível) ou crie uma nova issue com detalhes do erro.

---

Com estes passos, seu ambiente de desenvolvimento deve estar pronto para você começar a trabalhar no projeto Vigenda!
