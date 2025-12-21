package source

import (
	"context"
	"encoding/json"
	"log"
)

type Source struct {
	ID  string
	URL string
}

type Repository interface {
	Save(ctx context.Context, src *Source) error
}

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

func (s *Service) Create(ctx context.Context, src *Source) error {
	// 1. Save to DB
	if err := s.repo.Save(ctx, src); err != nil {
		return err
	}

	// 2. Publish to NSQ
	payload, _ := json.Marshal(map[string]string{"url": src.URL})
	if err := s.pub.Publish("ingest", payload); err != nil {
		log.Printf("Failed to publish ingest event: %v", err)
		// We don't fail the request if publishing fails, but in a real system we might want outbox pattern
	}
	
	return nil
}
