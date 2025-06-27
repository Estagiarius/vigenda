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
	if criteria.SubjectID == 0 {
		return nil, fmt.Errorf("SubjectID não pode ser zero")
	}
	if criteria.EasyCount < 0 || criteria.MediumCount < 0 || criteria.HardCount < 0 {
		return nil, fmt.Errorf("contagem de questões não pode ser negativa")
	}
	if criteria.EasyCount == 0 && criteria.MediumCount == 0 && criteria.HardCount == 0 {
		return nil, fmt.Errorf("pelo menos uma contagem de dificuldade deve ser maior que zero")
	}

	// Converter service.ProofCriteria para repository.ProofCriteria
	repoCriteria := repository.ProofCriteria{
		SubjectID:   criteria.SubjectID,
		Topic:       criteria.Topic,
		EasyCount:   criteria.EasyCount,
		MediumCount: criteria.MediumCount,
		HardCount:   criteria.HardCount,
	}

	// Use a nova função GetQuestionsByCriteriaProofGeneration do repositório
	// que é projetada para buscar o número exato de questões por dificuldade em uma única chamada ou transação.
	allQuestions, err := s.questionRepo.GetQuestionsByCriteriaProofGeneration(ctx, repoCriteria)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar questões para a prova: %w", err)
	}

	// Verificar se o número de questões retornadas corresponde ao solicitado.
	// A responsabilidade de retornar erro se não houver questões suficientes pode ser do repositório ou do serviço.
	// Se GetQuestionsByCriteriaProofGeneration já garante a contagem ou erro, esta verificação pode ser simplificada.
	// Assumindo que GetQuestionsByCriteriaProofGeneration tenta obter o máximo possível e não erro se menos.

	// Validar se o número total de questões corresponde ao solicitado
	// Esta é uma verificação simples; uma mais robusta verificaria cada dificuldade.
	expectedTotal := criteria.EasyCount + criteria.MediumCount + criteria.HardCount
	if len(allQuestions) < expectedTotal {
		// Isso indica que não foram encontradas questões suficientes para todos os critérios.
		// O comportamento aqui pode variar: retornar erro, ou retornar o que foi encontrado com um aviso.
		// Para ser rigoroso, vamos retornar um erro se não pudermos atender ao pedido completo.
		// No entanto, a lógica no repositório GetQuestionsByCriteriaProofGeneration já itera
		// por dificuldade e busca o limite. Se uma busca por dificuldade não encontra o suficiente,
		// o resultado final terá menos questões.
		// O repositório não retorna erro por encontrar menos, apenas o serviço pode decidir isso.

		// Contar quantas questões de cada dificuldade foram realmente retornadas.
		counts := make(map[string]int)
		for _, q := range allQuestions {
			counts[q.Difficulty]++
		}

		if counts["facil"] < criteria.EasyCount {
			return nil, fmt.Errorf("não há questões fáceis suficientes (solicitado: %d, disponível: %d)", criteria.EasyCount, counts["facil"])
		}
		if counts["media"] < criteria.MediumCount {
			return nil, fmt.Errorf("não há questões médias suficientes (solicitado: %d, disponível: %d)", criteria.MediumCount, counts["media"])
		}
		if counts["dificil"] < criteria.HardCount {
			return nil, fmt.Errorf("não há questões difíceis suficientes (solicitado: %d, disponível: %d)", criteria.HardCount, counts["dificil"])
		}
	}

	// A ordem das questões já é randomizada pelo repositório (ORDER BY RANDOM()).
	// Se uma ordem específica (ex: todas fáceis, depois médias, depois difíceis) ou uma nova randomização
	// do conjunto total for necessária, ela seria implementada aqui.
	// Por agora, a ordem do repositório é mantida.

	return allQuestions, nil
}
