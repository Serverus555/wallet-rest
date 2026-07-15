package http

import (
	"context"
	"errors"
	"strconv"
	"wallet-rest/gen/http"
	"wallet-rest/internal/domain"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Handler struct {
	usecase UseCase
	logger  zerolog.Logger
}

type UseCase interface {
	Deposit(ctx context.Context, id uuid.UUID, amount int64) error
	Withdraw(ctx context.Context, id uuid.UUID, amount int64) error
	GetBalance(ctx context.Context, id uuid.UUID) (int64, error)
}

func New(usecase UseCase, logger zerolog.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) Transaction(ctx context.Context, request http.TransactionRequestObject) (http.TransactionResponseObject, error) {
	var err error

	switch request.Body.OperationType {
	case http.DEPOSIT:
		err = h.usecase.Deposit(ctx, request.Body.WalletId, request.Body.Amount)
	case http.WITHDRAW:
		err = h.usecase.Withdraw(ctx, request.Body.WalletId, request.Body.Amount)
	default:
		err = domain.ErrInvalidOperation
	}

	if err != nil {
		if domainErr, ok := errors.AsType[*domain.DomainError](err); ok {
			zerolog.Ctx(ctx).Info().Err(domainErr).Msg("Handler.Transaction domain error")
			return http.Transaction400TextResponse(domainErr.Error()), nil
		}

		zerolog.Ctx(ctx).Error().Err(err).Msg("Handler.Transaction unknown error")
		return nil, err
	}

	return http.Transaction200Response{}, nil
}

func (h *Handler) GetBalance(ctx context.Context, request http.GetBalanceRequestObject) (http.GetBalanceResponseObject, error) {
	balance, err := h.usecase.GetBalance(ctx, request.Id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Handler.GetBalance unknown error")
		return nil, err
	}

	return http.GetBalance200TextResponse(strconv.FormatInt(balance, 10)), nil
}
