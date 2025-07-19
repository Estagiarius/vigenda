package repository

import (
	"context"
	"database/sql"
	"fmt"
	"vigenda/internal/models"
)

// subjectRepository é a implementação da interface SubjectRepository para um banco de dados SQL.
type subjectRepository struct {
	db *sql.DB
}

// NewSubjectRepository cria uma nova instância de SubjectRepository.
func NewSubjectRepository(db *sql.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

// Create insere uma nova disciplina no banco de dados.
func (r *subjectRepository) Create(ctx context.Context, subject *models.Subject) error {
	query := `INSERT INTO subjects (user_id, name) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, subject.UserID, subject.Name)
	if err != nil {
		return fmt.Errorf("falha ao inserir disciplina: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("falha ao obter o ID da disciplina inserida: %w", err)
	}
	subject.ID = id
	return nil
}

// GetByID busca uma disciplina pelo seu ID.
func (r *subjectRepository) GetByID(ctx context.Context, id int64) (*models.Subject, error) {
	query := `SELECT id, user_id, name FROM subjects WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	var subject models.Subject
	err := row.Scan(&subject.ID, &subject.UserID, &subject.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("disciplina com ID %d não encontrada", id)
		}
		return nil, fmt.Errorf("falha ao buscar disciplina por ID: %w", err)
	}
	return &subject, nil
}

// GetByUserID busca todas as disciplinas associadas a um ID de usuário.
func (r *subjectRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Subject, error) {
	query := `SELECT id, user_id, name FROM subjects WHERE user_id = ?`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar disciplinas por usuário: %w", err)
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.ID, &subject.UserID, &subject.Name); err != nil {
			return nil, fmt.Errorf("falha ao escanear disciplina: %w", err)
		}
		subjects = append(subjects, subject)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro nas linhas do resultado da busca de disciplinas: %w", err)
	}

	return subjects, nil
}

// Update atualiza os dados de uma disciplina existente no banco de dados.
func (r *subjectRepository) Update(ctx context.Context, subject *models.Subject) error {
	query := `UPDATE subjects SET name = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, subject.Name, subject.ID)
	if err != nil {
		return fmt.Errorf("falha ao atualizar disciplina: %w", err)
	}
	return nil
}

// Delete remove uma disciplina do banco de dados pelo seu ID.
func (r *subjectRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM subjects WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("falha ao deletar disciplina: %w", err)
	}
	return nil
}

// GetOrCreateByNameAndUser busca uma disciplina pelo nome e ID do usuário.
// Se não existir, cria uma nova.
func (r *subjectRepository) GetOrCreateByNameAndUser(ctx context.Context, name string, userID int64) (models.Subject, error) {
	// Tenta buscar primeiro
	query := `SELECT id, user_id, name FROM subjects WHERE name = ? AND user_id = ?`
	row := r.db.QueryRowContext(ctx, query, name, userID)

	var subject models.Subject
	err := row.Scan(&subject.ID, &subject.UserID, &subject.Name)

	if err == nil {
		// Encontrou, retorna
		return subject, nil
	}

	if err != sql.ErrNoRows {
		// Erro inesperado
		return models.Subject{}, fmt.Errorf("falha ao buscar disciplina por nome e usuário: %w", err)
	}

	// Não encontrou (sql.ErrNoRows), então cria
	newSubject := models.Subject{
		UserID: userID,
		Name:   name,
	}

	createQuery := `INSERT INTO subjects (user_id, name) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, createQuery, newSubject.UserID, newSubject.Name)
	if err != nil {
		return models.Subject{}, fmt.Errorf("falha ao criar nova disciplina em GetOrCreate: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Subject{}, fmt.Errorf("falha ao obter ID da nova disciplina em GetOrCreate: %w", err)
	}
	newSubject.ID = id

	return newSubject, nil
}
