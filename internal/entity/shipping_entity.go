package entity

import "gorm.io/gorm"

type Shipping struct {
	gorm.Model
	ShippingUUID string `gorm:"type:char(36);uniqueIndex;not null"`
	OrderID      uint   `gorm:"not null"`
	Order        Order  `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Address      string `gorm:"size:255;not null"`
	City         string `gorm:"size:255;not null"`
	Province     string `gorm:"size:255;not null"`
	PostalCode   string `gorm:"size:10;not null"`
	Status       string `gorm:"type:enum('pending', 'shipped', 'delivered', 'cancelled');default:'pending';not null"`
}