package domain

var (
	ErrInvalidOperation  = NewError("operation type not found")
	ErrInsufficientFunds = NewError("insufficient funds")
	ErrInvalidAmount     = NewError("invalid amount")
)

type DomainError struct {
	Msg string
}

func (e DomainError) Error() string {
	return e.Msg
}

func NewError(msg string) error {
	return &DomainError{Msg: msg}
}
