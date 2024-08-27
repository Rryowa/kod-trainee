package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"kod/internal/middleware"
	"kod/internal/models/config"
	"kod/internal/service"
	"kod/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var m middleware.Middleware

type API struct {
	server      *http.Server
	noteService *service.NoteService
	middleware  *middleware.Middleware
	zapLogger   *zap.SugaredLogger
}

func NewAPI(s storage.Storage, l *zap.SugaredLogger, cfg *config.HttpConfig) *API {
	return &API{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
		},
		noteService: service.NewNoteService(s),
		middleware:  middleware.NewMiddleware(),
		zapLogger:   l,
	}
}

func (a *API) Run() {
	router := mux.NewRouter()
	router.Host("api")
	router.Use(a.middleware.AuthMiddleware)
	router.HandleFunc("/", a.noteService.HandleWelcome)
	router.HandleFunc("/notes/{user_id}", a.noteService.HandleGetNotes).Methods("GET")
	router.HandleFunc("/notes", a.noteService.HandleAddNote).Methods("POST")
	a.server.Handler = router

	//authMiddleware :=

	idleConnClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := a.server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnClosed)
	}()

	if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnClosed
}