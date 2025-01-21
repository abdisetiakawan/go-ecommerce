package usecase

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
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
	UUIDHelper *helper.UUIDHelper
}

func NewSellerUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, sellerRepos *repository.SellerRepository, uuid *helper.UUIDHelper) *StoreUseCase {
	return &StoreUseCase{
		DB:              db,
		Log:             log,
		Validate:        validate,
		SellerRepository: sellerRepos,
		UUIDHelper: uuid,
	}
}

func (u *StoreUseCase) Create(ctx context.Context, request *model.RegisterStore) (*model.StoreResponse, error) {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return nil, err
	}
	// check if seller has a store
    hasStore, err := u.SellerRepository.HasStore(u.DB, request.ID)
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
	
	if err := u.SellerRepository.StoreRepository.Create(u.DB, store); err != nil {
		u.Log.Warnf("Failed to create store: %+v", err)
		return nil, err
	}
	
	return converter.StoreToResponse(store), nil
}

func (u *StoreUseCase) CreateProduct(ctx context.Context, request *model.RegisterProduct) (*model.ProductResponse, error) {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return nil, err
	}
	// check if seller has a store
	var store entity.Store
	if err := u.SellerRepository.CheckStore(u.DB, &store, request.AuthID); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("Store not found")
			return nil, model.ErrStoreNotFound
		}
		u.Log.Warnf("Failed to check store: %+v", err)
		return nil, model.ErrInternalServer
	}
	product := &entity.Product{
		ProductUUID: u.UUIDHelper.Generate(),
		StoreID: store.ID,
		ProductName: request.ProductName,
		Description: request.Description,
		Price: request.Price,
		Stock: request.Stock,
		Category: request.Category,
	}
	if err := u.SellerRepository.ProductRepository.Create(u.DB, product); err != nil {
		u.Log.Warnf("Failed to create product: %+v", err)
		return nil, err
	}

	return converter.ProductToResponse(product), nil
}

func (u *StoreUseCase) GetStore(ctx context.Context, id uint) (*model.StoreResponse, error) {
	var store entity.Store
	if err := u.SellerRepository.StoreRepository.FindByID(u.DB, &store, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("Store not found")
			return nil, model.ErrNotFound
		}
		u.Log.Warnf("Failed to get store: %+v", err)
		return nil, model.ErrInternalServer
	}
	return converter.StoreToResponse(&store), nil
}

func (u *StoreUseCase) Update(ctx context.Context, request *model.UpdateStore) (*model.StoreResponse, error) {
	var store entity.Store
	if err := u.SellerRepository.StoreRepository.FindByID(u.DB, &store, request.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("Store not found")
			return nil, model.ErrNotFound
		}
		u.Log.Warnf("Failed to get store: %+v", err)
		return nil, model.ErrInternalServer
	}
	if request.StoreName != "" {
		store.StoreName = request.StoreName
	}
	if request.Description != "" {
		store.Description = request.Description
	}
	if err := u.SellerRepository.StoreRepository.Update(u.DB, &store); err != nil {
		u.Log.Warnf("Failed to update store: %+v", err)
		return nil, err
	}
	return converter.StoreToResponse(&store), nil
}