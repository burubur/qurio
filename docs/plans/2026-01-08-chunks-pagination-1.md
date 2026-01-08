# Plan: Chunk Pagination and Count Fix

**Date:** 2026-01-08
**Scope:** Backend pagination for chunks, total count implementation, and frontend "Load More" functionality.
**Reference:** Aligns with Technical Constitution (Testability-First, I/O Isolation).

### Task 1: Backend Store - Implement Count and Pagination

**Files:**
- Modify: `apps/backend/features/source/source.go`
- Modify: `apps/backend/internal/adapter/weaviate/store.go`
- Modify: `apps/backend/features/source/mock_store.go` (if exists, or any mocks)

**Requirements:**
- **Acceptance Criteria**
  1. `CountChunksBySource` returns accurate count of chunks for a given `sourceID`.
  2. `GetChunks` accepts `limit` and `offset` and returns filtered results.

- **Functional Requirements**
  1. Update `ChunkStore` interface to include `CountChunksBySource` and updated `GetChunks`.
  2. Implement Weaviate `Aggregate` with `where` filter for counting.
  3. Implement `WithLimit` and `WithOffset` in `GetChunks`.

- **Test Coverage**
  - [Unit] `weaviate/store_test.go` - Test `CountChunksBySource` and `GetChunks` with pagination (mocking client or using integration test).
  - Test data fixtures: Mock Weaviate response.

**Step 1: Write failing test**
Create `apps/backend/internal/adapter/weaviate/store_pagination_test.go` checking interface compliance or functionality.

**Step 2: Verify test fails**
Run `go test ./apps/backend/internal/adapter/weaviate/...`

**Step 3: Write minimal implementation**
- Update interface in `source.go`.
- Implement methods in `store.go`.

**Step 4: Verify test passes**
Run `go test ./apps/backend/internal/adapter/weaviate/...`

### Task 2: Backend Service - Integrate Pagination & Optimization

**Files:**
- Modify: `apps/backend/features/source/source.go`
- Test: `apps/backend/features/source/service_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `Service.Get` returns correct `TotalChunks` regardless of returned `Chunks` length.
  2. `Service.Get` respects `limit` and `offset` arguments.
  3. `Service.Get` accepts `includeChunks` bool; if false, returns empty `Chunks` slice (optimization for polling).

- **Functional Requirements**
  1. Update `Get` signature to accept pagination params and `includeChunks` flag.
  2. Call `CountChunksBySource`.
  3. If `includeChunks` is true, call `store.GetChunks` with limit/offset.
  4. If `includeChunks` is false, skip `store.GetChunks` (return empty).

- **Test Coverage**
  - [Unit] `Service.Get` calls store with correct params.
  - [Unit] `Service.Get` skips store call when `includeChunks=false`.

**Step 1: Write failing test**
Update `service_test.go` to test pagination and the exclusion flag.

**Step 2: Verify test fails**
Run `go test ./apps/backend/features/source/...`

**Step 3: Write minimal implementation**
Update `Get` method in `source.go`.

**Step 4: Verify test passes**
Run `go test ./apps/backend/features/source/...`

### Task 3: Backend Handler - Parse Query Params

**Files:**
- Modify: `apps/backend/features/source/handler.go`
- Test: `apps/backend/features/source/handler_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `GET /sources/{id}?limit=10&offset=5` passes values to service.
  2. `GET /sources/{id}?exclude_chunks=true` passes `includeChunks=false` to service.
  3. Defaults to `limit=100`, `offset=0`, `includeChunks=true`.

- **Functional Requirements**
  1. Parse query string for `limit`, `offset`, `exclude_chunks`.
  2. Pass to `service.Get`.

- **Test Coverage**
  - [Unit] `Handler.Get` extracts query params correctly.

**Step 1: Write failing test**
Update `handler_test.go`.

**Step 2: Verify test fails**
Run `go test ./apps/backend/features/source/...`

**Step 3: Write minimal implementation**
Update `handler.go`.

**Step 4: Verify test passes**
Run `go test ./apps/backend/features/source/...`

### Task 4: Frontend Store - Support Pagination & Metadata Polling

**Files:**
- Modify: `apps/frontend/src/features/sources/source.store.ts`

**Requirements:**
- **Acceptance Criteria**
  1. `getSource(id, params)` accepts `limit`, `offset`, `exclude_chunks`.
  2. `fetchChunks(id, offset)` specifically fetches chunks and appends to state.
  3. `pollSourceStatus(id)` fetches *only* metadata (exclude_chunks=true) and updates status/total_chunks without touching the chunks list.

- **Functional Requirements**
  1. Refactor `getSource` to handle params.
  2. Add `fetchChunks` action.
  3. Update `chunks` state handling (append vs replace).

**Step 1: Write failing test**
Create/Update `source.store.spec.ts`.

**Step 2: Verify test fails**
Run `npm run test:unit`

**Step 3: Write minimal implementation**
Update `source.store.ts`.

**Step 4: Verify test passes**
Run `npm run test:unit`

**Step 4: Verify test passes**
Run `npm run test:unit`

### Reliability & Performance Strategy

**1. Data Integrity**
- **Prevention of Silent Truncation:** All chunks are now accessible via pagination, removing the hardcoded 100-item limit.
- **Accurate Counts:** The `total_chunks` field is now derived from a database aggregate query, not the length of the returned slice, ensuring 100% accuracy even when `exclude_chunks=true`.

**2. UI Stability**
- **State Preservation:** By separating status polling (`exclude_chunks=true`) from data fetching, the user's scroll position and loaded chunks are preserved during background updates.
- **Deterministic Polling:** Polling now strictly updates metadata (status, progress), eliminating race conditions where a poll might overwrite a "Load More" action.

**3. Performance & Safety**
- **Reduced Network Load:** Polling requests now return <1KB JSON payloads (metadata only) instead of potential MBs of chunk text, significantly reducing bandwidth and client CPU usage.
- **Safe Defaults:** The backend enforces a default `limit=100` if unspecified, protecting the system from OOM errors on large datasets.
- **Offset Trade-off:** We utilize Offset Pagination (`limit/offset`). While simple and effective for this use case, note that if new chunks are inserted *exactly* while a user is paginating, minor duplicates or skips could theoretically occur. This is an acceptable trade-off for an admin-facing view compared to the complexity of cursor-based pagination for this specific data model.