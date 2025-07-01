package service

import (
	"context"
	// "database/sql" // No longer directly needed by service stubs
	"encoding/csv" // For CSV processing in stub
	"fmt"
	"strings"
	"time"
	"vigenda/internal/models"
	"vigenda/internal/repository" // Assuming stubs might need repository interfaces
)

// --- Stub Service Implementations ---
// These are minimal implementations to allow the application to compile and run for testing CLI commands.
// They use stub repositories.

// StubTaskService
type stubTaskService struct {
	// Allow using the interface type for flexibility, though concrete stub might be passed.
	taskRepo repository.TaskRepository
}

// NewStubTaskService creates a new stub instance of TaskService.
func NewStubTaskService(taskRepo repository.TaskRepository) TaskService { // Changed param type
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

func (s *stubTaskService) GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error) {
	fmt.Printf("[StubTaskService] GetTaskByID called for TaskID: %d\n", taskID)
	return s.taskRepo.GetTaskByID(ctx, taskID)
}

func (s *stubTaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	fmt.Printf("[StubTaskService] UpdateTask called for TaskID: %d, Title: %s\n", task.ID, task.Title)
	return s.taskRepo.UpdateTask(ctx, task)
}

func (s *stubTaskService) DeleteTask(ctx context.Context, taskID int64) error {
	fmt.Printf("[StubTaskService] DeleteTask called for TaskID: %d\n", taskID)
	return s.taskRepo.DeleteTask(ctx, taskID)
}

func (s *stubTaskService) ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error) {
	fmt.Printf("[StubTaskService] ListActiveTasksByClass called for ClassID: %d\n", classID)
	allTasks, err := s.taskRepo.GetTasksByClassID(ctx, classID)
	if err != nil {
		return nil, err
	}
	activeTasks := []models.Task{}
	for _, task := range allTasks {
		if !task.IsCompleted {
			activeTasks = append(activeTasks, task)
		}
	}
	return activeTasks, nil
}

func (s *stubTaskService) MarkTaskAsCompleted(ctx context.Context, taskID int64) error {
	fmt.Printf("[StubTaskService] MarkTaskAsCompleted called for TaskID: %d\n", taskID)
	return s.taskRepo.MarkTaskCompleted(ctx, taskID)
}

func (s *stubTaskService) ListAllActiveTasks(ctx context.Context) ([]models.Task, error) {
	fmt.Printf("[StubTaskService] ListAllActiveTasks called\n")
	allTasks, err := s.taskRepo.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
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
	classRepo repository.ClassRepository
}

// NewStubClassService creates a new stub instance of ClassService.
func NewStubClassService(classRepo repository.ClassRepository) ClassService {
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
			if len(record) < 2 {
				continue
			}

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
			_, err := s.classRepo.AddStudent(ctx, &student)
			if err != nil {
				fmt.Printf("Error adding student %s: %v\n", student.FullName, err)
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
	cls, err := s.classRepo.GetClassByID(ctx, classID)
	if err != nil {
		return models.Class{}, err
	}
	if cls == nil {
		return models.Class{}, fmt.Errorf("stub class repo returned nil class for ID %d", classID)
	}
	return *cls, nil
}

func (s *stubClassService) ListAllClasses(ctx context.Context) ([]models.Class, error) {
	fmt.Printf("[StubClassService] ListAllClasses called\n")
	return []models.Class{
		{ID: 1, UserID: 1, SubjectID: 101, Name: "Turma Stub A"},
		{ID: 2, UserID: 1, SubjectID: 102, Name: "Turma Stub B"},
	}, nil
}

func (s *stubClassService) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
	fmt.Printf("[StubClassService] GetStudentsByClassID called for ClassID: %d\n", classID)
	return s.classRepo.GetStudentsByClassID(ctx, classID)
}

// StubAssessmentService
type stubAssessmentService struct {
	assessmentRepo repository.AssessmentRepository
}

// NewStubAssessmentService creates a new stub instance of AssessmentService.
func NewStubAssessmentService(assessmentRepo repository.AssessmentRepository) AssessmentService {
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
	return 7.5, nil
}

func (s *stubAssessmentService) ListAllAssessments(ctx context.Context) ([]models.Assessment, error) {
	fmt.Printf("[StubAssessmentService] ListAllAssessments called\n")
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
}

func NewStubQuestionService(qRepo repository.QuestionRepository, sRepo repository.SubjectRepository) QuestionService {
	return &stubQuestionService{questionRepo: qRepo, subjectRepo: sRepo}
}

func (s *stubQuestionService) AddQuestionsFromJSON(ctx context.Context, jsonData []byte) (int, error) {
	fmt.Printf("[StubQuestionService] AddQuestionsFromJSON called with %d bytes of JSON data\n", len(jsonData))
	return 5, nil
}

func (s *stubQuestionService) GenerateTest(ctx context.Context, criteria TestCriteria) ([]models.Question, error) {
	fmt.Printf("[StubQuestionService] GenerateTest called with criteria: %+v\n", criteria)
	return []models.Question{}, nil
}

// StubProofService
type stubProofService struct {
	questionRepo repository.QuestionRepository
}

func NewStubProofService(qRepo repository.QuestionRepository) ProofService {
	return &stubProofService{questionRepo: qRepo}
}

func (s *stubProofService) GenerateProof(ctx context.Context, criteria ProofCriteria) ([]models.Question, error) {
	fmt.Printf("[StubProofService] GenerateProof called with criteria: %+v\n", criteria)
	qCriteria := repository.QuestionQueryCriteria{
		SubjectID: criteria.SubjectID,
		Topic:     criteria.Topic,
	}
	return s.questionRepo.GetQuestionsByCriteria(ctx, qCriteria)
}
