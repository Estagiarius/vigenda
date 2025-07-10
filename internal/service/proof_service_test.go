package service_test

import (
	"context"
	"errors"
	"testing"
	"vigenda/internal/models"
	"vigenda/internal/repository"
	"vigenda/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuestionRepository é uma implementação mock da interface QuestionRepository.
type MockQuestionRepository struct {
	mock.Mock
}

func (m *MockQuestionRepository) GetQuestionsByCriteria(ctx context.Context, criteria repository.QuestionQueryCriteria) ([]models.Question, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockQuestionRepository) AddQuestion(ctx context.Context, question *models.Question) (int64, error) {
	args := m.Called(ctx, question)
	return args.Get(0).(int64), args.Error(1)
}

// GetQuestionsByCriteriaProofGeneration is the mock implementation
func (m *MockQuestionRepository) GetQuestionsByCriteriaProofGeneration(ctx context.Context, criteria repository.ProofCriteria) ([]models.Question, error) {
	// This mock implementation needs to be smart based on criteria, or tests need to set up different On().Return() calls.
	// For now, let's assume tests will use .On to specify behavior for different criteria.
	// The `GetQuestionsByCriteria` is already used by tests for this, which is fine.
	// This method is on the interface, so the mock needs it.
	// It can delegate to GetQuestionsByCriteria or have its own mock logic if the service calls it distinctly.
	// The service `GenerateProof` actually calls `GetQuestionsByCriteria` multiple times with different `QuestionQueryCriteria`.
	// So, this specific mock method might not be directly "called" by the tests if the tests mock `GetQuestionsByCriteria`.
	// However, to satisfy the interface, it must exist.
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Question), args.Error(1)
}

func TestProofService_GenerateProof(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockQuestionRepository)
		proofService := service.NewProofService(mockRepo)

		topic := "Algebra"
		criteria := service.ProofCriteria{
			SubjectID:   1,
			Topic:       &topic,
			EasyCount:   1,
			MediumCount: 1,
			HardCount:   1,
		}

		easyQ := []models.Question{{ID: 1, Difficulty: "facil", SubjectID: 1, Topic: topic}}
		mediumQ := []models.Question{{ID: 2, Difficulty: "media", SubjectID: 1, Topic: topic}}
		hardQ := []models.Question{{ID: 3, Difficulty: "dificil", SubjectID: 1, Topic: topic}}

		allExpectedQuestions := append(easyQ, mediumQ...)
		allExpectedQuestions = append(allExpectedQuestions, hardQ...)

		// Expect a single call to GetQuestionsByCriteriaProofGeneration
		mockRepo.On("GetQuestionsByCriteriaProofGeneration", ctx, repository.ProofCriteria{
			SubjectID:   1,
			Topic:       &topic,
			EasyCount:   1,
			MediumCount: 1,
			HardCount:   1,
		}).Return(allExpectedQuestions, nil).Once()

		questions, err := proofService.GenerateProof(ctx, criteria)

		assert.NoError(t, err)
		assert.Len(t, questions, 3)
		// Individual contains assertions are still good if the order is not guaranteed by the mock return,
		// but if allExpectedQuestions are returned in a specific order, direct comparison might be possible.
		assert.Contains(t, questions, easyQ[0])
		assert.Contains(t, questions, mediumQ[0])
		assert.Contains(t, questions, hardQ[0])
		mockRepo.AssertExpectations(t)
	})

	t.Run("error_no_difficulty_count", func(t *testing.T) {
		mockRepo := new(MockQuestionRepository)
		proofService := service.NewProofService(mockRepo)
		criteria := service.ProofCriteria{SubjectID: 1} // No counts

		_, err := proofService.GenerateProof(ctx, criteria)
		assert.Error(t, err)
		assert.EqualError(t, err, "pelo menos uma contagem de dificuldade deve ser maior que zero")
	})

	t.Run("error_fetching_easy_questions", func(t *testing.T) {
		mockRepo := new(MockQuestionRepository)
		proofService := service.NewProofService(mockRepo)
		topic := "Geometry"
		criteria := service.ProofCriteria{SubjectID: 1, Topic: &topic, EasyCount: 1}

		// Simulate an error when the repository's GetQuestionsByCriteriaProofGeneration is called
		mockRepo.On("GetQuestionsByCriteriaProofGeneration", ctx, repository.ProofCriteria{
			SubjectID:   1,
			Topic:       &topic,
			EasyCount:   1,
			MediumCount: 0, // Assuming 0 for other counts if not specified in this test case
			HardCount:   0,
		}).Return(nil, errors.New("db error general fetch")).Once()

		_, err := proofService.GenerateProof(ctx, criteria)
		assert.Error(t, err)
		// The service wraps the error from GetQuestionsByCriteriaProofGeneration
		assert.Contains(t, err.Error(), "erro ao buscar questões para a prova: db error general fetch")
		mockRepo.AssertExpectations(t)
	})

	t.Run("error_not_enough_easy_questions", func(t *testing.T) {
		mockRepo := new(MockQuestionRepository)
		proofService := service.NewProofService(mockRepo)
		topic := "History"
		criteria := service.ProofCriteria{SubjectID: 1, Topic: &topic, EasyCount: 2}

		// Simulate returning only one easy question when two were requested.
		// Other difficulties (medium, hard) are not requested in this specific criteria.
		easyQ := []models.Question{{ID: 1, Difficulty: "facil", SubjectID: 1, Topic: topic}}
		mockRepo.On("GetQuestionsByCriteriaProofGeneration", ctx, repository.ProofCriteria{
			SubjectID:   1,
			Topic:       &topic,
			EasyCount:   2,
			MediumCount: 0,
			HardCount:   0,
		}).Return(easyQ, nil).Once() // Return only one easy question

		_, err := proofService.GenerateProof(ctx, criteria)
		assert.Error(t, err)
		// The service logic should detect this discrepancy.
		assert.EqualError(t, err, "não há questões fáceis suficientes (solicitado: 2, disponível: 1)")
		mockRepo.AssertExpectations(t)
	})

	t.Run("success_only_medium_questions", func(t *testing.T) {
		mockRepo := new(MockQuestionRepository)
		proofService := service.NewProofService(mockRepo)
		criteria := service.ProofCriteria{
			SubjectID:   2,
			MediumCount: 2,
		}

		mediumQ := []models.Question{
			{ID: 10, Difficulty: "media", SubjectID: 2},
			{ID: 11, Difficulty: "media", SubjectID: 2},
		}

		mockRepo.On("GetQuestionsByCriteriaProofGeneration", ctx, repository.ProofCriteria{
			SubjectID:   2,
			Topic:       nil,
			EasyCount:   0,
			MediumCount: 2,
			HardCount:   0,
		}).Return(mediumQ, nil).Once()

		questions, err := proofService.GenerateProof(ctx, criteria)

		assert.NoError(t, err)
		assert.Len(t, questions, 2)
		assert.Contains(t, questions, mediumQ[0])
		assert.Contains(t, questions, mediumQ[1])
		mockRepo.AssertExpectations(t)
	})

	t.Run("error_fetching_medium_questions", func(t *testing.T) {
		mockRepo := new(MockQuestionRepository)
		proofService := service.NewProofService(mockRepo)
		criteria := service.ProofCriteria{SubjectID: 1, MediumCount: 1}

		mockRepo.On("GetQuestionsByCriteriaProofGeneration", ctx, repository.ProofCriteria{
			SubjectID:   1,
			Topic:       nil,
			EasyCount:   0,
			MediumCount: 1,
			HardCount:   0,
		}).Return(nil, errors.New("db error medium fetch")).Once()

		_, err := proofService.GenerateProof(ctx, criteria)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao buscar questões para a prova: db error medium fetch")
		mockRepo.AssertExpectations(t)
	})

	t.Run("error_not_enough_medium_questions", func(t *testing.T) {
		mockRepo := new(MockQuestionRepository)
		proofService := service.NewProofService(mockRepo)
		criteria := service.ProofCriteria{SubjectID: 1, MediumCount: 2}

		mediumQ := []models.Question{{ID: 1, Difficulty: "media", SubjectID: 1}}
		mockRepo.On("GetQuestionsByCriteriaProofGeneration", ctx, repository.ProofCriteria{
			SubjectID:   1,
			Topic:       nil,
			EasyCount:   0,
			MediumCount: 2,
			HardCount:   0,
		}).Return(mediumQ, nil).Once()

		_, err := proofService.GenerateProof(ctx, criteria)
		assert.Error(t, err)
		assert.EqualError(t, err, "não há questões médias suficientes (solicitado: 2, disponível: 1)")
		mockRepo.AssertExpectations(t)
	})

	t.Run("error_fetching_hard_questions", func(t *testing.T) {
		mockRepo := new(MockQuestionRepository)
		proofService := service.NewProofService(mockRepo)
		criteria := service.ProofCriteria{SubjectID: 1, HardCount: 1}

		mockRepo.On("GetQuestionsByCriteriaProofGeneration", ctx, repository.ProofCriteria{
			SubjectID:   1,
			Topic:       nil,
			EasyCount:   0,
			MediumCount: 0,
			HardCount:   1,
		}).Return(nil, errors.New("db error hard fetch")).Once()

		_, err := proofService.GenerateProof(ctx, criteria)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao buscar questões para a prova: db error hard fetch")
		mockRepo.AssertExpectations(t)
	})

	t.Run("error_not_enough_hard_questions", func(t *testing.T) {
		mockRepo := new(MockQuestionRepository)
		proofService := service.NewProofService(mockRepo)
		criteria := service.ProofCriteria{SubjectID: 1, HardCount: 2}

		hardQ := []models.Question{{ID: 1, Difficulty: "dificil", SubjectID: 1}}
		mockRepo.On("GetQuestionsByCriteriaProofGeneration", ctx, repository.ProofCriteria{
			SubjectID:   1,
			Topic:       nil,
			EasyCount:   0,
			MediumCount: 0,
			HardCount:   2,
		}).Return(hardQ, nil).Once()

		_, err := proofService.GenerateProof(ctx, criteria)
		assert.Error(t, err)
		assert.EqualError(t, err, "não há questões difíceis suficientes (solicitado: 2, disponível: 1)")
		mockRepo.AssertExpectations(t)
	})

}
