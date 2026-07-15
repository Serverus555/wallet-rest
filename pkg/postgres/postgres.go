package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	DBName   string `env:"DB_NAME,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	MinConns int    `env:"DB_MIN_CONNS"`
	MaxConns int    `env:"DB_MAX_CONNS"`
}

func InitPool(ctx context.Context, c Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("user=%s password=%s port=%s host=%s dbname=%s pool_min_conns=%d pool_max_conns=%d",
		c.User,
		c.Password,
		c.Port,
		c.Host,
		c.DBName,
		c.MinConns,
		c.MaxConns,
	)

	cfg, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
