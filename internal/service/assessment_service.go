package service

import (
	"context"
	"fmt"
	"vigenda/internal/models"
	"vigenda/internal/repository"
)

// GradingSheet é uma struct que agrupa todas as informações necessárias para a tela de lançamento de notas.
type GradingSheet struct {
	Assessment models.Assessment
	Students   []models.Student
	Grades     map[int64]models.Grade // Mapeia StudentID para a nota existente.
}

// AssessmentService define a interface para a lógica de negócios relacionada a avaliações e notas.
// Esta é uma versão atualizada da interface para incluir o método GetGradingSheet.
type AssessmentService interface {
	CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error)
	EnterGrades(ctx context.Context, assessmentID int64, studentGrades map[int64]float64) error
	CalculateClassAverage(ctx context.Context, classID int64) (float64, error)
	ListAllAssessments(ctx context.Context) ([]models.Assessment, error)
	GetGradingSheet(ctx context.Context, assessmentID int64) (*GradingSheet, error)
	GetAssessmentByID(ctx context.Context, assessmentID int64) (*models.Assessment, error)
	DeleteAssessment(ctx context.Context, assessmentID int64) error
	UpdateAssessment(ctx context.Context, assessmentID int64, name string, classID int64, term int, weight float64) (models.Assessment, error)
}


type assessmentServiceImpl struct {
	assessmentRepo repository.AssessmentRepository
	classRepo      repository.ClassRepository
}

// NewAssessmentService cria uma nova instância de AssessmentService.
func NewAssessmentService(
	assessmentRepo repository.AssessmentRepository,
	classRepo repository.ClassRepository,
) AssessmentService {
	return &assessmentServiceImpl{
		assessmentRepo: assessmentRepo,
		classRepo:      classRepo,
	}
}

func (s *assessmentServiceImpl) GetAssessmentByID(ctx context.Context, assessmentID int64) (*models.Assessment, error) {
	if assessmentID <= 0 {
		return nil, fmt.Errorf("assessment ID must be positive")
	}
	return s.assessmentRepo.GetAssessmentByID(ctx, assessmentID)
}

func (s *assessmentServiceImpl) DeleteAssessment(ctx context.Context, assessmentID int64) error {
	if assessmentID <= 0 {
		return fmt.Errorf("assessment ID must be positive")
	}
	// Optional: check ownership before deleting
	return s.assessmentRepo.DeleteAssessment(ctx, assessmentID)
}

func (s *assessmentServiceImpl) UpdateAssessment(ctx context.Context, assessmentID int64, name string, classID int64, term int, weight float64) (models.Assessment, error) {
	if assessmentID <= 0 {
		return models.Assessment{}, fmt.Errorf("assessment ID must be positive")
	}
	if name == "" {
		return models.Assessment{}, fmt.Errorf("assessment name cannot be empty")
	}
	// ... other validations ...

	assessmentToUpdate, err := s.assessmentRepo.GetAssessmentByID(ctx, assessmentID)
	if err != nil {
		return models.Assessment{}, fmt.Errorf("failed to get assessment to update: %w", err)
	}

	assessmentToUpdate.Name = name
	assessmentToUpdate.ClassID = classID
	assessmentToUpdate.Term = term
	assessmentToUpdate.Weight = weight

	err = s.assessmentRepo.UpdateAssessment(ctx, assessmentToUpdate)
	if err != nil {
		return models.Assessment{}, fmt.Errorf("failed to update assessment: %w", err)
	}
	return *assessmentToUpdate, nil
}


func (s *assessmentServiceImpl) CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error) {
	if name == "" {
		return models.Assessment{}, fmt.Errorf("assessment name cannot be empty")
	}
	if classID <= 0 {
		return models.Assessment{}, fmt.Errorf("class ID must be positive")
	}
	if term <= 0 {
		return models.Assessment{}, fmt.Errorf("term must be positive")
	}
	if weight <= 0 {
		return models.Assessment{}, fmt.Errorf("weight must be positive")
	}

	assessment := models.Assessment{
		ClassID: classID,
		Name:    name,
		Term:    term,
		Weight:  weight,
	}

	id, err := s.assessmentRepo.CreateAssessment(ctx, &assessment)
	if err != nil {
		return models.Assessment{}, fmt.Errorf("service.CreateAssessment: %w", err)
	}
	assessment.ID = id
	return assessment, nil
}

func (s *assessmentServiceImpl) EnterGrades(ctx context.Context, assessmentID int64, studentGrades map[int64]float64) error {
	if assessmentID <= 0 {
		return fmt.Errorf("assessment ID cannot be zero")
	}
	if len(studentGrades) == 0 {
		return fmt.Errorf("no grades provided")
	}

	assessment, err := s.assessmentRepo.GetAssessmentByID(ctx, assessmentID)
	if err != nil {
		return fmt.Errorf("service.EnterGrades: validating assessment: %w", err)
	}
	if assessment == nil {
		return fmt.Errorf("assessment with ID %d not found", assessmentID)
	}

	for studentID, gradeVal := range studentGrades {
		if studentID <= 0 {
			return fmt.Errorf("student ID must be positive")
		}
		grade := models.Grade{
			AssessmentID: assessmentID,
			StudentID:    studentID,
			Grade:        gradeVal,
		}
		if err := s.assessmentRepo.EnterGrade(ctx, &grade); err != nil {
			return fmt.Errorf("service.EnterGrades: entering grade for student %d: %w", studentID, err)
		}
	}
	return nil
}

func (s *assessmentServiceImpl) GetGradingSheet(ctx context.Context, assessmentID int64) (*GradingSheet, error) {
	if assessmentID <= 0 {
		return nil, fmt.Errorf("assessment ID must be positive")
	}

	// 1. Get the assessment details
	assessment, err := s.assessmentRepo.GetAssessmentByID(ctx, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assessment details: %w", err)
	}

	// 2. Get all students for the assessment's class
	students, err := s.classRepo.GetStudentsByClassID(ctx, assessment.ClassID)
	if err != nil {
		return nil, fmt.Errorf("failed to get students for class %d: %w", assessment.ClassID, err)
	}

	// 3. Get all existing grades for this assessment
	existingGrades, err := s.assessmentRepo.GetGradesByAssessmentID(ctx, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing grades for assessment %d: %w", assessmentID, err)
	}

	// 4. Build the map of grades for easy lookup
	gradesMap := make(map[int64]models.Grade)
	for _, grade := range existingGrades {
		gradesMap[grade.StudentID] = grade
	}

	return &GradingSheet{
		Assessment: *assessment,
		Students:   students,
		Grades:     gradesMap,
	}, nil
}

func (s *assessmentServiceImpl) CalculateClassAverage(ctx context.Context, classID int64) (float64, error) {
	if classID <= 0 {
		return 0, fmt.Errorf("class ID cannot be zero")
	}

	grades, assessments, students, err := s.assessmentRepo.GetGradesByClassID(ctx, classID)
	if err != nil {
		return 0, fmt.Errorf("service.CalculateClassAverage: fetching data: %w", err)
	}

	if len(students) == 0 {
		return 0, nil // No students, average is 0
	}
	if len(assessments) == 0 {
		return 0, nil // No assessments, average is 0
	}

	studentAverages := make(map[int64]struct {
		totalWeightedGrade float64
		totalWeight        float64
	})

	assessmentMap := make(map[int64]models.Assessment)
	for _, a := range assessments {
		assessmentMap[a.ID] = a
	}

	for _, g := range grades {
		assessment, okA := assessmentMap[g.AssessmentID]
		if !okA {
			continue
		}
		studentInfo, studentExists := findStudent(students, g.StudentID)
		if !studentExists || studentInfo.Status != "ativo" {
			continue
		}

		sa := studentAverages[g.StudentID]
		sa.totalWeightedGrade += g.Grade * assessment.Weight
		sa.totalWeight += assessment.Weight
		studentAverages[g.StudentID] = sa
	}

	var overallClassTotal float64
	var activeStudentsCount int

	for _, student := range students {
		if student.Status != "ativo" {
			continue
		}
		activeStudentsCount++
		sa, ok := studentAverages[student.ID]
		if ok && sa.totalWeight > 0 {
			overallClassTotal += sa.totalWeightedGrade / sa.totalWeight
		}
	}

	if activeStudentsCount == 0 {
		return 0, nil
	}

	return overallClassTotal / float64(activeStudentsCount), nil
}

func findStudent(students []models.Student, studentID int64) (models.Student, bool) {
	for _, s := range students {
		if s.ID == studentID {
			return s, true
		}
	}
	return models.Student{}, false
}

func (s *assessmentServiceImpl) ListAllAssessments(ctx context.Context) ([]models.Assessment, error) {
	assessments, err := s.assessmentRepo.ListAllAssessments(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.ListAllAssessments: %w", err)
	}
	return assessments, nil
}
