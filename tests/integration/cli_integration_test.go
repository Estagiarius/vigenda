package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"runtime"
	"fmt"
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

// runCLI executes the compiled CLI command with the given arguments.
func runCLI(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command(binPath, args...)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
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

	// For the dashboard, the command is just `vigenda` (no arguments)
	// This requires the database to be in a specific state to produce the golden output.
	// For now, we'll just check if the command runs and compare against the golden file.
	// Setting up the database state is a more complex task for later or needs to be handled by test fixtures.

	// stdout, stderr, err := runCLI(t) // No args for dashboard
	// if err != nil {
	//  // Handle error, perhaps stderr has useful info
	// 	t.Fatalf("CLI execution failed: %v\nStderr: %s", err, stderr)
	// }
	// if stderr != "" {
	// 	t.Logf("Stderr output: %s", stderr) // Log non-fatal stderr
	// }
	// assertGoldenFile(t, stdout, "golden_files/dashboard_output.txt")
	t.Log("TestDashboardOutput: Placeholder - requires database setup for meaningful comparison.")
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
	// stdout, stderr, err := runCLI(t, "tarefa", "listar", "--turma", "Turma 9A")
	// if err != nil {
	// 	t.Fatalf("CLI execution failed for 'tarefa listar': %v\nStderr: %s", err, stderr)
	// }
	// if stderr != "" {
	// 	t.Logf("Stderr output for 'tarefa listar': %s", stderr)
	// }
	// assertGoldenFile(t, stdout, "golden_files/tarefa_listar_turma_output.txt")
	t.Log("TestTarefaListarTurmaOutput: Placeholder - requires database setup for meaningful comparison.")
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
