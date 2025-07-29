package main

import (
	"context"
	configs "demo/internal/config"
	"demo/internal/http/handlers/students"
	"demo/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Load configuration
	slog.Info("Loading configuration")
	cfg := configs.MustLoad()

	// Initialize storage
	slog.Info("Initializing storage")
	storage, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("Failed to initialize storage", slog.String("env", cfg.Env), slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize HTTP server
	slog.Info("Setting up HTTP server")
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", students.New(storage))
	router.HandleFunc("GET /api/students/{id}", students.GetById(storage))
	router.HandleFunc("PUT /api/students/{id}", students.Update(storage))
	router.HandleFunc("DELETE /api/students/{id}", students.Delete(storage))
	router.HandleFunc("GET /api/students", students.GetAll(storage))

	// Start HTTP server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		slog.Info("Starting server", slog.String("address", cfg.Addr))
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("Failed to start server", slog.String("error", err.Error()))
		}
	}()

	<-done

	slog.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server gracefully", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}
