// handlers/health.go
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type healthCheck struct {
	Health bool `json:"health"`
}

func (h *Handler) Healthcheck(c *echo.Context) error {
	return c.JSON(http.StatusOK, healthCheck{
		Health: true,
	})
}
