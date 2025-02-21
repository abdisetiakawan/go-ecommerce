package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func NewRedisClient(viper *viper.Viper) *redis.Client{
	addr := viper.GetString("REDIS_ADDR")
	password := viper.GetString("REDIS_PASS")
	db := viper.GetInt("REDIS_DB")

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