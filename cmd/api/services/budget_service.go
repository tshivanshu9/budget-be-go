package services

import (
	"errors"
	"strings"
	"time"

	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/custom_errors"
	"github.com/tshivanshu9/budget-be/internal/models"
	"gorm.io/gorm"
)

type BudgetService struct {
	DB *gorm.DB
}

func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{DB: db}
}

func (budgetService *BudgetService) Create(payload *requests.StoreBudgetRequest, userId uint) (*models.BudgetModel, error) {
	slug := strings.ToLower(payload.Title)
	slug = strings.ReplaceAll(slug, " ", "-")

	model := &models.BudgetModel{
		Title:       payload.Title,
		Amount:      payload.Amount,
		UserId:      userId,
		Slug:        slug,
		Description: payload.Description,
	}

	if payload.Date == "" {
		currentDate := time.Now()
		model.Date = currentDate
		model.Month = uint8(currentDate.Month())
		model.Year = uint16(currentDate.Year())
	} else {
		parsedDate, err := time.Parse("2006-01-02", payload.Date)
		if err != nil {
			return nil, err
		}
		model.Date = parsedDate
		model.Month = uint8(parsedDate.Month())
		model.Year = uint16(parsedDate.Year())
	}
	retrievedBudget, err := budgetService.budgetExistsForYearMonthSlugUserid(userId, model.Slug, model.Year, model.Month)
	if err != nil {
		var notFoundErr *custom_errors.NotFoundError
		if errors.As(err, &notFoundErr) {
			result := budgetService.DB.Create(model)
			if result.Error != nil {
				return nil, errors.New("failed to create budget")
			}
			return model, nil
		}
		return nil, err
	}
	return retrievedBudget, nil
}

func (budgetService *BudgetService) budgetExistsForYearMonthSlugUserid(userId uint, slug string, year uint16, month uint8) (*models.BudgetModel, error) {
	var budget models.BudgetModel
	result := budgetService.DB.Where("user_id = ? AND slug = ? AND year = ? AND month = ?", userId, slug, year, month).First(&budget)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, custom_errors.NewNotFoundError("budget not found")
		}
		return nil, result.Error
	}
	return &budget, nil
}

func (budgetService *BudgetService) List(userId uint, pagination *common.Pagination) (*common.Pagination, error) {
	var budgets []*models.BudgetModel
	result := budgetService.DB.
		Scopes(common.WhereUserIdScope(userId)).
		Scopes(pagination.Paginate()).
		Preload("Categories").Find(&budgets)
	if result.Error != nil {
		return nil, errors.New("failed to fetch budgets")
	}
	pagination.Items = budgets
	return pagination, nil
}
