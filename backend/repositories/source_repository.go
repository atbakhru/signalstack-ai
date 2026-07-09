package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"signalstack-ai/backend/models"
)

type SourceRepository struct {
	Pool *pgxpool.Pool
}

func (r SourceRepository) List(ctx context.Context) ([]models.Source, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, name, api_name, base_url, enabled, created_at
		FROM sources
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSources(rows)
}

func (r SourceRepository) ListEnabled(ctx context.Context) ([]models.Source, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, name, api_name, base_url, enabled, created_at
		FROM sources
		WHERE enabled = TRUE
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSources(rows)
}

func (r SourceRepository) GetByID(ctx context.Context, id int64) (models.Source, error) {
	var source models.Source
	err := r.Pool.QueryRow(ctx, `
		SELECT id, name, api_name, base_url, enabled, created_at
		FROM sources
		WHERE id = $1
	`, id).Scan(&source.ID, &source.Name, &source.APIName, &source.BaseURL, &source.Enabled, &source.CreatedAt)
	return source, err
}

func (r SourceRepository) GetByName(ctx context.Context, name string) (models.Source, error) {
	var source models.Source
	err := r.Pool.QueryRow(ctx, `
		SELECT id, name, api_name, base_url, enabled, created_at
		FROM sources
		WHERE name = $1
	`, name).Scan(&source.ID, &source.Name, &source.APIName, &source.BaseURL, &source.Enabled, &source.CreatedAt)
	return source, err
}

func (r SourceRepository) ToggleEnabled(ctx context.Context, id int64) (models.Source, error) {
	var source models.Source
	err := r.Pool.QueryRow(ctx, `
		UPDATE sources
		SET enabled = NOT enabled
		WHERE id = $1
		RETURNING id, name, api_name, base_url, enabled, created_at
	`, id).Scan(&source.ID, &source.Name, &source.APIName, &source.BaseURL, &source.Enabled, &source.CreatedAt)
	return source, err
}

func scanSources(rows pgx.Rows) ([]models.Source, error) {
	sources := make([]models.Source, 0)
	for rows.Next() {
		var source models.Source
		if err := rows.Scan(&source.ID, &source.Name, &source.APIName, &source.BaseURL, &source.Enabled, &source.CreatedAt); err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	return sources, rows.Err()
}
