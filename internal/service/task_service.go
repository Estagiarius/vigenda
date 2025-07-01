package service

import (
	"context"
	"errors"
	"fmt"
	"database/sql" // Added for sql.ErrNoRows
	"errors"
	"fmt"
	"os"
	"strings" // Added for strings.Contains
	"time"
	"vigenda/internal/models"
	"vigenda/internal/repository" // Added import
)

// taskServiceImpl é a implementação concreta de TaskService.
// Ela encapsula a lógica de negócios relacionada a tarefas e interage
// com a camada de repositório para persistência de dados.
type taskServiceImpl struct {
	repo repository.TaskRepository // repo é a instância do repositório de tarefas.
}

// NewTaskService cria e retorna uma nova instância de TaskService.
// Recebe um repository.TaskRepository como dependência para interagir com a camada de dados.
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskServiceImpl{
		repo: repo,
	}
}

// logError é uma função auxiliar interna para logar erros.
// Atualmente, ela imprime para stderr. Em uma aplicação de produção,
// isso seria substituído por um sistema de logging mais robusto (ex: slog, zerolog).
func logError(format string, args ...interface{}) {
	// TODO: Substituir por um logger estruturado em futuras iterações.
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
}

// handleErrorAndCreateBugTask é uma função auxiliar para tratar erros inesperados
// ocorridos durante operações críticas do serviço. Ela loga o erro original
// e tenta criar uma nova tarefa no sistema para rastrear o bug, facilitando
// a depuração e correção.
//
// Parâmetros:
//   - ctx: O contexto da requisição.
//   - originalError: O erro que ocorreu.
//   - bugTitlePrefix: Um prefixo para o título da tarefa de bug a ser criada.
//   - bugDescriptionArgs: Argumentos para formatar a descrição da tarefa de bug.
//     O primeiro argumento deve ser a string de formato, seguido pelos valores.
//
// A criação da tarefa de bug em si pode falhar (ex: problema de conexão com DB),
// nesse caso, um log crítico adicional é gerado.
func (s *taskServiceImpl) handleErrorAndCreateBugTask(ctx context.Context, originalError error, bugTitlePrefix string, bugDescriptionArgs ...interface{}) {
	logError("%s: %v", bugTitlePrefix, originalError) // Loga o erro original.

	// Formata o título e a descrição para a tarefa de bug.
	bugTitle := fmt.Sprintf("[BUG][AUTO][PRIORITY_PENDING] %s", bugTitlePrefix)
	var bugDescription string
	if len(bugDescriptionArgs) > 0 {
		format, ok := bugDescriptionArgs[0].(string)
		if ok && len(bugDescriptionArgs) > 1 {
			bugDescription = fmt.Sprintf("[PRIORITY_PENDING] Error encountered: %v. Details: %s", originalError, fmt.Sprintf(format, bugDescriptionArgs[1:]...))
		} else {
			bugDescription = fmt.Sprintf("[PRIORITY_PENDING] Error encountered: %v. No additional details provided.", originalError)
		}
	} else {
		bugDescription = fmt.Sprintf("[PRIORITY_PENDING] Error encountered: %v.", originalError)
	}

	// UserID 0 é usado para tarefas de sistema/bugs. ClassID é nil, pois bugs geralmente não são específicos de uma turma.
	systemUserID := int64(0)
	bugTask := models.Task{
		UserID:      systemUserID,
		ClassID:     nil,
		Title:       bugTitle,
		Description: bugDescription,
		IsCompleted: false,
	}

	// Tenta criar a tarefa de bug usando o método interno do repositório.
	_, creationErr := s.repo.CreateTask(ctx, &bugTask) // Assume que CreateTask não está dentro de uma transação aqui ou é seguro chamar.
	if creationErr != nil {
		logError("CRITICAL: Failed to create bug task for '%s': %v. Original error: %v", bugTitle, creationErr, originalError)
	} else {
		logError("SYSTEM: Bug task created successfully for: %s", bugTitle)
	}
}

// createTaskInternal é uma versão interna de CreateTask usada principalmente por
// handleErrorAndCreateBugTask para evitar recursão infinita na criação de tarefas de bug.
// Este método chama diretamente o repositório sem a lógica de tratamento de erro
// que poderia, por sua vez, tentar criar outra tarefa de bug.
// Para operações normais de criação de tarefas pelo usuário, use CreateTask.
//
// Retorna a tarefa criada ou um erro se a criação no repositório falhar.
func (s *taskServiceImpl) createTaskInternal(ctx context.Context, userID int64, classID *int64, title, description string, dueDate *time.Time) (models.Task, error) {
	task := models.Task{
		UserID:      userID,
		ClassID:     classID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		IsCompleted: false,
	}

	id, err := s.repo.CreateTask(ctx, &task)
	if err != nil {
		// Envolve o erro do repositório para fornecer mais contexto.
		return models.Task{}, fmt.Errorf("repository failed to create task: %w", err)
	}
	task.ID = id // Atribui o ID gerado pelo banco de dados à struct da tarefa.
	return task, nil
}

// CreateTask cria uma nova tarefa no sistema.
// Valida se o título da tarefa não está vazio.
// Em caso de falha na criação (ex: erro no banco de dados), loga o erro e
// tenta criar uma tarefa de bug para rastreamento.
//
// Parâmetros:
//   - ctx: O contexto da requisição.
//   - title: O título da tarefa (obrigatório).
//   - description: A descrição da tarefa (opcional).
//   - classID: O ID da turma à qual a tarefa está associada (opcional).
//   - dueDate: A data de vencimento da tarefa (opcional).
//
// Retorna a tarefa criada ou um erro.
func (s *taskServiceImpl) CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error) {
	if title == "" {
		err := errors.New("task title cannot be empty")
		// Erros de validação não devem criar tarefas de bug automaticamente, pois são erros esperados do usuário.
		logError("CreateTask validation failed: %v", err)
		return models.Task{}, err
	}

	// Em uma aplicação real, UserID viria do contexto de autenticação/autorização.
	// Para este exemplo, um UserID fixo é usado.
	// TODO: Integrar com sistema de autenticação para obter UserID real.
	userID := int64(1)

	task, err := s.createTaskInternal(ctx, userID, classID, title, description, dueDate)
	if err != nil {
		// Erro inesperado durante a criação interna (ex: falha no DB vindo do repositório).
		s.handleErrorAndCreateBugTask(ctx, err, "Task Creation Failure", "Attempted to create task with title '%s'. UserID: %d", title, userID)
		return models.Task{}, err // Retorna o erro original para o chamador.
	}
	return task, nil
}

// UpdateTask atualiza uma tarefa existente no sistema.
// Valida se o título da tarefa não está vazio.
// Em caso de falha na atualização (ex: erro no banco de dados), loga o erro e
// tenta criar uma tarefa de bug para rastreamento.
//
// Retorna um erro se a atualização falhar.
func (s *taskServiceImpl) UpdateTask(ctx context.Context, task *models.Task) error {
	if task.Title == "" {
		err := errors.New("task title cannot be empty for update")
		logError("UpdateTask validation failed for Task ID %d: %v", task.ID, err)
		return err
	}

	// UserID should ideally be checked against the authenticated user,
	// but for now, we assume the task object carries the correct UserID from when it was fetched.
	// Or, it could be re-fetched here to ensure integrity if UserID is part of the update logic.

	err := s.repo.UpdateTask(ctx, task)
	if err != nil {
		// Check if the error is due to "not found" or "no change" which might not warrant a bug task.
		if strings.Contains(err.Error(), "no task found") || strings.Contains(err.Error(), "no values changed") {
			logError("UpdateTask: Failed to update task ID %d: %v", task.ID, err)
		} else {
			// For other unexpected errors.
			s.handleErrorAndCreateBugTask(ctx, err, "Task Update Failure", "Attempted to update task with ID %d, Title '%s'. UserID: %d", task.ID, task.Title, task.UserID)
		}
		return err // Retorna o erro original para o chamador.
	}
	return nil
}

// DeleteTask remove uma tarefa do sistema.
// Em caso de falha (ex: tarefa não encontrada ou erro no DB), loga o erro
// e tenta criar uma tarefa de bug, a menos que o erro seja simplesmente "não encontrado".
//
// Retorna um erro se a operação falhar.
func (s *taskServiceImpl) DeleteTask(ctx context.Context, taskID int64) error {
	err := s.repo.DeleteTask(ctx, taskID)
	if err != nil {
		if strings.Contains(err.Error(), "no task found") { // Check for specific error text from repo
			logError("DeleteTask: Failed to delete task ID %d: %v", taskID, err)
		} else {
			s.handleErrorAndCreateBugTask(ctx, err, "Task Deletion Failure", "Attempted to delete taskID %d", taskID)
		}
		return err // Retorna o erro original.
	}
	return nil
}

// ListActiveTasksByClass recupera todas as tarefas ativas (não concluídas)
// associadas a um ID de turma específico.
// Em caso de falha na listagem, loga o erro e tenta criar uma tarefa de bug.
//
// Retorna uma lista de tarefas ativas ou um erro.
func (s *taskServiceImpl) ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error) {
	tasks, err := s.repo.GetTasksByClassID(ctx, classID)
	if err != nil {
		s.handleErrorAndCreateBugTask(ctx, err, "Task Listing By Class Failure", "Attempted to list tasks for classID %d", classID)
		return nil, err // Retorna o erro original.
	}

	// Filtra as tarefas para retornar apenas as ativas (IsCompleted == false).
	// Idealmente, o repositório poderia ter um método como GetActiveTasksByClassID
	// para evitar a filtragem na camada de serviço.
	activeTasks := make([]models.Task, 0, len(tasks)) // Prealoca slice com capacidade.
	for _, task := range tasks {
		if !task.IsCompleted {
			activeTasks = append(activeTasks, task)
		}
	}
	return activeTasks, nil
}

// ListAllActiveTasks recupera todas as tarefas ativas (não concluídas) de todos os usuários.
// TODO: Em um sistema multiusuário, isso deveria ser filtrado pelo UserID do contexto.
// Em caso de falha na listagem, loga o erro e tenta criar uma tarefa de bug.
//
// Retorna uma lista de todas as tarefas ativas ou um erro.
func (s *taskServiceImpl) ListAllActiveTasks(ctx context.Context) ([]models.Task, error) {
	// TODO: Adicionar filtragem por UserID quando a autenticação estiver implementada.
	// Por enquanto, busca todas as tarefas, o que pode não ser ideal.
	tasks, err := s.repo.GetAllTasks(ctx) // Assumindo que GetAllTasks existe no repositório.
	if err != nil {
		s.handleErrorAndCreateBugTask(ctx, err, "Global Task Listing Failure", "Attempted to list all tasks")
		return nil, err
	}

	activeTasks := make([]models.Task, 0, len(tasks))
	for _, task := range tasks {
		if !task.IsCompleted {
			activeTasks = append(activeTasks, task)
		}
	}
	return activeTasks, nil
}

// MarkTaskAsCompleted marca uma tarefa específica como concluída.
// O ID da tarefa é usado para identificar a tarefa a ser atualizada.
// Em caso de falha (ex: tarefa não encontrada ou erro no DB), loga o erro
// e tenta criar uma tarefa de bug.
//
// Retorna um erro se a operação falhar.
func (s *taskServiceImpl) MarkTaskAsCompleted(ctx context.Context, taskID int64) error {
	err := s.repo.MarkTaskCompleted(ctx, taskID)
	if err != nil {
		// Considera-se que tentar completar uma tarefa inexistente pode ser um bug
		// ou um problema de integridade de dados, justificando uma tarefa de bug.
		s.handleErrorAndCreateBugTask(ctx, err, "Task Completion Failure", "Attempted to complete taskID %d", taskID)
		return err // Retorna o erro original.
	}
	return nil
}

// GetTaskByID recupera uma tarefa específica pelo seu ID.
// Em caso de falha (ex: tarefa não encontrada ou erro no DB), loga o erro
// e tenta criar uma tarefa de bug.
//
// Retorna a tarefa encontrada ou um erro.
func (s *taskServiceImpl) GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		// Tratar erro de "não encontrado" de forma diferente de outros erros de DB.
		// Não criar bug task se for apenas "não encontrado" pois pode ser um ID inválido fornecido pelo usuário.
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "no task found") { // Check for specific error text from repo
			logError("GetTaskByID: Task not found with ID %d: %v", taskID, err)
			// Retornar um erro específico que a camada de TUI possa interpretar como "não encontrado"
			return nil, fmt.Errorf("tarefa com ID %d não encontrada: %w", taskID, err)
		}
		// Para outros erros (inesperados), criar bug task.
		s.handleErrorAndCreateBugTask(ctx, err, "Task Retrieval Failure", "Attempted to retrieve taskID %d", taskID)
		return nil, err // Retorna o erro original.
	}
	return task, nil
}
