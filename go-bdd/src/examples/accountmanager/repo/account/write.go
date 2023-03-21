package account

const (
	CreateAccount = `
	INSERT INTO accounts (email, balance) VALUES ($1, $2) RETURNING id
	`
)
