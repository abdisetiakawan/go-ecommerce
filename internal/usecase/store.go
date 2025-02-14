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

// RegisterStore registers a new store and returns the newly created store in the response.
// It first validates the request body and checks if the user is a seller.
// If the user is not a seller, it returns a 403 error.
// If the request body is invalid, it returns a 400 error.
// If the store cannot be created, it returns a 500 error.
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

// GetStore retrieves a store by the user ID.
//
// Parameters:
//   * ctx: context.Context - Context for the request.
//   * id: uint - User ID associated with the store.
//
// Returns:
//   * model.StoreResponse: Store information if retrieval is successful.
//   * error: Propagates error from the repository layer if retrieval fails.

func (uc *StoreUseCase) GetStore(ctx context.Context, id uint) (*model.StoreResponse, error) {
	store, err := uc.storeRepo.FindStoreByUserID(id)
	if err != nil {
		return nil, err
	}
	return converter.StoreToResponse(&store), nil
}

// UpdateStore updates the information of a store associated with the given user ID.
// It first validates the request structure. If validation fails, it returns an error.
// It retrieves the store by user ID, and updates the store's name and description based on the request.
// If the store does not exist, or if any error occurs during the update process, it returns an error.
// 
// Parameters:
//   * ctx: context.Context - Context for the request.
//   * request: *model.UpdateStore - Request containing store update details.
//
// Returns:
//   * model.StoreResponse: The updated store information.
//   * error: If an error occurs during validation, retrieval, or update.

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
