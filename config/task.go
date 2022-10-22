package config

type taskConfig struct {
	Home string   `yaml:"home"`
	Urls []string `yaml:"urls"`
}

var Task *taskConfig
