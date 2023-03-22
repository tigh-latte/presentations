package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/bank/api/rest"
	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/bank/repo"
	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/bank/service"
	"github.com/labstack/echo"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	db, err := sql.Open("postgres", "postgresql://dev:dev@0.0.0.0/dev?sslmode=disable")
	if err != nil {
		return errors.Wrap(err, "failed to connect to database")
	}

	e := echo.New()
	g := e.Group("/api/v1")

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if err == nil {
			return
		}

		fmt.Println(err)

		e.DefaultHTTPErrorHandler(err, c)
	}

	repo := repo.NewAccountRepo(db)
	svc := service.NewAccountService(repo)
	rest.NewRESTAccount(svc).Register(g)

	go func() {
		if err := e.Start(":8081"); err != nil {
			panic(err)
		}
	}()

	q := make(chan os.Signal, 2)
	signal.Notify(q, os.Interrupt, syscall.SIGTERM)
	<-q

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}
