# Proposal: Integration Testing for Ingestion Worker with Testcontainers

**Date:** 2026-01-21
**Status:** Proposed

## 1. Context & Problem
Currently, the `ingestion-worker` (Python) relies heavily on **Unit Tests** (`pytest`) with mocked dependencies (`unittest.mock`). While this ensures the internal logic (handlers, retry calculations) is correct, it leaves a critical gap in verifying the **infrastructure integration**:

*   **NSQ Interaction:** We mock `nsq.Reader` and `nsq.Writer`. We don't automatically verify that the worker actually connects, consumes, acknowledges, requeues, or publishes messages correctly against a real NSQ daemon.
*   **Event Loop Behavior:** The interaction between `pynsq` (which relies on Tornado's IOLoop) and `asyncio` is complex. Mocks cannot catch deadlock regressions or loop integration issues.
*   **Regression Risk:** Bugs related to network timeouts, connection drops, or protocol mismatches can slip through unit tests and are currently only caught during manual `docker compose` testing or in production.

## 2. Proposed Solution
Introduce **Integration Tests** using **Testcontainers for Python**. This allows us to spin up disposable, real infrastructure (NSQ Lookupd, NSQD) within the test suite, run the actual worker code against it, and assert on the side effects.

### Key Components
1.  **Library:** `testcontainers-python` (specifically the `testcontainers[nsq]` or generic Docker support).
2.  **Scope:** Tests will reside in `apps/ingestion-worker/tests/integration/`.
3.  **CI/CD:** These tests will run in the CI pipeline (GitHub Actions), which supports nested Docker containers (Service Containers).

## 3. Implementation Plan

### Phase 1: Dependencies & Configuration
Add the following to `apps/ingestion-worker/requirements.txt` (dev dependencies):
```text
testcontainers>=3.7.1
```

### Phase 2: Test Fixtures (`conftest.py`)
Create a Pytest fixture that manages the lifecycle of the NSQ container.

```python
# apps/ingestion-worker/tests/concurrency_fixture.py (conceptual)
import pytest
from testcontainers.core.container import DockerContainer

@pytest.fixture(scope="session")
def nsq_container():
    # Spin up nsqd
    container = DockerContainer("nsqio/nsq:latest") \
        .with_command("/nsqd --lookupd-tcp-address=nsqlookupd:4160") \
        .with_exposed_ports(4150, 4151)
    
    with container as nsqd:
        yield nsqd
```

### Phase 3: The Integration Test Pattern
The test will follow this flow:
1.  **Setup:** Start NSQ container.
2.  **Arrange:** Publish a "poison pill" or "test task" message to `ingest.task.web` using a real `nsq.Writer` or HTTP request.
3.  **Act:** Start the Worker's `main()` function (or a refactored `Worker` class) in a separate thread/process or async task.
4.  **Assert:**
    *   **Success:** Poll the `ingest.result` topic to see if the success payload was published.
    *   **Retry:** Publish a task that is guaranteed to fail (e.g., unreachable URL), and query the NSQ stats endpoint (`http://localhost:4151/stats`) to verify `requeue_count` increments.

### Phase 4: Refactoring `main.py`
To make `main.py` testable, we need to decouple the startup logic from the global scope.
*   **Current:** `main()` creates global readers/writers and blocks forever.
*   **Required:** Refactor into a `Worker` class with `start()` and `stop()` methods, allowing tests to gracefully shut down the worker after assertions.

## 4. Benefits
*   **Confidence:** Guarantees that the worker talks the correct protocol to NSQ.
*   **Loop Safety:** Verifies that the `uvloop` + `Tornado` + `pynsq` bridge is stable.
*   **Refactoring Safety:** Allows us to safely upgrade dependencies (e.g., `pynsq` or `crawl4ai`) knowing integration points are covered.

## 5. Risks & Mitigation
*   **Complexity:** Testing async code that spans threads/processes is tricky. **Mitigation:** Use `pytest-asyncio` and keep the worker implementation in the same event loop if possible, or use a subprocess.
*   **Slowness:** Docker containers take seconds to start. **Mitigation:** Use `scope="session"` for the container fixture to pay the startup cost only once per test suite run.

## 6. Recommendation
Approve this proposal as a **Technical Debt / Infrastructure** task to be scheduled after the current "Ingestion Error" bug fix is deployed.
