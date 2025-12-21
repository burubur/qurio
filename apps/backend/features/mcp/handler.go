package mcp

import (
	"context"
	"encoding/json"
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

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Minimal implementation for 'tools/call' method
	if req.Method == "tools/call" {
		var params CallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return
		}

		if params.Name == "search" {
			var args SearchArgs
			if err := json.Unmarshal(params.Arguments, &args); err != nil {
				return
			}

			results, _ := h.retriever.Search(r.Context(), args.Query)
			
			// Format as MCP Tool Result
			textResult := "No results"
			if len(results) > 0 {
				textResult = results[0] // Simplify
			}

			response := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: ToolResult{
					Content: []ToolContent{
						{Type: "text", Text: textResult},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
