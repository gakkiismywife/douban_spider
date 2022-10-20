package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var vp *viper.Viper

func init() {
	vp = viper.New()

	vp.AddConfigPath("config/")
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")

	err := vp.ReadInConfig()
	if err != nil {
		fmt.Println("[config]viper.ReadInConfig err", err)
		return
	}

	Wechat = new(wechatConfig)
	_ = vp.UnmarshalKey("wechat", Wechat)
	Proxy = new(proxyConfig)
	_ = vp.UnmarshalKey("proxy", Proxy)
}
