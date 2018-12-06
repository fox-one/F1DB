package controller

import (
	"context"
	"fmt"

	account "github.com/fox-one/f1db/account"
	util "github.com/fox-one/f1db/util"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	ID string `json:"id"`
}

// LoginHandler handle login request
func LoginHandler(ctx context.Context) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		req := LoginRequest{}
		ginCtx.ShouldBindJSON(&req)
		_, err := account.Login(ctx, req.ID)
		if err == nil {
			ginCtx.JSON(200, gin.H{
				"code": "0",
				"data": gin.H{
					"id": req.ID,
				},
			})
		} else {
			util.RespError(ginCtx, 403, 3, "Login Error")
		}
	}
}

// RegisterHandler handle register request
func RegisterHandler(ctx context.Context, pk string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		userID, err := account.Register(ctx, pk)
		if err == nil {
			ginCtx.JSON(200, gin.H{
				"code": "0",
				"data": gin.H{
					"id": userID,
				},
			})
		} else {
			util.RespError(ginCtx, 403, 3, fmt.Sprintf("Register Error: %s", err))
		}
	}
}
