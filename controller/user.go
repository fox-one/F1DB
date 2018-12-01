package controller

import (
	account "github.com/fox-one/f1db/account"
	util "github.com/fox-one/f1db/util"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// LoginHandler: handle login request
func LoginHandler(ctx *gin.Context) {
	req := LoginRequest{}
	ctx.ShouldBindJSON(&req)
	user := account.AuthUser(req.ID, req.Key)
	if user != nil {
		ctx.JSON(200, gin.H{
			"code": "0",
			"data": gin.H{
				"id":  user.ID,
				"key": user.Key,
			},
		})
	} else {
		util.RespError(ctx, 403, 3, "Invalid Key")
	}
}
