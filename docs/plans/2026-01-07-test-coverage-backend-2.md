# Implementation Plan - Backend Test Coverage & Hardening (Phase 2)

This plan implements targeted test coverage improvements and hardening for critical backend components as identified in the `bug-testcoverage-backend-2` specification. It focuses on unit testing "unsealed rivets" (middleware, worker hardening, retrieval logic) and simulating adapter failures.

## Requirements Analysis

### Scope
- **Domain:** Backend Testing (Unit & Mock-based Integration)
- **Goal:** Improve test coverage and resilience handling.
- **Deliverables:**
  - Middleware Tests (`internal/middleware`)
  - Worker Hardening Tests (`internal/worker`)
  - Retrieval Service Tests (`internal/retrieval`)
  - Adapter Failure Tests (`internal/adapter`)

### Gap Analysis
- **Nouns:** `CorrelationID` (Middleware), `Poison Pill` (Worker Message), `Timeout` (Worker Context), `SearchOptions` (Retrieval), `GraphQL Error` (Weaviate).
- **Verbs:** `Inject` (Headers), `Handle` (Malformed JSON), `Rerank` (Search Results), `Simulate` (Network Failures).
- **Exclusions:** Full integration tests for Worker (covered in Plan 1).

### Knowledge Enrichment
- **RAG Queries:**
  - *Already performed in analysis phase:* Confirmed `httptest` usage for adapter simulation.
  - *Codebase Analysis:* Confirmed `link_discovery.go` is pure. Confirmed `store_integration_test.go` exists.

---

## Tasks

### Task 1: Middleware Unit Tests (CorrelationID)

**Files:**
- Modify: `apps/backend/internal/middleware/correlation_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Verify `X-Correlation-ID` header is preserved if present in request.
  2. Verify new UUID is generated if header is missing.
  3. Verify `X-Correlation-ID` is set in response header.
  4. Verify ID is populated in Request Context.

- **Test Coverage**
  - [Unit] `CorrelationID` middleware function.

**Step 1: Write failing test**
```go
// apps/backend/internal/middleware/correlation_test.go
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
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/middleware/... -v`
Expected: PASS (Logic likely exists, verifying coverage) or FAIL if logic is buggy.

**Step 3: Write minimal implementation**
(If validation fails, update `correlation.go` to use `r.Header.Get`)

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/middleware/... -v`
Expected: PASS

### Task 2: Worker Hardening Tests (Poison Pill)

**Files:**
- Modify: `apps/backend/internal/worker/result_consumer_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `HandleMessage` returns `nil` (ack) for malformed JSON body (Poison Pill).
  2. `HandleMessage` logs error but does not crash.
  3. `HandleMessage` handles empty body gracefully.

- **Test Coverage**
  - [Unit] `HandleMessage` - Malformed JSON input.

**Step 1: Write failing test**
```go
// apps/backend/internal/worker/result_consumer_test.go
func TestResultConsumer_HandleMessage_PoisonPill(t *testing.T) {
	// Mocks
	mockStore := new(MockVectorStore) // Assume mocks exist or generate them
	// ... setup other mocks ...
	consumer := NewResultConsumer(nil, mockStore, nil, nil, nil, nil, nil)

	// Malformed JSON
	msg := &nsq.Message{Body: []byte("{invalid-json")}
	
	// Should not panic, should return nil to ack
	err := consumer.HandleMessage(msg)
	assert.NoError(t, err)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/worker/... -v`
Expected: PASS (Current implementation likely handles it, this ensures it stays handled).

**Step 3: Write minimal implementation**
Ensure `json.Unmarshal` error check returns `nil`.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/worker/... -v`
Expected: PASS

### Task 3: Retrieval Service Table-Driven Tests

**Files:**
- Modify: `apps/backend/internal/retrieval/service_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `Search` combines results from `VectorStore` and `Reranker`.
  2. Metadata (Author, CreatedAt) is preserved in final results.
  3. Filters are passed correctly to `VectorStore`.

- **Test Coverage**
  - [Unit] `Service.Search` - Table-driven scenarios.

**Step 1: Write failing test**
```go
// apps/backend/internal/retrieval/service_test.go
func TestService_Search(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		opts     *SearchOptions
		mockSetup func(*MockEmbedder, *MockVectorStore, *MockReranker)
		wantLen  int
		wantErr  bool
	}{
		{
			name: "Success with Reranker",
			query: "test",
			opts: &SearchOptions{Limit: &[]int{5}[0]},
			mockSetup: func(e *MockEmbedder, s *MockVectorStore, r *MockReranker) {
				e.On("Embed", mock.Anything, "test").Return([]float32{0.1}, nil)
				s.On("Search", mock.Anything, "test", []float32{0.1}, float32(0.5), 5, mock.Anything).
					Return([]SearchResult{{Content: "B", Score: 0.8}, {Content: "A", Score: 0.9}}, nil)
				r.On("Rerank", mock.Anything, "test", []string{"B", "A"}).
					Return([]int{1, 0}, nil) // Swap order
			},
			wantLen: 2,
		},
	}
	// ... Run tests ...
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/retrieval/... -v`

**Step 3: Write minimal implementation**
Implement mocks and table loop.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/retrieval/... -v`

### Task 4: Adapter Failure Tests (Gemini)

**Files:**
- Modify: `apps/backend/internal/adapter/gemini/client_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `Embed` returns error on HTTP 503 (Service Unavailable).
  2. `Embed` returns error on HTTP 429 (Too Many Requests).

- **Test Coverage**
  - [Unit] `DynamicEmbedder.Embed` - Network failure scenarios.

**Step 1: Write failing test**
```go
// apps/backend/internal/adapter/gemini/client_test.go
func TestDynamicEmbedder_Embed_ServerFailure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()
	
	// Setup Embedder pointing to ts.URL
	// ...
	
	_, err := embedder.Embed(context.Background(), "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "503")
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/adapter/gemini/... -v`

**Step 3: Write minimal implementation**
(Implementation likely handles it via standard library, ensuring coverage).

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/adapter/gemini/... -v`

### Task 5: Adapter Failure Tests (Weaviate GraphQL)

**Files:**
- Modify: `apps/backend/internal/adapter/weaviate/store_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `Search` returns error when GraphQL response contains "errors" field (even with 200 OK).

- **Test Coverage**
  - [Unit] `Store.Search` - GraphQL error scenario.

**Step 1: Write failing test**
```go
// apps/backend/internal/adapter/weaviate/store_test.go
func TestStore_Search_GraphQLError(t *testing.T) {
	server := newMockWeaviateServer(t, func(r *http.Request, body map[string]interface{}) {
		// Verify request...
	})
	// Override handler to return GraphQL error
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []interface{}{
				map[string]interface{}{"message": "syntax error"},
			},
		})
	})
	defer server.Close()

	store := newTestStore(t, server)
	_, err := store.Search(context.Background(), "test", nil, 0.5, 10, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "syntax error")
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/adapter/weaviate/... -v`

**Step 3: Write minimal implementation**
Ensure `store.go` checks for `errors` in GraphQL response.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/adapter/weaviate/... -v`
