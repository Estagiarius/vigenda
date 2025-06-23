package service

import (
	"context"
	"vigenda/internal/models"
)

type assessmentServiceImpl struct {
	// Add repository dependencies here
}

// NewAssessmentService creates a new instance of AssessmentService.
func NewAssessmentService() AssessmentService {
	return &assessmentServiceImpl{}
}

func (s *assessmentServiceImpl) CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error) {
	// TODO: Implement actual logic
	return models.Assessment{}, nil
}

func (s *assessmentServiceImpl) EnterGrades(ctx context.Context, assessmentID int64, studentGrades map[int64]float64) error {
	// TODO: Implement actual logic
	return nil
}

func (s *assessmentServiceImpl) CalculateClassAverage(ctx context.Context, classID int64) (float64, error) {
	// TODO: Implement actual logic
	return 0.0, nil
}
