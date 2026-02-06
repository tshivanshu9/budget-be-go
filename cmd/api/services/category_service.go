package services

import (
	"errors"
	"strings"

	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/internal/custom_errors"
	"github.com/tshivanshu9/budget-be/internal/models"
	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

func (categoryService *CategoryService) List() ([]*models.CategoryModel, error) {
	var categories []*models.CategoryModel
	result := categoryService.db.Find(&categories)
	if result.Error != nil {
		return nil, errors.New("failed to fetch categories")
	}
	return categories, nil
}

func (categoryService *CategoryService) Create(data *requests.CreateCategoryRequest) (*models.CategoryModel, error) {
	slug := strings.ToLower(data.Name)
	slug = strings.ReplaceAll(slug, " ", "-")
	category := &models.CategoryModel{
		Name:     data.Name,
		Slug:     slug,
		IsCustom: data.IsCustome,
	}
	result := categoryService.db.Where("slug = ?", slug).FirstOrCreate(category)
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
	result := categoryService.db.Where("id = ?", id).First(&category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, custom_errors.NewNotFoundError("category not found")
		}
		return nil, errors.New("failed to fetch category")
	}
	return &category, nil
}

func (categoryService *CategoryService) DeleteById(id uint) error {
	result := categoryService.db.Where("id = ?", id).Delete(&models.CategoryModel{})
	if result.Error != nil {
		return errors.New("failed to delete category")
	}
	return nil
}
