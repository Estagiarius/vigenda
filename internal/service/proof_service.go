// Package service contém as interfaces e implementações da lógica de negócio da aplicação.
package service

import (
	"context"
	"fmt"
	"vigenda/internal/models"
	"vigenda/internal/repository"
)

type proofServiceImpl struct {
	questionRepo repository.QuestionRepository
	// Outras dependências, como SubjectRepository, podem ser adicionadas aqui.
}

// NewProofService cria uma nova instância de ProofService.
// É necessário passar um QuestionRepository para interagir com o banco de dados de questões.
func NewProofService(qr repository.QuestionRepository) ProofService {
	return &proofServiceImpl{
		questionRepo: qr,
	}
}

// GenerateProof gera uma prova (lista de questões) com base nos critérios fornecidos.
// Ele busca questões do repositório que correspondam aos critérios de disciplina,
// tópico (opcional) e distribui a quantidade de questões por dificuldade.
func (s *proofServiceImpl) GenerateProof(ctx context.Context, criteria ProofCriteria) ([]models.Question, error) {
	var proofQuestions []models.Question

	if criteria.EasyCount == 0 && criteria.MediumCount == 0 && criteria.HardCount == 0 {
		return nil, fmt.Errorf("pelo menos uma contagem de dificuldade deve ser maior que zero")
	}

	// Buscar questões fáceis
	if criteria.EasyCount > 0 {
		easyQuestions, err := s.questionRepo.GetQuestionsByCriteria(ctx, repository.QuestionQueryCriteria{
			SubjectID:  criteria.SubjectID,
			Topic:      criteria.Topic,
			Difficulty: "facil",
			Limit:      criteria.EasyCount,
		})
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar questões fáceis: %w", err)
		}
		if len(easyQuestions) < criteria.EasyCount {
			return nil, fmt.Errorf("não há questões fáceis suficientes para o critério (solicitado: %d, disponível: %d)", criteria.EasyCount, len(easyQuestions))
		}
		proofQuestions = append(proofQuestions, easyQuestions...)
	}

	// Buscar questões médias
	if criteria.MediumCount > 0 {
		mediumQuestions, err := s.questionRepo.GetQuestionsByCriteria(ctx, repository.QuestionQueryCriteria{
			SubjectID:  criteria.SubjectID,
			Topic:      criteria.Topic,
			Difficulty: "media",
			Limit:      criteria.MediumCount,
		})
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar questões médias: %w", err)
		}
		if len(mediumQuestions) < criteria.MediumCount {
			return nil, fmt.Errorf("não há questões médias suficientes para o critério (solicitado: %d, disponível: %d)", criteria.MediumCount, len(mediumQuestions))
		}
		proofQuestions = append(proofQuestions, mediumQuestions...)
	}

	// Buscar questões difíceis
	if criteria.HardCount > 0 {
		hardQuestions, err := s.questionRepo.GetQuestionsByCriteria(ctx, repository.QuestionQueryCriteria{
			SubjectID:  criteria.SubjectID,
			Topic:      criteria.Topic,
			Difficulty: "dificil",
			Limit:      criteria.HardCount,
		})
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar questões difíceis: %w", err)
		}
		if len(hardQuestions) < criteria.HardCount {
			return nil, fmt.Errorf("não há questões difíceis suficientes para o critério (solicitado: %d, disponível: %d)", criteria.HardCount, len(hardQuestions))
		}
		proofQuestions = append(proofQuestions, hardQuestions...)
	}

	// TODO: Adicionar lógica para randomizar a ordem das questões, se necessário.
	// TODO: Considerar o que fazer se houver menos questões disponíveis do que o solicitado para uma dificuldade específica.
	// A implementação atual retorna erro, o que é uma abordagem. Outra seria retornar o máximo possível com um aviso.

	return proofQuestions, nil
}
