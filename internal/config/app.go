package config

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/seller"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/user"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/route"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Jwt      *helper.JwtHelper
	UserUUID *helper.UUIDHelper
}

func Bootstrap(config *BootstrapConfig) {
	authRepository := repository.NewAuthRepository(config.Log)
	authUseCase := usecase.NewAuthUseCase(config.DB, config.Log, config.Validate, authRepository, config.Jwt, config.UserUUID)
	authController := http.NewAuthController(authUseCase, config.Log)

	sellerRepository := repository.NewSellerRepository(config.Log, config.DB)
	sellerUseCase := usecase.NewSellerUseCase(config.DB, config.Log, config.Validate, sellerRepository, config.UserUUID)
	sellerController := seller.NewSellerController(sellerUseCase, config.Log)

	userRepository := repository.NewUserRepository(config.Log)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)
	userController := user.NewUserController(userUseCase, config.Log)
	
	AuthMiddleware := middleware.NewAuth(config.Config)
	routeConfig := &route.RouteConfig{
		App: config.App,
		AuthController: authController,
		SellerController: sellerController,
		UserController: userController,
		AuthMiddleware: AuthMiddleware,
	}
	routeConfig.Setup()
}