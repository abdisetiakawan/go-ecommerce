package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/sirupsen/logrus"
)

type AuthRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewAuthRepository(log *logrus.Logger) *AuthRepository {
	return &AuthRepository{
		Log: log,
	}
}