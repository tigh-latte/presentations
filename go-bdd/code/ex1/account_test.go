package ex1_test

import (
	"testing"

	"github.com/tigh-latte/presentations/go-bdd/code/ex1"
	"github.com/tigh-latte/presentations/go-bdd/code/ex1/errs"
	"gotest.tools/assert"
)

func TestAccount_Withdrawal(t *testing.T) {
	tests := map[string]struct {
		balance  int64
		withdraw int64

		expBalance int64
		expErr     error
	}{
		"successful withdraw": {
			balance:    100,
			withdraw:   50,
			expBalance: 50,
		},
		"not enough funds": {
			balance:  50,
			withdraw: 100,

			expBalance: 50,
			expErr:     errs.ErrInsufficientFunds,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			account := ex1.Account{
				Balance: test.balance,
			}
			_, err := account.Withdraw(test.withdraw)

			assert.Equal(t, test.expBalance, account.Balance)
			assert.Equal(t, test.expErr, err)
		})
	}
}
