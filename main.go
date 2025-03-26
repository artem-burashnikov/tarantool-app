package main

import (
	"tarantool-app/internal/infrastructure/config"
	"tarantool-app/internal/infrastructure/http"
	"tarantool-app/internal/infrastructure/logging"
	"tarantool-app/internal/repository"
	"tarantool-app/internal/usecase"
)

func main() {
	cfg := config.MustLoad()

	log := logging.NewLogger(cfg.AppEnv)
	defer log.SugaredLogger.Sync()

	tt := repository.NewTarantoolRepository()

	crudUsecase := usecase.NewCRUD(tt)

	apiHandler := http.NewRequestHandler(crudUsecase, log)

	r := http.NewGinRouter(cfg.AppEnv, log, apiHandler)

	r.Run(":" + cfg.HTTPServer.Port)
}
