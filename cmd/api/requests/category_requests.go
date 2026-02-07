package requests

type CreateCategoryRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=200"`
	IsCustom bool   `default:"true" json:"is_custom"`
}
