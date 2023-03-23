package ex1

import (
	"github.com/tigh-latte/presentations/go-bdd/code/ex1/errs"
)

type Account struct {
	Balance int64
}

func (a *Account) Withdraw(amount int64) (int64, error) {
	if a.Balance < amount {
		return 0, errs.ErrInsufficientFunds
	}

	a.Balance = a.Balance - amount

	return amount, nil
}
