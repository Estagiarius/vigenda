package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"vigenda/internal/models"
	"vigenda/internal/repository/stubs" // Using generated mock

	"go.uber.org/mock/gomock"
)

func TestClassServiceImpl_CreateClass(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	// No need for mockSubjectRepo yet, as it's not used in CreateClass logic directly for now

	classService := NewClassService(mockClassRepo, nil) // Pass nil for subjectRepo if not used

	ctx := context.Background()
	className := "Test Class"
	subjectID := int64(1)
	userID := int64(1) // Assuming UserID 1 for tests, consistent with service logic

	expectedClass := models.Class{
		ID:        1,
		UserID:    userID,
		SubjectID: subjectID,
		Name:      className,
		CreatedAt: time.Now(), // Approximate, repo layer sets this
		UpdatedAt: time.Now(), // Approximate, repo layer sets this
	}

	mockClassRepo.EXPECT().
		CreateClass(gomock.Any(), gomock.Any()). // gomock.Any() for context and class pointer
		DoAndReturn(func(ctx context.Context, class *models.Class) (int64, error) {
			if class.Name != className || class.SubjectID != subjectID || class.UserID != userID {
				return 0, errors.New("input class data mismatch")
			}
			return expectedClass.ID, nil // Return the ID of the created class
		}).Times(1)

	mockClassRepo.EXPECT().
		GetClassByID(gomock.Any(), expectedClass.ID).
		Return(&expectedClass, nil).Times(1)

	createdClass, err := classService.CreateClass(ctx, className, subjectID)
	if err != nil {
		t.Fatalf("CreateClass failed: %v", err)
	}

	if createdClass.Name != expectedClass.Name || createdClass.SubjectID != expectedClass.SubjectID || createdClass.ID != expectedClass.ID {
		t.Errorf("CreateClass returned %+v, want %+v", createdClass, expectedClass)
	}
}

func TestClassServiceImpl_UpdateClass(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	classService := NewClassService(mockClassRepo, nil)

	ctx := context.Background()
	classID := int64(1)
	newName := "Updated Test Class"
	newSubjectID := int64(2)
	userID := int64(1)

	originalClass := models.Class{ID: classID, UserID: userID, Name: "Old Name", SubjectID: int64(1)}
	updatedClassModel := models.Class{ID: classID, UserID: userID, Name: newName, SubjectID: newSubjectID}

	mockClassRepo.EXPECT().GetClassByID(ctx, classID).Return(&originalClass, nil).Times(1)
	mockClassRepo.EXPECT().UpdateClass(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, class *models.Class) error {
			if class.ID != classID || class.Name != newName || class.SubjectID != newSubjectID {
				return errors.New("update data mismatch")
			}
			return nil
		}).Times(1)
	mockClassRepo.EXPECT().GetClassByID(ctx, classID).Return(&updatedClassModel, nil).Times(1)

	updatedClass, err := classService.UpdateClass(ctx, classID, newName, newSubjectID)
	if err != nil {
		t.Fatalf("UpdateClass failed: %v", err)
	}
	if updatedClass.Name != newName || updatedClass.SubjectID != newSubjectID {
		t.Errorf("UpdateClass returned wrong data: got %+v", updatedClass)
	}
}

func TestClassServiceImpl_DeleteClass(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	classService := NewClassService(mockClassRepo, nil)

	ctx := context.Background()
	classID := int64(1)
	userID := int64(1) // Assumed from context/auth in service

	// Mock GetClassByID to simulate class existence check
	mockClassRepo.EXPECT().GetClassByID(ctx, classID).Return(&models.Class{ID: classID, UserID: userID}, nil).Times(1)
	// Mock DeleteClass
	mockClassRepo.EXPECT().DeleteClass(ctx, classID, userID).Return(nil).Times(1)

	err := classService.DeleteClass(ctx, classID)
	if err != nil {
		t.Fatalf("DeleteClass failed: %v", err)
	}
}

func TestClassServiceImpl_AddStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	classService := NewClassService(mockClassRepo, nil)

	ctx := context.Background()
	classID := int64(1)
	fullName := "Test Student"
	enrollmentID := "TS001"
	status := "ativo"
	studentID := int64(10) // Mocked ID for the new student

	expectedStudent := models.Student{
		ID:           studentID,
		ClassID:      classID,
		FullName:     fullName,
		EnrollmentID: enrollmentID,
		Status:       status,
	}

	mockClassRepo.EXPECT().AddStudent(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, student *models.Student) (int64, error) {
			if student.ClassID != classID || student.FullName != fullName {
				return 0, errors.New("student data mismatch for add")
			}
			return studentID, nil // Return new student ID
		}).Times(1)

	mockClassRepo.EXPECT().GetStudentByID(ctx, studentID).Return(&expectedStudent, nil).Times(1)

	addedStudent, err := classService.AddStudent(ctx, classID, fullName, enrollmentID, status)
	if err != nil {
		t.Fatalf("AddStudent failed: %v", err)
	}
	if addedStudent.ID != studentID || addedStudent.FullName != fullName {
		t.Errorf("AddStudent returned %+v, want ID %d and Name %s", addedStudent, studentID, fullName)
	}
}

func TestClassServiceImpl_UpdateStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	classService := NewClassService(mockClassRepo, nil)

	ctx := context.Background()
	studentID := int64(1)
	classID := int64(5) // Student belongs to this class
	newFullName := "Updated Student Name"
	newEnrollmentID := "US002"
	newStatus := "inativo"

	originalStudent := models.Student{ID: studentID, ClassID: classID, FullName: "Old Name", Status: "ativo"}
	updatedStudentModel := models.Student{ID: studentID, ClassID: classID, FullName: newFullName, EnrollmentID: newEnrollmentID, Status: newStatus}

	mockClassRepo.EXPECT().GetStudentByID(ctx, studentID).Return(&originalStudent, nil).Times(1)
	mockClassRepo.EXPECT().UpdateStudent(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, student *models.Student) error {
			if student.ID != studentID || student.FullName != newFullName || student.Status != newStatus {
				return errors.New("student update data mismatch")
			}
			return nil
		}).Times(1)
	mockClassRepo.EXPECT().GetStudentByID(ctx, studentID).Return(&updatedStudentModel, nil).Times(1)

	updatedStudent, err := classService.UpdateStudent(ctx, studentID, newFullName, newEnrollmentID, newStatus)
	if err != nil {
		t.Fatalf("UpdateStudent failed: %v", err)
	}
	if updatedStudent.FullName != newFullName || updatedStudent.Status != newStatus {
		t.Errorf("UpdateStudent returned wrong data: got %+v", updatedStudent)
	}
}

func TestClassServiceImpl_DeleteStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	classService := NewClassService(mockClassRepo, nil)

	ctx := context.Background()
	studentID := int64(1)
	classID := int64(5) // Student belongs to this class

	// Mock GetStudentByID to simulate student existence and get its classID
	mockClassRepo.EXPECT().GetStudentByID(ctx, studentID).Return(&models.Student{ID: studentID, ClassID: classID}, nil).Times(1)
	// Mock DeleteStudent
	mockClassRepo.EXPECT().DeleteStudent(ctx, studentID, classID).Return(nil).Times(1)

	err := classService.DeleteStudent(ctx, studentID)
	if err != nil {
		t.Fatalf("DeleteStudent failed: %v", err)
	}
}

func TestClassServiceImpl_ListAllClasses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	classService := NewClassService(mockClassRepo, nil)
	ctx := context.Background()

	expectedClasses := []models.Class{
		{ID: 1, Name: "Class A", SubjectID: 101, UserID: 1},
		{ID: 2, Name: "Class B", SubjectID: 102, UserID: 1},
	}
	mockClassRepo.EXPECT().ListAllClasses(ctx).Return(expectedClasses, nil).Times(1)

	classes, err := classService.ListAllClasses(ctx)
	if err != nil {
		t.Fatalf("ListAllClasses failed: %v", err)
	}
	if len(classes) != len(expectedClasses) {
		t.Errorf("ListAllClasses returned %d classes, want %d", len(classes), len(expectedClasses))
	}
}

func TestClassServiceImpl_GetStudentsByClassID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	classService := NewClassService(mockClassRepo, nil)
	ctx := context.Background()
	classID := int64(1)

	expectedStudents := []models.Student{
		{ID: 1, FullName: "Student 1", ClassID: classID, Status: "ativo"},
		{ID: 2, FullName: "Student 2", ClassID: classID, Status: "inativo"},
	}
	mockClassRepo.EXPECT().GetStudentsByClassID(ctx, classID).Return(expectedStudents, nil).Times(1)

	students, err := classService.GetStudentsByClassID(ctx, classID)
	if err != nil {
		t.Fatalf("GetStudentsByClassID failed: %v", err)
	}
	if len(students) != len(expectedStudents) {
		t.Errorf("GetStudentsByClassID returned %d students, want %d", len(students), len(expectedStudents))
	}
}

func TestClassServiceImpl_ImportStudentsFromCSV(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	classService := NewClassService(mockClassRepo, nil)

	ctx := context.Background()
	classID := int64(1)
	csvData := `enrollment_id,full_name,status
S001,Student One,ativo
S002,Student Two,inativo
,Student Three,
S004,Student Four,transferido
`

	// Expect AddStudent to be called for each valid row
	mockClassRepo.EXPECT().AddStudent(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, s *models.Student) (int64, error) {
			return 1, nil
		}).Times(4)
	mockClassRepo.EXPECT().GetStudentByID(ctx, gomock.Any()).Return(&models.Student{}, nil).AnyTimes()

	importedCount, err := classService.ImportStudentsFromCSV(ctx, classID, []byte(csvData))
	if err != nil {
		t.Fatalf("ImportStudentsFromCSV failed: %v", err)
	}

	if importedCount != 4 {
		t.Errorf("ImportStudentsFromCSV returned count %d, want %d", importedCount, 4)
	}
}

func TestUpdateStudentStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	classService := NewClassService(mockClassRepo, nil)

	ctx := context.Background()
	studentID := int64(1)
	newStatus := "transferido"

	// Mock UpdateStudentStatus in repository
	mockClassRepo.EXPECT().UpdateStudentStatus(ctx, studentID, newStatus).Return(nil).Times(1)

	err := classService.UpdateStudentStatus(ctx, studentID, newStatus)
	if err != nil {
		t.Fatalf("UpdateStudentStatus failed: %v", err)
	}
}
