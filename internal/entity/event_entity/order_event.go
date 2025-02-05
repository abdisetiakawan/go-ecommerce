package evententity

import (
	"encoding/json"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"gorm.io/gorm"
)

type OrderEvent struct {
	gorm.Model
	EventUUID string `gorm:"type:char(36);uniqueIndex;not null"`
	OrderID   uint   `gorm:"not null"`
	Order	entity.Order `gorm:"foreignKey:OrderID"`
	EventType string `gorm:"type:enum('order_created', 'payment_processed', 'shipping_processed');not null"`
	Status       string         `gorm:"type:enum('pending','completed','failed');not null"`
	PaymentData json.RawMessage `gorm:"type:json"`
	ShippingData json.RawMessage `gorm:"type:json"`
	Error string `gorm:"type:text"`
}