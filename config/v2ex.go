package config

type V2exConfig struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

var V2ex *V2exConfig
