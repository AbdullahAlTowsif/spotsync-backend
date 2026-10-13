package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"spotsync-backend/internal/config"
	"spotsync-backend/internal/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadEnv()
	if err != nil {
		return err
	}

	db, err := config.ConnectDB(cfg)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	// From Phase 4 on, db is passed into repositories here (manual DI).

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(cfg)
	return srv.Run(ctx)
}
