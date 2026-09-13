package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/26081yogesh/students-api-golang/internal/config"
)

func main() {
    // 1. Load config file
    cfg := config.MustLoad()

    // 2. Setup Structured Logger (slog) based on Environment
    var logger *slog.Logger
    if cfg.Env == "dev" {
        logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    } else {
        logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    }
    slog.SetDefault(logger)

    // 3. Setup database (coming soon)

    // 4. Setup router
    router := http.NewServeMux()
    router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Welcome to Students API"))
    })

    // Corrected slog usage: clean message + key-value attribute without trailing spaces
    slog.Info("server started", slog.String("address", cfg.Addr), slog.String("env", cfg.Env))

    // 5. Setup server
    server := http.Server{
        Addr:    cfg.Addr,
        Handler: router,
    }

    // 6. Graceful shutdown setup
    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        slog.Info("starting server...")
        err := server.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
            log.Fatalf("failed to start server: %v", err)
        }
    }()

    // Block until we receive our signal
    <-done
    slog.Info("shutting down the server gracefully...")

    // Give outstanding requests 5 seconds to finish
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := server.Shutdown(ctx)
    if err != nil {
        slog.Error("failed to shutdown the server", slog.String("error", err.Error()))
    }

    slog.Info("server shutdown successfully")
}