package services

import (
	"github.com/tshivanshu9/budget-be/cmd/api/requests"
	"github.com/tshivanshu9/budget-be/internal/models"
	"gorm.io/gorm"
)

type WalletService struct {
	DB *gorm.DB
}

func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{DB: db}
}

func (w *WalletService) Create(data *requests.CreateWalletRequest, userId uint) (*models.WalletModel, error) {
	wallet := &models.WalletModel{
		UserId:  userId,
		Balance: data.Amount,
		Name:    data.Name,
	}
	result := w.DB.Create(wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return wallet, nil
}

func (w *WalletService) WalletExistsForNameAndUserId(name string, userId uint) (*models.WalletModel, error) {
	var wallet models.WalletModel
	result := w.DB.Where("user_id = ? AND name = ?", userId, name).First(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (w *WalletService) GenerateDefaultWallets(id uint) ([]*models.WalletModel, error) {
	walletNames := []string{"Cash", "Bank"}
	wallets := make([]*models.WalletModel, 0)
	for _, name := range walletNames {
		walletExists, err := w.WalletExistsForNameAndUserId(name, id)
		if err != nil {
			return nil, err
		}
		if walletExists != nil {
			wallets = append(wallets, walletExists)
			continue
		}
		walletReq := requests.CreateWalletRequest{
			Name:   name,
			Amount: 0,
		}
		walletExists, err = w.Create(&walletReq, id)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, walletExists)
	}
	return wallets, nil
}

func (w *WalletService) ListWalletsForUser(userId uint) ([]*models.WalletModel, error) {
	var wallets []*models.WalletModel
	result := w.DB.Where("user_id = ?", userId).Find(&wallets)
	if result.Error != nil {
		return nil, result.Error
	}
	return wallets, nil
}

func (w *WalletService) GetWalletByIdAndUserId(id uint, userId uint) (*models.WalletModel, error) {
	var wallet models.WalletModel
	result := w.DB.Where("id = ? AND user_id = ?", id, userId).First(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (w *WalletService) IncrementWalletBalance(wallet *models.WalletModel, amount float64) error {
	result := w.DB.Model(wallet).Update("balance", gorm.Expr("balance + ?", amount))
	return result.Error
}

func (w *WalletService) DecrementWalletBalance(wallet *models.WalletModel, amount float64) error {
	result := w.DB.Model(wallet).Update("balance", gorm.Expr("balance - ?", amount))
	return result.Error
}
