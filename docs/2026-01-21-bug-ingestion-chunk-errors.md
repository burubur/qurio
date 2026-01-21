## Analysis & Resolution (Jan 21 2026)

### Root Cause Analysis
1.  **Resource Exhaustion:** High concurrency (50 threads x 10 replicas) combined with spawning a **new browser instance for every URL** caused massive CPU/RAM overhead.
2.  **Redundant Probes:** The worker was aggressively probing `llms.txt` for every single page visited, doubling the crawl load.
3.  **Timeouts:**
    - `Page.goto` timeout (60s) was insufficient for heavy pages (e.g. `vuejs.org/llms.txt` with >3000 chunks) under load.
    - **NSQ Connection Drop:** Heavy parsing blocked the Python Event Loop for >60s, causing the NSQ client to miss heartbeats. The server would then close the connection ("Stream is closed"), causing the worker to fail the task acknowledgment (`FIN`) or requeue (`REQ`), leading to "stuck" pending tasks.

### Fixes Applied
1.  **Global Persistent Browser:** Refactored `apps/ingestion-worker/main.py` to initialize a single `AsyncWebCrawler` instance at startup and reuse it for all tasks. This significantly reduced overhead.
2.  **Optimized Logic:** Removed the per-page `llms.txt` check.
3.  **Configuration Tuning:**
    - Increased `CRAWLER_PAGE_TIMEOUT` to **120s** (from 60s) in `apps/ingestion-worker/config.py`.
    - Increased `NSQ_HEARTBEAT_INTERVAL` to **60s** (from 30s) in `apps/ingestion-worker/main.py` to prevent disconnects during heavy CPU operations.

### Outcome
- Worker stability restored.
- "Stuck" tasks at the end of ingestion queues were resolved by preventing NSQ disconnects.
- System can now handle heavy single-page documentations (like Vue.js) without crashing or timing out.
