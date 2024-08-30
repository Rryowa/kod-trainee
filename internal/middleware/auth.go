package middleware

import (
	"go.uber.org/zap"
	"kod/internal/models"
	"kod/internal/service"
	"net/http"
)

type Middleware struct {
	sessionService *service.SessionService
	zapLogger      *zap.SugaredLogger
}

func NewMiddleware(ss *service.SessionService, l *zap.SugaredLogger) *Middleware {
	return &Middleware{sessionService: ss, zapLogger: l}
}

// AuthMiddleware extracts user from cookie and pass it to context
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := m.sessionService.GetCookieValue(r)
		if err != nil {
			m.zapLogger.Errorf("unauthorized: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := m.sessionService.ValidateToken(tokenString)
		if err != nil {
			m.zapLogger.Error(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userCtx := &models.User{
			Id:       claims.UserId,
			Username: claims.UserName,
		}

		r = service.SetUserContext(r, userCtx)

		next.ServeHTTP(w, r)
	})
}