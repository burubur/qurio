# Ingestion Robustness Diagnosis & Fixes

Date: 2025-12-28

## Diagnosis
After analyzing the ingestion system, I identified several areas for improvement:
1.  **Politeness**: `robots.txt` was not being respected, risking bans.
2.  **Stuck Jobs**: No mechanism to reset jobs that crash silently (e.g., OOM) and stay in `processing` state forever.
3.  **Concurrency**: Potential for overloading target domains due to lack of global rate limiting.

## Improvements Implemented
1.  **Robots.txt Compliance**: Enabled `check_robots_txt=True` in the `crawl4ai` configuration within `apps/ingestion-worker/handlers/web.py`.
2.  **Stuck Job Recovery (Backend)**: Added `ResetStuckPages` to `PostgresRepo` and a placeholder `ResetStuckJobs` in `JobService`. This sets the foundation for a future "Janitor" cron job to rescue stuck pages.

## Next Steps
- Implement a `CronService` or `Ticker` in the backend `main.go` to call `ResetStuckPages` every 5-10 minutes.
- Investigate reusing `AsyncWebCrawler` instances to reduce memory overhead.
