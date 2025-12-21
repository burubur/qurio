package worker_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"qurio/apps/backend/internal/worker"
)

type MockEmbedder struct {
	mock.Mock
}

func (m *MockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	return args.Get(0).([]float32), args.Error(1)
}

type MockVectorStore struct {
	mock.Mock
}

func (m *MockVectorStore) StoreChunk(ctx context.Context, chunk worker.Chunk) error {
	args := m.Called(ctx, chunk)
	return args.Error(0)
}

type MockFetcher struct {
	mock.Mock
}

func (m *MockFetcher) Fetch(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

func TestHandleMessage(t *testing.T) {
	embedder := new(MockEmbedder)
	store := new(MockVectorStore)
	fetcher := new(MockFetcher)
	
	h := worker.NewIngestHandler(fetcher, embedder, store)
	
	url := "https://example.com"
	msgBody, _ := json.Marshal(map[string]string{"url": url})
	msg := &nsq.Message{Body: msgBody, ID: nsq.MessageID([16]byte{1})}
	
	// Mock expectation: Fetch -> Embed -> Store
	// Note: Chunking logic is simplified here (1 chunk per content) for MVP
	fetcher.On("Fetch", mock.Anything, url).Return("content", nil)
	embedder.On("Embed", mock.Anything, "content").Return([]float32{0.1, 0.2}, nil)
	store.On("StoreChunk", mock.Anything, mock.MatchedBy(func(c worker.Chunk) bool {
		return c.Content == "content"
	})).Return(nil)
	
	err := h.HandleMessage(msg)
	assert.NoError(t, err)
	
	fetcher.AssertExpectations(t)
	embedder.AssertExpectations(t)
	store.AssertExpectations(t)
}
