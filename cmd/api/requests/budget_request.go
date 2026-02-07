package requests

type StoreBudgetRequest struct {
	Categories  []uint  `json:"categories" validate:"required,min=1,dive"`
	Amount      float64 `json:"amount" validate:"required,numeric,min=1"`
	Date        string  `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Title       string  `json:"title" validate:"required,min=2,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=2,max=500"`
}

type UpdateBudgetRequest struct {
	Categories  []uint  `json:"categories,omitempty" validate:"omitempty,min=1,dive"`
	Amount      float64 `json:"amount,omitempty" validate:"omitempty,numeric,min=1"`
	Date        string  `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Title       string  `json:"title,omitempty" validate:"omitempty,min=2,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=2,max=500"`
}
