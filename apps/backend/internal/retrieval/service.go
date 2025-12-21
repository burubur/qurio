package retrieval

import (
	"context"
)

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

type VectorStore interface {
	Search(ctx context.Context, vector []float32) ([]string, error)
}

type Service struct {
	embedder Embedder
	store    VectorStore
}

func NewService(e Embedder, s VectorStore) *Service {
	return &Service{embedder: e, store: s}
}

func (s *Service) Search(ctx context.Context, query string) ([]string, error) {
	vec, err := s.embedder.Embed(ctx, query)
	if err != nil {
		return nil, err
	}
	return s.store.Search(ctx, vec)
}
