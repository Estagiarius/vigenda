package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"vigenda/internal/models"
)

type lessonRepositoryImpl struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) LessonRepository {
	return &lessonRepositoryImpl{db: db}
}

func (r *lessonRepositoryImpl) CreateLesson(ctx context.Context, lesson *models.Lesson) (int64, error) {
	// TODO: Adicionar UserID à tabela lessons se as lições forem por usuário.
	// Atualmente, models.Lesson não tem UserID, mas ClassID sim, que tem UserID.
	// A propriedade pode ser inferida pela Class.
	query := `INSERT INTO lessons (class_id, title, plan_content, scheduled_at, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, lesson.ClassID, lesson.Title, lesson.PlanContent, lesson.ScheduledAt, now, now)
	if err != nil {
		return 0, fmt.Errorf("lessonRepository.CreateLesson: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("lessonRepository.CreateLesson: error getting last insert ID: %w", err)
	}
	return id, nil
}

func (r *lessonRepositoryImpl) GetLessonByID(ctx context.Context, lessonID int64) (*models.Lesson, error) {
	query := `SELECT id, class_id, title, plan_content, scheduled_at, created_at, updated_at
              FROM lessons WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, lessonID)
	lesson := &models.Lesson{}
	// Adicionar CreatedAt e UpdatedAt ao Scan
	err := row.Scan(&lesson.ID, &lesson.ClassID, &lesson.Title, &lesson.PlanContent, &lesson.ScheduledAt, &lesson.CreatedAt, &lesson.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lesson with ID %d not found", lessonID)
		}
		return nil, fmt.Errorf("lessonRepository.GetLessonByID: %w", err)
	}
	return lesson, nil
}

func (r *lessonRepositoryImpl) GetLessonsByClassID(ctx context.Context, classID int64) ([]models.Lesson, error) {
	query := `SELECT id, class_id, title, plan_content, scheduled_at, created_at, updated_at
              FROM lessons WHERE class_id = ? ORDER BY scheduled_at ASC`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, fmt.Errorf("lessonRepository.GetLessonsByClassID: %w", err)
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		lesson := models.Lesson{}
		// Adicionar CreatedAt e UpdatedAt ao Scan
		if err := rows.Scan(&lesson.ID, &lesson.ClassID, &lesson.Title, &lesson.PlanContent, &lesson.ScheduledAt, &lesson.CreatedAt, &lesson.UpdatedAt); err != nil {
			return nil, fmt.Errorf("lessonRepository.GetLessonsByClassID: scanning row: %w", err)
		}
		lessons = append(lessons, lesson)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("lessonRepository.GetLessonsByClassID: iterating rows: %w", err)
	}
	return lessons, nil
}

// GetLessonsByDateRange busca lições para um usuário dentro de um intervalo de datas.
// O userID é usado para filtrar lições associadas a turmas desse usuário.
func (r *lessonRepositoryImpl) GetLessonsByDateRange(ctx context.Context, userID int64, startDate time.Time, endDate time.Time) ([]models.Lesson, error) {
	// Esta query assume que a tabela 'classes' tem uma coluna 'user_id'.
	// E que 'lessons' está ligada a 'classes' por 'class_id'.
	// O UserID em models.Lesson não existe diretamente, então filtramos via Class.UserID.
	// Se UserID for 0, podemos interpretar como "buscar para todos os usuários" ou ajustar a query.
	// Por agora, se UserID > 0, filtramos.
	var query string
	var args []interface{}

	// Garantir que endDate seja o fim do dia para incluir todas as lições da data final.
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())


	if userID > 0 {
		query = `SELECT l.id, l.class_id, l.title, l.plan_content, l.scheduled_at, l.created_at, l.updated_at
                 FROM lessons l
                 JOIN classes c ON l.class_id = c.id
                 WHERE c.user_id = ? AND l.scheduled_at >= ? AND l.scheduled_at <= ?
                 ORDER BY l.scheduled_at ASC`
		args = append(args, userID, startDate, endDate)
	} else { // Se userID for 0 ou negativo, busca para todas as turmas (comportamento de admin/sistema)
		query = `SELECT id, class_id, title, plan_content, scheduled_at, created_at, updated_at
                 FROM lessons
                 WHERE scheduled_at >= ? AND scheduled_at <= ?
                 ORDER BY scheduled_at ASC`
		args = append(args, startDate, endDate)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("lessonRepository.GetLessonsByDateRange: %w", err)
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		lesson := models.Lesson{}
		if err := rows.Scan(&lesson.ID, &lesson.ClassID, &lesson.Title, &lesson.PlanContent, &lesson.ScheduledAt, &lesson.CreatedAt, &lesson.UpdatedAt); err != nil {
			return nil, fmt.Errorf("lessonRepository.GetLessonsByDateRange: scanning row: %w", err)
		}
		lessons = append(lessons, lesson)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("lessonRepository.GetLessonsByDateRange: iterating rows: %w", err)
	}
	return lessons, nil
}

func (r *lessonRepositoryImpl) UpdateLesson(ctx context.Context, lesson *models.Lesson) error {
	query := `UPDATE lessons SET class_id = ?, title = ?, plan_content = ?, scheduled_at = ?, updated_at = ?
              WHERE id = ?`
	// Idealmente, verificaríamos também o UserID da turma associada para garantir propriedade,
	// mas isso é mais lógico na camada de serviço.
	_, err := r.db.ExecContext(ctx, query, lesson.ClassID, lesson.Title, lesson.PlanContent, lesson.ScheduledAt, time.Now(), lesson.ID)
	if err != nil {
		return fmt.Errorf("lessonRepository.UpdateLesson: %w", err)
	}
	return nil
}

func (r *lessonRepositoryImpl) DeleteLesson(ctx context.Context, lessonID int64) error {
	query := `DELETE FROM lessons WHERE id = ?`
	// Similar ao Update, a verificação de propriedade (via UserID da Class)
	// seria melhor na camada de serviço antes de chamar o delete.
	_, err := r.db.ExecContext(ctx, query, lessonID)
	if err != nil {
		return fmt.Errorf("lessonRepository.DeleteLesson: %w", err)
	}
	return nil
}
