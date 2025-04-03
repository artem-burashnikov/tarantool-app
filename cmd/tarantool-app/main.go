package main

import (
	"tarantool-app/internal/app"
)

// @title           Tarantool Key-Value API
// @version         1.0
// @description     This API provides CRUD operations for managing key-value pairs in the Tarantool database.

// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT

// @host        localhost:8080
// @BasePath    /

func main() {
	configPath := "app_config.yaml"
	app.Run(configPath)
}
