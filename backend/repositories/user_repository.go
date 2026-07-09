package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"signalstack-ai/backend/models"
)

type UserRepository struct {
	Pool *pgxpool.Pool
}

func (r UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	var saved models.User
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password_hash, created_at
	`, user.Name, user.Email, user.PasswordHash).Scan(&saved.ID, &saved.Name, &saved.Email, &saved.PasswordHash, &saved.CreatedAt)
	return saved, err
}

func (r UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.Pool.QueryRow(ctx, `
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

func (r UserRepository) GetByID(ctx context.Context, id int64) (models.User, error) {
	var user models.User
	err := r.Pool.QueryRow(ctx, `
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

func IsNotFound(err error) bool {
	return err == pgx.ErrNoRows
}
