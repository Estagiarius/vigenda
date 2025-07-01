package database

import (
	"database/sql"
	"log"

	// Importa o driver SQLite3
	_ "github.com/mattn/go-sqlite3"
)

// SeedData popula o banco de dados com dados de exemplo se estiver vazio.
func SeedData(db *sql.DB) error {
	// Verificar se o banco de dados já foi populado (ex: checando a tabela users)
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar a tabela users: %v", err)
		return err
	}

	// Se já houver usuários, assume-se que os dados de exemplo já existem ou não são necessários.
	if count > 0 {
		log.Println("Banco de dados já contém dados. Não é necessário popular com exemplos.")
		return nil
	}

	log.Println("Populando o banco de dados com dados de exemplo...")

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback em caso de erro

	// 1. Inserir Usuário de Exemplo
	// NOTA: Em um cenário real, a senha deve ser um hash.
	// Para simplicidade, usaremos uma senha placeholder, mas isso NÃO é seguro.
	// A aplicação real deve ter um mecanismo de hashing de senhas.
	res, err := tx.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", "demo_user", "hashed_password_placeholder")
	if err != nil {
		log.Printf("Erro ao inserir usuário de exemplo: %v", err)
		return err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		log.Printf("Erro ao obter ID do usuário de exemplo: %v", err)
		return err
	}

	// 2. Inserir Disciplinas de Exemplo
	stmtSubject, err := tx.Prepare("INSERT INTO subjects (user_id, name) VALUES (?, ?)")
	if err != nil {
		log.Printf("Erro ao preparar statement para subjects: %v", err)
		return err
	}
	defer stmtSubject.Close()

	subjects := []struct {
		Name string
		ID   int64 // Para armazenar o ID da disciplina inserida
	}{
		{Name: "Matemática"},
		{Name: "História"},
	}

	for i, s := range subjects {
		res, err = stmtSubject.Exec(userID, s.Name)
		if err != nil {
			log.Printf("Erro ao inserir disciplina '%s': %v", s.Name, err)
			return err
		}
		subjectID, err := res.LastInsertId()
		if err != nil {
			log.Printf("Erro ao obter ID da disciplina '%s': %v", s.Name, err)
			return err
		}
		subjects[i].ID = subjectID // Armazena o ID para uso posterior
	}

	// 3. Inserir Turmas de Exemplo
	stmtClass, err := tx.Prepare("INSERT INTO classes (user_id, subject_id, name) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Erro ao preparar statement para classes: %v", err)
		return err
	}
	defer stmtClass.Close()

	classes := []struct {
		SubjectID int64 // ID da disciplina correspondente
		Name      string
		ID        int64 // Para armazenar o ID da turma inserida
	}{
		{SubjectID: subjects[0].ID, Name: "Turma A - Matemática (Manhã)"}, // Matemática
		{SubjectID: subjects[0].ID, Name: "Turma B - Matemática (Tarde)"},  // Matemática
		{SubjectID: subjects[1].ID, Name: "Turma Única - História (Noite)"}, // História
	}

	for i, c := range classes {
		res, err = stmtClass.Exec(userID, c.SubjectID, c.Name)
		if err != nil {
			log.Printf("Erro ao inserir turma '%s': %v", c.Name, err)
			return err
		}
		classID, err := res.LastInsertId()
		if err != nil {
			log.Printf("Erro ao obter ID da turma '%s': %v", c.Name, err)
			return err
		}
		classes[i].ID = classID // Armazena o ID para uso posterior
	}

	// 4. Inserir Alunos de Exemplo
	stmtStudent, err := tx.Prepare("INSERT INTO students (class_id, full_name, enrollment_id, status) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("Erro ao preparar statement para students: %v", err)
		return err
	}
	defer stmtStudent.Close()

	studentsData := []struct {
		ClassID      int64
		FullName     string
		EnrollmentID string
		Status       string
	}{
		// Alunos da Turma 1 (Matemática - Turma A)
		{ClassID: classes[0].ID, FullName: "Alice Silva", EnrollmentID: "M001", Status: "ativo"},
		{ClassID: classes[0].ID, FullName: "Bruno Costa", EnrollmentID: "M002", Status: "ativo"},
		{ClassID: classes[0].ID, FullName: "Carla Dias", EnrollmentID: "M003", Status: "transferido"},
		// Alunos da Turma 2 (Matemática - Turma B)
		{ClassID: classes[1].ID, FullName: "Daniel Oliveira", EnrollmentID: "M101", Status: "ativo"},
		{ClassID: classes[1].ID, FullName: "Eduarda Ferreira", EnrollmentID: "M102", Status: "ativo"},
		// Alunos da Turma 3 (História - Turma Única)
		{ClassID: classes[2].ID, FullName: "Fernando Almeida", EnrollmentID: "H201", Status: "ativo"},
		{ClassID: classes[2].ID, FullName: "Gabriela Santos", EnrollmentID: "H202", Status: "inativo"},
		{ClassID: classes[2].ID, FullName: "Heitor Lima", EnrollmentID: "H203", Status: "ativo"},
	}

	for _, s := range studentsData {
		_, err = stmtStudent.Exec(s.ClassID, s.FullName, s.EnrollmentID, s.Status)
		if err != nil {
			log.Printf("Erro ao inserir aluno '%s': %v", s.FullName, err)
			return err
		}
	}

	log.Println("Dados de exemplo inseridos com sucesso.")
	return tx.Commit()
}
