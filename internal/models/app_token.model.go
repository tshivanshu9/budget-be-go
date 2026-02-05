package models

import "time"

type AppTokenModel struct {
	BaseModel
	Token     string    `json:"-" gorm:"type:varchar(255);index"`
	TargetId  uint      `json:"-" gorm:"not null;index"`
	Type      string    `json:"-" gorm:"type:varchar(255);not null;index"`
	Used      bool      `json:"-" gorm:"type:boolean;not null;index"`
	ExpiresAt time.Time `json:"-" gorm:"index;not null;"`
}

func (AppTokenModel) TableName() string {
	return "app_tokens"
}
