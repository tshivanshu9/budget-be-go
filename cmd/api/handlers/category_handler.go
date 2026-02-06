package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/cmd/api/services"
	"github.com/tshivanshu9/budget-be/common"
)

func (h *Handler) ListCategoriesHandler(c *echo.Context) error {
	categoryService := services.NewCategoryService(h.DB)
	categories, err := categoryService.List()
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Categories list fetched successfully", categories)
}

func (h *Handler) CreateCategoryHandler(c *echo.Context) error {
	payload := new(requests.CreateCategoryRequest)
	err := h.BindRequestBody(c, payload)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	categoryService := services.NewCategoryService(h.DB)
	category, err := categoryService.Create(payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Category created successfully", category)
}

func (h *Handler) DeleteCategoryHandler(c *echo.Context) error {
	var categoryId requests.IDParamRequest
	err := (&echo.DefaultBinder{}).Bind(c, &categoryId)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	categoryService := services.NewCategoryService(h.DB)
	err = categoryService.DeleteById(categoryId.Id)

	if err != nil {
		return common.SendNotFoundResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Category deleted successfully", nil)
}
