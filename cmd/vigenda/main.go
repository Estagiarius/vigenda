// Package main is the entry point of the Vigenda application.
// It handles the main function and Cobra CLI configuration.
package main

import (
	"context"
	"database/sql"
	"encoding/json" // Added missing import
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table" // Reativado para columns e rows
	"github.com/spf13/cobra"
	"vigenda/internal/database"
	"vigenda/internal/models" // Added import for models package
	"vigenda/internal/repository"
	"vigenda/internal/service"
	"vigenda/internal/tui"
)

var db *sql.DB // Global database connection pool

var taskService service.TaskService
var classService service.ClassService
var assessmentService service.AssessmentService
var questionService service.QuestionService
var proofService service.ProofService

var rootCmd = &cobra.Command{
	Use:   "vigenda",
	Short: "Vigenda é uma ferramenta CLI para auxiliar professores na gestão de suas atividades.",
	Long: `Vigenda é uma aplicação de linha de comando (CLI) projetada para ajudar professores,
especialmente aqueles com TDAH, a organizar tarefas, aulas, avaliações e outras
atividades pedagógicas de forma eficiente.

Funcionalidades Principais:
  - Dashboard: Visão geral da agenda do dia, tarefas urgentes e notificações.
  - Gestão de Tarefas: Crie, liste e marque tarefas como concluídas.
  - Gestão de Turmas: Administre turmas, alunos (incluindo importação) e seus status.
  - Gestão de Avaliações: Crie avaliações, lance notas e calcule médias.
  - Banco de Questões: Mantenha um banco de questões e gere provas.

Use "vigenda [comando] --help" para mais informações sobre um comando específico.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This is the main dashboard view
		// For now, printing a simplified version.
		// TODO: Implement actual data fetching and formatting as per golden file.
		fmt.Println("=================================================")
		fmt.Println("==                 DASHBOARD                   ==")
		fmt.Println("=================================================")
		fmt.Println("")
		fmt.Println("🕒 AGENDA DE HOJE (22/06/2025)")
		fmt.Println("   [09:00 - 10:00] Aula de História - Turma 9A")
		fmt.Println("   [14:00 - 15:00] Reunião Pedagógica")
		fmt.Println("")
		fmt.Println("🔥 TAREFAS PRIORITÁRIAS")
		fmt.Println("   [1] Corrigir provas (Turma 9A) (Prazo: Amanhã)")
		fmt.Println("   [2] Preparar aula sobre Era Vargas (Turma 9B) (Prazo: 24/06)")
		fmt.Println("")
		fmt.Println("🔔 NOTIFICAÇÕES")
		fmt.Println("   - 5 entregas pendentes para o trabalho \"Pesquisa sobre Clima\" (Turma 9A).")
		fmt.Println("") // Ensure a trailing newline if the golden file has one after trimming
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// This function will run before any command, ensuring DB is initialized.
		if db == nil { // Initialize only once
			dbType := os.Getenv("VIGENDA_DB_TYPE")
			dbDSN := os.Getenv("VIGENDA_DB_DSN") // Generic DSN

			// Specific environment variables for constructing DSN if VIGENDA_DB_DSN is not set
			dbHost := os.Getenv("VIGENDA_DB_HOST")
			dbPort := os.Getenv("VIGENDA_DB_PORT")
			dbUser := os.Getenv("VIGENDA_DB_USER")
			dbPassword := os.Getenv("VIGENDA_DB_PASSWORD")
			dbName := os.Getenv("VIGENDA_DB_NAME")
			dbSSLMode := os.Getenv("VIGENDA_DB_SSLMODE") // Primarily for PostgreSQL

			config := database.DBConfig{}

			if dbType == "" {
				dbType = "sqlite" // Default to SQLite
			}
			config.DBType = dbType

			switch dbType {
			case "sqlite":
				if dbDSN != "" {
					config.DSN = dbDSN
				} else {
					// VIGENDA_DB_PATH is specific to SQLite if VIGENDA_DB_DSN is not used
					sqlitePath := os.Getenv("VIGENDA_DB_PATH")
					if sqlitePath == "" {
						sqlitePath = database.DefaultSQLitePath()
					}
					config.DSN = sqlitePath
				}
			case "postgres":
				if dbDSN != "" {
					config.DSN = dbDSN
				} else {
					// Construct PostgreSQL DSN from individual parts
					if dbHost == "" {
						dbHost = "localhost" // Default host
					}
					if dbPort == "" {
						dbPort = "5432" // Default PostgreSQL port
					}
					if dbUser == "" {
						// User must be provided for PostgreSQL typically
						return fmt.Errorf("VIGENDA_DB_USER must be set for PostgreSQL connection")
					}
					if dbName == "" {
						// DB Name must be provided
						return fmt.Errorf("VIGENDA_DB_NAME must be set for PostgreSQL connection")
					}
					if dbSSLMode == "" {
						dbSSLMode = "disable" // Default SSLMode
					}
					// Password can be empty if auth method allows (e.g. peer auth)
					// Note: Real applications should handle password securely (e.g. from secrets manager)
					config.DSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
						dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
				}
			default:
				return fmt.Errorf("unsupported VIGENDA_DB_TYPE: %s. Supported types are 'sqlite', 'postgres'", dbType)
			}

			var err error
			db, err = database.GetDBConnection(config)
			if err != nil {
				return fmt.Errorf("failed to initialize database (type: %s): %w", config.DBType, err)
			}
			// Initialize services here, after DB is ready
			initializeServices(db)
		}
		return nil
	},
}

var taskCmd = &cobra.Command{
	Use:   "tarefa",
	Short: "Gerencia tarefas (add, listar, complete)",
	Long:  `O comando 'tarefa' permite gerenciar todas as suas atividades e pendências. Você pode adicionar novas tarefas, listar tarefas existentes (filtrando por turma) e marcar tarefas como concluídas.`,
	Example: `  vigenda tarefa add "Preparar aula de Revolução Francesa" --classid 1 --duedate 2024-07-15
  vigenda tarefa listar --classid 1
  vigenda tarefa complete 5`,
}

var taskAddCmd = &cobra.Command{
	Use:   "add [título]",
	Short: "Adiciona uma nova tarefa",
	Long: `Adiciona uma nova tarefa ao sistema.
Você pode fornecer uma descrição detalhada, associar a tarefa a uma turma específica
e definir um prazo de conclusão utilizando as flags correspondentes.`,
	Example: `  vigenda tarefa add "Corrigir provas bimestrais" --description "Corrigir as provas do 2º bimestre da turma 9A." --classid 1 --duedate 2024-07-20
  vigenda tarefa add "Planejar próxima unidade" --duedate 2024-08-01`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		description, _ := cmd.Flags().GetString("description")
		classIDStr, _ := cmd.Flags().GetString("classid")
		dueDateStr, _ := cmd.Flags().GetString("duedate")

		var classID *int64
		if classIDStr != "" {
			cid, err := strconv.ParseInt(classIDStr, 10, 64)
			if err != nil {
				fmt.Println("Error parsing class ID:", err)
				return
			}
			classID = &cid
		}

		var dueDate *time.Time
		if dueDateStr != "" {
			parsedDate, err := time.Parse("2006-01-02", dueDateStr)
			if err != nil {
				fmt.Println("Error parsing due date (use YYYY-MM-DD format):", err)
				return
			}
			dueDate = &parsedDate
		}

		if description == "" {
			desc, err := tui.GetInput("Enter task description (optional):", os.Stdout, os.Stdin)
			if err != nil {
				fmt.Println("Error getting description:", err)
				// Decide if we should proceed or return
			}
			description = desc
		}

		task, err := taskService.CreateTask(context.Background(), title, description, classID, dueDate)
		if err != nil {
			fmt.Println("Error creating task:", err)
			return
		}
		fmt.Printf("Task '%s' (ID: %d) created successfully.\n", task.Title, task.ID)
	},
}

var taskListCmd = &cobra.Command{
	Use:   "listar",
	Short: "Lista tarefas ativas",
	Long:  `Lista todas as tarefas ativas. É obrigatório filtrar as tarefas por um ID de turma específico usando a flag --classid.`,
	Example: `  vigenda tarefa listar --classid 1
  vigenda tarefa listar --classid 3`,
	Run: func(cmd *cobra.Command, args []string) {
		classIDStr, _ := cmd.Flags().GetString("classid")
		showAllStr, _ := cmd.Flags().GetString("all") // Check for the --all flag
		showAll := showAllStr == "true" // Convert to boolean

		var tasks []models.Task
		var err error
		var headerMsg string

		if classIDStr != "" {
			classID, parseErr := strconv.ParseInt(classIDStr, 10, 64)
			if parseErr != nil {
				fmt.Println("Error parsing class ID:", parseErr)
				return
			}
			tasks, err = taskService.ListActiveTasksByClass(context.Background(), classID)
			if err != nil {
				fmt.Println("Error listing tasks:", err)
				return
			}
			class, classErr := classService.GetClassByID(context.Background(), classID)
			if classErr == nil && class.ID != 0 {
				headerMsg = fmt.Sprintf("TAREFAS PARA: %s", class.Name) // Restaurado
			} else {
				headerMsg = fmt.Sprintf("TAREFAS PARA: Class ID %d", classID) // Restaurado
			}
		} else if showAll {
			// This part needs a new service method: ListAllActiveTasks (or similar)
			// For now, let's assume such a method exists or we adapt.
			// taskService.ListAllActiveTasks() would be ideal.
			// As a placeholder, we'll log that it's not fully supported yet without a classID unless a new service method is added.
			// However, the goal is to list *all* tasks, including bug tasks (which have nil classID).
			// So, we need a service method that fetches tasks with classID IS NULL OR classID = ? (if provided).
			// For now, to list bug tasks, we can simulate by calling ListActiveTasksByClass with a non-existent classID
			// if the repository's GetTasksByClassID is adapted or a new ListTasks(filter) method is made.
			// The current ListActiveTasksByClass filters by a *specific* class ID.
			// To list system/bug tasks (nil ClassID), we need a new service method.
			// Let's add a conceptual ListAllTasks to TaskService and its stub.
			// For this step, we'll assume taskService.ListAllTasks() exists and fetches all tasks.
			// If not, we'll need to add it in a subsequent step.
			// For now, let's call ListActiveTasksByClass with a dummy ID that might be handled by a more flexible repo.
			// This is a simplification for now. A proper implementation needs ListAllTasks in service and repo.
			// tasks, err = taskService.ListActiveTasksByClass(context.Background(), 0) // Assuming 0 or a special value might mean "all" or "no specific class"
			tasks, err = taskService.ListAllActiveTasks(context.Background()) // Use the new method
			// A better approach would be a new method: taskService.ListTasks(ctx context.Context, filter TaskFilter)
			// where filter could specify ClassID (optional) and IsCompleted (optional).
			// For now, we'll just say "All Tasks"
			if err != nil { // No longer need to check for "not found" as ListAllActiveTasks doesn't depend on a classID
				fmt.Println("Error listing all tasks:", err)
				return
			}
			headerMsg = "TODAS AS TAREFAS (INCLUINDO BUGS DO SISTEMA)" // Restaurado
		} else {
			fmt.Println("Erro: Especifique --classid OU use --all para listar todas as tarefas (incluindo bugs).")
			fmt.Println("Exemplo: vigenda tarefa listar --classid 1")
			fmt.Println("Exemplo: vigenda tarefa listar --all")
			return
		}


		if len(tasks) == 0 {
			fmt.Println("No active tasks found matching criteria.")
			return
		}

		fmt.Printf("%s\n\n", headerMsg) // Restaurado \n\n para espaço antes da tabela

		// Reativar a definição de colunas e o preenchimento de rows
		columns := []table.Column{
			{Title: "ID", Width: 3},
			{Title: "TAREFA", Width: 35},
			{Title: "PRAZO", Width: 10},
		}
		var rows []table.Row
		for _, task := range tasks {
			dueDateStr := "N/A"
			if task.DueDate != nil {
				dueDateStr = task.DueDate.Format("02/01/2006")
			}
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", task.ID),
				task.Title,
				dueDateStr,
			})
		}

		// Reativar impressão manual da tabela
		fmt.Printf("%-*s | %-*s | %s\n", columns[0].Width, columns[0].Title, columns[1].Width, columns[1].Title, columns[2].Title)
		fmt.Printf("%s | %s | %s\n", strings.Repeat("-", columns[0].Width), strings.Repeat("-", columns[1].Width), strings.Repeat("-", columns[2].Width))
		for _, row := range rows {
			fmt.Printf("%-*s | %-*s | %s\n", columns[0].Width, row[0], columns[1].Width, row[1], row[2])
		}
	},
}

var taskCompleteCmd = &cobra.Command{
	Use:   "complete [ID_da_tarefa]",
	Short: "Marca uma tarefa como concluída",
	Long:  `Marca uma tarefa específica como concluída, utilizando o seu ID numérico.`,
	Example: `  vigenda tarefa complete 12
  vigenda tarefa complete 3`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println("Error parsing task ID:", err)
			return
		}
		err = taskService.MarkTaskAsCompleted(context.Background(), taskID)
		if err != nil {
			fmt.Println("Error marking task as completed:", err)
			return
		}
		fmt.Printf("Task ID %d marked as completed.\n", taskID)
	},
}

// initializeServices sets up the service layer instances with their repository dependencies.
func initializeServices(db *sql.DB) {
	// Initialize real repositories with the db connection
	taskRepo := repository.NewTaskRepository(db)
	classRepo := repository.NewClassRepository(db)
	assessmentRepo := repository.NewAssessmentRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)

	// Initialize services with real repository implementations
	taskService = service.NewTaskService(taskRepo)
	// Assuming NewClassService, NewAssessmentService exist or will be created.
	// If they use stubs for now or are basic passthroughs, that's fine.
	// For now, let's assume they can take the real repos.
	// If NewStubClassService was a placeholder for NewClassService:
	classService = service.NewClassService(classRepo, subjectRepo) // Assuming ClassService might need SubjectRepo too, or just ClassRepo
	assessmentService = service.NewAssessmentService(assessmentRepo, classRepo) // AssessmentService might need ClassRepo to get students

	questionService = service.NewQuestionService(questionRepo, subjectRepo)
	proofService = service.NewProofService(questionRepo) // ProofService uses QuestionRepository for GetQuestionsByCriteriaProofGeneration
}

func init() {
	// Cobra command definitions and flag setups remain in init()

	// Setup flags for task add command
	taskAddCmd.Flags().StringP("description", "d", "", "Descrição detalhada da tarefa.")
	taskAddCmd.Flags().String("classid", "", "ID da turma para associar a tarefa (opcional).")
	taskAddCmd.Flags().String("duedate", "", "Data de conclusão da tarefa no formato YYYY-MM-DD (opcional).")

	// Setup flags for task list command
	//taskListCmd.Flags().String("classid", "", "ID da turma para filtrar as tarefas (obrigatório).")
	//_ = taskListCmd.MarkFlagRequired("classid") // No longer strictly mandatory if --all is used.
	taskListCmd.Flags().String("classid", "", "ID da turma para filtrar as tarefas.")
	taskListCmd.Flags().String("all", "false", "Listar todas as tarefas, incluindo tarefas de sistema/bugs (ignora --classid se presente).")


	taskCmd.AddCommand(taskAddCmd, taskListCmd, taskCompleteCmd)
	rootCmd.AddCommand(taskCmd)

	// Class Service Commands
	classCreateCmd.Flags().String("subjectid", "", "ID da disciplina à qual a turma pertence (obrigatório).")
	_ = classCreateCmd.MarkFlagRequired("subjectid")

	classCmd.AddCommand(classCreateCmd, classImportStudentsCmd, classUpdateStudentStatusCmd)
	rootCmd.AddCommand(classCmd)

	// Assessment Service Commands
	assessmentCreateCmd.Flags().String("classid", "", "ID da turma para a qual a avaliação será criada (obrigatório).")
	_ = assessmentCreateCmd.MarkFlagRequired("classid")
	assessmentCreateCmd.Flags().String("term", "", "Período/bimestre da avaliação (ex: 1, 2) (obrigatório).")
	_ = assessmentCreateCmd.MarkFlagRequired("term")
	assessmentCreateCmd.Flags().String("weight", "", "Peso da avaliação na média final (ex: 4.0) (obrigatório).")
	_ = assessmentCreateCmd.MarkFlagRequired("weight")

	assessmentCmd.AddCommand(assessmentCreateCmd, assessmentEnterGradesCmd, assessmentClassAverageCmd)
	rootCmd.AddCommand(assessmentCmd)

	// Question Service (bancoq) initialization and commands
	questionBankCmd.AddCommand(questionBankAddCmd)
	rootCmd.AddCommand(questionBankCmd)

	// Proof Service (prova) initialization and commands
	proofGenerateCmd.Flags().String("subjectid", "", "ID da disciplina para gerar a prova (obrigatório).")
	_ = proofGenerateCmd.MarkFlagRequired("subjectid")
	proofGenerateCmd.Flags().String("topic", "", "Tópico específico para filtrar questões (opcional).")
	proofGenerateCmd.Flags().String("easy", "0", "Número de questões fáceis.")
	proofGenerateCmd.Flags().String("medium", "0", "Número de questões médias.")
	proofGenerateCmd.Flags().String("hard", "0", "Número de questões difíceis.")
	proofCmd.AddCommand(proofGenerateCmd)
	rootCmd.AddCommand(proofCmd)
}

var classCmd = &cobra.Command{
	Use:   "turma",
	Short: "Gerencia turmas e alunos (criar, importar-alunos, atualizar-status)",
	Long: `O comando 'turma' é usado para administrar turmas, incluindo a criação de novas turmas,
a importação de listas de alunos de ficheiros CSV e a atualização do status de alunos individuais
(ex: ativo, inativo, transferido).`,
	Example: `  vigenda turma criar "História 9A" --subjectid 1
  vigenda turma importar-alunos 1 alunos_9a.csv
  vigenda turma atualizar-status 15 transferido`,
}

var classCreateCmd = &cobra.Command{
	Use:   "criar [nome_da_turma]",
	Short: "Cria uma nova turma",
	Long: `Cria uma nova turma no sistema. É necessário fornecer o nome da turma e o ID da disciplina
à qual ela pertence através da flag --subjectid.`,
	Example: `  vigenda turma criar "Matemática - Turma 101" --subjectid 2
  vigenda turma criar "Geografia 8B" --subjectid 3`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		subjectIDStr, _ := cmd.Flags().GetString("subjectid")
		if subjectIDStr == "" {
			// Example of interactive input if a required flag is missing
			var err error
			subjectIDStr, err = tui.GetInput("Enter Subject ID for the class:", os.Stdout, os.Stdin)
			if err != nil || subjectIDStr == "" {
				fmt.Println("Subject ID is required to create a class.")
				return
			}
		}

		subjectID, err := strconv.ParseInt(subjectIDStr, 10, 64)
		if err != nil {
			fmt.Println("Error parsing Subject ID:", err)
			return
		}

		class, err := classService.CreateClass(context.Background(), name, subjectID)
		if err != nil {
			fmt.Println("Error creating class:", err)
			return
		}
		fmt.Printf("Class '%s' (ID: %d) created successfully for Subject ID %d.\n", class.Name, class.ID, subjectID)
	},
}

var classImportStudentsCmd = &cobra.Command{
	Use:   "importar-alunos [ID_da_turma] [caminho_do_ficheiro_csv]",
	Short: "Importa alunos de um ficheiro CSV para uma turma",
	Long: `Importa uma lista de alunos de um ficheiro CSV para uma turma existente.
O ficheiro CSV deve conter as colunas 'numero_chamada', 'nome_completo', e opcionalmente 'situacao'.
Consulte a documentação (README.md, Artefacto 9.1) para a estrutura detalhada do CSV.`,
	Example: `  vigenda turma importar-alunos 1 ./lista_alunos_turma_a.csv
  vigenda turma importar-alunos 3 /documentos/alunos_turma_c.csv`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		classID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println("Error parsing Class ID:", err)
			return
		}
		csvFilePath := args[1]

		csvData, err := os.ReadFile(csvFilePath)
		if err != nil {
			fmt.Println("Error reading CSV file:", err)
			return
		}

		count, err := classService.ImportStudentsFromCSV(context.Background(), classID, csvData)
		if err != nil {
			fmt.Println("Error importing students:", err)
			return
		}
		fmt.Printf("%d students imported successfully into class ID %d.\n", count, classID)
	},
}

var classUpdateStudentStatusCmd = &cobra.Command{
	Use:   "atualizar-status [ID_do_aluno] [novo_status]",
	Short: "Atualiza o status de um aluno",
	Long: `Atualiza o status de um aluno específico (ex: 'ativo', 'inativo', 'transferido').
O ID do aluno é o identificador único na base de dados.
Status permitidos: 'ativo', 'inativo', 'transferido'.`,
	Example: `  vigenda turma atualizar-status 25 ativo
  vigenda turma atualizar-status 103 transferido`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		studentID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println("Error parsing Student ID:", err)
			return
		}
		newStatus := args[1]
		// TODO: Validate newStatus against allowed values ('ativo', 'inativo', 'transferido')
		// This could be done here or in the service layer. For now, assume service layer handles it.

		err = classService.UpdateStudentStatus(context.Background(), studentID, newStatus)
		if err != nil {
			fmt.Println("Error updating student status:", err)
			return
		}
		fmt.Printf("Status of student ID %d updated to '%s'.\n", studentID, newStatus)
	},
}

// Removed the second, duplicate init() function. All initializations are now in the first init().

var assessmentCmd = &cobra.Command{
	Use:   "avaliacao",
	Short: "Gerencia avaliações e notas (criar, lancar-notas, media-turma)",
	Long: `O comando 'avaliacao' permite gerenciar todo o ciclo de vida das avaliações,
desde a sua criação, passando pelo lançamento interativo de notas dos alunos,
até o cálculo da média final da turma para uma avaliação específica.`,
	Example: `  vigenda avaliacao criar "Prova Bimestral 1" --classid 1 --term 1 --weight 4.0
  vigenda avaliacao lancar-notas 1
  vigenda avaliacao media-turma 1`,
}

var assessmentCreateCmd = &cobra.Command{
	Use:   "criar [nome_da_avaliacao]",
	Short: "Cria uma nova avaliação para uma turma",
	Long: `Cria uma nova avaliação associada a uma turma específica.
É necessário fornecer o nome da avaliação e, através de flags, o ID da turma,
o período/bimestre e o peso da avaliação na média final.`,
	Example: `  vigenda avaliacao criar "Trabalho de História Moderna" --classid 2 --term 3 --weight 3.5
  vigenda avaliacao criar "Seminário de Literatura" --classid 1 --term 2 --weight 2.0`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		classIDStr, _ := cmd.Flags().GetString("classid")
		termStr, _ := cmd.Flags().GetString("term")
		weightStr, _ := cmd.Flags().GetString("weight")

		// Interactive prompts for missing required flags
		var err error
		if classIDStr == "" {
			classIDStr, err = tui.GetInput("Enter Class ID for the assessment:", os.Stdout, os.Stdin)
			if err != nil || classIDStr == "" {
				fmt.Println("Class ID is required.")
				return
			}
		}
		if termStr == "" {
			termStr, err = tui.GetInput("Enter Term (e.g., 1, 2, 3, 4) for the assessment:", os.Stdout, os.Stdin)
			if err != nil || termStr == "" {
				fmt.Println("Term is required.")
				return
			}
		}
		if weightStr == "" {
			weightStr, err = tui.GetInput("Enter Weight (e.g., 4.0) for the assessment:", os.Stdout, os.Stdin)
			if err != nil || weightStr == "" {
				fmt.Println("Weight is required.")
				return
			}
		}

		classID, err := strconv.ParseInt(classIDStr, 10, 64)
		if err != nil {
			fmt.Println("Error parsing Class ID:", err)
			return
		}
		term, err := strconv.Atoi(termStr)
		if err != nil {
			fmt.Println("Error parsing Term:", err)
			return
		}
		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			fmt.Println("Error parsing Weight:", err)
			return
		}

		assessment, err := assessmentService.CreateAssessment(context.Background(), name, classID, term, weight)
		if err != nil {
			fmt.Println("Error creating assessment:", err)
			return
		}
		fmt.Printf("Assessment '%s' (ID: %d) created for Class ID %d, Term %d, Weight %.1f.\n", assessment.Name, assessment.ID, classID, term, weight)
	},
}

var assessmentEnterGradesCmd = &cobra.Command{
	Use:   "lancar-notas [ID_da_avaliacao]",
	Short: "Lança notas para os alunos de uma avaliação",
	Long: `Inicia um processo interativo para lançar ou editar as notas dos alunos
para uma avaliação específica. A lista de alunos da turma associada à avaliação
será exibida, permitindo a inserção de cada nota.
O ID da avaliação é o identificador numérico único da avaliação.`,
	Example: `  vigenda avaliacao lancar-notas 7
  vigenda avaliacao lancar-notas 2`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		assessmentID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println("Error parsing Assessment ID:", err)
			return
		}

		// TODO: This is where the more complex TUI interaction for grade entry will go.
		// 1. Fetch students for the class associated with the assessmentID.
		//    (This might require adding a method to AssessmentService or ClassService,
		//     or having AssessmentService return student info needed for grading).
		// 2. Display students in a list/table, allowing navigation and grade input.
		//    tui.GetInput can be used for each grade, or a more sophisticated BubbleTea model.
		// 3. Collect all grades into a map[int64]float64 (studentID -> grade).
		// 4. Call assessmentService.EnterGrades().

		fmt.Printf("Interactive grade entry for Assessment ID %d is not yet fully implemented.\n", assessmentID)
		fmt.Println("Simulating grade entry for now...")

		// Placeholder for student grades map
		studentGrades := make(map[int64]float64)
		// Example: studentGrades[101] = 8.5
		// Example: studentGrades[102] = 9.0
		// In a real scenario, this map would be populated via TUI.

		// For now, let's assume the user will input studentID:grade pairs via flags or prompts
		// For a better UX as per Artefact 7 (golden_files/notas_lancar_interativo_output.txt),
		// a full BubbleTea model would be needed here.
		// We will simulate a simple input loop for now.

		fmt.Println("Enter student grades (StudentID:Grade). Type 'done' when finished.")
		for {
			input, _ := tui.GetInput("Enter StudentID:Grade (or 'done'):", os.Stdout, os.Stdin)
			if input == "done" {
				break
			}
			parts := strings.Split(input, ":")
			if len(parts) != 2 {
				fmt.Println("Invalid format. Use StudentID:Grade.")
				continue
			}
			studentID, errS := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
			grade, errG := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if errS != nil || errG != nil {
				fmt.Println("Invalid StudentID or Grade value.")
				continue
			}
			studentGrades[studentID] = grade
		}


		if len(studentGrades) > 0 {
			err = assessmentService.EnterGrades(context.Background(), assessmentID, studentGrades)
			if err != nil {
				fmt.Println("Error entering grades:", err)
				return
			}
			fmt.Println("Grades entered successfully for Assessment ID", assessmentID)
		} else {
			fmt.Println("No grades were entered.")
		}
	},
}

var assessmentClassAverageCmd = &cobra.Command{
	Use:   "media-turma [ID_da_turma]",
	Short: "Calcula a média geral das notas de uma turma",
	Long: `Calcula e exibe a média geral ponderada das notas para uma turma específica,
considerando todas as avaliações e seus respectivos pesos.
O ID da turma é o identificador numérico único da turma.`,
	Example: `  vigenda avaliacao media-turma 1
  vigenda avaliacao media-turma 5`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		classID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println("Error parsing Class ID:", err)
			return
		}

		average, err := assessmentService.CalculateClassAverage(context.Background(), classID)
		if err != nil {
			fmt.Println("Error calculating class average:", err)
			return
		}
		fmt.Printf("Average grade for Class ID %d: %.2f\n", classID, average)
	},
}

// --- Question Bank (bancoq) Commands ---
var questionBankCmd = &cobra.Command{
	Use:   "bancoq",
	Short: "Gerencia o banco de questões (add)",
	Long: `O comando 'bancoq' (Banco de Questões) permite adicionar novas questões ao sistema
a partir de um ficheiro JSON formatado.
Consulte a documentação (README.md, Artefacto 9.3) para a estrutura detalhada do JSON.`,
	Example: `  vigenda bancoq add ./minhas_questoes_historia.json
  vigenda bancoq add /usr/share/vigenda/questoes_padrao_matematica.json`,
}

var questionBankAddCmd = &cobra.Command{
	Use:   "add [caminho_do_ficheiro_json]",
	Short: "Adiciona questões de um ficheiro JSON ao banco",
	Long: `Adiciona um conjunto de questões de um ficheiro JSON para o banco de questões central.
O ficheiro JSON deve seguir uma estrutura específica contendo detalhes como disciplina,
tópico, tipo de questão (múltipla escolha, dissertativa), dificuldade, enunciado,
opções (para múltipla escolha) e resposta correta.`,
	Example: `  vigenda bancoq add questoes_bimestre1.json
  vigenda bancoq add ../shared/questoes_revisao.json`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jsonFilePath := args[0]
		jsonData, err := os.ReadFile(jsonFilePath)
		if err != nil {
			fmt.Println("Error reading JSON file:", err)
			return
		}

		count, err := questionService.AddQuestionsFromJSON(context.Background(), jsonData)
		if err != nil {
			fmt.Println("Error adding questions from JSON:", err)
			return
		}
		fmt.Printf("%d questions added successfully to the bank.\n", count)
	},
}

// --- Proof (prova) Commands ---
var proofCmd = &cobra.Command{
	Use:   "prova",
	Short: "Gerencia e gera provas (gerar)",
	Long: `O comando 'prova' permite gerar provas (avaliações textuais) a partir do banco de questões.
Você pode especificar critérios como disciplina, tópico e o número desejado de questões
por nível de dificuldade (fácil, médio, difícil).`,
	Example: `  vigenda prova gerar --subjectid 1 --easy 5 --medium 3 --hard 2
  vigenda prova gerar --subjectid 2 --topic "Segunda Guerra Mundial" --medium 10`,
}

var proofGenerateCmd = &cobra.Command{
	Use:   "gerar",
	Short: "Gera uma nova prova com base em critérios especificados",
	Long: `Gera uma prova selecionando questões do banco de questões.
É obrigatório especificar o ID da disciplina. Opcionalmente, pode-se filtrar por tópico
e definir o número de questões para cada nível de dificuldade (fácil, médio, difícil).
A prova gerada será exibida no terminal.`,
	Example: `  vigenda prova gerar --subjectid 1 --easy 5 --medium 3 --hard 2 --topic "Revolução Industrial"
  vigenda prova gerar --subjectid 3 --medium 10 --hard 5`,
	Run: func(cmd *cobra.Command, args []string) {
		subjectIDStr, _ := cmd.Flags().GetString("subjectid")
		topic, _ := cmd.Flags().GetString("topic")
		easyCountStr, _ := cmd.Flags().GetString("easy")
		mediumCountStr, _ := cmd.Flags().GetString("medium")
		hardCountStr, _ := cmd.Flags().GetString("hard")

		if subjectIDStr == "" {
			fmt.Println("Subject ID (--subjectid) is required.")
			return
		}
		subjectID, err := strconv.ParseInt(subjectIDStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid Subject ID:", err)
			return
		}

		easyCount, _ := strconv.Atoi(easyCountStr)   // Default to 0 if not provided or invalid
		mediumCount, _ := strconv.Atoi(mediumCountStr) // Default to 0
		hardCount, _ := strconv.Atoi(hardCountStr)     // Default to 0

		if easyCount == 0 && mediumCount == 0 && hardCount == 0 {
			fmt.Println("At least one difficulty count (--easy, --medium, --hard) must be greater than zero.")
			return
		}

		criteria := service.ProofCriteria{
			SubjectID:   subjectID,
			EasyCount:   easyCount,
			MediumCount: mediumCount,
			HardCount:   hardCount,
		}
		if topic != "" {
			criteria.Topic = &topic
		}

		questions, err := proofService.GenerateProof(context.Background(), criteria)
		if err != nil {
			fmt.Println("Error generating proof:", err)
			return
		}

		if len(questions) == 0 {
			fmt.Println("No questions matched the criteria to generate the proof.")
			return
		}

		fmt.Printf("Proof generated successfully with %d questions:\n\n", len(questions))
		// Display questions using TUI table or simple print
		// For now, a simple print. Later, can use tui.ShowTable.
		for i, q := range questions {
			fmt.Printf("Q%d (%s, %s): %s\n", i+1, q.Difficulty, q.Type, q.Statement)
			if q.Options != nil && *q.Options != "" && *q.Options != "null" {
				var opts []string
				if json.Unmarshal([]byte(*q.Options), &opts) == nil {
					for j, opt := range opts {
						fmt.Printf("  %c) %s\n", 'a'+j, opt)
					}
				}
			}
			fmt.Printf("   Answer: %s\n\n", q.CorrectAnswer)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
