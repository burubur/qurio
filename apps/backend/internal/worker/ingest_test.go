package worker_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"qurio/apps/backend/internal/worker"
	"qurio/apps/backend/internal/crawler"
)

// Mocks
type MockCrawler struct { mock.Mock }
func (m *MockCrawler) Crawl(startURL string) ([]crawler.Page, error) {
	args := m.Called(startURL)
	return args.Get(0).([]crawler.Page), args.Error(1)
}

type MockProcessor struct { mock.Mock }
func (m *MockProcessor) Process(ctx context.Context, filename string, content []byte) (string, error) {
	args := m.Called(ctx, filename, content)
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
func (m *MockUpdater) UpdateBodyHash(ctx context.Context, id, hash string) error {
	args := m.Called(ctx, id, hash)
	return args.Error(0)
}

func createMessage(body []byte) *nsq.Message {
	msg := nsq.NewMessage(nsq.MessageID{'1','2','3','4','5','6','7','8','9','0','1','2','3','4','5','6'}, body)
	return msg
}

func TestHandleMessage_Success(t *testing.T) {
	mc := new(MockCrawler)
	mp := new(MockProcessor)
	e := new(MockEmbedder)
	s := new(MockStore)
	p := new(MockProducer)
	u := new(MockUpdater)

	factory := func(cfg crawler.Config) (worker.Crawler, error) {
		return mc, nil
	}
	
	h := worker.NewIngestHandler(factory, mp, e, s, p, u)

	payload := []byte(`{"url":"http://test.com", "id":"123"}`)
	msg := createMessage(payload)

	// Expect Status Update: Processing
	u.On("UpdateStatus", mock.Anything, "123", "processing").Return(nil)
	u.On("UpdateBodyHash", mock.Anything, "123", mock.Anything).Return(nil)
	
	pages := []crawler.Page{{URL: "http://test.com", Content: "raw html"}}
	mc.On("Crawl", "http://test.com").Return(pages, nil)
	
	mp.On("Process", mock.Anything, "http://test.com", []byte("raw html")).Return("content", nil)
	e.On("Embed", mock.Anything, "content").Return([]float32{0.1, 0.2}, nil)
	s.On("StoreChunk", mock.Anything, mock.MatchedBy(func(c worker.Chunk) bool {
		return c.Content == "content" && c.SourceURL == "http://test.com" && c.SourceID == "123"
	})).Return(nil)

	// Expect Status Update: Completed
	u.On("UpdateStatus", mock.Anything, "123", "completed").Return(nil)

	err := h.HandleMessage(msg)
	assert.NoError(t, err)
}

func TestHandleMessage_MultiChunk(t *testing.T) {
	mc := new(MockCrawler)
	mp := new(MockProcessor)
	e := new(MockEmbedder)
	s := new(MockStore)
	p := new(MockProducer)
	u := new(MockUpdater)
	
	factory := func(cfg crawler.Config) (worker.Crawler, error) {
		return mc, nil
	}
	h := worker.NewIngestHandler(factory, mp, e, s, p, u)

	payload := []byte(`{"url":"http://test.com", "id":"123"}`)
	msg := createMessage(payload)

	u.On("UpdateStatus", mock.Anything, "123", "processing").Return(nil)
	u.On("UpdateBodyHash", mock.Anything, "123", mock.Anything).Return(nil)

	pages := []crawler.Page{
		{URL: "http://test.com", Content: "page1"},
		{URL: "http://test.com/sub", Content: "page2"},
	}
	mc.On("Crawl", "http://test.com").Return(pages, nil)

	mp.On("Process", mock.Anything, "http://test.com", []byte("page1")).Return("content1", nil)
	mp.On("Process", mock.Anything, "http://test.com/sub", []byte("page2")).Return("content2", nil)
	
	e.On("Embed", mock.Anything, mock.Anything).Return([]float32{0.1}, nil)
	
	s.On("StoreChunk", mock.Anything, mock.MatchedBy(func(c worker.Chunk) bool {
		return c.SourceID == "123"
	})).Return(nil).Times(2)

	u.On("UpdateStatus", mock.Anything, "123", "completed").Return(nil)

	err := h.HandleMessage(msg)
	assert.NoError(t, err)
	s.AssertNumberOfCalls(t, "StoreChunk", 2)
}

func TestHandleMessage_Retry(t *testing.T) {
	mc := new(MockCrawler)
	mp := new(MockProcessor)
	e := new(MockEmbedder)
	s := new(MockStore)
	p := new(MockProducer)
	u := new(MockUpdater)
	
	factory := func(cfg crawler.Config) (worker.Crawler, error) {
		return mc, nil
	}
	h := worker.NewIngestHandler(factory, mp, e, s, p, u)

	payload := []byte(`{"url":"http://test.com", "id":"123"}`)
	msg := createMessage(payload)
	msg.Attempts = 1

	u.On("UpdateStatus", mock.Anything, "123", "processing").Return(nil)
	// Expect Failed status
	u.On("UpdateStatus", mock.Anything, "123", "failed").Return(nil)

	mc.On("Crawl", "http://test.com").Return([]crawler.Page{}, errors.New("network error"))

	err := h.HandleMessage(msg)
	assert.Error(t, err) // Should return error to trigger requeue
}

func TestHandleMessage_DLQ(t *testing.T) {
	mc := new(MockCrawler)
	mp := new(MockProcessor)
	e := new(MockEmbedder)
	s := new(MockStore)
	p := new(MockProducer)
	u := new(MockUpdater)
	
	factory := func(cfg crawler.Config) (worker.Crawler, error) {
		return mc, nil
	}
	h := worker.NewIngestHandler(factory, mp, e, s, p, u)

	payload := []byte(`{"url":"http://test.com", "id":"123"}`)
	msg := createMessage(payload)
	msg.Attempts = 4 // Exceeds limit of 3

	// Should update to failed
	u.On("UpdateStatus", mock.Anything, "123", "failed").Return(nil)
	p.On("Publish", "ingestion_dlq", payload).Return(nil)

	err := h.HandleMessage(msg)
	assert.NoError(t, err) // Should return nil to Ack original message
	p.AssertExpectations(t)
	// Crawler should NOT be called
	mc.AssertNotCalled(t, "Crawl", mock.Anything)
}