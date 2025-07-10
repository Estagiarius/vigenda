// Conteúdo original do arquivo internal/service/question_service.go
// Se houver algum conteúdo pré-existente, ele deve ser mantido ou mesclado.
// Por enquanto, assumindo que o arquivo está vazio ou contém apenas o cabeçalho do pacote.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"vigenda/internal/models"
	"vigenda/internal/repository"
)

// questionServiceImpl é a implementação de QuestionService.
type questionServiceImpl struct {
	questionRepo repository.QuestionRepository
	subjectRepo  repository.SubjectRepository // Para validar SubjectID, se necessário
}

// NewQuestionService cria uma nova instância de QuestionService.
// O SubjectRepository é opcional e pode ser nil se a validação de SubjectID não for feita aqui.
func NewQuestionService(qr repository.QuestionRepository, sr repository.SubjectRepository) QuestionService {
	return &questionServiceImpl{
		questionRepo: qr,
		subjectRepo:  sr,
	}
}

// AddQuestionsFromJSON processa um payload JSON de questões e as adiciona ao banco de dados.
func (s *questionServiceImpl) AddQuestionsFromJSON(ctx context.Context, jsonData []byte) (int, error) {
	var questions []struct {
		SubjectName   string `json:"disciplina"` // Campo para nome da disciplina no JSON
		Topic         string `json:"topico"`
		Type          string `json:"tipo"`
		Difficulty    string `json:"dificuldade"`
		Statement     string `json:"enunciado"`
		Options       any    `json:"opcoes"` // Pode ser []string ou nil
		CorrectAnswer string `json:"resposta_correta"`
		UserID        int64  `json:"user_id"` // Assumindo que o UserID virá no JSON ou será obtido do contexto
	}

	if err := json.Unmarshal(jsonData, &questions); err != nil {
		return 0, fmt.Errorf("falha ao decodificar JSON: %w", err)
	}

	if len(questions) == 0 {
		return 0, fmt.Errorf("nenhuma questão fornecida no JSON")
	}

	addedCount := 0
	var questionsToAdd []models.Question

	for i, qJSON := range questions {
		// Validação básica dos campos da questão do JSON
		if qJSON.SubjectName == "" {
			return addedCount, fmt.Errorf("questão %d: 'disciplina' é obrigatório no JSON", i)
		}
		if qJSON.Statement == "" {
			return addedCount, fmt.Errorf("questão %d: 'enunciado' é obrigatório no JSON", i)
		}
		if qJSON.Difficulty == "" {
			return addedCount, fmt.Errorf("questão %d: 'dificuldade' é obrigatório no JSON", i)
		}
		if qJSON.Type == "" {
			return addedCount, fmt.Errorf("questão %d: 'tipo' é obrigatório no JSON", i)
		}
		if qJSON.CorrectAnswer == "" {
			return addedCount, fmt.Errorf("questão %d: 'resposta_correta' é obrigatório no JSON", i)
		}

		// TODO: Obter SubjectID a partir de SubjectName e UserID.
		// Isso requer que SubjectRepository tenha um método como GetByNameAndUser.
		// Por simplicidade, vamos assumir um SubjectID fixo ou que ele seja passado de outra forma.
		// Numa implementação real:
		// subject, err := s.subjectRepo.GetByNameAndUser(ctx, qJSON.SubjectName, qJSON.UserID)
		// if err != nil {
		//    return addedCount, fmt.Errorf("questão %d: disciplina '%s' não encontrada para o usuário: %w", i, qJSON.SubjectName, err)
		// }
		// currentSubjectID := subject.ID
		// Por enquanto, usaremos um valor placeholder ou assumiremos que SubjectID está no JSON.
		// Como o modelo `models.Question` espera `SubjectID` e não `SubjectName`,
		// esta etapa é crucial. Para o propósito desta tarefa, vamos assumir que o JSON
		// deveria ter `subject_id` ou que o `QuestionRepository` lida com a conversão/criação.
		// Para manter a conformidade com `models.Question` e o schema, vamos assumir que `SubjectID`
		// é o que deve ser usado e que o JSON deve fornecer isso, ou uma camada anterior o resolve.
		// Se o JSON realmente vem com "disciplina" (nome), então o `models.Question` precisaria ser
		// ajustado ou um DTO intermediário seria usado.
		//
		// Baseado no Artefacto 9.3, o JSON tem "disciplina" (nome).
		// O schema da tabela `questions` tem `subject_id`.
		// Isso implica uma busca/criação de Subject.
		// Para este passo, vamos simular que o SubjectID é 1.
		// Em um cenário real, o SubjectRepository seria usado aqui.
		// Se s.subjectRepo for nil, esta parte seria pulada ou tratada de forma diferente.

		var subjectID int64
		// UserID para GetOrCreateByNameAndUser:
		// Se qJSON.UserID não for fornecido ou for 0, precisamos de um UserID padrão ou do contexto.
		// Vamos assumir que qJSON.UserID é o UserID do proprietário da questão/disciplina.
		// Se UserID é um conceito global para o usuário da CLI, ele viria do contexto.
		// Para este exemplo, vamos usar qJSON.UserID, mas com uma verificação.
		currentUserID := qJSON.UserID
		if currentUserID == 0 {
			// Tentar obter UserID do contexto ou usar um padrão. Para agora, erro se não fornecido.
			// Em uma aplicação real, o UserID viria do contexto de autenticação.
			// Se o sistema é monousuário sem autenticação explícita, pode-se usar um UserID fixo (ex: 1).
			// Para este exemplo, vamos assumir que o JSON deve fornecer o UserID.
			return addedCount, fmt.Errorf("questão %d: UserID não fornecido no JSON e necessário para obter/criar disciplina", i)
		}

		if s.subjectRepo != nil {
			subject, err := s.subjectRepo.GetOrCreateByNameAndUser(ctx, qJSON.SubjectName, currentUserID)
			if err != nil {
				return addedCount, fmt.Errorf("questão %d: erro ao obter/criar disciplina '%s' para usuário %d: %w", i, qJSON.SubjectName, currentUserID, err)
			}
			subjectID = subject.ID
		} else {
			// Fallback ou erro se subjectRepo não estiver disponível, mas é necessário.
			return addedCount, fmt.Errorf("questão %d: SubjectRepository não está disponível para resolver disciplina '%s'", i, qJSON.SubjectName)
		}

		var optionsStr *string
		if qJSON.Type == "multipla_escolha" {
			if qJSON.Options == nil {
				return addedCount, fmt.Errorf("questão %d: 'opcoes' é obrigatório para tipo 'multipla_escolha'", i)
			}
			optionsBytes, err := json.Marshal(qJSON.Options)
			if err != nil {
				return addedCount, fmt.Errorf("questão %d: falha ao serializar 'opcoes': %w", i, err)
			}
			s := string(optionsBytes)
			if s == "null" || s == "[]" || s == "" { // "null" se qJSON.Options era nil e foi serializado
				return addedCount, fmt.Errorf("questão %d: 'opcoes' não pode ser vazio para tipo 'multipla_escolha'", i)
			}
			optionsStr = &s
		}

		questionModel := models.Question{
			// UserID: qJSON.UserID, // Assumindo que UserID está no JSON ou vem do contexto
			SubjectID:     subjectID, // Usar o ID obtido/simulado
			Topic:         qJSON.Topic,
			Type:          qJSON.Type,
			Difficulty:    qJSON.Difficulty,
			Statement:     qJSON.Statement,
			Options:       optionsStr,
			CorrectAnswer: qJSON.CorrectAnswer,
		}
		// O UserID deve ser preenchido, idealmente vindo do contexto da requisição (usuário logado)
		// ou se cada questão no JSON pode pertencer a um usuário diferente (menos comum para este cenário).
		// Para este exemplo, vamos assumir que o UserID é um valor fixo ou já está no qJSON.
		// Se qJSON.UserID for 0 e UserID for obrigatório, adicionar validação.
		// Aqui, vamos assumir que UserID está no modelo Question e será preenchido.
		// Se UserID não está no JSON de entrada, ele precisaria ser injetado de outra forma.
		// Por exemplo, se todas as questões são para o mesmo usuário:
		// O UserID deve ser preenchido, idealmente vindo do contexto da requisição (usuário logado)
		// ou se cada questão no JSON pode pertencer a um usuário diferente (menos comum para este cenário).
		// Para este exemplo, vamos assumir que o UserID é um valor fixo ou já está no qJSON.
		// Se qJSON.UserID for 0 e UserID for obrigatório, adicionar validação.
		// Aqui, vamos assumir que UserID está no modelo Question e será preenchido.
		// Se UserID não está no JSON de entrada, ele precisaria ser injetado de outra forma.
		// Por exemplo, se todas as questões são para o mesmo usuário:
		questionModel.UserID = qJSON.UserID // Corrigido: Atribuir UserID do JSON

		questionsToAdd = append(questionsToAdd, questionModel)
	}

	// Idealmente, o repositório teria um método AddManyQuestions para eficiência.
	// Se não, iteramos e adicionamos uma por uma.
	// O mock do repositório pode ser ajustado para esperar AddQuestion para cada item.
	for i := range questionsToAdd {
		// Devemos passar um ponteiro para o item individual do slice questionsToAdd
		// já que AddQuestion provavelmente espera *models.Question
		_, err := s.questionRepo.AddQuestion(ctx, &questionsToAdd[i])
		if err != nil {
			// Retornar o número de questões adicionadas com sucesso antes do erro.
			return addedCount, fmt.Errorf("falha ao adicionar questão %d (%s) ao repositório: %w", addedCount, questionsToAdd[i].Statement, err)
		}
		addedCount++
	}

	return addedCount, nil
}

// GenerateTest gera uma lista de questões para um teste com base nos critérios.
func (s *questionServiceImpl) GenerateTest(ctx context.Context, criteria TestCriteria) ([]models.Question, error) {
	var testQuestions []models.Question

	if criteria.EasyCount == 0 && criteria.MediumCount == 0 && criteria.HardCount == 0 {
		return nil, fmt.Errorf("pelo menos uma contagem de dificuldade deve ser maior que zero para GenerateTest")
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
			return nil, fmt.Errorf("GenerateTest: erro ao buscar questões fáceis: %w", err)
		}
		if len(easyQuestions) < criteria.EasyCount {
			return nil, fmt.Errorf("GenerateTest: não há questões fáceis suficientes (solicitado: %d, disponível: %d)", criteria.EasyCount, len(easyQuestions))
		}
		testQuestions = append(testQuestions, easyQuestions...)
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
			return nil, fmt.Errorf("GenerateTest: erro ao buscar questões médias: %w", err)
		}
		if len(mediumQuestions) < criteria.MediumCount {
			return nil, fmt.Errorf("GenerateTest: não há questões médias suficientes (solicitado: %d, disponível: %d)", criteria.MediumCount, len(mediumQuestions))
		}
		testQuestions = append(testQuestions, mediumQuestions...)
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
			return nil, fmt.Errorf("GenerateTest: erro ao buscar questões difíceis: %w", err)
		}
		if len(hardQuestions) < criteria.HardCount {
			return nil, fmt.Errorf("GenerateTest: não há questões difíceis suficientes (solicitado: %d, disponível: %d)", criteria.HardCount, len(hardQuestions))
		}
		testQuestions = append(testQuestions, hardQuestions...)
	}

	return testQuestions, nil
}
