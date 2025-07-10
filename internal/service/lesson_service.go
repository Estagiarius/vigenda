package service

import (
	"context"
	"fmt"
	"time"
	"vigenda/internal/models"
	"vigenda/internal/repository"
)

type lessonServiceImpl struct {
	lessonRepo repository.LessonRepository
	classRepo  repository.ClassRepository // Para verificar se a turma pertence ao usuário
}

func NewLessonService(lessonRepo repository.LessonRepository, classRepo repository.ClassRepository) LessonService {
	return &lessonServiceImpl{lessonRepo: lessonRepo, classRepo: classRepo}
}

// validateUserOwnsClass verifica se uma turma pertence a um usuário específico.
// Retorna um erro se a turma não for encontrada ou não pertencer ao usuário.
// UserID 0 é tratado como um superusuário/sistema que tem acesso a todas as turmas.
func (s *lessonServiceImpl) validateUserOwnsClass(ctx context.Context, userID int64, classID int64) (*models.Class, error) {
	if userID == 0 { // UserID 0 pode ser um "usuário sistema" ou admin, pula a checagem de propriedade.
		class, err := s.classRepo.GetClassByID(ctx, classID)
		if err != nil {
			return nil, fmt.Errorf("turma com ID %d não encontrada: %w", classID, err)
		}
		return class, nil
	}

	class, err := s.classRepo.GetClassByID(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("turma com ID %d não encontrada: %w", classID, err)
	}
	if class.UserID != userID {
		return nil, fmt.Errorf("acesso negado: turma %d não pertence ao usuário %d", classID, userID)
	}
	return class, nil
}


func (s *lessonServiceImpl) CreateLesson(ctx context.Context, classID int64, title string, planContent string, scheduledAt time.Time) (models.Lesson, error) {
	if title == "" {
		return models.Lesson{}, fmt.Errorf("título da lição não pode ser vazio")
	}
	// TODO: Obter UserID do contexto ctx quando a autenticação estiver implementada.
	// Por enquanto, vamos assumir um UserID fixo ou que a validação de propriedade será feita em outro lugar se necessário.
	// Para o dashboard, podemos usar UserID 1.
	userID := int64(1) // Placeholder

	if _, err := s.validateUserOwnsClass(ctx, userID, classID); err != nil {
		return models.Lesson{}, fmt.Errorf("CreateLesson: %w", err)
	}

	lesson := models.Lesson{
		ClassID:     classID,
		Title:       title,
		PlanContent: planContent,
		ScheduledAt: scheduledAt,
	}
	id, err := s.lessonRepo.CreateLesson(ctx, &lesson)
	if err != nil {
		return models.Lesson{}, fmt.Errorf("lessonService.CreateLesson: %w", err)
	}
	lesson.ID = id
	// Para obter CreatedAt/UpdatedAt, precisaríamos de um GetLessonByID ou o repo populá-los.
	// Vamos buscar para garantir.
	createdLesson, err := s.lessonRepo.GetLessonByID(ctx, id)
	if err != nil {
		return models.Lesson{}, fmt.Errorf("lessonService.CreateLesson: buscar lição criada: %w", err)
	}
	return *createdLesson, nil
}

func (s *lessonServiceImpl) GetLessonByID(ctx context.Context, lessonID int64) (models.Lesson, error) {
	lesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return models.Lesson{}, fmt.Errorf("lessonService.GetLessonByID: %w", err)
	}
	// TODO: Validar propriedade da lição via turma (UserID)
	// userID := int64(1) // Placeholder
	// if _, err := s.validateUserOwnsClass(ctx, userID, lesson.ClassID); err != nil {
	// 	 return models.Lesson{}, fmt.Errorf("GetLessonByID: %w", err)
	// }
	return *lesson, nil
}

func (s *lessonServiceImpl) GetLessonsByClassID(ctx context.Context, classID int64) ([]models.Lesson, error) {
	// TODO: Validar propriedade da turma (UserID)
	// userID := int64(1) // Placeholder
	// if _, err := s.validateUserOwnsClass(ctx, userID, classID); err != nil {
	// 	 return nil, fmt.Errorf("GetLessonsByClassID: %w", err)
	// }
	lessons, err := s.lessonRepo.GetLessonsByClassID(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("lessonService.GetLessonsByClassID: %w", err)
	}
	return lessons, nil
}

func (s *lessonServiceImpl) GetLessonsForDate(ctx context.Context, userID int64, date time.Time) ([]models.Lesson, error) {
	// Define o início e o fim do dia para a data fornecida.
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	// O userID é passado para o repositório, que pode usá-lo para filtrar classes pertencentes ao usuário.
	// Se userID for 0, o repositório pode interpretar como "buscar para todos os usuários".
	lessons, err := s.lessonRepo.GetLessonsByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("lessonService.GetLessonsForDate: %w", err)
	}
	return lessons, nil
}

func (s *lessonServiceImpl) UpdateLesson(ctx context.Context, lessonID int64, title string, planContent string, scheduledAt time.Time) (models.Lesson, error) {
	if title == "" {
		return models.Lesson{}, fmt.Errorf("título da lição não pode ser vazio")
	}

	existingLesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return models.Lesson{}, fmt.Errorf("lessonService.UpdateLesson: lição não encontrada: %w", err)
	}

	// TODO: Validar propriedade da lição via turma (UserID)
	userID := int64(1) // Placeholder
	if _, err := s.validateUserOwnsClass(ctx, userID, existingLesson.ClassID); err != nil {
		return models.Lesson{}, fmt.Errorf("UpdateLesson: %w", err)
	}

	existingLesson.Title = title
	existingLesson.PlanContent = planContent
	existingLesson.ScheduledAt = scheduledAt
	// ClassID não deve mudar aqui. Se precisar mudar a turma, seria outra operação.

	if err := s.lessonRepo.UpdateLesson(ctx, existingLesson); err != nil {
		return models.Lesson{}, fmt.Errorf("lessonService.UpdateLesson: %w", err)
	}
	// Retorna a lição atualizada buscando-a novamente para garantir consistência
	updatedLesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return models.Lesson{}, fmt.Errorf("lessonService.UpdateLesson: buscar lição atualizada: %w", err)
	}
	return *updatedLesson, nil
}

func (s *lessonServiceImpl) DeleteLesson(ctx context.Context, lessonID int64) error {
	existingLesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return fmt.Errorf("lessonService.DeleteLesson: lição não encontrada: %w", err)
	}
	// TODO: Validar propriedade da lição via turma (UserID)
	userID := int64(1) // Placeholder
	if _, err := s.validateUserOwnsClass(ctx, userID, existingLesson.ClassID); err != nil {
		return fmt.Errorf("DeleteLesson: %w", err)
	}

	if err := s.lessonRepo.DeleteLesson(ctx, lessonID); err != nil {
		return fmt.Errorf("lessonService.DeleteLesson: %w", err)
	}
	return nil
}
