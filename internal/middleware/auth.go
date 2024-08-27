package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, "Forbidden", http.StatusUnauthorized)
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte("secret"), nil
		})

		if err != nil {
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}