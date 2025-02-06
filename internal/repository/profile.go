package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	DB *gorm.DB
}

func NewProfileRepository(DB *gorm.DB) interfaces.ProfileRepository {
	return &ProfileRepository{DB}
}

func (r *ProfileRepository) GetIDProfileByUserID(userID uint) (uint, error) {
	var profileID uint
	if err := r.DB.Model(&entity.Profile{}).Select("id").Where("user_id = ?", userID).Find(&profileID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, model.ErrNotFound
		}
		return 0, err
	}
	return profileID, nil
}

func (r *ProfileRepository) CheckIDProfileByUserID(userID uint) (bool, error) {
	var count int64
	if err := r.DB.Model(&entity.Profile{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, model.ErrNotFound
		}
		return false, err
	}
	return count > 0, nil
}

func (r *ProfileRepository) GetProfileByUserID(userID uint) (*entity.Profile, error) {
	var profile entity.Profile
	if err := r.DB.Where("user_id = ?", userID).Find(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) CreateProfile(entity *entity.Profile) error {
	return r.DB.Create(entity).Error
}

func (r *ProfileRepository) UpdateProfile(entity *entity.Profile) error {
	if entity.ID == 0 {
		return model.ErrNotFound
	}

	if err := r.DB.Model(&entity).Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return err
	}
	return nil
}