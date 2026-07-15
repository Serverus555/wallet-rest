package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UseCaseMocks struct {
	repo    *MockRepository
	cache   *MockCache
	useCase *UseCase
}

func Test_TransactionDeleteCache(t *testing.T) {
	mocks := setupUseCase(t)

	mocks.cache.EXPECT().Delete(mock.Anything).Once()
	mocks.repo.EXPECT().Withdraw(mock.Anything, mock.Anything, mock.Anything).Return(0, nil)

	mocks.useCase.Withdraw(context.Background(), uuid.New(), 1000)
}

func Test_GetBalanceSaveCache(t *testing.T) {
	mocks := setupUseCase(t)

	mocks.cache.EXPECT().Get(mock.Anything).Return(0, false).Once()
	mocks.repo.EXPECT().GetBalance(mock.Anything, mock.Anything).Return(500, nil)
	mocks.cache.EXPECT().Put(mock.Anything, int64(500)).Once()

	mocks.useCase.GetBalance(context.Background(), uuid.New())
}

func Test_GetBalanceUseCache(t *testing.T) {
	mocks := setupUseCase(t)
	balance := int64(500)

	mocks.cache.EXPECT().Get(mock.Anything).Return(balance, true).Once()

	result, _ := mocks.useCase.GetBalance(context.Background(), uuid.New())
	assert.Equal(t, balance, result)
}

func setupUseCase(t *testing.T) UseCaseMocks {
	repo := NewMockRepository(t)
	cache := NewMockCache(t)

	useCase := UseCase{
		repo:   repo,
		cache:  cache,
		logger: zerolog.Nop(),
	}

	return UseCaseMocks{
		repo:    repo,
		cache:   cache,
		useCase: &useCase,
	}
}
