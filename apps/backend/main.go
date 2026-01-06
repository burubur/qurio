package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"qurio/apps/backend/internal/app"
	"qurio/apps/backend/internal/config"
	"qurio/apps/backend/internal/logger"

	"github.com/nsqio/go-nsq"
)

func main() {
	// Initialize structured logger
	logger := slog.New(logger.NewContextHandler(slog.NewJSONHandler(os.Stdout, nil)))
	slog.SetDefault(logger)

	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// 2. Bootstrap Infrastructure (DB, Weaviate, NSQ Producer, Migrations)
	deps, err := app.Bootstrap(context.Background(), cfg)
	if err != nil {
		slog.Error("bootstrap failed", "error", err)
		os.Exit(1)
	}
	defer deps.DB.Close()

	// 3. Initialize App
	application, err := app.New(cfg, deps.DB, deps.VectorStore, deps.NSQProducer, logger)
	if err != nil {
		slog.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}

	// 4. Worker (Result Consumer) Setup
	nsqCfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("ingest.result", "backend", nsqCfg)
	if err != nil {
		slog.Error("failed to create NSQ consumer for results", "error", err)
	} else {
		// Use AddConcurrentHandlers
		consumer.AddConcurrentHandlers(nsq.HandlerFunc(func(m *nsq.Message) error {
			return application.ResultConsumer.HandleMessage(m)
		}), cfg.IngestionConcurrency)
		
		// Connect to Lookupd
		if err := consumer.ConnectToNSQLookupd(cfg.NSQLookupd); err != nil {
			slog.Error("failed to connect to NSQLookupd", "error", err)
		} else {
			slog.Info("NSQ Result Consumer connected", "concurrency", cfg.IngestionConcurrency)
		}
	}

	// Background Janitor
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := application.SourceService.ResetStuckPages(context.Background()); err != nil {
					slog.Error("failed to reset stuck pages", "error", err)
				}
			}
		}
	}()

	// 5. Start Server
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := application.Run(ctx); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}