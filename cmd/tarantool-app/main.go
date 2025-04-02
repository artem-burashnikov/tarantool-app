package main

import (
	"os"
	"tarantool-app/internal/app"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	app.Run(configPath)
}
