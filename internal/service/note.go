package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"io"
	"kod/internal/models"
	"kod/internal/storage"
	"net/http"
	"strconv"
	"time"
)

const defaultLimit = 10

type NoteService struct {
	storage storage.Storage
}

func NewNoteService(s storage.Storage) *NoteService {
	return &NoteService{storage: s}
}

func (ns *NoteService) AddNote(r *http.Request, note *models.Note) (models.Note, error) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "service.AddNote")
	defer span.Finish()

	if err := ns.checkGrammar(ctx, note.Title); err != nil {
		return models.Note{}, fmt.Errorf("title grammar check failed: %v", err)
	}
	if err := ns.checkGrammar(ctx, note.Text); err != nil {
		return models.Note{}, fmt.Errorf("text grammar check failed: %v", err)
	}

	user, err := GetUserFromContext(ctx)
	if err != nil {
		return models.Note{}, err
	}

	note.UserId = user.Id
	note.UserName = user.Username
	note.CreatedAt = time.Now()

	return ns.storage.AddNote(ctx, note)
}

func (ns *NoteService) GetNotes(r *http.Request) ([]models.Note, error) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "service.GetNotes")
	defer span.Finish()

	pageParam := r.URL.Query().Get("p")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}
	offset := (page - 1) * defaultLimit
	limit := defaultLimit

	user, err := GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return ns.storage.GetNotes(ctx, user.Id, offset, limit)
}

func (ns *NoteService) checkGrammar(ctx context.Context, text string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.checkGrammar")
	defer span.Finish()

	if len(text) > 10000 {
		return fmt.Errorf("text too long")
	}
	apiURL := "https://speller.yandex.net/services/spellservice.json/checkText"
	requestBody := fmt.Sprintf("text=%s?options=%d?format=%s", text, 4, "plain")
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create grammar check request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform grammar check: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read grammar check response: %w", err)
	}

	var grammarErrors []map[string]interface{}
	err = json.Unmarshal(body, &grammarErrors)
	if err != nil {
		return fmt.Errorf("failed to parse grammar check response: %w", err)
	}

	if len(grammarErrors) > 0 {
		return fmt.Errorf("%d grammar errors found in the note", len(grammarErrors))
	}

	return nil
}