package interfaces

import "github.com/abdisetiakawan/go-ecommerce/internal/entity"

type PaymentRepository interface {
	CreatePayment(payment *entity.Payment) error
	UpdatePayment(payment *entity.Payment) error
}
