CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS classes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    subject_id INTEGER NOT NULL,
    name TEXT NOT NULL, -- Ex: "Turma 9A - 2025"
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(subject_id) REFERENCES subjects(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_id INTEGER NOT NULL,
    full_name TEXT NOT NULL,
    enrollment_id TEXT, -- Número de Matrícula/Chamada
    status TEXT NOT NULL DEFAULT 'ativo', -- Valores permitidos: 'ativo', 'inativo', 'transferido'
    FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS lessons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    plan_content TEXT, -- Conteúdo do plano de aula em Markdown
    scheduled_at TIMESTAMP NOT NULL,
    FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS assessments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_id INTEGER NOT NULL,
    name TEXT NOT NULL, -- Ex: "Prova Bimestral 1"
    term INTEGER NOT NULL, -- Ex: 1, 2, 3, 4 (para o bimestre)
    weight REAL NOT NULL, -- Ex: 4.0
    assessment_date DATE,
    FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS grades (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    assessment_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    grade REAL NOT NULL,
    FOREIGN KEY(assessment_id) REFERENCES assessments(id) ON DELETE CASCADE,
    FOREIGN KEY(student_id) REFERENCES students(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    class_id INTEGER, -- Uma tarefa pode estar associada a uma turma específica
    title TEXT NOT NULL,
    description TEXT,
    due_date TIMESTAMP,
    is_completed BOOLEAN NOT NULL DEFAULT 0,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    subject_id INTEGER NOT NULL,
    topic TEXT,
    type TEXT NOT NULL, -- 'multipla_escolha' ou 'dissertativa'
    difficulty TEXT NOT NULL, -- 'facil', 'media', 'dificil'
    statement TEXT NOT NULL,
    options TEXT, -- JSON array como string para multipla escolha
    correct_answer TEXT NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(subject_id) REFERENCES subjects(id) ON DELETE CASCADE
);
