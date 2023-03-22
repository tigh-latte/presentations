package integration_test

import (
	"flag"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestMain(m *testing.M) {
	flag.Parse()

	suite := godog.TestSuite{
		Name: "bank",
		Options: &godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  flag.Args(),
		},
		ScenarioInitializer: initScenario,
	}

	os.Exit(suite.Run())
}

func initScenario(sc *godog.ScenarioContext) {
	sc.Step(`^I am unauthenticated`, func() {
		// unauthenticate
	})
}
