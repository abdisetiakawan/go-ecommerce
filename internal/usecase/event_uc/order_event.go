package eventuc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	evententity "github.com/abdisetiakawan/go-ecommerce/internal/entity/event_entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	eventmodel "github.com/abdisetiakawan/go-ecommerce/internal/model/event_model"
	eventRepo "github.com/abdisetiakawan/go-ecommerce/internal/repository/event_repository/interfaces"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/event_uc/interfaces"
	"github.com/avast/retry-go"
	"gorm.io/gorm"
)

type OrderEventUseCase struct {
	db        *gorm.DB
	eventRepo eventRepo.OrderEventRepository
	kafka     *helper.KafkaProducer
}

func NewOrderEventEvent(db *gorm.DB, eventRepo eventRepo.OrderEventRepository, kafka *helper.KafkaProducer) interfaces.OrderEventUseCase {
	return &OrderEventUseCase{
		db:        db,
		eventRepo: eventRepo,
		kafka:     kafka,
	}
}

func (uc *OrderEventUseCase) ProcessOrderEvent(ctx context.Context, event *evententity.OrderEvent) error {
	var paymentData eventmodel.PaymentMessage
	var shippingData eventmodel.ShippingMessage

	if err := json.Unmarshal(event.PaymentData, &paymentData); err != nil {
		return err
	}
	if err := json.Unmarshal(event.ShippingData, &shippingData); err != nil {
		return err
	}

	paymentMessage := &eventmodel.PaymentMessage{
		PaymentUUID: paymentData.PaymentUUID,
		OrderID:     event.OrderID,
		Amount:      paymentData.Amount,
		Method:      paymentData.Method,
		Status:      paymentData.Status,
	}

	err := retry.Do(func() error {
		return uc.kafka.SendMessage(ctx, paymentMessage, "create_payment_topic")
	}, retry.Attempts(3), retry.Delay(2*time.Second))

	if err != nil {
		event.Status = "failed"
		event.Error = fmt.Sprintf("Payment processing failed: %v", err)
		uc.eventRepo.UpdateOrderEvent(event)
		return err
	}

	shippingMessage := &eventmodel.ShippingMessage{
		ShippingUUID: shippingData.ShippingUUID,
		OrderID:      event.OrderID,
		Address:      shippingData.Address,
		City:         shippingData.City,
		Province:     shippingData.Province,
		PostalCode:   shippingData.PostalCode,
		Status:       shippingData.Status,
	}

	err = retry.Do(func() error {
		return uc.kafka.SendMessage(ctx, shippingMessage, "create_shipping_topic")
	}, retry.Attempts(3), retry.Delay(2*time.Second))

	if err != nil {
		event.Status = "failed"
		event.Error = fmt.Sprintf("Shipping processing failed: %v", err)
		uc.eventRepo.UpdateOrderEvent(event)
		return err
	}

	event.Status = "completed"
	return uc.eventRepo.UpdateOrderEvent(event)
}

func (uc *OrderEventUseCase) RetryFailedEvents(ctx context.Context) error {
	events, err := uc.eventRepo.GetFailedEvents(24 * time.Hour)
	if err != nil {
		return err
	}

	for _, event := range events {
		go func(e evententity.OrderEvent) {
			if err := uc.ProcessOrderEvent(ctx, &e); err != nil {
				log.Println(err)
			}
		}(event)
	}
	return nil
}

func (uc *OrderEventUseCase) CancelOrderEvent(ctx context.Context, event *evententity.OrderEvent) error {
	var paymentStatus eventmodel.PaymentMessage
	var shippingStatus eventmodel.ShippingMessage

	if err := json.Unmarshal(event.PaymentData, &paymentStatus); err != nil {
		return err
	}
	if err := json.Unmarshal(event.ShippingData, &shippingStatus); err != nil {
		return err
	}

	paymentMessage := &eventmodel.PaymentMessage{
		OrderID:     event.OrderID,
		Status:      paymentStatus.Status,
	}

	err := retry.Do(func() error {
		return uc.kafka.SendMessage(ctx, paymentMessage, "cancel_payment_topic")
	}, retry.Attempts(3), retry.Delay(2*time.Second))

	if err != nil {
		event.Status = "failed"
		event.Error = fmt.Sprintf("Payment processing failed: %v", err)
		uc.eventRepo.UpdateOrderEvent(event)
		return err
	}

	shippingMessage := &eventmodel.ShippingMessage{
		OrderID:     event.OrderID,
		Status:       shippingStatus.Status,
	}

	err = retry.Do(func() error {
		return uc.kafka.SendMessage(ctx, shippingMessage, "cancel_shipping_topic")
	}, retry.Attempts(3), retry.Delay(2*time.Second))

	if err != nil {
		event.Status = "failed"
		event.Error = fmt.Sprintf("Shipping processing failed: %v", err)
		uc.eventRepo.UpdateOrderEvent(event)
		return err
	}

	event.Status = "completed"
	return uc.eventRepo.UpdateOrderEvent(event)
}

func (uc *OrderEventUseCase) CheckoutOrderEvent(ctx context.Context, event *evententity.OrderEvent) error {
	var paymentStatus eventmodel.PaymentMessage
	
	if err := json.Unmarshal(event.PaymentData, &paymentStatus); err != nil {
		return err
	}

	paymentMessage := &eventmodel.PaymentMessage{
		OrderID:     event.OrderID,
		Status:      paymentStatus.Status,
	}

	err := retry.Do(func() error {
		return uc.kafka.SendMessage(ctx, paymentMessage, "checkout_payment_topic")
	}, retry.Attempts(3), retry.Delay(2*time.Second))

	if err != nil {
		event.Status = "failed"
		event.Error = fmt.Sprintf("Payment processing failed: %v", err)
		uc.eventRepo.UpdateOrderEvent(event)
		return err
	}

	event.Status = "completed"
	return uc.eventRepo.UpdateOrderEvent(event)
}