package source

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
)

type Source struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	ContentHash string `json:"-"`
	BodyHash    string `json:"-"`
	Status      string `json:"status"`
}

type Repository interface {
	Save(ctx context.Context, src *Source) error
	ExistsByHash(ctx context.Context, hash string) (bool, error)
	Get(ctx context.Context, id string) (*Source, error)
	List(ctx context.Context) ([]Source, error)
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateBodyHash(ctx context.Context, id, hash string) error
	SoftDelete(ctx context.Context, id string) error
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
	// 0. Compute Hash
	hash := sha256.Sum256([]byte(src.URL))
	src.ContentHash = fmt.Sprintf("%x", hash)

	// 1. Check Duplicate
	exists, err := s.repo.ExistsByHash(ctx, src.ContentHash)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Duplicate detected")
	}

	// 2. Save to DB
	if err := s.repo.Save(ctx, src); err != nil {
		return err
	}

	// 3. Publish to NSQ
	payload, _ := json.Marshal(map[string]string{"url": src.URL, "id": src.ID})
	if err := s.pub.Publish("ingest", payload); err != nil {
		log.Printf("Failed to publish ingest event: %v", err)
	}
	
	return nil
}

func (s *Service) List(ctx context.Context) ([]Source, error) {
	return s.repo.List(ctx)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *Service) ReSync(ctx context.Context, id string) error {
	src, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"url":    src.URL,
		"id":     src.ID,
		"resync": true,
	})
	if err := s.pub.Publish("ingest", payload); err != nil {
		log.Printf("Failed to publish resync event: %v", err)
		return err
	}
	return nil
}
