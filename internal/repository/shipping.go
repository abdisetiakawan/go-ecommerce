package repository

import (
	"context"
	"encoding/json"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/event"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ShippingRepository struct {
	DB *gorm.DB
	kafka *helper.KafkaConsumer
}

func NewShippingRepository(DB *gorm.DB, kafka *helper.KafkaConsumer) interfaces.ShippingRepository {
	return &ShippingRepository{DB, kafka}
}

func (r *ShippingRepository) CreateShipping() error {
	consumer, err := r.kafka.Consume(context.Background(), "create_shipping_topic")
	if err != nil {
		logrus.WithError(err).Error("Failed to consume shipping topic")
		return err
	}
	defer consumer.Close()

	for {
		select {
		case msg := <-consumer.Messages():
			var shippingMessage event.ShippingMessage
			err := json.Unmarshal(msg.Value, &shippingMessage)
			if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal shipping message")
				continue
			}

			shipping := &entity.Shipping{
				ShippingUUID: shippingMessage.ShippingUUID,
				OrderID: shippingMessage.OrderID,
				Address: shippingMessage.Address,
				City: shippingMessage.City,
				Province: shippingMessage.Province,
				PostalCode: shippingMessage.PostalCode,
				Status: shippingMessage.Status,
			}

			if err := r.DB.Create(shipping).Error; err != nil {
				logrus.WithError(err).Error("Failed to create shipping")
				continue
			}

		case err := <-consumer.Errors():
			logrus.WithError(err).Error("Failed to consume shipping topic")
		}
	}
}

func (r *ShippingRepository) UpdateShipping(shipping *entity.Shipping) error {
	return r.DB.Save(shipping).Error
}