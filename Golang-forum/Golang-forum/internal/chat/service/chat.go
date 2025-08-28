package service

import (
	"context"
	"errors"
	"time"

	"golang-forum/internal/chat/models"
	"golang-forum/internal/chat/repository"
	user "golang-forum/internal/user/transport/grpc/gen/api"
)

type ChatService struct {
	repo      *repository.ChatRepository
	msgChan   chan *models.Message
	authGrpc  user.AuthServiceClient
	adminGrpc user.AdminServiceClient
}

func (s *ChatService) GetUserInfoByID(ctx context.Context, userID int) (*user.GetUserInfoByIDResponse, error) {
	return s.authGrpc.GetUserInfoByID(ctx, &user.GetUserInfoByIDRequest{UserId: int32(userID)})
}

func NewChatService(repo *repository.ChatRepository, authGrpc user.AuthServiceClient, adminGrpc user.AdminServiceClient) *ChatService {
	svc := &ChatService{
		repo:      repo,
		msgChan:   make(chan *models.Message, 100),
		authGrpc:  authGrpc,
		adminGrpc: adminGrpc,
	}
	go svc.StartMessageCleaner(24 * time.Hour)
	return svc
}

func (s *ChatService) SendMessage(ctx context.Context, userID int, text string) (*models.Message, error) {
	// Получаем email пользователя по userID
	userInfo, err := s.authGrpc.GetUserInfoByID(ctx, &user.GetUserInfoByIDRequest{UserId: int32(userID)})
	if err != nil {
		return nil, errors.New("failed to get user info")
	}

	// Проверяем бан по email
	isBannedResp, err := s.adminGrpc.IsBanned(ctx, &user.IsBannedRequest{Email: userInfo.Email})
	if err != nil {
		return nil, err
	}
	if isBannedResp.Banned {
		return nil, errors.New("user is banned")
	}

	message := &models.Message{
		UserID:    userID,
		Text:      text,
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveMessageRepo(ctx, message); err != nil {
		return nil, err
	}

	fullMsg, err := s.repo.GetMessageByIDRepo(ctx, message.ID)
	if err != nil {
		return nil, err
	}

	s.msgChan <- fullMsg
	return fullMsg, nil
}

func (s *ChatService) GetMessages(ctx context.Context, limit int) ([]*models.MessageResponse, error) {
	messages, err := s.repo.GetMessagesRepo(ctx, limit)
	if err != nil {
		return nil, err
	}

	var responses []*models.MessageResponse
	for _, msg := range messages {
		userInfo, err := s.authGrpc.GetUserInfoByID(ctx, &user.GetUserInfoByIDRequest{UserId: int32(msg.UserID)})
		if err != nil {
			continue
		}

		responses = append(responses, &models.MessageResponse{
			ID:        msg.ID,
			Text:      msg.Text,
			CreatedAt: msg.CreatedAt,
			Username:  userInfo.Username,
			Email:     userInfo.Email,
		})
	}

	return responses, nil
}

func (s *ChatService) DeleteMessage(ctx context.Context, messageID int, isAdmin bool) error {
	msg, err := s.repo.GetMessageByIDRepo(ctx, messageID)
	if err != nil {
		return err
	}

	if !isAdmin {
		return errors.New("permission denied: only admin can delete messages")
	}

	return s.repo.DeleteMessageRepo(ctx, msg)
}

func (s *ChatService) StartMessageCleaner(olderThan time.Duration) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		_ = s.repo.DeleteOldMessagesRepo(ctx, olderThan)
	}
}

func (s *ChatService) MessageChannel() <-chan *models.Message {
	return s.msgChan
}

func (s *ChatService) Close() {
	close(s.msgChan)
}
