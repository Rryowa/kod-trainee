package util

import (
	"errors"
)

var (
	ErrNoteNotFound          = errors.New("note not found")
	ErrInvalidRequestPayload = errors.New("invalid request payload")
	ErrEmptyNote             = errors.New("note cant be empty")
	ErrInvalidUserId         = errors.New("invalid user id")
	ErrUserIdNotProvided     = errors.New("user id not provided")
)