package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"vigenda/internal/models"
	"vigenda/internal/repository" // Added import
)

// Implementação do TaskService
type taskServiceImpl struct {
	repo repository.TaskRepository
}

// NewTaskService cria uma nova instância de TaskService.
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskServiceImpl{
		repo: repo,
	}
}

// logError é um helper para logar um erro e retornar.
// No futuro, isso pode ser expandido para usar um logger mais sofisticado.
func logError(format string, args ...interface{}) {
	// Simplesmente imprime para stderr por enquanto.
	// Em uma aplicação real, usaríamos um pacote de logging.
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
}

// handleErrorAndCreateBugTask lida com um erro, loga-o e tenta criar uma tarefa de bug.
// Note que a criação da tarefa de bug em si pode falhar, o que não é tratado recursivamente aqui para evitar loops.
func (s *taskServiceImpl) handleErrorAndCreateBugTask(ctx context.Context, originalError error, bugTitlePrefix string, bugDescriptionArgs ...interface{}) {
	logError("%s: %v", bugTitlePrefix, originalError)

	// Tenta criar uma tarefa para o bug.
	// O UserID para tarefas de sistema/bug pode ser um valor especial (ex: 0 ou um ID de usuário de sistema configurado)
	// ou podemos omiti-lo se o modelo permitir. Assumindo que UserID 0 é para tarefas do sistema.
	// ClassID pode ser nil se o bug não for específico de uma turma.
	bugTitle := fmt.Sprintf("[BUG] %s", bugTitlePrefix)
	bugDescription := fmt.Sprintf("Error encountered: %v. Details: %s", originalError, fmt.Sprintf(bugDescriptionArgs[0].(string), bugDescriptionArgs[1:]...))

	// O UserID para tarefas de sistema/bug pode ser um valor especial (ex: 0 ou um ID de usuário de sistema configurado)
	// Usaremos UserID 0 para tarefas do sistema. ClassID pode ser nil.
	systemUserID := int64(0)
	bugTask := models.Task{
		UserID:      systemUserID,
		ClassID:     nil, // Bugs geralmente não são específicos de uma turma, a menos que o contexto sugira
		Title:       bugTitle,
		Description: bugDescription,
		IsCompleted: false,
	}

	_, creationErr := s.repo.CreateTask(ctx, &bugTask)
	if creationErr != nil {
		// Log an even more critical error if creating the bug task itself fails.
		logError("CRITICAL: Failed to create bug task for '%s': %v. Original error: %v", bugTitle, creationErr, originalError)
	} else {
		logError("SYSTEM: Bug task created successfully for: %s", bugTitle)
	}
}

// CreateTaskInternal is used by handleErrorAndCreateBugTask to avoid recursive bug reporting.
// It directly calls the repository without the surrounding error handling that might create another bug task.
// For normal operations, use CreateTask.
func (s *taskServiceImpl) createTaskInternal(ctx context.Context, userID int64, classID *int64, title, description string, dueDate *time.Time) (models.Task, error) {
	task := models.Task{
		UserID:      userID, // Este deve vir do contexto de autenticação ou similar
		ClassID:     classID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		IsCompleted: false,
	}

	id, err := s.repo.CreateTask(ctx, &task)
	if err != nil {
		return models.Task{}, fmt.Errorf("repository failed to create task: %w", err)
	}
	task.ID = id
	return task, nil
}

func (s *taskServiceImpl) CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error) {
	if title == "" {
		err := errors.New("task title cannot be empty")
		logError("CreateTask validation failed: %v", err) // No bug task for validation errors
		return models.Task{}, err
	}

	// Assume UserID 1 for now, in a real app this would come from auth context
	userID := int64(1)
	task, err := s.createTaskInternal(ctx, userID, classID, title, description, dueDate)
	if err != nil {
		// This is an unexpected error from the internal creation process (e.g., DB error from repo)
		s.handleErrorAndCreateBugTask(ctx, err, "Task Creation Failure", "Attempted to create task with title '%s'. UserID: %d", title, userID)
		return models.Task{}, err // Return the original error
	}
	return task, nil
}

func (s *taskServiceImpl) ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error) {
	tasks, err := s.repo.GetTasksByClassID(ctx, classID)
	if err != nil {
		s.handleErrorAndCreateBugTask(ctx, err, "Task Listing Failure", "Attempted to list tasks for classID %d", classID)
		return nil, err // Return the original error
	}
	// Filter for active tasks (IsCompleted == false)
	// Though the stub GetTasksByClassID doesn't filter, a real repo might.
	// Or, the repo method could be GetActiveTasksByClassID. For now, filter here.
	activeTasks := []models.Task{}
	for _, task := range tasks {
		if !task.IsCompleted {
			activeTasks = append(activeTasks, task)
		}
	}
	return activeTasks, nil
}

func (s *taskServiceImpl) ListAllActiveTasks(ctx context.Context) ([]models.Task, error) {
	tasks, err := s.repo.GetAllTasks(ctx) // Assumes GetAllTasks exists in repo
	if err != nil {
		s.handleErrorAndCreateBugTask(ctx, err, "Task Listing Failure", "Attempted to list all tasks")
		return nil, err
	}
	activeTasks := []models.Task{}
	for _, task := range tasks {
		if !task.IsCompleted {
			activeTasks = append(activeTasks, task)
		}
	}
	return activeTasks, nil
}

func (s *taskServiceImpl) MarkTaskAsCompleted(ctx context.Context, taskID int64) error {
	err := s.repo.MarkTaskCompleted(ctx, taskID)
	if err != nil {
		// Consider if "not found" is a bug or an expected error.
		// For now, let's assume if we try to complete a non-existent task, it's a situation worth logging as a potential bug/issue.
		s.handleErrorAndCreateBugTask(ctx, err, "Task Completion Failure", "Attempted to complete taskID %d", taskID)
		return err // Return the original error
	}
	return nil
}

// Adicionar importações necessárias no topo do arquivo:
// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"os"
// 	"time"
// 	"vigenda/internal/models"
// )
