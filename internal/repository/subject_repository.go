package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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
	queryGet := `SELECT id, user_id, name, created_at, updated_at FROM subjects WHERE name = ? AND user_id = ?`
	row := r.db.QueryRowContext(ctx, queryGet, name, userID)

	var subject models.Subject
	err := row.Scan(&subject.ID, &subject.UserID, &subject.Name, &subject.CreatedAt, &subject.UpdatedAt)

	if err == nil {
		// Subject found
		return subject, nil
	}

	if err != sql.ErrNoRows {
		// An actual error occurred during scan or query
		return models.Subject{}, fmt.Errorf("subjectRepository.GetOrCreateByNameAndUser: getting subject: %w", err)
	}

	// Subject not found, create it
	queryCreate := `INSERT INTO subjects (user_id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, queryCreate, userID, name, now, now)
	if err != nil {
		return models.Subject{}, fmt.Errorf("subjectRepository.GetOrCreateByNameAndUser: creating subject: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Subject{}, fmt.Errorf("subjectRepository.GetOrCreateByNameAndUser: getting last insert ID: %w", err)
	}

	return models.Subject{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
