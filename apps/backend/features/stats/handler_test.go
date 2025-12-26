package stats

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockSourceRepo struct{}
func (m *MockSourceRepo) Count(ctx context.Context) (int, error) { return 10, nil }

type MockJobRepo struct{}
func (m *MockJobRepo) Count(ctx context.Context) (int, error) { return 5, nil }

type MockVectorStore struct{}
func (m *MockVectorStore) CountChunks(ctx context.Context) (int, error) { return 100, nil }

func TestHandler_GetStats(t *testing.T) {
	h := NewHandler(&MockSourceRepo{}, &MockJobRepo{}, &MockVectorStore{})

	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()

	h.GetStats(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := body["data"]; !ok {
		t.Error("Response missing 'data' field")
	}
	
	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Data field is not a map")
	}

	if data["sources"].(float64) != 10 {
		t.Errorf("Expected sources 10, got %v", data["sources"])
	}
}
