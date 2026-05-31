package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("Creates task")
	}
}

func GetTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("id")
		fmt.Println("userId", userId)
	}
}
