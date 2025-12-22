# Tech Stack

**Core Services:**
- **Orchestration/API:** Go 1.24+ (Standard Library `net/http`)
- **Frontend:** Vue 3 + TypeScript + Vite + Pinia + Vue Router
- **Document Processing:** Python 3.10+ (Docling)
- **Vector Database:** Weaviate OSS (v1.24+)
- **Metadata Database:** PostgreSQL (v15+)
- **Task Queue:** NSQ

**AI & ML:**
- **Embeddings:** Google Gemini Embedding API (`gemini-embedding-001`) - *Dynamically Configured*
- **Reranking:** Jina AI v2, Cohere Rerank v3 - *Dynamically Configured*
- **Protocol:** Model Context Protocol (MCP) JSON-RPC 2.0

**Infrastructure:**
- **Containerization:** Docker, Docker Compose
- **Migrations:** `golang-migrate`
