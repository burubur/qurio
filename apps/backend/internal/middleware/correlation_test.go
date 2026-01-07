package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorrelationID(t *testing.T) {
	handler := CorrelationID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.Context().Value(CorrelationKey).(string)
		if !ok || id == "" {
			t.Error("correlation id missing from context")
		}
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("X-Correlation-ID") == "" {
		t.Error("header missing")
	}
}

func TestCorrelationID_PreservesExisting(t *testing.T) {
	existingID := "existing-trace-id"
	handler := CorrelationID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.Context().Value(CorrelationKey).(string)
		if !ok || id != existingID {
			t.Errorf("expected context id %s, got %s", existingID, id)
		}
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Correlation-ID", existingID)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("X-Correlation-ID") != existingID {
		t.Errorf("expected header %s, got %s", existingID, w.Header().Get("X-Correlation-ID"))
	}
}

func TestGetCorrelationID(t *testing.T) {
	ctx := context.Background()
	if id := GetCorrelationID(ctx); id != "unknown" {
		t.Errorf("expected unknown, got %s", id)
	}

	ctx = WithCorrelationID(ctx, "test-id")
	if id := GetCorrelationID(ctx); id != "test-id" {
		t.Errorf("expected test-id, got %s", id)
	}
}