package models

import (
	"context"

	"github.com/google/uuid"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenService interface {
	CreateAccessToken(user *User, expired int) (accessToken string, err error)
	CreateRefreshToken(user *User, expired int) (refreshToken string, err error)
	ExtractIPFromRefreshToken(token string, originalToken string) (string, error) //Payload токенов должен содержать сведения об ip адресе клиента, которому он был выдан. В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера.
	ExtractIDFromToken(requestedToken string) (uuid.UUID, error)
	ExtractEmailFromRefreshToken(originalToken string) (string, error)
	IsAuthorized(token string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (user *User, err error)
}
