package mcp_test

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"qurio/apps/backend/features/mcp"
)

type MockRetriever struct {
	mock.Mock
}

func (m *MockRetriever) Search(ctx context.Context, query string) ([]string, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]string), args.Error(1)
}

func TestMCPCall(t *testing.T) {
	retriever := new(MockRetriever)
	handler := mcp.NewHandler(retriever)
	
	retriever.On("Search", mock.Anything, "test").Return([]string{"result"}, nil)
	
	reqBody := []byte(`{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search","arguments":{"query":"test"}},"id":1}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/mcp", bytes.NewBuffer(reqBody))
	
	handler.ServeHTTP(w, r)
	
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "result")
	retriever.AssertExpectations(t)
}
