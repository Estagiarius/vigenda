# Perguntas Frequentes (FAQ) - Vigenda

Este documento responde a algumas perguntas comuns sobre o Vigenda.

## Perguntas Gerais

**P1: O que é o Vigenda?**
R: Vigenda é uma aplicação de linha de comando (CLI) com uma Interface de Texto do Usuário (TUI) robusta, desenvolvida em Go. Seu objetivo é ajudar professores e estudantes a gerenciar suas atividades acadêmicas, como tarefas, aulas, disciplinas, avaliações, banco de questões e mais.

**P2: Qual a principal forma de interagir com o Vigenda?**
R: A principal forma de interação é através da TUI (Interface de Texto do Usuário), que é acessada executando `vigenda` (ou `go run ./cmd/vigenda/main.go`) sem subcomandos. A TUI oferece um menu completo para todas as funcionalidades. Alguns subcomandos CLI também estão disponíveis para acesso rápido a certas operações.

**P3: O Vigenda é gratuito?**
R: A informação sobre a licença do Vigenda deve ser especificada no arquivo `LICENSE` do projeto ou no `README.md` principal. Consulte esses arquivos para detalhes.

## Instalação e Configuração

**P4: Quais são os pré-requisitos para instalar o Vigenda?**
R: Você precisará de Go (versão 1.23.0 ou superior) e GCC (ou um compilador C compatível). Para instruções detalhadas, consulte o [**INSTALLATION.MD**](../../INSTALLATION.MD).

**P5: Como configuro o banco de dados? O Vigenda suporta PostgreSQL?**
R: Por padrão, Vigenda usa SQLite, e o arquivo de banco de dados é criado automaticamente. Você pode configurar o Vigenda para usar PostgreSQL ou um caminho diferente para o arquivo SQLite através de variáveis de ambiente. Consulte a seção "Configuração da Base de Dados" no [**Manual do Usuário**](../user_manual/README.md) para todos os detalhes, incluindo nomes das variáveis de ambiente (`VIGENDA_DB_TYPE`, `VIGENDA_DB_DSN`, etc.).

**P6: As migrações do banco de dados para PostgreSQL são automáticas?**
R: Não. Para PostgreSQL, as migrações de esquema devem ser gerenciadas externamente. O Vigenda aplicará o esquema inicial automaticamente apenas para SQLite se o banco de dados parecer vazio.

## Uso da Aplicação

**P7: Como crio uma nova disciplina ou turma?**
R: A criação e gerenciamento detalhado de disciplinas, turmas, alunos e aulas é feita primariamente através da **TUI principal**. Execute `vigenda` sem argumentos para acessar o menu e navegar até as opções correspondentes.

**P8: Posso importar uma lista de alunos para uma turma?**
R: Sim! Use o comando CLI `vigenda turma importar-alunos ID_DA_TURMA CAMINHO_DO_ARQUIVO_CSV`. O arquivo CSV deve seguir um formato específico. Veja a seção "Formatos de Ficheiros de Importação" no [**Manual do Usuário**](../user_manual/README.md).

**P9: Como funciona o lançamento de notas?**
R: O lançamento de notas é feito interativamente. Primeiro, crie uma avaliação (via TUI ou com `vigenda avaliacao criar`). Depois, use o comando `vigenda avaliacao lancar-notas ID_DA_AVALIACAO`. O sistema listará os alunos da turma para você inserir cada nota.

**P10: O Vigenda pode gerar provas a partir de um banco de questões?**
R: Sim. Você pode adicionar questões ao banco usando `vigenda bancoq add CAMINHO_DO_ARQUIVO_JSON` e depois gerar provas com `vigenda prova gerar --subjectid ID_DA_DISCIPLINA ...`. Consulte o [**Manual do Usuário**](../user_manual/README.md) para as opções de formatação JSON e geração de provas.

**P11: Os diagramas do sistema (casos de uso, arquitetura) estão disponíveis?**
R: Sim, os diagramas de caso de uso em formato PlantUML estão em `docs/diagrams/`. A [**Especificação Técnica**](../../TECHNICAL_SPECIFICATION.MD) também contém diagramas de arquitetura e fluxo de dados em formato Mermaid. Para visualizá-los, você precisará de um renderizador PlantUML/Mermaid.

## Solução de Problemas

**P12: Estou tendo problemas com CGo ou GCC durante a instalação/build. O que fazer?**
R: Certifique-se de que o GCC (ou MinGW-w64 no Windows) está instalado corretamente e no PATH do seu sistema. O [**INSTALLATION.MD**](../../INSTALLATION.MD) tem uma seção de "Solução de Problemas Comuns" que cobre isso.

**P13: Onde encontro os logs da aplicação?**
R: O Vigenda cria um arquivo de log (ex: `vigenda.log`) no diretório de configuração do usuário (como `~/.config/vigenda/` no Linux) ou no diretório atual se o primeiro não for acessível. Consulte `cmd/vigenda/main.go` (função `setupLogging`) para a lógica exata de determinação do caminho.

## Contribuição

**P14: Como posso contribuir para o Vigenda?**
R: Ótimo! Por favor, consulte o [**CONTRIBUTING.MD**](../../CONTRIBUTING.MD) para diretrizes sobre o processo de desenvolvimento, estilo de código, mensagens de commit e como submeter Pull Requests. Se você é um agente de IA, o [**AGENTS.md**](../../AGENTS.md) também tem informações úteis.

---

Se sua pergunta não foi respondida aqui, por favor, consulte o [**Manual do Usuário**](../user_manual/README.md) ou considere abrir uma "Issue" no rastreador de issues do projeto (se disponível).
