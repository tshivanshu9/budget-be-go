package models

import (
	"time"
)

type BudgetModel struct {
	BaseModel
	Title       string           `json:"title" gorm:"type:varchar(255);not null"`
	Description *string          `json:"description" gorm:"type:varchar(500)"`
	Slug        string           `json:"slug" gorm:"index;type:varchar(255);not null;uniqueIndex:unique_user_id_slug_year_month"`
	UserId      uint             `json:"user_id" gorm:"not null;column:user_id;uniqueIndex:unique_user_id_slug_year_month"`
	Amount      float64          `json:"amount" gorm:"type:decimal(10,2);not null"`
	Categories  []*CategoryModel `json:"categories" gorm:"many2many:budget_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Date        time.Time        `json:"date" gorm:"type:date;not null"`
	Month       uint8            `json:"month" gorm:"type:TINYINT UNSIGNED;not null;index:idx_month_year;uniqueIndex:unique_user_id_slug_year_month"`
	Year        uint16           `json:"year" gorm:"type:INT UNSIGNED;not null;index:idx_month_year;uniqueIndex:unique_user_id_slug_year_month"`
}

func (BudgetModel) TableName() string {
	return "budgets"
}
