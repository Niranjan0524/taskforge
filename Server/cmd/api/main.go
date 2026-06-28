package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Niranjan0524/taskforge/server/internals/Storage/redisStore"
	"github.com/Niranjan0524/taskforge/server/internals/config"
	"github.com/Niranjan0524/taskforge/server/internals/handlers"
	"github.com/Niranjan0524/taskforge/server/internals/handlers/webSockets"
	"github.com/Niranjan0524/taskforge/server/internals/worker"
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
	router.Use(corsMiddleware())
	// router.Use(gin.Recovery())

	fmt.Println("API server")

	ctx := context.Background()
	redisOptions, err := config.RedisOptionsFromEnv()
	if err != nil {
		log.Fatal("Invalid redis config", err)
	}
	rdb := redis.NewClient(redisOptions)
	store := redisStore.NewRedisStore(rdb)

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to redis", err)
	}

	log.Println("Connected to redis")
	go webSockets.StartTaskStatusSubscriber(ctx, rdb)
	workerPool := worker.NewWorkerPool(store, 2, rdb)
	go func() {
		if err := workerPool.Start(ctx); err != nil {
			log.Println("worker stopped:", err)
		}
	}()

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
	authorized.DELETE("/api/tasks/cancel/:id", handlers.CancelTask(store))
	router.GET("/ws", webSockets.WebSocketHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		allowedOrigin := os.Getenv("ORIGIN_URL")

		if allowedOrigin == "" || originAllowed(origin, allowedOrigin) {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Writer.Header().Set("Vary", "Origin")
		}

		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func originAllowed(origin string, allowedOrigins string) bool {
	if origin == "" {
		return true
	}

	for _, allowed := range strings.Split(allowedOrigins, ",") {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}

	return false
}
