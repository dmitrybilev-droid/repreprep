package repository_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"golang-forum/internal/user/repository"
	"golang-forum/pkg/db"
	
	"github.com/rs/zerolog"
)

var (
	testAdminRepo *repository.AdminRepository
	testUserRepo  *repository.UserRepository
	conn          *sql.DB
)

func TestMain(m *testing.M) {
	conn = db.NewTestDB()
	logger := zerolog.New(os.Stdout)
	testUserRepo = repository.NewUserRepository(conn, logger)
	testAdminRepo = repository.NewAdminRepository(conn, logger)
	code := m.Run()
	cleanupTestTables(conn)
	os.Exit(code)
}

func cleanupTestTables(conn *sql.DB) {
	conn.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	conn.Exec("TRUNCATE TABLE banned_users RESTART IDENTITY CASCADE")
}

func TestCreateAndGetUser(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	username := "testuser"
	pass := "hashedpass"

	err := testUserRepo.CreateUserRepo(ctx, email, username, pass, false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, err := testUserRepo.GetUserByEmailRepo(ctx, email)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if user.Email != email || user.Username != username {
		t.Fatalf("Expected user with email %s and username %s, got %v", email, username, user)
	}
}

func TestBanAndUnbanUser(t *testing.T) {
	ctx := context.Background()
	email := "ban@example.com"
	username := "banuser"
	pass := "banhash"

	err := testUserRepo.CreateUserRepo(ctx, email, username, pass, false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	err = testAdminRepo.BanUserRepo(ctx, email)
	if err != nil {
		t.Fatalf("Failed to ban user: %v", err)
	}

	banned, err := testAdminRepo.GetBannedUserByEmailRepo(ctx, email)
	if err != nil {
		t.Fatalf("Failed to get banned user: %v", err)
	}
	if banned.Email != email {
		t.Errorf("Expected banned email %s, got %s", email, banned.Email)
	}

	err = testAdminRepo.UnBanUserRepo(ctx, email)
	if err != nil {
		t.Fatalf("Failed to unban user: %v", err)
	}

	_, err = testAdminRepo.GetBannedUserByEmailRepo(ctx, email)
	if err == nil {
		t.Error("Expected error when getting unbanned user, got nil")
	}
}
