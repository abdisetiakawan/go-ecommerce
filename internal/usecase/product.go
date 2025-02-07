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

type ProductUseCase struct {
	db          *gorm.DB
	val         *validator.Validate
	productRepo repo.ProductRepository
	storeRepo   repo.StoreRepository
	uuid        *helper.UUIDHelper
}

func NewProductUseCase(db *gorm.DB, validate *validator.Validate, productRepo repo.ProductRepository, storeRepo repo.StoreRepository, uuid *helper.UUIDHelper) interfaces.ProductUseCase {
	return &ProductUseCase{
		db:          db,
		val:         validate,
		productRepo: productRepo,
		storeRepo: storeRepo,
		uuid:        uuid,
	}
}


func (u *ProductUseCase) CreateProduct(ctx context.Context, request *model.RegisterProduct) (*model.ProductResponse, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, err
	}
	storeID, err := u.storeRepo.GetStoreIDByUserID(request.AuthID)
	if err != nil {
		return nil, err
	}
	product := &entity.Product{
		ProductUUID: u.uuid.Generate(),
		StoreID: storeID,
		ProductName: request.ProductName,
		Description: request.Description,
		Price: request.Price,
		Stock: request.Stock,
		Category: request.Category,
	}
	if err := u.productRepo.CreateProduct(product); err != nil {
		return nil, err
	}

	return converter.ProductToResponse(product), nil
}


func (u *ProductUseCase) GetProducts(ctx context.Context, request *model.GetProductsRequest) ([]model.ProductResponse, int64, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, 0, err
	}
	products, total, err := u.productRepo.GetProducts(request)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]model.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *converter.ProductsToResponse(&product)
	}
	return responses, total, nil
}

func (u *ProductUseCase) GetProductById(ctx context.Context, request *model.GetProductRequest) (*model.ProductResponse, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, err
	}
	response, err := u.productRepo.GetProductById(request.UserID, request.ProductUUID)
	if err != nil {
		return nil, err
	}
	return converter.ProductToResponse(response), nil
}

func (u *ProductUseCase) UpdateProduct(ctx context.Context, request *model.UpdateProduct) (*model.ProductResponse, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, err
	}
	product, err := u.productRepo.GetProductById(request.UserID, request.ProductName)
	if err != nil {
		return nil, err
	}
	if request.ProductName != "" {
		product.ProductName = request.ProductName
	}
	if request.Description != "" {
		product.Description = request.Description
	}
	if request.Price != 0 {
		product.Price = request.Price
	}
	if request.Stock != 0 {
		product.Stock = request.Stock
	}
	if request.Category != "" {
		product.Category = request.Category
	}
	if err := u.productRepo.UpdateProduct(product); err != nil {
		return nil, err
	}
	return converter.ProductToResponse(product), nil
}

func (u *ProductUseCase) DeleteProduct(ctx context.Context, request *model.DeleteProductRequest) error {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return err
	}
	product, err := u.productRepo.GetProductById(request.UserID, request.ProductUUID)
	if err != nil {
		return err
	}
	if err := u.productRepo.DeleteProduct(product); err != nil {
		return err
	}
	return nil
}