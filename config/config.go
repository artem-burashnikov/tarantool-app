package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

//Note that `env` and `env-default` are config parser package specific tags.

type (
	Config struct {
		App        `yaml:"app"`
		HTTPServer `yaml:"http"`
		Log        `yaml:"log"`
		Storage
	}

	App struct {
		Name    string `yaml:"name" env:"APP_NAME" env-required:"true"`
		Version string `yaml:"version" env:"APP_VERSION" env-required:"true"`
	}

	HTTPServer struct {
		Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080" env-required:"true"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" env-required:"true"`
	}

	Storage struct {
		URI string `env:"TT_URI" env-required:"true"`
	}
)

func Load(configPath string) (Config, error) {
	if configPath == "" {
		return Config{}, fmt.Errorf("CONFIG_PATH is not set")
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
	// TODO: use some other config parser.
	err = cleanenv.ReadConfig(configPath, &cfg)

	return cfg, err
}
