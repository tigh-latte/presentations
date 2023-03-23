package middleware

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/tigh-latte/presentations/go-bdd/src/examples/bank"
	"github.com/tigh-latte/presentations/go-bdd/src/examples/bank/errs"
)

func Auth(repo bank.APIKeyRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			bearer := e.Request().Header.Get("Authorization")
			if bearer == "" {
				return errs.ErrUnauthorized
			}
			key, ok := strings.CutPrefix(bearer, "Bearer ")
			if !ok {
				return errs.ErrUnauthorized
			}
			apiKey, err := repo.APIKey(e.Request().Context(), bank.GetAPIKeyArgs{APIKey: key})
			if err != nil {
				return errs.ErrUnauthorized
			}
			if apiKey.APIKey != key {
				return errs.ErrUnauthorized
			}
			return next(e)
		}
	}
}
