package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BuyerRepository struct {
	OrderRepository *Repository[entity.Order]
	ProductRepository *Repository[entity.Product]
	PaymentRepository *Repository[entity.Payment]
	ShippingRepository *Repository[entity.Shipping]
	Log *logrus.Logger
}

func NewBuyerRepository(log *logrus.Logger, db *gorm.DB) *BuyerRepository {
	return &BuyerRepository{
		OrderRepository: &Repository[entity.Order]{DB: db},
		ProductRepository: &Repository[entity.Product]{DB: db},
		PaymentRepository: &Repository[entity.Payment]{DB: db},
		ShippingRepository: &Repository[entity.Shipping]{DB: db},
		Log: log,
	}
}