package retrieval

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueryLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewQueryLogger(&buf)

	entry := QueryLogEntry{
		Query:      "test",
		Duration:   100 * time.Millisecond,
		NumResults: 5,
	}

	logger.Log(entry)

	var output map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &output)
	assert.NoError(t, err)
	assert.Equal(t, "test", output["query"])
	assert.Equal(t, 5.0, output["num_results"])
	assert.Equal(t, 100.0, output["latency_ms"])
}
