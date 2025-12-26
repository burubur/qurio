package worker

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
	"qurio/apps/backend/features/job"
	"qurio/apps/backend/internal/middleware"
	"qurio/apps/backend/internal/text"
)

type ResultConsumer struct {
	embedder      Embedder
	store         VectorStore
	updater       SourceStatusUpdater
	jobRepo       job.Repository
	sourceFetcher SourceFetcher
}

func NewResultConsumer(e Embedder, s VectorStore, u SourceStatusUpdater, j job.Repository, sf SourceFetcher) *ResultConsumer {
	return &ResultConsumer{embedder: e, store: s, updater: u, jobRepo: j, sourceFetcher: sf}
}

func (h *ResultConsumer) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	var payload struct {
		SourceID      string `json:"source_id"`
		Content       string `json:"content"`
		URL           string `json:"url"`
		Status        string `json:"status,omitempty"` // "success" or "failed"
		Error         string `json:"error,omitempty"`
		CorrelationID string `json:"correlation_id,omitempty"`
	}
	if err := json.Unmarshal(m.Body, &payload); err != nil {
		slog.Error("invalid message format", "error", err)
		return nil // Don't retry invalid messages
	}

	correlationID := payload.CorrelationID
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	ctx := context.Background()
	ctx = middleware.WithCorrelationID(ctx, correlationID)
	
	if payload.Status == "failed" {
		slog.ErrorContext(ctx, "ingestion failed", "source_id", payload.SourceID, "error", payload.Error, "correlationId", correlationID)
		if err := h.updater.UpdateStatus(ctx, payload.SourceID, "failed"); err != nil {
			slog.WarnContext(ctx, "failed to update status to failed", "error", err, "correlationId", correlationID)
		}

		// Save Failed Job
		sType, sURL, err := h.sourceFetcher.GetSourceDetails(ctx, payload.SourceID)
		if err != nil {
			slog.ErrorContext(ctx, "failed to fetch source details for failed job", "error", err, "correlationId", correlationID)
		} else {
			jobPayload := map[string]interface{}{
				"type": sType,
				"id":   payload.SourceID,
			}
			if sType == "file" {
				jobPayload["path"] = sURL
			} else {
				jobPayload["url"] = sURL
			}

			pBytes, _ := json.Marshal(jobPayload)

			failedJob := &job.Job{
				SourceID: payload.SourceID,
				Handler:  sType,
				Payload:  json.RawMessage(pBytes),
				Error:    payload.Error,
			}
			if err := h.jobRepo.Save(ctx, failedJob); err != nil {
				slog.ErrorContext(ctx, "failed to save failed job", "error", err, "correlationId", correlationID)
			}
		}

		return nil
	}

	slog.InfoContext(ctx, "received result", "source_id", payload.SourceID, "content_len", len(payload.Content), "correlationId", correlationID)

	// 0. Delete Old Chunks (Idempotency)
	if payload.URL != "" {
		if err := h.store.DeleteChunksByURL(ctx, payload.SourceID, payload.URL); err != nil {
			slog.ErrorContext(ctx, "failed to delete old chunks", "error", err, "source_id", payload.SourceID, "url", payload.URL, "correlationId", correlationID)
			return err // Retry on error to ensure consistency
		}
	}

	// 1. Update Hash
	hash := sha256.Sum256([]byte(payload.Content))
	hashStr := fmt.Sprintf("%x", hash)
	if err := h.updater.UpdateBodyHash(ctx, payload.SourceID, hashStr); err != nil {
		slog.WarnContext(ctx, "failed to update body hash", "error", err, "correlationId", correlationID)
	}

	// 2. Chunk
	chunks := text.Chunk(payload.Content, 512, 50)
	if len(chunks) == 0 {
		slog.WarnContext(ctx, "no chunks generated", "source_id", payload.SourceID, "correlationId", correlationID)
		_ = h.updater.UpdateStatus(ctx, payload.SourceID, "completed")
		return nil
	}

	// 3. Embed & Store
	for i, c := range chunks {
		err := func() error {
			// Pass correlation ID to child context (though Context already has it via Value)
			// But creating new context with timeout might strip values if not careful?
			// context.WithTimeout derives from parent, so values are preserved.
			embedCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()

			vector, err := h.embedder.Embed(embedCtx, c)
			if err != nil {
				slog.ErrorContext(ctx, "embed failed", "error", err, "correlationId", correlationID)
				return err
			}

			chunk := Chunk{
				Content:    c,
				Vector:     vector,
				SourceID:   payload.SourceID,
				SourceURL:  payload.URL,
				ChunkIndex: i,
			}

			if err := h.store.StoreChunk(embedCtx, chunk); err != nil {
				slog.ErrorContext(ctx, "store failed", "error", err, "correlationId", correlationID)
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}

	slog.InfoContext(ctx, "stored chunks", "count", len(chunks), "source_id", payload.SourceID, "correlationId", correlationID)
	
	// 4. Update Status
	if err := h.updater.UpdateStatus(ctx, payload.SourceID, "completed"); err != nil {
		slog.WarnContext(ctx, "failed to update status", "error", err, "correlationId", correlationID)
	}

	return nil
}
