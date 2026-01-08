package config_test

import (
	"os"
	"testing"
	"qurio/apps/backend/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set env var directly to test envconfig logic
	os.Setenv("DB_HOST", "test-host")
	defer os.Unsetenv("DB_HOST")

	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.Equal(t, "test-host", cfg.DBHost)
}

func TestLoadConfig_FromEnvFile(t *testing.T) {
	// Create a temp .env file
	content := []byte("DB_HOST=loaded-from-file")
	err := os.WriteFile(".env", content, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".env")

	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.Equal(t, "loaded-from-file", cfg.DBHost)
}

func TestLoadConfig_RerankAPIKey(t *testing.T) {
	os.Setenv("RERANK_API_KEY", "test-key")
	defer os.Unsetenv("RERANK_API_KEY")

	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.Equal(t, "test-key", cfg.RerankAPIKey)
}