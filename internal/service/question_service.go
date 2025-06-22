package service

import (
	"context"
	"vigenda/internal/models"
)

// QuestionService define os métodos para o banco de questões e geração de provas.
type QuestionService interface {
	AddQuestionsFromJSON(ctx context.Context, jsonData []byte) (int, error)
	GenerateTest(ctx context.Context, criteria TestCriteria) ([]models.Question, error)
}

// TestCriteria define os parâmetros para a geração de uma prova.
type TestCriteria struct {
	SubjectID   int64
	Topic       string
	EasyCount   int
	MediumCount int
	HardCount   int
}
