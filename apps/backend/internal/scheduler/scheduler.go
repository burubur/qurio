package scheduler

import (
	"context"
	"log/slog"
	"time"

	"qurio/apps/backend/features/source"
)

type Scheduler struct {
	sourceRepo source.Repository
	service    *source.Service
	stop       chan struct{}
}

func New(repo source.Repository, service *source.Service) *Scheduler {
	return &Scheduler{
		sourceRepo: repo,
		service:    service,
		stop:       make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	// Run every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		defer ticker.Stop()
		slog.Info("Scheduler started")
		for {
			select {
			case <-ticker.C:
				s.checkAndSync(context.Background())
			case <-s.stop:
				slog.Info("Scheduler stopped")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}

func (s *Scheduler) checkAndSync(ctx context.Context) {
	sources, err := s.sourceRepo.ListSyncDue(ctx)
	if err != nil {
		slog.Error("Scheduler: failed to list sources", "error", err)
		return
	}

	for _, src := range sources {
		if s.isDue(src) {
			slog.Info("Scheduler: triggering sync", "source_id", src.ID, "schedule", src.SyncSchedule)

			// Trigger re-sync
			if err := s.service.ReSync(ctx, src.ID); err != nil {
				slog.Error("Scheduler: failed to resync source", "id", src.ID, "error", err)
				continue
			}

			// Update LastSyncedAt
			if err := s.sourceRepo.UpdateLastSyncedAt(ctx, src.ID, time.Now()); err != nil {
				slog.Error("Scheduler: failed to update last_synced_at", "id", src.ID, "error", err)
			}
		}
	}
}

func (s *Scheduler) isDue(src source.Source) bool {
	if src.LastSyncedAt == nil {
		return true // Never synced, sync now? Or wait? Let's say sync now.
	}

	last := *src.LastSyncedAt
	now := time.Now()

	switch src.SyncSchedule {
	case "minute":
		return now.Sub(last) >= time.Minute
	case "hourly":
		return now.Sub(last) >= time.Hour
	case "daily":
		return now.Sub(last) >= 24*time.Hour
	// Add more intervals if needed
	default:
		// Default daily
		return now.Sub(last) >= 24*time.Hour
	}
}
