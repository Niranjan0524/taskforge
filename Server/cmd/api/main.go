package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Niranjan0524/taskforge/server/internals/Storage/redisStore"
	"github.com/Niranjan0524/taskforge/server/internals/handlers"
	"github.com/Niranjan0524/taskforge/server/internals/handlers/webSockets"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	go webSockets.WsHub.Run()
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
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	store := redisStore.NewRedisStore(rdb)

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to redis", err)
	}

	log.Println("Connected to redis")
	go webSockets.StartTaskStatusSubscriber(ctx, rdb)

	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message":  "Working",
			"ClientIp": ctx.ClientIP(),
		})
	})

	authorized := router.Group("/", handlers.ValidateUser())

	authorized.POST("/api/task", handlers.CreateTask(store))
	authorized.GET("/api/task/:id/status", handlers.GetTaskStatus(store))
	authorized.GET("/api/task/:id", handlers.GetTask(store))
	authorized.GET("/api/tasks", handlers.GetAllTasks(store))
	authorized.DELETE("/api/tasks/:id", handlers.DeleteTask(store))
	router.GET("/ws", webSockets.WebSocketHandler)

	if err := router.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
