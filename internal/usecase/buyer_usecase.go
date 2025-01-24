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
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BuyerUseCase struct {
    DB *gorm.DB
    Log *logrus.Logger
    Validate *validator.Validate
    BuyerRepository *repository.BuyerRepository
    UUIDHelper *helper.UUIDHelper
}

func NewBuyerUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, buyerRepos *repository.BuyerRepository, uuid *helper.UUIDHelper) interfaces.BuyerUseCase {
    return &BuyerUseCase{
        DB: db,
        Log: log,
        Validate: validate,
        BuyerRepository: buyerRepos,
        UUIDHelper: uuid,
    }
}

func (u *BuyerUseCase) CreateOrder(ctx context.Context, input *model.CreateOrder) (*model.OrderResponse, error) {
    tx := u.DB.WithContext(ctx).Begin()
    defer tx.Rollback()

    if err := helper.ValidateStruct(u.Validate, u.Log, input); err != nil {
        return nil, err
    }
    var productUUIDs []string
    for _, item := range input.Items {
        productUUIDs = append(productUUIDs, item.ProductUUID)
    }

    // Validate that all products belong to the same store
    storeID, err := u.BuyerRepository.FindStoreByProductUUIDs(tx, productUUIDs)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            u.Log.Warn("One or more products not found")
            return nil, model.NewApiError(fiber.StatusNotFound, "One or more products not found", nil)
        }
        if err == model.ErrBadRequest {
            u.Log.Warn("Products belong to different stores")
            return nil, model.NewApiError(fiber.StatusConflict, "Products must belong to the same store", nil)
        }
        u.Log.WithError(err).Error("Failed to validate products")
        return nil, model.ErrInternalServer
    }

    var totalPrice float64
    var orderItems []entity.OrderItem
    for _, item := range input.Items {
        var product entity.Product
        if err := u.BuyerRepository.ProductRepository.FindByUUID(tx, &product, item.ProductUUID, "product_uuid"); err != nil {
            u.Log.WithError(err).Errorf("Product with UUID %s not found in store %d", item.ProductUUID, storeID)
            return nil, model.NewApiError(fiber.StatusNotFound, fmt.Sprintf("Product with UUID %s not found", item.ProductUUID), nil)
        }

        if product.Stock < item.Quantity {
            u.Log.Warnf("Product %s has insufficient stock", product.ProductName)
            return nil, model.NewApiError(fiber.StatusConflict, fmt.Sprintf("Product %s has insufficient stock", product.ProductName), nil)
        }

        product.Stock -= item.Quantity
        if err := u.BuyerRepository.ProductRepository.Update(tx, &product); err != nil {
            u.Log.WithError(err).Error("Failed to update product stock")
            return nil, model.ErrInternalServer
        }

        itemTotal := float64(item.Quantity) * product.Price
        totalPrice += itemTotal

        orderItems = append(orderItems, entity.OrderItem{
            OrderItemUUID: u.UUIDHelper.Generate(),
            ProductID:    product.ID,
            Quantity:     item.Quantity,
            TotalPrice:   itemTotal,
        })
    }

    order := &entity.Order{
        OrderUUID:  u.UUIDHelper.Generate(),
        UserID:     input.UserID,
        Status:     "pending",
        TotalPrice: totalPrice,
        Items:      orderItems,
    }

    if err := u.BuyerRepository.OrderRepository.Create(tx, order); err != nil {
        u.Log.WithError(err).Error("Failed to create order")
        return nil, model.ErrInternalServer
    }

    payment := &entity.Payment{
        PaymentUUID: u.UUIDHelper.Generate(),
        OrderID:     order.ID,
        Amount:      totalPrice,
        Status:      "pending",
        Method:      input.Payments.PaymentMethod,
    }
    if err := u.BuyerRepository.PaymentRepository.Create(tx, payment); err != nil {
        u.Log.WithError(err).Error("Failed to create payment")
        return nil, model.ErrInternalServer
    }

    shipping := &entity.Shipping{
        ShippingUUID: u.UUIDHelper.Generate(),
        OrderID:      order.ID,
        Address:      input.ShippingAddress.Address,
        City:         input.ShippingAddress.City,
        Province:     input.ShippingAddress.Province,
        PostalCode:   input.ShippingAddress.PostalCode,
        Status:       "pending",
    }
    if err := u.BuyerRepository.ShippingRepository.Create(tx, shipping); err != nil {
        u.Log.WithError(err).Error("Failed to create shipping")
        return nil, model.ErrInternalServer
    }

    order.Payment = payment
    order.Shipping = shipping

    if err := tx.Commit().Error; err != nil {
        u.Log.WithError(err).Error("Failed to commit transaction")
        return nil, model.ErrInternalServer
    }

    return converter.OrderToResponse(order), nil
}


func (u *BuyerUseCase) GetOrders(ctx context.Context, request *model.SearchOrderRequest) ([]model.ListOrderResponse, int64, error) {
    if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
        return nil, 0, err
    }
    tasks, total, err := u.BuyerRepository.GetOrders(u.DB, request)
    if err != nil {
        u.Log.WithError(err).Error("Failed to get orders")
        return nil, 0, model.ErrInternalServer
    }
    responses := make([]model.ListOrderResponse, len(tasks))
    for i, task := range tasks {
        responses[i] = *converter.OrdersToResponse(&task)
    }
    return responses, total, nil
}

func(u *BuyerUseCase) GetOrder(ctx context.Context, request *model.GetOrderDetails) (*model.OrderResponse, error) {
    if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
        return nil, err
    }
    order, err := u.BuyerRepository.GetOrder(u.DB, request)
    if err != nil {
        u.Log.WithError(err).Error("Failed to get order")
        return nil, err
    }
    return converter.OrderToResponse(order), nil
}

func (u *BuyerUseCase) CancelOrder(ctx context.Context, request *model.CancelOrderRequest) (*model.OrderResponse, error) {
    tx := u.DB.WithContext(ctx).Begin()
    defer tx.Rollback()

    if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
        return nil, err
    }
    
    order, err := u.BuyerRepository.GetOrder(tx, &model.GetOrderDetails{
        OrderUUID: request.OrderUUID,
        UserID:    request.UserID,
    })
    if err != nil {
        u.Log.WithError(err).Error("Failed to get order")
        return nil, err
    }

    if order.Status == "completed" || order.Status == "cancelled" || order.Status == "shipped" || order.Status == "processed" {
        return nil, model.NewApiError(fiber.StatusConflict, fmt.Sprintf("Order with ID %s cannot be cancelled, current status is %s", request.OrderUUID, order.Status), nil)
    }
    
    order.Status = "cancelled"
    if err := u.BuyerRepository.OrderRepository.Update(tx, order); err != nil {
        u.Log.WithError(err).Error("Failed to update order")
        return nil, model.ErrInternalServer
    }

    if order.Payment != nil {
        order.Payment.Status = "cancelled"
        if err := u.BuyerRepository.PaymentRepository.Update(tx, order.Payment); err != nil {
            u.Log.WithError(err).Error("Failed to update payment")
            return nil, model.ErrInternalServer
        }
    }

    if order.Shipping != nil {
        order.Shipping.Status = "cancelled"
        if err := u.BuyerRepository.ShippingRepository.Update(tx, order.Shipping); err != nil {
            u.Log.WithError(err).Error("Failed to update shipping")
            return nil, model.ErrInternalServer
        }
    }

    for _, item := range order.Items {
        var product entity.Product
        if err := u.BuyerRepository.ProductRepository.FindByID(tx, &product, item.ProductID); err != nil {
            u.Log.WithError(err).Error("Failed to find product")
            return nil, model.ErrInternalServer
        }
        product.Stock += item.Quantity
        if err := u.BuyerRepository.ProductRepository.Update(tx, &product); err != nil {
            u.Log.WithError(err).Error("Failed to update product")
            return nil, model.ErrInternalServer
        }
    }
    
    if err := tx.Commit().Error; err != nil {
        u.Log.WithError(err).Error("Failed to commit transaction")
        return nil, model.ErrInternalServer
    }
    
    return converter.OrderToResponse(order), nil
}

func (u *BuyerUseCase) CheckoutOrder(ctx context.Context, request *model.CheckoutOrderRequest) (*model.OrderResponse, error) {
    tx := u.DB.WithContext(ctx).Begin()
    defer tx.Rollback()

    if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
        return nil, err
    }

    order, err := u.BuyerRepository.GetOrder(tx, &model.GetOrderDetails{
        OrderUUID: request.OrderUUID,
        UserID:    request.UserID,
    })
    if err != nil {
        u.Log.WithError(err).Error("Failed to get order")
        return nil, err
    }

    if order.Status != "pending" {
        return nil, model.NewApiError(fiber.StatusConflict, fmt.Sprintf("Order with ID %s cannot be checked out, current status is %s", request.OrderUUID, order.Status), nil)
    }
    order.Status = "processed"
    if err := u.BuyerRepository.OrderRepository.Update(tx, order); err != nil {
        u.Log.WithError(err).Error("Failed to update order")
        return nil, model.ErrInternalServer
    }
    if order.Payment != nil {
        order.Payment.Status = "paid"
        if err := u.BuyerRepository.PaymentRepository.Update(tx, order.Payment); err != nil {
            u.Log.WithError(err).Error("Failed to update payment")
            return nil, model.ErrInternalServer
        }
    }

    if err := tx.Commit().Error; err != nil {
        u.Log.WithError(err).Error("Failed to commit transaction")
        return nil, model.ErrInternalServer
    }
    return converter.OrderToResponse(order), nil
}