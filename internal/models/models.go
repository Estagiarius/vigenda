import "time"

// Package models defines all Go structs used in the Vigenda application.
// These structs represent entities like Task, Lesson, Student, Class, etc.
package models

// Task represents a task in the system.
type Task struct{}

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
type Question struct{}
