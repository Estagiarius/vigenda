// Package main is the entry point of the Vigenda application.
// It handles the main function and Cobra CLI configuration.
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
	"vigenda/internal/service"
	"vigenda/internal/tui"

	"strings"
	"github.com/spf13/cobra"
	"github.com/charmbracelet/bubbles/table"
	"encoding/json" // Added missing import
)

// Removed duplicate import block

var taskService service.TaskService           // Será inicializado em init() ou main()
var classService service.ClassService         // Será inicializado em init() ou main()
var assessmentService service.AssessmentService // Será inicializado em init() ou main()
var questionService service.QuestionService     // Será inicializado em init() ou main()
var proofService service.ProofService           // Será inicializado em init() ou main()

var rootCmd = &cobra.Command{
	Use:   "vigenda",
	Short: "Vigenda is a CLI tool for teachers with ADHD.",
	Long:  `Vigenda helps teachers manage tasks, classes, assessments, and more, directly from the command line.`,
}

var taskCmd = &cobra.Command{
	Use:   "tarefa",
	Short: "Manage tasks",
	Long:  `Commands for creating, listing, and managing tasks.`,
}

var taskAddCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new task",
	Long:  `Adds a new task. You can optionally provide a description, class ID, and due date using flags.`,
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

		// TODO: Obter descrição de forma interativa se não fornecida via flag, usando tui.GetInput
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
	Short: "List tasks",
	Long:  `Lists tasks. Can be filtered by class ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		classIDStr, _ := cmd.Flags().GetString("classid")
		if classIDStr == "" {
			// TODO: Implement listing all tasks if desired, or prompt for class ID
			fmt.Println("Please specify a class ID using --classid to list tasks.")
			// For now, let's assume we always need a class ID for listing.
			// Alternatively, prompt for class ID using tui.GetInput or a selection list.
			return
		}

		classID, err := strconv.ParseInt(classIDStr, 10, 64)
		if err != nil {
			fmt.Println("Error parsing class ID:", err)
			return
		}

		tasks, err := taskService.ListActiveTasksByClass(context.Background(), classID)
		if err != nil {
			fmt.Println("Error listing tasks:", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("No active tasks found for this class.")
			return
		}

		columns := []table.Column{
			{Title: "ID", Width: 5},
			{Title: "Title", Width: 30},
			{Title: "Description", Width: 40},
			{Title: "Due Date", Width: 12},
		}
		var rows []table.Row
		for _, task := range tasks {
			dueDateStr := "N/A"
			if task.DueDate != nil {
				dueDateStr = task.DueDate.Format("2006-01-02")
			}
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", task.ID),
				task.Title,
				task.Description,
				dueDateStr,
			})
		}
		tui.ShowTable(columns, rows, os.Stdout)
	},
}

var taskCompleteCmd = &cobra.Command{
	Use:   "complete [taskID]",
	Short: "Mark a task as completed",
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

func init() {
	// TODO: Initialize services properly, likely with repository dependencies.
	// For now, using a placeholder initialization for TaskService.
	// This will need to be replaced with the actual service implementation
	// that connects to a repository.
	taskService = service.NewTaskService(/* pass repository here */)

	// Setup flags for task add command
	taskAddCmd.Flags().StringP("description", "d", "", "Description of the task")
	taskAddCmd.Flags().String("classid", "", "Class ID to associate the task with")
	taskAddCmd.Flags().String("duedate", "", "Due date of the task (YYYY-MM-DD)")

	// Setup flags for task list command
	taskListCmd.Flags().String("classid", "", "Filter tasks by Class ID")


	taskCmd.AddCommand(taskAddCmd, taskListCmd, taskCompleteCmd)
	rootCmd.AddCommand(taskCmd)

	// Class Service Commands
	classService = service.NewClassService(/* pass repository here */)
	classCreateCmd.Flags().String("subjectid", "", "Subject ID for the class") // Flag for class create
	classCmd.AddCommand(classCreateCmd, classImportStudentsCmd, classUpdateStudentStatusCmd)
	rootCmd.AddCommand(classCmd)

	// Assessment Service Commands
	assessmentService = service.NewAssessmentService(/* pass repository here */)
	assessmentCreateCmd.Flags().String("classid", "", "Class ID for the assessment")
	assessmentCreateCmd.Flags().String("term", "", "Term for the assessment (e.g., 1, 2)")
	assessmentCreateCmd.Flags().String("weight", "", "Weight of the assessment (e.g., 4.0)")
	assessmentCmd.AddCommand(assessmentCreateCmd, assessmentEnterGradesCmd, assessmentClassAverageCmd)
	rootCmd.AddCommand(assessmentCmd)

	// Question Service (bancoq) initialization and commands
	questionService = service.NewQuestionService(nil, nil /* TODO: Pass real repositories */) // Placeholder
	// No flags for questionBankAddCmd as it takes a positional argument.
	questionBankCmd.AddCommand(questionBankAddCmd)
	rootCmd.AddCommand(questionBankCmd)

	// Proof Service (prova) initialization and commands
	proofService = service.NewProofService(nil /* TODO: Pass real QuestionRepository */) // Placeholder
	proofGenerateCmd.Flags().String("subjectid", "", "Subject ID for the proof (required)")
	proofGenerateCmd.Flags().String("topic", "", "Topic to filter questions by (optional)")
	proofGenerateCmd.Flags().String("easy", "0", "Number of easy questions")
	proofGenerateCmd.Flags().String("medium", "0", "Number of medium questions")
	proofGenerateCmd.Flags().String("hard", "0", "Number of hard questions")
	proofCmd.AddCommand(proofGenerateCmd)
	rootCmd.AddCommand(proofCmd)
}

var classCmd = &cobra.Command{
	Use:   "turma",
	Short: "Manage classes and students",
}

var classCreateCmd = &cobra.Command{
	Use:   "criar [name]",
	Short: "Create a new class",
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
	Use:   "importar-alunos [classID] [csvFilePath]",
	Short: "Import students from a CSV file into a class",
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
	Use:   "atualizar-status [studentID] [newStatus]",
	Short: "Update the status of a student",
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
	Short: "Manage assessments and grades",
}

var assessmentCreateCmd = &cobra.Command{
	Use:   "criar [name]",
	Short: "Create a new assessment",
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
	Use:   "lancar-notas [assessmentID]",
	Short: "Enter grades for an assessment",
	Long:  "Interactively enter grades for students for a given assessment. Students and their current grades (if any) will be listed.",
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
	Use:   "media-turma [classID]",
	Short: "Calculate the average grade for a class",
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
	Short: "Manage the question bank",
}

var questionBankAddCmd = &cobra.Command{
	Use:   "add [jsonFilePath]",
	Short: "Add questions from a JSON file to the bank",
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
	Short: "Manage and generate proofs (tests)",
}

var proofGenerateCmd = &cobra.Command{
	Use:   "gerar",
	Short: "Generate a new proof based on criteria",
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
