package config

import (
	"billing/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Listen ListenConn `yaml:"listen"`
	Db     DbConn     `yaml:"db"`
}

type ListenConn struct {
	Host string `yaml:"host" env-default:"127.0.0.1"`
	Port string `yaml:"port" env-default:"8080"`
}

type DbConn struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host" env-default:"120.0.0.1"`
	Port     string `yaml:"port" env-default:"5432"`
	DbName   string `yaml:"db_name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Sslmode  string `yaml:"sslmode"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("./config/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
