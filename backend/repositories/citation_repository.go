package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"signalstack-ai/backend/models"
)

type CitationRepository struct {
	Pool *pgxpool.Pool
}

func (r CitationRepository) CreateMany(ctx context.Context, citations []models.Citation) ([]models.Citation, error) {
	saved := make([]models.Citation, 0, len(citations))
	for _, citation := range citations {
		var inserted models.Citation
		err := r.Pool.QueryRow(ctx, `
			INSERT INTO citations (message_id, document_id, chunk_id, source_name, title, url, relevance_score)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, message_id, document_id, chunk_id, source_name, title, url, relevance_score, created_at
		`, citation.MessageID, citation.DocumentID, citation.ChunkID, citation.SourceName, citation.Title, citation.URL, citation.RelevanceScore).Scan(
			&inserted.ID, &inserted.MessageID, &inserted.DocumentID, &inserted.ChunkID, &inserted.SourceName, &inserted.Title, &inserted.URL, &inserted.RelevanceScore, &inserted.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		saved = append(saved, inserted)
	}
	return saved, nil
}
