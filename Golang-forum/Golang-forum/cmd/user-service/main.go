package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chatgen "golang-forum/internal/chat/transport/grpc/gen/api"
	"golang-forum/internal/user/repository"
	"golang-forum/internal/user/service"
	userGrpc "golang-forum/internal/user/transport/grpc"
	usergen "golang-forum/internal/user/transport/grpc/gen/api"
	userHttp "golang-forum/internal/user/transport/http"
	"golang-forum/pkg/config"
	"golang-forum/pkg/db"
	"golang-forum/pkg/logger"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func main() {
	log := logger.New()
	cfg := config.LoadUserConfig()

	// Подключение к БД
	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer dbConn.Close()

	// Инициализация зависимостей
	userRepo := repository.NewUserRepository(dbConn, log)
	adminRepo := repository.NewAdminRepository(dbConn, log)
	authService := service.NewAuthService(userRepo)
	adminService := service.NewAdminService(adminRepo, nil)

	// gRPC-сервера
	authServer := userGrpc.NewAuthServer(authService, cfg.JWTSecret)
	adminServer := userGrpc.NewAdminServer(adminService, cfg.JWTSecret, nil)

	// Запуск gRPC
	go startGRPCServer(authServer, adminServer, cfg.AuthGRPCAddr, log)

	time.Sleep(2 * time.Second)

	// Подключение к chat-service по gRPC
	conn, err := grpc.Dial(cfg.ChatGRPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to chat-service")
	}
	defer conn.Close()
	chatClient := chatgen.NewChatServiceClient(conn)

	adminService.UpdateChatClient(chatClient)
	adminServer.UpdateChatClient(chatClient)

	// HTTP-маршруты
	mux := http.NewServeMux()
	userHttp.RegisterAuthRoutes(mux, authServer)
	userHttp.RegisterAdminRoutes(mux, adminServer)

	// Запуск HTTP
	go func() {
		log.Info().Msg("HTTP server started on :8080")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	waitForShutdown(log)
}

func startGRPCServer(authServer *userGrpc.AuthServer, adminServer *userGrpc.AdminServer, addr string, log zerolog.Logger) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	s := grpc.NewServer()
	usergen.RegisterAuthServiceServer(s, authServer)
	usergen.RegisterAdminServiceServer(s, adminServer)

	log.Info().Msgf("gRPC server started on %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("gRPC server failed")
	}
}

func waitForShutdown(log zerolog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")
}
