package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"kod/internal/storage"
	"time"
)

type AuthService struct {
	storage storage.Storage
	log     *zap.SugaredLogger
}

func NewAuthService(s storage.Storage) *NoteService {
	return &NoteService{storage: s}
}

// Middleware puts userId into context!
func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func createToken(userID, username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}