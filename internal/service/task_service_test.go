package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings" // Added import
	"testing"
	"vigenda/internal/models"
	"vigenda/internal/repository"

	_ "github.com/mattn/go-sqlite3" // DB driver
)

// MockTaskRepository for testing TaskService
type MockTaskRepository struct {
	CreateTaskFunc        func(ctx context.Context, task *models.Task) (int64, error)
	GetTaskByIDFunc       func(ctx context.Context, id int64) (*models.Task, error)
	GetTasksByClassIDFunc func(ctx context.Context, classID int64) ([]models.Task, error)
	GetAllTasksFunc       func(ctx context.Context) ([]models.Task, error)
	MarkTaskCompletedFunc func(ctx context.Context, taskID int64) error

	// Store created bug tasks for verification
	CreatedBugTasks []models.Task
}

func (m *MockTaskRepository) CreateTask(ctx context.Context, task *models.Task) (int64, error) {
	if m.CreateTaskFunc != nil {
		// Check if this is a bug task being created
		if task.UserID == 0 && task.Title[0:5] == "[BUG]" {
			m.CreatedBugTasks = append(m.CreatedBugTasks, *task)
		}
		return m.CreateTaskFunc(ctx, task)
	}
	return 0, errors.New("CreateTaskFunc not implemented in mock")
}

func (m *MockTaskRepository) GetTaskByID(ctx context.Context, id int64) (*models.Task, error) {
	if m.GetTaskByIDFunc != nil {
		return m.GetTaskByIDFunc(ctx, id)
	}
	return nil, errors.New("GetTaskByIDFunc not implemented in mock")
}

func (m *MockTaskRepository) GetTasksByClassID(ctx context.Context, classID int64) ([]models.Task, error) {
	if m.GetTasksByClassIDFunc != nil {
		return m.GetTasksByClassIDFunc(ctx, classID)
	}
	return nil, errors.New("GetTasksByClassIDFunc not implemented in mock")
}

func (m *MockTaskRepository) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	if m.GetAllTasksFunc != nil {
		return m.GetAllTasksFunc(ctx)
	}
	return nil, errors.New("GetAllTasksFunc not implemented in mock")
}

func (m *MockTaskRepository) MarkTaskCompleted(ctx context.Context, taskID int64) error {
	if m.MarkTaskCompletedFunc != nil {
		return m.MarkTaskCompletedFunc(ctx, taskID)
	}
	return errors.New("MarkTaskCompletedFunc not implemented in mock")
}

func TestTaskService_CreateTask(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("successful task creation", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{} // Reset
		expectedID := int64(1)
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) {
			return expectedID, nil
		}

		title := "Test Task"
		desc := "Test Description"
		task, err := taskService.CreateTask(ctx, title, desc, nil, nil)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if task.ID != expectedID {
			t.Errorf("Expected task ID %d, got %d", expectedID, task.ID)
		}
		if task.Title != title {
			t.Errorf("Expected task title %s, got %s", title, task.Title)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks to be created, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("repository error on create", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{} // Reset
		repoError := errors.New("database is down")
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) {
			// Simulate different behavior for bug task creation vs normal task creation
			if task.UserID == 0 && task.Title[0:5] == "[BUG]" { // This is the bug task itself
				return 2, nil // Bug task created successfully
			}
			return 0, repoError // Original task creation fails
		}

		_, err := taskService.CreateTask(ctx, "Another Task", "Desc", nil, nil)

		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, repoError) { // Check if original error is returned
			t.Errorf("Expected error %v, got %v", repoError, err)
		}

		if len(mockRepo.CreatedBugTasks) == 0 {
			t.Errorf("Expected a bug task to be created, but none found")
		} else {
			bugTask := mockRepo.CreatedBugTasks[0]
			expectedBugTitle := "[BUG] Task Creation Failure"
			if bugTask.Title != expectedBugTitle {
				t.Errorf("Expected bug task title '%s', got '%s'", expectedBugTitle, bugTask.Title)
			}
			// Check description contains original error
			if !strings.Contains(bugTask.Description, repoError.Error()) {
				t.Errorf("Bug task description should contain original error '%s', got '%s'", repoError.Error(), bugTask.Description)
			}
		}
	})

	t.Run("validation error (empty title)", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{} // Reset
		_, err := taskService.CreateTask(ctx, "", "Desc", nil, nil)
		if err == nil {
			t.Errorf("Expected validation error for empty title, got nil")
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks for validation error, got %d", len(mockRepo.CreatedBugTasks))
		}
	})
}

func TestTaskService_ListActiveTasksByClass(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()
	classID := int64(1)

	t.Run("successful listing", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		mockTasks := []models.Task{
			{ID: 1, Title: "Task 1", IsCompleted: false, ClassID: &classID},
			{ID: 2, Title: "Task 2", IsCompleted: true, ClassID: &classID}, // Should be filtered out
			{ID: 3, Title: "Task 3", IsCompleted: false, ClassID: &classID},
		}
		mockRepo.GetTasksByClassIDFunc = func(ctx context.Context, cID int64) ([]models.Task, error) {
			if cID == classID {
				return mockTasks, nil
			}
			return nil, fmt.Errorf("unexpected classID: %d", cID)
		}

		tasks, err := taskService.ListActiveTasksByClass(ctx, classID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("Expected 2 active tasks, got %d", len(tasks))
		}
		if tasks[0].Title != "Task 1" || tasks[1].Title != "Task 3" {
			t.Errorf("Unexpected tasks returned: %+v", tasks)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("repository error on listing", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		repoError := errors.New("repo list error")
		mockRepo.GetTasksByClassIDFunc = func(ctx context.Context, cID int64) ([]models.Task, error) {
			return nil, repoError
		}
		// Simulate bug task creation success
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) {
			return 100, nil
		}


		_, err := taskService.ListActiveTasksByClass(ctx, classID)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, repoError) {
			t.Errorf("Expected error %v, got %v", repoError, err)
		}
		if len(mockRepo.CreatedBugTasks) == 0 {
			t.Errorf("Expected a bug task to be created")
		} else if !strings.Contains(mockRepo.CreatedBugTasks[0].Title, "[BUG] Task Listing Failure") {
			t.Errorf("Incorrect bug task title: %s", mockRepo.CreatedBugTasks[0].Title)
		}
	})
}

func TestTaskService_ListAllActiveTasks(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("successful listing all", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		nilClassID := (*int64)(nil)
		classID1 := int64(1)
		mockTasks := []models.Task{
			{ID: 1, Title: "Task 1", IsCompleted: false, ClassID: &classID1},
			{ID: 2, Title: "System Task (Bug)", IsCompleted: false, ClassID: nilClassID, UserID: 0},
			{ID: 3, Title: "Task 3 Completed", IsCompleted: true, ClassID: &classID1},
		}
		mockRepo.GetAllTasksFunc = func(ctx context.Context) ([]models.Task, error) {
			return mockTasks, nil
		}

		tasks, err := taskService.ListAllActiveTasks(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("Expected 2 active tasks, got %d (tasks: %+v)", len(tasks), tasks)
		}
		// Check titles to ensure correct tasks are returned
		titles := []string{}
		for _, task := range tasks {
			titles = append(titles, task.Title)
		}
		if !contains(titles, "Task 1") || !contains(titles, "System Task (Bug)") {
			t.Errorf("Expected 'Task 1' and 'System Task (Bug)', got %v", titles)
		}

		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("repository error on listing all", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		repoError := errors.New("repo list all error")
		mockRepo.GetAllTasksFunc = func(ctx context.Context) ([]models.Task, error) {
			return nil, repoError
		}
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) { return 101, nil } // Bug task creation

		_, err := taskService.ListAllActiveTasks(ctx)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, repoError) {
			t.Errorf("Expected error %v, got %v", repoError, err)
		}
		if len(mockRepo.CreatedBugTasks) == 0 {
			t.Errorf("Expected a bug task to be created")
		} else if !strings.Contains(mockRepo.CreatedBugTasks[0].Title, "[BUG] Task Listing Failure") {
			t.Errorf("Incorrect bug task title: %s", mockRepo.CreatedBugTasks[0].Title)
		}
	})
}


func TestTaskService_MarkTaskAsCompleted(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()
	taskID := int64(1)

	t.Run("successful completion", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		mockRepo.MarkTaskCompletedFunc = func(ctx context.Context, tID int64) error {
			if tID == taskID {
				return nil
			}
			return fmt.Errorf("unexpected taskID: %d", tID)
		}

		err := taskService.MarkTaskAsCompleted(ctx, taskID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("repository error on completion", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		repoError := errors.New("repo mark completed error")
		mockRepo.MarkTaskCompletedFunc = func(ctx context.Context, tID int64) error {
			return repoError
		}
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) { return 102, nil } // Bug task creation

		err := taskService.MarkTaskAsCompleted(ctx, taskID)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, repoError) {
			t.Errorf("Expected error %v, got %v", repoError, err)
		}
		if len(mockRepo.CreatedBugTasks) == 0 {
			t.Errorf("Expected a bug task to be created")
		} else if !strings.Contains(mockRepo.CreatedBugTasks[0].Title, "[BUG] Task Completion Failure") {
			t.Errorf("Incorrect bug task title: %s", mockRepo.CreatedBugTasks[0].Title)
		}
	})
}


// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper function to set up a temporary DB for tests that might need it (though most use mocks)
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	// Use in-memory SQLite for testing; ensure a unique name if tests run in parallel
	// or ensure cleanup. For simple sequential tests, "file::memory:" is often fine.
	// For parallel, "file::memory:?cache=shared" or unique temp files.
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Run migrations if your stubs or actual repositories need the schema
	schema, err := os.ReadFile("../../internal/database/migrations/001_initial_schema.sql")
	if err != nil {
		db.Close()
		t.Fatalf("Failed to read schema file: %v", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		db.Close()
		t.Fatalf("Failed to apply schema: %v", err)
	}

	return db
}

func TestMain(m *testing.M) {
	// Setup tasks before running tests, e.g., DB connection or loading config
	// For these tests, we primarily use mocks, so global setup might not be strictly needed here.
	// However, if any test relies on a real DB (even sqlite in-memory for stubs),
	// it could be initialized here or per-test.

	exitCode := m.Run()

	// Teardown tasks after tests are done
	os.Exit(exitCode)
}
// Add import for strings at the top of the file
// import "strings"
