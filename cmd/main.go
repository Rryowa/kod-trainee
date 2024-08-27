package main

import (
	"context"
	"kod/internal/storage/postgres"
	"kod/internal/util"
)

func main() {
	ctx := context.Background()
	zapLogger := util.NewZapLogger()

	repository := postgres.NewPostgresRepository(ctx, util.NewDbConfig(), zapLogger)

}