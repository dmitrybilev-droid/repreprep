package grpc

import (
	"context"

	"golang-forum/internal/chat/service"
	gen "golang-forum/internal/chat/transport/grpc/gen/api"
)

type ChatGRPCServer struct {
	gen.UnimplementedChatServiceServer
	chatService *service.ChatService
}

func NewChatGRPCServer(chatService *service.ChatService) *ChatGRPCServer {
	return &ChatGRPCServer{
		chatService: chatService,
	}
}

func (s *ChatGRPCServer) DeleteMessage(ctx context.Context, req *gen.DeleteMessageRequest) (*gen.DeleteMessageResponse, error) {
	err := s.chatService.DeleteMessage(ctx, int(req.MessageId), req.IsAdmin)
	if err != nil {
		return &gen.DeleteMessageResponse{Success: false}, err
	}
	return &gen.DeleteMessageResponse{Success: true}, nil
}
