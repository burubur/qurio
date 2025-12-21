# Task Completion Definition

A task is considered complete when:
1. **Code Changes:** All requested features or fixes are implemented according to the PRD specifications.
2. **Verification:**
    - The system builds and starts successfully (`docker-compose up` or individual service builds).
    - Relevant tests pass (if available).
    - Linter checks pass (`golangci-lint`, `eslint`, `ruff`/`flake8`).
3. **Documentation:** Any new environment variables or configuration options are documented (e.g., in `.env.example`).
4. **Cleanup:** Temporary files created during development are removed.
