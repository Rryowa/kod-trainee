package service

import (
	"context"
	"kod/internal/models"
	"kod/internal/storage"
)

type NoteService struct {
	storage storage.Storage
}

func NewNoteService(s storage.Storage) *NoteService {
	return &NoteService{storage: s}
}

func (ns *NoteService) AddNote(ctx context.Context, note *models.Note) (models.Note, error) {
	return ns.storage.AddNote(ctx, *note)
}

func (ns *NoteService) GetNotes(ctx context.Context, userId int) ([]models.Note, error) {
	//TODO: replace 0 and 10 with query
	return ns.storage.GetNotes(ctx, userId, 0, 10)
}