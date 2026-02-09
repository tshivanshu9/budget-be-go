package services

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/tshivanshu9/budget-be/internal/models"
	"gorm.io/gorm"
)

type AppTokenService struct {
	db *gorm.DB
}

func NewAppTokenService(db *gorm.DB) *AppTokenService {
	return &AppTokenService{db: db}
}

func (appTokenService *AppTokenService) getToken() int {
	rand.Seed(time.Now().UnixNano())
	min := 10000
	max := 99999
	return rand.Intn(max-min+1) + min
}

func (appTokenService *AppTokenService) GenerateResetPasswordToken(user *models.UserModel) (*models.AppTokenModel, error) {
	tokenCreated := models.AppTokenModel{
		TargetId:  user.ID,
		Type:      "reset_password",
		Token:     strconv.Itoa(appTokenService.getToken()),
		Used:      false,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	result := appTokenService.db.Create(&tokenCreated)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tokenCreated, nil
}

func (appTokenService *AppTokenService) ValidateResetPasswordToken(user *models.UserModel, token string) (*models.AppTokenModel, error) {
	var retrievedToken models.AppTokenModel
	result := appTokenService.db.Where("target_id = ? AND type = ? and token = ?", user.ID, "reset_password", token).First(&retrievedToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid password reset token")
		}
		return nil, result.Error
	}

	if retrievedToken.Used {
		return nil, errors.New("invalid password reset token")
	}

	if retrievedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("password reset token has expired")
	}

	return &retrievedToken, nil
}

func (appTokenService *AppTokenService) InvalidateToken(userId uint, token *models.AppTokenModel) {
	appTokenService.db.Model(&models.AppTokenModel{}).Where("id = ? AND target_id = ?", token.ID, userId).Update("used", true)
}
