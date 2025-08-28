package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
}

type BannedUser struct {
	ID       int       `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	BannedAt time.Time `json:"banned_at"`
}
