package interfaces

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"gorm.io/gorm"
)

type StoreRepository interface {
	GetStoreIDByUserID(userID uint) (uint, error)
	FindStoreByUserID(userID uint) (entity.Store, error)
	HasStore(db *gorm.DB, userID uint) (bool, error)
	CreateStore(store *entity.Store) error
	UpdateStore(store *entity.Store) error
}
