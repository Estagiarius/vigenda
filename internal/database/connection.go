package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// GetDBConnection establishes a connection to the SQLite database.
// It also ensures the database schema is initialized.
// The dataSourceName is typically the path to the SQLite file.
func GetDBConnection(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database at %s: %w", dataSourceName, err)
	}

	// Ping the database to ensure the connection is live.
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database at %s: %w", dataSourceName, err)
	}

	// Apply migrations
	// In a real app, this would be more sophisticated, handling versioning.
	// For this project, we apply the initial schema if the tables don't exist.
	// A simple check: if the 'users' table exists, assume schema is applied.
	// This is a simplification. Proper migration management is preferred.
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='users';")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to check for existing tables: %w", err)
	}
	defer rows.Close()

	if !rows.Next() { // If 'users' table doesn't exist, apply schema
		// Attempt to find schema relative to executable if running built binary,
		// or relative to project structure if running with 'go run'.
		// This path needs to be robust. For tests, it's relative to test execution.
		// For the actual app, it's relative to the binary or `go run` CWD.

		// Try to determine the correct path to the schema file.
		// This is tricky because the CWD changes depending on how the app is run.
		// For `go test ./...` from project root, `internal/database/migrations/...` would be accessible.
		// For the built binary, it might need to be bundled or placed in a known relative path.

		// Simplified approach: Assume schema is in a known relative path from CWD
		// This might need adjustment based on build/deployment structure.
		// Let's assume CWD is project root when running `go run cmd/vigenda/main.go`
		// or when tests run the binary (tests set their own CWD or binary is in temp dir).

		// Path from project root: internal/database/migrations/001_initial_schema.sql
		primarySchemaPath := "internal/database/migrations/001_initial_schema.sql"
		// Fallback for tests where CWD might be e.g. tests/integration, or if binary is in a different relative location.
		alternativeSchemaPath := filepath.Join("..", "..", "internal", "database", "migrations", "001_initial_schema.sql")

		schemaBytes, err := os.ReadFile(primarySchemaPath)
		if err != nil {
			// Try alternative path
			schemaBytes, err = os.ReadFile(alternativeSchemaPath)
			if err != nil {
				db.Close()
				return nil, fmt.Errorf("failed to read schema file at %s or %s: %w", primarySchemaPath, alternativeSchemaPath, err)
			}
		}

		_, err = db.Exec(string(schemaBytes))
		if err != nil {
			db.Close()
			// The error message here could refer to primarySchemaPath or alternativeSchemaPath
			// depending on which one was successfully read. For simplicity, just state the action failed.
			return nil, fmt.Errorf("failed to apply initial schema (tried %s and %s): %w", primarySchemaPath, alternativeSchemaPath, err)
		}
		// fmt.Println("Database schema initialized.") // Removed noisy print
	}

	return db, nil
}

// DefaultDbPath returns the default path for the SQLite database file.
// It places it in the user's config directory or defaults to "vigenda.db" in CWD.
func DefaultDbPath() string {
	configDir, err := os.UserConfigDir()
	if err == nil {
		appConfigDir := filepath.Join(configDir, "vigenda")
		if err := os.MkdirAll(appConfigDir, 0750); err == nil {
			return filepath.Join(appConfigDir, "vigenda.db")
		}
	}
	// Fallback to current working directory if user config dir fails
	return "vigenda.db"
}
