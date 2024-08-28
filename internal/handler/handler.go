package handler

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"html/template"
	"io"
	"kod/internal/models"
	"kod/internal/service"
	"kod/internal/util"
	"kod/static"
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
	var note *models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		c.zapLogger.Error(util.ErrInvalidRequestPayload)
		util.WriteJSON(w, http.StatusBadRequest, util.ErrInvalidRequestPayload)
		return
	}
	if len(note.Text) == 0 {
		c.zapLogger.Error(util.ErrEmptyNote)
		util.WriteJSON(w, http.StatusBadRequest, util.ErrEmptyNote)
		return
	}

	userIdCtx, err := util.GetUserIdFromContext(r.Context())
	note.UserId = userIdCtx

	newNote, err := c.noteService.AddNote(r.Context(), note)
	if err != nil {
		c.zapLogger.Error(err)
		util.WriteJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	tmpl, err := template.New("addnote").Parse(static.AddNoteTemplate)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, note)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	c.zapLogger.Infof("HandleAddNote Ok")
	util.WriteJSON(w, http.StatusCreated, newNote)
}

func (c *Handler) HandleGetNotes(w http.ResponseWriter, r *http.Request) {
	userIdCtx, err := util.GetUserIdFromContext(r.Context())
	if err != nil {
		c.zapLogger.Error(err)
		util.WriteJSON(w, http.StatusUnauthorized, err.Error())
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//userId, err := strconv.Atoi(userIdCtx)
	//if err != nil || userId <= 0 {
	//	ns.zapLogger.Errorf("Error converting user id to int: %s", err)
	//	util.WriteJSON(w, http.StatusBadRequest, util.ErrInvalidUserId)
	//	return
	//}

	//TODO: replace 0 and 10 with query
	notes, err := c.noteService.GetNotes(r.Context(), userIdCtx)
	if err != nil {
		c.zapLogger.Errorf("Error getting notes: %s", err)
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	tmpl, err := template.New("getnotes").Parse(static.AddNoteTemplate)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, notes)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	c.zapLogger.Infof("HandleGetNotes Ok")
	util.WriteJSON(w, http.StatusOK, notes)
}

// HandleSignUpPage GET
func (c *Handler) HandleSignUpPage(w http.ResponseWriter, r *http.Request) {
	c.HandleAuthorizedUser(w, r, static.SignupTemplate)
}

// HandleSignUp POST
func (c *Handler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	userForm, err := util.ParseUserForm(r)
	if err != nil {
		c.zapLogger.Error(err)
		util.HandleHttpError(w, err.Error(), http.StatusBadRequest)

		//TODO: do i need WriteJson func?
		//util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	token, err := c.userService.SignUp(r.Context(), userForm)
	if err != nil {
		c.zapLogger.Errorf("Error SingUp: %s", err)
		util.HandleHttpError(w, err.Error(), http.StatusInternalServerError)
		//util.WriteJSON(w, http.StatusInternalServerError, err)
		return
	}
	//TODO:?
	//util.WriteJSONToken(w, token)

	c.zapLogger.Infof("HandleSignUp Ok")

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, static.LoginTemplate)

	util.WriteJSON(w, http.StatusCreated, token)
	//http.Redirect(w, r, "/notes/get", http.StatusSeeOther)
}

// HandleLogInPage GET
func (c *Handler) HandleLogInPage(w http.ResponseWriter, r *http.Request) {
	c.HandleAuthorizedUser(w, r, static.LoginTemplate)
}

// HandleLogIn POST
func (c *Handler) HandleLogIn(w http.ResponseWriter, r *http.Request) {
	userForm, err := util.ParseUserForm(r)
	if err != nil {
		c.zapLogger.Error(err)
		util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	token, err := c.userService.LogIn(r.Context(), userForm)
	if err != nil {
		c.zapLogger.Errorf("Error login: %s", err)
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	//util.WriteJSONToken(w, token)

	c.zapLogger.Infof("HandleLogIn Ok")
	util.WriteJSON(w, http.StatusOK, token)
	http.Redirect(w, r, "/notes/get", http.StatusSeeOther)
}

func (c *Handler) HandleAuthorizedUser(w http.ResponseWriter, r *http.Request, tmpl string) {
	tokenString, err := util.GetJwtToken(r)
	if err != nil {
		io.WriteString(w, tmpl)
		return
	}
	token, err := util.ValidateJwt(tokenString)
	if err != nil {
		io.WriteString(w, tmpl)
		return
	}

	userName, _ := token.Claims.(jwt.MapClaims)["userName"].(string)
	tmplLogged, err := template.New("loggedin").Parse(static.LoggedInTemplate)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmplLogged.Execute(w, userName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

//TODO: LOGOUT

func (c *Handler) HandleWelcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, static.IndexBegin)
	defer io.WriteString(w, static.IndexEnd)
	c.zapLogger.Infof("HandleWelcome Ok %v", r.Body)
	w.Write([]byte("Hello!"))
}