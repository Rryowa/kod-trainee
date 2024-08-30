package main

import (
	"context"
	"kod/internal/api"
	"kod/internal/handler"
	"kod/internal/middleware"
	"kod/internal/service"
	"kod/internal/storage/postgres"
	"kod/internal/util"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	zapLogger := util.NewZapLogger()
	dbConfig := util.NewDbConfig()
	httpCfg := util.NewHttpConfig()
	sesConfig := util.NewSessionConfig()

	storage := postgres.NewPostgresRepository(ctx, dbConfig, zapLogger)

	noteService := service.NewNoteService(storage)
	sessionService := service.NewSessionService(sesConfig)
	userService := service.NewUserService(storage, sessionService)

	middlewareService := middleware.NewMiddleware(sessionService, zapLogger)

	handlerController := handler.NewHandler(noteService, userService, zapLogger)

	app := api.NewAPI(handlerController, middlewareService, zapLogger, httpCfg)

	app.Run(ctx)
}