package repository

import (
	"context"
	"database/sql"
	"fmt"
	"vigenda/internal/models"
)

// AssessmentRepository define a interface para operações de acesso a dados relacionadas a 'assessments' (avaliações) e 'grades' (notas).
// Esta é uma versão atualizada para suportar as novas funcionalidades do serviço.
type AssessmentRepository interface {
	CreateAssessment(ctx context.Context, assessment *models.Assessment) (int64, error)
	GetAssessmentByID(ctx context.Context, assessmentID int64) (*models.Assessment, error)
	EnterGrade(ctx context.Context, grade *models.Grade) error
	GetGradesByClassID(ctx context.Context, classID int64) ([]models.Grade, []models.Assessment, []models.Student, error)
	ListAllAssessments(ctx context.Context) ([]models.Assessment, error)
	GetGradesByAssessmentID(ctx context.Context, assessmentID int64) ([]models.Grade, error)
	UpdateAssessment(ctx context.Context, assessment *models.Assessment) error
	DeleteAssessment(ctx context.Context, assessmentID int64) error
}

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
	query := `SELECT id, class_id, name, term, weight, assessment_date FROM assessments WHERE id = ?`
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

func (r *assessmentRepository) UpdateAssessment(ctx context.Context, assessment *models.Assessment) error {
	query := `UPDATE assessments SET name = ?, class_id = ?, term = ?, weight = ?, assessment_date = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, assessment.Name, assessment.ClassID, assessment.Term, assessment.Weight, assessment.AssessmentDate, assessment.ID)
	if err != nil {
		return fmt.Errorf("assessmentRepository.UpdateAssessment: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("assessmentRepository.UpdateAssessment: checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("assessmentRepository.UpdateAssessment: no assessment found with ID %d", assessment.ID)
	}
	return nil
}

func (r *assessmentRepository) DeleteAssessment(ctx context.Context, assessmentID int64) error {
	query := `DELETE FROM assessments WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, assessmentID)
	if err != nil {
		return fmt.Errorf("assessmentRepository.DeleteAssessment: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("assessmentRepository.DeleteAssessment: checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("assessmentRepository.DeleteAssessment: no assessment found with ID %d", assessmentID)
	}
	return nil
}

func (r *assessmentRepository) EnterGrade(ctx context.Context, grade *models.Grade) error {
	// Upsert logic: update if grade for student/assessment exists, otherwise insert.
	query := `INSERT INTO grades (assessment_id, student_id, grade) VALUES (?, ?, ?)
              ON CONFLICT(assessment_id, student_id) DO UPDATE SET grade = excluded.grade`
	_, err := r.db.ExecContext(ctx, query, grade.AssessmentID, grade.StudentID, grade.Grade)
	if err != nil {
		return fmt.Errorf("assessmentRepository.EnterGrade: %w", err)
	}
	return nil
}

func (r *assessmentRepository) GetGradesByAssessmentID(ctx context.Context, assessmentID int64) ([]models.Grade, error) {
	query := `SELECT id, assessment_id, student_id, grade FROM grades WHERE assessment_id = ?`
	rows, err := r.db.QueryContext(ctx, query, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("assessmentRepository.GetGradesByAssessmentID: query failed: %w", err)
	}
	defer rows.Close()

	var grades []models.Grade
	for rows.Next() {
		var grade models.Grade
		if err := rows.Scan(&grade.ID, &grade.AssessmentID, &grade.StudentID, &grade.Grade); err != nil {
			return nil, fmt.Errorf("assessmentRepository.GetGradesByAssessmentID: scan failed: %w", err)
		}
		grades = append(grades, grade)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("assessmentRepository.GetGradesByAssessmentID: rows error: %w", err)
	}
	return grades, nil
}

// GetGradesByClassID is a complex query to fetch all related data for average calculation.
func (r *assessmentRepository) GetGradesByClassID(ctx context.Context, classID int64) ([]models.Grade, []models.Assessment, []models.Student, error) {
	// 1. Get all assessments for the class
	assessmentsQuery := `SELECT id, class_id, name, term, weight, assessment_date FROM assessments WHERE class_id = ?`
	assessmentRows, err := r.db.QueryContext(ctx, assessmentsQuery, classID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("repo.GetGradesByClassID: fetching assessments: %w", err)
	}
	defer assessmentRows.Close()

	var assessments []models.Assessment
	var assessmentIDs []interface{}
	for assessmentRows.Next() {
		var a models.Assessment
		if err := assessmentRows.Scan(&a.ID, &a.ClassID, &a.Name, &a.Term, &a.Weight, &a.AssessmentDate); err != nil {
			return nil, nil, nil, fmt.Errorf("repo.GetGradesByClassID: scanning assessment: %w", err)
		}
		assessments = append(assessments, a)
		assessmentIDs = append(assessmentIDs, a.ID)
	}
	if err = assessmentRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("repo.GetGradesByClassID: assessment rows error: %w", err)
	}
	if len(assessments) == 0 {
		return []models.Grade{}, []models.Assessment{}, []models.Student{}, nil
	}

	// 2. Get all grades for those assessments
	gradesQuery := `SELECT id, assessment_id, student_id, grade FROM grades WHERE assessment_id IN (?` + sqlRepeat(len(assessmentIDs)-1) + `)`
	gradeRows, err := r.db.QueryContext(ctx, gradesQuery, assessmentIDs...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("repo.GetGradesByClassID: fetching grades: %w", err)
	}
	defer gradeRows.Close()

	var grades []models.Grade
	for gradeRows.Next() {
		var g models.Grade
		if err := gradeRows.Scan(&g.ID, &g.AssessmentID, &g.StudentID, &g.Grade); err != nil {
			return nil, nil, nil, fmt.Errorf("repo.GetGradesByClassID: scanning grade: %w", err)
		}
		grades = append(grades, g)
	}
	if err = gradeRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("repo.GetGradesByClassID: grade rows error: %w", err)
	}

	// 3. Get all students for the class
	studentsQuery := `SELECT id, class_id, enrollment_id, full_name, status, created_at, updated_at FROM students WHERE class_id = ?`
	studentRows, err := r.db.QueryContext(ctx, studentsQuery, classID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("repo.GetGradesByClassID: fetching students: %w", err)
	}
	defer studentRows.Close()

	var students []models.Student
	for studentRows.Next() {
		var s models.Student
		var enrollmentID sql.NullString
		if err := studentRows.Scan(&s.ID, &s.ClassID, &enrollmentID, &s.FullName, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, nil, nil, fmt.Errorf("repo.GetGradesByClassID: scanning student: %w", err)
		}
		s.EnrollmentID = enrollmentID.String
		students = append(students, s)
	}
	if err = studentRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("repo.GetGradesByClassID: student rows error: %w", err)
	}

	return grades, assessments, students, nil
}

func (r *assessmentRepository) ListAllAssessments(ctx context.Context) ([]models.Assessment, error) {
	query := `SELECT id, class_id, name, term, weight, assessment_date FROM assessments ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("assessmentRepository.ListAllAssessments: query failed: %w", err)
	}
	defer rows.Close()

	var assessments []models.Assessment
	for rows.Next() {
		var assessment models.Assessment
		err := rows.Scan(
			&assessment.ID,
			&assessment.ClassID,
			&assessment.Name,
			&assessment.Term,
			&assessment.Weight,
			&assessment.AssessmentDate,
		)
		if err != nil {
			return nil, fmt.Errorf("assessmentRepository.ListAllAssessments: scan failed: %w", err)
		}
		assessments = append(assessments, assessment)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("assessmentRepository.ListAllAssessments: rows error: %w", err)
	}
	return assessments, nil
}

// sqlRepeat is a helper to generate placeholders for IN clauses.
func sqlRepeat(count int) string {
	if count <= 0 {
		return ""
	}
	return ", ?"
}
