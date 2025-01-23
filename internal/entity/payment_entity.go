package entity

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	PaymentUUID string  `gorm:"type:char(36);uniqueIndex;not null"`
	OrderID     uint    `gorm:"not null"`
	Order       Order   `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount  	float64 `gorm:"not null"`
	Status      string  `gorm:"type:enum('pending', 'paid', 'cancelled');default:'pending';not null"`
	Method      string  `gorm:"type:enum('cash', 'transfer');not null"`
}