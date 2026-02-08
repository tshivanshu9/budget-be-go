package main

import (
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/models"
)

func main() {
	db, err := common.NewMysql()
	if err != nil {
		panic("Failed to connect to database!")
	}

	err = db.AutoMigrate(&models.UserModel{}, &models.AppTokenModel{}, &models.CategoryModel{}, &models.BudgetModel{}, &models.WalletModel{})
	if err != nil {
		panic(err)
	}
}
