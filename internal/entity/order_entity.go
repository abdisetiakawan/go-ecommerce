package entity

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderUUID string `gorm:"type:char(36);uniqueIndex;not null"`
	UserID 	  uint 	 `gorm:"not null"`
	User 	  User 	 `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status    string `gorm:"type:enum('pending', 'paid', 'shipped', 'completed', 'cancelled');default:'pending';not null"`
	TotalPrice float64 `gorm:"not null"`

	Items []OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payment *Payment `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Shipping *Shipping `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}