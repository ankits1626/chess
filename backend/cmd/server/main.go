package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "chess-coach/backend/internal/config"
    "chess-coach/backend/internal/db"
    "chess-coach/backend/internal/server"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        slog.Error("failed to load config", "error", err)
        os.Exit(1)
    }

    // Initialize logger
    logLevel := slog.LevelInfo
    if cfg.IsDevelopment() {
        logLevel = slog.LevelDebug
    }
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: logLevel,
    }))
    slog.SetDefault(logger)

    logger.Info("loaded configuration",
        "environment", cfg.Environment,
        "port", cfg.Port,
    )

    // Connect to database
    ctx := context.Background()
    pool, err := db.Connect(ctx, cfg.DatabaseURL)
    if err != nil {
        logger.Error("database connection failed", "error", err)
        os.Exit(1)
    }
    defer pool.Close()

    logger.Info("connected to database")

    // Create queries instance
    queries := db.New(pool)

    // Create server
    srv := server.New(logger, queries)
    srv.SetupRoutes()

    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
        <-sigChan

        logger.Info("received shutdown signal")
        if err := srv.Shutdown(); err != nil {
            logger.Error("shutdown error", "error", err)
        }
    }()

    // Start server
    if err := srv.Start(cfg.Port); err != nil {
        logger.Error("server error", "error", err)
        os.Exit(1)
    }
}