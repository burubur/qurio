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
)

type MockRetriever struct { mock.Mock }
func (m *MockRetriever) Search(ctx context.Context, query string) ([]string, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]string), args.Error(1)
}

func TestHandleSearch(t *testing.T) {
	reqBody := `{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search","arguments":{"query":"test query"}},"id":1}`
	req := httptest.NewRequest("POST", "/mcp", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()
	
	r := new(MockRetriever)
	r.On("Search", mock.Anything, "test query").Return([]string{"result1"}, nil)
	
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
	assert.Equal(t, "result1", resp.Result.Content[0].Text)
}
