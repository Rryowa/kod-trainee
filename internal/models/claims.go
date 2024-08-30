package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserId    int       `json:"user_id"`
	UserName  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
	jwt.RegisteredClaims
}