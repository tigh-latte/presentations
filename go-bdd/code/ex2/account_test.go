package ex2_test

import (
	"os"
	"testing"
	"time"

	"github.com/Tigh-Gherr/presentations/go-bdd/code/ex2"
	"github.com/cucumber/godog"
	"github.com/pkg/errors"
)

func TestMain(m *testing.M) {
	s := godog.TestSuite{
		ScenarioInitializer: initScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features"},
		},
	}

	os.Exit(s.Run())
}

type TestContext struct {
	Account *ex2.Account
	Error   error
}

func initScenario(sc *godog.ScenarioContext) {
	ctx := &TestContext{}

	sc.Step(`^I have an account with £(\d+)$`, ctx.iHaveAnAccountWith)
	sc.Step(`^I withdraw £(\d+)$`, ctx.iWithdraw)
	sc.Step(`^an error should state "([^"]+)"$`, ctx.anErrorShouldState)
	sc.Step(`^my remaining balance should be £(\d+)$`, ctx.myRemainingBalanceShouldBe)

	sc.BeforeStep(func(step *godog.Step) {
		time.Sleep(500 * time.Millisecond)
	})
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
