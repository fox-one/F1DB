package account

import (
	"context"

	util "github.com/fox-one/f1db/util"
	"github.com/gin-gonic/gin"
)

type Session struct {
	UserID string
	Key    string
	Token  string
}

var userSessions map[string]*Session

func Login(ctx context.Context, userID string, userKey string) (*User, error) {
	user := AuthUser(userID, userKey)
	if userSessions == nil {
		userSessions = make(map[string]*Session)
	}
	if user != nil {
		resp, err := GetBroker().Login(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		ses := new(Session)
		ses.UserID = userID
		ses.Token = resp.Token
		ses.Key = userKey
		userSessions[userID] = ses
	}
	return user, nil
}

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Request.Header.Get("-x-user-id")) == 0 || len(ctx.Request.Header.Get("-x-user-key")) == 0 {
			util.RespError(ctx, 403, 2, "Invalid auth info")
			ctx.Abort()
			return
		}
		userID := ctx.Request.Header.Get("-x-user-id")
		userKey := ctx.Request.Header.Get("-x-user-key")
		if user := AuthUser(userID, userKey); user != nil {
			_, err := Login(ctx, userID, userKey)
			if err != nil {
				util.RespError(ctx, 403, 2, "Failed to login the DAG Network gateway")
				ctx.Abort()
			} else {
				ctx.Next()
			}
		} else {
			util.RespError(ctx, 403, 2, "Incorrect user ID or key")
			ctx.Abort()
		}
	}
}

func GetSession(userID string) *Session {
	if userSessions == nil {
		userSessions = make(map[string]*Session)
	}
	if val, ok := userSessions[userID]; ok {
		return val
	}
	return nil
}

func CurrentSession(c *gin.Context) *Session {
	if userID := c.GetHeader("-x-user-id"); userID != "" {
		return GetSession(userID)
	}
	return nil
}
