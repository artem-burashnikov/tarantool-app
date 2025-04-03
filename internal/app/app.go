package app

import (
	"tarantool-app/config"
	"tarantool-app/internal/infrastructure/http"
	"tarantool-app/internal/repository"
	"tarantool-app/internal/usecases"
	"tarantool-app/internal/utils"
)

func Run(configPath string) {
	cfg := utils.Must(config.Load(configPath))

	log := NewLogger(cfg.App.Environment)
	defer log.Sync()

	tt := utils.Must(repository.NewTarantoolRepository(&cfg, &log))
	defer tt.Close()

	usecase := usecases.NewUserUseCase(tt, log)

	apiHandler := http.NewRequestHandler(usecase, &log)

	r := http.NewGinRouter(cfg.App.Environment, log, &apiHandler)

	if err := r.Run(":" + cfg.HTTPServer.Port); err != nil {
		log.Fatal("Failed to start HTTP server",
			"error", err,
		)
	}
}
