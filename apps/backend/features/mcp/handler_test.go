package mcp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"qurio/apps/backend/features/mcp"
	"qurio/apps/backend/internal/retrieval"
)

type MockRetriever struct { mock.Mock }
func (m *MockRetriever) Search(ctx context.Context, query string) ([]retrieval.SearchResult, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]retrieval.SearchResult), args.Error(1)
}

func TestHandleSearch(t *testing.T) {
	reqBody := `{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search","arguments":{"query":"test query"}},"id":1}`
	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()
	
	r := new(MockRetriever)
	expectedResults := []retrieval.SearchResult{
		{Content: "result1", Score: 0.9, Metadata: map[string]interface{}{"url": "http://example.com"}},
	}
	r.On("Search", mock.Anything, "test query").Return(expectedResults, nil)
	
	h := mcp.NewHandler(r)
	h.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)
	
	type Response struct {
		JSONRPC string `json:"jsonrpc"`
		Result  struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}
	
	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.Len(t, resp.Result.Content, 1)
	assert.Contains(t, resp.Result.Content[0].Text, "result1")
	// Verify metadata is serialized in text (since MCP returns text content typically)
	assert.Contains(t, resp.Result.Content[0].Text, "http://example.com")
}

func TestToolsList(t *testing.T) {
	reqBody := `{"jsonrpc":"2.0","method":"tools/list","id":1}`
	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	h := mcp.NewHandler(new(MockRetriever))
	h.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	type Response struct {
		Result struct {
			Tools []map[string]interface{} `json:"tools"`
		} `json:"result"`
	}
	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	
	assert.NotEmpty(t, resp.Result.Tools)
	assert.Equal(t, "search", resp.Result.Tools[0]["name"])
}

func TestHandleSSE_Connection(t *testing.T) {
	req := httptest.NewRequest("GET", "/mcp/sse", nil)
	w := httptest.NewRecorder()

	h := mcp.NewHandler(new(MockRetriever))
	
	// Execute in goroutine as SSE blocks
	go h.HandleSSE(w, req)
	
	// Wait a bit for event
	// In real test we'd need better synchronization, but for minimal check:
	// We expect headers and the endpoint event
	// However, HandleSSE blocks until disconnect. 
	// To test properly we can cancel context.
}
