package main

import (
	"tarantool-app/internal/infrastructure/config"
	"tarantool-app/internal/infrastructure/http"
	"tarantool-app/internal/infrastructure/logging"
	"tarantool-app/internal/repository"
	"tarantool-app/internal/usecase"
)

func main() {
	// Load server and repository configuration.
	// This must succeed.
	cfg := config.MustLoad()

	// Initialize logger.
	log := logging.NewLogger(cfg.AppEnv)
	// Sync (flush buffer) before exiting.
	defer func() {
		if err := log.SugaredLogger.Sync(); err != nil {
			log.Warn("Logger sync error",
				"error", err,
			)
		}
	}()

	// Establish connection to the repository.
	// In this case the actuatl initialization is performed in `configs/tarantool_storage_init.lua`.
	tt, err := repository.NewTarantoolRepository(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize repository. Aborting.")
	}
	defer tt.Close()

	// Inject repository into usecase layer.
	crudUsecase := usecase.NewCRUD(tt, log)

	// Inject usecases into handlers.
	apiHandler := http.NewRequestHandler(crudUsecase, log)

	r := http.NewGinRouter(cfg.AppEnv, log, apiHandler)

	// Run the application.
	if err := r.Run(":" + cfg.HTTPServer.Port); err != nil {
		log.Fatal("Failed to start HTTP server",
			"error", err,
		)
	}
}
