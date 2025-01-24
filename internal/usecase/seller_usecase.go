package usecase

import (
	"context"
	"fmt"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SellerUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	SellerRepository *repository.SellerRepository
	UUIDHelper *helper.UUIDHelper
}

func NewSellerUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, sellerRepos *repository.SellerRepository, uuid *helper.UUIDHelper) interfaces.SellerUseCase {
	return &SellerUseCase{
		DB:              db,
		Log:             log,
		Validate:        validate,
		SellerRepository: sellerRepos,
		UUIDHelper: uuid,
	}
}

func (u *SellerUseCase) Create(ctx context.Context, request *model.RegisterStore) (*model.StoreResponse, error) {
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

func (u *SellerUseCase) CreateProduct(ctx context.Context, request *model.RegisterProduct) (*model.ProductResponse, error) {
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

func (u *SellerUseCase) GetStore(ctx context.Context, id uint) (*model.StoreResponse, error) {
	var store entity.Store
	if err := u.SellerRepository.StoreRepository.FindByUserID(u.DB, &store, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("Store not found")
			return nil, model.ErrNotFound
		}
		u.Log.Warnf("Failed to get store: %+v", err)
		return nil, model.ErrInternalServer
	}
	return converter.StoreToResponse(&store), nil
}

func (u *SellerUseCase) Update(ctx context.Context, request *model.UpdateStore) (*model.StoreResponse, error) {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return nil, err
	}
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

func (u *SellerUseCase) GetProducts (ctx context.Context, request *model.GetProductsRequest) ([]model.ProductResponse, int64, error) {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return nil, 0, err
	}
	products, total, err := u.SellerRepository.GetProducts(u.DB, request)
	if err != nil {
		u.Log.WithError(err).Errorf("Failed to get products")
		return nil, 0, err
	}
	responses := make([]model.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *converter.ProductsToResponse(&product)
	}
	return responses, total, nil
}

func (u *SellerUseCase) GetProduct (ctx context.Context, request *model.GetProductRequest) (*model.ProductResponse, error) {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return nil, err
	}
	response, err := u.SellerRepository.GetProduct(u.DB, request)
	if err != nil {
		u.Log.WithError(err).Errorf("Failed to get product")
		return nil, err
	}
	return converter.ProductToResponse(response), nil
}

func (u *SellerUseCase) UpdateProduct (ctx context.Context, request *model.UpdateProduct) (*model.ProductResponse, error) {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return nil, err
	}
	var product entity.Product
	if err := u.SellerRepository.CheckProduct(u.DB, &product, request.UserID, request.ProductUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("Product not found")
			return nil, model.NewApiError(fiber.StatusNotFound, fmt.Sprintf("Product with id %s not found", request.ProductUUID), nil)
		}
		u.Log.WithError(err).Errorf("Failed to update product")
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
	if err := u.SellerRepository.ProductRepository.Update(u.DB, &product); err != nil {
		u.Log.WithError(err).Errorf("Failed to update product")
		return nil, err
	}
	return converter.ProductToResponse(&product), nil
}

func (u *SellerUseCase) DeleteProduct (ctx context.Context, request *model.DeleteProductRequest) error {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return err
	}
	var product entity.Product
	if err := u.SellerRepository.CheckProduct(u.DB, &product, request.UserID, request.ProductUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("Product not found")
			return model.NewApiError(fiber.StatusNotFound, fmt.Sprintf("Product with id %s not found", request.ProductUUID), nil)
		}
		u.Log.WithError(err).Errorf("Failed to delete product")
		return err
	}
	if err := u.SellerRepository.ProductRepository.SoftDelete(u.DB, product.ID); err != nil {
		u.Log.WithError(err).Errorf("Failed to delete product")
		return err
	}
	return nil
}

func (u *SellerUseCase) GetOrder(ctx context.Context, request *model.GetOrderDetails) (*model.OrderResponse, error) {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return nil, err
	}
	var store entity.Store
	if err := u.SellerRepository.StoreRepository.FindByUserID(u.DB, &store, request.UserID); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("Store not found")
			return nil, model.ErrStoreNotFound
		}
		return nil, err
	}
	order, err := u.SellerRepository.GetOrder(u.DB, request.OrderUUID, store.ID)
	if err != nil {
		u.Log.WithError(err).Errorf("Failed to get order")
		return nil, err
	}
	return converter.OrderToResponse(order), nil
}

func (u *SellerUseCase) GetOrders(ctx context.Context, request *model.SearchOrderRequestBySeller) ([]model.OrdersResponseForSeller, int64, error) {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return nil, 0, err
	}
	var store entity.Store
	if err := u.SellerRepository.StoreRepository.FindByUserID(u.DB, &store, request.UserID); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("Store not found")
			return nil, 0, model.ErrStoreNotFound
		}
		return nil, 0, err
	}
	request.StoreID = store.ID
	orders, total, err := u.SellerRepository.GetOrders(u.DB, request)
	if err != nil {
		u.Log.WithError(err).Errorf("Failed to get orders")
		return nil, 0, err
	}
	responses := make([]model.OrdersResponseForSeller, len(orders))
	for i, order := range orders {
		responses[i] = *converter.OrderToResponseForSeller(&order)
	}
	return responses, total, nil
}

func (c *SellerUseCase) UpdateShippingStatus(ctx context.Context, request *model.UpdateShippingStatusRequest) (*model.OrderResponse, error) {
    tx := c.DB.WithContext(ctx).Begin()
    defer tx.Rollback()

    var store entity.Store
    if err := c.SellerRepository.StoreRepository.FindByUserID(tx, &store, request.UserID); err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, model.ErrStoreNotFound
        }
        return nil, err
    }

    order, err := c.SellerRepository.GetOrder(tx, request.OrderUUID, store.ID)
    if err != nil {
        return nil, err
    }

    switch {
    case order.Payment.Status != "paid":
        return nil, model.NewApiError(fiber.StatusConflict, 
            fmt.Sprintf("Cannot update shipping. Payment status is %s", order.Payment.Status), nil)
    
    case order.Shipping.Status == "delivered":
        return nil, model.NewApiError(fiber.StatusConflict, fmt.Sprintf("Cannot update shipping. Order is already %s", order.Status), nil)
    
    case request.Status == "shipped" && order.Shipping.Status != "pending":
        return nil, model.NewApiError(fiber.StatusConflict, 
            fmt.Sprintf("Cannot ship order. Current shipping status is %s", order.Shipping.Status), nil)

	case request.Status == "delivered" && order.Shipping.Status == "pending": 
		return nil, model.NewApiError(fiber.StatusConflict, 
			fmt.Sprintf("Cannot deliver order. Current shipping status is %s", order.Shipping.Status), nil)
			
	case request.Status == order.Shipping.Status: 
		return nil, model.NewApiError(fiber.StatusConflict, 
			fmt.Sprintf("Cannot update shipping. Current shipping status is %s", order.Shipping.Status), nil)
    }

	order.Shipping.Status = request.Status
	if request.Status == "shipped" {
		order.Status = "shipped"
	}
	if request.Status == "delivered" {
		order.Status = "completed"
	}

    if err := c.SellerRepository.OrderRepository.Update(tx, order); err != nil {
        return nil, err
    }
	if err := c.SellerRepository.ShippingRepository.Update(tx, order.Shipping); err != nil {
		return nil, err
	}

    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return converter.OrderToResponse(order), nil
}