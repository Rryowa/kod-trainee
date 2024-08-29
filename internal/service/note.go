package service

import (
	"context"
	"kod/internal/models"
	"kod/internal/storage"
	"kod/internal/util"
	"net/http"
	"time"
)

type NoteService struct {
	storage storage.Storage
}

func NewNoteService(s storage.Storage) *NoteService {
	return &NoteService{storage: s}
}

func (ns *NoteService) AddNote(r *http.Request, note *models.Note) (models.Note, error) {
	userCtx, err := util.GetUserFromContext(r)
	if err != nil {
		return models.Note{}, err
	}
	note.UserId = userCtx.Id
	note.CreatedAt = time.Now()

	return ns.storage.AddNote(r.Context(), *note)
}

func (ns *NoteService) GetNotes(ctx context.Context, userId int) ([]models.Note, error) {
	//TODO: replace 0 and 10 with query
	return ns.storage.GetNotes(ctx, userId, 0, 10)
}