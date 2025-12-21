# Project Overview

**Status:** MVP Implementation Complete (2025-12-21)

## Components
1.  **Backend (Go):** Fully implemented with clean architecture.
    - Features: Source Management, MCP Endpoint, Ingestion Worker.
    - Infrastructure: Postgres (Metadata), Weaviate (Vectors), NSQ (Messaging).
    - Adapters: Gemini (Embeddings), Docling (OCR), Weaviate (Store).
2.  **Frontend (Vue 3):**
    - Features: Source List, Add Source Form.
    - Tech: Vite, TypeScript, Vitest.
3.  **Services:**
    - Docling: Python FastAPI service for OCR.

## Verification
- CI/CD: All unit tests passed (Backend + Frontend + Docling).
- Infrastructure: Docker Compose ready.
