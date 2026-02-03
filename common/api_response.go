package common

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type ValidationError struct {
	Error     string `json:"error"`
	Key       string `json:"key"`
	Condition string `json:"condition"`
}

type JSONSuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JSONFailedValidationResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Errors  []*ValidationError `json:"errors"`
}

type JSONErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SendSuccessResponse(c *echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, JSONSuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendFailedValidationResponse(c *echo.Context, errors []*ValidationError) error {
	return c.JSON(http.StatusUnprocessableEntity, JSONFailedValidationResponse{
		Success: false,
		Message: "Validation Failed",
		Errors:  errors,
	})
}

func SendErrorResponse(c *echo.Context, message string, statusCode int) error {
	return c.JSON(statusCode, JSONErrorResponse{
		Success: false,
		Message: message,
	})
}

func SendBadRequestResponse(c *echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusBadRequest)
}

func SendInternalServerErrorResponse(c *echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusInternalServerError)
}

func SendUnauthorizedResponse(c *echo.Context, message *string) error {
	if message == nil {
		defaultMessage := "Unauthorized"
		message = &defaultMessage
	}
	return SendErrorResponse(c, *message, http.StatusUnauthorized)
}
