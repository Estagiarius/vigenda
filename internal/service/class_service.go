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
