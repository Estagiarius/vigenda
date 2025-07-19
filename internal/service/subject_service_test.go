package service

import (
	"context"
	"fmt"
	"testing"
	"vigenda/internal/models"
	"vigenda/internal/repository/stubs"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSubjectServiceImpl_CreateSubject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := stubs.NewMockSubjectRepository(ctrl)
	subjectService := NewSubjectService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, subject *models.Subject) error {
			subject.ID = 1
			return nil
		}).Times(1)

		createdSubject, err := subjectService.CreateSubject(ctx, 1, "Matemática")

		assert.NoError(t, err)
		assert.Equal(t, int64(1), createdSubject.ID)
		assert.Equal(t, "Matemática", createdSubject.Name)
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := subjectService.CreateSubject(ctx, 1, "")
		assert.Error(t, err)
		assert.Equal(t, "O nome da disciplina não pode ser vazio", err.Error())
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(fmt.Errorf("db error")).Times(1)

		_, err := subjectService.CreateSubject(ctx, 1, "História")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Falha ao criar disciplina no repositório")
	})
}

func TestSubjectServiceImpl_ListSubjectsByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := stubs.NewMockSubjectRepository(ctrl)
	subjectService := NewSubjectService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedSubjects := []models.Subject{
			{ID: 1, UserID: 1, Name: "Matemática"},
			{ID: 2, UserID: 1, Name: "Português"},
		}
		mockRepo.EXPECT().GetByUserID(ctx, int64(1)).Return(expectedSubjects, nil).Times(1)

		subjects, err := subjectService.ListSubjectsByUser(ctx, 1)

		assert.NoError(t, err)
		assert.Len(t, subjects, 2)
		assert.Equal(t, expectedSubjects, subjects)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().GetByUserID(ctx, int64(1)).Return(nil, fmt.Errorf("db error")).Times(1)

		_, err := subjectService.ListSubjectsByUser(ctx, 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Falha ao listar disciplinas")
	})
}
