package models

import "time"

type TransactionModel struct {
	BaseModel
	ParentId    *uint          `json:"-" gorm:"omitnull;column:parent_id"`
	Title       *string        `json:"title" gorm:"type:varchar(200)"`
	Description *string        `json:"description" gorm:"type:varchar(500)"`
	UserId      uint           `json:"user_id" gorm:"not null;column:user_id"`
	CategoryId  *uint          `json:"category_id" gorm:"column:category_id"`
	WalletId    uint           `json:"wallet_id" gorm:"not null;column:wallet_id"`
	Amount      float64        `json:"amount" gorm:"not null;column:amount;type:double precision"`
	Date        time.Time      `json:"date" gorm:"not null;type:datetime"`
	Month       uint8          `json:"month" gorm:"not null;column:month;type:TINYINT UNSIGNED"`
	Year        uint16         `json:"year" gorm:"not null;column:year;type:INT UNSIGNED;index:idx_month_year"`
	Type        string         `json:"type" gorm:"not null;column:type;type:varchar(100);index:idx_month_year"`
	IsReversal  bool           `json:"is_reversal" gorm:"not null;column:is_reversal;default:false"`
	Category    *CategoryModel `json:"category" gorm:"foreignKey:CategoryId;constraint:OnDelete:CASCADE"`
	Wallet      *WalletModel   `json:"wallet" gorm:"foreignKey:WalletId;constraint:OnDelete:CASCADE"`
	User        *UserModel     `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}
