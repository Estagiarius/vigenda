// Package repository define as interfaces para a camada de acesso a dados.
// Estas interfaces abstraem as operações de banco de dados para os serviços.
package repository

import (
	"context"
	"vigenda/internal/models"
	"time" // Movido para o final do bloco de import para forçar reavaliação
)

// QuestionQueryCriteria define os critérios para buscar questões.
type QuestionQueryCriteria struct {
	SubjectID  int64
	Topic      *string // Ponteiro para permitir valor nulo/opcional
	Difficulty string
	Limit      int
	// Outros campos como Type podem ser adicionados aqui se necessário.
}

// ProofCriteria defines the criteria for fetching questions for proof generation at the repository level.
type ProofCriteria struct {
	SubjectID   int64
	Topic       *string
	EasyCount   int
	MediumCount int
	HardCount   int
}

// QuestionRepository define a interface para operações de acesso a dados de questões.
type QuestionRepository interface {
	GetQuestionsByCriteria(ctx context.Context, criteria QuestionQueryCriteria) ([]models.Question, error)
	AddQuestion(ctx context.Context, question *models.Question) (int64, error)
	GetQuestionsByCriteriaProofGeneration(ctx context.Context, criteria ProofCriteria) ([]models.Question, error) // Changed to use repository.ProofCriteria
	// GetQuestionByID(ctx context.Context, id int64) (models.Question, error)
	// UpdateQuestion(ctx context.Context, question *models.Question) error
	// DeleteQuestion(ctx context.Context, id int64) error
	// AddManyQuestions(ctx context.Context, questions []models.Question) (int, error) // Para importação em lote
}

// SubjectRepository define a interface para operações de acesso a dados de disciplinas.
// Necessário para validar a existência de uma disciplina ao adicionar questões, por exemplo.
type SubjectRepository interface {
	// GetSubjectByID(ctx context.Context, id int64) (models.Subject, error)
	// GetSubjectByName(ctx context.Context, name string) (models.Subject, error) // Pode ser útil
	GetOrCreateByNameAndUser(ctx context.Context, name string, userID int64) (models.Subject, error) // Adicionado para QuestionService
	// ... outros métodos CRUD para Subject
}

// Outras interfaces de repositório (TaskRepository, ClassRepository, etc.) seriam definidas aqui.
// Por exemplo:

type TaskRepository interface {
	CreateTask(ctx context.Context, task *models.Task) (int64, error)
	GetTaskByID(ctx context.Context, id int64) (*models.Task, error)
	GetTasksByClassID(ctx context.Context, classID int64) ([]models.Task, error) // Added for listing tasks
	GetAllTasks(ctx context.Context) ([]models.Task, error)                       // New method
	MarkTaskCompleted(ctx context.Context, taskID int64) error                   // Added for completing tasks
	// UpdateTask updates an existing task in the database.
	UpdateTask(ctx context.Context, task *models.Task) error                     // Added for updating tasks
	// DeleteTask removes a task from the database by its ID.
	DeleteTask(ctx context.Context, taskID int64) error                          // Added for deleting tasks
	GetUpcomingActiveTasks(ctx context.Context, userID int64, fromDate time.Time, limit int) ([]models.Task, error)
	// ... outros métodos
}

//go:generate mockgen -source=repository.go -destination=stubs/class_repository_mock.go -package=stubs ClassRepository
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

// LessonRepository define a interface para operações de persistência de aulas/lições.
type LessonRepository interface {
	CreateLesson(ctx context.Context, lesson *models.Lesson) (int64, error)
	GetLessonByID(ctx context.Context, lessonID int64) (*models.Lesson, error)
	GetLessonsByClassID(ctx context.Context, classID int64) ([]models.Lesson, error)
	// GetLessonsByDateRange busca lições para um usuário dentro de um intervalo de datas.
	// O userID pode ser usado para filtrar lições pertencentes a um usuário específico.
	GetLessonsByDateRange(ctx context.Context, userID int64, startDate time.Time, endDate time.Time) ([]models.Lesson, error)
	UpdateLesson(ctx context.Context, lesson *models.Lesson) error
	DeleteLesson(ctx context.Context, lessonID int64) error
}

type AssessmentRepository interface {
	CreateAssessment(ctx context.Context, assessment *models.Assessment) (int64, error)
	GetAssessmentByID(ctx context.Context, assessmentID int64) (*models.Assessment, error)
	// GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) // Moved to ClassRepository
	EnterGrade(ctx context.Context, grade *models.Grade) error
	GetGradesByClassID(ctx context.Context, classID int64) ([]models.Grade, []models.Assessment, []models.Student, error) // For calculating class average
	ListAllAssessments(ctx context.Context) ([]models.Assessment, error) // Novo método adicionado
	// GetAssessmentWithGrades(ctx context.Context, assessmentID int64) (*models.AssessmentWithGrades, error) // Example for a more complex query
}

// Estas são apenas exemplos e devem ser expandidas conforme necessário
// para cobrir todas as operações de dados requeridas pelos serviços.
// A TASK-I-01 é responsável pela implementação concreta destes repositórios.
//
// A struct models.Subject também precisaria ser definida em `internal/models/models.go`:
//
// package models
//
// type Subject struct {
//  ID     int64  `json:"id"`
//  UserID int64  `json:"user_id"`
//  Name   string `json:"name"`
// }
//
// Esta definição é importante para o QuestionService ao adicionar questões,
// para garantir que a disciplina referenciada exista.
