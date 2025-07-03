package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log" // Adicionado para logging
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

func (r *classRepository) UpdateClass(ctx context.Context, class *models.Class) error {
	query := `UPDATE classes SET name = ?, subject_id = ?, updated_at = ?
              WHERE id = ? AND user_id = ?`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, class.Name, class.SubjectID, now, class.ID, class.UserID)
	if err != nil {
		return fmt.Errorf("classRepository.UpdateClass: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("classRepository.UpdateClass: checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("classRepository.UpdateClass: no class found with ID %d or user mismatch", class.ID)
	}
	return nil
}

func (r *classRepository) DeleteClass(ctx context.Context, classID int64, userID int64) error {
	query := `DELETE FROM classes WHERE id = ? AND user_id = ?`
	result, err := r.db.ExecContext(ctx, query, classID, userID)
	if err != nil {
		return fmt.Errorf("classRepository.DeleteClass: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("classRepository.DeleteClass: checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("classRepository.DeleteClass: no class found with ID %d or user mismatch", classID)
	}
	return nil
}

func (r *classRepository) AddStudent(ctx context.Context, student *models.Student) (int64, error) {
	query := `INSERT INTO students (class_id, enrollment_id, full_name, status, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	var enrollmentID sql.NullString
	if student.EnrollmentID != "" {
		enrollmentID.String = student.EnrollmentID
		enrollmentID.Valid = true
	}
	result, err := r.db.ExecContext(ctx, query, student.ClassID, enrollmentID, student.FullName, student.Status, now, now)
	if err != nil {
		return 0, fmt.Errorf("classRepository.AddStudent: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("classRepository.AddStudent: failed to get last insert ID: %w", err)
	}
	return id, nil
}

func (r *classRepository) GetStudentByID(ctx context.Context, studentID int64) (*models.Student, error) {
	query := `SELECT id, class_id, enrollment_id, full_name, status, created_at, updated_at
			  FROM students WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, studentID)
	student := &models.Student{}
	var enrollmentID sql.NullString
	err := row.Scan(
		&student.ID,
		&student.ClassID,
		&enrollmentID,
		&student.FullName,
		&student.Status,
		&student.CreatedAt,
		&student.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("classRepository.GetStudentByID: no student found with ID %d", studentID)
		}
		return nil, fmt.Errorf("classRepository.GetStudentByID: %w", err)
	}
	if enrollmentID.Valid {
		student.EnrollmentID = enrollmentID.String
	} else {
		student.EnrollmentID = ""
	}
	return student, nil
}

func (r *classRepository) UpdateStudent(ctx context.Context, student *models.Student) error {
	query := `UPDATE students SET full_name = ?, enrollment_id = ?, status = ?, updated_at = ?
              WHERE id = ? AND class_id = ?` // Assuming class_id cannot be changed this way
	now := time.Now()
	var enrollmentID sql.NullString
	if student.EnrollmentID != "" {
		enrollmentID.String = student.EnrollmentID
		enrollmentID.Valid = true
	}
	result, err := r.db.ExecContext(ctx, query, student.FullName, enrollmentID, student.Status, now, student.ID, student.ClassID)
	if err != nil {
		return fmt.Errorf("classRepository.UpdateStudent: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("classRepository.UpdateStudent: checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("classRepository.UpdateStudent: no student found with ID %d or class_id mismatch", student.ID)
	}
	return nil
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

func (r *classRepository) DeleteStudent(ctx context.Context, studentID int64, classID int64) error {
	query := `DELETE FROM students WHERE id = ? AND class_id = ?`
	result, err := r.db.ExecContext(ctx, query, studentID, classID)
	if err != nil {
		return fmt.Errorf("classRepository.DeleteStudent: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("classRepository.DeleteStudent: checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("classRepository.DeleteStudent: no student found with ID %d in class ID %d", studentID, classID)
	}
	return nil
}

func (r *classRepository) ListAllClasses(ctx context.Context) ([]models.Class, error) {
	log.Println("Repository: classRepository.ListAllClasses - Chamado.")
	query := `SELECT id, user_id, subject_id, name, created_at, updated_at FROM classes ORDER BY name ASC`
	log.Printf("Repository: classRepository.ListAllClasses - Executando query: %s", query)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Repository: classRepository.ListAllClasses - Erro ao executar query: %v", err)
		return nil, fmt.Errorf("classRepository.ListAllClasses: query failed: %w", err)
	}
	defer rows.Close()

	var classes []models.Class
	log.Println("Repository: classRepository.ListAllClasses - Lendo linhas do resultado...")
	for rows.Next() {
		var class models.Class
		if err := rows.Scan(&class.ID, &class.UserID, &class.SubjectID, &class.Name, &class.CreatedAt, &class.UpdatedAt); err != nil {
			log.Printf("Repository: classRepository.ListAllClasses - Erro ao escanear linha: %v", err)
			return nil, fmt.Errorf("classRepository.ListAllClasses: scan failed: %w", err)
		}
		classes = append(classes, class)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Repository: classRepository.ListAllClasses - Erro após iteração das linhas: %v", err)
		return nil, fmt.Errorf("classRepository.ListAllClasses: rows error: %w", err)
	}

	log.Printf("Repository: classRepository.ListAllClasses - Query bem-sucedida. %d turmas lidas.", len(classes))
	return classes, nil
}

func (r *classRepository) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
	query := `SELECT id, class_id, enrollment_id, full_name, status, created_at, updated_at
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
			&student.CreatedAt,
			&student.UpdatedAt,
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

// GetTodaysLessonsByUserID retrieves all lessons associated with a given userID
// that are scheduled for the specified date (ignoring the time part).
// It joins lessons with classes to filter by user_id.
func (r *classRepository) GetTodaysLessonsByUserID(ctx context.Context, userID int64, today time.Time) ([]models.Lesson, error) {
	// Query needs to join with classes to filter by user_id,
	// as lessons table itself does not have user_id.
	// It filters lessons scheduled for the given date (today).
	query := `
        SELECT l.id, l.class_id, l.title, l.plan_content, l.scheduled_at
        FROM lessons l
        JOIN classes c ON l.class_id = c.id
        WHERE c.user_id = ?
          AND DATE(l.scheduled_at) = DATE(?)
        ORDER BY l.scheduled_at ASC`

	// today.Format("2006-01-02") could be used if DATE(?) doesn't work as expected with full timestamp.
	// However, SQLite's DATE() function should correctly extract the date part from a timestamp.
	rows, err := r.db.QueryContext(ctx, query, userID, today)
	if err != nil {
		return nil, fmt.Errorf("classRepository.GetTodaysLessonsByUserID: query failed: %w", err)
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		lesson := models.Lesson{}
		var planContent sql.NullString // plan_content can be NULL

		err := rows.Scan(
			&lesson.ID,
			&lesson.ClassID,
			&lesson.Title,
			&planContent,
			&lesson.ScheduledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("classRepository.GetTodaysLessonsByUserID: scanning lesson: %w", err)
		}
		if planContent.Valid {
			lesson.PlanContent = planContent.String
		} else {
			lesson.PlanContent = "" // Default to empty string if NULL
		}
		lessons = append(lessons, lesson)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("classRepository.GetTodaysLessonsByUserID: iterating rows: %w", err)
	}
	return lessons, nil
}
