package config

import (
	"log"
	"os"
	"tarantool-app/internal/config"

	"github.com/ilyakaznacheev/cleanenv"
)

// Configuration loader interface.
type Loader interface {
	Load() (*config.Config, error)
}

// Configuration loader model.
type ConfigLoader struct {
	configPath string
}

// Configuration loader constructor.
// Checks if file at `configPath` exists and returns `ConfigLoader` model on success.
func NewConfigLoader(configPath string) (*ConfigLoader, error) {
	if _, err := os.Stat(configPath); os.IsExist(err) {
		return nil, err
	}
	return &ConfigLoader{configPath: configPath}, nil
}

// Configuration loader interface implementation.
// Parses configuration file using `cleanenv` package.
func (l *ConfigLoader) Load() (*config.Config, error) {
	var cfg config.Config
	if err := cleanenv.ReadConfig(l.configPath, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Configuration loader wrapper.
// This is only called once at the start of the application and must succeed.
func MustLoad() *config.Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	loader, err := NewConfigLoader(configPath)
	if err != nil {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	cfg, err := loader.Load()
	if err != nil {
		log.Fatalf("Could not read config: %s", err)
	}
	return cfg
}
