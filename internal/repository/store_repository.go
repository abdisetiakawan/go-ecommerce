package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/sirupsen/logrus"
)

type StoreRepository struct {
	Repository[entity.Store]
	Log *logrus.Logger
}

func NewStoreRepository(log *logrus.Logger) *StoreRepository {
	return &StoreRepository{
		Log: log,
	}
}