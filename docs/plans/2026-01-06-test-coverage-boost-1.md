### Task 1: Define Interfaces & Mocks

**Files:**
- Create: `apps/backend/internal/app/interfaces.go`
- Create: `apps/backend/internal/app/mocks_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Define `Database` interface covering methods used by `app` (e.g., `Ping`, schema checks if any).
  2. Define `VectorStore` interface covering `EnsureSchema` and client methods used.
  3. Define `TaskPublisher` interface covering `Publish`.
  4. Create `MockDatabase`, `MockVectorStore`, `MockTaskPublisher` in `mocks_test.go` implementing these interfaces.

- **Functional Requirements**
  1. Interfaces must match the method signatures of the concrete types currently used in `app.go`.

- **Non-Functional Requirements**
  None for this task.

- **Test Coverage**
  - None (Interfaces and mocks are dependencies for other tests).

**Step 1: Write failing test**
*Skipped (Definition task)*

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
```go
// apps/backend/internal/app/interfaces.go
package app

import (
	"context"
	"database/sql"
)

type Database interface {
	PingContext(ctx context.Context) error
    // Add other methods as needed by app.go
}

type VectorStore interface {
	EnsureSchema(ctx context.Context) error
    // Add other methods
}

type TaskPublisher interface {
	Publish(topic string, body []byte) error
}

// apps/backend/internal/app/mocks_test.go
package app

type MockDatabase struct {
	PingErr error
}
func (m *MockDatabase) PingContext(ctx context.Context) error { return m.PingErr }

// ... Implement other mocks
```

**Step 4: Verify test passes**
*Skipped*


### Task 2: Refactor `app` Package (I/O Isolation)

**Files:**
- Modify: `apps/backend/internal/app/app.go`

**Requirements:**
- **Acceptance Criteria**
  1. `App` struct fields changed from concrete types to interfaces defined in Task 1.
  2. `New` constructor signature changed to accept interfaces.
  3. `Run` method implemented to handle server startup.

- **Functional Requirements**
  1. `New` should initialize the App with provided dependencies.
  2. `Run` should start the HTTP server and block until context cancellation.

- **Non-Functional Requirements**
  - Dependency Injection principle (Rule 1).

- **Test Coverage**
  - Covered by Task 4.

**Step 1: Write failing test**
*Skipped (Refactoring existing code, tests will break then be fixed in Task 4)*

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
```go
// apps/backend/internal/app/app.go

type App struct {
	db    Database
	vec   VectorStore
	queue TaskPublisher
    // ...
}

func New(logger *slog.Logger, db Database, vec VectorStore, queue TaskPublisher, cfg Config) *App {
    // ...
}

func (a *App) Run(ctx context.Context) error {
    // Move logic from main.go:
    // 1. Setup router
    // 2. Start server
    // 3. Handle shutdown
}
```

**Step 4: Verify test passes**
*Skipped*


### Task 3: Refactor `main.go`

**Files:**
- Modify: `apps/backend/main.go`

**Requirements:**
- **Acceptance Criteria**
  1. `main.go` initializes concrete `sql.DB`, `weaviate.Client`, `nsq.Producer`.
  2. `main.go` calls `app.New` passing these concrete types (casted/compliant with interfaces).
  3. `main.go` calls `app.Run`.

- **Functional Requirements**
  1. Application starts up exactly as before.

- **Non-Functional Requirements**
  - `main.go` should be "thin" (wiring only).

- **Test Coverage**
  - Manual verification (Application startup).

**Step 1: Write failing test**
*Skipped (Main entry point)*

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
```go
// apps/backend/main.go
func main() {
    // ... init concrete deps ...
    application := app.New(logger, db, vec, producer, cfg)
    if err := application.Run(ctx); err != nil {
        logger.Error("application failed", "error", err)
        os.Exit(1)
    }
}
```

**Step 4: Verify test passes**
*Manual execution*


### Task 4: Test `app` Package

**Files:**
- Modify: `apps/backend/internal/app/app_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `TestNew` uses `MockDatabase`, `MockVectorStore`, etc.
  2. Verify `New` does not panic and sets fields correctly.
  3. `TestRun` (if feasible to mock listeners) or at least verify `Run` setup logic.

- **Functional Requirements**
  1. 95% coverage of `app.go`.

- **Non-Functional Requirements**
  - Fast execution (no real I/O).

- **Test Coverage**
  - `TestNew_Success`
  - `TestNew_NilDependency_ShouldPanic` (if we enforce checks)

**Step 1: Write failing test**
```go
func TestNew(t *testing.T) {
    mockDB := &MockDatabase{}
    app := New(logger, mockDB, ...)
    if app == nil {
        t.Fatal("expected app, got nil")
    }
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/app/... -v`

**Step 3: Write minimal implementation**
*Implied by Task 2, this task confirms it working with mocks*

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/app/... -v`


### Task 5: Refactor Source Handlers

**Files:**
- Modify: `apps/backend/features/source/handler.go`
- Modify: `apps/backend/features/source/handler_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `Get` handler checks for `ErrNoRows` (or service equivalent) and returns 404.
  2. `Delete` handler checks for 404.
  3. Tests cover these paths.

- **Functional Requirements**
  1. `GET /sources/{id}` returns 404 if not found.

- **Non-Functional Requirements**
  - Standard JSON error envelope.

- **Test Coverage**
  - `TestGet_NotFound`
  - `TestDelete_NotFound`

**Step 1: Write failing test**
```go
// apps/backend/features/source/handler_test.go
func TestGet_NotFound(t *testing.T) {
    // Setup mock service to return ErrNotFound
    // Call handler
    // Expect 404, will likely get 500
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/source/... -v`

**Step 3: Write minimal implementation**
```go
// apps/backend/features/source/handler.go
if errors.Is(err, sql.ErrNoRows) {
    writeError(w, http.StatusNotFound, "Source not found")
    return
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/source/... -v`


### Task 6: Refactor Job Handlers

**Files:**
- Modify: `apps/backend/features/job/handler.go`
- Modify: `apps/backend/features/job/handler_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `Get`/`List` handlers handle not found or empty correctly.
  2. Tests cover these paths.

- **Functional Requirements**
  1. Correct HTTP status codes.

- **Non-Functional Requirements**
  - Consistent error handling.

- **Test Coverage**
  - `TestRetry_NotFound`

**Step 1: Write failing test**
```go
func TestRetry_NotFound(t *testing.T) {
    // Mock service to return ErrNotFound
    // Expect 404
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/job/... -v`

**Step 3: Write minimal implementation**
```go
if errors.Is(err, job.ErrJobNotFound) {
    writeError(w, http.StatusNotFound, "Job not found")
    return
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/job/... -v`
