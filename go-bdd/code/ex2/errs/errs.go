package errs

type AccountError string

const (
	ErrInsufficientFunds AccountError = "insufficient funds"
	ErrEmptyDeposit      AccountError = "empty transaction"
)

func (e AccountError) Error() string {
	return string(e)
}
