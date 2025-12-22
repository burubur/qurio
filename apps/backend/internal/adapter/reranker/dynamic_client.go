package reranker

import (
	"context"
	"fmt"
	"sync"

	"qurio/apps/backend/internal/settings"
)

type DynamicClient struct {
	settingsSvc *settings.Service
	client      *Client
	currentKey  string
	currentProv string
	mu          sync.RWMutex
}

func NewDynamicClient(svc *settings.Service) *DynamicClient {
	return &DynamicClient{settingsSvc: svc}
}

func (c *DynamicClient) Rerank(ctx context.Context, query string, docs []string) ([]int, error) {
	s, err := c.settingsSvc.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	if s.RerankProvider == "none" || s.RerankProvider == "" {
		// Return original order
		indices := make([]int, len(docs))
		for i := range indices {
			indices[i] = i
		}
		return indices, nil
	}

	client := c.getClient(s.RerankProvider, s.RerankAPIKey)
	return client.Rerank(ctx, query, docs)
}

func (c *DynamicClient) getClient(provider, key string) *Client {
	c.mu.RLock()
	if c.client != nil && c.currentKey == key && c.currentProv == provider {
		defer c.mu.RUnlock()
		return c.client
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double check
	if c.client != nil && c.currentKey == key && c.currentProv == provider {
		return c.client
	}

	client := NewClient(provider, key)
	c.client = client
	c.currentKey = key
	c.currentProv = provider
	return client
}
