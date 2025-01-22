package repository

import (
	"fmt"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SellerRepository struct {
	StoreRepository *Repository[entity.Store]
	ProductRepository *Repository[entity.Product]
	OrderRepository *Repository[entity.Order]
	Log *logrus.Logger
}

func NewSellerRepository(log *logrus.Logger, db *gorm.DB) *SellerRepository {
	return &SellerRepository{
		StoreRepository: &Repository[entity.Store]{DB: db},
		ProductRepository: &Repository[entity.Product]{DB: db},
		OrderRepository: &Repository[entity.Order]{DB: db},
		Log: log,
	}
}

func (r *SellerRepository) HasStore(db *gorm.DB, userID uint) (bool, error) {
    var count int64
    err := db.Model(&entity.Store{}).Where("user_id = ?", userID).Count(&count).Error
    return count > 0, err
}

func (r *SellerRepository) CheckStore(db *gorm.DB, store *entity.Store, authID uint) error {
	return db.Where("user_id = ?", authID).Take(store).Error
}

func (r *SellerRepository) GetProducts(db *gorm.DB,	request *model.GetProductsRequest) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := db.Model(&entity.Product{}).
		Preload("Store", func(db *gorm.DB) *gorm.DB {
			return db.Where("user_id = ?", request.UserID)
		})

	if request.Search != "" {
		searchTerm := "%" + request.Search + "%"
		query = query.Where("product_name LIKE ? OR description LIKE ?", searchTerm, searchTerm)
	}
	if request.Category != "" {
		query = query.Where("category = ?", request.Category)
	}
	if request.PriceMin > 0 {
		query = query.Where("price >= ?", request.PriceMin)
	}
	if request.PriceMax > 0 {
		query = query.Where("price <= ?", request.PriceMax)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (request.Page - 1) * request.Limit
	if err := query.Limit(request.Limit).Offset(offset).Find(&products).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, model.ErrNotFound
		}
		return nil, 0, err
	}

	return products, total, nil
}

func (r *SellerRepository) GetProduct(db *gorm.DB, request *model.GetProductRequest) (*entity.Product, error) {
	var product entity.Product
	if err := db.Model(&entity.Product{}).Preload("Store", func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", request.UserID)
	}).Take(&product, "product_uuid = ?", request.ProductUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewApiError(fiber.StatusNotFound, fmt.Sprintf("Product with ID %s not found", request.ProductUUID), nil)
		}
		return nil, err
	}
	return &product, nil
}

func (r *SellerRepository) CheckProduct(db *gorm.DB, product *entity.Product, user_id uint, product_uuid string) error {
	return db.Model(&entity.Product{}).
		Preload("Store", func(db *gorm.DB) *gorm.DB {
			return db.Where("user_id = ?", user_id)
		}).
		Where(&entity.Product{
			ProductUUID: product_uuid,
		}).
		Take(product).
		Error
}

func (r *SellerRepository) GetOrder(db *gorm.DB, order_uuid string, store_id uint) (*entity.Order, error) {
	var order entity.Order
	if err := db.Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
		return db.Where("store_id = ?", store_id)
	}).
		Preload("Payment").
		Preload("Shipping").
		Where("order_uuid = ?", order_uuid).
		Take(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewApiError(fiber.StatusNotFound, fmt.Sprintf("Order with UUID %s not found", order_uuid), nil)
		}
		return nil, err
	}

	return &order, nil
}
