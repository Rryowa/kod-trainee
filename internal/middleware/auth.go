package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"kod/internal/models"
	"kod/internal/util"
	"net/http"
	"strconv"
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
		w.Header().Add("Vary", "Authorization")

		tokenString, err := util.GetAuthHeader(r)
		if err != nil {
			m.zapLogger.Errorf("unauthorized: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := util.ValidateToken(tokenString)
		if err != nil {
			m.zapLogger.Error(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.zapLogger.Error("failed to get claims")
			http.Error(w, "failed to get claims", http.StatusUnauthorized)
			return
		}
		userIdStr, ok := claims["userId"].(string)
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			m.zapLogger.Error(err)
			http.Error(w, "user id is not integer", http.StatusUnauthorized)
		}
		userCtx := &models.User{
			Id:       userId,
			Username: claims["userName"].(string),
		}

		r = util.SetUserContext(r, userCtx)

		next.ServeHTTP(w, r)
	})
}