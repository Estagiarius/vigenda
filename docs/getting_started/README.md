# Guia de Introdução ao Vigenda

Bem-vindo ao Vigenda! Este guia rápido ajudará você a instalar, configurar e começar a usar as funcionalidades básicas do Vigenda.

## Sumário

1.  [O que é o Vigenda?](#o-que-e-o-vigenda)
2.  [Instalação](#instalacao)
    *   [Pré-requisitos](#pre-requisitos)
    *   [Compilando a Partir do Código Fonte](#compilando-a-partir-do-codigo-fonte)
    *   [Usando o Script de Build (Cross-Compilation)](#usando-o-script-de-build-cross-compilation)
3.  [Configuração Inicial](#configuracao-inicial)
    *   [Base de Dados](#base-de-dados)
4.  [Primeiros Passos: Usando o Vigenda](#primeiros-passos-usando-o-vigenda)
    *   [Visualizando o Dashboard](#visualizando-o-dashboard)
    *   [Criando sua Primeira Turma](#criando-sua-primeira-turma)
    *   [Adicionando sua Primeira Tarefa](#adicionando-sua-primeira-tarefa)
    *   [Listando Tarefas](#listando-tarefas)
    *   [Obtendo Ajuda](#obtendo-ajuda)
5.  [Próximos Passos](#proximos-passos)

## 1. O que é o Vigenda?

O Vigenda é uma aplicação de linha de comando (CLI) desenvolvida para auxiliar professores na organização de suas atividades diárias, como gerenciamento de tarefas, turmas, avaliações e muito mais. Ele é projetado para ser uma ferramenta focada e eficiente.

## 2. Instalação

Siga os passos abaixo para instalar o Vigenda em seu sistema.

### Pré-requisitos

Antes de compilar o Vigenda, certifique-se de que possui:

*   **Go**: Versão 1.23 ou superior. Verifique com `go version`.
*   **GCC**: Um compilador C.
    *   Debian/Ubuntu: `sudo apt-get install gcc`
    *   macOS: Xcode Command Line Tools (geralmente instalado).
    *   Windows: MinGW/TDM-GCC.

### Compilando a Partir do Código Fonte

1.  **Obtenha o código-fonte:**
    Se você tem os arquivos do projeto, navegue até o diretório raiz do projeto (onde o arquivo `go.mod` está localizado).

2.  **Compile o projeto:**
    ```bash
    go build -o vigenda ./cmd/vigenda/
    ```
    Isso criará um executável chamado `vigenda` (ou `vigenda.exe` no Windows) no diretório atual.

3.  **(Opcional) Adicione ao PATH:**
    Para usar o `vigenda` de qualquer diretório, mova o executável para um diretório em seu PATH (ex: `/usr/local/bin`) ou adicione o diretório do Vigenda ao seu PATH.

### Usando o Script de Build (Cross-Compilation)

O projeto inclui o script `build.sh` para compilar para diferentes sistemas.

1.  **Torne o script executável:**
    ```bash
    chmod +x build.sh
    ```
2.  **Execute o script:**
    ```bash
    ./build.sh
    ```
    Os binários estarão no diretório `dist/`. Consulte o `AGENTS.md` ou o [Manual do Usuário](../user_manual/README.md#instalacao) para detalhes sobre pré-requisitos de cross-compilation.

## 3. Configuração Inicial

### Base de Dados

Por padrão, o Vigenda usa uma base de dados **SQLite**, que é um arquivo chamado `vigenda.db`. Este arquivo será criado automaticamente:
*   No diretório de configuração do usuário (ex: `~/.config/vigenda/vigenda.db` no Linux).
*   Ou no diretório atual de onde você executa o `vigenda`, se o diretório de configuração não for acessível.

**Para usuários avançados:** Se você deseja usar um local diferente para o arquivo SQLite ou usar PostgreSQL, consulte a seção [Configuração da Base de Dados](../user_manual/README.md#configuracao-da-base-de-dados) no Manual do Usuário.

Por enquanto, nenhuma ação é necessária para começar com a configuração padrão.

## 4. Primeiros Passos: Usando o Vigenda

Com o Vigenda compilado (vamos assumir que ele está no seu diretório atual como `./vigenda`), você pode começar a usá-lo.

### Visualizando o Dashboard

Execute o Vigenda sem argumentos para ver o dashboard:
```bash
./vigenda
```
Isso mostrará um resumo de suas atividades e tarefas.

### Criando sua Primeira Turma

Antes de adicionar tarefas ou avaliações, você provavelmente precisará de uma turma.
(Nota: Para criar turmas, o sistema precisa de "Disciplinas" pré-existentes ou uma forma de criá-las. Assumindo que existe uma Disciplina com ID 1 para este exemplo.)

Use o comando `turma criar`:
```bash
./vigenda turma criar "Minha Primeira Turma (Ex: História 9A)" --subjectid 1
```
Substitua `"Minha Primeira Turma (Ex: História 9A)"` pelo nome desejado e `1` pelo ID da disciplina correta. O sistema informará o ID da turma criada. Anote-o para os próximos passos.

### Adicionando sua Primeira Tarefa

Agora, adicione uma tarefa para a turma que você acabou de criar. Se sua turma recém-criada recebeu o ID 1, por exemplo:

```bash
./vigenda tarefa add "Preparar plano de aula para a primeira semana" --classid 1 --duedate AAAA-MM-DD
```
Substitua `1` pelo ID real da sua turma e `AAAA-MM-DD` pela data de entrega desejada (ex: `2024-08-30`).

### Listando Tarefas

Para ver as tarefas da sua turma:
```bash
./vigenda tarefa listar --classid 1
```
Substitua `1` pelo ID da sua turma.

### Obtendo Ajuda

Para qualquer comando, você pode obter ajuda detalhada adicionando a flag `--help`:
```bash
./vigenda tarefa add --help
./vigenda turma criar --help
```

## 5. Próximos Passos

Você aprendeu o básico para instalar e começar a usar o Vigenda!

*   Para um mergulho profundo em todas as funcionalidades, consulte o **[Manual do Usuário](../user_manual/README.md)**.
*   Tem perguntas específicas? Confira nosso **[FAQ](../faq/README.md)**.
*   Para exemplos práticos, explore nossos **[Tutoriais](../tutorials/README.md)**.

Explore os outros comandos para gerenciar alunos, avaliações, notas e até mesmo gerar provas a partir de um banco de questões!
```bash
# Exemplos para explorar:
./vigenda turma importar-alunos --help
./vigenda avaliacao criar --help
./vigenda bancoq add --help
./vigenda prova gerar --help
```
