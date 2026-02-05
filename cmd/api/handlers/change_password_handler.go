package handlers

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/cmd/api/services"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/models"
)

func (h *Handler) ChangeUserPassword(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}

	payload := new(requests.ChangePasswordRequest)
	err := h.BindRequestBody(c, payload)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	fmt.Println(user)

	if !common.ComparePasswordHash(payload.CurrentPassword, user.Password) {
		return common.SendBadRequestResponse(c, "Current password is incorrect")
	}

	userService := services.NewUserService(h.DB)
	err = userService.ChangeUserPassword(payload.Password, user.ID)
	if err != nil {
		fmt.Println(err)
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Password changed successfully", nil)
}
