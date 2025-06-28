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

-   Nenhuma outra dependência de sistema é explicitamente necessária além de Go e GCC para la compilação e execução básica do Vigenda.

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

-   **Erro `command not found: go`:**
    -   **Causa:** O diretório de instalação do Go (contendo o executável `go`) não está no `PATH` do sistema.
    -   **Solução:**
        -   Verifique se o Go está instalado corretamente executando `which go` (Linux/macOS) ou `where go` (Windows).
        -   Adicione o diretório `bin` da sua instalação Go (ex: `/usr/local/go/bin`, `C:\Go\bin`, `~/go/bin`) à variável de ambiente `PATH`. Pode ser necessário reiniciar o terminal ou o sistema.
        -   Certifique-se de que `$GOROOT` está definido (geralmente `/usr/local/go` ou `C:\Go`) e `$GOROOT/bin` está no `PATH`.

-   **Problemas de Permissão:**
    -   **Causa:** Tentativa de instalar ferramentas globais (ex: `go install ...@latest` em diretórios do sistema) ou executar operações que exigem privilégios elevados sem `sudo` (Linux/macOS) ou como Administrador (Windows).
    -   **Solução:**
        -   Para instalações globais de ferramentas Go em `$GOPATH/bin` ou `~/go/bin`, geralmente não é necessário `sudo` se `$GOPATH` ou `~/go` estiverem em seu diretório home.
        -   Use `sudo` para comandos de sistema como `apt install` ou `dnf groupinstall`.
        -   No Windows, execute o prompt de comando ou PowerShell como Administrador se necessário.

-   **Dependências CGo Falhando (ex: `mattn/go-sqlite3`):**
    -   **Causa Comum:** Falta de um compilador C (GCC/Clang/MinGW) no sistema ou o compilador não está acessível no `PATH`. `go-sqlite3` requer CGO habilitado.
    -   **Soluções:**
        -   **Linux:** `sudo apt update && sudo apt install build-essential` (Debian/Ubuntu) ou `sudo dnf groupinstall "Development Tools"` (Fedora).
        -   **macOS:** Instale as Xcode Command Line Tools: `xcode-select --install`.
        -   **Windows:** Instale o MinGW-w64 (certifique-se de obter uma versão que corresponda à sua arquitetura Go, ex: 64-bit) e adicione o diretório `bin` do MinGW (ex: `C:\msys64\mingw64\bin`) ao `PATH` do sistema.
        -   **Cross-Compilação:** Se estiver compilando para um SO/arquitetura diferente do seu sistema host (ex: Windows a partir do Linux), você precisará de um cross-compilador C.
            -   Para Windows: `sudo apt install mingw-w64` no Linux. O script `build.sh` tenta usar `CC=x86_64-w64-mingw32-gcc`.
            -   Para macOS a partir do Linux: Isso é mais complexo e geralmente requer OSXCross ou um SDK do macOS.
        -   **Verifique `CGO_ENABLED`:** Por padrão, é `1`. Se foi desabilitado (`CGO_ENABLED=0`), `go-sqlite3` não compilará. Execute `go env CGO_ENABLED` para verificar.

-   **Erro de Conexão com o Banco de Dados SQLite:**
    -   **Causa:** Geralmente não é um problema de "conexão" como em bancos de dados cliente-servidor, mas sim de acesso ao arquivo.
    -   **Soluções:**
        -   **Permissões de Arquivo/Diretório:** Verifique se a aplicação tem permissão para ler/escrever no diretório onde o arquivo `vigenda.db` (ou similar) é/será criado. O local exato pode ser determinado pela lógica em `internal/database/database.go` e `internal/config/config.go`.
        -   **Arquivo de Banco de Dados Corrompido:** Em casos raros, o arquivo `.db` pode se corromper. Tente remover ou renomear o arquivo (faça backup se contiver dados importantes!) para permitir que a aplicação o recrie.
        -   **Caminho do Banco de Dados Incorreto:** Verifique a lógica em `internal/config/config.go` para garantir que o caminho do banco de dados está sendo determinado corretamente para o seu sistema operacional.

-   **Problemas com `go mod tidy` ou `go build`:**
    -   **Causa: Conflitos de versão ou `go.sum` desatualizado.**
        -   **Solução:** Tente `go clean -modcache` para limpar o cache de módulos e depois execute `go mod tidy` novamente.
    -   **Causa: Problemas de rede ao baixar dependências.**
        -   **Solução:** Verifique sua conexão com a internet. Se estiver atrás de um proxy, configure as variáveis de ambiente `GOPROXY`, `HTTP_PROXY`, `HTTPS_PROXY`.
    -   **Causa: Versão do Go incompatível com `go.mod`.**
        -   **Solução:** Certifique-se de que a versão do Go instalada (verificada com `go version`) é compatível com a versão especificada no `go.mod` (ex: `go 1.23.0`). Atualize o Go se necessário.

-   **Problemas com `golangci-lint` ou `goimports`:**
    -   **Causa: Ferramenta não instalada ou não encontrada no `PATH`.**
        -   **Solução:** Siga as instruções de instalação na seção 3.2 deste documento ou na documentação oficial das ferramentas. Verifique se `$(go env GOPATH)/bin` ou `~/go/bin` está no seu `PATH`.
    -   **Causa: `golangci-lint` reporta muitos erros.**
        -   **Solução:** Corrija os erros de linting. Alguns podem ser corrigidos automaticamente com `golangci-lint run --fix ./...`. Para outros, você precisará editar o código manualmente. Considere configurar um arquivo `.golangci.yml` para personalizar os linters e suas regras.

Se você encontrar outros problemas não listados aqui, verifique as `Issues` do projeto no GitHub (se o projeto estiver hospedado lá e tiver um rastreador de issues público) ou crie uma nova issue com detalhes do erro, incluindo seu sistema operacional, versão do Go e os passos para reproduzir o problema.

---

Com estes passos, seu ambiente de desenvolvimento deve estar pronto para você começar a trabalhar no projeto Vigenda!
