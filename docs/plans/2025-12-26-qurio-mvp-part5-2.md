# Implementation Plan - Bug Fixes & Inconsistencies

**Date:** 2025-12-26
**Feature:** MVP Part 5.2 (Bug Fixes & Technical Debt)
**Status:** Planned

## 1. Requirements Analysis

### Scope
Address 5 critical inconsistencies identified in `docs/2025-12-26-bugs-inconsistencies.md` to bring the codebase into compliance with the Technical Constitution. Focus on API response standardization, error handling/tracing, structured logging, data integrity (schema), and resource management.

### Gap Analysis
- **API Response Envelope:** `Job` and `Stats` features return raw JSON. -> **Fix:** Wrap in `{"data": ...}`.
- **Error Handling:** `Job` and `Stats` use `http.Error` (text). -> **Fix:** Use JSON envelope with correlation ID.
- **Logging:** `Job` feature has no logs. `Ingestion Worker` mixes logging/structlog. -> **Fix:** Add `slog` and clean up Python logging.
- **Data Integrity:** Weaviate `text` schema prevents exact match deletion. -> **Fix:** Change to `string`.
- **Resource Management:** `Job` service retry blocks indefinitely. -> **Fix:** Add timeout to NSQ publish.

### Exclusions
- Refactoring `Source` or `Settings` features (they are already correct).
- Changing the actual business logic of jobs/stats (only the interface/plumbing).

---

## 2. Knowledge Enrichment

### Reference Patterns (from Codebase Investigation)
- **API Envelope:** `apps/backend/features/source/handler.go` uses `{"data": response, "meta": meta}`.
- **Error Helper:** Local `writeError` method in handlers:
  ```go
  func (h *Handler) writeError(w http.ResponseWriter, err error, code string, status int, traceID string) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(status)
      json.NewEncoder(w).Encode(map[string]interface{}{
          "error": map[string]string{
              "code":    code,
              "message": err.Error(),
          },
          "correlationId": traceID,
      })
  }
  ```
- **Correlation ID:** `middleware.GetCorrelationID(ctx)`.
- **Weaviate Schema:** `text` = tokenized, `string` = exact match (needed for ID/URL filtering).

---

## 3. Implementation Tasks

### Task 1: Fix Job Feature API & Logging

**Files:**
- Modify: `apps/backend/features/job/handler.go`
- Modify: `apps/backend/features/job/service.go`

**Requirements:**
- **Acceptance Criteria**
  1. `GET /jobs` returns `{"data": [...], "meta": {"count": N}}`.
  2. Errors return JSON with `code`, `message`, `correlationId`.
  3. All public methods log start/finish with `slog` and correlation ID.
  4. `Retry` method in service uses a 5-second timeout for NSQ publish.

- **Test Coverage**
  - [Integration] `GET /jobs` verifies JSON structure.
  - [Unit] `Retry` service method mocks `Publish` delay to verify timeout error.

**Step 1: Write failing test (Service Timeout)**
Create `apps/backend/features/job/service_test.go`:
```go
func TestRetry_Timeout(t *testing.T) {
    // Mock publisher that hangs
    // Call Retry
    // Expect error "context deadline exceeded" or similar
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/job/... -v`

**Step 3: Implement Fixes**
- **handler.go:** Add `writeError`. Wrap `List` response. Add `slog`.
- **service.go:** Wrap `s.pub.Publish` in a goroutine + select with `time.After(5 * time.Second)`.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/job/... -v`

---

### Task 2: Fix Stats Feature API & Logging

**Files:**
- Modify: `apps/backend/features/stats/handler.go`

**Requirements:**
- **Acceptance Criteria**
  1. `GET /stats` returns `{"data": {...}}`.
  2. Errors return JSON with `code`, `message`, `correlationId`.
  3. `GetStats` logs start/finish with `slog` and correlation ID.

- **Test Coverage**
  - [Integration] `GET /stats` verifies JSON structure.

**Step 1: Write failing test**
Create/Update `apps/backend/features/stats/handler_test.go` to assert `data` envelope.

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/stats/... -v`

**Step 3: Implement Fixes**
- **handler.go:** Add `writeError`. Wrap `GetStats` response. Add `slog`.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/stats/... -v`

---

### Task 3: Fix Worker Trace Propagation

**Files:**
- Modify: `apps/backend/internal/worker/result_consumer.go`

**Requirements:**
- **Acceptance Criteria**
  1. `HandleMessage` extracts `correlationId` from NSQ message body (if available) or generates new one.
  2. Context passed to `embedder` and `store` contains the `correlationId`.
  3. Logs include the `correlationId`.

- **Test Coverage**
  - [Unit] `HandleMessage` verifies context contains ID from mock message.

**Step 1: Write failing test**
Create `apps/backend/internal/worker/result_consumer_test.go` checking context metadata.

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/worker/... -v`

**Step 3: Implement Fixes**
- Parse message body to get ID (assuming message structure allows). If not, at least ensure `middleware.WithCorrelationID(ctx, id)` is called before passing context down.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/worker/... -v`

---

### Task 4: Fix Vector Schema Types

**Files:**
- Modify: `apps/backend/internal/vector/schema.go`

**Requirements:**
- **Acceptance Criteria**
  1. `sourceId` and `url` properties are defined as `DataTypeString` (or equivalent for exact match) instead of `DataTypeText`.
  2. Re-syncing does not duplicate chunks (verified via logic, or manual e2e if possible).

- **Test Coverage**
  - [Unit] Verify `ensureSchema` definition uses correct types.

**Step 1: Write failing test**
Inspect `schema.go` or write a test that checks the schema definition struct.

**Step 2: Verify test fails**
(Manual inspection or test runner)

**Step 3: Implement Fixes**
- Change `DataTypeText` to `DataTypeString` for `sourceId` and `url`.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/vector/... -v`

---

### Task 5: Clean Up Python Logging

**Files:**
- Modify: `apps/ingestion-worker/main.py`

**Requirements:**
- **Acceptance Criteria**
  1. No import of `logging` (standard library).
  2. Only `structlog` is used for logging.
  3. Logs are structured JSON in production (or consistent with `config.py`).

- **Test Coverage**
  - [Manual] Run worker, check logs are JSON/structured and not double-logged.

**Step 1: Verify current state**
`grep "import logging" apps/ingestion-worker/main.py`

**Step 2: Implement Fixes**
- Remove `import logging` and `logging.basicConfig`.
- Ensure `structlog.configure` handles stdout.

**Step 3: Verify fix**
Run worker locally (if environment permits) or check imports.
