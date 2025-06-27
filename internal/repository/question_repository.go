package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"vigenda/internal/models"
)

type questionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) AddQuestion(ctx context.Context, question *models.Question) (int64, error) {
	query := `INSERT INTO questions (user_id, subject_id, topic, type, difficulty, statement, options, correct_answer, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()

	var optionsJSON sql.NullString
	if question.Options != nil && *question.Options != "" {
		// Assuming question.Options is already a JSON string. If it's a slice/map, marshal it here.
		// For safety, let's ensure it's a valid JSON string or handle potential errors.
		// The model currently has Options as *string, implying it could be pre-formatted JSON.
		optionsJSON.String = *question.Options
		optionsJSON.Valid = true
	}

	result, err := r.db.ExecContext(ctx, query,
		question.UserID,
		question.SubjectID,
		question.Topic,
		question.Type,
		question.Difficulty,
		question.Statement,
		optionsJSON,
		question.CorrectAnswer,
		now,
		now,
	)
	if err != nil {
		return 0, fmt.Errorf("questionRepository.AddQuestion: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("questionRepository.AddQuestion: failed to get last insert ID: %w", err)
	}
	return id, nil
}

func (r *questionRepository) GetQuestionsByCriteria(ctx context.Context, criteria QuestionQueryCriteria) ([]models.Question, error) {
	baseQuery := `SELECT id, user_id, subject_id, topic, type, difficulty, statement, options, correct_answer, created_at, updated_at
                  FROM questions WHERE subject_id = ?`
	args := []interface{}{criteria.SubjectID}

	if criteria.Topic != nil && *criteria.Topic != "" {
		baseQuery += " AND topic = ?"
		args = append(args, *criteria.Topic)
	}
	if criteria.Difficulty != "" {
		baseQuery += " AND difficulty = ?"
		args = append(args, criteria.Difficulty)
	}

	// Order by random if no specific order is needed, or add other ordering options
	baseQuery += " ORDER BY RANDOM()" // SQLite specific for random ordering

	if criteria.Limit > 0 {
		baseQuery += " LIMIT ?"
		args = append(args, criteria.Limit)
	}

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("questionRepository.GetQuestionsByCriteria: querying: %w", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		q := models.Question{}
		var topic sql.NullString
		var options sql.NullString // This will be a JSON string from the DB

		err := rows.Scan(
			&q.ID,
			&q.UserID,
			&q.SubjectID,
			&topic,
			&q.Type,
			&q.Difficulty,
			&q.Statement,
			&options,
			&q.CorrectAnswer,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("questionRepository.GetQuestionsByCriteria: scanning question: %w", err)
		}
		if topic.Valid {
			q.Topic = topic.String
		}
		if options.Valid {
			// The 'options' field in DB is TEXT, storing a JSON array string.
			// We need to ensure q.Options (which is *string) gets this JSON string.
			optStr := options.String
			q.Options = &optStr

			// If the model's Options field were []string, we would unmarshal here:
			// var opts []string
			// if err := json.Unmarshal([]byte(options.String), &opts); err == nil {
			//  q.Options = opts // Assuming q.Options is []string
			// } else {
			// Handle error or set to nil
			// }
		}
		questions = append(questions, q)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("questionRepository.GetQuestionsByCriteria: iterating rows: %w", err)
	}

	// If fewer questions were found than requested (specifically for difficulty-based generation),
	// this needs to be handled by the service layer. The repository just returns what it finds.
	return questions, nil
}

// Helper to marshal options if they are passed as a slice of strings
func marshalOptions(options []string) (sql.NullString, error) {
	if len(options) == 0 {
		return sql.NullString{}, nil
	}
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return sql.NullString{}, fmt.Errorf("failed to marshal options to JSON: %w", err)
	}
	return sql.NullString{String: string(jsonBytes), Valid: true}, nil
}

// Helper to unmarshal options if they are stored as JSON string in DB
func unmarshalOptions(jsonString sql.NullString) ([]string, error) {
	if !jsonString.Valid || jsonString.String == "" {
		return nil, nil
	}
	var options []string
	err := json.Unmarshal([]byte(jsonString.String), &options)
	if err != nil {
		// Try to unmarshal if it's a simple string rather than a JSON array string
		// This case might occur if data was inserted manually or incorrectly
		var singleOption string
		if errSingle := json.Unmarshal([]byte(jsonString.String), &singleOption); errSingle == nil {
			// It was a single string, e.g. "\"Option A\""
			// Check if it's a list of strings separated by a delimiter, e.g. "Option A,Option B"
			// This is heuristic and depends on how non-JSON data might have been stored.
			// For now, we assume it's a valid JSON array string or null.
			// If it's just a plain string, and not a JSON array, this unmarshal will fail.
			// The application should ensure 'options' are always stored as JSON array strings.
			return nil, fmt.Errorf("options field contains non-JSON-array string: %s. Error: %w", jsonString.String, err)

		}
		return nil, fmt.Errorf("failed to unmarshal options from JSON '%s': %w", jsonString.String, err)
	}
	return options, nil
}


// GetQuestionsByCriteria_ProofGeneration is a more specialized version for proof generation
// that fetches a specific number of questions for each difficulty level.
func (r *questionRepository) GetQuestionsByCriteriaProofGeneration(ctx context.Context, criteria service.ProofCriteria) ([]models.Question, error) {
	var allQuestions []models.Question

	difficulties := []struct {
		Level string
		Count int
	}{
		{"facil", criteria.EasyCount},
		{"media", criteria.MediumCount},
		{"dificil", criteria.HardCount},
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("questionRepository.GetQuestionsByCriteriaProofGeneration: beginning transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	for _, diff := range difficulties {
		if diff.Count == 0 {
			continue
		}

		queryBuilder := strings.Builder{}
		queryBuilder.WriteString(`SELECT id, user_id, subject_id, topic, type, difficulty, statement, options, correct_answer, created_at, updated_at
                                 FROM questions WHERE subject_id = ? AND difficulty = ?`)
		args := []interface{}{criteria.SubjectID, diff.Level}

		if criteria.Topic != nil && *criteria.Topic != "" {
			queryBuilder.WriteString(" AND topic = ?")
			args = append(args, *criteria.Topic)
		}
		queryBuilder.WriteString(" ORDER BY RANDOM() LIMIT ?")
		args = append(args, diff.Count)

		rows, err := tx.QueryContext(ctx, queryBuilder.String(), args...)
		if err != nil {
			return nil, fmt.Errorf("questionRepository.GetQuestionsByCriteriaProofGeneration: querying for %s: %w", diff.Level, err)
		}

		currentDifficultyQuestions := 0
		for rows.Next() {
			q := models.Question{}
			var topic sql.NullString
			var options sql.NullString

			errScan := rows.Scan(
				&q.ID, &q.UserID, &q.SubjectID, &topic, &q.Type, &q.Difficulty,
				&q.Statement, &options, &q.CorrectAnswer, &q.CreatedAt, &q.UpdatedAt,
			)
			if errScan != nil {
				rows.Close() // Important to close rows before returning from loop
				return nil, fmt.Errorf("questionRepository.GetQuestionsByCriteriaProofGeneration: scanning for %s: %w", diff.Level, errScan)
			}
			if topic.Valid {
				q.Topic = topic.String
			}
			if options.Valid {
				optStr := options.String
				q.Options = &optStr
			}
			allQuestions = append(allQuestions, q)
			currentDifficultyQuestions++
		}
		if errRows := rows.Err(); errRows != nil {
			rows.Close()
			return nil, fmt.Errorf("questionRepository.GetQuestionsByCriteriaProofGeneration: iterating rows for %s: %w", diff.Level, errRows)
		}
		rows.Close() // Close rows for this difficulty level

		if currentDifficultyQuestions < diff.Count {
			// Not enough questions for this difficulty.
			// The service layer will decide how to handle this (e.g., return error, or fewer questions).
			// This repository returns what it finds.
			// We can add a log or a specific error type if needed.
			// For now, the service.GenerateProof will check the total count.
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("questionRepository.GetQuestionsByCriteriaProofGeneration: committing transaction: %w", err)
	}

	return allQuestions, nil
}
