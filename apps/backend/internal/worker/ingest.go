package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nsqio/go-nsq"
)

type Chunk struct {
	Content   string
	Vector    []float32
	SourceURL string
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

type IngestHandler struct {
	fetcher  Fetcher
	embedder Embedder
	store    VectorStore
}

func NewIngestHandler(f Fetcher, e Embedder, s VectorStore) *IngestHandler {
	return &IngestHandler{fetcher: f, embedder: e, store: s}
}

func (h *IngestHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	var payload struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(m.Body, &payload); err != nil {
		log.Printf("Invalid message format: %v", err)
		return nil // Don't retry invalid messages
	}

	ctx := context.Background()

	// 1. Fetch
	content, err := h.fetcher.Fetch(ctx, payload.URL)
	if err != nil {
		log.Printf("Fetch failed: %v", err)
		return err // Retry
	}

	// 2. Chunk (Simplified: 1 chunk for now)
	
	// 3. Embed
	vector, err := h.embedder.Embed(ctx, content)
	if err != nil {
		log.Printf("Embed failed: %v", err)
		return err
	}

	// 4. Store
	chunk := Chunk{
		Content:   content,
		Vector:    vector,
		SourceURL: payload.URL,
	}
	if err := h.store.StoreChunk(ctx, chunk); err != nil {
		log.Printf("Store failed: %v", err)
		return err
	}

	log.Printf("Ingested: %s", payload.URL)
	return nil
}
