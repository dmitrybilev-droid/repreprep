package repository

import (
	"context"
	"database/sql"
	"time"

	"golang-forum/internal/chat/models"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) SaveMessageRepo(ctx context.Context, msg *models.Message) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO messages (user_id, text_message) VALUES ($1, $2) RETURNING id, created_at`,
		msg.UserID, msg.Text,
	).Scan(&msg.ID, &msg.CreatedAt)
}

func (r *ChatRepository) DeleteMessageRepo(ctx context.Context, msg *models.Message) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM messages WHERE id = $1`,
		msg.ID,
	)
	return err
}

func (r *ChatRepository) GetMessagesRepo(ctx context.Context, limit int) ([]*models.Message, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, text_message, created_at 
         FROM messages 
         ORDER BY created_at DESC 
         LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Text, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (r *ChatRepository) GetMessageByIDRepo(ctx context.Context, ID int) (*models.Message, error) {
	var msg models.Message
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, text_message, created_at 
         FROM messages 
         WHERE id = $1`,
		ID,
	).Scan(&msg.ID, &msg.UserID, &msg.Text, &msg.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *ChatRepository) DeleteOldMessagesRepo(ctx context.Context, olderThan time.Duration) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM messages 
         WHERE created_at < $1`,
		time.Now().Add(-olderThan),
	)
	return err
}
