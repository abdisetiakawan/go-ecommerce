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
	PaymentConsumer *helper.KafkaConsumer
	ShippingConsumer *helper.KafkaConsumer
	OrderConsumer *helper.KafkaConsumer
}

func Bootstrap(config *BootstrapConfig) {
    // Event
    orderEventRepo := eventrepository.NewOrderEventRepository(config.DB)
    orderEventUC := eventuc.NewOrderEventEvent(config.DB, orderEventRepo, config.KafkaProducer)

    // Create separate consumers for each topic/operation
    createPaymentConsumer, err := helper.NewKafkaConsumer(config.Config, "payment-create-consumer")
    if err != nil {
        log.Fatal(err)
        panic(err)
    }
    
    cancelPaymentConsumer, err := helper.NewKafkaConsumer(config.Config, "payment-cancel-consumer")
    if err != nil {
        log.Fatal(err)
        panic(err)
    }
    
    checkoutPaymentConsumer, err := helper.NewKafkaConsumer(config.Config, "payment-checkout-consumer") 
    if err != nil {
        log.Fatal(err)
        panic(err)
    }
    
    createShippingConsumer, err := helper.NewKafkaConsumer(config.Config, "shipping-create-consumer")
    if err != nil {
        log.Fatal(err)
        panic(err)
    }
    
    cancelShippingConsumer, err := helper.NewKafkaConsumer(config.Config, "shipping-cancel-consumer")
    if err != nil {
        log.Fatal(err)
        panic(err)
    }
    
    orderStatusConsumer, err := helper.NewKafkaConsumer(config.Config, "order-status-consumer")
    if err != nil {
        log.Fatal(err)
        panic(err)
    }

    // Create repositories with dedicated consumers
    createPaymentRepo := repository.NewPaymentConsumerHandler(config.DB, createPaymentConsumer)
    cancelPaymentRepo := repository.NewPaymentConsumerHandler(config.DB, cancelPaymentConsumer)
    checkoutPaymentRepo := repository.NewPaymentConsumerHandler(config.DB, checkoutPaymentConsumer)
    
    createShippingRepo := repository.NewShippingConsumerHandler(config.DB, createShippingConsumer)
    cancelShippingRepo := repository.NewShippingConsumerHandler(config.DB, cancelShippingConsumer)
    
    orderStatusRepo := repository.NewOrderConsumerHandler(config.DB, orderStatusConsumer)

    // Other repositories and usecases
    userRepository := repository.NewUserRepository(config.DB)
    profileRepository := repository.NewProfileRepository(config.DB)
    orderRepository := repository.NewOrderRepository(config.DB)
    productRepository := repository.NewProductRepository(config.DB)
    storeRepository := repository.NewStoreRepository(config.DB)
    shippingRepository := repository.NewShippingRepository(config.DB)

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

    // Start consumer goroutines with dedicated consumers
    go func() {
        ctx := context.Background()
        if err := createPaymentRepo.CreatePayment(ctx); err != nil {
            log.Error(err)
        }
    }()
    
    go func() {
        ctx := context.Background()
        if err := createShippingRepo.CreateShipping(ctx); err != nil {
            log.Error(err)
        }
    }()
    
    go func() {
        ctx := context.Background()
        if err := cancelPaymentRepo.CancelPayment(ctx); err != nil {
            log.Error(err)
        }
    }()
    
    go func() {
        ctx := context.Background()
        if err := cancelShippingRepo.CancelShipping(ctx); err != nil {
            log.Error(err)
        }
    }()
    
    go func() {
        ctx := context.Background()
        if err := checkoutPaymentRepo.CheckoutPayment(ctx); err != nil {
            log.Error(err)
        }
    }()
    
    go func() {
        ctx := context.Background()
        if err := orderStatusRepo.ChangeOrderStatus(ctx); err != nil {
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