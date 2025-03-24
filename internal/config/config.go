// Structures below mirror a config.yaml defined in tarantool-app/config.

package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-default:"prod"`
	Storage    Storage    `yaml:"storage"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

type Storage struct {
	Host        string `yaml:"storage_host" env:"STORAGE_HOST" env-required:"true"`
	Port        string `yaml:"storage_port" env:"STORAGE_PORT" env-required:"true"`
	Credentials struct {
		User     string `yaml:"user" env:"STORAGE_USER" env-required:"true"`
		Password string `yaml:"password" env:"STORAGE_PASSWORD" env-required:"true"`
	}
}

type HTTPServer struct {
	Host       string        `yaml:"app_host" env:"HOST" env-default:"localhost"`
	Port       string        `yaml:"app_port" env:"PORT" env-default:"8080"`
	Timeout    time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Could not read config: %s", err)
	}

	return &cfg
}
