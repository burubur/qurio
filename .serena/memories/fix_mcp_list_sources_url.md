# Fix: Added URL to `qurio_list_sources` response

**Date:** 2026-01-13
**Issue:** The `qurio_list_sources` MCP tool was returning a list of sources without the `url` field, although the field exists in the database and is useful for the user.
**Resolution:** Updated `apps/backend/features/mcp/handler.go` to include the `URL` field in the `SimpleSource` struct and populate it from the source object.

## Changes
- Modified `apps/backend/features/mcp/handler.go`:
    - Added `URL` field to `SimpleSource` struct.
    - Populated `URL` field in the response generation loop.

## Verification
- Ran `go test -v ./apps/backend/features/mcp/...` which passed.
