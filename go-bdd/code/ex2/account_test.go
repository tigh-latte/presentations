package ex2_test

import (
	"testing"

	"github.com/Tigh-Gherr/presentations/go-bdd/code/ex2"
	"github.com/cucumber/godog"
	"github.com/pkg/errors"
)

type TestContext struct {
	Account *ex2.Account
	Error   error
}

func TestAccount_Withdrawal(t *testing.T) {
	s := godog.TestSuite{
		ScenarioInitializer: InitScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if s.Run() != 0 {
		t.Fatal("holy hell")
	}
}

func InitScenario(sc *godog.ScenarioContext) {
	ctx := &TestContext{}

	sc.Step(`^I have an account with £(\d+)$`, ctx.iHaveAnAccountWith)
	sc.Step(`^I withdraw £(\d+)$`, ctx.iWithdraw)
	sc.Step(`^an error should state "([^"]+)"$`, ctx.anErrorShouldState)
	sc.Step(`^my remaining balance should be £(\d+)$`, ctx.myRemainingBalanceShouldBe)
}

func (t *TestContext) iHaveAnAccountWith(balance int64) {
	t.Account = &ex2.Account{
		Balance: balance,
	}
}

func (t *TestContext) iWithdraw(amount int64) {
	_, t.Error = t.Account.Withdraw(amount)
}

func (t *TestContext) myRemainingBalanceShouldBe(remaining int64) error {
	if t.Account.Balance != remaining {
		return errors.Errorf("unexpected balance: want %d got %d", remaining, t.Account.Balance)
	}

	return nil
}

func (t *TestContext) anErrorShouldState(err string) error {
	if t.Error == nil {
		return errors.Errorf("expecting error '%s', got nil", err)
	}

	if t.Error.Error() != err {
		return errors.Errorf("expecting error '%s', got '%s'", err, t.Error.Error())
	}

	return nil
}
