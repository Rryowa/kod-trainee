package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"kod/internal/models"
	"kod/internal/storage"
	"net/http"
	"strings"
)

type UserService struct {
	storage        storage.Storage
	sessionService *SessionService
}

func NewUserService(s storage.Storage, ss *SessionService) *UserService {
	return &UserService{storage: s, sessionService: ss}
}

// SignUp Hashes password and adds user to db
func (us *UserService) SignUp(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := us.storage.GetUser(ctx, user.Username)
	if err == nil {
		return nil, errors.New(fmt.Sprintf("user already exists: %s", user.Username))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Username = strings.ToLower(user.Username)
	user.Password = string(hash)

	newUser, err := us.storage.AddUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// LogIn Validates user's password, creates a jwt token and cookie
func (us *UserService) LogIn(r *http.Request, userRequest *models.User) (*http.Cookie, error) {
	user, err := us.storage.GetUser(r.Context(), userRequest.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := us.sessionService.CreateToken(&user)
	if err != nil {
		return nil, err
	}

	return us.sessionService.CreateCookie(token)
}

func (us *UserService) LogOut() (*http.Cookie, error) {
	return us.sessionService.DeleteCookie()
}