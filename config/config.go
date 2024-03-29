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
	Task = new(taskConfig)
	_ = vp.UnmarshalKey("task", Task)
	Database = new(DatabaseConfig)
	_ = vp.UnmarshalKey("database", Database)
	Service = new(ServiceConfig)
	_ = vp.UnmarshalKey("service", Service)
	RedisSetting = new(RedisConfig)
	_ = vp.UnmarshalKey("redis", RedisSetting)
	V2ex = new(V2exConfig)
	_ = vp.UnmarshalKey("v2ex", V2ex)
}
