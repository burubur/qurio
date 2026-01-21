# Backend Test Coverage Improvement Plan
Date: 2026-01-21

## Objective
Increase test coverage for the `apps/backend` service, specifically targeting functions with low coverage, edge cases, and infrastructure resilience. Excludes MCP interface.

## identified Gaps & Proposed Tests

### 1. Ingestion Worker - Link Discovery
**File:** `apps/backend/internal/worker/link_discovery.go`
**Current Status:** Basic success path covered. Lacks robust edge case handling.
**Proposed Unit Tests:**
- [ ] `TestDiscoverLinks_MalformedHTML`: Handle unclosed tags, mixed attributes.
- [ ] `TestDiscoverLinks_ComplexStructures`: Nested divs, relative paths vs absolute paths.
- [ ] `TestDiscoverLinks_Filtering`: Verify ignore patterns, same-domain restrictions.
- [ ] `TestDiscoverLinks_EdgeCases`:
    - Non-HTTP schemes (mailto:, tel:, javascript:)
    - Circular redirects (if logic exists to follow them, otherwise detect link)
    - Unicode/Non-ASCII URLs
    - Extremely long URLs
- [ ] `TestExtractMetadata_Variations`:
    - Only OG tags
    - Only Twitter Card tags
    - Standard meta description/title only
    - Fallback logic (body text excerpt)

### 2. Ingestion Worker - Result Consumer
**File:** `apps/backend/internal/worker/result_consumer.go`
**Current Status:** Central coordination logic. Needs granular state/failure testing.
**Proposed Unit Tests (Mocked):**
- [ ] `TestHandleMessage_Success`: Happy path (DB + Vector + NSQ Ack).
- [ ] `TestHandleMessage_VectorFailure`: DB success, Vector fail -> Verify Rollback/Compensating transaction or Error Log.
- [ ] `TestHandleMessage_DBFailure`: DB fail -> Verify Nack.
- [ ] `TestHandleMessage_InvalidJSON`: Payload corruption handling.
- [ ] `TestHandleMessage_StaleData`: Race conditions (if versioning exists).

### 3. Retrieval Logger (File Query Logger)
**File:** `apps/backend/internal/retrieval/logger.go`
**Current Status:** File operations. Lacks concurrency/error tests.
**Proposed Tests:**
- [ ] `TestFileQueryLogger_ConcurrentWrites`: Parallel execution to verify thread safety.
- [ ] `TestFileQueryLogger_FileErrors`: Simulate permission denied or disk full (using interface/mocking or temp dir restrictions).
- [ ] `TestFileQueryLogger_Rotation`: If log rotation exists, verify it triggers correctly.

### 4. Reranker Dynamic Client
**File:** `apps/backend/internal/adapter/reranker/dynamic_client.go`
**Current Status:** Dynamic configuration.
**Proposed Unit Tests:**
- [ ] `TestRerank_ProviderError_5xx`: Handle upstream server errors.
- [ ] `TestRerank_ProviderError_4xx`: Handle client errors (auth, bad request).
- [ ] `TestRerank_EmptyInput`: Validation.

### 5. Infrastructure / Bootstrap
**File:** `apps/backend/internal/app/bootstrap.go`
**Current Status:** Basic integration.
**Proposed Integration Tests (Testcontainers):**
- [ ] `TestBootstrap_DependencyFailures`:
    - Postgres down at startup.
    - Weaviate down at startup.
    - NSQ down at startup.
- [ ] `TestBootstrap_Recovery`: Dependencies become available after initial delay (if retry logic exists).

## Execution Strategy
1.  **Phase 1:** Implement Unit Tests for `link_discovery.go` (High value, low dependency).
2.  **Phase 2:** Implement Unit Tests for `result_consumer.go` (Critical path reliability).
3.  **Phase 3:** Implement Unit Tests for `logger.go` & `dynamic_client.go`.
4.  **Phase 4:** Implement Integration Tests for Bootstrap/Infra resilience.
