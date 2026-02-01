package middlewares

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

func CustomMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	fmt.Println("inside the custom middleware")
	return func(c *echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
	}
}
