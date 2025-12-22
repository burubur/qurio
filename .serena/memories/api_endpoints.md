# API Endpoints

Base URL: `/api` (Proxied via Nginx to Backend :8081)

## Sources
| Method | Endpoint | Description | Payload/Params |
| :--- | :--- | :--- | :--- |
| `GET` | `/sources` | List all active sources | - |
| `POST` | `/sources` | Create new source | `{"url": "string"}` |
| `DELETE` | `/sources/{id}` | Soft delete source | - |
| `POST` | `/sources/{id}/resync` | Trigger re-ingestion | - |

## Settings
| Method | Endpoint | Description | Payload/Params |
| :--- | :--- | :--- | :--- |
| `GET` | `/settings` | Get current config | - |
| `PUT` | `/settings` | Update config | `{"gemini_api_key": "...", "rerank_provider": "...", "rerank_api_key": "..."}` |

## MCP (Model Context Protocol)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/mcp` | Legacy JSON-RPC 2.0 Endpoint |
| `GET` | `/mcp/sse` | SSE Transport Connection (Yields Session ID) |
| `POST` | `/mcp/messages` | Send JSON-RPC Messages (Requires `?sessionId=...`) |

## Health
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Service health check |
