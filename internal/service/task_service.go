// Package service contém as implementações concretas das interfaces de serviço.
// Este arquivo específico implementa a interface TaskService.
package service

import (
	"context"
	"database/sql" // Necessário para sql.ErrNoRows ao verificar erros do repositório.
	"errors"       // Para criar erros de validação padrão.
	"fmt"
	"os"      // Usado temporariamente para logError.
	"strings" // Usado para verificar mensagens de erro específicas.
	"time"
	"vigenda/internal/models"
	"vigenda/internal/repository"
)

// taskServiceImpl é a implementação concreta da interface TaskService.
// Ela encapsula a lógica de negócios relacionada a tarefas, como validação de dados
// e coordenação com a camada de repositório para persistência.
type taskServiceImpl struct {
	repo repository.TaskRepository // repo é a instância do repositório de tarefas, usada para interagir com o banco de dados.
}

// NewTaskService é uma função construtora que cria e retorna uma nova instância de TaskService.
// Recebe um repository.TaskRepository como dependência (injeção de dependência),
// o que permite que a camada de serviço seja desacoplada da implementação específica do banco de dados
// e facilita os testes unitários através do uso de mocks de repositório.
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskServiceImpl{
		repo: repo,
	}
}

// logError é uma função auxiliar interna para registrar erros.
// Atualmente, imprime para stderr. Em uma aplicação de produção,
// isso seria substituído por um sistema de logging mais robusto e configurável
// (ex: usando pacotes como slog, zerolog, ou zap).
// TODO: Substituir por um logger estruturado em futuras iterações.
func logError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "SERVICE_ERROR: "+format+"\n", args...)
}

// handleErrorAndCreateBugTask é uma função utilitária para tratar erros inesperados
// que ocorrem durante operações críticas do serviço. Ela:
// 1. Loga o erro original para fins de depuração.
// 2. Tenta criar uma nova tarefa no sistema com um título e descrição padronizados
//    para rastrear o bug. Isso ajuda a garantir que erros internos sejam visíveis
//    e possam ser investigados posteriormente.
// A UserID 0 é usada para indicar que é uma tarefa de sistema/bug.
// ClassID é nil, pois bugs geralmente não são específicos de uma turma.
// A criação da tarefa de bug em si pode falhar (ex: problema de conexão com DB),
// nesse caso, um log crítico adicional é gerado.
func (s *taskServiceImpl) handleErrorAndCreateBugTask(ctx context.Context, originalError error, bugTitlePrefix string, bugDescriptionArgs ...interface{}) {
	logError("%s: %v", bugTitlePrefix, originalError) // Loga o erro original.

	bugTitle := fmt.Sprintf("[BUG][AUTO][PRIORITY_PENDING] %s", bugTitlePrefix)
	var bugDescription string
	if len(bugDescriptionArgs) > 0 {
		format, ok := bugDescriptionArgs[0].(string)
		if ok && len(bugDescriptionArgs) > 1 {
			bugDescription = fmt.Sprintf("[PRIORITY_PENDING] Error encountered: %v. Details: %s", originalError, fmt.Sprintf(format, bugDescriptionArgs[1:]...))
		} else {
			bugDescription = fmt.Sprintf("[PRIORITY_PENDING] Error encountered: %v. No additional details provided for bug task.", originalError)
		}
	} else {
		bugDescription = fmt.Sprintf("[PRIORITY_PENDING] Error encountered: %v.", originalError)
	}

	systemUserID := int64(0) // UserID 0 para tarefas de sistema/bugs.

	// Tenta criar a tarefa de bug.
	// Usa createTaskInternal para evitar recursão de tratamento de erro.
	_, creationErr := s.createTaskInternal(ctx, systemUserID, nil, bugTitle, bugDescription, nil)
	if creationErr != nil {
		logError("CRITICAL: Failed to create bug task for '%s': %v. Original error: %v", bugTitle, creationErr, originalError)
	} else {
		logError("SYSTEM: Bug task created successfully for: %s", bugTitle)
	}
}

// createTaskInternal é uma versão simplificada de CreateTask, usada internamente
// principalmente pela função handleErrorAndCreateBugTask para evitar loops infinitos
// de tratamento de erro caso a criação da tarefa de bug também falhe.
// Este método chama diretamente o repositório sem a lógica de validação de alto nível
// ou a criação automática de tarefas de bug do método CreateTask público.
// Para operações normais de criação de tarefas pelo usuário, o método público CreateTask deve ser usado.
// Retorna a tarefa criada (com ID preenchido) ou um erro se a criação no repositório falhar.
func (s *taskServiceImpl) createTaskInternal(ctx context.Context, userID int64, classID *int64, title, description string, dueDate *time.Time) (models.Task, error) {
	task := models.Task{
		UserID:      userID,
		ClassID:     classID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		IsCompleted: false, // Novas tarefas são sempre não concluídas.
	}

	id, err := s.repo.CreateTask(ctx, &task)
	if err != nil {
		return models.Task{}, fmt.Errorf("createTaskInternal: repositório falhou ao criar tarefa: %w", err)
	}
	task.ID = id // Atribui o ID gerado pelo banco de dados.
	return task, nil
}

// CreateTask é o método público para criar uma nova tarefa.
// Ele realiza validações (ex: título não pode ser vazio) antes de delegar
// a criação para o repositório através de createTaskInternal.
// Se um erro inesperado ocorrer durante a criação no repositório,
// uma tarefa de bug é automaticamente registrada pelo handleErrorAndCreateBugTask.
// O UserID é atualmente fixo (1), mas em uma aplicação real, viria do contexto de autenticação.
// TODO: Integrar com sistema de autenticação para obter UserID real do contexto.
func (s *taskServiceImpl) CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error) {
	if strings.TrimSpace(title) == "" {
		err := errors.New("título da tarefa não pode ser vazio")
		logError("CreateTask: falha de validação: %v", err)
		return models.Task{}, err
	}

	// UserID fixo para demonstração. Deve vir do contexto de autenticação.
	userID := int64(1)

	// Usa createTaskInternal para a lógica de criação real.
	task, err := s.createTaskInternal(ctx, userID, classID, title, description, dueDate)
	if err != nil {
		// Erro inesperado do repositório durante a criação.
		s.handleErrorAndCreateBugTask(ctx, err, "Falha na Criação de Tarefa", "Tentativa de criar tarefa com título '%s'. UserID: %d", title, userID)
		return models.Task{}, fmt.Errorf("CreateTask: falha ao criar tarefa: %w", err)
	}
	return task, nil
}

// GetUpcomingActiveTasks busca tarefas ativas futuras para um usuário.
// Validações básicas nos parâmetros podem ser adicionadas (ex: userID > 0, limit > 0).
// TODO: Adicionar validações para userID e limit, se necessário.
func (s *taskServiceImpl) GetUpcomingActiveTasks(ctx context.Context, userID int64, fromDate time.Time, limit int) ([]models.Task, error) {
	tasks, err := s.repo.GetUpcomingActiveTasks(ctx, userID, fromDate, limit)
	if err != nil {
		// Erros de listagem simples geralmente não criam tarefas de bug,
		// a menos que indiquem um problema sistêmico mais profundo.
		// O erro já vem formatado do repositório.
		logError("GetUpcomingActiveTasks: falha ao buscar tarefas futuras ativas para UserID %d: %v", userID, err)
		return nil, fmt.Errorf("serviço falhou ao buscar tarefas futuras ativas: %w", err)
	}
	return tasks, nil
}

// UpdateTask atualiza uma tarefa existente.
// Valida se o título da tarefa não está vazio.
// O UserID da tarefa deve corresponder ao usuário autenticado (verificação a ser implementada).
// Erros inesperados do repositório disparam a criação de uma tarefa de bug.
// TODO: Adicionar verificação de propriedade da tarefa (UserID).
func (s *taskServiceImpl) UpdateTask(ctx context.Context, task *models.Task) error {
	if strings.TrimSpace(task.Title) == "" {
		err := errors.New("título da tarefa não pode ser vazio para atualização")
		logError("UpdateTask: falha de validação para Tarefa ID %d: %v", task.ID, err)
		return err
	}

	// Idealmente, verificar se a tarefa pertence ao usuário autenticado antes de atualizar.
	// Ex: fetchedTask, _ := s.repo.GetTaskByID(ctx, task.ID); if fetchedTask.UserID != authenticatedUserID { return errors.New("unauthorized") }

	err := s.repo.UpdateTask(ctx, task)
	if err != nil {
		// Se o erro do repositório for "não encontrado" ou "sem alterações", apenas loga.
		// Outros erros (ex: falha de conexão com DB) criam uma tarefa de bug.
		if strings.Contains(err.Error(), "no task found") || strings.Contains(err.Error(), "no values changed") {
			logError("UpdateTask: falha ao atualizar Tarefa ID %d: %v", task.ID, err)
		} else {
			s.handleErrorAndCreateBugTask(ctx, err, "Falha na Atualização de Tarefa", "Tentativa de atualizar Tarefa ID %d, Título '%s'. UserID: %d", task.ID, task.Title, task.UserID)
		}
		return fmt.Errorf("UpdateTask: falha ao atualizar tarefa: %w", err)
	}
	return nil
}

// DeleteTask remove uma tarefa pelo seu ID.
// TODO: Adicionar verificação de propriedade da tarefa (UserID).
// Erros como "não encontrado" são logados, mas outros erros inesperados do repositório
// disparam a criação de uma tarefa de bug.
func (s *taskServiceImpl) DeleteTask(ctx context.Context, taskID int64) error {
	// Idealmente, verificar se a tarefa pertence ao usuário autenticado antes de deletar.
	err := s.repo.DeleteTask(ctx, taskID)
	if err != nil {
		if strings.Contains(err.Error(), "no task found") {
			logError("DeleteTask: falha ao deletar Tarefa ID %d: %v", taskID, err)
		} else {
			s.handleErrorAndCreateBugTask(ctx, err, "Falha na Deleção de Tarefa", "Tentativa de deletar Tarefa ID %d", taskID)
		}
		return fmt.Errorf("DeleteTask: falha ao deletar tarefa: %w", err)
	}
	return nil
}

// ListActiveTasksByClass retorna tarefas ativas (não concluídas) para uma turma específica.
// Erros inesperados do repositório disparam a criação de uma tarefa de bug.
// A filtragem para 'ativas' é feita aqui, mas poderia ser delegada ao repositório.
// TODO: Considerar mover a lógica de filtragem de 'ativas' para o repositório (ex: `repo.GetActiveTasksByClassID`).
func (s *taskServiceImpl) ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error) {
	tasks, err := s.repo.GetTasksByClassID(ctx, classID)
	if err != nil {
		s.handleErrorAndCreateBugTask(ctx, err, "Falha na Listagem de Tarefas por Turma", "Tentativa de listar tarefas para Turma ID %d", classID)
		return nil, fmt.Errorf("ListActiveTasksByClass: falha ao buscar tarefas: %w", err)
	}

	activeTasks := make([]models.Task, 0, len(tasks))
	for _, task := range tasks {
		if !task.IsCompleted {
			activeTasks = append(activeTasks, task)
		}
	}
	return activeTasks, nil
}

// ListAllTasks retorna todas as tarefas (pendentes e concluídas).
// TODO: Em um sistema multiusuário, filtrar pelo UserID do contexto autenticado.
// Erros inesperados do repositório disparam a criação de uma tarefa de bug.
func (s *taskServiceImpl) ListAllTasks(ctx context.Context) ([]models.Task, error) {
	// TODO: Adicionar filtragem por UserID quando a autenticação estiver implementada.
	// userID := ... // obter do contexto
	// tasks, err := s.repo.GetAllTasksByUserID(ctx, userID)
	tasks, err := s.repo.GetAllTasks(ctx) // Versão atual busca todas, globalmente.
	if err != nil {
		s.handleErrorAndCreateBugTask(ctx, err, "Falha na Listagem Global de Tarefas", "Tentativa de listar todas as tarefas")
		return nil, fmt.Errorf("ListAllTasks: falha ao buscar todas as tarefas: %w", err)
	}
	return tasks, nil
}

// ListAllActiveTasks retorna todas as tarefas ativas (não concluídas).
// TODO: Em um sistema multiusuário, filtrar pelo UserID do contexto autenticado.
// Erros inesperados do repositório disparam a criação de uma tarefa de bug.
// A filtragem para 'ativas' é feita aqui.
func (s *taskServiceImpl) ListAllActiveTasks(ctx context.Context) ([]models.Task, error) {
	// TODO: Adicionar filtragem por UserID.
	allTasks, err := s.repo.GetAllTasks(ctx) // Ou um método de repo mais específico se disponível.
	if err != nil {
		s.handleErrorAndCreateBugTask(ctx, err, "Falha na Listagem Global de Tarefas Ativas", "Tentativa de listar todas as tarefas ativas")
		return nil, fmt.Errorf("ListAllActiveTasks: falha ao buscar tarefas ativas: %w", err)
	}

	activeTasks := make([]models.Task, 0, len(allTasks))
	for _, task := range allTasks {
		if !task.IsCompleted {
			activeTasks = append(activeTasks, task)
		}
	}
	return activeTasks, nil
}

// MarkTaskAsCompleted marca uma tarefa como concluída.
// Erros, incluindo "tarefa não encontrada", disparam a criação de uma tarefa de bug,
// pois pode indicar um problema de consistência ou um ID inválido sendo passado.
func (s *taskServiceImpl) MarkTaskAsCompleted(ctx context.Context, taskID int64) error {
	// TODO: Adicionar verificação de propriedade da tarefa (UserID).
	err := s.repo.MarkTaskCompleted(ctx, taskID)
	if err != nil {
		s.handleErrorAndCreateBugTask(ctx, err, "Falha na Conclusão de Tarefa", "Tentativa de completar Tarefa ID %d", taskID)
		return fmt.Errorf("MarkTaskAsCompleted: falha ao marcar tarefa como concluída: %w", err)
	}
	return nil
}

// GetTaskByID recupera uma tarefa pelo seu ID.
// Se a tarefa não for encontrada (sql.ErrNoRows), um erro específico é retornado
// e uma tarefa de bug não é criada para este caso (considerado um erro esperado).
// Outros erros do repositório disparam a criação de uma tarefa de bug.
// TODO: Adicionar verificação de propriedade da tarefa (UserID).
func (s *taskServiceImpl) GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "no task found") {
			logError("GetTaskByID: Tarefa não encontrada com ID %d: %v", taskID, err)
			return nil, fmt.Errorf("tarefa com ID %d não encontrada", taskID) // Retorna erro amigável.
		}
		// Para outros erros inesperados do banco de dados:
		s.handleErrorAndCreateBugTask(ctx, err, "Falha na Recuperação de Tarefa", "Tentativa de recuperar Tarefa ID %d", taskID)
		return nil, fmt.Errorf("GetTaskByID: falha ao buscar tarefa: %w", err)
	}
	return task, nil
}
