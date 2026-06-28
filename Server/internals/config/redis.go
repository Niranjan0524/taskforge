package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func RedisOptionsFromEnv() (*redis.Options, error) {
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		return redis.ParseURL(redisURL)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	return &redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	}, nil
}
