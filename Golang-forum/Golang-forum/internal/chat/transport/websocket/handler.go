package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"

	"golang-forum/internal/chat/models"
	"golang-forum/internal/chat/service"
	"golang-forum/pkg/jwt"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	chatService *service.ChatService
	clients     map[*websocket.Conn]bool
	mutex       sync.Mutex
	jwtSecret   string
}

func NewWSHandler(chatService *service.ChatService, jwtSecret string) *WSHandler {
	return &WSHandler{
		chatService: chatService,
		clients:     make(map[*websocket.Conn]bool),
		jwtSecret:   jwtSecret,
	}
}

// Обработчик подключения WebSocket
func (h *WSHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	var userID int
	var isAuthorized bool
	token := r.URL.Query().Get("token")

	if token != "" {
		claims, err := jwt.ValidateJWT(token, h.jwtSecret)
		if err == nil && claims.UserID > 0 {
			userID = claims.UserID
			isAuthorized = true
			log.Printf("Authorized user: %d", userID)
		} else {
			log.Printf("Invalid token: %v", err)
		}
	} else {
		log.Println("Guest connected (no token)")
	}

	// Подключение WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Добавляем клиента
	h.mutex.Lock()
	h.clients[ws] = true
	h.mutex.Unlock()

	defer func() {
		h.mutex.Lock()
		delete(h.clients, ws)
		h.mutex.Unlock()
		ws.Close()
	}()

	// Отправляем историю сообщений
	ctx := context.Background()
	messages, err := h.chatService.GetMessages(ctx, 50)
	if err == nil {
		if err := ws.WriteJSON(messages); err != nil {
			log.Printf("Failed to write history to client: %v", err)
		}
	}

	// Читаем входящие сообщения
	for {
		var msg struct {
			Text string `json:"text"`
		}
		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("Read error from client: %v", err)
			h.mutex.Lock()
			delete(h.clients, ws)
			h.mutex.Unlock()
			break
		}

		// Сохраняем сообщение в БД
		if !isAuthorized {
			ws.WriteMessage(websocket.TextMessage, []byte("Вы не авторизованы для отправки сообщений."))
			continue
		}
		message, err := h.chatService.SendMessage(ctx, userID, msg.Text)

		if err != nil {
			log.Printf("Failed to send message: %v", err)
			ws.WriteMessage(websocket.TextMessage, []byte("Ошибка отправки сообщения: "+err.Error()))
			continue
		}
		log.Printf("Message sent: %s", message.Text)

	}
}

// Рассылка новых сообщений всем клиентам
func (h *WSHandler) BroadcastMessages() {
	for msg := range h.chatService.MessageChannel() {
		userInfo, err := h.chatService.GetUserInfoByID(context.Background(), msg.UserID)
		if err != nil {
			log.Printf("Failed to get user info: %v", err)
			continue
		}

		response := models.MessageResponse{
			ID:        msg.ID,
			Text:      msg.Text,
			CreatedAt: msg.CreatedAt,
			Username:  userInfo.Username,
			Email:     userInfo.Email,
		}

		h.mutex.Lock()
		for client := range h.clients {
			if err := client.WriteJSON(response); err != nil {
				log.Printf("Failed to write to client: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mutex.Unlock()
	}
}

// Запуск WebSocket-сервера
func (h *WSHandler) Start(addr string) {
	http.HandleFunc("/ws", h.HandleConnections)
	go h.BroadcastMessages()
	log.Printf("WebSocket server started on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start WebSocket server: %v", err)
	}
}
