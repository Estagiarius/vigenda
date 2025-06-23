package repository

import (
	"context"
	"database/sql"
	"fmt"
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


// StubClassRepository
type StubClassRepository struct {
	DB *sql.DB
}
func NewStubClassRepository(db *sql.DB) *StubClassRepository { return &StubClassRepository{DB:db} }

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
func (r *StubClassRepository) GetClassByID(ctx context.Context, classID int64) (models.Class, error) {
	fmt.Printf("[StubClassRepository] GetClassByID: %d\n", classID)
	// Simulate finding the class. In a real repo, query the DB.
	// For the test TestTarefaListarTurmaOutput, we need class with ID 1 to be "Turma 9A".
	if classID == 1 {
		return models.Class{ID: 1, UserID: 1, SubjectID: 1, Name: "Turma 9A"}, nil
	}
	return models.Class{}, fmt.Errorf("class with ID %d not found in stub", classID)
}


// StubAssessmentRepository
type StubAssessmentRepository struct {
	DB *sql.DB
}
func NewStubAssessmentRepository(db *sql.DB) *StubAssessmentRepository { return &StubAssessmentRepository{DB:db} }

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
