package util

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidRequestPayload = errors.New("invalid request payload")
	ErrEmptyNote             = errors.New("note cant be empty")
	ErrEmptyUsername         = errors.New("empty username")
	ErrEmptyPassword         = errors.New("empty password")
)

func HandleHttpError(w http.ResponseWriter, msg string, status int) {
	http.Error(w, msg, status)
}