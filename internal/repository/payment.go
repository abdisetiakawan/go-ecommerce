package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(DB *gorm.DB) interfaces.PaymentRepository {
	return &PaymentRepository{DB}
}

func (r *PaymentRepository) CreatePayment(payment *entity.Payment) error {
	return r.DB.Create(payment).Error
}

func (r *PaymentRepository) UpdatePayment(payment *entity.Payment) error {
	return r.DB.Save(payment).Error
}