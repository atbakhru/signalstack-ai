package models

import "time"

type Citation struct {
	ID             int64     `json:"id"`
	MessageID      int64     `json:"message_id"`
	DocumentID     int64     `json:"document_id"`
	ChunkID        int64     `json:"chunk_id"`
	SourceName     string    `json:"source_name"`
	Title          string    `json:"title"`
	URL            string    `json:"url"`
	RelevanceScore float64   `json:"relevance_score"`
	CreatedAt      time.Time `json:"created_at"`
}
