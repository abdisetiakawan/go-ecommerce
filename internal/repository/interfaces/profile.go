package interfaces

import "github.com/abdisetiakawan/go-ecommerce/internal/entity"

type ProfileRepository interface {
	GetIDProfileByUserID(userID uint) (uint, error)
	CheckIDProfileByUserID(userID uint) (bool, error)
	GetProfileByUserID(userID uint) (*entity.Profile, error)
	CreateProfile(profile *entity.Profile) error
	UpdateProfile(profile *entity.Profile) error
}
