package worker_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"qurio/apps/backend/internal/worker"
)

// Mocks
type MockFetcher struct { mock.Mock }
func (m *MockFetcher) Fetch(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

type MockEmbedder struct { mock.Mock }
func (m *MockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	return args.Get(0).([]float32), args.Error(1)
}

type MockStore struct { mock.Mock }
func (m *MockStore) StoreChunk(ctx context.Context, chunk worker.Chunk) error {
	args := m.Called(ctx, chunk)
	return args.Error(0)
}

type MockProducer struct { mock.Mock }
func (m *MockProducer) Publish(topic string, body []byte) error {
	args := m.Called(topic, body)
	return args.Error(0)
}

type MockUpdater struct { mock.Mock }
func (m *MockUpdater) UpdateStatus(ctx context.Context, id, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func createMessage(body []byte) *nsq.Message {
	msg := nsq.NewMessage(nsq.MessageID{'1','2','3','4','5','6','7','8','9','0','1','2','3','4','5','6'}, body)
	return msg
}

func TestHandleMessage_Success(t *testing.T) {
	f := new(MockFetcher)
	e := new(MockEmbedder)
	s := new(MockStore)
	p := new(MockProducer)
	u := new(MockUpdater)
	h := worker.NewIngestHandler(f, e, s, p, u)

	payload := []byte(`{"url":"http://test.com", "id":"123"}`)
	msg := createMessage(payload)

	// Expect Status Update: Processing
	u.On("UpdateStatus", mock.Anything, "123", "processing").Return(nil)
	
	f.On("Fetch", mock.Anything, "http://test.com").Return("content", nil)
	e.On("Embed", mock.Anything, "content").Return([]float32{0.1, 0.2}, nil)
	s.On("StoreChunk", mock.Anything, mock.MatchedBy(func(c worker.Chunk) bool {
		return c.Content == "content" && c.SourceURL == "http://test.com"
	})).Return(nil)

	// Expect Status Update: Completed
	u.On("UpdateStatus", mock.Anything, "123", "completed").Return(nil)

	err := h.HandleMessage(msg)
	assert.NoError(t, err)
}

func TestHandleMessage_MultiChunk(t *testing.T) {
	f := new(MockFetcher)
	e := new(MockEmbedder)
	s := new(MockStore)
	p := new(MockProducer)
	u := new(MockUpdater)
	h := worker.NewIngestHandler(f, e, s, p, u)

	payload := []byte(`{"url":"http://test.com", "id":"123"}`)
	msg := createMessage(payload)

	u.On("UpdateStatus", mock.Anything, "123", "processing").Return(nil)

	contentBuilder := strings.Builder{}
	for i := 0; i < 520; i++ {
		contentBuilder.WriteString("word ")
	}
	content := contentBuilder.String()

	f.On("Fetch", mock.Anything, "http://test.com").Return(content, nil)
	
	e.On("Embed", mock.Anything, mock.AnythingOfType("string")).Return([]float32{0.1}, nil).Times(2)
	s.On("StoreChunk", mock.Anything, mock.MatchedBy(func(c worker.Chunk) bool {
		return len(c.Content) > 0 && c.SourceURL == "http://test.com"
	})).Return(nil).Times(2)

	u.On("UpdateStatus", mock.Anything, "123", "completed").Return(nil)

	err := h.HandleMessage(msg)
	assert.NoError(t, err)
	e.AssertNumberOfCalls(t, "Embed", 2)
	s.AssertNumberOfCalls(t, "StoreChunk", 2)
}

func TestHandleMessage_Retry(t *testing.T) {
	f := new(MockFetcher)
	e := new(MockEmbedder)
	s := new(MockStore)
	p := new(MockProducer)
	u := new(MockUpdater)
	h := worker.NewIngestHandler(f, e, s, p, u)

	payload := []byte(`{"url":"http://test.com", "id":"123"}`)
	msg := createMessage(payload)
	msg.Attempts = 1

	u.On("UpdateStatus", mock.Anything, "123", "processing").Return(nil)
	f.On("Fetch", mock.Anything, "http://test.com").Return("", errors.New("network error"))

	err := h.HandleMessage(msg)
	assert.Error(t, err) // Should return error to trigger requeue
}

func TestHandleMessage_DLQ(t *testing.T) {
	f := new(MockFetcher)
	e := new(MockEmbedder)
	s := new(MockStore)
	p := new(MockProducer)
	u := new(MockUpdater)
	h := worker.NewIngestHandler(f, e, s, p, u)

	payload := []byte(`{"url":"http://test.com", "id":"123"}`)
	msg := createMessage(payload)
	msg.Attempts = 4 // Exceeds limit of 3

	// Should update to failed
	u.On("UpdateStatus", mock.Anything, "123", "failed").Return(nil)
	p.On("Publish", "ingestion_dlq", payload).Return(nil)

	err := h.HandleMessage(msg)
	assert.NoError(t, err) // Should return nil to Ack original message
	p.AssertExpectations(t)
	// Fetch/Embed/Store should NOT be called
	f.AssertNotCalled(t, "Fetch", mock.Anything, mock.Anything)
}
