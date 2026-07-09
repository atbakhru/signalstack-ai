package models

import "time"

type DocumentChunk struct {
	ID         int64     `json:"id"`
	DocumentID int64     `json:"document_id"`
	SourceID   int64     `json:"source_id"`
	ChunkIndex int       `json:"chunk_index"`
	Content    string    `json:"content"`
	TokenCount int       `json:"token_count"`
	Embedding  []float32 `json:"embedding,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
