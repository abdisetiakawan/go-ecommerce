package usecase

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	repo "github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StoreUseCase struct {
	db *gorm.DB
	log *logrus.Logger
	val *validator.Validate
	storeRepo repo.StoreRepository
	uuid *helper.UUIDHelper
}

func NewStoreUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, storeRepo repo.StoreRepository, uuid *helper.UUIDHelper) interfaces.StoreUseCase {
	return &StoreUseCase{
		db: db,
		log: log,
		val: validate,
		storeRepo: storeRepo,
		uuid: uuid,
	}
}

func (uc *StoreUseCase) RegisterStore(ctx context.Context, request *model.RegisterStore) (*model.StoreResponse, error) {
	if err := helper.ValidateStruct(uc.val, uc.log, request); err != nil {
		return nil, err
	}
	// check if seller has a store
    hasStore, err := uc.storeRepo.HasStore(uc.db, request.ID)
    if err != nil {
        uc.log.Warnf("Failed to check if seller has store: %+v", err)
        return nil, model.ErrInternalServer
    }
    if hasStore {
        uc.log.Warnf("Seller already has a store")
        return nil, model.ErrConflict
    }
	store := &entity.Store{
		UserID: request.ID,
		StoreName: request.StoreName,
		Description: request.Description,
	}
	
	if err := uc.storeRepo.CreateStore(store); err != nil {
		uc.log.Warnf("Failed to create store: %+v", err)
		return nil, err
	}
	
	return converter.StoreToResponse(store), nil
}

func (uc *StoreUseCase) GetStore(ctx context.Context, id uint) (*model.StoreResponse, error) {
	store, err := uc.storeRepo.FindStoreByUserID(id)
	if err != nil {
		return nil, err
	}
	return converter.StoreToResponse(&store), nil
}

func (uc *StoreUseCase) UpdateStore(ctx context.Context, request *model.UpdateStore) (*model.StoreResponse, error) {
	if err := helper.ValidateStruct(uc.val, uc.log, request); err != nil {
		return nil, err
	}
	store, err := uc.storeRepo.FindStoreByUserID(request.ID)
	if err != nil {
		return nil, err
	}
	if request.StoreName != "" {
		store.StoreName = request.StoreName
	}
	if request.Description != "" {
		store.Description = request.Description
	}
	if err := uc.storeRepo.UpdateStore(&store); err != nil {
		uc.log.Warnf("Failed to update store: %+v", err)
		return nil, err
	}
	return converter.StoreToResponse(&store), nil
}
