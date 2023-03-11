package config

type ServiceConfig struct {
	MaxRetryTimes int8 `yaml:"maxRetryTimes"`
}

var Service *ServiceConfig
