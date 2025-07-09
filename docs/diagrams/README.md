# Diagramas de Casos de Uso do Sistema Vigenda

Este diretório contém os diagramas de casos de uso para o sistema Vigenda, divididos em pacotes funcionais para melhor clareza e organização.

O ator principal em todos os diagramas é o **Usuário**.

## Diagramas Individuais

1.  **[Autenticação](./use_case_auth.puml)**
    *   Descreve os casos de uso relacionados ao processo de autenticação do usuário no sistema, como registro e login.
    *   [Visualizar `use_case_auth.puml`](./use_case_auth.puml)

2.  **[Gerenciamento Acadêmico](./use_case_academic_management.puml)**
    *   Detalha as funcionalidades centrais de gerenciamento de disciplinas, turmas, alunos e planos de aula.
    *   [Visualizar `use_case_academic_management.puml`](./use_case_academic_management.puml)

3.  **[Avaliação e Desempenho](./use_case_assessment_performance.puml)**
    *   Cobre os casos de uso relacionados à criação e gerenciamento de avaliações, banco de questões, geração de provas, lançamento de notas e acompanhamento do progresso dos alunos.
    *   [Visualizar `use_case_assessment_performance.puml`](./use_case_assessment_performance.puml)

4.  **[Ferramentas de Produtividade](./use_case_productivity_tools.puml)**
    *   Apresenta as funcionalidades voltadas para a produtividade do usuário, como gerenciamento de tarefas pessoais ou acadêmicas e sessões de foco.
    *   [Visualizar `use_case_productivity_tools.puml`](./use_case_productivity_tools.puml)

## Visão Geral Simplificada (Relação entre Pacotes)

Para entender como esses pacotes se conectam em um nível mais alto, considere o seguinte diagrama de visão geral simplificado. Ele omite os detalhes `<<include>>` e foca nas principais interações entre os grandes blocos funcionais.

```plantuml
@startuml
left to right direction
skinparam packageStyle rectangle

actor Usuário

rectangle "Sistema Vigenda" {

  package "Autenticação" {
    usecase "Autenticar Usuário" as UC_Auth
  }

  package "Gerenciamento Acadêmico" {
    usecase "Gerenciar Disciplinas" as UC_Subjects
    usecase "Gerenciar Turmas" as UC_Classes
    usecase "Gerenciar Alunos" as UC_Students
    usecase "Gerenciar Planos de Aula" as UC_Lessons
  }

  package "Avaliação e Desempenho" {
    usecase "Gerenciar Avaliações" as UC_Assessments
    usecase "Gerenciar Banco de Questões" as UC_Questions
    usecase "Gerar Provas" as UC_Proofs
    usecase "Gerenciar Notas" as UC_Grades
    usecase "Acompanhar Progresso" as UC_Progress
  }

  package "Ferramentas de Produtividade" {
    usecase "Gerenciar Tarefas" as UC_Tasks
    usecase "Gerenciar Sessões de Foco" as UC_Focus
  }

  Usuário -- UC_Auth
  Usuário -- UC_Subjects
  Usuário -- UC_Classes
  Usuário -- UC_Students
  Usuário -- UC_Lessons
  Usuário -- UC_Assessments
  Usuário -- UC_Questions
  Usuário -- UC_Proofs
  Usuário -- UC_Grades
  Usuário -- UC_Progress
  Usuário -- UC_Tasks
  Usuário -- UC_Focus

  UC_Classes ..> UC_Subjects : depende de
  UC_Students ..> UC_Classes : opera em
  UC_Lessons ..> UC_Classes : opera em

  UC_Assessments ..> UC_Classes : opera em
  UC_Questions ..> UC_Subjects : associado a
  UC_Proofs ..> UC_Questions : utiliza
  UC_Proofs ..> UC_Assessments : pode gerar
  UC_Grades ..> UC_Students : para
  UC_Grades ..> UC_Assessments : referente a
  UC_Progress ..> UC_Grades : analisa
  UC_Progress ..> UC_Classes : para

  UC_Focus ..> UC_Tasks : <<extend>>
}
@enduml
```

**Como Visualizar os Diagramas `.puml`:**

Copie o conteúdo de qualquer arquivo `.puml` e cole-o em um renderizador PlantUML online, como:
*   [PlantUML Online Server](http://www.plantuml.com/plantuml)
*   Ou utilize uma extensão PlantUML no seu editor de código (ex: VS Code).

Estes diagramas fornecem uma visão estruturada das funcionalidades do sistema Vigenda sob a perspectiva do usuário.
