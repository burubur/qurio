package config

import "github.com/kelseyhightower/envconfig"

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
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
