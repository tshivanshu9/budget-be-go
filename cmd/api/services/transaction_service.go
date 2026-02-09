package services

import (
	"errors"
	"time"

	"github.com/tshivanshu9/budget-be/cmd/api/requests"
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
