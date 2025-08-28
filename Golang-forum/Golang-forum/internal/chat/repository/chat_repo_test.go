// package repository_test

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/playboi9/golang-forum/internal/chat/models"
// 	"github.com/playboi9/golang-forum/internal/chat/repository"
// 	"github.com/playboi9/golang-forum/pkg/db"
// )

// var chatRepo *repository.ChatRepository

// func TestMain(m *testing.M) {
// 	conn := db.NewTestDB()
// 	chatRepo = repository.NewChatRepository(conn)
// 	code := m.Run()
// 	cleanupTestTables(conn) // Очистка таблиц после тестов
// 	os.Exit(code)
// }

// func cleanupTestTables(conn *sql.DB) {
// 	conn.Exec("TRUNCATE TABLE messages RESTART IDENTITY CASCADE")
// 	conn.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
// 	conn.Exec("TRUNCATE TABLE banned_users RESTART IDENTITY CASCADE")
// }

// func createTestUser(conn *sql.DB, id int) error {
// 	_, err := conn.Exec(`INSERT INTO users (id, username, password) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING`, id, fmt.Sprintf("testuser%d", id), "testpass")
// 	return err
// }

// func TestSaveAndGetMessage(t *testing.T) {
// 	ctx := context.Background()
// 	if err := createTestUser(chatRepo.DB(), 1); err != nil {
// 		t.Fatalf("Failed to create test user: %v", err)
// 	}
// 	msg := &models.Message{
// 		UserID: 1, // убедись, что пользователь с ID 1 существует
// 		Text:   "hello test",
// 	}

// 	err := chatRepo.SaveMessageRepo(ctx, msg)
// 	if err != nil {
// 		t.Fatalf("SaveMessageRepo failed: %v", err)
// 	}

// 	got, err := chatRepo.GetMessageByIDRepo(ctx, msg.ID)
// 	if err != nil {
// 		t.Fatalf("GetMessageByIDRepo failed: %v", err)
// 	}

// 	if got.Text != msg.Text || got.UserID != msg.UserID {
// 		t.Errorf("expected %+v, got %+v", msg, got)
// 	}
// }

// func TestDeleteOldMessages(t *testing.T) {
// 	ctx := context.Background()
// 	if err := createTestUser(chatRepo.DB(), 1); err != nil {
// 		t.Fatalf("Failed to create test user: %v", err)
// 	}
// 	msg := &models.Message{
// 		UserID: 1,
// 		Text:   "old message",
// 	}
// 	err := chatRepo.SaveMessageRepo(ctx, msg)
// 	if err != nil {
// 		t.Fatalf("SaveMessageRepo failed: %v", err)
// 	}

// 	// вручную установим старую дату
// 	_, err = chatRepo.DB().Exec(`UPDATE messages SET created_at = NOW() - INTERVAL '2 days' WHERE id = $1`, msg.ID)
// 	if err != nil {
// 		t.Fatalf("Failed to update created_at: %v", err)
// 	}

// 	err = chatRepo.DeleteOldMessagesRepo(ctx, 24*time.Hour)
// 	if err != nil {
// 		t.Fatalf("DeleteOldMessagesRepo failed: %v", err)
// 	}

// 	_, err = chatRepo.GetMessageByIDRepo(ctx, msg.ID)
// 	if err == nil {
// 		t.Errorf("Expected message to be deleted, but it still exists")
// 	}
// }

// func TestGetMessagesRepo(t *testing.T) {
// 	ctx := context.Background()
// 	if err := createTestUser(chatRepo.DB(), 1); err != nil {
// 		t.Fatalf("Failed to create test user: %v", err)
// 	}

// 	for i := 0; i < 5; i++ {
// 		msg := &models.Message{
// 			UserID: 1,
// 			Text:   fmt.Sprintf("Message %d", i),
// 		}
// 		err := chatRepo.SaveMessageRepo(ctx, msg)
// 		if err != nil {
// 			t.Fatalf("Failed to insert message: %v", err)
// 		}
// 	}

// 	msgs, err := chatRepo.GetMessagesRepo(ctx, 3)
// 	if err != nil {
// 		t.Fatalf("GetMessagesRepo failed: %v", err)
// 	}
// 	if len(msgs) != 3 {
// 		t.Errorf("Expected 3 messages, got %d", len(msgs))
// 	}
// }

package repository_test