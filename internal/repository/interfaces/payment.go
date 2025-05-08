package interfaces

import "github.com/abdisetiakawan/go-ecommerce/internal/entity"

type PaymentRepository interface {
	UpdatePayment(payment *entity.Payment) error
}
