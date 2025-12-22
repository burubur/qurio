package reranker_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"qurio/apps/backend/internal/adapter/reranker"
)

func TestRerank_Mock(t *testing.T) {
	c := reranker.NewClient("other", "key")
	docs := []string{"doc1", "doc2"}
	
	reranked, err := c.Rerank(context.Background(), "query", docs)
	
	assert.NoError(t, err)
	assert.Equal(t, docs, reranked)
}

func TestRerank_Jina(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			w.WriteHeader(401)
			return
		}
		
		// Return docs reversed: B (0.9), A (0.1)
		// Input was A, B. So index 1 is B, index 0 is A.
		w.WriteHeader(200)
		w.Write([]byte(`{"results": [{"index": 1, "relevance_score": 0.9}, {"index": 0, "relevance_score": 0.1}]}`))
	}))
	defer server.Close()
	
	c := reranker.NewClient("jina", "test-key")
	c.SetBaseURL(server.URL)
	
	docs := []string{"A", "B"}
	reranked, err := c.Rerank(context.Background(), "q", docs)
	
	assert.NoError(t, err)
	assert.Equal(t, []string{"B", "A"}, reranked)
}