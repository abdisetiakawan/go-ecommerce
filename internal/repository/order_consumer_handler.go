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

// OrderConsumerHandler implement interfaces.OrderConsumer dan sarama.ConsumerGroupHandler
type OrderConsumerHandler struct {
	db    *gorm.DB
	kafka *helper.KafkaConsumer
}

// Konstruktor
func NewOrderConsumerHandler(db *gorm.DB, kafka *helper.KafkaConsumer) interfaces.OrderConsumer {
	return &OrderConsumerHandler{db: db, kafka: kafka}
}

// ChangeOrderStatus mulai konsumsi topic
func (h *OrderConsumerHandler) ChangeOrderStatus(ctx context.Context) error {
	topics := []string{"change_order_topic"}
	return h.kafka.Consume(ctx, topics, h)
}

// --- Sarama ConsumerGroupHandler methods ---
func (h *OrderConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *OrderConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *OrderConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var m eventmodel.OrderMessage
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			logrus.WithError(err).Error("Failed to unmarshal order message")
			continue
		}
		if err := h.db.Model(&entity.Order{}).
			Where("id = ?", m.OrderID).
			Update("status", m.Status).Error; err != nil {
			logrus.WithError(err).Error("Failed to update order status")
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
