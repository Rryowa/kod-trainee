package middleware

import (
	"errors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
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

// AuthMiddleware extracts user from cookie, validates and pass it to context
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, _ := opentracing.StartSpanFromContext(r.Context(), "middleware.Auth")
		defer span.Finish()

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

// RateLimit allows 2 requests per second and 4 requests in a single ‘burst’
func (m *Middleware) RateLimit(next http.Handler) http.Handler {

	limiter := rate.NewLimiter(1, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			err := errors.New("rate limit exceeded")
			m.zapLogger.Error(err)
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}