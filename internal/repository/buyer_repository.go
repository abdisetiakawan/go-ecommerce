package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
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

func (r *BuyerRepository) FilterOrders(request *model.SearchOrderRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("user_id = ?", request.UserID)
		if status := request.Status; status != "" {
			db = db.Where("status = ?", status)
		}
		return db
	}
}