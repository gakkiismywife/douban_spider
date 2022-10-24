package config

type taskConfig struct {
	Home     string   `yaml:"home"`
	Urls     []string `yaml:"urls"`
	Interval int      `yaml:"interval"`
	Seconds  int      `yaml:"seconds"`
}

var Task *taskConfig
