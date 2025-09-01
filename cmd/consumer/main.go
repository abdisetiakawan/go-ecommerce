package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/abdisetiakawan/go-ecommerce/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, logger)

	consumers, err := config.BootstrapConsumers(&config.ConsumerBootstrapConfig{
		DB:     db,
		Config: viperConfig,
	})
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down consumer...")

	for _, c := range consumers {
		if err := c.Close(); err != nil {
			logger.Error(err)
		}
	}
	logger.Info("Consumer stopped")
}
