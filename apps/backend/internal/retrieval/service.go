package retrieval

import (
	"context"
)

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

type VectorStore interface {
	Search(ctx context.Context, query string, vector []float32, alpha float32) ([]string, error)
}

type Reranker interface {
	Rerank(ctx context.Context, query string, docs []string) ([]string, error)
}

type Service struct {
	embedder Embedder
	store    VectorStore
	reranker Reranker
}

func NewService(e Embedder, s VectorStore, r Reranker) *Service {
	return &Service{embedder: e, store: s, reranker: r}
}

func (s *Service) Search(ctx context.Context, query string) ([]string, error) {
	// 1. Embed Query
	vec, err := s.embedder.Embed(ctx, query)
	if err != nil {
		return nil, err
	}

	// 2. Hybrid Search (BM25 + Vector)
	// Default alpha = 0.5
	docs, err := s.store.Search(ctx, query, vec, 0.5)
	if err != nil {
		return nil, err
	}

	// 3. Rerank (if configured)
	if s.reranker != nil && len(docs) > 0 {
		return s.reranker.Rerank(ctx, query, docs)
	}

	return docs, nil
}