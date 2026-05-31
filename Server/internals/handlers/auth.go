package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ValidateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("Auth handler")
		ctx.Next()
	}
}
