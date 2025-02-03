package interfaces

import "github.com/abdisetiakawan/go-ecommerce/internal/entity"

type ShippingRepository interface {
	CreateShipping() error
	UpdateShipping(shipping *entity.Shipping) error
}
