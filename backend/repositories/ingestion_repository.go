package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"signalstack-ai/backend/models"
)

type IngestionRepository struct {
	Pool *pgxpool.Pool
}

func (r IngestionRepository) CreateRun(ctx context.Context, sourceID int64) (models.IngestionRun, error) {
	var run models.IngestionRun
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO ingestion_runs (source_id, status)
		VALUES ($1, 'running')
		RETURNING id, source_id, status, documents_fetched, documents_inserted, chunks_created, embeddings_created, started_at, completed_at, error_message
	`, sourceID).Scan(&run.ID, &run.SourceID, &run.Status, &run.DocumentsFetched, &run.DocumentsInserted, &run.ChunksCreated, &run.EmbeddingsCreated, &run.StartedAt, &run.CompletedAt, &run.ErrorMessage)
	return run, err
}

func (r IngestionRepository) FinishRun(ctx context.Context, runID int64, status string, documentsFetched, documentsInserted, chunksCreated, embeddingsCreated int, errorMessage string) (models.IngestionRun, error) {
	var run models.IngestionRun
	completedAt := time.Now().UTC()
	err := r.Pool.QueryRow(ctx, `
		UPDATE ingestion_runs
		SET status = $2,
			documents_fetched = $3,
			documents_inserted = $4,
			chunks_created = $5,
			embeddings_created = $6,
			completed_at = $7,
			error_message = $8
		WHERE id = $1
		RETURNING id, source_id, status, documents_fetched, documents_inserted, chunks_created, embeddings_created, started_at, completed_at, error_message
	`, runID, status, documentsFetched, documentsInserted, chunksCreated, embeddingsCreated, completedAt, errorMessage).Scan(&run.ID, &run.SourceID, &run.Status, &run.DocumentsFetched, &run.DocumentsInserted, &run.ChunksCreated, &run.EmbeddingsCreated, &run.StartedAt, &run.CompletedAt, &run.ErrorMessage)
	return run, err
}

func (r IngestionRepository) ListRuns(ctx context.Context) ([]models.IngestionRun, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, source_id, status, documents_fetched, documents_inserted, chunks_created, embeddings_created, started_at, completed_at, error_message
		FROM ingestion_runs
		ORDER BY started_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	runs := make([]models.IngestionRun, 0)
	for rows.Next() {
		var run models.IngestionRun
		if err := rows.Scan(&run.ID, &run.SourceID, &run.Status, &run.DocumentsFetched, &run.DocumentsInserted, &run.ChunksCreated, &run.EmbeddingsCreated, &run.StartedAt, &run.CompletedAt, &run.ErrorMessage); err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, rows.Err()
}

func (r IngestionRepository) GetRunByID(ctx context.Context, id int64) (models.IngestionRun, error) {
	var run models.IngestionRun
	err := r.Pool.QueryRow(ctx, `
		SELECT id, source_id, status, documents_fetched, documents_inserted, chunks_created, embeddings_created, started_at, completed_at, error_message
		FROM ingestion_runs
		WHERE id = $1
	`, id).Scan(&run.ID, &run.SourceID, &run.Status, &run.DocumentsFetched, &run.DocumentsInserted, &run.ChunksCreated, &run.EmbeddingsCreated, &run.StartedAt, &run.CompletedAt, &run.ErrorMessage)
	return run, err
}
