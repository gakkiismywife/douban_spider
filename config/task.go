package config

type taskConfig struct {
	Home     string   `yaml:"home"`
	Urls     []string `yaml:"urls"`
	Interval int      `yaml:"interval"`
	Seconds  int      `yaml:"seconds"`
    FilterWords []string `yaml:"filterWords"`
}

var Task *taskConfig
