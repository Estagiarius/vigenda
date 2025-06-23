package service

import (
	"context"
	"time"
	"vigenda/internal/models"
)

// Implementação do TaskService (exemplo, precisa ser preenchido)
type taskServiceImpl struct {
	// Adicionar dependências de repositório aqui
}

// NewTaskService cria uma nova instância de TaskService.
func NewTaskService(/* dependências do repositório */) TaskService {
	return &taskServiceImpl{
		// inicializar dependências
	}
}

func (s *taskServiceImpl) CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error) {
	// TODO: Implementar lógica
	return models.Task{}, nil
}

func (s *taskServiceImpl) ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error) {
	// TODO: Implementar lógica
	return nil, nil
}

func (s *taskServiceImpl) MarkTaskAsCompleted(ctx context.Context, taskID int64) error {
	// TODO: Implementar lógica
	return nil
}
