package config

import (
	"github.com/spf13/viper"
	"sync"
)

type cfg struct {
	Hostname       string
	GithubToken    string
	DevMode        bool
	KubeconfigPath string
}

var doOnce sync.Once
var config *cfg

func GetConfig() *cfg {
	doOnce.Do(func() {
		config = &cfg{}
		viper.AutomaticEnv()
		config.Hostname = viper.GetString("HOSTNAME")
		config.GithubToken = viper.GetString("GITHUB_TOKEN")
		config.DevMode = viper.GetBool("DEV_MODE")
		if config.DevMode {
			config.KubeconfigPath = viper.GetString("KUBECONFIG_PATH")
		}
	})
	return config
}
