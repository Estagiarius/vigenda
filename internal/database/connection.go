package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings" // Added import for strings.Join

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
		var schemaBytes []byte
		var readErr error
		var attemptedPaths []string

		// Path 1: Relative to executable (for distributed binary)
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			schemaPathFromExec := filepath.Join(execDir, "..", "internal", "database", "migrations", "001_initial_schema.sql")
			attemptedPaths = append(attemptedPaths, schemaPathFromExec)
			schemaBytes, readErr = os.ReadFile(schemaPathFromExec)
		}

		// Path 2: Relative to CWD (for 'go run' from project root)
		if readErr != nil {
			schemaPathFromCwd := "internal/database/migrations/001_initial_schema.sql"
			attemptedPaths = append(attemptedPaths, schemaPathFromCwd)
			schemaBytes, readErr = os.ReadFile(schemaPathFromCwd)
		}

		// Path 3: Fallback relative path (for tests or other structures)
		if readErr != nil {
			alternativeSchemaPath := filepath.Join("..", "..", "internal", "database", "migrations", "001_initial_schema.sql")
			attemptedPaths = append(attemptedPaths, alternativeSchemaPath)
			schemaBytes, readErr = os.ReadFile(alternativeSchemaPath)
		}

		if readErr != nil {
			db.Close()
			// Use strings.Join for better readability of attempted paths in the error message
			return nil, fmt.Errorf("failed to read schema file. Attempted paths: [%s]. Last error: %w", strings.Join(attemptedPaths, ", "), readErr)
		}

		_, err = db.Exec(string(schemaBytes))
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to apply initial schema (tried paths: [%s]): %w", strings.Join(attemptedPaths, ", "), err)
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
