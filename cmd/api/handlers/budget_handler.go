package handlers

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/cmd/api/services"
	"github.com/tshivanshu9/budget-be/common"
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

	categories, err := categoryService.GetCategoriesByIds(payload.Categories)
	if err != nil {
		fmt.Println(err)
		return common.SendInternalServerErrorResponse(c, "Budget creation failed, try again later")
	}

	err = budgetService.DB.Model(budget).Association("Categories").Replace(categories)
	if err != nil {
		fmt.Println(err)
		return common.SendInternalServerErrorResponse(c, "Budget creation failed, try again later")
	}

	budget.Categories = categories

	return common.SendSuccessResponse(c, "budget created successfully", budget)
}
