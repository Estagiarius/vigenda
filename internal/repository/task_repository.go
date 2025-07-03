package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time" // Added import for time
	"vigenda/internal/models"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *models.Task) (int64, error) {
	// Removed created_at, updated_at from INSERT as they are not in the tasks schema
	query := `INSERT INTO tasks (user_id, class_id, title, description, due_date, is_completed)
              VALUES (?, ?, ?, ?, ?, ?)`
	// now := time.Now() // Not used

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
		return 0, fmt.Errorf("taskRepository.CreateTask: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("taskRepository.CreateTask: failed to get last insert ID: %w", err)
	}
	return id, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int64) (*models.Task, error) {
	// Removed created_at, updated_at from SELECT
	query := `SELECT id, user_id, class_id, title, description, due_date, is_completed
              FROM tasks WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	task := &models.Task{}
	var classID sql.NullInt64
	var description sql.NullString
	var dueDate sql.NullTime

	// Removed task.CreatedAt, task.UpdatedAt from Scan
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
	// Removed created_at, updated_at from SELECT
	query := `SELECT id, user_id, class_id, title, description, due_date, is_completed
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

		// Removed task.CreatedAt, task.UpdatedAt from Scan
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
	// Removed created_at, updated_at from SELECT
	query := `SELECT id, user_id, class_id, title, description, due_date, is_completed FROM tasks`
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

		// Removed task.CreatedAt, task.UpdatedAt from Scan
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
	// Removed updated_at from UPDATE as it's not in the tasks schema
	query := `UPDATE tasks SET is_completed = ? WHERE id = ?`
	// now := time.Now() // Not used
	result, err := r.db.ExecContext(ctx, query, true, taskID)
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

// GetUpcomingTasksByUserID retrieves a list of tasks for a given user that are not completed,
// have a due date in the future, ordered by the due date, and limited by the specified limit.
func (r *taskRepository) GetUpcomingTasksByUserID(ctx context.Context, userID int64, limit int) ([]models.Task, error) {
	query := `
        SELECT id, user_id, class_id, title, description, due_date, is_completed
        FROM tasks
        WHERE user_id = ? AND is_completed = FALSE AND due_date IS NOT NULL AND due_date >= ?
        ORDER BY due_date ASC
        LIMIT ?`

	// Using time.Now() in UTC to compare with stored dates, assuming dates are stored in UTC or consistently.
	// For local time considerations, this might need adjustment based on how dates are handled application-wide.
	// SQLite typically stores DATETIME as TEXT in UTC if not specified otherwise.
	// Using time.Now().Format("2006-01-02 15:04:05") for SQLite compatibility if it expects text.
	// However, mattn/go-sqlite3 driver handles time.Time correctly, so direct time.Time should be fine.
	now := time.Now() // Using local time, assuming DB stores it consistently or driver handles conversion.
					  // If explicit UTC is needed: time.Now().UTC()

	rows, err := r.db.QueryContext(ctx, query, userID, now, limit)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.GetUpcomingTasksByUserID: query failed: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		task := models.Task{}
		var classID sql.NullInt64
		var description sql.NullString
		var dueDate sql.NullTime // due_date is NOT NULL in the query

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
			return nil, fmt.Errorf("taskRepository.GetUpcomingTasksByUserID: scanning task: %w", err)
		}
		if classID.Valid {
			task.ClassID = &classID.Int64
		}
		if description.Valid {
			task.Description = description.String
		}
		if dueDate.Valid { // Should always be valid due to query `due_date IS NOT NULL`
			task.DueDate = &dueDate.Time
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("taskRepository.GetUpcomingTasksByUserID: iterating rows: %w", err)
	}
	return tasks, nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, taskID int64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("taskRepository.DeleteTask: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taskRepository.DeleteTask: checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("taskRepository.DeleteTask: no task found with ID %d", taskID)
	}
	return nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `UPDATE tasks SET user_id = ?, class_id = ?, title = ?, description = ?, due_date = ?, is_completed = ?
              WHERE id = ?`

	var classID sql.NullInt64
	if task.ClassID != nil {
		classID.Int64 = *task.ClassID
		classID.Valid = true
	} else {
		classID.Valid = false // Ensure it's explicitly NULL if task.ClassID is nil
	}

	var dueDate sql.NullTime
	if task.DueDate != nil {
		dueDate.Time = *task.DueDate
		dueDate.Valid = true
	} else {
		dueDate.Valid = false // Ensure it's explicitly NULL if task.DueDate is nil
	}

	result, err := r.db.ExecContext(ctx, query, task.UserID, classID, task.Title, task.Description, dueDate, task.IsCompleted, task.ID)
	if err != nil {
		return fmt.Errorf("taskRepository.UpdateTask: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taskRepository.UpdateTask: checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("taskRepository.UpdateTask: no task found with ID %d, or no values changed", task.ID)
	}
	return nil
}
