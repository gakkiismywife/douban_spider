package config

type wechatConfig struct {
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
}

var Wechat *wechatConfig
