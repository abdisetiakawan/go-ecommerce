package entity

import "gorm.io/gorm"

type  OrderItem struct {
	gorm.Model
	OrderItemUUID string `gorm:"type:char(36);uniqueIndex;not null"`
	OrderID       uint   `gorm:"not null"`
	Order         Order  `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductID     uint   `gorm:"not null"`
	Product       Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quantity      int    `gorm:"not null"`
	TotalPrice    float64 `gorm:"not null"`
}