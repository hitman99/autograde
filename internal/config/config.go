package config

import (
	"github.com/spf13/viper"
	"sync"
)

type cfg struct {
	Hostname    string
	GithubToken string
}

var doOnce sync.Once

func GetConfig() *cfg {
	cfg := &cfg{}
	doOnce.Do(func() {
		viper.AutomaticEnv()
		cfg.Hostname = viper.GetString("HOSTNAME")
		cfg.GithubToken = viper.GetString("GITHUB_TOKEN")
	})
	return cfg
}
