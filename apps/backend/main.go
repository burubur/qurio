package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"qurio/apps/backend/features/mcp"
	"qurio/apps/backend/features/source"
	"qurio/apps/backend/internal/adapter/docling"
	"qurio/apps/backend/internal/adapter/gemini"
	"qurio/apps/backend/internal/adapter/reranker"
	wstore "qurio/apps/backend/internal/adapter/weaviate"
	"qurio/apps/backend/internal/config"
	"qurio/apps/backend/internal/retrieval"
	"qurio/apps/backend/internal/vector"
	"qurio/apps/backend/internal/settings"
	"qurio/apps/backend/internal/worker"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nsqio/go-nsq"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// 2. Database Connection
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("failed to open db connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Retry connection
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		slog.Warn("failed to ping db, retrying...", "attempt", i+1, "max_attempts", 10)
		time.Sleep(2 * time.Second)
	}

	if err := db.Ping(); err != nil {
		slog.Error("failed to ping db after retries", "error", err)
		os.Exit(1)
	}

	// 3. Run Migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("failed to create migration driver", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		slog.Error("failed to create migration instance", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations applied successfully")

	// 4. Weaviate Connection & Schema
	wCfg := weaviate.Config{
		Host:   cfg.WeaviateHost,
		Scheme: cfg.WeaviateScheme,
	}
	wClient, err := weaviate.NewClient(wCfg)
	if err != nil {
		slog.Error("failed to create weaviate client", "error", err)
		os.Exit(1)
	}

	wAdapter := vector.NewWeaviateClientAdapter(wClient)
	
	// Retry Weaviate Schema Ensure
	for i := 0; i < 10; i++ {
		if err := vector.EnsureSchema(context.Background(), wAdapter); err == nil {
			slog.Info("weaviate schema ensured")
			break
		}
		slog.Warn("failed to ensure weaviate schema, retrying...", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}

	if err := vector.EnsureSchema(context.Background(), wAdapter); err != nil {
		slog.Error("failed to ensure weaviate schema after retries", "error", err)
		os.Exit(1)
	}

	// 5. Initialize Adapters & Services
	doclingClient := docling.NewClient(cfg.DoclingURL)
	vecStore := wstore.NewStore(wClient)

	// NSQ Producer
	nsqCfg := nsq.NewConfig()
	nsqProducer, err := nsq.NewProducer(cfg.NSQDHost, nsqCfg)
	if err != nil {
		slog.Error("failed to create NSQ producer", "error", err)
		os.Exit(1)
	}

	// Feature: Source
	sourceRepo := source.NewPostgresRepo(db)
	sourceService := source.NewService(sourceRepo, nsqProducer)
	sourceHandler := source.NewHandler(sourceService)

	// Feature: Settings
	settingsRepo := settings.NewPostgresRepo(db)
	settingsService := settings.NewService(settingsRepo)
	settingsHandler := settings.NewHandler(settingsService)

	// Adapters: Dynamic
	geminiEmbedder := gemini.NewDynamicEmbedder(settingsService)
	rerankerClient := reranker.NewDynamicClient(settingsService)

	// Middleware: CORS
	enableCORS := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}

	// Routes
	http.HandleFunc("POST /sources", enableCORS(sourceHandler.Create))
	http.HandleFunc("GET /sources", enableCORS(sourceHandler.List))
	http.HandleFunc("DELETE /sources/{id}", enableCORS(sourceHandler.Delete))
	http.HandleFunc("POST /sources/{id}/resync", enableCORS(sourceHandler.ReSync))

	http.HandleFunc("GET /settings", enableCORS(settingsHandler.GetSettings))
	http.HandleFunc("PUT /settings", enableCORS(settingsHandler.UpdateSettings))

	// Feature: Retrieval & MCP
	retrievalService := retrieval.NewService(geminiEmbedder, vecStore, rerankerClient)
	mcpHandler := mcp.NewHandler(retrievalService)
	http.Handle("/mcp", mcpHandler)

	// Worker (Ingest)
	ingestHandler := worker.NewIngestHandler(doclingClient, geminiEmbedder, vecStore, nsqProducer, sourceRepo)
	
	nsqCfg = nsq.NewConfig()
	consumer, err := nsq.NewConsumer("ingest", "channel", nsqCfg)
	if err != nil {
		slog.Error("failed to create NSQ consumer", "error", err)
	} else {
		consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
			return ingestHandler.HandleMessage(m)
		}))
		// Connect to Lookupd
		if err := consumer.ConnectToNSQLookupd(cfg.NSQLookupd); err != nil {
			slog.Error("failed to connect to NSQLookupd", "error", err)
		} else {
			slog.Info("NSQ Consumer connected")
		}
	}

	// 6. Start Server
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	slog.Info("server starting", "port", 8081)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}