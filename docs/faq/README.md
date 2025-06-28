# FAQ - Perguntas Frequentes sobre o Vigenda

Esta seção responde a algumas das perguntas mais comuns sobre o Vigenda.

## Geral

**P1: O que é o Vigenda?**
R: O Vigenda é uma aplicação de linha de comando (CLI) projetada para ajudar professores a organizar tarefas, aulas, avaliações e outras atividades pedagógicas de forma eficiente.

**P2: Para quem o Vigenda é destinado?**
R: É especialmente útil para professores que buscam uma ferramenta de organização focada e direta, incluindo aqueles com TDAH que podem se beneficiar de interfaces menos distrativas.

## Instalação e Configuração

**P3: Como faço para instalar o Vigenda?**
R: A instalação envolve compilar o código-fonte. Você precisará de Go (versão 1.23+) e GCC instalados.
   1. Obtenha os arquivos do projeto.
   2. Navegue até o diretório raiz do projeto.
   3. Execute `go build -o vigenda ./cmd/vigenda/`.
   Para instruções detalhadas, consulte o [Guia de Introdução](../getting_started/README.md#instalacao) ou o [Manual do Usuário](../user_manual/README.md#instalacao).

**P4: Quais sistemas operacionais são suportados?**
R: O Vigenda pode ser compilado para Linux, Windows e macOS. O script `build.sh` fornecido facilita a compilação para Linux e Windows a partir de um ambiente Linux. Compilar para macOS é preferencialmente feito em uma máquina macOS.

**P5: Preciso configurar uma base de dados complexa?**
R: Não para começar. Por padrão, o Vigenda usa SQLite, que cria um arquivo de banco de dados simples (`vigenda.db`) automaticamente no seu diretório de configuração ou no diretório atual. Nenhuma configuração adicional é necessária para o modo padrão.

**P6: Como configuro o Vigenda para usar PostgreSQL?**
R: Você pode configurar o Vigenda para usar PostgreSQL definindo variáveis de ambiente como `VIGENDA_DB_TYPE="postgres"`, `VIGENDA_DB_HOST`, `VIGENDA_DB_USER`, `VIGENDA_DB_PASSWORD`, `VIGENDA_DB_NAME`, etc., ou usando uma DSN completa com `VIGENDA_DB_DSN`. Consulte a seção [Configuração da Base de Dados](../user_manual/README.md#configuracao-da-base-de-dados) no Manual do Usuário para todos os detalhes.

**P7: O que significa DSN?**
R: DSN significa "Data Source Name". É uma string que contém todas as informações necessárias para se conectar a uma base de dados, como tipo de banco, endereço do servidor, nome de usuário, senha, nome da base, etc.

**P8: Onde o arquivo do banco de dados SQLite é salvo por padrão?**
R: O Vigenda tenta salvar em um diretório de configuração específico do usuário (ex: `~/.config/vigenda/vigenda.db` no Linux). Se não conseguir, ele salvará no diretório de onde o `vigenda` foi executado. Você pode especificar um caminho customizado com a variável de ambiente `VIGENDA_DB_PATH`.

## Uso

**P9: Esqueci como adicionar uma tarefa, onde encontro ajuda?**
R: Você pode usar a flag `--help` com qualquer comando ou subcomando para ver as opções disponíveis. Por exemplo:
   ```bash
   ./vigenda tarefa add --help
   ./vigenda tarefa --help
   ./vigenda --help
   ```
   O [Manual do Usuário](../user_manual/README.md) também detalha todos os comandos.

**P10: Como importo uma lista de alunos para uma turma?**
R: Use o comando `vigenda turma importar-alunos <ID_DA_TURMA> <ARQUIVO_CSV>`. O arquivo CSV precisa ter colunas como `nome_completo` e, opcionalmente, `numero_chamada` e `situacao`. Veja o formato detalhado no [Manual do Usuário](../user_manual/README.md#importacao-de-alunos-csv).

**P11: Posso gerar provas com questões de diferentes níveis de dificuldade?**
R: Sim! O comando `vigenda prova gerar` permite especificar quantas questões fáceis, médias ou difíceis você quer, usando opções como `--easy <num>`, `--medium <num>`, `--hard <num>`.

**P12: O que é o Dashboard Interativo?**
R: Ao executar `./vigenda` sem nenhum subcomando, você verá o Dashboard. Ele oferece uma visão geral rápida de suas tarefas do dia, eventos e outras informações relevantes.

## Dados e Sincronização

**P13: Meus dados do Vigenda ficam salvos localmente?**
R: Sim. Se você estiver usando SQLite (o padrão), todos os seus dados são salvos em um único arquivo `.db` no seu computador. Se estiver usando PostgreSQL, os dados ficam no servidor PostgreSQL que você configurou.

**P14: Posso usar o Vigenda em várias máquinas sincronizando os dados?**
R: *   Com **SQLite**: A sincronização é manual. Você precisaria copiar o arquivo `vigenda.db` entre as máquinas ou usar um serviço de sincronização de arquivos (como Dropbox, Google Drive, Syncthing), mas tenha cuidado com conflitos se o arquivo for modificado em dois lugares ao mesmo tempo.
    *   Com **PostgreSQL**: Sim, se você configurar o Vigenda em várias máquinas para se conectar ao mesmo servidor PostgreSQL, todos acessarão os mesmos dados em tempo real. Esta é a abordagem recomendada para uso em múltiplos dispositivos.

## Solução de Problemas

**P15: Recebi um erro ao tentar compilar. O que devo fazer?**
R: Verifique se você tem Go e GCC (ou o compilador C apropriado para seu sistema) instalados e nas versões corretas, conforme listado nos [pré-requisitos](../getting_started/README.md#pre-requisitos). Se o erro for relacionado a uma dependência específica (como `go-sqlite3`), certifique-se de que as ferramentas de compilação C estão funcionando.

**P16: O comando `vigenda` não é encontrado no terminal.**
R: Isso geralmente significa que o local do executável `vigenda` não está no PATH do seu sistema. Você pode:
    1.  Navegar até o diretório onde `vigenda` foi compilado e executá-lo com `./vigenda`.
    2.  Mover o executável `vigenda` para um diretório que já esteja no seu PATH (ex: `/usr/local/bin` no Linux/macOS).
    3.  Adicionar o diretório do `vigenda` ao seu PATH permanentemente.

---

Não encontrou sua pergunta aqui? Consulte o [Manual do Usuário](../user_manual/README.md) para informações mais detalhadas ou os [Tutoriais](../tutorials/README.md) para exemplos práticos.
