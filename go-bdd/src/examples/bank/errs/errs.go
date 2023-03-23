package errs

type Error string

const (
	ErrUnauthorized Error = "unauthorized"

	ErrBadRequest Error = "bad request"

	ErrInsufficientFunds Error = "insufficient funds"

	ErrConflict Error = "resource conflict"
)

func (e Error) Error() string {
	return string(e)
}
