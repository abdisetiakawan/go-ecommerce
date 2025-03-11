package repository

import (
	"context"
	"encoding/json"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	eventmodel "github.com/abdisetiakawan/go-ecommerce/internal/model/event_model"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderRepository struct {
    kafka *helper.KafkaConsumer
	DB *gorm.DB
}

func NewOrderRepository(DB *gorm.DB, kafka *helper.KafkaConsumer) interfaces.OrderRepository {
	return &OrderRepository{kafka, DB}
}

func (r *OrderRepository) FindStoreByProductUUIDs(productUUIDs []string) (uint, error) {
    var storeID uint
    rows, err := r.DB.Model(&entity.Product{}).
        Select("store_id").
        Where("product_uuid IN ?", productUUIDs).
        Group("store_id").
        Rows()
    if err != nil {
        return 0, err
    }
    defer rows.Close()

    count := 0
    for rows.Next() {
        count++
        if count > 1 {
            return 0, model.ErrBadRequest
        }
        rows.Scan(&storeID)
    }

    if count == 0 {
        return 0, gorm.ErrRecordNotFound
    }

    return storeID, nil
}


func (r *OrderRepository) UpdateOrder(order *entity.Order) error {
	return r.DB.Save(order).Error
}

func (r *OrderRepository) CreateOrder(order *entity.Order) error {
	return r.DB.Create(order).Error
}

func (r *OrderRepository) GetOrdersByBuyer(request *model.SearchOrderRequest) ([]entity.Order, int64, error) {
    filteredQuery := r.DB.Model(&entity.Order{}).Scopes(r.FilterOrders(request))
    
    var total int64
    if err := filteredQuery.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    var orders []entity.Order
    if err := filteredQuery.Offset((request.Page - 1) * request.Limit).Limit(request.Limit).Find(&orders).Error; err != nil {
        return nil, 0, err
    }
    
    return orders, total, nil
}

func (r *OrderRepository) FilterOrders(request *model.SearchOrderRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("user_id = ?", request.UserID)
		if status := request.Status; status != "" {
			db = db.Where("status = ?", status)
		}
		return db
	}
}

func (r *OrderRepository) GetOrderByIdByBuyer(request *model.GetOrderDetails) (*entity.Order, error) {
	var order entity.Order
	if err := r.DB.Preload("Items", func (db *gorm.DB) *gorm.DB {
		return db.Order("order_items.created_at ASC")
	}).
    Preload("Items.Product").
	Preload("Payment").
	Preload("Shipping").
	Where(&entity.Order{
		OrderUUID: request.OrderUUID,
		UserID:    request.UserID,
	}).
	Take(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

    for i := range order.Items {
        order.Items[i].ProductName = order.Items[i].Product.ProductName
    }
    
	return &order, nil
}

func (r *OrderRepository) GetOrdersBySeller(request *model.SearchOrderRequestBySeller) ([]entity.Order, int64, error) {
    var orders []entity.Order
    var total int64

    subquery := r.DB.Model(&entity.OrderItem{}).
        Select("DISTINCT order_id").
        Joins("JOIN products ON order_items.product_id = products.id").
        Where("products.store_id = ?", request.StoreID)

    query := r.DB.Model(&entity.Order{}).
        Preload("Items.Product").
        Preload("Payment").
        Preload("Shipping").
        Where("id IN (?)", subquery)

    if request.Status != "" {
        query = query.Where("status = ?", request.Status)
    }

    if err := query.Count(&total).Find(&orders).Error; err != nil {
        return nil, 0, err
    }

    return orders, total, nil
}

func (r *OrderRepository) GetOrderBySeller(order_uuid string, store_id uint) (*entity.Order, error) {
    var order entity.Order

    if err := r.DB.Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
        return r.DB.Where("store_id = ?", store_id)
    }).
        Preload("Payment").
        Preload("Shipping").
        Where("order_uuid = ?", order_uuid).
        Joins("JOIN order_items ON order_items.order_id = orders.id").
        Joins("JOIN products ON products.id = order_items.product_id").
        Where("products.store_id = ?", store_id).
        Take(&order).Error; err != nil {

        if err == gorm.ErrRecordNotFound {
            return nil, model.ErrNotFound
        }
        return nil, err
    }

    return &order, nil
}

func (r *OrderRepository) ChangeOrderStatus() error {
    consumer, err := r.kafka.Consume(context.Background(), "change_order_topic")
    if err != nil {
        return err
    }
    defer consumer.Close()

    for {
        select {
        case msg := <-consumer.Messages():
            var orderMessage eventmodel.OrderMessage
            err := json.Unmarshal(msg.Value, &orderMessage)
            if err != nil {
                logrus.WithError(err).Error("Failed to unmarshal order message")
                continue
            }
            if err := r.DB.Model(&entity.Order{}).
                Where("id = ?", orderMessage.OrderID).
                Update("status", orderMessage.Status).Error; err != nil {
                logrus.WithError(err).Error("Failed to update order status")
            }

        case err := <-consumer.Errors():
            logrus.WithError(err).Error("Failed to consume order topic")
        }
    }
}