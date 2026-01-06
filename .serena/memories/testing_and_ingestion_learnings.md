# Testing and Ingestion Learnings (Updated 2026-01-06)

## Ingestion Worker
- **Architecture**: Hybrid async/sync model. Web crawling is async (crawl4ai), file processing is sync (docling) offloaded to `pebble.ProcessPool`.
- **Testing**: Heavy reliance on mocking (pebble, crawl4ai).
- **Recent Improvements (2026-01-06)**:
    - **Metadata Extraction**: Logic extracted to pure functions (`extract_metadata_from_doc`). Handled edge cases (callables, NoneTypes) defensively.
    - **Zombie Tasks**: `touch_loop` now uses `asyncio.wait_for(event.wait())` for immediate exit, preventing zombie processes on cancellation.
    - **Concurrency**: Global `WORKER_SEMAPHORE` (8) enforced in `main.py` for all task types.
    - **Error Handling**: `correlation_id` added to all NSQ failure payloads.

## Testing Strategy Updates
- **Metadata**: Use `pytest.mark.parametrize` for table-driven testing of extraction logic.
- **Concurrency**: explicit semaphore saturation tests required.
- **Logging**: Must verify stdlib bridge to structlog.
