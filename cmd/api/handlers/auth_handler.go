package handlers

import (
	"errors"
	"fmt"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/cmd/api/services"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/mailer"
	"gorm.io/gorm"
)

func (h *Handler) RegisterHandler(c *echo.Context) error {
	payload := new(requests.RegisterUserRequest)
	err := (&echo.DefaultBinder{}).Bind(c, payload)
	if err != nil {
		fmt.Println(err)
		return common.SendBadRequestResponse(*c, err.Error())
	}

	fmt.Println(payload)
	validationErrors := h.ValidateBodyRequest(c, *payload)
	fmt.Println(validationErrors)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(*c, validationErrors)
	}

	userService := services.NewUserService(h.DB)
	user, err := userService.GetUserByEmail(payload.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) == false && user != nil {
		return common.SendBadRequestResponse(*c, "User with this email already exists")
	}

	createdUser, err := userService.RegisterUser(payload)
	if err != nil {
		fmt.Println(err)
		return common.SendInternalServerErrorResponse(*c, "User registration failed")
	}

	mailData := mailer.EmailData{
		Subject: "Welcome to " + os.Getenv("APP_NAME"),
		Meta: struct {
			FirstName string
			LoginLink string
		}{
			FirstName: payload.FirstName,
			LoginLink: "#",
		},
	}

	err = h.Mailer.Send(payload.Email, "welcome.html", mailData)
	if err != nil {
		fmt.Println(err)
	}
	return common.SendSuccessResponse(*c, "User registration successful", createdUser)
}
