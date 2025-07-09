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

// applySQLiteSchema reads and applies the initial schema from the embedded .sql file.
func applySQLiteSchema(db *sql.DB) error {
	schemaBytes, err := migrationsFS.ReadFile("migrations/001_initial_schema.sql")
	if err != nil {
		// It's useful to know which file failed, even if there's only one for now.
		return fmt.Errorf("failed to read embedded schema file migrations/001_initial_schema.sql: %w", err)
	}

	_, err = db.Exec(string(schemaBytes))
	if err != nil {
		return fmt.Errorf("failed to apply initial schema from embedded file: %w", err)
	}
	// fmt.Println("SQLite database schema initialized from embedded file.") // Keep this commented or use a logger
	return nil
}

// DefaultSQLitePath returns the default path for the SQLite database file.
// It places it in the user's config directory or defaults to "vigenda.db" in CWD.
// This function is now correctly named and used.
// To avoid import cycle, this function cannot use DefaultSQLitePath_database from database.go directly.
// It needs to have its own implementation or database.go's version needs to be made available differently.
// For now, we will duplicate the logic slightly, assuming this is acceptable.
// A better solution might involve a shared utility package or rethinking package responsibilities.
func DefaultSQLitePath() string {
	// This implementation detail should ideally be centralized.
	// For now, replicating the logic from database.go's DefaultSQLitePath_database
	// to avoid import cycles or major refactoring.
	// This is a common issue when splitting database setup logic.
	// In a real scenario, further refactoring might be needed.
	userConfigDir, err := os.UserConfigDir()
	if err == nil {
		appConfigDir := filepath.Join(userConfigDir, "vigenda")
		// Ensure the directory exists. Using 0755 for broader compatibility.
		if mkdirErr := os.MkdirAll(appConfigDir, 0755); mkdirErr == nil {
			return filepath.Join(appConfigDir, "vigenda.db")
		}
		// If MkdirAll fails, log it or handle as appropriate, then fallback.
		// log.Printf("Warning: Failed to create config directory %s: %v. Using current directory.", appConfigDir, mkdirErr)
	} else {
		// log.Printf("Warning: Could not determine user config directory: %v. Using current directory.", err)
	}
	return "vigenda.db" // Fallback to current working directory
}


// DefaultDbPath is added back for compatibility if it's used elsewhere,
// but it should ideally be replaced by DefaultSQLitePath for clarity.
// For now, it can just call DefaultSQLitePath.
func DefaultDbPath() string {
	return DefaultSQLitePath()
}
