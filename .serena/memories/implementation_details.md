# Implementation Details

## Backend Architecture (Go)
The backend follows a **Feature-Based Architecture** (`apps/backend/features/`), grouping logic by domain rather than technical layer.

### Core Features
- **Source (`features/source`)**: Manages ingestion sources (Web/File).
  - Uses `PostgresRepo` for metadata and state (`source_pages` table).
  - Publishes tasks to NSQ (`ingest.task`).
  - Handles page-level status tracking (Pending -> Processing -> Completed/Failed).
- **Job (`features/job`)**: Manages failed ingestion tasks (DLQ).
  - **Dead Letter Queue**: Failed worker tasks are saved to `failed_jobs` table via `JobRepository`.
  - **Retry Mechanism**: `POST /jobs/{id}/retry` re-publishes the `original_payload` to NSQ.
- **MCP (`features/mcp`)**: Implements Model Context Protocol.
  - Supports both SSE (`/mcp/sse`) and JSON-RPC (`/mcp/messages`).
  - Integrates with `retrieval` service for RAG.

### Ingestion Worker (Python)
The worker is a distributed consumer built with `pynsq`, `asyncio`, and `crawl4ai`.

- **Reliability**:
  - **Robust Touch Loop**: Runs in background to keep NSQ connection alive. Cancels main task if connection drops (`StreamClosedError`).
  - **Robots.txt**: Enforced via `crawl4ai` config to ensure politeness.
  - **Error Reporting**: Captures `original_payload` on failure and sends to `ingest.result` with `status: failed`.
- **Handlers**:
  - `web.py`: Uses `AsyncWebCrawler` (Chromium) + `LLMContentFilter` (Gemini Flash) for content extraction.
  - `file.py`: Uses `docling` for local file conversion (PDF/Docx).

## Frontend Architecture (Vue 3)
Built with Vite, Pinia, and TailwindCSS.

- **Dashboard**: Displays real-time stats (Sources, Docs, Failed Jobs).
- **Failed Jobs Manager**:
  - View (`/jobs`): Lists failed tasks with error details and JSON payload inspection.
  - Action: Manual retry trigger.
- **Proxy**: Dev server proxies `/api` -> `localhost:8081`.

## Data Flow
1. **Ingestion**: User -> Backend (Create Source) -> NSQ (`ingest.task`) -> Worker (Crawl) -> NSQ (`ingest.result`) -> Backend (ResultConsumer).
2. **Failure**: Worker (Error) -> NSQ (`ingest.result` w/ Error) -> Backend -> `failed_jobs` table.
3. **Retry**: User -> Backend (Retry Endpoint) -> NSQ (`ingest.task` w/ Original Payload).
