package app

import (
	"context"
	"database/sql"

	"qurio/apps/backend/internal/retrieval"
	"qurio/apps/backend/internal/worker"
)

// Database defines the SQL database operations required by the application.
type Database interface {
	PingContext(ctx context.Context) error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// VectorStore defines the vector database operations required by the application.
type VectorStore interface {
	StoreChunk(ctx context.Context, chunk worker.Chunk) error
	DeleteChunksByURL(ctx context.Context, sourceID, url string) error
	DeleteChunksBySourceID(ctx context.Context, sourceID string) error
	Search(ctx context.Context, query string, vector []float32, alpha float32, limit int, searchFilters map[string]interface{}) ([]retrieval.SearchResult, error)
	GetChunks(ctx context.Context, sourceID string) ([]worker.Chunk, error)
	GetChunksByURL(ctx context.Context, url string) ([]retrieval.SearchResult, error)
	CountChunks(ctx context.Context) (int, error)
	EnsureSchema(ctx context.Context) error
}

// TaskPublisher defines the message queue operations required by the application.
type TaskPublisher interface {
	Publish(topic string, body []byte) error
}
