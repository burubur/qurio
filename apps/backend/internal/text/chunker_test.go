package text_test

import (
	"testing"
	"qurio/apps/backend/internal/text"
	"strings"
)

func TestChunk(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		size        int
		overlap     int
		expectedLen int
	}{
		{
			name:        "Empty input",
			input:       "",
			size:        10,
			overlap:     0,
			expectedLen: 0,
		},
		{
			name:        "Small input (no split)",
			input:       "hello world",
			size:        10,
			overlap:     0,
			expectedLen: 1,
		},
		{
			name:        "Exact split",
			input:       "one two three four",
			size:        2,
			overlap:     0,
			expectedLen: 2, // "one two", "three four"
		},
		{
			name:        "Overlap split",
			input:       "one two three four",
			size:        3,
			overlap:     1,
			expectedLen: 2, // "one two three", "three four" (step = 2: 0->3, 2->4)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := text.Chunk(tt.input, tt.size, tt.overlap)
			if len(chunks) != tt.expectedLen {
				t.Errorf("Chunk() length = %v, want %v. Chunks: %v", len(chunks), tt.expectedLen, chunks)
			}
		})
	}
}

func TestChunk_LongText(t *testing.T) {
	// Generate 1000 words
	input := ""
	for i := 0; i < 1000; i++ {
		input += "word "
	}
	
	// Chunk size 100, overlap 0 -> 10 chunks
	chunks := text.Chunk(input, 100, 0)
	if len(chunks) != 10 {
		t.Errorf("Expected 10 chunks, got %d", len(chunks))
	}
}

func TestChunk_OverlapLogic(t *testing.T) {
	input := "1 2 3 4 5 6"
	// Size 3, Overlap 1 -> Step 2
	// [1 2 3] (0-3)
	// [3 4 5] (2-5)
	// [5 6]   (4-6)
	chunks := text.Chunk(input, 3, 1)
	if len(chunks) != 3 {
		t.Errorf("Expected 3 chunks, got %d", len(chunks))
	}
	if !strings.Contains(chunks[0], "1 2 3") {
		t.Errorf("Chunk 0 mismatch: %s", chunks[0])
	}
	if !strings.Contains(chunks[1], "3 4 5") {
		t.Errorf("Chunk 1 mismatch: %s", chunks[1])
	}
}
