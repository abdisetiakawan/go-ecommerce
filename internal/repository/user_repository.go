package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.Profile]
	AuthRepository *Repository[entity.User]
	Logger *logrus.Logger
}

func NewUserRepository(log *logrus.Logger, db *gorm.DB) *UserRepository {
	return &UserRepository{
		AuthRepository: &Repository[entity.User]{DB: db},
		Logger: log,
	}
}

