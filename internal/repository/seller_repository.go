package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SellerRepository struct {
	StoreRepository *Repository[entity.Store]
	ProductRepository *Repository[entity.Product]
	Log *logrus.Logger
}

func NewSellerRepository(log *logrus.Logger, db *gorm.DB) *SellerRepository {
	return &SellerRepository{
		StoreRepository: &Repository[entity.Store]{DB: db},
		ProductRepository: &Repository[entity.Product]{DB: db},
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
