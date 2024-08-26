package models

import (
	"context"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginService interface {
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (user *User, err error)
	CreateAccessToken(user *User, expired int) (accessToken string, err error)
	CreateRefreshToken(user *User, expired int) (refreshToken string, err error)
	RevokeTokens(ctx context.Context, user *User) error // на всякий
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
}
