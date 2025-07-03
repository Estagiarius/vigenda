// Package models defines all Go structs used in the Vigenda application.
// These structs represent entities like Task, Lesson, Student, Class, Subject,
// Assessment, Grade, Question, and User. They are primarily used for data
// transfer between layers (service, repository) and for database interaction.
package models

import "time"

// User represents a user of the Vigenda application.
// Users can own subjects, classes, tasks, and questions.
type User struct {
	ID           int64  `json:"id"`                    // ID é o identificador único do usuário, gerado pelo banco de dados.
	Username     string `json:"username"`              // Username é o nome de usuário para login, deve ser único.
	PasswordHash string `json:"-"`                     // PasswordHash é o hash da senha do usuário, não exposto em JSON.
}

// Subject represents a subject or discipline (e.g., Mathematics, History).
// Each subject is associated with a user.
type Subject struct {
	ID     int64  `json:"id"`      // ID é o identificador único da disciplina.
	UserID int64  `json:"user_id"` // UserID é o ID do usuário proprietário desta disciplina.
	Name   string `json:"name"`    // Name é o nome da disciplina.
}

// Class represents a specific class or group of students within a subject
// (e.g., "Calculus I - Section A", "History 101 - Morning Class").
type Class struct {
	ID        int64  `json:"id"`         // ID é o identificador único da turma.
	UserID    int64  `json:"user_id"`    // UserID é o ID do usuário proprietário desta turma (indiretamente via Subject).
	SubjectID int64  `json:"subject_id"` // SubjectID é o ID da disciplina à qual esta turma pertence.
	Name      string `json:"name"`       // Name é o nome da turma (ex: "Turma 9A - 2025").
}

// Student represents a student enrolled in a specific class.
type Student struct {
	ID           int64  `json:"id"`            // ID é o identificador único do estudante.
	ClassID      int64  `json:"class_id"`      // ClassID é o ID da turma à qual o estudante pertence.
	FullName     string `json:"full_name"`     // FullName é o nome completo do estudante.
	EnrollmentID string `json:"enrollment_id"` // EnrollmentID é o número de matrícula ou de chamada (opcional).
	Status       string `json:"status"`        // Status indica a situação do estudante (ex: 'ativo', 'inativo', 'transferido').
}

// Lesson represents a planned lesson for a class.
type Lesson struct {
	ID           int64     `json:"id"`            // ID é o identificador único da aula.
	ClassID      int64     `json:"class_id"`      // ClassID é o ID da turma para a qual a aula é planejada.
	Title        string    `json:"title"`         // Title é o título da aula.
	PlanContent  string    `json:"plan_content"`  // PlanContent contém o conteúdo do plano de aula, preferencialmente em Markdown.
	ScheduledAt  time.Time `json:"scheduled_at"`  // ScheduledAt é a data e hora agendada para a aula.
}

// Assessment represents an assessment or evaluation (e.g., test, quiz, project) for a class.
type Assessment struct {
	ID             int64      `json:"id"`               // ID é o identificador único da avaliação.
	ClassID        int64      `json:"class_id"`         // ClassID é o ID da turma para a qual a avaliação é aplicada.
	Name           string     `json:"name"`             // Name é o nome da avaliação (ex: "Prova Bimestral 1").
	Term           int        `json:"term"`             // Term indica o período da avaliação (ex: 1 para o primeiro bimestre/trimestre).
	Weight         float64    `json:"weight"`           // Weight é o peso da avaliação na composição da nota final.
	AssessmentDate *time.Time `json:"assessment_date,omitempty"` // AssessmentDate é a data de aplicação da avaliação (ponteiro para permitir nulo).
}

// Grade represents a grade received by a student for a specific assessment.
type Grade struct {
	ID           int64   `json:"id"`            // ID é o identificador único do registro de nota.
	AssessmentID int64   `json:"assessment_id"` // AssessmentID é o ID da avaliação à qual esta nota se refere.
	StudentID    int64   `json:"student_id"`    // StudentID é o ID do estudante que recebeu a nota.
	Grade        float64 `json:"grade"`         // Grade é a nota numérica obtida pelo estudante.
}

// Task represents a task or to-do item.
// Tasks can be personal (associated only with a user) or related to a specific class.
type Task struct {
	ID          int64      `json:"id"`                    // ID é o identificador único da tarefa.
	UserID      int64      `json:"user_id"`               // UserID é o ID do usuário proprietário desta tarefa.
	ClassID     *int64     `json:"class_id,omitempty"`    // ClassID (opcional) é o ID da turma à qual esta tarefa pode estar associada. Ponteiro para permitir nulo.
	Title       string     `json:"title"`                 // Title é o título da tarefa.
	Description string     `json:"description,omitempty"` // Description fornece detalhes adicionais sobre a tarefa (opcional).
	DueDate     *time.Time `json:"due_date,omitempty"`    // DueDate é a data e hora de vencimento da tarefa (opcional). Ponteiro para permitir nulo.
	IsCompleted bool       `json:"is_completed"`          // IsCompleted indica se a tarefa foi concluída.
}

// Question represents a question stored in the question bank.
// Questions are associated with a user and a subject, and can be used to create assessments.
type Question struct {
	ID            int64   `json:"id"`             // ID é o identificador único da questão.
	UserID        int64   `json:"user_id"`        // UserID é o ID do usuário proprietário desta questão.
	SubjectID     int64   `json:"subject_id"`     // SubjectID é o ID da disciplina à qual esta questão está relacionada.
	Topic         string  `json:"topic"`          // Topic é um tópico específico dentro da disciplina (opcional).
	Type          string  `json:"type"`           // Type indica o tipo da questão (ex: 'multipla_escolha', 'dissertativa').
	Difficulty    string  `json:"difficulty"`     // Difficulty indica o nível de dificuldade (ex: 'facil', 'media', 'dificil').
	Statement     string  `json:"statement"`      // Statement é o enunciado da questão.
	Options       *string `json:"options"`        // Options (opcional) contém um array JSON (como string) para questões de múltipla escolha. Ex: `["Opção A", "Opção B"]`. Ponteiro para permitir nulo.
	CorrectAnswer string  `json:"correct_answer"` // CorrectAnswer armazena a resposta correta. Para múltipla escolha, pode ser o texto da opção ou um índice.
}

// ModelError é um tipo customizado para erros específicos da camada de modelo.
type ModelError string

// Error implementa a interface error para ModelError.
func (e ModelError) Error() string {
	return string(e)
}

const (
	ErrClassNotFound ModelError = "class not found"
)
