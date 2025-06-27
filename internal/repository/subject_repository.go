package repository

import (
	"context"
	"database/sql"
	"fmt"
	// "time" // Not used anymore
	"vigenda/internal/models"
)

type subjectRepository struct {
	db *sql.DB
}

func NewSubjectRepository(db *sql.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) GetOrCreateByNameAndUser(ctx context.Context, name string, userID int64) (models.Subject, error) {
	// Try to find the subject first
	// Removed created_at, updated_at from SELECT
	queryGet := `SELECT id, user_id, name FROM subjects WHERE name = ? AND user_id = ?`
	row := r.db.QueryRowContext(ctx, queryGet, name, userID)

	var subject models.Subject
	// Removed subject.CreatedAt, subject.UpdatedAt from Scan
	err := row.Scan(&subject.ID, &subject.UserID, &subject.Name)

	if err == nil {
		// Subject found
		return subject, nil
	}

	if err != sql.ErrNoRows {
		// An actual error occurred during scan or query
		return models.Subject{}, fmt.Errorf("subjectRepository.GetOrCreateByNameAndUser: getting subject: %w", err)
	}

	// Subject not found, create it
	// Removed created_at, updated_at from INSERT
	queryCreate := `INSERT INTO subjects (user_id, name) VALUES (?, ?)`
	// now := time.Now() // Not used
	result, err := r.db.ExecContext(ctx, queryCreate, userID, name)
	if err != nil {
		return models.Subject{}, fmt.Errorf("subjectRepository.GetOrCreateByNameAndUser: creating subject: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Subject{}, fmt.Errorf("subjectRepository.GetOrCreateByNameAndUser: getting last insert ID: %w", err)
	}

	// Removed CreatedAt, UpdatedAt from struct literal
	return models.Subject{
		ID:     id,
		UserID: userID,
		Name:   name,
	}, nil
}
