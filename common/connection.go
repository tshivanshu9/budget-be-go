package common

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysql()(*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	database := os.Getenv("DB_DATABASE")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	fmt.Println(dsn)

	if err != nil {
		return nil, err
	}

	fmt.Println("Database connection successful!")

	return db, nil
}