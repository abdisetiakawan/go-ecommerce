package interfaces

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, request *model.RegisterProduct) (*model.ProductResponse, error) 
	GetProducts(ctx context.Context, request *model.GetProductsRequest) ([]model.ProductResponse, int64, error)
	GetProductById(ctx context.Context, request *model.GetProductRequest) (*model.ProductResponse, error)
	UpdateProduct(ctx context.Context, request *model.UpdateProduct) (*model.ProductResponse, error)
	DeleteProduct(ctx context.Context, request *model.DeleteProductRequest) error
}