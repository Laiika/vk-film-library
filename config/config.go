package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"time"
	"vk-film-library/pkg/logger"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
	Postgres   `yaml:"postgres"`
	JWT        `yaml:"jwt"`
}

type HTTPServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port" default:"8080"`
}

type Postgres struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port" default:"5432"`
	Database string `yaml:"dbname"`
}

type JWT struct {
	SignKey  string        `yaml:"sign_key"`
	TokenTTL time.Duration `yaml:"token_ttl"`
}

var instance *Config
var once sync.Once

func GetConfig(logger *logger.Logger) *Config {
	once.Do(func() {
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("./", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
