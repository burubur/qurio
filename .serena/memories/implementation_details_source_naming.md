# Source Naming Refactor (Implemented)

**Date:** 2026-01-13
**Status:** Completed

## Changes
1.  **Mandatory Naming:** All sources (Web and File) now require a `name` field at creation time.
2.  **API Update:**
    - `POST /sources`: JSON body must include `name`.
    - `POST /sources/upload`: Multipart form must include `name` field.
3.  **Persistence:** The `name` is saved to the `sources` table in Postgres immediately upon creation.
4.  **Propagation:** The Ingestion Worker retrieves the name from the DB (via `SourceFetcher`) and attaches it to chunks sent to `ingest.embed`.
5.  **Frontend:** `SourceForm.vue` includes a required Name input. `source.store.ts` handles the transmission.

## Verification
-   **Backend Tests:** Added `handler_create_test.go` and `handler_upload_test.go` to enforce validation. Updated `handler_integration_test.go` and `topic_integration_test.go`.
-   **Worker Tests:** Verified `ResultConsumer` uses the configured name in `result_consumer_test.go`.
-   **E2E Tests:** Updated `ingestion.spec.ts` and `failure-retry.spec.ts` to provide names during source creation.

## Technical Details
-   `Handler.Create` and `Handler.Upload` return 400 Bad Request if `name` is missing.
-   `Service.Upload` signature changed to `Upload(ctx, path, hash, name)`.
-   `ResultConsumer` fetches source config (including name) before processing chunks.
