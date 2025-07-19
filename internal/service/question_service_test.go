package service_test

import (
	"context"
	"testing"
	"vigenda/internal/models"
	"vigenda/internal/repository"
	"vigenda/internal/repository/stubs"
	"vigenda/internal/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestQuestionService_AddQuestionsFromJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockQuestionRepo := stubs.NewMockQuestionRepository(ctrl)
		mockSubjectRepo := stubs.NewMockSubjectRepository(ctrl)
		questionService := service.NewQuestionService(mockQuestionRepo, mockSubjectRepo)

		jsonData := []byte(`[
			{
				"disciplina": "Matemática",
				"topico": "Álgebra",
				"tipo": "dissertativa",
				"dificuldade": "media",
				"enunciado": "Resolva a equação x + 5 = 10.",
				"resposta_correta": "x = 5",
				"user_id": 1
			},
			{
				"disciplina": "História",
				"topico": "Revolução Francesa",
				"tipo": "multipla_escolha",
				"dificuldade": "facil",
				"enunciado": "Qual evento marcou o início da Revolução Francesa?",
				"opcoes": ["Queda da Bastilha", "Tomada do Poder por Napoleão"],
				"resposta_correta": "Queda da Bastilha",
				"user_id": 1
			}
		]`)

		mockSubjectRepo.EXPECT().GetOrCreateByNameAndUser(ctx, "Matemática", int64(1)).Return(models.Subject{ID: 1, Name: "Matemática", UserID: 1}, nil).Times(1)
		mockSubjectRepo.EXPECT().GetOrCreateByNameAndUser(ctx, "História", int64(1)).Return(models.Subject{ID: 2, Name: "História", UserID: 1}, nil).Times(1)

		mockQuestionRepo.EXPECT().AddQuestion(ctx, gomock.Any()).Return(int64(1), nil).Times(2)

		count, err := questionService.AddQuestionsFromJSON(ctx, jsonData)

		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("error_invalid_json", func(t *testing.T) {
		mockQuestionRepo := stubs.NewMockQuestionRepository(ctrl)
		mockSubjectRepo := stubs.NewMockSubjectRepository(ctrl)
		questionService := service.NewQuestionService(mockQuestionRepo, mockSubjectRepo)
		jsonData := []byte(`[{"disciplina": "Matemática", "topico": "Álgebra"`)

		_, err := questionService.AddQuestionsFromJSON(ctx, jsonData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "falha ao decodificar JSON")
	})
}

func TestQuestionService_GenerateTest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := stubs.NewMockQuestionRepository(ctrl)
		questionService := service.NewQuestionService(mockRepo, nil)

		topic := "Calculus"
		criteria := service.TestCriteria{
			SubjectID:   1,
			Topic:       &topic,
			EasyCount:   1,
			MediumCount: 1,
			HardCount:   1,
		}

		easyQ := []models.Question{{ID: 1, Difficulty: "facil", SubjectID: 1, Topic: topic}}
		mediumQ := []models.Question{{ID: 2, Difficulty: "media", SubjectID: 1, Topic: topic}}
		hardQ := []models.Question{{ID: 3, Difficulty: "dificil", SubjectID: 1, Topic: topic}}

		mockRepo.EXPECT().GetQuestionsByCriteria(ctx, repository.QuestionQueryCriteria{SubjectID: 1, Topic: &topic, Difficulty: "facil", Limit: 1}).Return(easyQ, nil).Times(1)
		mockRepo.EXPECT().GetQuestionsByCriteria(ctx, repository.QuestionQueryCriteria{SubjectID: 1, Topic: &topic, Difficulty: "media", Limit: 1}).Return(mediumQ, nil).Times(1)
		mockRepo.EXPECT().GetQuestionsByCriteria(ctx, repository.QuestionQueryCriteria{SubjectID: 1, Topic: &topic, Difficulty: "dificil", Limit: 1}).Return(hardQ, nil).Times(1)

		questions, err := questionService.GenerateTest(ctx, criteria)

		assert.NoError(t, err)
		assert.Len(t, questions, 3)
		assert.Contains(t, questions, easyQ[0])
		assert.Contains(t, questions, mediumQ[0])
		assert.Contains(t, questions, hardQ[0])
	})
}
