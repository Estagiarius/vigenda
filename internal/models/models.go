// Package models defines all Go structs used in the Vigenda application.
// These structs represent entities like Task, Lesson, Student, Class, etc.
package models

import "time"

// Task represents a task in the system.
type Task struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	ClassID     *int64     `json:"class_id,omitempty"` // Pointer to allow null
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`   // Pointer to allow null
	IsCompleted bool       `json:"is_completed"`
}

// Class represents a class or turma.
type Class struct {
	ID        int64
	UserID    int64
	SubjectID int64
	Name      string
}

// Student represents a student in a class.
type Student struct {
	ID           int64
	ClassID      int64
	FullName     string
	EnrollmentID string
	Status       string // 'ativo', 'inativo', 'transferido'
}

// Assessment represents an assessment or evaluation.
type Assessment struct {
	ID             int64
	ClassID        int64
	Name           string
	Term           int
	Weight         float64
	AssessmentDate *time.Time
}

// Grade represents a grade given to a student for an assessment.
type Grade struct {
	ID           int64
	AssessmentID int64
	StudentID    int64
	Grade        float64
}

// Question represents a question in the question bank.
type Question struct {
	ID            int64   `json:"id"`
	UserID        int64   `json:"user_id"`
	SubjectID     int64   `json:"subject_id"`
	Topic         string  `json:"topic"`
	Type          string  `json:"type"` // 'multipla_escolha' ou 'dissertativa'
	Difficulty    string  `json:"difficulty"` // 'facil', 'media', 'dificil'
	Statement     string  `json:"statement"`
	Options       *string `json:"options"` // JSON array como string para multipla escolha
	CorrectAnswer string  `json:"correct_answer"`
}

// Subject represents a subject or discipline.
type Subject struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
}

// Custom error for when a class is not found
type ModelError string

func (e ModelError) Error() string {
	return string(e)
}

const (
	ErrClassNotFound ModelError = "class not found"
)
