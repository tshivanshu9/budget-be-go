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

func (budgetService *BudgetService) GetById(id uint) (*models.BudgetModel, error) {
	var budget models.BudgetModel
	result := budgetService.DB.Scopes(common.WhereIdScope(id)).First(&budget)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, custom_errors.NewNotFoundError("budget not found")
		}
		return nil, errors.New("failed to fetch budget")
	}
	return &budget, nil
}

func (budgetService *BudgetService) countForYearMonthSlugUserIdExcludeBudgetId(userId uint, slug string, year uint16, month uint8, budgetId uint) (int64, error) {
	var count int64
	result := budgetService.DB.Model(&models.BudgetModel{}).Where("user_id = ? AND slug = ? AND year = ? AND month = ? AND id <> ?", userId, slug, year, month, budgetId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (budgetService *BudgetService) Update(budget *models.BudgetModel, payload *requests.UpdateBudgetRequest) error {
	if payload.Date != "" {
		timeParsed, err := time.Parse(time.DateOnly, payload.Date)
		if err != nil {
			return errors.New("invalid date format, expected YYYY-MM-DD")
		}
		budget.Date = timeParsed
		budget.Month = uint8(timeParsed.Month())
		budget.Year = uint16(timeParsed.Year())
	}
	if payload.Amount > 0 {
		budget.Amount = payload.Amount
	}
	if payload.Description != nil {
		budget.Description = payload.Description
	}
	if payload.Title != "" {
		budget.Title = payload.Title
		slug := strings.ToLower(payload.Title)
		slug = strings.ReplaceAll(slug, " ", "-")
		budget.Slug = slug
	}

	count, err := budgetService.countForYearMonthSlugUserIdExcludeBudgetId(budget.UserId, budget.Slug, budget.Year, budget.Month, budget.ID)
	if err != nil {
		return errors.New("failed to update budget")
	}
	if count > 0 {
		return errors.New("another budget with same title exists for the month, please choose different title or month")
	}

	result := budgetService.DB.Model(&budget).Updates(budget)
	if result.Error != nil {
		return errors.New("failed to update budget")
	}
	return nil
}
