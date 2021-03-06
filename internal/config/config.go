package config

import (
    "github.com/spf13/viper"
    "sync"
    "time"
)

type cfg struct {
    Hostname        string
    GithubToken     string
    DevMode         bool
    KubeconfigPath  string
    Configmap       string
    AdminToken      string
    CheckInterval   time.Duration
    RedisAddress    string
    RedisStudKey    string
    KubeApiServer   string
    KubeApiServerCA string
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
        config.Configmap = viper.GetString("CONFIGMAP_NAME")
        config.AdminToken = viper.GetString("ADMIN_TOKEN")
        config.CheckInterval = viper.GetDuration("CHECK_INTERVAL")
        config.RedisAddress = viper.GetString("REDIS_ADDRESS")
        config.RedisStudKey = viper.GetString("REDIS_STUDENTS_KEY")
        config.KubeApiServer = viper.GetString("KUBE_API_SERVER")
        config.KubeApiServerCA = viper.GetString("KUBE_API_SERVER_CA")
        if config.DevMode {
            config.KubeconfigPath = viper.GetString("KUBECONFIG_PATH")
        }
    })
    return config
}
