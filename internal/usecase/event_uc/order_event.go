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

// NewOrderEventEvent creates a new instance of OrderEventUseCase.
// It takes database connection, event repository and kafka producer as an argument.
// The returned OrderEventUseCase is ready to use and contains all the necessary dependencies.

func NewOrderEventEvent(db *gorm.DB, eventRepo eventRepo.OrderEventRepository, kafka *helper.KafkaProducer) interfaces.OrderEventUseCase {
	return &OrderEventUseCase{
		db:        db,
		eventRepo: eventRepo,
		kafka:     kafka,
	}
}

// ProcessOrderEvent is a function that takes a context and an event entity as arguments.
// It attempts to process the payment and shipping message by sending the message to kafka topic.
// If the message sending process fails, it will retry up to 3 times with 2 seconds delay between each attempt.
// If the retry also fails, it will update the status of the event entity to "failed" and return an error.
// If the message sending process is successful, it will update the status of the event entity to "completed".
// The function will return an error if the event entity status is not "pending" or if there is an error when updating the event entity.
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

// RetryFailedEvents retrieves all order events with a "failed" status from the past 24 hours
// and attempts to reprocess each event concurrently. For each event, it calls the
// ProcessOrderEvent function. If an error occurs during processing, it logs the error
// message. The function returns an error if the retrieval of failed events fails.

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

// CancelOrderEvent takes a context and an event entity as arguments.
// It attempts to cancel the payment and shipping message by sending the message to kafka topic.
// If the message sending process fails, it will retry up to 3 times with 2 seconds delay between each attempt.
// If the retry also fails, it will update the status of the event entity to "failed" and return an error.
// If the message sending process is successful, it will update the status of the event entity to "completed".
// The function will return an error if the event entity status is not "pending" or if there is an error when updating the event entity.
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

// CheckoutOrderEvent takes a context and an event entity as arguments.
// It attempts to send a message to kafka topic to mark the order as "paid".
// If the message sending process fails, it will retry up to 3 times with 2 seconds delay between each attempt.
// If the retry also fails, it will update the status of the event entity to "failed" and return an error.
// If the message sending process is successful, it will update the status of the event entity to "completed".
// The function will return an error if there is an error when updating the event entity.
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

// ChangeOrderStatusUC takes a context and an event entity as arguments.
// It attempts to send a message to kafka topic to change the order status.
// If the message sending process fails, it will retry up to 3 times with 2 seconds delay between each attempt.
// If the retry also fails, it will update the status of the event entity to "failed" and return an error.
// If the message sending process is successful, it will update the status of the event entity to "completed".
// The function will return an error if there is an error when updating the event entity.
func (uc *OrderEventUseCase) ChangeOrderStatusUC(ctx context.Context,event *evententity.OrderEvent) error {
	var orderStatus eventmodel.OrderMessage
	if err := json.Unmarshal(event.OrderData, &orderStatus); err != nil {
		return err
	}

	orderMessage := &eventmodel.OrderMessage{
		OrderID:     event.OrderID,
		Status:      orderStatus.Status,
	}

	err := retry.Do(func() error {
		return uc.kafka.SendMessage(ctx, orderMessage, "change_order_topic")
	}, retry.Attempts(3), retry.Delay(2*time.Second))

	if err != nil {
		event.Status = "failed"
		event.Error = fmt.Sprintf("Order processing failed: %v", err)
		uc.eventRepo.UpdateOrderEvent(event)
		return err
	}

	event.Status = "completed"
	return uc.eventRepo.UpdateOrderEvent(event)
}