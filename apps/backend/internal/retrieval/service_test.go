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
func (m *MockStore) Search(ctx context.Context, query string, vector []float32, alpha float32) ([]string, error) {
	args := m.Called(ctx, query, vector, alpha)
	return args.Get(0).([]string), args.Error(1)
}

func TestSearch(t *testing.T) {
	e := new(MockEmbedder)
	s := new(MockStore)
	svc := retrieval.NewService(e, s, nil)

	ctx := context.Background()
	e.On("Embed", ctx, "test").Return([]float32{0.1}, nil)
	// Verify alpha is 0.5
	s.On("Search", ctx, "test", []float32{0.1}, float32(0.5)).Return([]string{"result"}, nil)

	res, err := svc.Search(ctx, "test")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}
