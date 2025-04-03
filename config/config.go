package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Note that `env` and `env-default` are config `cleanenv` package specific tags.

type Config struct {
	App        AppConfig        `yaml:"app"`
	HTTPServer HTTPServerConfig `yaml:"http_server"`
	Storage    Storage
}

type AppConfig struct {
	Environment string `yaml:"environment" env:"APP_ENV" env-required:"true"`
	Name        string `yaml:"name" env:"APP_NAME" env-required:"true"`
	Version     string `yaml:"version" env:"APP_VERSION" env-required:"true"`
}

type HTTPServerConfig struct {
	Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080" env-required:"true"`
}

type Storage struct {
	Address     string `env:"TT_URI" env-required:"true"`
	Credentials StorageCredentials
}

type StorageCredentials struct {
	Username string `env:"TT_USER" env-required:"true"`
	Password string `env:"TT_PASSWORD" env-required:"true"`
}

func Load(configPath string) (Config, error) {
	if configPath == "" {
		return Config{}, fmt.Errorf("CONFIG_PATH environment variable must be set")
	}

	fileInfo, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("Config file does not exist: %q", configPath)
		}
		return Config{}, err
	}
	if fileInfo.IsDir() {
		return Config{}, fmt.Errorf("config path %q is a directory, not a file", configPath)
	}

	var cfg Config
	err = cleanenv.ReadConfig(configPath, &cfg)

	return cfg, err
}
