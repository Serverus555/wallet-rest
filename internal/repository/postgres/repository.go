package repository

import (
	"context"
	"errors"
	"fmt"
	"wallet-rest/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

var (
	balanceConstraintName = "chk_balance_positive"
)

type Repository struct {
	pool   *pgxpool.Pool
	logger zerolog.Logger
}

// Тут нет блокировок, работаем исходя из atomic insert/update

func New(pool *pgxpool.Pool, logger zerolog.Logger) *Repository {
	return &Repository{pool: pool, logger: logger}
}

func (r *Repository) Deposit(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	if amount < 1 {
		return 0, domain.ErrInvalidAmount
	}

	const sql = `INSERT INTO wallets (id, balance)
					VALUES ($1, $2)
					ON CONFLICT (id) DO UPDATE
					SET balance = wallets.balance + $2
					RETURNING balance`

	var balance int64
	err := r.pool.QueryRow(ctx, sql, id, amount).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("Repository.Deposit Scan: %w", err)
	}

	return balance, nil
}

func (r *Repository) Withdraw(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	if amount < 1 {
		return 0, domain.ErrInvalidAmount
	}
	// Проверка на недостаточно средств реализована в constraint ради скорости
	const sql = `UPDATE wallets SET balance = balance - $2 WHERE id = $1 RETURNING balance`

	var balance int64
	err := r.pool.QueryRow(ctx, sql, id, amount).Scan(&balance)

	if err != nil {
		switch pgErr, isPgErr := errors.AsType[*pgconn.PgError](err); {
		case errors.Is(err, pgx.ErrNoRows):
			// От Scan если ничего не обновили
			return 0, domain.ErrInsufficientFunds
		case isPgErr && pgErr.ConstraintName == balanceConstraintName:
			// От postgres при нарушении constraint
			return 0, domain.ErrInsufficientFunds

		default:
			return 0, fmt.Errorf("Repository.Withdraw Scan: %w", err)
		}
	}

	return balance, nil
}

func (r *Repository) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	const sql = `SELECT balance FROM wallets WHERE id = $1;`
	var balance int64
	err := r.pool.QueryRow(ctx, sql, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("Repository.GetBalance Scan: %w", err)
	}
	return balance, nil
}
