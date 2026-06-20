package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Niranjan0524/taskforge/server/internals/Storage/redisStore"
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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer rdb.Close()

	store := redisStore.NewRedisStore(rdb)

	workerPool := worker.NewWorkerPool(store, 2, rdb)

	log.Println("TaskForge worker started")

	if err := workerPool.Start(ctx); err != nil {
		log.Println("worker stopped:", err)
	}

}
