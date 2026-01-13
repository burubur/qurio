# Plan: Source Naming Refactor (Mandatory Name)

## Scope
Enforce mandatory naming for all sources at creation time.
1.  **API:** `POST /sources` (Web) and `POST /sources/upload` (File) require a `name` field.
2.  **Validation:** Reject requests without a valid name.
3.  **Persistence:** Save `name` immediately to the database.
4.  **Ingestion:** Worker uses the pre-defined name for all chunks.
5.  **Frontend & E2E:** Update UI and tests to send `name`.

## Gap Analysis
*   **Noun "ID"**: Mapped to `Source.ID`.
*   **Noun "Name"**: Mapped to `Source.Name`. Currently unused/empty.
*   **Noun "Type"**: Mapped to `Source.Type` (file/web).
*   **Noun "URL"**: Mapped to `Source.URL`.
*   **Verb "Naming"**: SHIFT responsibility from "Ingestion Worker" to "User/API".
*   **Constraint "Mandatory"**: Add validation checks in Handlers.

## Tasks

### Task 1: Enforce Name in Web Source Creation

**Files:**
- Modify: `apps/backend/features/source/handler.go` (Handler.Create)
- Modify: `apps/backend/features/source/source.go` (Service.Create)
- Test: `apps/backend/features/source/handler_create_test.go` (New)

**Requirements:**
- **Acceptance Criteria**
  1. `POST /sources` returns 400 Bad Request if `name` is missing or empty.
  2. `Service.Create` saves the provided `name` to the `sources` table.

- **Functional Requirements**
  1. Parse `name` from JSON body.
  2. Validate `name != ""`.
  3. Pass `name` to `Service.Create`.

- **Test Coverage**
  - [Integration] `TestCreateSource_MissingName`: Assert 400.
  - [Integration] `TestCreateSource_Success`: Assert `name` is saved in DB.

**Step 1: Write failing test**
Create `apps/backend/features/source/handler_create_test.go`:
```go
package source_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
    "encoding/json"

	"github.com/stretchr/testify/assert"
    "qurio/apps/backend/features/source"
    "qurio/apps/backend/internal/testutil" // Assuming mock helpers exist
)

func TestCreateSource_MissingName(t *testing.T) {
    // Setup Mock Service (Using any typical mock pattern from project)
    mockService := new(testutil.MockSourceService)
    h := source.NewHandler(mockService)

    body := []byte(`{"url":"https://example.com","type":"web"}`) // No Name
    req := httptest.NewRequest("POST", "/sources", bytes.NewBuffer(body))
    w := httptest.NewRecorder()

    h.Create(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

**Step 2: Verify test fails**
Run test. Expect failure (currently 201 or other error).

**Step 3: Write minimal implementation**
1. Update Request struct in `Handler.Create`.
2. Add validation check.
3. Ensure `Source` struct passed to service has `Name`.

**Step 4: Verify test passes**
Run test. Expect pass.

### Task 2: Enforce Name in File Upload

**Files:**
- Modify: `apps/backend/features/source/handler.go` (Handler.Upload)
- Modify: `apps/backend/features/source/source.go` (Service.Upload)
- Test: `apps/backend/features/source/handler_upload_test.go` (New)

**Requirements:**
- **Acceptance Criteria**
  1. `POST /sources/upload` returns 400 Bad Request if `name` form field is missing/empty.
  2. `Service.Upload` accepts `name` argument.
  3. `Service.Upload` saves `name` to `sources` table.

- **Functional Requirements**
  1. `r.FormValue("name")`.
  2. Validate `name != ""`.
  3. Pass to `service.Upload(ctx, path, hash, name)`.

- **Test Coverage**
  - [Integration] `TestUpload_MissingName`: Multipart request without name -> 400.
  - [Integration] `TestUpload_Success`: Multipart request with name -> DB has name.

**Step 1: Write failing test**
In `apps/backend/features/source/handler_upload_test.go`:
```go
func TestUpload_MissingName(t *testing.T) {
    // Setup Multipart Writer
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, _ := writer.CreateFormFile("file", "test.txt")
    part.Write([]byte("content"))
    writer.Close()

    req := httptest.NewRequest("POST", "/sources/upload", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    w := httptest.NewRecorder()

    h.Upload(w, req)
    assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

**Step 2: Verify test fails**
Run test. Expect fail.

**Step 3: Write minimal implementation**
1. Get `name` from form.
2. Validate.
3. Update `Service.Upload` signature to accept `name`.
4. Update call site in Handler.

**Step 4: Verify test passes**
Run test. Expect pass.

### Task 3: Verify Worker Propagation (Validation)

**Files:**
- Test: `apps/backend/internal/worker/result_consumer_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Verify that `ResultConsumer` uses the name returned by `SourceFetcher` (which comes from DB).
  2. No code changes needed in Worker.

- **Test Coverage**
  - [Unit] `TestHandleMessage_UsesConfiguredName`:
    - Mock `SourceFetcher.GetSourceConfig` to return "My Manual Name".
    - Send payload.
    - Assert `Publish` to `ingest.embed` contains `"source_name":"My Manual Name"`.

### Task 4: Update Frontend (SourceForm)

**Files:**
- Modify: `apps/frontend/src/features/sources/SourceForm.vue`
- Modify: `apps/frontend/src/features/sources/source.store.ts`

**Requirements:**
- **Acceptance Criteria**
  1. `SourceForm.vue`: "File Upload" tab has a "Name" input field (Required).
  2. `source.store.ts`: `uploadSource` function accepts `name` parameter.
  3. `source.store.ts`: `uploadSource` appends `name` to FormData.

- **Functional Requirements**
  1. Add `<input v-model="form.name" required />`.
  2. Update store action.

- **Test Coverage**
  - Manual verification via UI or basic component test (if setup).

### Task 5: Update E2E Tests

**Files:**
- Modify: `apps/e2e/tests/ingestion.spec.ts`

**Requirements:**
- **Acceptance Criteria**
  1. E2E tests provide `name` during file upload.
  2. E2E tests provide `name` during web crawl (if applicable).

- **Functional Requirements**
  1. `page.fill('input[name="name"]', fileName)` before clicking upload.

- **Test Coverage**
  - Run `npx playwright test ingestion.spec.ts`.

**Step 1: Write failing test**
Run `npx playwright test ingestion.spec.ts` -> Should fail after Backend changes (Task 1 & 2).

**Step 2: Fix test**
Add `await page.getByLabel('Name').fill(fileName);`

**Step 3: Verify test passes**
Run `npx playwright test ingestion.spec.ts` -> Pass.
