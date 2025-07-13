package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"vigenda/internal/models"
	"vigenda/internal/repository/stubs"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAssessmentServiceImpl_CreateAssessment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAssessmentRepo := stubs.NewMockAssessmentRepository(ctrl)
	service := NewAssessmentService(mockAssessmentRepo, nil)

	ctx := context.Background()
	assessment := &models.Assessment{
		Name:    "Test Assessment",
		ClassID: 1,
		Term:    1,
		Weight:  2.5,
	}

	mockAssessmentRepo.EXPECT().
		CreateAssessment(ctx, gomock.Any()).
		Return(int64(1), nil).
		Times(1)

	created, err := service.CreateAssessment(ctx, assessment.Name, assessment.ClassID, assessment.Term, assessment.Weight)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), created.ID)
	assert.Equal(t, assessment.Name, created.Name)
}

func TestAssessmentServiceImpl_EnterGrades(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAssessmentRepo := stubs.NewMockAssessmentRepository(ctrl)
	service := NewAssessmentService(mockAssessmentRepo, nil)

	ctx := context.Background()
	assessmentID := int64(1)
	studentGrades := map[int64]float64{
		10: 8.5,
		11: 9.0,
	}

	mockAssessmentRepo.EXPECT().
		GetAssessmentByID(ctx, assessmentID).
		Return(&models.Assessment{ID: assessmentID, ClassID: 1}, nil).
		Times(1)

	mockAssessmentRepo.EXPECT().
		EnterGrade(ctx, gomock.Any()).
		Return(nil).
		Times(len(studentGrades))

	err := service.EnterGrades(ctx, assessmentID, studentGrades)
	assert.NoError(t, err)
}

func TestAssessmentServiceImpl_GetGradingSheet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAssessmentRepo := stubs.NewMockAssessmentRepository(ctrl)
	mockClassRepo := stubs.NewMockClassRepository(ctrl)
	service := NewAssessmentService(mockAssessmentRepo, mockClassRepo)

	ctx := context.Background()
	assessmentID := int64(1)
	classID := int64(5)

	mockAssessment := models.Assessment{ID: assessmentID, Name: "Final Exam", ClassID: classID}
	mockStudents := []models.Student{
		{ID: 101, FullName: "Alice", ClassID: classID, Status: "ativo"},
		{ID: 102, FullName: "Bob", ClassID: classID, Status: "ativo"},
	}
	mockGrades := []models.Grade{
		{StudentID: 101, AssessmentID: assessmentID, Grade: 88.0},
	}

	mockAssessmentRepo.EXPECT().GetAssessmentByID(ctx, assessmentID).Return(&mockAssessment, nil)
	mockClassRepo.EXPECT().GetStudentsByClassID(ctx, classID).Return(mockStudents, nil)
	mockAssessmentRepo.EXPECT().GetGradesByAssessmentID(ctx, assessmentID).Return(mockGrades, nil)

	sheet, err := service.GetGradingSheet(ctx, assessmentID)

	assert.NoError(t, err)
	assert.NotNil(t, sheet)
	assert.Equal(t, "Final Exam", sheet.Assessment.Name)
	assert.Len(t, sheet.Students, 2)
	assert.Len(t, sheet.Grades, 1)
	assert.Equal(t, 88.0, sheet.Grades[101].Grade)
	_, exists := sheet.Grades[102]
	assert.False(t, exists, "Bob should not have a grade yet")
}

func TestAssessmentServiceImpl_CalculateClassAverage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAssessmentRepo := stubs.NewMockAssessmentRepository(ctrl)
	service := NewAssessmentService(mockAssessmentRepo, nil)

	ctx := context.Background()
	classID := int64(1)

	mockAssessments := []models.Assessment{
		{ID: 1, ClassID: classID, Name: "P1", Weight: 4.0},
		{ID: 2, ClassID: classID, Name: "T1", Weight: 6.0},
	}
	mockStudents := []models.Student{
		{ID: 10, ClassID: classID, FullName: "Alice", Status: "ativo"},
		{ID: 11, ClassID: classID, FullName: "Bob", Status: "ativo"},
		{ID: 12, ClassID: classID, FullName: "Charlie", Status: "inativo"},
	}
	mockGrades := []models.Grade{
		// Alice: (10 * 4 + 8 * 6) / 10 = (40 + 48) / 10 = 8.8
		{AssessmentID: 1, StudentID: 10, Grade: 10},
		{AssessmentID: 2, StudentID: 10, Grade: 8},
		// Bob: (7 * 4 + 9 * 6) / 10 = (28 + 54) / 10 = 8.2
		{AssessmentID: 1, StudentID: 11, Grade: 7},
		{AssessmentID: 2, StudentID: 11, Grade: 9},
		// Charlie (inactive)
		{AssessmentID: 1, StudentID: 12, Grade: 5},
	}

	mockAssessmentRepo.EXPECT().
		GetGradesByClassID(ctx, classID).
		Return(mockGrades, mockAssessments, mockStudents, nil).
		Times(1)

	// Expected: (8.8 + 8.2) / 2 = 17.0 / 2 = 8.5
	expectedAverage := 8.5
	average, err := service.CalculateClassAverage(ctx, classID)

	assert.NoError(t, err)
	assert.InDelta(t, expectedAverage, average, 0.01)
}

func TestAssessmentServiceImpl_CalculateClassAverage_NoActiveStudents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAssessmentRepo := stubs.NewMockAssessmentRepository(ctrl)
	service := NewAssessmentService(mockAssessmentRepo, nil)
	ctx := context.Background()
	classID := int64(1)

	mockAssessments := []models.Assessment{{ID: 1, ClassID: classID, Weight: 1.0}}
	mockStudents := []models.Student{{ID: 1, ClassID: classID, Status: "inativo"}}
	mockGrades := []models.Grade{{AssessmentID: 1, StudentID: 1, Grade: 10}}

	mockAssessmentRepo.EXPECT().GetGradesByClassID(ctx, classID).Return(mockGrades, mockAssessments, mockStudents, nil)

	avg, err := service.CalculateClassAverage(ctx, classID)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, avg)
}

func TestAssessmentServiceImpl_CalculateClassAverage_NoAssessments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAssessmentRepo := stubs.NewMockAssessmentRepository(ctrl)
	service := NewAssessmentService(mockAssessmentRepo, nil)
	ctx := context.Background()
	classID := int64(1)

	mockStudents := []models.Student{{ID: 1, ClassID: classID, Status: "ativo"}}

	mockAssessmentRepo.EXPECT().GetGradesByClassID(ctx, classID).Return([]models.Grade{}, []models.Assessment{}, mockStudents, nil)

	avg, err := service.CalculateClassAverage(ctx, classID)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, avg)
}
