# Project Overview

Qurio is an open-source, localhost-only context operating system for AI agents. It ingests documentation (web, files) and exposes it via MCP.

**Current Status:** MVP Implementation Phase.
- **Completed:** 
    - Core Architecture (Go/Weaviate/Postgres/Vue)
    - Source CRUD (URL-based with Depth/Exclusion config)
    - Source Details View (Chunk visualization)
    - Recursive Web Crawler (Custom Go implementation)
    - Settings Management
    - Retrieval (MCP Endpoint, Hybrid Search)
    - Frontend Polling & Error Handling
- **In Progress:** 
    - File Uploads Integration
- **Next:**
    - Query Observability (Logging)
    - Advanced Filtering (Date, Author)

**Critical Constraints:**
- Localhost only (no auth)
- Testability-First (Mock adapters)
- Pure Go backend (mostly)
