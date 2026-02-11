package services

import (
	"time"

	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/internal/models"
	"gorm.io/gorm"
)

type TransferService struct {
	DB *gorm.DB
}

const transfer = "transfer"

func NewTransferService(db *gorm.DB) *TransferService {
	return &TransferService{DB: db}
}

func (ts *TransferService) Transfer(sourceWallet *models.WalletModel, destinationWallet *models.WalletModel, amount float64, userId uint) error {
	categoryService := NewCategoryService(ts.DB)
	category, err := categoryService.FindBySlug(transfer, false)
	if err != nil {
		return err
	}

	currentDate := time.Now()

	title := "Account transfer"
	sourceTransactionData := requests.CreateTransactionRequest{
		CategoryId:  &category.ID,
		WalletId:    sourceWallet.ID,
		Amount:      amount,
		Title:       &title,
		Date:        currentDate.Format(time.DateOnly),
		Description: &title,
		Type:        EXPENSE,
	}

	destinationTransactionData := requests.CreateTransactionRequest{
		CategoryId:  &category.ID,
		WalletId:    destinationWallet.ID,
		Amount:      amount,
		Title:       &title,
		Date:        currentDate.Format(time.DateOnly),
		Description: &title,
		Type:        INCOME,
	}

	err = ts.DB.Transaction(func(tx *gorm.DB) error {
		transactionService := NewTransactionService(tx)
		walletService := NewWalletService(tx)
		txErr := walletService.DecrementWalletBalance(sourceWallet, amount)
		if txErr != nil {
			return txErr
		}

		_, txErr = transactionService.Create(&sourceTransactionData, userId, false, &currentDate)
		if txErr != nil {
			return txErr
		}

		txErr = walletService.IncrementWalletBalance(destinationWallet, amount)
		if txErr != nil {
			return txErr
		}

		_, txErr = transactionService.Create(&destinationTransactionData, userId, false, &currentDate)
		if txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
