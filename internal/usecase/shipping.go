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

// NewShippingUseCase creates a new instance of ShippingUseCase.
// 
// Parameters:
//   - db: gorm.DB - Database connection for handling data operations.
//   - validate: validator.Validate - Validator instance for input validation.
//   - shippingRepo: repo.ShippingRepository - Repository for accessing shipping data.
//   - storeRepo: repo.StoreRepository - Repository for accessing store data.
//   - orderRepo: repo.OrderRepository - Repository for accessing order data.
//   - uuid: helper.UUIDHelper - Helper for generating UUIDs.
//   - orderEvent: ordereventUC.OrderEventUseCase - Use case for handling order events.
//
// Returns:
//   - interfaces.ShippingUseCase: A new ShippingUseCase instance with all necessary dependencies.

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

// UpdateShippingStatus updates the shipping status of an order.
//
// This function will start a transaction to update the order shipping status. If the order status is not "pending", it will return a 409 error.
//
// If the order status is "pending", it will check if the payment status is "paid". If the payment status is not "paid", it will return a 409 error.
//
// If the payment status is "paid", it will check if the shipping status is "pending". If the shipping status is not "pending", it will return a 409 error.
//
// If the shipping status is "pending", it will check if the request status is "shipped" or "delivered". If the request status is not "shipped" or "delivered", it will return a 400 error.
//
// If the request status is "shipped", it will update the shipping status to "shipped" and the order status to "shipped". If the request status is "delivered", it will update the shipping status to "delivered" and the order status to "completed".
//
// After updating the order, it will publish an event to kafka topic to update the order status.
//
// If there is an error when committing the transaction or publishing the event, it will rollback the transaction and return an error.
//
// Returns:
//
//	* 200 OK: model.OrderResponse if order shipping status is updated successfully.
//
// Errors:
//
//	* 400 Bad Request: if the request status is not "shipped" or "delivered".
//	* 409 Conflict: if the order status is not "pending" or the payment status is not "paid" or the shipping status is not "pending".
//	* Propagates error from use case layer if update fails.
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