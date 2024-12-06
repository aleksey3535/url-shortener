package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handler"
	"url-shortener/internal/http-server/middleware"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"	
)

func main() {
	// init config  cleanenv
	cfg := config.MustLoad()

	// init logger	slog
	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	// init storage	postgres
	db := postgres.MustLoad(cfg)
	postgres.MustRunMigrates(log, db)
	storage := storage.New(db)
	// init handler gorilla/mux
	middleware := middleware.New(log)
	handler := handler.New(log, middleware, storage)
	log.Info("URL-SHORTENER started", slog.String("url:", cfg.Address))
	if err := http.ListenAndServe(cfg.HttpServer.Address, handler.InitRoutes()); err != nil {
		panic(err)
	}
	// init router	chi, render

	// run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level:  slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log 

}