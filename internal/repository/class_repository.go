package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"vigenda/internal/models"
)

type classRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) CreateClass(ctx context.Context, class *models.Class) (int64, error) {
	query := `INSERT INTO classes (user_id, subject_id, name, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, class.UserID, class.SubjectID, class.Name, now, now)
	if err != nil {
		return 0, fmt.Errorf("classRepository.CreateClass: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("classRepository.CreateClass: failed to get last insert ID: %w", err)
	}
	return id, nil
}

func (r *classRepository) GetClassByID(ctx context.Context, id int64) (*models.Class, error) {
	query := `SELECT id, user_id, subject_id, name, created_at, updated_at
              FROM classes WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	class := &models.Class{}
	err := row.Scan(
		&class.ID,
		&class.UserID,
		&class.SubjectID,
		&class.Name,
		&class.CreatedAt,
		&class.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("classRepository.GetClassByID: no class found with ID %d", id)
		}
		return nil, fmt.Errorf("classRepository.GetClassByID: %w", err)
	}
	return class, nil
}

func (r *classRepository) AddStudent(ctx context.Context, student *models.Student) (int64, error) {
	query := `INSERT INTO students (class_id, user_id, call_number, full_name, status, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()

	var callNumber sql.NullInt64
	if student.CallNumber != 0 { // Assuming 0 is not a valid call number, adjust if it can be
		callNumber.Int64 = int64(student.CallNumber)
		callNumber.Valid = true
	}

	result, err := r.db.ExecContext(ctx, query, student.ClassID, student.UserID, callNumber, student.FullName, student.Status, now, now)
	if err != nil {
		return 0, fmt.Errorf("classRepository.AddStudent: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("classRepository.AddStudent: failed to get last insert ID: %w", err)
	}
	return id, nil
}

func (r *classRepository) UpdateStudentStatus(ctx context.Context, studentID int64, status string) error {
	query := `UPDATE students SET status = ?, updated_at = ? WHERE id = ?`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, status, now, studentID)
	if err != nil {
		return fmt.Errorf("classRepository.UpdateStudentStatus: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("classRepository.UpdateStudentStatus: checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("classRepository.UpdateStudentStatus: no student found with ID %d", studentID)
	}
	return nil
}
