package repository

import (
	"context"
	"vigenda/internal/models"
)

//go:generate mockgen -source=class_repo.go -destination=stubs/class_repository_mock.go -package=stubs ClassRepository
type ClassRepository interface {
	CreateClass(ctx context.Context, class *models.Class) (int64, error)
	GetClassByID(ctx context.Context, id int64) (*models.Class, error)
	AddStudent(ctx context.Context, student *models.Student) (int64, error)
	UpdateStudentStatus(ctx context.Context, studentID int64, status string) error
	ListAllClasses(ctx context.Context) ([]models.Class, error)
	GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error)
	UpdateClass(ctx context.Context, class *models.Class) error
	DeleteClass(ctx context.Context, classID int64, userID int64) error
	GetStudentByID(ctx context.Context, studentID int64) (*models.Student, error)
	UpdateStudent(ctx context.Context, student *models.Student) error
	DeleteStudent(ctx context.Context, studentID int64, classID int64) error
}
