package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DBConfig holds all parameters needed to connect to a database.
type DBConfig struct {
	DBType string // "sqlite" or "postgres"
	DSN    string // Data Source Name, specific to DBType
}

// GetDBConnection establishes a connection to the specified database.
// It also ensures the database schema is initialized for SQLite.
// For PostgreSQL, schema migration is assumed to be handled externally.
func GetDBConnection(config DBConfig) (*sql.DB, error) {
	var driverName string
	var dsn string

	switch config.DBType {
	case "sqlite", "": // Default to SQLite if DBType is empty
		driverName = "sqlite3"
		dsn = config.DSN
		if dsn == "" {
			dsn = DefaultSQLitePath() // Use default SQLite path if DSN is empty
		}
	case "postgres":
		driverName = "postgres"
		dsn = config.DSN
		if dsn == "" {
			// Construct DSN from environment variables for PostgreSQL
			// Or use sensible defaults. This part will be refined in main.go or a config loader.
			// For now, assume DSN is provided or constructed before calling this.
			return nil, fmt.Errorf("PostgreSQL DSN is empty; please configure VIGENDA_DB_DSN or individual VIGENDA_DB_* variables")
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.DBType)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database (type: %s, dsn: %s): %w", driverName, dsn, err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database (type: %s, dsn: %s): %w", driverName, dsn, err)
	}

	// Schema migration for SQLite
	if driverName == "sqlite3" {
		// Check if the 'users' table exists to determine if schema needs to be applied.
		// This is a simplified check.
		var tableName string
		query := "SELECT name FROM sqlite_master WHERE type='table' AND name='users';"
		err := db.QueryRow(query).Scan(&tableName)
		if err == sql.ErrNoRows { // 'users' table doesn't exist, apply schema
			if err := applySQLiteSchema(db); err != nil {
				db.Close()
				return nil, fmt.Errorf("failed to apply SQLite schema: %w", err)
			}
		} else if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to check for existing SQLite tables: %w", err)
		}
	}
	// For PostgreSQL, schema migrations are assumed to be handled by external tools
	// like goose, migrate, or flyway.

	return db, nil
}

// applySQLiteSchema reads and applies the initial schema from the .sql file.
func applySQLiteSchema(db *sql.DB) error {
	var schemaBytes []byte
	var readErr error
	var attemptedPaths []string

	// Path 1: Relative to executable (for distributed binary)
	execPath, err := os.Executable()
	if err == nil {
		// Path 1: Relative to CWD (works for `go run` and if binary is run from project root)
		schemaPath1 := "internal/database/migrations/001_initial_schema.sql"
		attemptedPaths = append(attemptedPaths, schemaPath1)
		schemaBytes, readErr = os.ReadFile(schemaPath1)

		// Path 2: Relative to executable's directory (for distributed binaries)
		if readErr != nil && err == nil { // err is from os.Executable()
			execDir := filepath.Dir(execPath)
			// Assuming migrations are bundled relative to the executable
			// e.g. <exec_dir>/migrations/001_initial_schema.sql
			// or <exec_dir>/../internal/database/migrations/001_initial_schema.sql
			// For this project structure, if binary is in root:
			schemaPath2 := filepath.Join(execDir, "internal/database/migrations/001_initial_schema.sql")
			// If binary is in cmd/vigenda/vigenda:
			// schemaPath2 := filepath.Join(execDir, "..", "..", "internal", "database", "migrations", "001_initial_schema.sql")
			attemptedPaths = append(attemptedPaths, schemaPath2)
			schemaBytes, readErr = os.ReadFile(schemaPath2)

			// One more attempt for cmd/vigenda/vigenda structure
			if readErr != nil {
				schemaPath3 := filepath.Join(execDir, "..", "..", "internal", "database", "migrations", "001_initial_schema.sql")
				attemptedPaths = append(attemptedPaths, schemaPath3)
				schemaBytes, readErr = os.ReadFile(schemaPath3)
			}
		}
	} else {
		// Fallback if os.Executable() failed (e.g. in some test environments)
		schemaPathCwdOnly := "internal/database/migrations/001_initial_schema.sql"
		attemptedPaths = append(attemptedPaths, schemaPathCwdOnly)
		schemaBytes, readErr = os.ReadFile(schemaPathCwdOnly)
	}


	if readErr != nil {
		return fmt.Errorf("failed to read schema file. Attempted paths: [%s]. Last error: %w", strings.Join(attemptedPaths, ", "), readErr)
	}

	_, err = db.Exec(string(schemaBytes))
	if err != nil {
		return fmt.Errorf("failed to apply initial schema (tried paths: [%s]): %w", strings.Join(attemptedPaths, ", "), err)
	}
	// fmt.Println("SQLite database schema initialized.") // Keep this commented or use a logger
	return nil
}

// DefaultSQLitePath returns the default path for the SQLite database file.
// It places it in the user's config directory or defaults to "vigenda.db" in CWD.
// This function is now correctly named and used.
func DefaultSQLitePath() string {
	configDir, err := os.UserConfigDir()
	if err == nil {
		appConfigDir := filepath.Join(configDir, "vigenda")
		if err := os.MkdirAll(appConfigDir, 0750); err == nil { // Ensure directory exists
			return filepath.Join(appConfigDir, "vigenda.db")
		}
	}
	// Fallback to current working directory if user config dir fails
	return "vigenda.db"
}

// DefaultDbPath is added back for compatibility if it's used elsewhere,
// but it should ideally be replaced by DefaultSQLitePath for clarity.
// For now, it can just call DefaultSQLitePath.
func DefaultDbPath() string {
	return DefaultSQLitePath()
}
