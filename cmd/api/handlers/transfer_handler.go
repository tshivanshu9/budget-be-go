package handlers

import (
	"sync"

	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/cmd/api/services"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/models"
)

func (h *Handler) TransferHanlder(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}

	payload := new(requests.TransferRequest)
	err := h.BindRequestBody(c, payload)
	if err != nil {
		return common.SendBadRequestResponse(c, "invalid payload")
	}

	validationErrors := h.ValidateBodyRequest(c, payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	walletService := services.NewWalletService(h.DB)

	var sourceWallet, destinationWallet *models.WalletModel
	var sourceErr, destinationErr error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		sourceWallet, sourceErr = walletService.GetWalletByIdAndUserId(payload.SourceWalletId, user.ID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		destinationWallet, destinationErr = walletService.GetWalletByIdAndUserId(payload.DestinationWalletId, user.ID)
	}()

	wg.Wait()

	if sourceErr != nil || destinationErr != nil {
		return common.SendBadRequestResponse(c, "invalid wallet ids")
	}

	if payload.Amount >= sourceWallet.Balance {
		return common.SendBadRequestResponse(c, "insufficient wallet balance")
	}

	transferService := services.NewTransferService(h.DB)
	err = transferService.Transfer(sourceWallet, destinationWallet, payload.Amount, user.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "transfer failed")
	}

	return common.SendSuccessResponse(c, "transfer successful", nil)
}
