package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"kod/internal/util"
	"net/http"
)

type Middleware struct {
	zapLogger *zap.SugaredLogger
}

func NewMiddleware(l *zap.SugaredLogger) *Middleware {
	return &Middleware{zapLogger: l}
}

// AuthMiddleware extracts user_id from request and pass it to context
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := util.GetJwtToken(r)
		if err != nil {
			m.zapLogger.Errorf("unauthorized: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		token, err := util.ValidateJwt(tokenString)
		if err != nil {
			m.zapLogger.Error(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.zapLogger.Error("failed to get claims")
			util.HandleHttpError(w, "failed to get claims", http.StatusUnauthorized)
			return
		}
		userId, ok := claims["userId"].(int)
		userName, ok := claims["userName"].(string)

		ctx := context.WithValue(r.Context(), "userId", userId)
		ctx = context.WithValue(ctx, "userName", userName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}