package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Niranjan0524/taskforge/server/internals/Storage/redisStore"
	"github.com/Niranjan0524/taskforge/server/internals/config"
	"github.com/Niranjan0524/taskforge/server/internals/worker"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	redisOptions, err := config.RedisOptionsFromEnv()
	if err != nil {
		log.Fatal("Invalid redis config", err)
	}

	rdb := redis.NewClient(redisOptions)
	defer rdb.Close()

	store := redisStore.NewRedisStore(rdb)

	workerPool := worker.NewWorkerPool(store, 2, rdb)

	log.Println("TaskForge worker started")

	if err := workerPool.Start(ctx); err != nil {
		log.Println("worker stopped:", err)
	}

}
