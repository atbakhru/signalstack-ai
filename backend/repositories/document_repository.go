package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"signalstack-ai/backend/models"
)

type DocumentRepository struct {
	Pool *pgxpool.Pool
}

type DocumentListFilter struct {
	SourceID *int64
	Search   string
	Limit    int
	Offset   int
}

func (r DocumentRepository) Upsert(ctx context.Context, sourceID int64, document models.NormalizedDocument, rawS3Key, contentHash string) (models.Document, bool, error) {
	var saved models.Document
	var inserted bool
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO documents (
			source_id, external_id, title, summary, url, author, published_at, raw_s3_key, content_hash
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (source_id, external_id) DO UPDATE
		SET title = EXCLUDED.title,
			summary = EXCLUDED.summary,
			url = EXCLUDED.url,
			author = EXCLUDED.author,
			published_at = EXCLUDED.published_at,
			raw_s3_key = EXCLUDED.raw_s3_key,
			content_hash = EXCLUDED.content_hash,
			updated_at = NOW()
		RETURNING id, source_id, external_id, title, summary, url, author, published_at, raw_s3_key, content_hash, created_at, updated_at, (xmax = 0)
	`, sourceID, document.ExternalID, document.Title, document.Summary, document.URL, document.Author, document.PublishedAt, rawS3Key, contentHash).Scan(
		&saved.ID, &saved.SourceID, &saved.ExternalID, &saved.Title, &saved.Summary, &saved.URL, &saved.Author, &saved.PublishedAt, &saved.RawS3Key, &saved.ContentHash, &saved.CreatedAt, &saved.UpdatedAt, &inserted,
	)
	return saved, inserted, err
}

func (r DocumentRepository) GetByID(ctx context.Context, id int64) (models.Document, error) {
	var document models.Document
	err := r.Pool.QueryRow(ctx, `
		SELECT id, source_id, external_id, title, summary, url, author, published_at, raw_s3_key, content_hash, created_at, updated_at
		FROM documents
		WHERE id = $1
	`, id).Scan(&document.ID, &document.SourceID, &document.ExternalID, &document.Title, &document.Summary, &document.URL, &document.Author, &document.PublishedAt, &document.RawS3Key, &document.ContentHash, &document.CreatedAt, &document.UpdatedAt)
	return document, err
}

func (r DocumentRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.Pool.Exec(ctx, `DELETE FROM documents WHERE id = $1`, id)
	return err
}

func (r DocumentRepository) List(ctx context.Context, filter DocumentListFilter) ([]models.Document, error) {
	query := strings.Builder{}
	query.WriteString(`
		SELECT id, source_id, external_id, title, summary, url, author, published_at, raw_s3_key, content_hash, created_at, updated_at
		FROM documents
		WHERE 1=1
	`)
	args := make([]any, 0)
	argIndex := 1

	if filter.SourceID != nil {
		query.WriteString(fmt.Sprintf(" AND source_id = $%d", argIndex))
		args = append(args, *filter.SourceID)
		argIndex++
	}
	if filter.Search != "" {
		query.WriteString(fmt.Sprintf(" AND (title ILIKE $%d OR COALESCE(summary, '') ILIKE $%d OR COALESCE(author, '') ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}
	query.WriteString(" ORDER BY published_at DESC NULLS LAST, id DESC")
	if filter.Limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT $%d", argIndex))
		args = append(args, filter.Limit)
		argIndex++
	}
	if filter.Offset > 0 {
		query.WriteString(fmt.Sprintf(" OFFSET $%d", argIndex))
		args = append(args, filter.Offset)
	}

	rows, err := r.Pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := make([]models.Document, 0)
	for rows.Next() {
		var document models.Document
		if err := rows.Scan(&document.ID, &document.SourceID, &document.ExternalID, &document.Title, &document.Summary, &document.URL, &document.Author, &document.PublishedAt, &document.RawS3Key, &document.ContentHash, &document.CreatedAt, &document.UpdatedAt); err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	return documents, rows.Err()
}
