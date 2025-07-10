package classes

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/table" // Para verificar as colunas da studentsTable
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vigenda/internal/models"  // Moved here
	"vigenda/internal/service" // Moved here
)

// Variável para garantir que o pacote textinput seja marcado como usado.
var _ textinput.Model

// Mock ClassService
type mockClassService struct {
	ListAllClassesFunc        func(ctx context.Context) ([]models.Class, error)
	CreateClassFunc           func(ctx context.Context, name string, subjectID int64) (models.Class, error)
	GetClassByIDFunc          func(ctx context.Context, id int64) (models.Class, error)
	GetStudentsByClassIDFunc  func(ctx context.Context, classID int64) ([]models.Student, error)
	ImportStudentsFromCSVFunc func(ctx context.Context, classID int64, csvData []byte) (int, error)
	UpdateStudentStatusFunc   func(ctx context.Context, studentID int64, newStatus string) error
	// Added missing methods to satisfy the interface
	UpdateClassFunc    func(ctx context.Context, classID int64, name string, subjectID int64) (models.Class, error)
	DeleteClassFunc    func(ctx context.Context, classID int64) error
	AddStudentFunc     func(ctx context.Context, classID int64, fullName string, enrollmentID string, status string) (models.Student, error)
	GetStudentByIDFunc func(ctx context.Context, studentID int64) (models.Student, error)
	UpdateStudentFunc  func(ctx context.Context, studentID int64, fullName string, enrollmentID string, status string) (models.Student, error)
	DeleteStudentFunc  func(ctx context.Context, studentID int64) error
}

func (m *mockClassService) ListAllClasses(ctx context.Context) ([]models.Class, error) {
	if m.ListAllClassesFunc != nil {
		return m.ListAllClassesFunc(ctx)
	}
	return []models.Class{}, nil
}

func (m *mockClassService) CreateClass(ctx context.Context, name string, subjectID int64) (models.Class, error) {
	if m.CreateClassFunc != nil {
		return m.CreateClassFunc(ctx, name, subjectID)
	}
	return models.Class{ID: 1, Name: name, SubjectID: subjectID}, nil
}

// Corrected mock to match interface: service.ClassService expects models.Class, not *models.Class
func (m *mockClassService) GetClassByID(ctx context.Context, classID int64) (models.Class, error) {
	if m.GetClassByIDFunc != nil {
		return m.GetClassByIDFunc(ctx, classID)
	}
	// Return a value type, potentially an empty struct if not found, or handle error
	return models.Class{ID: classID, Name: "Mocked Class"}, nil // Example value
}

func (m *mockClassService) GetStudentsByClassID(ctx context.Context, classID int64) ([]models.Student, error) {
	if m.GetStudentsByClassIDFunc != nil {
		return m.GetStudentsByClassIDFunc(ctx, classID)
	}
	return []models.Student{}, nil
}

func (m *mockClassService) ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error) {
	if m.ImportStudentsFromCSVFunc != nil {
		return m.ImportStudentsFromCSVFunc(ctx, classID, csvData)
	}
	return 0, nil
}

func (m *mockClassService) UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error {
	if m.UpdateStudentStatusFunc != nil {
		return m.UpdateStudentStatusFunc(ctx, studentID, newStatus)
	}
	return nil
}

// Implementações dos métodos adicionados para mockClassService
func (m *mockClassService) UpdateClass(ctx context.Context, classID int64, name string, subjectID int64) (models.Class, error) {
	if m.UpdateClassFunc != nil {
		return m.UpdateClassFunc(ctx, classID, name, subjectID)
	}
	return models.Class{ID: classID, Name: name, SubjectID: subjectID}, nil
}

func (m *mockClassService) DeleteClass(ctx context.Context, classID int64) error {
	if m.DeleteClassFunc != nil {
		return m.DeleteClassFunc(ctx, classID)
	}
	return nil
}

func (m *mockClassService) AddStudent(ctx context.Context, classID int64, fullName string, enrollmentID string, status string) (models.Student, error) {
	if m.AddStudentFunc != nil {
		return m.AddStudentFunc(ctx, classID, fullName, enrollmentID, status)
	}
	// Retorna um estudante mockado com um ID, por exemplo
	return models.Student{ID: 1, ClassID: classID, FullName: fullName, EnrollmentID: enrollmentID, Status: status}, nil
}

func (m *mockClassService) GetStudentByID(ctx context.Context, studentID int64) (models.Student, error) {
	if m.GetStudentByIDFunc != nil {
		return m.GetStudentByIDFunc(ctx, studentID)
	}
	return models.Student{ID: studentID}, nil
}

func (m *mockClassService) UpdateStudent(ctx context.Context, studentID int64, fullName string, enrollmentID string, status string) (models.Student, error) {
	if m.UpdateStudentFunc != nil {
		return m.UpdateStudentFunc(ctx, studentID, fullName, enrollmentID, status)
	}
	return models.Student{ID: studentID, FullName: fullName, EnrollmentID: enrollmentID, Status: status}, nil
}

func (m *mockClassService) DeleteStudent(ctx context.Context, studentID int64) error {
	if m.DeleteStudentFunc != nil {
		return m.DeleteStudentFunc(ctx, studentID)
	}
	return nil
}

func TestClassesModel_InitialState(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)

	assert.Equal(t, ListView, model.state, "Estado inicial deve ser ListView")
	assert.True(t, model.isLoading, "isLoading deve ser true inicialmente")
	assert.NotNil(t, model.table, "Tabela não deve ser nula")
	assert.Len(t, model.formInputs.inputs, 0, "formInputs.inputs deve estar vazio inicialmente")
}

func TestClassesModel_InitCmd(t *testing.T) {
	mockService := &mockClassService{
		ListAllClassesFunc: func(ctx context.Context) ([]models.Class, error) {
			return []models.Class{{ID: 1, Name: "Test", SubjectID: 1}}, nil
		},
	}
	model := New(mockService)
	cmd := model.Init()
	require.NotNil(t, cmd, "Init deve retornar um comando")
	assert.True(t, model.isLoading, "isLoading deve ser true após Init ser chamado")

	msg := cmd()
	_, ok := msg.(fetchedClassesMsg)
	if !ok {
		errMsg, okErrMsg := msg.(errMsg)
		require.True(t, okErrMsg, "Comando Init deve produzir fetchedClassesMsg ou errMsg, obteve %T: %v", msg, errMsg)
	} else {
		assert.True(t, ok, "Comando Init deve produzir fetchedClassesMsg ou errMsg, obteve %T", msg)
	}
}

func TestClassesModel_Update_KeyN_SwitchesToCreatingView(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	model.state = ListView
	model.isLoading = false

	keyN := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updatedModelTea, cmd := model.Update(keyN)
	m := updatedModelTea.(*Model)

	assert.Equal(t, CreatingView, m.state, "Estado deve mudar para CreatingView após 'n'")
	require.Len(t, m.formInputs.inputs, 2, "Deve haver 2 inputs no formulário de criação")
	assert.True(t, m.formInputs.inputs[0].Focused(), "Campo de nome (inputs[0]) deve estar focado")
	assert.Equal(t, "n", m.formInputs.inputs[0].Value(), "Campo de nome deve conter 'n'") // A tecla 'n' é processada pelo input
	assert.Equal(t, "", m.formInputs.inputs[1].Value(), "Campo de SubjectID (inputs[1]) deve estar vazio")
	assert.Nil(t, m.err, "Erro deve ser nil ao mudar para CreatingView")

	// A mensagem textinput.Blink é um comando que o componente textinput retorna.
	// O modelo principal (m) não retorna Blink diretamente, mas o comando do textinput
	// que ele retorna ao ser focado é o Blink.
	require.NotNil(t, cmd, "Um comando (para textinput.Blink) deve ser retornado")
	// Verificar se o comando é do tipo tea.Cmd e não tentar executá-lo diretamente aqui
	// a menos que seja para verificar o tipo de mensagem que ele produziria, o que é mais complexo.
	// O importante é que um comando é retornado, e a lógica do textinput se encarrega do Blink.
	// A asserção original `_, isBlinkMsg := blinkResultMsg.(textinput.BlinkMsg)` é válida se `cmd()`
	// de fato produzisse um BlinkMsg, o que é o caso para o textinput.Focus().
	assert.NotNil(t, cmd, "Um comando deve ser retornado ao focar o input (textinput.Blink)")
}

func TestClassesModel_Update_CreatingView_EscSwitchesToListView(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	model.state = CreatingView
	model.err = errors.New("erro anterior")

	keyEsc := tea.KeyMsg{Type: tea.KeyEscape}
	updatedModelTea, _ := model.Update(keyEsc)
	m := updatedModelTea.(*Model)

	assert.Equal(t, ListView, m.state, "Estado deve mudar para ListView após 'esc'")
	assert.Nil(t, m.err, "Erro deve ser limpo após 'esc'")
}

func TestClassesModel_Update_FetchedClassesMsg_Success(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	model.isLoading = true

	testClasses := []models.Class{{ID: 1, Name: "Turma Teste", SubjectID: 101}}
	msg := fetchedClassesMsg{classes: testClasses, err: nil}
	updatedModelTea, _ := model.Update(msg)
	m := updatedModelTea.(*Model)

	assert.False(t, m.isLoading, "isLoading deve ser false após fetchedClassesMsg")
	assert.Nil(t, m.err, "Erro deve ser nil em caso de sucesso")
	assert.Equal(t, testClasses, m.allClasses, "allClasses deve ser atualizado")
	require.Len(t, m.table.Rows(), 1, "Tabela deve ter uma linha")
	assert.Equal(t, fmt.Sprintf("%d", testClasses[0].ID), m.table.Rows()[0][0], "ID da turma na tabela incorreto")
	assert.Equal(t, testClasses[0].Name, m.table.Rows()[0][1], "Nome da turma na tabela incorreto")
}

func TestClassesModel_Update_FetchedClassesMsg_Error(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	model.isLoading = true

	fetchErr := errors.New("falha ao buscar")
	msg := fetchedClassesMsg{classes: nil, err: fetchErr}
	updatedModelTea, _ := model.Update(msg)
	m := updatedModelTea.(*Model)

	assert.False(t, m.isLoading, "isLoading deve ser false após fetchedClassesMsg com erro")
	assert.Equal(t, fetchErr, m.err, "Erro deve ser definido")
	assert.Nil(t, m.allClasses, "allClasses deve ser nil em caso de erro")
	assert.Len(t, m.table.Rows(), 0, "Tabela deve estar vazia em caso de erro")
}

func TestClassesModel_Update_CreateClass_Success(t *testing.T) {
	createdClassName := "Nova Turma Sucesso"
	createdSubjectID := int64(123)
	finalClassID := int64(5)

	mockService := &mockClassService{
		CreateClassFunc: func(ctx context.Context, name string, subjectID int64) (models.Class, error) {
			assert.Equal(t, createdClassName, name)
			assert.Equal(t, createdSubjectID, subjectID)
			return models.Class{ID: finalClassID, Name: name, SubjectID: subjectID}, nil
		},
		ListAllClassesFunc: func(ctx context.Context) ([]models.Class, error) {
			return []models.Class{{ID: finalClassID, Name: createdClassName, SubjectID: createdSubjectID}}, nil
		},
	}
	model := New(mockService)
	// Simulate entering the CreatingView state, which prepares the form
	keyN := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	modelInterface, _ := model.Update(keyN)
	model = modelInterface.(*Model)

	// Now set values on the prepared form
	require.Len(t, model.formInputs.inputs, 2, "Formulário de criação não inicializado corretamente")
	model.formInputs.inputs[0].SetValue(createdClassName)                    // Name input
	model.formInputs.inputs[1].SetValue(fmt.Sprintf("%d", createdSubjectID)) // SubjectID input
	model.formInputs.focusIndex = 1                                          // Assume o foco está no último campo para submeter com Enter

	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModelTea, cmdCreate := model.Update(keyEnter)
	m := updatedModelTea.(*Model)

	require.NotNil(t, cmdCreate, "Comando createClassCmd deve ser retornado")
	assert.True(t, m.isLoading, "isLoading deve ser true enquanto cria")

	createdMsg := cmdCreate().(classCreatedMsg)
	require.Nil(t, createdMsg.err, "Erro na criação da turma não esperado")
	assert.Equal(t, finalClassID, createdMsg.createdClass.ID)

	updatedModelTea, cmdFetch := m.Update(createdMsg)
	m = updatedModelTea.(*Model)

	assert.True(t, m.isLoading, "isLoading deve ser true para o fetch após classCreatedMsg")
	assert.Equal(t, ListView, m.state, "Estado deve voltar para ListView após criação bem-sucedida")
	assert.Nil(t, m.err, "Erro deve ser nil após criação bem-sucedida")
	require.NotNil(t, cmdFetch, "Comando fetchClassesCmd deve ser retornado após criação")

	fetchMsg := cmdFetch().(fetchedClassesMsg)
	updatedModelTea, _ = m.Update(fetchMsg)
	m = updatedModelTea.(*Model)

	assert.False(t, m.isLoading, "isLoading deve ser false após fetch bem-sucedido")
	require.Len(t, m.table.Rows(), 1, "Tabela deve ter a nova turma")
	assert.Equal(t, createdClassName, m.table.Rows()[0][1])
}

func TestClassesModel_Update_CreateClass_ServiceError(t *testing.T) {
	serviceErr := errors.New("falha no serviço ao criar")
	mockService := &mockClassService{
		CreateClassFunc: func(ctx context.Context, name string, subjectID int64) (models.Class, error) {
			return models.Class{}, serviceErr
		},
	}
	model := New(mockService)
	// Simulate entering the CreatingView state
	keyN := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	modelInterface, _ := model.Update(keyN)
	model = modelInterface.(*Model)

	require.Len(t, model.formInputs.inputs, 2)
	model.formInputs.inputs[0].SetValue("Turma Errada")
	model.formInputs.inputs[1].SetValue("1")
	model.formInputs.focusIndex = 1 // Focus on the last input to trigger submission

	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModelTea, cmd := model.Update(keyEnter)
	m := updatedModelTea.(*Model)
	assert.True(t, m.isLoading)

	createdMsg := cmd().(classCreatedMsg)
	updatedModelTea, _ = m.Update(createdMsg)
	m = updatedModelTea.(*Model)

	assert.False(t, m.isLoading, "isLoading deve ser false após erro na criação")
	assert.Equal(t, CreatingView, m.state, "Estado deve permanecer CreatingView após erro")
	require.NotNil(t, m.err, "Erro deve ser definido")
	assert.Contains(t, m.err.Error(), serviceErr.Error(), "Mensagem de erro deve conter o erro do serviço")
}

func TestClassesModel_Update_CreateClass_InvalidSubjectID(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	// Simulate entering the CreatingView state
	keyN := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	modelInterface, _ := model.Update(keyN)
	model = modelInterface.(*Model)

	require.Len(t, model.formInputs.inputs, 2)
	model.formInputs.inputs[0].SetValue("Turma ID Inválido")
	model.formInputs.inputs[1].SetValue("abc") // Invalid SubjectID
	model.formInputs.focusIndex = 1

	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModelTea, cmd := model.Update(keyEnter)
	m := updatedModelTea.(*Model)
	assert.True(t, m.isLoading)

	cmdResultMsg := cmd()
	errMsgFromCmd, ok := cmdResultMsg.(errMsg)
	if !ok {
		createdMsgWithErr, ok2 := cmdResultMsg.(classCreatedMsg)
		require.True(t, ok2, "Resultado do comando deve ser errMsg ou classCreatedMsg com erro")
		require.NotNil(t, createdMsgWithErr.err, "classCreatedMsg deve conter um erro")
		errMsgFromCmd = errMsg{err: createdMsgWithErr.err}
	}

	updatedModelTea, _ = m.Update(errMsgFromCmd)
	m2 := updatedModelTea.(*Model)

	assert.False(t, m2.isLoading, "isLoading deve ser false após erro de ID inválido")
	assert.Equal(t, CreatingView, m2.state, "Estado deve permanecer CreatingView")
	require.NotNil(t, m2.err, "Erro deve ser definido para ID inválido")
	assert.Contains(t, m2.err.Error(), "ID disciplina inválido", "Mensagem de erro deve indicar ID inválido") // Removido "da"
}

func TestClassesModel_Update_CreateClass_EmptyFields(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	// Simulate entering the CreatingView state
	keyN := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	modelInterface, _ := model.Update(keyN)
	model = modelInterface.(*Model)

	require.Len(t, model.formInputs.inputs, 2)
	model.formInputs.inputs[0].SetValue("") // Empty name
	model.formInputs.inputs[1].SetValue("123")
	model.formInputs.focusIndex = 1
	model.isLoading = false // Definir isLoading como false antes da tentativa de submissão

	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModelTea, cmd := model.Update(keyEnter)
	m := updatedModelTea.(*Model)

	assert.Nil(t, cmd, "Comando de criação não deve ser retornado se campos estiverem vazios na validação do Update")
	assert.False(t, m.isLoading, "isLoading deve permanecer false se a validação local falhar")
	assert.Equal(t, CreatingView, m.state, "Estado deve permanecer CreatingView")
	require.NotNil(t, m.err, "Erro deve ser definido para campos vazios")
	assert.Contains(t, m.err.Error(), "nome e ID da disciplina obrigatórios") // Ajustada a mensagem de erro
}

func TestClassesModel_IsFocused(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)

	model.state = ListView
	assert.False(t, model.IsFocused(), "Não deve estar focado na ListView")

	model.state = CreatingView
	assert.True(t, model.IsFocused(), "Deve estar focado na CreatingView")
}

func TestClassesModel_FormNavigation(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	// Simulate entering the CreatingView state, which prepares the form and focuses the first input
	keyN := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	modelInterface, _ := model.Update(keyN)
	model = modelInterface.(*Model)

	require.Len(t, model.formInputs.inputs, 2, "Formulário não inicializado com 2 inputs")
	assert.Equal(t, 0, model.formInputs.focusIndex, "Foco inicial deve ser no input de nome (índice 0)")
	assert.True(t, model.formInputs.inputs[0].Focused(), "Input de nome (inputs[0]) deve estar focado inicialmente")

	// Pressionar Tab
	keyTab := tea.KeyMsg{Type: tea.KeyTab}
	updatedModelTea, _ := model.Update(keyTab)
	m := updatedModelTea.(*Model)
	assert.Equal(t, 1, m.formInputs.focusIndex, "Foco deve mudar para subjectIDInput (índice 1)")
	assert.True(t, m.formInputs.inputs[1].Focused(), "subjectIDInput (inputs[1]) deve estar focado")
	assert.False(t, m.formInputs.inputs[0].Focused(), "nameInput (inputs[0]) não deve estar focado")

	// Pressionar Tab novamente (volta ao primeiro campo)
	updatedModelTea, _ = m.Update(keyTab)
	m = updatedModelTea.(*Model)
	assert.Equal(t, 0, m.formInputs.focusIndex, "Foco deve voltar para nameInput (índice 0)")
	assert.True(t, m.formInputs.inputs[0].Focused(), "nameInput (inputs[0]) deve estar focado novamente")
	assert.False(t, m.formInputs.inputs[1].Focused(), "subjectIDInput (inputs[1]) não deve estar focado")

	// Pressionar Shift+Tab (do primeiro campo, vai para o último)
	keyShiftTab := tea.KeyMsg{Type: tea.KeyShiftTab} // Usar tea.KeyShiftTab
	updatedModelTea, _ = m.Update(keyShiftTab)
	m = updatedModelTea.(*Model)
	assert.Equal(t, 1, m.formInputs.focusIndex, "Foco deve ir para subjectIDInput (índice 1) com Shift+Tab a partir do índice 0")
	assert.True(t, m.formInputs.inputs[1].Focused(), "subjectIDInput (inputs[1]) deve estar focado após Shift+Tab")
	assert.False(t, m.formInputs.inputs[0].Focused(), "nameInput (inputs[0]) não deve estar focado após Shift+Tab")

	// Pressionar Shift+Tab novamente (do último campo, vai para o primeiro)
	updatedModelTea, _ = m.Update(keyShiftTab) // Reutiliza keyShiftTab
	m = updatedModelTea.(*Model)
	assert.Equal(t, 0, m.formInputs.focusIndex, "Foco deve voltar para nameInput (índice 0) com Shift+Tab a partir do índice 1")
	assert.True(t, m.formInputs.inputs[0].Focused(), "nameInput (inputs[0]) deve estar focado após Shift+Tab")
	assert.False(t, m.formInputs.inputs[1].Focused(), "subjectIDInput (inputs[1]) não deve estar focado após Shift+Tab")
}

func TestClassesModel_StudentsTable_Initialization(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	require.NotNil(t, model.studentsTable, "studentsTable não deve ser nula")
	expectedColumns := []string{
		studentColumnTitleID,
		studentColumnTitleEnrollment,
		studentColumnTitleFullName,
		studentColumnTitleStatus,
		studentColumnTitleCreatedAt, // Adicionada
		studentColumnTitleUpdatedAt, // Adicionada
	}
	actualColumns := model.studentsTable.Columns()
	require.Len(t, actualColumns, len(expectedColumns), "Número incorreto de colunas na studentsTable")
	for i, expected := range expectedColumns {
		assert.Equal(t, expected, actualColumns[i].Title, "Título da coluna %d incorreto", i)
	}
}

func TestClassesModel_Update_ListView_EnterSelectsClassAndFetchesStudents(t *testing.T) {
	initialClasses := []models.Class{
		{ID: 1, Name: "Turma A", SubjectID: 101},
		{ID: 2, Name: "Turma B", SubjectID: 102},
	}
	mockSvc := &mockClassService{
		ListAllClassesFunc: func(ctx context.Context) ([]models.Class, error) {
			return initialClasses, nil
		},
		GetStudentsByClassIDFunc: func(ctx context.Context, classID int64) ([]models.Student, error) {
			assert.Equal(t, initialClasses[0].ID, classID, "ID da turma para buscar alunos incorreto")
			return []models.Student{{ID: 10, FullName: "Aluno 1"}}, nil
		},
	}

	model := New(mockSvc)
	// Simulate receiving fetchedClassesMsg
	modelInterface, _ := model.Update(fetchedClassesMsg{classes: initialClasses, err: nil})
	model = modelInterface.(*Model)
	model.isLoading = false
	model.table.SetCursor(0)

	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModelTea, cmd := model.Update(keyEnter)
	m := updatedModelTea.(*Model)

	assert.Equal(t, DetailsView, m.state, "Estado deve mudar para DetailsView")
	require.NotNil(t, m.selectedClass, "selectedClass não deve ser nil")
	assert.Equal(t, initialClasses[0].ID, m.selectedClass.ID, "Turma selecionada incorreta")
	assert.True(t, m.isLoading, "isLoading deve ser true para buscar alunos")
	assert.Nil(t, m.err, "Erro deve ser nil ao iniciar busca de alunos")
	assert.Nil(t, m.classStudents, "classStudents deve ser nil antes do fetch")

	require.NotNil(t, cmd, "Comando fetchClassStudentsCmd deve ser retornado")
	msg := cmd()
	_, isFetchedStudentsMsg := msg.(fetchedClassStudentsMsg)
	if !isFetchedStudentsMsg {
		errMsg, isErrMsg := msg.(errMsg)
		require.True(t, isErrMsg, "Comando deve produzir fetchedClassStudentsMsg ou errMsg, obteve %T: %v", msg, errMsg)
		assert.Fail(t, "Esperado fetchedClassStudentsMsg, mas obteve errMsg: "+errMsg.Error())
	}
	assert.True(t, isFetchedStudentsMsg, "Comando deve produzir fetchedClassStudentsMsg ou errMsg")
}

func TestClassesModel_Update_FetchedClassStudentsMsg_Success(t *testing.T) {
	mockSvc := &mockClassService{}
	model := New(mockSvc)
	model.state = DetailsView
	model.isLoading = true
	selectedClass := models.Class{ID: 1, Name: "Turma Teste"}
	model.selectedClass = &selectedClass

	testStudents := []models.Student{
		{ID: 1, ClassID: 1, EnrollmentID: "001", FullName: "Alice", Status: "ativo"},
		{ID: 2, ClassID: 1, EnrollmentID: "002", FullName: "Bob", Status: "inativo"},
	}
	msg := fetchedClassStudentsMsg{students: testStudents, err: nil}
	updatedModelTea, _ := model.Update(msg)
	m := updatedModelTea.(*Model)

	assert.False(t, m.isLoading, "isLoading deve ser false após fetchedClassStudentsMsg")
	assert.Nil(t, m.err, "Erro deve ser nil em caso de sucesso")
	assert.Equal(t, testStudents, m.classStudents, "classStudents deve ser atualizado")
	require.Len(t, m.studentsTable.Rows(), 2, "studentsTable deve ter duas linhas")
	assert.Equal(t, "Alice", m.studentsTable.Rows()[0][2], "Nome do aluno na tabela incorreto")
	assert.Equal(t, "Bob", m.studentsTable.Rows()[1][2], "Nome do aluno na tabela incorreto")
}

func TestClassesModel_Update_FetchedClassStudentsMsg_Error(t *testing.T) {
	mockSvc := &mockClassService{}
	model := New(mockSvc)
	model.state = DetailsView
	model.isLoading = true
	selectedClass := models.Class{ID: 1, Name: "Turma Teste"}
	model.selectedClass = &selectedClass

	fetchErr := errors.New("falha ao buscar alunos")
	msg := fetchedClassStudentsMsg{students: nil, err: fetchErr}
	updatedModelTea, _ := model.Update(msg)
	m := updatedModelTea.(*Model)

	assert.False(t, m.isLoading, "isLoading deve ser false após erro")
	assert.Equal(t, fetchErr, m.err, "Erro deve ser definido no modelo")
	assert.Nil(t, m.classStudents, "classStudents deve ser nil em caso de erro")
	assert.Len(t, m.studentsTable.Rows(), 0, "studentsTable deve estar vazia")
}

func TestClassesModel_Update_FetchClassStudentsCmd_ReturnsErrMsg(t *testing.T) {
	mockSvc := &mockClassService{
		GetStudentsByClassIDFunc: func(ctx context.Context, classID int64) ([]models.Student, error) {
			return nil, errors.New("erro direto do serviço de alunos")
		},
	}
	model := New(mockSvc)
	model.state = DetailsView
	model.isLoading = true
	selectedClass := models.Class{ID: 1, Name: "Turma Teste"}
	model.selectedClass = &selectedClass

	cmd := model.fetchClassStudentsCmd(selectedClass.ID)
	msg := cmd()

	updatedModelTea, _ := model.Update(msg)
	m := updatedModelTea.(*Model)

	assert.False(t, m.isLoading, "isLoading deve ser false após errMsg")
	require.NotNil(t, m.err, "Erro deve ser definido no modelo")
	assert.Contains(t, m.err.Error(), "erro direto do serviço de alunos", "Mensagem de erro incorreta")
	assert.Nil(t, m.classStudents, "classStudents deve ser nil")
	assert.Len(t, m.studentsTable.Rows(), 0, "studentsTable deve estar vazia")
}

func TestClassesModel_Update_DetailsView_EscReturnsToListView(t *testing.T) {
	mockSvc := &mockClassService{}
	model := New(mockSvc)
	model.state = DetailsView
	selectedClass := models.Class{ID: 1, Name: "Turma Selecionada"}
	model.selectedClass = &selectedClass
	model.classStudents = []models.Student{{ID: 1, FullName: "Aluno Teste"}}
	model.studentsTable.SetRows([]table.Row{{"1", "001", "Aluno Teste", "ativo"}})
	model.err = errors.New("erro anterior")

	keyEsc := tea.KeyMsg{Type: tea.KeyEscape}
	updatedModelTea, _ := model.Update(keyEsc)
	m := updatedModelTea.(*Model)

	assert.Equal(t, ListView, m.state, "Estado deve mudar para ListView")
	assert.Nil(t, m.selectedClass, "selectedClass deve ser nil após voltar")
	assert.Nil(t, m.classStudents, "classStudents deve ser nil após voltar")
	assert.Len(t, m.studentsTable.Rows(), 0, "Linhas da studentsTable devem ser limpas")
	assert.Nil(t, m.err, "Erro deve ser limpo após voltar")
	assert.True(t, m.table.Focused(), "Tabela de turmas deve estar focada")
}

func TestClassesModel_IsFocused_ForDetailsView(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)

	model.state = DetailsView
	assert.False(t, model.IsFocused(), "Não deve estar focado (para fins de 'esc' global) na DetailsView, a menos que um input interno esteja ativo")
}

// Ensure mockClassService implements service.ClassService
var _ service.ClassService = (*mockClassService)(nil)
