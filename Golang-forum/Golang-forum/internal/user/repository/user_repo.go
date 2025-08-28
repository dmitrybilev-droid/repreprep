package repository

import (
	"context"
	"database/sql"

	"golang-forum/internal/user/models"
	"github.com/rs/zerolog"
)

type UserRepository struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewUserRepository(db *sql.DB, logger zerolog.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

func (r *UserRepository) CreateUserRepo(ctx context.Context, email, username, passwordHash string, isAdmin bool) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, username, password_hash, is_admin, created_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)`,
		email, username, passwordHash, isAdmin,
	)
	r.logger.Info().Msgf("User created with email: %s, username: %s, isAdmin: %t", email, username, isAdmin)
	return err
}

func (r *UserRepository) GetUserByEmailRepo(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, username, password_hash, is_admin, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to get user by email: %s", email)
		return nil, err
	}
	r.logger.Info().Msgf("User get by email: %s", email)
	return &user, nil
}

func (r *UserRepository) GetUserByIDRepo(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, username, password_hash, is_admin, created_at FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to get user by id")
		return nil, err
	}
	r.logger.Info().Msgf("User get by id: %d", userID)
	return &user, nil
}

