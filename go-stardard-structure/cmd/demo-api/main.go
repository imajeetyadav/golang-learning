package main

import (
	"context"
	configs "demo/internal/config"
	"demo/internal/http/handlers/students"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := configs.MustLoad()

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", students.New())

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
