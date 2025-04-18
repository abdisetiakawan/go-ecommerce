package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
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

type OrderUseCase struct {
	db        *gorm.DB
	val       *validator.Validate
	orderRepo repo.OrderRepository
    productRepo repo.ProductRepository
    storeRepo repo.StoreRepository
	orderEvent ordereventUC.OrderEventUseCase
	uuid      *helper.UUIDHelper
}

func NewOrderUseCase(db *gorm.DB, validate *validator.Validate, orderRepo repo.OrderRepository, productRepo repo.ProductRepository, storeRepo repo.StoreRepository, uuid *helper.UUIDHelper, orderEvent ordereventUC.OrderEventUseCase) interfaces.OrderUseCase {
	return &OrderUseCase{
		db:        db,
		val:       validate,
		orderRepo: orderRepo,
        productRepo: productRepo,
        storeRepo: storeRepo,
		uuid:      uuid,
		orderEvent: orderEvent,
	}
}

// CreateOrder validates the input, verifies product availability and stock,
// and creates an order with associated order items. It ensures all products
// belong to the same store and reduces stock for each product. The method
// then creates an order event with payment and shipping data and commits the
// transaction. If any step fails, the transaction is rolled back, and an
// appropriate error is returned. The function launches an asynchronous
// process to handle the order event and returns the order response upon success.

func (uc *OrderUseCase) CreateOrder(ctx context.Context, input *model.CreateOrder) (*model.OrderResponse, error) {
	tx := uc.db.Begin()
    defer tx.Rollback()
	if err := helper.ValidateStruct(uc.val, input); err != nil {
		return nil, err
	}

	var productUUIDs []string
	for _, item := range input.Items {
		productUUIDs = append(productUUIDs, item.ProductUUID)
	}

	// val that all products belong to the same store
	_, err := uc.orderRepo.FindStoreByProductUUIDs(productUUIDs)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewApiError(fiber.StatusNotFound, "One or more products not found", nil)
		}
		if err == model.ErrBadRequest {
			return nil, model.NewApiError(fiber.StatusConflict, "Products must belong to the same store", nil)
		}
		return nil, model.ErrInternalServer
	}

	var totalPrice float64
	var orderItems []entity.OrderItem
	for _, item := range input.Items {
		product, err := uc.productRepo.FindProductByUUID(item.ProductUUID)
		if err != nil {
			return nil, err
		}

		if product.Stock < item.Quantity {
			return nil, model.NewApiError(fiber.StatusConflict, fmt.Sprintf("Product %s has insufficient stock", product.ProductName), nil)
		}

		product.Stock -= item.Quantity
        if err := uc.productRepo.UpdateProduct(&product); err != nil {
            return nil, model.ErrInternalServer
        }

		itemTotal := float64(item.Quantity) * product.Price
		totalPrice += itemTotal

		orderItems = append(orderItems, entity.OrderItem{
			OrderItemUUID: uc.uuid.Generate(),
			ProductID:    product.ID,
			Quantity:     item.Quantity,
			TotalPrice:   itemTotal,
		})
	}

	order := &entity.Order{
		OrderUUID:  uc.uuid.Generate(),
		UserID:     input.UserID,
		Status:     "pending",
		TotalPrice: totalPrice,
		Items:      orderItems,
	}

	if err := uc.orderRepo.CreateOrder(order); err != nil {
		return nil, model.ErrInternalServer
	}

	paymentData, err := json.Marshal(eventmodel.PaymentMessage{
		OrderID:     order.ID,
		PaymentUUID: uc.uuid.Generate(),
		Amount:      totalPrice,
		Method:      input.Payments.PaymentMethod,
		Status:      "pending",
	})
	if err != nil {
		return nil, model.ErrInternalServer
	}
	
	shippingData, err := json.Marshal(eventmodel.ShippingMessage{
		ShippingUUID: uc.uuid.Generate(),
		OrderID:       order.ID,
		Address:      input.ShippingAddress.Address,
		City:        input.ShippingAddress.City,
		Province:    input.ShippingAddress.Province,
		PostalCode:  input.ShippingAddress.PostalCode,
		Status:      "pending",
	})
	if err != nil {
		return nil, model.ErrInternalServer
	}
	
	orderEvent := &evententity.OrderEvent{
		EventUUID:    uc.uuid.Generate(),
		OrderID:      order.ID,
		EventType:    "order_created",
		Status:       "pending",
		PaymentData:  paymentData,
		ShippingData: shippingData,
	}
	

    if err := tx.Create(orderEvent).Error; err != nil {
        return nil, model.ErrInternalServer
    }

    if err := tx.Commit().Error; err != nil {
        return nil, model.ErrInternalServer
    }

	go uc.orderEvent.ProcessOrderEvent(ctx, orderEvent)

	var paymentMessage eventmodel.PaymentMessage
	if err := json.Unmarshal(paymentData, &paymentMessage); err != nil {
		return nil, model.ErrInternalServer
	}

	var shippingMessage eventmodel.ShippingMessage
	if err := json.Unmarshal(shippingData, &shippingMessage); err != nil {
		return nil, model.ErrInternalServer
	}

	return converter.CreateOrderToResponse(&paymentMessage, &shippingMessage, order), nil
}


// GetOrdersByBuyer returns a list of orders made by a buyer, given a search request.
//
// The search request can filter by status, start date, end date, page, and limit.
//
// The response is a list of ListOrderResponse, which only contains the basic information of the order.
//
// The total is the total number of orders that match the search request.
func (uc *OrderUseCase) GetOrdersByBuyer(ctx context.Context, request *model.SearchOrderRequest) ([]model.ListOrderResponse, int64, error) {
    if err := helper.ValidateStruct(uc.val, request); err != nil {
        return nil, 0, err
    }
    tasks, total, err := uc.orderRepo.GetOrdersByBuyer(request)
    if err != nil {
        return nil, 0, model.ErrInternalServer
    }
    responses := make([]model.ListOrderResponse, len(tasks))
    for i, task := range tasks {
        responses[i] = *converter.OrdersToResponse(&task)
    }
    return responses, total, nil
}

// GetOrderByIdByBuyer returns a single order by order UUID and user ID.
//
// The response is a OrderResponse, which contains the detailed information of the order.
//
// If the order is not found, it returns a 404 error.
func(uc *OrderUseCase) GetOrderByIdByBuyer(ctx context.Context, request *model.GetOrderDetails) (*model.OrderResponse, error) {
    if err := helper.ValidateStruct(uc.val, request); err != nil {
        return nil, err
    }
    order, err := uc.orderRepo.GetOrderByIdByBuyer(request)
    if err != nil {
        return nil, err
    }
    return converter.OrderToResponse(order), nil
}

    // CancelOrder cancels an order with the given order UUID and user ID.
    //
    // If the order is not found, it returns a 404 error.
    //
    // If the order status is not "pending", it returns a 409 error.
    //
    // The function will return an error if there is an error when updating the order in the database.
    //
    // The function will also publish an event to kafka topic to cancel the order payment and shipping.
func (uc *OrderUseCase) CancelOrder(ctx context.Context, request *model.CancelOrderRequest) (*model.OrderResponse, error) {
    tx := uc.db.WithContext(ctx).Begin()
    defer tx.Rollback()

    if err := helper.ValidateStruct(uc.val, request); err != nil {
        return nil, err
    }
    
    order, err := uc.orderRepo.GetOrderByIdByBuyer(&model.GetOrderDetails{
        OrderUUID: request.OrderUUID,
        UserID:    request.UserID,
    })
    if err != nil {
        return nil, err
    }

    if order == nil {
        return nil, model.NewApiError(fiber.StatusNotFound, "Order not found", nil)
    }

    if order.Status == "completed" || order.Status == "cancelled" || order.Status == "shipped" || order.Status == "processed" {
        return nil, model.NewApiError(fiber.StatusConflict, fmt.Sprintf("Order with ID %s cannot be cancelled, current status is %s", request.OrderUUID, order.Status), nil)
    }
    
    order.Status = "cancelled"
    if err := tx.Model(&entity.Order{}).Where("id = ?", order.ID).Update("status", "cancelled").Error; err != nil {
        return nil, model.ErrInternalServer
    }

    paymentStatus, err := json.Marshal(eventmodel.PaymentMessage{
        OrderID: order.ID,
        Status:  "cancelled",
    })
    if err != nil {
        return nil, model.ErrInternalServer
    }

    shippingStatus, err := json.Marshal(eventmodel.ShippingMessage{
        OrderID: order.ID,
        Status:  "cancelled",
    })
    if err != nil {
        return nil, model.ErrInternalServer
    }

    orderEvent := &evententity.OrderEvent{
        EventUUID:    uc.uuid.Generate(),
        OrderID:      order.ID,
        EventType:    "order_cancelled",
        Status:       "pending",
        PaymentData:  paymentStatus,
        ShippingData: shippingStatus,
    }

    if err := tx.Create(orderEvent).Error; err != nil {
        return nil, model.ErrInternalServer
    }

    for _, item := range order.Items {
        if err := tx.Exec("UPDATE products SET stock = stock + ? WHERE id = ?", item.Quantity, item.ProductID).Error; err != nil {
            return nil, model.ErrInternalServer
        }
    }

    if err := tx.Commit().Error; err != nil {
        return nil, model.ErrInternalServer
    }

    go uc.orderEvent.CancelOrderEvent(ctx, orderEvent)
	order.Shipping.Status = "cancelled"
	order.Payment.Status = "cancelled"
    return converter.OrderToResponse(order), nil
}

    // CheckoutOrder checks out an order. It will only work if the order status is "pending". If the order status is not "pending", it will return an error.
    // It will update the order status to "processed" and create an order event with type "payment_processed" and status "pending".
    // It will then commit the transaction and send the order event to kafka topic to be processed.
    // If there is an error when committing the transaction or sending the order event, it will rollback the transaction and return an error.
func (uc *OrderUseCase) CheckoutOrder(ctx context.Context, request *model.CheckoutOrderRequest) (*model.OrderResponse, error) {
    tx := uc.db.WithContext(ctx).Begin()
    defer tx.Rollback()

    if err := helper.ValidateStruct(uc.val, request); err != nil {
        return nil, err
    }

    order, err := uc.orderRepo.GetOrderByIdByBuyer(&model.GetOrderDetails{
        OrderUUID: request.OrderUUID,
        UserID:    request.UserID,
    })
    if err != nil {
        return nil, err
    }

    if order.Status != "pending" {
        return nil, model.NewApiError(fiber.StatusConflict, fmt.Sprintf("Order with ID %s cannot be checked out, current status is %s", request.OrderUUID, order.Status), nil)
    }
    order.Status = "processed"
    if err := uc.orderRepo.UpdateOrder(order); err != nil {
        return nil, model.ErrInternalServer
    }
	paymentStatus, err := json.Marshal(eventmodel.PaymentMessage{
		OrderID: order.ID,
		Status:  "paid",
	})
	if err != nil {
		return nil, model.ErrInternalServer
	}
	orderEvent := &evententity.OrderEvent{
		EventUUID: uc.uuid.Generate(),
		OrderID: order.ID,
		EventType: "payment_processed",
		Status: "pending",
		PaymentData: paymentStatus,
	}
	if err := tx.Create(orderEvent).Error; err != nil {
		return nil, model.ErrInternalServer
	}

    if err := tx.Commit().Error; err != nil {
        return nil, model.ErrInternalServer
    }

	go uc.orderEvent.CheckoutOrderEvent(ctx, orderEvent)
	order.Payment.Status = "paid"
    return converter.OrderToResponse(order), nil
}

// GetOrdersBySeller retrieves a list of orders for a specific seller based on the provided search request.
// 
// The function first validates the input request structure. It then identifies the store associated
// with the seller's user ID and fetches orders linked to that store. The search request can filter
// orders by status, page, and limit.
//
// Returns:
// - A slice of OrdersResponseForSeller containing the details of each order.
// - The total number of orders that match the search criteria.
// - An error, if any issue occurs during the process.

func (u *OrderUseCase) GetOrdersBySeller(ctx context.Context, request *model.SearchOrderRequestBySeller) ([]model.OrdersResponseForSeller, int64, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, 0, err
	}
	store, err := u.storeRepo.FindStoreByUserID(request.UserID)
    if err != nil {
        return nil, 0, err
    }
	request.StoreID = store.ID
	orders, total, err := u.orderRepo.GetOrdersBySeller(request)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]model.OrdersResponseForSeller, len(orders))
	for i, order := range orders {
		responses[i] = *converter.OrderToResponseForSeller(&order)
	}
	return responses, total, nil
}

// GetOrderBySeller retrieves a single order by order UUID and user ID.
// 
// The function first validates the input request structure. It then identifies the store associated
// with the seller's user ID and fetches the order linked to that store.
// 
// Returns:
// - A OrderResponse containing the details of the order.
// - An error, if any issue occurs during the process.
func (u *OrderUseCase) GetOrderBySeller(ctx context.Context, request *model.GetOrderDetails) (*model.OrderResponse, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, err
	}
	store, err := u.storeRepo.FindStoreByUserID(request.UserID)
    if err != nil {
        return nil, err
    }
	order, err := u.orderRepo.GetOrderBySeller(request.OrderUUID, store.ID)
	if err != nil {
		return nil, err
	}
	return converter.OrderToResponse(order), nil
}