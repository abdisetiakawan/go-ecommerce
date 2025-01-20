package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductUUID string  `gorm:"type:char(36);uniqueIndex;not null"`
	StoreID     uint    `gorm:"not null"`
	Store       Store   `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductName string  `gorm:"size:255;not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null"`
	Category    string  `gorm:"type:enum('clothes', 'electronics', 'accessories');not null"`

	OrderItems []OrderItem `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
