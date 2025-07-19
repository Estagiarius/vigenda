// Package repository define as interfaces para a camada de acesso a dados (Data Access Layer - DAL).
// Estas interfaces abstraem as operações de banco de dados (CRUD - Create, Read, Update, Delete, e outras consultas específicas)
// para as entidades do domínio. As implementações concretas destes repositórios (ex: TaskSQLRepository)
// interagem diretamente com o banco de dados, permitindo que a lógica de negócios nos serviços
// permaneça agnóstica em relação à tecnologia de persistência específica.
package repository

import (
	"context"
	"vigenda/internal/models"
	"time"
)

// QuestionQueryCriteria define os critérios para buscar questões no banco de dados.
// É usado principalmente pelo QuestionRepository.
type QuestionQueryCriteria struct {
	SubjectID  int64   // SubjectID é o ID da disciplina para filtrar as questões. Obrigatório.
	Topic      *string // Topic (opcional) é um tópico específico dentro da disciplina. Ponteiro para permitir valor nulo.
	Difficulty string  // Difficulty (opcional) filtra questões por nível de dificuldade (ex: "facil", "media", "dificil").
	Limit      int     // Limit (opcional) especifica o número máximo de questões a serem retornadas.
	// Type (opcional) poderia ser adicionado para filtrar por tipo de questão (ex: 'multipla_escolha').
}

// ProofCriteria define os critérios para buscar questões especificamente para a geração de provas.
// Difere de QuestionQueryCriteria por permitir a especificação de contagens por dificuldade.
type ProofCriteria struct {
	SubjectID   int64   // SubjectID é o ID da disciplina para a qual a prova será gerada. Obrigatório.
	Topic       *string // Topic (opcional) filtra questões por um tópico específico.
	EasyCount   int     // EasyCount é o número desejado de questões fáceis.
	MediumCount int     // MediumCount é o número desejado de questões médias.
	HardCount   int     // HardCount é o número desejado de questões difíceis.
}


// TaskRepository define a interface para operações de acesso a dados relacionadas a 'tasks' (tarefas).
type TaskRepository interface {
	// CreateTask adiciona uma nova tarefa ao banco de dados e retorna seu ID.
	CreateTask(ctx context.Context, task *models.Task) (int64, error)
	// GetTaskByID recupera uma tarefa específica por seu ID. Retorna nil se não encontrada.
	GetTaskByID(ctx context.Context, id int64) (*models.Task, error)
	// GetTasksByClassID recupera todas as tarefas associadas a um ClassID específico.
	GetTasksByClassID(ctx context.Context, classID int64) ([]models.Task, error)
	// GetAllTasks recupera todas as tarefas do banco de dados (pode precisar de filtragem por usuário em implementações reais).
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	// MarkTaskCompleted marca uma tarefa como concluída.
	MarkTaskCompleted(ctx context.Context, taskID int64) error
	// UpdateTask atualiza os detalhes de uma tarefa existente no banco de dados.
	UpdateTask(ctx context.Context, task *models.Task) error
	// DeleteTask remove uma tarefa do banco de dados pelo seu ID.
	DeleteTask(ctx context.Context, taskID int64) error
	// GetUpcomingActiveTasks recupera tarefas ativas (não concluídas) de um usuário específico
	// com data de vencimento a partir de 'fromDate', limitadas por 'limit'.
	GetUpcomingActiveTasks(ctx context.Context, userID int64, fromDate time.Time, limit int) ([]models.Task, error)
}

// LessonRepository define a interface para operações de persistência relacionadas a 'lessons' (aulas/lições).
type LessonRepository interface {
	// CreateLesson adiciona uma nova aula/lição ao banco de dados e retorna seu ID.
	CreateLesson(ctx context.Context, lesson *models.Lesson) (int64, error)
	// GetLessonByID recupera uma aula/lição específica por seu ID.
	GetLessonByID(ctx context.Context, lessonID int64) (*models.Lesson, error)
	// GetLessonsByClassID recupera todas as aulas/lições associadas a um ClassID específico.
	GetLessonsByClassID(ctx context.Context, classID int64) ([]models.Lesson, error)
	// GetLessonsByDateRange busca aulas/lições para um usuário dentro de um intervalo de datas.
	// O userID é usado para garantir que apenas as aulas do usuário sejam retornadas.
	GetLessonsByDateRange(ctx context.Context, userID int64, startDate time.Time, endDate time.Time) ([]models.Lesson, error)
	// UpdateLesson atualiza os detalhes de uma aula/lição existente.
	UpdateLesson(ctx context.Context, lesson *models.Lesson) error
	// DeleteLesson remove uma aula/lição do banco de dados pelo seu ID.
	DeleteLesson(ctx context.Context, lessonID int64) error
}

// AssessmentRepository define a interface para operações de acesso a dados relacionadas a 'assessments' (avaliações) e 'grades' (notas).
type AssessmentRepository interface {
	// CreateAssessment adiciona uma nova avaliação ao banco de dados e retorna seu ID.
	CreateAssessment(ctx context.Context, assessment *models.Assessment) (int64, error)
	// GetAssessmentByID recupera uma avaliação específica por seu ID.
	GetAssessmentByID(ctx context.Context, assessmentID int64) (*models.Assessment, error)
	// EnterGrade registra ou atualiza a nota de um aluno para uma avaliação.
	EnterGrade(ctx context.Context, grade *models.Grade) error
	// GetGradesByClassID recupera todas as notas, avaliações e alunos de uma turma específica.
	// Usado para calcular a média da turma, pois necessita de todas essas informações.
	GetGradesByClassID(ctx context.Context, classID int64) ([]models.Grade, []models.Assessment, []models.Student, error)
	// ListAllAssessments recupera todas as avaliações (pode precisar de filtragem por usuário ou turma).
	ListAllAssessments(ctx context.Context) ([]models.Assessment, error)
	// DeleteAssessment remove uma avaliação e suas notas associadas (via ON DELETE CASCADE no DB).
	DeleteAssessment(ctx context.Context, assessmentID int64) error
	// FindAssessmentByNameAndClass busca uma avaliação específica pelo nome e ID da turma.
	FindAssessmentByNameAndClass(ctx context.Context, name string, classID int64) (*models.Assessment, error)
	// GetGradesByAssessmentID recupera todas as notas para uma avaliação específica.
	GetGradesByAssessmentID(ctx context.Context, assessmentID int64) ([]models.Grade, error)
	// GetAssessmentWithGrades (Comentado) poderia ser um exemplo de consulta mais complexa,
	// retornando uma avaliação junto com todas as suas notas associadas.
	// GetAssessmentWithGrades(ctx context.Context, assessmentID int64) (*models.AssessmentWithGrades, error)
}
