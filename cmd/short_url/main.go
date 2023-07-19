package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"short_url/internal/config"
	"short_url/internal/http-server/handlers/save"
	"short_url/internal/lib/logger/sl"
	"short_url/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.LoadConfig()
	fmt.Println(cfg)
	// TODO: init logger: slog
	log := setUpLogger(cfg.Env)
//	log = log.With(slog.String("storage_path", cfg.StoragePath))
	
	log.Debug("Debug message", slog.String("env", cfg.Env))
	log.Info("Info message")
	log.Warn("Warning message")
	log.Error("Error message")
	
	// TODO: init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
//		log.Error("Failed to init storage: ", err)
		log.Error("Failed to init storage: ", sl.Err(err))
		os.Exit(1)
	}
	
	// TODO: router: chi, chi/render
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Post("/url", save.New(log, storage))
	
//	http.ListenAndServe(":3000", r)
	// TODO: run server
	log.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr: cfg.Address,
		Handler: r,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Error("server was killed")
	}
}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
} 