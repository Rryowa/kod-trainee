package api

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"kod/internal/handler"
	"kod/internal/middleware"
	"kod/internal/models/config"
	"net/http"
	"time"
)

type API struct {
	server     *http.Server
	controller *handler.Handler
	middleware *middleware.Middleware
	zapLogger  *zap.SugaredLogger
}

func NewAPI(c *handler.Handler, m *middleware.Middleware, l *zap.SugaredLogger, hc *config.HttpConfig) *API {
	return &API{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", hc.Host, hc.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
		},
		controller: c,
		middleware: m,
		zapLogger:  l,
	}
}

func (a *API) Run() {
	a.zapLogger.Infof("Listening on: %v", a.server.Addr)
	router := mux.NewRouter()
	//router.Host("notes")
	router.HandleFunc("/", a.controller.HandleWelcome)
	router.HandleFunc("/signup", a.controller.HandleSignUp).Methods("POST")
	router.HandleFunc("/login", a.controller.HandleLogIn).Methods("POST")

	authRouter := router.PathPrefix("/notes").Subrouter()
	authRouter.Use(a.middleware.AuthMiddleware)
	authRouter.HandleFunc("/get", a.controller.HandleGetNotes).Methods("GET")
	authRouter.HandleFunc("/add", a.controller.HandleAddNote).Methods("POST")
	a.server.Handler = router

	//TODO: gracefull shutdown
	if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		a.zapLogger.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}