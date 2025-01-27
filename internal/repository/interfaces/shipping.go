package interfaces

import "github.com/abdisetiakawan/go-ecommerce/internal/entity"

type ShippingRepository interface {
	CreateShipping(shipping *entity.Shipping) error
	UpdateShipping(shipping *entity.Shipping) error
}
