# Backend Testing & Resilience Improvements (Jan 21, 2026)

## Link Discovery Hardening
- **Logic**: Updated `DiscoverLinks` to strictly enforce `http/https` schemes, rejecting `ftp`, `mailto`, etc.
- **Testing**: Added `TestDiscoverLinks_Comprehensive` (Table-driven) covering 15+ scenarios including Unicode, Fragments, Exclusions, and weird schemes.

## Result Consumer Reliability
- **Partial Failure**: `ResultConsumer` now returns error (triggering NSQ retry) if `BulkCreatePages` fails, ensuring no discovered links are lost even if Embedding succeeded.
- **Testing**: Added `TestResultConsumer_HandleMessage_PartialFailure_BulkCreatePages`.

## Infrastructure Resilience
- **Config**: Added `BOOTSTRAP_RETRY_ATTEMPTS` (default 10) and `BOOTSTRAP_RETRY_DELAY_SECONDS` (default 2).
- **Testing**: Added `bootstrap_resilience_test.go` using Testcontainers to verify behavior when DB or Weaviate is down/unreachable.

## Utility Improvements
- **QueryLogger**: Added `sync.Mutex` for thread-safety. Verified with concurrent test.
- **Reranker**: Improved error reporting to include response body from Jina/Cohere APIs.
