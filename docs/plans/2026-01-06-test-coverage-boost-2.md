# Backend Test Coverage Boost - Execution Plan

### Task 1: Decouple main.go - Bootstrap Logic

**Files:**
- Create: `apps/backend/internal/app/bootstrap.go`
- Modify: `apps/backend/main.go:28-150`
- Test: `apps/backend/internal/app/bootstrap_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `Bootstrap(cfg)` function initializes DB, Weaviate, NSQ, and Migrations.
  2. Retry logic for DB and Weaviate is encapsulated within `Bootstrap`.
  3. `main.go` logic is reduced to Config Load -> Bootstrap -> App.New -> App.Run.

- **Functional Requirements**
  1. Infrastructure initialization must remain identical in behavior (retries, timeouts).

- **Non-Functional Requirements**
  1. Log messages must remain consistent.
  2. Error handling must cause immediate failure (panic or exit in main, error return in bootstrap).

- **Test Coverage**
  - [Unit] `Bootstrap` handles connection failures.

**Step 1: Write failing test**
```go
package app_test

import (
	"context"
	"testing"
	"qurio/apps/backend/internal/app"
	"qurio/apps/backend/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestBootstrap_ConfigurationError(t *testing.T) {
	cfg := &config.Config{
		DBHost: "invalid-host",
	}
	deps, err := app.Bootstrap(context.Background(), cfg)
	assert.Error(t, err)
	assert.Nil(t, deps)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/app/... -v`
Expected: FAIL (Undefined `Bootstrap`)

**Step 3: Write minimal implementation**
```go
package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"net"

	"qurio/apps/backend/internal/config"
	wstore "qurio/apps/backend/internal/adapter/weaviate"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nsqio/go-nsq"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

type Dependencies struct {
	DB          *sql.DB
	VectorStore VectorStore
	NSQProducer *nsq.Producer
}

func Bootstrap(ctx context.Context, cfg *config.Config) (*Dependencies, error) {
	// Database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// Retry loop
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		slog.Warn("failed to ping db, retrying...", "attempt", i+1)
		time.Sleep(2 * time.Second)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	// Migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("migration driver error: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("migration instance error: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("migration up error: %w", err)
	}

	// Weaviate
	wCfg := weaviate.Config{Host: cfg.WeaviateHost, Scheme: cfg.WeaviateScheme}
	wClient, err := weaviate.NewClient(wCfg)
	if err != nil {
		return nil, fmt.Errorf("weaviate client error: %w", err)
	}
	vecStore := wstore.NewStore(wClient)
	
	// Ensure Schema Retry
	for i := 0; i < 10; i++ {
		if err := vecStore.EnsureSchema(ctx); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err := vecStore.EnsureSchema(ctx); err != nil {
		return nil, fmt.Errorf("weaviate schema error: %w", err)
	}

	// NSQ Producer
	nsqCfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(cfg.NSQDHost, nsqCfg)
	if err != nil {
		return nil, fmt.Errorf("nsq producer error: %w", err)
	}
    
    // Topic pre-creation (Logic from main.go)
    createTopics(cfg.NSQDHost)

	return &Dependencies{
		DB:          db,
		VectorStore: vecStore,
		NSQProducer: producer,
	}, nil
}

func createTopics(nsqdHost string) {
    nsqHttpURL := fmt.Sprintf("http://%s:4151/topic/create?topic=ingest.task", "nsqd")
	nsqResultURL := fmt.Sprintf("http://%s:4151/topic/create?topic=ingest.result", "nsqd")
	
	host, _, _ := net.SplitHostPort(nsqdHost)
	if host != "" {
		nsqHttpURL = fmt.Sprintf("http://%s:4151/topic/create?topic=ingest.task", host)
		nsqResultURL = fmt.Sprintf("http://%s:4151/topic/create?topic=ingest.result", host)
	}
	
	go func() {
		time.Sleep(2 * time.Second)
		http.Post(nsqHttpURL, "application/json", nil)
		http.Post(nsqResultURL, "application/json", nil)
	}()
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/app/... -v`
Expected: PASS

### Task 2: Refactor main.go

**Files:**
- Modify: `apps/backend/main.go`

**Requirements:**
- **Acceptance Criteria**
  1. `main.go` uses `app.Bootstrap`.
  2. Application starts successfully.

- **Functional Requirements**
  1. No loss of functionality.

**Step 1: Write failing test**
(N/A - Refactoring main entry point)

**Step 2: Verify test fails**
(N/A)

**Step 3: Write minimal implementation**
```go
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"qurio/apps/backend/internal/app"
	"qurio/apps/backend/internal/config"
	"qurio/apps/backend/internal/logger"
)

func main() {
	logger := slog.New(logger.NewContextHandler(slog.NewJSONHandler(os.Stdout, nil)))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	deps, err := app.Bootstrap(context.Background(), cfg)
	if err != nil {
		slog.Error("bootstrap failed", "error", err)
		os.Exit(1)
	}
	defer deps.DB.Close()

	application, err := app.New(cfg, deps.DB, deps.VectorStore, deps.NSQProducer, logger)
	if err != nil {
		slog.Error("app init failed", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := application.Run(ctx); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
```

**Step 4: Verify test passes**
Run: `go build ./apps/backend`
Expected: Success

### Task 3: ResultConsumer Deep Testing

**Files:**
- Modify: `apps/backend/internal/worker/result_consumer_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Tests cover `HandleMessage` with invalid JSON.
  2. Tests cover `HandleMessage` with missing fields.
  3. Tests cover `HandleMessage` with dependency errors.

**Step 1: Write failing test**
```go
func TestResultConsumer_HandleMessage_InvalidJSON(t *testing.T) {
    // Setup minimal mocks
    consumer := &ResultConsumer{Logger: slog.Default()}
    msg := &nsq.Message{Body: []byte("{invalid-json")}
    
    // Should return nil to ack message (don't retry poison pill)
    // Or return error if retry desired. Assuming we swallow bad JSON to avoid loop.
    err := consumer.HandleMessage(msg)
    assert.NoError(t, err) // Should handle gracefully
}

func TestResultConsumer_HandleMessage_DependencyError(t *testing.T) {
    mockEmbedder := new(MockEmbedder)
    mockEmbedder.On("Embed", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("api error"))
    
    consumer := NewResultConsumer(mockEmbedder, ...)
    // Create valid message
    msg := &nsq.Message{Body: []byte(`{"job_id":"1", "content":"test"}`)}
    
    // Should return error to trigger NSQ requeue
    err := consumer.HandleMessage(msg)
    assert.Error(t, err)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/worker/... -v`

**Step 3: Write minimal implementation**
(Existing implementation should handle these, tests verify behavior. Adjust code if tests fail unexpectedly.)

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/worker/... -v`

### Task 4: MCP Tools Table-Driven Tests

**Files:**
- Modify: `apps/backend/features/mcp/handler_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Table-driven test for `processRequest` covering all tools.

**Step 1: Write failing test**
```go
func TestHandler_ProcessRequest_Table(t *testing.T) {
    // Define Mocks
    mockRetriever := new(MockRetriever)
    mockSourceMgr := new(MockSourceManager)
    handler := NewHandler(mockRetriever, mockSourceMgr)

    // Define Cases
    tests := []struct {
        name    string
        req     JSONRPCRequest
        setup   func()
        wantRes func(*JSONRPCResponse) bool
        wantErr bool
    }{
        {
            name: "Initialize",
            req:  JSONRPCRequest{Method: "initialize", ID: 1},
            setup: func() {},
            wantRes: func(r *JSONRPCResponse) bool {
                return r.Result.(map[string]interface{})["protocolVersion"] == "2024-11-05"
            },
        },
        {
            name: "List Tools",
            req:  JSONRPCRequest{Method: "tools/list", ID: 2},
            setup: func() {},
            wantRes: func(r *JSONRPCResponse) bool {
                 res := r.Result.(ListToolsResult)
                 return len(res.Tools) > 0
            },
        },
        {
             name: "Call Unknown Tool",
             req: JSONRPCRequest{Method: "tools/call", Params: json.RawMessage(`{"name": "unknown"}`), ID: 3},
             setup: func() {},
             wantRes: func(r *JSONRPCResponse) bool {
                 return r.Error != nil && r.Error.(map[string]interface{})["code"].(int) == ErrMethodNotFound
             },
        },
        // Add specific tool calls (Search, ReadPage) with mocks
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            res := handler.processRequest(context.Background(), tt.req)
            if tt.wantErr {
                assert.Nil(t, res)
            } else {
                assert.NotNil(t, res)
                assert.True(t, tt.wantRes(res))
            }
        })
    }
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/mcp/... -v`

**Step 3: Write minimal implementation**
(Implement the test suite logic fully in the file)

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/mcp/... -v`

### Task 5: Feature Handlers & Adapters

**Files:**
- Modify: `apps/backend/internal/adapter/weaviate/store_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Weaviate adapter handles network errors gracefully.

**Step 1: Write failing test**
```go
func TestStore_Search_NetworkError(t *testing.T) {
    // 1. Start a server that always fails
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer server.Close()

    // 2. Configure Weaviate client to point to this server
    cfg := weaviate.Config{
        Host:   server.URL[7:], // Strip http://
        Scheme: "http",
    }
    client, _ := weaviate.NewClient(cfg)
    store := NewStore(client)

    // 3. Call Search
    _, err := store.Search(context.Background(), "test", []float32{0.1}, 0.5, 10, nil)
    
    // 4. Expect Error
    assert.Error(t, err)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/adapter/weaviate/... -v`

**Step 3: Write minimal implementation**
(Add the test case to `store_test.go`)

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/adapter/weaviate/... -v`