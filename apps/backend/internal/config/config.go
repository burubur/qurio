package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBHost string `envconfig:"DB_HOST" default:"postgres"`
	DBPort int    `envconfig:"DB_PORT" default:"5432"`
	DBUser string `envconfig:"DB_USER" default:"qurio"`
	DBPass string `envconfig:"DB_PASS" default:"password"`
	DBName string `envconfig:"DB_NAME" default:"qurio"`

	WeaviateHost   string `envconfig:"WEAVIATE_HOST" default:"localhost:8080"`
	WeaviateScheme string `envconfig:"WEAVIATE_SCHEME" default:"http"`

	GeminiKey  string `envconfig:"GEMINI_KEY"`
	DoclingURL string `envconfig:"DOCLING_URL" default:"http://docling:8000"`
	NSQLookupd string `envconfig:"NSQ_LOOKUPD" default:"nsqlookupd:4161"`
	NSQDHost   string `envconfig:"NSQD_HOST" default:"nsqd:4150"`
	
	RerankProvider string `envconfig:"RERANK_PROVIDER" default:"jina"` // jina, cohere
	RerankAPIKey   string `envconfig:"RERANK_API_KEY"`
}

func Load() (*Config, error) {
	// Try loading .env from current dir and repo root
	// Ignore errors, as env vars might be set in the shell
	_ = godotenv.Load(".env")
	
	// Try finding root .env (assuming 2 levels up if in apps/backend)
	cwd, _ := os.Getwd()
	rootEnv := filepath.Join(cwd, "../../.env")
	_ = godotenv.Load(rootEnv)
	
	// Also try just "../.env" or simple heuristic if needed, but the above covers common cases

	var cfg Config
	err := envconfig.Process("", &cfg)
	
	// Log warning if critical keys are missing, to help debug
	if cfg.GeminiKey == "" {
		log.Println("WARNING: GEMINI_KEY is missing. Ingestion will be disabled.")
	}
	
	return &cfg, err
}
