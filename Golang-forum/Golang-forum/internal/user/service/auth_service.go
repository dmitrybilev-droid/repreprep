package service

import (
	"context"
	"errors"
	"fmt"

	"golang-forum/internal/user/models"
	"golang-forum/internal/user/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string, isAdmin bool) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	return s.repo.CreateUserRepo(ctx, email, username, string(hashedPassword), isAdmin)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmailRepo(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	return s.repo.GetUserByIDRepo(ctx, userID)
}
