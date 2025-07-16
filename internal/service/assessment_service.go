package service

import (
	"context"
	"fmt"
	"vigenda/internal/models"
	"vigenda/internal/repository" // Added import
)

type assessmentServiceImpl struct {
	assessmentRepo repository.AssessmentRepository
	classRepo      repository.ClassRepository // Added classRepo for fetching students
}

// NewAssessmentService creates a new instance of AssessmentService.
// It now accepts AssessmentRepository and ClassRepository as dependencies.
func NewAssessmentService(
	assessmentRepo repository.AssessmentRepository,
	classRepo repository.ClassRepository,
) AssessmentService {
	return &assessmentServiceImpl{
		assessmentRepo: assessmentRepo,
		classRepo:      classRepo,
	}
}

func (s *assessmentServiceImpl) CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error) {
	if name == "" {
		return models.Assessment{}, fmt.Errorf("assessment name cannot be empty")
	}
	if classID == 0 {
		return models.Assessment{}, fmt.Errorf("class ID cannot be zero")
	}
	// Bypassing validation for the special "Nota Final" assessment
	if name != FinalGradeAssessmentName {
		if term <= 0 {
			return models.Assessment{}, fmt.Errorf("term must be positive")
		}
		if weight <= 0 {
			return models.Assessment{}, fmt.Errorf("weight must be positive")
		}
	}
	// TODO: Validate classID exists using s.classRepo.GetClassByID(ctx, classID)

	// Assuming UserID 1 for now
	// userID := int64(1) // UserID is not part of models.Assessment

	assessment := models.Assessment{
		// UserID:  userID, // Removed
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
	if assessmentID == 0 {
		return fmt.Errorf("assessment ID cannot be zero")
	}
	if len(studentGrades) == 0 {
		return fmt.Errorf("no grades provided")
	}

	// Optional: Validate assessmentID exists
	assessment, err := s.assessmentRepo.GetAssessmentByID(ctx, assessmentID)
	if err != nil {
		return fmt.Errorf("service.EnterGrades: validating assessment: %w", err)
	}
	if assessment == nil {
		return fmt.Errorf("assessment with ID %d not found", assessmentID)
	}

	// Optional: Validate studentIDs exist within the assessment's class
	// This would involve fetching students for assessment.ClassID and checking existence.
	// For now, we assume student IDs are valid and belong to the correct class.

	// Assuming UserID 1 for now
	// userID := int64(1) // UserID is not part of models.Grade

	for studentID, gradeVal := range studentGrades {
		if studentID == 0 {
			return fmt.Errorf("student ID cannot be zero in grades map")
		}
		// Basic grade validation (e.g., 0-10, or whatever scale)
		if gradeVal < 0 || gradeVal > 100 { // Assuming a 0-100 scale for placeholder
			// return fmt.Errorf("invalid grade value %.2f for student %d. Must be between 0 and 100", gradeVal, studentID)
			// For now, let's allow any float. Specific validation can be added.
		}

		grade := models.Grade{
			AssessmentID: assessmentID,
			StudentID:    studentID,
			// UserID:       userID, // Removed: User who entered the grade is not stored in Grade model
			Grade:        gradeVal,
		}
		if err := s.assessmentRepo.EnterGrade(ctx, &grade); err != nil {
			return fmt.Errorf("service.EnterGrades: entering grade for student %d: %w", studentID, err)
		}
	}
	return nil
}

func (s *assessmentServiceImpl) CalculateClassAverage(ctx context.Context, classID int64, terms []int) (map[int64]float64, error) {
	if classID == 0 {
		return nil, fmt.Errorf("class ID cannot be zero")
	}

	grades, allAssessments, students, err := s.assessmentRepo.GetGradesByClassID(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("service.CalculateClassAverage: fetching data: %w", err)
	}

	if len(students) == 0 {
		return nil, fmt.Errorf("no students found in class %d to calculate average", classID)
	}

	// Create a set for quick lookup of terms to include
	includeTerms := make(map[int]bool)
	filterByTerm := len(terms) > 0
	if filterByTerm {
		for _, t := range terms {
			includeTerms[t] = true
		}
	}

	// Filter assessments based on terms
	var assessments []models.Assessment
	for _, a := range allAssessments {
		if !filterByTerm || includeTerms[a.Term] {
			assessments = append(assessments, a)
		}
	}

	if len(assessments) == 0 {
		return nil, fmt.Errorf("no assessments found for class %d matching the specified terms", classID)
	}

	assessmentMap := make(map[int64]models.Assessment)
	for _, a := range assessments {
		assessmentMap[a.ID] = a
	}

	studentTotals := make(map[int64]struct {
		totalWeightedGrade float64
		totalWeight        float64
	})

	for _, g := range grades {
		assessment, ok := assessmentMap[g.AssessmentID]
		if !ok {
			// This grade is for an assessment filtered out by term, so skip it
			continue
		}

		st := studentTotals[g.StudentID]
		st.totalWeightedGrade += g.Grade * assessment.Weight
		st.totalWeight += assessment.Weight
		studentTotals[g.StudentID] = st
	}

	studentAverages := make(map[int64]float64)
	for _, student := range students {
		if student.Status != "ativo" {
			continue // Only calculate for active students
		}
		totals, ok := studentTotals[student.ID]
		if ok && totals.totalWeight > 0 {
			studentAverages[student.ID] = totals.totalWeightedGrade / totals.totalWeight
		} else {
			studentAverages[student.ID] = 0 // Student has no grades for the selected terms
		}
	}

	return studentAverages, nil
}

// Helper function to find a student in a slice (if needed, not strictly necessary with map lookups)
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

func (s *assessmentServiceImpl) DeleteAssessment(ctx context.Context, assessmentID int64) error {
	if assessmentID == 0 {
		return fmt.Errorf("assessment ID cannot be zero")
	}
	return s.assessmentRepo.DeleteAssessment(ctx, assessmentID)
}

func (s *assessmentServiceImpl) GetStudentsForGrading(ctx context.Context, assessmentID int64) ([]models.Student, *models.Assessment, error) {
	if assessmentID == 0 {
		return nil, nil, fmt.Errorf("assessment ID cannot be zero")
	}

	// 1. Get the assessment details
	assessment, err := s.assessmentRepo.GetAssessmentByID(ctx, assessmentID)
	if err != nil {
		return nil, nil, fmt.Errorf("service.GetStudentsForGrading: failed to get assessment: %w", err)
	}
	if assessment == nil {
		return nil, nil, fmt.Errorf("assessment with ID %d not found", assessmentID)
	}

	// 2. Get students from the assessment's class
	students, err := s.classRepo.GetStudentsByClassID(ctx, assessment.ClassID)
	if err != nil {
		return nil, nil, fmt.Errorf("service.GetStudentsForGrading: failed to get students for class ID %d: %w", assessment.ClassID, err)
	}

	return students, assessment, nil
}

const FinalGradeAssessmentName = "Nota Final"

// getOrCreateFinalAssessment finds or creates a special assessment for final grades.
func (s *assessmentServiceImpl) getOrCreateFinalAssessment(ctx context.Context, classID int64) (int64, error) {
	// 1. Try to find the existing final assessment
	assessment, err := s.assessmentRepo.FindAssessmentByNameAndClass(ctx, FinalGradeAssessmentName, classID)
	if err != nil {
		return 0, fmt.Errorf("failed to find final assessment: %w", err)
	}

	// 2. If it exists, return its ID
	if assessment != nil {
		return assessment.ID, nil
	}

	// 3. If not, create it
	// Using term 0 and weight 0 to signify it's a special, non-weighted assessment
	newAssessment, err := s.CreateAssessment(ctx, FinalGradeAssessmentName, classID, 0, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to create final assessment: %w", err)
	}
	return newAssessment.ID, nil
}

func (s *assessmentServiceImpl) EnterFinalGrades(ctx context.Context, classID int64, finalGrades map[int64]float64) error {
	if classID == 0 {
		return fmt.Errorf("class ID cannot be zero")
	}

	// 1. Get the special assessment ID for final grades
	finalAssessmentID, err := s.getOrCreateFinalAssessment(ctx, classID)
	if err != nil {
		return err
	}

	// 2. Use the existing EnterGrades logic with the special assessment ID
	return s.EnterGrades(ctx, finalAssessmentID, finalGrades)
}
