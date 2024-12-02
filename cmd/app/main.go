package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/qvvan/test-jwt/internal/app/api"
	v1 "github.com/qvvan/test-jwt/internal/app/api/v1"
	"github.com/qvvan/test-jwt/internal/app/repository"
	"github.com/qvvan/test-jwt/internal/config"
	"github.com/qvvan/test-jwt/pkg/client/postgresql"
	"github.com/qvvan/test-jwt/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.LogLevel)

	fmt.Println(cfg.PgDSN)
	db, err := postgresql.NewClient(ctx, cfg.PgDSN, 5, 5*time.Second)
	if err != nil {
		log.Error("failed to connect to database", slog.Any("err", err))
		os.Exit(1)
	}

	factory := repository.NewFactory(db, log)

	manager := v1.NewManager(factory, log)

	r := api.InitRouters(cfg, manager)

	if errRun := r.Run(cfg.HttpServer.Address); errRun != nil {
		log.Error("failed to run server http server", slog.Any("err", errRun))
	}
}
