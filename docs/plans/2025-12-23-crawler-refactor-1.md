--- 
name: technical-constitution
description: Refactor crawler and document processing to a Python worker using crawl4ai and docling, communicating via NSQ.
---

# Implementation Plan - Crawler & Docling Refactor

**Ref:** `2025-12-23-crawler-refactor-1`  
**Feature:** Ingestion Worker (Python) & NSQ Architecture  
**Status:** Draft

## 1. Scope
Refactor the current synchronous/monolithic Go crawler and Docling API into a distributed architecture. Introduce a Python-based `ingestion-worker` that handles both web crawling (via `crawl4ai`) and document processing (via `docling`). Communication will be handled asynchronously via NSQ.

### Architecture Change
**Current:**
`Go Backend` -> `Internal Crawler` (Sequential) -> `Docling API` (HTTP) -> `Embed/Store`

**New:**
1. `Go Backend` -> Publish `ingest.task` -> `NSQ`
2. `Ingestion Worker (Python)` -> Consume `ingest.task` -> `Crawl4AI` / `Docling` -> Publish `ingest.result` -> `NSQ`
3. `Go Backend` -> Consume `ingest.result` -> `Embed/Store`

## 2. Requirements
- **Functional:**
    - Support Web Crawling (URL) with `crawl4ai` (Markdown output).
    - Support Document Processing (File) with `docling` (Markdown output).
    - Asynchronous processing via NSQ.
    - Result payload must include raw content/markdown for embedding.
- **Non-Functional:**
    - Asyncio-based Python worker.
    - Dockerized with necessary dependencies (Playwright, PyTorch/Docling).
    - Error handling: Retry via NSQ or Dead Letter Queue (DLQ) strategy.

## 3. Tasks

### Task 1: Scaffold Ingestion Worker
**Files:**
- Create: `apps/ingestion-worker/Dockerfile`
- Create: `apps/ingestion-worker/requirements.txt`
- Create: `apps/ingestion-worker/main.py` (Skeleton)
- Create: `apps/ingestion-worker/.dockerignore`

**Requirements:**
- **Base Image:** `python:3.10-slim-buster` (or similar compatible with Playwright & Docling).
- **Dependencies:** `asyncnsq`, `crawl4ai`, `docling`, `uvloop` (optional).
- **System Deps:** `libsnappy-dev` (for asyncnsq), Playwright browsers.

**Requirements Enrichment:**
- Search: `crawl4ai dockerfile` (Used example from search).
- Search: `asyncnsq` (Confirmed asyncio support).

**Step 1: Write failing test**
*N/A - Infrastructure task. Verification via Docker build.*

**Step 3: Implementation**
```dockerfile
FROM python:3.10-slim-buster

ENV PYTHONUNBUFFERED=1
ENV PLAYWRIGHT_BROWSERS_PATH="/ms-playwright"

# System dependencies for Playwright & NSQ (snappy)
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libsnappy-dev \
    git \
    # Playwright deps (simplified list, use official script if possible)
    libwoff-dev libharfbuzz-dev libicu-dev libgirepository1.0-dev \
    libcairo2-dev libjpeg-dev libpng-dev libtool libnss3 libxss1 \
    libasound2 libatk-bridge2.0-0 libgtk-3-0 libgbm-dev libxkbcommon-x11-0 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Install Playwright browsers
RUN playwright install --with-deps chromium

COPY . .

CMD ["python", "main.py"]
```

### Task 2: Implement Python NSQ Consumer
**Files:**
- Modify: `apps/ingestion-worker/main.py`

**Requirements:**
- Connect to `nsqlookupd` (or `nsqd`).
- Subscribe to `ingest.task` topic, `worker` channel.
- Decode JSON payload.
- Dispatch to handler (placeholder).
- Publish result to `ingest.result`.

**Step 1: Write failing test**
Create `tests/test_nsq.py` using `pytest-asyncio` mocking `asyncnsq`.

**Step 3: Implementation**
Use `asyncnsq` to create a `Reader` and `Writer`.

### Task 3: Implement Web Crawler (Crawl4AI)
**Files:**
- Create: `apps/ingestion-worker/handlers/web.py`
- Modify: `apps/ingestion-worker/main.py` (Integration)

**Requirements:**
- Input: URL, Depth, Exclusions.
- Logic: Use `AsyncWebCrawler` from `crawl4ai`.
- Output: Markdown string.

**Step 1: Write failing test**
Mock `AsyncWebCrawler` and assert handler calls it.

**Step 3: Implementation**
```python
from crawl4ai import AsyncWebCrawler

async def handle_web_task(url: str, max_depth: int = 1):
    async with AsyncWebCrawler(verbose=True) as crawler:
        result = await crawler.arun(url=url)
        return result.markdown
```

### Task 4: Implement Document Processor (Docling)
**Files:**
- Create: `apps/ingestion-worker/handlers/file.py`
- Modify: `apps/ingestion-worker/main.py` (Integration)

**Requirements:**
- Input: File path (shared volume?) or URL to file. *Correction:* For MVP, `docling` might need the file. If `ingest.task` comes from Go, where is the file?
- *Strategy:* Go backend saves upload to `tmp` (volume shared) or Object Storage (MinIO).
- *Constraint:* Localhost only. Shared Volume `/tmp/qurio-uploads` is easiest.
- Update `docker-compose.yml` to share volume.

**Step 1: Write failing test**
Mock `DocumentConverter`.

**Step 3: Implementation**
```python
from docling.document_converter import DocumentConverter

converter = DocumentConverter()

def handle_file_task(file_path: str):
    # This might need to be run in run_in_executor if blocking
    res = converter.convert(file_path)
    return res.document.export_to_markdown()
```

### Task 5: Refactor Go Backend (Producer)
**Files:**
- Modify: `apps/backend/features/source/source.go`
- Modify: `apps/backend/internal/worker/ingest.go` (Delete old logic)

**Requirements:**
- `Source.Create` / `ReSync`: Publish `ingest.task` with `{ "type": "web", "url": "...", "id": "..." }`.
- Remove `internal/crawler`.
- Remove `internal/adapter/docling`.

**Step 1: Write failing test**
Test `Source.Create` calls `publisher.Publish`.

### Task 6: Refactor Go Backend (Result Consumer)
**Files:**
- Modify: `apps/backend/internal/worker/result_consumer.go` (New)
- Modify: `apps/backend/main.go` (Wire up)

**Requirements:**
- Consume `ingest.result` topic.
- Payload: `{ "source_id": "...", "content": "..." }`.
- Logic: Chunk -> Embed -> Store (Re-use existing logic from old `ingest.go`).

**Step 1: Write failing test**
Test `HandleMessage` parses result and calls `Store`.

### Task 7: Integration & Cleanup
**Files:**
- Modify: `docker-compose.yml`
- Delete: `services/docling` folder.

**Requirements:**
- Add `ingestion-worker` service.
- Mount shared volume for uploads (if needed).
- Ensure `nsqlookupd` connection.

**Verification:**
- `docker-compose up --build`
- Add a Source -> Verify logs in `ingestion-worker` -> Verify embeddings in Weaviate.

