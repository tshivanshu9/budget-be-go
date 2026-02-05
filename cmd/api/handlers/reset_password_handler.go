package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/cmd/api/services"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/mailer"
	"gorm.io/gorm"
)

func (h *Handler) ForgotPasswordHandler(c *echo.Context) error {
	payload := new(requests.ForgotPasswordRequest)
	err := h.BindRequestBody(c, payload)
	if err != nil {
		fmt.Println(err)
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)
	user, err := userService.GetUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "User with this email does not exist")
		}
		return common.SendInternalServerErrorResponse(c, "Error retrieving user")
	}

	token, err := appTokenService.GenerateResetPasswordToken(user)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occured, try again later")
	}

	encodedEmail := base64.RawURLEncoding.EncodeToString([]byte(user.Email))
	frontendUrl, err := url.Parse(payload.FrontendURL)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid frontend URL")
	}

	query := url.Values{}
	query.Set("email", encodedEmail)
	query.Set("token", token.Token)

	frontendUrl.RawQuery = query.Encode()

	mailData := mailer.EmailData{
		Subject: "Welcome to " + os.Getenv("APP_NAME"),
		Meta: struct {
			Token       string
			FrontendUrl string
		}{
			Token:       token.Token,
			FrontendUrl: frontendUrl.String(),
		},
	}

	err = h.Mailer.Send(payload.Email, "forgot-password.html", mailData)
	if err != nil {
		fmt.Println(err)
	}
	return common.SendSuccessResponse(c, "Password reset email sent successfully", nil)
}

func (h *Handler) ResetPasswordHandler(c *echo.Context) error {
	payload := new(requests.ResetPasswordRequest)
	err := h.BindRequestBody(c, payload)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	email, err := base64.RawURLEncoding.DecodeString(payload.Meta)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occured, try again later")
	}

	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)

	retrievedUser, err := userService.GetUserByEmail(string(email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "Invalid password reset token")
		}
		return common.SendInternalServerErrorResponse(c, "An error occured, try again later")
	}

	retrievedToken, err := appTokenService.ValidateResetPasswordToken(retrievedUser, payload.Token)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	err = userService.ChangeUserPassword(payload.Password, retrievedUser.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occured, try again later")
	}

	appTokenService.InvalidateToken(retrievedUser.ID, retrievedToken)
	return common.SendSuccessResponse(c, "Password Reset Successfully", nil)
}
