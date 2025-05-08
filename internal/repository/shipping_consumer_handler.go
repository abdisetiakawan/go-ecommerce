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

type ShippingConsumerHandler struct {
    db    *gorm.DB
    kafka *helper.KafkaConsumer
}

func NewShippingConsumerHandler(db *gorm.DB, kafka *helper.KafkaConsumer) interfaces.ShippingConsumer {
    return &ShippingConsumerHandler{db: db, kafka: kafka}
}

func (h *ShippingConsumerHandler) CreateShipping(ctx context.Context) error {
    return h.kafka.Consume(ctx, []string{"create_shipping_topic"}, h)
}

func (h *ShippingConsumerHandler) CancelShipping(ctx context.Context) error {
    return h.kafka.Consume(ctx, []string{"cancel_shipping_topic"}, h)
}

func (h *ShippingConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ShippingConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ShippingConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var shippingMessage eventmodel.ShippingMessage
        if err := json.Unmarshal(msg.Value, &shippingMessage); err != nil {
            logrus.WithError(err).Error("Failed to unmarshal shipping message")
            continue
        }

        switch msg.Topic {
        case "create_shipping_topic":
            shipping := &entity.Shipping{
                ShippingUUID: shippingMessage.ShippingUUID,
                OrderID:      shippingMessage.OrderID,
                Address:      shippingMessage.Address,
                City:        shippingMessage.City,
                Province:    shippingMessage.Province,
                PostalCode:  shippingMessage.PostalCode,
                Status:      shippingMessage.Status,
            }
            if err := h.db.Create(shipping).Error; err != nil {
                logrus.WithError(err).Error("Failed to create shipping")
                continue
            }

        case "cancel_shipping_topic":
            if err := h.db.Model(&entity.Shipping{}).
                Where("order_id = ?", shippingMessage.OrderID).
                Update("status", "cancelled").Error; err != nil {
                logrus.WithError(err).Error("Failed to cancel shipping")
                continue
            }
        }

        session.MarkMessage(msg, "")
    }
    return nil
}