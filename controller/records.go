package controller

import (
	"context"
	"fmt"

	"github.com/fox-one/f1db/account"
	"github.com/fox-one/f1db/config"
	storage "github.com/fox-one/f1db/storage"
	util "github.com/fox-one/f1db/util"
	fxAccount "github.com/fox-one/foxgo/account"
	"github.com/gin-gonic/gin"
)

/*
NewRecordRequest is used in handler create new recoard
*/
type NewRecordRequest struct {
	ItemType string `json:"type"`
	Brief    string `json:"brief"`
	Content  string `json:"content"`
}

/*
KeepRecordRequest is used to keep an exsited record
*/
type KeepRecordRequest struct {
	Quota string `json:"quota"`
}

// GetRecordHandler: handle create content request
func NewRecordHandler(ctx *gin.Context) {
	req := NewRecordRequest{}
	ctx.ShouldBindJSON(&req)
	ses := account.CurrentSession(ctx)
	userID := ses.UserID
	if req.Brief == "" {
		if len(req.Content) > 32 {
			req.Brief = req.Content[:32]
		} else {
			req.Brief = req.Content
		}
	}
	item, err := storage.WriteItem(userID, req.ItemType, req.Brief, req.Content)
	if err != nil {
		util.RespError(ctx, 500, 1, fmt.Sprintf("Failed to write item. %s", err))
		return
	}
	ctx.JSON(200, gin.H{
		"code": "0",
		"data": item.Response(),
	})
}

// GetRecordHandler: handle get content request
func GetRecordHandler(ctx *gin.Context) {
	hash := ctx.Param("hash")
	item, err := storage.ReadRecord(hash)
	if err != nil {
		util.RespError(ctx, 500, 1, fmt.Sprintf("Failed to read item. %s", err))
		return
	}
	ctx.JSON(200, gin.H{
		"code": "0",
		"data": item.Response(),
	})
}

// KeepRecordHandler handle transfer and keep the specified file to storage
func KeepRecordHandler(pk string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var item *storage.Item
		hash := ctx.Param("hash")
		req := KeepRecordRequest{}
		quota := config.GetConfig().General.DefaultQuota
		ctx.ShouldBindJSON(&req)

		if hash == "" {
			util.RespError(ctx, 400, 1, fmt.Sprintf("Invalid arguments."))
			return
		}
		if req.Quota != "" {
			quota = req.Quota
		}
		item, err = storage.ReadRecord(hash)
		if err != nil {
			util.RespError(ctx, 500, 1, fmt.Sprintf("Failed to read record. %s", err))
			return
		}
		// prepare to transfer
		ses := account.CurrentSession(ctx)
		if ses == nil {
			util.RespError(ctx, 405, 1, fmt.Sprintf("Need to auth. %s", err))
			return
		}
		pin := fxAccount.NewPin(config.GetConfig().General.Pin, pk)
		item, err = storage.KeepItem(ctx, *item, ses.Token, pin, quota)
		if err != nil {
			util.RespError(ctx, 500, 10, fmt.Sprintf("Failed to generate transication: %v", err))
			return
		}
		ctx.JSON(200, gin.H{
			"code": "0",
			"data": item.Response(),
		})
	}
}

func GetSnapshotHandler(ctx context.Context) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		sid := ginCtx.Param("snapshot_id")
		item, err := storage.ReadItem(ctx, sid)
		if err != nil {
			util.RespError(ginCtx, 500, 1, fmt.Sprintf("Failed to read item. %s", err))
			return
		}
		ginCtx.JSON(200, gin.H{
			"code": "0",
			"data": item.Response(),
		})
	}
}
