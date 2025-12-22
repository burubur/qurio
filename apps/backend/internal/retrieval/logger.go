package retrieval

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type QueryLogEntry struct {
	Timestamp   time.Time     `json:"timestamp"`
	Query       string        `json:"query"`
	NumResults  int           `json:"num_results"`
	Duration    time.Duration `json:"duration_ns"`
	LatencyMs   int64         `json:"latency_ms"`
	CorrelationID string      `json:"correlation_id"`
}

type QueryLogger struct {
	writer io.Writer
}

func NewQueryLogger(w io.Writer) *QueryLogger {
	return &QueryLogger{writer: w}
}

func NewFileQueryLogger(path string) (*QueryLogger, error) {
	// Ensure directory exists
	// But path might be "data/logs/query.log", so dirname is "data/logs"
	// For now assume caller handles it or just try open
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	mw := io.MultiWriter(os.Stdout, f)
	return NewQueryLogger(mw), nil
}

func (l *QueryLogger) Log(entry QueryLogEntry) {
	entry.Timestamp = time.Now()
	entry.LatencyMs = entry.Duration.Milliseconds()
	json.NewEncoder(l.writer).Encode(entry)
}
