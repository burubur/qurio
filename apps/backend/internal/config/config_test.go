package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"qurio/apps/backend/internal/config"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("WEAVIATE_HOST", "localhost:8080")
	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "localhost:8080", cfg.WeaviateHost)
}
