package models

import "time"

type IngestionRun struct {
	ID                int64      `json:"id"`
	SourceID          int64      `json:"source_id"`
	Status            string     `json:"status"`
	DocumentsFetched  int        `json:"documents_fetched"`
	DocumentsInserted int        `json:"documents_inserted"`
	ChunksCreated     int        `json:"chunks_created"`
	EmbeddingsCreated int        `json:"embeddings_created"`
	StartedAt         time.Time  `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	ErrorMessage      string     `json:"error_message,omitempty"`
}
