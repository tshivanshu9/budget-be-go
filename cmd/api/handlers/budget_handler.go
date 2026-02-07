package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/cmd/api/services"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/custom_errors"
	"github.com/tshivanshu9/budget-be/internal/models"
)

func (h *Handler) CreateBudgetHandler(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}

	payload := new(requests.StoreBudgetRequest)
	err := h.BindRequestBody(c, payload)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	budgetService := services.NewBudgetService(h.DB)
	categoryService := services.NewCategoryService(h.DB)

	budget, err := budgetService.Create(payload, user.ID)
	if err != nil {
		fmt.Println(err)
		return common.SendInternalServerErrorResponse(c, "Budget creation failed, try again later")
	}

	if budget.ID == 0 {
		return common.SendInternalServerErrorResponse(c, "Budget creation failed")
	}

	categories, err := categoryService.GetCategoriesByIds(payload.Categories)
	if err != nil {
		fmt.Println(err)
		return common.SendInternalServerErrorResponse(c, "Budget creation failed, try again later")
	}

	err = budgetService.DB.Model(budget).Association("Categories").Replace(categories)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Budget creation failed, try again later")
	}

	budget.Categories = categories

	return common.SendSuccessResponse(c, "budget created successfully", budget)
}

func (h *Handler) ListBudgetsHandler(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}

	budgetService := services.NewBudgetService(h.DB)
	var budgets []*models.BudgetModel
	paginator := common.NewPaginator(budgets, c.Request(), h.DB)
	paginatedBudgets, err := budgetService.List(user.ID, paginator)

	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Budgets list fetched successfully", paginatedBudgets)
}

func (h *Handler) UpdateBudgetHandler(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}

	budgetIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid budget id")
	}
	budgetId := uint(budgetIdInt)

	payload := new(requests.UpdateBudgetRequest)
	err = h.BindRequestBody(c, payload)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	budgetService := services.NewBudgetService(h.DB)

	budget, err := budgetService.GetById(budgetId)
	if err != nil {
		var notFoundErr *custom_errors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return common.SendNotFoundResponse(c, "Budget not found")
		}
		return common.SendInternalServerErrorResponse(c, "Failed to fetch budget")
	}

	if budget.UserId != user.ID {
		return common.SendBadRequestResponse(c, "Invalid budget id")
	}

	err = budgetService.Update(budget, payload)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, errors.New("another budget with same title exists for the month, please choose different title or month")) {
			return common.SendBadRequestResponse(c, err.Error())
		}
		return common.SendInternalServerErrorResponse(c, "Failed to update budget")
	}

	if payload.Categories != nil {
		categoryService := services.NewCategoryService(h.DB)
		categories, err := categoryService.GetCategoriesByIds(payload.Categories)
		if err != nil {
			fmt.Println(err)
			return common.SendInternalServerErrorResponse(c, "Failed to update budget")
		}

		err = budgetService.DB.Model(budget).Association("Categories").Replace(categories)
		if err != nil {
			fmt.Println(err)
			return common.SendInternalServerErrorResponse(c, "Failed to update budget")
		}
		budget.Categories = categories
	}

	return common.SendSuccessResponse(c, "Budget updated successfully", budget)
}

func (h *Handler) DeleteBudgetHandler(c *echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, nil)
	}

	budgetIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid budget id")
	}
	budgetId := uint(budgetIdInt)

	budgetService := services.NewBudgetService(h.DB)

	budget, err := budgetService.GetById(budgetId)
	if budget.UserId != user.ID {
		return common.SendBadRequestResponse(c, "Invalid budget id")
	}
	if err != nil {
		var notFoundErr *custom_errors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return common.SendNotFoundResponse(c, "Budget not found")
		}
		return common.SendInternalServerErrorResponse(c, "Failed to fetch budget")
	}

	err = budgetService.Delete(budget)
	if err != nil {
		fmt.Println(err)
		return common.SendInternalServerErrorResponse(c, "Failed to delete budget")
	}

	return common.SendSuccessResponse(c, "Budget deleted successfully", nil)
}
