package repository

import (
	"context"
	"database/sql"
	"fmt"
	// "time" // Not strictly needed here anymore
	"vigenda/internal/models"
)

type assessmentRepository struct {
	db *sql.DB
}

func NewAssessmentRepository(db *sql.DB) AssessmentRepository {
	return &assessmentRepository{db: db}
}

func (r *assessmentRepository) CreateAssessment(ctx context.Context, assessment *models.Assessment) (int64, error) {
	query := `INSERT INTO assessments (class_id, name, term, weight, assessment_date)
              VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, assessment.ClassID, assessment.Name, assessment.Term, assessment.Weight, assessment.AssessmentDate)
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
	query := `SELECT id, class_id, name, term, weight, assessment_date
              FROM assessments WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, assessmentID)
	assessment := &models.Assessment{}
	err := row.Scan(
		&assessment.ID,
		&assessment.ClassID,
		&assessment.Name,
		&assessment.Term,
		&assessment.Weight,
		&assessment.AssessmentDate,
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
	query := `SELECT id, class_id, enrollment_id, full_name, status
              FROM students WHERE class_id = ? AND status = 'ativo' ORDER BY full_name`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, fmt.Errorf("assessmentRepository.GetStudentsByClassID: %w", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		student := models.Student{}
		var enrollmentID sql.NullString
		err := rows.Scan(
			&student.ID,
			&student.ClassID,
			&enrollmentID,
			&student.FullName,
			&student.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("assessmentRepository.GetStudentsByClassID: scanning student: %w", err)
		}
		if enrollmentID.Valid {
			student.EnrollmentID = enrollmentID.String
		}
		students = append(students, student)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("assessmentRepository.GetStudentsByClassID: iterating rows: %w", err)
	}
	return students, nil
}

func (r *assessmentRepository) EnterGrade(ctx context.Context, grade *models.Grade) error {
	query := `INSERT INTO grades (assessment_id, student_id, grade)
              VALUES (?, ?, ?)
              ON CONFLICT(assessment_id, student_id) DO UPDATE SET
              grade = excluded.grade`
	_, err := r.db.ExecContext(ctx, query, grade.AssessmentID, grade.StudentID, grade.Grade)
	if err != nil {
		return fmt.Errorf("assessmentRepository.EnterGrade: %w", err)
	}
	return nil
}

func (r *assessmentRepository) GetGradesByClassID(ctx context.Context, classID int64) ([]models.Grade, []models.Assessment, []models.Student, error) {
	assessmentsQuery := `SELECT id, class_id, name, term, weight, assessment_date FROM assessments WHERE class_id = ?`
	assessmentRows, err := r.db.QueryContext(ctx, assessmentsQuery, classID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: fetching assessments: %w", err)
	}
	defer assessmentRows.Close()

	var assessments []models.Assessment
	assessmentMap := make(map[int64]models.Assessment)
	for assessmentRows.Next() {
		var a models.Assessment
		if err := assessmentRows.Scan(&a.ID, &a.ClassID, &a.Name, &a.Term, &a.Weight, &a.AssessmentDate); err != nil {
			return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: scanning assessment: %w", err)
		}
		assessments = append(assessments, a)
		assessmentMap[a.ID] = a
	}
	if err = assessmentRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: iterating assessments: %w", err)
	}

	studentsQuery := `SELECT id, class_id, enrollment_id, full_name, status FROM students WHERE class_id = ?`
	studentRows, err := r.db.QueryContext(ctx, studentsQuery, classID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: fetching students: %w", err)
	}
	defer studentRows.Close()

	var students []models.Student
	studentMap := make(map[int64]models.Student)
	for studentRows.Next() {
		var s models.Student
		var enrollmentID sql.NullString
		if err := studentRows.Scan(&s.ID, &s.ClassID, &enrollmentID, &s.FullName, &s.Status); err != nil {
			return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: scanning student: %w", err)
		}
		if enrollmentID.Valid {
			s.EnrollmentID = enrollmentID.String
		}
		students = append(students, s)
		studentMap[s.ID] = s
	}
	if err = studentRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: iterating students: %w", err)
	}

	gradesQuery := `
        SELECT g.id, g.assessment_id, g.student_id, g.grade
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
		if err := gradeRows.Scan(&g.ID, &g.AssessmentID, &g.StudentID, &g.Grade); err != nil {
			return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: scanning grade: %w", err)
		}
		grades = append(grades, g)
	}
	if err = gradeRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("assessmentRepository.GetGradesByClassID: iterating grades: %w", err)
	}

	return grades, assessments, students, nil
}

func (r *assessmentRepository) ListAllAssessments(ctx context.Context) ([]models.Assessment, error) {
	query := `SELECT id, class_id, name, term, weight, assessment_date FROM assessments`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("assessmentRepository.ListAllAssessments: query failed: %w", err)
	}
	defer rows.Close()

	var assessments []models.Assessment
	for rows.Next() {
		var asm models.Assessment
		if err := rows.Scan(&asm.ID, &asm.ClassID, &asm.Name, &asm.Term, &asm.Weight, &asm.AssessmentDate); err != nil {
			return nil, fmt.Errorf("assessmentRepository.ListAllAssessments: scan failed: %w", err)
		}
		assessments = append(assessments, asm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("assessmentRepository.ListAllAssessments: rows error: %w", err)
	}

	return assessments, nil
}

func (r *assessmentRepository) DeleteAssessment(ctx context.Context, assessmentID int64) error {
	query := `DELETE FROM assessments WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, assessmentID)
	if err != nil {
		return fmt.Errorf("assessmentRepository.DeleteAssessment: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("assessmentRepository.DeleteAssessment: could not get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("assessmentRepository.DeleteAssessment: no assessment found with ID %d", assessmentID)
	}
	return nil
}
