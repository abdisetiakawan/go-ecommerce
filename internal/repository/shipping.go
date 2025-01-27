package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"gorm.io/gorm"
)

type ShippingRepository struct {
	DB *gorm.DB
}

func NewShippingRepository(DB *gorm.DB) *ShippingRepository {
	return &ShippingRepository{DB}
}

func (r *ShippingRepository) CreateShipping(shipping *entity.Shipping) error {
	return r.DB.Create(shipping).Error
}

func (r *ShippingRepository) UpdateShipping(shipping *entity.Shipping) error {
	return r.DB.Save(shipping).Error
}