package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("")

	// Set config type and name
	v.SetConfigType("env")
	v.SetConfigName(".env")

	// Add config path
	v.AddConfigPath(".")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Set default values
	v.SetDefault("KAFKA_CONSUMER_GROUP", "ecommerce-group")
	v.SetDefault("KAFKA_CLIENT_ID", "ecommerce-service")
	v.SetDefault("KAFKA_VERSION", "2.8.1")
	v.SetDefault("KAFKA_AUTO_OFFSET_RESET", "latest")
	v.SetDefault("KAFKA_RETRY_MAX", 3)
	v.SetDefault("KAFKA_REQUEST_TIMEOUT_MS", 5000)

	return v
}
