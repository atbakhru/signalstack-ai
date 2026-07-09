package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"signalstack-ai/backend/models"
)

type ChatRepository struct {
	Pool *pgxpool.Pool
}

func (r ChatRepository) CreateSession(ctx context.Context, session models.ChatSession) (models.ChatSession, error) {
	var saved models.ChatSession
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO chat_sessions (user_id, title)
		VALUES ($1, $2)
		RETURNING id, user_id, title, created_at
	`, session.UserID, session.Title).Scan(&saved.ID, &saved.UserID, &saved.Title, &saved.CreatedAt)
	return saved, err
}

func (r ChatRepository) ListSessions(ctx context.Context, userID int64) ([]models.ChatSession, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, user_id, title, created_at
		FROM chat_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]models.ChatSession, 0)
	for rows.Next() {
		var session models.ChatSession
		if err := rows.Scan(&session.ID, &session.UserID, &session.Title, &session.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, rows.Err()
}

func (r ChatRepository) GetSessionByID(ctx context.Context, id int64) (models.ChatSession, error) {
	var session models.ChatSession
	err := r.Pool.QueryRow(ctx, `
		SELECT id, user_id, title, created_at
		FROM chat_sessions
		WHERE id = $1
	`, id).Scan(&session.ID, &session.UserID, &session.Title, &session.CreatedAt)
	return session, err
}

func (r ChatRepository) DeleteSession(ctx context.Context, id, userID int64) error {
	return r.Pool.QueryRow(ctx, `
		DELETE FROM chat_sessions
		WHERE id = $1 AND user_id = $2
		RETURNING id
	`, id, userID).Scan(new(int64))
}

func (r ChatRepository) ListMessages(ctx context.Context, sessionID int64) ([]models.ChatMessage, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, session_id, user_id, role, content, created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at ASC, id ASC
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]models.ChatMessage, 0)
	for rows.Next() {
		var message models.ChatMessage
		if err := rows.Scan(&message.ID, &message.SessionID, &message.UserID, &message.Role, &message.Content, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (r ChatRepository) CreateMessage(ctx context.Context, message models.ChatMessage) (models.ChatMessage, error) {
	var saved models.ChatMessage
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO chat_messages (session_id, user_id, role, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, session_id, user_id, role, content, created_at
	`, message.SessionID, message.UserID, message.Role, message.Content).Scan(&saved.ID, &saved.SessionID, &saved.UserID, &saved.Role, &saved.Content, &saved.CreatedAt)
	return saved, err
}
