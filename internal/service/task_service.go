package service

import (
	"context"
	"time"
	"vigenda/internal/models"
)

// TaskService define os métodos para a gestão de tarefas.
type TaskService interface {
	CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error)
	ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error)
	MarkTaskAsCompleted(ctx context.Context, taskID int64) error
}
