package models

type CategoryModel struct {
	BaseModel
	Name     string `json:"name" gorm:"type:varchar(200);not null"`
	Slug     string `json:"slug" gorm:"type:varchar(200);not null;unique"`
	IsCustom bool   `json:"is_custom" gorm:"type:boolean;default:false"`
}

func (CategoryModel) TableName() string {
	return "categories"
}
