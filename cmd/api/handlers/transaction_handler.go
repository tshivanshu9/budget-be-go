package handlers

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/cmd/api/filters"
	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/cmd/api/services"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/models"
	"gorm.io/gorm"
)

func (h *Handler) CreateTransactionHandler(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}
	payload := new(requests.CreateTransactionRequest)
	err := h.BindRequestBody(c, payload)

	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	walletService := services.NewWalletService(h.DB)
	transactionService := services.NewTransactionService(h.DB)
	formattedDate, err := transactionService.FormatDate(payload.Date)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid date format, expected YYYY-MM-DD")
	}

	wallet, err := walletService.GetWalletByIdAndUserId(payload.WalletId, user.ID)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid wallet")
	}

	var category *models.CategoryModel
	if payload.CategoryId != nil {
		categoryService := services.NewCategoryService(h.DB)
		retrievedCategory, err := categoryService.GetById(*payload.CategoryId)
		if err != nil {
			return common.SendBadRequestResponse(c, "Invalid category")
		}
		category = retrievedCategory
	}

	var transaction *models.TransactionModel

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		walletService := services.NewWalletService(tx)
		transactionService.DB = tx
		if payload.Type == services.INCOME {
			txErr := walletService.IncrementWalletBalance(wallet, payload.Amount)
			if txErr != nil {
				return errors.New("Transaction creation failed, try again later")
			}
		}
		if payload.Type == services.EXPENSE {
			txErr := walletService.DecrementWalletBalance(wallet, payload.Amount)
			if txErr != nil {
				return errors.New("Transaction creation failed, try again later")
			}
		}

		createdTransaction, txErr := transactionService.Create(payload, user.ID, false, formattedDate)
		if txErr != nil {
			return errors.New("Transaction creation failed, try again later")
		}
		transaction = createdTransaction
		transaction.Category = category
		transaction.Wallet = wallet
		return nil
	})

	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Transaction creation failed, try again later")
	}

	return common.SendSuccessResponse(c, "Transaction created successfully", transaction)
}

func (h *Handler) ReverseTransactionHandler(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}

	transactionIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid transaction id")
	}
	transactionId := uint(transactionIdInt)
	transactionService := services.NewTransactionService(h.DB)
	transaction, err := transactionService.FindById(transactionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendBadRequestResponse(c, "Transaction not found")
		}
		return common.SendInternalServerErrorResponse(c, "Failed to fetch transaction, try again later")
	}

	if transaction.UserId != user.ID {
		return common.SendUnauthorizedResponse(c, nil)
	}
	if transaction.IsReversal {
		return common.SendBadRequestResponse(c, "Cannot reverse an already reversed transaction")
	}

	transactionWithParentId, err := transactionService.FindByParentId(transaction.ID)
	if err == nil && transactionWithParentId != nil {
		return common.SendBadRequestResponse(c, "Cannot reverse a transaction that has already been reversed")
	}
	err = transactionService.Reverse(transaction)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Transaction reversed successfully", nil)
}

func (h *Handler) ListTransactionsForUserHandler(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}

	filter := new(filters.TransactionFilter)
	err := (&echo.DefaultBinder{}).Bind(c, filter)
	if err != nil {
		return common.SendBadRequestResponse(c, "invalid query params")
	}

	err = filter.ValidateDates()
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	query := filter.ApplyFilters(h.DB)
	transactionService := services.NewTransactionService(query)

	var transactions []*models.TransactionModel
	pagination := common.NewPaginator(transactions, c.Request(), query)
	paginatedTransactions, err := transactionService.List(transactions, user.ID, pagination)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Failed to fetch transactions, try again later")
	}
	return common.SendSuccessResponse(c, "Transactions fetched successfully", paginatedTransactions)
}
