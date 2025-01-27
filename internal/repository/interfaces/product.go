package interfaces

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type ProductRepository interface {
	CreateProduct(product *entity.Product) error
	GetProducts(request *model.GetProductsRequest) ([]entity.Product, int64, error)
	GetProductById(userID uint, productUUID string) (*entity.Product, error)
	UpdateProduct(product *entity.Product) error
	DeleteProduct(product *entity.Product) error
	FindProductByUUID(productUUID string) (entity.Product, error)
	FindProductByID(productID uint) (entity.Product, error)
}