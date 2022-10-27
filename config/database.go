package config

type DatabaseConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	DbName      string `yaml:"dbName"`
	Charset     string `yaml:"charset"`
	ParseTime   string `yaml:"parseTime"`
	MaxOpen     int    `yaml:"maxOpen"`
	MaxIdle     int    `yaml:"maxIdle"`
	MaxIdleTime int    `yaml:"maxIdleTime"`
}

var Database *DatabaseConfig
