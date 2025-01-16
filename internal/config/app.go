package config

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http"
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


	routeConfig := &route.RouteConfig{
		App: config.App,
		AuthController: authController,
	}
	routeConfig.Setup()
}