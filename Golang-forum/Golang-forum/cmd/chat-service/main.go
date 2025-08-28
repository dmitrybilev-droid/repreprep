package main

import (
	"net/http"

	"golang-forum/internal/chat/repository"
	"golang-forum/internal/chat/service"
	"golang-forum/internal/chat/transport/websocket"
	usergen "golang-forum/internal/user/transport/grpc/gen/api"
	"golang-forum/pkg/config"
	"golang-forum/pkg/db"
	"github.com/playboi9/golang-forum-pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadChatConfig()
	log := logger.New()
	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to DB")
	}

	// gRPC клиент к user-service
	conn, err := grpc.Dial(cfg.AuthGRPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to user-service gRPC")
	}
	defer conn.Close()

	authClient := usergen.NewAuthServiceClient(conn)
	adminClient := usergen.NewAdminServiceClient(conn)

	// Инициализация сервисов
	chatRepo := repository.NewChatRepository(dbConn)
	chatSvc := service.NewChatService(chatRepo, authClient, adminClient)

	wsHandler := websocket.NewWSHandler(chatSvc, cfg.JWTSecret)

	// Запуск рассылки сообщений
	go wsHandler.BroadcastMessages()

	// HTTP-роутинг
	http.HandleFunc("/ws", wsHandler.HandleConnections)
	log.Info().Msgf("Chat service started on :8081")
	log.Fatal().Err(http.ListenAndServe(":8081", nil))
}
