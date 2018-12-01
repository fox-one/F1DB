package util

import (
	"github.com/gin-gonic/gin"
)

// RespError is a method to respond error info to client
func RespError(ctx *gin.Context, httpCode int, code int, message string) {
	ctx.JSON(httpCode, gin.H{
		"code":    code,
		"message": message,
		"data":    gin.H{},
	})
}
