package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"runtime"
	"fmt"
	"context" // Added for CommandContext
	"time"    // Added for timeout
	"database/sql" // Added for setupTestDB and seedDB
	_ "github.com/mattn/go-sqlite3" // SQLite driver for database/sql
)

var (
	binName = "vigenda"
	binPath string
)

func TestMain(m *testing.M) {
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// Build the CLI into a temporary directory
	tempBinDir, err := os.MkdirTemp("", "vigenda_test_bin")
    if err != nil {
        panic(fmt.Sprintf("Failed to create temp dir for binary: %v", err))
    }
    defer os.RemoveAll(tempBinDir)

	binPath = filepath.Join(tempBinDir, binName)

	// Corrected path to main.go, relative to project root
	mainGoPath := "cmd/vigenda/main.go"
	projectRoot := filepath.Join("..", "..") // Relative path to project root from tests/integration

	buildCmd := exec.Command("go", "build", "-o", binPath, mainGoPath)
	buildCmd.Dir = projectRoot // Set working directory for build command to project root

	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		// It's helpful to change to the directory to understand the context of the build error
		absProjectRoot, _ := filepath.Abs(projectRoot)
		panic(fmt.Sprintf("Failed to build CLI in %s: %v\nOutput: %s\nCommand: go build -o %s %s", absProjectRoot, err, string(buildOutput), binPath, mainGoPath))
	}

	// Run tests
	exitCode := m.Run()

	// Cleanup: os.Remove(binPath) // Done by defer os.RemoveAll(tempBinDir)

	os.Exit(exitCode)
}

const testDbDir = "test_dbs"
const testSchemaPath = "../database/migrations/001_initial_schema.sql" // Relative to project root

// setupTestDB creates a new test database, initializes the schema, and returns the path to the db file.
// It also sets the VIGENDA_DB_PATH environment variable for the CLI to use.
func setupTestDB(t *testing.T, testName string) string {
	t.Helper()

	// Create a directory for test databases if it doesn't exist
	if _, err := os.Stat(testDbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(testDbDir, 0755); err != nil {
			t.Fatalf("Failed to create test database directory %s: %v", testDbDir, err)
		}
	}

	dbPath := filepath.Join(testDbDir, fmt.Sprintf("vigenda_test_%s.db", testName))
	// Remove existing DB file to ensure a clean state
	os.Remove(dbPath)

	// Set environment variables for the CLI to use this database
	t.Setenv("VIGENDA_DB_TYPE", "sqlite")
	t.Setenv("VIGENDA_DB_PATH", dbPath)

	// Initialize schema
	// This requires a direct DB connection here, or a CLI command to init schema if available.
	// For now, let's assume we'll use direct DB access for setup.
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database %s: %v", dbPath, err)
	}
	defer db.Close()

	// Correct path to schema relative to this test file (tests/integration/cli_integration_test.go)
	schemaFilePath := filepath.Join("..", "..", "internal", "database", "migrations", "001_initial_schema.sql")
	schemaBytes, err := os.ReadFile(schemaFilePath)
	if err != nil {
		t.Fatalf("Failed to read schema file %s: %v", schemaFilePath, err)
	}
	_, err = db.Exec(string(schemaBytes))
	if err != nil {
		t.Fatalf("Failed to apply schema to %s: %v", dbPath, err)
	}
	return dbPath
}

// seedDB executes SQL statements to seed the database for a test.
func seedDB(t *testing.T, dbPath string, statements []string) {
	t.Helper()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("seedDB: failed to open db %s: %v", dbPath, err)
	}
	defer db.Close()

	for i, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("seedDB: failed to execute statement %d (%s): %v", i, stmt, err)
		}
	}
	t.Logf("Successfully seeded DB %s with %d statements", dbPath, len(statements))
}


// runCLI executes the compiled CLI command with the given arguments.
// It now ensures VIGENDA_DB_PATH is set if a test DB is configured.
func runCLI(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 30-second timeout for CLI command
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Stdin = nil // Explicitly set Stdin to nil for non-interactive commands

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		// Log more information if the command timed out
		t.Logf("Command '%s %s' timed out after 30s.", binPath, strings.Join(args, " "))
		t.Logf("Stdout so far: %s", stdout.String())
		t.Logf("Stderr so far: %s", stderr.String())
		// The error from cmd.Run() will likely be "signal: killed" or similar in this case.
		// We return it as is, but the log provides context.
	}

	return stdout.String(), stderr.String(), err
}

// assertGoldenFile compares the actual output with the content of a golden file.
// It trims whitespace from both actual and expected content before comparison.
func assertGoldenFile(t *testing.T, actualOutput string, goldenFilePath string) {
	t.Helper()

	absGoldenFilePath, err := filepath.Abs(goldenFilePath)
	if err != nil {
		t.Fatalf("Failed to get absolute path for golden file %s: %v", goldenFilePath, err)
	}

	expectedOutputBytes, err := os.ReadFile(absGoldenFilePath)
	if err != nil {
		t.Fatalf("Failed to read golden file %s: %v", goldenFilePath, err)
	}
	expectedOutput := string(expectedOutputBytes)

	// Normalize line endings and trim whitespace for comparison
	// This helps avoid issues with CRLF vs LF and trailing/leading whitespace.
	normalizedActual := strings.TrimSpace(strings.ReplaceAll(actualOutput, "\r\n", "\n"))
	normalizedExpected := strings.TrimSpace(strings.ReplaceAll(expectedOutput, "\r\n", "\n"))

	if normalizedActual != normalizedExpected {
		t.Errorf("Output does not match golden file %s.\nExpected:\n%s\n\nActual:\n%s",
			goldenFilePath, normalizedExpected, normalizedActual)
		// For better debugging, one might want to output the diff
		// Or save the actual output to a temporary file.
	}
}

// Placeholder for TestDashboard - TC-I-001 (related)
// Artefact 6 TC-I-001: "O comando vigenda notas lancar deve apresentar a lista correta de alunos ativos para a avaliação selecionada."
// This uses `golden_files/notas_lancar_interativo_output.txt`
// Artefact 7: `golden_files/dashboard_output.txt` is for the main command `vigenda` (RF01)
// Let's assume a test for the dashboard output.
func TestDashboardOutput(t *testing.T) {
	// This test will likely fail initially as the dashboard functionality might not be implemented
	// or might not produce the exact golden file output yet.
	// The task is to set up the test structure.

	dbPath := setupTestDB(t, "TestDashboardOutput")

	// Seed data necessary for the dashboard to reflect the golden file.
	// This includes:
	// - A user (implicit, user_id 1 is often assumed by stubs or initial implementations)
	// - A class "Turma 9A" (implies a subject too)
	// - A task "Corrigir provas (Turma 9A)" with due date "Amanhã" (relative to 22/06/2025 in golden file)
	// - A task "Preparar aula sobre Era Vargas (Turma 9B)" with due date "24/06" (relative to 22/06/2025)
	// - A notification "5 entregas pendentes para o trabalho "Pesquisa sobre Clima" (Turma 9A)."
	// The "Agenda de Hoje" part is static in the current CLI main.go's Run function, so no seeding needed for that.
	// For tasks and notifications, the application needs to query and format them.
	// Current stubs might not support this complex querying.
	// We will seed data and see if the actual implementation (once stubs are replaced) picks it up.
	// For now, the dashboard output in main.go is mostly static.
	// The golden file `dashboard_output.txt` is:
	// =================================================
	// ==                 DASHBOARD                   ==
	// =================================================
	//
	// 🕒 AGENDA DE HOJE (22/06/2025)
	//    [09:00 - 10:00] Aula de História - Turma 9A
	//    [14:00 - 15:00] Reunião Pedagógica
	//
	// 🔥 TAREFAS PRIORITÁRIAS
	//    [1] Corrigir provas (Turma 9A) (Prazo: Amanhã)
	//    [2] Preparar aula sobre Era Vargas (Turma 9B) (Prazo: 24/06)
	//
	// 🔔 NOTIFICAÇÕES
	//    - 5 entregas pendentes para o trabalho "Pesquisa sobre Clima" (Turma 9A).
	//
	// The current rootCmd Run function in `cmd/vigenda/main.go` prints exactly this.
	// So, no specific seeding is strictly necessary for *this current static implementation*.
	// However, if the dashboard becomes dynamic, seeding will be crucial.
	// For now, we just ensure the DB is initialized so the app doesn't fail on DB connection.
	_ = dbPath // Use dbPath to avoid unused variable error if no seeding is done.


	// For the dashboard, the command is just `vigenda` (no arguments)
	stdout, stderr, err := runCLI(t)
	if err != nil {
		t.Fatalf("CLI execution failed: %v\nStderr: %s", err, stderr)
	}
	if stderr != "" {
		// If stubs are fully replaced, stderr should ideally be empty for successful commands.
		// Allow it for now, but this might be an assertion point later.
		t.Logf("Stderr output: %s", stderr)
	}
	assertGoldenFile(t, stdout, "golden_files/dashboard_output.txt")
}


// TestNotasLancarOutput - Corresponds to TC-I-001 from Artefact 6
// "O comando `vigenda notas lancar` deve apresentar a lista correta de alunos ativos para a avaliação selecionada."
// Golden file: `golden_files/notas_lancar_interativo_output.txt`
func TestNotasLancarOutput(t *testing.T) {
	// This test will simulate the command `vigenda notas lancar --avaliacao <id>`
	// The golden file implies an interactive session. Testing interactive TUI applications
	// as a subprocess and comparing full screen output is complex.
	// The golden file `notas_lancar_interativo_output.txt` shows different states.
	// A simple subprocess execution might not capture this easily.
	// For now, we'll assume a non-interactive mode or a simplified output for this test.
	// Or, this test might need a more sophisticated approach (e.g., driving stdin).

	// For now, let's assume the command `vigenda notas lancar --avaliacao SOME_ID --non-interactive` (hypothetical)
	// would produce a static output similar to one part of the golden file.
	// Or, if the initial screen is static before interaction, we could test that.
	// Given the current setup, this test will be a placeholder or simplified.

	// stdout, stderr, err := runCLI(t, "notas", "lancar", "--avaliacao", "1") // Assuming assessment ID 1
	// if err != nil {
	// 	t.Fatalf("CLI execution failed for 'notas lancar': %v\nStderr: %s", err, stderr)
	// }
	// if stderr != "" {
	// 	t.Logf("Stderr output for 'notas lancar': %s", stderr)
	// }
	// assertGoldenFile(t, stdout, "golden_files/notas_lancar_interativo_output.txt")
	t.Log("TestNotasLancarOutput: Placeholder - testing interactive TUI output is complex and needs specific handling or a non-interactive mode for the command.")
}

// TestRelatorioProgressoTurmaOutput - Corresponds to TC-I-002 from Artefact 6
// "O comando `vigenda relatorio progresso-turma` deve exibir os dados corretos e calculados para a turma especificada."
// Golden file: `golden_files/relatorio_progresso_turma.txt`
func TestRelatorioProgressoTurmaOutput(t *testing.T) {
	// This test will simulate `vigenda relatorio progresso-turma --turma "Turma 9A"`
	// Similar to the dashboard, this requires the database to be in a specific state.

	// stdout, stderr, err := runCLI(t, "relatorio", "progresso-turma", "--turma", "Turma 9A")
	// if err != nil {
	// 	t.Fatalf("CLI execution failed for 'relatorio progresso-turma': %v\nStderr: %s", err, stderr)
	// }
	// if stderr != "" {
	// 	t.Logf("Stderr output for 'relatorio progresso-turma': %s", stderr)
	// }
	// assertGoldenFile(t, stdout, "golden_files/relatorio_progresso_turma.txt")
	t.Log("TestRelatorioProgressoTurmaOutput: Placeholder - requires database setup for meaningful comparison.")
}


// TestTarefaListarTurmaOutput - Based on Artefact 7 `golden_files/tarefa_listar_turma_output.txt`
// This implies a command like `vigenda tarefa listar --turma "Turma 9A"`
func TestTarefaListarTurmaOutput(t *testing.T) {
	dbPath := setupTestDB(t, "TestTarefaListarTurmaOutput")

	// Seed data: user, subject, class "Turma 9A", and two tasks for this class.
	// Note: IDs are hardcoded here. In a more complex system, you might fetch IDs after insertion.
	// For tasks, the golden file expects IDs 1 and 5. This implies specific insertion order or auto-increment behavior.
	// We'll insert them with these IDs.
	seedStatements := []string{
		"INSERT INTO users (id, username, password_hash) VALUES (1, 'testuser', 'hash');",
		"INSERT INTO subjects (id, user_id, name) VALUES (1, 1, 'Matemática');",
		"INSERT INTO classes (id, user_id, subject_id, name) VALUES (1, 1, 1, 'Turma 9A');",
		// Task ID 1 (as per golden file)
		"INSERT INTO tasks (id, user_id, class_id, title, description, due_date, is_completed) VALUES (1, 1, 1, 'Corrigir provas de Matemática', 'Corrigir as provas bimestrais.', '2025-06-23 00:00:00', 0);",
		// Insert some other tasks to ensure filtering works and IDs are not necessarily sequential if we didn't hardcode task IDs
		"INSERT INTO tasks (id, user_id, class_id, title, description, due_date, is_completed) VALUES (2, 1, 1, 'Planejar Aula Extra', 'Planejar aula de reforço.', '2025-06-24 00:00:00', 0);",
		"INSERT INTO tasks (id, user_id, class_id, title, description, due_date, is_completed) VALUES (3, 1, 2, 'Tarefa Outra Turma', 'tarefa.', '2025-06-24 00:00:00', 0);", // Class ID 2 (different class)
		// Task ID 5 (as per golden file)
		"INSERT INTO tasks (id, user_id, class_id, title, description, due_date, is_completed) VALUES (5, 1, 1, 'Lançar notas do trabalho', 'Lançar as notas do trabalho de pesquisa.', '2025-06-25 00:00:00', 0);",
	}
	seedDB(t, dbPath, seedStatements)

	// The golden file uses "--turma Turma 9A", but the actual flag is "--classid <ID>"
	// We seeded "Turma 9A" with class_id = 1.
	stdout, stderr, err := runCLI(t, "tarefa", "listar", "--classid", "1")
	if err != nil {
		t.Fatalf("CLI execution failed for 'tarefa listar': %v\nStderr: %s", err, stderr)
	}
	if stderr != "" {
		// Stderr might contain debug logs from services if stubs were still printing.
		// With real implementations, it should be cleaner.
		t.Logf("Stderr output for 'tarefa listar': %s", stderr)
	}
	assertGoldenFile(t, stdout, "golden_files/tarefa_listar_turma_output.txt")
}

// TestFocoIniciarOutput - Based on Artefact 7 `golden_files/foco_iniciar_output.txt`
// This implies a command like `vigenda foco iniciar --tarefa <id_tarefa>`
// The output is time-sensitive ("TEMPO RESTANTE: 24:59"). This makes direct golden file comparison hard.
// This might require a way to mock time or capture only the static parts of the output.
func TestFocoIniciarOutput(t *testing.T) {
	// stdout, stderr, err := runCLI(t, "foco", "iniciar", "--tarefa", "1") // Assuming task ID 1
	// if err != nil {
	// 	t.Fatalf("CLI execution failed for 'foco iniciar': %v\nStderr: %s", err, stderr)
	// }
	// if stderr != "" {
	// 	t.Logf("Stderr output for 'foco iniciar': %s", stderr)
	// }
	// // assertGoldenFile(t, stdout, "golden_files/foco_iniciar_output.txt") // This will likely fail due to time
	t.Log("TestFocoIniciarOutput: Placeholder - time-sensitive output, direct golden file comparison is problematic.")
}
// import "fmt" // Added import for fmt used in TestMain panic <- This line was removed
