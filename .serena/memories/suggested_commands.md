# Suggested Commands

## Deployment
- **Start System:** `docker-compose up -d`
- **Stop System:** `docker-compose down`
- **View Logs:** `docker-compose logs -f`

## Development
- **Backend (Go):** `cd apps/backend && go run main.go`
- **Frontend (Vue):** `cd apps/frontend && npm run dev`
- **Ingestion Worker (Python):** `cd apps/ingestion-worker && source venv/bin/activate && python main.py`

## Testing
- **Backend:** `cd apps/backend && go test ./...`
- **Frontend:** `cd apps/frontend && npm run test`
- **Ingestion Worker:** `cd apps/ingestion-worker && PYTHONPATH=. ./venv/bin/pytest`
- **E2E Tests:** `cd apps/e2e && npx playwright test`

## Verification
- **Health Check:** `curl http://localhost:8081/health`
- **MCP Endpoint:** `http://localhost:8081/mcp`
- **Admin UI:** `http://localhost:3000`

## Utilities
- **Linting (Go):** `golangci-lint run`
- **Linting (TS):** `npm run lint`
- **Formatting:** `go fmt ./...`, `prettier --write .`
