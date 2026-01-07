### 2026-01-07: Backend Test Coverage Improvements
- **Bootstrap Refactoring:**
    - Extracted retry logic for Weaviate schema initialization into `EnsureSchemaWithRetry` (pure function with retry loop).
    - Updated `Bootstrap` to use `EnsureSchemaWithRetry`.
    - Added unit tests for retry logic using `MockVectorStore`.
- **Integration Testing:**
    - Implemented `apps/backend/internal/app/bootstrap_integration_test.go` using `Testcontainers`.
    - Updated `IntegrationSuite` to expose container configuration via `GetAppConfig`.
    - Added support for running migrations in `IntegrationSuite` or deferring to `Bootstrap` (via `SkipMigrations` flag).
    - Added `MigrationPath` to `config.Config` to allow tests to override migration location.
- **Application Structure:**
    - Refactored `apps/backend/main.go` to extract `run(ctx, cfg, logger)` function.
    - Added `apps/backend/smoke_test.go` (package `main`) to verify full application startup and wiring by running `run` against Testcontainers.
    - Enhanced `app.New` unit tests to verify route registration for key endpoints.
- **Feature Hardening:**
    - **Source:** Implemented `Exclusions` regex validation in `Service.Create`. Added `POST /sources/upload` integration test verifying file persistence to `QURIO_UPLOAD_DIR`.
    - **Job:** Verified `Retry` timeout logic (5s limit on NSQ publish). Validated `failed_jobs` cascade delete via integration tests.
    - **MCP:** Hardened JSON-RPC handler with table-driven tests for edge cases (empty params, invalid values). Validated SSE session establishment and Correlation ID propagation in integration tests.

### 2026-01-07: Ingestion Worker Optimization
- **Dead Code Removal:**
    - Removed unused imports (`json`, `httpx`, `signal`, `ProcessPoolExecutor`) in `handlers/web.py` and `handlers/file.py`.
    - Removed unused `result_content` variable in `main.py`.
    - Removed unused `exclusions` parameter in `handle_web_task` (web handler).
- **Concurrency Optimization:**
    - Removed redundant `CONCURRENCY_LIMIT` semaphore (size 8) in `handle_file_task`. The global `WORKER_SEMAPHORE` (size 8) in `main.py` already throttles concurrency at the entry point.
