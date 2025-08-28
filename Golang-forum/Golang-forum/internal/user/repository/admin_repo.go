package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"golang-forum/internal/user/models"
	"github.com/rs/zerolog"
)

type AdminRepository struct {
	db       *sql.DB
	logger   zerolog.Logger
	userRepo *UserRepository
}

func NewAdminRepository(db *sql.DB, logger zerolog.Logger) *AdminRepository {
	return &AdminRepository{db: db, logger: logger, userRepo: NewUserRepository(db, logger)}
}

func (r *AdminRepository) GetAllUsersRepo(ctx context.Context) ([]*models.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, email, username, password_hash, is_admin, created_at FROM users`)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to get all users")
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt); err != nil {
			r.logger.Error().Err(err).Msg("Failed to scan user row")
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while iterating over users")
		return nil, err
	}
	if len(users) == 0 {
		r.logger.Info().Msg("No users found")
		return nil, nil
	}
	r.logger.Info().Msgf("Total users found: %d", len(users))
	return users, nil
}

func (r *AdminRepository) GetAllBannedUsersRepo(ctx context.Context) ([]*models.BannedUser, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, email, username, banned_at FROM banned_users`)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to get all banned users")
		return nil, err
	}
	defer rows.Close()

	var bannedUsers []*models.BannedUser
	for rows.Next() {
		var bannedUser models.BannedUser
		if err := rows.Scan(&bannedUser.ID, &bannedUser.Email, &bannedUser.Username, &bannedUser.BannedAt); err != nil {
			r.logger.Error().Err(err).Msg("Failed to scan banned user row")
			return nil, err
		}
		bannedUsers = append(bannedUsers, &bannedUser)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while iterating over banned users")
		return nil, err
	}
	if len(bannedUsers) == 0 {
		r.logger.Info().Msg("No banned users found")
		return nil, nil
	}
	r.logger.Info().Msgf("Total banned users found: %d", len(bannedUsers))
	return bannedUsers, nil
}

func (r *AdminRepository) GetBannedUserByIDRepo(ctx context.Context, userID int) (*models.BannedUser, error) {
	var bannedUser models.BannedUser
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, username, banned_at FROM banned_users WHERE id = $1`,
		userID,
	).Scan(&bannedUser.ID, &bannedUser.Email, &bannedUser.Username, &bannedUser.BannedAt)
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to get banned user by id")
		return nil, err
	}
	r.logger.Info().Msgf("Banned user found: %s", bannedUser.Email)
	return &bannedUser, nil
}

func (r *AdminRepository) GetBannedUserByEmailRepo(ctx context.Context, email string) (*models.BannedUser, error) {
	var bannedUser models.BannedUser
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, banned_at FROM banned_users WHERE email = $1`,
		email,
	).Scan(&bannedUser.ID, &bannedUser.Email, &bannedUser.BannedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	r.logger.Info().Msgf("Banned user found: %s", bannedUser.Email)
	return &bannedUser, nil
}

func (r *AdminRepository) BanUserRepo(ctx context.Context, email string) error {
	user, err := r.userRepo.GetUserByEmailRepo(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	_, err = r.db.ExecContext(ctx,
		`INSERT INTO banned_users (id, email, username, banned_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`,
		user.ID, user.Email, user.Username,
	)
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to ban user: %s", email)
		return fmt.Errorf("failed to ban user: %v", err)
	}
	r.logger.Info().Msgf("User banned successfully: %s", email)
	return err
}

func (r *AdminRepository) UnBanUserRepo(ctx context.Context, email string) error {
	user, err := r.userRepo.GetUserByEmailRepo(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	_, err = r.db.ExecContext(ctx,
		`DELETE FROM banned_users WHERE email = $1`,
		user.Email,
	)
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to unban user: %s", email)
		return fmt.Errorf("failed to unban user: %v", err)
	}
	r.logger.Info().Msgf("User unbanned successfully: %s", email)
	return err
}

func (r *AdminRepository) DeleteUserRepo(ctx context.Context, email string) error {
	user, err := r.userRepo.GetUserByEmailRepo(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	_, err = r.db.ExecContext(ctx,
		`DELETE FROM users WHERE email = $1`,
		user.Email,
	)
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to delete user: %s", email)
	}
	return err
}
