package handlers

import (
	"fmt"
	"net/http"
	"time"

	storage "github.com/Niranjan0524/taskforge/server/internals/Storage"
	"github.com/Niranjan0524/taskforge/server/internals/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTask(store storage.Storage) gin.HandlerFunc {

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

func GetTask(store storage.Storage) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId := ctx.Param("id")
		fmt.Println("userId", userId)

		taskErr, task := store.GetTask(ctx.Request.Context(), userId)

		if taskErr != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
			return
		}
		fmt.Println("Task", task)

		ctx.JSON(http.StatusOK, task)
	}
}

func GetAllTasks(store storage.Storage) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		allTasks, err := store.GetAllTasks(ctx.Request.Context())

		if err != nil {
			fmt.Println("Error in fetching all tasks", err)
			ctx.JSON(http.StatusInternalServerError, err)
		}

		ctx.JSON(http.StatusOK, allTasks)
	}
}

func DeleteTask(store storage.Storage) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		taskId := ctx.Param("id")
		err := store.DeleteTask(ctx.Request.Context(), taskId)

		if err != nil {

			fmt.Println("error deleting task", err)
			ctx.JSON(http.StatusInternalServerError, err)
		}

		ctx.JSON(http.StatusOK, "Task Deleted Successfully")
	}
}
