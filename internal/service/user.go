package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"kod/internal/models"
	"kod/internal/storage"
	"kod/internal/util"
	"strings"
)

type UserService struct {
	storage storage.Storage
	//httpConfig *config.HttpConfig
}

func NewUserService(s storage.Storage) *UserService {
	return &UserService{storage: s}
}

func (us *UserService) SignUp(ctx context.Context, user *models.User) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Username = strings.ToLower(user.Username)
	user.Password = string(hash)

	newUser, err := us.storage.AddUser(ctx, user)
	if err != nil {
		return "", err
	}

	return util.CreateJWT(newUser.Id, newUser.Username)
}

func (us *UserService) LogIn(ctx context.Context, userRequest *models.User) (string, error) {
	user, err := us.storage.GetUser(ctx, userRequest.Id)
	if err != nil {
		//TODO: smth with norows
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