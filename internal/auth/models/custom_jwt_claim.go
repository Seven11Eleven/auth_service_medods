package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTCustomClaims struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	ID        uuid.UUID `json:"id"`
	IPAddress string    `json:"ip_address"`
	jwt.RegisteredClaims
}

type CustomRefreshClaims struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
