package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Niranjan0524/taskforge/server/internals/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	fmt.Println("API server")

	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message":  "Working",
			"ClientIp": ctx.ClientIP(),
		})
	})

	authorized := router.Group("/", handlers.ValidateUser())

	authorized.POST("/api/task", handlers.CreateTask())
	authorized.GET("/api/task/:id", handlers.GetTask())

	if err := router.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
