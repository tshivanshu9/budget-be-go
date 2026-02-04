package handlers

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/internal/mailer"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Mailer mailer.Mailer
}

func (h *Handler) BindRequestBody(c *echo.Context, payload interface{}) error {
	if err := (&echo.DefaultBinder{}).Bind(c, payload); err != nil {
		fmt.Println(err)
		return fmt.Errorf("Failed to Bind request body: %w", err)
	}
	return nil
}
