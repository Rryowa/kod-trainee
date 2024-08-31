package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"kod/internal/models"
	"kod/internal/models/config"
	"net/http"
	"time"
)

const nameOfUserStruct = "user"

type SessionService struct {
	cfg *config.SessionConfig
}

func NewSessionService(c *config.SessionConfig) *SessionService {
	return &SessionService{cfg: c}
}

func (s *SessionService) CreateToken(user *models.User) (string, error) {
	claims := models.Claims{
		UserId:   user.Id,
		UserName: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JwtTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(s.cfg.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *SessionService) ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JwtSecret), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.New("that's not even a token")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, errors.New("invalid signature")
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errors.New("token is expired")
		case !token.Valid:
			return nil, errors.New("invalid token")
		default:
			return nil, fmt.Errorf("couldn't handle this token: %v", err)
		}
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	return claims, nil
}

func (s *SessionService) CreateCookie(token string) (*http.Cookie, error) {
	cookie := &http.Cookie{
		Name:     s.cfg.CookieName,
		Value:    token,
		Expires:  time.Now().Add(s.cfg.CookieTTL),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	// Encode the cookie value using base64.
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	// Check the total length of the cookie contents. Return the ErrValueTooLong
	if len(cookie.String()) > 4096 {
		return nil, errors.New("cookie value too long")
	}

	return cookie, nil
}

func (s *SessionService) GetCookieValue(r *http.Request) (string, error) {
	cookie, err := r.Cookie(s.cfg.CookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", errors.New("cookie expired")
		}
		return "", errors.New("server error")
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", errors.New("invalid cookie value")
	}

	return string(value), nil
}

func (s *SessionService) DeleteCookie() (*http.Cookie, error) {
	emptyCookie := &http.Cookie{
		Name:     s.cfg.CookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	return emptyCookie, nil
}

func GetUserFromContext(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(nameOfUserStruct).(*models.User)
	if !ok {
		return nil, errors.New("there is no user in context")
	}
	return user, nil
}

func SetUserContext(r *http.Request, userCtx *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), nameOfUserStruct, userCtx)
	return r.WithContext(ctx)
}