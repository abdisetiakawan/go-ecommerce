package config

import (
	"log"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("")

	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err != nil {
		log.Println("info: .env file not found, using environment variables")
	}
	return v
}
