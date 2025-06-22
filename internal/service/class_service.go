package service

import (
	"context"
	"vigenda/internal/models"
)

// ClassService define os métodos para a gestão de turmas e alunos.
type ClassService interface {
	CreateClass(ctx context.Context, name string, subjectID int64) (models.Class, error)
	ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error)
	UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error
}

type classServiceImpl struct {
	// Add repository dependencies here
}

// NewClassService creates a new instance of ClassService.
func NewClassService() ClassService {
	return &classServiceImpl{}
}

func (s *classServiceImpl) CreateClass(ctx context.Context, name string, subjectID int64) (models.Class, error) {
	// TODO: Implement actual logic
	return models.Class{}, nil
}

func (s *classServiceImpl) ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error) {
	// TODO: Implement actual logic
	return 0, nil
}

func (s *classServiceImpl) UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error {
	// TODO: Implement actual logic
	return nil
}
