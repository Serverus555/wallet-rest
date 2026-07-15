package main

import (
	"context"
	"os"
	"wallet-rest/config"
	"wallet-rest/internal/app"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

//go:generate go tool mockery

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = logger
	ctx := context.Background()

	cfg := config.Config{}
	if err := envconfig.Process(ctx, &cfg); err != nil {
		logger.Fatal().Err(err).Msg("envconfig.Process")
	}

	err := app.Run(ctx, cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("app.Run")
	}
}
