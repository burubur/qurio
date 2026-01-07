### Task 1: Middleware - Trace Chain Integrity

**Files:**
- Create/Modify: `apps/backend/internal/middleware/correlation_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Middleware injects `X-Correlation-ID` header if missing.
  2. Middleware preserves existing `X-Correlation-ID` header.
  3. `GetCorrelationID` helper correctly extracts ID from context.
  4. ContextHandler (if exists) logs include the ID.

- **Functional Requirements**
  1. Request handling must not fail if ID is missing.

- **Non-Functional Requirements**
  1. Zero allocation (or minimal) for ID generation path if possible.

- **Test Coverage**
  - [Unit] `CorrelationID` middleware - missing vs present header.
  - [Unit] `GetCorrelationID` - extraction logic.

**Step 1: Write failing test**
```go
package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"qurio/apps/backend/internal/middleware"
)

func TestCorrelationID_Middleware(t *testing.T) {
	tests := []struct {
		name           string
		incomingHeader string
		expectHeader   bool
		expectSameID   bool
	}{
		{
			name:           "Should Generate ID When Missing",
			incomingHeader: "",
			expectHeader:   true,
			expectSameID:   false,
		},
		{
			name:           "Should Preserve Existing ID",
			incomingHeader: "test-correlation-id-123",
			expectHeader:   true,
			expectSameID:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.incomingHeader != "" {
				req.Header.Set("X-Correlation-ID", tt.incomingHeader)
			}
			rec := httptest.NewRecorder()

			handler := middleware.CorrelationID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				id := middleware.GetCorrelationID(r.Context())
				if tt.expectHeader {
					assert.NotEmpty(t, id)
				}
				if tt.expectSameID {
					assert.Equal(t, tt.incomingHeader, id)
				}
			}))

			handler.ServeHTTP(rec, req)

			// Check Response Header
			respHeader := rec.Header().Get("X-Correlation-ID")
			if tt.expectHeader {
				assert.NotEmpty(t, respHeader)
			}
			if tt.expectSameID {
				assert.Equal(t, tt.incomingHeader, respHeader)
			}
		})
	}
}
```

**Step 2: Verify test fails**
Run: `go test apps/backend/internal/middleware/correlation_test.go -v`
Expected: FAIL (if middleware logic is missing or broken)

**Step 3: Write minimal implementation**
```go
// apps/backend/internal/middleware/correlation.go
package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const correlationIDKey contextKey = "correlationID"
const CorrelationHeader = "X-Correlation-ID"

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(CorrelationHeader)
		if id == "" {
			id = uuid.New().String()
		}

		// Set header for response
		w.Header().Set(CorrelationHeader, id)

		// Add to context
		ctx := context.WithValue(r.Context(), correlationIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDKey).(string); ok {
		return id
	}
	return ""
}
```

**Step 4: Verify test passes**
Run: `go test apps/backend/internal/middleware/correlation_test.go -v`
Expected: PASS

---

### Task 2: Retrieval Service - Table-Driven Unit Tests

**Files:**
- Create/Modify: `apps/backend/internal/retrieval/service_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `Search` method logic is verified in isolation.
  2. Alpha tuning parameters are respected (overrides defaults).
  3. Metadata (Author, CreatedAt) is preserved in results.
  4. Reranker is called if configured.

- **Test Coverage**
  - [Unit] `Service.Search` - Table driven:
    - Default settings
    - Custom Alpha/Limit
    - Reranker enabled/disabled
    - Embedder error
    - Store error

**Step 1: Write failing test**
```go
package retrieval_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"qurio/apps/backend/internal/retrieval"
	"qurio/apps/backend/internal/settings"
)

// Mocks defined here or imported from internal/testutils
// ... (MockEmbedder, MockVectorStore, MockReranker, MockSettingsRepo)

func TestService_Search_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		opts          *retrieval.SearchOptions
		setupMocks    func(*MockEmbedder, *MockVectorStore, *MockReranker, *MockSettingsRepo)
		expectedLen   int
		expectedTitle string
		expectError   bool
	}{
		{
			name:  "Basic Search Success",
			query: "golang",
			opts:  nil,
			setupMocks: func(e *MockEmbedder, s *MockVectorStore, r *MockReranker, set *MockSettingsRepo) {
				set.On("Get", mock.Anything).Return(&settings.Settings{SearchAlpha: 0.5, SearchTopK: 10}, nil)
				e.On("Embed", mock.Anything, "golang").Return([]float32{0.1, 0.2}, nil)
				s.On("Search", mock.Anything, "golang", []float32{0.1, 0.2}, float32(0.5), 10, mock.Anything).
					Return([]retrieval.SearchResult{{Content: "Go is great", Metadata: map[string]interface{}{"title": "Go Guide"}}}, nil)
			},
			expectedLen:   1,
			expectedTitle: "Go Guide",
			expectError:   false,
		},
		{
			name:  "Alpha Override",
			query: "python",
			opts:  &retrieval.SearchOptions{Alpha: ptr(0.9)},
			setupMocks: func(e *MockEmbedder, s *MockVectorStore, r *MockReranker, set *MockSettingsRepo) {
				set.On("Get", mock.Anything).Return(&settings.Settings{SearchAlpha: 0.5}, nil)
				e.On("Embed", mock.Anything, "python").Return([]float32{0.3}, nil)
				// Expect Alpha 0.9 passed to store
				s.On("Search", mock.Anything, "python", []float32{0.3}, float32(0.9), mock.Anything, mock.Anything).
					Return([]retrieval.SearchResult{}, nil)
			},
			expectedLen: 0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Mocks
			me := new(MockEmbedder)
			ms := new(MockVectorStore)
			mr := new(MockReranker)
			mset := new(MockSettingsRepo)
			
			tt.setupMocks(me, ms, mr, mset)
			
			svc := retrieval.NewService(me, ms, nil, settings.NewService(mset), nil) // Reranker nil for now
			
			results, err := svc.Search(context.Background(), tt.query, tt.opts)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, results, tt.expectedLen)
				if tt.expectedLen > 0 {
					assert.Equal(t, tt.expectedTitle, results[0].Title)
				}
			}
		})
	}
}

func ptr(v float32) *float32 { return &v }
```

**Step 2: Verify test fails**
Run: `go test apps/backend/internal/retrieval/service_test.go -v`
Expected: FAIL (or PASS if already implemented, this task expands coverage)

**Step 3: Write minimal implementation**
(Existing implementation in `service.go` should already pass this, this task is about *adding* the test coverage as requested)

**Step 4: Verify test passes**
Run: `go test apps/backend/internal/retrieval/service_test.go -v`
Expected: PASS

---

### Task 3: Worker - Poison Pill Unit Test

**Files:**
- Modify: `apps/backend/internal/worker/result_consumer_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Worker must not panic on malformed JSON.
  2. Worker must not panic on missing required fields (source_id, url).
  3. Worker must ACK the message (to remove from queue) even if processing fails (fail fast/discard).

- **Test Coverage**
  - [Unit] `HandleMessage` - Malformed JSON payload.
  - [Unit] `HandleMessage` - Valid JSON but missing SourceID.

**Step 1: Write failing test**
```go
func TestResultConsumer_HandleMessage_PoisonPill(t *testing.T) {
	// Setup consumer with mocks...
	// ...
	
	tests := []struct {
		name    string
		payload []byte
	}{
		{"Malformed JSON", []byte(`{"source_id": "1", "content": ...`)}, // Truncated
		{"Missing SourceID", []byte(`{"url": "http://example.com"}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &nsq.Message{Body: tt.payload}
			// Should return nil (ACK) even on error to prevent infinite redelivery of poison
			// Or return error if retry policy handles it? 
			// Spec says: "Worker acknowledges (acks) and drops these bad messages" -> Return nil
			err := consumer.HandleMessage(msg)
			assert.NoError(t, err)
		})
	}
}
```

**Step 2: Verify test fails**
Run: `go test apps/backend/internal/worker/result_consumer_test.go -v`
Expected: FAIL (if it currently returns error or panics)

**Step 3: Write minimal implementation**
```go
// In HandleMessage
var res IngestionResult
if err := json.Unmarshal(msg.Body, &res); err != nil {
    slog.Error("malformed message, dropping", "error", err)
    return nil // ACK to drop
}
if res.SourceID == "" || res.URL == "" {
    slog.Error("missing required fields, dropping", "sourceID", res.SourceID)
    return nil // ACK to drop
}
```

**Step 4: Verify test passes**
Run: `go test apps/backend/internal/worker/result_consumer_test.go -v`
Expected: PASS

---

### Task 4: Worker - Unskip Integration Test

**Files:**
- Modify: `apps/backend/internal/worker/integration_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `TestIngestIntegration` runs without skipping.
  2. Full flow (Produce NSQ -> Consume -> Weaviate -> Postgres) passes.

- **Test Coverage**
  - [Integration] `TestIngestIntegration` - Full happy path.

**Step 1: Write failing test**
Remove `t.Skip("skipping integration test")` from `apps/backend/internal/worker/integration_test.go`.

**Step 2: Verify test fails**
Run: `go test apps/backend/internal/worker/integration_test.go -v`
Expected: FAIL (likely due to environment setup if not running in proper containerized env, or logic bugs)

**Step 3: Write minimal implementation**
Fix any bugs revealed by the test (e.g. timeout issues, schema mismatch, or missing docker services). 
*Note: User must ensure `go test` is run in an environment with access to Weaviate/Postgres/NSQ (e.g. via `verify_infra.sh` or local docker-compose).*

**Step 4: Verify test passes**
Run: `go test apps/backend/internal/worker/integration_test.go -v`
Expected: PASS
