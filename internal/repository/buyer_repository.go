package repository

import (
	"fmt"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BuyerRepository struct {
	OrderRepository *Repository[entity.Order]
	ProductRepository *Repository[entity.Product]
	PaymentRepository *Repository[entity.Payment]
	ShippingRepository *Repository[entity.Shipping]
	Log *logrus.Logger
}

func NewBuyerRepository(log *logrus.Logger, db *gorm.DB) *BuyerRepository {
	return &BuyerRepository{
		OrderRepository: &Repository[entity.Order]{DB: db},
		ProductRepository: &Repository[entity.Product]{DB: db},
		PaymentRepository: &Repository[entity.Payment]{DB: db},
		ShippingRepository: &Repository[entity.Shipping]{DB: db},
		Log: log,
	}
}

func (r *BuyerRepository) GetOrders(db *gorm.DB, request *model.SearchOrderRequest) ([]entity.Order, int64, error) {
	filteredQuery := db.Model(&entity.Order{}).Scopes(r.FilterOrders(request))

	var orders []entity.Order
	if err := filteredQuery.Offset((request.Page - 1) * request.Limit).Limit(request.Limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := filteredQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *BuyerRepository) GetOrder(db *gorm.DB, request *model.GetOrderDetails) (*entity.Order, error) {
	var order entity.Order
	if err := db.Preload("Items", func (db *gorm.DB) *gorm.DB {
		return db.Order("order_items.created_at ASC")
	}).
	Preload("Payment").
	Preload("Shipping").
	Where(&entity.Order{
		OrderUUID: request.OrderUUID,
		UserID:    request.UserID,
	}).
	Take(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewApiError(fiber.StatusNotFound, fmt.Sprintf("Order with ID %s not found", request.OrderUUID), nil)
		}
		return nil, err
	}
	return &order, nil
}

func (r *BuyerRepository) FilterOrders(request *model.SearchOrderRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("user_id = ?", request.UserID)
		if status := request.Status; status != "" {
			db = db.Where("status = ?", status)
		}
		return db
	}
}