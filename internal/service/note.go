package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kod/internal/models"
	"kod/internal/storage"
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

	if err := ns.checkGrammar(note.Title); err != nil {
		return models.Note{}, fmt.Errorf("title grammar check failed: %v", err)
	}
	if err := ns.checkGrammar(note.Text); err != nil {
		return models.Note{}, fmt.Errorf("text grammar check failed: %v", err)
	}

	user, err := GetUserFromContext(r)
	if err != nil {
		return models.Note{}, err
	}

	note.UserId = user.Id
	note.UserName = user.Username
	note.CreatedAt = time.Now()

	return ns.storage.AddNote(r.Context(), note)
}

func (ns *NoteService) GetNotes(r *http.Request) ([]models.Note, error) {
	user, err := GetUserFromContext(r)
	if err != nil {
		return nil, err
	}

	return ns.storage.GetNotes(r.Context(), user.Id, 0, 10)
}

func (ns *NoteService) checkGrammar(text string) error {
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

	// Read and parse the response
	body, err := io.ReadAll(resp.Body) // Use io.ReadAll instead of ioutil.ReadAll
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