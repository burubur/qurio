package mcp

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Retriever interface {
	Search(ctx context.Context, query string) ([]string, error)
}

type Handler struct {
	retriever Retriever
}

func NewHandler(r Retriever) *Handler {
	return &Handler{retriever: r}
}

// JSON-RPC Request types
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

type CallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type SearchArgs struct {
	Query string `json:"query"`
}

// JSON-RPC Response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type ToolResult struct {
	Content []ToolContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

const (
	ErrParse          = -32700
	ErrInvalidRequest = -32600
	ErrMethodNotFound = -32601
	ErrInvalidParams  = -32602
	ErrInternal       = -32603
)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("mcp request received", "method", r.Method, "path", r.URL.Path)
	
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, nil, ErrParse, "Parse error")
		return
	}

	// Minimal implementation for 'tools/call' method
	if req.Method == "tools/call" {
		var params CallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			slog.Warn("invalid params structure", "error", err)
			h.writeError(w, req.ID, ErrInvalidParams, "Invalid params")
			return
		}

		if params.Name == "search" {
			var args SearchArgs
			if err := json.Unmarshal(params.Arguments, &args); err != nil {
				slog.Warn("invalid search arguments", "error", err)
				h.writeError(w, req.ID, ErrInvalidParams, "Invalid search arguments")
				return
			}

			results, err := h.retriever.Search(r.Context(), args.Query)
			if err != nil {
				slog.Error("search failed", "error", err)
				// Return tool error result, not protocol error, if the tool execution failed
				// OR return Internal Error depending on strictness. 
				// MCP usually prefers returning a ToolResult with isError=true for tool failures.
				response := JSONRPCResponse{
					JSONRPC: "2.0",
					ID:      req.ID,
					Result: ToolResult{
						Content: []ToolContent{{Type: "text", Text: "Error: " + err.Error()}},
						IsError: true,
					},
				}
				json.NewEncoder(w).Encode(response)
				return
			}
			
			textResult := "No results"
			if len(results) > 0 {
				textResult = results[0] // Simplify
			}

			slog.Info("tool execution completed", "tool", "search", "result_count", len(results))

			response := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: ToolResult{
					Content: []ToolContent{
						{Type: "text", Text: textResult},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		slog.Warn("method not found", "method", params.Name)
		h.writeError(w, req.ID, ErrMethodNotFound, "Method not found: "+params.Name)
		return
	}
	
	slog.Warn("unknown jsonrpc method", "method", req.Method)
	h.writeError(w, req.ID, ErrMethodNotFound, "Method not found")
}

func (h *Handler) writeError(w http.ResponseWriter, id interface{}, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	// JSON-RPC errors are usually 200 OK at HTTP level, containing the error object
	// But some implementations use 400/500. We'll use 200 to be safe with clients 
	// that parse the body regardless of status, or 400/500 if strict HTTP semantics are needed.
	// Standard JSON-RPC over HTTP typically uses 200 OK.
	w.WriteHeader(http.StatusOK) 

	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: map[string]interface{}{
			"code":    code,
			"message": message,
		},
		ID: id,
	}
	json.NewEncoder(w).Encode(resp)
}
