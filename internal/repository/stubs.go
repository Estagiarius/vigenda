package repository

import (
	"context"
	"database/sql"
	"errors" // Added import for errors
	"fmt"
	"time" // Added import for time
	"vigenda/internal/models"
)

// --- Stub Implementations ---
// These are minimal implementations to allow the application to compile and run for testing CLI commands.
// They do not perform real database operations beyond what's needed for basic schema interaction.

// StubQuestionRepository
type StubQuestionRepository struct {
	DB *sql.DB
}

func NewStubQuestionRepository(db *sql.DB) QuestionRepository {
	return &StubQuestionRepository{DB: db}
}

func (r *StubQuestionRepository) AddQuestion(ctx context.Context, question *models.Question) (int64, error) {
	// In a real implementation, this would insert into the DB.
	fmt.Printf("[StubQuestionRepository] AddQuestion called for: %s\n", question.Statement)
	return 1, nil // Dummy ID
}

func (r *StubQuestionRepository) GetQuestionsByCriteria(ctx context.Context, criteria QuestionQueryCriteria) ([]models.Question, error) {
	fmt.Printf("[StubQuestionRepository] GetQuestionsByCriteria called with: %+v\n", criteria)
	return []models.Question{}, nil // Empty list
}

func (r *StubQuestionRepository) GetQuestionsByCriteriaProofGeneration(ctx context.Context, criteria ProofCriteria) ([]models.Question, error) {
	fmt.Printf("[StubQuestionRepository] GetQuestionsByCriteriaProofGeneration called with: %+v\n", criteria)
	// Return a mix of questions based on criteria, or just an empty list for a stub
	return []models.Question{}, nil
}

// StubSubjectRepository
type StubSubjectRepository struct {
	DB *sql.DB
}

func NewStubSubjectRepository(db *sql.DB) SubjectRepository {
	return &StubSubjectRepository{DB: db}
}

func (r *StubSubjectRepository) GetOrCreateByNameAndUser(ctx context.Context, name string, userID int64) (models.Subject, error) {
	fmt.Printf("[StubSubjectRepository] GetOrCreateByNameAndUser called for: %s, UserID: %d\n", name, userID)
	// Simulate finding or creating
	return models.Subject{ID: 1, Name: name, UserID: userID}, nil
}

// StubTaskRepository
type StubTaskRepository struct {
	DB *sql.DB
}

func NewStubTaskRepository(db *sql.DB) TaskRepository { // Now implements TaskRepository
	return &StubTaskRepository{DB: db}
}
func (r *StubTaskRepository) CreateTask(ctx context.Context, task *models.Task) (int64, error) {
	fmt.Printf("[StubTaskRepository] CreateTask: %s\n", task.Title)
	// Basic INSERT for testing if DB connection works
	stmt, err := r.DB.PrepareContext(ctx, "INSERT INTO tasks (user_id, class_id, title, description, due_date, is_completed) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare insert task: %w", err)
	}
	defer stmt.Close()

	var classID sql.NullInt64
	if task.ClassID != nil {
		classID.Int64 = *task.ClassID
		classID.Valid = true
	}

	var dueDate sql.NullTime
	if task.DueDate != nil {
		dueDate.Time = *task.DueDate
		dueDate.Valid = true
	}


	res, err := stmt.ExecContext(ctx, task.UserID, classID, task.Title, task.Description, dueDate, task.IsCompleted)
	if err != nil {
		return 0, fmt.Errorf("failed to execute insert task: %w", err)
	}
	id, _ := res.LastInsertId()
	return id, nil
}

// GetTaskByID implements repository.TaskRepository.
func (r *StubTaskRepository) GetTaskByID(ctx context.Context, id int64) (*models.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (r *StubTaskRepository) GetTasksByClassID(ctx context.Context, classID int64) ([]models.Task, error) {
	fmt.Printf("[StubTaskRepository] GetTasksByClassID: %d\n", classID)
	// Basic SELECT for testing
	rows, err := r.DB.QueryContext(ctx, "SELECT id, user_id, class_id, title, description, due_date, is_completed FROM tasks WHERE class_id = ?", classID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var classID sql.NullInt64
		var dueDate sql.NullTime
		var description sql.NullString

		err := rows.Scan(&t.ID, &t.UserID, &classID, &t.Title, &description, &dueDate, &t.IsCompleted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}
		if classID.Valid {
			t.ClassID = &classID.Int64
		}
		if description.Valid {
			t.Description = description.String
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *StubTaskRepository) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	fmt.Printf("[StubTaskRepository] GetAllTasks\n")
	rows, err := r.DB.QueryContext(ctx, "SELECT id, user_id, class_id, title, description, due_date, is_completed FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to query all tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var classID sql.NullInt64
		var dueDate sql.NullTime
		var description sql.NullString

		err := rows.Scan(&t.ID, &t.UserID, &classID, &t.Title, &description, &dueDate, &t.IsCompleted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task row for GetAllTasks: %w", err)
		}
		if classID.Valid {
			t.ClassID = &classID.Int64
		}
		if description.Valid {
			t.Description = description.String
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *StubTaskRepository) MarkTaskCompleted(ctx context.Context, taskID int64) error {
	fmt.Printf("[StubTaskRepository] MarkTaskCompleted: %d\n", taskID)
	_, err := r.DB.ExecContext(ctx, "UPDATE tasks SET is_completed = 1 WHERE id = ?", taskID)
	return err
}

func (r *StubTaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	fmt.Printf("[StubTaskRepository] UpdateTask: ID %d, Title: %s\n", task.ID, task.Title)
	// In a real stub, you might update an in-memory list or simulate DB update.
	// For now, just log and return nil.
	return nil
}

func (r *StubTaskRepository) DeleteTask(ctx context.Context, taskID int64) error {
	fmt.Printf("[StubTaskRepository] DeleteTask: ID %d\n", taskID)
	// In a real stub, you might remove from an in-memory list or simulate DB delete.
	// For now, just log and return nil.
	return nil
}

func (r *StubTaskRepository) GetUpcomingTasksByUserID(ctx context.Context, userID int64, limit int) ([]models.Task, error) {
	fmt.Printf("[StubTaskRepository] GetUpcomingTasksByUserID called for UserID: %d, Limit: %d\n", userID, limit)
	// Return an empty list or a predefined list for stubbing purposes
	// This stub currently doesn't interact with r.DB for this method.
	return []models.Task{}, nil
}


// StubClassRepository
type StubClassRepository struct {
	DB *sql.DB
}

// Ensure StubClassRepository implements ClassRepository
var _ ClassRepository = (*StubClassRepository)(nil)

func NewStubClassRepository(db *sql.DB) ClassRepository { // Return interface type
	return &StubClassRepository{DB: db}
}


func (r *StubClassRepository) CreateClass(ctx context.Context, class *models.Class) (int64, error) {
	fmt.Printf("[StubClassRepository] CreateClass: %s\n", class.Name)
	// Example: INSERT INTO classes (user_id, subject_id, name) VALUES (?, ?, ?);
	// Returning dummy ID for now.
	return 1, nil
}
func (r *StubClassRepository) AddStudent(ctx context.Context, student *models.Student) (int64, error) {
    fmt.Printf("[StubClassRepository] AddStudent: %s to class %d\n", student.FullName, student.ClassID)
    return 1, nil // Dummy student ID
}
func (r *StubClassRepository) UpdateStudentStatus(ctx context.Context, studentID int64, status string) error {
    fmt.Printf("[StubClassRepository] UpdateStudentStatus for %d to %s\n", studentID, status)
    return nil
}
func (r *StubClassRepository) GetClassByID(ctx context.Context, classID int64) (*models.Class, error) { // Changed to return *models.Class
	fmt.Printf("[StubClassRepository] GetClassByID: %d\n", classID)
	// Simulate finding the class. In a real repo, query the DB.
	// For the test TestTarefaListarTurmaOutput, we need class with ID 1 to be "Turma 9A".
	if classID == 1 {
		return &models.Class{ID: 1, UserID: 1, SubjectID: 1, Name: "Turma 9A"}, nil // Return as pointer
	}
	return nil, fmt.Errorf("class with ID %d not found in stub", classID) // Return nil for not found
}

// Missing methods from ClassRepository for StubClassRepository
func (r *StubClassRepository) ListAllClasses(ctx context.Context) ([]models.Class, error) {
	fmt.Println("[StubClassRepository] ListAllClasses called")
	return []models.Class{}, nil
}

func (r *StubClassRepository) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
	fmt.Printf("[StubClassRepository] GetStudentsByClassID for class %d called\n", classID)
	return []models.Student{}, nil
}

func (r *StubClassRepository) UpdateClass(ctx context.Context, class *models.Class) error {
	fmt.Printf("[StubClassRepository] UpdateClass for class %d called\n", class.ID)
	return nil
}

func (r *StubClassRepository) DeleteClass(ctx context.Context, classID int64, userID int64) error {
	fmt.Printf("[StubClassRepository] DeleteClass for class %d, user %d called\n", classID, userID)
	return nil
}
func (r *StubClassRepository) GetStudentByID(ctx context.Context, studentID int64) (*models.Student, error) {
	fmt.Printf("[StubClassRepository] GetStudentByID for student %d called\n", studentID)
	return nil, errors.New("GetStudentByID not implemented in stub")
}
func (r *StubClassRepository) UpdateStudent(ctx context.Context, student *models.Student) error {
	fmt.Printf("[StubClassRepository] UpdateStudent for student %d called\n", student.ID)
	return nil
}
func (r *StubClassRepository) DeleteStudent(ctx context.Context, studentID int64, classID int64) error {
	fmt.Printf("[StubClassRepository] DeleteStudent for student %d, class %d called\n", studentID, classID)
	return nil
}

func (r *StubClassRepository) GetTodaysLessonsByUserID(ctx context.Context, userID int64, today time.Time) ([]models.Lesson, error) {
	fmt.Printf("[StubClassRepository] GetTodaysLessonsByUserID called for UserID: %d, Date: %s\n", userID, today.Format("2006-01-02"))
	// Return an empty list or a predefined list for stubbing purposes
	return []models.Lesson{}, nil
}


// StubAssessmentRepository
type StubAssessmentRepository struct {
	DB *sql.DB
}

// Ensure StubAssessmentRepository implements AssessmentRepository
var _ AssessmentRepository = (*StubAssessmentRepository)(nil)

func NewStubAssessmentRepository(db *sql.DB) AssessmentRepository { // Return interface type
	return &StubAssessmentRepository{DB: db}
}

func (r *StubAssessmentRepository) CreateAssessment(ctx context.Context, assessment *models.Assessment) (int64, error) {
	fmt.Printf("[StubAssessmentRepository] CreateAssessment: %s\n", assessment.Name)
	return 1, nil // Dummy ID
}
func (r *StubAssessmentRepository) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
    fmt.Printf("[StubAssessmentRepository] GetStudentsByClassID: %d\n", classID)
    // Return a dummy list of students for testing grade entry
    return []models.Student{
        {ID: 101, FullName: "Alice Wonderland", ClassID: classID, Status: "ativo"},
        {ID: 102, FullName: "Bob The Builder", ClassID: classID, Status: "ativo"},
        {ID: 103, FullName: "Charlie Brown (Transferido)", ClassID: classID, Status: "transferido"},
    }, nil
}
func (r *StubAssessmentRepository) GetAssessmentByID(ctx context.Context, assessmentID int64) (*models.Assessment, error) {
    fmt.Printf("[StubAssessmentRepository] GetAssessmentByID: %d\n", assessmentID)
    // Return a dummy assessment
    return &models.Assessment{ID: assessmentID, Name: "Dummy Assessment", ClassID: 1, Term: 1, Weight: 1.0}, nil
}
func (r *StubAssessmentRepository) EnterGrade(ctx context.Context, grade *models.Grade) error {
    fmt.Printf("[StubAssessmentRepository] EnterGrade for student %d, assessment %d: %.2f\n", grade.StudentID, grade.AssessmentID, grade.Grade)
    return nil
}
func (r *StubAssessmentRepository) GetGradesByClassID(ctx context.Context, classID int64) ([]models.Grade, []models.Assessment, []models.Student, error) {
    fmt.Printf("[StubAssessmentRepository] GetGradesByClassID: %d\n", classID)
    // Return dummy data for average calculation
    return []models.Grade{}, []models.Assessment{}, []models.Student{}, nil
}

func (r *StubAssessmentRepository) ListAllAssessments(ctx context.Context) ([]models.Assessment, error) {
	fmt.Println("[StubAssessmentRepository] ListAllAssessments called")
	return []models.Assessment{}, nil
}
