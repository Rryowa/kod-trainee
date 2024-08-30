package handler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"kod/internal/models"
	"kod/internal/service"
	"kod/internal/util"
	"net/http"
)

type Handler struct {
	noteService *service.NoteService
	userService *service.UserService
	zapLogger   *zap.SugaredLogger
}

func NewHandler(ns *service.NoteService, us *service.UserService, l *zap.SugaredLogger) *Handler {
	return &Handler{
		noteService: ns,
		userService: us,
		zapLogger:   l,
	}
}

func (c *Handler) HandleAddNote(w http.ResponseWriter, r *http.Request) {
	//Middleware already verified a token
	var note models.Note
	if err := util.DecodeJSONBody(r, &note); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	newNote, err := c.noteService.AddNote(r, &note)
	if err != nil {
		c.zapLogger.Error(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	util.WriteJSON(w, newNote)
}

func (c *Handler) HandleGetNotes(w http.ResponseWriter, r *http.Request) {
	//Middleware already verified a token
	//TODO: replace 0 and 10 with query
	notes, err := c.noteService.GetNotes(r)
	if err != nil {
		c.zapLogger.Errorf("Error getting notes: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, notes)
}

func (c *Handler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := util.DecodeJSONBody(r, &user); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newUser, err := c.userService.SignUp(r.Context(), &user)
	if err != nil {
		c.zapLogger.Errorf("Error SingUp: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, newUser)
}

func (c *Handler) HandleLogIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := util.DecodeJSONBody(r, &user); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cookie, err := c.userService.LogIn(r, &user)
	if err != nil {
		c.zapLogger.Errorf("Error LogIn: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
}

func (c *Handler) HandleLogOut(w http.ResponseWriter, r *http.Request) {
	emptyCookie, err := c.userService.LogOut()
	if err != nil {
		c.zapLogger.Errorf("Error LogOut: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r = r.WithContext(context.Background())

	http.SetCookie(w, emptyCookie)
}