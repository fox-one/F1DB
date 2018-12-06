package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

/*
Config is used in user system.
*/
type Config struct {
	General GeneralCfg
	Server  ServerCfg
	Ipfs    IpfsCfg
	Broker  BrokerCfg
}

/*
ServerCfg is used in user system.
*/
type ServerCfg struct {
	Host string `yaml:"host" json:"host"`
}

/*
GeneralCfg is used in user system.
*/
type GeneralCfg struct {
	QuotaID         string `yaml:"quota_id"`
	DefaultQuota    string `yaml:"default_quota"`
	CollectorUserID string `yaml:"collector_user_id"`
	Pin             string `yaml:"pin"`
}

/*
BrokerCfg is used in user system.
*/
type BrokerCfg struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
	DevMode   bool   `yaml:"dev_mode"`
}

/*
IpfsCfg is used in user system.
*/
type IpfsCfg struct {
	URL string `yaml:"url" json:"url"`
}

var cfg *Config

func Init() (*Config, error) {
	cfg = new(Config)
	bytes, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(bytes), cfg)
	if err != nil {
		panic(err)
	}
	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}
