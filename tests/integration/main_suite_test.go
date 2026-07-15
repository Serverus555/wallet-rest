package integration

import (
	"context"
	"net"
	"os"
	"testing"
	"time"
	"wallet-rest/config"
	"wallet-rest/gen/http"
	"wallet-rest/internal/app"
	"wallet-rest/pkg/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions
	ctx    context.Context
	cancel context.CancelFunc

	pgPool *pgxpool.Pool
	client *http.ClientWithResponses
	cfg    config.Config
}

func (s *Suite) SetupSuite() {
	s.Assertions = s.Require()
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	err := godotenv.Load("../../test.env")
	s.NoError(err)

	cfg := config.Config{}
	err = envconfig.Process(ctx, &cfg)
	s.NoError(err)
	s.cfg = cfg

	s.pgPool, err = postgres.InitPool(ctx, cfg.Postgres)
	s.NoError(err)
	s.client, err = http.NewClientWithResponses("http://localhost:" + cfg.Http.Port)
	s.NoError(err)

	go func() {
		err := app.Run(ctx, cfg, initLogger(false))
		if err != nil {
			cancel()
		}
	}()

	waitServerStart(ctx, cfg.Http.Port)
}
func (s *Suite) SetupTest() {
	_, err := s.pgPool.Exec(s.ctx, "TRUNCATE TABLE wallets")
	s.NoError(err)
}
func (s *Suite) TearDownTest() {}

func (s *Suite) TearDownSuite() {
	s.cancel()
	s.pgPool.Close()
}

func waitServerStart(ctx context.Context, port string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conn, err := net.Dial("tcp", "localhost:"+port)
			if err == nil {
				_ = conn.Close()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *Suite) selectBalance(id uuid.UUID) int64 {
	var balance int64
	err := s.pgPool.QueryRow(s.ctx, "SELECT balance FROM wallets WHERE id = $1", id).Scan(&balance)
	s.NoError(err)
	return balance
}

func (s *Suite) insertWallet(id uuid.UUID, amount int64) {
	_, err := s.pgPool.Exec(s.ctx, "INSERT INTO wallets (id, balance) VALUES ($1, $2)", id, amount)
	s.NoError(err)
}

// realLogger = true для включения логов во время тестов
func initLogger(realLogger bool) zerolog.Logger {
	if !realLogger {
		return zerolog.Nop()
	}
	// For debug
	out := zerolog.ConsoleWriter{Out: os.Stdout}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(out).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	log.Logger = logger
	return logger
}
