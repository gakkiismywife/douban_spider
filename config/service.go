package config

type ServiceConfig struct {
	MaxRetryTimes int `yaml:"maxRetryTimes"`
}

var Service *ServiceConfig 