package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RefreshToken string    `json:"refresh_token"`
	Salt         string    `json:"salt"`
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CheckRefreshTokenExists(ctx context.Context, hashedToken string) (bool, error)
	SaveRefreshToken(id uuid.UUID, hashedToken string) error
	GetUserByUsername(ctx context.Context, username string) (user *User, err error)
	GetRefreshToken(ctx context.Context, email string) (hashedToken string, err error)
	DeleteUserRefreshTokenByEmail(ctx context.Context, email string) (err error)
	GetUserByEmail(ctx context.Context, email string) (user *User,  err error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
}
