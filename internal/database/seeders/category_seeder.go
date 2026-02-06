package main

import (
	"fmt"
	"strings"

	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/models"
)

func main() {
	db, err := common.NewMysql()
	if err != nil {
		panic("Failed to connect to database")
	}

	categoryNames := []string{
		"Food & Dining",
		"Transportation",
		"Utilities",
		"Entertainment",
		"Healthcare",
		"Personal Care",
		"Education",
		"Shopping",
		"Travel",
		"Savings & Investments",
	}

	categories := make([]models.CategoryModel, len(categoryNames))
	for i, name := range categoryNames {
		categories[i] = models.CategoryModel{
			Name:     name,
			Slug:     strings.ToLower(strings.ReplaceAll(name, " ", "-")),
			IsCustom: false,
		}
	}

	result := db.Create(categories)
	if result.Error != nil {
		fmt.Println("Failed to seed categories:", result.Error)
		panic("Failed to seed categories")
	}
	fmt.Println("Categories seeded successfully")
}
