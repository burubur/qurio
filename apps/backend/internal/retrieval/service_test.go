package retrieval_test

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"qurio/apps/backend/internal/retrieval"
)

type MockEmbedder struct { mock.Mock }
func (m *MockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	return args.Get(0).([]float32), args.Error(1)
}

type MockStore struct { mock.Mock }
func (m *MockStore) Search(ctx context.Context, query string, vector []float32, alpha float32) ([]retrieval.SearchResult, error) {
	args := m.Called(ctx, query, vector, alpha)
	return args.Get(0).([]retrieval.SearchResult), args.Error(1)
}

type MockReranker struct { mock.Mock }
func (m *MockReranker) Rerank(ctx context.Context, query string, docs []string) ([]int, error) {
	args := m.Called(ctx, query, docs)
	return args.Get(0).([]int), args.Error(1)
}

func TestSearch_WithReranker(t *testing.T) {
	e := new(MockEmbedder)
	s := new(MockStore)
	r := new(MockReranker)
	svc := retrieval.NewService(e, s, r, nil)

	ctx := context.Background()
	e.On("Embed", ctx, "test").Return([]float32{0.1}, nil)
	
	initialResults := []retrieval.SearchResult{
		{Content: "A", Score: 0.5},
		{Content: "B", Score: 0.6},
	}
	s.On("Search", ctx, "test", []float32{0.1}, float32(0.5)).Return(initialResults, nil)
	
	// Reranker swaps them: [1, 0]
	r.On("Rerank", ctx, "test", []string{"A", "B"}).Return([]int{1, 0}, nil)

	res, err := svc.Search(ctx, "test")
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "B", res[0].Content)
	assert.Equal(t, "A", res[1].Content)
}

func TestSearch(t *testing.T) {
	e := new(MockEmbedder)
	s := new(MockStore)
	svc := retrieval.NewService(e, s, nil, nil)

	ctx := context.Background()
	e.On("Embed", ctx, "test").Return([]float32{0.1}, nil)
	
	expected := []retrieval.SearchResult{
		{Content: "result", Score: 0.9, Metadata: map[string]interface{}{"source": "doc1"}},
	}
	// Verify alpha is 0.5
	s.On("Search", ctx, "test", []float32{0.1}, float32(0.5)).Return(expected, nil)

	res, err := svc.Search(ctx, "test")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "doc1", res[0].Metadata["source"])
}
