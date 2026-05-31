package handlers

import (
	"fmt"
	"net/http"
	"time"

	storage "github.com/Niranjan0524/taskforge/server/internals/Storage"
	redisStore "github.com/Niranjan0524/taskforge/server/internals/Storage/redisStore"
	"github.com/Niranjan0524/taskforge/server/internals/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func CreateTask(redis *redis.Client) gin.HandlerFunc {
	store := redisStore.NewRedisStore(redis)

	return func(ctx *gin.Context) {
		fmt.Println("Creates task")

		var req types.CreateTaskRequest

		err := ctx.ShouldBindJSON(&req)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request Body",
			})
			return
		}

		newTask := storage.Task{
			ID:         uuid.New().String(),
			Type:       req.Type,
			Payload:    req.Payload,
			Priority:   req.Priority,
			Status:     "pending",
			RetryCount: 0,
			MaxRetries: req.MaxRetries,
			CreatedAt:  time.Now().UTC(),
		}

		if newTask.MaxRetries == 0 {
			newTask.MaxRetries = 3
		}

		err2 := store.CreateTask(ctx.Request.Context(), newTask)
		if err2 != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create task",
			})
			return
		}

		ctx.JSON(http.StatusCreated, newTask)
	}
}

func GetTask(redis *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("id")
		fmt.Println("userId", userId)
	}
}
