package main

import (
	"fmt"

	"github.com/abdisetiakawan/go-ecommerce/internal/config"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
)

func main() {
    viperConfig := config.NewViper()

    logger := config.NewLogger(viperConfig)
    db := config.NewDatabase(viperConfig, logger)
    validator := config.NewValidator(viperConfig)
    jwt := helper.NewJWTHelper(viperConfig)
    uuid := helper.NewUUIDHelper()
    
    app := config.NewFiber(viperConfig)

    config.Bootstrap(&config.BootstrapConfig{DB: db, App: app, Log: logger, Validate: validator, Config: viperConfig, Jwt: jwt, UserUUID: uuid})

    port := viperConfig.GetInt("web.port")
    logger.Infof("Starting server on port %d", port)
    err := app.Listen(fmt.Sprintf(":%d", port))
    if err != nil {
        logger.Fatalf("Failed to start server: %v", err)
    }
}