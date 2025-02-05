package usecase

import (
	"context"
	"fmt"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/event"
	repo "github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderUseCase struct {
	db        *gorm.DB
	log       *logrus.Logger
	val       *validator.Validate
	orderRepo repo.OrderRepository
    productRepo repo.ProductRepository
    paymentRepo repo.PaymentRepository
    shippingRepo repo.ShippingRepository
    storeRepo repo.StoreRepository
	uuid      *helper.UUIDHelper
	kafka *helper.KafkaProducer
}

func NewOrderUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, orderRepo repo.OrderRepository, productRepo repo.ProductRepository, paymentRepo repo.PaymentRepository, shippingRepo repo.ShippingRepository, storeRepo repo.StoreRepository, uuid *helper.UUIDHelper, kafka *helper.KafkaProducer) interfaces.OrderUseCase {
	return &OrderUseCase{
		db:        db,
		log:       log,
		val:       validate,
		orderRepo: orderRepo,
        productRepo: productRepo,
        paymentRepo: paymentRepo,
        shippingRepo: shippingRepo,
        storeRepo: storeRepo,
		uuid:      uuid,
		kafka: kafka,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, input *model.CreateOrder) (*model.OrderResponse, error) {
	if err := helper.ValidateStruct(uc.val, uc.log, input); err != nil {
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
			uc.log.Warn("One or more products not found")
			return nil, model.NewApiError(fiber.StatusNotFound, "One or more products not found", nil)
		}
		if err == model.ErrBadRequest {
			uc.log.Warn("Products belong to different stores")
			return nil, model.NewApiError(fiber.StatusConflict, "Products must belong to the same store", nil)
		}
		uc.log.WithError(err).Error("Failed to validate products")
		return nil, model.ErrInternalServer
	}

	var totalPrice float64
	var orderItems []entity.OrderItem
	for _, item := range input.Items {
		product, err := uc.productRepo.FindProductByUUID(item.ProductUUID)
		if err != nil {
			uc.log.WithError(err).Error("Failed to find product")
			return nil, err
		}

		if product.Stock < item.Quantity {
			uc.log.Warnf("Product %s has insufficient stock", product.ProductName)
			return nil, model.NewApiError(fiber.StatusConflict, fmt.Sprintf("Product %s has insufficient stock", product.ProductName), nil)
		}

		product.Stock -= item.Quantity
        if err := uc.productRepo.UpdateProduct(&product); err != nil {
            uc.log.WithError(err).Error("Failed to update product stock")
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
		uc.log.WithError(err).Error("Failed to create order")
		return nil, model.ErrInternalServer
	}

	paymentMessage := &event.PaymentMessage{
		Status:  "pending",
		PaymentUUID: uc.uuid.Generate(),
		OrderID: order.ID,
		Amount:  totalPrice,
		Method:  input.Payments.PaymentMethod,
	}
	shippingMessage := &event.ShippingMessage{
		ShippingUUID: uc.uuid.Generate(),
		Status:       "pending",
		OrderID: order.ID,
		Address: input.ShippingAddress.Address,
		City:    input.ShippingAddress.City,
		Province: input.ShippingAddress.Province,
		PostalCode: input.ShippingAddress.PostalCode,
	}

	if err := uc.kafka.SendMessage(ctx, paymentMessage, "create_payment_topic"); err != nil {
		uc.log.WithError(err).Error("Failed to send payment message to Kafka")
		return nil, model.ErrInternalServer
	}

	if err := uc.kafka.SendMessage(ctx, shippingMessage, "create_shipping_topic"); err != nil {
		uc.log.WithError(err).Error("Failed to send shipping message to Kafka")
		return nil, model.ErrInternalServer
	}

	return converter.CreateOrderToResponse(paymentMessage, shippingMessage, order), nil
}


func (uc *OrderUseCase) GetOrdersByBuyer(ctx context.Context, request *model.SearchOrderRequest) ([]model.ListOrderResponse, int64, error) {
    if err := helper.ValidateStruct(uc.val, uc.log, request); err != nil {
        return nil, 0, err
    }
    tasks, total, err := uc.orderRepo.GetOrdersByBuyer(request)
    if err != nil {
        uc.log.WithError(err).Error("Failed to get orders")
        return nil, 0, model.ErrInternalServer
    }
    responses := make([]model.ListOrderResponse, len(tasks))
    for i, task := range tasks {
        responses[i] = *converter.OrdersToResponse(&task)
    }
    return responses, total, nil
}

func(uc *OrderUseCase) GetOrderByIdByBuyer(ctx context.Context, request *model.GetOrderDetails) (*model.OrderResponse, error) {
    if err := helper.ValidateStruct(uc.val, uc.log, request); err != nil {
        return nil, err
    }
    order, err := uc.orderRepo.GetOrderByIdByBuyer(request)
    if err != nil {
        uc.log.WithError(err).Error("Failed to get order")
        return nil, err
    }
    return converter.OrderToResponse(order), nil
}

func (uc *OrderUseCase) CancelOrder(ctx context.Context, request *model.CancelOrderRequest) (*model.OrderResponse, error) {
    tx := uc.db.WithContext(ctx).Begin()
    defer tx.Rollback()

    if err := helper.ValidateStruct(uc.val, uc.log, request); err != nil {
        return nil, err
    }
    
    order, err := uc.orderRepo.GetOrderByIdByBuyer(&model.GetOrderDetails{
        OrderUUID: request.OrderUUID,
        UserID:    request.UserID,
    })
    if err != nil {
        uc.log.WithError(err).Error("Failed to get order")
        return nil, err
    }

    if order.Status == "completed" || order.Status == "cancelled" || order.Status == "shipped" || order.Status == "processed" {
        return nil, model.NewApiError(fiber.StatusConflict, fmt.Sprintf("Order with ID %s cannot be cancelled, current status is %s", request.OrderUUID, order.Status), nil)
    }
    
    order.Status = "cancelled"
    if err := uc.orderRepo.UpdateOrder(order); err != nil {
        uc.log.WithError(err).Error("Failed to update order")
        return nil, model.ErrInternalServer
    }

    if order.Payment != nil {
        order.Payment.Status = "cancelled"
        if err := uc.paymentRepo.UpdatePayment(order.Payment); err != nil {
            uc.log.WithError(err).Error("Failed to update payment")
            return nil, model.ErrInternalServer
        }
    }

    if order.Shipping != nil {
        order.Shipping.Status = "cancelled"
        if err := uc.shippingRepo.UpdateShipping(order.Shipping); err != nil {
            uc.log.WithError(err).Error("Failed to update shipping")
            return nil, model.ErrInternalServer
        }
    }

    for _, item := range order.Items {
        product, err := uc.productRepo.FindProductByID(item.ProductID)
        if err != nil {
            uc.log.WithError(err).Error("Failed to get product")
            return nil, model.ErrInternalServer
        }
        product.Stock += item.Quantity
        if err := uc.productRepo.UpdateProduct(&product); err != nil {
            uc.log.WithError(err).Error("Failed to update product")
            return nil, model.ErrInternalServer
        }
    }
    
    if err := tx.Commit().Error; err != nil {
        uc.log.WithError(err).Error("Failed to commit transaction")
        return nil, model.ErrInternalServer
    }
    
    return converter.OrderToResponse(order), nil
}

func (uc *OrderUseCase) CheckoutOrder(ctx context.Context, request *model.CheckoutOrderRequest) (*model.OrderResponse, error) {
    tx := uc.db.WithContext(ctx).Begin()
    defer tx.Rollback()

    if err := helper.ValidateStruct(uc.val, uc.log, request); err != nil {
        return nil, err
    }

    order, err := uc.orderRepo.GetOrderByIdByBuyer(&model.GetOrderDetails{
        OrderUUID: request.OrderUUID,
        UserID:    request.UserID,
    })
    if err != nil {
        uc.log.WithError(err).Error("Failed to get order")
        return nil, err
    }

    if order.Status != "pending" {
        return nil, model.NewApiError(fiber.StatusConflict, fmt.Sprintf("Order with ID %s cannot be checked out, current status is %s", request.OrderUUID, order.Status), nil)
    }
    order.Status = "processed"
    if err := uc.orderRepo.UpdateOrder(order); err != nil {
        uc.log.WithError(err).Error("Failed to update order")
        return nil, model.ErrInternalServer
    }
    if order.Payment != nil {
        order.Payment.Status = "paid"
        if err := uc.paymentRepo.UpdatePayment(order.Payment); err != nil {
            uc.log.WithError(err).Error("Failed to update payment")
            return nil, model.ErrInternalServer
        }
    }

    if err := tx.Commit().Error; err != nil {
        uc.log.WithError(err).Error("Failed to commit transaction")
        return nil, model.ErrInternalServer
    }
    return converter.OrderToResponse(order), nil
}

func (u *OrderUseCase) GetOrdersBySeller(ctx context.Context, request *model.SearchOrderRequestBySeller) ([]model.OrdersResponseForSeller, int64, error) {
	if err := helper.ValidateStruct(u.val, u.log, request); err != nil {
		return nil, 0, err
	}
	store, err := u.storeRepo.FindStoreByUserID(request.UserID)
    if err != nil {
        return nil, 0, err
    }
	request.StoreID = store.ID
	orders, total, err := u.orderRepo.GetOrdersBySeller(request)
	if err != nil {
		u.log.WithError(err).Errorf("Failed to get orders")
		return nil, 0, err
	}
	responses := make([]model.OrdersResponseForSeller, len(orders))
	for i, order := range orders {
		responses[i] = *converter.OrderToResponseForSeller(&order)
	}
	return responses, total, nil
}

func (u *OrderUseCase) GetOrderBySeller(ctx context.Context, request *model.GetOrderDetails) (*model.OrderResponse, error) {
	if err := helper.ValidateStruct(u.val, u.log, request); err != nil {
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