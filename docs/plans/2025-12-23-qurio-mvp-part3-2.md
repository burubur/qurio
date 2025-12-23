# Implementation Plan - MVP Part 3.2: Ingestion & Crawler Integration

**Scope:** Integrate the recursive web crawler into the ingestion worker, enabling depth control and exclusion rules. Update Frontend to support these configurations.

**Gap Analysis:**
- **Worker:** Currently fetches single URL; needs to use `crawler` package.
- **Source Config:** DB and API missing `max_depth` and `exclusions`.
- **Frontend:** Missing inputs for advanced crawl settings.

**Exclusions:**
- **File Upload:** Deferred to Part 3.3 to keep this plan atomic to Web Crawling.

***

### Task 1: Database Migration for Source Config

**Files:**
- Create: `apps/backend/migrations/000005_add_source_config.up.sql`
- Modify: `apps/backend/features/source/source.go` (Struct update)
- Modify: `apps/backend/features/source/repo.go` (Scan/Save update)
- Test: `apps/backend/features/source/repo_test.go` (Verify new fields)

**Requirements:**
- **Acceptance Criteria**
  1. `sources` table has `max_depth` (int, default 0) and `exclusions` (text/json, default empty).
  2. `Source` struct includes `MaxDepth` and `Exclusions`.
  3. Repository saves and retrieves these fields correctly.

- **Test Coverage**
  - [Integration] `Repo.Save` preserves depth/exclusions.
  - [Integration] `Repo.Get` returns depth/exclusions.

**Step 1: Write failing test**
```go
// apps/backend/features/source/repo_test.go
func TestSaveAndGet_WithConfig(t *testing.T) {
    repo := setupTestRepo(t)
    src := &source.Source{
        URL: "http://example.com",
        MaxDepth: 2,
        Exclusions: []string{"/blog", "/login"},
    }
    
    err := repo.Save(context.Background(), src)
    assert.NoError(t, err)
    
    saved, err := repo.Get(context.Background(), src.ID)
    assert.NoError(t, err)
    assert.Equal(t, 2, saved.MaxDepth)
    assert.Contains(t, saved.Exclusions, "/blog")
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/source/... -v`
Expected: Fail due to missing struct fields and DB columns.

**Step 3: Write minimal implementation**
```sql
-- apps/backend/migrations/000005_add_source_config.up.sql
ALTER TABLE sources ADD COLUMN max_depth INTEGER DEFAULT 0;
ALTER TABLE sources ADD COLUMN exclusions TEXT DEFAULT ''; -- CSV or JSON
```

```go
// apps/backend/features/source/source.go
type Source struct {
    // ... existing ...
    MaxDepth   int      `json:"max_depth"`
    Exclusions []string `json:"exclusions"`
}

// apps/backend/features/source/repo.go
// Update Save:
// INSERT INTO sources ... (..., max_depth, exclusions) VALUES (..., $3, $4)
// exclusions stored as JSON string or comma-separated
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/source/... -v`

***

### Task 2: Backend Source Logic Update

**Files:**
- Modify: `apps/backend/features/source/source.go` (Service.Create)
- Test: `apps/backend/features/source/source_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `Service.Create` includes config in the NSQ payload.
  2. Payload format: `{"url": "...", "id": "...", "max_depth": 2, "exclusions": [...]}`.

- **Test Coverage**
  - [Unit] `Create` publishes correct JSON payload to `ingest` topic.

**Step 1: Write failing test**
```go
// apps/backend/features/source/source_test.go
func TestCreate_PublishesConfig(t *testing.T) {
    mockPub := new(MockPublisher)
    svc := source.NewService(mockRepo, mockPub)
    
    src := &source.Source{URL: "http://test.com", MaxDepth: 3}
    
    mockPub.On("Publish", "ingest", mock.MatchedBy(func(body []byte) bool {
        var p map[string]interface{}
        json.Unmarshal(body, &p)
        return p["max_depth"] == float64(3)
    })).Return(nil)
    
    svc.Create(context.Background(), src)
    mockPub.AssertExpectations(t)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/source/... -v`

**Step 3: Write minimal implementation**
```go
// apps/backend/features/source/source.go
func (s *Service) Create(ctx context.Context, src *Source) error {
    // ... save ...
    payload, _ := json.Marshal(map[string]interface{}{
        "url":        src.URL,
        "id":         src.ID,
        "max_depth":  src.MaxDepth,
        "exclusions": src.Exclusions,
    })
    return s.pub.Publish("ingest", payload)
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/source/... -v`

***

### Task 3: Worker Integration with Crawler

**Files:**
- Modify: `apps/backend/internal/worker/ingest.go`
- Test: `apps/backend/internal/worker/ingest_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Worker unmarshals `max_depth` and `exclusions`.
  2. Worker initializes `crawler.New(config)`.
  3. Worker iterates over `crawler.Crawl()` results (multiple pages) instead of single `fetcher.Fetch`.

- **Test Coverage**
  - [Unit] `HandleMessage` invokes crawler and processes multiple pages.

**Step 1: Write failing test**
```go
// apps/backend/internal/worker/ingest_test.go
func TestHandleMessage_RecursiveCrawl(t *testing.T) {
    // Setup MockCrawler that returns 2 pages
    // Verify store.StoreChunk is called for BOTH pages
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/worker/... -v`

**Step 3: Write minimal implementation**
```go
// apps/backend/internal/worker/ingest.go
// Remove Fetcher interface, use crawler directly or wrapped
func (h *IngestHandler) HandleMessage(m *nsq.Message) error {
    // ... unmarshal payload ...
    cfg := crawler.Config{
        MaxDepth:   int(payload["max_depth"].(float64)),
        Exclusions: toSlice(payload["exclusions"]),
    }
    c, _ := crawler.New(cfg)
    pages, _ := c.Crawl(payload.URL)
    
    for _, page := range pages {
        // Chunk, Embed, Store loop for EACH page
    }
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/worker/... -v`

***

### Task 4: Frontend Source Form Update

**Files:**
- Modify: `apps/frontend/src/features/sources/SourceForm.vue`
- Modify: `apps/frontend/src/features/sources/source.store.ts`
- Test: `apps/frontend/src/features/sources/SourceForm.spec.ts`

**Requirements:**
- **Acceptance Criteria**
  1. Form includes Number input for "Crawl Depth" (0-5).
  2. Form includes Textarea for "Exclusions" (one per line).
  3. Submit payload includes `max_depth` and `exclusions`.

- **Test Coverage**
  - [Unit] Form emits submit event with new fields.

**Step 1: Write failing test**
```typescript
// apps/frontend/src/features/sources/SourceForm.spec.ts
it('submits depth and exclusions', async () => {
    // fill inputs
    // assert emitted payload
})
```

**Step 2: Verify test fails**
Run: `npm run test:unit apps/frontend/src/features/sources/SourceForm.spec.ts`

**Step 3: Write minimal implementation**
```vue
<!-- SourceForm.vue -->
<Input v-model="form.maxDepth" type="number" label="Depth" />
<Textarea v-model="form.exclusions" label="Exclusions (regex)" />
```

**Step 4: Verify test passes**
Run: `npm run test:unit apps/frontend/src/features/sources/SourceForm.spec.ts`
