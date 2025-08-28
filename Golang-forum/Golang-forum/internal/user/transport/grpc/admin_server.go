package grpc

import (
	"context"

	chatgen "golang-forum/internal/chat/transport/grpc/gen/api"
	"golang-forum/internal/user/service"
	gen "golang-forum/internal/user/transport/grpc/gen/api"
	"golang-forum/pkg/jwt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminServer struct {
	gen.UnimplementedAdminServiceServer
	adminService *service.AdminService
	jwtSecret    string
	chatClient   chatgen.ChatServiceClient
}

func NewAdminServer(adminService *service.AdminService, jwtSecret string, chatClient chatgen.ChatServiceClient) *AdminServer {
	return &AdminServer{
		adminService: adminService,
		jwtSecret:    jwtSecret,
		chatClient:   chatClient,
	}
}
func (s *AdminServer) UpdateChatClient(client chatgen.ChatServiceClient) {
	s.chatClient = client
}

func (s *AdminServer) GetUserList(ctx context.Context, req *gen.GetUserListRequest) (*gen.GetUserListResponse, error) {
	claims, err := jwt.ValidateJWTFromContext(ctx, s.jwtSecret)
	if err != nil || !claims.IsAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin rights required")
	}
	users, err := s.adminService.GetAllUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var userList []*gen.UserInfo
	for _, user := range users {
		userList = append(userList, &gen.UserInfo{
			Id:       int32(user.ID),
			Email:    user.Email,
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		})
	}
	return &gen.GetUserListResponse{
		Users: userList,
		Total: int32(len(users)),
	}, nil
}

func (s *AdminServer) IsBanned(ctx context.Context, req *gen.IsBannedRequest) (*gen.IsBannedResponse, error) {
	isBanned, err := s.adminService.IsUserBanned(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.IsBannedResponse{Banned: isBanned}, nil
}

func (s *AdminServer) BanUser(ctx context.Context, req *gen.BanRequest) (*gen.BanResponse, error) {
	claims, err := jwt.ValidateJWTFromContext(ctx, s.jwtSecret)
	if err != nil || !claims.IsAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin rights required")
	}

	if err := s.adminService.BanUser(ctx, req.Email); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.BanResponse{Success: true}, nil
}

func (s *AdminServer) UnBanUser(ctx context.Context, req *gen.UnBanRequest) (*gen.UnBanResponse, error) {
	claims, err := jwt.ValidateJWTFromContext(ctx, s.jwtSecret)
	if err != nil || !claims.IsAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin rights required")
	}
	if err := s.adminService.UnBanUser(ctx, req.Email); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.UnBanResponse{Success: true}, nil
}

func (s *AdminServer) DeleteMessage(ctx context.Context, req *gen.DeleteMessageRequest) (*gen.DeleteMessageResponse, error) {
	claims, err := jwt.ValidateJWTFromContext(ctx, s.jwtSecret)
	if err != nil || !claims.IsAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin rights required")
	}

	err = s.adminService.DeleteMessage(ctx, int(req.MessageId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.DeleteMessageResponse{Success: true}, nil
}
