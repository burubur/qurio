# Parallel Ingestion Implementation

**Date:** 2026-01-12
**Status:** Implemented

## Overview
The ingestion pipeline has been refactored to support parallel processing of Web and File sources using dedicated NSQ topics and worker instances. This prevents long-running web crawls from blocking fast file ingestion tasks.

## Changes
1.  **Split Topics:**
    *   `ingest.task.web`: For web crawling tasks.
    *   `ingest.task.file`: For file processing tasks.
    *   `ingest.result`: Shared result topic (unchanged).

2.  **Backend Refactoring:**
    *   `apps/backend/internal/config/topics.go`: Defined constants for topics.
    *   `SourceService` and `JobService`: Updated to route tasks to the correct topic based on `source.Type`.
    *   `ResultConsumer`: Updated to publish discovered links to `ingest.task.web`.
    *   `Bootstrap`: Updated to pre-create all topics.

3.  **Infrastructure:**
    *   `docker-compose.yml`: Replaced single `ingestion-worker` with `ingestion-worker-web` and `ingestion-worker-file`.
    *   Workers are configured via `NSQ_TOPIC_INGEST` environment variable.

## Verification
*   **Unit Tests:** Backend tests updated to verify correct topic selection.
*   **Integration Test:** `apps/backend/internal/worker/topic_integration_test.go` verifies end-to-end routing with real NSQ (via Testcontainers).
*   **Configuration:** Verified `apps/ingestion-worker/config.py` respects environment variables.

## How to Run
```bash
docker compose up -d
```
This will start both worker services. Scaling can be done independently:
```bash
docker compose up -d --scale ingestion-worker-web=3
```
