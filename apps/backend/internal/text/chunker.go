package text

import "strings"

// Chunk splits text into chunks of approximately 'size' tokens (words) with 'overlap' tokens.
// This is a simple heuristic where 1 word ≈ 1 token.
func Chunk(text string, size, overlap int) []string {
	words := strings.Fields(text)
	var chunks []string
	if len(words) == 0 {
		return chunks
	}

	if size <= 0 {
		size = 512 // Default
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= size {
		overlap = size - 1 // Ensure progress
	}

	step := size - overlap
	if step < 1 {
		step = 1
	}

	for i := 0; i < len(words); i += step {
		end := i + size
		if end > len(words) {
			end = len(words)
		}
		
		chunk := strings.Join(words[i:end], " ")
		chunks = append(chunks, chunk)

		if end == len(words) {
			break
		}
	}
	return chunks
}
