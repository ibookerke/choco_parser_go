package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ibookerke/choco_parser_go/internal/config"
)

func main() {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	conf, err := config.Get()
	if err != nil {
		slog.Error("couldn't get config", "err", err)
		return
	}

	slogHandler := slog.Handler(slog.NewTextHandler(os.Stdout, nil))
	if !conf.Project.Debug {
		slogHandler = slog.NewJSONHandler(os.Stdout, nil)
	}

	logger := slog.New(slogHandler).With("svc", conf.Project.ServiceName)

	pool, err := pgx.NewPgxPool(ctx, conf.Database.DSN)
	if err != nil {
		logger.Error("couldn't create pgx pool", "err", err)
		return
	}
	defer pool.Close()

	// migrating database scheme using migrate library
	m, err := migrate.New("file://migrations", conf.Database.DSN)
	if err != nil {
		logger.Error("couldn't create migrate instance", "err", err)
		return
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("couldn't migrate database", "err", err)
		return
	}

}
