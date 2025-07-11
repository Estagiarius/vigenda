// Package service define as interfaces para a camada de lógica de negócios da aplicação Vigenda.
// Cada interface de serviço (ex: TaskService, ClassService) agrupa um conjunto de operações
// de negócios relacionadas a uma entidade ou funcionalidade específica. As implementações
// concretas desses serviços (geralmente em arquivos separados, ex: task_service.go)
// orquestram chamadas para a camada de repositório, aplicam regras de negócios,
// validam dados e preparam dados para a camada de apresentação (TUI/CLI).
package service

import (
	"context"
	"time"
	"vigenda/internal/models"
)

// TaskService define a interface para a lógica de negócios relacionada a tarefas.
// Suas implementações orquestram operações como criação, listagem, atualização e conclusão de tarefas.
type TaskService interface {
	// CreateTask cria uma nova tarefa com os detalhes fornecidos.
	// Retorna a tarefa criada ou um erro em caso de falha (ex: validação, erro no repositório).
	CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error)
	// ListActiveTasksByClass retorna uma lista de tarefas ativas (não concluídas) para um ID de turma específico.
	ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error)
	// ListAllActiveTasks retorna uma lista de todas as tarefas ativas (não concluídas) no sistema.
	// Em um sistema multiusuário, isso normalmente seria filtrado pelo usuário autenticado.
	ListAllActiveTasks(ctx context.Context) ([]models.Task, error)
	// ListAllTasks retorna uma lista de todas as tarefas (pendentes e concluídas) no sistema.
	// Em um sistema multiusuário, isso normalmente seria filtrado pelo usuário autenticado.
	ListAllTasks(ctx context.Context) ([]models.Task, error)
	// MarkTaskAsCompleted marca uma tarefa específica como concluída.
	MarkTaskAsCompleted(ctx context.Context, taskID int64) error
	// GetTaskByID recupera os detalhes de uma tarefa específica pelo seu ID.
	GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error)
	// UpdateTask atualiza os detalhes de uma tarefa existente.
	// O objeto models.Task fornecido deve conter o ID da tarefa a ser atualizada e os novos valores.
	UpdateTask(ctx context.Context, task *models.Task) error
	// DeleteTask remove uma tarefa do sistema pelo seu ID.
	DeleteTask(ctx context.Context, taskID int64) error
	// GetUpcomingActiveTasks recupera uma lista limitada de tarefas ativas futuras para um usuário específico.
	GetUpcomingActiveTasks(ctx context.Context, userID int64, fromDate time.Time, limit int) ([]models.Task, error)
}

// ClassService define a interface para a lógica de negócios relacionada a turmas e alunos.
type ClassService interface {
	// CreateClass cria uma nova turma associada a uma disciplina.
	CreateClass(ctx context.Context, name string, subjectID int64) (models.Class, error)
	// ImportStudentsFromCSV importa alunos de um arquivo CSV para uma turma existente.
	// Retorna o número de alunos importados com sucesso ou um erro.
	ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error)
	// UpdateStudentStatus atualiza o status de um aluno (ex: 'ativo', 'inativo', 'transferido').
	UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error
	// GetClassByID recupera os detalhes de uma turma específica.
	GetClassByID(ctx context.Context, classID int64) (models.Class, error)
	// ListAllClasses retorna uma lista de todas as turmas.
	// Em um sistema multiusuário, isso seria filtrado pelo usuário autenticado.
	ListAllClasses(ctx context.Context) ([]models.Class, error)
	// GetStudentsByClassID retorna uma lista de todos os alunos de uma turma específica.
	GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error)
	// UpdateClass atualiza os detalhes de uma turma existente.
	UpdateClass(ctx context.Context, classID int64, name string, subjectID int64) (models.Class, error)
	// DeleteClass remove uma turma e seus dados associados (alunos, avaliações, etc., dependendo do ON DELETE CASCADE do DB).
	DeleteClass(ctx context.Context, classID int64) error
	// AddStudent adiciona um novo aluno a uma turma.
	AddStudent(ctx context.Context, classID int64, fullName string, enrollmentID string, status string) (models.Student, error)
	// GetStudentByID recupera os detalhes de um aluno específico.
	GetStudentByID(ctx context.Context, studentID int64) (models.Student, error)
	// UpdateStudent atualiza os detalhes de um aluno.
	UpdateStudent(ctx context.Context, studentID int64, fullName string, enrollmentID string, status string) (models.Student, error)
	// DeleteStudent remove um aluno de uma turma.
	DeleteStudent(ctx context.Context, studentID int64) error
}

// AssessmentService define a interface para a lógica de negócios relacionada a avaliações e notas.
type AssessmentService interface {
	// CreateAssessment cria uma nova avaliação para uma turma.
	CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error)
	// EnterGrades registra ou atualiza as notas de múltiplos alunos para uma avaliação específica.
	// studentGrades é um mapa onde a chave é o StudentID e o valor é a nota.
	EnterGrades(ctx context.Context, assessmentID int64, studentGrades map[int64]float64) error
	// CalculateClassAverage calcula a média ponderada das notas de uma turma.
	// O cálculo considera todas as avaliações e seus pesos para a turma especificada.
	CalculateClassAverage(ctx context.Context, classID int64) (float64, error)
	// ListAllAssessments retorna uma lista de todas as avaliações.
	// Em um sistema multiusuário, isso seria filtrado pelo usuário ou turma.
	ListAllAssessments(ctx context.Context) ([]models.Assessment, error)
}

// QuestionService define a interface para a lógica de negócios relacionada ao banco de questões.
// Atualmente, inclui a importação de questões e a geração de provas (que foi movida para ProofService).
// TODO: Avaliar se GenerateTest deve ser mantido aqui ou se ProofService é suficiente.
type QuestionService interface {
	// AddQuestionsFromJSON importa questões de um payload JSON para o banco de dados.
	// Retorna o número de questões importadas com sucesso ou um erro.
	AddQuestionsFromJSON(ctx context.Context, jsonData []byte) (int, error)
	// GenerateTest gera uma lista de questões baseada nos critérios fornecidos.
	// Esta funcionalidade pode estar sobreposta com ProofService.GenerateProof.
	GenerateTest(ctx context.Context, criteria TestCriteria) ([]models.Question, error)
}

// TestCriteria define os parâmetros para a geração de uma prova/teste.
// Usado por QuestionService.GenerateTest.
type TestCriteria struct {
	SubjectID   int64   // SubjectID é o ID da disciplina para filtrar as questões.
	Topic       *string // Topic (opcional) é um tópico específico dentro da disciplina.
	EasyCount   int     // EasyCount é o número desejado de questões fáceis.
	MediumCount int     // MediumCount é o número desejado de questões médias.
	HardCount   int     // HardCount é o número desejado de questões difíceis.
}

// ProofService define a interface para a lógica de negócios específica da geração de provas.
type ProofService interface {
	// GenerateProof gera uma lista de questões para uma prova com base nos critérios fornecidos.
	// Os critérios incluem o ID da disciplina, tópico opcional e contagem de questões por dificuldade.
	GenerateProof(ctx context.Context, criteria ProofCriteria) ([]models.Question, error)
}

// ProofCriteria define os parâmetros para a geração de uma prova.
// Usado por ProofService.GenerateProof.
type ProofCriteria struct {
	SubjectID   int64   // SubjectID é o ID da disciplina para a qual a prova será gerada.
	Topic       *string // Topic (opcional) filtra questões por um tópico específico.
	EasyCount   int     // EasyCount é o número desejado de questões fáceis.
	MediumCount int     // MediumCount é o número desejado de questões médias.
	HardCount   int     // HardCount é o número desejado de questões difíceis.
}

// LessonService define a interface para a lógica de negócios relacionada a aulas/lições.
type LessonService interface {
	// CreateLesson cria uma nova aula/lição para uma turma.
	CreateLesson(ctx context.Context, classID int64, title string, planContent string, scheduledAt time.Time) (models.Lesson, error)
	// GetLessonByID recupera os detalhes de uma aula/lição específica.
	GetLessonByID(ctx context.Context, lessonID int64) (models.Lesson, error)
	// GetLessonsByClassID retorna uma lista de todas as aulas/lições de uma turma específica.
	GetLessonsByClassID(ctx context.Context, classID int64) ([]models.Lesson, error)
	// GetLessonsForDate busca aulas/lições para um usuário em uma data específica.
	// O UserID é usado para filtrar; em um sistema multiusuário, viria do contexto de autenticação.
	GetLessonsForDate(ctx context.Context, userID int64, date time.Time) ([]models.Lesson, error)
	// UpdateLesson atualiza os detalhes de uma aula/lição existente.
	UpdateLesson(ctx context.Context, lessonID int64, title string, planContent string, scheduledAt time.Time) (models.Lesson, error)
	// DeleteLesson remove uma aula/lição do sistema.
	DeleteLesson(ctx context.Context, lessonID int64) error
}

// TODO: Adicionar SubjectService interface para gerenciar CRUD de Disciplinas.
// Exemplo:
// type SubjectService interface {
//    CreateSubject(ctx context.Context, userID int64, name string) (models.Subject, error)
//    GetSubjectByID(ctx context.Context, subjectID int64) (models.Subject, error)
//    ListSubjectsByUser(ctx context.Context, userID int64) ([]models.Subject, error)
//    UpdateSubject(ctx context.Context, subjectID int64, name string) (models.Subject, error)
//    DeleteSubject(ctx context.Context, subjectID int64) error
// }
