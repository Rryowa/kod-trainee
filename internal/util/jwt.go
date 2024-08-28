package util

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

func CreateJWT(userId int, userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"userName":  userName,
		"expiresAt": time.Now().Add(2 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJwt(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.New("that's not even a token")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, errors.New("invalid signature")
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, errors.New("token is either expired or not active yet")
		default:
			return nil, fmt.Errorf("couldn't handle this token: %v", err)
		}
	}

	return token, nil
}

func GetJwtToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		return "", errors.New("empty token from request")
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	return tokenString, nil
}