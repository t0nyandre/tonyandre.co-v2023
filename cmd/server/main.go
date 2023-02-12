package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/t0nyandre/tonyandre.co-v2023/internal/config"
	"github.com/t0nyandre/tonyandre.co-v2023/internal/healthcheck"
	"github.com/t0nyandre/tonyandre.co-v2023/internal/user"
	"github.com/t0nyandre/tonyandre.co-v2023/pkg/database/postgres"
	"github.com/t0nyandre/tonyandre.co-v2023/pkg/logger"
)

var appConfig = flag.String("config", "./config/local.json", "path to config file")

func main() {
	flag.Parse()

	log := logger.NewLogger()

	cfg, err := config.Load(*appConfig, log)
	if err != nil {
		log.Fatalw("Failed to load config", "error", err)
	}

	// Connect to database
	db, err := postgres.NewPostgres(log, cfg)
	if err != nil {
		log.Fatalw("Failed to connect to database", "database", cfg.PostgresDb, "error", err)
	}

	router := chi.NewRouter()
	router.Use(logger.LoggingMiddleware(log))
	router.Mount("/healthcheck", healthcheck.RegisterHandlers(cfg))
	router.Mount("/v1/users", user.RegisterHandlers(user.NewService(user.NewRepository(db, log), log), log))

	log.Infow("Server successfully up and running", "host", cfg.AppHost, "port", cfg.AppPort)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%v", cfg.AppHost, cfg.AppPort), router); err != nil {
		log.Fatalw("Server failed to start", "error", err)
	}
}
