package config

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ConsumerBootstrapConfig struct {
	DB     *gorm.DB
	Config *viper.Viper
}

func BootstrapConsumers(cfg *ConsumerBootstrapConfig) ([]*helper.KafkaConsumer, error) {
	createPaymentConsumer, err := helper.NewKafkaConsumer(cfg.Config, "payment-create-consumer")
	if err != nil {
		return nil, err
	}
	cancelPaymentConsumer, err := helper.NewKafkaConsumer(cfg.Config, "payment-cancel-consumer")
	if err != nil {
		return nil, err
	}
	checkoutPaymentConsumer, err := helper.NewKafkaConsumer(cfg.Config, "payment-checkout-consumer")
	if err != nil {
		return nil, err
	}
	createShippingConsumer, err := helper.NewKafkaConsumer(cfg.Config, "shipping-create-consumer")
	if err != nil {
		return nil, err
	}
	cancelShippingConsumer, err := helper.NewKafkaConsumer(cfg.Config, "shipping-cancel-consumer")
	if err != nil {
		return nil, err
	}
	orderStatusConsumer, err := helper.NewKafkaConsumer(cfg.Config, "order-status-consumer")
	if err != nil {
		return nil, err
	}

	createPaymentRepo := repository.NewPaymentConsumerHandler(cfg.DB, createPaymentConsumer)
	cancelPaymentRepo := repository.NewPaymentConsumerHandler(cfg.DB, cancelPaymentConsumer)
	checkoutPaymentRepo := repository.NewPaymentConsumerHandler(cfg.DB, checkoutPaymentConsumer)
	createShippingRepo := repository.NewShippingConsumerHandler(cfg.DB, createShippingConsumer)
	cancelShippingRepo := repository.NewShippingConsumerHandler(cfg.DB, cancelShippingConsumer)
	orderStatusRepo := repository.NewOrderConsumerHandler(cfg.DB, orderStatusConsumer)

	go func() {
		ctx := context.Background()
		if err := createPaymentRepo.CreatePayment(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		ctx := context.Background()
		if err := cancelPaymentRepo.CancelPayment(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		ctx := context.Background()
		if err := checkoutPaymentRepo.CheckoutPayment(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		ctx := context.Background()
		if err := createShippingRepo.CreateShipping(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		ctx := context.Background()
		if err := cancelShippingRepo.CancelShipping(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		ctx := context.Background()
		if err := orderStatusRepo.ChangeOrderStatus(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	return []*helper.KafkaConsumer{
		createPaymentConsumer,
		cancelPaymentConsumer,
		checkoutPaymentConsumer,
		createShippingConsumer,
		cancelShippingConsumer,
		orderStatusConsumer,
	}, nil
}
