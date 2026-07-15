package integration

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"wallet-rest/gen/http"
	"wallet-rest/internal/domain"

	"github.com/google/uuid"
)

func (s *Suite) Test_DepositOk() {
	id := uuid.New()
	var amount int64 = 1000

	_, _ = s.transaction(id, amount, string(http.DEPOSIT), true)
	s.Equal(amount, s.selectBalance(id))
}

func (s *Suite) Test_WithdrawOk() {
	id := uuid.New()
	s.insertWallet(id, 1000)
	s.transaction(id, 400, string(http.WITHDRAW), true)
	s.Equal(int64(600), s.selectBalance(id))
}

func (s *Suite) Test_WithdrawAllSum() {
	id := uuid.New()
	var amount int64 = 1000
	s.insertWallet(id, 1000)
	_, _ = s.transaction(id, amount, string(http.WITHDRAW), true)
	s.Equal(int64(0), s.selectBalance(id))
}

func (s *Suite) Test_WithdrawInsufficientFunds() {
	id := uuid.New()
	var balance int64 = 1000
	s.insertWallet(id, balance)
	resp, err := s.transaction(id, 1001, string(http.WITHDRAW), false)
	s.NoError(err)
	s.Equal(balance, s.selectBalance(id))
	s.Equal(400, resp.StatusCode())
	s.Equal(domain.ErrInsufficientFunds.Error(), string(resp.Body))
}

func (s *Suite) Test_ConcurrencyOperations() {
	id := uuid.New()

	iterations := 20
	var wg sync.WaitGroup
	wg.Add(iterations)

	var expectedBalance atomic.Int64

	exec := func(op string, amount int64) {
		defer wg.Done()
		resp, _ := s.transaction(id, amount, op, false)
		switch op {
		case string(http.DEPOSIT):
			expectedBalance.Add(amount)
		case string(http.WITHDRAW):
			// Если произошла попытка снятия до появления минимальной суммы, то игнорируем вклад этого снятия в итоговую сумму
			if string(resp.Body) != domain.ErrInsufficientFunds.Error() {
				expectedBalance.Add(-amount)
			}
		}
	}

	for range iterations {
		var op string
		var amount int64

		switch rand.Intn(2) {
		// +1: Минимальный amount
		case 0:
			op = string(http.DEPOSIT)
			amount = rand.Int63n(1000) + 1
		//case 1:
		default:
			op = string(http.WITHDRAW)
			// Int63n: n <= 0 -> err
			amount = rand.Int63n(200) + 1
		}
		go exec(op, amount)
	}
	wg.Wait()

	s.Equal(expectedBalance.Load(), s.selectBalance(id))
}

func (s *Suite) transaction(id uuid.UUID, amount int64, op string, failOnErr bool) (*http.TransactionResponse, error) {
	req := http.TransactionJSONRequestBody{
		Amount:        amount,
		OperationType: http.TransactionInputOperationType(op),
		WalletId:      id,
	}
	resp, err := s.client.TransactionWithResponse(s.ctx, req)
	if failOnErr {
		s.NoError(err)
		s.Equal(200, resp.StatusCode())
	}
	return resp, err
}
