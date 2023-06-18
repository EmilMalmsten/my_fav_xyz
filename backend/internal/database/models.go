package database

import (
	"database/sql"
	"time"
)

type DbConfig struct {
	database *sql.DB
}

type Toplist struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	UserID      int           `json:"user_id"`
	CreatedAt   time.Time     `json:"created_at"`
	Items       []ToplistItem `json:"items"`
}

type ToplistItem struct {
	ID          int    `json:"id"`
	ListID      int    `json:"listID"`
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
}

type Revocation struct {
	Token     string    `json:"token"`
	RevokedAt time.Time `json:"revoked_at"`
}
