package config

type proxyConfig struct {
	Id       string `yaml:"id"`
	Secret   string `yaml:"secret"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Tunnel   string `yaml:"tunnel"`
}

var Proxy *proxyConfig
