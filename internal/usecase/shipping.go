package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	evententity "github.com/abdisetiakawan/go-ecommerce/internal/entity/event_entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	eventmodel "github.com/abdisetiakawan/go-ecommerce/internal/model/event_model"
	repo "github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	ordereventUC "github.com/abdisetiakawan/go-ecommerce/internal/usecase/event_uc/interfaces"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber"
	"gorm.io/gorm"
)

type ShippingUseCase struct {
	db           *gorm.DB
	shippingRepo repo.ShippingRepository
	storeRepo    repo.StoreRepository
	orderRepo    repo.OrderRepository
	uuid         *helper.UUIDHelper
	orderEvent   ordereventUC.OrderEventUseCase
}

func NewShippingUseCase(db *gorm.DB, validate *validator.Validate, shippingRepo repo.ShippingRepository, storeRepo repo.StoreRepository, orderRepo repo.OrderRepository, uuid *helper.UUIDHelper, orderEvent ordereventUC.OrderEventUseCase) interfaces.ShippingUseCase {
	return &ShippingUseCase{
		db:           db,
		shippingRepo: shippingRepo,
		storeRepo: storeRepo,
		orderRepo: orderRepo,
		uuid:         uuid,
		orderEvent:   orderEvent,
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
	var eventType string
	if request.Status == "shipped" {
		eventType = "shipping_processed"
		order.Status = "shipped"
	} else if request.Status == "delivered" {
		eventType = "order_delivered"
		order.Status = "completed"
	} else {
		return nil, model.NewApiError(fiber.StatusBadRequest,
			fmt.Sprintf("Invalid status: %s", request.Status), nil)
	}
	
	if err := c.shippingRepo.UpdateShipping(order.Shipping); err != nil {
		return nil, err
	}
	
	orderStatus, err := json.Marshal(eventmodel.OrderMessage{
		OrderID: order.ID,
		Status:  order.Status,
	})
	if err != nil {
		return nil, err
	}
	orderEvent := &evententity.OrderEvent{
		EventUUID:  c.uuid.Generate(),
		OrderID:    order.ID,
		Status:     "pending",
		EventType:  eventType,
		OrderData: orderStatus,
	}
	if err := tx.Create(orderEvent).Error; err != nil {
		return nil, err
	}

	go c.orderEvent.ChangeOrderStatusUC(ctx, orderEvent)

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return converter.OrderToResponse(order), nil
}