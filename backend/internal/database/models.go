package database

import (
	"database/sql"
	"time"
)

type DbConfig struct {
	database *sql.DB
}

type Toplist struct {
	ToplistID   int           `json:"toplist_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	UserID      int           `json:"user_id"`
	Username    string        `json:"username"`
	CreatedAt   time.Time     `json:"created_at"`
	Items       []ToplistItem `json:"items"`
	Views       int           `json:"views"`
	LikeCount   int           `json:"like_count"`
	LikeIDs     []int         `json:"like_ids"`
}

type ToplistItem struct {
	ItemID      int    `json:"item_id"`
	ListID      int    `json:"list_id"`
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       []byte `json:"-"`
	ImagePath   string `json:"image_path"`
}

type User struct {
	ID                         int       `json:"id"`
	Email                      string    `json:"email"`
	HashedPassword             string    `json:"hashed_password"`
	CreatedAt                  time.Time `json:"created_at"`
	PasswordResetToken         string    `json:"password_reset_token"`
	PasswordResetTokenExpireAt time.Time `json:"password_reset_token_expire_at"`
	Username                   string    `json:"username"`
}

type Revocation struct {
	Token     string    `json:"token"`
	RevokedAt time.Time `json:"revoked_at"`
}
