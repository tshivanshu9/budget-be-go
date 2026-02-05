package services

import (
	"errors"
	"fmt"

	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (userService *UserService) RegisterUser(userRequest *requests.RegisterUserRequest) (*models.UserModel, error) {
	hashedPassword, err := common.HashPassword(userRequest.Password)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("user registration failed")
	}

	createdUser := models.UserModel{
		Email:     userRequest.Email,
		Password:  hashedPassword,
		FirstName: &userRequest.FirstName,
		LastName:  &userRequest.LastName,
	}

	result := userService.db.Create(&createdUser)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, errors.New("user registration failed")
	}
	return &createdUser, nil
}

func (userService *UserService) GetUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	result := userService.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (userService *UserService) ChangeUserPassword(newPassword string, userId uint) error {
	hashedPassword, err := common.HashPassword(newPassword)
	if err != nil {
		fmt.Println(err)
		return errors.New("Change password failed")
	}

	result := userService.db.Model(&models.UserModel{}).Where("id = ?", userId).Update("password", hashedPassword)
	if result.Error != nil {
		return errors.New("password change failed")
	}

	return nil
}
