package main

import (
	"github.com/tshivanshu9/budget-be-go/common"
	"github.com/tshivanshu9/budget-be-go/internal/models"
)

func main() {
	db, err := common.NewMysql()
	if err != nil {
		panic("Failed to connect to database!")
	}

	err = db.AutoMigrate(&models.UserModel{})
	if err != nil {
		panic(err)
	}
}
