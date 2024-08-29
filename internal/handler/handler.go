package handler

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"kod/internal/models"
	"kod/internal/service"
	"kod/internal/util"
	"net/http"
)

type Handler struct {
	noteService service.NoteService
	userService service.UserService
	zapLogger   *zap.SugaredLogger
}

func NewHandler(ns *service.NoteService, us *service.UserService, l *zap.SugaredLogger) *Handler {
	return &Handler{
		noteService: *ns,
		userService: *us,
		zapLogger:   l,
	}
}

func (c *Handler) HandleAddNote(w http.ResponseWriter, r *http.Request) {
	//Middleware already verified a token
	var note models.Note
	if err := util.DecodeJSONBody(w, r, &note); err != nil {
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

	c.zapLogger.Infof("HandleAddNote Ok")
	fmt.Fprintf(w, "Note: %+v\n", newNote)
	util.WriteJSON(w, newNote)
}

func (c *Handler) HandleGetNotes(w http.ResponseWriter, r *http.Request) {
	//Middleware already verified a token
	userCtx, err := util.GetUserFromContext(r)
	if err != nil {
		c.zapLogger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//userId, err := strconv.Atoi(userIdCtx)
	//if err != nil || userId <= 0 {
	//	ns.zapLogger.Errorf("Error converting user id to int: %s", err)
	//	util.WriteJSON(w, http.StatusBadRequest, util.ErrInvalidUserId)
	//	return
	//}

	//TODO: replace 0 and 10 with query
	notes, err := c.noteService.GetNotes(r.Context(), userCtx.Id)
	if err != nil {
		c.zapLogger.Errorf("Error getting notes: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.zapLogger.Infof("HandleGetNotes Ok")
	fmt.Fprintf(w, "Notes: %+v\n", notes)
	util.WriteJSON(w, notes)
}

func (c *Handler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := util.DecodeJSONBody(w, r, &user); err != nil {
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

	c.zapLogger.Infof("HandleSignUp Ok")

	//fmt.Fprintf(w, "User: %+v\n", newUser)
	util.WriteJSON(w, newUser)
}

func (c *Handler) HandleLogIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := util.DecodeJSONBody(w, r, &user); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	token, err := c.userService.LogIn(r, &user)
	if err != nil {
		c.zapLogger.Errorf("Error LogIn: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, token)
	//TODO: add redirect to /get
}

//TODO: LOGOUT

func (c *Handler) HandleWelcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Message: %+v", "Hello!")
}