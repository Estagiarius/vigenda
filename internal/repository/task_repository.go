package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"vigenda/internal/models"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *models.Task) (int64, error) {
	query := `INSERT INTO tasks (user_id, class_id, title, description, due_date, is_completed, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()

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

	result, err := r.db.ExecContext(ctx, query, task.UserID, classID, task.Title, task.Description, dueDate, task.IsCompleted, now, now)
	if err != nil {
		return 0, fmt.Errorf("taskRepository.CreateTask: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("taskRepository.CreateTask: failed to get last insert ID: %w", err)
	}
	return id, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int64) (*models.Task, error) {
	query := `SELECT id, user_id, class_id, title, description, due_date, is_completed, created_at, updated_at
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
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("taskRepository.GetTaskByID: no task found with ID %d", id)
		}
		return nil, fmt.Errorf("taskRepository.GetTaskByID: %w", err)
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

func (r *taskRepository) GetTasksByClassID(ctx context.Context, classID int64) ([]models.Task, error) {
	query := `SELECT id, user_id, class_id, title, description, due_date, is_completed, created_at, updated_at
              FROM tasks WHERE class_id = ?`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.GetTasksByClassID: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		task := models.Task{}
		var cID sql.NullInt64 // Renamed to avoid conflict with parameter
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
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("taskRepository.GetTasksByClassID: scanning task: %w", err)
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
		return nil, fmt.Errorf("taskRepository.GetTasksByClassID: iterating rows: %w", err)
	}
	return tasks, nil
}

func (r *taskRepository) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	query := `SELECT id, user_id, class_id, title, description, due_date, is_completed, created_at, updated_at FROM tasks`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.GetAllTasks: %w", err)
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
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("taskRepository.GetAllTasks: scanning task: %w", err)
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
		return nil, fmt.Errorf("taskRepository.GetAllTasks: iterating rows: %w", err)
	}
	return tasks, nil
}

func (r *taskRepository) MarkTaskCompleted(ctx context.Context, taskID int64) error {
	query := `UPDATE tasks SET is_completed = ?, updated_at = ? WHERE id = ?`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, true, now, taskID)
	if err != nil {
		return fmt.Errorf("taskRepository.MarkTaskCompleted: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taskRepository.MarkTaskCompleted: checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("taskRepository.MarkTaskCompleted: no task found with ID %d or task already completed", taskID)
	}
	return nil
}
