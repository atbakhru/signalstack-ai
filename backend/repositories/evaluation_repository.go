package repositories

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"

	"signalstack-ai/backend/models"
)

type EvaluationRepository struct {
	Pool *pgxpool.Pool
}

func (r EvaluationRepository) ListQuestions(ctx context.Context) ([]models.EvaluationQuestion, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, question, expected_source, expected_keywords, created_at
		FROM evaluation_questions
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := make([]models.EvaluationQuestion, 0)
	for rows.Next() {
		var question models.EvaluationQuestion
		if err := rows.Scan(&question.ID, &question.Question, &question.ExpectedSource, &question.ExpectedKeywords, &question.CreatedAt); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, rows.Err()
}

func (r EvaluationRepository) ListResults(ctx context.Context) ([]models.EvaluationResult, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, question_id, retrieved_document_ids, top_k_accuracy, answer_latency_ms, retrieval_latency_ms, citation_count, created_at
		FROM evaluation_results
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]models.EvaluationResult, 0)
	for rows.Next() {
		var result models.EvaluationResult
		var retrieved []byte
		if err := rows.Scan(&result.ID, &result.QuestionID, &retrieved, &result.TopKAccuracy, &result.AnswerLatencyMS, &result.RetrievalLatencyMS, &result.CitationCount, &result.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(retrieved, &result.RetrievedDocumentIDs)
		results = append(results, result)
	}
	return results, rows.Err()
}

func (r EvaluationRepository) CreateResult(ctx context.Context, result models.EvaluationResult) (models.EvaluationResult, error) {
	var saved models.EvaluationResult
	encoded, err := json.Marshal(result.RetrievedDocumentIDs)
	if err != nil {
		return saved, err
	}
	err = r.Pool.QueryRow(ctx, `
		INSERT INTO evaluation_results (question_id, retrieved_document_ids, top_k_accuracy, answer_latency_ms, retrieval_latency_ms, citation_count)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, question_id, retrieved_document_ids, top_k_accuracy, answer_latency_ms, retrieval_latency_ms, citation_count, created_at
	`, result.QuestionID, encoded, result.TopKAccuracy, result.AnswerLatencyMS, result.RetrievalLatencyMS, result.CitationCount).Scan(&saved.ID, &saved.QuestionID, &encoded, &saved.TopKAccuracy, &saved.AnswerLatencyMS, &saved.RetrievalLatencyMS, &saved.CitationCount, &saved.CreatedAt)
	if err != nil {
		return saved, err
	}
	_ = json.Unmarshal(encoded, &saved.RetrievedDocumentIDs)
	return saved, nil
}
