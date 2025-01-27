package interfaces

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type ShippingUseCase interface {
	UpdateShippingStatus(ctx context.Context, request *model.UpdateShippingStatusRequest) (*model.OrderResponse, error)
}