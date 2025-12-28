# Ingestion DLQ Fix

Date: 2025-12-28

## Context
The "Ingestion Error Handling" plan was largely implemented but lacked the final integration step to save failed jobs to the persistent storage for retries.

## Changes
1. **Ingestion Worker (`apps/ingestion-worker/main.py`)**: Updated to include the `original_payload` in the `ingest.result` message when a task fails.
2. **Backend Result Consumer (`apps/backend/internal/worker/result_consumer.go`)**: Updated to read `original_payload` and save a `FailedJob` record via `jobRepo.Save` when a failure is reported.

## Result
Failed ingestion tasks are now persisted in the `failed_jobs` table and can be retried using the `POST /jobs/{id}/retry` endpoint.
