# Guia de Introdução ao Vigenda

Bem-vindo ao Vigenda! Este guia ajudará você a instalar o Vigenda, configurar o ambiente e dar os primeiros passos para gerenciar suas atividades acadêmicas.

## 1. Instalação

Para instruções detalhadas de instalação, incluindo pré-requisitos de sistema (Go, GCC) e como lidar com dependências, consulte o documento principal [**INSTALLATION.MD**](../../INSTALLATION.MD).

Resumidamente, os passos são:

1.  **Pré-requisitos:**
    *   Tenha o **Go** (versão 1.23.0 ou superior) instalado.
    *   Tenha o **GCC** (ou um compilador C compatível) instalado para dependências CGO.
    *   Tenha o **Git** instalado (se for clonar o repositório).

2.  **Obtenha o Projeto:**
    *   Se você tem acesso ao repositório Git:
        ```bash
        git clone [URL_DO_REPOSITORIO_GIT] vigenda
        cd vigenda
        ```
    *   Caso contrário, navegue até o diretório onde os arquivos do projeto foram extraídos.

3.  **Instale as Dependências Go:**
    Dentro do diretório raiz do projeto (`vigenda`):
    ```bash
    go mod tidy
    ```

## 2. Executando o Vigenda pela Primeira Vez

A principal forma de interagir com o Vigenda é através da sua Interface de Texto do Usuário (TUI).

### Iniciando a TUI

Para iniciar a aplicação no modo TUI:
```bash
# A partir do diretório raiz do projeto
go run ./cmd/vigenda/main.go
```
Isso compilará e executará a aplicação, abrindo o menu principal interativo.

**Alternativamente, você pode construir um binário primeiro:**
```bash
# Construir o binário (ex: vigenda_cli)
go build -o vigenda_cli ./cmd/vigenda/main.go

# Executar o binário
./vigenda_cli
```
Os binários para diferentes plataformas também podem ser gerados usando o script `./build.sh` (consulte `INSTALLATION.MD`).

### Navegando na TUI

*   Use as **teclas de seta (Cima/Baixo)** para navegar pelas opções do menu.
*   Pressione **Enter** para selecionar uma opção.
*   Pressione **Esc (Escape)** para voltar ao menu anterior ou sair de uma visualização/formulário.

A TUI é autoexplicativa e guiará você através da criação de disciplinas, turmas, lançamento de tarefas, etc.

## 3. Exemplo de Uso de um Comando CLI

Embora a TUI seja central, alguns comandos CLI estão disponíveis para acesso rápido.

**Exemplo: Listar todas as tarefas (de sistema e de todas as turmas)**
```bash
# A partir do diretório raiz do projeto
go run ./cmd/vigenda/main.go tarefa listar --all
```
Ou, se você construiu o binário `./vigenda_cli`:
```bash
./vigenda_cli tarefa listar --all
```

Para ver todos os comandos disponíveis e suas opções:
```bash
go run ./cmd/vigenda/main.go --help
# ou
./vigenda_cli --help
```
E para um comando específico:
```bash
go run ./cmd/vigenda/main.go tarefa --help
```

## 4. Próximos Passos

*   **Explore a TUI:** Dedique um tempo para navegar por todas as opções do menu principal da TUI. Esta é a forma mais completa de usar o Vigenda.
*   **Consulte o Manual do Usuário:** Para um guia detalhado de todas as funcionalidades, comandos CLI, formatos de arquivo de importação e configuração, leia o [**Manual do Usuário**](../user_manual/README.md).
*   **Veja os Tutoriais:** Para exemplos práticos de cenários de uso, como gerenciar avaliações ou o banco de questões, explore nossos [**Tutoriais**](../tutorials/).

Se você é um desenvolvedor interessado em contribuir ou entender a arquitetura, consulte o [**Guia do Desenvolvedor**](../developer/README.md) e a [**Especificação Técnica**](../../TECHNICAL_SPECIFICATION.MD).

Aproveite o Vigenda para organizar sua vida acadêmica!
