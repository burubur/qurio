package worker

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nsqio/go-nsq"
	"qurio/apps/backend/internal/crawler"
	"qurio/apps/backend/internal/text"
)

type Chunk struct {
	Content    string
	Vector     []float32
	SourceURL  string
	SourceID   string
	ChunkIndex int
}

type Crawler interface {
	Crawl(startURL string) ([]crawler.Page, error)
}

type CrawlerFactory func(cfg crawler.Config) (Crawler, error)

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

type VectorStore interface {
	StoreChunk(ctx context.Context, chunk Chunk) error
}

type ContentProcessor interface {
	Process(ctx context.Context, filename string, content []byte) (string, error)
}

type Producer interface {
	Publish(topic string, body []byte) error
}

type SourceStatusUpdater interface {
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateBodyHash(ctx context.Context, id, hash string) error
}

type IngestHandler struct {
	crawlerFactory CrawlerFactory
	processor      ContentProcessor
	embedder       Embedder
	store          VectorStore
	producer       Producer
	updater        SourceStatusUpdater
}

func NewIngestHandler(cf CrawlerFactory, cp ContentProcessor, e Embedder, s VectorStore, p Producer, u SourceStatusUpdater) *IngestHandler {
	return &IngestHandler{crawlerFactory: cf, processor: cp, embedder: e, store: s, producer: p, updater: u}
}

func (h *IngestHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	var payload struct {
		URL        string   `json:"url"`
		ID         string   `json:"id"`
		MaxDepth   int      `json:"max_depth"`
		Exclusions []string `json:"exclusions"`
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

	// 1. Configure Crawler
	cfg := crawler.Config{
		MaxDepth:   payload.MaxDepth,
		Exclusions: payload.Exclusions,
	}
	c, err := h.crawlerFactory(cfg)
	if err != nil {
		slog.Error("failed to create crawler", "error", err)
		return err
	}

	// 2. Crawl
	pages, err := c.Crawl(payload.URL)
	if err != nil {
		slog.Error("crawl failed", "url", payload.URL, "error", err)
		// Mark failed if crawl fails completely
		_ = h.updater.UpdateStatus(ctx, payload.ID, "failed")
		return err
	}

	if len(pages) == 0 {
		slog.Warn("no pages found", "url", payload.URL)
		_ = h.updater.UpdateStatus(ctx, payload.ID, "completed")
		return nil
	}

	totalChunks := 0

	failedPages := 0
	for _, page := range pages {
		// 3. Process Content (Docling)
		content, err := h.processor.Process(ctx, page.URL, []byte(page.Content))
		if err != nil {
			slog.Warn("process failed", "url", page.URL, "error", err)
			failedPages++
			continue
		}

		// 1.5 Update Hash (Only for root URL)
		if payload.ID != "" && page.URL == payload.URL {
			hash := sha256.Sum256([]byte(content))
			hashStr := fmt.Sprintf("%x", hash)
			if err := h.updater.UpdateBodyHash(ctx, payload.ID, hashStr); err != nil {
				slog.Warn("failed to update body hash", "error", err)
			}
		}

		// 4. Chunk
		chunks := text.Chunk(content, 512, 50)
		if len(chunks) == 0 {
			continue
		}

		for i, c := range chunks {
			// 5. Embed
			vector, err := h.embedder.Embed(ctx, c)
			if err != nil {
				slog.Error("embed failed", "error", err)
				// Determine if we should fail the whole batch or just this chunk?
				// For now, fail hard on embed/store errors to trigger retry of the message
				return err
			}

			// 6. Store
			chunk := Chunk{
				Content:    c,
				Vector:     vector,
				SourceURL:  page.URL,
				SourceID:   payload.ID,
				ChunkIndex: i,
			}
			if err := h.store.StoreChunk(ctx, chunk); err != nil {
				slog.Error("store failed", "error", err)
				return err
			}
		}
		totalChunks += len(chunks)
	}

	slog.Info("ingested pages", "pages", len(pages), "failed_pages", failedPages, "total_chunks", totalChunks, "url", payload.URL)

	// Determine final status
	finalStatus := "completed"
	if failedPages > 0 && totalChunks == 0 {
		finalStatus = "failed"
	} else if failedPages > 0 {
		finalStatus = "completed_with_errors" // Or just completed for MVP
	}
	
	// Success
	if payload.ID != "" {
		if err := h.updater.UpdateStatus(ctx, payload.ID, finalStatus); err != nil {
			slog.Warn("failed to update status", "error", err)
			// Non-critical error, we still succeeded ingestion
		}
	}
	
	return nil
}