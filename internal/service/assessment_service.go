package service

import (
	"context"
	"vigenda/internal/models"
)

// AssessmentService define os métodos para a gestão de avaliações e notas.
type AssessmentService interface {
    CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error)
    EnterGrades(ctx context.Context, assessmentID int64, studentGrades map[int64]float64) error
    CalculateClassAverage(ctx context.Context, classID int64) (float64, error)
}
