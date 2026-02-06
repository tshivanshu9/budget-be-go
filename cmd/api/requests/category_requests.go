package requests

type CreateCategoryRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=200"`
	IsCustome bool   `default:"true" json:"is_custom"`
}
