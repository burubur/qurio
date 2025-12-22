package source_test

import (
	"context"
	"fmt"
	"testing"
	"crypto/sha256"

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

func (m *MockRepo) ExistsByHash(ctx context.Context, hash string) (bool, error) {
	args := m.Called(ctx, hash)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) List(ctx context.Context) ([]source.Source, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]source.Source), args.Error(1)
}

func (m *MockRepo) UpdateStatus(ctx context.Context, id, status string) error {
	args := m.Called(ctx, id, status)
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
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(src.URL)))

	// Expect ExistsByHash -> false
	repo.On("ExistsByHash", mock.Anything, hash).Return(false, nil)
	
	// Expect Save -> success
	repo.On("Save", mock.Anything, mock.MatchedBy(func(s *source.Source) bool {
		return s.URL == src.URL && s.ContentHash == hash
	})).Return(nil)
	
	// Expect Publish -> success
	pub.On("Publish", "ingest", mock.Anything).Return(nil)
	
	err := svc.Create(context.Background(), src)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestCreateSource_Duplicate(t *testing.T) {
	repo := new(MockRepo)
	pub := new(MockPublisher)
	svc := source.NewService(repo, pub)
	
	src := &source.Source{URL: "https://duplicate.com"}
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(src.URL)))

	// Expect ExistsByHash -> true
	repo.On("ExistsByHash", mock.Anything, hash).Return(true, nil)
	
	err := svc.Create(context.Background(), src)
	
	assert.Error(t, err)
	assert.Equal(t, "Duplicate detected", err.Error())
	
	// Save and Publish should NOT be called
	repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	pub.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything)
}