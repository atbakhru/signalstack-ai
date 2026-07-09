package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"signalstack-ai/backend/models"
)

type ChunkMatch struct {
	Chunk           models.DocumentChunk `json:"chunk"`
	DocumentTitle   string               `json:"document_title"`
	SourceName      string               `json:"source_name"`
	URL             string               `json:"url"`
	SimilarityScore float64              `json:"similarity_score"`
}

type ChunkRepository struct {
	Pool *pgxpool.Pool
}

func (r ChunkRepository) InsertMany(ctx context.Context, chunks []models.DocumentChunk) ([]models.DocumentChunk, error) {
	inserted := make([]models.DocumentChunk, 0, len(chunks))
	for _, chunk := range chunks {
		var saved models.DocumentChunk
		var embedding any
		if len(chunk.Embedding) > 0 {
			embedding = VectorLiteral(chunk.Embedding)
		}
		err := r.Pool.QueryRow(ctx, `
			INSERT INTO document_chunks (document_id, source_id, chunk_index, content, token_count, embedding)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (document_id, chunk_index) DO UPDATE
			SET content = EXCLUDED.content,
				token_count = EXCLUDED.token_count,
				embedding = EXCLUDED.embedding
			RETURNING id, document_id, source_id, chunk_index, content, token_count, created_at
		`, chunk.DocumentID, chunk.SourceID, chunk.ChunkIndex, chunk.Content, chunk.TokenCount, embedding).Scan(
			&saved.ID, &saved.DocumentID, &saved.SourceID, &saved.ChunkIndex, &saved.Content, &saved.TokenCount, &saved.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		saved.Embedding = chunk.Embedding
		inserted = append(inserted, saved)
	}
	return inserted, nil
}

func (r ChunkRepository) ListByDocumentID(ctx context.Context, documentID int64) ([]models.DocumentChunk, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, document_id, source_id, chunk_index, content, token_count, created_at
		FROM document_chunks
		WHERE document_id = $1
		ORDER BY chunk_index
	`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chunks := make([]models.DocumentChunk, 0)
	for rows.Next() {
		var chunk models.DocumentChunk
		if err := rows.Scan(&chunk.ID, &chunk.DocumentID, &chunk.SourceID, &chunk.ChunkIndex, &chunk.Content, &chunk.TokenCount, &chunk.CreatedAt); err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, rows.Err()
}

func (r ChunkRepository) SearchByEmbedding(ctx context.Context, queryEmbedding []float32, topK int, sourceIDs []int64) ([]models.DocumentChunk, error) {
	if topK <= 0 {
		topK = 5
	}

	query := strings.Builder{}
	query.WriteString(`
		SELECT id, document_id, source_id, chunk_index, content, token_count, created_at
		FROM document_chunks
		WHERE embedding IS NOT NULL
	`)
	args := []any{VectorLiteral(queryEmbedding)}
	argIndex := 2
	if len(sourceIDs) > 0 {
		query.WriteString(fmt.Sprintf(" AND source_id = ANY($%d)", argIndex))
		args = append(args, sourceIDs)
		argIndex++
	}
	query.WriteString(fmt.Sprintf(" ORDER BY embedding <=> $1::vector LIMIT $%d", argIndex))
	args = append(args, topK)

	rows, err := r.Pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chunks := make([]models.DocumentChunk, 0)
	for rows.Next() {
		var chunk models.DocumentChunk
		if err := rows.Scan(&chunk.ID, &chunk.DocumentID, &chunk.SourceID, &chunk.ChunkIndex, &chunk.Content, &chunk.TokenCount, &chunk.CreatedAt); err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, rows.Err()
}

func (r ChunkRepository) SearchByEmbeddingWithMetadata(ctx context.Context, queryEmbedding []float32, topK int, sourceIDs []int64) ([]ChunkMatch, error) {
	if topK <= 0 {
		topK = 5
	}

	query := strings.Builder{}
	query.WriteString(`
		SELECT dc.id, dc.document_id, dc.source_id, dc.chunk_index, dc.content, dc.token_count, dc.created_at,
		       d.title, s.name, COALESCE(d.url, ''), 1 - (dc.embedding <=> $1::vector) AS similarity_score
		FROM document_chunks dc
		JOIN documents d ON d.id = dc.document_id
		JOIN sources s ON s.id = dc.source_id
		WHERE dc.embedding IS NOT NULL
	`)
	args := []any{VectorLiteral(queryEmbedding)}
	argIndex := 2
	if len(sourceIDs) > 0 {
		query.WriteString(fmt.Sprintf(" AND dc.source_id = ANY($%d)", argIndex))
		args = append(args, sourceIDs)
		argIndex++
	}
	query.WriteString(fmt.Sprintf(" ORDER BY dc.embedding <=> $1::vector LIMIT $%d", argIndex))
	args = append(args, topK)

	rows, err := r.Pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]ChunkMatch, 0)
	for rows.Next() {
		var match ChunkMatch
		if err := rows.Scan(
			&match.Chunk.ID,
			&match.Chunk.DocumentID,
			&match.Chunk.SourceID,
			&match.Chunk.ChunkIndex,
			&match.Chunk.Content,
			&match.Chunk.TokenCount,
			&match.Chunk.CreatedAt,
			&match.DocumentTitle,
			&match.SourceName,
			&match.URL,
			&match.SimilarityScore,
		); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, rows.Err()
}
