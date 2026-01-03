package text

import (
	"testing"
)

func TestChunkMarkdown_API(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantType ChunkType
	}{
		{
			name:     "Prose with API keywords",
			content:  "This endpoint uses the GET method to retrieve data from the URL.",
			wantType: ChunkTypeAPI,
		},
		{
			name:     "Prose without API keywords",
			content:  "This is a normal paragraph describing a cat.",
			wantType: ChunkTypeProse,
		},
		{
			name:     "Code block with http language",
			content:  "```http\nGET /api/v1/users\n```",
			wantType: ChunkTypeAPI,
		},
		{
			name:     "Code block with swagger language",
			content:  "```swagger\nswagger: '2.0'\n```",
			wantType: ChunkTypeAPI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := ChunkMarkdown(tt.content, 100, 0)
			if len(chunks) == 0 {
				t.Fatalf("expected chunks, got 0")
			}
			if chunks[0].Type != tt.wantType {
				t.Errorf("got Type %q, want %q", chunks[0].Type, tt.wantType)
			}
		})
	}
}
