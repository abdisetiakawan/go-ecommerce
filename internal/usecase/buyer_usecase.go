package usecase

import (
	"context"
	"fmt"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository"
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

func NewBuyerUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, buyerRepos *repository.BuyerRepository, uuid *helper.UUIDHelper) *BuyerUseCase {
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


    var totalPrice float64
    var orderItems []entity.OrderItem
    var product entity.Product

    for _, item := range input.Items {
        err := u.BuyerRepository.ProductRepository.FindByUUID(tx, &product, item.ProductUUID, "product_uuid")
        if err != nil {
            u.Log.WithError(err).Errorf("Product with ID %s not found", item.ProductUUID)
            return nil, model.NewApiError(fiber.StatusNotFound, fmt.Sprintf("Product with ID %s not found", item.ProductUUID), nil)
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

    // Create order
    order := &entity.Order{
        OrderUUID:  u.UUIDHelper.Generate(),
        UserID:     input.UserID,
        Status:     "pending",
        TotalPrice: totalPrice,
        Items:      orderItems,
        Payment:    nil,
    }

    if err := u.BuyerRepository.OrderRepository.Create(tx, order); err != nil {
        u.Log.WithError(err).Error("Failed to create order")
        return nil, model.ErrInternalServer
    }

    // Create payment
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

    // Create shipping
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