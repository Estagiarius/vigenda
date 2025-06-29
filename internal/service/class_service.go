package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"vigenda/internal/models"
	"vigenda/internal/repository" // Added import
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
	// Basic validation
	if name == "" {
		return models.Class{}, fmt.Errorf("class name cannot be empty")
	}
	if subjectID == 0 { // Assuming 0 is not a valid subject ID
		return models.Class{}, fmt.Errorf("subject ID cannot be zero")
	}

	// TODO: Validate if subjectID exists using subjectRepo if necessary.
	// For now, we assume subjectID is valid.

	// Assuming UserID 1 for now, this should come from context or auth
	userID := int64(1)

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
	return class, nil
}

func (s *classServiceImpl) ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error) {
	if classID == 0 {
		return 0, fmt.Errorf("class ID cannot be zero")
	}
	// TODO: Validate if classID exists

	reader := csv.NewReader(strings.NewReader(string(csvData)))
	// Skip header row
	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return 0, fmt.Errorf("CSV file is empty or only contains a header")
		}
		return 0, fmt.Errorf("failed to read CSV header: %w", err)
	}

	var importedCount int
	// Assuming UserID 1 for now, this should come from context or auth
	// userID := int64(1) // Removed as it's not used

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Log error and continue? Or stop? For now, stop.
			return importedCount, fmt.Errorf("error reading CSV record: %w", err)
		}

		if len(record) < 2 || record[1] == "" { // nome_completo is mandatory
			// Log or skip invalid record
			fmt.Printf("Skipping invalid record: %v (missing full name)\n", record)
			continue
		}

		var callNumber int
		if record[0] != "" {
			cn, errCN := strconv.Atoi(record[0])
			if errCN != nil {
				fmt.Printf("Skipping invalid record: %v (invalid call number: %s)\n", record, record[0])
				continue
			}
			callNumber = cn
		}

		fullName := record[1]
		status := "ativo" // Default status
		if len(record) > 2 && record[2] != "" {
			status = strings.ToLower(record[2])
			// Basic validation for status
			if status != "ativo" && status != "inativo" && status != "transferido" {
				fmt.Printf("Skipping invalid record: %v (invalid status: %s)\n", record, record[2])
				continue
			}
		}

		student := models.Student{
			ClassID:      classID,
			// UserID:     userID, // Removed: UserID is not part of models.Student
			EnrollmentID: strconv.Itoa(callNumber), // Converted callNumber to string for EnrollmentID
			FullName:     fullName,
			Status:       status,
		}

		_, err = s.classRepo.AddStudent(ctx, &student)
		if err != nil {
			// Log error and potentially continue, or return error immediately
			fmt.Printf("Failed to add student '%s': %v. Continuing with next records.\n", fullName, err)
			// return importedCount, fmt.Errorf("failed to add student '%s': %w", fullName, err)
		} else {
			importedCount++
		}
	}

	return importedCount, nil
}

func (s *classServiceImpl) UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error {
	if studentID == 0 {
		return fmt.Errorf("student ID cannot be zero")
	}
	newStatus = strings.ToLower(newStatus)
	if newStatus != "ativo" && newStatus != "inativo" && newStatus != "transferido" {
		return fmt.Errorf("invalid student status: %s. Allowed values are 'ativo', 'inativo', 'transferido'", newStatus)
	}

	err := s.classRepo.UpdateStudentStatus(ctx, studentID, newStatus)
	if err != nil {
		return fmt.Errorf("service.UpdateStudentStatus: %w", err)
	}
	return nil
}

func (s *classServiceImpl) GetClassByID(ctx context.Context, classID int64) (models.Class, error) {
	if classID == 0 {
		return models.Class{}, fmt.Errorf("class ID cannot be zero")
	}
	class, err := s.classRepo.GetClassByID(ctx, classID)
	if err != nil {
		// Wrap error for context, or return a specific service-level error
		return models.Class{}, fmt.Errorf("service.GetClassByID: %w", err)
	}
	if class == nil { // Should be handled by repo returning sql.ErrNoRows, which gets wrapped
		return models.Class{}, models.ErrClassNotFound // Or a more specific service error
	}
	return *class, nil
}

func (s *classServiceImpl) ListAllClasses(ctx context.Context) ([]models.Class, error) {
	fmt.Println("[LOG ClassService] ListAllClasses(): called")
	classes, err := s.classRepo.ListAllClasses(ctx)
	if err != nil {
		fmt.Printf("[LOG ClassService] ListAllClasses(): error from classRepo.ListAllClasses: %v\n", err)
		return nil, fmt.Errorf("service.ListAllClasses: %w", err)
	}
	fmt.Printf("[LOG ClassService] ListAllClasses(): success, returning %d classes\n", len(classes))
	return classes, nil
}

func (s *classServiceImpl) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
	if classID == 0 {
		return nil, fmt.Errorf("class ID cannot be zero when fetching students")
	}
	// TODO: Adicionar validação para verificar se a turma (classID) realmente existe, se necessário.
	//       Isso pode envolver uma chamada a s.classRepo.GetClassByID(ctx, classID) primeiro.
	//       Por enquanto, vamos assumir que o ID da turma é válido se não for zero.

	students, err := s.classRepo.GetStudentsByClassID(ctx, classID)
	if err != nil {
		// Não é necessário verificar sql.ErrNoRows aqui, pois o repositório já o trata
		// e retorna uma lista vazia se nenhum aluno for encontrado, o que é um resultado válido.
		return nil, fmt.Errorf("service.GetStudentsByClassID: failed to get students: %w", err)
	}
	// Se a lista estiver vazia, isso é um resultado válido (turma sem alunos).
	return students, nil
}
