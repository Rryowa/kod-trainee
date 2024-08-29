package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"kod/internal/models"
	"kod/internal/storage"
	"kod/internal/util"
	"net/http"
	"strings"
)

type UserService struct {
	storage storage.Storage
}

func NewUserService(s storage.Storage) *UserService {
	return &UserService{storage: s}
}

func (us *UserService) SignUp(ctx context.Context, user *models.User) (*models.User, error) {
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

func (us *UserService) LogIn(r *http.Request, userRequest *models.User) (string, error) {
	user, err := us.storage.GetUser(r.Context(), userRequest.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password)); err != nil {
		return "", errors.New("invalid password")
	}

	token, err := util.CreateJWT(user.Id, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}