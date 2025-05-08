package interfaces

import "github.com/abdisetiakawan/go-ecommerce/internal/entity"

type ShippingRepository interface {
    UpdateShipping(shipping *entity.Shipping) error
}