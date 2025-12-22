# Implementation Plan: MCP over SSE

**Objective:** Refactor MCP handler to support Server-Sent Events (SSE) transport, compliant with MCP specifications for clients like Claude Desktop.

**Architecture:**
- **GET /mcp/sse:** Establishes persistent connection.
    - Sends `endpoint` event with URL for posting messages.
    - Sends `id` event with session UUID.
- **POST /mcp/messages:** Accepts JSON-RPC messages associated with a session.
- **Session Management:** In-memory map `sessions map[string]chan string`.

**Implementation Steps:**

1.  **Refactor Handler Struct:**
    - Add `sessions` map and `sync.RWMutex`.
    - Extract `processRequest(JSONRPCRequest) JSONRPCResponse` logic from `ServeHTTP` to be reusable.

2.  **Implement `HandleSSE`:**
    - Generate Session ID.
    - Set headers (`Content-Type: text/event-stream`, `Cache-Control: no-cache`).
    - Register session channel.
    - Send initial `endpoint` event (pointing to `/mcp/messages?sessionId=...`).
    - Loop receiving from channel -> writing to ResponseWriter.
    - Cleanup on disconnect.

3.  **Implement `HandleMessage`:**
    - Parse `sessionId` query param.
    - Validate session exists.
    - Parse JSON-RPC body.
    - Process request (using extracted logic).
    - Send JSON-RPC response **into the session channel** (not the HTTP response body of the POST).
    - Return `202 Accepted` to the POST request.

4.  **Update `main.go`:**
    - Register new endpoints.
    - Keep (or deprecate) existing `POST /mcp` for direct simple clients if desired, or fully switch. *Decision: Keep basic POST at /mcp for backwards compatibility or simpler testing, add /mcp/sse and /mcp/messages.*

5.  **Verify:**
    - Integration Test with `curl` (SSE stream) + `curl` (POST message).
    - Update `handler_test.go`.

**Files to Modify:**
- `apps/backend/features/mcp/handler.go`
- `apps/backend/main.go`
- `apps/backend/features/mcp/handler_test.go`
