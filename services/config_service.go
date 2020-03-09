package services

import (
	"github.com/BurntSushi/toml"
)

type ConfigService struct {
	AppName string `toml:"app_name"`
	Version string
	PG      pgConfig
	HTTP    httpEndpointConfig `toml:"http"`
	Auth    authConfig
}

type authConfig struct {
	TokenTTL      int64 `toml:"token_ttl"`
	TokenTTLRenew int64 `toml:"token_renew_ttl"`
}

type pgConfig struct {
	Address  string
	Database string
	User     string
	Password string
	PoolSize int `toml:"pool_size"`
}

type httpEndpointConfig struct {
	Address   string
	RateLimit float64 `toml:"rate_limit"`
}

func (c *ConfigService) Load(file string) error {
	if _, err := toml.DecodeFile(file, &c); err != nil {
		return err
	}

	return nil
}
