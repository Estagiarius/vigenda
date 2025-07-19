package service

import (
	"context"
	"fmt"
	"vigenda/internal/models"
	"vigenda/internal/repository"
)

// subjectServiceImpl é a implementação concreta da interface SubjectService.
// Ela encapsula a lógica de negócios para o gerenciamento de disciplinas,
// dependendo de um repositório para a persistência dos dados.
type subjectServiceImpl struct {
	repo repository.SubjectRepository
	// Adicionar outros repositórios aqui se necessário para lógicas complexas,
	// como verificar se uma disciplina tem turmas antes de deletar.
}

// NewSubjectService cria uma nova instância de SubjectService.
// Recebe um SubjectRepository como dependência para interagir com a camada de dados.
func NewSubjectService(repo repository.SubjectRepository) SubjectService {
	return &subjectServiceImpl{
		repo: repo,
	}
}

// CreateSubject implementa a lógica para criar uma nova disciplina.
func (s *subjectServiceImpl) CreateSubject(ctx context.Context, userID int64, name string) (models.Subject, error) {
	// Validação básica
	if name == "" {
		return models.Subject{}, fmt.Errorf("O nome da disciplina não pode ser vazio")
	}

	subject := models.Subject{
		UserID: userID,
		Name:   name,
	}

	// Chama o repositório para persistir a disciplina
	err := s.repo.Create(ctx, &subject)
	if err != nil {
		return models.Subject{}, fmt.Errorf("Falha ao criar disciplina no repositório: %w", err)
	}

	return subject, nil
}

// GetSubjectByID implementa a lógica para buscar uma disciplina pelo ID.
func (s *subjectServiceImpl) GetSubjectByID(ctx context.Context, subjectID int64) (models.Subject, error) {
	subject, err := s.repo.GetByID(ctx, subjectID)
	if err != nil {
		return models.Subject{}, fmt.Errorf("Falha ao buscar disciplina: %w", err)
	}
	return *subject, nil
}

// ListSubjectsByUser implementa a lógica para listar as disciplinas de um usuário.
func (s *subjectServiceImpl) ListSubjectsByUser(ctx context.Context, userID int64) ([]models.Subject, error) {
	subjects, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("Falha ao listar disciplinas: %w", err)
	}
	return subjects, nil
}

// UpdateSubject implementa a lógica para atualizar uma disciplina.
func (s *subjectServiceImpl) UpdateSubject(ctx context.Context, subjectID int64, name string) (models.Subject, error) {
	// Validação
	if name == "" {
		return models.Subject{}, fmt.Errorf("O nome da disciplina não pode ser vazio")
	}

	// Busca a disciplina para garantir que ela existe antes de atualizar
	subject, err := s.repo.GetByID(ctx, subjectID)
	if err != nil {
		return models.Subject{}, fmt.Errorf("Disciplina não encontrada para atualização: %w", err)
	}

	subject.Name = name

	err = s.repo.Update(ctx, subject)
	if err != nil {
		return models.Subject{}, fmt.Errorf("Falha ao atualizar disciplina: %w", err)
	}

	return *subject, nil
}

// DeleteSubject implementa a lógica para deletar uma disciplina.
// TODO: Adicionar lógica para verificar se a disciplina tem turmas associadas antes de deletar.
func (s *subjectServiceImpl) DeleteSubject(ctx context.Context, subjectID int64) error {
	err := s.repo.Delete(ctx, subjectID)
	if err != nil {
		return fmt.Errorf("Falha ao deletar disciplina: %w", err)
	}
	return nil
}
