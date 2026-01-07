# Plan: Backend Test Coverage Boost - Part 6

## Context
While the backend has a 100% test pass rate, the statement coverage is at **67.9%**.
Critical components like `qurio/apps/backend/features/mcp`, `qurio/apps/backend/internal/adapter/weaviate`, and `qurio/apps/backend/internal/app` have coverage gaps in error handling and specific methods.

## Coverage Gaps Identification

| Component | Coverage | Missing Areas |
| :--- | :--- | :--- |
| **MCP Handler** | 70.1% | `HandleMessage` (52.8%) error paths, `writeError` (0%), `ServeHTTP` (70%) edge cases. |
| **Weaviate Adapter** | 47.3% | `GetChunks` (0%), `GetChunksByURL` (0%), `Search` (86.4%) edge cases. |
| **App Wiring** | 66.9% | `Run` (0%), `GetSourceDetails` (0%), `GetSourceConfig` (0%), `BulkCreatePages` (0%). |
| **Gemini Adapter** | 57.4% | `NewEmbedder` (0%), `Embed` (0%) in `embedder.go` (deprecated?), `Embed` in dynamic embedder error paths. |
| **Job Service** | 81.8% | `Count` (0%), `ResetStuckJobs` (0%). |
| **Reranker** | 60.2% | `NewDynamicClient` (0%), `Rerank` (0%) in dynamic client. |

## Tasks

### Task 1: Job Service Coverage
**Goal:** Cover `Count` and `ResetStuckJobs`.
- **Files:** `apps/backend/features/job/service_test.go`
- **Action:** Add unit tests for `Count` and `ResetStuckJobs` (even if ResetStuckJobs returns 0/nil currently, verify it calls what it needs or returns expected default).

### Task 2: Weaviate Adapter Coverage
**Goal:** Cover `GetChunks` and `GetChunksByURL`.
- **Files:** `apps/backend/internal/adapter/weaviate/store_test.go`
- **Action:** Add unit/integration tests for these methods. They are likely used by retrieval service but not tested directly in adapter tests.

### Task 3: MCP Handler Hardening
**Goal:** Improve `HandleMessage` and `writeError` coverage.
- **Files:** `apps/backend/features/mcp/handler_test.go`
- **Action:** Add test cases for `HandleMessage` where `json.Unmarshal` fails (malformed body), or `sessionId` is missing (already covered? check why 52.8%). The low coverage might be due to the async goroutine not being tracked well or error paths in the goroutine.

### Task 4: App Wiring Helper Methods
**Goal:** Cover `GetSourceDetails`, `GetSourceConfig`, `BulkCreatePages`.
- **Files:** `apps/backend/internal/app/app_test.go`
- **Action:** These seem to be helper methods in `app.go` that might be used by the worker or other components. If they are unused, **delete them**. If used, test them.
    - *Investigation:* Check usages of `GetSourceDetails`, `GetSourceConfig`, `BulkCreatePages` in `apps/backend`.

### Task 5: Cleanup Deprecated Code
**Goal:** Remove unused code to improve density.
- **Investigation:** `qurio/apps/backend/internal/adapter/gemini/embedder.go` has 0% coverage. Is it used? `dynamic_embedder.go` seems to be the active one.
- **Action:** If unused, delete `embedder.go`.
- **Investigation:** `qurio/apps/backend/internal/adapter/reranker/dynamic_client.go` has 0% coverage. Is it used?

## Execution Strategy
1.  **Investigate & Prune:** Check usages of 0% coverage methods. Delete if dead code.
2.  **Test Gaps:** Write tests for the remaining live methods with low coverage.
