package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log" // Adicionado para logging
	"strings"
	"vigenda/internal/models"
	"vigenda/internal/repository"
)

type classServiceImpl struct {
	classRepo   repository.ClassRepository
	subjectRepo repository.SubjectRepository // Added subjectRepo if needed for validation or other logic
}

// NewClassService creates a new instance of ClassService.
// It now accepts ClassRepository and SubjectRepository as dependencies.
func NewClassService(
	classRepo repository.ClassRepository,
	subjectRepo repository.SubjectRepository,
) ClassService {
	return &classServiceImpl{
		classRepo:   classRepo,
		subjectRepo: subjectRepo,
	}
}

func (s *classServiceImpl) CreateClass(ctx context.Context, name string, subjectID int64) (models.Class, error) {
	if name == "" {
		return models.Class{}, fmt.Errorf("class name cannot be empty")
	}
	if subjectID <= 0 {
		return models.Class{}, fmt.Errorf("subject ID must be positive")
	}

	// TODO: Validate if subjectID exists using subjectRepo if necessary.
	// For now, we assume subjectID is valid.

	// Assuming UserID 1 for now, this should come from context or auth
	// This should ideally be retrieved from the context after authentication/authorization.
	userID := int64(1) // Placeholder for actual User ID from context

	class := models.Class{
		UserID:    userID,
		SubjectID: subjectID,
		Name:      name,
	}

	id, err := s.classRepo.CreateClass(ctx, &class)
	if err != nil {
		return models.Class{}, fmt.Errorf("service.CreateClass: failed to create class: %w", err)
	}
	class.ID = id
	// To get CreatedAt and UpdatedAt, we should ideally fetch the created record.
	// However, for simplicity, we can assume the repository sets these.
	// If not, a GetClassByID call would be needed here.
	// For now, we'll rely on the repository to have populated these if possible,
	// or accept that they might be zero if CreateClass in repo doesn't return them.
	// The model.go changes ensure the fields exist.
	// The repository.go changes ensure they are set on insert.
	// So, a fresh GetClassByID might be redundant if CreateClass returns the full object or sets it.
	// Let's assume the ID is sufficient and the caller might re-fetch if they need timestamps immediately.
	// Or, even better, modify CreateClass in repo to return the full models.Class object.
	// For now, let's fetch it to ensure all fields are populated.
	createdClass, err := s.classRepo.GetClassByID(ctx, id)
	if err != nil {
		return models.Class{}, fmt.Errorf("service.CreateClass: failed to retrieve created class with id %d: %w", id, err)
	}
	return *createdClass, nil
}

func (s *classServiceImpl) UpdateClass(ctx context.Context, classID int64, name string, subjectID int64) (models.Class, error) {
	if classID <= 0 {
		return models.Class{}, fmt.Errorf("class ID must be positive")
	}
	if name == "" {
		return models.Class{}, fmt.Errorf("class name cannot be empty")
	}
	if subjectID <= 0 {
		return models.Class{}, fmt.Errorf("subject ID must be positive")
	}

	// Assuming UserID 1 for now, this should come from context or auth
	userID := int64(1) // Placeholder for actual User ID from context

	// Fetch the existing class to ensure it belongs to the user
	classToUpdate, err := s.classRepo.GetClassByID(ctx, classID)
	if err != nil {
		return models.Class{}, fmt.Errorf("service.UpdateClass: failed to get class: %w", err)
	}
	if classToUpdate.UserID != userID {
		return models.Class{}, fmt.Errorf("service.UpdateClass: class does not belong to user") // Or ErrForbidden
	}

	classToUpdate.Name = name
	classToUpdate.SubjectID = subjectID
	// UserID remains the same

	err = s.classRepo.UpdateClass(ctx, classToUpdate)
	if err != nil {
		return models.Class{}, fmt.Errorf("service.UpdateClass: failed to update class: %w", err)
	}

	updatedClass, err := s.classRepo.GetClassByID(ctx, classID)
	if err != nil {
		return models.Class{}, fmt.Errorf("service.UpdateClass: failed to retrieve updated class: %w", err)
	}
	return *updatedClass, nil
}

func (s *classServiceImpl) DeleteClass(ctx context.Context, classID int64) error {
	if classID <= 0 {
		return fmt.Errorf("class ID must be positive")
	}
	// Assuming UserID 1 for now
	userID := int64(1) // Placeholder

	// Optional: Check if class exists and belongs to user before deleting
	_, err := s.classRepo.GetClassByID(ctx, classID) // Check existence
	if err != nil {
		return fmt.Errorf("service.DeleteClass: failed to get class or class not found: %w", err)
	}
	// The repository's DeleteClass should handle the userID check for ownership.

	err = s.classRepo.DeleteClass(ctx, classID, userID)
	if err != nil {
		return fmt.Errorf("service.DeleteClass: failed to delete class: %w", err)
	}
	return nil
}

func (s *classServiceImpl) AddStudent(ctx context.Context, classID int64, fullName string, enrollmentID string, status string) (models.Student, error) {
	if classID <= 0 {
		return models.Student{}, fmt.Errorf("class ID must be positive")
	}
	if fullName == "" {
		return models.Student{}, fmt.Errorf("student full name cannot be empty")
	}
	if status == "" {
		status = "ativo" // Default status
	} else {
		status = strings.ToLower(status)
		if status != "ativo" && status != "inativo" && status != "transferido" {
			return models.Student{}, fmt.Errorf("invalid student status: %s. Allowed: 'ativo', 'inativo', 'transferido'", status)
		}
	}

	// Optional: Validate if classID exists and belongs to the user
	// _, err := s.GetClassByID(ctx, classID) // Assuming UserID check is within GetClassByID or not needed here
	// if err != nil {
	// 	return models.Student{}, fmt.Errorf("service.AddStudent: class not found or not accessible: %w", err)
	// }


	student := models.Student{
		ClassID:      classID,
		FullName:     fullName,
		EnrollmentID: enrollmentID,
		Status:       status,
	}

	studentID, err := s.classRepo.AddStudent(ctx, &student)
	if err != nil {
		return models.Student{}, fmt.Errorf("service.AddStudent: failed to add student: %w", err)
	}
	// student.ID = studentID // Set ID for the returned object

	// Fetch the newly created student to get all fields, including ID, CreatedAt, UpdatedAt
	newStudent, err := s.classRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		return models.Student{}, fmt.Errorf("service.AddStudent: failed to retrieve newly added student: %w", err)
	}

	return *newStudent, nil
}

func (s *classServiceImpl) GetStudentByID(ctx context.Context, studentID int64) (models.Student, error) {
	if studentID <= 0 {
		return models.Student{}, fmt.Errorf("student ID must be positive")
	}
	student, err := s.classRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		return models.Student{}, fmt.Errorf("service.GetStudentByID: %w", err)
	}
	// TODO: Add ownership check if necessary (e.g., student's class belongs to user)
	return *student, nil
}

func (s *classServiceImpl) UpdateStudent(ctx context.Context, studentID int64, fullName string, enrollmentID string, status string) (models.Student, error) {
	if studentID <= 0 {
		return models.Student{}, fmt.Errorf("student ID must be positive")
	}
	if fullName == "" {
		return models.Student{}, fmt.Errorf("student full name cannot be empty")
	}
	if status == "" {
		return models.Student{}, fmt.Errorf("student status cannot be empty")
	}
	status = strings.ToLower(status)
	if status != "ativo" && status != "inativo" && status != "transferido" {
		return models.Student{}, fmt.Errorf("invalid student status: %s", status)
	}

	// Fetch existing student to check ownership (e.g., via class) and to get ClassID
	studentToUpdate, err := s.classRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		return models.Student{}, fmt.Errorf("service.UpdateStudent: failed to get student: %w", err)
	}

	// Optional: Check if the class of the student belongs to the current user
	// _, err = s.GetClassByID(ctx, studentToUpdate.ClassID) // This would involve userID check
	// if err != nil {
	// 	 return models.Student{}, fmt.Errorf("service.UpdateStudent: class not found or not accessible: %w", err)
	// }

	studentToUpdate.FullName = fullName
	studentToUpdate.EnrollmentID = enrollmentID
	studentToUpdate.Status = status
	// ClassID should not be changed here. If student needs to move class, that's a different operation.

	err = s.classRepo.UpdateStudent(ctx, studentToUpdate)
	if err != nil {
		return models.Student{}, fmt.Errorf("service.UpdateStudent: failed to update student: %w", err)
	}

	updatedStudent, err := s.classRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		return models.Student{}, fmt.Errorf("service.UpdateStudent: failed to retrieve updated student: %w", err)
	}
	return *updatedStudent, nil
}

func (s *classServiceImpl) DeleteStudent(ctx context.Context, studentID int64) error {
	if studentID <= 0 {
		return fmt.Errorf("student ID must be positive")
	}

	// Fetch student to get ClassID for repository delete call (which needs classID for ownership check)
	student, err := s.classRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		return fmt.Errorf("service.DeleteStudent: failed to get student or student not found: %w", err)
	}

	// Optional: Verify class ownership by the current user
	// _, err = s.GetClassByID(ctx, student.ClassID) // This involves userID
	// if err != nil {
	// 	return fmt.Errorf("service.DeleteStudent: class not found or not accessible: %w", err)
	// }
	// The repository's DeleteStudent method takes classID, it might perform an ownership check
	// or the service ensures the user owns the class before calling.
	// For now, we assume the repo layer might need classID for its own reasons (like the current FK in delete).

	err = s.classRepo.DeleteStudent(ctx, studentID, student.ClassID)
	if err != nil {
		return fmt.Errorf("service.DeleteStudent: failed to delete student: %w", err)
	}
	return nil
}

func (s *classServiceImpl) ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error) {
	if classID <= 0 {
		return 0, fmt.Errorf("class ID must be positive")
	}
	// TODO: Validate if classID exists and belongs to the user
	// _, err := s.GetClassByID(ctx, classID)
	// if err != nil {
	// 	return 0, fmt.Errorf("service.ImportStudentsFromCSV: class not found or not accessible: %w", err)
	// }

	reader := csv.NewReader(strings.NewReader(string(csvData)))
	// Skip header
	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return 0, fmt.Errorf("CSV is empty or header-only")
		}
		return 0, fmt.Errorf("failed to read CSV header: %w", err)
	}

	var importedCount int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV record, processed %d students: %v", importedCount, err)
			return importedCount, fmt.Errorf("error reading CSV record: %w", err)
		}

		if len(record) < 2 || strings.TrimSpace(record[1]) == "" {
			log.Printf("Skipping invalid CSV record (missing full name): %v", record)
			continue
		}

		enrollmentID := strings.TrimSpace(record[0]) // nº chamada / enrollment_id
		fullName := strings.TrimSpace(record[1])
		status := "ativo" // Default
		if len(record) > 2 && strings.TrimSpace(record[2]) != "" {
			status = strings.ToLower(strings.TrimSpace(record[2]))
			if status != "ativo" && status != "inativo" && status != "transferido" {
				log.Printf("Skipping invalid CSV record (invalid status '%s'): %v", status, record)
				continue
			}
		}
		// Call AddStudent for each valid record
		_, err = s.AddStudent(ctx, classID, fullName, enrollmentID, status)
		if err != nil {
			// Log and continue, or stop? For now, log and attempt to continue.
			log.Printf("Failed to add student '%s' from CSV: %v. Continuing...", fullName, err)
			// If strict import is needed, return error here:
			// return importedCount, fmt.Errorf("failed to import student '%s': %w", fullName, err)
		} else {
			importedCount++
		}
	}
	return importedCount, nil
}

func (s *classServiceImpl) UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error {
	if studentID <= 0 {
		return fmt.Errorf("student ID must be positive")
	}
	newStatus = strings.ToLower(strings.TrimSpace(newStatus))
	if newStatus != "ativo" && newStatus != "inativo" && newStatus != "transferido" {
		return fmt.Errorf("invalid student status: '%s'. Allowed: 'ativo', 'inativo', 'transferido'", newStatus)
	}

	// Optional: Check ownership of student's class
	// student, err := s.GetStudentByID(ctx, studentID)
	// if err != nil {
	// 	return fmt.Errorf("service.UpdateStudentStatus: student not found: %w", err)
	// }
	// _, err = s.GetClassByID(ctx, student.ClassID) // UserID check
	// if err != nil {
	// 	return fmt.Errorf("service.UpdateStudentStatus: class not found or not accessible: %w", err)
	// }


	err := s.classRepo.UpdateStudentStatus(ctx, studentID, newStatus)
	if err != nil {
		return fmt.Errorf("service.UpdateStudentStatus: %w", err)
	}
	return nil
}

func (s *classServiceImpl) GetClassByID(ctx context.Context, classID int64) (models.Class, error) {
	if classID <= 0 {
		return models.Class{}, fmt.Errorf("class ID must be positive")
	}
	// Assuming UserID 1 for now
	// userID := int64(1) // Placeholder

	class, err := s.classRepo.GetClassByID(ctx, classID)
	if err != nil {
		return models.Class{}, fmt.Errorf("service.GetClassByID: %w", err)
	}
	// if class.UserID != userID {
	// 	 return models.Class{}, models.ErrClassNotFound // Or a specific "forbidden" error
	// }
	return *class, nil
}

func (s *classServiceImpl) ListAllClasses(ctx context.Context) ([]models.Class, error) {
	log.Println("Service: classServiceImpl.ListAllClasses - Chamado.")
	// Assuming UserID 1 for now
	// userID := int64(1) // Placeholder
	// The repository method ListAllClasses currently does not filter by userID.
	// If it needs to, the method signature and implementation in repository should change,
	// or filter here after fetching all. For now, it lists all classes in the system.
	// If user-specific listing is needed, this needs adjustment.
	// For now, let's assume it's listing all classes the current user *could* see,
	// or it's an admin-like function. The TUI implies user-specific.
	// The repository ListAllClasses currently doesn't take userID.
	// This implies that either:
	// 1. The TUI will filter based on a UserID known to it (less likely for a service method).
	// 2. The repository needs to be updated to filter by UserID.
	// 3. The service layer should fetch all and then filter if UserID is available (inefficient).
	// Given the existing repo, ListAllClasses returns ALL classes.
	// The CreateClass sets UserID=1. So, for now, this will list all classes,
	// and if only UserID=1 creates classes, it effectively lists classes for UserID=1.

	log.Println("Service: classServiceImpl.ListAllClasses - Chamando repositório para listar turmas.")
	classes, err := s.classRepo.ListAllClasses(ctx)
	if err != nil {
		log.Printf("Service: classServiceImpl.ListAllClasses - Erro ao listar turmas do repositório: %v", err)
		return nil, fmt.Errorf("service.ListAllClasses: %w", err)
	}
	log.Printf("Service: classServiceImpl.ListAllClasses - Repositório retornou %d turmas.", len(classes))
	return classes, nil
}

func (s *classServiceImpl) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
	if classID <= 0 {
		return nil, fmt.Errorf("class ID must be positive when fetching students")
	}
	// Optional: Check if classID exists and belongs to the user
	// _, err := s.GetClassByID(ctx, classID) // This involves userID check
	// if err != nil {
	// 	return nil, fmt.Errorf("service.GetStudentsByClassID: class not found or not accessible: %w", err)
	// }

	students, err := s.classRepo.GetStudentsByClassID(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("service.GetStudentsByClassID: failed to get students: %w", err)
	}
	return students, nil
}
