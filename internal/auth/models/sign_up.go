package models

import "context"

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignUpService interface {
	RegisterUser(ctx context.Context, user *User) (err error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CreateAccessToken(user *User, expired int) (accessToken string, err error)
	CreateRefreshToken(user *User, expired int) (refreshToken string, err error)
}