package config

import (
	"sync"

	"github.com/IvanKyrylov/cost-management-api/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	Postgres struct {
		Username string `yaml:"username" env-required:"true"`
		Password string `yaml:"password" env-required:"true"`
		Host     string `yaml:"host" env-required:"true"`
		Port     string `yaml:"port" env-required:"true"`
		DBName   string `yaml:"dbname" env-required:"true"`
		SslMode  string `yaml:"sslmode" env-required:"true"`
	} `yaml:"postgres" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logging.CommonLog.Println("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logging.CommonLog.Println(help)
			logging.ErrorLog.Fatal(err)
		}
	})
	return instance
}
