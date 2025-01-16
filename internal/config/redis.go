package config

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedisClient(viper *viper.Viper, log *logrus.Logger) *redis.Client{
	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return client
}