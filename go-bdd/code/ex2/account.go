package ex2

import "github.com/Tigh-Gherr/presentations/go-bdd/code/ex2/errs"

type Account struct {
	Balance int64
}

func (a *Account) Withdraw(amount int64) (int64, error) {
	if a.Balance < amount {
		return 0, errs.ErrInsufficientFunds
	}

	a.Balance -= amount

	return amount, nil
}
