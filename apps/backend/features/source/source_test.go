package source_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"qurio/apps/backend/features/source"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Save(ctx context.Context, src *source.Source) error {
	args := m.Called(ctx, src)
	return args.Error(0)
}

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(topic string, body []byte) error {
	args := m.Called(topic, body)
	return args.Error(0)
}

func TestCreateSource(t *testing.T) {
	repo := new(MockRepo)
	pub := new(MockPublisher)
	svc := source.NewService(repo, pub)
	
	src := &source.Source{URL: "https://example.com"}
	repo.On("Save", mock.Anything, src).Return(nil)
	pub.On("Publish", "ingest", mock.Anything).Return(nil)
	
	err := svc.Create(context.Background(), src)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}
