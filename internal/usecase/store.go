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
	"gorm.io/gorm"
)

type StoreUseCase struct {
	db *gorm.DB
	val *validator.Validate
	storeRepo repo.StoreRepository
	uuid *helper.UUIDHelper
}

func NewStoreUseCase(db *gorm.DB, validate *validator.Validate, storeRepo repo.StoreRepository, uuid *helper.UUIDHelper) interfaces.StoreUseCase {
	return &StoreUseCase{
		db: db,
		val: validate,
		storeRepo: storeRepo,
		uuid: uuid,
	}
}

func (uc *StoreUseCase) RegisterStore(ctx context.Context, request *model.RegisterStore) (*model.StoreResponse, error) {
	if err := helper.ValidateStruct(uc.val, request); err != nil {
		return nil, err
	}
	// check if seller has a store
    hasStore, err := uc.storeRepo.HasStore(uc.db, request.ID)
    if err != nil {
        return nil, model.ErrInternalServer
    }
    if hasStore {
        return nil, model.ErrConflict
    }
	store := &entity.Store{
		UserID: request.ID,
		StoreName: request.StoreName,
		Description: request.Description,
	}
	
	if err := uc.storeRepo.CreateStore(store); err != nil {
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
	if err := helper.ValidateStruct(uc.val, request); err != nil {
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
		return nil, err
	}
	return converter.StoreToResponse(&store), nil
}
