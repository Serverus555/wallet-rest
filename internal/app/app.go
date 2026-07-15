package app

import (
	"context"
	"embed"
	"fmt"
	"wallet-rest/config"
	"wallet-rest/internal/handler/http"
	"wallet-rest/internal/repository/cache"
	repository "wallet-rest/internal/repository/postgres"
	"wallet-rest/internal/usecase"
	"wallet-rest/pkg/httpserver"
	"wallet-rest/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Run(ctx context.Context, cfg config.Config, logger zerolog.Logger) error {
	pgPool, err := postgres.InitPool(ctx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer pgPool.Close()

	if err = migrate(pgPool); err != nil {
		return err
	}

	repo := repository.New(pgPool, logger)

	repoCache, err := cache.New(cfg.Cache)
	if err != nil {
		logger.Error().Err(err).Msg("app.Run: cache fallback to noop")
		repoCache = cache.NewNoop()
	}

	useCase := usecase.New(repo, repoCache, logger)
	handler := http.New(useCase, logger)
	// block
	httpserver.Run(ctx, handler, cfg.Http, logger)

	return nil
}

func migrate(pgPool *pgxpool.Pool) error {
	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return fmt.Errorf("migrate goose.SetDialect %w", err)
	}
	goose.SetBaseFS(embedMigrations)

	dbConn := stdlib.OpenDBFromPool(pgPool)

	if err := goose.Up(dbConn, "migrations"); err != nil {
		return fmt.Errorf("migrate goose.Up %w", err)
	}

	return nil
}
