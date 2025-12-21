# API Endpoints

## Backend (:8081)

### System
- `GET /health`: Health check (returns 200 OK).

### Source Management
- `POST /sources`: Create a new source.
  - Body: `{"url": "https://example.com"}`
  - Returns: Source object with ID.

### MCP (Model Context Protocol)
- `POST /mcp`: JSON-RPC 2.0 endpoint.
  - Method: `tools/call`
  - Tool: `search`
  - Arguments: `{"query": "search term"}`

## Docling Service (:8000)
- `POST /process`: OCR processing.
  - Multipart Form: `file`
  - Returns: `{"text": "markdown content"}`
