// Package service contém as interfaces e implementações da lógica de negócio da aplicação.
// As interfaces definem os contratos para os serviços, enquanto as implementações
// concretas interagem com a camada de repositório e outras dependências.
package service

import (
	"context"
	"time"
	"vigenda/internal/models"
)

// TaskService define os métodos para a gestão de tarefas.
type TaskService interface {
	CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error)
	ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error) // Returns only active tasks for a class
	ListAllActiveTasks(ctx context.Context) ([]models.Task, error)       // Fetches ALL active tasks (pending)
	ListAllTasks(ctx context.Context) ([]models.Task, error)             // Fetches ALL tasks (pending and completed)
	MarkTaskAsCompleted(ctx context.Context, taskID int64) error
	GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error) // Added for viewing details
	// UpdateTask updates an existing task.
	UpdateTask(ctx context.Context, task *models.Task) error             // Added for updating tasks
	// DeleteTask removes a task by its ID.
	DeleteTask(ctx context.Context, taskID int64) error                  // Added for deleting tasks
}

// ClassService define os métodos para a gestão de turmas e alunos.
type ClassService interface {
	CreateClass(ctx context.Context, name string, subjectID int64) (models.Class, error)
	ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error)
	UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error
	GetClassByID(ctx context.Context, classID int64) (models.Class, error)
	ListAllClasses(ctx context.Context) ([]models.Class, error)
	GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error)
	UpdateClass(ctx context.Context, classID int64, name string, subjectID int64) (models.Class, error)
	DeleteClass(ctx context.Context, classID int64) error
	AddStudent(ctx context.Context, classID int64, fullName string, enrollmentID string, status string) (models.Student, error)
	GetStudentByID(ctx context.Context, studentID int64) (models.Student, error)
	UpdateStudent(ctx context.Context, studentID int64, fullName string, enrollmentID string, status string) (models.Student, error)
	DeleteStudent(ctx context.Context, studentID int64) error
}

// AssessmentService define os métodos para a gestão de avaliações e notas.
type AssessmentService interface {
	CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error)
	EnterGrades(ctx context.Context, assessmentID int64, studentGrades map[int64]float64) error
	CalculateClassAverage(ctx context.Context, classID int64) (float64, error)
	ListAllAssessments(ctx context.Context) ([]models.Assessment, error) // Novo método adicionado
}

// QuestionService define os métodos para o banco de questões e geração de provas.
type QuestionService interface {
	AddQuestionsFromJSON(ctx context.Context, jsonData []byte) (int, error)
	GenerateTest(ctx context.Context, criteria TestCriteria) ([]models.Question, error)
}

// TestCriteria define os parâmetros para a geração de uma prova.
type TestCriteria struct {
	SubjectID   int64
	Topic       *string // Alterado para ponteiro para consistência com ProofCriteria
	EasyCount   int
	MediumCount int
	HardCount   int
}

// ProofService define os métodos para a geração de provas.
type ProofService interface {
	GenerateProof(ctx context.Context, criteria ProofCriteria) ([]models.Question, error)
}

// ProofCriteria define os parâmetros para a geração de uma prova.
type ProofCriteria struct {
	SubjectID   int64
	Topic       *string
	EasyCount   int
	MediumCount int
	HardCount   int
}

// Adicionar outras interfaces de serviço aqui: SubjectService, LessonService, etc.
