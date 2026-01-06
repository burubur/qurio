package app

import (
	"context"

	"qurio/apps/backend/internal/retrieval"
	"qurio/apps/backend/internal/worker"
)

// MockVectorStore implements VectorStore for testing.
type MockVectorStore struct {
	StoreChunkErr        error
	DeleteChunksByURLErr error
	DeleteChunksErr      error
	SearchRes            []retrieval.SearchResult
	SearchErr            error
	GetChunksRes         []worker.Chunk
	GetChunksErr         error
	GetChunksByURLRes    []retrieval.SearchResult
	GetChunksByURLErr    error
	CountChunksRes       int
	CountChunksErr       error
	EnsureSchemaErr      error
}

func (m *MockVectorStore) StoreChunk(ctx context.Context, chunk worker.Chunk) error {
	return m.StoreChunkErr
}

func (m *MockVectorStore) DeleteChunksByURL(ctx context.Context, sourceID, url string) error {
	return m.DeleteChunksByURLErr
}

func (m *MockVectorStore) DeleteChunksBySourceID(ctx context.Context, sourceID string) error {
	return m.DeleteChunksErr
}

func (m *MockVectorStore) Search(ctx context.Context, query string, vector []float32, alpha float32, limit int, searchFilters map[string]interface{}) ([]retrieval.SearchResult, error) {
	return m.SearchRes, m.SearchErr
}

func (m *MockVectorStore) GetChunks(ctx context.Context, sourceID string) ([]worker.Chunk, error) {
	return m.GetChunksRes, m.GetChunksErr
}

func (m *MockVectorStore) GetChunksByURL(ctx context.Context, url string) ([]retrieval.SearchResult, error) {
	return m.GetChunksByURLRes, m.GetChunksByURLErr
}

func (m *MockVectorStore) CountChunks(ctx context.Context) (int, error) {
	return m.CountChunksRes, m.CountChunksErr
}

func (m *MockVectorStore) EnsureSchema(ctx context.Context) error {
	return m.EnsureSchemaErr
}

// MockTaskPublisher implements TaskPublisher for testing.
type MockTaskPublisher struct {
	PublishErr error
}

func (m *MockTaskPublisher) Publish(topic string, body []byte) error {
	return m.PublishErr
}

// MockDatabase is not needed as we can use sqlmock to generate a *sql.DB that satisfies the interface.
// However, if we need a custom struct for some reason, we would face issues returning *sql.Row.
