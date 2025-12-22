# Implementation Plan - MVP Part 3.1: Retrieval Pipeline & Agent Integration

**Scope:** Implement Advanced Retrieval (Hybrid Search with Reranking), Full MCP Compliance (tools/list), and Query Observability.

**Gap Analysis:**
- **Retrieval Service:** Missing metadata in return types; `alpha` parameter not exposed.
- **MCP Endpoint:** Missing `tools/list` implementation; `tools/call` response lacks metadata; logging is basic.
- **Query Logging:** Missing structured file logging (FR-5.3).
- **Frontend:** Missing Reranker configuration UI.

**Exclusions:**
- **Web Crawler:** Deferred to next plan (MVP Part 3.2).
- **Docling OCR:** Deferred to next plan.

***

### Task 1: Refactor Retrieval Types & VectorStore

**Files:**
- Modify: `apps/backend/internal/retrieval/service.go`
- Modify: `apps/backend/internal/adapter/weaviate/store.go` (Signature update)
- Modify: `apps/backend/features/source/source.go` (If it uses VectorStore, but likely distinct)
- Test: `apps/backend/internal/retrieval/service_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `SearchResult` struct defined with `Content`, `Score`, `Metadata`.
  2. `VectorStore.Search` returns `[]SearchResult` instead of `[]string`.
  3. `Service.Search` returns `[]SearchResult`.

- **Test Coverage**
  - [Unit] `Service.Search` - mocks VectorStore returning `SearchResult`s.
  - [Integration] Weaviate adapter returns populated `SearchResult`s.

**Step 1: Write failing test**
```go
// apps/backend/internal/retrieval/service_test.go
package retrieval_test

import (
	"context"
	"testing"
	"qurio/apps/backend/internal/retrieval"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStore struct { mock.Mock }
func (m *MockStore) Search(ctx context.Context, query string, vector []float32, alpha float32) ([]retrieval.SearchResult, error) {
	args := m.Called(ctx, query, vector, alpha)
	return args.Get(0).([]retrieval.SearchResult), args.Error(1)
}

func TestSearch_ReturnsMetadata(t *testing.T) {
	mockStore := new(MockStore)
	mockEmbedder := new(MockEmbedder) // Assume existing
	svc := retrieval.NewService(mockEmbedder, mockStore, nil)

	mockEmbedder.On("Embed", mock.Anything, "test").Return([]float32{0.1}, nil)
	expected := []retrieval.SearchResult{
		{Content: "test content", Score: 0.9, Metadata: map[string]interface{}{"source": "doc1"}},
	}
	mockStore.On("Search", mock.Anything, "test", []float32{0.1}, float32(0.5)).Return(expected, nil)

	results, err := svc.Search(context.Background(), "test")
	assert.NoError(t, err)
	assert.Equal(t, "doc1", results[0].Metadata["source"])
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/retrieval/... -v`
Expected: Fail due to undefined `SearchResult` and signature mismatch.

**Step 3: Write minimal implementation**
```go
// apps/backend/internal/retrieval/service.go
package retrieval

type SearchResult struct {
	Content  string                 `json:"content"`
	Score    float32                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}

type VectorStore interface {
	Search(ctx context.Context, query string, vector []float32, alpha float32) ([]SearchResult, error)
}

// Update Service.Search signature and implementation to pass through results
func (s *Service) Search(ctx context.Context, query string) ([]SearchResult, error) {
    // ... embedding ...
    docs, err := s.store.Search(ctx, query, vec, 0.5)
    // ...
    // Reranking logic needs update to handle SearchResult slice (Task 2)
    // For now, minimal pass-through
    return docs, err
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/retrieval/... -v`

***

### Task 2: Implement Configurable Reranking Logic

**Files:**
- Modify: `apps/backend/internal/retrieval/service.go`
- Test: `apps/backend/internal/retrieval/service_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `Service.Search` re-orders `[]SearchResult` based on `Reranker` output.
  2. Reranking preserves metadata.

- **Test Coverage**
  - [Unit] `Service.Search` with MockReranker verifies order change.

**Step 1: Write failing test**
```go
// apps/backend/internal/retrieval/service_test.go
// Add to existing test file
func TestSearch_WithReranker(t *testing.T) {
    // ... setup mocks ...
    initialResults := []retrieval.SearchResult{
        {Content: "A", Score: 0.5},
        {Content: "B", Score: 0.6},
    }
    mockStore.On("Search", ...).Return(initialResults, nil)
    
    // Reranker returns indices [1, 0] (swaps them)
    mockReranker.On("Rerank", ..., []string{"A", "B"}).Return([]int{1, 0}, nil)

    svc := retrieval.NewService(mockEmbedder, mockStore, mockReranker)
    results, _ := svc.Search(context.Background(), "test")
    
    assert.Equal(t, "B", results[0].Content)
    assert.Equal(t, "A", results[1].Content)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/retrieval/... -v`

**Step 3: Write minimal implementation**
```go
// apps/backend/internal/retrieval/service.go
// Update Search method
func (s *Service) Search(ctx context.Context, query string) ([]SearchResult, error) {
    // ... embed & store search ...
    
    if s.reranker != nil && len(docs) > 0 {
        // Extract content for reranker
        contents := make([]string, len(docs))
        for i, d := range docs {
            contents[i] = d.Content
        }

        indices, err := s.reranker.Rerank(ctx, query, contents)
        if err != nil {
             return nil, err 
        }

        reranked := make([]SearchResult, len(indices))
        for i, idx := range indices {
            if idx < len(docs) {
                reranked[i] = docs[idx]
                // Ideally, update score here if reranker provides it
            }
        }
        return reranked, nil
    }
    return docs, nil
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/retrieval/... -v`

***

### Task 3: Implement Query Logger (FR-5.3)

**Files:**
- Create: `apps/backend/internal/retrieval/logger.go`
- Modify: `apps/backend/internal/retrieval/service.go` (Integrate logger)
- Test: `apps/backend/internal/retrieval/logger_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Queries logged to stdout AND file `data/logs/query.log`.
  2. Log format: JSON with timestamp, query, latency, result count.

- **Test Coverage**
  - [Unit] `LogQuery` writes to provided io.Writer (buffer).

**Step 1: Write failing test**
```go
// apps/backend/internal/retrieval/logger_test.go
package retrieval

import (
    "bytes"
    "testing"
    "time"
    "encoding/json"
    "github.com/stretchr/testify/assert"
)

func TestQueryLogger(t *testing.T) {
    var buf bytes.Buffer
    logger := NewQueryLogger(&buf) // Inject writer
    
    entry := QueryLogEntry{
        Query: "test",
        Duration: 100 * time.Millisecond,
        NumResults: 5,
    }
    
    logger.Log(entry)
    
    var output map[string]interface{}
    err := json.Unmarshal(buf.Bytes(), &output)
    assert.NoError(t, err)
    assert.Equal(t, "test", output["query"])
    assert.Equal(t, 5.0, output["num_results"])
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/retrieval/... -v`

**Step 3: Write minimal implementation**
```go
// apps/backend/internal/retrieval/logger.go
package retrieval

import (
    "encoding/json"
    "io"
    "os"
    "time"
)

type QueryLogEntry struct {
    Timestamp   time.Time `json:"timestamp"`
    Query       string    `json:"query"`
    NumResults  int       `json:"num_results"`
    Duration    time.Duration `json:"duration_ns"` // or ms
    LatencyMs   int64     `json:"latency_ms"`
}

type QueryLogger struct {
    writer io.Writer
}

func NewQueryLogger(w io.Writer) *QueryLogger {
    return &QueryLogger{writer: w}
}

func NewFileQueryLogger(path string) (*QueryLogger, error) {
    f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    mw := io.MultiWriter(os.Stdout, f)
    return NewQueryLogger(mw), nil
}

func (l *QueryLogger) Log(entry QueryLogEntry) {
    entry.Timestamp = time.Now()
    entry.LatencyMs = entry.Duration.Milliseconds()
    json.NewEncoder(l.writer).Encode(entry)
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/retrieval/... -v`

***

### Task 4: Enhance MCP Handler (tools/list & Metadata)

**Files:**
- Modify: `apps/backend/features/mcp/handler.go`
- Test: `apps/backend/features/mcp/handler_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `tools/list` returns available tools (search).
  2. `tools/call` ("search") returns `SearchResult` content and populates metadata.

- **Test Coverage**
  - [Unit] `ServeHTTP` handles `tools/list`.
  - [Unit] `ServeHTTP` handles `tools/call` and includes metadata in response.

**Step 1: Write failing test**
```go
// apps/backend/features/mcp/handler_test.go
func TestMCP_ToolsList(t *testing.T) {
    // ... setup handler ...
    reqBody := `{"jsonrpc":"2.0", "method":"tools/list", "id":1}`
    req, _ := http.NewRequest("POST", "/mcp", strings.NewReader(reqBody))
    rr := httptest.NewRecorder()
    
    handler.ServeHTTP(rr, req)
    
    var resp JSONRPCResponse
    json.Unmarshal(rr.Body.Bytes(), &resp)
    
    result := resp.Result.(map[string]interface{})
    tools := result["tools"].([]interface{})
    assert.NotEmpty(t, tools)
    assert.Equal(t, "search", tools[0].(map[string]interface{})["name"])
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/mcp/... -v`

**Step 3: Write minimal implementation**
```go
// apps/backend/features/mcp/handler.go

// Add Tool definition structs
type Tool struct {
    Name        string      `json:"name"`
    Description string      `json:"description"`
    InputSchema interface{} `json:"inputSchema"`
}

type ListToolsResult struct {
    Tools []Tool `json:"tools"`
}

// In ServeHTTP
if req.Method == "tools/list" {
    response := JSONRPCResponse{
        JSONRPC: "2.0",
        ID:      req.ID,
        Result: ListToolsResult{
            Tools: []Tool{
                {
                    Name: "search",
                    Description: "Search documentation",
                    InputSchema: map[string]interface{}{
                        "type": "object",
                        "properties": map[string]interface{}{
                            "query": map[string]string{"type": "string"},
                        },
                        "required": []string{"query"},
                    },
                },
            },
        },
    }
    // write response
    return
}

// Update tools/call to use retrieval.SearchResult
// ... map SearchResult fields to ToolContent metadata ...
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/mcp/... -v`

***

### Task 5: Frontend Reranker Settings

**Files:**
- Modify: `apps/frontend/src/features/settings/Settings.vue`
- Modify: `apps/frontend/src/features/settings/settings.store.ts`
- Test: `apps/frontend/src/features/settings/Settings.spec.ts`

**Requirements:**
- **Acceptance Criteria**
  1. Dropdown for "Reranker Provider" (None, Jina AI, Cohere).
  2. Input field for "Reranker API Key" (visible only if provider != None).

- **Test Coverage**
  - [Unit] Store updates correctly.
  - [Unit] Component renders API key input conditionally.

**Step 1: Write failing test**
```typescript
// apps/frontend/src/features/settings/Settings.spec.ts
it('shows api key input when reranker is selected', async () => {
  const wrapper = mount(Settings, { ... });
  await wrapper.find('select[name="reranker"]').setValue('jinaai');
  expect(wrapper.find('input[name="rerankerApiKey"]').exists()).toBe(true);
});
```

**Step 2: Verify test fails**
Run: `npm run test:unit apps/frontend/src/features/settings/Settings.spec.ts`

**Step 3: Write minimal implementation**
```vue
<!-- apps/frontend/src/features/settings/Settings.vue -->
<template>
  <!-- ... -->
  <Select v-model="settings.rerankProvider">
    <SelectTrigger><SelectValue placeholder="Select provider" /></SelectTrigger>
    <SelectContent>
      <SelectItem value="none">None</SelectItem>
      <SelectItem value="jinaai">Jina AI</SelectItem>
      <SelectItem value="cohere">Cohere</SelectItem>
    </SelectContent>
  </Select>
  
  <div v-if="settings.rerankProvider !== 'none'">
     <Input v-model="settings.rerankApiKey" placeholder="API Key" />
  </div>
</template>
```

**Step 4: Verify test passes**
Run: `npm run test:unit apps/frontend/src/features/settings/Settings.spec.ts`
