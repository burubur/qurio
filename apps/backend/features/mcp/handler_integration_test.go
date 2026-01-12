package mcp_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"qurio/apps/backend/features/mcp"
	"qurio/apps/backend/features/source"
	"qurio/apps/backend/internal/adapter/weaviate"
	"qurio/apps/backend/internal/middleware"
	"qurio/apps/backend/internal/retrieval"
	"qurio/apps/backend/internal/settings"
	"qurio/apps/backend/internal/testutils"
	"qurio/apps/backend/internal/worker"
)

// MockEmbedder
type MockEmbedder struct {
	mock.Mock
}

func (m *MockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	// Return a dummy vector matching Weaviate schema dimension if strict, 
	// or just any vector if vectorizer is 'none' and we don't enforce dim check in test
	return []float32{0.1, 0.2, 0.3}, nil
}

func TestMCPHandler_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	s := testutils.NewIntegrationSuite(t)
	s.Setup()
	defer s.Teardown()

	ctx := context.Background()

	// 1. Setup Dependencies
	vectorStore := weaviate.NewStore(s.Weaviate)
	require.NoError(t, vectorStore.EnsureSchema(ctx))

	embedder := new(MockEmbedder)
	settingsRepo := settings.NewPostgresRepo(s.DB)
	settingsSvc := settings.NewService(settingsRepo)
	retrievalSvc := retrieval.NewService(embedder, vectorStore, nil, settingsSvc, nil)
	sourceRepo := source.NewPostgresRepo(s.DB)

	handler := mcp.NewHandler(retrievalSvc, sourceRepo)

	// 2. Seed Data
	// A. Source in DB
	src := &source.Source{
		Type:        "web",
		URL:         "http://example.com",
		ContentHash: "hash-mcp",
		Status:      "completed",
		Name:        "MCP Test Source",
	}
	err := sourceRepo.Save(ctx, src)
	require.NoError(t, err)

	_, err = sourceRepo.BulkCreatePages(ctx, []source.SourcePage{{
		SourceID: src.ID,
		URL:      src.URL,
		Status:   "completed",
		Depth:    0,
	}})
	require.NoError(t, err)

	// B. Chunks in Weaviate
	chunk := worker.Chunk{
		SourceID:   src.ID,
		SourceURL:  src.URL,
		Content:    "The quick brown fox jumps over the lazy dog.",
		ChunkIndex: 0,
		Title:      "Fox Page",
		Type:       "web",
		Vector:     []float32{0.1, 0.2, 0.3},
	}
	err = vectorStore.StoreChunk(ctx, chunk)
	require.NoError(t, err)

	// 3. Test qurio_search via JSON-RPC
	searchArgs := mcp.SearchArgs{
		Query: "fox",
	}
	argsBytes, _ := json.Marshal(searchArgs)
	
	callParams := mcp.CallParams{
		Name:      "qurio_search",
		Arguments: argsBytes,
	}
	paramsBytes, _ := json.Marshal(callParams)

	reqBody := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsBytes,
		ID:      1,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/mcp/messages?sessionId=test-session", strings.NewReader(string(bodyBytes)))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp mcp.JSONRPCResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Nil(t, resp.Error)
	
	resultMap, ok := resp.Result.(map[string]interface{})
	require.True(t, ok)
	
	contentList, ok := resultMap["content"].([]interface{})
	require.True(t, ok)
	require.NotEmpty(t, contentList)
	
	firstContent := contentList[0].(map[string]interface{})
	text := firstContent["text"].(string)
	
	assert.Contains(t, text, "Fox Page")
	assert.Contains(t, text, "The quick brown fox")

	// 4. Test qurio_read_page
	readArgs := mcp.FetchPageArgs{
		URL: src.URL,
	}
	readArgsBytes, _ := json.Marshal(readArgs)
	callParamsRead := mcp.CallParams{
		Name:      "qurio_read_page",
		Arguments: readArgsBytes,
	}
	paramsReadBytes, _ := json.Marshal(callParamsRead)
	reqBodyRead := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsReadBytes,
		ID:      2,
	}
	bodyReadBytes, _ := json.Marshal(reqBodyRead)
	
	reqRead := httptest.NewRequest("POST", "/mcp", strings.NewReader(string(bodyReadBytes)))
	rrRead := httptest.NewRecorder()
	
	handler.ServeHTTP(rrRead, reqRead)
	
	assert.Equal(t, http.StatusOK, rrRead.Code)
	var respRead mcp.JSONRPCResponse
	json.Unmarshal(rrRead.Body.Bytes(), &respRead)
	assert.Nil(t, respRead.Error)

	resMapRead := respRead.Result.(map[string]interface{})
	contentListRead := resMapRead["content"].([]interface{})
	textRead := contentListRead[0].(map[string]interface{})["text"].(string)
	
	assert.Contains(t, textRead, "The quick brown fox")
}

func TestHandler_SSE_Correlation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	s := testutils.NewIntegrationSuite(t)
	s.Setup()
	defer s.Teardown()

	// Dependencies
	vectorStore := weaviate.NewStore(s.Weaviate)
	embedder := new(MockEmbedder)
	settingsRepo := settings.NewPostgresRepo(s.DB)
	settingsSvc := settings.NewService(settingsRepo)
	retrievalSvc := retrieval.NewService(embedder, vectorStore, nil, settingsSvc, nil)
	sourceRepo := source.NewPostgresRepo(s.DB)
	handler := mcp.NewHandler(retrievalSvc, sourceRepo)

	// Create test server to avoid Race Condition with ResponseRecorder
	mux := http.NewServeMux()
	mux.HandleFunc("/mcp/sse", handler.HandleSSE)
	mux.HandleFunc("/mcp/messages", handler.HandleMessage)
	
	// Wrap with middleware
	server := httptest.NewServer(middleware.CorrelationID(mux))
	defer server.Close()

	// 1. Establish SSE Connection
	resp, err := http.Get(server.URL + "/mcp/sse")
	require.NoError(t, err)
	defer resp.Body.Close()

	// 2. Read Session ID from stream
	var sessionID string
	buffer := make([]byte, 1024)
	timeout := time.After(2 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			n, err := resp.Body.Read(buffer)
			if err != nil {
				break
			}
			chunk := string(buffer[:n])
			if strings.Contains(chunk, "event: id") {
				parts := strings.Split(chunk, "event: id\ndata: ")
				if len(parts) > 1 {
					rest := parts[1]
					idPart := strings.Split(rest, "\n")[0]
					sessionID = strings.TrimSpace(idPart)
					if sessionID != "" {
						done <- true
						return
					}
				}
			}
		}
	}()

	select {
	case <-timeout:
		t.Fatal("Timeout waiting for SSE session ID")
	case <-done:
		// success
	}
	
	assert.NotEmpty(t, sessionID)

	// 3. Send Message with Correlation ID (Verify synchronous error response includes it)
	correlationID := "test-correlation-id-123"
	
	client := &http.Client{}
	reqErr, _ := http.NewRequest("POST", server.URL+"/mcp/messages?sessionId="+sessionID, strings.NewReader("invalid-json"))
	reqErr.Header.Set("X-Correlation-ID", correlationID)

	respErr, err := client.Do(reqErr)
	require.NoError(t, err)
	defer respErr.Body.Close()
	
	assert.Equal(t, http.StatusBadRequest, respErr.StatusCode)
	var errResp map[string]interface{}
	
	bodyBytes, _ := io.ReadAll(respErr.Body)
	err = json.Unmarshal(bodyBytes, &errResp)
	require.NoError(t, err)
	
	assert.Equal(t, correlationID, errResp["correlationId"])
}