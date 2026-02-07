package services

import (
	"errors"
	"strings"

	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/custom_errors"
	"github.com/tshivanshu9/budget-be/internal/models"
	"gorm.io/gorm"
)

type CategoryService struct {
	Db *gorm.DB
}

func NewCategoryService(Db *gorm.DB) *CategoryService {
	return &CategoryService{Db: Db}
}

func (categoryService *CategoryService) List(categories []*models.CategoryModel, pagination *common.Pagination) (*common.Pagination, error) {
	result := categoryService.Db.Scopes(pagination.Paginate()).Find(&categories)
	if result.Error != nil {
		return nil, errors.New("failed to fetch categories")
	}

	pagination.Items = categories
	return pagination, nil
}

func (categoryService *CategoryService) Create(data *requests.CreateCategoryRequest) (*models.CategoryModel, error) {
	slug := strings.ToLower(data.Name)
	slug = strings.ReplaceAll(slug, " ", "-")
	category := &models.CategoryModel{
		Name:     data.Name,
		Slug:     slug,
		IsCustom: data.IsCustom,
	}
	result := categoryService.Db.Where("slug = ?", slug).FirstOrCreate(category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return category, nil
		}
		return nil, errors.New("failed to create category")
	}
	return category, nil
}

func (categoryService *CategoryService) GetById(id uint) (*models.CategoryModel, error) {
	var category models.CategoryModel
	result := categoryService.Db.Where("id = ?", id).First(&category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, custom_errors.NewNotFoundError("category not found")
		}
		return nil, errors.New("failed to fetch category")
	}
	return &category, nil
}

func (categoryService *CategoryService) DeleteById(id uint) error {
	result := categoryService.Db.Where("id = ?", id).Delete(&models.CategoryModel{})
	if result.Error != nil {
		return errors.New("failed to delete category")
	}
	return nil
}

func (categoryService *CategoryService) GetCategoriesByIds(ids []uint) ([]*models.CategoryModel, error) {
	if len(ids) == 0 {
		return []*models.CategoryModel{}, nil
	}
	var categories []*models.CategoryModel
	result := categoryService.Db.Where("id IN ?", ids).Find(&categories)
	if result.Error != nil {
		return nil, errors.New("failed to fetch categories")
	}
	return categories, nil
}
