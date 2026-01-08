Plan created for fixing chunks pagination limit.
Plan file: docs/plans/2026-01-08-chunks-pagination-1.md
Key changes:
- Backend: Add pagination and CountChunksBySource to Weaviate store.
- Backend: Update Service/Handler to support limit/offset.
- Frontend: Update Store and View to support "Load More".