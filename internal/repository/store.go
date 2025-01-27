package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"gorm.io/gorm"
)

type StoreRepository struct {
	DB *gorm.DB
}

func NewStoreRepository(DB *gorm.DB) interfaces.StoreRepository {
	return &StoreRepository{DB}
}

func (r *StoreRepository) GetStoreIDByUserID(userID uint) (uint, error) {
	var storeID uint
	if err := r.DB.Model(&entity.Store{}).Select("id").Where("user_id = ?", userID).Find(&storeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, model.ErrNotFound
		}
		return 0, err
	}
	return storeID, nil
}

func (r *StoreRepository) FindStoreByUserID(userID uint) (entity.Store, error) {
	var store entity.Store
	if err := r.DB.Where("user_id = ?", userID).First(&store).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Store{}, model.ErrNotFound
		}
		return entity.Store{}, err
	}
	return store, nil
}

func (r *StoreRepository) HasStore(db *gorm.DB, userID uint) (bool, error) {
    var count int64
    err := db.Model(&entity.Store{}).Where("user_id = ?", userID).Count(&count).Error
    return count > 0, err
}

func (r *StoreRepository) CreateStore(store *entity.Store) error {
	return r.DB.Create(store).Error
}

func (r *StoreRepository) UpdateStore(store *entity.Store) error {
	return r.DB.Save(store).Error
}