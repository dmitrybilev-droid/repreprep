package grpc

import (
	"context"

	"golang-forum/internal/user/service"
	gen "golang-forum/internal/user/transport/grpc/gen/api"
	"golang-forum/pkg/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	gen.UnimplementedAuthServiceServer
	authService *service.AuthService
	jwtSecret   string
}

func NewAuthServer(authService *service.AuthService, jwtSecret string) *AuthServer {
	return &AuthServer{
		authService: authService,
		jwtSecret:   jwtSecret,
	}
}
func (s *AuthServer) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	if req.Email == "" || req.Password == "" || req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "email, username, and password are required")
	}

	err := s.authService.Register(ctx, req.Email, req.Username, req.Password, req.IsAdmin)
	if err != nil {
		return &gen.RegisterResponse{
			Error: err.Error(),
		}, nil
	}

	loginResp, err := s.Login(ctx, &gen.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return &gen.RegisterResponse{
			Error: "auto-login failed after registration",
		}, nil
	}

	return &gen.RegisterResponse{
		Email:    loginResp.Email,
		Username: loginResp.Username,
		Token:    loginResp.Token,
		IsAdmin:  loginResp.IsAdmin,
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	user, err := s.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &gen.LoginResponse{
			Error: "invalid credentials",
		}, status.Error(codes.Unauthenticated, "authentication failed")
	}

	token, err := jwt.GenerateJWT(user.ID, user.IsAdmin, s.jwtSecret)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &gen.LoginResponse{
		Token:    token,
		Email:    user.Email,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *gen.TokenRequest) (*gen.TokenResponse, error) {
	claims, err := jwt.ValidateJWT(req.Token, s.jwtSecret)
	if err != nil {
		return &gen.TokenResponse{Valid: false}, nil
	}

	return &gen.TokenResponse{
		Valid:   true,
		UserId:  int32(claims.UserID),
		IsAdmin: claims.IsAdmin,
	}, nil
}

func (s *AuthServer) GetUserInfoByID(ctx context.Context, req *gen.GetUserInfoByIDRequest) (*gen.GetUserInfoByIDResponse, error) {
	user, err := s.authService.GetUserByID(ctx, int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &gen.GetUserInfoByIDResponse{
		Email:    user.Email,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}, nil
}
