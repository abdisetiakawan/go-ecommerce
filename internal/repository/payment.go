package repository

import (
	"context"
	"encoding/json"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	eventmodel "github.com/abdisetiakawan/go-ecommerce/internal/model/event_model"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	kafka *helper.KafkaConsumer
	DB *gorm.DB
}

func NewPaymentRepository(DB *gorm.DB, kafka *helper.KafkaConsumer) interfaces.PaymentRepository {
	return &PaymentRepository{kafka,DB}
}

func (r *PaymentRepository) CreatePayment() error {
    consumer, err := r.kafka.Consume(context.Background(), "create_payment_topic")
    if err != nil {
        return err
    }
    defer consumer.Close()

    for {
        select {
        case msg := <-consumer.Messages():
            var paymentMessage eventmodel.PaymentMessage
            err := json.Unmarshal(msg.Value, &paymentMessage)
            if err != nil {
                logrus.WithError(err).Error("Failed to unmarshal payment message")
                continue
            }

            payment := &entity.Payment{
                PaymentUUID: paymentMessage.PaymentUUID,
                Amount:      paymentMessage.Amount,
                Method:      paymentMessage.Method,
                OrderID:     paymentMessage.OrderID,
                Status:      paymentMessage.Status,
            }

            if err := r.DB.Create(payment).Error; err != nil {
                logrus.WithError(err).Error("Failed to create payment")
                continue
            }

        case err := <-consumer.Errors():
            logrus.WithError(err).Error("Failed to consume payment topic")
        }
    }
}

func (r *PaymentRepository) UpdatePayment(payment *entity.Payment) error {
	return r.DB.Save(payment).Error
}

func (r *PaymentRepository) CancelPayment() error {
    consumer, err := r.kafka.Consume(context.Background(), "cancel_payment_topic")
    if err != nil {
        return err
    }
    defer consumer.Close()

    for {
        select {
        case msg := <-consumer.Messages():
            var paymentMessage eventmodel.PaymentMessage
            err := json.Unmarshal(msg.Value, &paymentMessage)
            if err != nil {
                logrus.WithError(err).Error("Failed to unmarshal payment message")
                continue
            }

            if err := r.DB.Model(&entity.Payment{}).
                Where("order_id = ?", paymentMessage.OrderID).
                Update("status", "cancelled").Error; err != nil {
                logrus.WithError(err).Error("Failed to update order status")
            }
        

        case err := <-consumer.Errors():
            logrus.WithError(err).Error("Failed to consume payment topic")
        }
    }
}

func (r *PaymentRepository) CheckoutPayment() error {
    consumer, err := r.kafka.Consume(context.Background(), "checkout_payment_topic")
    if err != nil {
        return err
    }
    defer consumer.Close()

    for {
        select {
        case msg := <-consumer.Messages():
            var paymentMessage eventmodel.PaymentMessage
            err := json.Unmarshal(msg.Value, &paymentMessage)
            if err != nil {
                logrus.WithError(err).Error("Failed to unmarshal payment message")
                continue
            }

            if err := r.DB.Model(&entity.Payment{}).
                Where("order_id = ?", paymentMessage.OrderID).
                Update("status", paymentMessage.Status).Error; err != nil {
                logrus.WithError(err).Error("Failed to update order status")
            }
        

        case err := <-consumer.Errors():
            logrus.WithError(err).Error("Failed to consume payment topic")
        }
    }
}