package repository

import (
	"context"
	"database/sql"
	"fmt"
	// "time" // Not used anymore
	"vigenda/internal/models"
)

type classRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) CreateClass(ctx context.Context, class *models.Class) (int64, error) {
	// Removed created_at, updated_at as they are not in the 'classes' table schema
	query := `INSERT INTO classes (user_id, subject_id, name)
              VALUES (?, ?, ?)`
	// now := time.Now() // Not used
	result, err := r.db.ExecContext(ctx, query, class.UserID, class.SubjectID, class.Name)
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
	// Removed created_at, updated_at
	query := `SELECT id, user_id, subject_id, name
              FROM classes WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	class := &models.Class{}
	err := row.Scan(
		&class.ID,
		&class.UserID,
		&class.SubjectID,
		&class.Name,
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
	// Corrected query: call_number -> enrollment_id. Removed user_id, created_at, updated_at.
	query := `INSERT INTO students (class_id, enrollment_id, full_name, status)
              VALUES (?, ?, ?, ?)`
	// now := time.Now() // Not used

	var enrollmentID sql.NullString
	if student.EnrollmentID != "" {
		enrollmentID.String = student.EnrollmentID
		enrollmentID.Valid = true
	}
	// student.UserID is not in the 'students' table schema
	result, err := r.db.ExecContext(ctx, query, student.ClassID, enrollmentID, student.FullName, student.Status)
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
	// Removed updated_at as it's not in the 'students' table schema
	query := `UPDATE students SET status = ? WHERE id = ?`
	// now := time.Now() // Not used
	result, err := r.db.ExecContext(ctx, query, status, studentID)
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

func (r *classRepository) ListAllClasses(ctx context.Context) ([]models.Class, error) {
	query := `SELECT id, user_id, subject_id, name FROM classes`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("classRepository.ListAllClasses: query failed: %w", err)
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var class models.Class
		if err := rows.Scan(&class.ID, &class.UserID, &class.SubjectID, &class.Name); err != nil {
			return nil, fmt.Errorf("classRepository.ListAllClasses: scan failed: %w", err)
		}
		classes = append(classes, class)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("classRepository.ListAllClasses: rows error: %w", err)
	}

	return classes, nil
}

func (r *classRepository) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
	query := `SELECT id, class_id, enrollment_id, full_name, status
              FROM students
              WHERE class_id = ?
              ORDER BY full_name ASC` // Ordenar por nome completo
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, fmt.Errorf("classRepository.GetStudentsByClassID: query failed: %w", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		var enrollmentID sql.NullString // enrollment_id pode ser NULL
		if err := rows.Scan(
			&student.ID,
			&student.ClassID,
			&enrollmentID, // Scan para sql.NullString
			&student.FullName,
			&student.Status,
		); err != nil {
			return nil, fmt.Errorf("classRepository.GetStudentsByClassID: scan failed: %w", err)
		}
		if enrollmentID.Valid {
			student.EnrollmentID = enrollmentID.String
		} else {
			student.EnrollmentID = "" // Ou algum outro valor padrão se preferir
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("classRepository.GetStudentsByClassID: rows error: %w", err)
	}

	return students, nil
}
