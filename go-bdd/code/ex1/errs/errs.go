package errs

type Ex1Err string

const (
	ErrInsufficientFunds Ex1Err = "insufficient funds"
)

func (e Ex1Err) Error() string {
	return string(e)
}
