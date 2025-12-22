package worker

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nsqio/go-nsq"
	"qurio/apps/backend/internal/text"
)

type Chunk struct {
	Content    string
	Vector     []float32
	SourceURL  string
	ChunkIndex int
}

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

type VectorStore interface {
	StoreChunk(ctx context.Context, chunk Chunk) error
}

type Fetcher interface {
	Fetch(ctx context.Context, url string) (string, error)
}

type Producer interface {
	Publish(topic string, body []byte) error
}

type SourceStatusUpdater interface {
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateBodyHash(ctx context.Context, id, hash string) error
}

type IngestHandler struct {
	fetcher  Fetcher
	embedder Embedder
	store    VectorStore
	producer Producer
	updater  SourceStatusUpdater
}

func NewIngestHandler(f Fetcher, e Embedder, s VectorStore, p Producer, u SourceStatusUpdater) *IngestHandler {
	return &IngestHandler{fetcher: f, embedder: e, store: s, producer: p, updater: u}
}

func (h *IngestHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	var payload struct {
		URL string `json:"url"`
		ID  string `json:"id"`
	}
	if err := json.Unmarshal(m.Body, &payload); err != nil {
		slog.Error("invalid message format", "error", err)
		return nil // Don't retry invalid messages
	}

	ctx := context.Background()

	// 0. Retry Limit / DLQ
	if m.Attempts > 3 {
		slog.Warn("message exceeded max attempts", "id", m.ID, "attempts", m.Attempts, "action", "dlq")
		h.updater.UpdateStatus(ctx, payload.ID, "failed") // Mark as failed
		if err := h.producer.Publish("ingestion_dlq", m.Body); err != nil {
			slog.Error("failed to publish to DLQ", "error", err)
			return err // Retry publishing to DLQ
		}
		return nil // Ack original message
	}

	// Update status to processing
	if payload.ID != "" {
		_ = h.updater.UpdateStatus(ctx, payload.ID, "processing")
	}

	// 1. Fetch
	content, err := h.fetcher.Fetch(ctx, payload.URL)
	if err != nil {
		slog.Error("fetch failed", "url", payload.URL, "error", err)
		// Don't mark failed yet, let NSQ retry
		return err 
	}

	// 1.5 Update Hash
	if payload.ID != "" {
		hash := sha256.Sum256([]byte(content))
		hashStr := fmt.Sprintf("%x", hash)
		if err := h.updater.UpdateBodyHash(ctx, payload.ID, hashStr); err != nil {
			slog.Warn("failed to update body hash", "error", err)
		}
	}

	// 2. Chunk
	chunks := text.Chunk(content, 512, 50)
	if len(chunks) == 0 {
		slog.Warn("no chunks generated", "url", payload.URL)
		_ = h.updater.UpdateStatus(ctx, payload.ID, "completed") // Or warning?
		return nil
	}

	for i, c := range chunks {
		// 3. Embed
		vector, err := h.embedder.Embed(ctx, c)
		if err != nil {
			slog.Error("embed failed", "error", err)
			return err
		}

		// 4. Store
		chunk := Chunk{
			Content:    c,
			Vector:     vector,
			SourceURL:  payload.URL,
			ChunkIndex: i,
		}
		if err := h.store.StoreChunk(ctx, chunk); err != nil {
			slog.Error("store failed", "error", err)
			return err
		}
	}

	slog.Info("ingested chunks", "count", len(chunks), "url", payload.URL)
	
	// Success
	if payload.ID != "" {
		if err := h.updater.UpdateStatus(ctx, payload.ID, "completed"); err != nil {
			slog.Warn("failed to update status", "error", err)
			// Non-critical error, we still succeeded ingestion
		}
	}
	
	return nil
}