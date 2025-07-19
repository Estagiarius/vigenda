package repository

import (
	"context"
	"vigenda/internal/models"
)

//go:generate mockgen -source=question_repo.go -destination=stubs/question_repository_mock.go -package=stubs QuestionRepository
type QuestionRepository interface {
	GetQuestionsByCriteria(ctx context.Context, criteria QuestionQueryCriteria) ([]models.Question, error)
	AddQuestion(ctx context.Context, question *models.Question) (int64, error)
	GetQuestionsByCriteriaProofGeneration(ctx context.Context, criteria ProofCriteria) ([]models.Question, error)
}
