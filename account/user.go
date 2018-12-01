package account

import (
	"context"
	"fmt"
	"log"

	fxBroker "github.com/fox-one/broker"
	config "github.com/fox-one/f1db/config"
	fxWallet "github.com/fox-one/foxgo/wallet"
	uuid "github.com/satori/go.uuid"
)

/*
User is used in user system.
*/
type User struct {
	ID      string `yaml:"id" json:"id"`
	Key     string `yaml:"key" json:"key"`
	MixinID string
}

var authUser *User
var userMap map[string]User
var userTokenMap map[string]string

func InitUserDb(ctx context.Context) error {
	var err error
	var stat *fxWallet.Status
	userMap = make(map[string]User)
	userTokenMap = make(map[string]string)
	cfg := config.GetConfig()
	for _, user := range cfg.Users {
		userMap[user.Key] = User{
			ID:      user.ID,
			Key:     user.Key,
			MixinID: "",
		}
		_, err = Login(ctx, user.ID, user.Key)
		if err != nil {
			log.Printf("Login user %s Failed\n", user.ID)
		} else {
			ses := GetSession(user.ID)
			if stat, err = fxWallet.GetWalletStatus(ctx, ses.Token); err == nil {
				fmt.Printf("%v\n", stat)
				// set mixin user id
				userMap[user.Key] = User{
					ID:      user.ID,
					Key:     user.Key,
					MixinID: stat.MixinId,
				}
			} else {
				log.Printf("get wallet stat for user %d failed: %s\n", user.ID, err.Error())
			}
			log.Printf("Login user %s:\n- Key: %s\n- Token %s\n- Mixin ID: %s\n", user.ID, user.Key, ses.Token, stat.MixinId)
		}
	}
	return nil
}

func getUser(userId string, key string) *User {
	if val, ok := userMap[key]; ok {
		if val.ID == userId {
			return &val
		}
	}
	return nil
}

func AuthUser(userId string, key string) *User {
	if user := getUser(userId, key); user != nil {
		return user
	}
	return nil
}

func Register(ctx context.Context, pk string) (*User, error) {
	var user *User
	var uid uuid.UUID
	var err error
	var userResp *fxBroker.UserResponse
	uid, err = uuid.NewV4()
	if err != nil {
		return nil, err
	}
	uidStr := uid.String()
	user = new(User)
	userResp, err = GetBroker().Register(ctx, uidStr, uidStr, "")
	if err != nil {
		return nil, err
	}

	pin := config.GetConfig().General.Pin
	if err := UpdatePin(ctx, userResp.Token, EmptyPin, NewPin(pin, pk)); err != nil {
		return nil, err
	}
	user.ID = userResp.User.Id
	user.Key = uidStr
	return user, nil
}
