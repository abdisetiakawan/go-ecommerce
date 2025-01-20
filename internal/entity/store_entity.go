package entity

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StoreName   string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Products []Product `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
