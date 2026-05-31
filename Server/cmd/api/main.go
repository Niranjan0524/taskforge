package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Niranjan0524/taskforge/server/internals/handlers"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	fmt.Println("API server")

	ctx := context.Background()
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to redis", err)
	}

	log.Println("Connected to redis")

	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message":  "Working",
			"ClientIp": ctx.ClientIP(),
		})
	})

	authorized := router.Group("/", handlers.ValidateUser())

	authorized.POST("/api/task", handlers.CreateTask(rdb))
	authorized.GET("/api/task/:id", handlers.GetTask(rdb))

	if err := router.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
