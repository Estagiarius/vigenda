// Package main é o ponto de entrada da aplicação Vigenda.
// Ele configura a CLI usando o pacote Cobra, inicializa serviços,
// e gerencia o ciclo de vida da aplicação, incluindo logging e conexão com banco de dados.
// Se nenhum subcomando for fornecido, a TUI principal (gerenciada por internal/app) é iniciada.
package main

import (
	"context"
	"database/sql"
	"encoding/json" // Adicionado para exibir opções de questões de múltipla escolha.
	"fmt"
	"log" // Usado para logging em arquivo.
	"os"
	"path/filepath" // Para manipulação de caminhos de arquivo para logging.
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table" // Usado para formatação de tabelas na CLI.
	"github.com/spf13/cobra"                // Framework CLI.
	"vigenda/internal/app"                  // Pacote da TUI principal.
	"vigenda/internal/database"             // Pacote de configuração e conexão com DB.
	"vigenda/internal/models"               // Structs de modelo do domínio.
	"vigenda/internal/repository"           // Interfaces e implementações de repositório.
	"vigenda/internal/service"              // Interfaces e implementações de serviço.
	"vigenda/internal/tui"                  // Componentes TUI reutilizáveis (ex: prompt).
)

// db é a conexão global com o banco de dados, inicializada em PersistentPreRunE.
var db *sql.DB

// logFile é o descritor de arquivo para o arquivo de log global.
// É necessário para fechar o arquivo corretamente no final da execução da aplicação.
var logFile *os.File

// Variáveis globais para os serviços da aplicação.
// São inicializadas em PersistentPreRunE após a conexão com o banco de dados ser estabelecida.
// Esta abordagem com variáveis globais para serviços é comum em CLIs simples com Cobra,
// mas para aplicações maiores ou mais complexas, a injeção de dependência via construtores
// para os comandos Cobra (se possível) ou um container de DI seriam alternativas.
var (
	taskService       service.TaskService
	classService      service.ClassService
	assessmentService service.AssessmentService
	questionService   service.QuestionService
	proofService      service.ProofService
	// lessonService é declarado separadamente abaixo devido à ordem de inicialização.
)

// rootCmd é o comando raiz da aplicação Vigenda.
// Quando executado sem subcomandos, ele inicia a Interface de Texto do Usuário (TUI).
// PersistentPreRunE é usado para inicializar o logging, a conexão com o banco de dados
// e os serviços antes da execução de qualquer comando (incluindo o Run do rootCmd ou subcomandos).
var rootCmd = &cobra.Command{
	Use:   "vigenda",
	Short: "Vigenda é uma CLI para auxiliar na gestão de atividades acadêmicas.",
	Long: `Vigenda é uma aplicação de linha de comando (CLI) com uma robusta Interface de Texto do Usuário (TUI),
projetada para ajudar professores e estudantes a organizar tarefas, aulas, avaliações e outras
atividades pedagógicas de forma eficiente.

A principal forma de interação é através da TUI, iniciada executando 'vigenda' sem subcomandos.
Subcomandos CLI também estão disponíveis para acesso direto a funcionalidades específicas.

Funcionalidades Principais (acessíveis majoritariamente via TUI):
  - Painel de Controle: Visão geral da agenda, tarefas urgentes e notificações.
  - Gestão de Tarefas: Crie, liste e marque tarefas como concluídas.
  - Gestão de Disciplinas e Turmas: Administre disciplinas, turmas e alunos.
  - Gestão de Aulas: Planeje e visualize aulas.
  - Gestão de Avaliações: Crie avaliações, lance notas e calcule médias.
  - Banco de Questões e Geração de Provas: Mantenha um banco de questões e gere provas.
  - Ferramentas de Produtividade: Como sessões de foco (funcionalidade futura).

Use "vigenda [comando] --help" para mais informações sobre um subcomando específico.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Inicia a aplicação TUI principal.
		// PersistentPreRunE garante que todos os serviços necessários já foram inicializados.
		// Os serviços são passados para a TUI para que ela possa interagir com a lógica de negócios.
		app.StartApp(taskService, classService, assessmentService, questionService, proofService, lessonService)
	},
	// PersistentPreRunE é executado antes do Run de qualquer comando (rootCmd ou subcomandos).
	// É usado aqui para garantir que o logging e a conexão com o banco de dados,
	// bem como a inicialização dos serviços, ocorram uma única vez.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Configura o logging para arquivo.
		// Se falhar, um aviso é impresso no stderr, mas a aplicação tenta continuar.
		if err := setupLogging(); err != nil {
			fmt.Fprintf(os.Stderr, "AVISO: Falha ao configurar o logging para arquivo: %v. Logando para stderr.\n", err)
		}

		// Inicializa a conexão com o banco de dados e os serviços, se ainda não o foram.
		if db == nil { // Garante a inicialização única.
			// Determina a configuração do banco de dados com base em variáveis de ambiente.
			// Prioriza VIGENDA_DB_DSN, depois VIGENDA_DB_TYPE, e então variáveis específicas
			// como VIGENDA_DB_PATH (para SQLite) ou VIGENDA_DB_HOST/USER/etc. (para PostgreSQL).
			dbType := os.Getenv("VIGENDA_DB_TYPE")
			dbDSN := os.Getenv("VIGENDA_DB_DSN")

			dbHost := os.Getenv("VIGENDA_DB_HOST")
			dbPort := os.Getenv("VIGENDA_DB_PORT")
			dbUser := os.Getenv("VIGENDA_DB_USER")
			dbPassword := os.Getenv("VIGENDA_DB_PASSWORD")
			dbName := os.Getenv("VIGENDA_DB_NAME")
			dbSSLMode := os.Getenv("VIGENDA_DB_SSLMODE")

			config := database.DBConfig{}
			if dbType == "" {
				dbType = "sqlite" // Padrão para SQLite.
			}
			config.DBType = dbType

			switch dbType {
			case "sqlite":
				if dbDSN != "" {
					config.DSN = dbDSN
				} else {
					sqlitePath := os.Getenv("VIGENDA_DB_PATH")
					if sqlitePath == "" {
						sqlitePath = database.DefaultSQLitePath() // Caminho padrão definido em internal/database.
					}
					config.DSN = sqlitePath
				}
			case "postgres":
				if dbDSN != "" {
					config.DSN = dbDSN
				} else {
					// Constrói DSN para PostgreSQL a partir de variáveis de ambiente individuais.
					// Validações para campos obrigatórios (usuário, nome do banco) são feitas.
					if dbHost == "" {
						dbHost = "localhost"
					}
					if dbPort == "" {
						dbPort = "5432"
					}
					if dbUser == "" {
						return fmt.Errorf("VIGENDA_DB_USER deve ser definida para conexão PostgreSQL")
					}
					if dbName == "" {
						return fmt.Errorf("VIGENDA_DB_NAME deve ser definida para conexão PostgreSQL")
					}
					if dbSSLMode == "" {
						dbSSLMode = "disable" // Padrão SSLMode.
					}
					config.DSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
						dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
				}
			default:
				return fmt.Errorf("tipo de banco de dados não suportado VIGENDA_DB_TYPE: %s. Tipos suportados: 'sqlite', 'postgres'", dbType)
			}

			var err error
			db, err = database.GetDBConnection(config) // Obtém a conexão com o DB.
			if err != nil {
				// Loga o erro fatal, pois a aplicação não pode funcionar sem DB.
				// O log irá para o arquivo, se configurado, e para stderr.
				log.Fatalf("CRÍTICO: Falha ao inicializar o banco de dados (tipo: %s): %v", config.DBType, err)
				// return fmt.Errorf(...) // log.Fatalf já encerra a aplicação.
			}
			initializeServices(db) // Inicializa todos os serviços da aplicação.
		}
		return nil
	},
}

// taskCmd é o comando pai para todas as operações relacionadas a tarefas.
// Ele agrupa subcomandos como 'add', 'listar', 'complete'.
var taskCmd = &cobra.Command{
	Use:   "tarefa",
	Short: "Gerencia tarefas (add, listar, complete)",
	Long: `O comando 'tarefa' permite gerenciar todas as suas atividades e pendências.
Você pode adicionar novas tarefas, listar tarefas existentes (filtrando por turma ou todas)
e marcar tarefas como concluídas. Muitas dessas funcionalidades também estão disponíveis
de forma mais interativa através da TUI principal (executando 'vigenda' sem subcomandos).`,
	Example: `  vigenda tarefa add "Preparar aula de Revolução Francesa" --classid 1 --duedate 2024-07-15
  vigenda tarefa listar --classid 1
  vigenda tarefa listar --all
  vigenda tarefa complete 5`,
}

// taskAddCmd define o subcomando 'vigenda tarefa add'.
// Permite adicionar uma nova tarefa com título, descrição opcional, ID de turma e data de vencimento.
// Se a descrição não for fornecida via flag, e a entrada for um TTY, ela será solicitada interativamente.
var taskAddCmd = &cobra.Command{
	Use:   "add [título]",
	Short: "Adiciona uma nova tarefa",
	Long: `Adiciona uma nova tarefa ao sistema.
O título da tarefa é obrigatório.
Você pode fornecer uma descrição detalhada, associar a tarefa a uma turma específica (usando --classid)
e definir um prazo de conclusão (usando --duedate no formato AAAA-MM-DD).
Se a descrição não for fornecida pela flag --description, ela será solicitada interativamente.`,
	Example: `  vigenda tarefa add "Corrigir provas bimestrais" --description "Corrigir as provas do 2º bimestre da turma 9A." --classid 1 --duedate 2024-07-20
  vigenda tarefa add "Planejar próxima unidade" --duedate 2024-08-01`,
	Args:  cobra.ExactArgs(1), // Requer exatamente um argumento (o título da tarefa).
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		description, _ := cmd.Flags().GetString("description")
		classIDStr, _ := cmd.Flags().GetString("classid")
		dueDateStr, _ := cmd.Flags().GetString("duedate")

		var classID *int64
		if classIDStr != "" {
			cid, err := strconv.ParseInt(classIDStr, 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao parsear ID da turma: %v\n", err)
				return
			}
			classID = &cid
		}

		var dueDate *time.Time
		if dueDateStr != "" {
			parsedDate, err := time.Parse("2006-01-02", dueDateStr) // Formato AAAA-MM-DD.
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao parsear data de conclusão (use o formato AAAA-MM-DD): %v\n", err)
				return
			}
			dueDate = &parsedDate
		}

		// Se a descrição não foi fornecida via flag, tenta obter interativamente.
		if description == "" {
			// Utiliza tui.GetInput, que lida com TTY e entrada redirecionada.
			desc, err := tui.GetInput("Digite a descrição da tarefa (opcional):", os.Stdout, os.Stdin)
			if err != nil {
				// Se GetInput retornar erro (ex: usuário cancelou), não prossegue.
				// Não é necessariamente um erro fatal para a aplicação, apenas para este comando.
				fmt.Fprintf(os.Stderr, "Falha ao obter descrição: %v\n", err)
				// Decide-se não prosseguir se a obtenção interativa falhar ou for cancelada.
				// Alternativamente, poderia prosseguir com descrição vazia se o erro não for crítico.
				// return // Comentado para permitir criação de tarefa sem descrição se o prompt falhar.
			}
			description = desc // Usa a descrição obtida, que pode ser vazia se o usuário não digitou nada.
		}

		// Chama o serviço para criar a tarefa.
		// O UserID é gerenciado internamente pelo serviço (atualmente fixo, mas deveria vir do contexto de auth).
		task, err := taskService.CreateTask(context.Background(), title, description, classID, dueDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar tarefa: %v\n", err)
			log.Printf("ERRO CMD: Falha ao criar tarefa '%s': %v", title, err) // Loga o erro também.
			return
		}
		fmt.Printf("Tarefa '%s' (ID: %d) criada com sucesso.\n", task.Title, task.ID)
		log.Printf("INFO CMD: Tarefa '%s' (ID: %d) criada.", task.Title, task.ID)
	},
}

// taskListCmd define o subcomando 'vigenda tarefa listar'.
// Lista tarefas ativas. Requer a flag --classid para filtrar por turma ou --all para listar todas as tarefas.
// A saída é formatada como uma tabela simples no console.
var taskListCmd = &cobra.Command{
	Use:   "listar",
	Short: "Lista tarefas ativas",
	Long:  `Lista tarefas ativas.
Use a flag --classid para filtrar tarefas de uma turma específica.
Use a flag --all para listar todas as tarefas ativas de todas as turmas e tarefas do sistema (bugs).
Uma dessas duas flags (--classid ou --all) é obrigatória.`,
	Example: `  vigenda tarefa listar --classid 1
  vigenda tarefa listar --all`,
	Run: func(cmd *cobra.Command, args []string) {
		classIDStr, _ := cmd.Flags().GetString("classid")
		showAllStr, _ := cmd.Flags().GetString("all")
		showAll, _ := strconv.ParseBool(showAllStr) // Converte para booleano, erro é ignorado (false por padrão).

		var tasks []models.Task
		var err error
		var headerMsg string

		if classIDStr != "" {
			classID, parseErr := strconv.ParseInt(classIDStr, 10, 64)
			if parseErr != nil {
				fmt.Fprintf(os.Stderr, "Erro ao parsear ID da turma: %v\n", parseErr)
				return
			}
			// Lista tarefas ativas para a turma especificada.
			tasks, err = taskService.ListActiveTasksByClass(context.Background(), classID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao listar tarefas da turma: %v\n", err)
				log.Printf("ERRO CMD: Falha ao listar tarefas da turma %d: %v", classID, err)
				return
			}
			// Tenta obter o nome da turma para o cabeçalho.
			class, classErr := classService.GetClassByID(context.Background(), classID)
			if classErr == nil && class.ID != 0 {
				headerMsg = fmt.Sprintf("TAREFAS ATIVAS PARA: %s (ID: %d)", class.Name, classID)
			} else {
				headerMsg = fmt.Sprintf("TAREFAS ATIVAS PARA: Turma ID %d", classID)
			}
		} else if showAll {
			// Lista todas as tarefas ativas (incluindo de sistema/bugs).
			tasks, err = taskService.ListAllActiveTasks(context.Background())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao listar todas as tarefas ativas: %v\n", err)
				log.Printf("ERRO CMD: Falha ao listar todas as tarefas ativas: %v", err)
				return
			}
			headerMsg = "TODAS AS TAREFAS ATIVAS (INCLUINDO BUGS DO SISTEMA)"
		} else {
			// Nenhuma flag válida fornecida.
			fmt.Fprintln(os.Stderr, "Erro: Especifique --classid <ID> para listar tarefas de uma turma OU use --all para listar todas as tarefas ativas.")
			fmt.Fprintln(os.Stderr, "Exemplo: vigenda tarefa listar --classid 1")
			fmt.Fprintln(os.Stderr, "Exemplo: vigenda tarefa listar --all")
			return
		}

		if len(tasks) == 0 {
			fmt.Println("Nenhuma tarefa ativa encontrada para os critérios fornecidos.")
			return
		}

		fmt.Printf("\n%s\n\n", headerMsg) // Adiciona nova linha antes do cabeçalho.

		// Define colunas para a tabela de saída.
		columns := []table.Column{
			{Title: "ID", Width: 4},
			{Title: "TAREFA (TÍTULO)", Width: 40},
			{Title: "DESCRIÇÃO", Width: 50},
			{Title: "PRAZO", Width: 12},
			{Title: "TURMA ID", Width: 10},
		}
		var rows []table.Row
		for _, task := range tasks {
			dueDateStr := "N/A"
			if task.DueDate != nil {
				dueDateStr = task.DueDate.Format("02/01/2006")
			}
			classIDDisplay := "N/A"
			if task.ClassID != nil {
				classIDDisplay = fmt.Sprintf("%d", *task.ClassID)
			}
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", task.ID),
				task.Title,
				task.Description, // Adicionada coluna de descrição
				dueDateStr,
				classIDDisplay,
			})
		}

		// Imprime a tabela manualmente para melhor controle de formatação e sem dependência de tui.ShowTable aqui.
		// Imprime cabeçalho.
		for i, col := range columns {
			fmt.Printf("%-*s", col.Width+2, col.Title) // +2 para padding e separador "|"
			if i < len(columns)-1 {
				fmt.Print("| ")
			}
		}
		fmt.Println()
		// Imprime linha separadora.
		for i, col := range columns {
			fmt.Printf("%s", strings.Repeat("-", col.Width+2))
			if i < len(columns)-1 {
				fmt.Print("+-"); // Separador mais robusto
			}
		}
		fmt.Println()
		// Imprime linhas de dados.
		for _, row := range rows {
			for i, cell := range row {
				// Trunca células se excederem a largura da coluna para evitar quebra de layout.
				content := fmt.Sprintf("%s", cell)
				if len(content) > columns[i].Width {
					content = content[:columns[i].Width-3] + "..."
				}
				fmt.Printf("%-*s", columns[i].Width+2, content)
				if i < len(columns)-1 {
					fmt.Print("| ")
				}
			}
			fmt.Println()
		}
	},
}

// taskCompleteCmd define o subcomando 'vigenda tarefa complete'.
// Marca uma tarefa existente como concluída usando seu ID.
var taskCompleteCmd = &cobra.Command{
	Use:     "complete [ID_da_tarefa]",
	Short:   "Marca uma tarefa como concluída",
	Long:    `Marca uma tarefa específica como concluída, utilizando o seu ID numérico.`,
	Example: `  vigenda tarefa complete 12`,
	Args:    cobra.ExactArgs(1), // Requer exatamente um argumento (ID da tarefa).
	Run: func(cmd *cobra.Command, args []string) {
		taskID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao parsear ID da tarefa: %v\n", err)
			return
		}
		err = taskService.MarkTaskAsCompleted(context.Background(), taskID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao marcar tarefa como concluída: %v\n", err)
			log.Printf("ERRO CMD: Falha ao completar tarefa ID %d: %v", taskID, err)
			return
		}
		fmt.Printf("Tarefa ID %d marcada como concluída.\n", taskID)
		log.Printf("INFO CMD: Tarefa ID %d completada.", taskID)
	},
}

// initializeServices configura as instâncias dos serviços da aplicação,
// injetando as dependências de repositório necessárias.
// É chamada após a conexão com o banco de dados ser estabelecida.
func initializeServices(db *sql.DB) {
	// Inicializa os repositórios concretos com a conexão db.
	taskRepo := repository.NewTaskRepository(db)
	classRepo := repository.NewClassRepository(db)
	assessmentRepo := repository.NewAssessmentRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	subjectRepo := repository.NewSubjectRepository(db) // Assumindo que existe e é necessário.
	lessonRepo := repository.NewLessonRepository(db)   // Assumindo que existe.

	// Inicializa os serviços com suas implementações de repositório.
	taskService = service.NewTaskService(taskRepo)
	classService = service.NewClassService(classRepo, subjectRepo) // SubjectRepo pode ser necessário para validações.
	assessmentService = service.NewAssessmentService(assessmentRepo, classRepo) // ClassRepo pode ser usado para buscar alunos.
	questionService = service.NewQuestionService(questionRepo, subjectRepo)
	proofService = service.NewProofService(questionRepo)
	lessonService = service.NewLessonService(lessonRepo, classRepo) // ClassRepo para validação de propriedade da turma.
	log.Println("INFO CMD: Todos os serviços foram inicializados.")
}

// lessonService é uma variável global para o serviço de aulas,
// declarada separadamente para ser acessível pelo Run do rootCmd
// ao iniciar a TUI principal.
var lessonService service.LessonService

// init configura todos os comandos e flags da CLI usando Cobra.
// Esta função é chamada automaticamente pelo Go na inicialização do pacote.
func init() {
	// Configuração das flags para 'tarefa add'.
	taskAddCmd.Flags().StringP("description", "d", "", "Descrição detalhada da tarefa.")
	taskAddCmd.Flags().String("classid", "", "ID da turma para associar a tarefa (opcional).")
	taskAddCmd.Flags().String("duedate", "", "Data de conclusão da tarefa no formato AAAA-MM-DD (opcional).")

	// Configuração das flags para 'tarefa listar'.
	taskListCmd.Flags().String("classid", "", "ID da turma para filtrar as tarefas.")
	taskListCmd.Flags().Bool("all", false, "Listar todas as tarefas ativas (ignora --classid se presente).")

	// Adiciona subcomandos ao comando 'tarefa'.
	taskCmd.AddCommand(taskAddCmd, taskListCmd, taskCompleteCmd)
	// Adiciona o comando 'tarefa' ao comando raiz.
	rootCmd.AddCommand(taskCmd)

	// Configuração e adição do comando 'turma' e seus subcomandos.
	// O comando 'turma criar' foi removido da CLI, pois a criação de turmas
	// agora é primariamente feita via TUI principal.
	classCmd.AddCommand(classImportStudentsCmd, classUpdateStudentStatusCmd)
	rootCmd.AddCommand(classCmd)

	// Configuração das flags para 'avaliacao criar'.
	assessmentCreateCmd.Flags().String("classid", "", "ID da turma para a qual a avaliação será criada (obrigatório).")
	_ = assessmentCreateCmd.MarkFlagRequired("classid") // Marca a flag como obrigatória.
	assessmentCreateCmd.Flags().String("term", "", "Período/bimestre da avaliação (ex: 1, 2) (obrigatório).")
	_ = assessmentCreateCmd.MarkFlagRequired("term")
	assessmentCreateCmd.Flags().String("weight", "", "Peso da avaliação na média final (ex: 4.0) (obrigatório).")
	_ = assessmentCreateCmd.MarkFlagRequired("weight")

	// Adiciona subcomandos ao comando 'avaliacao'.
	assessmentCmd.AddCommand(assessmentCreateCmd, assessmentEnterGradesCmd, assessmentClassAverageCmd)
	rootCmd.AddCommand(assessmentCmd)

	// Configuração e adição do comando 'bancoq' (banco de questões).
	questionBankCmd.AddCommand(questionBankAddCmd)
	rootCmd.AddCommand(questionBankCmd)

	// Configuração das flags e adição do comando 'prova'.
	proofGenerateCmd.Flags().String("subjectid", "", "ID da disciplina para gerar a prova (obrigatório).")
	_ = proofGenerateCmd.MarkFlagRequired("subjectid")
	proofGenerateCmd.Flags().String("topic", "", "Tópico específico para filtrar questões (opcional).")
	proofGenerateCmd.Flags().String("easy", "0", "Número de questões fáceis.")
	proofGenerateCmd.Flags().String("medium", "0", "Número de questões médias.")
	proofGenerateCmd.Flags().String("hard", "0", "Número de questões difíceis.")
	proofCmd.AddCommand(proofGenerateCmd)
	rootCmd.AddCommand(proofCmd)
}

// setupLogging configura o sistema de logging para escrever em um arquivo.
// O arquivo de log é nomeado 'vigenda.log' e é colocado no diretório de configuração
// do usuário (ex: ~/.config/vigenda/ no Linux) ou no diretório de trabalho atual
// como fallback.
// Retorna um erro se a configuração do log falhar.
func setupLogging() error {
	logDir := ""
	userConfigDir, err := os.UserConfigDir()
	if err == nil {
		logDir = filepath.Join(userConfigDir, "vigenda")
	} else {
		// Se não conseguir obter o diretório de config do usuário, tenta o diretório atual.
		cwd, errCwd := os.Getwd()
		if errCwd == nil {
			logDir = cwd // Usa o diretório atual como logDir.
		} else {
			// Em caso de falha total em determinar um diretório, loga um aviso.
			// O log ainda será configurado para stderr pelo pacote log padrão.
			log.Printf("AVISO: Não foi possível determinar o diretório de configuração do usuário ou o diretório de trabalho atual para logs. Tentando logar no diretório atual se possível, ou stderr.")
			// Não retorna erro aqui, pois o log padrão para stderr ainda funcionará.
		}
	}

	// Se um logDir foi determinado (mesmo que seja o diretório atual),
	// tenta criar o subdiretório 'vigenda' se não for o diretório atual direto.
	// Isso é para o caso de userConfigDir ter sucesso, mas o subdiretório 'vigenda' não existir.
	if logDir != "" && logDir != "." && logDir != mustGetwd() {
		// Tenta criar o diretório de log (ex: ~/.config/vigenda/).
		if err := os.MkdirAll(logDir, 0755); err != nil {
			// Se não conseguir criar o diretório específico, tenta logar no diretório atual como último recurso.
			log.Printf("AVISO: Não foi possível criar o diretório de log %s: %v. Tentando logar no diretório de trabalho atual.", logDir, err)
			logDir = "." // Define para o diretório de trabalho atual.
		}
	}
	// Se logDir ainda estiver vazio (caso extremo), define para diretório atual.
	if logDir == "" {
		logDir = "."
	}

	logFilePath := filepath.Join(logDir, "vigenda.log")

	// Abre o arquivo de log. Cria se não existir, anexa se existir.
	var errOpen error
	logFile, errOpen = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errOpen != nil {
		return fmt.Errorf("falha ao abrir arquivo de log %s: %w", logFilePath, errOpen)
	}

	// Configura a saída do log para o arquivo.
	log.SetOutput(logFile)
	// Adiciona flags para incluir data, hora, microssegundos e arquivo:linha no log.
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	log.Println("INFO: Logging inicializado para arquivo:", logFilePath)
	return nil
}

// mustGetwd é uma função auxiliar para obter o diretório de trabalho atual (CWD).
// Retorna "." em caso de erro, simplificando a lógica de fallback para o logging.
func mustGetwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		// Em um cenário real, um tratamento mais robusto pode ser necessário.
		// Para a lógica de fallback do log, se CWD falhar, "." é um fallback razoável.
		return "."
	}
	return cwd
}

// classCmd é o comando pai para todas as operações relacionadas a turmas e alunos.
// A criação de turmas é primariamente feita via TUI.
var classCmd = &cobra.Command{
	Use:   "turma",
	Short: "Gerencia turmas e alunos (importar-alunos, atualizar-status)",
	Long: `O comando 'turma' é usado para administrar funcionalidades relacionadas a turmas,
como a importação de listas de alunos de ficheiros CSV e a atualização do status de alunos.
A criação e edição detalhada de turmas é feita através da interface interativa principal
(executando 'vigenda' sem subcomandos).`,
	Example: `  vigenda turma importar-alunos 1 alunos_9a.csv
  vigenda turma atualizar-status 15 transferido`,
}

// classImportStudentsCmd define o subcomando 'vigenda turma importar-alunos'.
// Importa uma lista de alunos de um arquivo CSV para uma turma existente.
var classImportStudentsCmd = &cobra.Command{
	Use:   "importar-alunos [ID_da_turma] [caminho_do_ficheiro_csv]",
	Short: "Importa alunos de um ficheiro CSV para uma turma",
	Long: `Importa uma lista de alunos de um ficheiro CSV para uma turma existente.
A turma deve ser criada previamente através da TUI principal.
O ficheiro CSV deve conter as colunas 'numero_chamada' (opcional), 'nome_completo' (obrigatório),
e 'situacao' (opcional; padrões para 'ativo'). Consulte a documentação para a estrutura detalhada.`,
	Example: `  vigenda turma importar-alunos 1 ./lista_alunos_turma_a.csv
  vigenda turma importar-alunos 3 /documentos/alunos_turma_c.csv`,
	Args:  cobra.ExactArgs(2), // Requer ID da turma e caminho do CSV.
	Run: func(cmd *cobra.Command, args []string) {
		classID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao parsear ID da turma: %v\n", err)
			return
		}
		csvFilePath := args[1]

		csvData, err := os.ReadFile(csvFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler arquivo CSV '%s': %v\n", csvFilePath, err)
			return
		}

		count, err := classService.ImportStudentsFromCSV(context.Background(), classID, csvData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao importar alunos: %v\n", err)
			log.Printf("ERRO CMD: Falha ao importar alunos para turma ID %d: %v", classID, err)
			return
		}
		fmt.Printf("%d aluno(s) importado(s) com sucesso para a turma ID %d.\n", count, classID)
		log.Printf("INFO CMD: %d aluno(s) importado(s) para turma ID %d.", count, classID)
	},
}

// classUpdateStudentStatusCmd define o subcomando 'vigenda turma atualizar-status'.
// Atualiza o status de um aluno existente (ex: 'ativo', 'inativo', 'transferido').
var classUpdateStudentStatusCmd = &cobra.Command{
	Use:   "atualizar-status [ID_do_aluno] [novo_status]",
	Short: "Atualiza o status de um aluno",
	Long: `Atualiza o status de um aluno específico (ex: 'ativo', 'inativo', 'transferido').
O ID do aluno é o identificador numérico único na base de dados.
Status permitidos: 'ativo', 'inativo', 'transferido'.`,
	Example: `  vigenda turma atualizar-status 25 ativo
  vigenda turma atualizar-status 103 transferido`,
	Args:  cobra.ExactArgs(2), // Requer ID do aluno e novo status.
	Run: func(cmd *cobra.Command, args []string) {
		studentID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao parsear ID do aluno: %v\n", err)
			return
		}
		newStatus := args[1]
		// A validação do valor de newStatus é feita na camada de serviço.

		err = classService.UpdateStudentStatus(context.Background(), studentID, newStatus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao atualizar status do aluno: %v\n", err)
			log.Printf("ERRO CMD: Falha ao atualizar status do aluno ID %d para '%s': %v", studentID, newStatus, err)
			return
		}
		fmt.Printf("Status do aluno ID %d atualizado para '%s'.\n", studentID, newStatus)
		log.Printf("INFO CMD: Status do aluno ID %d atualizado para '%s'.", studentID, newStatus)
	},
}

// assessmentCmd é o comando pai para operações relacionadas a avaliações.
var assessmentCmd = &cobra.Command{
	Use:   "avaliacao",
	Short: "Gerencia avaliações e notas (criar, lancar-notas, media-turma)",
	Long: `O comando 'avaliacao' permite gerenciar o ciclo de vida das avaliações,
desde a sua criação para uma turma existente, o lançamento interativo de notas,
até o cálculo da média da turma. A criação de turmas é feita via TUI.`,
	Example: `  vigenda avaliacao criar "Prova Bimestral 1" --classid 1 --term 1 --weight 4.0
  vigenda avaliacao lancar-notas 1
  vigenda avaliacao media-turma 1`,
}

// assessmentCreateCmd define o subcomando 'vigenda avaliacao criar'.
// Cria uma nova avaliação para uma turma, especificando nome, período e peso.
var assessmentCreateCmd = &cobra.Command{
	Use:   "criar [nome_da_avaliacao]",
	Short: "Cria uma nova avaliação para uma turma",
	Long: `Cria uma nova avaliação associada a uma turma específica.
A turma deve existir e ser identificada pela flag --classid.
É necessário fornecer o nome da avaliação, o período/bimestre (--term) e o peso (--weight).
Se alguma flag obrigatória não for fornecida, será solicitada interativamente.`,
	Example: `  vigenda avaliacao criar "Trabalho de História Moderna" --classid 2 --term 3 --weight 3.5
  vigenda avaliacao criar "Seminário de Literatura" --classid 1 --term 2 --weight 2.0`,
	Args:  cobra.ExactArgs(1), // Requer o nome da avaliação como argumento.
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		classIDStr, _ := cmd.Flags().GetString("classid")
		termStr, _ := cmd.Flags().GetString("term")
		weightStr, _ := cmd.Flags().GetString("weight")

		var err error // Declarar err aqui para ser acessível em todos os blocos if.

		// Solicita interativamente as flags obrigatórias se não fornecidas.
		if classIDStr == "" {
			classIDStr, err = tui.GetInput("Digite o ID da Turma para a avaliação:", os.Stdout, os.Stdin)
			if err != nil || classIDStr == "" {
				fmt.Fprintln(os.Stderr, "ID da Turma é obrigatório.")
				return
			}
		}
		if termStr == "" {
			termStr, err = tui.GetInput("Digite o Período/Bimestre (ex: 1, 2) para a avaliação:", os.Stdout, os.Stdin)
			if err != nil || termStr == "" {
				fmt.Fprintln(os.Stderr, "Período/Bimestre é obrigatório.")
				return
			}
		}
		if weightStr == "" {
			weightStr, err = tui.GetInput("Digite o Peso (ex: 4.0) para a avaliação:", os.Stdout, os.Stdin)
			if err != nil || weightStr == "" {
				fmt.Fprintln(os.Stderr, "Peso é obrigatório.")
				return
			}
		}

		classID, errConv := strconv.ParseInt(classIDStr, 10, 64)
		if errConv != nil {
			fmt.Fprintf(os.Stderr, "Erro ao parsear ID da Turma: %v\n", errConv)
			return
		}
		term, errConv := strconv.Atoi(termStr)
		if errConv != nil {
			fmt.Fprintf(os.Stderr, "Erro ao parsear Período/Bimestre: %v\n", errConv)
			return
		}
		weight, errConv := strconv.ParseFloat(weightStr, 64)
		if errConv != nil {
			fmt.Fprintf(os.Stderr, "Erro ao parsear Peso: %v\n", errConv)
			return
		}

		assessment, errService := assessmentService.CreateAssessment(context.Background(), name, classID, term, weight)
		if errService != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar avaliação: %v\n", errService)
			log.Printf("ERRO CMD: Falha ao criar avaliação '%s': %v", name, errService)
			return
		}
		fmt.Printf("Avaliação '%s' (ID: %d) criada para Turma ID %d, Período %d, Peso %.1f.\n", assessment.Name, assessment.ID, classID, term, weight)
		log.Printf("INFO CMD: Avaliação '%s' (ID: %d) criada.", assessment.Name, assessment.ID)
	},
}

// assessmentEnterGradesCmd define o subcomando 'vigenda avaliacao lancar-notas'.
// Inicia um processo interativo para lançar notas para uma avaliação específica.
// TODO: A implementação atual da Run deste comando é um placeholder e simula a entrada.
//       Uma TUI mais completa para lançamento de notas seria ideal, possivelmente
//       integrada ao módulo `internal/app/assessments`.
var assessmentEnterGradesCmd = &cobra.Command{
	Use:   "lancar-notas [ID_da_avaliacao]",
	Short: "Lança notas para os alunos de uma avaliação",
	Long: `Inicia um processo interativo para lançar ou editar as notas dos alunos
para uma avaliação específica. A lista de alunos da turma associada à avaliação
será exibida, permitindo a inserção de cada nota.
O ID da avaliação é o identificador numérico único da avaliação.
NOTA: A interatividade deste comando CLI é básica. Para uma experiência completa,
use a funcionalidade de lançamento de notas na TUI principal.`,
	Example: `  vigenda avaliacao lancar-notas 7`,
	Args:  cobra.ExactArgs(1), // Requer o ID da avaliação.
	Run: func(cmd *cobra.Command, args []string) {
		assessmentID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao parsear ID da Avaliação: %v\n", err)
			return
		}

		// A implementação atual é um placeholder e simula a entrada de notas.
		// Uma TUI dedicada (possivelmente em `internal/app/assessments`) seria melhor.
		fmt.Printf("Entrada de notas interativa para Avaliação ID %d (CLI placeholder).\n", assessmentID)
		fmt.Println("Para uma experiência completa, use a TUI principal.")
		fmt.Println("Simulando entrada de notas (digite StudentID:Nota, ou 'done' para finalizar):")

		studentGrades := make(map[int64]float64)
		for {
			input, _ := tui.GetInput("Digite IDdoAluno:Nota (ou 'done'):", os.Stdout, os.Stdin)
			if strings.ToLower(input) == "done" {
				break
			}
			parts := strings.Split(input, ":")
			if len(parts) != 2 {
				fmt.Fprintln(os.Stderr, "Formato inválido. Use IDdoAluno:Nota.")
				continue
			}
			studentID, errS := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
			grade, errG := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if errS != nil || errG != nil {
				fmt.Fprintln(os.Stderr, "ID do Aluno ou Nota inválido(a).")
				continue
			}
			studentGrades[studentID] = grade
		}

		if len(studentGrades) > 0 {
			err = assessmentService.EnterGrades(context.Background(), assessmentID, studentGrades)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao lançar notas: %v\n", err)
				log.Printf("ERRO CMD: Falha ao lançar notas para Avaliação ID %d: %v", assessmentID, err)
				return
			}
			fmt.Printf("Notas lançadas com sucesso para Avaliação ID %d.\n", assessmentID)
			log.Printf("INFO CMD: Notas lançadas para Avaliação ID %d.", assessmentID)
		} else {
			fmt.Println("Nenhuma nota foi lançada.")
		}
	},
}

// assessmentClassAverageCmd define o subcomando 'vigenda avaliacao media-turma'.
// Calcula e exibe a média geral ponderada das notas para uma turma específica.
var assessmentClassAverageCmd = &cobra.Command{
	Use:   "media-turma [ID_da_turma]",
	Short: "Calcula a média geral das notas de uma turma",
	Long: `Calcula e exibe a média geral ponderada das notas para uma turma específica,
considerando todas as avaliações e seus respectivos pesos para aquela turma.
O ID da turma é o identificador numérico único da turma.`,
	Example: `  vigenda avaliacao media-turma 1`,
	Args:  cobra.ExactArgs(1), // Requer o ID da turma.
	Run: func(cmd *cobra.Command, args []string) {
		classID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao parsear ID da Turma: %v\n", err)
			return
		}

		average, err := assessmentService.CalculateClassAverage(context.Background(), classID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao calcular média da turma: %v\n", err)
			log.Printf("ERRO CMD: Falha ao calcular média para Turma ID %d: %v", classID, err)
			return
		}
		fmt.Printf("Média geral para Turma ID %d: %.2f\n", classID, average)
	},
}

// questionBankCmd é o comando pai para operações relacionadas ao banco de questões.
var questionBankCmd = &cobra.Command{
	Use:   "bancoq",
	Short: "Gerencia o banco de questões (add)",
	Long: `O comando 'bancoq' (Banco de Questões) permite adicionar novas questões ao sistema
a partir de um ficheiro JSON formatado. A criação de disciplinas (às quais as questões
são associadas) é feita via TUI principal.`,
	Example: `  vigenda bancoq add ./minhas_questoes_historia.json`,
}

// questionBankAddCmd define o subcomando 'vigenda bancoq add'.
// Adiciona questões de um arquivo JSON para o banco de questões.
var questionBankAddCmd = &cobra.Command{
	Use:   "add [caminho_do_ficheiro_json]",
	Short: "Adiciona questões de um ficheiro JSON ao banco",
	Long: `Adiciona um conjunto de questões de um ficheiro JSON para o banco de questões central.
O ficheiro JSON deve seguir uma estrutura específica, e as disciplinas referenciadas
nas questões devem existir (criadas via TUI). Consulte a documentação para o formato do JSON.`,
	Example: `  vigenda bancoq add questoes_bimestre1.json`,
	Args:  cobra.ExactArgs(1), // Requer o caminho do arquivo JSON.
	Run: func(cmd *cobra.Command, args []string) {
		jsonFilePath := args[0]
		jsonData, err := os.ReadFile(jsonFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler arquivo JSON '%s': %v\n", jsonFilePath, err)
			return
		}

		count, err := questionService.AddQuestionsFromJSON(context.Background(), jsonData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao adicionar questões do JSON: %v\n", err)
			log.Printf("ERRO CMD: Falha ao adicionar questões do JSON '%s': %v", jsonFilePath, err)
			return
		}
		fmt.Printf("%d questão(ões) adicionada(s) com sucesso ao banco.\n", count)
		log.Printf("INFO CMD: %d questão(ões) adicionada(s) do JSON '%s'.", count, jsonFilePath)
	},
}

// proofCmd é o comando pai para operações relacionadas à geração de provas.
var proofCmd = &cobra.Command{
	Use:   "prova",
	Short: "Gerencia e gera provas (gerar)",
	Long: `O comando 'prova' permite gerar provas textuais a partir do banco de questões.
Você pode especificar critérios como disciplina, tópico e o número desejado de questões
por nível de dificuldade. Disciplinas devem ser criadas via TUI.`,
	Example: `  vigenda prova gerar --subjectid 1 --easy 5 --medium 3 --hard 2`,
}

// proofGenerateCmd define o subcomando 'vigenda prova gerar'.
// Gera uma nova prova com base em critérios como disciplina, tópico e contagem de dificuldades.
var proofGenerateCmd = &cobra.Command{
	Use:   "gerar",
	Short: "Gera uma nova prova com base em critérios especificados",
	Long: `Gera uma prova selecionando questões do banco de questões.
É obrigatório especificar o ID da disciplina (--subjectid).
Opcionalmente, pode-se filtrar por tópico (--topic) e definir o número de questões
para cada nível de dificuldade (--easy, --medium, --hard).
A prova gerada será exibida no terminal ou pode ser salva em um arquivo com --output.`,
	Example: `  vigenda prova gerar --subjectid 1 --easy 5 --medium 3 --hard 2 --topic "Revolução Industrial" --output prova.txt
  vigenda prova gerar --subjectid 3 --medium 10 --hard 5`,
	Run: func(cmd *cobra.Command, args []string) {
		subjectIDStr, _ := cmd.Flags().GetString("subjectid")
		topic, _ := cmd.Flags().GetString("topic")
		easyCountStr, _ := cmd.Flags().GetString("easy")
		mediumCountStr, _ := cmd.Flags().GetString("medium")
		hardCountStr, _ := cmd.Flags().GetString("hard")
		// outputFilePath, _ := cmd.Flags().GetString("output") // Descomentar se for usar

		if subjectIDStr == "" {
			fmt.Fprintln(os.Stderr, "Erro: ID da disciplina (--subjectid) é obrigatório.")
			return
		}
		subjectID, err := strconv.ParseInt(subjectIDStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro: ID da disciplina inválido: %v\n", err)
			return
		}

		easyCount, _ := strconv.Atoi(easyCountStr)
		mediumCount, _ := strconv.Atoi(mediumCountStr)
		hardCount, _ := strconv.Atoi(hardCountStr)

		if easyCount == 0 && mediumCount == 0 && hardCount == 0 {
			fmt.Fprintln(os.Stderr, "Erro: Pelo menos uma contagem de dificuldade (--easy, --medium, --hard) deve ser maior que zero.")
			return
		}

		criteria := service.ProofCriteria{
			SubjectID:   subjectID,
			EasyCount:   easyCount,
			MediumCount: mediumCount,
			HardCount:   hardCount,
		}
		if topic != "" {
			criteria.Topic = &topic // Define o tópico se fornecido.
		}

		questions, err := proofService.GenerateProof(context.Background(), criteria)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao gerar prova: %v\n", err)
			log.Printf("ERRO CMD: Falha ao gerar prova para Subject ID %d: %v", subjectID, err)
			return
		}

		if len(questions) == 0 {
			fmt.Println("Nenhuma questão encontrada para os critérios fornecidos. A prova não pôde ser gerada.")
			return
		}

		// TODO: Implementar salvamento em arquivo se outputFilePath for fornecido.
		//       Por enquanto, imprime no console.
		fmt.Printf("\n--- Prova Gerada (%d Questões) ---\n\n", len(questions))
		for i, q := range questions {
			fmt.Printf("**Questão %d (%s, %s)**\n", i+1, q.Difficulty, q.Type)
			if q.Topic != "" {
				fmt.Printf("Tópico: %s\n", q.Topic)
			}
			fmt.Printf("%s\n", q.Statement)
			if q.Options != nil && *q.Options != "" && *q.Options != "null" {
				var opts []string
				if json.Unmarshal([]byte(*q.Options), &opts) == nil {
					for j, opt := range opts {
						fmt.Printf("  %c) %s\n", 'a'+j, opt)
					}
				}
			}
			fmt.Printf("   Resposta Correta: %s\n\n", q.CorrectAnswer)
		}
		fmt.Println("--- Fim da Prova ---")
		log.Printf("INFO CMD: Prova gerada com %d questões para Subject ID %d.", len(questions), subjectID)
	},
}

// main é a função principal da aplicação.
// Ela configura e executa o comando raiz do Cobra.
// Garante que o arquivo de log seja fechado corretamente ao final da execução.
func main() {
	// PersistentPreRunE (no rootCmd) já chama setupLogging.
	// É crucial fechar logFile no final.
	if err := rootCmd.Execute(); err != nil {
		// Erros de execução de comando Cobra são geralmente impressos pelo Cobra.
		// Logar adicionalmente aqui pode ser redundante se o Cobra já o faz,
		// mas garante que vá para o arquivo de log.
		log.Printf("ERRO FATAL: Falha ao executar rootCmd: %v", err)
		// Cobra já imprime o erro no Stderr, então não precisamos fazer isso aqui.
		// fmt.Fprintln(os.Stderr, "Erro ao executar comando:", err)
		if logFile != nil {
			logFile.Close()
		}
		os.Exit(1) // Sai com código de erro.
	}

	// Se Execute() for bem-sucedido e a aplicação terminar normalmente.
	if logFile != nil {
		log.Println("INFO: Aplicação finalizada com sucesso. Fechando arquivo de log.")
		logFile.Close()
	}
}

[end of cmd/vigenda/main.go]
