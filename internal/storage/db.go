package storage

import (
	"context"
	"kod/internal/models"
)

type Storage interface {
	// AddNote adds a note to db
	AddNote(ctx context.Context, note models.Note) (models.Note, error)
	// GetNotes returns list of notes, or ErrDoesNotExist
	GetNotes(ctx context.Context, userId int, offset, limit int) ([]models.Note, error)
	UserStorage
}

type UserStorage interface {
	AddUser(ctx context.Context, user *models.User) (models.User, error)
	GetUser(ctx context.Context, userName string) (models.User, error)
}