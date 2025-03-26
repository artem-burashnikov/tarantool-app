// Configuration loader interface and implementation.
// Interface may be used for testing.

package config

import (
	"log"
	"os"
	"tarantool-app/internal/config"

	"github.com/ilyakaznacheev/cleanenv"
)

type Loader interface {
	Load() (*config.Config, error)
}

type YAMLconfigLoader struct {
	configPath string
}

// Configuration loader constructor.
func NewYAMLconfigLoader(configPath string) (*YAMLconfigLoader, error) {
	return &YAMLconfigLoader{configPath: configPath}, nil
}

// Satisfies Loader interface.
// Checks if file exists and parses it using `cleanenv` package.
func (l *YAMLconfigLoader) Load() (*config.Config, error) {
	if _, err := os.Stat(l.configPath); os.IsExist(err) {
		return nil, err
	}

	var cfg config.Config
	if err := cleanenv.ReadConfig(l.configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Configuration loader wrapper.
// It is the very first function to be called at the start of an application and must succeed.
// The logger is not ready yet, using plaintext stdout.
func MustLoad() *config.Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	loader, err := NewYAMLconfigLoader(configPath)
	if err != nil {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	cfg, err := loader.Load()
	if err != nil {
		log.Fatalf("Could not read config: %s", err)
	}
	return cfg
}
