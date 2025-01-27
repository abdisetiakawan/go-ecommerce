package interfaces

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, input *model.CreateOrder) (*model.OrderResponse, error)
	GetOrders(ctx context.Context, request *model.SearchOrderRequest) ([]model.ListOrderResponse, int64, error)
	GetOrder(ctx context.Context, request *model.GetOrderDetails) (*model.OrderResponse, error)
	CancelOrder(ctx context.Context, request *model.CancelOrderRequest) (*model.OrderResponse, error)
	CheckoutOrder(ctx context.Context, request *model.CheckoutOrderRequest) (*model.OrderResponse, error)
	GetOrdersBySeller(ctx context.Context, request *model.SearchOrderRequestBySeller) ([]model.OrdersResponseForSeller, int64, error)
}