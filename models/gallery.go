package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Gallery struct {
	ID          int            `gorm:"primaryKey" json:"id"`
	UserID      int            `json:"user_id"`
	Photo       *string        `gorm:"size:255" json:"image"`
	Title       string         `gorm:"size:50" json:"title"`
	Description *string        `gorm:"size:255" json:"description"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func CreateNewGallery(gallery *Gallery) error {
	return DB().Create(gallery).Error
}

func GetAllGallery(userId int, limit, offset int) []Gallery {
	var galleries []Gallery
	DB().Model(&Gallery{}).Where("user_id = ?", userId).Order(clause.OrderByColumn{
		Column: clause.Column{Name: "created_at"},
		Desc:   true,
	}).Limit(limit).Offset(offset).Find(&galleries)
	return galleries
}

func GetGalleryByID(id int) (Gallery, error) {
	var gallery Gallery
	err := DB().Model(&Gallery{}).Where("id = ?", id).First(&gallery).Error
	return gallery, err
}
