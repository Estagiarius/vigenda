package repository

import (
	"context"
	"vigenda/internal/models"
)

//go:generate mockgen -source=subject_repo.go -destination=stubs/subject_repository_mock.go -package=stubs SubjectRepository
type SubjectRepository interface {
	Create(ctx context.Context, subject *models.Subject) error
	GetByID(ctx context.Context, id int64) (*models.Subject, error)
	GetByUserID(ctx context.Context, userID int64) ([]models.Subject, error)
	Update(ctx context.Context, subject *models.Subject) error
	Delete(ctx context.Context, id int64) error
	GetOrCreateByNameAndUser(ctx context.Context, name string, userID int64) (models.Subject, error)
}
