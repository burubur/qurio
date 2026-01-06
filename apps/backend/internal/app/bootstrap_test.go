package app_test

import (
	"context"
	"testing"
	"qurio/apps/backend/internal/app"
	"qurio/apps/backend/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestBootstrap_ConfigurationError(t *testing.T) {
	cfg := &config.Config{
		DBHost: "invalid-host",
	}
	deps, err := app.Bootstrap(context.Background(), cfg)
	assert.Error(t, err)
	assert.Nil(t, deps)
}
