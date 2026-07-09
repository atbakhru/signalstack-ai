package models

import "time"

type EvaluationResult struct {
	ID                   int64     `json:"id"`
	QuestionID           int64     `json:"question_id"`
	RetrievedDocumentIDs []int64   `json:"retrieved_document_ids"`
	TopKAccuracy         float64   `json:"top_k_accuracy"`
	AnswerLatencyMS      int       `json:"answer_latency_ms"`
	RetrievalLatencyMS   int       `json:"retrieval_latency_ms"`
	CitationCount        int       `json:"citation_count"`
	CreatedAt            time.Time `json:"created_at"`
}
