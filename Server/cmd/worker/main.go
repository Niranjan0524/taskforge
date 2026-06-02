package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Niranjan0524/taskforge/server/internals/Storage/redisStore"
	"github.com/Niranjan0524/taskforge/server/internals/worker"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	store := redisStore.NewRedisStore(rdb)

	workerPool := worker.NewWorkerPool(store, 3)

	log.Println("TaskForge worker started")

	if err := workerPool.Start(ctx); err != nil {
		log.Println("worker stopped:", err)
	}

}
