package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SellerRepository struct {
	Repository[entity.Store]
	Log *logrus.Logger
}

func NewSellerRepository(log *logrus.Logger) *SellerRepository {
	return &SellerRepository{
		Log: log,
	}
}

func (r *SellerRepository) HasStore(db *gorm.DB, userID uint) (bool, error) {
    var count int64
    err := db.Model(&entity.Store{}).Where("user_id = ?", userID).Count(&count).Error
    return count > 0, err
}