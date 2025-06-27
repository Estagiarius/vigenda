package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"vigenda/internal/models"
)

type assessmentRepository struct {
	db *sql.DB
}

func NewAssessmentRepository(db *sql.DB) AssessmentRepository {
	return &assessmentRepository{db: db}
}

func (r *assessmentRepository) CreateAssessment(ctx context.Context, assessment *models.Assessment) (int64, error) {
	query := `INSERT INTO assessments (user_id, class_id, name, term, weight, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, assessment.UserID, assessment.ClassID, assessment.Name, assessment.Term, assessment.Weight, now, now)
	if err != nil {
		return 0, fmt.Errorf("assessmentRepository.CreateAssessment: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("assessmentRepository.CreateAssessment: failed to get last insert ID: %w", err)
	}
	return id, nil
}

func (r *assessmentRepository) GetAssessmentByID(ctx context.Context, assessmentID int64) (*models.Assessment, error) {
	query := `SELECT id, user_id, class_id, name, term, weight, created_at, updated_at
              FROM assessments WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, assessmentID)
	assessment := &models.Assessment{}
	err := row.Scan(
		&assessment.ID,
		&assessment.UserID,
		&assessment.ClassID,
		&assessment.Name,
		&assessment.Term,
		&assessment.Weight,
		&assessment.CreatedAt,
		&assessment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("assessmentRepository.GetAssessmentByID: no assessment found with ID %d", assessmentID)
		}
		return nil, fmt.Errorf("assessmentRepository.GetAssessmentByID: %w", err)
	}
	return assessment, nil
}

func (r *assessmentRepository) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
	query := `SELECT id, class_id, user_id, call_number, full_name, status, created_at, updated_at
              FROM students WHERE class_id = ? AND status = 'ativo' ORDER BY full_name` // Assuming we only want active students for grading
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, fmt.Errorf("assessmentRepository.GetStudentsByClassID: %w", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		student := models.Student{}
		var callNumber sql.NullInt64
		err := rows.Scan(
			&student.ID,
			&student.ClassID,
			&student.UserID,
			&callNumber,
			&student.FullName,
			&student.Status,
			&student.CreatedAt,
			&student.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("assessmentRepository.GetStudentsByClassID: scanning student: %w", err)
		}
		if callNumber.Valid {
			student.CallNumber = int(callNumber.Int64)
		}
		students = append(students, student)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("assessmentRepository.GetStudentsByClassID: iterating rows: %w", err)
	}
	return students, nil
}

func (r *assessmentRepository) EnterGrade(ctx context.Context, grade *models.Grade) error {
	// Use INSERT OR REPLACE (UPSERT) to handle new grades and updates
	query := `INSERT INTO grades (assessment_id, student_id, user_id, grade, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?)
              ON CONFLICT(assessment_id, student_id) DO UPDATE SET
              grade = excluded.grade,
              updated_at = excluded.updated_at`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, grade.AssessmentID, grade.StudentID, grade.UserID, grade.Grade, now, now)
	if err != nil {
		return fmt.Errorf("assessmentRepository.EnterGrade: %w", err)
	}
	return nil
}

func (r *assessmentRepository) GetGradesByClassID(ctx context.Context, classID int64) ([]models.Grade, []models.Assessment, []models.Student, error) {
	// This is a complex query. It needs to fetch all grades for all assessments in a class,
	// along with assessment details (for weight) and student details.

	// 1. Get all assessments for the class
	assessmentsQuery := `SELECT id, user_id, class_id, name, term, weight FROM assessments WHERE class_id = ?`
	assessmentRows, err := r.db.QueryContext(ctx, assessmentsQuery, classID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: fetching assessments: %w", err)
	}
	defer assessmentRows.Close()

	var assessments []models.Assessment
	assessmentMap := make(map[int64]models.Assessment)
	for assessmentRows.Next() {
		var a models.Assessment
		if err := assessmentRows.Scan(&a.ID, &a.UserID, &a.ClassID, &a.Name, &a.Term, &a.Weight); err != nil {
			return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: scanning assessment: %w", err)
		}
		assessments = append(assessments, a)
		assessmentMap[a.ID] = a
	}
	if err = assessmentRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: iterating assessments: %w", err)
	}

	// 2. Get all students for the class
	studentsQuery := `SELECT id, class_id, user_id, full_name, status FROM students WHERE class_id = ?`
	studentRows, err := r.db.QueryContext(ctx, studentsQuery, classID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: fetching students: %w", err)
	}
	defer studentRows.Close()

	var students []models.Student
	studentMap := make(map[int64]models.Student)
	for studentRows.Next() {
		var s models.Student
		if err := studentRows.Scan(&s.ID, &s.ClassID, &s.UserID, &s.FullName, &s.Status); err != nil {
			return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: scanning student: %w", err)
		}
		students = append(students, s)
		studentMap[s.ID] = s
	}
	if err = studentRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: iterating students: %w", err)
	}

	// 3. Get all grades for those assessments and students
	// This assumes student_id and assessment_id in 'grades' table are foreign keys
	// to 'students' and 'assessments' tables respectively.
	gradesQuery := `
        SELECT g.id, g.assessment_id, g.student_id, g.user_id, g.grade
        FROM grades g
        JOIN assessments a ON g.assessment_id = a.id
        WHERE a.class_id = ?`
	gradeRows, err := r.db.QueryContext(ctx, gradesQuery, classID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: fetching grades: %w", err)
	}
	defer gradeRows.Close()

	var grades []models.Grade
	for gradeRows.Next() {
		var g models.Grade
		if err := gradeRows.Scan(&g.ID, &g.AssessmentID, &g.StudentID, &g.UserID, &g.Grade); err != nil {
			return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: scanning grade: %w", err)
		}
		grades = append(grades, g)
	}
	if err = gradeRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: iterating grades: %w", err)
	}

	return grades, assessments, students, nil
}
