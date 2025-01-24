package interfaces

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type SellerUseCase interface {
	Create(ctx context.Context, request *model.RegisterStore) (*model.StoreResponse, error)
	CreateProduct(ctx context.Context, request *model.RegisterProduct) (*model.ProductResponse, error)
	GetStore(ctx context.Context, id uint) (*model.StoreResponse, error)
	Update(ctx context.Context, request *model.UpdateStore) (*model.StoreResponse, error)
	GetProducts(ctx context.Context, request *model.GetProductsRequest) ([]model.ProductResponse, int64, error)
	GetProduct(ctx context.Context, request *model.GetProductRequest) (*model.ProductResponse, error)
	UpdateProduct(ctx context.Context, request *model.UpdateProduct) (*model.ProductResponse, error)
	DeleteProduct(ctx context.Context, request *model.DeleteProductRequest) error
	GetOrder(ctx context.Context, request *model.GetOrderDetails) (*model.OrderResponse, error)
	GetOrders(ctx context.Context, request *model.SearchOrderRequestBySeller) ([]model.OrdersResponseForSeller, int64, error)
	UpdateShippingStatus(ctx context.Context, request *model.UpdateShippingStatusRequest) (*model.OrderResponse, error)
}