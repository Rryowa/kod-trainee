package main

import (
	"context"
	"kod/internal/api"
	"kod/internal/handler"
	"kod/internal/middleware"
	"kod/internal/service"
	"kod/internal/storage/postgres"
	"kod/internal/util"
)

func main() {
	ctx := context.Background()
	zapLogger := util.NewZapLogger()
	httpCfg := util.NewHttpConfig()

	storage := postgres.NewPostgresRepository(ctx, util.NewDbConfig(), zapLogger)
	noteService := service.NewNoteService(storage)
	userService := service.NewUserService(storage)
	midService := middleware.NewMiddleware(zapLogger)
	ctrl := handler.NewHandler(noteService, userService, zapLogger)
	app := api.NewAPI(ctrl, midService, zapLogger, httpCfg)

	app.Run()
}