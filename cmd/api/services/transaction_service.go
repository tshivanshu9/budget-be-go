package services

import (
	"errors"
	"time"

	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/models"
	"gorm.io/gorm"
)

type TransactionService struct {
	DB *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{DB: db}
}

const INCOME = "income"
const EXPENSE = "expense"

func (t *TransactionService) Create(payload *requests.CreateTransactionRequest, userId uint, isReversal bool, date *time.Time) (*models.TransactionModel, error) {
	transaction := &models.TransactionModel{
		UserId:      userId,
		Amount:      payload.Amount,
		Description: payload.Description,
		Title:       payload.Title,
		Type:        payload.Type,
		Date:        *date,
		Month:       uint8(date.Month()),
		Year:        uint16(date.Year()),
		IsReversal:  isReversal,
		WalletId:    payload.WalletId,
		CategoryId:  payload.CategoryId,
	}
	if isReversal {
		transaction.ParentId = payload.ParentId
	}
	result := t.DB.Create(transaction)
	if result.Error != nil {
		return nil, result.Error
	}

	budgetService := NewBudgetService(t.DB)
	if transaction.Type == EXPENSE {
		budgetService.DecrementBudgetBalance(t.DB, payload.CategoryId, payload.Amount, userId)
	} else if transaction.Type == INCOME && isReversal {
		budgetService.IncrementBudgetBalance(t.DB, payload.CategoryId, payload.Amount, userId)
	}
	return transaction, nil
}

func (t *TransactionService) FormatDate(date string) (*time.Time, error) {
	currentTime := time.Now()
	if date == "" {
		return &currentTime, nil
	}
	suppliedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}
	suppliedDateTime := time.Date(suppliedDate.Year(), suppliedDate.Month(), suppliedDate.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond(), time.UTC)
	return &suppliedDateTime, nil
}

func (transactionService *TransactionService) FindById(id uint) (*models.TransactionModel, error) {
	var transaction models.TransactionModel
	result := transactionService.DB.Joins("Wallet").Joins("Category").Where("transactions.id = ?", id).First(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

func (transactionService *TransactionService) FindByParentId(id uint) (*models.TransactionModel, error) {
	var transaction models.TransactionModel
	result := transactionService.DB.Where("parent_id = ?", id).First(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

func (t *TransactionService) Reverse(transaction *models.TransactionModel) error {
	transactionDate := time.Now()
	description := "Reversal of transaction ID " + string(rune(transaction.ID))
	transactionRequest := requests.CreateTransactionRequest{
		ParentId:    &transaction.ID,
		WalletId:    transaction.WalletId,
		Description: &description,
		Date:        transactionDate.Format(time.DateOnly),
		Amount:      transaction.Amount,
	}
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		walletService := NewWalletService(tx)
		if transaction.Type == INCOME {
			transaction.Type = EXPENSE
			txErr := walletService.DecrementWalletBalance(transaction.Wallet, transaction.Amount)
			if txErr != nil {
				return errors.New("Transaction reversal failed, try again later")
			}
		}
		if transaction.Type == EXPENSE {
			transaction.Type = INCOME
			txErr := walletService.IncrementWalletBalance(transaction.Wallet, transaction.Amount)
			if txErr != nil {
				return errors.New("Transaction reversal failed, try again later")
			}
		}
		_, txErr := t.Create(&transactionRequest, transaction.UserId, true, &transactionDate)
		if txErr != nil {
			return errors.New("Transaction reversal failed, try again later")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (transactionService *TransactionService) List(transactions []*models.TransactionModel, userId uint, pagination *common.Pagination) (*common.Pagination, error) {
	result := transactionService.DB.Scopes(pagination.Paginate(), common.WhereUserIdScope(userId)).Preload("Category").Preload("Wallet").Find(&transactions)
	if result.Error != nil {
		return nil, errors.New("failed to fetch transactions")
	}

	pagination.Items = transactions
	return pagination, nil
}
