package usecase

import (
	"context"
	"fmt"
	"wallet-rest/internal/domain"
	cache2 "wallet-rest/internal/repository/cache"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

//mockery:generate: true
type Repository interface {
	Deposit(ctx context.Context, id uuid.UUID, amount int64) (int64, error)
	Withdraw(ctx context.Context, id uuid.UUID, amount int64) (int64, error)
	GetBalance(ctx context.Context, id uuid.UUID) (int64, error)
}

type Cache interface {
	Put(key uuid.UUID, value int64)
	Get(key uuid.UUID) (int64, bool)
	Delete(key uuid.UUID)
}

type UseCase struct {
	repo   Repository
	cache  cache2.Cache
	logger zerolog.Logger
}

func New(repo Repository, cache Cache, logger zerolog.Logger) *UseCase {
	return &UseCase{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (u *UseCase) Deposit(ctx context.Context, id uuid.UUID, amount int64) error {
	if amount < 1 {
		return domain.ErrInvalidAmount
	}

	_, err := u.repo.Deposit(ctx, id, amount)
	if err != nil {
		return fmt.Errorf("UseCase.Deposit repo: %w", err)
	}

	u.cache.Delete(id)
	return nil
}

func (u *UseCase) Withdraw(ctx context.Context, id uuid.UUID, amount int64) error {
	if amount < 1 {
		return domain.ErrInvalidAmount
	}

	_, err := u.repo.Withdraw(ctx, id, amount)
	if err != nil {
		return fmt.Errorf("UseCase.Withdraw repo: %w", err)
	}

	u.cache.Delete(id)
	return nil
}

// GetBalance Метод возвращает баланс "для справки", а не для критического функционала
func (u *UseCase) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	var err error

	balance, cached := u.cache.Get(id)
	if !cached {
		balance, err = u.repo.GetBalance(ctx, id)
		if err != nil {
			return 0, fmt.Errorf("UseCase.GetBalance repo: %w", err)
		}
		u.cache.Put(id, balance)
	}
	return balance, err
}
