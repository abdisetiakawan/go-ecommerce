package config

import (
	"fmt"

	"github.com/spf13/viper"
)



func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error to load config.json %w", err))
	}
	return config
}