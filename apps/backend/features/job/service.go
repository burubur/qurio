package job

import (
	"context"
	"fmt"
	"time"
)

type EventPublisher interface {
	Publish(topic string, body []byte) error
}

type Service struct {
	repo Repository
	pub  EventPublisher
}

func NewService(repo Repository, pub EventPublisher) *Service {
	return &Service{repo: repo, pub: pub}
}

func (s *Service) List(ctx context.Context) ([]Job, error) {
	return s.repo.List(ctx)
}

func (s *Service) Retry(ctx context.Context, id string) error {
	// 1. Get Job
	job, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	// 2. Publish to NSQ with timeout
	done := make(chan error, 1)
	go func() {
		done <- s.pub.Publish("ingest.task", job.Payload)
	}()

	select {
	case err := <-done:
		if err != nil {
			return err
		}
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for NSQ publish")
	case <-ctx.Done():
		return ctx.Err()
	}

	// 3. Delete Job
	return s.repo.Delete(ctx, id)
}

func (s *Service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}
