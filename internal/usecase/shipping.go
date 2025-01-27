package usecase

import (
	"context"
	"fmt"

	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	repo "github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ShippingUseCase struct {
	db           *gorm.DB
	log          *logrus.Logger
	shippingRepo repo.ShippingRepository
	storeRepo    repo.StoreRepository
	orderRepo    repo.OrderRepository
	uuid         *helper.UUIDHelper
}

func NewShippingUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, shippingRepo repo.ShippingRepository, storeRepo repo.StoreRepository, orderRepo repo.OrderRepository, uuid *helper.UUIDHelper) interfaces.ShippingUseCase {
	return &ShippingUseCase{
		db:           db,
		log:          log,
		shippingRepo: shippingRepo,
		storeRepo: storeRepo,
		orderRepo: orderRepo,
		uuid:         uuid,
	}
}

func (c *ShippingUseCase) UpdateShippingStatus(ctx context.Context, request *model.UpdateShippingStatusRequest) (*model.OrderResponse, error) {
	tx := c.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	store, err := c.storeRepo.FindStoreByUserID(request.UserID)
	if err != nil {
		return nil, err
	}

	order, err := c.orderRepo.GetOrderBySeller(request.OrderUUID, store.ID)
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

	if err := c.orderRepo.UpdateOrder(order); err != nil {
		return nil, err
	}
	if err := c.shippingRepo.UpdateShipping(order.Shipping); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return converter.OrderToResponse(order), nil
}