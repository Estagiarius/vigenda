// Package repository define as interfaces para a camada de acesso a dados.
// Estas interfaces abstraem as operações de banco de dados para os serviços.
package repository

import (
	"context"
	"vigenda/internal/models"
)

// QuestionQueryCriteria define os critérios para buscar questões.
type QuestionQueryCriteria struct {
	SubjectID  int64
	Topic      *string // Ponteiro para permitir valor nulo/opcional
	Difficulty string
	Limit      int
	// Outros campos como Type podem ser adicionados aqui se necessário.
}

// QuestionRepository define a interface para operações de acesso a dados de questões.
type QuestionRepository interface {
	GetQuestionsByCriteria(ctx context.Context, criteria QuestionQueryCriteria) ([]models.Question, error)
	AddQuestion(ctx context.Context, question *models.Question) (int64, error)
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
	// ... outros métodos
}

// type ClassRepository interface {
//    CreateClass(ctx context.Context, class *models.Class) (int64, error)
//    GetClassByID(ctx context.Context, id int64) (*models.Class, error)
//    // ... outros métodos
// }
//
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
