package rest

import (
	"net/http"
	"strconv"

	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/bank"
	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/bank/errs"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type accountRest struct {
	svc bank.AccountService
}

func NewRESTAccount(svc bank.AccountService) bank.RESTResource {
	return &accountRest{svc: svc}
}

func (a *accountRest) Register(e *echo.Group) {
	e.GET("/accounts", a.list)
	e.POST("/accounts", a.create)
	e.GET("/accounts/:id", a.get)

	e.GET("/accounts/:id/balance", a.getBalance)
	e.GET("/accounts/:id/action", a.postAction)
}

func (a *accountRest) list(c echo.Context) error {
	accs, err := a.svc.List(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, accs)
}

func (a *accountRest) get(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return errs.ErrBadRequest
	}

	acc, err := a.svc.GetByID(c.Request().Context(), bank.GetAccountArgs{
		ID: id,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return c.JSON(http.StatusOK, acc)
}

func (a *accountRest) create(e echo.Context) error {
	acc, err := unmarshal[bank.Account](e.Request().Body)
	if err != nil {
		return err
	}

	acc, err = a.svc.Create(e.Request().Context(), acc)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusCreated, acc)
}

func (a *accountRest) getBalance(e echo.Context) error {
	id, err := strconv.ParseInt(e.Param("id"), 10, 64)
	if err != nil {
		return errs.ErrBadRequest
	}

	acc, err := a.svc.GetByID(e.Request().Context(), bank.GetAccountArgs{
		ID: id,
	})
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, struct {
		Balance int64 `json:"balance"`
	}{
		Balance: acc.Balance,
	})
}

func (a *accountRest) postAction(e echo.Context) error {
	id, err := strconv.ParseInt(e.Param("id"), 10, 64)
	if err != nil {
		return errs.ErrBadRequest
	}

	acc, err := a.svc.GetByID(e.Request().Context(), bank.GetAccountArgs{
		ID: id,
	})
	if err != nil {
		return err
	}

	_ = acc
	return e.NoContent(http.StatusNoContent)
}
