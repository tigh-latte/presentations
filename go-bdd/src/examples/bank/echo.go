package bank

import "github.com/labstack/echo"

type RESTResource interface {
	Register(e *echo.Group)
}
