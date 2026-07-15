package httpserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	genhttp "wallet-rest/gen/http"
	httphandler "wallet-rest/internal/handler/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type HttpConfig struct {
	Port string `env:"HTTP_PORT,required"`
}

func Run(ctx context.Context, handler *httphandler.Handler, config HttpConfig, logger zerolog.Logger) {
	router := gin.Default()

	strictHandler := genhttp.NewStrictHandler(handler, nil)
	genhttp.RegisterHandlers(router, strictHandler)

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		// Закрытие по SIG
		logger.Info().Msg("Http server shutdown signal")
		gracefulCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		if err := srv.Shutdown(gracefulCtx); err != nil {
			logger.Warn().Err(err).Msg("Http server: shutdown timeout?")
		}
		cancel()
	case <-ctx.Done():
		// Закрытие по контексту
		logger.Info().Msg("Http server shutdown by context")
		gracefulCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := srv.Shutdown(gracefulCtx); err != nil {
			logger.Warn().Err(err).Msg("Http server: shutdown by context timeout?")
		}
		cancel()
	case err := <-errChan:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("Http server: crashed")
		}
	}

	logger.Info().Msg("Http server: Stopped")
}
