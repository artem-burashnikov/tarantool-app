package main

import (
	"tarantool-app/internal/app"
)

func main() {
	configPath := "app_config.yaml"
	app.Run(configPath)
}
