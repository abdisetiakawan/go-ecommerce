package usecase

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StoreUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	SellerRepository *repository.SellerRepository
}

func NewSellerUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, storeRepos *repository.SellerRepository) *StoreUseCase {
	return &StoreUseCase{
		DB:              db,
		Log:             log,
		Validate:        validate,
		SellerRepository: storeRepos,
	}
}

func (u *StoreUseCase) Create(ctx context.Context, request *model.RegisterStore) (*model.StoreResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Failed to validate request body: %+v", err)
		return nil, err
	}
	// check if seller has a store
    hasStore, err := u.SellerRepository.HasStore(tx, request.ID)
    if err != nil {
        u.Log.Warnf("Failed to check if seller has store: %+v", err)
        return nil, model.ErrInternalServer
    }
    if hasStore {
        u.Log.Warnf("Seller already has a store")
        return nil, model.ErrConflict
    }
	store := &entity.Store{
		UserID: request.ID,
		StoreName: request.StoreName,
		Description: request.Description,
	}
	
	if err := u.SellerRepository.Create(tx, store); err != nil {
		u.Log.Warnf("Failed to create store: %+v", err)
		return nil, err
	}
	
	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, err
	}
	
	return converter.StoreToResponse(store), nil
}