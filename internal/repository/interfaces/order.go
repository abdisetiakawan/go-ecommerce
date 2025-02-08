package interfaces

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type OrderRepository interface {
	FindStoreByProductUUIDs(productUUIDs []string) (uint, error)
	UpdateOrder(order *entity.Order) error
	CreateOrder(order *entity.Order) error
	GetOrdersByBuyer(request *model.SearchOrderRequest) ([]entity.Order, int64, error)
	GetOrderByIdByBuyer(request *model.GetOrderDetails) (*entity.Order, error)
	GetOrdersBySeller(request *model.SearchOrderRequestBySeller) ([]entity.Order, int64, error)
	GetOrderBySeller(orderUUID string, storeID uint) (*entity.Order, error)
	ChangeOrderStatus() error
}
