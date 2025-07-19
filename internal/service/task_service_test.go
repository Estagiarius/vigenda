package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings" // Added import
	"testing"
	"time"
	"vigenda/internal/models"
	// "vigenda/internal/repository" // repository package is not directly used by this test file, types are defined locally or through models.

	_ "github.com/mattn/go-sqlite3" // DB driver
)

// MockTaskRepository for testing TaskService
type MockTaskRepository struct {
	CreateTaskFunc        func(ctx context.Context, task *models.Task) (int64, error)
	GetTaskByIDFunc       func(ctx context.Context, id int64) (*models.Task, error)
	GetTasksByClassIDFunc func(ctx context.Context, classID int64) ([]models.Task, error)
	GetAllTasksFunc       func(ctx context.Context) ([]models.Task, error)
	MarkTaskCompletedFunc func(ctx context.Context, taskID int64) error
	UpdateTaskFunc        func(ctx context.Context, task *models.Task) error // Added
	DeleteTaskFunc        func(ctx context.Context, taskID int64) error    // Added

	// Store created bug tasks for verification
	CreatedBugTasks []models.Task
}

func (m *MockTaskRepository) GetUpcomingActiveTasks(ctx context.Context, userID int64, fromDate time.Time, limit int) ([]models.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTaskRepository) CreateTask(ctx context.Context, task *models.Task) (int64, error) {
	if task.UserID == 0 && strings.HasPrefix(task.Title, "[BUG]") {
		m.CreatedBugTasks = append(m.CreatedBugTasks, *task)
	}
	if m.CreateTaskFunc != nil {
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

func (m *MockTaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	if m.UpdateTaskFunc != nil {
		return m.UpdateTaskFunc(ctx, task)
	}
	return errors.New("UpdateTaskFunc not implemented in mock")
}

func (m *MockTaskRepository) DeleteTask(ctx context.Context, taskID int64) error {
	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(ctx, taskID)
	}
	return errors.New("DeleteTaskFunc not implemented in mock")
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
			expectedBugTitle := "[BUG][AUTO][PRIORITY_PENDING] Falha na Criação de Tarefa"
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
		} else if !strings.Contains(mockRepo.CreatedBugTasks[0].Title, "[BUG][AUTO][PRIORITY_PENDING] Falha na Listagem de Tarefas por Turma") {
			t.Errorf("Incorrect bug task title: %s", mockRepo.CreatedBugTasks[0].Title)
		}
	})
}

func TestTaskService_ListAllTasks(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("successful listing all tasks (active and completed)", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		nilClassID := (*int64)(nil)
		classID1 := int64(1)
		mockTasks := []models.Task{
			{ID: 1, Title: "Task 1 Active", IsCompleted: false, ClassID: &classID1},
			{ID: 2, Title: "System Task Active (Bug)", IsCompleted: false, ClassID: nilClassID, UserID: 0},
			{ID: 3, Title: "Task 3 Completed", IsCompleted: true, ClassID: &classID1},
			{ID: 4, Title: "Task 4 Active", IsCompleted: false, ClassID: &classID1},
		}
		mockRepo.GetAllTasksFunc = func(ctx context.Context) ([]models.Task, error) {
			return mockTasks, nil
		}

		tasks, err := taskService.ListAllTasks(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(tasks) != 4 { // Should return all tasks, filtering is up to caller
			t.Errorf("Expected 4 tasks, got %d (tasks: %+v)", len(tasks), tasks)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("repository error on listing all tasks", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		repoError := errors.New("repo list all error")
		mockRepo.GetAllTasksFunc = func(ctx context.Context) ([]models.Task, error) {
			return nil, repoError
		}
		// Simulate bug task creation success
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) { return 101, nil }

		_, err := taskService.ListAllTasks(ctx)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, repoError) {
			t.Errorf("Expected error %v, got %v", repoError, err)
		}
		if len(mockRepo.CreatedBugTasks) == 0 {
			t.Errorf("Expected a bug task to be created")
		} else if !strings.Contains(mockRepo.CreatedBugTasks[0].Title, "[BUG][AUTO][PRIORITY_PENDING] Falha na Listagem Global de Tarefas") { // Adjusted title
			t.Errorf("Incorrect bug task title: %s", mockRepo.CreatedBugTasks[0].Title)
		}
	})
}

func TestTaskService_ListAllActiveTasks(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("successful listing all active tasks", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		nilClassID := (*int64)(nil)
		classID1 := int64(1)
		mockTasksFromRepo := []models.Task{
			{ID: 1, Title: "Task 1 Active", IsCompleted: false, ClassID: &classID1},
			{ID: 2, Title: "System Task Active (Bug)", IsCompleted: false, ClassID: nilClassID, UserID: 0},
			{ID: 3, Title: "Task 3 Completed", IsCompleted: true, ClassID: &classID1}, // Should be filtered out
			{ID: 4, Title: "Task 4 Active", IsCompleted: false, ClassID: &classID1},
			{ID: 5, Title: "Task 5 Completed", IsCompleted: true, ClassID: nilClassID, UserID: 0}, // Should be filtered out
		}
		mockRepo.GetAllTasksFunc = func(ctx context.Context) ([]models.Task, error) {
			return mockTasksFromRepo, nil
		}

		tasks, err := taskService.ListAllActiveTasks(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(tasks) != 3 {
			t.Errorf("Expected 3 active tasks, got %d (tasks: %+v)", len(tasks), tasks)
		}
		// Check titles to ensure correct tasks are returned
		titles := []string{}
		for _, task := range tasks {
			titles = append(titles, task.Title)
		}
		// expectedTitles := []string{"Task 1 Active", "System Task Active (Bug)", "Task 4 Active"} // This line was unused
		for _, et := range []string{"Task 1 Active", "System Task Active (Bug)", "Task 4 Active"} {
			if !contains(titles, et) {
				t.Errorf("Expected title '%s' to be in active tasks, but not found. Got: %v", et, titles)
			}
		}
		if contains(titles, "Task 3 Completed") || contains(titles, "Task 5 Completed") {
			t.Errorf("Completed tasks should not be in the active list. Got: %v", titles)
		}

		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("repository error on listing all active tasks", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		repoError := errors.New("repo list all active error")
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
		} else if !strings.Contains(mockRepo.CreatedBugTasks[0].Title, "[BUG][AUTO][PRIORITY_PENDING] Falha na Listagem Global de Tarefas Ativas") { // Corrected expected bug title
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
		} else if !strings.Contains(mockRepo.CreatedBugTasks[0].Title, "[BUG][AUTO][PRIORITY_PENDING] Task Completion Failure") {
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
// This function is not currently used by the tests above as they use mocks primarily.
// It's kept for potential future tests or if stubs required more complex DB interactions.
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
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

func TestTaskService_GetTaskByID(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()
	taskID := int64(1)
	expectedTask := &models.Task{ID: taskID, Title: "Test Task", UserID: 1}

	t.Run("successful get", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		mockRepo.GetTaskByIDFunc = func(ctx context.Context, id int64) (*models.Task, error) {
			if id == taskID {
				return expectedTask, nil
			}
			return nil, fmt.Errorf("unexpected taskID: %d", id)
		}

		task, err := taskService.GetTaskByID(ctx, taskID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if task == nil || task.ID != taskID {
			t.Errorf("Expected task with ID %d, got %+v", taskID, task)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("task not found", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		// Simulate repository returning a specific "not found" error
		// The service layer's GetTaskByID is expected to wrap this or return a specific error.
		// Here, the mock directly returns an error that the service layer interprets.
		notFoundError := fmt.Errorf("taskRepository.GetTaskByID: no task found with ID %d", taskID)
		mockRepo.GetTaskByIDFunc = func(ctx context.Context, id int64) (*models.Task, error) {
			return nil, notFoundError
		}

		_, err := taskService.GetTaskByID(ctx, taskID)
		if err == nil {
			t.Errorf("Expected an error for not found, got nil")
		}
		// Check if the error message from the service indicates "not found"
		// This depends on how GetTaskByID in taskServiceImpl formats its error.
		// Based on current impl: fmt.Errorf("tarefa com ID %d não encontrada: %w", taskID, err)
		if !strings.Contains(err.Error(), fmt.Sprintf("tarefa com ID %d não encontrada", taskID)) {
			t.Errorf("Expected 'not found' error message, got: %v", err)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks for 'not found' error, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("repository error on get", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		repoError := errors.New("repo get error")
		mockRepo.GetTaskByIDFunc = func(ctx context.Context, id int64) (*models.Task, error) {
			return nil, repoError
		}
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) { return 103, nil }

		_, err := taskService.GetTaskByID(ctx, taskID)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, repoError) { // Service should return the original repo error for unexpected ones
			t.Errorf("Expected error %v, got %v", repoError, err)
		}
		if len(mockRepo.CreatedBugTasks) == 0 {
			t.Errorf("Expected a bug task to be created for unexpected repo error")
		} else if !strings.Contains(mockRepo.CreatedBugTasks[0].Title, "Task Retrieval Failure") {
			t.Errorf("Incorrect bug task title: %s", mockRepo.CreatedBugTasks[0].Title)
		}
	})
}


func TestTaskService_UpdateTask(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()
	taskToUpdate := &models.Task{ID: 1, Title: "Updated Title", UserID: 1}

	t.Run("successful update", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		mockRepo.UpdateTaskFunc = func(ctx context.Context, task *models.Task) error {
			if task.ID == taskToUpdate.ID && task.Title == taskToUpdate.Title {
				return nil
			}
			return fmt.Errorf("unexpected task data in update: %+v", task)
		}

		err := taskService.UpdateTask(ctx, taskToUpdate)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("validation error (empty title)", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		invalidTask := &models.Task{ID: 1, Title: "", UserID: 1}
		err := taskService.UpdateTask(ctx, invalidTask)
		if err == nil {
			t.Errorf("Expected validation error for empty title, got nil")
		}
		if !strings.Contains(err.Error(), "task title cannot be empty") {
			t.Errorf("Expected error message about empty title, got: %v", err)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks for validation error, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("repository error on update (task not found by repo)", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		// This specific error message comes from the repository implementation
		repoNotFoundError := fmt.Errorf("taskRepository.UpdateTask: no task found with ID %d", taskToUpdate.ID)
		mockRepo.UpdateTaskFunc = func(ctx context.Context, task *models.Task) error {
			return repoNotFoundError
		}
		// No bug task should be created for "not found" or "no change" by service
		// CreateTaskFunc should not be called in this scenario by handleErrorAndCreateBugTask
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) {
			// This function should not be called in this test case.
			// If it is, the test will fail.
			return 0, nil
		}


		err := taskService.UpdateTask(ctx, taskToUpdate)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, repoNotFoundError) {
			t.Errorf("Expected error %v, got %v", repoNotFoundError, err)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug task for repo 'not found' error, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("unexpected repository error on update", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		repoError := errors.New("repo update general error")
		mockRepo.UpdateTaskFunc = func(ctx context.Context, task *models.Task) error {
			return repoError
		}
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) { return 104, nil }


		err := taskService.UpdateTask(ctx, taskToUpdate)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, repoError) {
			t.Errorf("Expected error %v, got %v", repoError, err)
		}
		if len(mockRepo.CreatedBugTasks) == 0 {
			t.Errorf("Expected a bug task to be created for unexpected repo error")
		} else if !strings.Contains(mockRepo.CreatedBugTasks[0].Title, "Task Update Failure") {
			t.Errorf("Incorrect bug task title: %s", mockRepo.CreatedBugTasks[0].Title)
		}
	})
}

func TestTaskService_DeleteTask(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()
	taskID := int64(1)

	t.Run("successful delete", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		mockRepo.DeleteTaskFunc = func(ctx context.Context, id int64) error {
			if id == taskID {
				return nil
			}
			return fmt.Errorf("unexpected taskID for delete: %d", id)
		}

		err := taskService.DeleteTask(ctx, taskID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug tasks, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("repository error on delete (task not found by repo)", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		repoNotFoundError := fmt.Errorf("taskRepository.DeleteTask: no task found with ID %d", taskID)
		mockRepo.DeleteTaskFunc = func(ctx context.Context, id int64) error {
			return repoNotFoundError
		}
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) {
			// This function should not be called in this test case.
			return 0, nil
		}

		err := taskService.DeleteTask(ctx, taskID)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, repoNotFoundError) {
			t.Errorf("Expected error %v, got %v", repoNotFoundError, err)
		}
		if len(mockRepo.CreatedBugTasks) != 0 {
			t.Errorf("Expected no bug task for repo 'not found' error during delete, got %d", len(mockRepo.CreatedBugTasks))
		}
	})

	t.Run("unexpected repository error on delete", func(t *testing.T) {
		mockRepo.CreatedBugTasks = []models.Task{}
		repoError := errors.New("repo delete general error")
		mockRepo.DeleteTaskFunc = func(ctx context.Context, id int64) error {
			return repoError
		}
		mockRepo.CreateTaskFunc = func(ctx context.Context, task *models.Task) (int64, error) { return 105, nil }

		err := taskService.DeleteTask(ctx, taskID)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, repoError) {
			t.Errorf("Expected error %v, got %v", repoError, err)
		}
		if len(mockRepo.CreatedBugTasks) == 0 {
			t.Errorf("Expected a bug task to be created for unexpected repo error")
		} else if !strings.Contains(mockRepo.CreatedBugTasks[0].Title, "Task Deletion Failure") {
			t.Errorf("Incorrect bug task title: %s", mockRepo.CreatedBugTasks[0].Title)
		}
	})
}
