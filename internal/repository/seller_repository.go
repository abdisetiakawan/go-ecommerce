package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
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