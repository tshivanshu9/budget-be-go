package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/cmd/api/services"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/models"
)

func (h *Handler) CreateWalletHandler(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}
	payload := new(requests.CreateWalletRequest)
	err := h.BindRequestBody(c, payload)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	walletService := services.NewWalletService(h.DB)
	wallet, _ := walletService.WalletExistsForNameAndUserId(payload.Name, user.ID)
	if wallet != nil {
		return common.SendBadRequestResponse(c, "wallet already exists")
	}

	wallet, err = walletService.Create(payload, user.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "wallet creation failed, try again later")
	}
	return common.SendSuccessResponse(c, "wallet created successfully", wallet)
}

func (h *Handler) GenerateDefaultWalletsHandler(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}

	walletService := services.NewWalletService(h.DB)
	wallets, err := walletService.GenerateDefaultWallets(user.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Failed to generate default wallets, try again later")
	}
	return common.SendSuccessResponse(c, "Default wallets generated successfully", wallets)
}
