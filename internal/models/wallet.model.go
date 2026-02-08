package models

type WalletModel struct {
	BaseModel
	UserId  uint      `json:"user_id" gorm:"not null;column:user_id;uniqueIndex:unique_userid_name"`
	Balance float64   `json:"balance" gorm:"not null;default:0.0;type:double precision"`
	Name    string    `json:"name" gorm:"not null;size:100;index;uniqueIndex:unique_userid_name"`
	Owner   UserModel `json:"user,omitempty" gorm:"foreignKey:UserId"`
}

func (w WalletModel) TableName() string {
	return "wallets"
}
