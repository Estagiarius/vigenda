
## **Guia de Implementação Completo para IA Generativa**

Projeto: Vigenda

Versão do Documento: Final 2.4

### **Lembrete Essencial para a Execução**

**Antes de iniciar cada nova fase de implementação, este documento deve ser consultado na sua totalidade.** Ele é a fonte única de verdade para todos os requisitos, designs e especificações do projeto. A aderência rigorosa a este guia é fundamental para o sucesso da implementação.

### **Artefacto 1: Plano de Projeto em Cascata com Tarefas para IA**

#### **1. Introdução à Metodologia**

Este projeto seguirá uma metodologia em cascata (Waterfall) estrita. Cada fase deve ser totalmente concluída, revista e validada antes que a fase seguinte seja iniciada. Esta estrutura é projetada para ser executada por uma IA Generativa, com tarefas definidas como instruções atómicas, claras e inequívocas para garantir uma implementação previsível e ordenada.

#### **2. Fases do Projeto**

1.  **Análise de Requisitos:** Consolidação e formalização dos requisitos.
    
2.  **Design do Sistema e Arquitetura:** Desenho completo do "blueprint" técnico.
    
3.  **Implementação (Codificação):** Escrita de todo o código-fonte da aplicação.
    
4.  **Testes e Integração:** Verificação e validação do software completo.
    
5.  **Implantação (Deployment):** Preparação e lançamento da versão final.
    
6.  **Manutenção:** Ciclo de vida pós-lançamento.
    

#### **3. Plano de Tarefas por Fase**

**FASE 1: Análise de Requisitos**

-   **Objetivo:** Validar e congelar o escopo completo do projeto.
    
-   Critério de Conclusão: O documento de requisitos é aprovado sem alterações pendentes.
    
    | ID da Tarefa | Instrução de Execução |
    
    | :--- | :--- |
    
    | TASK-R-01 | Analisar os Artefactos 2 a 9 deste documento. Extrair todos os requisitos funcionais e não-funcionais para uma lista de verificação. Validar que não há contradições. Sinalizar o fim da fase. |
    

**FASE 2: Design do Sistema e Arquitetura**

-   **Objetivo:** Produzir um plano técnico completo (blueprint) da aplicação.
    
-   **Critério de Conclusão:** Todos os artefactos de design listados no Artefacto 3 e 4 estão completos e aprovados.
    

**ID da Tarefa**

**Instrução de Execução**

`TASK-D-01`

`Com base no Artefacto 3.2, gerar o ficheiro de migração SQL inicial (`001_initial_schema.sql`) que cria todas as tabelas, colunas, tipos, chaves e índices definidos.`

`TASK-D-02`

`Com base no Artefacto 3.1, gerar a estrutura de diretórios completa do projeto Go. Gerar os ficheiros Go vazios dentro de cada pacote com os comentários de cabeçalho descrevendo a sua finalidade.`

`TASK-D-03`

`Com base no Artefacto 4, gerar os ficheiros de interface Go para cada serviço (TaskService, ClassService, etc.) dentro do pacote` /internal/service/`.`

**FASE 3: Implementação (Codificação)**

-   **Objetivo:** Escrever todo o código-fonte da aplicação conforme as especificações.
    
-   **Critério de Conclusão:** Todo o código está escrito, comentado e passa nos testes unitários.
    

**ID da Tarefa**

**Instrução de Execução**

`TASK-I-01`

`Implementar a camada de acesso à base de dados (`/internal/repository`) seguindo o Padrão Repositório. Para cada` struct `de modelo, criar as funções CRUD correspondentes. Implementar a lógica de conexão com o ficheiro SQLite. Adicionar testes unitários que simulem a base de dados (mock).`

`TASK-I-02`

`Implementar a estrutura principal da CLI em` /cmd/vigenda/`usando`Cobra`. Criar o comando raiz e registar todos os subcomandos como stubs (sem lógica). Implementar a lógica para carregar a configuração (`config.toml`).`

`TASK-I-03`

`Implementar a lógica de negócio completa para os serviços de Produtividade (TaskService, AgendaService, RoutineService) no pacote` /internal/service/`, interagindo com a camada de repositório. Adicionar testes unitários.`

`TASK-I-04`

`Implementar a lógica de negócio completa para os serviços de Gestão Académica (ClassService, AssessmentService) no pacote` /internal/service/`. Adicionar testes unitários.`

`TASK-I-05`

`Implementar a lógica de negócio completa para os serviços de Conteúdo Pedagógico (QuestionService, ProofService) no pacote` /internal/service/`. Adicionar testes unitários.`

`TASK-I-06`

`Implementar os componentes reutilizáveis da UI no pacote` /internal/tui/`usando`Bubble Tea `(ex: tabelas, prompts, barra de status, etc.).`

`TASK-I-07`

`Implementar a lógica de ligação (glue code) dentro dos comandos da CLI. Conectar os comandos` Cobra`para que chamem as funções da camada de serviço e usem os componentes do pacote`/internal/tui/ `para renderizar a saída.`

**FASE 4: Testes e Integração**

-   **Objetivo:** Identificar e corrigir defeitos no software como um todo.
    
-   **Critério de Conclusão:** A aplicação está funcional, estável e cumpre todos os requisitos definidos nos Casos de Teste.
    

**ID da Tarefa**

**Instrução de Execução**

`TASK-T-01`

`Com base no Artefacto 6, gerar os ficheiros de teste de integração automatizados. Os testes devem executar a CLI compilada como um subprocesso e validar a sua saída contra os ficheiros "golden" definidos no Artefacto 7.`

`TASK-T-02`

`Executar todos os testes de integração. Listar todas as falhas e as discrepâncias encontradas entre a saída real e os ficheiros "golden".`

`TASK-T-03`

`Corrigir todos os bugs encontrados na fase anterior. Cada correção deve ser validada executando novamente o teste de integração correspondente até que ele passe.`

**FASE 5: Implantação (Deployment)**

-   **Objetivo:** Empacotar e lançar a versão final do software.
    
-   **Critério de Conclusão:** O software está disponível para os utilizadores finais.
    

**ID da Tarefa**

**Instrução de Execução**

`TASK-P-01`

`Rever e finalizar a documentação de ajuda (`--help`) para todos os comandos, garantindo clareza e completude. Gerar um ficheiro` README.md `abrangente.`

`TASK-P-02`

`Criar um script de build (ex: um` Makefile`ou`build.sh`) que automatiza a compilação cruzada dos binários para Windows (amd64), macOS (amd64, arm64) e Linux (amd64).`

**FASE 6: Manutenção**

-   **Objetivo:** Dar suporte ao software em produção.
    
-   Critério de Conclusão: Processo contínuo.
    
    | ID da Tarefa | Instrução de Execução |
    
    | :--- | :--- |
    
    | TASK-M-01 | Monitorizar o sistema de reporte de bugs. Analisar, priorizar e criar novas tarefas de correção.|
    

### **Artefacto 2: Documento de Requisitos de Software (ERS) - Versão CLI**

-   **Propósito:** Fornecer uma ferramenta de organização e produtividade para professores com TDAH, focada na eficiência através de uma CLI.
    
-   **Escopo:** Uma aplicação CLI autónoma para Windows, Mac e Linux, com uma base de dados local SQLite.
    
-   **Requisitos Funcionais (RF):**
    
    -   **RF01 (Dashboard):** Ao executar o comando principal (`vigenda`), exibir um dashboard com agenda do dia (aulas por turma), tarefas urgentes e notificações.
        
    -   **RF02 (Agenda):** O comando `vigenda agenda` deve permitir a gestão de eventos, incluindo aulas associadas a turmas específicas.
        
    -   **RF03 (Tarefas):** O comando `vigenda tarefa` deve permitir a gestão de tarefas, que podem ser associadas a uma disciplina ou a uma turma específica. O comando `vigenda rotina` deve permitir a criação de tarefas recorrentes.
        
    -   **RF04 (Aulas e Turmas):** O conceito de "Turma" é central. O comando `vigenda turma` deve permitir a criação de turmas (ex: "História - Turma 9A") e a gestão de alunos dentro dessas turmas, incluindo importação em massa e gestão de status.
        
    -   **RF05 (Avaliações):** O comando `vigenda avaliacao` deve permitir a criação de avaliações (com pesos) para uma turma específica. O comando `vigenda notas` deve permitir o lançamento de notas por aluno para uma avaliação específica e o cálculo de médias por turma.
        
    -   **RF06 (Foco):** O comando `vigenda foco iniciar` deve iniciar uma sessão de trabalho cronometrada e sem distrações.
        
    -   **RF07 (Conteúdo):** O comando `vigenda bancoq` deve gerir um banco de questões por disciplina. O comando `vigenda prova` deve gerar avaliações textuais a partir deste banco.
        

### **Artefacto 3: Documento de Design de Arquitetura (DDA)**

#### **3.1. Arquitetura de Pacotes (Go)**

-   `/cmd/vigenda/`: Ponto de entrada da aplicação (função `main` e configuração da CLI `Cobra`).
    
-   `/internal/config/`: Lógica para ler e gerir o ficheiro de configuração `config.toml`.
    
-   `/internal/database/`: Código para a ligação com SQLite e execução de migrações.
    
-   `/internal/models/`: Definição de todas as `structs` Go (ex: `Task`, `Lesson`, `Student`, `Class`).
    
-   `/internal/repository/`: Implementação do Padrão Repositório (queries SQL).
    
-   `/internal/service/`: Camada de serviço com a lógica de negócio.
    
-   `/internal/tui/`: Implementação da Interface de Utilizador Textual (TUI), utilizando a biblioteca `Bubble Tea` e seus componentes.
    

#### **3.2. Esquema da Base de Dados (SQL para SQLite)**

```
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

```

### **Artefacto 4: Especificação da API de Serviço Interna (Go Interfaces)**

```
package service

import (
	"context"
	"time"
	"vigenda/internal/models"
)

// TaskService define os métodos para a gestão de tarefas.
type TaskService interface {
	CreateTask(ctx context.Context, title, description string, classID *int64, dueDate *time.Time) (models.Task, error)
	ListActiveTasksByClass(ctx context.Context, classID int64) ([]models.Task, error)
	MarkTaskAsCompleted(ctx context.Context, taskID int64) error
}

// ClassService define os métodos para a gestão de turmas e alunos.
type ClassService interface {
    CreateClass(ctx context.Context, name string, subjectID int64) (models.Class, error)
    ImportStudentsFromCSV(ctx context.Context, classID int64, csvData []byte) (int, error)
    UpdateStudentStatus(ctx context.Context, studentID int64, newStatus string) error
}

// AssessmentService define os métodos para a gestão de avaliações e notas.
type AssessmentService interface {
    CreateAssessment(ctx context.Context, name string, classID int64, term int, weight float64) (models.Assessment, error)
    EnterGrades(ctx context.Context, assessmentID int64, studentGrades map[int64]float64) error
    CalculateClassAverage(ctx context.Context, classID int64) (float64, error)
}

// QuestionService define os métodos para o banco de questões e geração de provas.
type QuestionService interface {
	AddQuestionsFromJSON(ctx context.Context, jsonData []byte) (int, error)
	GenerateTest(ctx context.Context, criteria TestCriteria) ([]models.Question, error)
}

// TestCriteria define os parâmetros para a geração de uma prova.
type TestCriteria struct {
	SubjectID   int64
	Topic       *string // Opcional: pode ser nil se não for filtrar por tópico específico
	EasyCount   int
	MediumCount int
	HardCount   int
}

// ProofService define os métodos para a geração de provas (potencialmente uma visão mais específica ou diferente de Test).
type ProofService interface {
	GenerateProof(ctx context.Context, criteria ProofCriteria) ([]models.Question, error)
}

// ProofCriteria define os parâmetros para a geração de uma prova.
type ProofCriteria struct {
	SubjectID   int64
	Topic       *string // Opcional: pode ser nil se não for filtrar por tópico específico
	EasyCount   int
	MediumCount int
	HardCount   int
}

// Adicionar outras interfaces de serviço aqui: SubjectService, LessonService, etc.

```

### **Artefacto 5: User Stories e Critérios de Aceitação**

-   **US-001: Criação de uma Nova Tarefa para uma Turma**
    
    -   **Como um** professor, **Eu quero** adicionar uma tarefa específica para a "Turma 9A", **Para que** eu possa organizar o trabalho daquela turma.
        
    -   **Critérios de Aceitação:** Dado que eu executo `vigenda tarefa add`, quando eu seleciono a "Turma 9A", então a nova tarefa deve ser criada e associada ao `id` correto da turma na base de dados.
        
-   **US-002: Lançamento de Notas para uma Avaliação**
    
    -   **Como um** professor, **Eu quero** criar uma avaliação para a "Turma 9A" e depois lançar as notas de cada aluno, **Para que** eu possa acompanhar o desempenho da turma.
        
    -   **Critérios de Aceitação:**
        
        1.  Dado que a "Turma 9A" existe, quando eu executo `vigenda avaliacao criar`, então uma nova avaliação é criada para essa turma.
            
        2.  Dado que a avaliação existe, quando eu executo `vigenda notas lancar --avaliacao <id>`, então a aplicação me mostra uma lista apenas com os alunos ativos da "Turma 9A".
            
        3.  Quando eu insiro as notas, então os registos são criados na tabela `grades`, ligando cada `aluno` à `avaliação` com a sua respetiva `nota`.
            
-   **US-003: Gestão de Alunos**
    
    -   **Como um** professor, **Eu quero** importar uma lista de alunos de um ficheiro CSV para a "Turma 9B" e depois marcar um deles como "transferido", **Para que** a minha lista de chamada esteja sempre atualizada.
        
    -   **Critérios de Aceitação:**
        
        1.  Dado que tenho um ficheiro `turma9b.csv`, quando eu executo `vigenda turma importar-alunos`, então todos os alunos são adicionados à "Turma 9B".
            
        2.  Dado que o aluno "Daniel Mendes" existe, quando eu executo `vigenda turma atualizar-status --aluno <id_daniel> --status transferido`, então o status do aluno é alterado.
            
        3.  Quando eu for lançar notas, então o nome "Daniel Mendes" deve aparecer na lista, mas indicado como "(Transferido)" e não deve ser editável.
            

### **Artefacto 6: Especificações de Casos de Teste (Unitários e de Integração)**

-   **Módulo: `AssessmentService` (Testes Unitários)**
    
    -   **TC-U-001 (Sucesso):** Deve calcular a média ponderada correta de uma turma com várias avaliações e pesos diferentes.
        
    -   **TC-U-002 (Lógica):** Deve ignorar alunos com status "inativo" ou "transferido" ao calcular a média da turma.
        
    -   **TC-U-003 (Falha):** Deve retornar um erro se o utilizador tentar lançar uma nota para um aluno que não pertence à turma da avaliação.
        
-   **Módulo: CLI (Testes de Integração)**
    
    -   **TC-I-001 (Sucesso):** O comando `vigenda notas lancar` deve apresentar a lista correta de alunos ativos para a avaliação selecionada.
        
    -   **TC-I-002 (Sucesso):** O comando `vigenda relatorio progresso-turma` deve exibir os dados corretos e calculados para a turma especificada.
        

### **Artefacto 7: Ficheiros "Golden" para Testes de UI (Expandido)**

-   **Ficheiro 1: `golden_files/dashboard_output.txt`**
    
    ```
    =================================================
    ==                 DASHBOARD                   ==
    =================================================
    
    🕒 AGENDA DE HOJE (22/06/2025)
       [09:00 - 10:00] Aula de História - Turma 9A
       [14:00 - 15:00] Reunião Pedagógica
    
    🔥 TAREFAS PRIORITÁRIAS
       [1] Corrigir provas (Turma 9A) (Prazo: Amanhã)
       [2] Preparar aula sobre Era Vargas (Turma 9B) (Prazo: 24/06)
    
    🔔 NOTIFICAÇÕES
       - 5 entregas pendentes para o trabalho "Pesquisa sobre Clima" (Turma 9A).
    
    
    ```
    
-   **Ficheiro 2: `golden_files/tarefa_listar_turma_output.txt`**
    
    ```
    $ vigenda tarefa listar --turma "Turma 9A"
    
    TAREFAS PARA: Turma 9A
    
    ID | TAREFA                            | PRAZO
    -- | --------------------------------- | ----------
    1  | Corrigir provas de Matemática     | 23/06/2025
    5  | Lançar notas do trabalho          | 25/06/2025
    
    
    ```
    
-   **Ficheiro 3: `golden_files/foco_iniciar_output.txt`** (Ecrã do modo foco)
    
    ```
    ======================================================================
    ==                          MODO FOCO                             ==
    ======================================================================
    
    TAREFA: Corrigir provas de Matemática (Turma 9A)
    
    TEMPO RESTANTE: 24:59
    
    (Pressione 'espaço' para pausar/retomar, 'q' para sair e concluir o ciclo)
    
    ```
    
-   **Ficheiro 4: `golden_files/notas_lancar_interativo_output.txt`** (Exemplo de interação)
    
    ```
    Lançamento de notas para a avaliação "Prova Bimestral 1" - Turma 9A
    Use as setas para navegar, 'enter' para editar a nota. Digite 'q' para sair.
    
    ALUNO                   NOTA
    -------------------     ----
    › Ana Beatriz Costa     8.5
      Bruno Dias            7.0
      Carla Esteves         [Pendente]
      Daniel Mendes (Transferido) --
    
    ---
    (Ao pressionar 'enter' em 'Carla Esteves')
    ---
    
    Lançamento de notas para a avaliação "Prova Bimestral 1" - Turma 9A
    Use as setas para navegar, 'enter' para editar a nota. Digite 'q' para sair.
    
    ALUNO                   NOTA
    -------------------     ----
      Ana Beatriz Costa     8.5
      Bruno Dias            7.0
    › Carla Esteves         › 9.0_
      Daniel Mendes (Transferido) --
    
    ```
    
-   **Ficheiro 5: `golden_files/relatorio_progresso_turma.txt`**
    
    ```
    =================================================
    ==       RELATÓRIO DE PROGRESSO - TURMA 9A     ==
    =================================================
    
    MÉDIA GERAL DA TURMA (Alunos Ativos): 8.2
    
    DESEMPENHO POR AVALIAÇÃO:
    - Prova Bimestral 1 (Peso 4): Média 8.5
    - Trabalho de Pesquisa (Peso 3): Média 7.8
    - Apresentação Oral (Peso 3): Média 8.3
    
    ALUNOS COM MAIOR DESEMPENHO:
    1. Carla Esteves (9.1)
    2. Ana Beatriz Costa (8.9)
    
    ALUNOS QUE NECESSITAM DE ATENÇÃO:
    1. Felipe Martins (6.5)
    2. Laura Santos (6.8)
    
    ```
    

### **Artefacto 8: Manifesto de Dependências e Ambiente de Desenvolvimento**

-   **Ambiente de Desenvolvimento Go:**
    
    -   **Instalação Requerida:** O ambiente de execução deve ter a versão **Go 1.18** instalada e configurada corretamente no `PATH` do sistema.
        
    -   **Verificação:** Antes de iniciar a implementação, executar o comando `go version` para confirmar que a saída corresponde à versão `go1.18`.
        
    -   **Dependências do Sistema:** Para a compilação cruzada (cross-compilation), pode ser necessário um compilador C (como GCC) para a dependência `go-sqlite3`. O ambiente deve estar preparado para isso.
        
-   **Ficheiro: `go.mod`**
    
    ```
    module vigenda
    
    go 1.18
    
    require (
        github.com/charmbracelet/bubbles v0.18.0
        github.com/charmbracelet/bubbletea v0.26.4
        github.com/charmbracelet/lipgloss v0.11.0
        github.com/mattn/go-sqlite3 v1.14.22
        github.com/spf13/cobra v1.8.0
        github.com/BurntSushi/toml v1.3.2
    )
    
    ```
    

### **Artefacto 9: Fluxos de Trabalho de Entrada de Dados**

Este artefacto detalha como as tarefas de entrada de dados complexas serão simplificadas para o utilizador.

-   **9.1. Inserção de Alunos em Turmas**
    
    -   **Problema:** Inserir 30 alunos um por um na linha de comando é impraticável.
        
    -   **Solução:** Importação em massa via CSV.
        
    -   **Comando:** `vigenda turma importar-alunos --id <id_da_turma> --arquivo /caminho/para/alunos.csv`
        
    -   **Estrutura do `alunos.csv`:**
        
        ```
        numero_chamada,nome_completo,situacao
        1,"Ana Beatriz Costa","ativo"
        2,"Bruno Dias","ativo"
        3,"Daniel Mendes","transferido"
        
        ```
        
    -   **Fluxo de Trabalho:**
        
        1.  O professor cria uma turma com `vigenda turma criar`.
            
        2.  Exporta a lista de alunos do sistema da escola para uma folha de cálculo.
            
        3.  Formata a folha para ter as colunas: `numero_chamada`, `nome_completo` e `situacao` (opcional, padrão 'ativo').
            
        4.  Salva como `alunos.csv`.
            
        5.  Executa o comando de importação. A aplicação processa o ficheiro e adiciona todos os alunos à turma de uma só vez, definindo o seu status.
            
-   **9.2. Criação de Avaliações**
    
    -   **Problema:** Definir uma nova avaliação com nome, peso e período requer vários dados.
        
    -   **Solução:** Um assistente interativo.
        
    -   **Comando:** `vigenda avaliacao criar`
        
    -   **Fluxo de Trabalho Interativo:**
        
        ```
        $ vigenda avaliacao criar
        ? Qual o nome da avaliação? › Prova Bimestral 2
        ? Para qual turma? (Use setas) › Turma 9A - 2025
        ? A qual período (bimestre) pertence? (1-4) › 2
        ? Qual o peso desta avaliação na média final? (ex: 4.0) › 4.0
        ✔ Avaliação "Prova Bimestral 2" criada com sucesso para a Turma 9A!
        
        ```
        
-   **9.3. Inserção de Questões no Banco**
    
    -   **Problema:** Questões, especialmente de múltipla escolha, são estruturas de dados complexas.
        
    -   **Solução:** Importação em massa via JSON.
        
    -   **Comando:** `vigenda bancoq add --arquivo /caminho/para/questoes.json`
        
    -   **Estrutura do `questoes.json`:**
        
        ```
        [
          {
            "disciplina": "História",
            "topico": "Revolução Francesa",
            "tipo": "multipla_escolha",
            "dificuldade": "media",
            "enunciado": "Qual destes eventos é considerado o estopim da Revolução Francesa?",
            "opcoes": [
              "A Queda da Bastilha",
              "A convocação dos Estados Gerais",
              "O Juramento da Quadra de Tênis"
            ],
            "resposta_correta": "A Queda da Bastilha"
          }
        ]
        
        ```
        
    -   **Fluxo de Trabalho:** O professor utiliza o seu editor de código ou de texto preferido para criar e gerir os seus ficheiros `.json` de questões e depois importa-os para a `Vigenda` com um único comando.
        
-   **9.4. Gestão do Status do Aluno**
    
    -   **Problema:** Um aluno é transferido no meio do ano e não deve mais aparecer nas listas de lançamento de notas, mas o seu histórico deve ser mantido.
        
    -   **Solução:** Um comando para atualizar o status do aluno de forma individual.
        
    -   **Comando:** `vigenda turma atualizar-status --aluno <id_do_aluno> --status inativo`
        
    -   **Fluxo de Trabalho:** O professor usa este comando para marcar alunos como `inativo` ou `transferido`. Estes alunos não aparecerão em novas listas de lançamento de notas, mas as suas notas e o seu registo permanecerão no sistema para consulta histórica.
