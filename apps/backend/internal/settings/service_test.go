package settings

import (
	"context"
	"testing"
)

type MockRepo struct {
	settings *Settings
	err      error
}

func (m *MockRepo) Get(ctx context.Context) (*Settings, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.settings, nil
}

func (m *MockRepo) Update(ctx context.Context, s *Settings) error {
	if m.err != nil {
		return m.err
	}
	m.settings = s
	return nil
}

func TestGetSettings(t *testing.T) {
	mockRepo := &MockRepo{
		settings: &Settings{RerankProvider: "jina", RerankAPIKey: "key"},
	}
	svc := NewService(mockRepo)

	s, err := svc.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.RerankProvider != "jina" {
		t.Errorf("expected jina, got %s", s.RerankProvider)
	}
}

func TestUpdateSettings(t *testing.T) {
	mockRepo := &MockRepo{
		settings: &Settings{RerankProvider: "none"},
	}
	svc := NewService(mockRepo)

	newSettings := &Settings{RerankProvider: "cohere", RerankAPIKey: "newkey"}
	err := svc.Update(context.Background(), newSettings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mockRepo.settings.RerankProvider != "cohere" {
		t.Errorf("expected cohere, got %s", mockRepo.settings.RerankProvider)
	}
}
