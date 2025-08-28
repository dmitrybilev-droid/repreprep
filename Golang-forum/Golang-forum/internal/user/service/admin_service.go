package service

import (
	"context"

	chatgen "golang-forum/internal/chat/transport/grpc/gen/api"
	"golang-forum/internal/user/models"
	"golang-forum/internal/user/repository"
)

type AdminService struct {
	repo       *repository.AdminRepository
	chatClient chatgen.ChatServiceClient
}

func NewAdminService(repo *repository.AdminRepository, chatClient chatgen.ChatServiceClient) *AdminService {
	return &AdminService{
		repo:       repo,
		chatClient: chatClient,
	}
}
func (s *AdminService) UpdateChatClient(client chatgen.ChatServiceClient) {
	s.chatClient = client
}

func (s *AdminService) IsUserBanned(ctx context.Context, email string) (bool, error) {
	user, err := s.repo.GetBannedUserByEmailRepo(ctx, email)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return true, nil
}

func (s *AdminService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	users, err := s.repo.GetAllUsersRepo(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (s *AdminService) GetAllBannedUsers(ctx context.Context) ([]*models.BannedUser, error) {
	bannedUsers, err := s.repo.GetAllBannedUsersRepo(ctx)
	if err != nil {
		return nil, err
	}
	return bannedUsers, nil
}

func (s *AdminService) BanUser(ctx context.Context, email string) error {
	bannedUser, err := s.repo.GetBannedUserByEmailRepo(ctx, email)
	if err != nil {
		return err
	}
	if bannedUser != nil {
		return nil
	}

	return s.repo.BanUserRepo(ctx, email)
}

func (s *AdminService) UnBanUser(ctx context.Context, email string) error {
	bannedUser, err := s.repo.GetBannedUserByEmailRepo(ctx, email)
	if err != nil {
		return err
	}
	if bannedUser == nil {
		return nil
	}

	return s.repo.UnBanUserRepo(ctx, email)
}

func (s *AdminService) DeleteMessage(ctx context.Context, messageID int) error {
	_, err := s.chatClient.DeleteMessage(ctx, &chatgen.DeleteMessageRequest{MessageId: int32(messageID)})
	if err != nil {
		return err
	}
	return nil
}
