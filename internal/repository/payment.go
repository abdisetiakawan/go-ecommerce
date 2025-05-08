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
	return &PaymentRepository{DB: DB}
}

func (r *PaymentRepository) UpdatePayment(payment *entity.Payment) error {
	return r.DB.Save(payment).Error
}

