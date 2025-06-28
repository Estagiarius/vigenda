package service

import (
	"context"
	"database/sql" // Required for DB field in stubs if they need it directly
	"fmt"
	"strings"
	"time"
	"encoding/csv" // For CSV processing in stub

	"vigenda/internal/models"
	"vigenda/internal/repository" // Assuming stubs might need repository interfaces
)

// --- Stub Service Implementations ---
// These are minimal implementations to allow the application to compile and run for testing CLI commands.
// They use stub repositories.

// StubTaskService
type stubTaskService struct {
	taskRepo *repository.StubTaskRepository // Using specific stub type
}

// NewStubTaskService creates a new stub instance of TaskService.
func NewStubTaskService(taskRepo *repository.StubTaskRepository) TaskService {
	return &stubTaskService{taskRepo: taskRepo}
}

func (s *stubTaskService) CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error) {
	fmt.Printf("[StubTaskService] CreateTask called: %s\n", title)
	task := models.Task{
		UserID:      1, // Assuming a default user ID for stubs
		Title:       title,
		Description: description,
		ClassID:     classID,
		DueDate:     dueDate,
		IsCompleted: false,
	}
	id, err := s.taskRepo.CreateTask(ctx, &task)
	if err != nil {
		return models.Task{}, err
	}
	task.ID = id
	return task, nil
}

func (s *stubTaskService) ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error) {
	fmt.Printf("[StubTaskService] ListActiveTasksByClass called for ClassID: %d\n", classID)
	return s.taskRepo.GetTasksByClassID(ctx, classID) // Assuming stub repo has this
}

func (s *stubTaskService) MarkTaskAsCompleted(ctx context.Context, taskID int64) error {
	fmt.Printf("[StubTaskService] MarkTaskAsCompleted called for TaskID: %d\n", taskID)
	return s.taskRepo.MarkTaskCompleted(ctx, taskID) // Assuming stub repo has this
}

func (s *stubTaskService) ListAllActiveTasks(ctx context.Context) ([]models.Task, error) {
	fmt.Printf("[StubTaskService] ListAllActiveTasks called\n")
	// In a real stub, you might want to mock this further or provide specific stub data.
	// For now, just pass through to the repository stub.
	allTasks, err := s.taskRepo.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	// The service layer is responsible for filtering active tasks
	activeTasks := []models.Task{}
	for _, task := range allTasks {
		if !task.IsCompleted {
			activeTasks = append(activeTasks, task)
		}
	}
	return activeTasks, nil
}

// StubClassService
type stubClassService struct {
	classRepo *repository.StubClassRepository
	// userRepo repository.UserRepository // If needed for user context
}

// NewStubClassService creates a new stub instance of ClassService.
func NewStubClassService(classRepo *repository.StubClassRepository /*, userRepo repository.UserRepository*/) ClassService {
	return &stubClassService{classRepo: classRepo}
}

func (s *stubClassService) CreateClass(ctx context.Context, name string, subjectID int64) (models.Class, error) {
	fmt.Printf("[StubClassService] CreateClass: %s, SubjectID: %d\n", name, subjectID)
	class := models.Class{
		UserID:    1, // Default user
		SubjectID: subjectID,
		Name:      name,
	}
	id, err := s.classRepo.CreateClass(ctx, &class)
	if err != nil {
		return models.Class{}, err
	}
	class.ID = id
	return class, nil
}

func (s *stubClassService) ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error) {
    fmt.Printf("[StubClassService] ImportStudentsFromCSV for ClassID: %d\n", classID)
    reader := csv.NewReader(strings.NewReader(string(csvData)))
    records, err := reader.ReadAll()
    if err != nil {
        return 0, fmt.Errorf("error reading CSV data: %w", err)
    }

    count := 0
    if len(records) > 1 { // Skip header row
        for _, record := range records[1:] {
            if len(record) < 2 { continue } // Expect at least numero_chamada, nome_completo

            status := "ativo"
            if len(record) >= 3 && record[2] != "" {
                status = record[2]
            }
            student := models.Student{
                ClassID:      classID,
                EnrollmentID: record[0],
                FullName:     record[1],
                Status:       status,
            }
            _, err := s.classRepo.AddStudent(ctx, &student) // Assumes AddStudent exists on stub repo
            if err != nil {
                fmt.Printf("Error adding student %s: %v\n", student.FullName, err)
                // Optionally decide whether to stop or continue
            } else {
                count++
            }
        }
    }
    return count, nil
}


func (s *stubClassService) UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error {
	fmt.Printf("[StubClassService] UpdateStudentStatus for StudentID %d to %s\n", studentID, newStatus)
	return s.classRepo.UpdateStudentStatus(ctx, studentID, newStatus)
}

func (s *stubClassService) GetClassByID(ctx context.Context, classID int64) (models.Class, error) {
	fmt.Printf("[StubClassService] GetClassByID: %d\n", classID)
	return s.classRepo.GetClassByID(ctx, classID)
}

func (s *stubClassService) ListAllClasses(ctx context.Context) ([]models.Class, error) {
	fmt.Printf("[StubClassService] ListAllClasses called\n")
	// Retorna uma lista pré-definida de turmas ou uma lista vazia para o stub
	return []models.Class{
		{ID: 1, UserID: 1, SubjectID: 101, Name: "Turma Stub A"},
		{ID: 2, UserID: 1, SubjectID: 102, Name: "Turma Stub B"},
	}, nil
}

func (s *stubClassService) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
	fmt.Printf("[StubClassService] GetStudentsByClassID called for ClassID: %d\n", classID)
	// Retorna uma lista vazia por padrão para o stub.
	// Se o StubClassRepository tiver um método GetStudentsByClassID, poderia chamá-lo:
	// return s.classRepo.GetStudentsByClassID(ctx, classID)
	return []models.Student{}, nil
}

// StubAssessmentService
type stubAssessmentService struct {
	assessmentRepo *repository.StubAssessmentRepository
	// classRepo repository.ClassRepository // Might be needed to fetch student lists
}

// NewStubAssessmentService creates a new stub instance of AssessmentService.
func NewStubAssessmentService(assessmentRepo *repository.StubAssessmentRepository) AssessmentService {
	return &stubAssessmentService{assessmentRepo: assessmentRepo}
}

func (s *stubAssessmentService) CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error) {
	fmt.Printf("[StubAssessmentService] CreateAssessment: %s for ClassID %d\n", name, classID)
	assessment := models.Assessment{
		ClassID: classID,
		Name:    name,
		Term:    term,
		Weight:  weight,
	}
	id, err := s.assessmentRepo.CreateAssessment(ctx, &assessment)
	if err != nil {
		return models.Assessment{}, err
	}
	assessment.ID = id
	return assessment, nil
}

func (s *stubAssessmentService) EnterGrades(ctx context.Context, assessmentID int64, studentGrades map[int64]float64) error {
	fmt.Printf("[StubAssessmentService] EnterGrades for AssessmentID %d: %+v\n", assessmentID, studentGrades)
	// In a real service, you'd fetch the assessment, validate students belong to the class, etc.
	for studentID, gradeVal := range studentGrades {
		grade := models.Grade{
			AssessmentID: assessmentID,
			StudentID:    studentID,
			Grade:        gradeVal,
		}
		if err := s.assessmentRepo.EnterGrade(ctx, &grade); err != nil {
			return fmt.Errorf("error entering grade for student %d: %w", studentID, err)
		}
	}
	return nil
}

func (s *stubAssessmentService) CalculateClassAverage(ctx context.Context, classID int64) (float64, error) {
	fmt.Printf("[StubAssessmentService] CalculateClassAverage for ClassID %d\n", classID)
	// This is a complex calculation in reality. Stub returns a fixed value.
	return 7.5, nil
}

func (s *stubAssessmentService) ListAllAssessments(ctx context.Context) ([]models.Assessment, error) {
	fmt.Printf("[StubAssessmentService] ListAllAssessments called\n")
	// Retorna uma lista pré-definida de avaliações ou uma lista vazia para o stub
	now := time.Now()
	return []models.Assessment{
		{ID: 1, ClassID: 1, Name: "Prova 1 Stub", Term: 1, Weight: 2, AssessmentDate: &now},
		{ID: 2, ClassID: 1, Name: "Trabalho 1 Stub", Term: 1, Weight: 1.5, AssessmentDate: &now},
	}, nil
}

// StubQuestionService
type stubQuestionService struct {
	questionRepo repository.QuestionRepository
	subjectRepo  repository.SubjectRepository
	db           *sql.DB // Or pass DB to repositories if they don't manage their own connection
}

// Constructor NewQuestionService is now expected to be used from actual question_service.go
// func NewQuestionService(db *sql.DB, qRepo repository.QuestionRepository, sRepo repository.SubjectRepository) QuestionService {
// 	return &stubQuestionService{db: db, questionRepo: qRepo, subjectRepo: sRepo}
// }

func (s *stubQuestionService) AddQuestionsFromJSON(ctx context.Context, jsonData []byte) (int, error) {
	fmt.Printf("[StubQuestionService] AddQuestionsFromJSON called with %d bytes of JSON data\n", len(jsonData))
	// Dummy implementation: parse JSON, call repository.AddQuestion for each.
	// For now, just return a dummy count.
	return 5, nil // Simulate 5 questions added
}

func (s *stubQuestionService) GenerateTest(ctx context.Context, criteria TestCriteria) ([]models.Question, error) {
	fmt.Printf("[StubQuestionService] GenerateTest called with criteria: %+v\n", criteria)
	// Dummy implementation: call repository.GetQuestionsByCriteria.
	return []models.Question{}, nil // Return empty list
}

// StubProofService
type stubProofService struct {
	questionRepo repository.QuestionRepository
	// db *sql.DB // if needed
}

// Constructor NewProofService is now expected to be used from actual proof_service.go
// func NewProofService(qRepo repository.QuestionRepository) ProofService {
// 	return &stubProofService{questionRepo: qRepo}
// }

func (s *stubProofService) GenerateProof(ctx context.Context, criteria ProofCriteria) ([]models.Question, error) {
	fmt.Printf("[StubProofService] GenerateProof called with criteria: %+v\n", criteria)
	// Convert ProofCriteria to QuestionQueryCriteria if necessary, or have repo support ProofCriteria
	// For now, directly call with a simplified mapping if possible or just return empty.

	// Example of mapping (simplified, might need more fields)
	qCriteria := repository.QuestionQueryCriteria{
		SubjectID: criteria.SubjectID,
		Topic:     criteria.Topic,
		// Difficulty and Limit would need to be handled based on counts
	}

	// This is a simplification. A real implementation would fetch easy, medium, hard questions
	// based on counts and combine them.
	if criteria.EasyCount > 0 {
		 qCriteria.Difficulty = "facil"
		 qCriteria.Limit = criteria.EasyCount
		 // fetch and add to list
	}
    // ... similar for medium and hard

	return s.questionRepo.GetQuestionsByCriteria(ctx, qCriteria)
}
