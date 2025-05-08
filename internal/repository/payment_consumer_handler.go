package repository

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	eventmodel "github.com/abdisetiakawan/go-ecommerce/internal/model/event_model"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentConsumerHandler struct {
    db    *gorm.DB
    kafka *helper.KafkaConsumer
}

func NewPaymentConsumerHandler(db *gorm.DB, kafka *helper.KafkaConsumer) interfaces.PaymentConsumer {
    return &PaymentConsumerHandler{db: db, kafka: kafka}
}

func (h *PaymentConsumerHandler) CreatePayment(ctx context.Context) error {
    return h.kafka.Consume(ctx, []string{"create_payment_topic"}, h)
}

func (h *PaymentConsumerHandler) CancelPayment(ctx context.Context) error {
    return h.kafka.Consume(ctx, []string{"cancel_payment_topic"}, h)
}

func (h *PaymentConsumerHandler) CheckoutPayment(ctx context.Context) error {
    return h.kafka.Consume(ctx, []string{"checkout_payment_topic"}, h)
}

// Sarama ConsumerGroupHandler interface implementation
func (h *PaymentConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *PaymentConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *PaymentConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var paymentMessage eventmodel.PaymentMessage
        if err := json.Unmarshal(msg.Value, &paymentMessage); err != nil {
            logrus.WithError(err).Error("Failed to unmarshal payment message")
            continue
        }

        switch msg.Topic {
        case "create_payment_topic":
            payment := &entity.Payment{
                PaymentUUID: paymentMessage.PaymentUUID,
                Amount:      paymentMessage.Amount,
                Method:      paymentMessage.Method,
                OrderID:     paymentMessage.OrderID,
                Status:      paymentMessage.Status,
            }
            if err := h.db.Create(payment).Error; err != nil {
                logrus.WithError(err).Error("Failed to create payment")
                continue
            }

        case "cancel_payment_topic":
            if err := h.db.Model(&entity.Payment{}).
                Where("order_id = ?", paymentMessage.OrderID).
                Update("status", "cancelled").Error; err != nil {
                logrus.WithError(err).Error("Failed to cancel payment")
                continue
            }

        case "checkout_payment_topic":
            if err := h.db.Model(&entity.Payment{}).
                Where("order_id = ?", paymentMessage.OrderID).
                Update("status", paymentMessage.Status).Error; err != nil {
                logrus.WithError(err).Error("Failed to checkout payment")
                continue
            }
        }

        session.MarkMessage(msg, "")
    }
    return nil
}