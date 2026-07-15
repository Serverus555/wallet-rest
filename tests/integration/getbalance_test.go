package integration

import (
	"strconv"

	"github.com/google/uuid"
)

func (s *Suite) Test_Ok() {
	id := uuid.New()
	var amount int64 = 1000

	s.insertWallet(id, 1000)

	resp, err := s.client.GetBalanceWithResponse(s.ctx, id)
	s.NoError(err)
	s.Equal(strconv.FormatInt(amount, 10), string(resp.Body))
}

func (s *Suite) Test_OkNoWallet() {
	id := uuid.New()

	resp, err := s.client.GetBalanceWithResponse(s.ctx, id)
	s.NoError(err)
	s.Equal("0", string(resp.Body))
}
