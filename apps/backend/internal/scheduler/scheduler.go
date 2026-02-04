package scheduler

import (
	"context"
	"log/slog"
	"qurio/apps/backend/features/source"
	"time"

	"github.com/robfig/cron/v3"
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
		return true
	}

	last := *src.LastSyncedAt
	now := time.Now()

	scheduleStr := src.SyncSchedule
	// Legacy mapping
	switch scheduleStr {
	case "minute":
		scheduleStr = "* * * * *"
	case "hourly":
		scheduleStr = "@hourly"
	case "daily":
		scheduleStr = "@daily"
	case "":
		scheduleStr = "@daily" // Default
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schedule, err := parser.Parse(scheduleStr)
	if err != nil {
		slog.Warn("Scheduler: invalid cron schedule, fallback to daily", "schedule", scheduleStr, "error", err)
		// Fallback to daily
		return now.Sub(last) >= 24*time.Hour
	}

	nextSyncTime := schedule.Next(last)
	return now.After(nextSyncTime)
}
