package errs

type Error string

const (
	ErrBadRequest Error = "bad request"

	ErrInsufficientFunds Error = "insufficient funds"
)

func (e Error) Error() string {
	return string(e)
}
