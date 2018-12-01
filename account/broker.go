package account

import (
	"github.com/fox-one/broker"
	config "github.com/fox-one/f1db/config"
	"github.com/fox-one/foxgo/request"
)

var globalBroker *broker.Broker

func SetupBroker(c config.BrokerCfg) {
	b, err := broker.New(broker.Config{
		AppId:     c.AppID,
		AppSecret: c.AppSecret,
		Develop:   c.DevMode,
	})
	if err != nil {
		panic(err)
	}

	globalBroker = &b

	if c.DevMode {
		request.ActiveDevEnvironment()
	}
}

func GetBroker() *broker.Broker {
	return globalBroker
}
