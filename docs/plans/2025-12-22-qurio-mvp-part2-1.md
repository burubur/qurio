# Implementation Plan - Qurio MVP Part 2

**Date:** 2025-12-22
**Source:** `docs/2025-12-21-qurio-mvp.md`
**Goal:** Implement Source Management (Dedupe, Re-sync, Delete) and Dynamic Configuration (Reranker Settings).

## ✓ Requirements Extracted

**Scope:**
-   **Backend:** SHA-256 Deduplication (FR-2.3), Re-sync Logic (FR-4.1), Source Deletion (FR-4.2), Dynamic Settings API (Story 6).
-   **Frontend:** Source Actions (Re-sync, Delete), Settings UI.

**Gap Analysis:**
-   **Nouns:** Settings (Table), Source Actions (Delete, Re-sync), Content Hash.
-   **Verbs:** Calculate Hash, Soft Delete, Update Settings, Trigger Re-sync.

## ✓ Knowledge Enrichment

**RAG Queries Executed:**
-   "PostgreSQL singleton settings table pattern" (Decision: Single row with columns for type safety).
-   "NSQ publish message for re-sync" (Standard `nsq.Producer`).

---

### Task 1: Backend - Dynamic Settings (Store & API)

**Files:**
-   Create: `apps/backend/migrations/000002_create_settings.up.sql`
-   Create: `apps/backend/internal/settings/service.go`
-   Create: `apps/backend/internal/settings/handler.go`
-   Modify: `apps/backend/main.go` (Register routes)
-   Test: `apps/backend/internal/settings/service_test.go`

**Requirements:**
-   **Functional:**
    -   Store `rerank_provider` (jina/cohere/none) and `rerank_api_key`.
    -   Ensure only ONE row exists (singleton).
    -   `GET /api/settings` returns current config.
    -   `PUT /api/settings` updates config.
-   **Test Coverage:**
    -   [Unit] `GetSettings` returns default if empty.
    -   [Integration] `UpdateSettings` persists changes.

**Step 1: Write failing test**
```go
// apps/backend/internal/settings/service_test.go
package settings_test

import (
	"context"
	"testing"
	"qurio/apps/backend/internal/settings"
)

func TestGetSettings_Default(t *testing.T) {
	repo := newMockRepo() // Empty repo
	svc := settings.NewService(repo)
	
	s, err := svc.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if s.RerankProvider != "none" {
		t.Errorf("Expected default 'none', got %s", s.RerankProvider)
	}
}
```

**Step 2: Verify test fails**
`go test ./apps/backend/internal/settings/...` -> FAIL

**Step 3: Write minimal implementation**
```go
// apps/backend/internal/settings/service.go
package settings

type Settings struct {
	RerankProvider string `json:"rerank_provider"` // none, jina, cohere
	RerankAPIKey   string `json:"rerank_api_key"`
}

func (s *Service) Get(ctx context.Context) (*Settings, error) {
	// Repo fetch, if error or empty return default
	set, err := s.repo.GetLast(ctx)
	if err != nil {
		return &Settings{RerankProvider: "none"}, nil
	}
	return set, nil
}
```

**Step 4: Verify test passes**
`go test ./apps/backend/internal/settings/...` -> PASS

---

### Task 2: Backend - Source Management (Dedupe & Actions)

**Files:**
-   Modify: `apps/backend/features/source/service.go` (Add Dedupe check, Soft Delete)
-   Modify: `apps/backend/features/source/handler.go` (Add Re-sync, Delete endpoints)
-   Modify: `apps/backend/internal/worker/ingest.go` (Calculate/Store Hash)
-   Modify: `apps/backend/migrations/000003_add_source_hash_deleted.up.sql`
-   Test: `apps/backend/features/source/service_test.go`

**Requirements:**
-   **Functional:**
    -   `Create`: Check SHA-256 hash. If exists, return duplicate error.
    -   `Delete`: Set `deleted_at` timestamp (Soft delete).
    -   `ReSync`: Publish NSQ message with `resync=true`.
-   **Test Coverage:**
    -   [Unit] `Create` fails if hash exists.
    -   [Unit] `Delete` updates timestamp.

**Step 1: Write failing test**
```go
// apps/backend/features/source/service_test.go
func TestCreate_Duplicate(t *testing.T) {
	repo := newMockRepo()
	repo.ExistsHash = true // Simulate existing hash
	svc := source.NewService(repo, nil)
	
	err := svc.Create(ctx, &source.Source{Hash: "abc"})
	if err == nil || err.Error() != "Duplicate detected" {
		t.Fatal("Expected duplicate error")
	}
}
```

**Step 2: Verify test fails**
`go test ./apps/backend/features/source/...` -> FAIL

**Step 3: Write minimal implementation**
```go
// apps/backend/features/source/service.go
func (s *Service) Create(ctx context.Context, src *Source) error {
	exists, _ := s.repo.CheckHash(ctx, src.Hash)
	if exists {
		return fmt.Errorf("Duplicate detected")
	}
	return s.repo.Save(ctx, src)
}
```

**Step 4: Verify test passes**
`go test ./apps/backend/features/source/...` -> PASS

---

### Task 3: Frontend - Settings Page (Story 6)

**Files:**
-   Create: `apps/frontend/src/features/settings/Settings.vue`
-   Create: `apps/frontend/src/features/settings/settings.store.ts`
-   Modify: `apps/frontend/src/App.vue` (Add nav link)
-   Test: `apps/frontend/src/features/settings/Settings.spec.ts`

**Requirements:**
-   **Functional:**
    -   Form to select Reranker (Jina/Cohere/None).
    -   Input for API Key.
    -   Save button calls `PUT /api/settings`.
-   **Test Coverage:**
    -   [Unit] Loading page fetches settings.
    -   [Unit] Save calls API.

**Step 1: Write failing test**
```typescript
// apps/frontend/src/features/settings/Settings.spec.ts
import { mount } from '@vue/test-utils'
import Settings from './Settings.vue'
import { createTestingPinia } from '@pinia/testing'

test('loads settings on mount', () => {
  const wrapper = mount(Settings, {
    global: { plugins: [createTestingPinia()] }
  })
  const store = useSettingsStore()
  expect(store.fetchSettings).toHaveBeenCalled()
})
```

**Step 2: Verify test fails**
`npm run test` -> FAIL

**Step 3: Write minimal implementation**
```typescript
// apps/frontend/src/features/settings/Settings.vue
<script setup>
import { onMounted } from 'vue'
import { useSettingsStore } from './settings.store'
const store = useSettingsStore()
onMounted(() => store.fetchSettings())
</script>
```

**Step 4: Verify test passes**
`npm run test` -> PASS

---

### Task 4: Frontend - Source Actions (Re-sync/Delete)

**Files:**
-   Modify: `apps/frontend/src/features/sources/SourceList.vue`
-   Modify: `apps/frontend/src/features/sources/source.store.ts`
-   Test: `apps/frontend/src/features/sources/SourceList.spec.ts`

**Requirements:**
-   **Functional:**
    -   Add "Re-sync" and "Delete" buttons to each row.
    -   Delete asks for confirmation (browser confirm ok for MVP).
-   **Test Coverage:**
    -   [Unit] Click Delete -> calls store.deleteSource.

**Step 1: Write failing test**
```typescript
// apps/frontend/src/features/sources/SourceList.spec.ts
test('calls delete when button clicked', async () => {
  const wrapper = mount(SourceList, ...)
  await wrapper.find('.delete-btn').trigger('click')
  expect(store.deleteSource).toHaveBeenCalled()
})
```

**Step 2: Verify test fails**
`npm run test` -> FAIL

**Step 3: Write minimal implementation**
```typescript
// SourceList.vue
<button class="delete-btn" @click="store.deleteSource(source.id)">Delete</button>
```

**Step 4: Verify test passes**
`npm run test` -> PASS
