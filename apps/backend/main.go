package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
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
	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Database Connection
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open db connection: %v", err)
	}
	defer db.Close()

	// Retry connection
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		log.Printf("Failed to ping db, retrying in 2s... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping db after retries: %v", err)
	}

	// 3. Run Migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	// 4. Weaviate Connection & Schema
	wCfg := weaviate.Config{
		Host:   cfg.WeaviateHost,
		Scheme: cfg.WeaviateScheme,
	}
	wClient, err := weaviate.NewClient(wCfg)
	if err != nil {
		log.Fatalf("Failed to create weaviate client: %v", err)
	}

	wAdapter := vector.NewWeaviateClientAdapter(wClient)
	
	// Retry Weaviate Schema Ensure
	for i := 0; i < 10; i++ {
		if err := vector.EnsureSchema(context.Background(), wAdapter); err == nil {
			log.Println("Weaviate schema ensured")
			break
		}
		log.Printf("Failed to ensure weaviate schema, retrying in 2s... (%d/10) Error: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err := vector.EnsureSchema(context.Background(), wAdapter); err != nil {
		log.Fatalf("Failed to ensure weaviate schema after retries: %v", err)
	}

	// 5. Initialize Adapters & Services
	doclingClient := docling.NewClient(cfg.DoclingURL)
	vecStore := wstore.NewStore(wClient)

	// NSQ Producer
	nsqCfg := nsq.NewConfig()
	nsqProducer, err := nsq.NewProducer(cfg.NSQDHost, nsqCfg)
	if err != nil {
		log.Fatalf("Failed to create NSQ producer: %v", err)
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
		log.Printf("Failed to create NSQ consumer: %v", err)
	} else {
		consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
			return ingestHandler.HandleMessage(m)
		}))
		// Connect to Lookupd
		if err := consumer.ConnectToNSQLookupd(cfg.NSQLookupd); err != nil {
			log.Printf("Failed to connect to NSQLookupd: %v", err)
		} else {
			log.Println("NSQ Consumer connected")
		}
	}

	// 6. Start Server
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server starting on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}