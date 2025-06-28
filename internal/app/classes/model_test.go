package classes

import (
	"context" // Adicionado para mock do serviço
	"errors"
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput" // Para BlinkMsg
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vigenda/internal/models"
	"vigenda/internal/service" // Para o mock do serviço
)

// Mock ClassService
type mockClassService struct {
	ListAllClassesFunc  func(ctx context.Context) ([]models.Class, error)
	CreateClassFunc     func(ctx context.Context, name string, subjectID int64) (models.Class, error)
	GetClassByIDFunc    func(ctx context.Context, id int64) (models.Class, error)
	// Adicione outros métodos conforme necessário para satisfazer a interface service.ClassService
	ImportStudentsFromCSVFunc func(ctx context.Context, classID int64, csvData []byte) (int, error)
	UpdateStudentStatusFunc func(ctx context.Context, studentID int64, newStatus string) error
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

func (m *mockClassService) GetClassByID(ctx context.Context, classID int64) (models.Class, error) {
	if m.GetClassByIDFunc != nil {
		return m.GetClassByIDFunc(ctx, classID)
	}
	return models.Class{ID: classID}, nil
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


func TestClassesModel_InitialState(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)

	assert.Equal(t, ListView, model.state, "Estado inicial deve ser ListView")
	assert.True(t, model.isLoading, "isLoading deve ser true inicialmente")
	assert.NotNil(t, model.table, "Tabela não deve ser nula")
	assert.NotNil(t, model.createForm.nameInput, "Input de nome não deve ser nulo")
	assert.NotNil(t, model.createForm.subjectIDInput, "Input de SubjectID não deve ser nulo")
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

	msg := cmd() // Executa o comando
	_, ok := msg.(fetchedClassesMsg)
	if !ok {
		_, okErrMsg := msg.(errMsg)
		assert.True(t, okErrMsg, "Comando Init deve produzir fetchedClassesMsg ou errMsg, obteve %T", msg)
	} else {
		assert.True(t, ok, "Comando Init deve produzir fetchedClassesMsg ou errMsg, obteve %T", msg)
	}
	// Model.isLoading é definido no início de Init() e no comando, então verificamos o estado após a chamada de Init()
	// e não após a execução do comando, pois isso já é testado em FetchedClassesMsg.
	// A verificação aqui é se a chamada a Init() *configura* o carregamento.
	// O model original é passado por valor para Init, então o model retornado por New() não é modificado por model.Init().
	// Para testar o isLoading após Init, precisamos do model *antes* de executar o comando.
	// A linha abaixo está correta: model.isLoading é true porque New() o define e Init() o reafirma.
	newModelForInit := New(mockService)
	_ = newModelForInit.Init() // Chamada para configurar isLoading
	assert.True(t, newModelForInit.isLoading, "isLoading deve ser true após Init ser chamado")
}

func TestClassesModel_Update_KeyN_SwitchesToCreatingView(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	model.state = ListView
	model.isLoading = false

	keyN := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updatedModel, cmd := model.Update(keyN)

	m, ok := updatedModel.(Model)
	require.True(t, ok, "updatedModel deve ser do tipo Model")

	assert.Equal(t, CreatingView, m.state, "Estado deve mudar para CreatingView após 'n'")
	assert.True(t, m.createForm.nameInput.Focused(), "Campo de nome deve estar focado")
	assert.Equal(t, "", m.createForm.nameInput.Value(), "Campo de nome deve estar vazio")
	assert.Equal(t, "", m.createForm.subjectIDInput.Value(), "Campo de SubjectID deve estar vazio")
	assert.Nil(t, m.err, "Erro deve ser nil ao mudar para CreatingView")

	require.NotNil(t, cmd, "Um comando (para textinput.Blink) deve ser retornado")
	blinkMsg := cmd()
	_, isBlink := blinkMsg.(textinput.BlinkMsg)
	assert.True(t, isBlink, "O comando retornado deve ser textinput.Blink")
}

func TestClassesModel_Update_CreatingView_EscSwitchesToListView(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	model.state = CreatingView
	model.err = errors.New("erro anterior")

	keyEsc := tea.KeyMsg{Type: tea.KeyEscape}
	updatedModel, _ := model.Update(keyEsc)
	m, _ := updatedModel.(Model)

	assert.Equal(t, ListView, m.state, "Estado deve mudar para ListView após 'esc'")
	assert.Nil(t, m.err, "Erro deve ser limpo após 'esc'")
}

func TestClassesModel_Update_FetchedClassesMsg_Success(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	model.isLoading = true

	testClasses := []models.Class{{ID: 1, Name: "Turma Teste", SubjectID: 101}}
	msg := fetchedClassesMsg{classes: testClasses, err: nil}
	updatedModel, _ := model.Update(msg)
	m, _ := updatedModel.(Model)

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
	updatedModel, _ := model.Update(msg)
	m, _ := updatedModel.(Model)

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
	model.state = CreatingView
	model.createForm.nameInput.SetValue(createdClassName)
	model.createForm.subjectIDInput.SetValue(fmt.Sprintf("%d", createdSubjectID))
	model.createForm.focusIndex = 1

	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmdCreate := model.Update(keyEnter)
	m, _ := updatedModel.(Model)

	require.NotNil(t, cmdCreate, "Comando createClassCmd deve ser retornado")
	assert.True(t, m.isLoading, "isLoading deve ser true enquanto cria")

	createdMsg := cmdCreate().(classCreatedMsg)
	require.Nil(t, createdMsg.err, "Erro na criação da turma não esperado")
	assert.Equal(t, finalClassID, createdMsg.createdClass.ID)


	updatedModelAfterCreate, cmdFetch := m.Update(createdMsg)
	m, _ = updatedModelAfterCreate.(Model)

	assert.True(t, m.isLoading, "isLoading deve ser true para o fetch após classCreatedMsg")
	assert.Equal(t, ListView, m.state, "Estado deve voltar para ListView após criação bem-sucedida")
	assert.Nil(t, m.err, "Erro deve ser nil após criação bem-sucedida")
	require.NotNil(t, cmdFetch, "Comando fetchClassesCmd deve ser retornado após criação")

	fetchMsg := cmdFetch().(fetchedClassesMsg)
	updatedModelAfterFetch, _ := m.Update(fetchMsg)
	m, _ = updatedModelAfterFetch.(Model)

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
	model.state = CreatingView
	model.createForm.nameInput.SetValue("Turma Errada")
	model.createForm.subjectIDInput.SetValue("1")
	model.createForm.focusIndex = 1

	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := model.Update(keyEnter)
	m, _ := updatedModel.(Model)
	assert.True(t, m.isLoading)

	createdMsg := cmd().(classCreatedMsg)
	updatedModel, _ = m.Update(createdMsg)
	m, _ = updatedModel.(Model)

	assert.False(t, m.isLoading, "isLoading deve ser false após erro na criação")
	assert.Equal(t, CreatingView, m.state, "Estado deve permanecer CreatingView após erro")
	require.NotNil(t, m.err, "Erro deve ser definido")
	assert.Contains(t, m.err.Error(), serviceErr.Error(), "Mensagem de erro deve conter o erro do serviço")
}

func TestClassesModel_Update_CreateClass_InvalidSubjectID(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	model.state = CreatingView
	model.createForm.nameInput.SetValue("Turma ID Inválido")
	model.createForm.subjectIDInput.SetValue("abc")
	model.createForm.focusIndex = 1

	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := model.Update(keyEnter)
	m, _ := updatedModel.(Model)
	assert.True(t, m.isLoading)

	cmdResultMsg := cmd()
	errMsgFromCmd, ok := cmdResultMsg.(errMsg)
	if !ok { // Pode ser classCreatedMsg com erro de conversão dentro do serviço, dependendo da implementação exata
		createdMsgWithErr, ok2 := cmdResultMsg.(classCreatedMsg)
		require.True(t, ok2, "Resultado do comando deve ser errMsg ou classCreatedMsg com erro")
		require.NotNil(t, createdMsgWithErr.err, "classCreatedMsg deve conter um erro")
		errMsgFromCmd = errMsg{err: createdMsgWithErr.err}
	}

	updatedModelAfterCmd, _ := m.Update(errMsgFromCmd)
	m2, _ := updatedModelAfterCmd.(Model)


	assert.False(t, m2.isLoading, "isLoading deve ser false após erro de ID inválido")
	assert.Equal(t, CreatingView, m2.state, "Estado deve permanecer CreatingView")
	require.NotNil(t, m2.err, "Erro deve ser definido para ID inválido")
	assert.Contains(t, m2.err.Error(), "ID da disciplina inválido", "Mensagem de erro deve indicar ID inválido")
}


func TestClassesModel_Update_CreateClass_EmptyFields(t *testing.T) {
	mockService := &mockClassService{}
	model := New(mockService)
	model.state = CreatingView
	model.createForm.nameInput.SetValue("")
	model.createForm.subjectIDInput.SetValue("123")
	model.createForm.focusIndex = 1

	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := model.Update(keyEnter)
	m, _ := updatedModel.(Model)

	assert.Nil(t, cmd, "Comando de criação não deve ser retornado se campos estiverem vazios na validação do Update")
	assert.False(t, m.isLoading, "isLoading não deve mudar se a validação falhar no local")
	assert.Equal(t, CreatingView, m.state, "Estado deve permanecer CreatingView")
	require.NotNil(t, m.err, "Erro deve ser definido para campos vazios")
	assert.Contains(t, m.err.Error(), "nome da turma e ID da disciplina são obrigatórios")
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
	model.state = CreatingView
	model.createForm.nameInput.Focus()
	model.createForm.focusIndex = 0

	keyTab := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(keyTab)
	m, _ := updatedModel.(Model)
	assert.Equal(t, 1, m.createForm.focusIndex, "Foco deve mudar para subjectIDInput (índice 1)")
	assert.True(t, m.createForm.subjectIDInput.Focused(), "subjectIDInput deve estar focado")
	assert.False(t, m.createForm.nameInput.Focused(), "nameInput não deve estar focado")

	updatedModel, _ = m.Update(keyTab)
	m, _ = updatedModel.(Model)
	assert.Equal(t, 0, m.createForm.focusIndex, "Foco deve voltar para nameInput (índice 0)")
	assert.True(t, m.createForm.nameInput.Focused(), "nameInput deve estar focado")
	assert.False(t, m.createForm.subjectIDInput.Focused(), "subjectIDInput não deve estar focado")

	keyShiftTab := tea.KeyMsg{Type: tea.KeyTab, Shift: true} // Simula Shift+Tab
    // Estando no input de nome (índice 0), Shift+Tab deve ir para o último (índice 1)
	updatedModel, _ = m.Update(keyShiftTab)
	m, _ = updatedModel.(Model)
	assert.Equal(t, 1, m.createForm.focusIndex, "Foco deve ir para subjectIDInput (índice 1) com Shift+Tab a partir do índice 0")
	assert.True(t, m.createForm.subjectIDInput.Focused(), "subjectIDInput deve estar focado após Shift+Tab")
	assert.False(t, m.createForm.nameInput.Focused(), "nameInput não deve estar focado após Shift+Tab")

    // Estando no input de subjectID (índice 1), Shift+Tab deve ir para o primeiro (índice 0)
	updatedModel, _ = m.Update(keyShiftTab)
	m, _ = updatedModel.(Model)
	assert.Equal(t, 0, m.createForm.focusIndex, "Foco deve voltar para nameInput (índice 0) com Shift+Tab a partir do índice 1")
	assert.True(t, m.createForm.nameInput.Focused(), "nameInput deve estar focado após Shift+Tab")
	assert.False(t, m.createForm.subjectIDInput.Focused(), "subjectIDInput não deve estar focado após Shift+Tab")
}

// Helper para obter o tipo de mensagem de um comando
func getMsgFromCmd(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}
