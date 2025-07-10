// Package repository contém as implementações concretas das interfaces de repositório
// definidas no pacote pai 'repository'. Este arquivo específico implementa o TaskRepository.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"vigenda/internal/models"
)

// taskRepository é a implementação concreta de TaskRepository.
// Ele usa uma conexão de banco de dados SQL (*sql.DB) para interagir com a tabela 'tasks'.
type taskRepository struct {
	db *sql.DB // db é a conexão com o banco de dados.
}

// NewTaskRepository cria e retorna uma nova instância de TaskRepository.
// Requer uma conexão de banco de dados (*sql.DB) como dependência.
func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

// CreateTask insere uma nova tarefa no banco de dados.
// Retorna o ID da tarefa recém-criada ou um erro.
// Os campos ClassID e DueDate são tratados como opcionais (NULLable no banco de dados).
func (r *taskRepository) CreateTask(ctx context.Context, task *models.Task) (int64, error) {
	query := `INSERT INTO tasks (user_id, class_id, title, description, due_date, is_completed)
              VALUES (?, ?, ?, ?, ?, ?)`

	var classID sql.NullInt64
	if task.ClassID != nil {
		classID.Int64 = *task.ClassID
		classID.Valid = true
	}

	var dueDate sql.NullTime
	if task.DueDate != nil {
		dueDate.Time = *task.DueDate
		dueDate.Valid = true
	}

	result, err := r.db.ExecContext(ctx, query, task.UserID, classID, task.Title, task.Description, dueDate, task.IsCompleted)
	if err != nil {
		return 0, fmt.Errorf("taskRepository.CreateTask: erro ao executar insert: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("taskRepository.CreateTask: erro ao obter ID do último insert: %w", err)
	}
	return id, nil
}

// GetTaskByID busca uma tarefa pelo seu ID.
// Retorna um ponteiro para models.Task ou nil se não encontrada, além de um erro.
func (r *taskRepository) GetTaskByID(ctx context.Context, id int64) (*models.Task, error) {
	query := `SELECT id, user_id, class_id, title, description, due_date, is_completed
              FROM tasks WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	task := &models.Task{}
	var classID sql.NullInt64
	var description sql.NullString
	var dueDate sql.NullTime

	err := row.Scan(
		&task.ID,
		&task.UserID,
		&classID,
		&task.Title,
		&description,
		&dueDate,
		&task.IsCompleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("taskRepository.GetTaskByID: nenhuma tarefa encontrada com ID %d: %w", id, err)
		}
		return nil, fmt.Errorf("taskRepository.GetTaskByID: erro ao escanear linha: %w", err)
	}

	if classID.Valid {
		task.ClassID = &classID.Int64
	}
	if description.Valid {
		task.Description = description.String
	}
	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	}
	return task, nil
}

// GetTasksByClassID busca todas as tarefas associadas a um ClassID específico.
// Retorna uma slice de models.Task ou um erro.
func (r *taskRepository) GetTasksByClassID(ctx context.Context, classID int64) ([]models.Task, error) {
	query := `SELECT id, user_id, class_id, title, description, due_date, is_completed
              FROM tasks WHERE class_id = ?`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.GetTasksByClassID: erro ao consultar tarefas por classID: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		task := models.Task{}
		var cID sql.NullInt64 // Usar nome diferente para evitar sombreamento
		var description sql.NullString
		var dueDate sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&cID,
			&task.Title,
			&description,
			&dueDate,
			&task.IsCompleted,
		)
		if err != nil {
			return nil, fmt.Errorf("taskRepository.GetTasksByClassID: erro ao escanear tarefa: %w", err)
		}
		if cID.Valid {
			task.ClassID = &cID.Int64
		}
		if description.Valid {
			task.Description = description.String
		}
		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("taskRepository.GetTasksByClassID: erro ao iterar linhas: %w", err)
	}
	return tasks, nil
}

// GetAllTasks busca todas as tarefas do banco de dados.
// Em uma aplicação real, isso provavelmente seria paginado ou filtrado por usuário.
// Retorna uma slice de models.Task ou um erro.
func (r *taskRepository) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	query := `SELECT id, user_id, class_id, title, description, due_date, is_completed FROM tasks`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.GetAllTasks: erro ao consultar todas as tarefas: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		task := models.Task{}
		var classID sql.NullInt64
		var description sql.NullString
		var dueDate sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&classID,
			&task.Title,
			&description,
			&dueDate,
			&task.IsCompleted,
		)
		if err != nil {
			return nil, fmt.Errorf("taskRepository.GetAllTasks: erro ao escanear tarefa: %w", err)
		}
		if classID.Valid {
			task.ClassID = &classID.Int64
		}
		if description.Valid {
			task.Description = description.String
		}
		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("taskRepository.GetAllTasks: erro ao iterar linhas: %w", err)
	}
	return tasks, nil
}

// MarkTaskCompleted atualiza o status de uma tarefa para concluída (is_completed = true).
// Retorna um erro se a tarefa não for encontrada ou se houver um problema na atualização.
func (r *taskRepository) MarkTaskCompleted(ctx context.Context, taskID int64) error {
	query := `UPDATE tasks SET is_completed = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, true, taskID)
	if err != nil {
		return fmt.Errorf("taskRepository.MarkTaskCompleted: erro ao executar update: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taskRepository.MarkTaskCompleted: erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("taskRepository.MarkTaskCompleted: nenhuma tarefa encontrada com ID %d ou tarefa já estava concluída", taskID)
	}
	return nil
}

// GetUpcomingActiveTasks busca tarefas ativas (não concluídas) de um usuário específico
// com data de vencimento a partir de 'fromDate', ordenadas pela data de vencimento e limitadas por 'limit'.
// Retorna uma slice de models.Task ou um erro.
func (r *taskRepository) GetUpcomingActiveTasks(ctx context.Context, userID int64, fromDate time.Time, limit int) ([]models.Task, error) {
	query := `
		SELECT id, user_id, class_id, title, description, due_date, is_completed
		FROM tasks
		WHERE user_id = ?
		  AND is_completed = false
		  AND date(due_date) >= date(?)
		ORDER BY due_date ASC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, userID, fromDate, limit)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.GetUpcomingActiveTasks: erro ao consultar tarefas: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		task := models.Task{}
		var classID sql.NullInt64
		var description sql.NullString
		var dueDate sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&classID,
			&task.Title,
			&description,
			&dueDate,
			&task.IsCompleted,
		)
		if err != nil {
			return nil, fmt.Errorf("taskRepository.GetUpcomingActiveTasks: erro ao escanear tarefa: %w", err)
		}
		if classID.Valid {
			task.ClassID = &classID.Int64
		}
		if description.Valid {
			task.Description = description.String
		}
		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("taskRepository.GetUpcomingActiveTasks: erro ao iterar linhas: %w", err)
	}
	return tasks, nil
}

// DeleteTask remove uma tarefa do banco de dados pelo seu ID.
// Retorna um erro se a tarefa não for encontrada ou se houver um problema na exclusão.
func (r *taskRepository) DeleteTask(ctx context.Context, taskID int64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("taskRepository.DeleteTask: erro ao executar delete: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taskRepository.DeleteTask: erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("taskRepository.DeleteTask: nenhuma tarefa encontrada com ID %d", taskID)
	}
	return nil
}

// UpdateTask atualiza todos os campos de uma tarefa existente no banco de dados.
// Retorna um erro se a tarefa não for encontrada ou se houver um problema na atualização.
func (r *taskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `UPDATE tasks SET user_id = ?, class_id = ?, title = ?, description = ?, due_date = ?, is_completed = ?
              WHERE id = ?`

	var classID sql.NullInt64
	if task.ClassID != nil {
		classID.Int64 = *task.ClassID
		classID.Valid = true
	} else {
		classID.Valid = false // Garante que será NULL se task.ClassID for nil
	}

	var dueDate sql.NullTime
	if task.DueDate != nil {
		dueDate.Time = *task.DueDate
		dueDate.Valid = true
	} else {
		dueDate.Valid = false // Garante que será NULL se task.DueDate for nil
	}

	result, err := r.db.ExecContext(ctx, query, task.UserID, classID, task.Title, task.Description, dueDate, task.IsCompleted, task.ID)
	if err != nil {
		return fmt.Errorf("taskRepository.UpdateTask: erro ao executar update: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taskRepository.UpdateTask: erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		// Pode ser que a tarefa não exista ou que os valores fornecidos para atualização
		// sejam idênticos aos valores já existentes no banco, resultando em 0 linhas afetadas.
		// Considerar retornar um erro específico para "not found" se for um requisito.
		return fmt.Errorf("taskRepository.UpdateTask: nenhuma tarefa encontrada com ID %d, ou nenhum valor foi alterado", task.ID)
	}
	return nil
}
