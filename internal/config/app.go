package config

import (
	"context"
	"time"

	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/route"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository"
	eventrepository "github.com/abdisetiakawan/go-ecommerce/internal/repository/event_repository"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase"
	eventuc "github.com/abdisetiakawan/go-ecommerce/internal/usecase/event_uc"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Validate *validator.Validate
	Config   *viper.Viper
	Jwt      *helper.JwtHelper
	UserUUID *helper.UUIDHelper
	KafkaProducer *helper.KafkaProducer
	KafkaConsumer *helper.KafkaConsumer
}

func Bootstrap(config *BootstrapConfig) {
	// Event
	orderEventRepo := eventrepository.NewOrderEventRepository(config.DB)
	orderEventUC := eventuc.NewOrderEventEvent(config.DB, orderEventRepo, config.KafkaProducer)

	userRepository := repository.NewUserRepository(config.DB)
	profileRepository := repository.NewProfileRepository(config.DB)
	orderRepository := repository.NewOrderRepository(config.DB, config.KafkaConsumer)
	productRepository := repository.NewProductRepository(config.DB)
	storeRepository := repository.NewStoreRepository(config.DB)
	shippingRepository := repository.NewShippingRepository(config.DB, config.KafkaConsumer)
	paymentRepository := repository.NewPaymentRepository(config.DB, config.KafkaConsumer)

	userUseCase := usecase.NewUserUseCase(config.DB, config.Validate, userRepository, config.UserUUID, config.Jwt)
	profileUseCase := usecase.NewProfileUseCase(config.DB, config.Validate, profileRepository)
	orderUseCase := usecase.NewOrderUseCase(config.DB, config.Validate, orderRepository, productRepository, storeRepository, config.UserUUID, orderEventUC)
	productUseCase := usecase.NewProductUseCase(config.DB, config.Validate, productRepository, storeRepository, config.UserUUID)
	storeUseCase := usecase.NewStoreUseCase(config.DB, config.Validate, storeRepository, config.UserUUID)
	shippingUseCase := usecase.NewShippingUseCase(config.DB, config.Validate, shippingRepository, storeRepository, orderRepository, config.UserUUID, orderEventUC)
	
	userController := http.NewUserController(userUseCase)
	profileController := http.NewProfileController(profileUseCase)
	orderController := http.NewOrderController(orderUseCase)
	productController := http.NewProductController(productUseCase)
	storeController := http.NewStoreController(storeUseCase)
	shippingController := http.NewShippingController(shippingUseCase)

	go func() {
		if err := paymentRepository.CreatePayment(); err != nil {
			log.Error(err)
		}
	}()
	go func() {
		if err := shippingRepository.CreateShipping(); err != nil {
			log.Error(err)
		}
	}()
	go func() {
		if err := paymentRepository.CancelPayment(); err != nil {
			log.Error(err)
		}
	}()
	go func() {
		if err := shippingRepository.CancelShipping(); err != nil {
			log.Error(err)
		}
	}()
	go func() {
		if err := paymentRepository.CheckoutPayment(); err != nil {
			log.Error(err)
		}
	}()
	go func() {
		if err := orderRepository.ChangeOrderStatus(); err != nil {
			log.Error(err)
		}
	}()
	go func() {
        ticker := time.NewTicker(5 * time.Minute)
        for range ticker.C {
            if err := orderEventUC.RetryFailedEvents(context.Background()); err != nil {
                logrus.WithError(err).Error("Failed to retry events")
            }
        }
    }()

	AuthMiddleware := middleware.NewAuth(config.Config)
	routeConfig := &route.RouteConfig{
		App: config.App,
		ProfileController: profileController,
		UserController: userController,
		OrderController: orderController,
		ProductController: productController,
		StoreController: storeController,
		ShippingController: shippingController,
		AuthMiddleware: AuthMiddleware,
	}
	routeConfig.Setup()
}