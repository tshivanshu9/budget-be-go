package handlers

import (
	"github.com/tshivanshu9/budget-be/internal/mailer"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Mailer mailer.Mailer
}
