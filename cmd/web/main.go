package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abdisetiakawan/go-ecommerce/internal/config"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

//go:embed docs/swagger.yaml
var swaggerYAML []byte

func main() {
	viperConfig := config.NewViper()

	logger := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, logger)
	validator := config.NewValidator()
	jwt := helper.NewJWTHelper(viperConfig)
	uuid := helper.NewUUIDHelper()

	kafkaProducer, err := helper.NewKafkaProducer(viperConfig)
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}

	app := config.NewFiber(viperConfig)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
		MaxAge:           3600,
	}))
	app.Get("/swagger.yaml", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/yaml")
		return c.Send(swaggerYAML)
	})

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger.yaml",
	}))

	config.Bootstrap(&config.BootstrapConfig{DB: db, App: app, Validate: validator, Config: viperConfig, Jwt: jwt, UserUUID: uuid, KafkaProducer: kafkaProducer})
	port := viperConfig.GetInt("PORT")
	logger.Infof("Starting server on port %d", port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Shutting down server...")

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Create error channel to collect shutdown errors
		shutdownErrors := make(chan error, 2)

		// Shutdown components concurrently
		go func() { shutdownErrors <- app.Shutdown() }()
		go func() { shutdownErrors <- kafkaProducer.Close() }()

		// Collect shutdown errors
		var shutdownError error
		for i := 0; i < 2; i++ {
			if err := <-shutdownErrors; err != nil {
				logger.Error(err)
				shutdownError = err
			}
		}

		// Check if shutdown completed within timeout
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				logger.Error("Shutdown timed out")
				os.Exit(1)
			}
		default:
			if shutdownError != nil {
				logger.Error("Error during shutdown")
				os.Exit(1)
			}
			logger.Info("Server gracefully stopped")
		}
	}()

	err = app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
