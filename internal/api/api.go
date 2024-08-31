package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"kod/internal/handler"
	"kod/internal/middleware"
	"kod/internal/models/config"
	"kod/telemetry"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const (
	shutdownTimeout = 5 * time.Second
)

type API struct {
	server        *http.Server
	controller    *handler.Handler
	middleware    *middleware.Middleware
	zapLogger     *zap.SugaredLogger
	telemetryAddr string
}

func NewAPI(c *handler.Handler, m *middleware.Middleware, l *zap.SugaredLogger, hc *config.HttpConfig) *API {
	return &API{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", hc.Host, hc.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
		},
		controller:    c,
		middleware:    m,
		zapLogger:     l,
		telemetryAddr: hc.TelemetryAddr,
	}
}

func (a *API) Run(ctxBackground context.Context) {
	ctx, stop := signal.NotifyContext(ctxBackground, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go telemetry.Listen(ctx, a.zapLogger, a.telemetryAddr)

	router := mux.NewRouter()
	router.HandleFunc("/signup", a.controller.HandleSignUp).Methods("POST")
	router.HandleFunc("/login", a.controller.HandleLogIn).Methods("POST")
	router.HandleFunc("/logout", a.controller.HandleLogOut).Methods("GET")

	authRouter := router.PathPrefix("/notes").Subrouter()
	authRouter.Use(a.middleware.AuthMiddleware)
	authRouter.HandleFunc("/get", a.controller.HandleGetNotes).Methods("GET")
	authRouter.HandleFunc("/add", a.controller.HandleAddNote).Methods("POST")
	a.server.Handler = router

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.zapLogger.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()
	a.zapLogger.Infof("Listening on: %v\n", a.server.Addr)

	<-ctx.Done()
	a.zapLogger.Info("Shutting down server...\n")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.zapLogger.Errorf("shutdown: %v", err)
	}

	longShutdown := make(chan struct{}, 1)

	go func() {
		time.Sleep(3 * time.Second)
		longShutdown <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		a.zapLogger.Errorf("server shutdown: %v", ctx.Err())
	case <-longShutdown:
		a.zapLogger.Infof("finished")
	}
}