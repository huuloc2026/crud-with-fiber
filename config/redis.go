package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisClient(cfg *RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		os.Exit(1)
	}

	return rdb
}
