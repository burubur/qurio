package retrieval

import (
	"context"
	"time"
)

type SearchResult struct {
	Content  string                 `json:"content"`
	Score    float32                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

type VectorStore interface {
	Search(ctx context.Context, query string, vector []float32, alpha float32) ([]SearchResult, error)
}

type Reranker interface {
	Rerank(ctx context.Context, query string, docs []string) ([]int, error)
}

type Service struct {
	embedder Embedder
	store    VectorStore
	reranker Reranker
	logger   *QueryLogger
}

func NewService(e Embedder, s VectorStore, r Reranker, l *QueryLogger) *Service {
	return &Service{embedder: e, store: s, reranker: r, logger: l}
}

func (s *Service) Search(ctx context.Context, query string) ([]SearchResult, error) {
	start := time.Now()
	var finalDocs []SearchResult
	var err error

	defer func() {
		if s.logger != nil && err == nil {
			s.logger.Log(QueryLogEntry{
				Query:      query,
				NumResults: len(finalDocs),
				Duration:   time.Since(start),
			})
		}
	}()

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
		// Extract content for reranker
		contents := make([]string, len(docs))
		for i, d := range docs {
			contents[i] = d.Content
		}

		indices, err := s.reranker.Rerank(ctx, query, contents)
		if err != nil {
			return nil, err
		}
		
		reranked := make([]SearchResult, len(indices))
		for i, idx := range indices {
			if idx < len(docs) {
				reranked[i] = docs[idx]
			}
		}
		finalDocs = reranked
		return reranked, nil
	}

	finalDocs = docs
	return docs, nil
}