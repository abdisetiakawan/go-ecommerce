package converter

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

func ProductToResponse(product *entity.Product) *model.ProductResponse {
	return &model.ProductResponse{
		StoreID:     product.StoreID,
		Store:       product.Store.StoreName,
		ProductUUID: product.ProductUUID,
		ProductName: product.ProductName,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ProductsToResponse(product *entity.Product) *model.ProductResponse {
	return &model.ProductResponse{
		StoreID:     product.StoreID,
		Store:       product.Store.StoreName,
		ProductUUID: product.ProductUUID,
		ProductName: product.ProductName,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		Store:       product.Store.StoreName,
	}
}