package config

type wechatConfig struct {
	Key     string `yaml:"key"`
	Secret  string `yaml:"secret"`
	AgentId string `yaml:"agentId"`
}

var Wechat *wechatConfig
