package ex2_test

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/Tigh-Gherr/presentations/go-bdd/code/ex2"
	"github.com/cucumber/godog"
	"github.com/pkg/errors"
)

type (
	accountKey struct{}
	errKey     struct{}
)

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

var opts = godog.Options{
	Format: "pretty",
	Paths:  []string{"features"},
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()
	s := godog.TestSuite{
		ScenarioInitializer: initScenario,
		Options:             &opts,
	}

	os.Exit(s.Run())
}

func initScenario(sc *godog.ScenarioContext) {
	sc.Step(`^I have an account with £(\d+)$`, iHaveAnAccountWith)
	sc.Step(`^I withdraw £(\d+)$`, iWithdraw)
	sc.Step(`^an error should state "([^"]+)"$`, anErrorShouldState)
	sc.Step(`^my remaining balance should be £(\d+)$`, myRemainingBalanceShouldBe)
	sc.Step(`^I deposit £(\d+)$`, iDeposit)
	sc.BeforeStep(func(step *godog.Step) {
		time.Sleep(500 * time.Millisecond)
	})
}

func Value[K, V any](ctx context.Context) V {
	return ctx.Value(*new(K)).(V)
}

func iHaveAnAccountWith(ctx context.Context, balance int64) context.Context {
	return context.WithValue(ctx, accountKey{}, &ex2.Account{
		Balance: balance,
	})
}

func iWithdraw(ctx context.Context, amount int64) context.Context {
	account := Value[accountKey, *ex2.Account](ctx)
	_, err := account.Withdraw(amount)
	return context.WithValue(ctx, errKey{}, err)
}

func myRemainingBalanceShouldBe(ctx context.Context, remaining int64) error {
	account := Value[accountKey, *ex2.Account](ctx)
	if account.Balance != remaining {
		return errors.Errorf("unexpected balance: want %d got %d", remaining, account.Balance)
	}

	return nil
}

func anErrorShouldState(ctx context.Context, errMsg string) error {
	err := Value[errKey, error](ctx)
	if err == nil {
		return errors.Errorf("expecting error '%s', got nil", errMsg)
	}

	if err.Error() != errMsg {
		return errors.Errorf("expecting error '%s', got '%s'", errMsg, err.Error())
	}

	return nil
}

func iDeposit(ctx context.Context, amount int64) context.Context {
	account := Value[accountKey, *ex2.Account](ctx)
	_, err := account.Deposit(amount)
	return context.WithValue(ctx, errKey{}, err)
}
