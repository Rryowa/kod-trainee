package service

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"kod/internal/models"
	"kod/internal/storage"
	"kod/internal/util"
	"net/http"
	"strconv"
)

type NoteService struct {
	storage storage.Storage
	log     *zap.SugaredLogger
}

func NewNoteService(s storage.Storage) *NoteService {
	return &NoteService{storage: s}
}

func (ns *NoteService) HandleAddNote(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var note *models.Note
	err = json.Unmarshal(body, &note)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, util.ErrInvalidRequestPayload)
		return
	}

	if len(note.Text) == 0 {
		util.WriteJSON(w, http.StatusBadRequest, util.ErrEmptyNote)
		return
	}

	newNote, err := ns.storage.AddNote(r.Context(), *note)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusCreated, newNote)
}

func (ns *NoteService) HandleGetNotes(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("user_id")
	if len(userIdStr) == 0 {
		util.WriteJSON(w, http.StatusBadRequest, util.ErrUserIdNotProvided)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId <= 0 {
		util.WriteJSON(w, http.StatusBadRequest, util.ErrInvalidUserId)
	}

	//TODO: replace 0 and 10 with query
	notes, err := ns.storage.GetNotes(r.Context(), userIdStr, 0, 10)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, notes)
}

func (ns *NoteService) HandleWelcome(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, "Hello!")
}