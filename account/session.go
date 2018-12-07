package account

import (
	"context"
	"log"

	fxBroker "github.com/fox-one/broker"
	config "github.com/fox-one/f1db/config"
	util "github.com/fox-one/f1db/util"
	fxWallet "github.com/fox-one/foxgo/wallet"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	UserID  string
	Token   string
	MixinID string
}

type Quota struct {
	QuotaID     string  `json:"quota_id"`
	QuotaAmount string  `json:"quota_amount"`
	PublicKey   string  `json:"public_key"`
	Balance     float64 `json:"balance"`
}

var userSessions map[string]*Session

func InitSession() {
	userSessions = make(map[string]*Session)
}

func Login(ctx context.Context, userID string) (string, error) {
	if userSessions == nil {
		userSessions = make(map[string]*Session)
	}
	if len(userID) != 0 {
		resp, err := GetBroker().Login(ctx, userID)
		if err != nil {
			return "", err
		}
		ses := new(Session)
		ses.UserID = userID
		ses.Token = resp.Token
		if stat, err := fxWallet.GetWalletStatus(ctx, ses.Token); err == nil {
			ses.MixinID = stat.MixinId
		} else {
			log.Printf("get wallet stat for user %d failed: %s\n", userID, err.Error())
		}
		log.Printf("User Login:\n- userID: %s\n- MixinID: %s\n", ses.UserID, ses.MixinID)
		userSessions[userID] = ses
	}
	return userID, nil
}

func Register(ctx context.Context, pk string) (string, error) {
	var uid uuid.UUID
	var err error
	var resp *fxBroker.UserResponse
	uid, err = uuid.NewV4()
	if err != nil {
		return "", err
	}
	uidStr := uid.String()
	resp, err = GetBroker().Register(ctx, uidStr, uidStr, "")
	if err != nil {
		return "", err
	}

	pin := config.GetConfig().General.Pin
	if err := UpdatePin(ctx, resp.Token, EmptyPin, NewPin(pin, pk)); err != nil {
		return "", err
	}
	ses := new(Session)
	ses.UserID = resp.User.Id
	ses.Token = resp.Token
	if stat, err := fxWallet.GetWalletStatus(ctx, ses.Token); err == nil {
		ses.MixinID = stat.MixinId
	} else {
		log.Printf("get wallet stat for user %d failed: %s\n", ses.UserID, err.Error())
	}
	userSessions[ses.UserID] = ses
	return ses.UserID, nil
}

func GetQuota(ctx context.Context, token string) (*Quota, error) {
	var assets fxWallet.Assets
	var detail *fxWallet.Asset
	var err error
	assets, err = fxWallet.GetAssets(ctx, token, true)
	if err != nil {
		return nil, err
	}
	ret := &Quota{
		QuotaID:     config.GetConfig().General.QuotaID,
		Balance:     0,
		PublicKey:   "",
		QuotaAmount: config.GetConfig().General.QuotaAmount,
	}
	for _, asset := range assets {
		if asset.AssetId == config.GetConfig().General.QuotaID {
			ret.Balance = asset.Balance
			break
		}
	}
	detail, err = fxWallet.GetAssetDetail(ctx, token, config.GetConfig().General.QuotaID)
	if err != nil {
		return nil, err
	}
	ret.PublicKey = detail.PublicKey
	return ret, nil
}

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Request.Header.Get("-x-user-id")
		if len(userID) == 0 {
			util.RespError(ctx, 403, 2, "Invalid auth info")
			ctx.Abort()
			return
		}
		_, err := Login(ctx, userID)
		if err != nil {
			util.RespError(ctx, 403, 2, "Failed to login the DAG Network gateway")
			ctx.Abort()
		} else {
			ctx.Next()
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
