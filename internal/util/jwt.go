package util

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreateJWT(userId int, userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userId),
		"userName":  userName,
		"expiresAt": time.Now().Add(2 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(t string) (*jwt.Token, error) {
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
		case !token.Valid:
			return nil, errors.New("token is expired")
		default:
			return nil, fmt.Errorf("couldn't handle this token: %v", err)
		}
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	fmt.Println(claims["expiresAt"])

	expiresAt := claims["expiresAt"].(int64)
	if expiresAt < time.Now().Unix() {
		return nil, errors.New("token is expired")
	}

	return token, nil
}

func GetAuthHeader(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		return "", errors.New("empty token from request")
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	return tokenString, nil
}