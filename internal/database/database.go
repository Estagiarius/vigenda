// Package database handles the connection to the SQLite database
// and the execution of database migrations.
package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	// Import the SQLite driver
	_ "github.com/mattn/go-sqlite3"
	// Import the PostgreSQL driver
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DBConfig_database holds the configuration for database connection.
// Renamed to avoid conflict with DBConfig in connection.go
type DBConfig_database struct {
	DBType string // "sqlite" or "postgres"
	DSN    string // Data Source Name
}

// GetDBConnection_database establishes a connection to the database based on the provided config.
// It also ensures that migrations are applied for SQLite.
// Renamed to avoid conflict
func GetDBConnection_database(config DBConfig_database) (*sql.DB, error) {
	var db *sql.DB
	var err error

	driverName := ""
	switch config.DBType {
	case "sqlite":
		driverName = "sqlite3"
		// Ensure the directory for the SQLite file exists
		dbDir := filepath.Dir(config.DSN)
		if _, statErr := os.Stat(dbDir); os.IsNotExist(statErr) {
			if mkdirErr := os.MkdirAll(dbDir, 0755); mkdirErr != nil {
				return nil, fmt.Errorf("failed to create database directory %s: %w", dbDir, mkdirErr)
			}
		}
	case "postgres":
		driverName = "postgres"
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.DBType)
	}

	db, err = sql.Open(driverName, config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database (type: %s, DSN: %s): %w", config.DBType, config.DSN, err)
	}

	err = db.Ping()
	if err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("failed to ping database (type: %s, DSN: %s): %w", config.DBType, config.DSN, err)
	}

	log.Printf("Successfully connected to %s database.", config.DBType)

	// Apply migrations only for SQLite, as specified in AGENTS.md
	if config.DBType == "sqlite" {
		err = applyMigrations_database(db) // Renamed
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to apply migrations for SQLite: %w", err)
		}
		// Seed data after applying migrations for SQLite, if the database is new/empty
		err = SeedData(db) // Assuming SeedData is correctly defined and doesn't conflict
		if err != nil {
			db.Close()
			// Log the error but don't necessarily fail the connection if seeding fails
			log.Printf("Warning: failed to seed database with example data: %v", err)
			// return nil, fmt.Errorf("failed to seed database: %w", err) // Optional: make seeding failure critical
		}

	}

	return db, nil
}

// applyMigrations_database executes all .sql files found in the migrations directory.
// Renamed to avoid conflict
func applyMigrations_database(db *sql.DB) error {
	log.Println("Applying database migrations for SQLite...")
	migrationFiles, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range migrationFiles {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			log.Printf("Applying migration: %s", file.Name())
			content, err := migrationsFS.ReadFile(filepath.Join("migrations", file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
			}

			_, err = db.Exec(string(content))
			if err != nil {
				statements := strings.Split(string(content), ";")
				tx, txErr := db.Begin()
				if txErr != nil {
					return fmt.Errorf("failed to begin transaction for migration %s: %w", file.Name(), txErr)
				}
				for _, stmt := range statements {
					trimmedStmt := strings.TrimSpace(stmt)
					if trimmedStmt == "" {
						continue
					}
					if _, execErr := tx.Exec(trimmedStmt); execErr != nil {
						tx.Rollback()
						return fmt.Errorf("error executing statement from migration file %s: \"%s\". Error: %w", file.Name(), trimmedStmt, execErr)
					}
				}
				if commitErr := tx.Commit(); commitErr != nil {
					return fmt.Errorf("failed to commit transaction for migration %s: %w", file.Name(), commitErr)
				}
				log.Printf("Migration %s applied successfully (executed statement by statement).", file.Name())
			} else {
				log.Printf("Migration %s applied successfully (executed as a single block).", file.Name())
			}
		}
	}
	log.Println("All migrations applied successfully.")
	return nil
}

// DefaultSQLitePath_database returns the default path for the SQLite database file.
// Renamed to avoid conflict
func DefaultSQLitePath_database() string {
	configDir, err := os.UserConfigDir()
	if err == nil {
		vigendaConfigDir := filepath.Join(configDir, "vigenda")
		if _, statErr := os.Stat(vigendaConfigDir); os.IsNotExist(statErr) {
			// Correctly use mkdirErr for the log message
			if mkdirErr := os.MkdirAll(vigendaConfigDir, 0755); mkdirErr == nil {
				return filepath.Join(vigendaConfigDir, "vigenda.db")
			} else { // Added else to handle MkdirAll error
				log.Printf("Warning: Failed to create config directory %s: %v. Using current directory.", vigendaConfigDir, mkdirErr)
			}
		} else if statErr == nil {
			return filepath.Join(vigendaConfigDir, "vigenda.db")
		} else {
			log.Printf("Warning: Error checking config directory %s: %v. Using current directory.", vigendaConfigDir, statErr)
		}
	} else {
		log.Printf("Warning: Could not determine user config directory: %v. Using current directory.", err)
	}
	return "vigenda.db"
}
