package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ProxyConfig struct {
	Id       string `yaml:"id"`
	Secret   string `yaml:"secret"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var Proxy *ProxyConfig

func init() {
	viper.AddConfigPath("config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("[config]viper.ReadInConfig err", err)
		return
	}
	Proxy = new(ProxyConfig)
	err = viper.UnmarshalKey("proxy", Proxy)
	if err != nil {
		fmt.Println("[config]viper.Unmarshal err", err)
		return
	}
}
