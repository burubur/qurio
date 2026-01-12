### Task 1: Refactor Backend to Support Split Ingestion Topics

**Files:**
- Create: `apps/backend/internal/config/topics.go`
- Modify: `apps/backend/features/source/source.go:140-305`
- Modify: `apps/backend/features/job/service.go:38-42`
- Modify: `apps/backend/internal/worker/result_consumer.go:266`
- Test: `apps/backend/features/source/source_test.go`
- Test: `apps/backend/features/job/service_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `SourceService.Create` publishes to `ingest.task.web` if type is "web" (or default).
  2. `SourceService.Create` publishes to `ingest.task.file` if type is "file".
  3. `SourceService.Upload` always publishes to `ingest.task.file`.
  4. `SourceService.ReSync` publishes to the correct topic based on source type.
  5. `JobService.Retry` parses payload and publishes to the correct topic.
  6. `ResultConsumer` (crawler) publishes discovered links to `ingest.task.web`.

- **Functional Requirements**
  1. Define constants for topics `ingest.task.web` and `ingest.task.file`.
  2. Parse JSON payload in Job Retry to extract `type`.

- **Non-Functional Requirements**
  - Maintain backward compatibility for existing tests where possible (by updating mocks).

- **Test Coverage**
  - [Unit] `SourceService.Create` - Verify `Publish` called with correct topic.
  - [Unit] `JobService.Retry` - Verify `Publish` called with correct topic based on payload.

**Step 1: Write failing test**
```go
// apps/backend/features/source/source_test.go

func TestCreate_FileSource_PublishesToFileTopic(t *testing.T) {
	// ... setup mocks ...
	src := &source.Source{Type: "file", URL: "/tmp/test.pdf"}
	// Expect Publish to "ingest.task.file"
	mockPub.On("Publish", "ingest.task.file", mock.Anything).Return(nil)
	
	err := service.Create(ctx, src)
	assert.NoError(t, err)
	mockPub.AssertExpectations(t)
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/source/...`
Expected: FAIL - mock expected "ingest.task.file" but got "ingest.task"

**Step 3: Write minimal implementation**
```go
// apps/backend/internal/config/topics.go
package config
const (
	TopicIngestWeb  = "ingest.task.web"
	TopicIngestFile = "ingest.task.file"
)

// apps/backend/features/source/source.go
// Update Publish calls to use config.TopicIngestWeb or config.TopicIngestFile based on src.Type
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/source/...`
Expected: PASS


### Task 2: Refactor Docker Compose for Independent Scaling

**Files:**
- Modify: `docker-compose.yml`

**Requirements:**
- **Acceptance Criteria**
  1. `ingestion-worker` service is replaced by `ingestion-worker-web` and `ingestion-worker-file`.
  2. `ingestion-worker-web` listens to `ingest.task.web`.
  3. `ingestion-worker-file` listens to `ingest.task.file`.
  4. Both workers have access to `qurio_uploads` volume.

- **Functional Requirements**
  1. Use `NSQ_TOPIC_INGEST` environment variable to configure topics.
  2. Allocate higher resources (or keep default) for file worker if needed (keeping default for now as per "scale independent" requirement, user can change later).

- **Non-Functional Requirements**
  - No downtime is not required (dev environment).

- **Test Coverage**
  - Manual verification via `docker compose config` and `docker compose up`.

**Step 1: Write failing test**
(Infrastructure change - manual verification)

**Step 2: Verify test fails**
N/A

**Step 3: Write minimal implementation**
```yaml
# docker-compose.yml
services:
  ingestion-worker-web:
    build: ./apps/ingestion-worker
    environment:
      - NSQ_TOPIC_INGEST=ingest.task.web
    # ...
  ingestion-worker-file:
    build: ./apps/ingestion-worker
    environment:
      - NSQ_TOPIC_INGEST=ingest.task.file
    # ...
```

**Step 4: Verify test passes**
Run: `docker compose config`
Expected: Valid configuration with two worker services.


### Task 3: Verify Worker Configuration

**Files:**
- Test: `apps/ingestion-worker/tests/test_config.py` (create if needed)

**Requirements:**
- **Acceptance Criteria**
  1. Verify `Settings` class correctly loads `NSQ_TOPIC_INGEST` from environment.

- **Functional Requirements**
  1. Pydantic `BaseSettings` handles this automatically, but a test confirms it.

- **Test Coverage**
  - [Unit] `test_settings_override`

**Step 1: Write failing test**
```python
# apps/ingestion-worker/tests/test_config.py
import os
from config import Settings

def test_topic_override():
    os.environ["NSQ_TOPIC_INGEST"] = "ingest.test.topic"
    settings = Settings()
    assert settings.nsq_topic_ingest == "ingest.test.topic"
    del os.environ["NSQ_TOPIC_INGEST"]
```

**Step 2: Verify test fails**
Run: `pytest apps/ingestion-worker/tests/test_config.py`
Expected: PASS (if Pydantic works as expected, this is a confirmation step. If it fails, we fix config.py)

**Step 3: Write minimal implementation**
(If validation fails, modify `apps/ingestion-worker/config.py` to ensure env var reading)

**Step 4: Verify test passes**
Run: `pytest apps/ingestion-worker/tests/test_config.py`
Expected: PASS


### Task 4: Update Worker Unit Tests & Integration Tests

**Files:**
- Modify: `apps/backend/internal/worker/result_consumer_test.go`
- Modify: `apps/backend/internal/worker/integration_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `ResultConsumer` publishes discovered links to `ingest.task.web`.
  2. Integration tests use correct topic for verification.

- **Functional Requirements**
  1. Update `MockPublisher` in `result_consumer_test.go` to expect `ingest.task.web`.
  2. Update integration test setup if it relies on specific topic subscription (if applicable).

- **Test Coverage**
  - [Unit] `ResultConsumer.HandleMessage` (Recursion case) - Verify topic is `ingest.task.web`.
  - [Integration] `TestIngestIntegration` - Ensure compatibility with topic split.

**Step 1: Write failing test**
```go
// apps/backend/internal/worker/result_consumer_test.go
// In TestResultConsumer_HandleMessage_Success_Recursion
// Update the mock expectation to fail if it receives the old topic
mockPub.On("Publish", "ingest.task.web", mock.Anything).Return(nil)
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/worker/...`
Expected: FAIL (publishing to "ingest.task")

**Step 3: Write minimal implementation**
```go
// apps/backend/internal/worker/result_consumer.go
// Update Publish call to use config.TopicIngestWeb
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/worker/...`
Expected: PASS


### Task 5: Integration Test with Testcontainers

**Files:**
- Create: `apps/backend/internal/worker/topic_integration_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Verify `SourceService` publishes to `ingest.task.web` for Web sources.
  2. Verify `SourceService` publishes to `ingest.task.file` for File sources.
  3. Use real NSQ, Postgres, and Weaviate via `IntegrationSuite`.

- **Functional Requirements**
  1. Create a consumer for `ingest.task.web` and verify message receipt.
  2. Create a consumer for `ingest.task.file` and verify message receipt.

- **Test Coverage**
  - [Integration] `TestTopicRouting`

**Step 1: Write failing test**
(Test file creation is the step here)

```go
// apps/backend/internal/worker/topic_integration_test.go
package worker_test

import (
    "context"
    "testing"
    "time"
    "encoding/json"

    "github.com/nsqio/go-nsq"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "qurio/apps/backend/features/source"
    "qurio/apps/backend/internal/testutils"
    "qurio/apps/backend/internal/worker"
    "qurio/apps/backend/internal/config"
)

func TestTopicRouting(t *testing.T) {
    s := testutils.NewIntegrationSuite(t)
    s.Setup()
    defer s.Teardown()

    ctx := context.Background()

    // 1. Setup Service
    repo := source.NewPostgresRepo(s.DB)
    // We need a ChunkStore and SettingsService for NewService, mocks or minimal impls are fine if not focus
    // But since we have IntegrationSuite, maybe we can use real ones or mocks if simpler.
    // For this test, we care about 'pub' (NSQ Producer).
    // SourceService uses 's.NSQ' which is the producer.
    
    // We need real ChunkStore for NewService? 
    // Let's use mocks for non-critical parts to keep test focused on NSQ
    
    svc := source.NewService(repo, s.NSQ, &worker.MockChunkStore{}, &worker.MockSettings{}) 

    // 2. Setup Consumers for verification
    webChan := make(chan *nsq.Message, 1)
    fileChan := make(chan *nsq.Message, 1)

    nsqCfg := nsq.NewConfig()
    
    webConsumer, _ := nsq.NewConsumer(config.TopicIngestWeb, "test-ch", nsqCfg)
    webConsumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
        webChan <- m
        return nil
    }))
    // Connect to the random port from testcontainers
    nsqAddr := s.NSQ.String() // This gives producer addr, we need to extract host:port or usage helper
    // s.GetAppConfig().NSQDHost gives the address
    appCfg := s.GetAppConfig()
    webConsumer.ConnectToNSQD(appCfg.NSQDHost)

    fileConsumer, _ := nsq.NewConsumer(config.TopicIngestFile, "test-ch", nsqCfg)
    fileConsumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
        fileChan <- m
        return nil
    }))
    fileConsumer.ConnectToNSQD(appCfg.NSQDHost)

    // 3. Action: Create Web Source
    webSrc := &source.Source{Type: "web", URL: "http://example.com/topic-test"}
    err := svc.Create(ctx, webSrc)
    require.NoError(t, err)

    // 4. Verify Web Topic
    select {
    case msg := <-webChan:
        var payload map[string]interface{}
        json.Unmarshal(msg.Body, &payload)
        assert.Equal(t, "web", payload["type"])
        assert.Equal(t, "http://example.com/topic-test", payload["url"])
        msg.Finish()
    case <-time.After(5 * time.Second):
        t.Fatal("Timeout waiting for web task")
    }

    // 5. Action: Create File Source
    // Upload calls repo.Save then Publish
    fileSrc, err := svc.Upload(ctx, "/tmp/test.pdf", "hash-topic-test")
    require.NoError(t, err)

    // 6. Verify File Topic
    select {
    case msg := <-fileChan:
        var payload map[string]interface{}
        json.Unmarshal(msg.Body, &payload)
        assert.Equal(t, "file", payload["type"])
        assert.Equal(t, "/tmp/test.pdf", payload["path"])
        msg.Finish()
    case <-time.After(5 * time.Second):
        t.Fatal("Timeout waiting for file task")
    }
}
```

**Step 2: Verify test fails**
Run: `go test -v apps/backend/internal/worker/topic_integration_test.go`
Expected: FAIL (Compilation error or Runtime fail if code not implemented)

**Step 3: Write minimal implementation**
(This task is purely verification, implementation is in Task 1. So we just ensure the test runs and passes with the code from Task 1)

**Step 4: Verify test passes**
Run: `go test -v apps/backend/internal/worker/topic_integration_test.go`
Expected: PASS
