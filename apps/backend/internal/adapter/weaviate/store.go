package weaviate

import (
	"context"
	"fmt"
	"qurio/apps/backend/internal/worker"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
)

type Store struct {
	client *weaviate.Client
}

func NewStore(client *weaviate.Client) *Store {
	return &Store{client: client}
}

func (s *Store) StoreChunk(ctx context.Context, chunk worker.Chunk) error {
	_, err := s.client.Data().Creator().
		WithClassName("DocumentChunk").
		WithProperties(map[string]interface{}{
			"content":    chunk.Content,
			"url":        chunk.SourceURL,
			"sourceId":   chunk.SourceURL,
			"chunkIndex": chunk.ChunkIndex,
		}).
		WithVector(chunk.Vector).
		Do(ctx)
	return err
}

func (s *Store) Search(ctx context.Context, query string, vector []float32, alpha float32) ([]string, error) {
	hybrid := s.client.GraphQL().HybridArgumentBuilder().
		WithQuery(query).
		WithVector(vector).
		WithAlpha(alpha)

	res, err := s.client.GraphQL().Get().
		WithClassName("DocumentChunk").
		WithHybrid(hybrid).
		WithLimit(5).
		WithFields(graphql.Field{Name: "content"}).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(res.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %v", res.Errors)
	}

	var results []string
	if data, ok := res.Data["Get"].(map[string]interface{}); ok {
		if chunks, ok := data["DocumentChunk"].([]interface{}); ok {
			for _, c := range chunks {
				if props, ok := c.(map[string]interface{}); ok {
					if content, ok := props["content"].(string); ok {
						results = append(results, content)
					}
				}
			}
		}
	}

	return results, nil
}
