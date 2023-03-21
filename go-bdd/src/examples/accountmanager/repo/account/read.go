package account

const (
	SelectAll = `SELECT id, balance, email FROM accounts`

	GetByID = `SELECT id, balance, email FROM accounts WHERE id = $1`
)
