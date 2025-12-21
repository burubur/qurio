# Implementation Plan - Qurio MVP

**Scope:** Complete MVP implementation including Backend (Go), Frontend (Vue 3), Docling Service (Python), and Infrastructure (Docker Compose).
**Status:** Approved
**Date:** 2025-12-21

---

## 1. Requirements Analysis

### 1.1 Scope
- **In Scope:** Monorepo setup, Docker Compose (Postgres, Weaviate, NSQ), Go Backend (API+Worker), Python Docling Service, Vue 3 Frontend.
- **Out of Scope:** Multi-tenancy, Auth, Cloud deployment.

### 1.2 Gap Analysis
- **Nouns:** Source, Document, Chunk, IngestionJob, MCPRequest, MCPResponse.
- **Verbs:** Create Source, Crawl URL, Upload File, Deduplicate, Embed (Gemini), Index (Weaviate), Search (Hybrid), Rerank, Log Query.

---

## 2. Implementation Tasks

### Task 1: Project Scaffold & Infrastructure
**Files:**
- Create: `docker-compose.yml`
- Create: `.env.example`
- Create: `go.work` (Go workspace)
- Create: `apps/backend/go.mod`
- Create: `apps/backend/main.go`
- Create: `services/docling/Dockerfile`
- Create: `services/docling/requirements.txt`

**Requirements:**
- **AC:** `docker-compose up` starts Postgres, Weaviate, NSQ, and placeholder services.
- **FR:** FR-1.1 Single Command Deployment.
- **NFR:** Localhost access only.

**Step 1: Write failing test (Infrastructure verification script)**
```bash
# verify_infra.sh
#!/bin/bash
# Check if services are reachable
curl -f http://localhost:8080/health || exit 1
curl -f http://localhost:3000 || exit 1
nc -z localhost 5432 || exit 1 # Postgres
nc -z localhost 8080 || exit 1 # Weaviate
```

**Step 2: Verify test fails**
Run: `./verify_infra.sh`
Expected: FAIL (Connection refused)

**Step 3: Implementation**
```yaml
# docker-compose.yml
version: '3.8'
services:
  weaviate:
    image: semitechnologies/weaviate:1.24.10
    ports: ["8080:8080"]
    environment:
      QUERY_DEFAULTS_LIMIT: 25
      AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED: 'true'
      PERSISTENCE_DATA_PATH: '/var/lib/weaviate'
      DEFAULT_VECTORIZER_MODULE: 'none'
  postgres:
    image: postgres:15-alpine
    ports: ["5432:5432"]
    environment:
      POSTGRES_USER: qurio
      POSTGRES_PASSWORD: password
      POSTGRES_DB: qurio
  nsqlookupd:
    image: nsqio/nsq
    command: /nsqlookupd
  nsqd:
    image: nsqio/nsq
    command: /nsqd --lookupd-tcp-address=nsqlookupd:4160
    depends_on: [nsqlookupd]
  docling:
    build: ./services/docling
    ports: ["8000:8000"]
  backend:
    build: ./apps/backend
    ports: ["8081:8081"] # API Port
    depends_on: [postgres, weaviate, nsqd]
  frontend:
    build: ./apps/frontend
    ports: ["3000:3000"]
```

**Step 4: Verify test passes**
Run: `docker-compose up -d && ./verify_infra.sh`
Expected: PASS

---

### Task 2: Backend Core - Configuration & Database
**Files:**
- Create: `apps/backend/internal/config/config.go`
- Create: `apps/backend/internal/config/config_test.go`
- Create: `apps/backend/migrations/000001_init_schema.up.sql`
- Modify: `apps/backend/main.go` (Add config load & migration)

**Requirements:**
- **AC:** Backend loads config from env, runs Postgres migrations on startup.
- **FR:** FR-1.1 (Config/Migrations).

**Step 1: Write failing test**
```go
// apps/backend/internal/config/config_test.go
package config_test

import (
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
    "qurio/apps/backend/internal/config"
)

func TestLoadConfig(t *testing.T) {
    os.Setenv("DB_HOST", "localhost")
    cfg, err := config.Load()
    assert.NoError(t, err)
    assert.Equal(t, "localhost", cfg.DBHost)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/config/...`
Expected: FAIL (undefined config package)

**Step 3: Implementation**
```go
// apps/backend/internal/config/config.go
package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
    DBHost string `envconfig:"DB_HOST" default:"postgres"`
    // ... other fields
}

func Load() (*Config, error) {
    var cfg Config
    err := envconfig.Process("", &cfg)
    return &cfg, err
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/config/...`
Expected: PASS

---

### Task 3: Weaviate Schema Initialization
**Files:**
- Create: `apps/backend/internal/vector/schema.go`
- Create: `apps/backend/internal/vector/schema_test.go`

**Requirements:**
- **AC:** Ensure `DocumentChunk` class exists in Weaviate on startup.
- **FR:** FR-2.5 (Contextual Embeddings - Schema).

**Step 1: Write failing test**
```go
// apps/backend/internal/vector/schema_test.go
package vector_test

import (
    "testing"
    "qurio/apps/backend/internal/vector"
    "github.com/stretchr/testify/mock"
)

type MockWeaviateClient struct {
    mock.Mock
}
// ... mock implementation

func TestEnsureSchema(t *testing.T) {
    mockClient := new(MockWeaviateClient)
    mockClient.On("ClassExists", "DocumentChunk").Return(false, nil)
    mockClient.On("CreateClass", mock.Anything).Return(nil)
    
    err := vector.EnsureSchema(mockClient)
    assert.NoError(t, err)
    mockClient.AssertExpectations(t)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/vector/...`
Expected: FAIL (undefined)

**Step 3: Implementation**
```go
// apps/backend/internal/vector/schema.go
package vector

import (
    "github.com/weaviate/weaviate-go-client/v5/weaviate"
    "github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
    "context"
)

func EnsureSchema(client *weaviate.Client) error {
    // Check if class exists
    exists, _ := client.Schema().ClassExistenceChecker().WithClassName("DocumentChunk").Do(context.Background())
    if !exists {
        // Create class definition
        // ... implementation using weaviate models
    }
    return nil
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/vector/...`
Expected: PASS

---

### Task 4: Docling Service (Python)
**Files:**
- Create: `services/docling/main.py`
- Create: `services/docling/test_main.py`

**Requirements:**
- **AC:** POST `/process` accepts file, returns Markdown.
- **FR:** FR-2.2 OCR Processing.

**Step 1: Write failing test**
```python
# services/docling/test_main.py
from fastapi.testclient import TestClient
from main import app

client = TestClient(app)

def test_process_document():
    files = {'file': ('test.txt', b'Hello World', 'text/plain')}
    response = client.post("/process", files=files)
    assert response.status_code == 200
    assert "text" in response.json()
```

**Step 2: Verify test fails**
Run: `pytest services/docling`
Expected: FAIL (Module not found)

**Step 3: Implementation**
```python
# services/docling/main.py
from fastapi import FastAPI, UploadFile
from docling.document_converter import DocumentConverter

app = FastAPI()
converter = DocumentConverter()

@app.post("/process")
async def process(file: UploadFile):
    # Save temp file, run docling, return text
    return {"text": "processed text"} 
```

**Step 4: Verify test passes**
Run: `pytest services/docling`
Expected: PASS

---

### Task 5: Source Management (Go Feature)
**Files:**
- Create: `apps/backend/features/source/source.go` (Domain)
- Create: `apps/backend/features/source/handler.go` (HTTP)
- Create: `apps/backend/features/source/repo.go` (DB)
- Create: `apps/backend/features/source/source_test.go`

**Requirements:**
- **AC:** CRUD API for Sources.
- **FR:** FR-4.2 Source Management.

**Step 1: Write failing test**
```go
// apps/backend/features/source/source_test.go
func TestCreateSource(t *testing.T) {
    repo := new(MockRepo)
    svc := source.NewService(repo)
    
    src := source.Source{URL: "https://example.com"}
    repo.On("Save", src).Return(nil)
    
    err := svc.Create(context.Background(), src)
    assert.NoError(t, err)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/source/...`
Expected: FAIL

**Step 3: Implementation**
```go
// apps/backend/features/source/source.go
type Source struct {
    ID string
    URL string
}
type Service interface {
    Create(ctx context.Context, src Source) error
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/source/...`
Expected: PASS

---

### Task 6: Ingestion Worker (NSQ Consumer)
**Files:**
- Create: `apps/backend/cmd/worker/main.go`
- Create: `apps/backend/internal/worker/ingest.go`
- Create: `apps/backend/internal/worker/ingest_test.go`

**Requirements:**
- **AC:** Consume message, fetch content, chunk, embed, index.
- **FR:** FR-2.1, FR-2.4, FR-2.5.

**Step 1: Write failing test**
```go
func TestHandleMessage(t *testing.T) {
    // Mock Embedding API, Weaviate
    worker := worker.NewIngestHandler(mockEmbedder, mockVectorDB)
    err := worker.HandleMessage(&nsq.Message{Body: []byte(`{"url":"https://example.com"}`)})
    assert.NoError(t, err)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/worker/...`
Expected: FAIL

**Step 3: Implementation**
```go
// apps/backend/internal/worker/ingest.go
func (h *IngestHandler) HandleMessage(m *nsq.Message) error {
    // 1. Fetch/Crawl
    // 2. Chunk
    // 3. Embed (Gemini)
    // 4. Index (Weaviate)
    return nil
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/worker/...`
Expected: PASS

---

### Task 7: Retrieval & MCP Endpoint
**Files:**
- Create: `apps/backend/features/mcp/handler.go`
- Create: `apps/backend/features/mcp/handler_test.go`
- Create: `apps/backend/internal/retrieval/service.go`

**Requirements:**
- **AC:** POST `/mcp` accepts JSON-RPC, returns Weaviate results.
- **FR:** FR-5.1, FR-5.2.

**Step 1: Write failing test**
```go
func TestMCPCall(t *testing.T) {
    req := []byte(`{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search","arguments":{"query":"test"}}}`)
    w := httptest.NewRecorder()
    r := httptest.NewRequest("POST", "/mcp", bytes.NewBuffer(req))
    
    mcpHandler.ServeHTTP(w, r)
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "result")
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/mcp/...`
Expected: FAIL

**Step 3: Implementation**
```go
// apps/backend/features/mcp/handler.go
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Parse JSON-RPC
    // Call retrieval.Search(query)
    // Return JSON-RPC response
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/mcp/...`
Expected: PASS

---

### Task 8: Frontend Implementation (Vue)
**Files:**
- Create: `apps/frontend/src/App.vue`
- Create: `apps/frontend/src/features/sources/SourceList.vue`
- Create: `apps/frontend/src/features/sources/SourceForm.vue`

**Requirements:**
- **AC:** View sources, add source.
- **FR:** FR-4.2 (UI).

**Step 1: Write failing test**
```typescript
// apps/frontend/src/features/sources/SourceList.spec.ts
import { mount } from '@vue/test-utils'
import SourceList from './SourceList.vue'

test('displays sources', () => {
  const wrapper = mount(SourceList, { props: { sources: [{ id: 1, url: 'test' }] } })
  expect(wrapper.text()).toContain('test')
})
```

**Step 2: Verify test fails**
Run: `npm test`
Expected: FAIL

**Step 3: Implementation**
```vue
<!-- apps/frontend/src/features/sources/SourceList.vue -->
<template>
  <div v-for="source in sources" :key="source.id">{{ source.url }}</div>
</template>
<script setup lang="ts">
defineProps<{ sources: any[] }>()
</script>
```

**Step 4: Verify test passes**
Run: `npm test`
Expected: PASS
