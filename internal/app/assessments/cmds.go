package assessments

import (
	"context"
	"fmt"
	"strconv"
	"vigenda/internal/models"
	"vigenda/internal/service"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Messages ---
type fetchedAssessmentsMsg struct {
	assessments []models.Assessment
	err         error
}
type assessmentCreatedMsg struct{ err error }
type assessmentUpdatedMsg struct{ err error }
type assessmentDeletedMsg struct{ err error }
type fetchedGradingSheetMsg struct {
	sheet *service.GradingSheet
	err   error
}
type gradesEnteredMsg struct{ err error }
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// --- Commands ---
func (m *Model) fetchAssessmentsCmd() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
	defer cancel()
	assessments, err := m.assessmentService.ListAllAssessments(ctx)
	return fetchedAssessmentsMsg{assessments: assessments, err: err}
}

func (m *Model) createAssessmentCmd(name string, classID int64, term int, weight float64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		_, err := m.assessmentService.CreateAssessment(ctx, name, classID, term, weight)
		return assessmentCreatedMsg{err: err}
	}
}

func (m *Model) updateAssessmentCmd(id int64, name string, classID int64, term int, weight float64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		_, err := m.assessmentService.UpdateAssessment(ctx, id, name, classID, term, weight)
		return assessmentUpdatedMsg{err: err}
	}
}

func (m *Model) deleteAssessmentCmd(id int64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		err := m.assessmentService.DeleteAssessment(ctx, id)
		return assessmentDeletedMsg{err: err}
	}
}

func (m *Model) fetchGradingSheetCmd(assessmentID int64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		sheet, err := m.assessmentService.GetGradingSheet(ctx, assessmentID)
		return fetchedGradingSheetMsg{sheet: sheet, err: err}
	}
}

func (m *Model) submitGradesCmd() tea.Cmd {
	if m.gradingSheet == nil {
		return func() tea.Msg { return errMsg{err: fmt.Errorf("grading sheet not loaded")} }
	}

	grades := make(map[int64]float64)
	for studentID, ti := range m.gradeInputs {
		gradeStr := ti.Value()
		if gradeStr == "" {
			// Check if a grade previously existed, if so, we might need to delete it.
			// For simplicity, we'll just ignore empty inputs.
			// A value of -1 could signify deletion.
			continue
		}
		grade, err := strconv.ParseFloat(gradeStr, 64)
		if err != nil {
			return func() tea.Msg { return errMsg{err: fmt.Errorf("nota inválida para aluno ID %d: '%s'", studentID, gradeStr)} }
		}
		grades[studentID] = grade
	}

	if len(grades) == 0 {
		return nil // Nothing to submit
	}

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
		defer cancel()
		err := m.assessmentService.EnterGrades(ctx, m.gradingSheet.Assessment.ID, grades)
		return gradesEnteredMsg{err: err}
	}
}
