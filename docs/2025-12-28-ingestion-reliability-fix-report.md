# Ingestion Worker Reliability Fixes - 2025-12-28

## Summary
Addressed reliability issues in the ingestion worker where NSQ connection drops (`StreamClosedError`) caused unhandled exceptions and potential "zombie" tasks (tasks that continue processing but cannot be acknowledged).

## Changes

### 1. Robust Touch Loop (`apps/ingestion-worker/main.py`)
- **Issue:** The original `touch_loop` did not handle exceptions when the NSQ connection was lost.
- **Fix:** Implemented a `try-except` block within the `touch_loop`.
- **Behavior:** If `message.touch()` fails with a fatal error (`nsq.Error`, `StreamClosedError`), the loop now catches it, logs a warning, and **cancels the main processing task** (`current_task.cancel()`). This ensures the worker stops processing a job it can no longer acknowledge.

### 2. Error Handling for Producer (`apps/ingestion-worker/main.py`)
- **Issue:** `producer.pub` and `message.finish` calls could raise exceptions if the connection was lost, crashing the worker or leaving the loop in an inconsistent state.
- **Fix:** Wrapped `producer.pub` and `message.finish` calls in `try-except` blocks to log errors gracefully without crashing the main loop.

### 3. Gemini Configuration (`apps/ingestion-worker/handlers/web.py`)
- **Verification:** Confirmed that `LLMConfig` is set with `temperature=1.0` to avoid infinite loops/warnings with Gemini 3 models.

### 4. Testing
- **New Test:** Created `apps/ingestion-worker/tests/test_worker_reliability.py`.
- **Scope:** Validates the logic pattern used in `main.py` where a `touch_loop` monitors the connection and triggers a callback (cancellation) upon failure.
- **Status:** Passed.

## Verification
- Ran unit tests using `pytest` in the local environment.
- Code changes applied to `main.py`.
